package game

type Card2321009ChainLightningScroll struct{ AlwaysActive }

func (Card2321009ChainLightningScroll) ID() string   { return "2321009" }
func (Card2321009ChainLightningScroll) Name() string { return "连锁闪电卷轴" }

const chainLightningScrollDrawChoiceID = "chain_lightning_scroll_draw"

func (Card2321009ChainLightningScroll) OnSpellHit(ctx *EffectContext) error {
	if !isOwnSpellHit(ctx) {
		return nil
	}
	searchCandidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card.Card.Number == "2321009"
	})
	if len(searchCandidates) == 0 {
		return DrawCards(1)(ctx)
	}
	candidates := []map[string]any{{
		"instance_id": chainLightningScrollDrawChoiceID,
		"name":        "抽1张牌",
		"zone":        "choice",
		"side":        "own",
		"can_select":  true,
	}}
	candidates = append(candidates, searchCandidates...)
	ctx.Engine.SetPendingAction(ctx.PlayerID, "chain_lightning_scroll_choice",
		"连锁闪电卷轴:抽1张牌或检索1张连锁闪电卷轴",
		candidates, 1, 1,
		func(selected []string) {
			choice := firstSelected(selected)
			if choice == chainLightningScrollDrawChoiceID {
				_ = DrawCards(1)(ctx)
				return
			}
			ctx.Engine.searchDeckCardToHand(ctx.PlayerID, choice)
		})
	return nil
}
