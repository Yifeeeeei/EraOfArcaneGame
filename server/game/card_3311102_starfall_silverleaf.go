package game

type Card3311102StarfallSilverleaf struct{ AlwaysActive }

func (Card3311102StarfallSilverleaf) ID() string { return "3311102" }

func (Card3311102StarfallSilverleaf) Name() string { return "星落之银叶" }

func (Card3311102StarfallSilverleaf) OnDiscard(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.ExtraData == nil {
		return nil
	}
	discardedPlayer, _ := ctx.ExtraData["discarded_player"].(int)
	if discardedPlayer != ctx.PlayerID || ctx.Target == ctx.Source {
		return nil
	}
	discarded := ctx.Target
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "starfall_silverleaf_store_discard",
		"星落之银叶:将弃置的卡牌放在此卡下方", []map[string]any{candidateInfo(discarded, "graveyard", "own")}, 1, 1,
		func(selected []string) {
			if !ctx.Engine.cardStillOnField(source) {
				return
			}
			ctx.Engine.placeCardUnder(source, discarded)
		})
	return nil
}

func (Card3311102StarfallSilverleaf) OnSpellHit(ctx *EffectContext) error {
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
			resetCardForPublicSpecialZone(card)
			ctx.Engine.State.Players[ctx.PlayerID].Deck = append(ctx.Engine.State.Players[ctx.PlayerID].Deck, card)
			ctx.Engine.shuffleDeck(ctx.PlayerID)
			ctx.Engine.drawCards(ctx.PlayerID, 1)
		})
	return nil
}
