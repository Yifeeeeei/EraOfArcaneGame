package game

import "eraofarcane/model"

type Card2521010ShiningCrystal struct{ AlwaysActive }

func (Card2521010ShiningCrystal) ID() string   { return "2521010" }
func (Card2521010ShiningCrystal) Name() string { return "闪耀水晶" }

func (Card2521010ShiningCrystal) OnSpellHit(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.Card == nil || !ctx.Target.Card.IsSkill() || ctx.Target.Card.Category != model.ElementLight {
		return nil
	}
	for _, unit := range affectedUnitsFromHit(ctx) {
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
	}
	return nil
}
