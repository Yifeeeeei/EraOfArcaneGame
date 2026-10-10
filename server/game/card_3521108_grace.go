package game

import (
	"eraofarcane/model"
)

type Card3521108Grace struct{ AlwaysActive }

func (Card3521108Grace) ID() string { return "3521108" }

func (Card3521108Grace) Name() string { return "恩典" }

func (Card3521108Grace) OnSpellCast(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.CurrentLife < maxLife(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "grace_heal_companion",
		"恩典:选择1个受伤友方伙伴回复2血", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findUnitByInstanceID(firstSelected(selected))
			if target == nil || target.OwnerID != ctx.PlayerID || target.Card == nil || !target.Card.IsCompanion() || target.Card.IsHero() || target.CurrentLife >= maxLife(target) {
				return
			}
			ctx.Engine.healUnit(target, 2, ctx.Source)
			if target.CurrentLife >= maxLife(target) {
				target.Statuses["max_life_bonus"]++
				ctx.Engine.gainLife(target, 1, ctx.Source)
				ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
			}
		})
	return nil
}
