package game

type Card1021012BlackMarketVendor struct{ AlwaysActive }

func (Card1021012BlackMarketVendor) ID() string { return "1021012" }

func (Card1021012BlackMarketVendor) Name() string { return "黑市商贩" }

func (Card1021012BlackMarketVendor) OnUltimate(ctx *EffectContext) error {
	candidates := append(
		ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool { return card.Card.IsItem() }),
		ctx.Engine.friendlyEquipment(ctx.PlayerID, nil)...,
	)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "discard_card",
		"选择1张手牌或装备区的道具弃置，抽2张牌",
		candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			if ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, selected[0]) {
				_ = DrawCards(2)(ctx)
			}
		})
	return nil
}
