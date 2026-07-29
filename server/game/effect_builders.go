package game

// effect_builders.go — 可复用的效果构建器
// 每个函数返回一个 EffectHandler，可直接用于 Register/RegisterActive
//
// 使用方式:
//   r.Register("1234567", TriggerOnEnter, DrawCards(2))
//   r.Register("1234567", TriggerOnDeath, DealDamageAuto(3))

import (
	"fmt"
	"math/rand"

	"eraofarcane/model"
)

// ══════════════════════════════════════
// 抽牌 / 充能
// ══════════════════════════════════════

// DrawCards 抽N张牌
func DrawCards(n int) EffectHandler {
	return func(ctx *EffectContext) error {
		ctx.Engine.drawCards(ctx.PlayerID, n)
		return nil
	}
}

// GainCharge 获得N点充能
func GainCharge(n int) EffectHandler {
	return func(ctx *EffectContext) error {
		ctx.Engine.addCharge(ctx.PlayerID, n)
		ctx.Engine.emit(GameEvent{
			Type: "effect_trigger", Player: -1,
			Data: map[string]any{
				"source": cardToInfo(ctx.Source), "effect": "charge", "amount": n,
			},
		})
		return nil
	}
}

// ══════════════════════════════════════
// 伤害
// ══════════════════════════════════════

// DealDamageToTarget 对指定目标造成N点伤害（需要ctx.Target）
func DealDamageToTarget(n int) EffectHandler {
	return func(ctx *EffectContext) error {
		if ctx.Target == nil {
			return nil
		}
		ctx.Engine.dealDamage(ctx.Target, n, ctx.OpponentID)
		ctx.Engine.emit(GameEvent{
			Type: "effect_trigger", Player: -1,
			Data: map[string]any{
				"source": cardToInfo(ctx.Source), "effect": "damage",
				"amount": n, "target": cardToInfo(ctx.Target),
			},
		})
		return nil
	}
}

// DealDamageAuto 对前排敌方造成N点伤害（无目标时自动选择）
func DealDamageAuto(n int) EffectHandler {
	return func(ctx *EffectContext) error {
		target := ctx.Target
		if target == nil {
			opponent := ctx.Engine.State.Players[ctx.OpponentID]
			target = findFrontRowUnit(opponent)
		}
		if target == nil {
			return nil
		}
		ownerOfTarget := target.OwnerID
		ctx.Engine.dealDamage(target, n, ownerOfTarget)
		ctx.Engine.emit(GameEvent{
			Type: "effect_trigger", Player: -1,
			Data: map[string]any{
				"source": cardToInfo(ctx.Source), "effect": "damage",
				"amount": n, "target": cardToInfo(target),
			},
		})
		return nil
	}
}

// DealDamageToRandomEnemy 对随机一个敌方单位造成N点伤害
func DealDamageToRandomEnemy(n int) EffectHandler {
	return func(ctx *EffectContext) error {
		opponent := ctx.Engine.State.Players[ctx.OpponentID]
		target := findRandomUnit(opponent)
		if target == nil {
			return nil
		}
		ctx.Engine.dealDamage(target, n, ctx.OpponentID)
		ctx.Engine.emit(GameEvent{
			Type: "effect_trigger", Player: -1,
			Data: map[string]any{
				"source": cardToInfo(ctx.Source), "effect": "damage",
				"amount": n, "target": cardToInfo(target),
			},
		})
		return nil
	}
}

// DealDamageToSelfHero 对自己的英雄造成N点伤害
func DealDamageToSelfHero(n int) EffectHandler {
	return func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		if ps.Hero != nil {
			ctx.Engine.dealDamage(ps.Hero, n, ctx.PlayerID)
		}
		return nil
	}
}

// DealDamageToEnemyHero 对敌方英雄造成N点伤害
func DealDamageToEnemyHero(n int) EffectHandler {
	return func(ctx *EffectContext) error {
		opponent := ctx.Engine.State.Players[ctx.OpponentID]
		if opponent.Hero != nil {
			ctx.Engine.dealDamage(opponent.Hero, n, ctx.OpponentID)
		}
		return nil
	}
}

// ══════════════════════════════════════
// 状态效果
// ══════════════════════════════════════

// ApplyStatusToTarget 对指定目标施加状态（需要ctx.Target）
func ApplyStatusToTarget(status string, amount int) EffectHandler {
	return func(ctx *EffectContext) error {
		if ctx.Target == nil {
			return nil
		}
		if !ctx.Engine.addStatus(ctx.Target, status, amount) {
			return nil
		}
		ctx.Engine.emit(GameEvent{
			Type: "effect_trigger", Player: -1,
			Data: map[string]any{
				"source": cardToInfo(ctx.Source), "effect": "apply_status",
				"status": status, "amount": amount, "target": cardToInfo(ctx.Target),
			},
		})
		return nil
	}
}

