"""Actual three-group defense with editable payment and one guarded submission."""
from frontend_card_visibility import *

async def run():
 async with async_playwright() as pw:
  browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
  room,pages,errors=await start(browser); actor=(await state(pages[0]))['current_turn']; a=pages[actor]; b=pages[1-actor]
  await add(room,actor,'3321002','skill_slot')
  await add(room,1-actor,'3321002','skill_slot')
  await add(room,1-actor,'2121011','hand')
  for owner in range(2):
   api('/api/test-room/elements',dict(room_id=room,player=owner,mode='set',elements={e:30 for e in ['火','水','气','地','光','暗','无']}))
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.elements["气"]===30')
  await a.locator('.my-board .slot').filter(has=a.locator('img[alt="雷击"]')).get_by_role('button',name='施放',exact=True).click()
  await a.locator('.opponent-board .unit-img').first.click()
  assert await a.locator('.context-menu').count()==0, 'target click must not open menu'
  await payment(a); await commit(a)
  await b.wait_for_function('window.__arcaneDebug.showDefenseWindow.value')
  before=(await state(b))['revision']; await b.evaluate('window.__sent=[]')
  for width,height in [(1440,1000),(1366,768),(1024,768)]:
   await b.set_viewport_size(dict(width=width,height=height))
   await b.locator('[data-defense-action="overexert"]').first.click()
   await b.locator('[data-defense-action="skill"]').filter(has_text='雷击').click()
   boost=b.locator('[data-defense-action="boost"]').filter(has_text='火流星卷轴')
   if not await boost.count(): boost=b.locator('[data-defense-action="boost"]:not(.disabled)').first
   await boost.click()
   for name in ['defenseSelected','defenseBoostSelected','defenseOverexertSelected']:
    assert await b.evaluate('n=>window.__arcaneDebug[n].value.length',name)==1,name
   await visible_in_viewport(b,'.defense-group > h3,.defense-actions,.fixed-inspector,#hand')
   heights=await b.locator('.defense-group-scroll').evaluate_all('(es)=>es.map(e=>e.clientHeight)')
   assert min(heights)>230,heights
   assert await b.locator('.defense-skill-img:visible').count()==0
   await b.locator('[data-defense-action="boost"].selected').hover()
   assert await b.locator('.fixed-inspector h3').inner_text()
   await b.screenshot(path=str(VISUAL/f'defense-complete-{width}.png'),animations='disabled')
   assert (await state(b))['revision']==before
   assert await b.evaluate('window.__sent.filter(m=>m.action!=="client_log").length')==0
   await b.get_by_role('button',name='重置选择',exact=True).click()
   assert await b.locator('.defense-group .selected').count()==0
  await b.locator('[data-defense-action="overexert"]').first.click()
  await b.locator('[data-defense-action="skill"]').filter(has_text='雷击').click()
  await b.locator('[data-defense-action="boost"]:not(.disabled)').first.click()
  # Remove one actually planned element, check invalid payment blocks submit, then restore it.
  row=b.locator('.defense-payment > span').filter(has_text='气').first
  await row.get_by_role('button',name='减少气').click()
  assert await b.locator('[data-defense-action="submit"]').is_disabled()
  await row.get_by_role('button',name='增加气').click()
  assert not await b.locator('[data-defense-action="submit"]').is_disabled()
  expected=await b.evaluate('({skills:[...window.__arcaneDebug.defenseSelected.value],boosts:[...window.__arcaneDebug.defenseBoostSelected.value],units:[...window.__arcaneDebug.defenseOverexertSelected.value]})')
  await b.locator('[data-defense-action="submit"]').evaluate('(e)=>{e.click();e.click()}')
  await b.wait_for_function('r=>window.__arcaneDebug.gameState.value.revision>r',arg=before)
  sent=await b.evaluate('window.__sent.filter(m=>m.action==="defend")')
  assert len(sent)==1,sent
  assert sent[0]['data']['skill_ids']==expected['skills'] and sent[0]['data']['boost_ids']==expected['boosts'] and sent[0]['data']['overexert_ids']==expected['units']
  assert sent[0]['data']['payment']
  assert (await state(b))['phase']=='main'
  assert not errors,errors
  print('PASS defense: three simultaneously visible groups; fixed reader/payment/submit; reset no send; payment correction; full spell+boost+overexert payload submitted once; target/menu separation; 1440/1366/1024')
  await browser.close()
if __name__=='__main__': asyncio.run(run())
