"""Two real Vue clients against a running local server. No window automation.
Run: python tests/frontend_match.py (Python Playwright + /usr/bin/chromium).
Fixtures use test-room APIs only. Gameplay goes through the real page/WebSocket.
"""
import asyncio, json, os, urllib.request, urllib.parse
from pathlib import Path
from playwright.async_api import async_playwright
BASE=os.environ.get('ARCANE_URL','http://localhost:9090')
OUT=Path('agent-data/frontend-review'); OUT.mkdir(parents=True,exist_ok=True)
DECK='4311003 // '+' '.join(n for n in ['1021001','1021002','1021004','1021005','1021006','1021007','1021008','1021009','1021010','1021011','1021012','1021013','1021014','1021015','1021016'] for _ in range(2))+' // 3321002 3001001 3001002 3021001 3021002 3021003 3021004 3021005 3021006 3021007'
def api(path,data=None):
    req=urllib.request.Request(BASE+path,data=None if data is None else json.dumps(data).encode(),headers={'Content-Type':'application/json'})
    return json.load(urllib.request.urlopen(req))
async def state(page): return await page.evaluate('window.__arcaneDebug.gameState.value')
async def commit(page):
    revision=(await state(page))['revision']
    await page.get_by_role('button',name='确认提交',exact=True).click()
    await page.wait_for_function('(r) => window.__arcaneDebug.gameState.value.revision > r',arg=revision)
async def payment(page):
    req=await page.evaluate('window.__arcaneDebug.paymentRequest.value')
    if not req: return
    names={'火':'火焰','水':'水纹','气':'大气','地':'地脉','光':'光辉','暗':'暗影','无':'奥术'}
    first=next((e for e,n in req['cost'].items() if n>0),None)
    if first:
        token=page.locator('.payment-token').filter(has=page.locator('span',has_text=names[first]))
        await token.click()
        await page.get_by_role('button',name='− '+names[first],exact=True).click()
        assert await token.locator('strong').inner_text()=='0'
        await token.click()
        await page.get_by_role('button',name='重置付款',exact=True).click()
        assert await token.locator('strong').inner_text()=='0'
    for elem,amount in req['cost'].items():
        for _ in range(amount):
            await page.locator('.payment-token').filter(has=page.locator('span',has_text=names[elem])).click()
    await page.get_by_role('button',name='确认支付',exact=True).click()
async def start(browser):
    room=api('/api/test-room/create',{})['room_id']; pages=[]; errors=[]
    for i in range(2):
        page=await browser.new_page(viewport={'width':1440,'height':1000},ignore_https_errors=True)
        page.set_default_timeout(12000)
        await page.add_init_script("""window.__sent=[]; const NativeWS=WebSocket; window.WebSocket=class extends NativeWS { constructor(...args){super(...args);window.__testSocket=this;} send(data){window.__sent.push(JSON.parse(data));super.send(data);} };""")
        page.on('console',lambda m: (errors.append(m.text),print('CONSOLE',m.text,flush=True)) if m.type=='error' and 'Failed to load resource' not in m.text else None)
        page.on('pageerror',lambda e: (errors.append(str(e)),print('PAGEERROR',str(e),flush=True)))
        await page.goto(BASE+'/game.html?'+urllib.parse.urlencode(dict(room=room,player_id=f'front-{i}',player_name=['星火','月影'][i],deck_code=DECK)),wait_until='domcontentloaded')
        pages.append(page)
    for page in pages:
        await page.get_by_role('button',name='保留手牌',exact=True).click()
        await commit(page)
    await pages[0].wait_for_function('window.__arcaneDebug.gameState.value.phase === "main"')
    return room,pages,errors
async def add(room,player,num,zone,**kwargs):
    api('/api/test-room/add-card',dict(room_id=room,player=player,card_number=num,zone=zone,count=1,**kwargs))
