package game

import "eraofarcane/model"

type Card2221004ManesStaff struct{}

func (Card2221004ManesStaff) ID() string   { return "2221004" }
func (Card2221004ManesStaff) Name() string { return "玛涅斯之杖" }

func (Card2221004ManesStaff) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.Target == nil || ctx.Target.Card == nil {
		return
	}
	if ctx.Target.Card.Category == model.ElementWater {
		stats.PowerBonus++
	}
}
