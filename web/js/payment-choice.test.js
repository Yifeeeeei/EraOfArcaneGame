const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const planner = require("./payment-choice.js");
const source = fs.readFileSync(require.resolve("./game.js"), "utf8");
const order = ["火", "水", "气", "地", "光", "暗", "无"];
const plain = (value) => JSON.parse(JSON.stringify(value));

function scope(elements = {}, flags = {}) {
  const context = {
    ELEMENT_ORDER: order,
    myState: { value: { elements, ...flags } },
    hasLightWildcardPayment: () => !!flags.light_wildcard,
    cloneElements: (input) =>
      Object.fromEntries(order.map((e) => [e, Number(input?.[e] || 0)])),
    ArcanePayments: planner,
    paymentRequest: { value: null },
    paymentSelection: { value: {} },
    totalCost: (cost) =>
      Object.values(cost || {}).reduce((sum, n) => sum + n, 0),
    calls: [],
  };
  context.sendAction = (...args) => {
    context.calls.push(args);
    return true;
  };
  vm.createContext(context);
  const start = source.indexOf("function calculatePaymentWithAvailable(");
  const end = source.indexOf("function totalCost(", start);
  vm.runInContext(source.slice(start, end), context);
  const sendStart = source.indexOf("function sendPaidAction(");
  vm.runInContext(
    source.slice(
      sendStart,
      source.indexOf("function paymentAvailable(", sendStart),
    ),
    context,
  );
  return context;
}

function analyze(available, cost, flags) {
  return plain(
    planner.analyze(
      available,
      cost,
      scope(available, flags).calculatePaymentWithAvailable,
    ),
  );
}

test("specific costs auto-pay the exact amount despite surplus or irrelevant elements", () => {
  assert.deepEqual(analyze({ 火: 2 }, { 火: 1 }), {
    kind: "unique",
    payment: { 火: 1 },
  });
  assert.deepEqual(analyze({ 火: 4, 地: 1 }, { 火: 2, 地: 1 }), {
    kind: "unique",
    payment: { 火: 2, 地: 1 },
  });
  assert.deepEqual(analyze({ 火: 4, 地: 8 }, { 火: 1 }), {
    kind: "unique",
    payment: { 火: 1 },
  });
});

test("arcane involvement asks only when a different payment is possible", () => {
  assert.equal(analyze({ 火: 1, 无: 1 }, { 火: 1 }).kind, "choice");
  assert.equal(analyze({ 火: 1, 地: 1 }, { 无: 1 }).kind, "choice");
  assert.deepEqual(analyze({ 无: 2 }, { 火: 1 }), {
    kind: "unique",
    payment: { 无: 1 },
  });
  assert.deepEqual(analyze({ 火: 2 }, { 无: 1 }), {
    kind: "unique",
    payment: { 火: 1 },
  });
  assert.deepEqual(analyze({ 火: 1, 地: 1 }, { 火: 1, 无: 1 }), {
    kind: "unique",
    payment: { 火: 1, 地: 1 },
  });
  assert.deepEqual(analyze({ 火: 1, 无: 1 }, { 火: 2 }), {
    kind: "unique",
    payment: { 火: 1, 无: 1 },
  });
});

test("active light-payment rules preserve alternatives in either direction", () => {
  assert.equal(
    analyze({ 火: 1, 光: 1 }, { 火: 1 }, { light_wildcard: true }).kind,
    "choice",
  );
  assert.deepEqual(analyze({ 光: 2 }, { 火: 1 }, { light_wildcard: true }), {
    kind: "unique",
    payment: { 光: 1 },
  });
  assert.equal(
    analyze({ 水: 1, 光: 1 }, { 光: 1 }, { light_cost_wildcard: true }).kind,
    "choice",
  );
  assert.deepEqual(
    analyze({ 水: 2 }, { 光: 1 }, { light_cost_wildcard: true }),
    { kind: "unique", payment: { 水: 1 } },
  );
});

test("unique payment uses the normal revision gate immediately, without another confirmation", () => {
  const c = scope({ 火: 2 });
  c.sendPaidAction(
    "summon",
    { instance_id: "unit", col: 0, row: 0 },
    { 火: 1 },
    "召唤",
  );
  assert.equal(c.calls.length, 1);
  assert.equal(c.calls[0][3], true);
  assert.deepEqual(plain(c.calls[0][1]), {
    instance_id: "unit",
    col: 0,
    row: 0,
    payment: { 火: 1 },
  });
  assert.equal(c.paymentRequest.value, null);
});

