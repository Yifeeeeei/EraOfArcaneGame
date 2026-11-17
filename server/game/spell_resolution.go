package game

import (
	"fmt"
)

// handleNoDefend handles when the defender chooses not to defend
func (e *Engine) handleNoDefend(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseDefenseWindow {
		return fmt.Errorf("not in defense window")
	}
	if e.State.PendingSpell == nil {
		return fmt.Errorf("no pending spell")
	}
	if playerID == e.State.PendingSpell.AttackerID {
		return fmt.Errorf("attacker cannot respond here")
	}

	e.resolvePendingSpellHit()
	return nil
}

func (e *Engine) resolvePendingSpellHit() {
	if e.State.PendingSpell == nil {
		return
	}
	spell := e.State.PendingSpell
	if e.resolveSpellHit(
		spell.AttackerID,
		spell.Skill,
		spell.Target,
		spell.BoostSkills,
		spell.ExtraTargets,
	) {
		return
	}

	e.removeStoredArchmageStaffSkillAfterUse(spell.AttackerID, spell.Skill)
	clearFiveRainbowBeamSelection(spell.Skill)
	continued := e.completePendingSpell(spell)
	if !continued && e.State.PendingAction == nil {
		e.State.Phase = PhaseMain
	}
	e.checkWinCondition()
}

func (e *Engine) cancelPendingSpell(playerID int, source *CardInstance, reason string) {
	if e.State.PendingSpell == nil {
		return
	}
	spell := e.State.PendingSpell
	e.emit(GameEvent{
		Type:   "spell_cancelled",
		Player: -1,
		Data: map[string]any{
			"player": playerID,
			"card":   cardToInfo(source),
			"reason": reason,
		},
	})
	if e.promptSpellMissOrCancelledCounters(spell.AttackerID, spell.Skill, spell.BoostSkills, spell.ExtraTargets, reason) {
		return
	}
	e.removeStoredArchmageStaffSkillAfterUse(spell.AttackerID, spell.Skill)
	clearFiveRainbowBeamSelection(spell.Skill)
	continued := e.completePendingSpell(spell)
	if !continued && e.State.PendingAction == nil {
		e.State.Phase = PhaseMain
	}
}

func (e *Engine) promptSpellMissOrCancelledCounters(attackerID int, skill *CardInstance, boostSkills []*CardInstance, extraTargets []SpellTarget, reason string) bool {
	if e == nil || e.State == nil || skill == nil || attackerID < 0 || attackerID >= len(e.State.Players) {
		return false
	}
	defenderID := 1 - attackerID
	if defenderID < 0 || defenderID >= len(e.State.Players) {
		return false
	}
	data := map[string]any{
		"attacker":      attackerID,
		"cast_player":   attackerID,
		"reason":        reason,
		"boost_skills":  boostSkills,
		"extra_targets": extraTargets,
	}
	counters := e.eligibleCounterTraps(defenderID, TriggerOnSpellMissOrCancelled, skill, data)
	if len(counters) == 0 {
		return false
	}
	return e.promptCounterTrapQueue(counters, TriggerOnSpellMissOrCancelled, skill, data, func() {
		if e.State.PendingSpell != nil && e.State.PendingSpell.Skill == skill {
			spell := e.State.PendingSpell
			e.removeStoredArchmageStaffSkillAfterUse(spell.AttackerID, spell.Skill)
			clearFiveRainbowBeamSelection(spell.Skill)
			continued := e.completePendingSpell(spell)
			if !continued && e.State.PendingAction == nil {
				e.State.Phase = PhaseMain
			}
		}
	})
}

func (e *Engine) spellAllowsDefense(attackerID int, skill *CardInstance, target SpellTarget) bool {
	return e.spellDefenderID(attackerID, skill, target) != attackerID
}

func spellTargetsContain(targets []SpellTarget, target SpellTarget) bool {
	for _, existing := range targets {
		if existing.Type != target.Type || existing.Position != target.Position {
			continue
		}
		if existing.OwnerID == nil && target.OwnerID == nil {
			return true
		}
		if existing.OwnerID != nil && target.OwnerID != nil && *existing.OwnerID == *target.OwnerID {
			return true
		}
	}
	return false
}

func (e *Engine) spellDefenderID(attackerID int, skill *CardInstance, target SpellTarget) int {
	defenderID := 1 - attackerID
	if target.OwnerID != nil && *target.OwnerID == attackerID {
		defenderID = attackerID
	}
	if skill != nil && skill.Card != nil {
		if friendly, ok := behaviorForNumber(skill.Card.Number).(FriendlySpellTargetBehavior); ok && friendly.HasActiveFriendlySpellTarget(skill) && friendly.AllowsFriendlySpellTarget() && target.Type == "unit" && target.Position.Valid() {
			if e.State.Players[attackerID].Units[target.Position.Col][target.Position.Row] != nil {
				defenderID = attackerID
			}
		}
	}
	if target.Type == "hero" {
		defenderID = attackerID
	}
	return defenderID
}

