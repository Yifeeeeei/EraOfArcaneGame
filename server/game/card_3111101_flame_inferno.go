package game

type Card3111101FlameInferno struct{ AlwaysActive }

func (Card3111101FlameInferno) ID() string { return "3111101" }

func (Card3111101FlameInferno) Name() string { return "火炎炼狱" }

func (Card3111101FlameInferno) OnSpellHit(ctx *EffectContext) error {
	if !isOwnSpellHit(ctx) {
		return nil
	}
	if ctx == nil || ctx.Engine == nil || ctx.ExtraData == nil {
		return nil
	}
	power := intFromData(ctx.ExtraData, "power", 0)
	burn := max(power, 0) / 8
	if burn <= 0 {
		return nil
	}
	for _, unit := range spellHitAffectedUnitsFromData(ctx) {
		if unit == nil || unit.CurrentLife <= 0 {
			continue
		}
		ctx.Engine.addStatus(unit, StatusBurn, burn)
	}
	return nil
}

func (Card3111101FlameInferno) OnDefend(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.ExtraData["attack_skill"] != ctx.Source || !boolFromData(ctx.ExtraData, "defense_success") {
		return nil
	}
	power := intFromData(ctx.ExtraData, "attack_power", 0)
	defense := intFromData(ctx.ExtraData, "defense_power", 0)
	burn := max(0, power/8-max(0, defense-power)/4)
	spell := ctx.Engine.State.PendingSpell
	if burn == 0 || spell == nil {
		return nil
	}
	for _, unit := range ctx.Engine.spellAffectedUnits(ctx.OpponentID, ctx.Source, spell.Target) {
		ctx.Engine.addStatus(unit, StatusBurn, burn)
	}
	return nil
}
