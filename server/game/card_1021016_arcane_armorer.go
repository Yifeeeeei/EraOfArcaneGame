package game

type Card1021016ArcaneArmorer struct{ AlwaysActive }

func (Card1021016ArcaneArmorer) ID() string   { return "1021016" }
func (Card1021016ArcaneArmorer) Name() string { return "奥术盔甲匠" }

func (Card1021016ArcaneArmorer) OnEnter(ctx *EffectContext) error {
	if ctx.Engine.hasAnyEquipment(ctx.PlayerID) {
		return nil
	}
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, isEquipmentItem)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "search_deck",
		"选择1张装备道具加入手牌",
		candidates, 1, 1,
		func(selected []string) {
			if len(selected) > 0 {
				ctx.Engine.searchDeckToHand(ctx.PlayerID, selected[0])
			}
		})
	return nil
}
