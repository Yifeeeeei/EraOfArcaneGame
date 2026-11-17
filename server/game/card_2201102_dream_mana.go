package game

import (
	"eraofarcane/model"
)

type Card2201102DreamMana struct{ AlwaysActive }

func (Card2201102DreamMana) ID() string { return "2201102" }

func (Card2201102DreamMana) Name() string { return "幻创之梦-幻能" }

func (Card2201102DreamMana) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementArcane: 3})
	return nil
}
