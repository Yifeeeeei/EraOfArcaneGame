package game

func (Card2321011TeleportRune) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnUnitEnter, TriggerOnConsume}
}

func (Card2321011TeleportRune) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Card != nil && ctx.Event.Card.Card.IsCompanion() && ctx.Event.Card.Position != nil &&
		len(ctx.Engine.emptyUnitPositionsForPlayer(ctx.Event.Card.OwnerID, ctx.Source.OwnerID)) > 0
}

type Card2321011TeleportRune struct{ AlwaysActive }

func (Card2321011TeleportRune) ID() string { return "2321011" }

func (Card2321011TeleportRune) Name() string { return "传送符文" }

func (Card2321011TeleportRune) OnUseItem(ctx *EffectContext) error {
	target := ctx.Target
	if target == nil || target.Card == nil || !target.Card.IsCompanion() || target.Position == nil {
		return nil
	}
	positions := ctx.Engine.emptyUnitPositionsForPlayer(target.OwnerID, ctx.PlayerID)
	if len(positions) == 0 {
		return nil
	}
	targetID := target.InstanceID
	targetOwner := target.OwnerID
	ctx.Engine.SetPendingAction(ctx.PlayerID, "teleport_rune_position",
		"Teleport Rune: choose another empty position", positions, 1, 1,
		func(selected []string) {
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok {
				return
			}
			ctx.Engine.moveUnitToPosition(targetOwner, targetID, pos)
		})
	return nil
}
