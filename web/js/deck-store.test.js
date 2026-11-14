const { test } = require("node:test");
const assert = require("node:assert/strict");
const DS = require("./deck-store.js");
test("Composer deck codes preserve zones, duplicates and the optional extra list", () => {
  const d = DS.parse(DS.SAMPLE);
  assert.equal(d.hero, "4311003");
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
  let data = DS.put(DS.load(storage), { name: "原卡组", code: DS.SAMPLE });
  const first = data.selectedId;
  data = DS.put(data, { name: "另一套", code: DS.SAMPLE });
  const second = data.selectedId;
  data = DS.put(data, { id: first, name: "重新命名", code: DS.SAMPLE });
  DS.write(data, storage);
  const loaded = DS.load(storage);
  assert.equal(loaded.decks.length, 2);
  assert.equal(loaded.selectedId, first);
  assert.equal(loaded.decks[0].name, "重新命名");
  assert.ok(loaded.decks.some((d) => d.id === second));
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
