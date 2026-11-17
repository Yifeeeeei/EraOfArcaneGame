package game

func (Card2521110AngelPrayer) AllowsNormalCounterUse() bool { return true }

func (Card2521110AngelPrayer) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnSpellCast}
}

func (Card2521110AngelPrayer) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnSpellCast && ctx.Event.PlayerID != ctx.Source.OwnerID &&
		!ctx.Event.BoostUse && (ctx.Event.Purpose == "" || ctx.Event.Purpose == skillPurposeAttack) &&
		ctx.Event.Card != nil && isCombatSpellCard(ctx.Event.Card.Card)
}

func (Card2521110AngelPrayer) ModifyCounterCost(_ *CounterContext, cost map[string]int) {
	clear(cost)
}

func (Card2521110AngelPrayer) ResolveCounter(ctx *CounterContext) {
	ctx.Engine.flipDeckMatchesToHandThen(ctx.Source.OwnerID, 1, 0, isLightSpirit, func(drawn []*CardInstance) {
		for _, card := range drawn {
			card.Statuses["入场费用光-1"]++
		}
	})
}
