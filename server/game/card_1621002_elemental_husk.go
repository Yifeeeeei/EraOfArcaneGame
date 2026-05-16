package game

import "eraofarcane/model"

type Card1621002ElementalHusk struct{}

func (Card1621002ElementalHusk) ID() string   { return "1621002" }
func (Card1621002ElementalHusk) Name() string { return "元素躯壳" }

func (Card1621002ElementalHusk) OnDeath(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].Elements[model.ElementArcane]++
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":  cardToInfo(ctx.Source),
		"effect":  "gain_element",
		"element": model.ElementArcane,
		"amount":  1,
	}})
	return nil
}
