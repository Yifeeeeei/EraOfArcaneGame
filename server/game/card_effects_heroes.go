package game

func registerCard4011001Skadi(r *EffectRegistry) {
	// 回合技:丢弃1张手牌,获得2点该卡牌属性种类的元素
	registerNoopActive(r, "4011001", TriggerPerTurn)
}

func registerCard4011002NoFace(r *EffectRegistry) {
	// 每当你让1张与你场上原有卡牌属性相同的卡牌入场,你受到1点伤害
	r.Register("4011002", TriggerOnUnitEnter, func(ctx *EffectContext) error {
		if ctx.Target == nil || ctx.Source == nil {
			return nil
		}
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		for _, c := range ctx.Engine.getAllFieldCards(ps) {
			if c == ctx.Target {
				continue
			}
			if c.Card.Category == ctx.Target.Card.Category && ps.Hero != nil {
				ctx.Engine.dealDamage(ps.Hero, 1, ctx.PlayerID)
				ctx.Engine.emit(GameEvent{
					Type:   "effect_trigger",
					Player: -1,
					Data: map[string]any{
						"source": cardToInfo(ctx.Source),
						"effect": "same_element_penalty",
						"damage": 1,
					},
				})
				return nil
			}
		}
		return nil
	})
}

func registerCard4011101PureSpirit(r *EffectRegistry) {
	// 你每使1张奥术以外的卡牌入场,使你的所有法术获得虚弱2
	r.Register("4011101", TriggerOnUnitEnter, func(ctx *EffectContext) error {
		if ctx.Target == nil || ctx.Target.Card.Category == "无" {
			return nil
		}
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		for i := 0; i < 5; i++ {
			if ps.Skills[i] != nil {
				ps.Skills[i].Statuses[StatusWeaken] += 2
			}
		}
		return nil
	})
}

func registerCard4111002WitchVerland(r *EffectRegistry) {
	// 回合技:你获得点燃1,然后本回合将此卡负载的1火变为1无
	r.RegisterActive("4111002", TriggerPerTurn, func(ctx *EffectContext) error {
		ctx.Source.Statuses[StatusBurn]++
		return nil
	})
}

func registerCard4111003Brahma(r *EffectRegistry) {
	// 绝技:本回合每当你的火焰法术命中,此卡获得负载+1火
	registerNoopActive(r, "4111003", TriggerUltimate)
}

func registerCard4111101Felin(r *EffectRegistry) {
	// 绝技:献祭1个火焰伙伴,下1个火焰卡牌减费
	registerNoopActive(r, "4111101", TriggerUltimate)
}

func registerCard4111102Kran(r *EffectRegistry) {
	// 绝技:选择火焰法术,未命中时获得其攻的火
	registerNoopActive(r, "4111102", TriggerUltimate)
}

func registerCard4111201RudokClark(r *EffectRegistry) {
	// 法术战斗威力超3点时获得护盾1和1火
	r.Register("4111201", TriggerOnDefend, func(ctx *EffectContext) error {
		return nil
	})
}

func registerCard4211001Bartel(r *EffectRegistry) {
	// 绝技:下一张手牌属性视为水
	registerNoopActive(r, "4211001", TriggerUltimate)
}

func registerCard4211003CrystalHeart(r *EffectRegistry) {
	// 绝技:本回合剩余时间内,技能区内法术获得冻结1
	registerNoopActive(r, "4211003", TriggerUltimate)
}

func registerCard4211101CoralBelly(r *EffectRegistry) {
	// 本局游戏你第一次使用法术攻击时,使该法术永久+3威
	r.Register("4211101", TriggerOnSpellCast, func(ctx *EffectContext) error {
		if ctx.Target != nil && ctx.Target.Card.IsSkill() && ctx.Source.Statuses["first_spell_used"] == 0 {
			ctx.Source.Statuses["first_spell_used"] = 1
		}
		return nil
	})
}

func registerCard4211102Sophia(r *EffectRegistry) {
	// 绝技:移除场上任意1个单位的1层冻结,对其造成2点伤害
	r.RegisterActive("4211102", TriggerUltimate, func(ctx *EffectContext) error {
		if ctx.Target != nil {
			if ctx.Target.Statuses[StatusFreeze] > 0 {
				ctx.Target.Statuses[StatusFreeze]--
			}
			ctx.Engine.dealDamage(ctx.Target, 2, ctx.Target.OwnerID)
		}
		ctx.Source.UltimateUsed = true
		return nil
	})
}

