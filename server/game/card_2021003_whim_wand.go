package game

type Card2021003WhimWand struct{}

func (Card2021003WhimWand) ID() string   { return "2021003" }
func (Card2021003WhimWand) Name() string { return "随心魔杖" }

func (Card2021003WhimWand) OnConsume(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkills(ctx.PlayerID, lowCostSkill)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "reset_skill",
		"选择1个使用花费小于3的法术重置",
		candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			card, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			resetCard(card)
		})
	return nil
}
