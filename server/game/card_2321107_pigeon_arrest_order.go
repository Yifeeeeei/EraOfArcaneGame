package game

type Card2321107PigeonArrestOrder struct{ AlwaysActive }

func (Card2321107PigeonArrestOrder) ID() string { return "2321107" }

func (Card2321107PigeonArrestOrder) Name() string { return "飞鸽拘捕令" }

func (Card2321107PigeonArrestOrder) OnSpellHit(ctx *EffectContext) error {
	if !isFriendlySpellHit(ctx) || !useTriggeredTurn(ctx.Source) {
		return nil
	}
	addGeneratedCardToPlayerHand(ctx, ctx.OpponentID, "2001102")
	return nil
}
