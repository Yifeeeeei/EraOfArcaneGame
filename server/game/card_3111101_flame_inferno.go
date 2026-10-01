package game

type Card3111101FlameInferno struct{ AlwaysActive }

func (Card3111101FlameInferno) ID() string { return "3111101" }

func (Card3111101FlameInferno) Name() string { return "火炎炼狱" }

func (Card3111101FlameInferno) OnSpellHit(ctx *EffectContext) error {
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
		unit.Statuses[StatusBurn] += burn
	}
	return nil
}

func (Card3111101FlameInferno) RequiredDefensePower(_ *CardInstance, power int) int {
	return power + max(power, 0)/8*4
}
