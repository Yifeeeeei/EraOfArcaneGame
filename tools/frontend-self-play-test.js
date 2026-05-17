const { chromium } = require('/Users/yifeichen/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/node_modules/playwright');
const fs = require('node:fs');

const BASE = process.env.EOA_BASE_URL || 'http://127.0.0.1:9090';
const OUT_PATH = process.env.EOA_SELF_PLAY_REPORT || 'tmp/frontend-self-play-report.json';
const MAX_TURNS = Number(process.env.EOA_SELF_PLAY_MAX_TURNS || '18');
const ONLY = new Set((process.env.EOA_SELF_PLAY_ONLY || '').split(',').map((x) => x.trim()).filter(Boolean));

const games = [
  {
    name: 'fire_vs_water',
    p1: {
      hero: '4111001',
      main: [
        '1121001', '1121001', '1121002', '1121002', '1121003', '1121003',
        '1121004', '1121004', '1121009', '1121009', '1121012', '1121012',
        '1121013', '1121013', '1121014', '1121014', '1111001', '1111003',
        '2111001', '2111002', '2121001', '2121002', '2121004', '2121005',
        '2121007', '2121012', '2021014', '2021015', '1021006', '1021015',
      ],
      skills: ['3101002', '3121002', '3121007', '3121008', '3121014', '3121015', '3021005', '3021007', '3001002', '3321002'],
    },
    p2: {
      hero: '4211002',
      main: [
        '1221001', '1221001', '1221004', '1221004', '1221005', '1221005',
        '1221006', '1221006', '1221008', '1221010', '1221011', '1221012',
        '1221013', '1221014', '1221015', '1211001', '1211002', '1211003',
        '2211002', '2221001', '2221002', '2221004', '2221005', '2221008',
        '2221010', '2221011', '2221012', '2221013', '1021006', '1021016',
      ],
      skills: ['3201001', '3201002', '3221003', '3221007', '3221008', '3221009', '3221010', '3221015', '3021005', '3321002'],
    },
  },
  {
    name: 'earth_vs_shadow',
    p1: {
      hero: '4411003',
      main: [
        '1421001', '1421003', '1421004', '1421007', '1421009', '1421010',
        '1421011', '1421012', '1421014', '1411001', '1411002', '1411003',
        '1401001', '1401002', '2411001', '2411002', '2421001', '2421004',
        '2421006', '2421007', '2421013', '1021006', '1021007', '1021010',
        '1021012', '1021014', '1021015', '1021016', '1021017', '2021015',
      ],
      skills: ['3421003', '3421008', '3421011', '3421012', '3421013', '3421014', '3421015', '3021005', '3021007', '3321002'],
    },
    p2: {
      hero: '4611003',
      main: [
        '1621001', '1621002', '1621003', '1621004', '1621005', '1621006',
        '1621009', '1621010', '1621011', '1621012', '1621013', '1611001',
        '1611002', '1611003', '2601001', '2601002', '2611001', '2611002',
        '2621002', '2621003', '2621004', '2621005', '2621006', '2621008',
        '2621010', '2621011', '2621012', '2621013', '1021006', '1021015',
      ],
      skills: ['3621002', '3621006', '3621007', '3621008', '3621010', '3621012', '3621013', '3621015', '3021005', '3321002'],
    },
  },
  {
    name: 'air_vs_light_arcane',
    p1: {
      hero: '4311002',
      main: [
        '1321001', '1321002', '1321003', '1321004', '1321005', '1321006',
        '1321010', '1321012', '1321013', '1321015', '1311001', '1311002',
        '1311003', '2311001', '2311002', '2321001', '2321002', '2321003',
        '2321005', '2321006', '2321007', '2321008', '2321009', '2321010',
        '2321011', '2321012', '1021006', '1021014', '1021015', '2021015',
      ],
      skills: ['3301001', '3321001', '3321002', '3321008', '3321012', '3321014', '3321015', '3021005', '3021007', '3001002'],
    },
    p2: {
      hero: '4511002',
      main: [
        '1011001', '1011002', '1011003', '1021006', '1021007', '1021008',
        '1021010', '1021012', '1021014', '1021015', '1021016', '1021017',
        '1511001', '1511002', '1521002', '1521005', '1521006', '1521007',
        '1521009', '1521011', '1521013', '1521014', '1521015', '2011001',
        '2011002', '2011003', '2021002', '2021012', '2511001', '2511002',
      ],
      skills: ['3501001', '3521001', '3521007', '3521011', '3521013', '3521014', '3021005', '3021011', '3001001', '3001002'],
    },
  },
];

