package game

type Card3621006DeadSoulBite struct{ AlwaysActive }

func (Card3621006DeadSoulBite) ID() string   { return "3621006" }
func (Card3621006DeadSoulBite) Name() string { return "死魂之噬" }

func (Card3621006DeadSoulBite) OnSpellHit(ctx *EffectContext) error {
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "dead_soul_bite_weaken",
		"选择最多3个敌方法术分配虚弱", candidates, 1, 3,
		func(selected []string) {
			for _, id := range selected {
				for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
					if skill != nil && skill.InstanceID == id {
						skill.Statuses[StatusWeaken]++
					}
				}
			}
		})
	return nil
}
