# 奥术纪元 - 线上版系统设计文档

## 1. 项目概述

将奥术纪元 (Era of Arcane) 卡牌游戏制作为线上对战版本。玩家通过浏览器进入游戏，导入卡组代码，匹配对手进行实时对战。所有卡牌对所有玩家免费开放。

### 技术栈
- **前端**: HTML/CSS/JS + Vue 3 (通过CDN引入，无需构建工具)
- **后端**: Go (Gin框架 + Gorilla WebSocket)
- **数据库**: PostgreSQL (用户/卡组/对局记录)
- **通信**: WebSocket (游戏内实时通信) + REST API (登录/卡组/匹配)
- **部署**: Docker容器，单实例部署至AWS EC2等（Phase 0，支持多房间并发）

### 为什么选择这些技术
- **Vue 3 CDN**: 无需Node构建链，单HTML文件即可运行，轻量且响应式
- **Go**: 高并发WebSocket处理性能优秀，编译为单二进制文件部署简单
- **PostgreSQL**: 云服务商通用支持 (AWS RDS, GCP Cloud SQL)
- **无Redis**: 匹配队列、会话状态、游戏状态全部放在Go进程内存中管理。单服务器部署场景下无需跨进程共享状态，减少组件和成本。如果未来需要多实例水平扩展再考虑引入

---

## 2. 系统架构

```
┌─────────────────────────────────────────────┐
│                   客户端 (浏览器)              │
│  ┌──────────┐  ┌───────────┐  ┌───────────┐ │
│  │  大厅界面  │  │  组卡/导入  │  │  对战界面  │ │
│  └────┬─────┘  └─────┬─────┘  └─────┬─────┘ │
│       │ REST          │ REST         │ WS     │
└───────┼───────────────┼──────────────┼───────┘
        │               │              │
┌───────▼───────────────▼──────────────▼───────┐
│                  Go 后端服务                    │
│  ┌──────────┐  ┌───────────┐  ┌───────────┐  │
│  │ REST API  │  │ 匹配服务   │  │ 游戏引擎   │  │
│  │ (Gin)    │  │(Matchmaker)│  │(GameEngine)│  │
│  └────┬─────┘  └─────┬─────┘  └─────┬─────┘  │
│       │               │              │         │
│  ┌────▼───────────────▼──────────────▼─────┐  │
│  │        GameManager (Go内存)              │  │
│  │  匹配队列 / 会话映射 / 游戏状态管理        │  │
│  └────┬────────────────────────────────────┘  │
└───────┼───────────────────────────────────────┘
        │
   ┌────▼─────┐
   │PostgreSQL │
   │用户/卡组DB │
   └──────────┘
```

---

## 3. 卡牌数据模型

### 3.1 卡牌编号规则 (7位数字)

基于对 `all_card_infos.json` 的分析，343张卡的编号结构：

```
位置:  1    2    3    4    5    6    7
含义: [类型][?] [?]  [?]  [?]  [版本相关] [序号]
```

编号首位与类型对应关系（从数据观察）：
- `1xxxxxx` → 伙伴 (Companion/生物)
- `2xxxxxx` → 道具 (Item)
- `3xxxxxx` → 技能 (Skill)
- `4xxxxxx` → 人物 (Hero)

### 3.2 卡牌JSON结构

```json
{
  "number": "4011001",        // 卡牌编号 (唯一ID)
  "type": "人物",             // 类型: 人物/伙伴/技能/道具
  "name": "卡牌名称",
  "category": "无",           // 分类/元素属性
  "tag": "",                  // 标签 (如"装备-饰物")
  "description": "效果描述",
  "quote": "风味文本",
  "elements_cost": {},        // 入场/使用费用 {"火焰": 2, "地脉": 1}
  "elements_gain": {},        // 负载/元素产出
  "elements_expense": {},     // 其他费用
  "version_number": "1",      // 版本号
  "version_name": "基础包",   // 所属扩展包
  "attack": -1,               // 攻击力 (-1表示无)
  "life": 6,                  // 生命值 (-1表示无)
  "duration": -1,             // 持续时间 (-1表示无)
  "power": -1,                // 威力 (-1表示无)
  "spawns": [],               // 衍生卡牌编号列表
  "output_path": "output/..." // 卡图路径
}
```

