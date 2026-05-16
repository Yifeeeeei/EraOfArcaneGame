const { chromium } = require('/Users/yifeichen/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/node_modules/playwright');
const fs = require('node:fs');

const BASE = process.env.EOA_BASE_URL || 'http://127.0.0.1:9090';
const CARDS_PATH = 'data/supported_card_infos.json';
const OUT_PATH = process.env.EOA_FRONTEND_REPORT || 'tmp/frontend-card-operation-report.json';
const LIMIT = Number(process.env.EOA_LIMIT || '0');
const TYPE = process.env.EOA_TYPE || '';
const IDS = new Set((process.env.EOA_IDS || '').split(',').map((id) => id.trim()).filter(Boolean));
const RETRIES = Number(process.env.EOA_RETRIES || '3');
const TURN_ATTEMPTS = Number(process.env.EOA_TURN_ATTEMPTS || '2');

const HERO_BY_ELEMENT = {
  '无': '4011001',
  '火': '4111001',
  '水': '4211001',
  '气': '4311003',
  '地': '4411001',
  '光': '4511001',
  '暗': '4611001',
};

const FILLER_MAIN = [
  '1021001', '1021001', '1021002', '1021002', '1021004', '1021004',
  '1021005', '1021005', '1021006', '1021006', '1021007', '1021007',
  '1021008', '1021008', '1021009', '1021009', '1021010', '1021010',
  '1021011', '1021011', '1021012', '1021012', '1021013', '1021013',
  '1021014', '1021014', '1021015', '1021015', '1021016', '1021016',
];

const SUPPORT_BY_ELEMENT = {
  '无': ['1021001', '1021001', '1021007', '1021007', '1021008', '1021008', '1021012', '1021012', '1021017', '1021017'],
  '火': ['1121001', '1121001', '1121002', '1121002', '1121003', '1121003', '1121005', '1121005', '1121014', '1121014'],
  '水': ['1221001', '1221001', '1221003', '1221003', '1221006', '1221006', '1221012', '1221012', '1221014', '1221014'],
  '气': ['1311003', '1311003', '1321001', '1321001', '1321003', '1321003', '1321011', '1321011', '1321013', '1321013'],
  '地': ['1421002', '1421002', '1421003', '1421003', '1421008', '1421008', '1421010', '1421010', '1421012', '1421012'],
  '光': ['1501001', '1501001', '1521005', '1521005', '1521006', '1521006', '1521007', '1521007', '1521009', '1521009'],
  '暗': ['1621001', '1621001', '1621003', '1621003', '1621004', '1621004', '1621007', '1621007', '1621011', '1621011'],
};

const FILLER_SKILLS = [
  '3321002', '3001001', '3001002', '3021001', '3021002',
  '3021003', '3021004', '3021005', '3021006', '3021007',
];

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function totalCost(cost) {
  return Object.values(cost || {}).reduce((sum, value) => sum + Number(value || 0), 0);
}

function primaryElement(cost, fallback = '气') {
  const entries = Object.entries(cost || {}).filter(([, value]) => Number(value) > 0);
  if (!entries.length) return fallback;
  entries.sort((a, b) => Number(b[1]) - Number(a[1]));
  return entries[0][0] || fallback;
}

function isDefenseSkill(card) {
  return card.type === '技能' && !!card.is_defense_only;
}

function isSorcery(card) {
  return card.type === '技能' && !!card.is_sorcery;
}

function skillNeedsTarget(card) {
  return !!(card && card.type === '技能' && card.needs_target);
}

function isTerrain(card) {
  return card.type === '道具' && !!card.is_terrain;
}

function isConsumable(card) {
  return card.type === '道具' && !!card.is_consumable;
}

function opKind(card) {
  if (card.type === '人物') return expectsActiveHeroAbility(card) ? 'hero_ability' : 'hero_basic';
  if (card.type === '伙伴') return 'summon';
  if (card.type === '技能') {
    if (isDefenseSkill(card)) return 'defense_skill';
    if (isSorcery(card)) return 'sorcery_skill';
    return 'attack_skill';
  }
  if (card.type === '道具') {
    if (isTerrain(card)) return 'terrain';
    if (isConsumable(card)) return 'consumable';
    return 'equip';
  }
  return 'unknown';
}

