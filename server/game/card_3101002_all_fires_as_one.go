package game

type Card3101002AllFiresAsOne struct{}

func (Card3101002AllFiresAsOne) ID() string   { return "3101002" }
func (Card3101002AllFiresAsOne) Name() string { return "万火合一术" }

func (Card3101002AllFiresAsOne) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	stats.DamageBonus += max(stats.PowerBonus, 0) / 5
}
