/* Real Vue/WebSocket lesson, from lobby entry through normal engine victory.
 * NODE_PATH=<playwright modules> ARCANE_URL=http://localhost:19095 node tests/frontend_tutorial.cjs
 */
const fs = require('node:fs');
const path = require('node:path');
const { chromium } = require('playwright');
const base = process.env.ARCANE_URL || 'http://localhost:9090';
const out = path.resolve('agent-data/frontend-tutorial');
fs.mkdirSync(out, { recursive: true });
const assert = (ok, why) => { if (!ok) throw Error(why); };

(async () => {
  const browser = await chromium.launch({ headless: true, executablePath: process.env.ARCANE_CHROME || 'C:/Program Files/Google/Chrome/Application/chrome.exe' });
  const page = await browser.newPage({ viewport: { width: 1920, height: 1080 }, ignoreHTTPSErrors: true });
  page.setDefaultTimeout(12000);
  const errors = [], wireErrors = [], reached = [];
  page.on('pageerror', error => errors.push(error.message));
  page.on('websocket', socket => socket.on('framereceived', frame => {
    try { const m = JSON.parse(frame.payload); if (m.type === 'error' || m.type === 'action_result' && !m.ok) wireErrors.push(m.message); } catch {}
  }));
  const state = () => page.evaluate(() => window.__arcaneDebug?.gameState.value);
  const course = async () => (await state())?.tutorial;
  async function waitStep(index) {
    await page.waitForFunction(index => window.__arcaneDebug?.gameState.value?.tutorial?.index !== index, index);
  }
  async function reviewAndPay(expected) {
    await page.waitForTimeout(70);
    const tokens = page.locator('.payment-panel [data-payment-element]');
    if (await tokens.count()) {
      for (const [element, amount] of Object.entries(expected.payment || {})) {
        for (let i = 0; i < amount; i++) await page.locator(`.payment-panel [data-payment-element="${element}"]`).click();
      }
      await page.getByRole('button', {name:'支付并提交',exact:true}).click();
    }
    if (await page.getByRole('button', {name:'确认提交',exact:true}).isVisible() && await page.getByRole('button', {name:'确认提交',exact:true}).isEnabled()) {
      await page.getByRole('button', {name:'确认提交',exact:true}).click();
    }
  }
  async function menu(cardID, label, scope = '.my-board') {
    const card = page.locator(`${scope} [data-instance-id="${cardID}"]`).first();
    await card.click();
    await page.getByRole('menuitem', {name:label,exact:true}).click();
  }
  try {
    await page.goto(base + '/');
    assert(await page.locator('.mode-choices button').count() === 3, 'home must contain exactly three mode entries');
    assert(await page.locator('#player-name, .room-browser, [data-solo-start], [data-tutorial-start]').count() === 0, 'home exposes mode setup or rooms');
    await page.screenshot({path:path.join(out,'lobby-home.png'),fullPage:true});
    await page.locator('[data-mode="multiplayer"]').click();
    await page.locator('.room-browser').waitFor();
    assert(await page.locator('#player-name').count() === 1, 'multiplayer lost deck/name setup');
    assert(await page.locator('[data-solo-start], [data-tutorial-start]').count() === 0, 'multiplayer contains other modes');
    await page.screenshot({path:path.join(out,'lobby-multiplayer.png'),fullPage:true});
    await page.reload();
    await page.locator('.room-browser').waitFor();
    await page.getByRole('button', {name:'← 返回大厅',exact:true}).click();
    await page.locator('.mode-choices').waitFor();
    await page.goBack();
    await page.locator('.room-browser').waitFor();
    await page.goForward();
    await page.locator('.mode-choices').waitFor();
    await page.locator('[data-mode="tutorial"]').click();
    assert(await page.locator('#player-name, .room-browser').count() === 0, 'tutorial requires multiplayer setup');
    await page.locator('[data-tutorial-start]').click();
    await page.waitForFunction(() => window.__arcaneDebug?.gameState.value?.tutorial);
    const lessonURL = page.url();
    let refreshedOpening = false, refreshedCombat = false, checkedHelp = false, checkedEmpty = false;
    for (let guard = 0; guard < 120; guard++) {
      const c = await course();
      if (c.completed) break;
      assert(!c.failure, c.failure || 'missing lesson');
      assert(!c.text.includes('本课'), 'outdated tutorial wording');
      if (c.waiting) { await waitStep(c.index); continue; }
      reached.push({ index: c.index, title: c.title, action: c.allowed.action });
      console.log(`${c.index}: ${c.title} [${c.allowed.action}]`);
      if (!refreshedOpening && c.index === 1) {
        await page.reload();
        await page.waitForFunction(() => window.__arcaneDebug?.gameState.value?.tutorial?.index === 1);
        refreshedOpening = true;
      }
      if (!refreshedCombat && c.title === '透支的卡还没有重置') {
        await page.reload();
        await page.waitForFunction(index => window.__arcaneDebug?.gameState.value?.tutorial?.index === index, c.index);
        refreshedCombat = true;
      }
      if (!checkedEmpty && c.title === '读懂一张伙伴牌') {
        await page.waitForTimeout(80);
        assert(await page.locator('.fixed-inspector .keyword-notes').count() === 0, 'blank apprentice gained help terms');
        checkedEmpty = true;
      }
      if (!checkedHelp && c.title === '用奥术学习防御法术') {
        await page.waitForTimeout(100);
        const notes = page.locator('.fixed-inspector .keyword-notes');
        assert((await notes.innerText()).includes('只能用于防御'), 'missing hover help');
        const style = await notes.evaluate(el => ({margin: parseFloat(getComputedStyle(el).marginTop), border: parseFloat(getComputedStyle(el).borderLeftWidth)}));
        assert(style.margin >= 16 && style.border >= 2, 'help lacks visual separation');
        const ids = await notes.locator('dt').allTextContents();
        await page.locator('.fixed-inspector button').filter({hasText:'原卡'}).click();
        assert(JSON.stringify(ids) === JSON.stringify(await page.locator('.card-detail-modal .keyword-notes dt').allTextContents()), 'hover/detail terms differ');
        await page.screenshot({path:path.join(out,'keyword-detail.png'),fullPage:true});
        await page.locator('.card-detail-modal .close-btn').click();
        checkedHelp = true;
      }
      const a = c.allowed, d = a.data || {};
      if (a.action === 'tutorial_continue') await page.locator('[data-tutorial-continue]').click();
      else if (a.action === 'mulligan') await page.locator('[data-tutorial-keep]').click();
      else if (a.action === 'consume') await menu(d.instance_id,'消耗');
      else if (a.action === 'summon') {
        await menu(d.instance_id,'选择召唤位置','.hand-cards');
        await page.locator(`.my-board .unit-cell[data-col="${d.col}"][data-row="${d.row}"]`).click();
        await reviewAndPay(d);
      } else if (a.action === 'equip') {
        await menu(d.instance_id,'使用 / 装备','.hand-cards');
        await reviewAndPay(d);
      } else if (a.action === 'learn_skill') {
        await page.locator(`.sp-card[data-instance-id="${d.instance_id}"]`).getByRole('button',{name:'学习',exact:true}).click();
        await reviewAndPay(d);
      } else if (a.action === 'cast_spell') {
        await menu(d.instance_id,'施放');
        await page.locator(`.opponent-board .unit-cell[data-col="${d.target_col}"][data-row="${d.target_row}"]`).click();
        await reviewAndPay(d);
      } else if (a.action === 'defend') {
        await page.setViewportSize({width:1280,height:720});
        const coach = page.locator('.tutorial-defense-coach');
        const focusCard = id => page.locator(`.my-board .tutorial-focus[data-instance-id="${id}"]`).first();
        assert((await coach.innerText()).includes('① 防御'), 'missing main-defense guidance');
        assert(await focusCard(d.skill_ids[0]).count(), 'main defense card not highlighted');
        await page.getByRole('button',{name:'选择烈焰护体',exact:true}).click();
        await coach.getByRole('heading', {name:'② 强化：本次跳过',exact:true}).waitFor();
        assert((await coach.innerText()).includes('② 强化'), 'missing skip-boost guidance');
        await page.screenshot({path:path.join(out,'defense-boost-guide.png')});
        await page.getByRole('button',{name:'跳过强化 → 透支',exact:true}).click();
        assert((await coach.innerText()).includes('③ 透支'), 'missing tap-payment guidance');
        assert(await page.getByRole('button',{name:'核对付款 →',exact:true}).isDisabled(), 'payment can advance before choosing tap source');
        assert(await focusCard(d.overexert_ids[0]).count(), 'tap source not highlighted');
        await page.screenshot({path:path.join(out,'defense-overexert-guide.png')});
        await page.getByRole('button',{name:'透支火焰面甲',exact:true}).click();
        assert((await coach.innerText()).includes('面甲已选中'), 'tap selection not acknowledged');
        await page.getByRole('button',{name:'核对付款 →',exact:true}).click();
        assert((await coach.innerText()).includes('④ 付款'), 'missing final payment guidance');
        await page.locator('[data-defense-action="submit"]').click();
        await page.setViewportSize({width:1920,height:1080});
      } else if (a.action === 'no_defend') await page.locator('[data-defense-action="no_defend"]').click();
      else if (a.action === 'end_turn') await page.locator('.btn-end-turn').click();
      else if (a.action === 'resolve_action') {
        await page.locator(`.pending-action-panel [data-instance-id="${d.selected[0]}"]`).click();
        await page.locator('.pending-action-panel').getByRole('button',{name:/确认/}).last().click();
        await reviewAndPay(d);
      } else throw Error('Unknown tutorial action: '+a.action);
      await waitStep(c.index);
    }
    const final = await state();
    assert(final.tutorial.completed && final.phase === 'game_over' && final.winner === 0, 'lesson did not finish with a real win');
    assert(final.you.hero.current_life === 5 && final.opponent.hero.current_life === 0, 'wrong final life');
    assert(refreshedOpening && refreshedCombat && checkedHelp && checkedEmpty, 'missing reload/help checks');
    assert(errors.length === 0 && wireErrors.length === 0, JSON.stringify({errors,wireErrors}));
    await page.screenshot({path:path.join(out,'completed.png'),fullPage:true});
    fs.writeFileSync(path.join(out,'result.json'),JSON.stringify({ok:true,room:final.game_id,reached,errors,wireErrors,lessonURL},null,2));
    console.log('PASS: real tutorial victory, both reload points, matching separated keyword help, no page/action errors.');
  } catch (error) {
    await page.screenshot({path:path.join(out,'failure.png'),fullPage:true});
    fs.writeFileSync(path.join(out,'failure.json'),JSON.stringify({state:await state(),errors,wireErrors,reached},null,2));
    throw error;
  } finally { await browser.close(); }
})().catch(error => {console.error(error);process.exitCode=1;});
