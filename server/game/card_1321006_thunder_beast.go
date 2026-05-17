package game

import "eraofarcane/model"

type Card1321006ThunderBeast struct{ AlwaysActive }

func (Card1321006ThunderBeast) ID() string   { return "1321006" }
func (Card1321006ThunderBeast) Name() string { return "雷霆兽" }

func (Card1321006ThunderBeast) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.Target == nil || ctx.Target.Card == nil {
		return
	}
	if ctx.Target.Card.Category == model.ElementAir {
		stats.DamageBonus++
	}
}
