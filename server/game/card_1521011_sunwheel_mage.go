package game

type Card1521011SunwheelMage struct{ AlwaysActive }

func (Card1521011SunwheelMage) ID() string   { return "1521011" }
func (Card1521011SunwheelMage) Name() string { return "日轮法师" }

func (Card1521011SunwheelMage) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkills(ctx.PlayerID, isLightSkill)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "reset_skill",
		"选择1个光辉法术重置",
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
