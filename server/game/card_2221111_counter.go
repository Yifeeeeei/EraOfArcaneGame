package game

func (Card2221111IceSoulSealForge) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnSpellCast}
}

func (Card2221111IceSoulSealForge) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnSpellCast && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.Card != nil && ctx.Event.Power > 10
}

type Card2221111IceSoulSealForge struct{ AlwaysActive }

func (Card2221111IceSoulSealForge) ID() string { return "2221111" }

func (Card2221111IceSoulSealForge) Name() string { return "冰魄印 淬" }

func (Card2221111IceSoulSealForge) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	original := spellPowerFromData(ctx.ExtraData)
	if original <= 10 {
		return nil
	}
	reduced := (original + 1) / 2
	if ctx.Engine.State.PendingSpell != nil && ctx.Engine.State.PendingSpell.Skill == ctx.Target {
		ctx.Engine.State.PendingSpell.TotalPower = reduced
	}
	if ctx.ExtraData != nil {
		ctx.ExtraData["power"] = reduced
	}
	ctx.Engine.emit(GameEvent{
		Type:   "spell_power_reduced",
		Player: -1,
		Data: map[string]any{
			"player":   ctx.PlayerID,
			"source":   cardToInfo(ctx.Source),
			"spell":    cardToInfo(ctx.Target),
			"original": original,
			"power":    reduced,
		},
	})
	return nil
}
