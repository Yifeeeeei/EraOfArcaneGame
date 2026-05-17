# 基础包卡牌效果覆盖表

> 这个文件由 `node tools/effect-coverage-report.js --write` 生成，用来追踪基础包卡牌效果是否有运行时代码、后端语义测试、以及前端操作记录。

## 总览

- 基础包卡牌总数：378
- 已实现专属行为：260
- 只依赖通用机制：56
- 白板/无效果：62
- 需复核：0
- 未实现：0
- 后端语义测试已点名：310
- 后端语义测试未点名：0
- 前端已有操作记录：378
- 高风险待验证：0

## 按类型统计

| 类型 | 已实现 | 通用机制 | 白板/无效果 | 需复核 | 未实现 |
|---|---:|---:|---:|---:|---:|
| 人物 | 20 | 0 | 0 | 0 | 0 |
| 伙伴 | 94 | 2 | 35 | 0 | 0 |
| 技能 | 55 | 41 | 13 | 0 | 0 |
| 道具 | 91 | 13 | 14 | 0 | 0 |

## 未实现机制分布

| 机制族 | 未实现卡数 |
|---|---:|
| 选择/待选动作 | 0 |
| 检索/看牌/洗牌/牌堆顶 | 0 |
| 绑定技能/衍生牌 | 0 |
| 费用修改 | 0 |
| 状态/标记/计数器 | 0 |
| 吞噬/献祭/死亡触发 | 0 |
| 范围/目标规则 | 0 |
| 主动能力 | 0 |
| 反制/触发监听 | 0 |

## 状态说明

- `已实现`：已在 `server/game/card_effects_catalog.go` 注册专属 Go 行为。
- `通用机制`：目前没有专属行为，但描述主要落在引擎已有关键词机制上。仍需要前端实测。
- `白板/无效果`：描述为空，可以按普通卡运行。
- `需复核`：描述不为空，但脚本无法判断是否需要代码。
- `未实现`：描述明显要求行为代码或尚不存在的通用机制。
- `后端语义测试` 的 `已点名` 表示测试文件里直接出现过该卡编号；它不是完整正确性的证明，但能区分“只被全量 smoke 扫过”和“有专门语义断言”。
- `前端操作` 来自 `tmp/frontend-card-operation-report*.json`，是历史浏览器操作记录；旧失败不一定代表当前仍失败，但必须重新验证。
- `风险` 主要用于排队：已实现但没有语义测试、未实现、需复核的卡优先处理。

## 下一批高风险清单

| 编号 | 名称 | 类型 | 运行时代码 | 后端语义测试 | 前端操作 | 描述 |
|---|---|---|---|---|---|---|

## 全量清单

