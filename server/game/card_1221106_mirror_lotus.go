package game

import (
	"eraofarcane/model"
)

type Card1221106MirrorLotus struct{ AlwaysActive }

func (Card1221106MirrorLotus) ID() string { return "1221106" }

func (Card1221106MirrorLotus) Name() string { return "镜花海之莲" }

func (Card1221106MirrorLotus) IsPrayerAbility() bool { return true }

func (Card1221106MirrorLotus) OnPerTurn(ctx *EffectContext) error {
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementWater, 1, ctx.Source)
	return nil
}
