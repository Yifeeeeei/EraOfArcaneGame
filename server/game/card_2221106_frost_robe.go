package game

import (
	"eraofarcane/model"
)

type Card2221106FrostRobe struct{ AlwaysActive }

func (Card2221106FrostRobe) ID() string { return "2221106" }

func (Card2221106FrostRobe) Name() string { return "凛冰法袍" }

func (Card2221106FrostRobe) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil || ctx.Source.UltimateUsed {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	if !ok || attacker == ctx.PlayerID || !friendlyWaterUnitTookSpellDamage(ctx) {
		return nil
	}
	frozen := 0
	for _, candidate := range ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card != nil && card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	}) {
		id, _ := candidate["instance_id"].(string)
		target := findEnemyCardCandidate(ctx.Engine, ctx.PlayerID, id, []map[string]any{candidate})
		if target != nil && ctx.Engine.addStatus(target, StatusFreeze, 1) {
			frozen++
		}
	}
	if frozen == 0 {
		return nil
	}
	ctx.Source.UltimateUsed = true
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"effect": "frost_robe_freeze_enemies",
		"count":  frozen,
	}})
	return nil
}

func friendlyWaterUnitTookSpellDamage(ctx *EffectContext) bool {
	if ctx == nil || ctx.ExtraData == nil {
		return false
	}
	actualDamage, _ := ctx.ExtraData["actual_friendly_damage_by_instance"].(map[string]int)
	if len(actualDamage) == 0 {
		return false
	}
	for _, unit := range spellHitAffectedUnitsFromData(ctx) {
		if unit == nil || unit.OwnerID != ctx.PlayerID || unit.Card == nil || unit.Card.Category != model.ElementWater {
			continue
		}
		if actualDamage[unit.InstanceID] > 0 {
			return true
		}
	}
	return false
}
