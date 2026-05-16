package game

import "eraofarcane/model"

type Card3221015Blizzard struct{}

func (Card3221015Blizzard) ID() string   { return "3221015" }
func (Card3221015Blizzard) Name() string { return "暴风雪" }

func (Card3221015Blizzard) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.Target == nil || ctx.Target.Card == nil {
		return
	}
	if ctx.Target.Card.Category != model.ElementWater && ctx.Target.Card.Category != model.ElementAir {
		return
	}
	stats.PowerBonus++
}

func (Card3221015Blizzard) OnSpellHit(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.Card == nil || !ctx.Target.Card.IsSkill() {
		return nil
	}
	if ctx.Target.Card.Category != model.ElementWater && ctx.Target.Card.Category != model.ElementAir {
		return nil
	}
	for _, unit := range affectedUnitsFromHit(ctx) {
		unit.Statuses[StatusFreeze]++
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(unit),
			"effect": "apply_status",
			"status": StatusFreeze,
			"amount": 1,
		}})
	}
	return nil
}
