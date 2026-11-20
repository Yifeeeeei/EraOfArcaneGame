const assert = require('node:assert/strict');
const test = require('node:test');
global.window = global;
require('../assets/card-keyword-index.js');
const { entries } = require('./card-keywords.js');

test('printed explanations ignore instance status and temporary text', () => {
  const plain = entries({ number: '1021107' });
  assert.deepEqual(entries({ number: '1021107', statuses: { 精通: 2, 石化: 1 }, description: '其他临时文字' }), plain);
  assert.deepEqual(entries({ number: '3121001', statuses: { 点燃: 3 } }), []);
  assert.deepEqual(entries({ number: '3121003', is_hidden: true }), []);
});
test('defense restrictions and references remain distinct', () => {
  const ids = (number) => entries({ number }).map((entry) => entry.id);
  assert(ids('3121012').includes('K13'));
  for (const number of ['3521013', '3121101', '3121105']) {
    assert(ids(number).includes('A07'));
    assert(!ids(number).includes('K13'));
  }
  assert(!ids('1121112').includes('K01'));
  assert(!ids('3121101').includes('K24'));
  assert(ids('3321015').includes('K32'));
});
