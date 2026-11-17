package game

type Card3521013Moonlight struct{ AlwaysActive }

func (Card3521013Moonlight) ID() string { return "3521013" }

func (Card3521013Moonlight) Name() string { return "月之辉" }

func (Card3521013Moonlight) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	purpose, _ := ctx.ExtraData["purpose"].(string)
	if skillPurpose(purpose) != skillPurposeDefend && skillPurpose(purpose) != skillPurposeDefenseBoost {
		return
	}
	stats.PowerBonus += 2
}
