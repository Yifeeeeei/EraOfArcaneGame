package game

type Card2611002DemonContract struct{ AlwaysActive }

func (Card2611002DemonContract) ID() string   { return "2611002" }
func (Card2611002DemonContract) Name() string { return "与恶魔的契约书" }

func (Card2611002DemonContract) OnUseItem(ctx *EffectContext) error {
	sacrifices := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion()
	})
	if len(sacrifices) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "demon_contract_sacrifice", "献祭1个友方单位", sacrifices, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			sacrifice := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
			if sacrifice == nil {
				return
			}
			ctx.Engine.destroyUnit(sacrifice, ctx.PlayerID)
			targets := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
				return card.Card.IsCompanion()
			})
			ctx.Engine.SetPendingAction(ctx.PlayerID, "demon_contract_destroy", "选择1个敌方伙伴消灭", targets, 1, 1,
				func(selected []string) {
					if len(selected) == 0 {
						return
					}
					target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.OpponentID], selected[0])
					if target != nil {
						ctx.Engine.destroyUnit(target, ctx.OpponentID)
					}
				})
		})
	return nil
}
