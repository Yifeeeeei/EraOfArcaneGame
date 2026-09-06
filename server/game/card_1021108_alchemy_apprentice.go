package game

import (
	"eraofarcane/model"
	"fmt"
	"strings"
)

type Card1021108AlchemyApprentice struct{ AlwaysActive }

func (Card1021108AlchemyApprentice) ID() string { return "1021108" }

func (Card1021108AlchemyApprentice) Name() string { return "炼金术学徒" }

func (Card1021108AlchemyApprentice) PerTurnLabel(*CardInstance) string {
	return "主动"
}

func (Card1021108AlchemyApprentice) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.canConsumeCard(ctx.Source) {
		return fmt.Errorf("炼金术学徒不能被消耗")
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps.Elements[model.ElementArcane] < 1 {
		return fmt.Errorf("炼金术学徒需要1点奥术元素")
	}
	choices := make([]map[string]any, 0, 12)
	for _, elem := range []string{model.ElementFire, model.ElementWater, model.ElementEarth, model.ElementAir, model.ElementLight, model.ElementShadow} {
		for i := 1; i <= 2; i++ {
			choices = append(choices, map[string]any{"instance_id": fmt.Sprintf("%s#%d", elem, i), "number": "1021108", "name": elem, "type": "选择", "zone": "choice", "side": "own", "element": elem})
		}
	}
	ctx.Source.IsHorizontal = true
	ps.Elements[model.ElementArcane]--
	ctx.Engine.SetPendingAction(ctx.PlayerID, "alchemy_apprentice_elements",
		"炼金术学徒:选择2点非奥术元素", choices, 2, 2,
		func(selected []string) {
			gain := make(map[string]int)
			seen := make(map[string]bool, len(selected))
			for _, id := range selected {
				if seen[id] {
					continue
				}
				seen[id] = true
				elem, _, ok := strings.Cut(id, "#")
				if ok && isNonArcaneElement(elem) {
					gain[elem]++
				}
			}
			if len(gain) > 0 {
				ps.GainElements(gain)
			}
		})
	return nil
}
