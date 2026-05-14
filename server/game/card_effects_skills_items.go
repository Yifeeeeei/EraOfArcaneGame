package game

func registerCard3621206ChaosDevour(r *EffectRegistry) {
	// 献祭你的1个生物:回复其生命值的生命,或获得光和暗
	r.Register("3621206", TriggerOnSpellHit, func(ctx *EffectContext) error {
		return nil
	})
}

func registerCard3621301ChargeRecycle(r *EffectRegistry) {
	// 冷却1.异能:每当1个单位死亡,充能1
	r.Register("3621301", TriggerOnFriendlyDeath, GainCharge(1))
	r.Register("3621301", TriggerOnEnemyDeath, GainCharge(1))
}

func registerCard2511101NinefoldRadiance(r *EffectRegistry) {
	// 绝技:双方将手牌全部丢弃,然后抽等量的牌
	r.RegisterActive("2511101", TriggerUltimate, func(ctx *EffectContext) error {
		for i := 0; i < 2; i++ {
			ps := ctx.Engine.State.Players[i]
			handSize := len(ps.Hand)
			ps.Graveyard = append(ps.Graveyard, ps.Hand...)
			ps.Hand = make([]*CardInstance, 0)
			drawn := ps.DrawCards(handSize)
			for _, c := range drawn {
				ctx.Engine.emit(GameEvent{
					Type:   "draw_card",
					Player: i,
					Data:   map[string]any{"card": cardToInfo(c)},
				})
			}
		}
		ctx.Source.UltimateUsed = true
		return nil
	})
}

func registerCard2521101BlessedLoneStar(r *EffectRegistry) {
	// 使1个友方伙伴获得负载+1光和+1血
	r.Register("2521101", TriggerOnEquip, func(ctx *EffectContext) error {
		if ctx.Target == nil {
			return nil
		}
		if ctx.Target.Card.ElementsGain == nil {
			ctx.Target.Card.ElementsGain = make(map[string]int)
		}
		ctx.Target.Card.ElementsGain["光"]++
		ctx.Target.CurrentLife++
		return nil
	})
}

func registerCard2521104GoldenDragonbone(r *EffectRegistry) {
	// 献祭:抽2张牌
	r.Register("2521104", TriggerOnEquip, func(ctx *EffectContext) error {
		return nil
	})
}
