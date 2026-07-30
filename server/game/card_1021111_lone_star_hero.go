package game

import "eraofarcane/model"

type Card1021111LoneStarHero struct{ AlwaysActive }

func (Card1021111LoneStarHero) ID() string   { return "1021111" }
func (Card1021111LoneStarHero) Name() string { return "孤星勇者" }

func (Card1021111LoneStarHero) ModifySelfCardPlayCost(ctx *EffectContext, cost map[string]int) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return
	}
	hand := ctx.Engine.State.Players[ctx.PlayerID].Hand
	otherHandCards := 0
	for _, card := range hand {
		if card != nil && card.InstanceID != ctx.Source.InstanceID {
			otherHandCards++
		}
	}
	if otherHandCards > 0 {
		cost[model.ElementArcane] += otherHandCards
	}
}

var _ SelfCardPlayCostModifier = Card1021111LoneStarHero{}
