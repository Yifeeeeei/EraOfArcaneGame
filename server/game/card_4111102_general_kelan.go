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
	if !success || !triggeredTurnAvailable(ctx.Source) || !deckHasMatch(ctx.Engine.State.Players[ctx.PlayerID], isFireCard) {
		return nil
	}
	ctx.Engine.SetTriggeredTurnAction(ctx.Source, ctx.PlayerID, "general_kelan_flip_fire_card",
		"大将军 克兰:是否翻取1张火焰卡牌并弃1张手牌", []map[string]any{candidateInfo(ctx.Source, "hero", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || !deckHasMatch(ctx.Engine.State.Players[ctx.PlayerID], isFireCard) || !useTriggeredTurn(ctx.Source) {
				return
			}
			ctx.Engine.flipDeckMatchesToHandThen(ctx.PlayerID, 1, 0, isFireCard, func(drawn []*CardInstance) {
				if len(drawn) == 0 {
					return
				}
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
		})
	return nil
}

func isFireCard(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.Category == model.ElementFire
}
