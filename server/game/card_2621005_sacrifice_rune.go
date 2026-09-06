package game

func (Card2621005SacrificeRune) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnFriendlyDeath, TriggerOnEnemyDeath}
}

func (Card2621005SacrificeRune) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Card != nil && ctx.Event.Card.Card.IsCompanion()
}

type Card2621005SacrificeRune struct{ AlwaysActive }

func (Card2621005SacrificeRune) ID() string { return "2621005" }

func (Card2621005SacrificeRune) Name() string { return "献祭符文" }

func (Card2621005SacrificeRune) OnFriendlyDeath(ctx *EffectContext) error {
	return DrawCards(2)(ctx)
}

func (Card2621005SacrificeRune) OnEnemyDeath(ctx *EffectContext) error {
	return DrawCards(2)(ctx)
}
