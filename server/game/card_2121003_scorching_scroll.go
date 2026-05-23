package game

type Card2121003ScorchingScroll struct{ AlwaysActive }

func (Card2121003ScorchingScroll) ID() string   { return "2121003" }
func (Card2121003ScorchingScroll) Name() string { return "灼烧卷轴" }

func (Card2121003ScorchingScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "scorching_scroll_burn",
		"灼烧卷轴:选择1个敌方单位施加点燃1", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := ctx.Engine.findCardOnField(ctx.Engine.State.Players[ctx.OpponentID], selected[0])
			if target != nil {
				target.Statuses[StatusBurn]++
			}
		})
	return nil
}
