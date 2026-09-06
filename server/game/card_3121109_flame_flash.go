package game

import (
	"eraofarcane/model"
)

type Card3121109FlameFlash struct{ AlwaysActive }

func (Card3121109FlameFlash) ID() string { return "3121109" }

func (Card3121109FlameFlash) Name() string { return "烈焰闪" }

func (Card3121109FlameFlash) OnSpellHit(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementFire: 3})
	return nil
}
