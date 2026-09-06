package game

import "eraofarcane/model"

type Card1021112ArcanePureBody struct{ AlwaysActive }

func (Card1021112ArcanePureBody) ID() string   { return "1021112" }
func (Card1021112ArcanePureBody) Name() string { return "奥术纯净体" }
func (Card1021112ArcanePureBody) PaymentConstraint(_ *CardInstance, purpose paymentPurpose, cost map[string]int) PaymentConstraint {
	if purpose == paymentPurposePlay {
		return PaymentConstraint{StrictElements: map[string]int{model.ElementArcane: totalElementCost(cost)}}
	}
	return PaymentConstraint{}
}
