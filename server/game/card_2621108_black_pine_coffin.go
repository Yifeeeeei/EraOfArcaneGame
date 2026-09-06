package game

import (
	"eraofarcane/model"
)

type Card2621108BlackPineCoffin struct{ AlwaysActive }

func (Card2621108BlackPineCoffin) ID() string { return "2621108" }

func (Card2621108BlackPineCoffin) Name() string { return "黑松棺木" }

func (Card2621108BlackPineCoffin) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, isLowCostShadowCompanion)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "black_pine_coffin_discard_shadow_companions",
		"黑松棺木:丢弃最多2张低费暗影伙伴并结算它们的遗言", candidates, 0, min(2, len(candidates)),
		func(selected []string) {
			discarded := ctx.Engine.discardSelectedHandCardsMatching(ctx.PlayerID, selected, 2, isLowCostShadowCompanion)
			for _, card := range discarded {
				ctx.Engine.triggerEffects(TriggerOnDeath, card, nil, map[string]any{
					"death_cause": "black_pine_coffin",
					"from_zone":   "hand",
				})
			}
		})
	return nil
}

func isLowCostShadowCompanion(card *CardInstance) bool {
	return card != nil && card.Card != nil &&
		card.Card.IsCompanion() &&
		card.Card.Category == model.ElementShadow &&
		totalElementCost(card.Card.ElementsCost) < 5
}
