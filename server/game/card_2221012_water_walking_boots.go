package game

import "eraofarcane/model"

type Card2221012WaterWalkingBoots struct{ AlwaysActive }

func (Card2221012WaterWalkingBoots) ID() string   { return "2221012" }
func (Card2221012WaterWalkingBoots) Name() string { return "水行之靴" }

func (Card2221012WaterWalkingBoots) OnTurnStart(ctx *EffectContext) error {
	hero := ctx.Engine.State.Players[ctx.PlayerID].Hero
	if hero == nil || hero.Position == nil {
		return nil
	}
	count := 0
	for _, unit := range adjacentUnits(ctx.Engine.State.Players[ctx.PlayerID], hero.Position) {
		if unit.Card.IsCompanion() && unit.Card.Category == model.ElementWater {
			count++
		}
	}
	if count >= 3 && ctx.Source.ElementsGainBonus[model.ElementWater] < 1 {
		ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementWater, 1, ctx.Source)
	}
	return nil
}
