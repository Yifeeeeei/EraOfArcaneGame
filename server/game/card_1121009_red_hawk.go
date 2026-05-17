package game

type Card1121009RedHawk struct{ AlwaysActive }

func (Card1121009RedHawk) ID() string   { return "1121009" }
func (Card1121009RedHawk) Name() string { return "赤鹰" }

func (Card1121009RedHawk) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, isFireCompanionWithCostAboveFour)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "search_deck",
		"选择1张入场花费大于4的火焰伙伴加入手牌",
		candidates, 1, 1,
		func(selected []string) {
			if len(selected) > 0 {
				ctx.Engine.searchDeckToHand(ctx.PlayerID, selected[0])
			}
		})
	return nil
}
