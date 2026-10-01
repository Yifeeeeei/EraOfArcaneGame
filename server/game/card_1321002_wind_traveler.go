package game

import (
	"eraofarcane/model"
)

type Card1321002WindTraveler struct{ AlwaysActive }

func (Card1321002WindTraveler) ID() string { return "1321002" }

func (Card1321002WindTraveler) Name() string { return "随风旅行者" }

func (Card1321002WindTraveler) OnEnter(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementAir: 1})
	return nil
}

func (Card1321002WindTraveler) OnDeath(ctx *EffectContext) error {
	return DrawCards(1)(ctx)
}
