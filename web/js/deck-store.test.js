const { test } = require("node:test");
const assert = require("node:assert/strict");
const DS = require("./deck-store.js");
const catalog = require("../assets/preconstructed-decks.json");
const recipe = catalog.decks.find((d) => d.id === "基础包/合一火");
const code = recipe.deckCode;
test("Composer deck codes preserve zones, duplicates and the optional extra list", () => {
  const d = DS.parse(code);
  assert.equal(d.hero, "4111001");
  assert.equal(d.main.length, 30);
  assert.equal(d.skills.length, 10);
  d.extra = ["1001101"];
  assert.deepEqual(DS.parse(DS.encode(d)), d);
  assert.deepEqual(DS.groups(["1021001", "1021001", "1321001"]), [
    { number: "1021001", count: 2 },
    { number: "1321001", count: 1 },
  ]);
  assert.throws(() => DS.parse("4311003 // 1021001"), /卡组码/);
  assert.throws(() => DS.parse("4311003 // <script> // 3321002"), /卡号/);
});
test("editing one local deck preserves other decks and selection across a reload", () => {
  let raw = null;
  const storage = {
    getItem: () => raw,
    setItem: (key, value) => (raw = value),
  };
  let data = DS.put(DS.load(storage), { name: "原卡组", code });
  const first = data.selectedId;
  data = DS.put(data, { name: "另一套", code });
  const second = data.selectedId;
  data = DS.put(data, { id: first, name: "重新命名", code });
  DS.write(data, storage);
  const loaded = DS.load(storage);
  assert.equal(loaded.decks.length, 2);
  assert.equal(loaded.selectedId, first);
  assert.equal(loaded.decks[0].name, "重新命名");
  assert.ok(loaded.decks.some((d) => d.id === second));
});
test("official recipes preserve all zones and repeated selection preserves saved decks", () => {
  for (const r of catalog.decks) {
    const parsed = DS.parse(r.deckCode);
    assert.equal(parsed.main.length, 30, r.id);
    assert.equal(parsed.skills.length, 10, r.id);
    assert.deepEqual(DS.parse(DS.encode(parsed)), parsed, r.id);
  }
  let data = DS.put({ decks: [] }, { name: "我的构筑", code: "4111001 // //" });
  const custom = data.decks[0];
  data = DS.useRecipe(data, recipe);
  const importedId = data.selectedId;
  assert.equal(data.decks[0].name, "合一火 · 基础包");
  assert.deepEqual(data.decks[1], custom);
  data = DS.useRecipe(data, recipe);
  assert.equal(data.decks.length, 2);
  assert.equal(data.selectedId, importedId);
  // A renamed matching build remains the player's record, including its name.
  data = DS.put(data, {
    id: importedId,
    name: "我的火焰",
    code: DS.encode(DS.parse(code)),
  });
  data = DS.useRecipe(data, recipe);
  assert.equal(data.decks[0].name, "我的火焰");
  // Editing the old import must not let a later recipe selection overwrite it.
  data = DS.put(data, {
    id: importedId,
    name: "火焰草稿",
    code: "4111001 // //",
  });
  data = DS.useRecipe(data, recipe);
  assert.equal(data.decks.length, 3);
  assert.ok(
    data.decks.some((d) => d.id === importedId && d.name === "火焰草稿"),
  );
});
test("damaged, unsupported or inaccessible storage is reported without overwriting it", () => {
  let writes = 0;
  for (const raw of [
    "broken JSON",
    "null",
    '{"version":1,"decks":[null]}',
    '{"version":9,"decks":[]}',
    '{"version":1,"decks":[{}]}',
  ]) {
    assert.throws(() =>
      DS.load({ getItem: () => raw, setItem: () => writes++ }),
    );
  }
  assert.equal(writes, 0);
  assert.throws(
    () =>
      DS.write(
        { decks: [] },
        {
          setItem: () => {
            throw Error("quota");
          },
        },
      ),
    /quota/,
  );
});
