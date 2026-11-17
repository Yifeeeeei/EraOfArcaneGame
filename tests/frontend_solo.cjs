/* Real Vue/WS checks. Run against a local Go server:
 * NODE_PATH=<playwright node_modules> ARCANE_URL=http://localhost:9090
 * ARCANE_CHROME=<optional executable> node tests/frontend_solo.cjs
 */
const fs = require('node:fs');
const path = require('node:path');
const { chromium } = require('playwright');
const base = process.env.ARCANE_URL || 'http://localhost:9090';
const out = path.resolve('agent-data/frontend-solo');
fs.mkdirSync(out, { recursive: true });
const recipes = JSON.parse(fs.readFileSync('web/assets/preconstructed-decks.json', 'utf8'));
const deck = recipes.decks.find(d => d.id === '王权纷争/快攻火').deckCode;
async function api(route, data) {
  const response = await fetch(base + route, data === undefined ? {} : {
    method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(data),
  });
  const body = await response.json();
  if (!response.ok) throw Error(JSON.stringify(body));
  return body;
}
function assert(ok, reason) { if (!ok) throw Error(reason); }
async function state(page) { return page.evaluate(() => window.__arcaneDebug?.gameState.value); }
async function commit(page) {
  await page.waitForFunction(() => !window.__arcaneDebug.actionBusy.value);
  if (await page.evaluate(() => !!window.__arcaneDebug.actionDraft.value)) {
    await page.getByRole('button', { name: '确认提交', exact: true }).click();
  }
}
async function keep(page) {
  await page.getByRole('button', { name: '保留手牌', exact: true }).click();
  await commit(page);
}
async function newClient(browser, room, id) {
  const page = await browser.newPage({ viewport: { width: 1440, height: 1000 }, ignoreHTTPSErrors: true });
  page.setDefaultTimeout(15000);
  await instrument(page);
  await page.goto(base + '/game.html?' + new URLSearchParams({ room, player_id: id, player_name: id, deck_code: deck }), { waitUntil: 'domcontentloaded' });
  return page;
}
const errors = [];
async function instrument(page) {
  page.on('pageerror', e => errors.push(e.message));
  await page.addInitScript(() => {
    window.__states = [];
    window.__sent = [];
    window.__wireErrors = [];
    const WS = window.WebSocket;
    window.WebSocket = class extends WS {
      constructor(...args) {
        super(...args);
        this.addEventListener('message', e => {
          const m = JSON.parse(e.data);
          if (m.type === 'game_event' && m.event.type === 'state_sync') {
            const s = m.event.data;
            window.__states.push({ revision: s.revision, phase: s.phase, turn: s.current_turn });
          }
          if (m.type === 'error' || m.type === 'action_result' && !m.ok) window.__wireErrors.push(m);
        });
      }
      send(data) { window.__sent.push(JSON.parse(data)); return super.send(data); }
    };
  });
}
async function consume(page) {
  const s = await state(page);
  await page.locator(`.my-board .unit-cell[data-instance-id="${s.you.hero.instance_id}"] .unit-img`).click();
  await page.locator('.context-menu').getByRole('menuitem', { name: '消耗', exact: true }).click();
  await commit(page);
  await page.waitForFunction(r => window.__arcaneDebug.gameState.value.revision > r, s.revision);
}
async function learnFireball(page) {
  const rev = (await state(page)).revision;
  await page.locator('.dock-tabs').getByRole('button', { name: '技能池', exact: true }).click();
  await page.locator('.sp-card').filter({ hasText: '火球术' }).getByRole('button', { name: '学习', exact: true }).click();
  await commit(page);
  await page.waitForFunction(r => window.__arcaneDebug.gameState.value.revision > r, rev);
  await page.locator('.dock-tabs button').first().click();
}
async function castFireball(page, target) {
  await page.waitForFunction(() => {
    const s = window.__arcaneDebug.gameState.value;
    return s.phase === 'main' && s.decision_player === s.you.player_id && !window.__arcaneDebug.actionBusy.value;
  });
  fs.writeFileSync(path.join(out, 'before-cast.json'), JSON.stringify(await state(page), null, 2));
  await page.screenshot({ path: path.join(out, 'before-cast.png') });
  await page.locator('.my-board .slot').filter({ has: page.locator('img[alt="火球术"]') }).click();
  await page.locator('.context-menu').getByRole('menuitem', { name: '施放', exact: true }).click();
  await page.locator(`.opponent-board .unit-cell[data-instance-id="${target.instance_id}"]`).click();
  await commit(page);
}
async function finishDiscard(page) {
  await page.waitForFunction(() => !window.__arcaneDebug.actionBusy.value);
  const pending = (await state(page)).pending_action;
  if (!pending || pending.type !== 'discard') return false;
  for (let i = 0; i < pending.min_select; i++) await page.locator('.pending-card').nth(i).click();
  await page.getByRole('button', { name: '确认选择', exact: true }).click();
  await commit(page);
  return true;
}
async function waitOwnMain(page) {
  const deadline = Date.now() + 30000;
  while (Date.now() < deadline) {
    const s = await state(page);
    if (s?.phase === 'main' && s.current_turn === s.you.player_id) return;
    if (s?.phase === 'game_over') throw Error('match ended before the intended UI checks');
    if (s?.decision_player === s.you.player_id) {
      if (s.phase === 'defense_window') {
        await page.getByRole('button', { name: '不防御', exact: true }).click();
        await commit(page);
      } else if (s.pending_action) {
        for (let i = 0; i < s.pending_action.min_select; i++) await page.locator('.pending-card').nth(i).click();
        await page.getByRole('button', { name: /^(确认选择|不选择 \/ 放弃发动)$/ }).click();
        await commit(page);
      }
    }
    await page.waitForTimeout(100);
  }
  throw Error('AI/human decision made no progress for 30s: ' + JSON.stringify(await state(page)));
}
async function solo(browser) {
  const page = await browser.newPage({ viewport: { width: 1440, height: 1000 }, ignoreHTTPSErrors: true });
  page.setDefaultTimeout(20000);
  await instrument(page);
  await page.addInitScript(code => {
    localStorage.setItem('eraofarcane.decks.v1', JSON.stringify({ version: 1, selectedId: 'review-fire', decks: [{ id: 'review-fire', name: '测试火预组', code }] }));
    localStorage.setItem('playerName', '练习旅法师');
  }, deck);
  await page.goto(base, { waitUntil: 'domcontentloaded' });
  const name = page.locator('#player-name');
  await name.fill('练习旅法师');
  await page.setViewportSize({ width: 1000, height: 900 });
  assert(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth), 'practice entry overflows minimum lobby width');
  await page.screenshot({ path: path.join(out, 'solo-lobby.png') });
  await page.setViewportSize({ width: 1440, height: 1000 });
  await page.locator('[data-solo-start]').click();
  await page.waitForURL('**/game.html?**');
  await keep(page);
  await waitOwnMain(page);
  let s = await state(page);
  const room = new URL(page.url()).searchParams.get('room');
  assert(s.bot_opponent && s.opponent.hero.number === '4111002', 'wrong AI opponent');
  assert(!s.opponent.hand && !s.opponent.skill_pool, 'hidden AI cards leaked');
  await page.screenshot({ path: path.join(out, 'solo-start.png') });
  await consume(page);
  await learnFireball(page);
  await page.getByRole('button', { name: '结束回合', exact: true }).click();
  await commit(page);
  await finishDiscard(page);
  await waitOwnMain(page);
  await page.reload({ waitUntil: 'domcontentloaded' });
  await page.waitForFunction(() => window.__arcaneDebug?.gameState.value?.bot_opponent);
  assert(new URL(page.url()).searchParams.get('room') === room, 'reload created a different match');
  const info = await api('/api/room/info?id=' + room);
  assert(info.players.filter(p => p.is_bot).length === 1, 'reload duplicated AI');
  await consume(page);
  s = await state(page);
  let target;
  for (let row = 0; row < 3 && !target; row++) for (let col = 0; col < 3 && !target; col++) target = s.opponent.units[col][row];
  await castFireball(page, target);
  await page.waitForFunction(() => window.__states.some(s => s.phase === 'defense_window'));
  await waitOwnMain(page);
  await page.screenshot({ path: path.join(out, 'solo-battle.png') });
  await page.getByRole('button', { name: '结束回合', exact: true }).click();
  await commit(page);
  const discarded = await finishDiscard(page);
  assert(discarded, 'did not exercise normal discard choice');
  await page.getByRole('button', { name: '投降', exact: true }).click();
  await commit(page);
  await page.waitForFunction(() => window.__arcaneDebug.gameState.value.phase === 'game_over');
  assert((await page.evaluate(() => window.__wireErrors)).length === 0, 'solo UI submitted invalid actions');
  await page.screenshot({ path: path.join(out, 'solo-result.png') });
  await page.close();
  console.log('PASS single player: lobby, AI, mulligan, consume, learn, automatic defense, discard, reconnect, game over');
  return room;
}
async function multiplayer(browser) {
  const room = (await api('/api/test-room/create', {})).room_id;
  const pages = [await newClient(browser, room, 'human-a'), await newClient(browser, room, 'human-b')];
  await keep(pages[0]); await keep(pages[1]);
  await pages[0].waitForFunction(() => window.__arcaneDebug.gameState.value.phase === 'main');
  const actor = (await state(pages[0])).current_turn;
  const a = pages[actor], b = pages[1 - actor];
  await consume(a);
  await api('/api/test-room/elements', { room_id: room, player: actor, mode: 'set', elements: { 火: 30 } });
  await a.waitForFunction(() => window.__arcaneDebug.gameState.value.you.elements['火'] === 30);
  await learnFireball(a);
  const spell = (await state(a)).you.skills.find(c => c?.number === '3121001');
  await api('/api/test-room/card-state', { room_id: room, player: actor, instance_id: spell.instance_id, horizontal: false });
  await a.waitForFunction(id => window.__arcaneDebug.gameState.value.you.skills.some(c => c?.instance_id === id && !c.is_horizontal), spell.instance_id);
  for (const player of [0, 1]) await api('/api/test-room/add-card', { room_id: room, player, card_number: '1121112', zone: 'hand', count: 1 });
  await castFireball(a, (await state(a)).opponent.hero);
  await b.waitForFunction(() => window.__arcaneDebug.showDefenseWindow.value);
  await b.getByRole('button', { name: '不防御', exact: true }).click();
  await commit(b);
  // Both preconstructed hands may contain Spark Moth, which legitimately
  // pauses a hit for either player's optional reveal before main resumes.
  for (let attempt = 0; attempt < 100; attempt++) {
    for (const page of pages) {
      const pending = (await state(page)).pending_action;
      if (!pending) continue;
      for (let i = 0; i < pending.min_select; i++) await page.locator('.pending-card').nth(i).click();
      await page.getByRole('button', { name: /^(确认选择|不选择 \/ 放弃发动)$/ }).click();
      await commit(page);
    }
    if ((await state(a)).phase === 'main') break;
    await a.waitForTimeout(100);
  }
  assert((await state(a)).phase === 'main', 'spell did not finish after owned hit choices');
  await api('/api/test-room/add-card', { room_id: room, player: actor, card_number: '1021001', zone: 'hand', count: 8 });
  await a.waitForFunction(() => window.__arcaneDebug.gameState.value.you.hand.length > 5);
  await a.getByRole('button', { name: '结束回合', exact: true }).click(); await commit(a);
  await a.waitForFunction(() => window.__arcaneDebug.gameState.value.pending_action?.type === 'discard');
  assert(await finishDiscard(a), 'multiplayer discard missing');
  await a.getByRole('button', { name: '求和', exact: true }).click(); await commit(a);
  await b.getByRole('button', { name: '接受和局', exact: true }).click(); await commit(b);
  for (const page of pages) {
    await page.waitForFunction(() => window.__arcaneDebug.gameState.value.phase === 'game_over');
    assert(!(await state(page)).bot_opponent, 'normal match became AI match');
    assert((await page.evaluate(() => window.__wireErrors)).length === 0, 'multiplayer submitted invalid actions');
    await page.close();
  }
  console.log('PASS two Vue clients: mulligan, consume, learn, cast, defense, discard, game over');
  return room;
}
(async () => {
  const browser = await chromium.launch({ headless: true, ...(process.env.ARCANE_CHROME ? { executablePath: process.env.ARCANE_CHROME } : {}) });
  try {
    const rooms = [await solo(browser), await multiplayer(browser)];
    assert(errors.length === 0, errors.join('\n'));
    fs.writeFileSync(path.join(out, 'result.json'), JSON.stringify({ passed: true, rooms, page_errors: errors }, null, 2));
  } finally { await browser.close(); }
})().catch(e => { console.error(e); process.exitCode = 1; });
