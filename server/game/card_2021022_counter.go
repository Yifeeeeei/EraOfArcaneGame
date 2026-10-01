package game

func (Card2021022CounterRune) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnUseItem}
}

func (Card2021022CounterRune) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnUseItem && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.Card != nil && ctx.Event.Card.Card != nil && ctx.Event.Card.Card.IsItem() && (hasCardTag(ctx.Event.Card.Card, "卷轴") || hasCardTag(ctx.Event.Card.Card, "符文"))
}

func (Card2021022CounterRune) ResolveCounter(ctx *CounterContext) { ctx.CancelItem() }

type Card2021022CounterRune struct{ AlwaysActive }

func (Card2021022CounterRune) ID() string { return "2021022" }

func (Card2021022CounterRune) Name() string { return "反制符文" }

func (Card2021022CounterRune) OnUseItem(ctx *EffectContext) error {
	emitBatchEffect(ctx, "counter_rune_ready")
	return nil
}
