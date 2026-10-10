package game

import (
	"eraofarcane/model"
)

type Card3021101ArcaneImpact struct{ AlwaysActive }

func (Card3021101ArcaneImpact) ID() string { return "3021101" }

func (Card3021101ArcaneImpact) Name() string { return "奥术冲击" }

func (Card3021101ArcaneImpact) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Source == nil || ctx.Source.Card == nil {
		return
	}
	if !costUsesOnlyArcane(ctx.Source.Card.ElementsCost) || !costUsesOnlyArcane(ctx.Source.Card.ElementsExpense) {
		return
	}
	stats.PowerBonus++
	stats.DamageBonus++
}

func costUsesOnlyArcane(cost map[string]int) bool {
	seen := false
	for elem, amount := range cost {
		if amount <= 0 {
			continue
		}
		if elem != model.ElementArcane {
			return false
		}
		seen = true
	}
	return seen
}

var _ SkillContributionModifier = Card3021101ArcaneImpact{}
