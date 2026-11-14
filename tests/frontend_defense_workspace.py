"""Three-step selection on original cards with invariant dock across all steps."""
from frontend_card_visibility import *
VISUAL=Path('agent-data/defense-steps');VISUAL.mkdir(parents=True,exist_ok=True)
async def run():
 async with async_playwright() as pw:
  browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
  room,pages,errors=await start(browser); actor=(await state(pages[0]))['current_turn'];a=pages[actor];b=pages[1-actor]
  await add(room,actor,'3321002','skill_slot');await add(room,1-actor,'3321002','skill_slot')
  for _ in range(9):await add(room,1-actor,'1021001','hand')
  await add(room,1-actor,'2121011','hand')
  for owner in range(2):api('/api/test-room/elements',dict(room_id=room,player=owner,mode='set',elements={e:30 for e in ['火','水','气','地','光','暗','无']}))
  await b.wait_for_function('window.__arcaneDebug.gameState.value.you.elements["气"]===30')
  boxes={}
  for w,h in [(1440,1000),(1366,768),(1024,768)]:
   await b.set_viewport_size(dict(width=w,height=h));boxes[w]=await b.locator('#hand').bounding_box();await b.screenshot(path=str(VISUAL/f'normal-{w}.png'))
  await a.locator('.my-board .slot').filter(has=a.locator('img[alt="雷击"]')).get_by_role('button',name='施放',exact=True).click()
  await a.locator('.opponent-board .unit-img').first.click();await payment(a);await commit(a)
  await b.wait_for_function('window.__arcaneDebug.showDefenseWindow.value')
  (VISUAL/'dock-boxes.json').write_text(json.dumps(boxes,indent=2))
  before=(await state(b))['revision'];await b.evaluate('window.__sent=[]')
  you=(await state(b))['you'];hero=you['hero']['instance_id'];skill=next(c for c in you['skills'] if c and c['number']=='3321002')['instance_id'];boost=next(c for c in you['hand'] if c['number']=='2121011')['instance_id']
  async def shot(name,w):
   assert await b.locator('#hand').bounding_box()==boxes[w],(name,boxes[w],await b.locator('#hand').bounding_box())
   await visible_in_viewport(b,'.defense-actions,.defense-actions > button,#hand,.fixed-inspector')
   await b.screenshot(path=str(VISUAL/f'{name}-{w}.png'),animations='disabled')
  for w,h in [(1440,1000),(1366,768),(1024,768)]:
   await b.set_viewport_size(dict(width=w,height=h));await shot('overexert',w)
   unit=b.locator(f'[data-defense-card="{hero}"]');assert await unit.get_attribute('data-defense-status')=='eligible'
   await unit.focus();await b.keyboard.press('Enter');await b.wait_for_function('id=>window.__arcaneDebug.defenseOverexertSelected.value.includes(id)',arg=hero)
   await b.get_by_role('button',name='完成透支，下一步',exact=True).click();await shot('defense',w)
   card=b.locator(f'[data-defense-card="{skill}"]');await card.click();assert await b.get_by_role('button',name='3 · 强化').get_attribute('aria-current')=='step'
   # The hand candidate is off the current page; locate before selecting the original card.
   index=next(i for i,c in enumerate(you['hand']) if c['instance_id']==boost);page=index//7
   if not await b.locator(f'[data-defense-card="{boost}"]').count():await b.get_by_role('button',name=f'手牌第 {page+1} 页有候选 · 定位',exact=True).click()
   boostcard=b.locator(f'[data-defense-card="{boost}"]');await boostcard.focus();await b.keyboard.press('Enter');await shot('boost',w)
   assert await b.locator('.context-menu').count()==0
   await b.get_by_role('button',name='完成强化，核对付款',exact=True).click();await shot('payment',w)
   await b.get_by_role('button',name='1 · 透支',exact=True).click();await unit.click();await shot('back',w)
   assert await b.evaluate('window.__arcaneDebug.defenseOverexertSelected.value.length')==0
   await b.get_by_role('button',name='2 · 防御牌',exact=True).click();await card.click();await b.get_by_role('button',name='完成强化，核对付款',exact=True).click()
   assert (await state(b))['revision']==before
   assert await b.evaluate('window.__sent.filter(m=>m.action!=="client_log").length')==0
   await b.get_by_role('button',name='取消本地选择',exact=True).click()
  # A server revision invalidates old selections even while the defense remains open.
  await b.locator(f'[data-defense-card="{hero}"]').click()
  api('/api/test-room/elements',dict(room_id=room,player=1-actor,mode='set',elements={e:30 for e in ['火','水','气','地','光','暗','无']}))
  await b.wait_for_function('window.__arcaneDebug.defenseOverexertSelected.value.length===0')
  await b.locator(f'[data-defense-card="{hero}"]').click()
  await b.get_by_role('button',name='完成透支，下一步',exact=True).click();await b.locator(f'[data-defense-card="{skill}"]').click()
  await b.locator(f'[data-defense-card="{boost}"]').click()
  await b.get_by_role('button',name='完成强化，核对付款',exact=True).click()
  summary=await b.locator('.defense-live-summary').inner_text()
  assert '防御威力：8' in summary.replace(' ','')
  for cost in ['大气1','奥术1','火焰2']:assert cost in summary,summary
  row=b.locator('.defense-payment > span').filter(has_text='气').first
  await row.get_by_role('button',name='减少气').click();assert await b.locator('[data-defense-action="submit"]').is_disabled()
  await row.get_by_role('button',name='增加气').click();assert not await b.locator('[data-defense-action="submit"]').is_disabled()
  before=(await state(b))['revision'];await b.locator('[data-defense-action="submit"]').evaluate('(e)=>{e.click();e.click()}')
  await b.wait_for_function('r=>window.__arcaneDebug.gameState.value.revision>r',arg=before)
  sent=await b.evaluate('window.__sent.filter(m=>m.action==="defend")');assert len(sent)==1
  assert sent[0]['data']['skill_ids']==[skill] and sent[0]['data']['boost_ids']==[boost] and sent[0]['data']['overexert_ids']==[hero]
  assert sum(sent[0]['data']['payment'].values())==4
  assert (await state(b))['phase']=='main'
  assert not errors,errors
  print('PASS steps: invariant dock at 1440/1366/1024 normal/overexert/defense/boost/back/payment; original-card mouse+keyboard; hand pagination; no intermediate sends; revision invalidation; editable payment; one submission')
  await browser.close()
if __name__=='__main__':asyncio.run(run())
