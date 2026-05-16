package game

import "eraofarcane/model"

type Card3121014FlameRekindle struct{}

func (Card3121014FlameRekindle) ID() string   { return "3121014" }
func (Card3121014FlameRekindle) Name() string { return "烈焰重燃" }

func (Card3121014FlameRekindle) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	count := ctx.Engine.State.Players[ctx.PlayerID].SpellsCastThisTurn[model.ElementFire]
	if count <= 0 {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	ps.GainElements(map[string]int{model.ElementFire: count})
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":   cardToInfo(ctx.Source),
		"effect":   "gain_element",
		"element":  model.ElementFire,
		"amount":   count,
		"elements": ps.Elements,
	}})
	return nil
}
