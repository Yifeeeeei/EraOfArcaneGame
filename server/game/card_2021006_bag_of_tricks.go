package game

type Card2021006BagOfTricks struct{ AlwaysActive }

func (Card2021006BagOfTricks) ID() string { return "2021006" }

func (Card2021006BagOfTricks) Name() string { return "百宝锦囊" }

func (Card2021006BagOfTricks) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, isConsumableCardInstance)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "bag_of_tricks_search_consumable",
		"献祭百宝锦囊,从卡组检索1张消耗品道具牌", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			if !ctx.Engine.sacrificeEquipment(ctx.PlayerID, ctx.Source.InstanceID) {
				return
			}
			ctx.Engine.searchDeckToHand(ctx.PlayerID, selected[0])
		})
	return nil
}
