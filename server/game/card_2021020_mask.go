package game

import "eraofarcane/model"

type Card2021020Mask struct{ AlwaysActive }

func (Card2021020Mask) ID() string   { return "2021020" }
func (Card2021020Mask) Name() string { return "假面" }
func (Card2021020Mask) OnEquip(ctx *EffectContext) error {
	hero := ctx.Engine.State.Players[ctx.PlayerID].Hero
	if hero == nil {
		return nil
	}
	total := totalLoad(hero)
	setElementsGain(hero, map[string]int{model.ElementArcane: total})
	return nil
}
