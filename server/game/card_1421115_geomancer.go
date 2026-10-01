package game

type Card1421115Geomancer struct{ AlwaysActive }

func (Card1421115Geomancer) ID() string { return "1421115" }

func (Card1421115Geomancer) Name() string { return "地卜行者" }

func (Card1421115Geomancer) OnEnter(ctx *EffectContext) error {
	ctx.Engine.drawCards(ctx.PlayerID, 1)
	return nil
}
