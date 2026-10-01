package game

type Card3021007ElementalEnchant struct{ AlwaysActive }

func (Card3021007ElementalEnchant) ID() string { return "3021007" }

func (Card3021007ElementalEnchant) Name() string { return "元素附魔" }

func (Card3021007ElementalEnchant) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := []map[string]any{
		statusChoice(StatusBurn, "点燃"),
		statusChoice(StatusFreeze, "冻结"),
		statusChoice(StatusStun, "眩晕"),
		statusChoice(StatusPetrify, "石化"),
		statusChoice(StatusWeaken, "虚弱"),
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "elemental_enchant_status",
		"选择下一次法术命中附加的负面效果", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
				Type:          TempModNextSpellHitStatus,
				RemainingUses: 1,
				Status:        selected[0],
				Amount:        1,
			})
		})
	return nil
}

func statusChoice(status string, name string) map[string]any {
	return map[string]any{
		"instance_id": status,
		"number":      "3021007",
		"name":        name,
		"type":        "状态",
		"zone":        "choice",
		"side":        "own",
	}
}
