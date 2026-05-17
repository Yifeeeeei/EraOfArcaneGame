package game

type Card1121014Firethorn struct{ AlwaysActive }

func (Card1121014Firethorn) ID() string   { return "1121014" }
func (Card1121014Firethorn) Name() string { return "火荆" }
func (Card1121014Firethorn) OnDeath(ctx *EffectContext) error {
	return ApplyStatusAuto(StatusBurn, 1)(ctx)
}
