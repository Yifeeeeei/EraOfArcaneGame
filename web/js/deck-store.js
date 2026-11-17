/* Local deck files contain construction data only, never live match state. */
(function (root) {
  "use strict";
  const KEY = "eraofarcane.decks.v1";
  function parse(code) {
    const parts = String(code || "")
      .split("//")
      .map((s) => s.trim());
    if (parts.length < 3 || parts.length > 4)
      throw Error("卡组码需要 人物 // 主卡组 // 技能池，可附加额外卡组");
    const values = parts.map((s) => (s ? s.split(/\s+/) : []));
    if (values.some((a) => a.some((n) => !/^\d{7}$/.test(n))))
      throw Error("卡号必须为 7 位数字");
    if (values[0].length > 1) throw Error("只能选择一张人物");
    return {
      hero: values[0][0] || "",
      main: values[1],
      skills: values[2],
      extra: values[3] || [],
    };
  }
  function encode(d) {
    const p = [
      d.hero || "",
      (d.main || []).join(" "),
      (d.skills || []).join(" "),
    ];
    if (d.extra?.length) p.push(d.extra.join(" "));
    return p.join(" // ");
  }
  function groups(ids) {
    const counts = new Map();
    for (const id of ids || []) counts.set(id, (counts.get(id) || 0) + 1);
    return [...counts].map(([number, count]) => ({ number, count }));
  }
  function load(storage = root.localStorage) {
    const raw = storage.getItem(KEY);
    if (!raw) return { version: 1, decks: [], selectedId: "" };
    let data;
    try {
      data = JSON.parse(raw);
    } catch {
      throw Error("本地卡组文件损坏，请先导出备份");
    }
    if (
      !data ||
      data.version !== 1 ||
      !Array.isArray(data.decks) ||
      data.decks.some(
        (d) =>
          !d ||
          typeof d.id !== "string" ||
          typeof d.code !== "string" ||
          typeof d.name !== "string",
      )
    )
      throw Error("本地卡组格式不受支持，请先导出备份");
    return data;
  }
  function write(data, storage = root.localStorage) {
    storage.setItem(KEY, JSON.stringify({ ...data, version: 1 }));
    return data;
  }
  function put(data, deck) {
    const record = {
      id:
        deck.id ||
        "d_" +
          Date.now().toString(36) +
          "_" +
          Math.random().toString(36).slice(2, 8),
      name: String(deck.name || "未命名卡组").slice(0, 60),
      code: String(deck.code),
      updatedAt: new Date().toISOString(),
    };
    return {
      ...data,
      decks: [record, ...data.decks.filter((d) => d.id !== record.id)],
      selectedId: record.id,
    };
  }
  function recipeName(recipe) {
    return `${recipe.name} · ${recipe.version.join(" / ")}`;
  }
  function useRecipe(data, recipe) {
    const code = encode(parse(recipe.deckCode));
    const existing = data.decks.find((deck) => {
      try {
        return encode(parse(deck.code)) === code;
      } catch {
        return false;
      }
    });
    return existing
      ? { ...data, selectedId: existing.id }
      : put(data, { name: recipeName(recipe), code });
  }
  const api = {
    KEY,
    parse,
    encode,
    groups,
    load,
    write,
    put,
    recipeName,
    useRecipe,
  };
  if (typeof module !== "undefined") module.exports = api;
  root.ArcaneDecks = api;
})(typeof window === "undefined" ? globalThis : window);
