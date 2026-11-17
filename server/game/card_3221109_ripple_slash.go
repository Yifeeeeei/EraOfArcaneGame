package game

type Card3221109RippleSlash struct{ AlwaysActive }

func (Card3221109RippleSlash) ID() string { return "3221109" }

func (Card3221109RippleSlash) HasSpellSourceModifier() {}

func (Card3221109RippleSlash) Name() string { return "波纹斩" }

func (Card3221109RippleSlash) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if !rippleSlashAuraActive(ctx) {
		return
	}
	stats.PowerBonus += 2
}

func (Card3221109RippleSlash) ModifySpellArea(ctx *EffectContext, area *SpellArea) {
	if !rippleSlashAuraActive(ctx) || area == nil || *area != SpellAreaSingle {
		return
	}
	*area = SpellAreaFrontRow
}

func rippleSlashAuraActive(ctx *EffectContext) bool {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Source != ctx.Target {
		return false
	}
	// The pending cast has already incremented the turn counter. It must not
	// grant its own aura while targets are previewed or damage is resolved.
	if spell := ctx.Engine.State.PendingSpell; spell != nil && spell.Skill == ctx.Source && spell.priorSameNumberCasts != nil {
		return *spell.priorSameNumberCasts > 0
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	return spellCastByNumberThisTurn(ps, "3221109") > 0
}

var _ SpellStatModifier = Card3221109RippleSlash{}

var _ SpellAreaModifier = Card3221109RippleSlash{}