### 3.3 卡组代码格式

```
英雄编号 // 主卡组(30张,空格分隔) // 技能池(10张) // 额外卡组(衍生牌)
```

解析示例：
```
4311003 // 1021001 1021001 ... // 3221209 3311201 ... // 2001201 2001202 ...
```

---

## 4. 核心后端设计：游戏引擎

### 4.1 架构思路：事件驱动 + 状态机

参考成熟TCG游戏引擎（炉石传说、万智牌Arena）的设计理念，采用**事件驱动架构**。这是TCG游戏的业界标准做法，原因是：

1. 卡牌效果本质上是"当X发生时，做Y"——天然适合事件模型
2. 关键词如"入场"、"遗言"、"祈咒"都是对特定事件的响应
3. 法术防御链、强化等需要嵌套事件处理
4. 效果之间的交互（如石化使效果无效）需要统一的拦截机制

```
┌─────────────────────────────────────────┐
│              GameEngine                  │
│                                          │
│  ┌──────────┐    ┌───────────────────┐  │
│  │ TurnFSM  │    │   EventBus        │  │
│  │(回合状态机)│    │  (事件总线)        │  │
│  │          │    │                   │  │
│  │ StartPhase│───►│ emit(Event)       │  │
│  │ MainPhase│    │   │               │  │
│  │ EndPhase │    │   ▼               │  │
│  │ Wait*    │    │ EventQueue        │  │
│  └──────────┘    │   │               │  │
│                  │   ▼               │  │
│                  │ EffectResolver    │  │
│                  │ (效果结算器)       │  │
│                  └───────────────────┘  │
│                                          │
│  ┌──────────────────────────────────┐   │
│  │         GameState (游戏状态)       │   │
│  │  Player1State / Player2State     │   │
│  │  - Hero, Units[3][3], Skills[5]  │   │
│  │  - Equipment[5], Hand[], Deck[]  │   │
│  │  - SkillPool[], Graveyard[]      │   │
│  │  - Elements (当前可用元素)         │   │
│  └──────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

### 4.2 回合状态机 (TurnFSM)

```
游戏开始
  │
  ▼
┌──────────────┐
│  GAME_START  │  双方抽4张牌，可重抽一次
└──────┬───────┘
       │ 双方确认
       ▼
┌──────────────┐
│ TURN_START   │  当前玩家抽1张牌，触发"祈咒"效果
└──────┬───────┘  重置所有卡牌(竖置)，恢复元素
       │
       ▼
┌──────────────┐
│ MAIN_PHASE   │  玩家执行操作（召唤/装备/施法/攻击/学习技能）
│              │  每次操作后检查胜负
└──┬───┬───────┘
   │   │
   │   │ 玩家选择结束回合
   │   ▼
   │ ┌──────────────┐
   │ │  END_PHASE   │  弃牌至上限，结算标记(点燃/冻结等)
   │ └──────┬───────┘  临时效果消失
   │        │
   │        ▼
   │  切换当前玩家 → 回到 TURN_START
   │
   │ 法术/攻击发生
   ▼
┌──────────────┐
│ COMBAT_PHASE │  法术战斗子流程（见4.3）
│ (嵌套状态)    │
└──────────────┘
```

**等待状态 (Wait States)**：当需要玩家做出选择时（如防御、选择目标），状态机进入等待状态，通过WebSocket等待玩家响应。设置超时机制（如90秒）。

### 4.3 法术战斗流程

```
进攻方释放法术
  │
  ▼
┌─────────────────┐
│ SPELL_DECLARED   │  进攻方选择法术+目标，支付费用
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ DEFENSE_WINDOW   │  防御方可选择：
│                  │  1. 不防御 → 法术命中
│                  │  2. 防御（施放防御法术，可强化）
│                  │  3. 透支获取额外费用
└────┬───────┬────┘
     │       │
 不防御    防御
     │       │
     ▼       ▼
┌────────┐ ┌──────────────┐
│命中结算 │ │ 比较威力      │
└────────┘ │ 防御方>=进攻方 │
           │  → 防御成功   │
           │ 防御方<进攻方  │
           │  → 法术命中   │
           └──────────────┘
