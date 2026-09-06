package game

import (
	"eraofarcane/model"
)

type Card1221014NorthSeaFlyingFish struct{ AlwaysActive }

func (Card1221014NorthSeaFlyingFish) ID() string { return "1221014" }

func (Card1221014NorthSeaFlyingFish) Name() string { return "北海飞鱼" }

func (Card1221014NorthSeaFlyingFish) OnPerTurn(ctx *EffectContext) error {
	setElementsGain(ctx.Source, map[string]int{model.ElementAir: 1})
	return nil
}

func (Card1221014NorthSeaFlyingFish) OnTurnEnd(ctx *EffectContext) error {
	clearElementsGainSet(ctx.Source)
	return nil
}
