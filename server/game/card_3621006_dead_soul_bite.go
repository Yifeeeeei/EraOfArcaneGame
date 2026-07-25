package game

type Card3621006DeadSoulBite struct{ AlwaysActive }

func (Card3621006DeadSoulBite) ID() string   { return "3621006" }
func (Card3621006DeadSoulBite) Name() string { return "死魂之噬" }

func (Card3621006DeadSoulBite) OnSpellHit(ctx *EffectContext) error {
	if !isOwnSpellHit(ctx) {
		return nil
	}
	if ctx.Target != nil && ctx.Target.Card != nil && ctx.Target.Card.IsSkill() {
		return nil
	}
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, canInstanceBeWeakened)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "dead_soul_bite_weaken",
		"死魂之噬:选择最多3个敌方法术分配3层虚弱", candidates, 1, 3,
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
			for id, amount := range allocations {
				for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
					if skill != nil && skill.InstanceID == id {
						ctx.Engine.addStatus(skill, StatusWeaken, amount)
					}
				}
			}
		})
	return nil
}
