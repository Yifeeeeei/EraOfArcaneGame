package game

type Card3321109SilverleafCyclone struct{ AlwaysActive }

func (Card3321109SilverleafCyclone) ID() string { return "3321109" }

func (Card3321109SilverleafCyclone) Name() string { return "银叶旋风" }

func (Card3321109SilverleafCyclone) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil || !ctx.Engine.State.CardEnteredGraveyardThisTurn {
		return
	}
	stats.PowerBonus = 6
}

var _ SkillContributionModifier = Card3321109SilverleafCyclone{}
