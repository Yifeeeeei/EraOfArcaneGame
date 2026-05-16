const { chromium } = require('/Users/yifeichen/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/node_modules/playwright');

const BASE = process.env.EOA_BASE_URL || 'http://127.0.0.1:9090';
const ATTEMPTS = Number(process.env.EOA_ATTEMPTS || '30');
const DECK = '4111001 // 2121001 2121001 1121001 1121001 1121002 1121002 1121003 1121003 1121005 1121005 1121014 1121014 1021001 1021001 1021002 1021002 1021004 1021004 1021005 1021005 1021006 1021006 1021007 1021007 1021008 1021008 1021009 1021009 1021010 1021010 // 3021001 3021002 3021003 3021004 3021005 3021006 3021007 3021008 3021009 3121001';

async function visible(locator, timeout = 500) {
  try {
    return await locator.isVisible({ timeout });
  } catch {
    return false;
  }
}

async function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function createRoom() {
  const resp = await fetch(`${BASE}/api/room/create`, { method: 'POST' });
  if (!resp.ok) throw new Error(`create room failed: ${resp.status}`);
  const data = await resp.json();
  return data.room_id;
}

async function join(page, room, suffix) {
	const url = `${BASE}/game.html?room=${room}&player_id=phoenix_${Date.now()}_${suffix}&player_name=${suffix}&deck_code=${encodeURIComponent(DECK)}`;
	await page.goto(url, { waitUntil: 'domcontentloaded' });
}

async function keepOrRedraw(page) {
  if (await visible(page.locator('.card-mini').filter({ hasText: '凤凰之羽' }).first(), 700)) {
    await page.getByRole('button', { name: '保留手牌' }).click();
    return true;
  }
  await page.getByRole('button', { name: '全部重抽' }).click();
  await sleep(700);
  return await visible(page.locator('.hand-card').filter({ hasText: '凤凰之羽' }).first(), 1200);
}

async function consumeUntilPlayable(page) {
  for (let i = 0; i < 8; i++) {
    const phoenix = page.locator('.hand-card').filter({ hasText: '凤凰之羽' }).first();
    if (await visible(phoenix, 500)) {
      const cls = await phoenix.getAttribute('class');
      if (!/\bunaffordable\b/.test(cls || '')) return;
    }
    const consume = page.locator('button.action-btn.consume').first();
    if (!(await visible(consume, 500))) return;
    await consume.click();
    await sleep(250);
  }
}

async function fireCount(page) {
	const text = await page.locator('.player-panel.is-active .elem-item').filter({ hasText: '火焰' }).locator('.elem-count').first().textContent().catch(() => '0');
	return Number(text || 0);
}

async function runOnce(browser) {
  const room = await createRoom();
  const c1 = await browser.newContext({ viewport: { width: 1440, height: 960 } });
  const c2 = await browser.newContext({ viewport: { width: 1440, height: 960 } });
  const p1 = await c1.newPage();
  const p2 = await c2.newPage();
	try {
		await join(p1, room, 'P1');
		await join(p2, room, 'P2');
		await p1.getByText('选择初始手牌', { exact: false }).waitFor({ state: 'visible', timeout: 12000 });
		await p2.getByText('选择初始手牌', { exact: false }).waitFor({ state: 'visible', timeout: 12000 });
		const hasPhoenix = await keepOrRedraw(p1);
    await p2.getByRole('button', { name: '保留手牌' }).click();
    if (!hasPhoenix) return { status: 'not_in_hand' };

    await p1.getByRole('button', { name: '结束回合' }).waitFor({ state: 'visible', timeout: 15000 });
    await consumeUntilPlayable(p1);
    const phoenix = p1.locator('.hand-card').filter({ hasText: '凤凰之羽' }).first();
    if (!(await visible(phoenix, 1000))) return { status: 'missing_card' };
    await phoenix.click();
    await sleep(250);
    await phoenix.click();
    await sleep(700);

	const slot = p1.locator('.slot.occupied').filter({ has: p1.locator('img[src*="2121001"]') }).first();
	if (!(await visible(slot, 1200))) return { status: 'equip_failed' };

	const ability = slot.locator('button.action-btn.ability').filter({ hasText: '回合技' }).first();
	if (!(await visible(ability, 1200))) return { status: 'ability_button_missing' };
	const beforeFire = await fireCount(p1);
	await ability.click();
	await sleep(700);

	const body = await p1.locator('body').innerText();
	if (!body.includes('凤凰之羽 使用回合技')) return { status: 'ability_log_missing' };
	const afterFire = await fireCount(p1);
	if (afterFire <= beforeFire) return { status: 'fire_not_gained', beforeFire, afterFire };
	return { status: 'pass' };
  } finally {
    await c1.close().catch(() => {});
    await c2.close().catch(() => {});
  }
}

(async () => {
  const browser = await chromium.launch({ headless: true });
  try {
    let last = null;
	for (let attempt = 1; attempt <= ATTEMPTS; attempt++) {
		last = await runOnce(browser);
		console.log(`${attempt}/${ATTEMPTS} ${last.status}`);
		if (last.status === 'pass') return;
	}
    process.exitCode = 1;
    console.error(`Phoenix feather scenario failed: ${JSON.stringify(last)}`);
  } finally {
    await browser.close();
  }
})();
