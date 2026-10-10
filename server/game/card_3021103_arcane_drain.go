package game

type Card3021103ArcaneDrain struct{ AlwaysActive }

func (Card3021103ArcaneDrain) ID() string { return "3021103" }

func (Card3021103ArcaneDrain) Name() string { return "奥能汲取" }

func (Card3021103ArcaneDrain) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	ctx.Engine.drawCards(ctx.PlayerID, 2)
	return nil
}

var _ OnSpellCastBehavior = Card3021103ArcaneDrain{}

func (Card3021103ArcaneDrain) PaymentConstraint(_ *CardInstance, purpose paymentPurpose, _ map[string]int) PaymentConstraint {
	return PaymentConstraint{DistinctOwnUse: purpose == paymentPurposeUse}
}
