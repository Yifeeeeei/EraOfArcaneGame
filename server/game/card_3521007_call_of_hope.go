package game

import "eraofarcane/model"

type Card3521007CallOfHope struct{ AlwaysActive }

func (Card3521007CallOfHope) ID() string   { return "3521007" }
func (Card3521007CallOfHope) Name() string { return "希望呼唤" }

func (Card3521007CallOfHope) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	ctx.Engine.drawFirstDeckMatch(ctx.PlayerID, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Card.Category == model.ElementLight
	})
	ctx.Engine.shuffleDeck(ctx.PlayerID)
	return nil
}
