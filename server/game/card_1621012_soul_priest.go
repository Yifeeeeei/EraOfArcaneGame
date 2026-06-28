package game

type Card1621012SoulPriest struct{ AlwaysActive }

func (Card1621012SoulPriest) ID() string   { return "1621012" }
func (Card1621012SoulPriest) Name() string { return "灵魂祭司" }

func (Card1621012SoulPriest) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sacrifice_unit",
		"选择1个友方伙伴献祭，抽2张牌",
		candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := ctx.Engine.findUnitOnGrid(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
			if target != nil && !target.Card.IsHero() {
				ctx.Engine.destroyUnitWithCause(target, ctx.PlayerID, DeathCauseSacrifice)
				_ = DrawCards(2)(ctx)
			}
		})
	return nil
}
