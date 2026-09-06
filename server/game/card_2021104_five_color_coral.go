package game

import (
	"eraofarcane/model"
)

type Card2021104FiveColorCoral struct{ AlwaysActive }

func (Card2021104FiveColorCoral) ID() string { return "2021104" }

func (Card2021104FiveColorCoral) Name() string { return "五色珊瑚" }

func (Card2021104FiveColorCoral) OnEnter(ctx *EffectContext) error {
	choices := elementChoiceCandidates("2021104", model.ElementFire, model.ElementWater, model.ElementEarth, model.ElementAir, model.ElementLight, model.ElementShadow)
	ctx.Engine.SetPendingAction(ctx.PlayerID, "five_color_coral_load",
		"五色珊瑚:选择2种不同的非奥术元素各获得1点负载", choices, 2, 2,
		func(selected []string) {
			seen := make(map[string]bool, len(selected))
			for _, elem := range selected {
				if isNonArcaneElement(elem) && !seen[elem] {
					ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, elem, 1, ctx.Source)
					seen[elem] = true
				}
			}
		})
	return nil
}
