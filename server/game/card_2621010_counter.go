package game

import ()

func (Card2621010DragIntoAbyss) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnFriendlyDeath}
}

func (Card2621010DragIntoAbyss) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnFriendlyDeath && ctx.Event.PlayerID == ctx.Source.OwnerID && ctx.Event.Card != nil && ctx.Event.Card.DamageTakenThisTurn > 0
}

type Card2621010DragIntoAbyss struct{ AlwaysActive }

func (Card2621010DragIntoAbyss) ID() string { return "2621010" }

func (Card2621010DragIntoAbyss) Name() string { return "拖入深渊" }

func (Card2621010DragIntoAbyss) OnUseItem(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.DamageTakenThisTurn <= 0 {
		return nil
	}
	damage := ctx.Target.DamageTakenThisTurn
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	ctx.Engine.SetPendingAction(ctx.PlayerID, "drag_into_abyss_target", "Drag Into Abyss: choose an enemy unit", candidates, 1, 1, func(selected []string) {
		target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
		if target != nil {
			ctx.DealDamage(target, damage)
		}
	})
	return nil
}