```

### 4.4 事件系统

```go
// 事件类型枚举
type EventType int
const (
    EventGameStart EventType = iota
    EventTurnStart
    EventTurnEnd
    EventDraw              // 抽牌
    EventSummon            // 召唤生物
    EventEnterField        // 入场（触发入场效果）
    EventLeaveField        // 离场
    EventDeath             // 死亡（触发遗言）
    EventSpellCast         // 施放法术
    EventSpellHit          // 法术命中
    EventSpellDefended     // 法术被防御
    EventAttack            // 单位攻击
    EventAttackHit         // 攻击命中
    EventDamage            // 造成伤害
    EventHeal              // 治疗
    EventEquip             // 装备道具
    EventLearnSkill        // 学习技能
    EventConsume           // 消耗（横置获取元素）
    EventSacrifice         // 献祭
    EventStatusApply       // 施加状态标记
    EventStatusRemove      // 移除状态标记
    EventElementGain       // 获取元素
    EventElementSpend      // 花费元素
)

// 事件结构
type GameEvent struct {
    Type       EventType
    Source     CardInstance   // 事件来源卡牌
    Target     CardInstance   // 目标卡牌/区域
    Player     PlayerID       // 触发玩家
    Data       map[string]any // 额外数据（伤害值、元素类型等）
    Cancelled  bool           // 是否被取消
}

// 效果监听器（每张卡牌的效果注册为监听器）
type EffectListener struct {
    CardID    string
    EventType EventType
    Priority  int            // 结算优先级
    Condition func(*GameState, *GameEvent) bool
    Execute   func(*GameState, *GameEvent) []GameEvent // 可产生新事件
}
```

### 4.5 事件处理流程

```
emit(Event)
    │
    ▼
EventQueue.Push(event)
    │
    ▼
while EventQueue not empty:
    │
    event = EventQueue.Pop()
    │
    ▼
    收集所有监听该事件的 EffectListener
    │
    ▼
    按优先级自动排序（见下方排序规则）
    │
    ▼
    依次执行 Listener:
      - 检查 Condition（如石化则跳过）
      - 执行 Execute
      - 产生的新事件加入 EventQueue
    │
    ▼
    检查胜负条件
```

### 4.6 同时事件的自动排序规则

线下规则中，同时发生的事件由当前回合玩家决定顺序。线上版为避免频繁打断流程等待玩家选择，改为**自动确定性排序**，规则如下：

```
同时触发的多个效果，依次按以下维度排序（优先级从高到低）：

1. 效果类别优先级（高→低）：
   - 死亡/离场结算 (如遗言)
   - 状态变化 (施加/移除标记)
   - 伤害/治疗
   - 召唤/入场
   - 抽牌/检索
   - 元素变化
   - 其他

2. 归属方：当前回合玩家的效果先于对手

3. 场上位置：按单位区从左到右、从前排到后排的固定顺序
   位置编号: (0,0) → (1,0) → (2,0) → (0,1) → (1,1) → (2,1) → (0,2) → (1,2) → (2,2)
   技能区/装备区按槽位编号 0→4

4. 卡牌入场时间：先入场的卡牌先结算
```

此规则完全确定，不需要玩家交互，且大多数情况下符合直觉（先死亡再入场、先伤害再治疗、回合方优先）。

这种设计确保：
- "入场"效果在 `EventEnterField` 时自动触发
- "遗言"效果在 `EventDeath` 时自动触发
- "祈咒"效果在 `EventTurnStart` 时自动触发
- "石化"通过在 Condition 检查中拦截来使效果无效
- "精通"通过在 `EventConsume` 上注册监听来叠加层数

### 4.6 卡牌效果的实现方式

每张卡的 `description` 字段描述了效果。由于卡牌数量有限（343张），我们对每张卡的效果进行**手动编码**而非脚本解析（更可靠、更可调试）。

```go
// 卡牌效果注册表
var CardEffectRegistry = map[string]CardEffectFactory{
    "4011001": NewHero4011001,  // 斯卡尔蒂
    "1321002": NewCompanion1321002,
    // ... 每张卡一个工厂函数
}