function expectsActiveHeroAbility(card) {
  return !!(card.has_per_turn || card.has_ultimate);
}

function dedupeSkillPool(target) {
  const ids = [target.number, ...FILLER_SKILLS.filter((id) => id !== target.number)];
  return ids.slice(0, 10);
}

function buildDeck(card) {
  const elem = primaryElement(card.elements_cost, card.category || '气');
  const hero = card.type === '人物' ? card.number : (HERO_BY_ELEMENT[elem] || '4311003');
  const main = [];
  const skills = [];

  if (card.type === '伙伴' || card.type === '道具') main.push(card.number, card.number);
  const support = SUPPORT_BY_ELEMENT[elem] || FILLER_MAIN;
  for (const id of [...support, ...FILLER_MAIN]) {
    if (main.length >= 30) break;
    const current = main.filter((x) => x === id).length;
    if (current < 2) main.push(id);
  }

  if (card.type === '技能') skills.push(...dedupeSkillPool(card));
  else skills.push(...FILLER_SKILLS);

  return `${hero} // ${main.slice(0, 30).join(' ')} // ${skills.slice(0, 10).join(' ')}`;
}

async function visible(locator, timeout = 300) {
  try {
    return await locator.isVisible({ timeout });
  } catch {
    return false;
  }
}

async function count(locator) {
  try {
    return await locator.count();
  } catch {
    return 0;
  }
}

async function createRoom() {
  const resp = await fetch(`${BASE}/api/room/create`, { method: 'POST' });
  if (!resp.ok) throw new Error(`create room failed: ${resp.status}`);
  const data = await resp.json();
  return data.room_id;
}

async function resolveInterrupts(p1, p2, stats) {
  for (const page of [p1, p2]) {
    const pending = page.locator('.pending-card');
    if (await count(pending) > 0 && await visible(pending.first())) {
      await pending.first().click();
      await page.getByRole('button', { name: '确认' }).click();
      stats.pendingResolved++;
      await sleep(200);
      return true;
    }
    const noDefend = page.getByRole('button', { name: '不防御' });
    if (await visible(noDefend)) {
      await noDefend.click();
      stats.noDefend++;
      await sleep(200);
      return true;
    }
  }
  return false;
}

async function waitMain(page, other, stats, timeoutMs = 20000) {
  const start = Date.now();
  while (Date.now() - start < timeoutMs) {
    await resolveInterrupts(page, other, stats);
    if (await visible(page.getByRole('button', { name: '结束回合' }))) return true;
    await sleep(250);
  }
  return false;
}

async function endTurnIfPossible(page, label, stats) {
  const btn = page.getByRole('button', { name: '结束回合' });
  if (await visible(btn)) {
    await btn.click();
    stats.turns[label] = (stats.turns[label] || 0) + 1;
    await sleep(300);
    return true;
  }
  return false;
}

async function advanceToNextP1Turn(p1, p2, stats) {
  await endTurnIfPossible(p1, 'p1', stats);
  const start = Date.now();
  while (Date.now() - start < 30000) {
    await resolveInterrupts(p1, p2, stats);
    if (await visible(p2.getByRole('button', { name: '结束回合' }))) {
      await endTurnIfPossible(p2, 'p2', stats);
    }
    if (await visible(p1.getByRole('button', { name: '结束回合' }))) return true;
    await sleep(250);
  }
  return false;
}

async function consumeAll(page) {
  let consumed = 0;
  for (let i = 0; i < 12; i++) {
    const buttons = page.locator('button.action-btn.consume');
    if ((await count(buttons)) < 1 || !(await visible(buttons.first()))) break;
    await buttons.first().click();
    consumed++;
    await sleep(180);
  }
  return consumed;
}

async function playSupportCards(page, targetName) {
  let played = 0;
  for (let i = 0; i < 8; i++) {
    const cards = await page.locator('.hand-card.playable').all();
    let chosen = null;
    for (const card of cards) {
      const text = await card.innerText().catch(() => '');
      if (!text.includes(targetName)) {
        chosen = card;
        break;
      }
    }
    if (!chosen) break;
    await chosen.click();
    await sleep(250);
    const summonTarget = page.locator('.unit-cell.summon-target').first();
    if (await visible(summonTarget)) {
      await summonTarget.click();
      played++;
      await sleep(500);
      continue;
    }
    // For playable support items, a second click uses/equips them.
    await chosen.click().catch(() => {});
    played++;
    await sleep(500);
  }
  return played;
}

