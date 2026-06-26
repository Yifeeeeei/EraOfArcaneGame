# Mechanism Validation Log

This file tracks manual frontend validation performed after the base-set checkpoint merge.

## Results

- Backend machine-marked prompt audit: PASS
  - Issue #59 / PR #58 checked the current machine-test list for explicit `PendingAction` prompt coverage.
  - `3021006 洞察之眼` now opens a selection prompt when cast instead of destroying the first enemy equipment/set card by default.
  - Added backend assertions for target/mode prompts on `预见`, `万灵药`, `黑市商贩`, `新生卷轴`, `血魔爆`, `元素附魔`, `专精法师`, `伦德萨尔`, `食腐者`, and `回魂术`.
  - `cd server && go test ./...` passes.
  - Remaining risk: browser manual testing should still watch chained prompts, cancel branches, and visual clarity of target/mode windows.
- Binding / generated skill cleanup: PASS
  - Frontend summoned `"风暴之女" 艾拉雅`.
  - Card detail showed `风暴之怒` under bound skills.
  - State confirmed `风暴之怒` was attached to the host and absent from the skill pool.
  - Opponent cast `陨石术`; defender clicked `不防御` in the frontend defense window.
  - After death, the battlefield slot was empty and the graveyard contained only 艾拉雅; the bound skill did not enter the skill pool or graveyard.
- Devour / persistent cost reduction (`风暴奇美拉`): PASS
  - Frontend summoned `风暴奇美拉`; the UI opened the devour selection before placement could finish.
  - Selected friendly `风息谷雷鸟` with load `大气3`; after confirmation the thunderbird entered the graveyard and the chimera occupied the chosen battlefield cell.
  - The chimera remained a passive unit with `引魔` and no per-turn ability button.
  - With only `1 大气` available, the player cast `气旋波` whose printed use cost is `2 大气`; state entered the defense window and the player's `大气` was reduced to 0, confirming the passive `大气` spell cost reduction.
- Revive / sacrifice / graveyard movement / payment / marker timing batch: PASS
  - `人鱼之泪`: frontend clicked equipment `绝技`, selected a dead `守护骑士`, and the knight revived onto the field with `current_life=1`; `人鱼之泪` left the equipment zone.
  - `回收小精灵`: frontend summoned it, selected a graveyard card in the pending-action modal, and that card left the graveyard while deck count increased by 1.
  - `灵魂祭司`: frontend clicked unit `绝技`, selected a friendly `巫师的学徒`, and the apprentice entered the graveyard while hand size increased by 2.
  - Arcane payment choice: with `1 火焰` and `1 水纹`, frontend summoning a `1 无` card opened the payment modal; choosing `火焰` spent only fire and left water untouched.
  - Cooldown timing: a horizontal `冷却1` skill stayed horizontal after its controller ended turn, while the `冷却` marker cleared. This confirms reset happens before marker settlement, so a cooldown card is not prematurely ready next turn.
- Targeting / search-shuffle / hero ultimate / defense-overexert batch: PASS
  - `雷电元素`: frontend summoned it while the opponent had a front-row `守护骑士`; the target gained `眩晕1`.
  - `眺望者商舰`: frontend clicked `回合技`, selected a deck `人鱼 菲尔`, then selected a hand card to shuffle back; the searched water card entered hand and deck count was net unchanged.
  - `掌门 穆伶`: frontend clicked hero `绝技`, selected one friendly and one enemy companion, paid the `大气` cost difference, and both companions returned to their owners' hands.
  - Defense overexert: opponent cast an attack spell; defender selected both a vertical load unit and a defense spell in the defense window. The unit became horizontal, the defense skill became horizontal, no leftover load was stored in the element pool, and the target took no damage.
- Mastery / negative marks / reaction sorcery / area targeting batch: PASS
  - `"知识古树" 深耕`: frontend summoned it while `成长的树人` and `森林守卫` were on the field; both cards received their own `精通` markers up to max and applied their own threshold bonuses, while player `charge` stayed 0.
  - `流沙法师` and `冰域恶魔`: frontend summon applied `石化1` and `冻结1` to the expected opponent cards.
  - `冰封消解`: frontend reaction in the defense window reduced pending spell power from 3 to 2, tapped the reaction skill, and paid `2 水纹`.
  - `风洞`: frontend reaction canceled a single-target spell before hit resolution; the target remained undamaged.
  - `玄冰阵`: after selecting the spell, the frontend exposed only the three front-row cells as legal core targets; resolving the spell froze the splash-cross affected units, including the hero behind the chosen front-row target.
- Deathrattle / complex passive item batch: PASS
  - `随风旅行者`: after being killed by a frontend-cast spell and choosing `不防御`, the traveler entered the graveyard and its deathrattle drew the card placed on top of the deck.
  - `统御者之冠`: with the crown equipped, a frontend-summoned friendly companion entered with `elements_gain={}`. The companion's own enter effect still resolved, so the crown only clears load and does not suppress unrelated enter effects.
