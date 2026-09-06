package game

type Card1321113CouncilMessenger struct{ AlwaysActive }

func (Card1321113CouncilMessenger) ID() string { return "1321113" }

func (Card1321113CouncilMessenger) Name() string { return "议庭传信鸽" }

func (Card1321113CouncilMessenger) OnEnter(ctx *EffectContext) error {
	addGeneratedCardToPlayerHand(ctx, ctx.OpponentID, "2001102")
	return nil
}
