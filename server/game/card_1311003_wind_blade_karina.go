package game

import (
	"eraofarcane/model"
)

type Card1311003WindBladeKarina struct{ AlwaysActive }

func (Card1311003WindBladeKarina) ID() string { return "1311003" }

func (Card1311003WindBladeKarina) Name() string { return "\"风刃\" 卡琳娜" }

func (Card1311003WindBladeKarina) ModifySkillUseCost(ctx *EffectContext, cost map[string]int) {
	if ctx.Source != nil && ctx.Source.Card.Category == model.ElementAir && skillNeedsTargetInstance(ctx.Source) && !cardHasPierce(ctx.Source) {
		cost[model.ElementAir]++
	}
}

func (Card1311003WindBladeKarina) SpellTargetGrant(_ *EffectContext, skill *CardInstance, _ SpellTarget) SpellTargetGrant {
	return SpellTargetGrant{Pierce: skill != nil && skill.Card != nil && skill.Card.Category == model.ElementAir && skillNeedsTargetInstance(skill)}
}
