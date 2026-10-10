package game

import (
	"fmt"
)

// handleDefend handles the defender's response to a spell
func (e *Engine) handleDefend(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseDefenseWindow {
		return fmt.Errorf("not in defense window")
	}
	if e.State.PendingSpell == nil {
		return fmt.Errorf("no pending spell")
	}
	if playerID == e.State.PendingSpell.AttackerID {
		return fmt.Errorf("attacker cannot defend")
	}

	defenseIDsRaw, _ := action.Data["skill_ids"].([]any)
	defenseScrollIDsRaw, _ := action.Data["scroll_ids"].([]any)
	boostIDsRaw, _ := action.Data["boost_ids"].([]any)
	overexertIDsRaw, _ := action.Data["overexert_ids"].([]any)

	ps := e.State.Players[playerID]

	defenseIDs := stringsFromAnySlice(defenseIDsRaw)
	defenseScrollIDs := stringsFromAnySlice(defenseScrollIDsRaw)
	boostIDs := stringsFromAnySlice(boostIDsRaw)
	overexertIDs := stringsFromAnySlice(overexertIDsRaw)
	defenseSkills, _, err := e.collectSkillUses(ps, defenseIDs, skillPurposeDefend, nil)
	if err != nil {
		return err
	}
	usedIDs := skillIDSet(defenseSkills)
	defenseScrolls, _, err := e.collectDefenseScrollUses(ps, defenseScrollIDs, usedIDs)
	if err != nil {
		return err
	}
	usedIDs = mergeSkillIDSet(usedIDs, skillIDSet(defenseScrolls))
	boostSkillIDs, boostScrollIDs := e.splitBoostIDs(ps, boostIDs)
	boostSkills, _, err := e.collectSkillUses(ps, boostSkillIDs, skillPurposeDefenseBoost, usedIDs)
	if err != nil {
		return err
	}
	usedIDs = mergeSkillIDSet(usedIDs, skillIDSet(boostSkills))
	boostScrolls, _, err := e.collectHandBoostScrollUses(ps, boostScrollIDs, usedIDs, skillPurposeDefenseBoost)
	if err != nil {
		return err
	}
	overexertUnits, err := e.collectOverexertUnits(ps, overexertIDs)
	if err != nil {
		return err
	}
	costPlan := newActionCostPlan(ps)
	totalCost := make(map[string]int)
	for _, defenseSkill := range defenseSkills {
		totalCost = mergeElementCosts(totalCost, costPlan.skillUseCost(e, defenseSkill, skillPurposeDefend, nil))
	}
	for _, defenseScroll := range defenseScrolls {
		totalCost = mergeElementCosts(totalCost, costPlan.cardPlayCost(e, defenseScroll))
	}
	for _, boostSkill := range boostSkills {
		totalCost = mergeElementCosts(totalCost, costPlan.skillUseCost(e, boostSkill, skillPurposeDefenseBoost, nil))
	}
	for _, boostScroll := range boostScrolls {
		totalCost = mergeElementCosts(totalCost, costPlan.cardPlayCost(e, boostScroll))
	}
	if !e.canPayCostWithOverexertOptions(ps, totalCost, overexertUnits, e.playerHasLightWildcard(ps)) {
		return fmt.Errorf("not enough elements for defense")
	}
	defensePowerSources := append([]*CardInstance{}, defenseSkills...)
	defensePowerSources = append(defensePowerSources, boostSkills...)
	defensePowerSacrifice, defensePowerSacrificeSource, defensePowerSacrificeBonus, err := e.validateSpellPowerSacrificeForSources(playerID, defensePowerSources, action)
	if err != nil {
		return err
	}
	if len(defenseSkills)+len(defenseScrolls)+len(boostSkills)+len(boostScrolls) > 0 {
		if !e.payDefenseCostWithOptions(ps, totalCost, action, overexertUnits, e.playerHasLightWildcard(ps)) {
			return fmt.Errorf("invalid payment")
		}
		e.triggerErebosSoulChainMarkedOverexert(playerID, overexertUnits)
		e.destroyFuyeDoomedAfterExert(overexertUnits)
		if defensePowerSacrifice != nil && defensePowerSacrificeBonus > 0 {
			e.destroyUnitWithCause(defensePowerSacrifice, playerID, DeathCauseSacrifice)
		}
		tapSkills(defenseSkills)
		tapSkills(boostSkills)
		usedSkills := append([]*CardInstance{}, defenseSkills...)
		usedSkills = append(usedSkills, boostSkills...)
		for _, defenseSkill := range defenseSkills {
			e.consumeNextSkillUseModifiersForPurpose(ps, defenseSkill, skillPurposeDefend)
		}
		e.moveHandSpellScrollsToGraveyard(ps, defenseScrolls, "defense")
		for _, boostSkill := range boostSkills {
			e.consumeNextSkillUseModifiersForPurpose(ps, boostSkill, skillPurposeDefenseBoost)
		}
		e.moveHandSpellScrollsToGraveyard(ps, boostScrolls, "defense")
		e.advanceMasteryForUsedSkills(playerID, usedSkills...)
	}

	defenseSources := append([]*CardInstance{}, defenseSkills...)
	defenseSources = append(defenseSources, defenseScrolls...)
	boostSources := append([]*CardInstance{}, boostSkills...)
	boostSources = append(boostSources, boostScrolls...)

	continueAfterDefenseSpellCounters := func() {
		if e.State.PendingSpell == nil {
			return
		}
		if e.promptDispelDefenseSpellIfEligible(e.State.PendingSpell.AttackerID, playerID, defenseSources, boostSources, len(overexertUnits), defensePowerSacrificeSource, defensePowerSacrificeBonus) {
			return
		}
		e.finishDefenseResolution(playerID, defenseSources, boostSources, len(overexertUnits), defensePowerSacrificeSource, defensePowerSacrificeBonus)
	}
	if e.promptDefenseSpellCastCounters(e.State.PendingSpell.AttackerID, playerID, defenseSources, boostSources, overexertUnits, continueAfterDefenseSpellCounters) {
		return nil
	}

	continueAfterDefenseSpellCounters()
	return nil
}

