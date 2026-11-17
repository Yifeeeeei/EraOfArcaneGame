package game

import (
	"eraofarcane/model"
)

type Card1021018ArcaneBulwark struct{ AlwaysActive }

func (Card1021018ArcaneBulwark) ID() string { return "1021018" }

func (Card1021018ArcaneBulwark) Name() string { return "奥术壁垒" }

func (Card1021018ArcaneBulwark) OnDeath(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.OpponentID].GainElements(map[string]int{model.ElementArcane: 2})
	return nil
}