// ApplyStatusToSelf 对自身施加状态
func ApplyStatusToSelf(status string, amount int) EffectHandler {
	return func(ctx *EffectContext) error {
		if !ctx.Engine.addStatus(ctx.Source, status, amount) {
			return nil
		}
		ctx.Engine.emit(GameEvent{
			Type: "effect_trigger", Player: -1,
			Data: map[string]any{
				"source": cardToInfo(ctx.Source), "effect": "apply_status_self",
				"status": status, "amount": amount,
			},
		})
		return nil
	}
}

// RemoveStatusFromTarget 移除目标的指定状态
func RemoveStatusFromTarget(status string) EffectHandler {
	return func(ctx *EffectContext) error {
		if ctx.Target == nil {
			return nil
		}
		delete(ctx.Target.Statuses, status)
		return nil
	}
}

// ══════════════════════════════════════
// 护盾 / 隐蔽
// ══════════════════════════════════════

// GainShield 自身获得护盾
func GainShield(amount int) EffectHandler {
	return func(ctx *EffectContext) error {
		ctx.Engine.gainPlayerShield(ctx.PlayerID, amount)
		return nil
	}
}

// GainStealth 自身获得隐蔽
func GainStealth(amount int) EffectHandler {
	return ApplyStatusToSelf(StatusStealth, amount)
}

// GiveShieldToTarget 给目标护盾
func GiveShieldToTarget(amount int) EffectHandler {
	return func(ctx *EffectContext) error {
		if ctx.Target != nil {
			ctx.Engine.gainPlayerShield(ctx.Target.OwnerID, amount)
		}
		return nil
	}
}

// ══════════════════════════════════════
// 数值修改
// ══════════════════════════════════════

// ModifySelfAttack 修改自身攻击力（可正可负）
func ModifySelfAttack(delta int) EffectHandler {
	return func(ctx *EffectContext) error {
		ctx.Source.CurrentAttack += delta
		if ctx.Source.CurrentAttack < 0 {
			ctx.Source.CurrentAttack = 0
		}
		return nil
	}
}

// ModifySelfLife 修改自身生命值
func ModifySelfLife(delta int) EffectHandler {
	return func(ctx *EffectContext) error {
		ctx.Source.CurrentLife += delta
		return nil
	}
}

// ModifyTargetAttack 修改目标攻击力
func ModifyTargetAttack(delta int) EffectHandler {
	return func(ctx *EffectContext) error {
		if ctx.Target == nil {
			return nil
		}
		ctx.Target.CurrentAttack += delta
		if ctx.Target.CurrentAttack < 0 {
			ctx.Target.CurrentAttack = 0
		}
		return nil
	}
}

// ModifyTargetLife 修改目标生命值
func ModifyTargetLife(delta int) EffectHandler {
	return func(ctx *EffectContext) error {
		if ctx.Target == nil {
			return nil
		}
		ctx.Target.CurrentLife += delta
		return nil
	}
}

// HealTarget 治疗目标N点生命
func HealTarget(n int) EffectHandler {
	return func(ctx *EffectContext) error {
		if ctx.Target == nil {
			return nil
		}
		ctx.Target.CurrentLife += n
		if ctx.Target.CurrentLife > ctx.Target.Card.Life {
			ctx.Target.CurrentLife = ctx.Target.Card.Life
		}
		return nil
	}
}

// HealSelf 治疗自身N点生命
func HealSelf(n int) EffectHandler {
	return func(ctx *EffectContext) error {
		ctx.Source.CurrentLife += n
		if ctx.Source.CurrentLife > ctx.Source.Card.Life {
			ctx.Source.CurrentLife = ctx.Source.Card.Life
		}
		return nil
	}
}

// HealHero 治疗己方英雄N点生命
func HealHero(n int) EffectHandler {
	return func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		if ps.Hero != nil {
			ps.Hero.CurrentLife += n
			if ps.Hero.CurrentLife > ps.Hero.Card.Life {
				ps.Hero.CurrentLife = ps.Hero.Card.Life
			}
		}
		return nil
	}
}

// ══════════════════════════════════════
// 元素
// ══════════════════════════════════════

// GainElements 获得指定元素
func GainElements(elements map[string]int) EffectHandler {
	return func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		ps.GainElements(elements)
		ctx.Engine.emit(GameEvent{
			Type: "effect_trigger", Player: ctx.PlayerID,
			Data: map[string]any{
				"source": cardToInfo(ctx.Source), "effect": "gain_elements",
				"elements": elements,
			},
		})
		return nil
	}
}

// ══════════════════════════════════════
// 组合效果
// ══════════════════════════════════════

