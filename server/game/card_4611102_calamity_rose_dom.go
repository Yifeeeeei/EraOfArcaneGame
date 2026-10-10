package game

type Card4611102CalamityRoseDom struct{ AlwaysActive }

func (Card4611102CalamityRoseDom) ID() string { return "4611102" }

func (Card4611102CalamityRoseDom) Name() string { return "灾厄玫瑰 多姆" }

func (Card4611102CalamityRoseDom) OnEnter(ctx *EffectContext) error {
	ctx.Engine.millTopDeckCards(ctx.PlayerID, 4)
	ctx.Engine.millTopDeckCards(ctx.OpponentID, 4)
	return nil
}
