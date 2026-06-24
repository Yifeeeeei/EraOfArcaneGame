package game

type Card2321002LightningRune struct{ AlwaysActive }

func (Card2321002LightningRune) ID() string   { return "2321002" }
func (Card2321002LightningRune) Name() string { return "闪电符文" }

func (Card2321002LightningRune) OnConsume(ctx *EffectContext) error {
	if ctx.Target == nil || (!ctx.Target.Card.IsHero() && !ctx.Target.Card.IsCompanion()) || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	applyStun := func(unit *CardInstance) {
		if unit == nil {
			return
		}
		if !ctx.Engine.addStatus(unit, StatusStun, 1) {
			return
		}
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(unit),
			"effect": "apply_status",
			"status": StatusStun,
			"amount": 1,
		}})
	}
	applyStun(ctx.Target)

	candidates := ctx.Engine.adjacentUnitCandidatesForCounter(ctx.PlayerID, ctx.Target)
	if len(candidates) == 0 {
		return nil
	}
	ownerID := ctx.Target.OwnerID
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lightning_rune_adjacent", "闪电符文:选择1个相邻单位晕眩", candidates, 1, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ownerID], selected[0])
		applyStun(target)
	})
	return nil
}
