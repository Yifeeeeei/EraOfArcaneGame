"""Actual blood-sand partial failure, two players and an observer; no rule mocks."""
from frontend_match import *
from frontend_reactions import stage
async def run():
 async with async_playwright() as pw:
  browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
  room,pages,errors=await start(browser); actor=(await state(pages[0]))['current_turn']; a=pages[actor]; b=pages[1-actor]
  await add(room,actor,'3411102','skill_slot')
  await add(room,actor,'1021002','unit',col=0,row=0)
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.units[0][0]?.number==="1021002"')
  skill=next(c for c in (await state(a))['you']['skills'] if c and c['number']=='3411102')
  unit=(await state(a))['you']['units'][0][0]
  api('/api/test-room/card-state',dict(room_id=room,player=actor,instance_id=unit['instance_id'],life=2))
  await a.wait_for_timeout(200)
  await stage(a,'use_ability',{'instance_id':skill['instance_id'],'ability':'per_turn'})
  watcher=await browser.new_page(viewport={'width':1440,'height':1000},ignore_https_errors=True)
  watcher.on('pageerror',lambda e:errors.append(str(e)))
  await watcher.goto(BASE+'/game.html?'+urllib.parse.urlencode(dict(room=room,player_id='failure-observer',role='spectator')))
  await watcher.wait_for_function('window.__arcaneDebug?.gameState.value?.pending_action')
  before=(await state(a))['revision']; did=(await state(a))['pending_action']['decision_id']
  payload={'selected':[unit['instance_id']],'payments':[{'instance_id':unit['instance_id'],'amount':2,'mode':'life'}]}
  # Intentionally malformed payment bypasses the UI capacity check, but submits
  # through the real review/confirm transport to the real existing card callback.
  await a.evaluate('d=>window.__arcaneDebug.sendAction("resolve_action",d)',payload)
  await a.get_by_role('button',name='确认提交',exact=True).click()
  for page in [a,b,watcher]: await page.wait_for_function('r=>window.__arcaneDebug.gameState.value.revision===r',arg=before+1)
  await a.locator('.action-review').wait_for(state='hidden')
  assert 'cannot pay life from selected unit' in await a.locator('.decision-error').inner_text()
  assert (await state(a))['you']['units'][0][0]['current_life']==1
  assert (await state(b))['opponent']['units'][0][0]['current_life']==1
  assert (await state(a))['pending_action']['decision_id']==did
  assert await a.locator('.pending-card.selected').count()==0
  assert (await state(b))['pending_action'] is None
  public=await state(watcher)
  assert 'candidates' not in public['pending_action'] and 'type' not in public['pending_action']
  assert not public['you'].get('hand') and not public['opponent'].get('hand')
  assert await watcher.locator('.pending-card').count()==0
  # Another rejected handler changes no life but invalidates the submitted draft.
  await a.evaluate('window.__arcaneDebug.sendAction("resolve_action",{selected:["invalid"]})')
  await a.get_by_role('button',name='确认提交',exact=True).click()
  for page in [a,b,watcher]: await page.wait_for_function('r=>window.__arcaneDebug.gameState.value.revision===r',arg=before+2)
  assert (await state(a))['you']['units'][0][0]['current_life']==1
  assert 'invalid selection' in await a.locator('.decision-error').inner_text()
  await a.locator('.action-review').wait_for(state='hidden')
  # A preflight rejection preserves this local draft and does not advance state.
  await a.evaluate('window.__arcaneDebug.sendAction("unknown_action",{})')
  await a.get_by_role('button',name='确认提交',exact=True).click()
  await a.wait_for_function('document.querySelector(".decision-error")?.textContent.includes("unknown action")')
  assert await a.locator('.action-review').count()==1
  assert (await state(a))['revision']==before+2
  await a.get_by_role('button',name='返回编辑',exact=True).click()
  await a.get_by_role('button',name='不选择 / 放弃发动',exact=True).click(); await commit(a)
  assert (await state(b))['pending_action']['decision_id'] != did
  await b.get_by_role('button',name='不选择 / 放弃发动',exact=True).click(); await commit(b)
  assert (await state(a))['phase']=='main'
  assert not errors,errors
  print('PASS failed action: real partial life payment, original errors, revision/decision identity, two-player+observer filtered sync, no-change rule rejection, preflight draft preservation, fresh decision resumes')
  await browser.close()
if __name__=='__main__': asyncio.run(run())
