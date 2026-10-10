package game

type Card3321105SweepingWind struct{ AlwaysActive }

func (Card3321105SweepingWind) ID() string { return "3321105" }

func (Card3321105SweepingWind) Name() string { return "风卷残云" }

func (Card3321105SweepingWind) DamageScope() DamageScope { return DamageAny }

func (Card3321105SweepingWind) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if ctx == nil || ctx.Engine == nil || event.Target == nil || event.Target.Position == nil || event.Target.CurrentLife != 1 {
		return nil
	}
	ctx.Engine.destroyUnit(event.Target, event.Target.OwnerID)
	return nil
}
