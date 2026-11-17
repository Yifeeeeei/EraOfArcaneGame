"""Bound fallback, hand defense, disconnect invalidation and no-defense semantics."""
from frontend_match import *
async def run():
 async with async_playwright() as pw:
  browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
  room,pages,errors=await start(browser);actor=(await state(pages[0]))['current_turn'];a=pages[actor];b=pages[1-actor]
  for owner in range(2):api('/api/test-room/elements',dict(room_id=room,player=owner,mode='set',elements={e:30 for e in ['火','水','气','地','光','暗','无']}))
  await add(room,actor,'2011001','hand');await add(room,actor,'2121011','hand');await add(room,1-actor,'3321002','skill_slot')
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.hand.some(c=>c.number==="2011001")')
  await a.locator('.hand-card[data-card-number="2011001"]').click();await a.get_by_role('menuitem',name='使用 / 装备',exact=True).click();await payment(a);await commit(a)
  await a.locator('.pending-card-name').filter(has_text='雷击').first.click();await a.get_by_role('button',name='确认选择',exact=True).click();await commit(a)
  host=next(c for c in (await state(a))['you']['equipment'] if c and c['number']=='2011001');bound=host['bound_skills'][0]
  await a.get_by_role('button',name='结束回合',exact=True).click();await commit(a)
  if (await state(a)).get('pending_action'):
   while (await state(a)).get('pending_action'):
    pa=(await state(a))['pending_action'];assert pa['type']=='discard',pa
    for _ in range(pa['min_select']):await a.locator('.pending-card:not(.selected)').first.click()
    await a.get_by_role('button',name='确认选择',exact=True).click();await commit(a)
  await b.wait_for_function('n=>window.__arcaneDebug.gameState.value.current_turn===n',arg=1-actor)
  await b.locator('.my-board .slot').filter(has=b.locator('img[alt="雷击"]')).get_by_role('button',name='施放',exact=True).click();await b.locator('.opponent-board .unit-img').first.click();await payment(b);await commit(b)
  await a.wait_for_function('window.__arcaneDebug.showDefenseWindow.value');before=await state(a);dock=await a.locator('#hand').bounding_box();await a.evaluate('window.__sent=[]')
  await a.get_by_role('button',name='跳过透支',exact=True).click()
  assert await a.locator('.defense-fallback button').count()==0
  original=a.locator('.bound-skill-action[data-defense-card="'+bound['instance_id']+'"]');assert await original.get_attribute('data-defense-status')=='eligible'
  assert '绑定于' in await original.get_attribute('title')
  await original.scroll_into_view_if_needed();await a.screenshot(path='agent-data/defense-steps/bound-choice-1440.png')
  await original.focus();await a.keyboard.press('Enter')
  await a.get_by_role('button',name='跳过强化，核对付款',exact=True).click()
  await a.screenshot(path='agent-data/defense-steps/bound-payment-1440.png')
  await a.get_by_role('button',name='2 · 防御牌',exact=True).click();scroll=a.locator('.hand-card[data-card-number="2121011"]');assert await scroll.get_attribute('data-defense-status')=='eligible';await scroll.click()
  assert not await a.evaluate('window.__arcaneDebug.defenseSelected.value.length')
  assert await a.evaluate('window.__arcaneDebug.defenseScrollSelected.value.length')==1
  assert await a.locator('#hand').bounding_box()==dock
  assert await a.evaluate('window.__sent.filter(m=>m.action!=="client_log").length')==0
  await a.evaluate('window.__testSocket.close()');await a.get_by_role('button',name='重连并同步',exact=True).wait_for()
  assert await a.evaluate('window.__arcaneDebug.defenseScrollSelected.value.length')==0
  await a.get_by_role('button',name='重连并同步',exact=True).click();await a.wait_for_function('window.__arcaneDebug?.showDefenseWindow.value')
  assert await a.get_by_role('button',name='1 · 透支').get_attribute('aria-current')=='step'
  await a.locator('[data-defense-card="'+before['you']['hero']['instance_id']+'"]').click()
  await a.get_by_role('button',name='不防御',exact=True).click();await commit(a)
  after=await state(a);assert after['you']['elements']==before['you']['elements'];assert not after['you']['hero']['is_horizontal']
  await b.get_by_role('button',name='结束回合',exact=True).click();await commit(b)
  await a.locator('.my-board .unit-img').first.click();assert await a.locator('.context-menu').is_visible()
  assert not errors,errors
  print('PASS edges: original bound bar keyboard, no duplicate fallback; hand scroll original-card selection; replacement; dock stable; no middle sends; disconnect invalidates; no-defense does not consume local overexert/payment; menus restored')
  await browser.close()
if __name__=='__main__':asyncio.run(run())
