package game

import (
	"eraofarcane/model"
)

type Card3621110BloodNourish struct{ AlwaysActive }

func (Card3621110BloodNourish) ID() string { return "3621110" }

func (Card3621110BloodNourish) Name() string { return "鲜血滋养" }

func (Card3621110BloodNourish) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	candidates := make([]map[string]any, 0)
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.Category == model.ElementShadow {
			candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_nourish_exile",
		"鲜血滋养:选择弃牌堆1张暗影卡牌移出游戏，获得2暗", candidates, 1, 1,
		func(selected []string) {
			for _, card := range ps.Graveyard {
				if card != nil && card.InstanceID == firstSelected(selected) && card.Card != nil && card.Card.Category == model.ElementShadow {
					if ctx.Engine.exileCard(ctx.PlayerID, card) {
						ps.GainElements(map[string]int{model.ElementShadow: 2})
					}
					return
				}
			}
		})
	return nil
}
