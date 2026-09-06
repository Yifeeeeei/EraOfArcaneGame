package game

import (
	"eraofarcane/model"
)

func (Card2321108ScatterAway) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnDamaged}
}

func (Card2321108ScatterAway) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnDamaged && ctx.Event.PlayerID == ctx.Source.OwnerID && ctx.Event.Card != nil &&
		ctx.Event.Card.Card != nil && ctx.Event.Card.Card.IsCompanion() &&
		ctx.Event.Card.Card.Category == model.ElementAir &&
		ctx.Event.Damage > 0
}

type Card2321108ScatterAway struct{ AlwaysActive }

func (Card2321108ScatterAway) ID() string { return "2321108" }

func (Card2321108ScatterAway) Name() string { return "散去" }

func (Card2321108ScatterAway) DamageScope() DamageScope { return DamageFriendly }

func (Card2321108ScatterAway) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if ctx == nil || ctx.Engine == nil || event.Target == nil || event.Target.Card == nil {
		return nil
	}
	if event.Target.OwnerID != ctx.PlayerID || !event.Target.Card.IsCompanion() || event.Target.Card.Category != model.ElementAir || event.Amount <= 0 {
		return nil
	}
	event.Target.Statuses[temporaryDamageAndNegativeImmunityUntilStatus] = ctx.Engine.State.TurnNumber + 1
	ctx.Engine.emit(GameEvent{
		Type:   "scatter_away_immunity",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(event.Target),
		},
	})
	return nil
}

var _ OnDamagedBehavior = Card2321108ScatterAway{}
