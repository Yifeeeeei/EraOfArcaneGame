/* Presentation coordinates only. Wire coordinates always remain unchanged. */
(function (root, factory) {
  const view = factory();
  if (typeof module === "object" && module.exports) module.exports = view;
  else root.ArcaneBoardView = view;
})(typeof window === "undefined" ? globalThis : window, function () {
  function cells(side, active = [], center = null) {
    const order = side === "enemy" ? [2, 1, 0] : [0, 1, 2];
    return order.flatMap((row) =>
      order.map((col) => ({
        key: `${col}-${row}`,
        col,
        row,
        active: active.some((p) => p.col === col && p.row === row),
        center: center?.col === col && center?.row === row,
      })),
    );
  }
  function square(origin) {
    return [0, 1].flatMap((dc) =>
      [0, 1].map((dr) => ({ col: origin.col + dc, row: origin.row + dr })),
    );
  }
  return { cells, square };
});