// Combine 组合多个效果，依次执行
func Combine(handlers ...EffectHandler) EffectHandler {
	return func(ctx *EffectContext) error {
		for _, h := range handlers {
			if err := h(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

// NoEffect 空效果（用于标记已确认无运行效果的卡）
func NoEffect() EffectHandler {
	return func(ctx *EffectContext) error {
		return nil
	}
}

// ══════════════════════════════════════
// 辅助函数
// ══════════════════════════════════════

// findRandomUnit 随机找一个单位
func findRandomUnit(ps *PlayerState) *CardInstance {
	var units []*CardInstance
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil {
				units = append(units, ps.Units[col][row])
			}
		}
	}
	if len(units) == 0 {
		return nil
	}
	return units[rand.Intn(len(units))]
}

// ══════════════════════════════════════
// 检索卡组
// ══════════════════════════════════════

// SearchDeckAndDraw 检索卡组拿取符合条件的卡（简化版：自动拿取）
// 完整版应该让玩家选择，这里先实现自动拿第一张匹配
func SearchDeckAndDraw(predicate func(*model.Card) bool) EffectHandler {
	return func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		for _, c := range ps.Deck {
			if predicate(c.Card) {
				// 从卡组移除
				ps.Deck = removeCardFromDeck(ps.Deck, c.InstanceID)
				// 加入手牌
				ps.Hand = append(ps.Hand, c)
				ctx.Engine.emit(GameEvent{
					Type: "search_draw", Player: ctx.PlayerID,
					Data: map[string]any{"card": cardToInfo(c)},
				})
				return nil
			}
		}
		return nil
	}
}

// ══════════════════════════════════════
// 弃牌
// ══════════════════════════════════════

// DiscardRandom 随机弃N张手牌
func DiscardRandom(n int) EffectHandler {
	return func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		for i := 0; i < n && len(ps.Hand) > 0; i++ {
			idx := rand.Intn(len(ps.Hand))
			card := ps.Hand[idx]
			ps.Hand = append(ps.Hand[:idx], ps.Hand[idx+1:]...)
			ctx.Engine.emit(GameEvent{
				Type: "discard", Player: ctx.PlayerID,
				Data: map[string]any{"card": cardToInfo(card)},
			})
		}
		return nil
	}
}

// DiscardSelf 弃掉此卡自身（用于献祭效果）
func DiscardSelf() EffectHandler {
	return func(ctx *EffectContext) error {
		if ctx.Source == nil {
			return nil
		}
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		// 从手牌中移除
		for i, c := range ps.Hand {
			if c.InstanceID == ctx.Source.InstanceID {
				ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
				ctx.Engine.emit(GameEvent{
					Type: "discard", Player: ctx.PlayerID,
					Data: map[string]any{"card": cardToInfo(c)},
				})
				return nil
			}
		}
		return nil
	}
}

// ══════════════════════════════════════
// 召唤
// ══════════════════════════════════════

// SummonToken 召唤衍生物到前排空位
func SummonToken(cardID string) EffectHandler {
	return func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		// 获取卡牌数据库
		cardDB := getCardDB()
		// 找前排空位
		for col := 0; col < 3; col++ {
			if ps.Units[col][0] == nil {
				card := cardDB[cardID]
				if card == nil {
					return fmt.Errorf("unknown card: %s", cardID)
				}
				// 获取当前回合数
				turn := 0
				if ctx.Engine.State != nil {
					turn = ctx.Engine.State.TurnNumber
				}
				instance := NewCardInstance(card, ctx.PlayerID, turn)
				instance.Position = &Position{Col: col, Row: 0}
				ps.Units[col][0] = instance
				ctx.Engine.ApplyKeywordOnEnter(instance)
				ctx.Engine.ApplySummonModifiersOnEnter(instance)
				ctx.Engine.emit(GameEvent{
					Type: "summon", Player: ctx.PlayerID,
					Data: map[string]any{"card": cardToInfo(instance)},
				})
				return nil
			}
		}
		return nil
	}
}

// ══════════════════════════════════════
// 献祭相关
// ══════════════════════════════════════

// SacrificeSelfAndDo 献祭自身并执行效果
func SacrificeSelfAndDo(effect EffectHandler) EffectHandler {
	return Combine(DiscardSelf(), effect)
}

// SacrificeTarget 献祭目标单位（造成等同于生命值的"即死"伤害）
func SacrificeTarget() EffectHandler {
	return func(ctx *EffectContext) error {
		if ctx.Target == nil {
			return nil
		}
		// 对目标造成等同于其当前生命的伤害（即死效果）
		damage := ctx.Target.CurrentLife
		ctx.Engine.dealDamage(ctx.Target, damage, ctx.OpponentID)
		return nil
	}
}

// ══════════════════════════════════════
// 辅助函数
// ══════════════════════════════════════

func removeCardFromDeck(deck []*CardInstance, instanceID string) []*CardInstance {
	for i, c := range deck {
		if c.InstanceID == instanceID {
			return append(deck[:i], deck[i+1:]...)
		}
	}
	return deck
}
