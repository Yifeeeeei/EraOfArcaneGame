package game

type Card3121002Burn struct{ AlwaysActive }

func (Card3121002Burn) ID() string   { return "3121002" }
func (Card3121002Burn) Name() string { return "焚烧" }

func (Card3121002Burn) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	target, _ := ctx.ExtraData["spell_target_unit"].(*CardInstance)
	if target == nil || target.Statuses[StatusBurn] <= 0 {
		return
	}
	stats.PowerBonus += 2
}
