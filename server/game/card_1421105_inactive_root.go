package game

import (
	"eraofarcane/model"
)

type Card1421105InactiveRoot struct{ AlwaysActive }

func (Card1421105InactiveRoot) ID() string { return "1421105" }

func (Card1421105InactiveRoot) Name() string { return "失活的根须" }

func (Card1421105InactiveRoot) IsPrayerAbility() bool { return true }

func (Card1421105InactiveRoot) OnPerTurn(ctx *EffectContext) error {
	if totalLoad(ctx.Source) == 0 {
		ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementEarth, 1, ctx.Source)
	}
	return nil
}