async def main():
 from frontend_card_visibility import decision_views
 async with async_playwright() as pw:
    browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
    room,pages,errors=await start(browser)
    actor=(await state(pages[0]))['current_turn']; a=pages[actor]; b=pages[1-actor]
    # Consume is selected on the actual battlefield and confirmed in the workspace.
    await a.locator('.my-board').get_by_role('button',name='消耗',exact=True).first.click()
    if await a.get_by_role('button',name='确认消耗',exact=True).count():
        await a.locator('.payment-token').first.click(); await a.get_by_role('button',name='确认消耗',exact=True).click()
    await commit(a)
    for i in range(2): api('/api/test-room/elements',dict(room_id=room,player=i,mode='set',elements={e:30 for e in ['火','水','气','地','光','暗','无']}))
    await a.wait_for_function('window.__arcaneDebug.gameState.value.you.elements["气"] === 30')
    # Learn via visible skill pool, exercising exact editable payment.
    await a.locator('.sp-card').filter(has_text='雷击').get_by_role('button',name='学习',exact=True).click()
    await payment(a); await commit(a)
    skill=next(c for c in (await state(a))['you']['skills'] if c and c['name']=='雷击')
    api('/api/test-room/card-state',dict(room_id=room,player=actor,instance_id=skill['instance_id'],horizontal=False))
    await add(room,1-actor,'1021011','unit',col=0,row=0)
    # cast from actual button, target actual enemy unit.
    await a.locator('.my-board .slot').filter(has=a.locator('img[alt="雷击"]')).get_by_role('button',name='施放',exact=True).click()
    await a.locator('.opponent-board .unit-cell').filter(has_text='屠魔者杀手').click()
    await payment(a); await commit(a)
    await b.wait_for_function('window.__arcaneDebug.showDefenseWindow.value')
    assert await a.locator('.defense-panel').count()==0
    await decision_views(b, 'defense')
    await b.screenshot(path=str(OUT/'defense-desktop.png'),full_page=True)
    await b.get_by_role('button',name='不防御',exact=True).click()
    await commit(b)
    # End-turn discard through actual candidate controls.
    for _ in range(8): await add(room,actor,'1021001','hand')
    await a.get_by_role('button',name='结束回合',exact=True).click(); await commit(a)
    pa=(await state(a))['pending_action']; assert pa and pa['type']=='discard'
    await a.screenshot(path=str(OUT/'discard-debug.png'),full_page=True)
    assert (await state(b)).get('pending_action') is None
    for i in range(pa['min_select']): await a.locator('.pending-card').nth(i).click()
    await a.get_by_role('button',name='确认选择',exact=True).click(); await commit(a)
    assert (await state(a))['current_turn'] != actor
    # Both clients see terminal outcome; no browser confirm dialog.
    await a.get_by_role('button',name='求和',exact=True).click(); await commit(a)
    await b.get_by_role('button',name='接受和局',exact=True).click(); await commit(b)
    for page in pages: await page.wait_for_function('window.__arcaneDebug.gameState.value.phase === "game_over"')
    assert not errors,errors
    print('PASS smoke: mulligan, consume, learn/payment, cast, defense handoff, discard, end-turn, game-over',flush=True)
    # Representative continuous dynamic selection, first target dies.
    room,pages,errors=await start(browser); actor=(await state(pages[0]))['current_turn']; a=pages[actor]; b=pages[1-actor]
    await add(room,actor,'3611101','skill_slot')
    for col in [0,2]: await add(room,actor,'1621109','unit',col=col,row=0)
    await add(room,1-actor,'1021001','unit',col=0,row=0)
    await add(room,1-actor,'1021002','unit',col=2,row=2)
    api('/api/test-room/elements',dict(room_id=room,player=actor,mode='set',elements={e:30 for e in ['火','水','气','地','光','暗','无']}))
    await a.wait_for_function('window.__arcaneDebug.gameState.value.you.elements["暗"] === 30')
    await a.locator('.my-board .slot').filter(has=a.locator('img[alt="红月"]')).get_by_role('button',name='施放',exact=True).click()
    await a.get_by_role('button',name='准备施放',exact=True).click()
    await payment(a); await commit(a)
    pa=(await state(a))['pending_action']; first=pa['decision_id']; assert pa['type']=='scarlet_wings_red_moon_damage'
    assert (await state(b)).get('pending_action') is None
    await a.screenshot(path=str(OUT/'choice-desktop.png'),full_page=True)
    candidate=a.locator('.pending-card').filter(has_text='巫师的学徒'); await candidate.click()
    # Single select replacement remains possible without deselecting first.
    await a.locator('.pending-card').filter(has_text='猩红之翼').first.click(); await candidate.click()
    await a.get_by_role('button',name='确认选择',exact=True).click(); await commit(a)
    pa=(await state(a))['pending_action']; assert pa['decision_id']>first
    assert not any(c.get('name')=='巫师的学徒' and c.get('side')=='enemy' for c in pa['candidates'])
    assert await a.locator('.pending-card.selected').count()==0
    await a.set_viewport_size({'width':390,'height':844})
    await a.wait_for_timeout(500)
    await a.screenshot(path=str(OUT/'choice-mobile.png'),full_page=True)
    from PIL import Image
    assert Image.open(OUT/'choice-mobile.png').width == 390
    assert await a.evaluate('Math.max(document.documentElement.scrollWidth,document.body.scrollWidth) <= 390'), 'mobile horizontal overflow'
    # Retry and duplicate click through central transport, with deliberate invalid draft.
    await a.evaluate('window.__arcaneDebug.sendAction("resolve_action", {selected:["invalid"]})')
    await a.get_by_role('button',name='确认提交',exact=True).click()
    await a.locator('.decision-error').wait_for()
    await a.locator('.action-review').wait_for(state='hidden')
    assert await a.locator('.action-review').count()==0
    assert await a.locator('.decision-error').inner_text()
    pa=(await state(a))['pending_action']
    candidate=next(c for c in pa['candidates'] if c.get('side')=='enemy')
    await a.locator(f'.pending-card[data-instance-id="{candidate["instance_id"]}"]').click()
    await a.get_by_role('button',name='确认选择',exact=True).click()
    rev=(await state(a))['revision']
    await a.get_by_role('button',name='确认提交',exact=True).evaluate('(b)=>{b.click();b.click()}')
    await a.wait_for_function('(r)=>window.__arcaneDebug.gameState.value.revision>r',arg=rev)
    assert (await state(a))['revision']==rev+1
    assert not errors,errors
    # Deliberately disconnect immediately after a confirmed action. Never replay it.
    await a.evaluate('window.__arcaneDebug.sendAction("offer_draw", {})')
    await a.get_by_role('button',name='确认提交',exact=True).evaluate('(b)=>{b.click();window.__testSocket.close()}')
    await a.get_by_role('button',name='重连并同步',exact=True).wait_for()
    assert await a.evaluate('window.__sent.filter(m=>m.action === "offer_draw").length')==1
    await a.get_by_role('button',name='重连并同步',exact=True).click()
    await a.wait_for_function('window.__arcaneDebug?.gameState.value')
    assert await a.evaluate('window.__sent.filter(m=>m.action !== "client_log").length')==0
    print('PASS disconnect: unknown-result action is not replayed after reconnect',flush=True)
    print('PASS continuous choices: two wings, target death refresh, replacement, private wait, retry, duplicate click, mobile',flush=True)
    (OUT/'result.json').write_text(json.dumps({'passed':True,'rooms':[room],'viewports':[[1440,1000],[390,844]],'page_errors':errors},ensure_ascii=False,indent=2))
    await browser.close()
if __name__ == '__main__': asyncio.run(main())
