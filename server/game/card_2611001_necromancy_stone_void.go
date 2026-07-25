package game

import "eraofarcane/model"

type Card2611001NecromancyStoneVoid struct{ AlwaysActive }

func (Card2611001NecromancyStoneVoid) ID() string   { return "2611001" }
func (Card2611001NecromancyStoneVoid) Name() string { return "死灵魔石 虚无" }

func (Card2611001NecromancyStoneVoid) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx.Target == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	if ctx.Source == nil || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
		return nil
	}
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
	ctx.Source.UsedThisTurn++
	return nil
}
