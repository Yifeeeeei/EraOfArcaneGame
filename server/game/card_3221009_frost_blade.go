package game

type Card3221009FrostBlade struct{ AlwaysActive }

func (Card3221009FrostBlade) ID() string   { return "3221009" }
func (Card3221009FrostBlade) Name() string { return "冰霜利刃" }

func (Card3221009FrostBlade) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	purpose, _ := ctx.ExtraData["purpose"].(string)
	if skillPurpose(purpose) != skillPurposeAttack && skillPurpose(purpose) != skillPurposeAttackBoost {
		return
	}
	stats.PowerBonus += 2
}
