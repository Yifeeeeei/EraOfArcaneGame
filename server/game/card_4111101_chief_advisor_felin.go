package game

import (
	"eraofarcane/model"
)

type Card4111101ChiefAdvisorFelin struct{ AlwaysActive }

func (Card4111101ChiefAdvisorFelin) ID() string { return "4111101" }

func (Card4111101ChiefAdvisorFelin) Name() string { return "首席顾问 费林" }

func (Card4111101ChiefAdvisorFelin) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, isFireCompanion)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "felin_sacrifice_fire_companion",
		"首席顾问 费林:献祭1个友方火焰伙伴", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || !isFireCompanion(target) {
				return
			}
			cost := copyElementCost(target.Card.ElementsCost)
			if totalElementCost(cost) <= 0 {
				return
			}
			ctx.Engine.destroyUnitWithCause(target, ctx.PlayerID, DeathCauseSacrifice)
			for _, elem := range model.AllElements {
				if cost[elem] <= 0 {
					continue
				}
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
					Type:          TempModNextFireCardPlayCostMinus,
					Element:       elem,
					Amount:        cost[elem],
					RemainingUses: 1,
					ExpiresTurn:   ctx.Engine.State.TurnNumber + 1,
				})
			}
		})
	return nil
}
