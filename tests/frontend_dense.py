"""Authoritative high-density board screenshots; run after frontend_match.py."""
import asyncio
from frontend_match import *
async def dense():
 async with async_playwright() as pw:
  browser=await pw.chromium.launch(executable_path='/usr/bin/chromium',headless=True,args=['--no-sandbox'])
  room,pages,errors=await start(browser)
  for owner in range(2):
   for col in range(3):
    for row in range(3):
     if (col,row)==(1,1): continue
     await add(room,owner,['1021014','1621109','1021002'][(col+row)%3],'unit',col=col,row=row)
   for num in ['3611101','3321002','3221001','3121001','3001001']: await add(room,owner,num,'skill_slot')
   for num in ['2021013','2311002','2321001','2011001','2011002']: await add(room,owner,num,'equipment')
  a=pages[(await state(pages[0]))['current_turn']]
  await a.wait_for_timeout(3000)
  await a.screenshot(path=str(OUT/'dense-desktop.png'),full_page=True)
  await a.set_viewport_size({'width':390,'height':844})
  await a.wait_for_timeout(500)
  assert await a.evaluate('Math.max(document.documentElement.scrollWidth,document.body.scrollWidth) <= 390')
  await a.screenshot(path=str(OUT/'dense-mobile.png'),full_page=True)
  from PIL import Image
  assert Image.open(OUT/'dense-mobile.png').width == 390
  assert not errors,errors
  print('PASS dense: both full battlefields, 10 skills, 10 equipment, desktop/mobile, no page overflow')
  await browser.close()
if __name__=='__main__': asyncio.run(dense())
