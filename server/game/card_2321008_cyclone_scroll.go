package game

type Card2321008CycloneScroll struct{ AlwaysActive }

func (Card2321008CycloneScroll) ID() string { return "2321008" }

func (Card2321008CycloneScroll) Name() string { return "旋风卷轴" }

func (Card2321008CycloneScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyEquipment(ctx.PlayerID, func(card *CardInstance) bool {
		return card.Card.IsItem() && totalElementCost(card.Card.ElementsCost) < 5
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "cyclone_scroll_destroy_equipment",
		"选择1个敌方入场花费小于5的装备道具摧毁", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Engine.destroyEnemyEquipment(ctx.PlayerID, selected[0])
		})
	return nil
}
