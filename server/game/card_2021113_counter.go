package game

func (Card2021113ArcaneBarrierScroll) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnSpellHitBeforeDamage}
}

func (Card2021113ArcaneBarrierScroll) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnSpellHitBeforeDamage && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.Card != nil && isSpellLikeCard(ctx.Event.Card.Card)
}
