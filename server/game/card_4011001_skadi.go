package game

import (
	"eraofarcane/model"
)

type Card4011001Skadi struct{ AlwaysActive }

func (Card4011001Skadi) ID() string { return "4011001" }

func (Card4011001Skadi) Name() string { return "\"南境百灵\" 斯卡尔蒂 罗佳" }

func (Card4011001Skadi) OnPerTurn(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool {
		if card == nil || card.Card == nil {
			return false
		}
		elem := card.Card.Category
		return elem == model.ElementArcane || ctx.Source.Statuses["斯卡蒂已用:"+elem] == 0
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "skadi_discard",
		"丢弃1张手牌，获得2点该属性元素", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			card, _ := ps.FindHandCard(selected[0])
			if card == nil || card.Card == nil {
				return
			}
			elem := card.Card.Category
			if elem != model.ElementArcane && ctx.Source.Statuses["斯卡蒂已用:"+elem] > 0 {
				return
			}
			if ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, selected[0]) {
				ps.Elements[elem] += 2
				if elem != model.ElementArcane {
					ctx.Source.Statuses["斯卡蒂已用:"+elem] = 1
				}
			}
		})
	return nil
}
