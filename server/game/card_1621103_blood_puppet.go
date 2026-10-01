package game

import ()

type Card1621103BloodPuppet struct{ AlwaysActive }

func (Card1621103BloodPuppet) ID() string { return "1621103" }

func (Card1621103BloodPuppet) Name() string { return "鲜血傀儡" }

func (Card1621103BloodPuppet) OnEnter(ctx *EffectContext) error {
	ctx.DealDamage(ctx.Engine.State.Players[ctx.PlayerID].Hero, 2)
	return nil
}
