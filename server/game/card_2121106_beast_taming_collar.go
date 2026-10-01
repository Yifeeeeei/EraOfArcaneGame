package game

import (
	"strings"
)

type Card2121106BeastTamingCollar struct{ AlwaysActive }

func (Card2121106BeastTamingCollar) ID() string { return "2121106" }

func (Card2121106BeastTamingCollar) Name() string { return "驯兽项圈" }

func (Card2121106BeastTamingCollar) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, isCollarEligibleFireCompanion)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "beast_taming_collar_target",
		"驯兽项圈:选择1个巫师以外的火焰伙伴", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || !isCollarEligibleFireCompanion(target) {
				return
			}
			clearStatusPrefix(ctx.Source, beastTamingCollarTargetPrefix)
			ctx.Source.Statuses[beastTamingCollarTargetPrefix+target.InstanceID] = 1
		})
	return nil
}

func (Card2121106BeastTamingCollar) PerTurnLabel(*CardInstance) string {
	return "消耗目标"
}

func (Card2121106BeastTamingCollar) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.equipmentInOwnerSlot(ctx.PlayerID, ctx.Source) {
		return nil
	}
	target := collarTarget(ctx.Engine, ctx.PlayerID, ctx.Source)
	if target == nil || !ctx.Engine.canConsumeCard(target) {
		return nil
	}
	ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, target, target.Card.ElementsCost, "2121106")
	return nil
}

const beastTamingCollarTargetPrefix = "驯兽项圈目标:"

func isCollarEligibleFireCompanion(card *CardInstance) bool {
	return isFireCompanion(card) && !hasCardTag(card.Card, "巫师")
}

func collarTarget(e *Engine, playerID int, collar *CardInstance) *CardInstance {
	if e == nil || collar == nil {
		return nil
	}
	for status, amount := range collar.Statuses {
		if amount <= 0 || !strings.HasPrefix(status, beastTamingCollarTargetPrefix) {
			continue
		}
		target, zone := e.findFriendlyCandidate(playerID, strings.TrimPrefix(status, beastTamingCollarTargetPrefix))
		if target != nil && zone == "unit" && isCollarEligibleFireCompanion(target) {
			return target
		}
	}
	return nil
}
