package game

import (
	"eraofarcane/model"
)

type Card1111103GreedyTyrant struct{ AlwaysActive }

func (Card1111103GreedyTyrant) ID() string { return "1111103" }

func (Card1111103GreedyTyrant) Name() string { return "贪婪暴君 卡姆 弗卡莱诺" }

func (Card1111103GreedyTyrant) ModifyGlobalCardPlayCost(ctx *EffectContext, card *CardInstance, cost map[string]int) {
	if ctx == nil || ctx.Engine == nil || card == nil || !cardInAnyPlayerHand(ctx.Engine, card) {
		return
	}
	cost[model.ElementArcane]++
}

func cardInAnyPlayerHand(engine *Engine, target *CardInstance) bool {
	if engine == nil || target == nil {
		return false
	}
	for _, player := range engine.State.Players {
		if player == nil {
			continue
		}
		for _, handCard := range player.Hand {
			if handCard != nil && handCard.InstanceID == target.InstanceID {
				return true
			}
		}
	}
	return false
}

var _ GlobalCardPlayCostModifier = Card1111103GreedyTyrant{}