// 示例：一张有入场效果的伙伴
func NewCompanion1321002(instance *CardInstance) []EffectListener {
    return []EffectListener{
        {
            EventType: EventEnterField,
            Condition: func(gs *GameState, e *GameEvent) bool {
                return e.Target.InstanceID == instance.InstanceID
            },
            Execute: func(gs *GameState, e *GameEvent) []GameEvent {
                // 入场效果逻辑
                return nil
            },
        },
    }
}
```

**分阶段实现策略**：先实现无效果的"白板"卡牌（仅有数值），让游戏核心流程跑通，再逐张添加卡牌效果。

---

## 5. 前端设计

### 5.1 页面结构

```
/                → 首页/大厅
/game/:roomId    → 对战房间
```

### 5.2 大厅界面

```
┌─────────────────────────────────────────────┐
│  奥术纪元 Era of Arcane         [玩家名称]    │
├─────────────────────────────────────────────┤
│                                             │
│  ┌─────────────────────────────────┐        │
│  │  卡组代码                        │        │
│  │  ┌─────────────────────────┐    │        │
│  │  │  粘贴卡组代码...          │    │        │
│  │  └─────────────────────────┘    │        │
│  │  [导入卡组]                      │        │
│  │                                 │        │
│  │  英雄: 斯卡尔蒂 (4011001)       │        │
│  │  主卡组: 30/30 ✓                │        │
│  │  技能池: 10/10 ✓                │        │
│  └─────────────────────────────────┘        │
│                                             │
│  ┌──────────────┐  ┌──────────────┐         │
│  │  快速匹配     │  │  创建房间     │         │
│  └──────────────┘  └──────────────┘         │
│                                             │
│  ┌──────────────┐                           │
│  │  房间号: ___  │  [加入房间]               │
│  └──────────────┘                           │
└─────────────────────────────────────────────┘
```

### 5.3 对战界面布局

参照原始设计图（rules页面image_1.png），每位玩家的场地布局如下：

```
单个玩家场地（原始设计）:
┌──────────────────────────────────────────────────┐
│                                                  │
│  元素    ┌──────┐  ┌─────┬─────┬─────┐  ┌─────┐ │
│  计数    │      │  │     │     │     │  │技能池│ │
│          │ 主卡组│  │     │     │     │  │ ②  │ │
│  火 0    │  ①  │  ├─────┼─────┼─────┤  ├─────┤ │
│  水 0    │      │  │     │[英雄]│     │  │弃牌堆│ │
│  地 0    └──────┘  │     │     │     │  │ ③  │ │
│  气 0              ├─────┼─────┼─────┤  └─────┘ │
│  光 0              │     │     │     │          │
│  暗 0              │     │     │     │          │
│  奥 0              └─────┴─────┴─────┘          │
│                      3×3 单位区                   │
│                                                  │
│  ┌──────┐  ⑥装备区(≤5)    ④技能区(≤5)   ┌─────┐ │
│  │手牌⑦ │  [装][装][ ][ ][ ]  [技][技][ ][ ][ ] │ │
│  │      │                               │闲置⑤│ │
│  └──────┘                               └─────┘ │
└──────────────────────────────────────────────────┘
```

线上版对战界面（双方镜像排列）：

```
┌──────────────────────────────────────────────────────────┐
│  对手信息: HP:6  手牌:5  卡组:22                           │
├──────────────────────────────────────────────────────────┤
│                                                          │
│         [技][ ][ ][ ][ ]  [装][ ][ ][ ][ ]               │
│          对手技能区④        对手装备区⑥                    │
│                                                          │
│  元素     ┌─────┬─────┬─────┐                            │
│  火 0     │     │     │     │     [技能池②] [弃牌堆③]     │
│  水 0     ├─────┼─────┼─────┤                            │
│  地 0     │     │[英雄]│     │     卡组剩余: 22            │
│  气 0     ├─────┼─────┼─────┤                            │
│  光 0     │     │     │     │  ← 对手前排                 │
│  暗 0     └─────┴─────┴─────┘                            │
│  奥 0       对手单位区 (3×3)                               │
│                                                          │
│  ══════════════ 战场分界线 ══════════════                  │
│                                                          │
│  火 3       我方单位区 (3×3)                               │
│  水 2     ┌─────┬─────┬─────┐                            │
│  地 0     │     │     │     │  ← 我方前排                 │
│  气 0     ├─────┼─────┼─────┤                            │
│  光 0     │     │[英雄]│     │     [技能池②] [弃牌堆③]     │
│  暗 0     ├─────┼─────┼─────┤                            │
│  奥 1     │     │     │     │     卡组剩余: 25            │
│           └─────┴─────┴─────┘                            │
│                                                          │
│         [技][技][ ][ ][ ]  [装][装][ ][ ][ ]             │
│          我方技能区④         我方装备区⑥                   │
│                                                          │
├──────────────────────────────────────────────────────────┤
│  我方信息: HP:6  元素总览              [结束回合]           │
├──────────────────────────────────────────────────────────┤
│  手牌: [卡1] [卡2] [卡3] [卡4] [卡5] [卡6]               │
├──────────────────────────────────────────────────────────┤
│  操作日志                                                 │
└──────────────────────────────────────────────────────────┘
```

**交互设计要点：**
- 卡牌悬停显示大图+完整描述
- 横置卡牌视觉上旋转90度，灰色遮罩
- 可操作的卡牌边框高亮（绿色=可召唤，蓝色=可消耗，红色=可攻击）
- 法术释放时：选择法术 → 选择目标 → 确认，目标区域高亮提示
- 防御阶段：屏幕提示"对方施放法术，是否防御？"，倒计时显示
- 元素用彩色图标显示，不同元素不同颜色

### 5.4 卡牌图片

卡牌图片URL通过 `output_path` 字段拼接基础URL获得：
```
https://yifeeeeei.github.io/ArcaneImages/{output_path}
```

---

## 6. 通信协议

### 6.1 REST API

```
POST   /api/auth/guest          # 游客登录，返回临时token
POST   /api/deck/import         # 导入卡组代码，验证合法性
GET    /api/deck/:id            # 获取已保存卡组
POST   /api/match/quick         # 快速匹配
POST   /api/room/create         # 创建房间
POST   /api/room/:id/join       # 加入房间
GET    /api/cards                # 获取所有卡牌数据（前端缓存）
```

### 6.2 WebSocket 消息协议

客户端→服务器 (Action):
```json
{"action": "summon",       "card_id": "1321002", "position": [1, 0]}
{"action": "consume",      "instance_id": "uuid-xxx"}
{"action": "cast_spell",   "skill_id": "3221209", "target": {"type": "unit", "pos": [0, 1]}}
{"action": "defend",       "skill_ids": ["3321002"], "boost_ids": ["3321007"]}
{"action": "no_defend"}
{"action": "attack",       "attacker_id": "uuid-xxx", "target_pos": [0, 1]}
{"action": "equip",        "card_id": "2321006"}
{"action": "learn_skill",  "card_id": "3321206", "replace_id": null}
{"action": "use_item",     "card_id": "2221205", "target": {"type": "unit", "pos": [1, 1]}}
{"action": "end_turn"}
{"action": "mulligan",     "card_ids": ["1321002", "2321006"]}
```

服务器→客户端 (Event Stream):
```json
{"event": "game_start",    "your_side": "player1", "initial_hand": [...]}
{"event": "turn_start",    "player": "player1", "drawn_card": "1321013"}
{"event": "card_summoned",  "player": "player2", "card": {...}, "position": [2, 0]}
{"event": "spell_cast",    "player": "player1", "spell": {...}, "target": {...}}
{"event": "defend_window",  "timeout": 30}
{"event": "damage_dealt",  "source": "...", "target": "...", "amount": 3}
{"event": "status_applied", "target": "...", "status": "burn", "stacks": 2}
{"event": "card_died",     "card": {...}, "triggered_effects": [...]}
{"event": "game_over",     "winner": "player1", "reason": "hero_killed"}
{"event": "state_sync",    "full_state": {...}}
```

**信息可见性规则：**
- 对手手牌：只发送数量，不发送具体卡牌
- 对手卡组：只发送剩余数量
- 对手技能池：不可见（除非规则另有规定）
- 符文（盖放的反制道具）：只发送"有一张盖放的牌"
- 所有公开信息（场上卡牌、弃牌堆）：双方都完整可见

---

## 7. 数据库设计

### PostgreSQL

```sql
-- 玩家表（支持游客和注册用户）
CREATE TABLE players (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(32) NOT NULL,
    is_guest    BOOLEAN DEFAULT true,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- 卡牌数据表（从all_card_infos.json导入）
CREATE TABLE cards (
    number      VARCHAR(8) PRIMARY KEY,  -- 卡牌编号
    type        VARCHAR(16) NOT NULL,     -- 人物/伙伴/技能/道具
    name        VARCHAR(128) NOT NULL,
    category    VARCHAR(32),
    tag         VARCHAR(64),
    description TEXT,
    quote       TEXT,
    elem_cost   JSONB DEFAULT '{}',
    elem_gain   JSONB DEFAULT '{}',
    elem_expense JSONB DEFAULT '{}',
    version_num VARCHAR(4),
    version_name VARCHAR(32),
    attack      INT DEFAULT -1,
    life        INT DEFAULT -1,
    duration    INT DEFAULT -1,
    power       INT DEFAULT -1,
    spawns      JSONB DEFAULT '[]',
    image_path  VARCHAR(256)
);

-- 已保存卡组
CREATE TABLE decks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id   UUID REFERENCES players(id),
    name        VARCHAR(64),
    deck_code   TEXT NOT NULL,           -- 原始卡组代码
    hero_id     VARCHAR(8) NOT NULL,
    main_deck   VARCHAR(8)[] NOT NULL,   -- 30张
    skill_pool  VARCHAR(8)[] NOT NULL,   -- 10张
    extra_deck  VARCHAR(8)[],            -- 衍生牌
    created_at  TIMESTAMP DEFAULT NOW()
);

