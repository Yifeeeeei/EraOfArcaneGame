# 奥术纪元 (Era of Arcane) — 开发进度与思路

## 项目概述
将巫师牌桌游制作成线上对战版本。Go后端 + Vue 3 CDN前端，WebSocket实时通信，无构建工具。

## 技术架构
- **后端**: Go (net/http + gorilla/websocket)，无框架
- **前端**: Vue 3 CDN + 纯CSS，内联在 `web/game.html`
- **卡牌数据**: 1,278张卡，从 `all_card_infos.json` 加载
- **卡牌图片**: `https://yifeeeeei.github.io/ArcaneImages/{output_path}`
- **卡牌数据API**: `https://yifeeeeei.github.io/ArcaneImages/output/all_card_infos.json`

## 关键文件
| 文件 | 职责 |
|------|------|
| `server/game/engine.go` | 核心游戏引擎，13个action handler，状态同步 |
| `server/game/state.go` | 游戏状态结构体，卡牌实例，元素支付逻辑 |
| `server/game/effect_system.go` | 效果注册表，18种触发器，效果执行框架 |
| `server/game/effect_parser.go` | 自动解析卡牌描述文本注册效果（11种模式） |
| `server/game/effect_cards.go` | 手动注册的70张卡效果（48英雄+14伙伴+5技能+3道具） |
| `server/game/effect_keywords.go` | 19个关键词实现（速攻/穿透/隐蔽/护盾等） |
| `web/game.html` | 完整前端（Vue 3 + 内联JS），所有游戏交互UI |
| `web/css/game.css` | 暗幻风格主题，弹性布局适配屏幕 |
| `server/api/ws.go` | WebSocket处理，消息路由 |
| `server/api/router.go` | HTTP路由，静态文件服务 |
| `server/match/room.go` | 房间管理，玩家匹配 |

## 设计决策
1. **不做EventQueue重构** — 保持同步执行模型，用PendingAction处理玩家选择
2. **边缘情况自由决定** — 保持逻辑模块化，方便后续调整
3. **自动解析器优先** — 尽量通过描述文本自动注册效果，减少手动编码

## 已完成工作

### 基础引擎 (完成)
- 完整回合流程：抽牌→出牌→结束回合→重置→状态结算
- 元素系统：7种元素 + 奥术万能替代
- 消耗机制：横置产出元素
- 召唤/装备/学习技能
- 法术战斗：施放→防御窗口→命中/防御成功
- 单位攻击：前排限制 + 法力范围
- 生命/伤害/死亡/遗言链
- Mulligan（重抽）
- 弃牌上限（自动弃最后一张，待改为玩家选择）

### 效果系统 (部分完成)
- 18种触发器类型
- 70张手动效果（大部分英雄）
- 自动解析器覆盖11种简单模式
- 19个关键词实现
- ~60-65%卡牌可正常运作（白板+关键词+简单效果）

### 前端 (部分完成)
- 暗幻主题，弹性布局适配屏幕
- 3x3单位格子 + 技能/装备栏
- 手牌栏 + 技能池
- 召唤/消耗/攻击/施法交互
- 卡牌详情弹窗
- Mulligan界面
- 操作日志

---

## 已完成的推进工作 (Session 2, 2026-03-22)

### Phase 1: 核心可玩性修复 ✅

#### 1A. 防御UI ✅
- 防御窗口展示可用防御技能（竖置、非冷却、非咒术）
- 支持选择多个防御技能 + 强化技能
- 实时显示威力对比和防御成功/失败预判
- 发送 `defend` action（skill_ids + boost_ids）

#### 1B. 英雄/单位技能UI ✅
- `cardToInfo()` 添加 `has_per_turn`/`has_ultimate`/`used_this_turn`/`ultimate_used`/`per_turn_limit`
- 添加 `is_defense_only`/`is_sorcery` 标记
- 添加 `elements_expense` 到卡牌序列化
- 英雄肖像旁添加回合技/绝技按钮
- 单位格子上添加技能按钮（检测眩晕/石化状态阻止使用）

#### 1C. 消耗品使用UI ✅
- 区分装备道具（equip）和消耗品（use_item），通过tag检测
- 手牌选择消耗品时显示"使用"按钮
- 技能栏显示防御/咒术标记badge

### Phase 2: 增强自动解析器 ✅
- 移除early return，允许一张卡注册多个效果
- 入场新模式: 护盾X, 隐蔽(X), 虚弱X, 自动目标front row
- 遗言新模式: 冻结X, 点燃X (对前排敌方)
- 新增祈咒解析: 抽牌/充能/伤害/点燃/冻结
- `registerStatusOnEnter` 支持无目标时自动选前排
- 新增 `findFrontRowUnit()` 辅助函数

### Phase 4B: 法术强化 ✅
- `handleCastSpell` 接受 `boost_ids` 参数
- 强化技能的威力加到主技能上
- 前端: 选定技能后显示"强化"按钮，可选多技能强化
- 前端: boost-selected 样式高亮已选强化技能

---

### Phase 3: PendingAction目标选择系统 ✅
- `PendingAction` 结构: type/playerID/prompt/candidates/min/max/callback
- `PhaseWaitingAction` 新阶段
- `SetPendingAction()` 暂停游戏等待玩家选择
- `handleResolveAction()` 处理玩家选择并回调
- `GetStateForPlayer()` 包含 pending_action 数据
- 前端: 通用选择弹窗（显示候选卡牌，支持单选/多选）
- **弃牌选择**: endTurn超手牌上限时弹出选择弹窗（不再自动弃最后一张）

### 前端增强 ✅
- 扩充事件日志: effect_trigger/ability_used/defense_attempt/pending_action/discard/terrain等
- 技能栏: 防御/咒术badge标记，强化按钮
- 防御专用技能在主动施放时隐藏

---

## 进行中工作

### Phase 2: 增强自动解析器
新增解析模式（预计覆盖额外200+张卡）：
- 护盾X、隐蔽X
- 检索1张XXX
- 负载修改
- 对敌方施加状态
- 祈咒效果
- 防御标记
- 范围标记（前排/纵列/方阵/溅射/全场）
- **修复**: 移除early return，允许一张卡注册多个效果

### Phase 3: 目标选择系统
- 后端 PendingAction 结构（select_target/select_card/discard）
- 前端通用选择弹窗
- 接入检索/献祭/弃牌等效果

### Phase 4: 法术战斗完善
- 4A: 区域目标（front_row/column/splash/all）
- 4B: 法术强化（handleCastSpell接受boost_ids）
- 4C: TriggerOnDefend

### Phase 5: 高级机制
- 反制（43张卡）、仪式（19张卡）、法宝（35张卡）、地形效果

### Phase 6: 收尾
- 剩余~100张复杂卡牌手动效果
- 英雄开局效果
- 费用修改器
- 弃牌选择

## 验证方式
```bash
# Go测试
cd server && go test ./game/ -v

# WebSocket集成测试
cd server && go test ./api/ -v

# 启动服务器
cd server && go run .

# 访问
open http://localhost:8080
```
