"""Layout regression for full card faces in the real two-client match.
Run against the local server: python tests/frontend_card_visibility.py.
Also imported by match/reactions so server-issued decisions are checked.
"""
from frontend_match import *
VISUAL = Path('agent-data/card-visibility'); VISUAL.mkdir(parents=True, exist_ok=True)
FACES = '.unit-img,.terrain-img,.slot-img,.hand-card-img,.sp-card > img,.pending-card-img,.defense-skill-img,.card-img-large,.spell-info-img'

async def assert_faces(page):
    await page.wait_for_function('''selector => [...document.querySelectorAll(selector)]
      .filter(i=>i.getClientRects().length).every(i=>i.complete && i.naturalWidth>0)''', arg=FACES, timeout=60000)
    problems = await page.evaluate('''selector => {
      const failures=[];
      for(const img of document.querySelectorAll(selector)) {
        if(!img.getClientRects().length) continue;
        const r=img.getBoundingClientRect(), css=getComputedStyle(img);
        const borderX=parseFloat(css.borderLeftWidth)+parseFloat(css.borderRightWidth);
        const borderY=parseFloat(css.borderTopWidth)+parseFloat(css.borderBottomWidth);
        if(css.visibility==='hidden') failures.push('hidden image '+img.className);
        if(Math.abs((r.width-borderX)*img.naturalHeight/img.naturalWidth-(r.height-borderY))>2)
          failures.push('aspect ratio '+img.className);
        const parent=img.parentElement, p=parent.getBoundingClientRect();
        if(r.top<p.top-1 || r.bottom>p.bottom+1) failures.push('clipped by card '+img.className);
        for(const other of parent.children) {
          if(other===img || !other.getClientRects().length) continue;
          const b=other.getBoundingClientRect();
          if(Math.min(r.right,b.right)-Math.max(r.left,b.left)>1 && Math.min(r.bottom,b.bottom)-Math.max(r.top,b.top)>1)
            failures.push('overlap '+img.className+' / '+other.className);
        }
      }
      if(document.documentElement.scrollWidth>innerWidth) failures.push('page horizontal overflow');
      return failures;
    }''', FACES)
    assert not problems, problems

async def decision_views(page, name):
    for width,height in [(1440,1000),(390,844)]:
        await page.set_viewport_size(dict(width=width,height=height))
        await assert_faces(page)
        panel=page.locator('#decision')
        await panel.evaluate('(e)=>{e.scrollTop=0}')
        await panel.screenshot(path=str(VISUAL/f'{name}-{width}.png'), animations='disabled')
        if width==1440:
            await panel.evaluate('(e)=>{e.scrollTop=e.scrollHeight}')
            await panel.screenshot(path=str(VISUAL/f'{name}-{width}-bottom.png'), animations='disabled')
            await panel.evaluate('(e)=>{e.scrollTop=0}')
    await page.set_viewport_size(dict(width=1440,height=1000))

async def run():
 async with async_playwright() as pw:
  browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
  room,pages,errors=await start(browser)
  for owner in range(2):
   for col in range(3):
    for row in range(3):
     if (col,row)!=(1,1): await add(room,owner,['1021014','1621109','1021002'][(col+row)%3],'unit',col=col,row=row)
   for num in ['3611101','3321002','3221001','3121001','3001001']: await add(room,owner,num,'skill_slot')
   for num in ['2021013','2311002','2321001','2011001','2011002']: await add(room,owner,num,'equipment')
   for _ in range(12): await add(room,owner,'1021001','hand')
  actor=(await state(pages[0]))['current_turn']; a=pages[actor]
  snapshot=await state(a)
  for card in [snapshot['you']['hero'],snapshot['you']['skills'][0]]:
   api('/api/test-room/card-state',dict(room_id=room,player=actor,instance_id=card['instance_id'],horizontal=True,statuses={'护盾':2},mode='replace'))
  await a.wait_for_function('window.__arcaneDebug.gameState.value.you.hero.is_horizontal')
  for width,height in [(1440,1000),(1024,768),(390,844)]:
   await a.set_viewport_size(dict(width=width,height=height))
   await assert_faces(a)
   # Every row and skill-pool action must be reachable by normal document scrolling.
   await a.locator('.sp-card').last.get_by_role('button',name='查看',exact=True).scroll_into_view_if_needed()
   assert await a.locator('.sp-card').last.get_by_role('button',name='查看',exact=True).is_visible()
   await a.evaluate('scrollTo(0,0)')
   await a.screenshot(path=str(VISUAL/f'dense-{width}.png'),full_page=True,animations='disabled')
   if width in (1440,390):
    await a.locator('.my-board').screenshot(path=str(VISUAL/f'board-{width}.png'),animations='disabled')
    await a.locator('.hand-card').last.scroll_into_view_if_needed()
    await a.locator('#hand').screenshot(path=str(VISUAL/f'hand-last-{width}.png'),animations='disabled')
    await a.locator('.hand-card').last.get_by_role('button',name='查看',exact=True).click()
    await assert_faces(a)
    await a.locator('.card-detail-modal').screenshot(path=str(VISUAL/f'inspector-{width}.png'),animations='disabled')
    last=a.locator('.card-detail-info > *').last
    await last.scroll_into_view_if_needed()
    assert await last.is_visible()
    await a.locator('.card-detail-modal .close-btn').click()
  assert not errors,errors
  await browser.close()
  print('PASS full card faces: dense two boards, 17+ hands, skills/equipment, last hand, inspector; 1440/1024/390; no sibling overlays or ratio distortion')
if __name__=='__main__': asyncio.run(run())
