const assert = require('node:assert/strict');
const test = require('node:test');
const { permits, defenseHint } = require('./tutorial.js');
test('reading acknowledgements are possible while all game writes stay locked', () => {
  const course = { reading: true, allowed: { action: 'tutorial_continue' } };
  assert(permits(course, 'tutorial_continue'));
  assert(!permits(course, 'end_turn'));
  assert(!permits({ ...course, waiting: true }, 'tutorial_continue'));
});
test('defense may select only the authored spell and tap source', () => {
  const course = { allowed: { action: 'defend', data: { skill_ids: ['shield'], overexert_ids: ['mask'] } } };
  assert(permits(course, 'defend', 'shield'));
  assert(permits(course, 'defend', 'mask'));
  assert(!permits(course, 'defend', 'hero'));
  assert(!permits(course, 'no_defend'));
  assert(permits(null, 'cast_spell', 'anything'));
});
test('defense guidance follows the selected card, skip, tap and final payment', () => {
  const course = { allowed: { action: 'defend', data: { skill_ids: ['shield'], overexert_ids: ['mask'] } } };
  assert.equal(defenseHint(course, 0).cardId, 'shield');
  assert.match(defenseHint(course, 1, ['shield']).instruction, /跳过强化/);
  assert.equal(defenseHint(course, 2, ['shield']).cardId, 'mask');
  assert.equal(defenseHint(course, 2, ['shield']).canReview, false);
  const selected = defenseHint(course, 2, ['shield'], ['mask']);
  assert.equal(selected.canReview, true);
  assert.equal(selected.cardId, undefined);
  assert.match(selected.instruction, /核对付款/);
  assert.equal(defenseHint(course, 3, ['shield'], ['mask']).selector, '[data-defense-action="submit"]');
  assert.equal(defenseHint({ ...course, reading: true }, 0), null);
  assert.equal(defenseHint(null, 0), null);
});
