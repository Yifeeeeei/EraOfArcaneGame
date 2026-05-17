package game

import "eraofarcane/model"

const phoenixFeatherCounter = "火焰标记物"

type Card2121001PhoenixFeather struct{ AlwaysActive }

func (Card2121001PhoenixFeather) ID() string   { return "2121001" }
func (Card2121001PhoenixFeather) Name() string { return "凤凰之羽" }

func (Card2121001PhoenixFeather) OnEnter(ctx *EffectContext) error {
	ctx.Source.Statuses[phoenixFeatherCounter] += 3
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":  cardToInfo(ctx.Source),
		"effect":  "add_counter",
		"counter": phoenixFeatherCounter,
		"amount":  3,
	}})
	return nil
}

func (Card2121001PhoenixFeather) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source.Statuses[phoenixFeatherCounter] <= 0 {
		return nil
	}
	ctx.Source.Statuses[phoenixFeatherCounter]--
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementFire: 1})
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":  cardToInfo(ctx.Source),
		"effect":  "gain_element",
		"element": model.ElementFire,
		"amount":  1,
	}})
	return nil
}
