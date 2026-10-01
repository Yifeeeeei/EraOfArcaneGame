package game

type Card1321114CouncilExecutor struct{ AlwaysActive }

func (Card1321114CouncilExecutor) ID() string { return "1321114" }

func (Card1321114CouncilExecutor) Name() string { return "议庭执行者" }

func (Card1321114CouncilExecutor) OnEnter(ctx *EffectContext) error {
	first := ctx.Engine.discardRandomHandCard(ctx.OpponentID)
	if first != nil && first.Card != nil && first.Card.Number == "2001102" {
		ctx.Engine.discardRandomHandCard(ctx.OpponentID)
	}
	return nil
}
