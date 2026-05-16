package game

import "eraofarcane/model"

type Card1221006AquaticTanuki struct{}

func (Card1221006AquaticTanuki) ID() string   { return "1221006" }
func (Card1221006AquaticTanuki) Name() string { return "水栖狸猫" }

func (Card1221006AquaticTanuki) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	count := 0
	for _, unit := range adjacentUnits(ctx.Engine.State.Players[ctx.PlayerID], ctx.Source.Position) {
		if unit.Card.IsCompanion() && unit.Card.Category == model.ElementWater {
			count++
		}
	}
	if count >= 2 && ctx.Source.ElementsGainBonus[model.ElementWater] < 1 {
		addElementsGainBonus(ctx.Source, model.ElementWater, 1)
	}
	return nil
}
