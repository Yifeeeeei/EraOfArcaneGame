package game

type Card4011001Skadi struct{ noopPerTurn }

func (Card4011001Skadi) ID() string   { return "4011001" }
func (Card4011001Skadi) Name() string { return "\"南境百灵\" 斯卡尔蒂 罗佳" }

type Card4011002NoFace struct{}

func (Card4011002NoFace) ID() string   { return "4011002" }
func (Card4011002NoFace) Name() string { return "\"无面\"" }
func (Card4011002NoFace) OnUnitEnter(ctx *EffectContext) error {
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
}

type Card4111002WitchVerland struct{}

func (Card4111002WitchVerland) ID() string   { return "4111002" }
func (Card4111002WitchVerland) Name() string { return "女巫 维兰德" }
func (Card4111002WitchVerland) OnPerTurn(ctx *EffectContext) error {
	ctx.Source.Statuses[StatusBurn]++
	return nil
}

type Card4111003Brahma struct{ noopUltimate }

func (Card4111003Brahma) ID() string   { return "4111003" }
func (Card4111003Brahma) Name() string { return "大祭司 梵天" }

type Card4211001Bartel struct{ noopUltimate }

func (Card4211001Bartel) ID() string   { return "4211001" }
func (Card4211001Bartel) Name() string { return "\"浪之人\" 巴特尔" }

type Card4211003CrystalHeart struct{ noopUltimate }

func (Card4211003CrystalHeart) ID() string   { return "4211003" }
func (Card4211003CrystalHeart) Name() string { return "凛冬城主 水晶心" }

type Card4311001Su struct{}

func (Card4311001Su) ID() string   { return "4311001" }
func (Card4311001Su) Name() string { return "雷术士 肃" }
func (Card4311001Su) OnUltimate(ctx *EffectContext) error {
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
	return nil
}

type Card4311003Muling struct{ noopUltimate }

func (Card4311003Muling) ID() string   { return "4311003" }
func (Card4311003Muling) Name() string { return "掌门 穆伶" }

type Card4411001Whitebeard struct{}

func (Card4411001Whitebeard) ID() string   { return "4411001" }
func (Card4411001Whitebeard) Name() string { return "森林隐士 白须" }
func (Card4411001Whitebeard) OnTurnStart(ctx *EffectContext) error {
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
			return nil
		}
	}
	return nil
}

type Card4511001Maris struct{ noopUltimate }

func (Card4511001Maris) ID() string   { return "4511001" }
func (Card4511001Maris) Name() string { return "圣使 玛丽斯 南森埃尔" }

type Card4611001Alice struct{}

func (Card4611001Alice) ID() string   { return "4611001" }
func (Card4611001Alice) Name() string { return "暗影学者 爱莉斯" }
func (Card4611001Alice) OnFriendlyDeath(*EffectContext) error {
	return nil
}

type Card4611002Fuye struct{}

func (Card4611002Fuye) ID() string   { return "4611002" }
func (Card4611002Fuye) Name() string { return "芙雅夫人" }
func (Card4611002Fuye) OnUltimate(ctx *EffectContext) error {
	if ctx.Target != nil {
		ctx.Target.CurrentAttack *= 2
		ctx.Target.Statuses["临时"] = 1
	}
	return nil
}
