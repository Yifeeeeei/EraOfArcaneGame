package game

func (Card2521002ShelterRune) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnSpellHitBeforeDamage}
}

func (Card2521002ShelterRune) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnSpellHitBeforeDamage && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.Card != nil && !isSorcerySkill(ctx.Event.Card.Card) && ctx.Event.Power < 10
}

func (Card2521002ShelterRune) ResolveCounter(ctx *CounterContext) { ctx.CancelHit() }

type Card2521002ShelterRune struct{ AlwaysActive }

func (Card2521002ShelterRune) ID() string { return "2521002" }

func (Card2521002ShelterRune) Name() string { return "庇护符文" }

func (Card2521002ShelterRune) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: "shelter_rune", Amount: 10, RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 2})
	return nil
}
