package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card2521007BlueCrystalLamp struct{ AlwaysActive }

func (Card2521007BlueCrystalLamp) ID() string { return "2521007" }

func (Card2521007BlueCrystalLamp) Name() string { return "蓝晶灯盏" }

func (Card2521007BlueCrystalLamp) OnUltimate(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	cost := map[string]int{model.ElementLight: 5}
	if !ctx.Engine.canPayCost(ps, cost) {
		return fmt.Errorf("光辉元素不足")
	}
	ctx.Engine.payCostForAction(ps, cost, ActionMessage{})
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementLight, 2, ctx.Source)
	return nil
}
