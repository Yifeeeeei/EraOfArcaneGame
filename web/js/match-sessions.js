/* Browser-local seats retain the original player ID when returning to a room. */
(function (root) {
  "use strict";
  const KEY = "eraofarcane.matches.v1", TTL = 24 * 60 * 60 * 1000;
  function list() {
    try {
      const rows = JSON.parse(root.localStorage.getItem(KEY) || "[]");
      return Array.isArray(rows) ? rows.filter(r => r && /^\d{4}$/.test(r.room) &&
        typeof r.player_id === "string" && r.player_id && typeof r.deck_code === "string" &&
        r.deck_code && Number.isFinite(r.saved_at) && Date.now() - r.saved_at < TTL) : [];
    } catch { return []; }
  }
  function write(rows) {
    try { root.localStorage.setItem(KEY, JSON.stringify(rows.slice(0, 6))); return true; }
    catch { return false; }
  }
  function save(seat) {
    if (!seat?.room || !seat.player_id || !seat.deck_code) return false;
    return write([{room:seat.room, player_id:seat.player_id, player_name:seat.player_name || "玩家",
      deck_code:seat.deck_code, saved_at:Date.now()},
      ...list().filter(r => r.room !== seat.room || r.player_id !== seat.player_id)]);
  }
  function remove(room, playerID) {
    write(list().filter(r => r.room !== room || r.player_id !== playerID));
  }
  function url(seat) {
    return "/game.html?" + new URLSearchParams({room:seat.room, player_id:seat.player_id,
      player_name:seat.player_name || "玩家", deck_code:seat.deck_code}).toString();
  }
  root.ArcaneMatchSessions = {KEY, list, save, remove, url};
})(typeof window === "undefined" ? globalThis : window);
