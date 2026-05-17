package game

type Card3621012SoulRecall struct{ AlwaysActive }

func (Card3621012SoulRecall) ID() string   { return "3621012" }
func (Card3621012SoulRecall) Name() string { return "回魂术" }

func (Card3621012SoulRecall) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	candidates := make([]map[string]any, 0)
	for _, card := range ps.Graveyard {
		if card != nil && card.Card.IsCompanion() {
			candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "soul_recall",
		"从弃牌堆选择最多2个伙伴移回手牌", candidates, 1, 2,
		func(selected []string) {
			for _, id := range selected {
				ctx.Engine.moveGraveyardCardToHand(ctx.PlayerID, id)
			}
		})
	return nil
}
