package game

// LethalSacrificeBehavior supplies eligible replacement costs. The engine owns
// offering the choice, revalidating it, paying the sacrifice and resuming the
// already-calculated damage when the owner declines.
type LethalSacrificeBehavior interface {
	HasActiveLethalSacrifice(*CardInstance) bool
	CanPreventLethalBySacrifice(*EffectContext, DamageEvent) bool
}

func (AlwaysActive) HasActiveLethalSacrifice(*CardInstance) bool { return true }

func (e *Engine) canPreventLethalBySacrifice(source, target *CardInstance, amount int, data map[string]any) bool {
	if source == nil || source.Card == nil || e.hasEffectiveStatus(source, StatusPetrify) {
		return false
	}
	behavior, ok := cardBehavior(source).(LethalSacrificeBehavior)
	if !ok || !behavior.HasActiveLethalSacrifice(source) {
		return false
	}
	ctx := &EffectContext{Engine: e, Source: source, Target: target, PlayerID: source.OwnerID,
		OpponentID: 1 - source.OwnerID, ExtraData: data}
	event := damageEventFromContext(ctx)
	event.Amount = amount
	return behavior.CanPreventLethalBySacrifice(ctx, event)
}
