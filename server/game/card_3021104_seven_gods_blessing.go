package game

type Card3021104SevenGodsBlessing struct{ AlwaysActive }

func (Card3021104SevenGodsBlessing) ID() string   { return "3021104" }
func (Card3021104SevenGodsBlessing) Name() string { return "七神加护" }

func (Card3021104SevenGodsBlessing) ModifySkillUseCost(ctx *EffectContext, cost map[string]int) {
	if !allLearnedSkillCategoriesDifferent(ctx) {
		return
	}
	if ctx.Source != nil && ctx.Source.Card != nil {
		reduceGenericCost(cost, ctx.Source.Card.Category, 1)
	}
}

func (Card3021104SevenGodsBlessing) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if !allLearnedSkillCategoriesDifferent(ctx) || ctx.Target == nil || ctx.Target.Card == nil || !hasCardTag(ctx.Target.Card, "法术") {
		return
	}
	stats.PowerBonus += 2
}

func allLearnedSkillCategoriesDifferent(ctx *EffectContext) bool {
	if ctx == nil || ctx.Engine == nil {
		return false
	}
	seen := map[string]bool{}
	for _, skill := range ctx.Engine.State.Players[ctx.PlayerID].Skills {
		if skill == nil || skill.Card == nil {
			continue
		}
		category := skill.Card.Category
		if seen[category] {
			return false
		}
		seen[category] = true
	}
	return true
}

var _ SkillUseCostModifier = Card3021104SevenGodsBlessing{}
var _ SpellStatModifier = Card3021104SevenGodsBlessing{}
