package game

import "fmt"

type Card2121108BurnoutScroll struct{ AlwaysActive }

func (Card2121108BurnoutScroll) ID() string { return "2121108" }

func (Card2121108BurnoutScroll) Name() string { return "燃烬卷轴" }

func (Card2121108BurnoutScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return isFireCompanion(card) && ctx.Engine.canConsumeCard(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "burnout_scroll_consume_fire_companion",
		"燃烬卷轴:消耗1个友方火焰伙伴并获得其入场花费的元素", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || !isFireCompanion(target) || !ctx.Engine.canConsumeCard(target) {
				return
			}
			ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, target, target.Card.ElementsCost, "2121108")
		})
	return nil
}

func (Card2121108BurnoutScroll) ValidateItemUse(ctx *EffectContext) error {
	e, playerID := ctx.Engine, ctx.PlayerID
	if len(e.friendlyUnits(playerID, false, func(unit *CardInstance) bool {
		return isFireCompanion(unit) && e.canConsumeCard(unit)
	})) == 0 {
		return fmt.Errorf("Burnout Scroll requires a ready friendly fire companion")
	}
	return nil
}
