package game

import (
	"eraofarcane/model"
)

type Card1021107GiftedYouth struct{ AlwaysActive }

func (Card1021107GiftedYouth) ID() string { return "1021107" }

func (Card1021107GiftedYouth) Name() string { return "天才少年" }

func (Card1021107GiftedYouth) MasteryMax() int { return 2 }

func (Card1021107GiftedYouth) OnMastery(ctx *EffectContext, level int) error {
	if level != 2 {
		return nil
	}
	source := ctx.Source
	choices := elementChoiceCandidates("1021107", model.ElementFire, model.ElementWater, model.ElementEarth, model.ElementAir, model.ElementLight, model.ElementShadow)
	ctx.Engine.SetPendingAction(ctx.PlayerID, "gifted_youth_mastery_load",
		"天才少年:选择获得的非奥术负载", choices, 1, 1,
		func(selected []string) {
			elem := firstSelected(selected)
			if !isNonArcaneElement(elem) || !ctx.Engine.cardStillOnField(source) {
				return
			}
			ctx.Engine.addElementsGainBonus(source, ctx.PlayerID, elem, 1, source)
		})
	return nil
}
