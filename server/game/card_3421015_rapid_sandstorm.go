package game

type Card3421015RapidSandstorm struct{ AlwaysActive }

func (Card3421015RapidSandstorm) ID() string   { return "3421015" }
func (Card3421015RapidSandstorm) Name() string { return "急袭沙暴" }

func (Card3421015RapidSandstorm) HasActiveSpellStatModifier(card *CardInstance) bool {
	return abilityDurationActive(card)
}

func (Card3421015RapidSandstorm) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.Target == nil || ctx.Target.Card == nil || !ctx.Target.Card.IsSkill() {
		return
	}
	if ctx.Target.Card.Power >= 5 {
		return
	}
	stats.PowerBonus -= 2
	stats.DamageBonus -= 2
}