| 编号 | 名称 | 类型 | 属性 | 运行时代码 | 后端语义测试 | 前端操作 | 风险 | 通用机制 | 描述 |
|---|---|---|---|---|---|---|---|---|---|
| 1011001 | 魔龙 奥瑞 | 伙伴 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report.json) | 低 | 引魔 | 引魔.绑定技能:破灭魔光 |
| 1011002 | 巫师之塔 通天阁 | 伙伴 | 无 | 已实现 | 已点名(11) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 | 引魔 | 引魔.你的法力范围变为全场 |
| 1011003 | 盟主 法罗兰克 | 伙伴 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-retry-expensive-neutral.json) | 低 |  | 入场:获得等同于所有相邻伙伴的负载.绑定技能:纯净奥术 |
| 1021001 | 巫师的学徒 | 伙伴 | 无 | 白板/无效果 | 已点名(128) | 已操作(frontend-card-operation-report-full.json) | 低 |  |  |
| 1021002 | 学院导师 | 伙伴 | 无 | 白板/无效果 | 已点名(36) | 已操作(frontend-card-operation-report-retry-expensive-neutral.json) | 低 |  |  |
| 1021003 | 誓约巫师 | 伙伴 | 无 | 白板/无效果 | 已点名(8) | 已操作(frontend-card-operation-report-retry-expensive-neutral.json) | 低 |  |  |
| 1021004 | 守护骑士 | 伙伴 | 无 | 白板/无效果 | 已点名(40) | 已操作(frontend-card-operation-report-full.json) | 低 |  |  |
| 1021005 | 内阁巫师 | 伙伴 | 无 | 白板/无效果 | 已点名(4) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 |  |  |
| 1021006 | 杂货商贩 | 伙伴 | 无 | 已实现 | 已点名(7) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:抽2张牌 |
| 1021007 | 回收小精灵 | 伙伴 | 无 | 已实现 | 已点名(15) | 已操作(frontend-card-operation-report-full.json) | 低 |  | 入场:将你弃牌堆的1张牌放到卡组顶 |
| 1021008 | 预见先知 | 伙伴 | 无 | 已实现 | 已点名(4) | 已操作(frontend-card-operation-report-full.json) | 低 |  | 回合开始抽牌前,你可以查看牌堆顶的1张牌,将其放回牌堆顶或牌堆底 |
| 1021009 | 竞技场虚像 | 伙伴 | 无 | 已实现 | 已点名(5) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 不会受到法术攻击以外的伤害 |
| 1021010 | 专精法师 | 伙伴 | 无 | 已实现 | 已点名(5) | 已操作(frontend-card-operation-report-retry-expensive-neutral.json) | 低 |  | 入场:选择1个属性,此卡的负载变为该种属性 |
| 1021011 | 屠魔者杀手 | 伙伴 | 无 | 通用机制 | 已点名(4) | 已操作(frontend-card-operation-report-companions.json) | 低 | 速攻 | 速攻. |
| 1021012 | 黑市商贩 | 伙伴 | 无 | 已实现 | 已点名(5) | 已操作(frontend-card-operation-report-full.json) | 低 |  | 绝技:从你的手牌或者装备区弃置1张道具牌,抽2张牌 |
| 1021013 | 屠魔者武士 | 伙伴 | 无 | 白板/无效果 | 已点名(16) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 |  |  |
| 1021014 | 急不可耐的小师弟 | 伙伴 | 无 | 已实现 | 已点名(4) | 已操作(frontend-card-operation-report-companions.json) | 低 | 速攻 | 入场:本回合你学习的下一个技能获得"速攻" |
| 1021015 | 精力充沛的大师兄 | 伙伴 | 无 | 已实现 | 已点名(5) | 已操作(frontend-card-operation-report-full.json) | 低 | 冷却 | 入场:本回合你施放的下一个技能不需要冷却 |
| 1021016 | 奥术盔甲匠 | 伙伴 | 无 | 已实现 | 已点名(5) | 已操作(frontend-card-operation-report-full.json) | 低 |  | 入场:如果你没有装备,检索1个装备道具 |
| 1021017 | 符文师 | 伙伴 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-full.json) | 低 |  | 入场:丢弃1张手牌,检索1个符文或卷轴 |
| 1111001 | 火龙 "辉煌" | 伙伴 | 火 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-full.json) | 低 | 引魔 | 吞噬:3\火.引魔.绑定技能:火焰吐息 |
| 1111002 | 炎狱大将军 狄斯托德 | 伙伴 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 每当对方召唤1个伙伴,使其获得点燃1和石化1 |
| 1111003 | 毕方 | 伙伴 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 | 引魔、状态容器 | 引魔.敌方单位受到的点燃伤害+1 |
| 1121001 | 火焰精灵 | 伙伴 | 火 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-full.json) | 低 | 状态容器、命中状态 | 每当此卡被消耗,获得点燃1 |
| 1121002 | 活泼的炉火 | 伙伴 | 火 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-full.json) | 低 |  | 入场:抽1张牌 |
| 1121003 | 锻石工匠 | 伙伴 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-full.json) | 低 |  | 消耗:使你的1个法术在本回合+2\威 |
| 1121004 | 凯尔特雄狮 | 伙伴 | 火 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-companions.json) | 低 | 法术攻威被动 | 你的所有法术+1\威 |
| 1121005 | 熔岩傀儡 | 伙伴 | 火 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1121006 | 熔岩烽蛇 | 伙伴 | 火 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1121007 | 至纯之火 | 伙伴 | 火 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1121008 | 炎狱使者 | 伙伴 | 火 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1121009 | 赤鹰 | 伙伴 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:检索1个入场花费大于4的火焰伙伴 |
| 1121010 | 火焰艺人 | 伙伴 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 绝技:消耗此卡,重置你的另1张人物牌以外的火焰牌 |
| 1121011 | 火山飞龙 | 伙伴 | 火 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1121012 | 火焰洞察者 | 伙伴 | 火 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 回合技:若有单位受到火焰伤害,抽1张牌 |
| 1121013 | 纵火者 | 伙伴 | 火 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 回合技:在你使用1个火焰法术后,使法力范围内的1个单位点燃1 |
| 1121014 | 火荆 | 伙伴 | 火 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 遗言:使法力范围内的1个敌人点燃1 |
| 1121015 | 火云法师 | 伙伴 | 火 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1211001 | 人鱼 菲尔 | 伙伴 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 祈咒:如果此卡相邻没有伙伴,检索1张水纹伙伴 |
| 1211002 | 深渊巨口 利维坦 | 伙伴 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 消耗:消灭法力范围内1个伙伴,下个你的回合不能使用此效果 |
| 1211003 | "雪女" 天户凌 | 伙伴 | 水 | 已实现 | 已点名(5) | 已操作(frontend-card-operation-report-companions.json) | 低 | 引魔、状态容器、命中状态 | 引魔.回合技3:在你检索1张水纹卡牌后,使1个前排敌人冻结1 |
| 1221001 | 海豚伙伴 | 伙伴 | 水 | 已实现 | 已点名(7) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 献祭:防止1个其他友方单位受到的1次致命伤害 |
| 1221002 | 冰原法师 | 伙伴 | 水 | 白板/无效果 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1221003 | 掠夺者海盗 | 伙伴 | 水 | 白板/无效果 | 已点名(3) | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1221004 | 寒霜傀儡 | 伙伴 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 入场:对法力范围内1个敌方伙伴冻结1 |
| 1221005 | 西境海妖 | 伙伴 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 祈咒:消耗法力范围内的1个伙伴 |
| 1221006 | 水栖狸猫 | 伙伴 | 水 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 在本卡相邻有2个及以上水纹伙伴时,本卡负载+1\水 |
| 1221007 | 冰原狼 | 伙伴 | 水 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1221008 | 冰域恶魔 | 伙伴 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 入场:对法力范围内的所有敌人冻结1 |
| 1221009 | 南海海怪 | 伙伴 | 水 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1221010 | 护壁者 | 伙伴 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:直到下个回合结束所有法术\攻变为0 |
| 1221011 | 凛冬城术士 | 伙伴 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 绝技:本回合你的下一次法术获得冻结1 |
| 1221012 | 龙王子裔 | 伙伴 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 精通2:检索1个水纹伙伴并使其入场花费减少1\水 |
| 1221013 | 唤雨师 | 伙伴 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 | 法术攻威被动 | 你的水纹和大气法术+1\威 |
| 1221014 | 北海飞鱼 | 伙伴 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 临时 | 回合技:负载临时改为1\气 |
| 1221015 | 眺望者商舰 | 伙伴 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 祈咒:检索1个水纹卡牌,然后将1张手牌洗回卡组 |
| 1311001 | 大鹏 | 伙伴 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:翻开卡组顶8张牌,抽取其中入场花费小于3的卡牌,重洗你的卡组.回合结束时丢弃这些这些被抽取的卡牌 |
| 1311002 | "风暴之女" 艾拉雅 | 伙伴 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 绑定技能:风暴之怒.在你的手牌数大于等于手牌上限时,风暴之怒视为已经生效 |
| 1311003 | "风刃" 卡琳娜 | 伙伴 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 穿透 | 你的没有穿透的大气技能获得穿透和使用花费+1\气(不需要选择目标的技能不受影响) |
| 1321001 | 渡鸦信使 | 伙伴 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 消耗:抽1张牌 |
| 1321002 | 随风旅行者 | 伙伴 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 获得元素 | 入场:获得2\气.遗言:抽1张牌 |
| 1321003 | 魔法蒲公英 | 伙伴 | 气 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 当你抽到此卡时,将其展示.入场:如果你在本回合抽到此卡,抽1张牌 |
| 1321004 | 雷电元素 | 伙伴 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 入场:使法力范围内1个伙伴晕眩1 |
| 1321005 | 驭风师 | 伙伴 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 | 获得元素 | 绝技:丢弃任意数量的手牌,每张牌使你获得1\气 |
| 1321006 | 雷霆兽 | 伙伴 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 法术攻威被动 | 你的大气法术+1\攻 |
| 1321007 | 工蜂骑士 | 伙伴 | 气 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1321008 | 风息奔马 | 伙伴 | 气 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1321009 | 风魔 | 伙伴 | 气 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1321010 | 风暴奇美拉 | 伙伴 | 气 | 已实现 | 已点名(5) | 已操作(frontend-card-operation-report-companions.json) | 低 | 引魔 | 引魔.吞噬:3\气.你的大气法术使用花费减少1\气 |
| 1321011 | 雷精灵 | 伙伴 | 气 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1321012 | 风灵媒师 | 伙伴 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 回合技:在你使用1个大气技能后,抽1张牌 |
| 1321013 | 传送法师 | 伙伴 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 回合技:移动1个友方单位 |
| 1321014 | 风息谷雷鸟 | 伙伴 | 气 | 白板/无效果 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1321015 | 风语者 | 伙伴 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 获得元素 | 回合技:当你丢弃手牌时,获得1\气 |
| 1401001 | 生命种子 | 伙伴 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 精通2:献祭此卡并从你的手牌中召唤1个地脉伙伴(无需花费),它会继承此卡的生命和负载加成 |
| 1401002 | 灵兽 辛柯 | 伙伴 | 地 | 已实现 | 已点名(4) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 当友方单位受到敌方伤害后,可以从卡组或手牌召唤此卡,无需入场花费 |
| 1411001 | "轮回不息" 大德鲁伊 烟尘 | 伙伴 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 绝技:当1个友方伙伴死亡时,召唤1个生命种子,它会继承该伙伴的所有生命和负载加成 |
| 1411002 | "知识古树" 深耕 | 伙伴 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:你的精通立刻达到最高 |
| 1411003 | 沙之魔巫 梭默 | 伙伴 | 地 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-companions.json) | 低 | 基础范围 | 你的没有范围效果的地脉法术获得范围:方阵 |
| 1421001 | 流沙法师 | 伙伴 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 入场:使1个无视范围的敌人石化1 |
| 1421002 | 祝祷祭师 | 伙伴 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 此卡和相邻单位不受负面状态影响(仍可处于负面状态) |
| 1421003 | 成长的树人 | 伙伴 | 地 | 已实现 | 已点名(8) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 精通2,4:负载+1\地或者+1\血 |
| 1421004 | 森林守卫 | 伙伴 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 精通1:+1\血.精通3:负载+1\地.精通5:+2\攻 |
| 1421005 | 磐石元素 | 伙伴 | 地 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1421006 | 林地变形者 | 伙伴 | 地 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1421007 | 高地泰坦 | 伙伴 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 |  | 未被强化的法术对本卡造成的伤害+1 |
| 1421008 | 岩山翼龙 | 伙伴 | 地 | 白板/无效果 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1421009 | 被祝福的少女 | 伙伴 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 祈咒:使1个相邻地脉伙伴获得负载+1\地 |
| 1421010 | 种植园丁 | 伙伴 | 地 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 你的卡牌每次获得负载,在此卡上放置1个标记.回合技:取除2个标记,抽1张牌 |
| 1421011 | 大长老 | 伙伴 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 精通1,3:下一次学习地脉技能的花费-2 |
| 1421012 | 林地飞鼠 | 伙伴 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 临时 | 回合技:负载临时改为1\气 |
| 1421013 | 岩山恐兽 | 伙伴 | 地 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1421014 | 风息谷旅商 | 伙伴 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:你的场上每有1个野兽,精灵或植物,抽1张牌(最多3张) |
| 1421015 | 苍绿之龙 | 伙伴 | 地 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1501001 | 孪生天使 | 伙伴 | 光 | 白板/无效果 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1511001 | 白袍大贤者 掌号使 | 伙伴 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 绝技:选择法力范围内的1个敌方伙伴,支付其入场花费,获得其控制权 |
| 1511002 | 大法师 伦德萨尔 | 伙伴 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场,遗言:使你的一个法术永久获得+3\威或+1\攻 |
| 1511003 | 天枢圣兽 珀伽索斯 | 伙伴 | 光 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-companions.json) | 低 | 引魔 | 引魔.对方使用法术攻击时:该法术对天枢圣兽 珀伽索斯以外的友方单位造成伤害变为0 |
| 1521001 | 治疗术士 | 伙伴 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 祈咒:使1个友方单位回复1\血 |
| 1521002 | 光铸泰坦 | 伙伴 | 光 | 已实现 | 已点名(4) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:抽2张牌.神秘和聚能法术对本卡造成的伤害+1 |
| 1521003 | 七神侍从 | 伙伴 | 光 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1521004 | 誓约之泉的守卫 | 伙伴 | 光 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1521005 | 双生天使 | 伙伴 | 光 | 已实现 | 已点名(7) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:将1张孪生天使置于你的手牌 |
| 1521006 | 生命之花 | 伙伴 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:使1个其他友方单位+1\血 |
| 1521007 | 虹之天使 | 伙伴 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 你的光辉元素可以当做任意元素使用 |
| 1521008 | 御座的圣翼 | 伙伴 | 光 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1521009 | 天马骑士 | 伙伴 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:检索1张独角天马 |
| 1521010 | 神护者 | 伙伴 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 免疫负面状态 |
| 1521011 | 日轮法师 | 伙伴 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 绝技:重置你的1个光辉法术 |
| 1521012 | 独角天马 | 伙伴 | 光 | 白板/无效果 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1521013 | 神火兽 | 伙伴 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 | 法术攻威被动 | 你的法术在攻击时+2\威 |
| 1521014 | 炬之女巫 | 伙伴 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 入场:本卡获得点燃2.祈咒:使1个相邻伙伴获得负载+1\光 |
| 1521015 | 烬之女巫 | 伙伴 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 入场:本卡获得点燃3.遗言:使你的1个法术永久+2\威 |
| 1611001 | "观察者" 欧柯茹 | 伙伴 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:查看卡组顶5张牌,你可以将其抽取或以任意顺序放回卡组顶部、底部,每抽取1张,对你的人物造成1点伤害 |
| 1611002 | 黑袍执行官 无心 | 伙伴 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 每当你献祭或吞噬1个伙伴,根据其生命值在此卡上放置暗影标记物.绝技:选择法力范围内的1个伙伴,取除其生命值数量的暗影标记物并将其消灭 |
| 1611003 | "穿心人" | 伙伴 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 入场:将1张衍生道具幻痛加入手牌.幻痛在触发时可以额外选择1个敌方法术 |
| 1621001 | 冥界信鸽 | 伙伴 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 遗言:抽1张牌 |
| 1621002 | 元素躯壳 | 伙伴 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 | 获得元素 | 遗言:获得1\无 |
| 1621003 | 恐惧魔 | 伙伴 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 吞噬:3\血 |
| 1621004 | 巫术祭司 | 伙伴 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 绝技:献祭你的1个伙伴,使另一个角色获得其生命值 |
| 1621005 | 诅咒魔像 | 伙伴 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 入场:使1个敌方法术获得虚弱2 |
| 1621006 | 梦魇 | 伙伴 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 每当其他友方单位死亡,此卡获得+1\血 |
| 1621007 | 巫师的人偶 | 伙伴 | 暗 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1621008 | 南境奴隶 | 伙伴 | 暗 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1621009 | 唤魔邪术士 | 伙伴 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 回合技:在你的1个伙伴死亡后,检索1个暗影造物或恶魔 |
| 1621010 | 恶魔尊主 | 伙伴 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 吞噬:4\暗 |
| 1621011 | 白骨骑士 | 伙伴 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 遗言:重新召唤此伙伴,并失去此遗言 |
| 1621012 | 灵魂祭司 | 伙伴 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-companions.json) | 低 |  | 绝技:献祭1个友方伙伴,抽2张牌 |
| 1621013 | 言灵 | 伙伴 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-companions.json) | 低 | 状态容器、命中状态 | 回合技:对方在1回合内使用第3个技能后,使那些技能虚弱1 |
| 1621014 | 恶魔仆从 | 伙伴 | 暗 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 1621015 | 人面枭 | 伙伴 | 暗 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-companions.json) | 低 |  |  |
| 2011001 | 大法师之杖 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 |  | 入场:从你的技能池将1个法术置于此卡上.绝技:花费元素使用此卡上的1个技能,然后将该卡牌从游戏中移除 |
| 2011002 | 统御者之冠 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 入场:此后你召唤的所有伙伴负载变为0 |
| 2011003 | 君王法袍 至贤 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 |  | 绝技:当敌方法术命中时,如果你场上的负载多于敌方,每多2点使本回合所有敌方法术-1\攻 |
| 2021001 | 秘法宝典 | 道具 | 无 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2021002 | 记忆项链 | 道具 | 无 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 你的技能槽位+1 |
| 2021003 | 随心魔杖 | 道具 | 无 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 |  | 消耗:将你的1个使用花费小于3的法术重置 |
| 2021004 | 巫师权杖 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 | 法术攻威被动 | 你的法术+1\威 |
| 2021005 | 瓶装元素 | 道具 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 获得元素 | 获得1\无 |
| 2021006 | 百宝锦囊 | 道具 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 献祭:从卡组检索1张消耗品道具牌 |
| 2021007 | 巫师齐射线列 | 道具 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 重置你的一个法术,下一次它的范围变成AOE:前排 |
| 2021008 | 魔法石 | 道具 | 无 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2021009 | 誓约之戒 | 道具 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 你的法术在攻击和强化攻击时-2\威 |
| 2021010 | 封印卷轴 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 如果敌方有4个及以上的技能,选择其中1个,使其直到下个回合结束不能使用 |
| 2021011 | 生命护符 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 入场:使1个友方角色+1\血 |
| 2021012 | 速写卷轴 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 释放1个你已经学习的法术并支付其使用花费,无需消耗 |
| 2021013 | 断绝之刃 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 | 防御限制、法术攻威被动 | 你的法术攻击和强化攻击时+2\威,你的法术无法用于防御 |
| 2021014 | 法力增强剂A型 | 道具 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 本回合你的下1次技能使用花费为0 |
| 2021015 | 法力增强剂C型 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 | 冷却 | 本回合你的法术使用花费为0,但在使用后获得冷却2 |
| 2021016 | 纹饰佩剑 | 道具 | 无 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2021017 | 旅行行囊 | 道具 | 无 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 你的道具槽位+3 |
| 2021018 | 奥术符文 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:当敌方使用法术时,使你的1个法术在本回合+3\威(敌方可以继续进行强化) |
| 2021019 | 诅咒卷轴 | 道具 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 抽2张牌,但在本回合结束时将那些牌丢弃 |
| 2021020 | 假面 | 道具 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 你的人物的负载变为等量的奥术元素\无 |
| 2021021 | 聚能卷轴 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 | 获得元素 | 在你的下个回合开始时获得3\无 |
| 2021022 | 反制符文 | 道具 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:当敌方使用卷轴或符文时,将其无效 |
| 2111001 | 火龙之心 | 道具 | 火 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 回合技:献祭最多3点\火,每1点火使下一次火焰法术获得+1\攻或者+3\威 |
| 2111002 | 努尔之眼 | 道具 | 火 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 | 获得元素、法术攻威被动 | 每当1个单位受到1点火焰伤害,放置1个火焰标记物.祈咒:移除此卡所有火焰标记物,根据数量执行以下效果.0个:摧毁此卡;1个:获得2\火;2个:本回合你的火焰法术+2\威;3个:本回合你的火焰法术+1\攻;4个及以上:造成2点火焰伤害(不放置标记物) |
| 2121001 | 凤凰之羽 | 道具 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 获得元素 | 入场:放置3个火焰标记物.回合技:取除1个火焰标记物,获得1\火 |
| 2121002 | 火焰符文 | 道具 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 反制:当有单位被消耗时,使其获得点燃1 |
| 2121003 | 灼烧卷轴 | 道具 | 火 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 点燃1 |
| 2121004 | 火焰箭 | 道具 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 献祭:对法力范围内1个敌人造成1点伤害 |
| 2121005 | 神炎魔咒药剂 | 道具 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态、法术攻威被动 | 直到下个回合结束你的法术+2\威,你的人物获得点燃1 |
| 2121006 | 火焰面甲 | 道具 | 火 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2121007 | 舞火战裙 | 道具 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 绝技:移除1个友方火焰单位所有负面状态 |
| 2121008 | 烈焰风暴卷轴 | 道具 | 火 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态、基础范围 | 范围:方阵.点燃1 |
| 2121009 | 烈焰障壁卷轴 | 道具 | 火 | 通用机制 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 防御限制 | 防御 |
| 2121010 | 炽火链鞭 | 道具 | 火 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2121011 | 火流星卷轴 | 道具 | 火 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 | 穿透 | 穿透 |
| 2121012 | 狱火符文 | 道具 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 反制:当敌方召唤1个伙伴时,使其获得晕眩2,石化2,点燃2 |
| 2121013 | 熔火战铠 | 道具 | 火 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2211001 | 人鱼之泪 | 道具 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 将此卡从游戏中移除:复活你的1个死亡伙伴但只有1\血 |
| 2211002 | 嗜魔弓 凛冬 | 道具 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 每当有玩家使用法术,可以花费1\水在此卡上放置1个水纹标记物.绑定技能:凛冬将至 |
| 2221001 | 冰霜之心 | 道具 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 敌方法术命中时献祭:该法术伤害变为0 |
| 2221002 | 冰霜符文 | 道具 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 反制:当有敌方伙伴被消耗时,使其冻结1 |
| 2221003 | 冰封卷轴 | 道具 | 水 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 使所有前排敌人冻结1 |
| 2221004 | 玛涅斯之杖 | 道具 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 法术攻威被动 | 你的水纹法术+1\威 |
| 2221005 | 精力药剂 | 道具 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:对方回合结束时,将你的全部法术重置 |
| 2221006 | 海之眷顾 | 道具 | 水 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2221007 | 凝霜手镯 | 道具 | 水 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2221008 | 水形之束卷轴 | 道具 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 命中:若目标为伙伴牌,消耗该伙伴 |
| 2221009 | 寒冰爆裂卷轴 | 道具 | 水 | 通用机制 | 已点名(4) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态、基础范围 | 范围:溅射.冻结1 |
| 2221010 | 潮涌符文 | 道具 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:当敌方在一个回合内抽第三张牌时,使你的1个水纹伙伴获得负载+2\水 |
| 2221011 | 恩惠之雨 | 道具 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:当1个友方单位受伤后,使所有友方单位回复2\血 |
| 2221012 | 水行之靴 | 道具 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 在你的人物与至少3个水纹伙伴相邻时,此卡负载+1\水 |
| 2221013 | 深寒诅咒卷轴 | 道具 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器 | 使法力范围内的1个敌方伙伴永久冻结 |
| 2311001 | 雷之源 | 道具 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 你的卡牌入场花费和使用花费减少1\气 |
| 2311002 | 唤雷震鼓 | 道具 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 每当你抽1张牌,可以将其展示并在此卡上放置1个标记.你无法在当回合使用展示过的卡牌.回合技:移除3个标记,本回合你的大气法术获得+1\攻或者晕眩1 |
| 2321001 | 风息罗盘 | 道具 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 | 临时 | 回合技3:当你抽到1张大气卡牌,你可以将其展示然后此卡临时获得负载1点\气 |
| 2321002 | 闪电符文 | 道具 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 反制:当有敌人被消耗时,使其与1个相邻单位晕眩1 |
| 2321003 | 雷暴卷轴 | 道具 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态、基础范围 | 范围:方阵.命中:使所有命中伙伴晕眩1 |
| 2321004 | 雷霆魔杖 | 道具 | 气 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2321005 | 唤风卷轴 | 道具 | 气 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 抽2张牌,但在下个你的回合开始不抽牌 |
| 2321006 | 瓶中闪电 | 道具 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态、获得元素 | 获得3\气,使1个友方单位获得晕眩2 |
| 2321007 | 风语之戒 | 道具 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 入场:抽1张牌 |
| 2321008 | 旋风卷轴 | 道具 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 摧毁敌方任意1个入场花费小于5的装备道具 |
| 2321009 | 连锁闪电卷轴 | 道具 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 命中:抽1张牌或者检索1张连锁闪电卷轴 |
| 2321010 | 幻术卷轴 | 道具 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:当敌方使用法术攻击时,重新排列你场上的所有单位,对方需要重新选择目标 |
| 2321011 | 传送符文 | 道具 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:当1个伙伴被召唤或消耗后,将其移动到另一位置 |
| 2321012 | 随风斗篷 | 道具 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 绝技:移动你的人物至另一位置 |
| 2321013 | 驭风杖 | 道具 | 气 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2411001 | 古树之心 | 道具 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 回合技:友方单位获得负载时使其+1\血,或获得生命时使其负载+1\地 |
| 2411002 | 裂地巨剑 阿托比斯 | 道具 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 消耗:本回合下一次法术获得+4\威且范围变为前排,或者+2\攻且范围变为纵列 |
| 2421001 | 知识古树的关怀 | 道具 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 | 获得元素 | 当你的卡牌达到精通时可以消耗此卡:抽1张牌并获得1\地 |
| 2421002 | 生长药水 | 道具 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 重置你的1个地脉伙伴 |
| 2421003 | 坚固卷轴 | 道具 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 直到下个回合结束,所有友方单位受到的法术伤害-1 |
| 2421004 | 德鲁伊水平测试 | 道具 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 你所有负载大于2的伙伴获得负载+1\地 |
| 2421005 | 石化卷轴 | 道具 | 地 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 使1个无视范围的单位石化2 |
| 2421006 | 磐藤胸甲 | 道具 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 入场:你的人物获得+2\血 |
| 2421007 | 寄生之触 | 道具 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 精通1:负载+1\地 |
| 2421008 | 巨石阵卷轴 | 道具 | 地 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 | 基础范围 | 范围:方阵 |
| 2421009 | 森林之矢卷轴 | 道具 | 地 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2421010 | 自然封印卷轴 | 道具 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 直到你下个回合的回合结束,所有法术\攻变为0 |
| 2421011 | 精灵铠 | 道具 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 祈咒:为你的人物回复1\血 |
| 2421012 | 地脉灵石 | 道具 | 地 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2421013 | 《地理学入门》 | 道具 | 地 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 你原始入场花费大于5的卡牌入场费用减少2\地,如果你装备了2张《地理学入门》,它们共计只能使入场费用减少3\地 |
| 2501001 | 桎梏 | 道具 | 光 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 |  | 当你抽到这张牌时(起始手牌除外),必须将其展示并丢弃,之后你可以再抽1张牌 |
| 2511001 | 万灵药 | 道具 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 | 获得元素 | 回复1个友方单位所有生命,或抽4张牌,或获得5\无,或重置你的1个技能 |
| 2511002 | 辉之盾 闪耀 | 道具 | 光 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 | 防御限制、状态容器、命中状态 | 你在防御时额外获得2\威.回合技:当你防御成功时,对法力范围内所有敌人造成晕眩1 |
| 2521001 | 生命药剂 | 道具 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 选择1个友方单位,使其回复2\血 |
| 2521002 | 庇护符文 | 道具 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:当1个威力小于10的敌方法术命中时,将其无效 |
| 2521003 | 净化卷轴 | 道具 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 移除1个友方卡牌所有负面状态或任意1个敌方卡牌所有标记物 |
| 2521004 | 神圣制裁卷轴 | 道具 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 | 咒术立即结算 | 反制:敌方使用咒术时,无效敌人的那个技能 |
| 2521005 | 新生卷轴 | 道具 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 |  | 复活你的1个死亡的光辉伙伴并花费相应的入场所需元素 |
| 2521006 | 绿玉权杖 | 道具 | 光 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2521007 | 蓝晶灯盏 | 道具 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 绝技:花费5\光,负载+2\光 |
| 2521008 | 惩戒之箭卷轴 | 道具 | 光 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 | 穿透 | 穿透 |
| 2521009 | 光之刃卷轴 | 道具 | 光 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 | 基础范围 | 范围:前排 |
| 2521010 | 闪耀水晶 | 道具 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 你的光辉法术获得晕眩1 |
| 2521011 | 闪光符文 | 道具 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 反制:当敌方使用技能时,使所有前排敌人晕眩1 |
| 2521012 | 幻彩颜料 | 道具 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 将你场上负载的最多4点\光变为\无 |
| 2521013 | 防护结界卷轴 | 道具 | 光 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 | 防御限制 | 防御 |
| 2601001 | 幻痛 | 道具 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 防御限制、状态容器、命中状态 | 回合技:当敌方使用法术防御成功后,使用于防御和强化防御的法术虚弱2 |
| 2601002 | 咒言书 | 道具 | 暗 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 入场:使所有敌方法术虚弱1 |
| 2611001 | 死灵魔石 虚无 | 道具 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 回合技:当1个友方伙伴死亡,此卡获得负载+1\暗 |
| 2611002 | 与恶魔的契约书 | 道具 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 献祭1个友方单位然后消灭法力范围内1个敌方伙伴,二者每相差1点\血必须额外支付2\暗.此卡在打出后洗回卡组 |
| 2621001 | 虚弱药剂 | 道具 | 暗 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器、命中状态 | 使敌方最多2个不同的法术虚弱2 |
| 2621002 | 巫毒娃娃 | 道具 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 在此卡上放置3个暗影标记物并选择法力范围内的2个伙伴,其一受到伤害时可以让另一者收到同等的伤害,并取除伤害数量的暗影标记物 |
| 2621003 | 杀戮本能 | 道具 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:当对手召唤1个伙伴时,对其造成2点伤害 |
| 2621004 | 暗影帷幕 | 道具 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 | 引魔 | 敌方法术命中时献祭:这个回合你的暗影伙伴不会受到法术伤害,但你的人物会获得引魔 |
| 2621005 | 献祭符文 | 道具 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:当1个伙伴死亡时,抽2张牌 |
| 2621006 | 亡魂项链 | 道具 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 获得元素 | 回合技:当你的1个伙伴死亡时,获得1\暗 |
| 2621007 | 安迪斯之镰 | 道具 | 暗 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-items.json) | 低 |  |  |
| 2621008 | 魂噬卷轴 | 道具 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器 | 命中:将3层虚弱分配给敌方法术 |
| 2621009 | 暗冥弹卷轴 | 道具 | 暗 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-items.json) | 低 | 基础范围 | 范围:溅射. |
| 2621010 | 拖入深渊 | 道具 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:当1个友方单位受到伤害且死亡后,对法力范围内的1个敌人造成等同于那个友方单位在本回合受到的全部伤害 |
| 2621011 | 狂乱符文 | 道具 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 反制:当1个具有攻击力的敌方伙伴消耗时,使那次消耗视为其对1个相邻单位的攻击 |
| 2621012 | 暗影披风 | 道具 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-items.json) | 低 |  | 入场:下一次命中的敌方法术伤害为0 |
| 2621013 | 巫术指环 | 道具 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-items.json) | 低 | 状态容器 | 回合技:当敌方1个法术受到虚弱时,使该虚弱层数+1 |
| 3001001 | 破灭魔光 | 技能 | 无 | 通用机制 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 | 基础范围 | 范围:前排.无法强化或被强化 |
| 3001002 | 纯净奥术 | 技能 | 无 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 花费最多10点同种元素,使下一次该属性法术威力上升那个数值 |
| 3021001 | 移形换影 | 技能 | 无 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 速攻 | 速攻.移动1个友方单位 |
| 3021002 | 预见 | 技能 | 无 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 查看牌堆顶3张牌,将其置于牌堆顶或牌堆底 |
| 3021003 | 冥想 | 技能 | 无 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 获得元素 | 获得1\无 |
| 3021004 | 刻印 | 技能 | 无 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却1.丢弃1张手牌,从卡组上方开始,抽取翻出的第1张卷轴或符文 |
| 3021005 | 奥术箭矢 | 技能 | 无 | 已实现 | 已点名(24) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 对法力范围内1个单位造成1点伤害 |
| 3021006 | 洞察之眼 | 技能 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却1.摧毁1张敌方盖放的卡牌 |
| 3021007 | 元素附魔 | 技能 | 无 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器 | 使你的下一次法术获得1点任意负面效果(点燃、冻结、麻痹、石化、虚弱) |
| 3021008 | 缴械 | 技能 | 无 | 已实现 | 已点名(5) | 已操作(frontend-card-operation-report-skills.json) | 低 | 速攻、冷却 | 速攻.冷却1.命中敌方单位:摧毁敌方的1个装备 |
| 3021009 | 昏睡 | 技能 | 无 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 速攻、状态容器、命中状态 | 速攻.命中:使目标伙伴晕眩1 |
| 3021010 | 解咒 | 技能 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却、防御限制 | 冷却1.当敌方使用防御法术时使用,将其无效 |
| 3021011 | 统御者的制裁 | 技能 | 无 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 | 穿透 | 穿透.此卡的学习和使用花费必须为同种元素 |
| 3021012 | 心炼 | 技能 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却1.使你的1个法术永久获得+3\威或者+1\攻 |
| 3101001 | 火焰吐息 | 技能 | 火 | 白板/无效果 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 |  |  |
| 3101002 | 万火合一术 | 技能 | 火 | 已实现 | 已点名(4) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 每5点\威获得+1\攻 |
| 3121001 | 火球术 | 技能 | 火 | 白板/无效果 | 已点名(25) | 已操作(frontend-card-operation-report-skills.json) | 低 |  |  |
| 3121002 | 焚烧 | 技能 | 火 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器 | 攻击时目标具有点燃则+2\威 |
| 3121003 | 炽热射线 | 技能 | 火 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器、命中状态 | 点燃2 |
| 3121004 | 燃烧大地 | 技能 | 火 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 基础范围 | 范围:前排 |
| 3121005 | 烈焰风暴 | 技能 | 火 | 通用机制 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器、命中状态、基础范围 | 范围:方阵.点燃1 |
| 3121006 | 陨石术 | 技能 | 火 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 穿透 | 穿透 |
| 3121007 | 激情之火 | 技能 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 速攻、冷却 | 速攻.冷却1.异能:你的火焰法术命中时抽1张牌 |
| 3121008 | 火焰结界 | 技能 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却、状态容器、命中状态 | 冷却1.异能:你的火焰法术获得点燃1和+2\威 |
| 3121009 | 爆焰一击 | 技能 | 火 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 |  |  |
| 3121010 | 岩浆爆发 | 技能 | 火 | 通用机制 | 已点名(4) | 已操作(frontend-card-operation-report-skills.json) | 低 | 穿透、状态容器、命中状态、基础范围 | 范围:方阵.穿透.点燃1 |
| 3121011 | 引燃 | 技能 | 火 | 通用机制 | 已点名(4) | 已操作(frontend-card-operation-report-skills.json) | 低 | 速攻、状态容器、命中状态 | 速攻.使1个敌方单位点燃1 |
| 3121012 | 烈焰护体 | 技能 | 火 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制 | 防御 |
| 3121013 | 烈焰反噬 | 技能 | 火 | 通用机制 | 已点名(4) | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制、状态容器、命中状态 | 防御.若防御成功,对敌方人物造成点燃1 |
| 3121014 | 烈焰重燃 | 技能 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 获得元素 | 本回合你每使用过1个火焰法术就获得1\火 |
| 3121015 | 焚风 | 技能 | 火 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 | 穿透 | 穿透.强化其他法术时使其获得穿透 |
| 3201001 | 百川归海 | 技能 | 水 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制 | 防御.若防御成功,获得等同于所有攻击法术的攻击力合计的\水 |
| 3201002 | 凛冬将至 | 技能 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 穿透、状态容器、命中状态、基础范围 | 范围:溅射.穿透.冻结1.使用时必须移除嗜魔弓 凛冬上的5个水纹标记物 |
| 3221001 | 冰雹术 | 技能 | 水 | 通用机制 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 | 基础范围 | 范围:方阵 |
| 3221002 | 冰锥术 | 技能 | 水 | 白板/无效果 | 已点名(4) | 已操作(frontend-card-operation-report-skills.json) | 低 |  |  |
| 3221003 | 激冻寒流 | 技能 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 强化其他水纹法术时+2\威 |
| 3221004 | 寒冰屏障 | 技能 | 水 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制 | 防御 |
| 3221005 | 玄冰阵 | 技能 | 水 | 通用机制 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器、命中状态、基础范围 | 范围:溅射.冻结1 |
| 3221006 | 海啸 | 技能 | 水 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 | 冷却、基础范围 | 冷却1.范围:全场 |
| 3221007 | 水占术 | 技能 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却1.查看牌堆顶4张牌并检索其中1张水纹卡牌,其余按任意顺序置于牌堆顶或牌堆底 |
| 3221008 | 冰封消解 | 技能 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却1.当对方使用法术时使用,使其中1个\威变为0 |
| 3221009 | 冰霜利刃 | 技能 | 水 | 已实现 | 已点名(4) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 攻击或强化攻击时+2\威 |
| 3221010 | 水幻影 | 技能 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却1.本回合你召唤的下一个水纹伙伴将会额外产生1个只有1\血的复制. |
| 3221011 | 幽影寒锋 | 技能 | 水 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 穿透 | 穿透 |
| 3221012 | 霜冻射线 | 技能 | 水 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器、命中状态 | 冻结2 |
| 3221013 | 猎潮 | 技能 | 水 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 |  |  |
| 3221014 | 坚冰领域 | 技能 | 水 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制、状态容器、命中状态 | 防御.若防御成功,使所有前排敌人冻结1 |
| 3221015 | 暴风雪 | 技能 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却、状态容器、命中状态 | 冷却1.异能:你的水纹和大气法术获得冻结1和+1\威 |
| 3301001 | 风暴之怒 | 技能 | 气 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 法术攻威被动 | 异能:展示你的所有手牌,每张使你的大气法术+1\威 |
| 3321001 | 闪电链 | 技能 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 可以额外选择1个无视范围的目标 |
| 3321002 | 雷击 | 技能 | 气 | 白板/无效果 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 |  |  |
| 3321003 | 静电脉冲 | 技能 | 气 | 通用机制 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器、命中状态 | 晕眩1 |
| 3321004 | 雷闪 | 技能 | 气 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 穿透 | 穿透 |
| 3321005 | 气旋波 | 技能 | 气 | 白板/无效果 | 已点名(17) | 已操作(frontend-card-operation-report-skills.json) | 低 |  |  |
| 3321006 | 雷暴术 | 技能 | 气 | 通用机制 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器、命中状态、基础范围 | 范围:方阵,晕眩1 |
| 3321007 | 源力之风 | 技能 | 气 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却1.补充手牌至手牌上限,每补1张花费1\气 |
| 3321008 | 风洞 | 技能 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却1.当敌方的非范围法术命中时使用,将其无效 |
| 3321009 | 宇宙飓风 | 技能 | 气 | 通用机制 | 已点名(1) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 | 穿透 | 穿透 |
| 3321010 | 涡旋屏障 | 技能 | 气 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制 | 防御 |
| 3321011 | 撕裂长空 | 技能 | 气 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 基础范围 | 范围:纵列 |
| 3321012 | 空天感应 | 技能 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 速攻、冷却 | 冷却1.速攻.异能:如果你的法术目标或区域中包含了非前排单位,使其获得+2\威 |
| 3321013 | 霹雳惊雷 | 技能 | 气 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 速攻、穿透 | 速攻.穿透 |
| 3321014 | 引雷 | 技能 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却、状态容器、命中状态 | 冷却1.丢弃1张手牌,使1个敌人晕眩1 |
| 3321015 | 静电屏障 | 技能 | 气 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制、状态容器、命中状态 | 防御.若防御失败,使1个前排敌人晕眩1 |
| 3421001 | 森林的庇护 | 技能 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制 | 防御.精通1:改为4\威;精通3:改为6\威 |
| 3421002 | 石化缠绕 | 技能 | 地 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器、命中状态 | 石化1. |
| 3421003 | 裂地重击 | 技能 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 精通1,3:获得+1\威和+1\攻 |
| 3421004 | 再生之力 | 技能 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 重置你的1张地脉伙伴 |
| 3421005 | 岩石壁垒 | 技能 | 地 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制 | 防御 |
| 3421006 | 天崩地裂 | 技能 | 地 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 |  |  |
| 3421007 | 大地震 | 技能 | 地 | 通用机制 | 已点名(4) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 | 状态容器、命中状态、基础范围 | 范围:方阵.晕眩1 |
| 3421008 | 联合施法 | 技能 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 强化其他法术时使其+1\攻 |
| 3421009 | 惧怖之颜 | 技能 | 地 | 通用机制 | 已点名(4) | 已操作(frontend-card-operation-report-skills.json) | 低 | 穿透、冷却、状态容器、命中状态 | 穿透.冷却1.使1个敌人石化2 |
| 3421010 | 大地穿刺 | 技能 | 地 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 穿透 | 穿透 |
| 3421011 | 自然生长 | 技能 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 选择你的1个负载小于4的地脉伙伴,它在本回合结束时获得负载+1\地 |
| 3421012 | 石破天惊 | 技能 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 | 穿透 | 穿透.你的伙伴每负载1点\地获得+1\威 |
| 3421013 | 大地共鸣 | 技能 | 地 | 已实现 | 已点名(4) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 | 冷却、基础范围 | 冷却1.范围:全场.你的场上每有1个负载或生命大于3的伙伴,获得+1\攻 |
| 3421014 | 千里流沙 | 技能 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 | 冷却、基础范围 | 范围:方阵.冷却1.若本卡攻击未命中,无需冷却 |
| 3421015 | 急袭沙暴 | 技能 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 | 速攻、冷却 | 冷却2.速攻.异能:所有原始威力小于5的法术-2\攻-2\威(最低为0) |
| 3501001 | 团结的希望 | 技能 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 从卡组上方开始将翻开5张牌,检索其中1张光辉伙伴,之后重洗卡组 |
| 3521001 | 治疗术 | 技能 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 回复2\血 |
| 3521002 | 神圣之火 | 技能 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 对友方单位不造成伤害,改为移除所有负面状态 |
| 3521003 | 神圣防护罩 | 技能 | 光 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制 | 防御 |
| 3521004 | 闪光魔术 | 技能 | 光 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器、命中状态 | 晕眩1 |
| 3521005 | 星陨 | 技能 | 光 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 | 穿透 | 穿透 |
| 3521006 | 光辉斩裂 | 技能 | 光 | 白板/无效果 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 |  |  |
| 3521007 | 希望呼唤 | 技能 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 从卡组上方开始将翻到的第1张光辉伙伴抽取,之后重洗卡组 |
| 3521008 | 光辉波动 | 技能 | 光 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器、命中状态、基础范围 | 范围:前排.晕眩1 |
| 3521009 | 幻彩流光 | 技能 | 光 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 |  |  |
| 3521011 | 光之庇护 | 技能 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 | 速攻、冷却 | 速攻.冷却2.选择1个伙伴,直到下个回合结束防止所有致命伤害 |
| 3521012 | 长虹贯日 | 技能 | 光 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 基础范围 | 范围:纵列 |
| 3521013 | 月之辉 | 技能 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制 | 用于防御或强化防御时+2威力 |
| 3521014 | 光之祝福 | 技能 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却1.使1个友方伙伴获得+1\血和负载+1\光 |
| 3521015 | 寂灭之光 | 技能 | 光 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-retry-nonpass.json) | 低 |  |  |
| 3621001 | 暗影冲击 | 技能 | 暗 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 |  |  |
| 3621002 | 噬血 | 技能 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 命中:使1个友方单位+2\血 |
| 3621003 | 死亡收割 | 技能 | 暗 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 基础范围 | 范围:前排 |
| 3621004 | 暗影箭 | 技能 | 暗 | 通用机制 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 | 穿透 | 穿透 |
| 3621005 | 暗冥弹 | 技能 | 暗 | 白板/无效果 | 通用/白板 | 已操作(frontend-card-operation-report-skills.json) | 低 |  |  |
| 3621006 | 死魂之噬 | 技能 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 状态容器 | 命中:将3层虚弱分配给敌方法术 |
| 3621007 | 安迪斯的惩罚 | 技能 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 每当友方单位受到1点伤害,下一次此技能获得+1\威 |
| 3621008 | 亡者之怒 | 技能 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 |  | 每当1个伙伴死亡,此法术获得+1\威 |
| 3621009 | 虚弱诅咒 | 技能 | 暗 | 通用机制 | 已点名(4) | 已操作(frontend-card-operation-report-skills.json) | 低 | 速攻、状态容器、命中状态 | 速攻.使1个敌方法术虚弱2 |
| 3621010 | 血魔爆 | 技能 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却1.献祭你的1个伙伴:对前排所有敌人造成该伙伴生命值的伤害 |
| 3621011 | 次元爆诞 | 技能 | 暗 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 穿透、基础范围 | 穿透.范围:方阵 |
| 3621012 | 回魂术 | 技能 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却1.从你的弃牌堆将最多2个伙伴移回手牌 |
| 3621013 | 亡灵护壁 | 技能 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制 | 防御.如果本回合或上个回合有友方单位死亡,此法术+2\威 |
| 3621014 | 业障 | 技能 | 暗 | 通用机制 | 已点名(3) | 已操作(frontend-card-operation-report-skills.json) | 低 | 防御限制、状态容器、命中状态 | 防御.若防御成功,使敌方攻击法术虚弱2 |
| 3621015 | 虹吸 | 技能 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-skills.json) | 低 | 冷却 | 冷却2.当敌方法术命中时使用,将造成的伤害改为回复生命值 |
| 4011001 | "南境百灵" 斯卡尔蒂 罗佳 | 人物 | 无 | 已实现 | 已点名(3) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 回合技:丢弃1张手牌,获得2点该卡牌属性种类的元素,对于奥术以外的属性这个效果只能使用1次 |
| 4011002 | "无面" | 人物 | 无 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 每当你让1张与你场上原有卡牌属性相同的卡牌入场,你受到1点伤害 |
| 4111001 | 掌门 龙卷火 | 人物 | 火 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 在游戏开始前,将万火合一术置于你的技能池 |
| 4111002 | 女巫 维兰德 | 人物 | 火 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-heroes.json) | 低 | 状态容器、命中状态 | 回合技:你获得点燃1,然后仅在本回合,将此卡负载的1\火变为1\无 |
| 4111003 | 大祭司 梵天 | 人物 | 火 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 绝技:本回合每当你的火焰法术命中,此卡获得负载+1\火 |
| 4211001 | "浪之人" 巴特尔 | 人物 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 绝技:本回合你将要打出的下一张手牌属性视为水,入场花费和负载的元素全部变为等量的\水 |
| 4211002 | 大贤者 沃尔波特 | 人物 | 水 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 在游戏开始前,将百川归海置于你的技能 |
| 4211003 | 凛冬城主 水晶心 | 人物 | 水 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-heroes.json) | 低 | 状态容器、命中状态 | 绝技:在本回合剩余时间内,你技能区内的法术获得冻结1 |
| 4311001 | 雷术士 肃 | 人物 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 绝技:丢弃2张大气手牌,对任意1名敌人造成1点伤害 |
| 4311002 | "渡鸦" 睿文 | 人物 | 气 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 你的起始手牌数与换牌机会+1 |
| 4311003 | 掌门 穆伶 | 人物 | 气 | 已实现 | 已点名(10) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 绝技:选择法力范围内双方各1个伙伴,花费它们入场费用差值的\气,将它们移回手牌 |
| 4411001 | 森林隐士 白须 | 人物 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 你可以用检索1张地属性野兽、植物或精灵来代替你首个回合的回合开始抽牌 |
| 4411002 | 大法师 安德鲁 | 人物 | 地 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 在游戏开始前,将灵兽 辛柯置于你的牌组 |
| 4411003 | 麦吉教授 | 人物 | 地 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 你的打出或学习的第一张原始费用大于5的卡牌费用减少2\地 |
| 4511001 | 圣使 玛丽斯 南森埃尔 | 人物 | 光 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-heroes.json) | 低 | 获得元素 | 绝技:当敌方造成伤害时,直到你的下个回合结束,每当你的单位受到对方伤害,获得2\光 |
| 4511002 | 神之眷子 爱里默 | 人物 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 在游戏开始前:将5张桎梏置于对手的牌组,当全部被解除(进入弃牌堆)时获得绝技:移除场上最多3张人物牌以外的卡牌 |
| 4511003 | 骑士团长 蕾曦娅 | 人物 | 光 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 在游戏开始前,用"团结的希望"替换你的技能池中的"希望呼唤" |
| 4611001 | 暗影学者 爱莉斯 | 人物 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 回合技:当1个你的伙伴死亡,使你的1个法术+1\威 |
| 4611002 | 芙雅夫人 | 人物 | 暗 | 已实现 | 已点名(2) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 绝技:使你的1个伙伴攻击和负载翻倍,但会在回合结束时死亡 |
| 4611003 | 咒言师 结影 | 人物 | 暗 | 已实现 | 已点名(1) | 已操作(frontend-card-operation-report-heroes.json) | 低 |  | 游戏开始前,将三张咒言书洗入你的卡组 |
