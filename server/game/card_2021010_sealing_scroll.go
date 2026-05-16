package game

type Card2021010SealingScroll struct{}

func (Card2021010SealingScroll) ID() string   { return "2021010" }
func (Card2021010SealingScroll) Name() string { return "封印卷轴" }

func (Card2021010SealingScroll) OnUseItem(ctx *EffectContext) error {
	opponentSkills := ctx.Engine.enemySkills(ctx.PlayerID, nil)
	if len(opponentSkills) < 4 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sealing_scroll",
		"选择1个敌方法术封印到下个回合结束", opponentSkills, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.OpponentID], selected[0])
			if target != nil {
				target.Statuses[StatusSeal] = 2
			}
		})
	return nil
}
