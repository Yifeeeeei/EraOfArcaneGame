package game

type Card1221015LookoutMerchantShip struct{ AlwaysActive }

func (Card1221015LookoutMerchantShip) ID() string   { return "1221015" }
func (Card1221015LookoutMerchantShip) Name() string { return "眺望者商舰" }

func (Card1221015LookoutMerchantShip) OnPerTurn(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, isWaterCard)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "merchant_ship_search",
		"检索1张水纹卡牌", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 || !ctx.Engine.searchDeckToHand(ctx.PlayerID, selected[0]) {
				return
			}
			hand := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
			ctx.Engine.SetPendingAction(ctx.PlayerID, "merchant_ship_shuffle_hand",
				"选择1张手牌洗回卡组", hand, 1, 1,
				func(selected []string) {
					if len(selected) == 0 {
						return
					}
					ctx.Engine.moveHandCardToDeckBottom(ctx.PlayerID, selected[0])
				})
		})
	return nil
}
