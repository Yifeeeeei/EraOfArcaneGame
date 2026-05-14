package game

import (
	"math/rand"
	"strings"
)

func registerCard1221004FrostPuppet(r *EffectRegistry) {
	// 入场:对法力范围内1个敌方伙伴冻结1
	r.Register("1221004", TriggerOnEnter, func(ctx *EffectContext) error {
		if ctx.Target != nil {
			ctx.Target.Statuses[StatusFreeze]++
			return nil
		}
		opponent := ctx.Engine.State.Players[ctx.OpponentID]
		frontRow := opponent.GetFrontRow()
		if frontRow < 0 {
			return nil
		}
		for col := 0; col < 3; col++ {
			if opponent.Units[col][frontRow] != nil {
				opponent.Units[col][frontRow].Statuses[StatusFreeze]++
				break
			}
		}
		return nil
	})
}

func registerCard1221203AbyssFogDemon(r *EffectRegistry) {
	// 入场:检索1张水属性或深渊地形牌
	r.Register("1221203", TriggerOnEnter, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		for i, c := range ps.Deck {
			if c.Card.Category == "水" || (HasKeyword(c.Card.Description, "深渊") && HasKeyword(c.Card.Description, "地形")) {
				ps.Hand = append(ps.Hand, c)
				ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
				shuffleDeck(ps.Deck)
				ctx.Engine.emit(GameEvent{
					Type:   "effect_trigger",
					Player: ctx.PlayerID,
					Data: map[string]any{
						"source": cardToInfo(ctx.Source),
						"effect": "search",
						"card":   cardToInfo(c),
					},
				})
				break
			}
		}
		return nil
	})
}

func registerCard1221207Reconstructor(r *EffectRegistry) {
	// 回合技:在你检索1张卡牌后,获得负载+1水
	registerNoopActive(r, "1221207", TriggerPerTurn)
}

func registerCard1221208HunterCruiser(r *EffectRegistry) {
	// 回合技:移动至其他位置
	registerNoopActive(r, "1221208", TriggerPerTurn)
}

func registerCard1321301Technician(r *EffectRegistry) {
	// 绝技:本回合每当洗回卡组1张牌,充能1
	registerNoopActive(r, "1321301", TriggerUltimate)
}

func registerCard1321304SkyWell(r *EffectRegistry) {
	// 引魔.你的具有威力的卡牌获得穿透
	r.Register("1321304", TriggerPassive, func(ctx *EffectContext) error {
		return nil
	})
}

func registerCard1321306TaijiMaster(r *EffectRegistry) {
	// 回合技:将场上或手牌1张卡洗回卡组,抽1张牌
	registerNoopActive(r, "1321306", TriggerPerTurn)
}

func registerCard1321308FloatingPilot(r *EffectRegistry) {
	// 入场:充能1
	r.Register("1321308", TriggerOnEnter, GainCharge(1))
}

func registerCard1321309TaijiHeir(r *EffectRegistry) {
	// 入场:从卡组检索1张调和、化劲、收势
	r.Register("1321309", TriggerOnEnter, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		targets := []string{"调和", "化劲", "收势"}
		for i, c := range ps.Deck {
			for _, name := range targets {
				if strings.Contains(c.Card.Name, name) {
					ps.Hand = append(ps.Hand, c)
					ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
					shuffleDeck(ps.Deck)
					ctx.Engine.emit(GameEvent{
						Type:   "effect_trigger",
						Player: ctx.PlayerID,
						Data: map[string]any{
							"source": cardToInfo(ctx.Source),
							"effect": "search",
							"card":   cardToInfo(c),
						},
					})
					return nil
				}
			}
		}
		return nil
	})
}

func registerCard1611103RobertBlackpine(r *EffectRegistry) {
	// 友方单位受友方伤害放标记,友方因友方死亡放2标记
	r.Register("1611103", TriggerOnFriendlyDeath, func(ctx *EffectContext) error {
		ctx.Source.Statuses["鲜血标记"] += 2
		return nil
	})
	r.RegisterActive("1611103", TriggerPerTurn, func(ctx *EffectContext) error {
		if ctx.Source.Statuses["鲜血标记"] < 3 {
			return nil
		}
		ctx.Source.Statuses["鲜血标记"] -= 3
		switch rand.Intn(3) {
		case 0:
			ctx.Source.CurrentLife++
		case 1:
			if ctx.Source.Card.ElementsGain == nil {
				ctx.Source.Card.ElementsGain = make(map[string]int)
			}
			ctx.Source.Card.ElementsGain["暗"]++
		case 2:
			ctx.Source.CurrentAttack++
		}
		return nil
	})
}

func registerCard1621103BloodPuppet(r *EffectRegistry) {
	// 入场:对你的人物造成2点伤害
	r.Register("1621103", TriggerOnEnter, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		if ps.Hero != nil {
			ctx.Engine.dealDamage(ps.Hero, 2, ctx.PlayerID)
		}
		return nil
	})
}

func registerCard1621112SilentHunter(r *EffectRegistry) {
	// 遗言:对任意1个敌人造成1点伤害
	r.Register("1621112", TriggerOnDeath, func(ctx *EffectContext) error {
		target := findAnyUnit(ctx.Engine.State.Players[ctx.OpponentID])
		if target != nil {
			ctx.Engine.dealDamage(target, 1, ctx.OpponentID)
		}
		return nil
	})
}

func registerCard1621113SilentPriest(r *EffectRegistry) {
	// 遗言:使1个友方伙伴负载+1暗
	r.Register("1621113", TriggerOnDeath, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				unit := ps.Units[col][row]
				if unit != nil && unit != ctx.Source && unit.Card.IsCompanion() {
					if unit.Card.ElementsGain == nil {
						unit.Card.ElementsGain = make(map[string]int)
					}
					unit.Card.ElementsGain["暗"]++
					return nil
				}
			}
		}
		return nil
	})
}

func registerCard1621114SoulSymbiote(r *EffectRegistry) {
	// 遗言:给最多2个法术放灵魂标记物
	r.Register("1621114", TriggerOnDeath, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		count := 0
		for i := 0; i < 5 && count < 2; i++ {
			if ps.Skills[i] != nil {
				ps.Skills[i].Statuses["灵魂标记"]++
				count++
			}
		}
		return nil
	})
}
