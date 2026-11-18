package game

import "eraofarcane/model"

type Card2121105FlameArrayScroll struct{ AlwaysActive }

func (Card2121105FlameArrayScroll) ID() string                               { return "2121105" }
func (Card2121105FlameArrayScroll) Name() string                             { return "魔焰阵卷轴" }
func (Card2121105FlameArrayScroll) HasActiveSpellTapCost(*CardInstance) bool { return true }

func (Card2121105FlameArrayScroll) SpellTapCostBonus(ctx *EffectContext, target *CardInstance) int {
	if target == nil || target.Card == nil || target.OwnerID != ctx.PlayerID ||
		!target.Card.IsCompanion() || target.Card.Category != model.ElementFire || !ctx.Engine.canConsumeCard(target) {
		return 0
	}
	if cost := totalElementCost(target.Card.ElementsCost); cost > 4 {
		return cost
	}
	return 0
}

func (Card2121105FlameArrayScroll) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	stats.PowerBonus += ctx.Source.spellTapPowerBonus
}

func (Card2121105FlameArrayScroll) ClearSkillUseState(card *CardInstance) {
	card.spellTapPowerBonus = 0
}
