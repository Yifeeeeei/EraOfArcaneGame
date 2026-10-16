"""Desktop fan interaction regression. Real two-client room; no desktop automation.
The final empty-hand check is explicitly an isolated presentation fixture, restored
without sending it to the server. All preceding checks use server-issued cards.
"""
from frontend_card_visibility import *

async def run():
 async with async_playwright() as pw:
  browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
  room,pages,errors=await start(browser)
  actor=(await state(pages[0]))['current_turn']; a=pages[actor]
  for _ in range(26): await add(room,actor,'1021001','hand')
  await add(room,actor,'3611101','hand') # Skill is inspectable but cannot be played from hand.
  await a.wait_for_function('window.__arcaneDebug.myHand.value.length>30')
  ids=[c['instance_id'] for c in (await state(a))['you']['hand']]
  await a.evaluate('window.__sent=[]')
  for width,height in [(1440,1000),(1366,768),(1024,768)]:
   await a.set_viewport_size(dict(width=width,height=height))
   await assert_faces(a)
   first=a.locator('.hand-card').first
   await first.focus(); await a.keyboard.press('Home'); await a.keyboard.press('Escape')
   # Repeated sweeps across exposed hit strips; raised artwork cannot catch the pointer.
   rects=await a.locator('.hand-card').evaluate_all('(es)=>es.map(e=>{const r=e.getBoundingClientRect();return {id:e.dataset.instanceId,x:r.x+8,y:r.y+70,rect:[r.x,r.y,r.width,r.height]}})')
   for _ in range(3):
    for r in rects:
     await a.mouse.move(r['x'],r['y'])
     await a.wait_for_timeout(35)
     assert await a.locator('.fan-preview img').get_attribute('alt')
     hit=await a.evaluate('p=>document.elementFromPoint(p.x,p.y)?.closest(".hand-card")?.dataset.instanceId',r)
     assert hit==r['id'],(hit,r)
   after=await a.locator('.hand-card').evaluate_all('(es)=>es.map(e=>{const r=e.getBoundingClientRect();return [r.x,r.y,r.width,r.height]})')
   assert after==[r['rect'] for r in rects], (after,[r['rect'] for r in rects])
   assert await a.evaluate('''()=>{const b=document.querySelector('#hand').getBoundingClientRect();return [...document.querySelectorAll('.fan-face')].every(e=>{const r=e.getBoundingClientRect();return r.left>=b.left && r.right<=b.right})}'''), 'outer fan face clipped'
   before=await a.locator('.fan-preview img').get_attribute('src')
   await a.locator('.fan-preview img').hover()
   assert await a.locator('.fan-preview img').get_attribute('src')==before
   assert await a.evaluate('''()=>{const a=document.querySelector('#arena').getBoundingClientRect(),b=document.querySelector('.fan-preview').getBoundingClientRect();return b.left>=a.right}''')
   assert await a.locator('.hand-card[aria-pressed=true]').count()==0
   # Arrow traversal reaches every instance, including identical printed cards, across pages.
   await first.focus(); await a.keyboard.press('Home')
   reached=[]
   for i in range(len(ids)):
    reached.append(await a.evaluate('document.activeElement.dataset.instanceId'))
    if i+1<len(ids): await a.keyboard.press('ArrowRight')
   assert reached==ids
   await a.keyboard.press('Enter')
   assert await a.locator('.hand-card[aria-pressed=true]').count()==0 # unplayable skill
   await a.keyboard.press('Home'); await a.keyboard.press('End')
   assert await a.evaluate('document.activeElement.dataset.instanceId')==ids[-1]
   await a.keyboard.press('Home')
   await a.get_by_role('button',name='平铺',exact=True).click()
   await a.locator('.hand-card').last.scroll_into_view_if_needed()
   await visible_in_viewport(a,'.hand-card:last-child')
   await a.get_by_role('button',name='扇形',exact=True).click()
   await a.get_by_role('slider',name='手牌间距').fill('22')
   await a.get_by_role('slider',name='手牌间距').fill('30')
   await a.locator('.hand-card').first.focus()
   await a.screenshot(path=str(VISUAL/f'fan-interaction-{width}.png'),animations='disabled')
   assert await a.evaluate('window.__sent.filter(m=>m.action!=="client_log").length')==0
  # With resources available, keyboard chooses a companion, never submits it.
  api('/api/test-room/elements',dict(room_id=room,player=actor,mode='set',elements={e:30 for e in ['火','水','气','地','光','暗','无']}))
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.elements["无"]===30')
  await a.locator('.hand-card').first.focus(); await a.keyboard.press('Home'); await a.keyboard.press('Enter')
  assert await a.locator('.hand-card[aria-pressed=true]').count()==1
  assert await a.locator('.fan-preview').count()==0 # does not displace active targeting
  await a.get_by_role('button',name='取消本地选择',exact=True).click()
  assert await a.evaluate('window.__sent.filter(m=>m.action!=="client_log").length')==0
  # The former per-card scroll boost remains available in the shared hand control.
  await add(room,actor,'3321002','skill_slot')
  await add(room,actor,'2121003','hand')
  await a.wait_for_function('window.__arcaneDebug.myHand.value.some(c=>c.number==="2121003")')
  await a.locator('.my-board .slot').filter(has=a.locator('img[alt="雷击"]')).get_by_role('button',name='施放',exact=True).click()
  await a.locator('.hand-card').first.focus(); await a.keyboard.press('End')
  await a.locator('.fan-boost').click()
  assert '取消强化' in await a.locator('.fan-boost').inner_text()
  assert await a.locator('.hand-card.boost-selected').count()==1
  await a.locator('.fan-boost').click()
  assert await a.locator('.hand-card.boost-selected').count()==0
  await a.get_by_role('button',name='取消本地选择',exact=True).click()
  assert await a.evaluate('window.__sent.filter(m=>m.action!=="client_log").length')==0
  # Presentation boundary only: zero cards after being on the last page clamps safely.
  await a.locator('.hand-card').first.focus(); await a.keyboard.press('End')
  await a.evaluate('window.__savedHand=window.__arcaneDebug.gameState.value.you.hand; window.__arcaneDebug.gameState.value.you.hand=[]')
  await a.get_by_text('暂无手牌',exact=True).wait_for()
  assert await a.locator('.hand-card').count()==0
  assert await a.get_by_role('button',name='下一页手牌').is_disabled()
  await a.evaluate('window.__arcaneDebug.gameState.value.you.hand=window.__savedHand; delete window.__savedHand')
  assert not errors,errors
  await browser.close()
  print('PASS fan: 31+ cards, duplicate instances, stable mouse hit strips, persistent preview, all-card keyboard traversal, disabled/selected cards, no hover/key submission, paging/spread/flat layout, scroll boost toggle, isolated empty-hand fixture; 1440/1366/1024')
if __name__=='__main__': asyncio.run(run())
