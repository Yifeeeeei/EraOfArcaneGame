package game

import (
	"eraofarcane/model"
)

type Card1121111LoneStarFireSeed struct{ AlwaysActive }

func (Card1121111LoneStarFireSeed) ID() string { return "1121111" }

func (Card1121111LoneStarFireSeed) Name() string { return "孤星火种" }

func (Card1121111LoneStarFireSeed) DamageScope() DamageScope { return DamageAny }

func (Card1121111LoneStarFireSeed) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if ctx == nil || ctx.Source == nil || ctx.Source.UltimateUsed || event.Target == nil || event.Target == ctx.Source || event.Target.Card == nil || !event.Target.Card.IsCompanion() {
		return nil
	}
	if event.Element != model.ElementFire && event.Status != StatusBurn {
		return nil
	}
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementFire, 1, ctx.Source)
	ctx.Source.UltimateUsed = true
	return nil
}
