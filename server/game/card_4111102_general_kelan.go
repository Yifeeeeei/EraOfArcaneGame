package game

import (
	"eraofarcane/model"
)

type Card4111102GeneralKelan struct{ AlwaysActive }

func (Card4111102GeneralKelan) ID() string { return "4111102" }

func (Card4111102GeneralKelan) Name() string { return "大将军 克兰" }

func (Card4111102GeneralKelan) OnDefend(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	success, _ := ctx.ExtraData["defense_success"].(bool)
	if !success || !useTriggeredTurn(ctx.Source) {
		return nil
	}
	ctx.Engine.flipDeckMatchesToHandThen(ctx.PlayerID, 1, 0, isFireCard, func(_ []*CardInstance) {
		candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
		if len(candidates) == 0 {
			return
		}
		ctx.Engine.SetPendingAction(ctx.PlayerID, "general_kelan_discard",
			"大将军 克兰:弃1张手牌", candidates, 1, 1,
			func(discardSelected []string) {
				ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, firstSelected(discardSelected))
			})
	})
	return nil
}

func isFireCard(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.Category == model.ElementFire
}
