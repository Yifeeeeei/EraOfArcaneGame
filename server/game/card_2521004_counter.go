package game

func (Card2521004HolySanctionScroll) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnSpellCast}
}

func (Card2521004HolySanctionScroll) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnSpellCast && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.Card != nil && isSorcerySkill(ctx.Event.Card.Card)
}

type Card2521004HolySanctionScroll struct{ AlwaysActive }

func (Card2521004HolySanctionScroll) ID() string { return "2521004" }

func (Card2521004HolySanctionScroll) Name() string { return "神圣制裁卷轴" }

func (Card2521004HolySanctionScroll) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: "holy_sanction", RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 2})
	return nil
}
