package game

type Card1621005CursedGolem struct{ AlwaysActive }

func (Card1621005CursedGolem) ID() string { return "1621005" }

func (Card1621005CursedGolem) Name() string { return "诅咒魔像" }

func (Card1621005CursedGolem) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, canInstanceBeWeakened)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "cursed_golem_weaken",
		"诅咒魔像:选择1个敌方法术虚弱2", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
				if skill != nil && skill.InstanceID == selected[0] {
					ctx.Engine.addStatus(skill, StatusWeaken, 2)
				}
			}
		})
	return nil
}
