package game

func (Card2621011FrenzyRune) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnConsume}
}

func (Card2621011FrenzyRune) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnConsume && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.Card != nil && ctx.Event.Card.Card.IsCompanion() && ctx.Event.Card.CurrentAttack > 0 && len(ctx.Engine.adjacentUnitCandidatesForCounter(ctx.Source.OwnerID, ctx.Event.Card)) > 0
}

type Card2621011FrenzyRune struct{ AlwaysActive }

func (Card2621011FrenzyRune) ID() string { return "2621011" }

func (Card2621011FrenzyRune) Name() string { return "狂乱符文" }

func (Card2621011FrenzyRune) OnUseItem(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.CurrentAttack <= 0 {
		return nil
	}
	attacker := ctx.Target
	candidates := ctx.Engine.adjacentUnitCandidatesForCounter(ctx.PlayerID, attacker)
	ctx.Engine.SetPendingAction(ctx.PlayerID, "frenzy_rune_target", "Frenzy Rune: choose an adjacent unit", candidates, 1, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[attacker.OwnerID], selected[0])
		ctx.Engine.resolveForcedUnitAttack(attacker.OwnerID, attacker, target, "frenzy_rune")
	})
	return nil
}
