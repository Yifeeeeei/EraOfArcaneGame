package game

type Card2121002FireRune struct{ AlwaysActive }

func (Card2121002FireRune) ID() string   { return "2121002" }
func (Card2121002FireRune) Name() string { return "火焰符文" }

func (Card2121002FireRune) OnConsume(ctx *EffectContext) error {
	if ctx.Target == nil || (!ctx.Target.Card.IsHero() && !ctx.Target.Card.IsCompanion()) {
		return nil
	}
	return ApplyStatusToTarget(StatusBurn, 1)(ctx)
}