func (e *Engine) promptDefenseSpellCastCounters(attackerID int, defenderID int, defenseSources []*CardInstance, boostSources []*CardInstance, overexertUnits []*CardInstance, afterDone func()) bool {
	type defenseSpellSource struct {
		card    *CardInstance
		purpose skillPurpose
	}
	sources := make([]defenseSpellSource, 0, len(defenseSources)+len(boostSources))
	for _, source := range defenseSources {
		if source != nil && source.Card != nil && isSpellLikeCard(source.Card) {
			sources = append(sources, defenseSpellSource{card: source, purpose: skillPurposeDefend})
		}
	}
	for _, source := range boostSources {
		if source != nil && source.Card != nil && isSpellLikeCard(source.Card) {
			sources = append(sources, defenseSpellSource{card: source, purpose: skillPurposeDefenseBoost})
		}
	}
	var promptNext func(int, bool)
	promptNext = func(index int, continuing bool) {
		for index < len(sources) {
			source := sources[index]
			index++
			if source.card.Statuses[iceSoulSealCancelledBoostStatus] > 0 {
				continue
			}
			cancelled := false
			power := e.effectiveSkillPowerForPurpose(defenderID, source.card, source.purpose)
			data := map[string]any{
				"cast_player":     defenderID,
				"attacker":        defenderID,
				"skill":           cardToInfo(source.card),
				"power":           power,
				"purpose":         string(source.purpose),
				"is_sorcery":      isSorcerySkill(source.card.Card),
				"defense_use":     true,
				"boost_use":       source.purpose == skillPurposeDefenseBoost,
				"cancel_boost":    &cancelled,
				"overexert_units": overexertUnits,
			}
			continueAfterFieldEffects := func() {
				counters := e.eligibleCounterTraps(attackerID, TriggerOnSpellCast, source.card, data)
				if e.promptCounterTrapQueue(counters, TriggerOnSpellCast, source.card, data, func() {
					promptNext(index, true)
				}) {
					return
				}
				promptNext(index, true)
			}
			if e.triggerSpellUseFieldEffectsWithContinuation(defenderID, source.card, data, continueAfterFieldEffects) {
				return
			}
			counters := e.eligibleCounterTraps(attackerID, TriggerOnSpellCast, source.card, data)
			if e.promptCounterTrapQueue(counters, TriggerOnSpellCast, source.card, data, func() {
				promptNext(index, true)
			}) {
				return
			}
		}
		if continuing && afterDone != nil {
			afterDone()
		}
	}
	promptNext(0, false)
	return e.State.PendingAction != nil
}

