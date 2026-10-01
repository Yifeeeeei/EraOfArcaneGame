package game

type Card2321111CounterWindHoleScroll struct{ AlwaysActive }

func (Card2321111CounterWindHoleScroll) ID() string { return "2321111" }

func (Card2321111CounterWindHoleScroll) Name() string { return "反击风洞卷轴" }

func (Card2321111CounterWindHoleScroll) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnSpellMissOrCancelled}
}

func (Card2321111CounterWindHoleScroll) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnSpellMissOrCancelled && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.Card != nil &&
		isSpellLikeCard(ctx.Event.Card.Card) && ctx.Engine.effectiveSpellArea(ctx.Event.Card) == SpellAreaSingle &&
		len(ctx.Engine.enemyUnits(ctx.Source.OwnerID, false, func(unit *CardInstance) bool { return unit != nil && unit.Card != nil && unit.Card.IsCompanion() })) > 0
}

func (Card2321111CounterWindHoleScroll) ResolveCounter(ctx *CounterContext) { ctx.ReflectSpell() }
