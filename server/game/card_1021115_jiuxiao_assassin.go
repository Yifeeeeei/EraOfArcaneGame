package game

type Card1021115JiuxiaoAssassin struct{ AlwaysActive }

func (Card1021115JiuxiaoAssassin) ID() string { return "1021115" }

func (Card1021115JiuxiaoAssassin) Name() string { return "九霄刺客" }

func (Card1021115JiuxiaoAssassin) OnEnter(ctx *EffectContext) error {
	addGeneratedCardToPlayerHand(ctx, ctx.OpponentID, "2001102")
	return nil
}

func (Card1021115JiuxiaoAssassin) OnDeath(ctx *EffectContext) error {
	addGeneratedCardsToPlayerDeck(ctx, ctx.OpponentID, "2001102", 4)
	return nil
}
