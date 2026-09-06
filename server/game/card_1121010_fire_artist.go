package game

type Card1121010FireArtist struct{ AlwaysActive }

func (Card1121010FireArtist) ID() string { return "1121010" }

func (Card1121010FireArtist) Name() string { return "火焰艺人" }

func (Card1121010FireArtist) OnUltimate(ctx *EffectContext) error {
	ctx.Source.IsHorizontal = true
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card != ctx.Source && isNonHeroFireCard(card)
	})
	candidates = append(candidates, ctx.Engine.friendlySkills(ctx.PlayerID, isNonHeroFireCard)...)
	candidates = append(candidates, ctx.Engine.friendlyEquipment(ctx.PlayerID, isNonHeroFireCard)...)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "reset_card",
		"选择另1张人物牌以外的火焰牌重置",
		candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			card, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			resetCard(card)
		})
	return nil
}
