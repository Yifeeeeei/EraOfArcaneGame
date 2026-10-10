package game

type Card1321112JiuxiaoContact struct{ AlwaysActive }

func (Card1321112JiuxiaoContact) ID() string { return "1321112" }

func (Card1321112JiuxiaoContact) Name() string { return "九霄接头人" }

func (Card1321112JiuxiaoContact) IsPrayerAbility() bool { return true }

func (Card1321112JiuxiaoContact) OnPerTurn(ctx *EffectContext) error {
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	if len(opponent.Hand) < ctx.Engine.handLimitForPlayer(opponent) {
		addGeneratedCardToPlayerHand(ctx, ctx.OpponentID, "2001102")
	}
	return nil
}
