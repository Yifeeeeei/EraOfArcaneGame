"""Counter -> nested choice -> second counter decline -> reaction -> defense."""
from frontend_match import *
async def stage(page,action,data):
 await page.evaluate('([action,data])=>window.__arcaneDebug.sendAction(action,data)',[action,data]); await commit(page)
async def choose(page,index=0):
 await page.locator('.pending-card').nth(index).locator('.pending-card-name').click()
 await page.get_by_role('button',name='确认选择',exact=True).click()
 await payment(page); await commit(page)
async def run():
 async with async_playwright() as pw:
  browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
  room,pages,errors=await start(browser); actor=(await state(pages[0]))['current_turn']; a=pages[actor]; b=pages[1-actor]
  await add(room,actor,'3321001','skill_slot')
  await add(room,1-actor,'3321002','skill_slot')
  await add(room,1-actor,'3221008','skill_slot')
  for _ in range(2): await add(room,1-actor,'2021018','hand')
  await add(room,1-actor,'1021001','unit',col=0,row=0)
  await add(room,1-actor,'1021002','unit',col=2,row=2)
  # Turn transfer and actual setting of both face-down counters through the review path.
  await a.get_by_role('button',name='结束回合',exact=True).click(); await commit(a)
  for _ in range(2):
   card=next(c for c in (await state(b))['you']['hand'] if c['number']=='2021018')
   await stage(b,'use_item',{'instance_id':card['instance_id']})
  await b.get_by_role('button',name='结束回合',exact=True).click(); await commit(b)
  for i in range(2): api('/api/test-room/elements',dict(room_id=room,player=i,mode='set',elements={e:30 for e in ['火','水','气','地','光','暗','无']}))
  api('/api/test-room/elements',dict(room_id=room,player=1-actor,mode='set',elements={**{e:0 for e in ['火','水','气','地','光','暗','无']},'水':2,'无':1}))
  await b.wait_for_function('window.__arcaneDebug.gameState.value.you.elements["无"]===1')
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.elements["气"]===30')
  await a.locator('.my-board .slot').filter(has=a.locator('img[alt="闪电链"]')).get_by_role('button',name='施放',exact=True).click()
  await a.locator('.opponent-board .unit-cell').filter(has_text='巫师的学徒').click()
  await a.locator('.opponent-board .unit-cell').filter(has_text='学院导师').click()
  await payment(a); await commit(a)
  await b.wait_for_function('window.__arcaneDebug.gameState.value.pending_action?.type === "counter_trigger"')
  assert (await state(a)).get('pending_action') is None
  # Unknown counter names are absent from the caster's DOM until revealed.
  assert '奥术符文' not in await a.locator('.decision-workspace').inner_text()
  await b.screenshot(path=str(OUT/'counter-desktop.png'),full_page=True)
  await choose(b)
  assert (await state(b))['pending_action']['type']=='arcane_rune_skill'
  await choose(b)
  if (await state(b)).get('pending_action'):
   await b.get_by_role('button',name='不选择 / 放弃发动',exact=True).click(); await commit(b)
  await b.wait_for_function('window.__arcaneDebug.showDefenseWindow.value')
  before=(await state(b))['revision']
  await b.locator('[data-defense-action="reaction"]').first.click()
  assert (await state(b))['revision']==before,'reaction click must not send immediately'
  await payment(b); await commit(b)
  if (await state(b)).get('pending_action'): await choose(b)
  assert (await state(b))['phase']=='defense_window'
  if (await state(b))['phase']=='defense_window':
   await b.locator('[data-defense-action="overexert"]').filter(has_text='掌门').click()
   await b.locator('[data-defense-action="skill"]').first.click()
   await b.get_by_role('button',name='确认防御',exact=True).click()
   await payment(b); await commit(b)
  await a.wait_for_function('window.__arcaneDebug.gameState.value.phase === "main"')
  assert (await state(b))['you']['hero']['is_horizontal']
  assert (await state(b))['you']['elements'].get('气',0)==0, 'excess overexert must not enter pool'
  assert not errors,errors
  print('PASS two clients: extra target -> counter/payment -> nested choice -> decline next counter -> explicit reaction -> paid defense -> original turn')
  await browser.close()
if __name__=='__main__': asyncio.run(run())
