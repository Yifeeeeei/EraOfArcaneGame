package game

type Card3021008Disarm struct{ AlwaysActive }

func (Card3021008Disarm) ID() string   { return "3021008" }
func (Card3021008Disarm) Name() string { return "缴械" }

func (Card3021008Disarm) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.ExtraData == nil || ctx.ExtraData["spell_source"] != ctx.Source {
		return nil
	}
	candidates := ctx.Engine.enemyEquipment(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "disarm_destroy_equipment",
		"选择1个敌方装备摧毁", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Engine.destroyEnemyEquipment(ctx.PlayerID, selected[0])
		})
	return nil
}
