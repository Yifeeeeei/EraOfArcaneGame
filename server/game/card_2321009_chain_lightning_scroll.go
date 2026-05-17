package game

type Card2321009ChainLightningScroll struct{ AlwaysActive }

func (Card2321009ChainLightningScroll) ID() string   { return "2321009" }
func (Card2321009ChainLightningScroll) Name() string { return "连锁闪电卷轴" }

func (Card2321009ChainLightningScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return DrawCards(1)(ctx)
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "chain_lightning_scroll",
		"选择1个敌方单位造成1点伤害,然后抽1张牌", candidates, 1, 1,
		func(selected []string) {
			target := findEnemyByID(ctx, selected)
			if target != nil {
				ctx.Engine.dealDamage(target, 1, ctx.OpponentID)
			}
			_ = DrawCards(1)(ctx)
		})
	return nil
}
