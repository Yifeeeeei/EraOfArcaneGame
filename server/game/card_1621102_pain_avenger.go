package game

type Card1621102PainAvenger struct{ AlwaysActive }

func (Card1621102PainAvenger) ID() string { return "1621102" }

func (Card1621102PainAvenger) Name() string { return "苦痛复仇者" }

func (Card1621102PainAvenger) DamageScope() DamageScope { return DamageSelf }

func (Card1621102PainAvenger) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if !useTriggeredTurn(ctx.Source) {
		return nil
	}
	ctx.Source.CurrentAttack++
	return nil
}
