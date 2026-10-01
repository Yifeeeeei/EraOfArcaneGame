package game

import (
	"eraofarcane/model"
)

type Card3421011NaturalGrowth struct{ AlwaysActive }

func (Card3421011NaturalGrowth) ID() string { return "3421011" }

func (Card3421011NaturalGrowth) Name() string { return "自然生长" }

func (Card3421011NaturalGrowth) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Card.Category == model.ElementEarth && card.IsHorizontal && totalLoad(card) < 4
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "natural_growth",
		"自然生长:选择1个横置且负载小于4的地脉伙伴，获得1地负载", candidates, 1, 1,
		func(selected []string) {
			card, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if card == nil || !card.IsHorizontal || totalLoad(card) >= 4 {
				return
			}
			ctx.Engine.addElementsGainBonus(card, ctx.PlayerID, model.ElementEarth, 1, ctx.Source)
		})
	return nil
}
