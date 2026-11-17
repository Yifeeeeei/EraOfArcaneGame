"""Real-client regression for leaving/resuming, Raven actions and incoming spell details.
Run with the server running; ARCANE_CHROMIUM may select a local Chrome executable.
"""
from frontend_match import *

async def run():
 async with async_playwright() as pw:
  browser = await pw.chromium.launch(executable_path=os.environ.get('ARCANE_CHROMIUM', '/usr/bin/chromium'), headless=True, args=['--no-sandbox'])
  room, pages, errors = await start(browser)
  actor = (await state(pages[0]))['current_turn']; a = pages[actor]; b = pages[1-actor]
  original_url = a.url; original_slot = await a.evaluate('window.__arcaneDebug.mySlot.value')
  original_hero = (await state(a))['you']['hero']['instance_id']
  revision = (await state(a))['revision']
  await a.locator('.match-status .brand').click()
  dialog = a.get_by_role('dialog', name='暂离对局', exact=True)
  await dialog.wait_for(); assert a.url == original_url
  await dialog.get_by_role('button', name='继续对局', exact=True).click()
  assert a.url == original_url and (await state(a))['revision'] == revision
  await a.locator('.match-status .brand').click()
  await a.get_by_role('button', name='暂离到大厅', exact=True).click()
  await a.wait_for_url(BASE + '/')
  await a.get_by_role('button', name='返回对局', exact=True).click()
  await a.wait_for_function('window.__arcaneDebug?.gameState.value?.phase === "main"')
  assert a.url == original_url
  assert await a.evaluate('window.__arcaneDebug.mySlot.value') == original_slot
  assert (await state(a))['you']['hero']['instance_id'] == original_hero
  assert (await state(a))['revision'] == revision
  assert await a.evaluate('window.__sent.filter(m=>m.action!=="client_log").length') == 0
  print('PASS leave confirmation, cancel and lobby resume with original identity/no replay', flush=True)

  await add(room, actor, '1321001', 'unit', col=0, row=0)
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.units[0][0]?.number === "1321001"')
  raven = (await state(a))['you']['units'][0][0]
  api('/api/test-room/card-state', dict(room_id=room, player=actor, instance_id=raven['instance_id'], horizontal=False))
  await a.wait_for_function('!window.__arcaneDebug.gameState.value.you.units[0][0].is_horizontal')
  await a.locator('.my-board .unit-img[alt="渡鸦信使"]').click()
  await a.get_by_role('menuitem', name='消耗', exact=True).wait_for()
  await a.get_by_role('menuitem', name='消耗抽牌', exact=True).wait_for()
  hand_count = len((await state(a))['you']['hand']); revision = (await state(a))['revision']
  await a.get_by_role('menuitem', name='消耗抽牌', exact=True).click()
  await a.wait_for_function('r=>window.__arcaneDebug.gameState.value.revision>r', arg=revision)
  assert len((await state(a))['you']['hand']) == hand_count + 1
  assert (await state(a))['you']['units'][0][0]['is_horizontal']
  assert await a.evaluate('window.__sent.filter(m=>m.action==="use_ability").length') == 1
  api('/api/test-room/card-state', dict(room_id=room, player=actor, instance_id=raven['instance_id'], horizontal=False))
  await a.wait_for_function('!window.__arcaneDebug.gameState.value.you.units[0][0].is_horizontal')
  await a.locator('.my-board .unit-img[alt="渡鸦信使"]').click()
  revision = (await state(a))['revision']; hand_count = len((await state(a))['you']['hand'])
  await a.get_by_role('menuitem', name='消耗', exact=True).click()
  await a.wait_for_function('r=>window.__arcaneDebug.gameState.value.revision>r', arg=revision)
  assert len((await state(a))['you']['hand']) == hand_count
  assert await a.evaluate('window.__sent.filter(m=>m.action==="consume").length') == 1
  print('PASS distinct Raven labels and correct ordinary consume / ability actions', flush=True)

  for owner in range(2):
   api('/api/test-room/elements', dict(room_id=room, player=owner, mode='set', elements={e:30 for e in ['火','水','气','地','光','暗','无']}))
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.elements["气"]===30')
  await a.locator('.dock-tabs').get_by_role('button', name='技能池', exact=True).click()
  await a.locator('.sp-card').filter(has_text='雷击').get_by_role('button', name='学习', exact=True).click()
  await payment(a); await commit(a)
  skill = next(c for c in (await state(a))['you']['skills'] if c and c['name']=='雷击')
  api('/api/test-room/card-state', dict(room_id=room, player=actor, instance_id=skill['instance_id'], horizontal=False, statuses={'冷却':0}))
  await add(room, 1-actor, '1021011', 'unit', col=0, row=0)
  await a.locator('.my-board .slot').filter(has=a.locator('img[alt="雷击"]')).click()
  await a.get_by_role('menuitem', name='施放', exact=True).click()
  await a.locator('.opponent-board .unit-img[alt="屠魔者杀手"]').click()
  await payment(a); await commit(a)
  await b.wait_for_function('window.__arcaneDebug.showDefenseWindow.value')
  revision = (await state(b))['revision']
  await b.evaluate('window.__sent=[]')
  for width, height in [(1440,1000), (1024,768)]:
   await b.set_viewport_size(dict(width=width, height=height))
   incoming = b.get_by_role('button', name='查看详情').filter(has_text='来袭：')
   await incoming.focus(); await b.keyboard.press('Enter')
   await b.get_by_role('dialog', name='雷击', exact=True).wait_for()
   assert '使用费用' in await b.locator('.card-detail-info').inner_text()
   await b.locator('.card-img-large').evaluate('(img)=>img.decode()')
   assert await b.locator('.card-img-large').evaluate('(img)=>{const r=img.getBoundingClientRect();return img.naturalHeight>img.naturalWidth && r.height>r.width && r.height<=innerHeight*.81}')
   await b.screenshot(path=str(OUT/f'reported-ux-spell-detail-{width}.png'))
   await b.locator('.card-detail-modal .close-btn').click()
   assert (await state(b))['revision'] == revision
   assert await b.evaluate('window.__sent.filter(m=>m.action!=="client_log").length') == 0
  await b.get_by_role('button', name='不防御', exact=True).click(); await commit(b)
  await a.wait_for_function('window.__arcaneDebug.gameState.value.phase === "main"')
  print('PASS incoming spell details via keyboard at 1440/1024; no selection/submission side effects; defense resolves', flush=True)

  for _ in range(6): await add(room, actor, '1021001', 'hand')
  await a.get_by_role('button', name='结束回合', exact=True).click(); await commit(a)
  await a.wait_for_function('window.__arcaneDebug.gameState.value.pending_action?.type === "discard"')
  pending = (await state(a))['pending_action']
  for i in range(pending['min_select']): await a.locator('.pending-card').nth(i).click()
  await a.get_by_role('button', name='确认选择', exact=True).click(); await commit(a)
  await a.wait_for_function('(actor)=>window.__arcaneDebug.gameState.value.current_turn!==actor', arg=actor)
  await a.get_by_role('button', name='求和', exact=True).click(); await commit(a)
  await b.get_by_role('button', name='接受和局', exact=True).click(); await commit(b)
  for page in pages:
   await page.wait_for_function('window.__arcaneDebug.gameState.value.phase === "game_over"')
   assert await page.evaluate('ArcaneMatchSessions.list().length') == 0
  assert not errors, errors
  print('PASS discard handoff, both clients game-over and saved-seat cleanup', flush=True)
  await browser.close()

if __name__ == '__main__': asyncio.run(run())