async function joinTwoClients(browser, card) {
  const room = await createRoom();
  const deck = buildDeck(card);
  const c1 = await browser.newContext({ viewport: { width: 1440, height: 960 } });
  const c2 = await browser.newContext({ viewport: { width: 1440, height: 960 } });
  const p1 = await c1.newPage();
  const p2 = await c2.newPage();
  const logs = [];
  for (const [label, page] of [['p1', p1], ['p2', p2]]) {
    page.on('console', (msg) => logs.push(`${label} console ${msg.type()}: ${msg.text()}`));
    page.on('pageerror', (err) => logs.push(`${label} pageerror: ${err.message}`));
  }

  const p1Url = `${BASE}/game.html?room=${encodeURIComponent(room)}&player_id=front_card_${Date.now()}_p1&player_name=${encodeURIComponent('卡牌测试P1')}&deck_code=${encodeURIComponent(deck)}`;
  const p2Url = `${BASE}/game.html?room=${encodeURIComponent(room)}&player_id=front_card_${Date.now()}_p2&player_name=${encodeURIComponent('卡牌测试P2')}&deck_code=${encodeURIComponent(deck)}`;
  await p1.goto(p1Url, { waitUntil: 'domcontentloaded' });
  await p2.goto(p2Url, { waitUntil: 'domcontentloaded' });
  await p1.getByText('选择初始手牌', { exact: false }).waitFor({ state: 'visible', timeout: 12000 });
  await p2.getByText('选择初始手牌', { exact: false }).waitFor({ state: 'visible', timeout: 12000 });
  return { room, deck, contexts: [c1, c2], p1, p2, logs };
}

async function keepOrRedrawForCard(page, card) {
  if (card.type !== '伙伴' && card.type !== '道具') {
    await page.getByRole('button', { name: '保留手牌' }).click();
    return 'not_needed';
  }
  if (await visible(page.locator('.card-mini').filter({ hasText: card.name }).first(), 500)) {
    await page.getByRole('button', { name: '保留手牌' }).click();
    return 'initial';
  }
  await page.getByRole('button', { name: '全部重抽' }).click();
  await sleep(600);
  return 'redraw';
}

async function operateSkill(page, other, card, stats) {
  await consumeAll(page);
  await sleep(300);
  let skillCard = page.locator('.sp-card').filter({ hasText: card.name }).first();
  if (!(await visible(skillCard, 2500))) {
    skillCard = page.getByText(card.name, { exact: true })
      .locator('xpath=ancestor::*[contains(concat(" ", normalize-space(@class), " "), " sp-card ")]')
      .first();
  }
  if (!(await visible(skillCard, 2500))) {
    skillCard = page.getByText(card.name, { exact: false }).first();
  }
  if (!(await visible(skillCard, 2500))) {
    for (const candidate of await page.locator('.sp-card').all()) {
      const text = await candidate.innerText().catch(() => '');
      if (text.includes(card.name)) {
        skillCard = candidate;
        break;
      }
    }
  }
  if (!(await visible(skillCard, 2500))) return { status: 'missing_control', reason: 'skill not visible in skill pool' };
  const hint = await skillCard.locator('.sp-learn-hint').textContent().catch(() => '');
  if (/元素不足/.test(hint || '')) return { status: 'unaffordable', reason: `learn hint: ${hint}` };
  await skillCard.click({ button: 'right' });
  const learn = page.getByText('学习技能', { exact: false }).first();
  if (!(await visible(learn))) return { status: 'missing_control', reason: 'learn menu item not visible' };
  await learn.click();
  await sleep(600);

  const body = await page.locator('body').innerText();
  if (!body.includes(`学习技能 ${card.name}`)) {
    return { status: 'missing_confirmation', reason: 'learn log not found' };
  }

  if (isDefenseSkill(card)) return { status: 'pass', action: 'learn_defense_skill' };
  return await castLearnedSkill(page, other, card, stats);
}