func mergeSkillIDSet(a map[string]bool, b map[string]bool) map[string]bool {
	merged := make(map[string]bool, len(a)+len(b))
	for id := range a {
		merged[id] = true
	}
	for id := range b {
		merged[id] = true
	}
	return merged
}

func (e *Engine) collectDefenseScrollUses(ps *PlayerState, ids []string, reserved map[string]bool) ([]*CardInstance, map[string]int, error) {
	scrolls := make([]*CardInstance, 0, len(ids))
	totalCost := make(map[string]int)
	costPlan := newActionCostPlan(ps)
	seen := make(map[string]bool)
	for id := range reserved {
		seen[id] = true
	}
	for _, id := range ids {
		if seen[id] {
			return nil, nil, fmt.Errorf("defense source %s selected more than once", id)
		}
		seen[id] = true
		card, _ := ps.FindHandCard(id)
		if card == nil {
			return nil, nil, fmt.Errorf("defense scroll not found: %s", id)
		}
		if !isSpellScrollCard(card.Card) || !canUseSkillForPurpose(card.Card, skillPurposeDefend) {
			return nil, nil, fmt.Errorf("card %s is not a defense spell scroll", id)
		}
		if err := e.validateHandSpellScrollForPurpose(card, skillPurposeDefend); err != nil {
			return nil, nil, fmt.Errorf("defense scroll %s cannot be used for %s: %w", id, skillPurposeDefend, err)
		}
		scrolls = append(scrolls, card)
		for elem, amount := range costPlan.cardPlayCost(e, card) {
			totalCost[elem] += amount
		}
	}
	return scrolls, totalCost, nil
}

func (e *Engine) splitBoostIDs(ps *PlayerState, ids []string) ([]string, []string) {
	skillIDs := make([]string, 0, len(ids))
	scrollIDs := make([]string, 0)
	for _, id := range ids {
		if card, _ := ps.FindHandCard(id); card != nil && isSpellScrollCard(card.Card) {
			scrollIDs = append(scrollIDs, id)
			continue
		}
		skillIDs = append(skillIDs, id)
	}
	return skillIDs, scrollIDs
}

func (e *Engine) collectHandBoostScrollUses(ps *PlayerState, ids []string, reserved map[string]bool, purpose skillPurpose) ([]*CardInstance, map[string]int, error) {
	scrolls := make([]*CardInstance, 0, len(ids))
	totalCost := make(map[string]int)
	seen := make(map[string]bool)
	costPlan := newActionCostPlan(ps)
	for id := range reserved {
		seen[id] = true
	}
	for _, id := range ids {
		if seen[id] {
			return nil, nil, fmt.Errorf("boost source %s selected more than once", id)
		}
		seen[id] = true
		card, _ := ps.FindHandCard(id)
		if card == nil {
			return nil, nil, fmt.Errorf("boost scroll not found: %s", id)
		}
		if !isSpellScrollCard(card.Card) {
			return nil, nil, fmt.Errorf("card %s is not a spell scroll", id)
		}
		if err := e.validateHandSpellScrollForPurpose(card, purpose); err != nil {
			return nil, nil, fmt.Errorf("boost scroll %s cannot be used for %s: %w", id, purpose, err)
		}
		scrolls = append(scrolls, card)
		for elem, amount := range costPlan.cardPlayCost(e, card) {
			totalCost[elem] += amount
		}
	}
	return scrolls, totalCost, nil
}

