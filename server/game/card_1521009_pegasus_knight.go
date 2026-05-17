package game

type Card1521009PegasusKnight struct{ AlwaysActive }

func (Card1521009PegasusKnight) ID() string   { return "1521009" }
func (Card1521009PegasusKnight) Name() string { return "天马骑士" }

func (Card1521009PegasusKnight) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card.Card.Number == "1521012"
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "search_deck",
		"选择1张独角天马加入手牌",
		candidates, 1, 1,
		func(selected []string) {
			if len(selected) > 0 {
				ctx.Engine.searchDeckToHand(ctx.PlayerID, selected[0])
			}
		})
	return nil
}
