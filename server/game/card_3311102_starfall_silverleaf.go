package game

type Card3311102StarfallSilverleaf struct{ AlwaysActive }

func (Card3311102StarfallSilverleaf) ID() string { return "3311102" }

func (Card3311102StarfallSilverleaf) Name() string { return "星落之银叶" }

func (Card3311102StarfallSilverleaf) OnDiscardBatch(ctx *EffectContext, discarded []*CardInstance) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	candidatesNow := func() []map[string]any {
		candidates := []map[string]any{}
		for _, card := range discarded {
			if card != nil && card != ctx.Source && cardInstanceInSlice(ctx.Engine.State.Players[ctx.PlayerID].Graveyard, card) {
				candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
			}
		}
		return candidates
	}
	candidates := candidatesNow()
	if len(candidates) == 0 {
		return nil
	}
	source := ctx.Source
	action := ctx.Engine.setPendingActionWithOptions(ctx.PlayerID, "starfall_silverleaf_store_discard",
		"星落之银叶:选择本次弃置的1张卡牌放在此卡下方", candidates, 1, 1, nil, false,
		func(selected []string) {
			if !ctx.Engine.cardStillOnField(source) {
				return
			}
			for _, card := range discarded {
				if card.InstanceID == firstSelected(selected) && cardInstanceInSlice(ctx.Engine.State.Players[ctx.PlayerID].Graveyard, card) {
					ctx.Engine.placeCardUnder(source, card)
					return
				}
			}
		}, nil, nil, nil, func() bool {
			return ctx.Engine.cardStillOnField(source) && !ctx.Engine.hasEffectiveStatus(source, StatusPetrify) && len(candidatesNow()) > 0
		})
	action.Refresh = func(p *PendingAction) { p.Candidates = candidatesNow() }
	return nil
}

func (Card3311102StarfallSilverleaf) OnSpellHit(ctx *EffectContext) error {
	if !isOwnSpellHit(ctx) {
		return nil
	}
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3311102" || len(ctx.Source.UnderCards) == 0 {
		return nil
	}
	candidates := make([]map[string]any, 0, len(ctx.Source.UnderCards))
	for _, card := range ctx.Source.UnderCards {
		if card != nil {
			candidates = append(candidates, candidateInfo(card, "under", "own"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "starfall_silverleaf_recycle_under_card",
		"星落之银叶:选择下方1张牌洗回卡组并抽1张", candidates, 1, 1,
		func(selected []string) {
			if !ctx.Engine.cardStillOnField(source) {
				return
			}
			card := ctx.Engine.detachCardFromKnownZones(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
			if card == nil {
				return
			}
			ctx.Engine.State.Players[ctx.PlayerID].Deck = append(ctx.Engine.State.Players[ctx.PlayerID].Deck, card)
			ctx.Engine.shuffleDeck(ctx.PlayerID)
			ctx.Engine.drawCards(ctx.PlayerID, 1)
		})
	return nil
}
