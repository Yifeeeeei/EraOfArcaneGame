const { test } = require("node:test");
const assert = require("node:assert/strict");
const vm = require("node:vm");
const fs = require("node:fs");
const path = require("node:path");

function loadUI() {
  const observed = new Set();
  const sandbox = {
    localStorage: { getItem: () => null },
    document: {
      documentElement: { classList: { toggle() {} }, dataset: {} },
      addEventListener() {},
    },
    addEventListener() {},
    ResizeObserver: class {
      observe(img) {
        observed.add(img);
      }
      unobserve(img) {
        observed.delete(img);
      }
    },
  };
  sandbox.window = sandbox;
  vm.createContext(sandbox);
  for (const file of ["../assets/card-art-index.js", "arcane-ui.js"])
    vm.runInContext(
      fs.readFileSync(path.join(__dirname, file), "utf8"),
      sandbox,
    );
  return { ...sandbox, observed };
}

test("every supported original has one bounded, reviewed focal point", () => {
  const { ArcaneArtworkIndex: art, ArcaneArtworkFocusIndex: focus } = loadUI();
  const cards = JSON.parse(
    fs.readFileSync(
      path.join(__dirname, "../../data/supported_card_infos.json"),
      "utf8",
    ),
  );
  assert.deepEqual(
    Object.keys(focus).sort(),
    cards.map((c) => c.number).sort(),
  );
  assert.deepEqual(Object.keys(art).sort(), Object.keys(focus).sort());
  for (const point of Object.values(focus)) {
    assert.equal(point.length, 2);
    assert.ok(point.every((n) => Number.isFinite(n) && n >= 0 && n <= 100));
  }
});

test("crop centers the original-image focal point and clamps at its edges", () => {
  const { ArcaneUI: ui } = loadUI();
  const position = (...args) =>
    ui
      .artFocusPosition(...args)
      .split(" ")
      .map(parseFloat);
  assert.ok(
    Math.abs(position([50, 25], 1000, 1000, 400, 100)[1] - 100 / 6) < 1e-10,
  );
  assert.deepEqual(position([50, 5], 1000, 1000, 400, 100), [50, 0]);
  assert.deepEqual(position([50, 95], 1000, 1000, 400, 100), [50, 100]);
  assert.deepEqual(position([50, 50], 1000, 1000, 400, 100), [50, 50]);
  assert.ok(
    Math.abs(position([25, 50], 1000, 1000, 100, 400)[0] - 100 / 6) < 1e-10,
  );
  assert.equal(ui.artFocusPosition([50, 25], 0, 0, 400, 100), "50% 40%");
  // A low lizard head, unlike most faces, must move the crop downwards.
  assert.ok(position([29, 57], 1254, 1254, 220, 60)[1] > 50);
});

test("hidden cards never expose focal data; fallback and unmount release observation", () => {
  const { ArcaneUI: ui, observed } = loadUI();
  for (const hidden of [{ hidden: true }, { is_hidden: true }]) {
    const attrs = ui.artAttrs({ number: "1021001", ...hidden });
    assert.equal(attrs["data-artwork"], null);
    assert.equal(attrs["data-art-focus-y"], null);
    assert.equal(attrs.src, "/assets/card-back.png");
  }
  const attrs = ui.artAttrs({ number: "1021001" });
  const properties = {};
  const img = {
    src: "https://example.test/drawings/1021001.webp",
    dataset: { artwork: "true", artFocusX: "48", artFocusY: "36" },
    naturalWidth: 1254,
    naturalHeight: 1254,
    clientWidth: 220,
    clientHeight: 60,
    style: {
      setProperty(k, v) {
        properties[k] = v;
      },
      removeProperty(k) {
        delete properties[k];
      },
    },
    hasAttribute: () => false,
    removeAttribute() {},
    getAttribute: () => img.src,
    setAttribute(k, v) {
      img[k] = v;
    },
  };
  attrs.onLoad({ target: img });
  assert.ok(observed.has(img));
  assert.ok(properties["--art-position"]);
  attrs.onVnodeUnmounted({ el: img });
  assert.ok(!observed.has(img));
  attrs.onLoad({ target: img });
  ui.imageError({ target: img });
  assert.ok(!observed.has(img));
  assert.equal(properties["--art-position"], undefined);
  assert.equal(img["data-art-fallback"], "true");
  const effectAttrs = ui.artAttrs({ number: "1021001" }, false);
  assert.ok(!("onVnodeUnmounted" in effectAttrs));
  effectAttrs.onLoad({ target: img });
  assert.ok(
    !observed.has(img),
    "temporary effect images do not outlive their DOM node in an observer",
  );
});
