package game

type Card2221002FrostRune struct{}

func (Card2221002FrostRune) ID() string   { return "2221002" }
func (Card2221002FrostRune) Name() string { return "冰霜符文" }

func (Card2221002FrostRune) OnConsume(ctx *EffectContext) error {
	if ctx.Target == nil || !ctx.Target.Card.IsCompanion() || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	return ApplyStatusToTarget(StatusFreeze, 1)(ctx)
}
