package game

type Card3321015StaticBarrier struct{ AlwaysActive }

func (Card3321015StaticBarrier) ID() string   { return "3321015" }
func (Card3321015StaticBarrier) Name() string { return "静电屏障" }

func (Card3321015StaticBarrier) OnDefend(ctx *EffectContext) error {
	success, _ := ctx.ExtraData["defense_success"].(bool)
	if success {
		return nil
	}
	attackerID, _ := ctx.ExtraData["attacker"].(int)
	attacker := ctx.Engine.State.Players[attackerID]
	row := attacker.GetFrontRow()
	if row < 0 {
		return nil
	}
	for col := 0; col < 3; col++ {
		unit := attacker.Units[col][row]
		if unit == nil {
			continue
		}
		if !ctx.Engine.addStatus(unit, StatusStun, 1) {
			continue
		}
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(unit),
			"effect": "apply_status",
			"status": StatusStun,
			"amount": 1,
		}})
		return nil
	}
	return nil
}
