package game

import (
	"eraofarcane/model"
)

type Card3121008FireBarrier struct{ AlwaysActive }

func (Card3121008FireBarrier) ID() string { return "3121008" }

func (Card3121008FireBarrier) Name() string { return "火焰结界" }

func (Card3121008FireBarrier) HasActiveSpellStatModifier(card *CardInstance) bool {
	return abilityDurationActive(card)
}

func (Card3121008FireBarrier) HasActiveSpellHit(card *CardInstance) bool {
	return abilityDurationActive(card)
}

func (Card3121008FireBarrier) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.Target == nil || ctx.Target.Card == nil || ctx.Target.Card.Category != model.ElementFire {
		return
	}
	stats.PowerBonus += 2
}

func (Card3121008FireBarrier) OnSpellHit(ctx *EffectContext) error {
	if !fireBarrierIncludesFireSpell(ctx) {
		return nil
	}
	for _, unit := range affectedUnitsFromHit(ctx) {
		ctx.Engine.addStatus(unit, StatusBurn, 1)
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(unit),
			"effect": "apply_status",
			"status": StatusBurn,
			"amount": 1,
		}})
	}
	return nil
}

func fireBarrierIncludesFireSpell(ctx *EffectContext) bool {
	if ctx == nil {
		return false
	}
	if ctx.Target != nil && ctx.Target.Card != nil && ctx.Target.Card.IsSkill() && ctx.Target.Card.Category == model.ElementFire {
		return true
	}
	if ctx.ExtraData == nil {
		return false
	}
	boosts, _ := ctx.ExtraData["boost_skills"].([]*CardInstance)
	for _, boost := range boosts {
		if boost != nil && boost.Card != nil && boost.Card.IsSkill() && boost.Card.Category == model.ElementFire {
			return true
		}
	}
	return false
}

func affectedUnitsFromHit(ctx *EffectContext) []*CardInstance {
	if ctx == nil || ctx.ExtraData == nil {
		return nil
	}
	units, _ := ctx.ExtraData["affected_units"].([]*CardInstance)
	return units
}
