package game

func (Card2021114GuardianRune) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerBeforeDamage}
}

func (Card2021114GuardianRune) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerBeforeDamage && ctx.Event.PlayerID == ctx.Source.OwnerID && ctx.Event.Card != nil &&
		ctx.Event.IsLethal()
}

func (Card2021114GuardianRune) ResolveCounter(ctx *CounterContext) {
	if ctx.Event.PreventDamage == nil || !ctx.Event.IsLethal() {
		return
	}
	*ctx.Event.PreventDamage = true
	ctx.Engine.emit(GameEvent{Type: "damage_prevented", Player: -1, Data: map[string]any{
		"source": cardToInfo(ctx.Source), "target": cardToInfo(ctx.Event.Card), "amount": ctx.Event.Damage, "reason": "guardian_rune",
	}})
}

type Card2021114GuardianRune struct{ AlwaysActive }

func (Card2021114GuardianRune) ID() string { return "2021114" }

func (Card2021114GuardianRune) Name() string { return "神护符文" }