func (e *Engine) validateHandSpellScrollForPurpose(card *CardInstance, purpose skillPurpose) error {
	if card == nil || card.Card == nil || !isSpellScrollCard(card.Card) {
		return fmt.Errorf("card is not a spell scroll")
	}
	if e.hasEffectiveStatus(card, StatusPetrify) {
		return fmt.Errorf("card is petrified")
	}
	if !canUseSkillForPurpose(card.Card, purpose) {
		return fmt.Errorf("spell scroll cannot be used for %s", purpose)
	}
	return e.validateSkillUsePermissionModifiers(card, purpose)
}

func (e *Engine) moveHandSpellScrollsToGraveyard(ps *PlayerState, cards []*CardInstance, use string) {
	for _, card := range cards {
		if card == nil {
			continue
		}
		_, idx := ps.FindHandCard(card.InstanceID)
		if idx < 0 {
			continue
		}
		e.notifyCardPlayCostPaid(ps, card)
		ps.RemoveFromHand(idx)
		e.addToGraveyard(ps.PlayerID, card)
		e.emit(GameEvent{
			Type:   "use_item",
			Player: -1,
			Data: map[string]any{
				"player":           ps.PlayerID,
				"card":             cardToInfo(card),
				"elements":         ps.Elements,
				"spell_scroll_use": use,
				"defense_scroll":   use == "defense",
				"attack_boost":     use == "attack_boost",
			},
		})
	}
}

func (e *Engine) finishDefenseResolution(playerID int, defenseSkills []*CardInstance, boostSkills []*CardInstance, overexerted int, powerBonusSource *CardInstance, powerBonus int) {
	defenseSkills = e.filterIceSoulSealCancelledBoosts(defenseSkills)
	boostSkills = e.filterIceSoulSealCancelledBoosts(boostSkills)
	totalDefPower := e.totalEffectiveSkillPower(playerID, defenseSkills, skillPurposeDefend) +
		e.totalEffectiveSkillPower(playerID, boostSkills, skillPurposeDefenseBoost)
	if powerBonusSource != nil && (cardInstanceInSlice(defenseSkills, powerBonusSource) || cardInstanceInSlice(boostSkills, powerBonusSource)) {
		totalDefPower += max(powerBonus, 0)
	}

	attackPower := e.State.PendingSpell.TotalPower

	e.emit(GameEvent{
		Type:   "defense_attempt",
		Player: -1,
		Data: map[string]any{
			"defender":      playerID,
			"defense_power": totalDefPower,
			"attack_power":  attackPower,
			"skills_used":   len(defenseSkills) + len(boostSkills),
			"overexerted":   overexerted,
		},
	})

	requiredPower := e.requiredDefensePowerForSpell(e.State.PendingSpell.Skill, attackPower)
	defenseSuccess := attackPower <= 0 || (totalDefPower >= requiredPower && len(defenseSkills) > 0)
	defendData := map[string]any{
		"defender":        playerID,
		"attacker":        e.State.PendingSpell.AttackerID,
		"defense_power":   totalDefPower,
		"attack_power":    attackPower,
		"required_power":  requiredPower,
		"defense_success": defenseSuccess,
		"attack_skill":    e.State.PendingSpell.Skill,
		"boost_skills":    e.State.PendingSpell.BoostSkills,
		"defense_skills":  defenseSkills,
		"defense_boosts":  boostSkills,
	}
	for _, defenseSkill := range defenseSkills {
		e.triggerEffects(TriggerOnDefend, defenseSkill, nil, defendData)
	}
	if e.State.PendingSpell != nil && e.State.PendingSpell.Skill != nil {
		e.triggerEffects(TriggerOnDefend, e.State.PendingSpell.Skill, nil, defendData)
	}
	e.triggerFieldEffectsWithData(TriggerOnDefend, playerID, e.State.PendingSpell.Skill, defendData)
	e.triggerFieldEffectsWithData(TriggerOnDefend, e.State.PendingSpell.AttackerID, e.State.PendingSpell.Skill, defendData)

	if defenseSuccess {
		// Defense successful
		e.emit(GameEvent{
			Type:   "defense_success",
			Player: -1,
			Data:   map[string]any{"defender": playerID},
		})
		e.consumeNextSpellAttackBonuses(e.State.Players[e.State.PendingSpell.AttackerID], e.State.PendingSpell.Skill)
		e.removeStoredArchmageStaffSkillAfterUse(e.State.PendingSpell.AttackerID, e.State.PendingSpell.Skill)
		clearFiveRainbowBeamSelection(e.State.PendingSpell.Skill)
	} else {
		// Defense failed, spell hits
		if e.resolveSpellHit(
			e.State.PendingSpell.AttackerID,
			e.State.PendingSpell.Skill,
			e.State.PendingSpell.Target,
			e.State.PendingSpell.BoostSkills,
			e.State.PendingSpell.ExtraTargets,
		) {
			return
		}
		e.removeStoredArchmageStaffSkillAfterUse(e.State.PendingSpell.AttackerID, e.State.PendingSpell.Skill)
		clearFiveRainbowBeamSelection(e.State.PendingSpell.Skill)
	}

	continued := e.completePendingSpell(e.State.PendingSpell)
	if !continued && e.State.PendingAction == nil {
		e.State.Phase = PhaseMain
	}
	e.checkWinCondition()
}

