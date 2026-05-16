package game

type Card3021004Engrave struct{}

func (Card3021004Engrave) ID() string   { return "3021004" }
func (Card3021004Engrave) Name() string { return "刻印" }

func (Card3021004Engrave) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card != ctx.Source
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "engrave_discard",
		"丢弃1张手牌,抽取牌堆中第一张卷轴或符文", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 || !ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, selected[0]) {
				return
			}
			ctx.Engine.drawFirstDeckMatch(ctx.PlayerID, isRuneOrScroll)
		})
	return nil
}
