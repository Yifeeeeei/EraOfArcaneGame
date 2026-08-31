const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');

const html = fs.readFileSync(require.resolve('../game.html'), 'utf8');
const match = html.match(/\/\* ACTION_COST_PLAN_START \*\/([\s\S]*?)\/\* ACTION_COST_PLAN_END \*\//);
assert.ok(match, 'game.html should contain the action cost planner');
const browserScope = { window: {} };
vm.runInNewContext(match[1], browserScope);
const { planActionCost } = browserScope.window.ArcaneActionCostPlan;

function plain(value) {
    return JSON.parse(JSON.stringify(value));
}

function card(instanceID, fields) {
    return { instance_id: instanceID, ...fields };
}

const discount = [{
    type: 'next_item_or_skill_cost_minus',
    element: '水',
    amount: 3,
    remaining_uses: 1,
}];

const hail = card('hail', {
    type: '技能',
    action_base_attack_cost: { 水: 2 },
    action_base_attack_boost_cost: { 水: 2 },
    effective_elements_cost: { 水: 2 },
});
const waterScroll = card('water-scroll', {
    type: '道具',
    is_consumable: true,
    action_base_play_cost: { 水: 2 },
    effective_elements_cost: {},
});
const iceScroll = card('ice-scroll', {
    type: '道具',
    is_consumable: true,
    action_base_play_cost: { 水: 4 },
    effective_elements_cost: { 水: 1 },
});

assert.deepEqual(plain(planActionCost([
    { kind: 'skill', purpose: 'attack', card: hail },
    { kind: 'card', card: waterScroll },
    { kind: 'card', card: iceScroll },
], discount)), { 水: 6 });

const defense = card('defense', {
    type: '技能',
    action_base_defense_cost: { 水: 1 },
});
assert.deepEqual(plain(planActionCost([
    { kind: 'skill', purpose: 'defend', card: defense },
    { kind: 'card', card: waterScroll },
], discount)), { 水: 2 });

const freeNextSkill = [{ type: 'next_skill_cost_zero', remaining_uses: 1 }];
assert.deepEqual(plain(planActionCost([
    { kind: 'skill', purpose: 'attack_boost', card: hail },
], freeNextSkill)), { 水: 2 });

console.log('game-costs tests passed');
