package game

type Card3121015BurningWind struct{ AlwaysActive }

func (Card3121015BurningWind) ID() string { return "3121015" }

func (Card3121015BurningWind) Name() string { return "焚风" }

func (Card3121015BurningWind) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData != nil {
		purpose, _ := ctx.ExtraData["purpose"].(string)
		if !isBoostPurpose(skillPurpose(purpose)) {
			return
		}
		stats.Pierce = true
	}
}
