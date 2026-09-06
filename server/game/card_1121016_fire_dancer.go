package game

type Card1121016FireDancer struct{ AlwaysActive }

func (Card1121016FireDancer) ID() string { return "1121016" }

func (Card1121016FireDancer) Name() string { return "舞火者" }

func (Card1121016FireDancer) OnEnter(ctx *EffectContext) error {
	grantFireNegativeStatusImmunity(ctx)
	return nil
}

func (Card1121016FireDancer) OnDeath(ctx *EffectContext) error {
	grantFireNegativeStatusImmunity(ctx)
	return nil
}
