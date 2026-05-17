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
		unit.Statuses[StatusStun]++
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(unit),
			"effect": "apply_status",
			"status": StatusStun,
			"amount": 1,
		}})
	}
	applyStun(ctx.Target)

	if ctx.Target.Position == nil {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.Target.OwnerID]
	for _, delta := range []Position{{Col: -1, Row: 0}, {Col: 1, Row: 0}, {Col: 0, Row: -1}, {Col: 0, Row: 1}} {
		pos := Position{Col: ctx.Target.Position.Col + delta.Col, Row: ctx.Target.Position.Row + delta.Row}
		if !pos.Valid() {
			continue
		}
		if adjacent := ps.Units[pos.Col][pos.Row]; adjacent != nil {
			applyStun(adjacent)
			break
		}
	}
	return nil
}
