package game

type Card3021004Engrave struct{ AlwaysActive }

func (Card3021004Engrave) ID() string   { return "3021004" }
func (Card3021004Engrave) Name() string { return "刻印" }

func (Card3021004Engrave) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	discardCandidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	searchCandidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, isRuneOrScroll)
	if len(discardCandidates) == 0 || len(searchCandidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "engrave_discard",
		"刻印:选择1张手牌弃置，然后检索1张卷轴或符文", discardCandidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 || !ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, selected[0]) {
				return
			}
			candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, isRuneOrScroll)
			if len(candidates) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "engrave_search",
				"刻印:从卡组检索1张卷轴或符文", candidates, 1, 1,
				func(selected []string) {
					if len(selected) == 0 {
						return
					}
					ctx.Engine.searchDeckToHand(ctx.PlayerID, selected[0])
				})
		})
	return nil
}
