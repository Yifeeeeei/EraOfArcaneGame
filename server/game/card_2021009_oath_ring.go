package game

type Card2021009OathRing struct{ AlwaysActive }

func (Card2021009OathRing) ID() string   { return "2021009" }
func (Card2021009OathRing) Name() string { return "誓约之戒" }

func (Card2021009OathRing) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	stats.PowerBonus -= 2
}
