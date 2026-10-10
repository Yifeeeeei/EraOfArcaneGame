"""Review original entry points using two players plus a read-only spectator."""
from frontend_match import *
async def run():
 async with async_playwright() as pw:
  browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
  room,pages,errors=await start(browser); actor=(await state(pages[0]))['current_turn']; a=pages[actor]
  revision=(await state(a))['revision']
  await a.evaluate("""r => { window.__staleAck=null; window.__testSocket.addEventListener('message', e=>{const m=JSON.parse(e.data);if(m.request_id==='quality-stale') window.__staleAck=m;}); window.__testSocket.send(JSON.stringify({action:'end_turn',data:{},request_id:'quality-stale',expected_revision:r-1})); }""", revision)
  await a.wait_for_function('window.__staleAck')
  assert not await a.evaluate('window.__staleAck.ok')
  assert (await state(a))['revision']==revision
  api('/api/test-room/elements',dict(room_id=room,player=actor,mode='set',elements={e:30 for e in ['火','水','气','地','光','暗','无']}))
  await add(room,actor,'2011001','equipment')
  await add(room,actor,'2011001','hand')
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.hand.some(c=>c.number==="2011001")')
  staff=next(c for c in (await state(a))['you']['hand'] if c['number']=='2011001')
  hand=a.locator('.hand-card').filter(has_text=staff['name'])
  await hand.click(); await hand.click()
  await a.get_by_role('heading',name='选择要替换的装备').wait_for()
  await a.locator('#decision .pending-card-name').first.click()
  await a.get_by_role('button',name='确认替换',exact=True).click()
  await payment(a)
  # A card inspection opened during a draft must remain closable.
  await hand.get_by_role('button',name='查看',exact=True).click()
  await a.locator('.card-detail-modal .close-btn').click()
  assert await a.locator('.card-detail-modal').count()==0
  await commit(a)
  await a.locator('.pending-card-name').filter(has_text='雷击').first.click()
  await a.get_by_role('button',name='确认选择',exact=True).click(); await commit(a)
  host=next(c for c in (await state(a))['you']['equipment'] if c and c['instance_id']==staff['instance_id'])
  bound=host['bound_skills'][0]
  await a.wait_for_timeout(200)
  await a.locator('.my-board .slot').filter(has=a.locator('img[alt="'+staff['name']+'"]')).locator('img').click()
  await a.locator('.bound-skill-chip').filter(has_text=bound['name']).first.click()
  assert await a.locator('.card-detail-info h3').inner_text()==bound['name']
  await a.locator('.card-detail-modal .close-btn').click()
  await add(room,actor,'3311102','skill_slot')
  for _ in range(8): await add(room,actor,'1021001','hand')
  await a.get_by_role('button',name='结束回合',exact=True).click(); await commit(a)
  assert (await state(a))['pending_action']['type']=='discard'
  watcher=await browser.new_page(viewport={'width':1440,'height':1000},ignore_https_errors=True)
  watcher.on('pageerror',lambda e:errors.append(str(e)))
  await watcher.goto(BASE+'/game.html?'+urllib.parse.urlencode(dict(room=room,player_id='quality-observer',role='spectator')))
  await watcher.wait_for_function('window.__arcaneDebug?.gameState.value?.phase==="waiting_action"')
  public=await state(watcher)
  assert not public['you'].get('hand') and not public['opponent'].get('hand')
  assert 'candidates' not in public['pending_action']
  assert await watcher.locator('.pending-card').count()==0
  assert await watcher.get_by_role('button',name='投降',exact=True).count()==0
  assert await watcher.locator('.turn-badge').count()==1
  count=(await state(a))['pending_action']['min_select']
  for i in range(count): await a.locator('.pending-card').nth(i).click()
  await a.get_by_role('button',name='确认选择',exact=True).click(); await commit(a)
  stored=0
  for _ in range(20):
   pending=(await state(a)).get('pending_action')
   if not pending: break
   assert pending['type']=='starfall_silverleaf_store_discard',pending['type']
   await a.locator('.pending-card-name').first.click()
   await a.get_by_role('button',name='确认选择',exact=True).click(); await commit(a); stored+=1
  assert stored>0
  star=next(c for c in (await state(a))['you']['skills'] if c and c['number']=='3311102')
  assert len(star['under_cards'])==stored
  await a.locator('.my-board .slot').filter(has=a.locator('img[alt="'+star['name']+'"]')).click()
  await a.locator('.under-card-list').wait_for(state='visible')
  await a.screenshot(path=str(OUT/'quality-under-cards.png'),full_page=True,animations='disabled')
  await a.locator('.under-card-list button').first.click()
  assert await a.locator('.card-detail-info h3').inner_text()==star['under_cards'][0]['name']
  await a.locator('.card-detail-modal .close-btn').click()
  await a.get_by_role('button',name='投降',exact=True).click()
  assert (await state(a))['phase']!='game_over'
  await commit(a)
  for page in [*pages,watcher]: await page.wait_for_function('window.__arcaneDebug.gameState.value.phase==="game_over"')
  assert not errors,errors
  print('PASS quality: equipment replacement, binding inspection, inspection close during draft, end-turn discard/under-card chain, under-card inspection, spectator privacy/turn badge, explicit surrender')
  await browser.close()
if __name__=='__main__': asyncio.run(run())
