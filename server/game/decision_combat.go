package game

func (e *Engine) decisionTargetVariants(player int, skill *CardInstance, target SpellTarget) []SpellTarget {
	if e.spellAreaForTarget(skill, target) != SpellAreaSquare {
		return []SpellTarget{target}
	}
	var targets []SpellTarget
	for col := 0; col < 2; col++ {
		for row := 0; row < 2; row++ {
			pos := Position{Col: col, Row: row}
			copy := target
			copy.AreaOrigin = &pos
			if e.validateSpellTarget(player, skill, copy) == nil {
				targets = append(targets, copy)
			}
		}
	}
	return targets
}

// This preview queries public, active reductions without committing their
// consumption. It does not predict optional counters, hidden cards or choices.
func (e *Engine) decisionSpellImpacts(player int, skill *CardInstance, raw int, units []*CardInstance) []DecisionImpact {
	var result []DecisionImpact
	for _, unit := range units {
		if unit == nil {
			continue
		}
		damage := raw
		data := (DamageRequest{Source: skill, Target: unit, Amount: raw, Kind: "spell", Element: skill.Card.Category, Spell: skill.Card}).triggerData()
		data["damage"] = damage
		ctx := &EffectContext{Engine: e, Source: unit, Target: unit, PlayerID: unit.OwnerID, OpponentID: 1 - unit.OwnerID, ExtraData: data}
		prevented := e.temporaryDamageAndNegativeImmunityActive(unit) || unit.Statuses[scatterAwayUntilPlayerStartStatus] > 0
		if h, ok := cardBehavior(unit).(DamagePreventionBehavior); ok && h.HasActiveDamagePrevention(unit) && h.PreventsDamage(ctx) {
			prevented = true
		}
		if e.fieldDamagePreventionSource(unit, unit.OwnerID, data) != nil {
			prevented = true
		}
		if prevented {
			damage = 0
		} else {
			event := damageEventFromContext(ctx)
			event.Amount = damage
			damage = e.adjustedDamage(event)
			if unit.Statuses[sturdyScrollShieldUntilStatus] >= e.State.TurnNumber {
				damage = max(0, damage-unit.Statuses[sturdyScrollShieldStatus])
			}
			if player != unit.OwnerID {
				damage = max(0, damage-e.State.Players[unit.OwnerID].Shield)
			}
			damage = e.planDamageModifiers(unit, damage, unit.OwnerID, data).Amount
			if damage >= unit.CurrentLife && (e.lightShelterActive(unit) || unit.Statuses["防止致命"] > 0 && unit.Statuses[lightShelterUntilStatus] == 0) {
				damage = 0
			}
		}
		impact := DecisionImpact{InstanceID: unit.InstanceID, Damage: damage}
		for status, amount := range e.explicitSpellHitStatuses(skill, unit) {
			if e.rejectsNegativeStatusApplication(unit, status) {
				continue
			}
			if status == StatusBurn {
				impact.Burn = amount
			}
			if status == StatusStun {
				impact.Stun = amount
			}
		}
		result = append(result, impact)
	}
	return result
}

func (e *Engine) decisionHasAura(card *CardInstance) bool {
	if card == nil || card.IsSetCounter || e.hasEffectiveStatus(card, StatusPetrify) {
		return false
	}
	h := cardBehavior(card)
	if m, ok := h.(SpellStatModifier); ok && m.HasActiveSpellStatModifier(card) {
		return true
	}
	if m, ok := h.(SkillUseCostModifier); ok && m.HasActiveSkillUseCostModifier(card) {
		return true
	}
	if m, ok := h.(FieldDamagePreventionBehavior); ok && m.HasActiveFieldDamagePrevention(card) {
		return true
	}
	if m, ok := h.(SlotGrantBehavior); ok && m.HasActiveSlotGrant(card) {
		return true
	}
	return false
}
