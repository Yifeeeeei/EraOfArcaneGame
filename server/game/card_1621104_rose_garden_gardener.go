package game

type Card1621104RoseGardenGardener struct{ AlwaysActive }

func (Card1621104RoseGardenGardener) ID() string { return "1621104" }

func (Card1621104RoseGardenGardener) Name() string { return "蔷薇花园园丁" }

func (Card1621104RoseGardenGardener) OnFriendlyDeath(ctx *EffectContext) error {
	return triggerRoseGardenGardener(ctx)
}

func (Card1621104RoseGardenGardener) OnEnemyDeath(ctx *EffectContext) error {
	return triggerRoseGardenGardener(ctx)
}

func triggerRoseGardenGardener(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || !triggeredTurnAvailable(ctx.Source) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card != nil && card.CurrentLife < maxLife(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	if !useTriggeredTurn(ctx.Source) {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "rose_garden_gardener_heal",
		"蔷薇花园园丁:选择1个友方单位回复2血", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findUnitByInstanceID(firstSelected(selected))
			if target != nil && target.OwnerID == ctx.PlayerID && target.CurrentLife < maxLife(target) {
				ctx.Engine.healUnit(target, 2, ctx.Source)
			}
		})
	return nil
}