-- 对局记录
CREATE TABLE games (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player1_id  UUID REFERENCES players(id),
    player2_id  UUID REFERENCES players(id),
    winner_id   UUID REFERENCES players(id),
    p1_hero     VARCHAR(8),
    p2_hero     VARCHAR(8),
    turns       INT,
    started_at  TIMESTAMP DEFAULT NOW(),
    ended_at    TIMESTAMP
);
```

### Go内存状态（替代Redis）

```go
// 在GameManager中维护，无需外部存储
type GameManager struct {
    matchQueue  []MatchEntry              // 匹配队列
    sessions    map[PlayerID]*Session     // 玩家会话 (WebSocket连接、所在房间等)
    activeGames map[GameID]*GameInstance  // 进行中的游戏状态 (断线重连时从这里恢复)
    rooms       map[RoomID]*Room          // 房间列表
    mu          sync.RWMutex
}
```

每个房间的GameInstance在独立goroutine中运行，通过channel接收玩家操作，天然支持多房间并发。单实例部署，服务重启时进行中的游戏会丢失，当前规模可以接受。

---

## 8. 项目目录结构

```
EraOfArcaneGame/
├── DESIGN.md
├── docker-compose.yml
├── Dockerfile
│
├── server/                     # Go后端
│   ├── main.go
│   ├── go.mod
│   ├── config/
│   │   └── config.go           # 配置加载
│   ├── api/
│   │   ├── router.go           # 路由注册
│   │   ├── auth.go             # 认证API
│   │   ├── deck.go             # 卡组API
│   │   ├── room.go             # 房间API
│   │   └── ws.go               # WebSocket处理
│   ├── game/
│   │   ├── engine.go           # 游戏引擎主体
│   │   ├── state.go            # GameState定义
│   │   ├── turn.go             # 回合状态机
│   │   ├── event.go            # 事件系统
│   │   ├── combat.go           # 法术战斗/单位攻击
│   │   ├── effect.go           # 效果处理器
│   │   ├── element.go          # 元素系统
│   │   ├── status.go           # 状态标记系统
│   │   ├── board.go            # 场地/区域管理
│   │   ├── card_instance.go    # 卡牌实例
│   │   ├── validation.go       # 操作合法性验证
│   │   └── visibility.go       # 信息可见性过滤
│   ├── cards/
│   │   ├── registry.go         # 卡牌效果注册表
│   │   ├── loader.go           # 从JSON加载卡牌数据
│   │   ├── heroes/             # 英雄效果
│   │   ├── companions/         # 伙伴效果
│   │   ├── skills/             # 技能效果
│   │   └── items/              # 道具效果
│   ├── match/
│   │   ├── matchmaker.go       # 匹配逻辑
│   │   └── room.go             # 房间管理
│   ├── model/
│   │   ├── card.go             # 卡牌数据模型
│   │   ├── deck.go             # 卡组模型
│   │   ├── player.go           # 玩家模型
│   │   └── game_record.go      # 对局记录
│   └── store/
│       └── postgres.go         # PostgreSQL操作
│
├── web/                        # 前端
│   ├── index.html              # 大厅页面
│   ├── game.html               # 对战页面
│   ├── css/
│   │   ├── main.css            # 通用样式
│   │   ├── lobby.css           # 大厅样式
│   │   └── game.css            # 对战样式
│   ├── js/
│   │   ├── api.js              # REST API调用
│   │   ├── ws.js               # WebSocket管理
│   │   ├── lobby.js            # 大厅逻辑（Vue app）
│   │   ├── game.js             # 对战主逻辑（Vue app）
│   │   ├── board.js            # 棋盘渲染
│   │   ├── card.js             # 卡牌组件
│   │   ├── combat-ui.js        # 战斗交互UI
│   │   └── animations.js       # 动画效果
│   └── assets/
│       └── elements/           # 元素图标
│
└── data/
    └── all_card_infos.json     # 卡牌数据源
