package game

import "eraofarcane/model"

type Card3301001StormFury struct{ AlwaysActive }

func (Card3301001StormFury) ID() string   { return "3301001" }
func (Card3301001StormFury) Name() string { return "风暴之怒" }

func (Card3301001StormFury) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.Target == nil || ctx.Target.Card == nil || ctx.Target.Card.Category != model.ElementAir {
		return
	}
	stats.PowerBonus += len(ctx.Engine.State.Players[ctx.PlayerID].Hand)
}
