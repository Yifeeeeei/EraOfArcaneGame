"""Layout regression for full card faces in the real two-client match.
Run against the local server: python tests/frontend_card_visibility.py.
Also imported by match/reactions so server-issued decisions are checked.
"""
from frontend_match import *
VISUAL = Path('agent-data/formation-review'); VISUAL.mkdir(parents=True, exist_ok=True)
FACES = '.unit-img,.terrain-img,.slot-img,.hand-card-img,.sp-card > img,.pending-card-img,.defense-skill-img,.card-img-large,.spell-info-img'

async def assert_faces(page):
    await page.wait_for_function('''selector => [...document.querySelectorAll(selector)]
      .filter(i=>i.getClientRects().length).every(i=>i.complete && i.naturalWidth>0)''', arg=FACES, timeout=60000)
    problems = await page.evaluate('''selector => {
      const failures=[];
      for(const img of document.querySelectorAll(selector)) {
        if(!img.getClientRects().length) continue;
        const r=img.getBoundingClientRect(), css=getComputedStyle(img);
        if(img.closest('.fan-face')) {
          if(Math.abs(parseFloat(css.width)*img.naturalHeight/img.naturalWidth-parseFloat(css.height))>2) failures.push('fan image ratio');
          continue; // Faces intentionally overlap, but hit boxes never move.
        }
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
      const layout=document.querySelector('.game-layout');
      if(layout.scrollWidth>layout.clientWidth+1) failures.push('layout horizontal overflow');
      return failures;
    }''', FACES)
    assert not problems, problems

async def visible_in_viewport(page, selector):
    await page.wait_for_function('''selector=>[...document.querySelectorAll(selector)].filter(e=>e.getClientRects().length).every(e=>{const r=e.getBoundingClientRect();return r.top>=0 && r.bottom<=innerHeight+1 && r.left>=0 && r.right<=innerWidth+1})''',arg=selector,timeout=3000)
    failures=await page.locator(selector).evaluate_all("""els=>els.filter(e=>e.getClientRects().length).flatMap(e=>{
      const r=e.getBoundingClientRect();
      if(r.top<0 || r.bottom>innerHeight+1 || r.left<0 || r.right>innerWidth+1) return ['outside viewport '+e.className];
      const hit=document.elementFromPoint((r.left+r.right)/2,(r.top+r.bottom)/2);
      return hit && (e.contains(hit)||hit===e) ? [] : ['covered '+e.className];
    })""")
    assert not failures,failures

async def dock_check(page):
    width=page.viewport_size['width']
    if width>820:
        for end in [False,True]:
            await page.locator('#arena').evaluate('(e,end)=>{e.scrollTop=end?e.scrollHeight:0}',end)
            await visible_in_viewport(page, '#hand')
            await page.locator('.decision-scroll').evaluate_all('(els,end)=>els.forEach(e=>e.scrollTop=end?e.scrollHeight:0)',end)
            await visible_in_viewport(page, '#decision .pending-actions,#decision .defense-actions,#decision .workspace-footer')
        assert await page.evaluate('document.documentElement.scrollHeight<=innerHeight+1')
    else:
        await page.locator('.mobile-jump a[href="#hand"]').click()
        await visible_in_viewport(page, '.mobile-jump,#hand')
        await page.locator('.mobile-jump a[href="#decision"]').click()
        await visible_in_viewport(page, '.mobile-jump,#decision .pending-actions,#decision .defense-actions,#decision .workspace-footer')
        await page.locator('.decision-scroll').evaluate_all('(els)=>els.forEach(e=>e.scrollTop=e.scrollHeight)')
        await visible_in_viewport(page, '#decision .pending-actions,#decision .defense-actions,#decision .workspace-footer')

async def reset_zones(page):
    await page.evaluate("""()=>document.querySelectorAll('.game-layout,#arena,.hand-cards,.skill-pool-grid,.slots,.decision-scroll').forEach(e=>{e.scrollTop=0;e.scrollLeft=0})""")

async def decision_views(page, name):
    for width,height in [(1440,1000),(1366,768),(1024,768)]:
        await page.set_viewport_size(dict(width=width,height=height))
        await assert_faces(page)
        await dock_check(page)
        await reset_zones(page)
        if width<=820: await page.locator('.mobile-jump a[href="#decision"]').click()
        await page.screenshot(path=str(VISUAL/f'{name}-{width}.png'), animations='disabled')
        await page.locator('.decision-scroll').evaluate_all('(els)=>els.forEach(e=>e.scrollTop=e.scrollHeight)')
        await page.screenshot(path=str(VISUAL/f'{name}-{width}-bottom.png'), animations='disabled')
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
  for width,height in [(1440,1000),(1366,768),(1024,768)]:
   await a.set_viewport_size(dict(width=width,height=height))
   await assert_faces(a)
   await dock_check(a)
   # Skill-pool actions remain reachable inside their own dock; no document scrolling.
   await a.locator('.dock-tabs').get_by_role('button',name='技能池',exact=True).click()
   await a.locator('.sp-card').last.get_by_role('button',name='查看',exact=True).scroll_into_view_if_needed()
   assert await a.locator('.sp-card').last.get_by_role('button',name='查看',exact=True).is_visible()
   await a.locator('.dock-tabs button').first.click()
   await reset_zones(a)
   dimensions=await a.evaluate('''()=>Object.fromEntries(['#hand','#arena','#decision','.unit-img','.hand-card-img'].map(s=>{const r=document.querySelector(s).getBoundingClientRect();return[s,{x:r.x,y:r.y,width:r.width,height:r.height}]}))''')
   (VISUAL/f'layout-{width}.json').write_text(json.dumps(dimensions,indent=2))
   await a.screenshot(path=str(VISUAL/f'dense-{width}.png'),full_page=True,animations='disabled')
   await a.locator('#arena').evaluate('(e)=>{e.scrollTop=e.scrollHeight}')
   await a.screenshot(path=str(VISUAL/f'dense-{width}-arena-bottom.png'),animations='disabled')
   await a.locator('.hand-card').first.focus()
   await a.keyboard.press('End')
   await a.screenshot(path=str(VISUAL/f'hand-preview-{width}.png'),animations='disabled')
   await a.locator('.fan-preview').get_by_role('button',name='查看完整详情').click()
   await assert_faces(a)
   await a.locator('.card-detail-modal').screenshot(path=str(VISUAL/f'inspector-{width}.png'),animations='disabled')
   await a.locator('.card-detail-modal .close-btn').click()
   await a.locator('.hand-card').first.focus()
   await a.keyboard.press('Home')
   await a.keyboard.press('Escape')
  assert not errors,errors
  await browser.close()
  print('PASS desktop formation: pinned hand/confirmation, two dense boards, skills/equipment, last hand, inspector; 1440/1366/1024; full card ratios')
if __name__=='__main__': asyncio.run(run())
