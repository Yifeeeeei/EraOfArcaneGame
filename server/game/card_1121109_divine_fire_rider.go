package game

import (
	"eraofarcane/model"
)

type Card1121109DivineFireRider struct{ AlwaysActive }

func (Card1121109DivineFireRider) ID() string { return "1121109" }

func (Card1121109DivineFireRider) Name() string { return "神火兽骑手" }

func (Card1121109DivineFireRider) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != ctx.Source && isFireCompanion(card) && ctx.Engine.canConsumeCard(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "divine_fire_rider_consume_companion",
		"神火兽骑手:消耗1个其他友方火焰伙伴", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || target == ctx.Source || zone != "unit" || !isFireCompanion(target) || !ctx.Engine.canConsumeCard(target) {
				return
			}
			bonus := totalElementCost(target.Card.ElementsCost)
			if bonus <= 0 {
				return
			}
			ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, target, ctx.Engine.effectiveElementsGain(target), "1121109")
			ctx.Engine.addNextElementSpellPowerBonus(ctx.PlayerID, model.ElementFire, bonus)
		})
	return nil
}
