package game

type Card3421008JointCasting struct{}

func (Card3421008JointCasting) ID() string   { return "3421008" }
func (Card3421008JointCasting) Name() string { return "联合施法" }

func (Card3421008JointCasting) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	purpose, _ := ctx.ExtraData["purpose"].(string)
	if !isBoostPurpose(skillPurpose(purpose)) {
		return
	}
	stats.DamageBonus++
}
