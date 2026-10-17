"""Desktop interaction review: menus, fixed reader, private card backs, direct actions."""
from frontend_card_visibility import *

async def run():
 async with async_playwright() as pw:
  browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
  room,pages,errors=await start(browser); actor=(await state(pages[0]))['current_turn']; a=pages[actor]; b=pages[1-actor]
  await add(room,actor,'1021001','unit',col=0,row=0)
  await add(room,actor,'2021018','hand')
  hero=(await state(a))['you']['hero']
  for width,height in [(1440,1000),(1366,768),(1024,768)]:
   await a.set_viewport_size(dict(width=width,height=height))
   card=a.locator('.my-board .unit-img[alt="'+hero['name']+'"]')
   await card.hover()
   await a.locator('.fixed-inspector h3').filter(has_text=hero['name']).wait_for()
   await card.click()
   await a.get_by_role('menuitem',name='放大查看',exact=True).click()
   assert await a.locator('.card-detail-info h3').inner_text()==hero['name']
   await a.locator('.card-detail-modal .close-btn').click()
   # Keyboard activation, escape restoration and clamped menus at each viewport edge.
   await card.focus(); await a.keyboard.press('Enter')
   assert await a.locator('.context-menu').is_visible()
   await a.keyboard.press('End'); await a.keyboard.press('Home'); await a.keyboard.press('Escape')
   assert await a.locator('.context-menu').count()==0
   for x,y in [(1,1),(width-1,1),(1,height-1),(width-1,height-1)]:
    await card.dispatch_event('click',{'clientX':x,'clientY':y})
    await a.wait_for_timeout(50)
    assert await a.locator('.context-menu').evaluate('(e)=>{const r=e.getBoundingClientRect();return r.left>=0&&r.top>=0&&r.right<=innerWidth&&r.bottom<=innerHeight}')
    await a.keyboard.press('Escape')
   await card.click()
   await a.screenshot(path=str(VISUAL/f'menu-reader-{width}.png'),animations='disabled')
   await a.keyboard.press('Escape')
  await a.set_viewport_size(dict(width=1440,height=1000))
  # An intentional menu command consumes directly; two synthetic clicks cannot send twice.
  unit=a.locator('.my-board .unit-img[alt="巫师的学徒"]')
  await unit.click(); await a.evaluate('window.__sent=[]')
  before=(await state(a))['revision']
  await a.get_by_role('menuitem',name='消耗',exact=True).evaluate('(e)=>{e.click();e.click()}')
  await a.wait_for_function('r=>window.__arcaneDebug.gameState.value.revision>r',arg=before)
  assert await a.evaluate('window.__sent.filter(m=>m.action==="consume").length')==1
  assert await a.locator('.action-review').count()==0
  assert (await state(a))['you']['units'][0][0]['is_horizontal']
  # Canceling a menu has no gameplay effect; a hidden counter is only shown as a back to the opponent.
  await a.locator('.hand-card').first.focus(); await a.keyboard.press('End')
  trap=a.locator('.hand-card[data-card-number="2021018"]')
  await trap.click(); await a.get_by_role('menuitem',name='使用 / 装备',exact=True).click()
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.equipment.some(c=>c?.is_set_counter)')
  await b.wait_for_function('window.__arcaneDebug.gameState.value.opponent.equipment.some(c=>c?.is_hidden)')
  assert await a.locator('.my-board .set-card-back').count()==1
  hidden=b.locator('.opponent-board .set-card-back')
  assert await hidden.count()==1
  assert await hidden.get_attribute('src')=='images/card-back.svg'
  await hidden.hover(); await hidden.click()
  assert '奥术符文' not in await b.locator('.fixed-inspector').inner_text()
  assert await b.locator('.card-detail-modal,.context-menu').count()==0
  # No leaked image request/name/card number in the hidden opponent slot.
  assert '2021018' not in await hidden.locator('..').inner_html()
  await b.screenshot(path=str(VISUAL/'private-card-back.png'),animations='disabled')
  assert not errors,errors
  await browser.close()
  print('PASS usability: own hero hover/enlarge; edge-clamped keyboard menus; escape; direct consume duplicate protection; owner/opponent card backs and hidden identity; 1440/1366/1024')
if __name__=='__main__': asyncio.run(run())