func (e *Engine) requiredDefensePowerForSpell(skill *CardInstance, attackPower int) int {
	required := attackPower
	if b, ok := cardBehavior(skill).(DefenseRequirementBehavior); ok && b.HasActiveDefenseRequirement(skill) {
		required = b.RequiredDefensePower(skill, attackPower)
	}
	return required
}

func (e *Engine) filterIceSoulSealCancelledBoosts(skills []*CardInstance) []*CardInstance {
	if len(skills) == 0 {
		return skills
	}
	filtered := skills[:0]
	for _, skill := range skills {
		if skill == nil || skill.Statuses[iceSoulSealCancelledBoostStatus] <= 0 {
			filtered = append(filtered, skill)
			continue
		}
		delete(skill.Statuses, iceSoulSealCancelledBoostStatus)
	}
	return filtered
}

func (e *Engine) promptDispelDefenseSpellIfEligible(attackerID int, defenderID int, defenseSkills []*CardInstance, boostSkills []*CardInstance, overexerted int, powerBonusSource *CardInstance, powerBonus int) bool {
	defenseOnlySkills := make([]*CardInstance, 0, len(defenseSkills))
	for _, skill := range defenseSkills {
		if skill != nil && skill.Card != nil && isDefenseOnlySkill(skill.Card) {
			defenseOnlySkills = append(defenseOnlySkills, skill)
		}
	}
	if len(defenseOnlySkills) == 0 {
		return false
	}

	dispel := e.findReadyDispelSkill(attackerID)
	if dispel == nil {
		return false
	}

	candidates := make([]map[string]any, 0, len(defenseOnlySkills))
	validTargets := make(map[string]*CardInstance, len(defenseOnlySkills))
	for _, skill := range defenseOnlySkills {
		candidate := candidateInfo(skill, "skill", "enemy")
		candidates = append(candidates, candidate)
		validTargets[skill.InstanceID] = skill
	}

	cost := e.effectiveSkillUseCost(e.State.Players[attackerID], dispel)
	e.SetPendingActionWithError(attackerID, "dispel_defense_spell",
		"解咒:选择1个防御法术无效",
		candidates, 0, 1, cost, true,
		func(selected []string, data map[string]any) error {
			if len(selected) == 0 {
				e.finishDefenseResolution(defenderID, defenseSkills, boostSkills, overexerted, powerBonusSource, powerBonus)
				return nil
			}
			cancelledSkill := validTargets[selected[0]]
			if cancelledSkill == nil {
				return fmt.Errorf("invalid defense spell selection")
			}
			if err := e.payAndUseDispel(attackerID, dispel, cost, data); err != nil {
				return err
			}
			e.emit(GameEvent{
				Type:   "spell_reaction",
				Player: -1,
				Data: map[string]any{
					"player":    attackerID,
					"card":      cardToInfo(dispel),
					"effect":    "cancel_defense_spell",
					"cancelled": cardToInfo(cancelledSkill),
				},
			})
			e.finishDefenseResolution(defenderID, withoutCardInstance(defenseSkills, cancelledSkill), boostSkills, overexerted, powerBonusSource, powerBonus)
			return nil
		})
	return e.State.PendingAction != nil && e.State.PendingAction.Type == "dispel_defense_spell"
}

