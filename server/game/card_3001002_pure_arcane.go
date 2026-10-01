package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card3001002PureArcane struct{ AlwaysActive }

func (Card3001002PureArcane) ID() string { return "3001002" }

func (Card3001002PureArcane) Name() string { return "纯净奥术" }

func (Card3001002PureArcane) NeedsSpellTarget() bool {
	return false
}

func (Card3001002PureArcane) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) {
		return nil
	}
	choices := ctx.Engine.pureArcaneChoices(ctx.PlayerID)
	if len(choices) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "pure_arcane_spend",
		"选择同种元素与数量,使下一次该属性法术+等量威力", choices, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			elem, amount, ok := parsePureArcaneChoice(selected[0])
			if !ok {
				return
			}
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			if ps.Elements[elem] < amount {
				return
			}
			ps.Elements[elem] -= amount
			ctx.Engine.addNextElementSpellPowerBonus(ctx.PlayerID, elem, amount)
		})
	return nil
}

func (Card3001002PureArcane) ResolvesSorceryHit(*CardInstance) bool { return false }

func (e *Engine) pureArcaneChoices(playerID int) []map[string]any {
	ps := e.State.Players[playerID]
	choices := make([]map[string]any, 0)
	for _, elem := range []string{model.ElementFire, model.ElementWater, model.ElementAir, model.ElementEarth, model.ElementLight, model.ElementShadow, model.ElementArcane} {
		available := ps.Elements[elem]
		if available <= 0 {
			continue
		}
		for amount := 1; amount <= min(available, 10); amount++ {
			id := elem + ":" + fmt.Sprintf("%d", amount)
			choices = append(choices, map[string]any{
				"instance_id": id,
				"number":      "3001002",
				"name":        elem + " " + fmt.Sprintf("%d", amount),
				"type":        "元素选择",
				"zone":        "choice",
				"side":        "own",
				"element":     elem,
				"amount":      amount,
			})
		}
	}
	return choices
}
