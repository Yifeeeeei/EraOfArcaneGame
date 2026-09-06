package game

type Card1401002SpiritBeastXinke struct{ AlwaysActive }

func (Card1401002SpiritBeastXinke) ID() string { return "1401002" }

func (Card1401002SpiritBeastXinke) Name() string { return "灵兽 辛柯" }

func (Card1401002SpiritBeastXinke) OnFriendlyDamagedFromHidden(ctx *EffectContext, event DamageEvent) error {
	if event.Matches(ctx.Source, DamageFriendly) && event.IsEnemyDamage(ctx.PlayerID) {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		candidates := make([]map[string]any, 0)
		for _, card := range ps.Hand {
			if card != nil && card.Card.Number == "1401002" {
				candidates = append(candidates, candidateInfo(card, "hand", "own"))
			}
		}
		candidates = append(candidates, ctx.Engine.friendlyDeckCards(ctx.PlayerID, func(card *CardInstance) bool { return card.Card.Number == "1401002" })...)
		ctx.Engine.SetPendingAction(ctx.PlayerID, "xinke_summon", "免费召唤灵兽 辛柯", candidates, 0, 1,
			func(selected []string) {
				if len(selected) > 0 {
					cardID := selected[0]
					positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
					if len(positions) == 0 {
						return
					}
					ctx.Engine.SetPendingAction(ctx.PlayerID, "xinke_summon_position", "选择灵兽 辛柯的入场位置", positions, 1, 1,
						func(posSelected []string) {
							if len(posSelected) == 0 {
								return
							}
							pos, ok := positionFromSelectionID(posSelected[0])
							if !ok {
								return
							}
							summonCardFreeFromHandOrDeckAtPosition(ctx, cardID, pos)
						})
				}
			})
	}
	return nil
}

func (Card1401002SpiritBeastXinke) HiddenDamageGroup(source *CardInstance) string {
	return source.Card.Number
}
