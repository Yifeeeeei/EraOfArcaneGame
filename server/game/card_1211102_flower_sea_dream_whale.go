package game

type Card1211102FlowerSeaDreamWhale struct{ AlwaysActive }

func (Card1211102FlowerSeaDreamWhale) ID() string { return "1211102" }

func (Card1211102FlowerSeaDreamWhale) Name() string { return "花海梦鲸" }

func (Card1211102FlowerSeaDreamWhale) OnEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for _, number := range []string{"2201101", "2201102", "2201103"} {
		card := getCardDB()[number]
		if card == nil {
			continue
		}
		ps.Deck = append(ps.Deck, ctx.Engine.newCardInstance(card, ctx.PlayerID, ctx.Engine.State.TurnNumber))
	}
	ctx.Engine.shuffleDeck(ctx.PlayerID)
	ctx.Engine.emit(GameEvent{
		Type:   "flower_sea_dream_whale_shuffle_dreams",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"source": cardToInfo(ctx.Source),
			"count":  3,
		},
	})
	return nil
}

func (Card1211102FlowerSeaDreamWhale) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.ExtraData == nil {
		return nil
	}
	castPlayer, _ := ctx.ExtraData["cast_player"].(int)
	if castPlayer != ctx.PlayerID || ctx.Target.OwnerID != ctx.PlayerID || ctx.Target.Card == nil || !hasCardTag(ctx.Target.Card, "创造") {
		return nil
	}
	ctx.Source.Statuses[flowerSeaDreamWhaleCreationCountStatus]++
	if ctx.Source.Statuses[flowerSeaDreamWhaleCreationCountStatus] < 2 {
		return nil
	}
	searchDeckToHandByPredicateWithResult(ctx, "flower_sea_dream_whale_search",
		"花海梦鲸:检索1张幻创之梦", isDreamCreationCardInstance,
		func(*CardInstance) {
			ctx.Source.Statuses[flowerSeaDreamWhaleCreationCountStatus] -= 2
			if ctx.Source.Statuses[flowerSeaDreamWhaleCreationCountStatus] < 0 {
				ctx.Source.Statuses[flowerSeaDreamWhaleCreationCountStatus] = 0
			}
		})
	return nil
}

var _ OnEnterBehavior = Card1211102FlowerSeaDreamWhale{}

var _ OnSpellCastBehavior = Card1211102FlowerSeaDreamWhale{}

const flowerSeaDreamWhaleCreationCountStatus = "flower_sea_dream_whale_creation_count"

func isDreamCreationCardInstance(card *CardInstance) bool {
	return card != nil && card.Card != nil && (card.Card.Number == "2201101" || card.Card.Number == "2201102" || card.Card.Number == "2201103")
}
