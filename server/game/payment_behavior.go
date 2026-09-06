package game

// PaymentConstraint describes requirements on the printed/own portion of a
// cost. It is queried during preparation and never spends resources.
type PaymentConstraint struct {
	StrictElements map[string]int
	DistinctOwnUse bool
}

type PaymentConstraintBehavior interface {
	PaymentConstraint(*CardInstance, paymentPurpose, map[string]int) PaymentConstraint
}

type OverexertPaymentEvent struct {
	Units     []*CardInstance
	Payment   map[string]int
	PoolSpent map[string]int
}

type OverexertPaymentBehavior interface {
	HasActiveOverexertPayment(*CardInstance) bool
	OnOverexertPayment(*EffectContext, OverexertPaymentEvent)
}

func (AlwaysActive) HasActiveOverexertPayment(*CardInstance) bool { return true }

func (e *Engine) notifyOverexertPayment(ps *PlayerState, units []*CardInstance, payment, poolSpent map[string]int) {
	if ps == nil || len(units) == 0 {
		return
	}
	event := OverexertPaymentEvent{Units: units, Payment: payment, PoolSpent: poolSpent}
	for _, source := range e.getAllFieldCards(ps) {
		if source == nil || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		if b, ok := cardBehavior(source).(OverexertPaymentBehavior); ok && b.HasActiveOverexertPayment(source) {
			b.OnOverexertPayment(e.skillContext(ps.PlayerID, source), event)
		}
	}
}
