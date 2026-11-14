const { test } = require("node:test");
const assert = require("node:assert/strict");
const P = require("./battle-presentation.js");

test("five-slot pages preserve actual slot indexes, gaps and bound-skill separation", () => {
  const cards = Array.from({ length: 7 }, (_, i) => ({ instance_id: `s${i}` }));
  cards[2] = null;
  const bound = Array.from({ length: 6 }, (_, i) => ({ instance_id: `b${i}` }));
  const pages = P.slotPages(cards, 7, bound);
  assert.equal(pages.length, 4);
  assert.ok(pages.every((page) => page.entries.length === 5));
  assert.equal(pages[0].entries[2].card, null);
  assert.equal(pages[1].entries[1].index, 6);
  assert.equal(pages[1].entries[1].card.instance_id, "s6");
  assert.equal(pages[1].entries[2].card, null);
  assert.equal(pages[2].kind, "bound");
  assert.equal(pages[3].entries[0].card.instance_id, "b5");
  assert.equal(P.cardPage(pages, "b5"), 3);
  assert.equal(P.cardPage(pages, "missing"), -1);
  assert.equal(P.slotPages([], 99, [], "装备").length, 20);
  assert.equal(
    P.slotPages(cards, 5).length,
    2,
    "cards awaiting a capacity cleanup remain visible",
  );
});

test("battle badges keep important states; full inspection preserves the other public labels", () => {
  const card = {
    is_horizontal: true,
    statuses: {
      祈祷之焰标记物: 3,
      精通: 2,
      冷却: 1,
      眩晕: 1,
      点燃: 2,
      internal_counter: 5,
    },
  };
  assert.deepEqual(
    P.statuses(card).map((s) => s.label),
    ["横置", "眩晕", "点燃2", "冷却", "精通2"],
  );
  assert.ok(P.statuses(card, true).some((s) => s.label === "祈祷之焰标记物3"));
  assert.ok(!P.statuses(card, true).some((s) => s.label.includes("internal")));
  assert.deepEqual(P.statuses({ ...card, is_hidden: true }, true), []);
  assert.equal(P.statuses(card, false, 2).length, 2);
});

test("targeting follows structured needs_target for all skill categories", () => {
  for (const is_sorcery of [true, false]) {
    assert.equal(
      P.needsTarget({ is_sorcery, needs_target: false, power: 8 }),
      false,
    );
    assert.equal(
      P.needsTarget({ is_sorcery, needs_target: true, power: 0 }),
      true,
    );
  }
  assert.equal(P.needsTarget(null), false);
});