test("ambiguous or unaffordable payments never submit automatically", () => {
  for (const elements of [{ 火: 1, 无: 1 }, { 地: 1 }]) {
    const c = scope(elements);
    c.sendPaidAction("summon", { instance_id: "unit" }, { 火: 1 }, "召唤");
    assert.equal(c.calls.length, 0);
    assert.deepEqual(plain(c.paymentRequest.value.cost), { 火: 1 });
  }
  const c = scope({ 火: 1 }, { strict_arcane: 1 });
  c.sendPaidAction("summon", { instance_id: "unit" }, { 无: 1 }, "召唤");
  assert.equal(
    c.calls.length,
    0,
    "dedicated arcane resources must not be ignored by auto-pay",
  );
});

test("free and combined costs use the same path and preserve caller callbacks", () => {
  const c = scope({ 火: 4, 地: 1 });
  const done = () => {};
  c.sendPaidAction("learn_skill", { instance_id: "skill" }, {}, "学习", done);
  assert.equal(c.calls[0][3], true);
  assert.equal(c.calls[0][2].onSuccess, done);
  assert.equal(c.paymentRequest.value, null);
  c.sendPaidAction("summon", { instance_id: "unit" }, { 火: 2, 地: 1 }, "召唤");
  assert.deepEqual(plain(c.calls[1][1].payment), { 火: 2, 地: 1 });
});

test("dedicated arcane keeps its restriction and offers meaningful alternatives", () => {
  const solve = scope().calculatePaymentWithAvailable;
  function quote(pool, cost, strict) {
    const context = planner.context(pool, solve, strict);
    return plain(planner.analyze(context.available, cost, context.solve));
  }
  assert.deepEqual(quote({}, { 无: 1 }, 2), {
    kind: "unique",
    payment: { 无: 1 },
  });
  assert.equal(quote({}, { 火: 1 }, 2).kind, "unavailable");
  assert.equal(quote({ 火: 1 }, { 无: 1 }, 1).kind, "choice");
  assert.equal(quote({ 火: 2 }, { 火: 1, 无: 1 }, 1).kind, "choice");
  assert.deepEqual(quote({ 火: 1 }, { 火: 1, 无: 1 }, 1), {
    kind: "unique",
    payment: { 火: 1, 无: 1 },
  });
  assert.deepEqual(quote({ 无: 1 }, { 无: 1 }, 1), {
    kind: "unique",
    payment: { 无: 1 },
  });
});

test("unique dedicated-arcane payment auto-submits while an alternative remains selectable", () => {
  const forced = scope({}, { strict_arcane: 2 });
  forced.sendPaidAction("summon", { instance_id: "unit" }, { 无: 1 }, "召唤");
  assert.equal(forced.calls.length, 1);
  assert.deepEqual(plain(forced.calls[0][1].payment), { 无: 1 });
  const choice = scope({ 火: 1 }, { strict_arcane: 1 });
  choice.sendPaidAction("summon", { instance_id: "unit" }, { 无: 1 }, "召唤");
  assert.equal(choice.calls.length, 0);
  assert.deepEqual(plain(choice.paymentRequest.value.available), {
    火: 1,
    无: 1,
  });
});

test("uniqueness matches exhaustive independent payment enumeration for 729 small pools/costs", () => {
  const keys = ["火", "地", "无"];
  const possibilities = [];
  for (let f = 0; f <= 2; f++)
    for (let e = 0; e <= 2; e++)
      for (let a = 0; a <= 2; a++) possibilities.push({ 火: f, 地: e, 无: a });
  const solve = scope().calculatePaymentWithAvailable;
  for (const available of possibilities)
    for (const cost of possibilities) {
      const valid = [];
      const total = cost.火 + cost.地 + cost.无;
      for (let f = 0; f <= available.火; f++)
        for (let e = 0; e <= available.地; e++)
          for (let a = 0; a <= available.无; a++) {
            if (
              f + e + a === total &&
              Math.max(0, cost.火 - f) + Math.max(0, cost.地 - e) <= a
            )
              valid.push({ 火: f, 地: e, 无: a });
          }
      const before = { ...available };
      const result = planner.analyze(available, cost, solve);
      assert.equal(
        result.kind,
        valid.length === 0
          ? "unavailable"
          : valid.length === 1
            ? "unique"
            : "choice",
        JSON.stringify({ available, cost }),
      );
      if (valid.length === 1)
        for (const e of keys) assert.equal(result.payment[e] || 0, valid[0][e]);
      assert.deepEqual(available, before);
    }
});
