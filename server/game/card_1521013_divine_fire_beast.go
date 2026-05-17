package game

type Card1521013DivineFireBeast struct{ AlwaysActive }

func (Card1521013DivineFireBeast) ID() string   { return "1521013" }
func (Card1521013DivineFireBeast) Name() string { return "神火兽" }

func (Card1521013DivineFireBeast) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	stats.PowerBonus += 2
}
