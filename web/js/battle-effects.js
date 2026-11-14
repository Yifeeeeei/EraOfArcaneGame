/* Decorative event consumer. Never dispatches, queues or delays game actions. */
(function (root) {
  "use strict";
  const colors = {
    火: "#ee9857",
    水: "#8fcbd7",
    气: "#b4e0c4",
    地: "#b7a072",
    光: "#e9d59a",
    暗: "#bd9ab7",
    无: "#b8bedb",
  };
  const kinds = {
    火: "fire",
    水: "water",
    气: "air",
    地: "earth",
    光: "light",
    暗: "dark",
    无: "arcane",
  };
  let previous = new Map(),
    timers = new Set();
  function later(fn, ms) {
    const t = setTimeout(() => {
      timers.delete(t);
      fn();
    }, ms);
    timers.add(t);
    return t;
  }
  function allowed() {
    return root.ArcaneUI?.preferences().effects !== "off" && !document.hidden;
  }
  function reduced() {
    return (
      root.ArcaneUI?.preferences().motion === "reduced" ||
      matchMedia("(prefers-reduced-motion:reduce)").matches
    );
  }
  function layer() {
    return document.querySelector(".vfx-layer");
  }
  function point(node) {
    if (!node) return null;
    const r = node.getBoundingClientRect();
    return {
      x: r.left + r.width / 2,
      y: r.top + r.height / 2,
      width: r.width,
      height: r.height,
    };
  }
  function byID(id) {
    if (!id) return null;
    const escaped = CSS.escape(String(id));
    return (
      document.querySelector(
        '.unit-cell[data-instance-id="' + escaped + '"]',
      ) ||
      document.querySelector('.slot[data-instance-id="' + escaped + '"]') ||
      document.querySelector('[data-defense-card="' + escaped + '"]')
    );
  }
  function atPosition(owner, col, row) {
    return document.querySelector(
      '.unit-cell[data-owner-id="' +
        Number(owner) +
        '"][data-col="' +
        Number(col) +
        '"][data-row="' +
        Number(row) +
        '"]',
    );
  }
  function center() {
    const r = document
      .querySelector(".battlefield-center")
      ?.getBoundingClientRect();
    return r
      ? { x: r.left + r.width / 2, y: r.top + r.height / 2 }
      : { x: innerWidth / 2, y: innerHeight / 2 };
  }
  function anchor(card) {
    if (!card) return center();
    const dom = point(byID(card.instance_id));
    if (dom) return dom;
    const old = previous.get(card.instance_id);
    if (old?.point) return old.point;
    if (card.position)
      return (
        point(
          atPosition(
            card.owner ?? card.owner_id,
            card.position.col,
            card.position.row,
          ),
        ) || center()
      );
    return center();
  }
  function add(className, p, color, text, ms = 850) {
    if (!allowed()) return;
    const parent = layer();
    if (!parent) return;
    while (parent.childElementCount >= 24) parent.firstChild.remove();
    const n = document.createElement("div");
    n.className = "arcane-vfx " + className + (reduced() ? " vfx-static" : "");
    n.style.left = p.x + "px";
    n.style.top = p.y + "px";
    n.style.setProperty("--fx", color || colors["无"]);
    if (text) n.textContent = text;
    parent.append(n);
    later(() => n.remove(), reduced() ? 450 : ms);
    return n;
  }
  function pulse(card, kind = "arcane", label = "") {
    const p = anchor(card);
    add(
      "vfx-ring vfx-" + kind,
      p,
      colors[Object.keys(kinds).find((k) => kinds[k] === kind)],
    );
    if (label) add("vfx-label", p, null, label, 1000);
  }
  function trail(from, to, element) {
    if (!allowed()) return;
    const a = from || center(),
      b = to || center(),
      dx = b.x - a.x,
      dy = b.y - a.y;
    const n = add(
      "vfx-trail vfx-" + (kinds[element] || "arcane"),
      a,
      colors[element],
      null,
      650,
    );
    if (!n) return;
    n.style.width = Math.hypot(dx, dy) + "px";
    n.style.setProperty("--angle", Math.atan2(dy, dx) + "rad");
    if (!reduced() && root.ArcaneUI.preferences().effects !== "low") {
      for (let i = 0; i < 5; i++) {
        const t = (i + 1) / 6;
        const p = { x: a.x + dx * t, y: a.y + dy * t };
        const spark = add("vfx-spark", p, colors[element], null, 650);
        spark?.style.setProperty("--delay", i * 35 + "ms");
      }
    }
  }
  function targetPoint(target, attacker, state, slot, spectator) {
    if (!target) return center();
    const owner = target.owner_id ?? 1 - Number(attacker);
    if (target.type === "hero") {
      const own = spectator ? owner === 0 : owner === slot;
      return anchor((own ? state?.you : state?.opponent)?.hero);
    }
    if (target.position)
      return (
        point(atPosition(owner, target.position.col, target.position.row)) ||
        center()
      );
    if (target.instance_id) return anchor(target);
    return center();
  }
  function legendary(card) {
    if (!allowed() || !card?.name) return;
    document.querySelector(".legendary-arrival")?.remove();
    const banner = document.createElement("div");
    banner.className =
      "legendary-arrival" + (reduced() ? " legend-static" : "");
    banner.setAttribute("role", "status");
    const halo = document.createElement("div");
    halo.className = "legendary-halo";
    halo.innerHTML = ArcaneUI.crest;
    halo.setAttribute("aria-hidden", "true");
    const label = document.createElement("span");
    label.className = "legendary-kicker";
    label.textContent = "LEGEND AWAKENS";
    const art = document.createElement("div");
    art.className = "legendary-art";
    const img = document.createElement("img");
    const attrs = ArcaneUI.artAttrs(card);
    for (const [key, value] of Object.entries(attrs)) {
      if (key === "onLoad") img.onload = value;
      else if (value !== null) img.setAttribute(key, value);
    }
    img.alt = "";
    img.onerror = ArcaneUI.imageError;
    art.append(img);
    const name = document.createElement("strong");
    name.textContent = card.name;
    const skip = document.createElement("button");
    skip.className = "mini-btn";
    skip.textContent = "跳过";
    skip.setAttribute("aria-label", "跳过登场特效");
    skip.onclick = () => banner.remove();
    const identity = document.createElement("span");
    identity.className = "legendary-identity";
    identity.textContent = [
      card.category === "无" ? "奥术" : card.category,
      card.type,
    ]
      .filter(Boolean)
      .join(" · ");
    banner.append(halo, label, art, name, identity, skip);
    document.body.append(banner);
    later(() => banner.remove(), reduced() ? 1200 : 1600);
  }
  function cardsIn(state, slot, spectator) {
    const result = new Map();
    for (const [side, p] of [
      ["you", state?.you],
      ["opponent", state?.opponent],
    ]) {
      if (!p) continue;
      const owner = spectator
        ? side === "you"
          ? 0
          : 1
        : side === "you"
          ? slot
          : 1 - slot;
      const push = (card, pos) => {
        if (card && !card.is_hidden && !card.hidden && card.instance_id)
          result.set(card.instance_id, {
            card,
            owner,
            point:
              point(byID(card.instance_id)) ||
              (pos ? point(atPosition(owner, pos.col, pos.row)) : null),
          });
      };
      for (let c = 0; c < 3; c++)
        for (let r = 0; r < 3; r++) push(p.units?.[c]?.[r], { col: c, row: r });
      for (const card of [...(p.skills || []), ...(p.equipment || [])])
        push(card);
    }
    return result;
  }
  function snapshot(next, old, slot, spectator) {
    const now = cardsIn(next, slot, spectator),
      oldMap = cardsIn(old, slot, spectator);
    // This is a visual profile for the released Red Moon card, not a rule
    // inferred from its text. Duration comes from authoritative instance state.
    const moon = [...now.values()].find(
      ({ card }) =>
        card.number === "3611101" && Number(card.statuses?.["异能持续"]) > 0,
    );
    const oldMoon = [...oldMap.values()].find(
      ({ card }) =>
        card.number === "3611101" && Number(card.statuses?.["异能持续"]) > 0,
    );
    document.documentElement.classList.toggle(
      "red-moon-table",
      !!moon && allowed(),
    );
    if (moon && !oldMoon && old && allowed()) {
      add("vfx-red-moon", center(), "#d18077", "红月", 1100);
    }
    if (old && next.revision !== old.revision && allowed())
      for (const [id, entry] of now) {
        const before = oldMap.get(id)?.card,
          card = entry.card;
        if (!before) continue;
        const delta = Number(card.current_life) - Number(before.current_life);
        if (delta && Number.isFinite(delta)) {
          add(
            "vfx-number " + (delta < 0 ? "vfx-loss" : "vfx-gain"),
            entry.point || center(),
            null,
            (delta > 0 ? "+" : "") + delta,
            1050,
          );
        }
        if (!before.is_horizontal && card.is_horizontal) {
          const element = byID(id);
          if (element && !reduced()) {
            element.classList.add("card-tapped");
            later(() => element.classList.remove("card-tapped"), 400);
          }
        }
        for (const [status, count] of Object.entries(card.statuses || {})) {
          if (
            Number(count) > Number(before.statuses?.[status] || 0) &&
            ["点燃", "冻结", "隐蔽", "石化", "护盾"].includes(status)
          )
            add(
              "vfx-state",
              entry.point || center(),
              status === "点燃" ? colors["火"] : colors["水"],
              status + " " + count,
              900,
            );
        }
      }
    previous = new Map([...oldMap, ...now]);
  }
  function event(evt, state, slot, spectator) {
    if (!evt || !allowed()) return;
    const d = evt.data || {};
    switch (evt.type) {
      case "summon": {
        const p =
          point(atPosition(d.player, d.position?.col, d.position?.row)) ||
          anchor(d.card);
        add("vfx-ring vfx-summon", p, colors[d.card?.category]);
        if (
          String(d.card?.tag || "")
            .split(/[-,，、\s]+/)
            .includes("传奇")
        )
          legendary(d.card);
        break;
      }
      case "equip":
      case "learn_skill":
      case "counter_revealed":
        pulse(
          d.card,
          "arcane",
          evt.type === "counter_revealed" ? d.card?.name : "",
        );
        break;
      case "spell_cast":
      case "spell_hit": {
        const skill = d.skill,
          attacker = d.attacker ?? d.cast_player;
        const p = targetPoint(d.target, attacker, state, slot, spectator);
        if (evt.type === "spell_cast") {
          const side = spectator ? attacker === 0 : attacker === slot;
          const source = anchor((side ? state?.you : state?.opponent)?.hero);
          trail(source, p, skill?.category);
        } else {
          add(
            "vfx-ring vfx-" + (kinds[skill?.category] || "arcane"),
            p,
            colors[skill?.category],
          );
        }
        break;
      }
      case "unit_attack":
        trail(anchor(d.attacker), anchor(d.target), d.attacker?.category);
        break;
      case "consume": {
        const card = previous.get(d.instance_id)?.card || {
          instance_id: d.instance_id,
        };
        const own = spectator ? d.player === 0 : d.player === slot;
        const pool = document.querySelector(
          own
            ? ".player-panel.is-active .elem-grid"
            : ".player-panel.is-opponent .elem-grid",
        );
        trail(
          anchor(card),
          point(pool),
          Object.keys(d.gained || d.elements || {})[0] || "无",
        );
        pulse(card, "arcane", "消耗");
        break;
      }
      case "defense_attempt":
        if (Number(d.overexerted) > 0)
          add(
            "vfx-state",
            center(),
            colors["水"],
            "透支 ×" + d.overexerted,
            850,
          );
        break;
      case "defense_success":
        add("vfx-shield", center(), colors["水"], "防御", 850);
        break;
      case "damage_prevented":
        pulse(d.target, "water", "阻止伤害");
        break;
      case "unit_destroyed": {
        const p = anchor(d.card);
        add("vfx-departure", p, colors[d.card?.category], d.card?.name, 800);
        break;
      }
      case "unit_returned":
        pulse(d.card, "air", "返回手牌");
        break;
      case "spell_reaction":
        pulse(d.card, "arcane", d.card?.name);
        break;
      case "turn_start":
        add(
          "vfx-turn",
          center(),
          colors["光"],
          "第 " + d.turn_number + " 回合",
          1000,
        );
        break;
      case "effect_trigger":
        if (d.effect === "gain_shield")
          pulse(d.source, "water", "护盾 +" + d.amount);
        break;
    }
  }
  function clear() {
    for (const t of timers) clearTimeout(t);
    timers.clear();
    layer()?.replaceChildren();
    document.querySelector(".legendary-arrival")?.remove();
    document.documentElement.classList.remove("red-moon-table");
    document
      .querySelectorAll(".card-tapped")
      .forEach((n) => n.classList.remove("card-tapped"));
  }
  root.ArcaneEffects = { event, snapshot, clear };
  root.addEventListener("arcane-preferences", clear);
  document.addEventListener("visibilitychange", () => {
    if (document.hidden) clear();
  });
})(window);
