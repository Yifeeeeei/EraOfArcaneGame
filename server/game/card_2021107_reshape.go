package game

type Card2021107Reshape struct{ AlwaysActive }

func (Card2021107Reshape) ID() string { return "2021107" }

func (Card2021107Reshape) Name() string { return "重塑" }

func (Card2021107Reshape) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.runResolution("reshape", func() { ctx.Engine.discardAllHandCards(ctx.PlayerID) }, func() { ctx.Engine.drawCards(ctx.PlayerID, 2) })
	return nil
}
