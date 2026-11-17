package game

import (
	"eraofarcane/model"
)

type Card1521106ChurchExorcist struct{ AlwaysActive }

func (Card1521106ChurchExorcist) ID() string { return "1521106" }

func (Card1521106ChurchExorcist) Name() string { return "教廷驱魔师" }

func (Card1521106ChurchExorcist) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, hasAnyNegativeStatus)
	candidates = append(candidates, ctx.Engine.friendlyEquipment(ctx.PlayerID, hasAnyNegativeStatus)...)
	candidates = append(candidates, ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, hasAnyNegativeStatus)...)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "church_exorcist_purify",
		"教廷驱魔师:选择1张友方卡牌移除全部负面状态", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findFriendlyCardIncludingBound(ctx.PlayerID, firstSelected(selected))
			removed := countNegativeStatusLayers(target)
			if removed <= 0 {
				return
			}
			clearNegativeStatuses(target)
			ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementLight: removed})
		})
	return nil
}

func countNegativeStatusLayers(card *CardInstance) int {
	if card == nil {
		return 0
	}
	total := 0
	for _, status := range negativeStatuses {
		if card.Statuses[status] > 0 {
			total += card.Statuses[status]
		}
	}
	return total
}
