package game

import "eraofarcane/model"

type Card1121111LoneStarFireSeed struct{ AlwaysActive }

func (Card1121111LoneStarFireSeed) ID() string   { return "1121111" }
func (Card1121111LoneStarFireSeed) Name() string { return "孤星火种" }

func (Card1121111LoneStarFireSeed) OnDamaged(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target == ctx.Source || ctx.Target.Card == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	if ctx.ExtraData == nil || (ctx.ExtraData["damage_element"] != model.ElementFire && ctx.ExtraData["status_damage"] != StatusBurn) {
		return nil
	}
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementFire, 1, ctx.Source)
	return nil
}
