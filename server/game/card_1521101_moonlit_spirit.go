package game

const moonlitSpiritAuraSpentStatus = "月霞光环已失效"

type Card1521101MoonlitSpirit struct{ AlwaysActive }

func (Card1521101MoonlitSpirit) ID() string   { return "1521101" }
func (Card1521101MoonlitSpirit) Name() string { return "月霞之灵" }

func (Card1521101MoonlitSpirit) HasActiveSpellStatModifier(card *CardInstance) bool {
	return card != nil && card.Statuses[moonlitSpiritAuraSpentStatus] <= 0
}

func (Card1521101MoonlitSpirit) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Target == nil || ctx.Target.Card == nil || !hasCardTag(ctx.Target.Card, "法术") {
		return
	}
	stats.PowerBonus += 2
}

func (Card1521101MoonlitSpirit) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || !isFriendlySpellCast(ctx) {
		return nil
	}
	attacker, _ := ctx.ExtraData["attacker"].(int)
	if attacker != ctx.PlayerID {
		return nil
	}
	ctx.Source.Statuses[moonlitSpiritAuraSpentStatus] = 1
	return nil
}

var _ SpellStatModifier = Card1521101MoonlitSpirit{}
var _ OnSpellCastBehavior = Card1521101MoonlitSpirit{}