func (e *Engine) findReadyDispelSkill(playerID int) *CardInstance {
	ps := e.State.Players[playerID]
	for _, skill := range ps.Skills {
		if skill == nil || skill.Card == nil {
			continue
		}
		if b, ok := cardBehavior(skill).(DefenseSpellCancellationBehavior); ok && b.CanCancelDefenseSpell(skill) && e.validateReadySkill(skill) == nil {
			return skill
		}
	}
	return nil
}

func (e *Engine) payAndUseDispel(playerID int, dispel *CardInstance, cost map[string]int, data map[string]any) error {
	if dispel == nil {
		return fmt.Errorf("dispel not found")
	}
	if err := e.validateReadySkill(dispel); err != nil {
		return err
	}
	overexertIDs := stringsFromAnySlice(anySliceFromData(data, "overexert_ids"))
	overexertUnits, err := e.collectOverexertUnits(e.State.Players[playerID], overexertIDs)
	if err != nil {
		return err
	}
	if !e.canPayCostWithOverexertOptions(e.State.Players[playerID], cost, overexertUnits, e.playerHasLightWildcard(e.State.Players[playerID])) {
		return fmt.Errorf("not enough elements")
	}
	if !e.payDefenseCostWithOptions(e.State.Players[playerID], cost, ActionMessage{Data: data}, overexertUnits, e.playerHasLightWildcard(e.State.Players[playerID])) {
		return fmt.Errorf("invalid payment")
	}
	e.destroyFuyeDoomedAfterExert(overexertUnits)
	dispel.IsHorizontal = true
	if !e.shouldSkipCooldown(e.State.Players[playerID], dispel) {
		e.ApplyKeywordOnSkillUse(dispel)
	}
	e.applySkillUseCooldownModifiers(e.State.Players[playerID], dispel)
	e.consumeNextSkillUseModifiers(e.State.Players[playerID], dispel)
	e.advanceMasteryForUsedSkills(playerID, dispel)
	return nil
}

func withoutCardInstance(cards []*CardInstance, removed *CardInstance) []*CardInstance {
	result := make([]*CardInstance, 0, len(cards))
	for _, card := range cards {
		if card != nil && removed != nil && card.InstanceID == removed.InstanceID {
			continue
		}
		result = append(result, card)
	}
	return result
}

func (e *Engine) collectOverexertUnits(ps *PlayerState, ids []string) ([]*CardInstance, error) {
	cards := make([]*CardInstance, 0, len(ids))
	seen := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			return nil, fmt.Errorf("card %s selected more than once", id)
		}
		seen[id] = true
		card := e.findUnitOnGrid(ps, id)
		if card == nil {
			for _, equipment := range ps.Equipment {
				if equipment != nil && equipment.InstanceID == id {
					card = equipment
					break
				}
			}
		}
		if card == nil || card.Card == nil {
			return nil, fmt.Errorf("overexert card not found: %s", id)
		}
		if !e.canConsumeCard(card) {
			return nil, fmt.Errorf("card cannot be overexerted: %s", id)
		}
		cards = append(cards, card)
	}
	return cards, nil
}
