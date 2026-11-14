/* Local deck files contain construction data only, never live match state. */
(function (root) {
  "use strict";
  const KEY = "eraofarcane.decks.v1";
  const SAMPLE =
    "4311003 // 1021001 1021001 1021004 1021004 1021011 1021011 1321002 1321002 1321001 1321001 1321003 1321003 1321007 1321007 1321013 1321013 1321011 1321011 1321010 1321010 1321008 1321008 1321004 1321004 1021006 1021006 1021007 1021007 1021012 1021012 // 3321002 3321003 3321005 3321010 3321013 3321015 3021008 3021009 3021005 3021003";
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
  const api = { KEY, SAMPLE, parse, encode, groups, load, write, put };
  if (typeof module !== "undefined") module.exports = api;
  root.ArcaneDecks = api;
})(typeof window === "undefined" ? globalThis : window);
