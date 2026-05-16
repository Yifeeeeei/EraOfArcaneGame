package game

import "eraofarcane/model"

type Card2611001NecromancyStoneVoid struct{}

func (Card2611001NecromancyStoneVoid) ID() string   { return "2611001" }
func (Card2611001NecromancyStoneVoid) Name() string { return "死灵魔石 虚无" }

func (Card2611001NecromancyStoneVoid) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx.Target == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	addElementsGainBonus(ctx.Source, model.ElementShadow, 1)
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":  cardToInfo(ctx.Source),
		"target":  cardToInfo(ctx.Target),
		"effect":  "load_bonus",
		"element": model.ElementShadow,
		"amount":  1,
	}})
	return nil
}
