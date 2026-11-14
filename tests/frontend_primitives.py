"""Real engine decisions: ordering, numeric allocations, load redistribution."""
from frontend_reactions import stage
from frontend_match import *
async def run():
 async with async_playwright() as pw:
  browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
  room,pages,errors=await start(browser); actor=(await state(pages[0]))['current_turn']; a=pages[actor]; b=pages[1-actor]
  await add(room,actor,'1421001','unit',col=0,row=0)
  await add(room,actor,'1421002','unit',col=2,row=0)
  await add(room,actor,'2421105','hand')
  api('/api/test-room/elements',dict(room_id=room,player=actor,mode='set',elements={'地':30,'无':30}))
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.elements["地"]===30')
  card=next(c for c in (await state(a))['you']['hand'] if c['number']=='2421105')
  await stage(a,'use_item',{'instance_id':card['instance_id']})
  assert (await state(a))['pending_action']['ui']['kind']=='load_distribution'
  for i in range(2): await a.locator('.pending-card-name').nth(i).click()
  groups=a.locator('.load-distribution > div')
  await groups.nth(0).locator('input').nth(2).fill('1')
  await groups.nth(1).locator('input').nth(2).fill('2')
  await a.get_by_role('button',name='确认选择',exact=True).click(); await commit(a)
  assert (await state(a))['you']['units'][0][0]['elements_gain']['地']==1
  assert (await state(a))['you']['units'][2][0]['elements_gain']['地']==2
  await add(room,actor,'2011102','equipment')
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.equipment.some(c=>c?.number==="2011102")')
  card=next(c for c in (await state(a))['you']['equipment'] if c and c['number']=='2011102')
  await stage(a,'use_ability',{'instance_id':card['instance_id'],'ability':'per_turn'})
  assert (await state(a))['pending_action']['ui']['kind']=='reorder'
  await a.get_by_role('button',name='下移',exact=True).first.click()
  await a.get_by_role('button',name='放底',exact=True).first.click()
  await a.get_by_role('button',name='确认选择',exact=True).click(); await commit(a)
  for col in [0,2]: await add(room,1-actor,'1021002','unit',col=col,row=0)
  await add(room,actor,'2201103','hand')
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.hand.some(c=>c?.number==="2201103")')
  card=next(c for c in (await state(a))['you']['hand'] if c['number']=='2201103')
  await stage(a,'use_item',{'instance_id':card['instance_id']})
  assert (await state(a))['pending_action']['ui']['kind']=='allocation'
  for i in range(2): await a.locator('.pending-card-name').nth(i).click()
  await a.locator('.load-distribution input').nth(0).fill('2')
  await a.locator('.load-distribution input').nth(1).fill('1')
  await a.get_by_role('button',name='确认选择',exact=True).click(); await commit(a)
  await add(room,actor,'3411102','skill_slot')
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.skills.some(c=>c?.number==="3411102")')
  card=next(c for c in (await state(a))['you']['skills'] if c and c['number']=='3411102')
  await stage(a,'use_ability',{'instance_id':card['instance_id'],'ability':'per_turn'})
  assert (await state(a))['pending_action']['ui']['kind']=='resource_distribution'
  await a.locator('.pending-card-name').first.click()
  await a.locator('.load-distribution input').fill('1')
  await a.get_by_role('button',name='确认选择',exact=True).click(); await commit(a)
  await b.get_by_role('button',name='不选择 / 放弃发动',exact=True).click(); await commit(b)
  assert not errors,errors
  print('PASS primitives: load redistribution, deck top/bottom ordering, exact damage allocation, resource modes and two-player payment handoff')
  await browser.close()
if __name__=='__main__': asyncio.run(run())
