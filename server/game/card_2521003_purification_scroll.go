package game

var negativeStatuses = []string{StatusBurn, StatusFreeze, StatusStun, StatusPetrify, StatusWeaken}

type Card2521003PurificationScroll struct{}

func (Card2521003PurificationScroll) ID() string   { return "2521003" }
func (Card2521003PurificationScroll) Name() string { return "净化卷轴" }

func (Card2521003PurificationScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return hasAnyNegativeStatus(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "purify_friendly",
		"选择1个友方卡牌移除所有负面状态", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			if target == nil || zone != "unit" {
				return
			}
			for _, status := range negativeStatuses {
				delete(target.Statuses, status)
			}
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"target": cardToInfo(target),
				"effect": "purify",
			}})
		})
	return nil
}

func hasAnyNegativeStatus(card *CardInstance) bool {
	if card == nil {
		return false
	}
	for _, status := range negativeStatuses {
		if card.Statuses[status] > 0 {
			return true
		}
	}
	return false
}
