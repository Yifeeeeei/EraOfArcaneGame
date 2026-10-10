package game

type Card3421001ForestShelter struct{ AlwaysActive }

func (Card3421001ForestShelter) ID() string { return "3421001" }

func (Card3421001ForestShelter) Name() string { return "森林的庇护" }

func (Card3421001ForestShelter) MasteryMax() int { return 3 }

func (Card3421001ForestShelter) OnMastery(ctx *EffectContext, level int) error {
	return nil
}

func (Card3421001ForestShelter) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData == nil || ctx.ExtraData["purpose"] != string(skillPurposeDefend) {
		return
	}
	switch mastery := ctx.Source.Statuses[StatusMastery]; {
	case mastery >= 3:
		stats.PowerBonus = 6
	case mastery >= 1:
		stats.PowerBonus = 4
	}
}
