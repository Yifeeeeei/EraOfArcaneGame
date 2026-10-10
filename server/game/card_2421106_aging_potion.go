package game

import (
	"eraofarcane/model"
)

type Card2421106AgingPotion struct{ AlwaysActive }

func (Card2421106AgingPotion) ID() string { return "2421106" }

func (Card2421106AgingPotion) Name() string { return "苍老药剂" }

func (Card2421106AgingPotion) OnUseItem(ctx *EffectContext) error {
	candidates := friendlyFieldCardsIncludingBound(ctx.Engine, ctx.PlayerID, func(card *CardInstance) bool {
		behavior, ok := masteryBehavior(card)
		return ok && card.Statuses[StatusMastery] < behavior.MasteryMax() && reducibleElementLoad(card, model.ElementEarth) > 0
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "aging_potion_mastery",
		"苍老药剂:移除1点地负载并立刻达到下一次精通", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findFriendlyCardIncludingBound(ctx.PlayerID, firstSelected(selected))
			behavior, ok := masteryBehavior(target)
			if target == nil || !ok || target.Statuses[StatusMastery] >= behavior.MasteryMax() || reducibleElementLoad(target, model.ElementEarth) <= 0 {
				return
			}
			ctx.Engine.reduceCardElementLoadWithTriggers(ctx.PlayerID, target, model.ElementEarth, 1, ctx.Source)
			ctx.Engine.advanceMastery(target, ctx.PlayerID, 1)
		})
	return nil
}
