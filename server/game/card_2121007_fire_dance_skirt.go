package game

type Card2121007FireDanceSkirt struct{ AlwaysActive }

func (Card2121007FireDanceSkirt) ID() string   { return "2121007" }
func (Card2121007FireDanceSkirt) Name() string { return "舞火战裙" }

func (Card2121007FireDanceSkirt) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Card.Category == "火" && hasAnyNegativeStatus(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "fire_dance_skirt_purify",
		"选择1个友方火焰单位移除所有负面状态", candidates, 1, 1,
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
