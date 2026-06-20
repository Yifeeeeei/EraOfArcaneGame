package game

type Card3101002AllFiresAsOne struct{ AlwaysActive }

func (Card3101002AllFiresAsOne) ID() string   { return "3101002" }
func (Card3101002AllFiresAsOne) Name() string { return "万火合一术" }

func (Card3101002AllFiresAsOne) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	power := stats.PowerBonus
	if finalPower, ok := ctx.ExtraData["final_power"].(int); ok {
		power = finalPower
	}
	stats.DamageBonus += max(power, 0) / 5
}
