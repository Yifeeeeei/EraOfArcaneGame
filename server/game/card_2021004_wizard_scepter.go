package game

type Card2021004WizardScepter struct{ AlwaysActive }

func (Card2021004WizardScepter) ID() string { return "2021004" }

func (Card2021004WizardScepter) Name() string { return "巫师权杖" }

func (Card2021004WizardScepter) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	stats.PowerBonus++
}
