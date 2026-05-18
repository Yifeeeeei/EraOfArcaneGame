package game

import (
	"fmt"

	"eraofarcane/model"
)

type Card2521007BlueCrystalLamp struct{ AlwaysActive }

func (Card2521007BlueCrystalLamp) ID() string   { return "2521007" }
func (Card2521007BlueCrystalLamp) Name() string { return "蓝晶灯盏" }

func (Card2521007BlueCrystalLamp) OnUltimate(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	cost := map[string]int{model.ElementLight: 5}
	if !ps.CanPayCost(cost) {
		return fmt.Errorf("光辉元素不足")
	}
	ps.PayCost(cost)
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementLight, 2, ctx.Source)
	return nil
}
