package game

import (
	"eraofarcane/model"
)

type Card3501001UnitedHope struct{ AlwaysActive }

func (Card3501001UnitedHope) ID() string { return "3501001" }

func (Card3501001UnitedHope) Name() string { return "团结的希望" }

func (Card3501001UnitedHope) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	limit := min(5, len(ps.Deck))
	candidates := make([]map[string]any, 0, limit)
	for i := 0; i < limit; i++ {
		card := ps.Deck[i]
		if card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementLight {
			continue
		}
		candidates = append(candidates, candidateInfo(card, "deck", "own"))
	}
	if len(candidates) == 0 {
		ctx.Engine.shuffleDeck(ctx.PlayerID)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "united_hope_search",
		"团结的希望:从卡组上方5张中选择1张光辉伙伴加入手牌", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				ctx.Engine.shuffleDeck(ctx.PlayerID)
				return
			}
			ctx.Engine.searchDeckToHand(ctx.PlayerID, selected[0])
			ctx.Engine.shuffleDeck(ctx.PlayerID)
		})
	return nil
}
