package game

type Card1021017Runemaster struct{ AlwaysActive }

func (Card1021017Runemaster) ID() string { return "1021017" }

func (Card1021017Runemaster) Name() string { return "符文师" }

func (Card1021017Runemaster) OnEnter(ctx *EffectContext) error {
	discardCandidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	searchCandidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, isRuneOrScroll)
	if len(discardCandidates) == 0 || len(searchCandidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "discard_for_search",
		"选择1张手牌弃置，然后检索1张符文或卷轴",
		discardCandidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 || !ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, selected[0]) {
				return
			}
			candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, isRuneOrScroll)
			if len(candidates) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "search_deck",
				"选择1张符文或卷轴加入手牌",
				candidates, 1, 1,
				func(selected []string) {
					if len(selected) > 0 {
						ctx.Engine.searchDeckToHand(ctx.PlayerID, selected[0])
					}
				})
		})
	return nil
}
