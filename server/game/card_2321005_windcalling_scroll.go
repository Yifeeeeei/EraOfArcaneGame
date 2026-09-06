package game

type Card2321005WindcallingScroll struct{ AlwaysActive }

func (Card2321005WindcallingScroll) ID() string { return "2321005" }

func (Card2321005WindcallingScroll) Name() string { return "唤风卷轴" }

func (Card2321005WindcallingScroll) OnUseItem(ctx *EffectContext) error {
	if err := DrawCards(2)(ctx); err != nil {
		return err
	}
	ctx.Engine.State.Players[ctx.PlayerID].SkipNextDraw = true
	return nil
}
