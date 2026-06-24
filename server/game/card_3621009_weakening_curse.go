package game

type Card3621009WeakeningCurse struct{ AlwaysActive }

func (Card3621009WeakeningCurse) ID() string   { return "3621009" }
func (Card3621009WeakeningCurse) Name() string { return "虚弱诅咒" }

func (Card3621009WeakeningCurse) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "weakening_curse",
		"选择1个敌方法术虚弱2", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
				if skill != nil && skill.InstanceID == selected[0] {
					ctx.Engine.addStatus(skill, StatusWeaken, 2)
					return
				}
			}
		})
	return nil
}
