package game

const (
	sturdyScrollShieldStatus      = "sturdy_scroll_shield"
	sturdyScrollShieldUntilStatus = "sturdy_scroll_shield_until"
)

type Card2421003SturdyScroll struct{ AlwaysActive }

func (Card2421003SturdyScroll) ID() string   { return "2421003" }
func (Card2421003SturdyScroll) Name() string { return "坚固卷轴" }
func (Card2421003SturdyScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sturdy_scroll_target", "坚固卷轴:选择1个友方单位免疫最多3点伤害", candidates, 1, 1, func(selected []string) {
		target, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
		if target == nil {
			return
		}
		target.Statuses[sturdyScrollShieldStatus] += 3
		target.Statuses[sturdyScrollShieldUntilStatus] = ctx.Engine.State.TurnNumber + 1
	})
	return nil
}
