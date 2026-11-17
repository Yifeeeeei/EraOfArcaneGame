const { test } = require("node:test");
const assert = require("node:assert/strict");
const view = require("./board-view.js");

test("front-right is top-right for you and bottom-left for your opponent", () => {
  const position = { col: 2, row: 0 };
  const own = view.cells("own", [position]);
  const enemy = view.cells("enemy", [position]);
  assert.equal(
    own.findIndex((c) => c.active),
    2,
  );
  assert.equal(
    enemy.findIndex((c) => c.active),
    6,
  );
  assert.deepEqual(
    own.map((c) => c.key).sort(),
    enemy.map((c) => c.key).sort(),
  );
});

test("a square preview rotates the display while preserving its four wire positions", () => {
  const area = view.square({ col: 1, row: 0 });
  const own = view.cells("own", area, { col: 1, row: 1 });
  const enemy = view.cells("enemy", area, { col: 1, row: 1 });
  assert.deepEqual(
    own
      .filter((c) => c.active)
      .map((c) => c.key)
      .sort(),
    enemy
      .filter((c) => c.active)
      .map((c) => c.key)
      .sort(),
  );
  assert.deepEqual(
    enemy.filter((c) => c.active).map((c) => c.key),
    ["2-1", "1-1", "2-0", "1-0"],
  );
  assert.equal(enemy.filter((c) => c.center).length, 1);
});
