package game

type Card3421014ThousandMileQuicksand struct{ AlwaysActive }

func (Card3421014ThousandMileQuicksand) ID() string { return "3421014" }

func (Card3421014ThousandMileQuicksand) Name() string { return "千里流沙" }

func (Card3421014ThousandMileQuicksand) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	stats.PowerBonus += 1
}
