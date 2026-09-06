package game

type Card2221112IceSoulSealCancel struct{ AlwaysActive }

func (Card2221112IceSoulSealCancel) ID() string { return "2221112" }

func (Card2221112IceSoulSealCancel) Name() string { return "冰魄印 消" }

func (Card2221112IceSoulSealCancel) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnSpellCast}
}

func (Card2221112IceSoulSealCancel) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnSpellCast && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.Card != nil &&
		ctx.Event.BoostUse && ctx.Event.Power < 5
}

func (Card2221112IceSoulSealCancel) ResolveCounter(ctx *CounterContext) { ctx.CancelBoost() }
