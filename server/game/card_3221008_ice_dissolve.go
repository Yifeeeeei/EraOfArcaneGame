package game

type Card3221008IceDissolve struct{ AlwaysActive }

func (Card3221008IceDissolve) ID() string { return "3221008" }

func (Card3221008IceDissolve) Name() string { return "冰封消解" }

func (Card3221008IceDissolve) CanReactToSpell(ctx *EffectContext, spell *SpellCast) bool {
	return ctx != nil && spell != nil && spell.AttackerID != ctx.PlayerID && spell.TotalPower > 0
}

func (Card3221008IceDissolve) OnSpellReaction(ctx *EffectContext, spell *SpellCast) error {
	sources := positiveSpellPowerSources(spell)
	if len(sources) == 0 {
		return nil
	}
	if len(sources) == 1 {
		applyIceDissolveToSource(ctx, spell, sources[0].InstanceID)
		return nil
	}
	candidates := make([]map[string]any, 0, len(sources))
	for _, source := range sources {
		candidates = append(candidates, map[string]any{
			"instance_id": source.InstanceID,
			"name":        source.CardName,
			"power":       source.Power,
			"is_main":     source.IsMain,
		})
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "ice_dissolve_power_source", "冰封消解:选择1个法术威力变为0", candidates, 1, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		applyIceDissolveToSource(ctx, spell, selected[0])
	})
	return nil
}
