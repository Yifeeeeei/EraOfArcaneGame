package game

import (
	"eraofarcane/model"
)

type Card2321103ThunderBreath struct{ AlwaysActive }

func (Card2321103ThunderBreath) ID() string { return "2321103" }

func (Card2321103ThunderBreath) Name() string { return "雷鸣之息" }

func (Card2321103ThunderBreath) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementAir: 1})
	return nil
}
