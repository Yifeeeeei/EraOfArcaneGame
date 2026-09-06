package game

type Card1021102SwordsmanshipTeacher struct{ AlwaysActive }

func (Card1021102SwordsmanshipTeacher) ID() string { return "1021102" }

func (Card1021102SwordsmanshipTeacher) Name() string { return "剑术师傅" }

func (Card1021102SwordsmanshipTeacher) OnEnter(ctx *EffectContext) error {
	candidates := adjacentFriendlyCompanions(ctx)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "swordsmanship_teacher_buff",
		"剑术师傅:选择1个相邻友方伙伴获得+1攻", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
			if target != nil && target.Card != nil && target.Card.IsCompanion() {
				target.AttackBonus++
			}
		})
	return nil
}