// resolveSpellHit applies spell damage to the target. It returns true when a
// pre-hit counter prompt delayed resolution.
func (e *Engine) resolveSpellHit(attackerID int, skill *CardInstance, target SpellTarget, boostSkills []*CardInstance, extraTargets []SpellTarget) bool {
	e.beginResolution()
	defer e.endResolution()

	defenderID := e.spellDefenderID(attackerID, skill, target)
	affectedUnits := e.spellAffectedUnits(defenderID, skill, target)
	for _, extraTarget := range extraTargets {
		if extraTarget.Type != "unit" || !extraTarget.Position.Valid() {
			continue
		}
		extraUnit := e.spellTargetUnitForCaster(attackerID, extraTarget)
		if extraUnit == nil {
			continue
		}
		if target.Type == "unit" && extraTarget.Position == target.Position {
			affectedUnits = append(affectedUnits, extraUnit)
			continue
		}
		alreadyIncluded := false
		for _, unit := range affectedUnits {
			if unit == extraUnit {
				alreadyIncluded = true
				break
			}
		}
		if !alreadyIncluded {
			affectedUnits = append(affectedUnits, extraUnit)
		}
	}
	if target.Type == "unit" && len(affectedUnits) == 0 {
		e.emit(GameEvent{
			Type:   "spell_miss",
			Player: -1,
			Data: map[string]any{
				"attacker": attackerID,
				"skill":    cardToInfo(skill),
				"target":   target,
				"reason":   "target_lost",
			},
		})
		if e.promptSpellMissOrCancelledCounters(attackerID, skill, boostSkills, extraTargets, "target_lost") {
			return true
		}
		return false
	}
	var targetUnit *CardInstance
	if target.Type == "hero" {
		targetUnit = e.State.Players[attackerID].Hero
	}
	if len(affectedUnits) > 0 {
		targetUnit = affectedUnits[0]
	}
	ctx := &EffectContext{
		Engine:     e,
		Source:     skill,
		Target:     targetUnit,
		PlayerID:   attackerID,
		OpponentID: defenderID,
		ExtraData:  map[string]any{"target": target},
	}
	dmg := max(skill.Card.Attack+skill.AttackBonus, 0)
	if override, ok := globalRegistry.SpellDamage(skill.Card.Number, ctx); ok {
		dmg = max(override, 0)
	}
	dmg = e.effectiveSpellDamage(attackerID, skill, dmg, boostSkills, affectedUnits)
	e.consumeNextElementSpellDamageBonus(e.State.Players[attackerID], skill)
	e.consumeAllSpellDamageZero(e.State.Players[attackerID], skill)
	e.consumeAllSpellDamageZero(e.State.Players[defenderID], skill)
	e.consumeFriendlySpellDamageMinus(e.State.Players[defenderID], skill)

	{
		totalPower := e.effectiveSpellPower(attackerID, skill, boostSkills, target)
		if e.State.PendingSpell != nil && e.State.PendingSpell.Skill == skill {
			totalPower = e.State.PendingSpell.TotalPower
		}
		hitCancelled := false
		hitData := map[string]any{
			"damage":           dmg,
			"power":            totalPower,
			"attacker":         attackerID,
			"spell_source":     skill,
			"target":           target,
			"affected_units":   affectedUnits,
			"boost_skills":     boostSkills,
			"cancel_spell_hit": &hitCancelled,
			"damage_ptr":       &dmg,
		}
		finishHit := func() bool {
			if hitCancelled {
				if e.promptSpellMissOrCancelledCounters(attackerID, skill, boostSkills, extraTargets, "hit_cancelled") {
					return true
				}
				return false
			}
			e.emit(GameEvent{
				Type:   "spell_hit",
				Player: -1,
				Data: map[string]any{
					"attacker": attackerID,
					"skill":    cardToInfo(skill),
					"target":   target,
					"damage":   dmg,
				},
			})

			hitData["skip_counter_traps"] = true
			hitData["timing"] = "before_damage"
			e.triggerEffects(TriggerOnSpellHitBeforeDamage, skill, targetUnit, hitData)
			e.triggerFieldEffectsWithData(TriggerOnSpellHitBeforeDamage, attackerID, skill, hitData)
			if !spellSuppressesOpponentResponses(skill) {
				e.triggerFieldEffectsWithData(TriggerOnSpellHitBeforeDamage, defenderID, skill, hitData)
			}
			actualSpellDamageByInstance := map[string]int{}
			actualFriendlySpellDamageByInstance := map[string]int{}
			e.runResolution("spell hit",
				func() {
					if hitCancelled {
						return
					}
					e.consumeNextSpellAttackBonuses(e.State.Players[attackerID], skill)
					if dmg > 0 {
						request := DamageRequest{Source: skill, Kind: "spell", Element: skill.Card.Category, Spell: skill.Card,
							BoostCount: len(boostSkills), ActualDamage: actualSpellDamageByInstance, ActualFriendlyDamage: actualFriendlySpellDamageByInstance}
						spellDamage := dmg
						if len(affectedUnits) > 1 {
							shieldTarget := targetUnit
							if shieldTarget == nil && len(affectedUnits) > 0 {
								shieldTarget = affectedUnits[0]
							}
							spellDamage = e.applyPlayerShieldDamage(shieldTarget, dmg, request.triggerData())
							request.SkipPlayerShield = true
						}
						for _, unit := range affectedUnits {
							request.Target, request.Amount = unit, spellDamage
							e.ApplyDamage(request)
						}
					}
				},
				func() {
					if hitCancelled {
						return
					}
					hitData["actual_damage_by_instance"] = actualSpellDamageByInstance
					hitData["actual_friendly_damage_by_instance"] = actualFriendlySpellDamageByInstance
					resolvedUnits := e.unitsStillOnField(affectedUnits)
					actualDamage := 0
					for _, damage := range actualSpellDamageByInstance {
						actualDamage += damage
					}
					e.recordSpellHitStats(attackerID, len(affectedUnits), actualDamage)
					resolvedTargetUnit := targetUnit
					if target.Type != "hero" && !e.unitStillOnField(resolvedTargetUnit) {
						resolvedTargetUnit = nil
					}
					hitData["affected_units"] = resolvedUnits
					e.applyGenericSpellEffects(attackerID, defenderID, skill, resolvedUnits, target)
					e.applyTemporarySpellHitStatus(attackerID, skill, resolvedUnits)

					hitData["timing"] = "after_damage"
					e.triggerEffects(TriggerOnSpellHit, skill, resolvedTargetUnit, hitData)
					e.triggerFieldEffectsWithData(TriggerOnSpellHit, attackerID, skill, hitData)
					if !spellSuppressesOpponentResponses(skill) {
						e.triggerFieldEffectsWithData(TriggerOnSpellHit, defenderID, skill, hitData)
					}
					e.triggerSparkMothAfterSpellHit(skill)
					if skill.Statuses[StatusNextFrontRowRange] > 0 {
						skill.Statuses[StatusNextFrontRowRange]--
					}
				},
				func() { e.finishMatchingSpellHit(attackerID, skill) },
			)
			// This frame owns completion even when every step was synchronous.
			return true

		}
		afterCounterWindow := func() {
			if finishHit() {
				return
			}
			if e.State.PendingSpell != nil && e.State.PendingSpell.Skill == skill {
				spell := e.State.PendingSpell
				e.removeStoredArchmageStaffSkillAfterUse(attackerID, skill)
				continued := e.completePendingSpell(spell)
				if !continued && e.State.PendingAction == nil {
					e.State.Phase = PhaseMain
				}
				e.checkWinCondition()
			}
		}
		if !spellSuppressesOpponentResponses(skill) && e.promptCounterTrapQueue(e.eligibleCounterTraps(defenderID, TriggerOnSpellHitBeforeDamage, skill, hitData), TriggerOnSpellHitBeforeDamage, skill, hitData, afterCounterWindow) {
			return true
		}
		return finishHit()
	}
}

func (e *Engine) unitsStillOnField(units []*CardInstance) []*CardInstance {
	result := make([]*CardInstance, 0, len(units))
	for _, unit := range units {
		if e.unitStillOnField(unit) {
			result = append(result, unit)
		}
	}
	return result
}

func (e *Engine) unitStillOnField(unit *CardInstance) bool {
	if unit == nil || unit.CurrentLife <= 0 || unit.OwnerID < 0 || unit.OwnerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[unit.OwnerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] == unit {
				return true
			}
		}
	}
	return false
}

// finishMatchingSpellHit never completes a replacement spell created by an
// effect. Its caller has already waited for damage and post-hit choices.
func (e *Engine) finishMatchingSpellHit(attackerID int, skill *CardInstance) {
	spell := e.State.PendingSpell
	if spell == nil || spell.Skill != skill {
		return
	}
	e.removeStoredArchmageStaffSkillAfterUse(attackerID, skill)
	clearFiveRainbowBeamSelection(skill)
	continued := e.completePendingSpell(spell)
	if !continued && e.State.PendingAction == nil {
		e.State.Phase = PhaseMain
	}
	e.checkWinCondition()
}
