package game

func (Card2321010IllusionScroll) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnSpellCast}
}

func (Card2321010IllusionScroll) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnSpellCast && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.Card != nil && !isSorcerySkill(ctx.Event.Card.Card) && ctx.Engine.State.PendingSpell != nil
}

type Card2321010IllusionScroll struct{ AlwaysActive }

func (Card2321010IllusionScroll) ID() string { return "2321010" }

func (Card2321010IllusionScroll) Name() string { return "幻术卷轴" }

func (Card2321010IllusionScroll) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.startIllusionScrollRearrange(ctx.PlayerID, ctx.Source, ctx.ExtraData)
	return nil
}