```

---

## 9. 开发计划（分阶段）

### Phase 0: 本地可运行测试版

**目标**：`go run .` 一条命令启动，浏览器打开即玩，无需安装任何外部依赖。

**实现方式**：
- 数据库使用 **SQLite**（`github.com/mattn/go-sqlite3`），单文件存储，零配置
- 前端静态文件由Go内嵌（`embed`包）或直接从`web/`目录serve
- 卡牌数据从 `all_card_infos.json` 启动时加载到内存
- 本地测试时同一台机器开两个浏览器标签页即可双人对战
- 通过环境变量或配置文件切换 SQLite / PostgreSQL：
  ```
  # 本地开发（默认）
  DB_DRIVER=sqlite
  DB_DSN=./data/arcane.db

  # 生产部署
  DB_DRIVER=postgres
  DB_DSN=postgres://user:pass@host:5432/arcane
  ```
- 使用 `database/sql` 标准接口 + 兼容SQL写法，同一套代码适配两种数据库

**本地启动流程**：
```bash
cd server
go run .
# 服务启动在 http://localhost:8080
# 浏览器打开两个标签页，各自导入卡组，创建/加入房间即可对战
```

### Phase 1: 基础骨架
- 搭建Go项目结构，初始化数据库
- 加载卡牌JSON数据
- 实现卡组代码导入/验证
- 前端大厅页面 + 卡组导入
- 房间创建/加入 + WebSocket连接

### Phase 2: 核心游戏引擎
- GameState 数据结构
- 回合状态机（开始/出牌/结束）
- 元素系统（负载→消耗→获取→花费）
- 基础操作：召唤生物、装备道具、学习技能
- 单位攻击（无法术效果的简单版本）

### Phase 3: 法术战斗系统
- 法术释放与目标选择
- 防御窗口与威力比较
- 法术强化机制
- 透支机制
- 咒术（不可防御）

### Phase 4: 事件系统与卡牌效果
- 事件总线实现
- 基础关键词：入场、遗言、速攻、临时
- 状态标记：点燃、冻结、眩晕、石化、虚弱
- 逐步实现各卡牌效果（按扩展包分批）

### Phase 5: 前端对战界面
- 完整棋盘渲染
- 卡牌交互（拖拽/点击放置）
- 法术战斗UI（防御选择、目标选择）
- 动画与视觉反馈
- 断线重连

### Phase 6: 完善与部署
- 匹配系统优化
- 对局回放
- Docker打包
- AWS部署配置

---

## 10. 关键设计决策备注

1. **所有卡牌免费开放**：无需抽卡/解锁系统，大幅简化用户系统
2. **手动编码卡牌效果**：343张卡数量可控，比实现DSL脚本引擎更可靠
3. **事件驱动架构**：TCG游戏的标准做法，天然支持"当X时触发Y"的卡牌逻辑
4. **服务端权威**：所有游戏逻辑在服务端执行，客户端只负责展示和输入，防作弊
5. **WebSocket + REST混合**：对战用WS实时通信，其余用REST，清晰分工
6. **断线重连**：Go内存保存游戏状态，玩家重连后发送 `state_sync` 恢复，无需额外组件
7. **Vue 3 CDN模式**：避免构建工具链，降低前端复杂度，专注游戏逻辑
