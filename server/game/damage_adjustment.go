package game

// DamageAdjustmentBehavior calculates incoming damage before shields and
// replacement windows. Implementations must be pure: querying an amount must
// never consume a resource, trigger another effect, or mutate game state.
// Field order is stable (player 0 then 1, each field's documented zone order).
type DamageAdjustmentBehavior interface {
	HasActiveDamageAdjustment(*CardInstance) bool
	DamageAdjustmentScope() DamageScope
	AdjustDamage(*EffectContext, DamageEvent, int) int
}

func (AlwaysActive) HasActiveDamageAdjustment(*CardInstance) bool { return true }

func (e *Engine) adjustedDamage(event DamageEvent) int {
	amount := event.Amount
	if amount <= 0 {
		return 0
	}
	for _, player := range e.State.Players {
		for _, observer := range e.getAllFieldCards(player) {
			if observer == nil || observer.Card == nil || observer.IsSetCounter || e.hasEffectiveStatus(observer, StatusPetrify) {
				continue
			}
			behavior, ok := cardBehavior(observer).(DamageAdjustmentBehavior)
			if !ok || !behavior.HasActiveDamageAdjustment(observer) || !event.Matches(observer, behavior.DamageAdjustmentScope()) {
				continue
			}
			ctx := &EffectContext{Engine: e, Source: observer, Target: event.Target, PlayerID: observer.OwnerID, OpponentID: 1 - observer.OwnerID}
			amount = max(0, behavior.AdjustDamage(ctx, event, amount))
		}
	}
	return amount
}
