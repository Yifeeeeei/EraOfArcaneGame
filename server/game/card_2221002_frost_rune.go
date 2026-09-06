package game

func (Card2221002FrostRune) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnConsume}
}

func (Card2221002FrostRune) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnConsume && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.Card != nil && ctx.Event.Card.Card.IsCompanion()
}

type Card2221002FrostRune struct{ AlwaysActive }

func (Card2221002FrostRune) ID() string { return "2221002" }

func (Card2221002FrostRune) Name() string { return "冰霜符文" }

func (Card2221002FrostRune) OnConsume(ctx *EffectContext) error {
	if ctx.Target == nil || !ctx.Target.Card.IsCompanion() || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	return ApplyStatusToTarget(StatusFreeze, 1)(ctx)
}
