/* Pure display helpers; gameplay eligibility remains engine-owned. */
(function (root) {
  const priority = [
    "横置",
    "石化",
    "眩晕",
    "冻结",
    "封印",
    "点燃",
    "虚弱",
    "冷却",
    "精通",
    "技能受限",
    "水纹",
  ];
  const aliases = {
    winter_bow_water_mark: "水纹",
    cannot_use_skill_until_turn: "技能受限",
  };
  function statuses(card, all = false, limit = 0) {
    if (!card || card.hidden || card.is_hidden) return [];
    const entries = card.is_horizontal ? [{ name: "横置", label: "横置" }] : [];
    for (const [key, raw] of Object.entries(card.statuses || {})) {
      if (/[A-Za-z_]/.test(key) && !aliases[key]) continue;
      const count = Number(raw),
        name = aliases[key] || key;
      if (
        count <= 0 ||
        !Number.isFinite(count) ||
        (!all && !priority.includes(name))
      )
        continue;
      entries.push({
        name,
        label:
          name +
          (count > 1 && key !== "cannot_use_skill_until_turn" ? count : ""),
      });
    }
    if (!all)
      entries.sort(
        (a, b) => priority.indexOf(a.name) - priority.indexOf(b.name),
      );
    return limit > 0 ? entries.slice(0, limit) : entries;
  }
  function slotPages(slots, capacity, bound = [], label = "技能") {
    const pages = [];
    function add(cards, count, kind, title) {
      for (let start = 0; start < count; start += 5) {
        pages.push({
          kind,
          start,
          label: `${title} ${start + 1}–${Math.min(start + 5, count)}`,
          entries: Array.from({ length: 5 }, (_, offset) => ({
            card: cards[start + offset] || null,
            index: start + offset,
            kind,
          })),
        });
      }
    }
    add(
      slots || [],
      Math.max(5, Number(capacity) || 5, slots?.length || 0),
      "normal",
      label,
    );
    if (bound.length) add(bound, bound.length, "bound", "绑定技能");
    return pages;
  }
  function cardPage(pages, id) {
    return pages.findIndex((page) =>
      page.entries.some((entry) => entry.card?.instance_id === id),
    );
  }
  function needsTarget(spell) {
    return !!spell?.needs_target;
  }
  const api = { statuses, slotPages, cardPage, needsTarget };
  if (typeof module !== "undefined") module.exports = api;
  root.ArcaneBattlePresentation = api;
})(typeof window === "undefined" ? globalThis : window);
