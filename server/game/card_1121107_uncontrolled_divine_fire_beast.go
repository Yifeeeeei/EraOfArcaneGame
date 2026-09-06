package game

type Card1121107UncontrolledDivineFireBeast struct{ AlwaysActive }

func (Card1121107UncontrolledDivineFireBeast) ID() string { return "1121107" }

func (Card1121107UncontrolledDivineFireBeast) Name() string { return "失控的神火兽" }

func (Card1121107UncontrolledDivineFireBeast) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	stats.PowerBonus += 2
}

func (Card1121107UncontrolledDivineFireBeast) ModifyEnemySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	stats.PowerBonus += 2
}
