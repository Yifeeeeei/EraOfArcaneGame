package game

import "eraofarcane/model"

type Card3501001UnitedHope struct{}

func (Card3501001UnitedHope) ID() string   { return "3501001" }
func (Card3501001UnitedHope) Name() string { return "团结的希望" }

func (Card3501001UnitedHope) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	limit := min(5, len(ps.Deck))
	for i := 0; i < limit; i++ {
		card := ps.Deck[i]
		if card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementLight {
			continue
		}
		ps.Hand = append(ps.Hand, card)
		ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
		ctx.Engine.shuffleDeck(ctx.PlayerID)
		ctx.Engine.emit(GameEvent{Type: "search_card", Player: ctx.PlayerID, Data: map[string]any{"card": cardToInfo(card)}})
		return nil
	}
	ctx.Engine.shuffleDeck(ctx.PlayerID)
	return nil
}
