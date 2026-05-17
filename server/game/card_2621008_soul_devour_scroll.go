package game

type Card2621008SoulDevourScroll struct{ AlwaysActive }

func (Card2621008SoulDevourScroll) ID() string   { return "2621008" }
func (Card2621008SoulDevourScroll) Name() string { return "魂噬卷轴" }

func (Card2621008SoulDevourScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "soul_devour_weaken",
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
