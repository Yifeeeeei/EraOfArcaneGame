package game

type Card1321107SkyCityThief struct{ AlwaysActive }

func (Card1321107SkyCityThief) ID() string { return "1321107" }

func (Card1321107SkyCityThief) Name() string { return "云霄城大盗" }

func (Card1321107SkyCityThief) OnEnter(ctx *EffectContext) error {
	ctx.Engine.discardRandomHandCard(ctx.PlayerID)
	ctx.Engine.discardRandomHandCard(ctx.OpponentID)
	return nil
}
