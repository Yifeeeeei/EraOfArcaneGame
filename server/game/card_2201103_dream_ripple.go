package game

type Card2201103DreamRipple struct{ AlwaysActive }

func (Card2201103DreamRipple) ID() string { return "2201103" }

func (Card2201103DreamRipple) Name() string { return "幻创之梦-波纹" }

func (Card2201103DreamRipple) OnUseItem(ctx *EffectContext) error {
	candidates := frontRowEnemyCandidates(ctx)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "dream_ripple_damage",
		"幻创之梦-波纹:选择前排敌人分配共计3点伤害", candidates, 1, min(3, len(candidates)),
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
				target := ctx.Engine.findUnitByInstanceID(id)
				if target == nil || target.OwnerID != ctx.OpponentID || target.Position == nil || !isCurrentFrontRowUnit(ctx.Engine.State.Players[ctx.OpponentID], target) {
					continue
				}
				ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: allocations[id], Kind: "effect", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
			}
		})
	return nil
}

func frontRowEnemyCandidates(ctx *EffectContext) []map[string]any {
	if ctx == nil || ctx.Engine == nil || ctx.OpponentID < 0 || ctx.OpponentID >= len(ctx.Engine.State.Players) {
		return nil
	}
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 || frontRow >= 3 {
		return nil
	}
	candidates := make([]map[string]any, 0, 3)
	for col := 0; col < 3; col++ {
		if unit := opponent.Units[col][frontRow]; unit != nil {
			candidates = append(candidates, candidateInfo(unit, "unit", "enemy"))
		}
	}
	return candidates
}

func isCurrentFrontRowUnit(ps *PlayerState, card *CardInstance) bool {
	if ps == nil || card == nil || card.Position == nil {
		return false
	}
	frontRow := ps.GetFrontRow()
	return frontRow >= 0 && card.Position.Row == frontRow
}
