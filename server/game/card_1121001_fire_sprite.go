package game

type Card1121001FireSprite struct{ AlwaysActive }

func (Card1121001FireSprite) ID() string   { return "1121001" }
func (Card1121001FireSprite) Name() string { return "火焰精灵" }

func (Card1121001FireSprite) OnConsume(ctx *EffectContext) error {
	if ctx.Source != nil {
		return ApplyStatusToSelf(StatusBurn, 1)(ctx)
	}
	return nil
}