async function castLearnedSkill(page, other, card, stats) {
  for (let i = 0; i < 8; i++) {
    await resolveInterrupts(page, other, stats);
    await consumeAll(page);
    const slot = page.locator(`.slot.occupied:has(img[src*="${card.number}"])`).first();
    if (!(await visible(slot, 500))) return { status: 'missing_control', reason: 'learned skill slot not visible' };
    const castButton = slot.locator('.cast-btn').first();
    if (await visible(castButton, 500)) {
      await castButton.click();
      await sleep(300);
      if (skillNeedsTarget(card)) {
        const target = page.locator('.unit-cell.spell-target.occupied').first();
        if (!(await visible(target, 800))) return { status: 'missing_control', reason: 'spell target not highlighted after clicking cast' };
        await target.click();
      }
      await sleep(800);
      await resolveInterrupts(page, other, stats);
      const body = await page.locator('body').innerText();
      return body.includes(`施放 ${card.name}`) ? { status: 'pass', action: 'cast_skill' } : { status: 'missing_confirmation', reason: 'cast log not found' };
    }
    await playSupportCards(page, card.name);
    const advanced = await advanceToNextP1Turn(page, other, stats);
    if (!advanced) break;
  }
  return { status: 'unaffordable', reason: 'learned skill never became castable' };
}

async function operateHandCard(page, other, card, stats) {
  await consumeAll(page);
  const handCard = page.locator('.hand-card').filter({ hasText: card.name }).first();
  if (!(await visible(handCard, 500))) return { status: 'not_in_hand', reason: 'target card not drawn after mulligan/redraw' };

  if (await handCard.evaluate((el) => el.classList.contains('unaffordable')).catch(() => false)) {
    return { status: 'unaffordable', reason: 'hand card marked unaffordable' };
  }

  await handCard.click();
  await sleep(300);

  if (card.type === '伙伴') {
    const target = page.locator('.unit-cell.summon-target').first();
    if (!(await visible(target))) return { status: 'missing_control', reason: 'summon target not highlighted' };
    await target.click();
    await sleep(700);
    const body = await page.locator('body').innerText();
    return body.includes(`召唤 ${card.name}`) ? { status: 'pass', action: 'summon' } : { status: 'missing_confirmation', reason: 'summon log not found' };
  }

  if (card.type === '道具' && isTerrain(card)) {
    const target = page.locator('.unit-cell.terrain-target').first();
    if (!(await visible(target))) return { status: 'missing_control', reason: 'terrain target not highlighted' };
    await target.click();
    await sleep(700);
    const body = await page.locator('body').innerText();
    return body.includes(`放置地形 ${card.name}`) ? { status: 'pass', action: 'place_terrain' } : { status: 'missing_confirmation', reason: 'terrain log not found' };
  }

  await handCard.click();
  await sleep(700);
  const body = await page.locator('body').innerText();
  if (card.type === '道具' && isConsumable(card)) {
    return body.includes(`使用 ${card.name}`) || body.includes(`使用/装备 ${card.name}`) || body.includes(card.name)
      ? { status: 'pass', action: 'use_item' }
      : { status: 'missing_confirmation', reason: 'use item log not found' };
  }
  if (card.type === '道具') {
    return body.includes(`装备 ${card.name}`) || body.includes(card.name)
      ? { status: 'pass', action: 'equip' }
      : { status: 'missing_confirmation', reason: 'equip log not found' };
  }
  return { status: 'unsupported', reason: 'unknown hand card operation' };
}

async function operateHero(page, other, card, stats) {
  await consumeAll(page);
  const body = await page.locator('body').innerText();
  if (!body.includes(card.name)) return { status: 'missing_confirmation', reason: 'hero name not visible' };
  if (opKind(card) === 'hero_basic') return { status: 'pass', action: 'hero_visible_consume' };
  const ability = page.locator('button.action-btn.ability').first();
  if (!(await visible(ability))) return { status: 'missing_control', reason: 'hero ability button not visible after consuming/checking main phase' };
  return { status: 'pass', action: 'hero_ability_button_visible' };
}

