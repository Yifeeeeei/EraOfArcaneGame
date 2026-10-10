package game

import (
	"fmt"
)

type Card1321105Illusionist struct{ AlwaysActive }

func (Card1321105Illusionist) ID() string { return "1321105" }

func (Card1321105Illusionist) Name() string { return "幻术师" }

func (Card1321105Illusionist) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() &&
			totalElementCost(card.Card.ElementsCost) < 6
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "illusionist_return_companion",
		"幻术师:选择1个入场花费小于6的友方伙伴移回手牌", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() || totalElementCost(target.Card.ElementsCost) >= 6 {
				return
			}
			gain := copyElementAmounts(effectiveElementsGain(target))
			ctx.Engine.returnUnitToHand(target, ctx.PlayerID)
			resetCardForHiddenZone(target)
			ctx.Engine.State.Players[ctx.PlayerID].GainElements(gain)
		})
	return nil
}

func (Card1321105Illusionist) ValidateAbility(ctx *EffectContext, trigger EffectTrigger) error {
	if trigger != TriggerUltimate {
		return nil
	}
	if len(ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(candidate *CardInstance) bool {
		return candidate != nil && candidate.Card != nil && candidate.Card.IsCompanion() &&
			totalElementCost(candidate.Card.ElementsCost) < 6
	})) == 0 {
		return fmt.Errorf("幻术师需要1个入场花费小于6的友方伙伴")
	}
	return nil
}

func copyElementAmounts(src map[string]int) map[string]int {
	copied := make(map[string]int, len(src))
	for elem, amount := range src {
		if amount > 0 {
			copied[elem] = amount
		}
	}
	return copied
}
