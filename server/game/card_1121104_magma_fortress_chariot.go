package game

type Card1121104MagmaFortressChariot struct{ AlwaysActive }

func (Card1121104MagmaFortressChariot) ID() string   { return "1121104" }
func (Card1121104MagmaFortressChariot) Name() string { return "熔岩堡战车" }

func (Card1121104MagmaFortressChariot) OnAttack(ctx *EffectContext) error {
	if ctx.Target == nil {
		return nil
	}
	ctx.Engine.addStatus(ctx.Target, StatusBurn, 1)
	return nil
}

var _ OnAttackBehavior = Card1121104MagmaFortressChariot{}
