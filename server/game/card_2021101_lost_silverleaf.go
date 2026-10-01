package game

type Card2021101LostSilverleaf struct{ AlwaysActive }

func (Card2021101LostSilverleaf) ID() string { return "2021101" }

func (Card2021101LostSilverleaf) Name() string { return "失落的银叶花" }

func (Card2021101LostSilverleaf) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.drawCards(ctx.PlayerID, 2)
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lost_silverleaf_discard",
		"失落的银叶花:弃1张手牌", candidates, 1, 1,
		func(selected []string) {
			if len(selected) > 0 {
				ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, selected[0])
			}
		})
	return nil
}