func registerCard4311001Su(r *EffectRegistry) {
	// 绝技:丢弃2张大气手牌,对任意1名敌人造成1点伤害
	r.RegisterActive("4311001", TriggerUltimate, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		discarded := 0
		for i := len(ps.Hand) - 1; i >= 0 && discarded < 2; i-- {
			if ps.Hand[i].Card.Category == "气" {
				ps.Graveyard = append(ps.Graveyard, ps.Hand[i])
				ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
				discarded++
			}
		}
		if discarded < 2 {
			return nil
		}
		target := findAnyUnit(ctx.Engine.State.Players[ctx.OpponentID])
		if target != nil {
			ctx.Engine.dealDamage(target, 1, ctx.OpponentID)
		}
		ctx.Source.UltimateUsed = true
		return nil
	})
}

func registerCard4311003Muling(r *EffectRegistry) {
	// 绝技:选择法力范围内双方各1个伙伴,花费差值气,移回手牌
	registerNoopActive(r, "4311003", TriggerUltimate)
}

func registerCard4311101Soland(r *EffectRegistry) {
	// 绝技:使驱动和聚能法术永久+1威,其他法术永久-2威
	registerNoopActive(r, "4311101", TriggerUltimate)
}

func registerCard4311102Fog(r *EffectRegistry) {
	// 绝技:双方各自召唤的下1个伙伴入场时获得隐蔽2
	registerNoopActive(r, "4311102", TriggerUltimate)
}

func registerCard4311201Lillian(r *EffectRegistry) {
	// 回合技:将此卡移动至另一位置
	registerNoopActive(r, "4311201", TriggerPerTurn)
}

func registerCard4311202Trachi(r *EffectRegistry) {
	// 绝技:花费2气,双方各抽2张牌,或双方随机弃2张牌
	registerNoopActive(r, "4311202", TriggerUltimate)
}

func registerCard4311302Yuling(r *EffectRegistry) {
	// 绝技:将手牌/场上2张牌洗回卡组,抽2张牌
	registerNoopActive(r, "4311302", TriggerUltimate)
}

func registerCard4411001Whitebeard(r *EffectRegistry) {
	// 检索1张地属性野兽/植物/精灵代替首回合抽牌
	r.Register("4411001", TriggerOnTurnStart, func(ctx *EffectContext) error {
		if ctx.Engine.State.TurnNumber != 1 {
			return nil
		}
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		for i, c := range ps.Deck {
			if c.Card.Category == "地" && hasTag(c.Card.Tag, "野兽", "植物", "精灵") {
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

func registerCard4411201Hisson(r *EffectRegistry) {
	// 对方一回合内召唤超2个时,每超1个-1地
	r.Register("4411201", TriggerOnUnitEnter, func(ctx *EffectContext) error {
		return nil
	})
}

func registerCard4411202Dorothy(r *EffectRegistry) {
	// 召唤第一个野兽/植物/精灵后抽1张牌
	r.Register("4411202", TriggerOnUnitEnter, func(ctx *EffectContext) error {
		if ctx.Target == nil {
			return nil
		}
		tag := ctx.Target.Card.Tag
		if !hasTag(tag, "野兽", "植物", "精灵") {
			return nil
		}
		key := "first_" + getFirstMatchingTag(tag, "野兽", "植物", "精灵")
		if ctx.Source.Statuses[key] != 0 {
			return nil
		}
		ctx.Source.Statuses[key] = 1
		drawn := ctx.Engine.State.Players[ctx.PlayerID].DrawCards(1)
		if len(drawn) > 0 {
			ctx.Engine.emit(GameEvent{
				Type:   "draw_card",
				Player: ctx.PlayerID,
				Data:   map[string]any{"card": cardToInfo(drawn[0])},
			})
		}
		return nil
	})
}

func registerCard4511001Maris(r *EffectRegistry) {
	// 绝技:当敌方造成伤害时,直到下回合结束,每当友方单位受伤获得2光
	registerNoopActive(r, "4511001", TriggerUltimate)
}

func registerCard4511101Sivar(r *EffectRegistry) {
	// 绝技:当受到3+伤害,下回合结束前所有友方不受伤
	registerNoopActive(r, "4511101", TriggerUltimate)
}

func registerCard4611001Alice(r *EffectRegistry) {
	// 回合技:当1个你的伙伴死亡,使你的1个法术+1威
	r.Register("4611001", TriggerOnFriendlyDeath, func(ctx *EffectContext) error {
		return nil
	})
}

func registerCard4611002Fuye(r *EffectRegistry) {
	// 绝技:使你的1个伙伴攻击和负载翻倍,但会在回合结束时死亡
	r.RegisterActive("4611002", TriggerUltimate, func(ctx *EffectContext) error {
		if ctx.Target != nil {
			ctx.Target.CurrentAttack *= 2
			ctx.Target.Statuses["临时"] = 1
		}
		ctx.Source.UltimateUsed = true
		return nil
	})
}

func registerCard4611202Yuexi(r *EffectRegistry) {
	// 每当使用暗影技能,查看卡组顶1张.回合技:可以重洗
	r.Register("4611202", TriggerOnSpellCast, func(ctx *EffectContext) error {
		return nil
	})
}
