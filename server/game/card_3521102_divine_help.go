package game

type Card3521102DivineHelp struct{ AlwaysActive }

func (Card3521102DivineHelp) ID() string { return "3521102" }

func (Card3521102DivineHelp) Name() string { return "神助" }

func (Card3521102DivineHelp) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	purpose, _ := ctx.ExtraData["purpose"].(string)
	if !isBoostPurpose(skillPurpose(purpose)) {
		return
	}
	if ctx.Target == nil || ctx.Target.Card == nil || !hasCardTag(ctx.Target.Card, "神秘") {
		return
	}
	stats.PowerBonus += 2
}
