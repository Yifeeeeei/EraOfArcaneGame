package game

import (
	"eraofarcane/model"
)

type Card2421108EmeraldFruit struct{ AlwaysActive }

func (Card2421108EmeraldFruit) ID() string { return "2421108" }

func (Card2421108EmeraldFruit) Name() string { return "翡翠果" }

func (Card2421108EmeraldFruit) OnEnter(ctx *EffectContext) error {
	targets := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "emerald_fruit_target",
		"翡翠果:选择1个友方伙伴获得负载", targets, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() {
				return
			}
			choices := elementChoiceCandidates("2421108", model.ElementFire, model.ElementWater, model.ElementAir, model.ElementLight, model.ElementShadow)
			ctx.Engine.SetPendingAction(ctx.PlayerID, "emerald_fruit_element",
				"翡翠果:选择除地与奥术外的1点负载", choices, 1, 1,
				func(selected []string) {
					elem := firstSelected(selected)
					if elem != model.ElementEarth && isNonArcaneElement(elem) {
						ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, elem, 1, ctx.Source)
					}
				})
		})
	return nil
}
