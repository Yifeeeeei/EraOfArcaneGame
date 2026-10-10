package game

type Card1021109ChurchEnvoy struct{ AlwaysActive }

func (Card1021109ChurchEnvoy) ID() string { return "1021109" }

func (Card1021109ChurchEnvoy) Name() string { return "教廷特使" }

func (Card1021109ChurchEnvoy) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, hasAnyNegativeStatus)
	candidates = append(candidates, ctx.Engine.friendlyEquipment(ctx.PlayerID, hasAnyNegativeStatus)...)
	candidates = append(candidates, ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, hasAnyNegativeStatus)...)
	if len(candidates) == 0 {
		return nil
	}
	allowed := make(map[string]bool, len(candidates))
	for _, candidate := range candidates {
		if id, _ := candidate["instance_id"].(string); id != "" {
			allowed[id] = true
		}
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "church_envoy_purify",
		"教廷特使:选择1张友方卡牌移除全部负面效果", candidates, 1, 1,
		func(selected []string) {
			id := firstSelected(selected)
			if !allowed[id] {
				return
			}
			clearNegativeStatuses(ctx.Engine.findFriendlyCardIncludingBound(ctx.PlayerID, id))
		})
	return nil
}