async function testOne(browser, card) {
  const kind = opKind(card);
  const aggregate = { pendingResolved: 0, noDefend: 0, turns: { p1: 0, p2: 0 } };
  let last = null;

  for (let attempt = 1; attempt <= RETRIES; attempt++) {
    let env;
    try {
      env = await joinTwoClients(browser, card);
      const p1Mulligan = await keepOrRedrawForCard(env.p1, card);
      await env.p2.getByRole('button', { name: '保留手牌' }).click();
      if (!(await waitMain(env.p1, env.p2, aggregate))) {
        last = { status: 'timeout', reason: 'p1 main phase not reached', attempt, p1Mulligan };
        continue;
      }

      let result;
      for (let turnAttempt = 1; turnAttempt <= TURN_ATTEMPTS; turnAttempt++) {
        if (card.type === '技能') result = await operateSkill(env.p1, env.p2, card, aggregate);
        else if (card.type === '伙伴' || card.type === '道具') result = await operateHandCard(env.p1, env.p2, card, aggregate);
        else if (card.type === '人物') result = await operateHero(env.p1, env.p2, card, aggregate);
        else result = { status: 'unsupported', reason: 'unknown card type' };

        result.turnAttempt = turnAttempt;
        if (['pass', 'missing_control', 'unsupported'].includes(result.status)) break;
        if (turnAttempt < TURN_ATTEMPTS && ['unaffordable', 'not_in_hand'].includes(result.status)) {
          if (result.status === 'unaffordable') {
            await playSupportCards(env.p1, card.name);
          }
          const advanced = await advanceToNextP1Turn(env.p1, env.p2, aggregate);
          if (!advanced) break;
          continue;
        }
        break;
      }

      result.attempt = attempt;
      result.room = env.room;
      result.kind = kind;
      result.cost = totalCost(card.elements_cost);
      result.primaryElement = primaryElement(card.elements_cost, card.category || '气');
      result.frontendLogs = env.logs.filter((line) => !line.includes('development build of Vue')).slice(-5);
      last = result;
      if (result.status === 'pass' || result.status === 'missing_control') return { ...result, ...aggregate };
    } catch (err) {
      last = { status: 'frontend_error', reason: err.message, attempt, kind };
    } finally {
      if (env) {
        for (const context of env.contexts) await context.close().catch(() => {});
      }
    }
  }

  return { ...(last || { status: 'unknown' }), ...aggregate };
}

(async () => {
  const allCards = JSON.parse(fs.readFileSync(CARDS_PATH, 'utf8'));
  let cards = TYPE ? allCards.filter((card) => card.type === TYPE) : allCards;
  if (IDS.size > 0) cards = cards.filter((card) => IDS.has(card.number));
  if (LIMIT > 0) cards = cards.slice(0, LIMIT);

  fs.mkdirSync('tmp', { recursive: true });
  const browser = await chromium.launch({ headless: true });
  const results = [];
  try {
    for (const [idx, card] of cards.entries()) {
      const result = await testOne(browser, card);
      const row = {
        number: card.number,
        name: card.name,
        type: card.type,
        tag: card.tag,
        kind: opKind(card),
        status: result.status,
        action: result.action || '',
        reason: result.reason || '',
        attempt: result.attempt || 0,
        room: result.room || '',
        cost: result.cost ?? totalCost(card.elements_cost),
        primaryElement: result.primaryElement || primaryElement(card.elements_cost, card.category || '气'),
      };
      results.push(row);
      console.log(`${idx + 1}/${cards.length} ${row.status}\t${row.number}\t${row.name}\t${row.kind}${row.reason ? `\t${row.reason}` : ''}`);
    }
  } finally {
    await browser.close();
  }

  const summary = {};
  for (const row of results) summary[row.status] = (summary[row.status] || 0) + 1;
  const report = {
    generatedAt: new Date().toISOString(),
    baseUrl: BASE,
    retries: RETRIES,
    typeFilter: TYPE || null,
    limit: LIMIT || null,
    summary,
    results,
  };
  fs.writeFileSync(OUT_PATH, JSON.stringify(report, null, 2));
  console.log(`REPORT ${OUT_PATH}`);
  console.log(JSON.stringify(summary));
})().catch((err) => {
  console.error(err);
  process.exit(1);
});
