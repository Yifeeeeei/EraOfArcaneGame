package game

import "eraofarcane/model"

type Card1421012WoodlandFlyingSquirrel struct{}

func (Card1421012WoodlandFlyingSquirrel) ID() string   { return "1421012" }
func (Card1421012WoodlandFlyingSquirrel) Name() string { return "林地飞鼠" }

func (Card1421012WoodlandFlyingSquirrel) OnPerTurn(ctx *EffectContext) error {
	setElementsGain(ctx.Source, map[string]int{model.ElementAir: 1})
	return nil
}
