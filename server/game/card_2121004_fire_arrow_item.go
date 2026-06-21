package game

type Card2121004FireArrowItem struct{ AlwaysActive }

func (Card2121004FireArrowItem) ID() string   { return "2121004" }
func (Card2121004FireArrowItem) Name() string { return "火焰箭" }

func (Card2121004FireArrowItem) OnUseItem(ctx *EffectContext) error {
	return fireArrowSelectAndDamage(ctx, false)
}

func (Card2121004FireArrowItem) OnUltimate(ctx *EffectContext) error {
	return fireArrowSelectAndDamage(ctx, true)
}

func fireArrowSelectAndDamage(ctx *EffectContext, sacrificeSelf bool) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "fire_arrow_damage",
		"献祭火焰箭:选择1个敌方单位造成1点伤害", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			if sacrificeSelf && !ctx.Engine.sacrificeEquipment(ctx.PlayerID, ctx.Source.InstanceID) {
				return
			}
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target == nil {
				return
			}
			ctx.Engine.dealDamage(target, 1, ctx.OpponentID)
		})
	return nil
}
