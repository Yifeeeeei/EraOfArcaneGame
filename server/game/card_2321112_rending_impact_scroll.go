package game

type Card2321112RendingImpactScroll struct{ AlwaysActive }

func (Card2321112RendingImpactScroll) ID() string { return "2321112" }

func (Card2321112RendingImpactScroll) Name() string { return "撕裂冲击卷轴" }

func (Card2321112RendingImpactScroll) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "2321112" || ctx.ExtraData == nil {
		return nil
	}
	candidates := make([]map[string]any, 0)
	for _, unit := range spellHitAffectedUnitsFromData(ctx) {
		if unit != nil && unit.OwnerID == ctx.OpponentID && ctx.Engine.unitStillOnField(unit) {
			candidates = append(candidates, candidateInfo(unit, "unit", "enemy"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "rending_impact_distribute_damage",
		"撕裂冲击卷轴:选择目标范围内单位分配共计3点伤害", candidates, 1, min(3, len(candidates)),
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			allocations := map[string]int{}
			order := make([]string, 0, len(selected))
			for _, id := range selected {
				if allocations[id] == 0 {
					order = append(order, id)
				}
				allocations[id]++
			}
			for remaining := 3 - len(selected); remaining > 0 && len(order) > 0; remaining-- {
				allocations[order[0]]++
			}
			for _, id := range order {
				target := findEnemyCardCandidate(ctx.Engine, ctx.PlayerID, id, candidates)
				if target == nil || !ctx.Engine.unitStillOnField(target) {
					continue
				}
				ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: allocations[id], Kind: "effect", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
			}
		})
	return nil
}
