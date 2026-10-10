package game

type Card3021006InsightEye struct{ AlwaysActive }

func (Card3021006InsightEye) ID() string { return "3021006" }

func (Card3021006InsightEye) Name() string { return "洞察之眼" }

func (Card3021006InsightEye) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.enemyEquipment(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "insight_eye_destroy_equipment",
		"选择1张敌方盖放的卡牌摧毁", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Engine.destroyEnemyEquipment(ctx.PlayerID, selected[0])
		})
	return nil
}
