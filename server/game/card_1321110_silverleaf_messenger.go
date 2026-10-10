package game

type Card1321110SilverleafMessenger struct{ AlwaysActive }

func (Card1321110SilverleafMessenger) ID() string { return "1321110" }

func (Card1321110SilverleafMessenger) Name() string { return "银叶信使" }

func (Card1321110SilverleafMessenger) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.Number == "2021101"
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "silverleaf_messenger_search",
		"银叶信使:检索1张失落的银叶花", candidates, 1, 1,
		func(selected []string) {
			ctx.Engine.searchDeckToHand(ctx.PlayerID, firstSelected(selected))
		})
	return nil
}
