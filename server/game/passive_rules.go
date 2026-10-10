package game

// Pure passive rule queries. Availability is evaluated against the instance;
// providers must not consume marks, uses or resources while answering a query.
type ShieldDecayBehavior interface {
	HasActiveShieldDecay(*CardInstance) bool
	PreventsShieldDecay(*CardInstance, int) bool
}
type FieldStatusSuppressionBehavior interface {
	HasActiveFieldStatusSuppression(*CardInstance) bool
	SuppressesFieldStatus(*EffectContext, *CardInstance, string) bool
}
type IntrinsicLoadBehavior interface {
	HasActiveIntrinsicLoad(*CardInstance) bool
	ModifyIntrinsicLoad(*CardInstance, bool, map[string]int)
}
type DeckSearchPermissionBehavior interface {
	CanBeFlippedOrSearched(*CardInstance) bool
}

func (AlwaysActive) HasActiveShieldDecay(*CardInstance) bool            { return true }
func (AlwaysActive) HasActiveFieldStatusSuppression(*CardInstance) bool { return true }
func (AlwaysActive) HasActiveIntrinsicLoad(*CardInstance) bool          { return true }

func (e *Engine) playerShieldDecayPrevented(ps *PlayerState) bool {
	if ps == nil {
		return false
	}
	for _, source := range e.getAllFieldCards(ps) {
		if source == nil || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		b, ok := cardBehavior(source).(ShieldDecayBehavior)
		if ok && b.HasActiveShieldDecay(source) && b.PreventsShieldDecay(source, ps.Shield) {
			return true
		}
	}
	return false
}

func (e *Engine) fieldSuppressesStatus(card *CardInstance, status string) bool {
	if card == nil || card.Position == nil {
		return false
	}
	for _, source := range e.getAllFieldCards(e.State.Players[card.OwnerID]) {
		// Use raw petrification here: asking whether petrification is suppressed by
		// a provider must not recursively ask that provider the same question.
		if source == nil || source.Statuses[StatusPetrify] > 0 {
			continue
		}
		b, ok := cardBehavior(source).(FieldStatusSuppressionBehavior)
		if ok && b.HasActiveFieldStatusSuppression(source) && b.SuppressesFieldStatus(e.skillContext(card.OwnerID, source), card, status) {
			return true
		}
	}
	return false
}
