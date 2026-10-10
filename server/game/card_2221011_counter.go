package game

func (Card2221011RainOfGrace) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnDamaged}
}

func (Card2221011RainOfGrace) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnDamaged && ctx.Event.PlayerID == ctx.Source.OwnerID
}

type Card2221011RainOfGrace struct{ AlwaysActive }

func (Card2221011RainOfGrace) ID() string { return "2221011" }

func (Card2221011RainOfGrace) Name() string { return "恩惠之雨" }

func (Card2221011RainOfGrace) OnUseItem(ctx *EffectContext) error {
	for _, unit := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		ctx.Engine.healUnit(unit, 2, ctx.Source)
	}
	return nil
}
