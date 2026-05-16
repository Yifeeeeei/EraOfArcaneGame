package game

import "eraofarcane/model"

type Card2621006SoulNecklace struct{}

func (Card2621006SoulNecklace) ID() string   { return "2621006" }
func (Card2621006SoulNecklace) Name() string { return "亡魂项链" }

func (Card2621006SoulNecklace) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx.Target == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	ps.GainElements(map[string]int{model.ElementShadow: 1})
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":   cardToInfo(ctx.Source),
		"target":   cardToInfo(ctx.Target),
		"effect":   "gain_element",
		"element":  model.ElementShadow,
		"amount":   1,
		"elements": ps.Elements,
	}})
	return nil
}
