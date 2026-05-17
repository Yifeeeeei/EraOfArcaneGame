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
	addElementsGainBonus(ctx.Source, model.ElementLight, 2)
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":   cardToInfo(ctx.Source),
		"effect":   "load_bonus",
		"element":  model.ElementLight,
		"amount":   2,
		"elements": ps.Elements,
	}})
	return nil
}
