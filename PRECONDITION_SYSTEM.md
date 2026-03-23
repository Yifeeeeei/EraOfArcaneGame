# 前提条件检查系统 (Precondition Check System)

## 概述

这是一个通用的、可扩展的前提条件检查机制，用于在使用卡牌效果前验证各种条件是否满足。

## 设计原则

1. **可组合**: 多个条件可以组合使用 (AllOf, AnyOf, Not)
2. **可扩展**: 新条件类型可以方便地添加
3. **可复用**: 常用条件作为预定义函数提供
4. **错误信息**: 检查失败时提供清晰的错误信息

## 核心类型

### Precondition 接口
```go
type Precondition interface {
    Check(ctx *ConditionContext) PreconditionResult
    String() string
}
```

### ConditionContext 上下文
```go
type ConditionContext struct {
    Engine     *Engine
    PlayerID   int
    SourceCard *CardInstance
    ExtraData  map[string]any
}
```

## 组合条件

### AllOf - 所有条件都满足
```go
AllOf{
    Conditions: []Precondition{
        HasFieldUnits(-2),      // 对手有场上单位
        HasCharge(2),           // 自己有2点充能
    },
}
```

### AnyOf - 至少一个条件满足
```go
AnyOf{
    Conditions: []Precondition{
        HasFieldUnits(-1),      // 自己有场上单位
        HasHandCards(-1, 3),    // 或者有3张以上手牌
    },
}
```

### Not - 条件不满足
```go
Not{
    Condition: HasStatus("冰冻", 0),  // 没有被冰冻
}
```

## 具体条件类型

### 场上单位检查
```go
HasFieldUnits(playerID)     // playerID: -1=自己, -2=对手, >=0=指定玩家
```

### 手牌检查
```go
HasHandCards(playerID, minCount)
```

### 充能检查
```go
HasCharge(minCharge)
```

### 状态检查
```go
SourceHasStatus(status, minAmount)  // 来源卡牌有状态
TargetHasStatus(status, minAmount)  // 目标卡牌有状态
```

## 快捷函数

提供了常用条件的快捷创建函数：

```go
// 检查场上单位
HasFieldUnits(playerID)  // -1=自己, -2=对手

// 检查手牌
HasHandCards(playerID, minCount)

// 检查充能
HasCharge(minCharge)

// 检查状态
SourceHasStatus(status, minAmount)
TargetHasStatus(status, minAmount)
```

## 使用示例

### 掌门穆伶绝技 (4311003)

绝技效果："选择法力范围内双方各1个伙伴，花费差值气，移回手牌"

```go
// 注册带前提条件的绝技
preconditions := []Precondition{
    // 需要对手有场上单位（"双方各1个"）
    HasFieldUnits(-2),
    // 需要自己有场上单位
    HasFieldUnits(-1),
}

GetEffectRegistry().RegisterWithPrecondition(
    "4311003",
    TriggerUltimate,
    preconditions,
    func(ctx *EffectContext) error {
        // 绝技效果实现
        // ... 需要PendingAction进行目标选择
        ctx.Source.UltimateUsed = true
        return nil
    },
    true, // isActive
)
```

### 组合条件示例

```go
// 需要一个友方单位有护盾，且自己有至少3点充能
preconditions := []Precondition{
    AllOf{
        Conditions: []Precondition{
            AnyOf{
                Conditions: []Precondition{
                    // 检查各个友方单位是否有护盾
                    // 这里需要循环检查
                },
            },
            HasCharge(3),
        },
    },
}
```

## 架构优势

1. **声明式**: 条件以声明式方式定义，易于阅读和维护
2. **可测试**: 每个条件都是独立的，可以单独测试
3. **可复用**: 常用条件可以作为常量或快捷函数提供
4. **可扩展**: 新的条件类型只需实现 Precondition 接口
5. **清晰的错误信息**: 检查失败时提供清晰的错误信息，帮助玩家理解为什么无法使用

## 与 EffectSystem 的集成

```go
// EffectRegistry 扩展
type EffectRegistry struct {
    effects              map[string][]*CardEffect
    conditionalAbilities  map[string][]*CardAbilityWithPrecondition
}

type CardAbilityWithPrecondition struct {
    Preconditions []Precondition
    Effect        EffectHandler
    Trigger       EffectTrigger
    IsActive      bool
}

// 检查前提条件
func (r *EffectRegistry) CheckPreconditions(
    cardNumber string,
    trigger EffectTrigger,
    ctx *ConditionContext,
) (bool, string) {
    // ... 检查所有条件
}
```

这个设计使得卡牌效果的前提条件检查变得清晰、可维护和可扩展。