function deckCode(deck) {
  return `${deck.hero} // ${deck.main.join(' ')} // ${deck.skills.join(' ')}`;
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function visible(locator, timeout = 250) {
  try {
    return await locator.isVisible({ timeout });
  } catch {
    return false;
  }
}

async function resolvePayment(page, stats) {
  const panel = page.locator('.payment-panel').first();
  if (!(await visible(panel, 150))) return false;

  const confirm = panel.getByRole('button', { name: '确认支付' });
  const preferredElement = await page.evaluate(() => {
    const d = window.__arcaneDebug;
    const req = d?.paymentRequest?.value;
    const state = d?.gameState?.value;
    const slot = d?.mySlot?.value;
    const cost = req?.cost || {};
    if (!Number(cost['奥术'] || 0) || !state || slot === undefined) return '';
    const total = Object.values(cost).reduce((sum, value) => sum + Number(value || 0), 0);
    const elements = state.players?.[slot]?.elements || {};
    return ['暗', '光', '气', '地', '水', '火', '无'].find((elem) => Number(elements[elem] || 0) >= total) || '';
  }).catch(() => '');
  for (let i = 0; i < 10; i++) {
    const disabled = await confirm.getAttribute('disabled').catch(() => null);
    if (disabled === null && await visible(confirm, 100)) break;
    const token = preferredElement
      ? panel.locator(`.payment-token:has(img[alt="${preferredElement}"]):not([disabled])`).first()
      : panel.locator('.payment-token:not([disabled])').first();
    if (!(await visible(token, 250))) break;
    await token.click().catch(() => {});
    await sleep(120);
  }

  if (await visible(confirm, 250)) {
    const disabled = await confirm.getAttribute('disabled').catch(() => null);
    if (disabled === null) {
      await confirm.click().catch(() => {});
      stats.paymentsResolved++;
      await sleep(250);
      return true;
    }
  }
  const cancel = panel.getByRole('button', { name: '取消' });
  if (await visible(cancel, 250)) {
    await cancel.click().catch(() => {});
    stats.paymentsCanceled++;
    await sleep(250);
    return true;
  }
  return false;
}

async function createRoom() {
  const resp = await fetch(`${BASE}/api/room/create`, { method: 'POST' });
  if (!resp.ok) throw new Error(`create room failed: ${resp.status}`);
  return (await resp.json()).room_id;
}

async function resolveInterrupts(pages, stats) {
  for (const page of pages) {
    if (await resolvePayment(page, stats)) return true;

    const devourPanel = page.locator('.pending-action-panel').filter({ hasText: '选择吞噬对象' }).first();
    if (await visible(devourPanel, 150)) {
      const candidates = devourPanel.locator('.pending-card');
      if ((await candidates.count().catch(() => 0)) === 0) {
        const cancel = devourPanel.getByRole('button', { name: '取消' });
        if (await visible(cancel, 250)) {
          await cancel.click().catch(() => {});
          stats.devoursCanceled++;
          await sleep(250);
          return true;
        }
      }
    }

    const pending = page.locator('.pending-card');
    if ((await pending.count()) > 0 && await visible(pending.first())) {
      const body = await page.locator('body').innerText().catch(() => '');
      const discardMatch = body.match(/需弃(\d+)张/);
      const muling = body.includes('双方各1个伙伴');
      const needed = discardMatch ? Number(discardMatch[1]) : (muling ? 2 : 1);
      const max = Math.min(await pending.count(), Math.max(1, needed));
      for (let i = 0; i < max; i++) {
        const card = pending.nth(i);
        if (await visible(card, 100)) await card.click().catch(() => {});
      }
      const confirm = page.getByRole('button', { name: '确认' });
      if (await visible(confirm, 500)) {
        await confirm.click().catch(() => {});
        stats.pendingResolved++;
        await sleep(250);
        return true;
      }
    }
    const noDefend = page.getByRole('button', { name: '不防御' });
    if (await visible(noDefend, 250)) {
      await noDefend.click().catch(() => {});
      stats.noDefend++;
      await sleep(250);
      return true;
    }
  }
  return false;
}

async function consumeSome(page, stats) {
  for (let i = 0; i < 8; i++) {
    const clicked = await clickRevealedButton(page, 'button.action-btn.consume');
    if (!clicked) break;
    stats.consumes++;
    await sleep(150);
  }
}

async function clickRevealedButton(page, selector) {
  const buttons = page.locator(selector);
  const total = await buttons.count().catch(() => 0);
  for (let i = 0; i < total; i++) {
    const button = buttons.nth(i);
    const owner = button.locator('xpath=ancestor::*[contains(concat(" ", normalize-space(@class), " "), " unit-cell ") or contains(concat(" ", normalize-space(@class), " "), " slot ")][1]');
    if ((await owner.count().catch(() => 0)) > 0) {
      await owner.first().hover().catch(() => {});
      await sleep(80);
    } else {
      await button.hover().catch(() => {});
      await sleep(80);
    }
    if (await visible(button, 250)) {
      await button.click().catch(() => {});
      return true;
    }
  }
  return false;
}

async function learnOneSkill(page, stats) {
  if ((await page.locator('.slot.occupied').count().catch(() => 0)) >= 5) return false;
  const skill = page.locator('.sp-card.learnable').first();
  if (!(await visible(skill, 300))) return false;
  const text = await skill.innerText().catch(() => '');
  if (/元素不足/.test(text)) return false;
  await skill.click({ button: 'right' }).catch(() => {});
  await sleep(150);
  const learn = page.getByText('学习技能', { exact: false }).first();
  if (!(await visible(learn, 500))) return false;
  await learn.click().catch(() => {});
  stats.learns++;
  await sleep(450);
  return true;
}

async function playOneHandCard(page, stats) {
  const cards = await page.locator('.hand-card.playable').all();
  for (const card of cards) {
    const cls = await card.getAttribute('class').catch(() => '');
    if (cls.includes('unaffordable')) continue;
    await card.click().catch(() => {});
    await sleep(200);
    const summon = page.locator('.unit-cell.summon-target').first();
    if (await visible(summon, 300)) {
      await summon.click().catch(() => {});
      stats.handPlays++;
      await sleep(600);
      return true;
    }
    const terrain = page.locator('.unit-cell.terrain-target').first();
    if (await visible(terrain, 300)) {
      await terrain.click().catch(() => {});
      stats.handPlays++;
      await sleep(600);
      return true;
    }
    await card.click().catch(() => {});
    stats.handPlays++;
    await sleep(600);
    return true;
  }
  return false;
}

async function useOneAbility(page, stats) {
  const clicked = await clickRevealedButton(page, 'button.action-btn.ability');
  if (!clicked) return false;
  stats.abilities++;
  await sleep(450);
  return true;
}

async function castOneSkill(page, stats) {
  const cast = page.locator('.cast-btn').first();
  if (!(await visible(cast, 250))) return false;
  await cast.click().catch(() => {});
  await sleep(200);
  const target = page.locator('.unit-cell.spell-target.occupied').first();
  if (await visible(target, 600)) {
    await target.click().catch(() => {});
    await sleep(200);
    const needsExtraTarget = await page.evaluate(() => {
      const d = window.__arcaneDebug;
      return !!d?.pendingExtraTargetCast?.value && d?.selectedSkill?.value?.number === '3321001';
    }).catch(() => false);
    if (needsExtraTarget) {
      const targets = page.locator('.unit-cell.spell-target.occupied');
      const count = await targets.count().catch(() => 0);
      const extraTarget = count > 1 ? targets.nth(1) : targets.first();
      if (await visible(extraTarget, 600)) {
        await extraTarget.click().catch(() => {});
      }
    }
  }
  stats.casts++;
  await sleep(700);
  return true;
}

async function attackOnce(page, stats) {
  const clicked = await clickRevealedButton(page, 'button.action-btn.attack');
  if (!clicked) return false;
  await sleep(150);
  const target = page.locator('.unit-cell.attack-target.occupied').first();
  if (!(await visible(target, 500))) return false;
  await target.click().catch(() => {});
  stats.attacks++;
  await sleep(500);
  return true;
}

async function gameOver(page) {
  const body = await page.locator('body').innerText().catch(() => '');
  return /游戏结束|获胜|胜利/.test(body);
}

async function playTurn(active, passive, label, stats) {
  const pages = [active, passive];
  if (!(await visible(active.getByRole('button', { name: '结束回合' }), 500))) return false;

  for (let i = 0; i < 3; i++) await resolveInterrupts(pages, stats);
  await attackOnce(active, stats);
  await resolveInterrupts(pages, stats);

  await consumeSome(active, stats);
  await resolveInterrupts(pages, stats);

  for (let i = 0; i < 2; i++) {
    if (!(await learnOneSkill(active, stats))) break;
    await resolveInterrupts(pages, stats);
  }

  await useOneAbility(active, stats);
  await resolveInterrupts(pages, stats);

  for (let i = 0; i < 3; i++) {
    if (!(await playOneHandCard(active, stats))) break;
    await resolveInterrupts(pages, stats);
  }

  for (let i = 0; i < 2; i++) {
    if (!(await castOneSkill(active, stats))) break;
    await resolveInterrupts(pages, stats);
  }

  const end = active.getByRole('button', { name: '结束回合' });
  if (await visible(end, 500)) {
    await end.click().catch(() => {});
    stats.turns[label]++;
    await sleep(450);
    return true;
  }
  return false;
}

async function waitForMain(page, other, stats, timeoutMs = 20000) {
  const start = Date.now();
  while (Date.now() - start < timeoutMs) {
    await resolveInterrupts([page, other], stats);
    if (await visible(page.getByRole('button', { name: '结束回合' }), 250)) return true;
    await sleep(250);
  }
  return false;
}

async function runOne(browser, spec, index) {
  const room = await createRoom();
  const contexts = [
    await browser.newContext({ viewport: { width: 1440, height: 960 } }),
    await browser.newContext({ viewport: { width: 1440, height: 960 } }),
  ];
  for (const context of contexts) {
    context.setDefaultTimeout(5000);
    context.setDefaultNavigationTimeout(10000);
  }
  const pages = [await contexts[0].newPage(), await contexts[1].newPage()];
  const logs = [];
  const errors = [];
  for (const [i, page] of pages.entries()) {
    page.on('console', (msg) => {
      const text = msg.text();
      if (msg.type() === 'error') errors.push(`p${i + 1} console: ${text}`);
      logs.push(`p${i + 1} ${msg.type()}: ${text}`);
    });
    page.on('pageerror', (err) => errors.push(`p${i + 1} pageerror: ${err.message}`));
  }
  const stats = {
    name: spec.name,
    room,
    turns: { p1: 0, p2: 0 },
    consumes: 0,
    learns: 0,
    handPlays: 0,
    casts: 0,
    attacks: 0,
    abilities: 0,
    pendingResolved: 0,
    noDefend: 0,
    paymentsResolved: 0,
    paymentsCanceled: 0,
    devoursCanceled: 0,
    gameOver: false,
    errors,
  };

  try {
    const urls = [
      `${BASE}/game.html?room=${encodeURIComponent(room)}&player_id=self_${index}_p1_${Date.now()}&player_name=${encodeURIComponent('SelfPlay1')}&deck_code=${encodeURIComponent(deckCode(spec.p1))}`,
      `${BASE}/game.html?room=${encodeURIComponent(room)}&player_id=self_${index}_p2_${Date.now()}&player_name=${encodeURIComponent('SelfPlay2')}&deck_code=${encodeURIComponent(deckCode(spec.p2))}`,
    ];
    await Promise.all([
      pages[0].goto(urls[0], { waitUntil: 'domcontentloaded' }),
      pages[1].goto(urls[1], { waitUntil: 'domcontentloaded' }),
    ]);
    await Promise.all([
      pages[0].getByText('选择初始手牌', { exact: false }).waitFor({ state: 'visible', timeout: 12000 }),
      pages[1].getByText('选择初始手牌', { exact: false }).waitFor({ state: 'visible', timeout: 12000 }),
    ]);
    await pages[0].getByRole('button', { name: '保留手牌' }).click();
    await pages[1].getByRole('button', { name: '保留手牌' }).click();
    await waitForMain(pages[0], pages[1], stats);

    for (let i = 0; i < MAX_TURNS; i++) {
      await resolveInterrupts(pages, stats);
      if (await gameOver(pages[0]) || await gameOver(pages[1])) {
        stats.gameOver = true;
        break;
      }
      if (await visible(pages[0].getByRole('button', { name: '结束回合' }), 250)) {
        await playTurn(pages[0], pages[1], 'p1', stats);
      } else if (await visible(pages[1].getByRole('button', { name: '结束回合' }), 250)) {
        await playTurn(pages[1], pages[0], 'p2', stats);
      } else if (!(await waitForMain(pages[0], pages[1], stats, 5000)) && !(await waitForMain(pages[1], pages[0], stats, 5000))) {
        stats.stalled = true;
        break;
      }
    }

    stats.gameOver = stats.gameOver || await gameOver(pages[0]) || await gameOver(pages[1]);
    const bodies = [
      await pages[0].locator('body').innerText().catch(() => ''),
      await pages[1].locator('body').innerText().catch(() => ''),
    ];
    stats.pageStates = [];
    for (const [pageIndex, page] of pages.entries()) {
      stats.pageStates.push({
        page: pageIndex + 1,
        pendingCards: await page.locator('.pending-card').count().catch(() => 0),
        buttons: await page.locator('button').evaluateAll((buttons) => buttons.map((button) => button.textContent.trim()).filter(Boolean).slice(0, 40)).catch(() => []),
        debug: await page.evaluate(() => {
          const d = window.__arcaneDebug;
          const state = d?.gameState?.value;
          const pending = d?.pendingAction?.value;
          return {
            phase: state?.phase,
            currentTurn: state?.current_turn,
            mySlot: d?.mySlot?.value,
            isMyTurn: d?.isMyTurn?.value,
            showPendingAction: d?.showPendingAction?.value,
            pendingType: pending?.type,
            pendingPrompt: pending?.prompt,
            pendingCandidates: pending?.candidates?.map((c) => ({ id: c.instance_id, name: c.name, zone: c.zone, side: c.side })).slice(0, 12),
          };
        }).catch((err) => ({ error: err.message })),
        tail: bodies[pageIndex].slice(-900),
      });
    }
    const body = bodies[0];
    stats.effectLogLines = (body.match(/effect|触发|获得|召唤|装备|施放|使用|伤害|冻结|点燃|眩晕|虚弱/g) || []).length;
    stats.bodyTail = body.slice(-1200);
    stats.recentConsole = logs.slice(-12);
    return stats;
  } finally {
    for (const context of contexts) await context.close().catch(() => {});
  }
}

(async () => {
  fs.mkdirSync('tmp', { recursive: true });
  const browser = await chromium.launch({ headless: true });
  const results = [];
  try {
    const selectedGames = ONLY.size > 0 ? games.filter((game) => ONLY.has(game.name)) : games;
    for (const [i, spec] of selectedGames.entries()) {
      const result = await runOne(browser, spec, i + 1);
      results.push(result);
      console.log(`${i + 1}/${selectedGames.length} ${result.name}\tturns=${result.turns.p1 + result.turns.p2}\thand=${result.handPlays}\tlearn=${result.learns}\tcast=${result.casts}\tattack=${result.attacks}\tability=${result.abilities}\tpending=${result.pendingResolved}\terrors=${result.errors.length}\tstalled=${!!result.stalled}\tgameOver=${!!result.gameOver}`);
    }
  } finally {
    await browser.close();
  }
  const report = { generatedAt: new Date().toISOString(), baseUrl: BASE, maxTurns: MAX_TURNS, results };
  fs.writeFileSync(OUT_PATH, JSON.stringify(report, null, 2));
  console.log(`REPORT ${OUT_PATH}`);
  const failed = results.filter((r) => r.errors.length || r.stalled);
  if (failed.length > 0) process.exitCode = 1;
})().catch((err) => {
  console.error(err);
  process.exit(1);
});
