package game

import "eraofarcane/model"

type Card3321007SourceWind struct{ AlwaysActive }

func (Card3321007SourceWind) ID() string   { return "3321007" }
func (Card3321007SourceWind) Name() string { return "源力之风" }

func (Card3321007SourceWind) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for len(ps.Hand) < ctx.Engine.State.HandLimit && len(ps.Deck) > 0 {
		if !ps.PayCost(map[string]int{model.ElementAir: 1}) {
			return nil
		}
		ctx.Engine.drawCards(ctx.PlayerID, 1)
	}
	return nil
}
