package game

import "eraofarcane/model"

type Card3221003FreezingCurrent struct{ AlwaysActive }

func (Card3221003FreezingCurrent) ID() string   { return "3221003" }
func (Card3221003FreezingCurrent) Name() string { return "激冻寒流" }

func (Card3221003FreezingCurrent) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	purpose, _ := ctx.ExtraData["purpose"].(string)
	if !isBoostPurpose(skillPurpose(purpose)) {
		return
	}
	if ctx.Target == nil || ctx.Target.Card == nil || ctx.Target.Card.Category != model.ElementWater {
		return
	}
	stats.PowerBonus += 2
}
