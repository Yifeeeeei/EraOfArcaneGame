package game

type Card1621004WitchPriest struct{ AlwaysActive }

func (Card1621004WitchPriest) ID() string   { return "1621004" }
func (Card1621004WitchPriest) Name() string { return "巫术祭司" }

func (Card1621004WitchPriest) OnUltimate(ctx *EffectContext) error {
	sacrifices := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != ctx.Source && card.Card.IsCompanion()
	})
	if len(sacrifices) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "witch_priest_sacrifice",
		"献祭1个伙伴", sacrifices, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			sacrifice := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
			if sacrifice == nil {
				return
			}
			life := max(sacrifice.CurrentLife, 0)
			ctx.Engine.destroyUnit(sacrifice, ctx.PlayerID)
			targets := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
				return card != sacrifice
			})
			ctx.Engine.SetPendingAction(ctx.PlayerID, "witch_priest_heal",
				"选择另一个角色获得生命值", targets, 1, 1,
				func(selected []string) {
					if len(selected) == 0 {
						return
					}
					target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
					if target != nil {
						target.CurrentLife += life
					}
				})
		})
	return nil
}
