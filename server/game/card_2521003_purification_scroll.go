package game

var negativeStatuses = []string{StatusBurn, StatusFreeze, StatusStun, StatusPetrify, StatusWeaken}

// These are counters explicitly created by card text. Do not include mastery,
// cooldown, seal, temporary flags, links, or other internal state.
var cardEffectMarkerStatuses = []string{
	phoenixFeatherCounter,
	arcaneCylinderCounter,
	fireBoxCounter,
	waterAriaCounter,
	windQuillCounter,
	forestStorageCounter,
	blessingStaffCounter,
	burierCounter,
	nurEyeFireMark,
	winterBowWaterMark,
	"火焰标记",
	"地脉标记",
	"暗影标记",
	"雷鼓标记",
}

type Card2521003PurificationScroll struct{ AlwaysActive }

func (Card2521003PurificationScroll) ID() string   { return "2521003" }
func (Card2521003PurificationScroll) Name() string { return "净化卷轴" }

func (Card2521003PurificationScroll) OnUseItem(ctx *EffectContext) error {
	purifyTargets := ctx.Engine.friendlyUnits(ctx.PlayerID, true, hasAnyNegativeStatus)
	markerTargets := enemyCardsWithMarkers(ctx)
	choices := make([]map[string]any, 0, 2)
	if len(purifyTargets) > 0 {
		choices = append(choices, map[string]any{
			"instance_id": "purify",
			"number":      "2521003",
			"name":        "移除1个友方卡牌所有负面状态",
			"type":        "选择",
			"zone":        "choice",
			"side":        "own",
		})
	}
	if len(markerTargets) > 0 {
		choices = append(choices, map[string]any{
			"instance_id": "remove_markers",
			"number":      "2521003",
			"name":        "移除任意1个敌方卡牌所有标记物",
			"type":        "选择",
			"zone":        "choice",
			"side":        "own",
		})
	}
	if len(choices) == 0 {
		return nil
	}

	ctx.Engine.SetPendingAction(ctx.PlayerID, "purification_scroll_mode", "净化卷轴:选择效果", choices, 1, 1,
		func(selected []string) {
			switch firstSelected(selected) {
			case "purify":
				ctx.Engine.SetPendingAction(ctx.PlayerID, "purify_friendly", "选择1个友方卡牌移除所有负面状态", purifyTargets, 1, 1,
					func(selected []string) {
						target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
						if target == nil || zone != "unit" {
							return
						}
						clearNegativeStatuses(target)
						ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
							"source": cardToInfo(ctx.Source),
							"target": cardToInfo(target),
							"effect": "purify",
						}})
					})
			case "remove_markers":
				ctx.Engine.SetPendingAction(ctx.PlayerID, "purification_scroll_remove_markers", "选择1个敌方卡牌移除所有标记物", markerTargets, 1, 1,
					func(selected []string) {
						target := findEnemyCardCandidate(ctx.Engine, ctx.PlayerID, firstSelected(selected), markerTargets)
						if target == nil {
							return
						}
						clearCardEffectMarkers(target)
						ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
							"source": cardToInfo(ctx.Source),
							"target": cardToInfo(target),
							"effect": "remove_markers",
						}})
					})
			}
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

func clearNegativeStatuses(card *CardInstance) {
	if card == nil {
		return
	}
	for _, status := range negativeStatuses {
		delete(card.Statuses, status)
	}
}

func hasAnyCardEffectMarker(card *CardInstance) bool {
	if card == nil {
		return false
	}
	for _, status := range cardEffectMarkerStatuses {
		if card.Statuses[status] > 0 {
			return true
		}
	}
	return false
}

func clearCardEffectMarkers(card *CardInstance) {
	if card == nil {
		return
	}
	for _, status := range cardEffectMarkerStatuses {
		delete(card.Statuses, status)
	}
}

func enemyCardsWithMarkers(ctx *EffectContext) []map[string]any {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, hasAnyCardEffectMarker)
	candidates = append(candidates, ctx.Engine.enemyEquipment(ctx.PlayerID, hasAnyCardEffectMarker)...)
	candidates = append(candidates, ctx.Engine.enemySkills(ctx.PlayerID, hasAnyCardEffectMarker)...)
	return candidates
}

func findEnemyCardCandidate(e *Engine, playerID int, instanceID string, candidates []map[string]any) *CardInstance {
	if e == nil || instanceID == "" || !candidateContains(candidates, instanceID) {
		return nil
	}
	ps := e.State.Players[1-playerID]
	if ps == nil {
		return nil
	}
	for _, card := range ps.Equipment {
		if card != nil && card.InstanceID == instanceID {
			return card
		}
	}
	for _, card := range ps.Skills {
		if card != nil && card.InstanceID == instanceID {
			return card
		}
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit != nil && unit.InstanceID == instanceID {
				return unit
			}
		}
	}
	return nil
}
