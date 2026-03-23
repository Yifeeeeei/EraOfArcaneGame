# 卡牌效果注册指南

本文档说明如何为奥术纪元游戏注册新的卡牌效果。

## 核心原则

1. **不使用自动解析** - 每张卡的效果都是显式代码编写的
2. **可复用效果构建器** - 使用 `effect_builders.go` 中的预定义函数
3. **显式注册** - 在 `card_registry.go` 中逐卡注册效果

## 文件结构

```
server/game/
├── effect_builders.go    # 效果构建器函数库
├── effect_system.go      # 效果系统框架
├── card_registry.go      # 卡牌显式注册（你在这里工作）
└── effect_cards.go       # 旧的手动效果（参考用）
```

## 快速开始

### 1. 找到要注册的卡牌

获取卡牌的：
- **卡牌编号** (如 `"1321301"`)
- **触发时机** (入场/遗言/回合开始等)
- **效果描述** (造成X伤害/抽X张牌等)

### 2. 选择合适的效果构建器

常用构建器：

| 构建器 | 效果 | 示例 |
|--------|------|------|
| `DrawCards(n)` | 抽N张牌 | `DrawCards(2)` |
| `DealDamageAuto(n)` | 对前排敌方造成N点伤害 | `DealDamageAuto(3)` |
| `DealDamageToTarget(n)` | 对指定目标造成N点伤害 | `DealDamageToTarget(2)` |
| `GainShield(n)` | 获得N点护盾 | `GainShield(3)` |
| `GainStealth(n)` | 获得N层隐蔽 | `GainStealth(2)` |
| `ApplyStatusAuto(status, n)` | 对前排敌方施加N层状态 | `ApplyStatusAuto("虚弱", 2)` |
| `ApplyStatusToTarget(status, n)` | 对目标施加N层状态 | `ApplyStatusToTarget("冻结", 1)` |
| `HealSelf(n)` | 治疗自身N点 | `HealSelf(3)` |
| `HealTarget(n)` | 治疗目标N点 | `HealTarget(2)` |
| `ModifySelfAttack(n)` | 自身攻击力+N | `ModifySelfAttack(2)` |
| `ModifyTargetAttack(n)` | 目标攻击力+N | `ModifyTargetAttack(-1)` |
| `GainCharge(n)` | 获得N点充能 | `GainCharge(1)` |
| `Combine(handlers...)` | 组合多个效果 | `Combine(DrawCards(1), DealDamageAuto(1))` |

### 3. 在 card_registry.go 中注册

打开 `server/game/card_registry.go`，在 `RegisterAllCardEffects` 函数中添加注册：

```go
func RegisterAllCardEffects(r *EffectRegistry) {
	// 已有注册...

	// ===== 新卡牌注册 =====
	// 卡牌编号: 1321301
	// 效果: 入场时抽1张牌
	r.Register("1321301", TriggerOnEnter, DrawCards(1))

	// 卡牌编号: 1321302
	// 效果: 入场时对前排敌方造成2点伤害
	r.Register("1321302", TriggerOnEnter, DealDamageAuto(2))

	// 卡牌编号: 1321303
	// 效果: 遗言：治疗己方英雄3点生命
	r.Register("1321303", TriggerOnDeath, HealHero(3))

	// 更多卡牌...
}
```

## 触发器类型

| 触发器 | 说明 |
|--------|------|
| `TriggerOnEnter` | 单位/装备入场时触发 |
| `TriggerOnDeath` | 单位死亡时触发（遗言） |
| `TriggerOnTurnStart` | 回合开始时触发 |
| `TriggerOnTurnEnd` | 回合结束时触发 |
| `TriggerOnAttack` | 攻击时触发 |
| `TriggerOnDefend` | 防御时触发 |
| `TriggerOnDamaged` | 受到伤害时触发 |
| `TriggerOnBoosted` | 被强化时触发 |

## 高级用法

### 组合效果

使用 `Combine()` 组合多个效果：

```go
// 抽1张牌，然后对前排造成1点伤害
r.Register("1321304", TriggerOnEnter, Combine(
	DrawCards(1),
	DealDamageAuto(1),
))

// 获得2点护盾，然后攻击力+1
r.Register("1321305", TriggerOnEnter, Combine(
	GainShield(2),
	ModifySelfAttack(1),
))
```

### 条件效果

对于复杂条件，可以编写自定义效果：

```go
// 如果己方英雄生命低于10，抽2张牌
r.Register("1321306", TriggerOnEnter, func(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps.Hero != nil && ps.Hero.CurrentLife < 10 {
		return DrawCards(2)(ctx)
	}
	return nil
})
```

## 最佳实践

1. **一张卡注册一次** - 每张卡的效果只用一行代码注册
2. **使用构建器优先** - 优先使用预定义构建器，保持代码简洁
3. **添加注释** - 为每张卡添加注释说明卡牌编号和效果
4. **分组管理** - 按卡牌类型或扩展包分组管理注册代码
5. **测试验证** - 注册新卡后启动游戏测试效果是否正常

## 常见问题

### Q: 如何注册需要选择目标的效果？
A: 目前自动构建器使用自动目标（前排敌方）。如果需要精确目标选择，需要编写自定义效果处理函数，使用 PendingAction 系统。

### Q: 如何注册有消耗的效果（如弃牌）？
A: 使用 `Combine(SacrificeSelfAndDo(...), ...)` 或 `DiscardRandom(n)` 等构建器。

### Q: 如果构建器不能满足需求怎么办？
A: 可以直接编写 EffectHandler 函数，参考 `effect_cards.go` 中的旧实现方式。

---

如有问题，请查阅 `effect_builders.go` 中的完整构建器列表，或参考已有的注册示例。