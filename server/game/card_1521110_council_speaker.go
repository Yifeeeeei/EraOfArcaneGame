package game

type Card1521110CouncilSpeaker struct{ AlwaysActive }

func (Card1521110CouncilSpeaker) ID() string { return "1521110" }

func (Card1521110CouncilSpeaker) Name() string { return "议庭言客" }

func (Card1521110CouncilSpeaker) OnEnter(ctx *EffectContext) error {
	addGeneratedCardsToPlayerDeck(ctx, ctx.OpponentID, "2001102", 4)
	return nil
}

func (Card1521110CouncilSpeaker) OnDeath(ctx *EffectContext) error {
	ctx.Engine.moveDeckCardToTop(ctx.OpponentID, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.Number == "2001102"
	})
	return nil
}

func (e *Engine) moveDeckCardToTop(playerID int, predicate func(*CardInstance) bool) *CardInstance {
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	for i, card := range ps.Deck {
		if card == nil || (predicate != nil && !predicate(card)) {
			continue
		}
		ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
		ps.Deck = append([]*CardInstance{card}, ps.Deck...)
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"card":   cardToInfo(card),
			"effect": "deck_card_to_top",
		}})
		return card
	}
	return nil
}
