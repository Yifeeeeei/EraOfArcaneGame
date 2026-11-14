/* Shared presentation only. Match rules and decisions remain in game.js. */
(function (root) {
  "use strict";
  const KEY = "eraofarcane.preferences.v1";
  const IMAGE_BASE = "https://yifeeeeei.github.io/ArcaneImages/";
  const defaults = { motion: "full", effects: "normal", largeText: false };
  function preferences() {
    try {
      return { ...defaults, ...JSON.parse(localStorage.getItem(KEY) || "{}") };
    } catch {
      return { ...defaults };
    }
  }
  function apply() {
    const p = preferences();
    document.documentElement.classList.toggle(
      "reduce-motion",
      p.motion === "reduced",
    );
    document.documentElement.classList.toggle("large-text", !!p.largeText);
    document.documentElement.dataset.effects = p.effects;
    return p;
  }
  function save(p) {
    try {
      localStorage.setItem(KEY, JSON.stringify(p));
      apply();
      root.dispatchEvent(new CustomEvent("arcane-preferences", { detail: p }));
      return true;
    } catch {
      toast("浏览器无法保存设置");
      return false;
    }
  }
  const crest =
    '<svg class="brand-mark" viewBox="0 0 48 48" aria-hidden="true"><circle cx="24" cy="24" r="21"/><g class="arcane-weave"><path d="M24 5C34 13 36 24 27 33C23 37 17 39 12 37C23 34 29 25 25 17L20 21L24 5Z"/><path d="M24 5C34 13 36 24 27 33C23 37 17 39 12 37C23 34 29 25 25 17L20 21L24 5Z" transform="rotate(120 24 24)"/><path d="M24 5C34 13 36 24 27 33C23 37 17 39 12 37C23 34 29 25 25 17L20 21L24 5Z" transform="rotate(240 24 24)"/></g></svg>';
  function toast(message) {
    document.querySelector(".arcane-toast")?.remove();
    const node = document.createElement("div");
    node.className = "arcane-toast";
    node.setAttribute("role", "status");
    node.textContent = message;
    document.body.append(node);
    setTimeout(() => node.remove(), 3200);
  }
  function dialog(title, content, afterOpen) {
    const previous = document.activeElement;
    const host = document.createElement("div");
    host.className = "arcane-dialog-backdrop";
    host.innerHTML =
      '<section class="arcane-dialog" role="dialog" aria-modal="true" aria-label="' +
      title +
      '"><button class="btn icon-btn dialog-close" aria-label="关闭">×</button><h2>' +
      title +
      "</h2>" +
      content +
      "</section>";
    function close() {
      host.remove();
      document.removeEventListener("keydown", keydown);
      previous?.focus();
    }
    function keydown(e) {
      if (e.key === "Escape") {
        e.preventDefault();
        close();
      }
      if (e.key === "Tab") {
        const nodes = [
          ...host.querySelectorAll("button,input,select,textarea,a[href]"),
        ].filter((n) => !n.disabled && n.getClientRects().length);
        const first = nodes[0],
          last = nodes.at(-1);
        if (e.shiftKey && document.activeElement === first) {
          e.preventDefault();
          last.focus();
        } else if (!e.shiftKey && document.activeElement === last) {
          e.preventDefault();
          first.focus();
        }
      }
    }
    host.addEventListener("click", (e) => {
      if (e.target === host || e.target.closest(".dialog-close")) close();
    });
    document.addEventListener("keydown", keydown);
    document.body.append(host);
    host.querySelector("button").focus();
    afterOpen?.(host, close);
    return close;
  }
  function settings() {
    const p = preferences();
    dialog(
      "界面设置",
      '<label class="setting-row">动态效果<select name="motion"><option value="full">标准</option><option value="reduced">减少动态效果</option></select></label><label class="setting-row">特效强度<select name="effects"><option value="normal">标准</option><option value="low">轻量</option><option value="off">关闭</option></select></label><label class="setting-row">较大文字<input type="checkbox" name="largeText"></label><div class="dialog-footer"><button class="btn btn-primary" data-save>保存</button></div>',
      (host, close) => {
        host.querySelector("[name=motion]").value = p.motion;
        host.querySelector("[name=effects]").value = p.effects;
        host.querySelector("[name=largeText]").checked = p.largeText;
        host.querySelector("[data-save]").onclick = () => {
          if (
            save({
              motion: host.querySelector("[name=motion]").value,
              effects: host.querySelector("[name=effects]").value,
              largeText: host.querySelector("[name=largeText]").checked,
            })
          ) {
            close();
            toast("设置已保存");
          }
        };
      },
    );
  }
  const elementNames = {
    火: "fire",
    水: "water",
    气: "air",
    地: "earth",
    光: "light",
    暗: "dark",
    无: "none",
  };
  function elementIcon(k) {
    return "/assets/elements/" + (elementNames[k] || "none") + ".png";
  }
  function image(card) {
    if (card?.hidden || card?.is_hidden) return "/assets/card-back.png";
    return card?.output_path
      ? IMAGE_BASE +
          card.output_path.replace(/\\/g, "/").replace(/\/{2,}/g, "/")
      : "/images/card-back.svg";
  }
  function artworkPath(card) {
    if (!card || card.hidden || card.is_hidden) return "";
    return root.ArcaneArtworkIndex?.[card.number] || "";
  }
  function artwork(card) {
    const path = artworkPath(card);
    return path ? IMAGE_BASE + path : image(card);
  }
  function artworkLoaded(event) {
    const img = event.target;
    if ((img.currentSrc || img.src).includes("/drawings/")) {
      img.removeAttribute("data-art-fallback");
    }
  }
  function artAttrs(card) {
    return {
      src: artwork(card),
      "data-artwork": artworkPath(card) ? "true" : null,
      "data-art-type": card?.type || null,
      "data-card-src": image(card),
      onLoad: artworkLoaded,
    };
  }
  function imageError(e) {
    const img = e.target;
    const url = img.getAttribute("src") || "";
    if (url.includes("yifeeeeei.github.io/ArcaneImages/")) {
      img.src = url.replace(
        "https://yifeeeeei.github.io/ArcaneImages/",
        "https://raw.githubusercontent.com/Yifeeeeei/ArcaneImages/main/",
      );
    } else if (img.dataset.artwork === "true" && url.includes("/drawings/")) {
      img.setAttribute("data-art-fallback", "true");
      img.src = img.dataset.cardSrc || "/images/card-back.svg";
    } else if (!url.includes("/images/card-back.svg")) {
      img.src = "/images/card-back.svg";
      img.classList.add("image-unavailable");
    }
  }
  function cost(cost) {
    const entries = Object.entries(cost || {}).filter(([, n]) => n > 0);
    return entries.length ? entries.map(([e, n]) => e + n).join(" ") : "免费";
  }
  async function request(url, options = {}) {
    const r = await fetch(url, options);
    let body;
    try {
      body = await r.json();
    } catch {
      throw Error("服务器返回无效响应");
    }
    if (!r.ok || body?.error) throw Error(body?.error || "请求失败");
    return body;
  }
  function download(name, text, type = "text/plain") {
    const link = document.createElement("a");
    const url = URL.createObjectURL(new Blob([text], { type }));
    link.href = url;
    link.download = name;
    link.click();
    setTimeout(() => URL.revokeObjectURL(url), 1000);
  }
  root.ArcaneUI = {
    preferences,
    apply,
    save,
    crest,
    toast,
    dialog,
    settings,
    elementIcon,
    image,
    artwork,
    artAttrs,
    artworkLoaded,
    imageError,
    cost,
    request,
    download,
  };
  apply();
  root.addEventListener("storage", (e) => {
    if (e.key === KEY) apply();
  });
  document.addEventListener("click", (e) => {
    if (e.target.closest("[data-open-settings]")) settings();
  });
})(window);
