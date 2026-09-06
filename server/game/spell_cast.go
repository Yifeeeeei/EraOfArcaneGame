package game

import (
	"fmt"
)

// handleCastSpell handles casting a spell
func (e *Engine) handleCastSpell(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}
	instanceID, _ := action.Data["instance_id"].(string)
	targetType, _ := action.Data["target_type"].(string)
	targetPos := Position{}
	var err error
	if targetType == "unit" {
		targetPos, err = requiredBoardPosition(action.Data, "target_col", "target_row")
		if err != nil {
			return err
		}
	}
	extraTargetPos, hasExtraTarget, err := optionalBoardPosition(action.Data, "extra_target_col", "extra_target_row")
	if err != nil {
		return err
	}
	boostIDsRaw, _ := action.Data["boost_ids"].([]any)

	ps := e.State.Players[playerID]

	skill := e.findSkill(ps, instanceID)
	if skill == nil {
		return fmt.Errorf("skill not found in skill area or bound skills")
	}

	if err := e.validateSkillForPurpose(skill, skillPurposeAttack); err != nil {
		return err
	}

	target := SpellTarget{
		Type:     targetType,
		Position: targetPos,
	}
	if ownerF, ok := action.Data["target_owner"].(float64); ok {
		ownerID := int(ownerF)
		target.OwnerID = &ownerID
	}
	targetUnit := e.spellTargetUnitForCaster(playerID, target)

	options, err := e.prepareSpellOptions(playerID, skill, target, action)
	if err != nil {
		return err
	}
	if len(boostIDsRaw) > 0 && !canSkillBeBoosted(skill) {
		return fmt.Errorf("skill cannot be boosted")
	}
	// Process boost skills (法术强化)
	boostIDs := stringsFromAnySlice(boostIDsRaw)
	boostSkillIDs, boostScrollIDs := e.splitBoostIDs(ps, boostIDs)
	boostSkills, _, err := e.collectSkillUses(ps, boostSkillIDs, skillPurposeAttackBoost, map[string]bool{instanceID: true})
	if err != nil {
		return err
	}
	usedBoostIDs := mergeSkillIDSet(map[string]bool{instanceID: true}, skillIDSet(boostSkills))
	boostScrolls, _, err := e.collectHandBoostScrollUses(ps, boostScrollIDs, usedBoostIDs, skillPurposeAttackBoost)
	if err != nil {
		return err
	}
	boostSources := append(append([]*CardInstance{}, boostSkills...), boostScrolls...)
	hasPierce := e.spellHasPierceWithBoosts(playerID, skill, boostSources) || options.Pierce
	if err := e.validateSpellTargetWithPierce(playerID, skill, target, hasPierce); err != nil {
		return err
	}
	extraTargets := make([]SpellTarget, 0, 1)
	if (options.AllowExtraTarget || e.hasNextDriveSpellExtraTarget(ps, skill)) && hasExtraTarget {
		extra := SpellTarget{Type: "unit", Position: extraTargetPos}
		if err := e.validateSpellExtraTargetForSkill(playerID, skill, target, extra); err != nil {
			return err
		}
		if extra.Position != target.Position || e.allowsSameSpellExtraTarget(ps, skill) {
			extraTargets = append(extraTargets, extra)
		}
	}
	if options.ExtraTargets != nil {
		additional, err := options.ExtraTargets(hasPierce)
		if err != nil {
			return err
		}
		extraTargets = append(extraTargets, additional...)
	}
	consumeNextExtraTargetModifier := e.hasNextSpellExtraTarget(ps, skill)
	extraTargets = e.expandSpellTargets(playerID, target, extraTargets)
	costPlan := newActionCostPlan(ps)
	cost := costPlan.skillUseCost(e, skill, skillPurposeAttack, map[string]any{
		"spell_target":      target,
		"spell_target_unit": targetUnit,
	})
	if options.ModifyCost != nil {
		options.ModifyCost(cost)
	}
	if err := e.validateOwnCost(playerID, skill, cost, action); err != nil {
		return err
	}
	totalCost := mergeElementCosts(cost)
	for _, boostSkill := range boostSkills {
		totalCost = mergeElementCosts(totalCost, costPlan.skillUseCost(e, boostSkill, skillPurposeAttackBoost, nil))
	}
	for _, boostScroll := range boostScrolls {
		totalCost = mergeElementCosts(totalCost, costPlan.cardPlayCost(e, boostScroll))
	}
	if !e.canPayCost(ps, totalCost) {
		return fmt.Errorf("not enough elements for boost skills")
	}
	powerSacrifice, powerSacrificeSource, powerSacrificeBonus, err := e.validateSpellPowerSacrificeForSources(playerID, append([]*CardInstance{skill}, boostSources...), action)
	if err != nil {
		return err
	}

	// Pay costs and set cards horizontal only after all validation succeeds.
	if !e.payCostForCardAction(ps, skill, cost, totalCost, paymentPurposeUse, action) {
		return fmt.Errorf("invalid payment")
	}
	if options.Commit != nil {
		options.Commit()
	}
	skill.IsHorizontal = true
	tapSkills(boostSkills)

	// Apply cooldown from keyword
	if !e.shouldSkipCooldown(ps, skill) {
		e.ApplyKeywordOnSkillUse(skill)
	}
	e.applySkillUseCooldownModifiers(ps, append([]*CardInstance{skill}, boostSkills...)...)
	e.notifySkillUseCommitted(playerID, skill)
	e.consumeNextSkillUseModifiers(ps, skill)
	for _, boostSkill := range boostSkills {
		e.consumeNextSkillUseModifiersForPurpose(ps, boostSkill, skillPurposeAttackBoost)
	}
	e.moveHandSpellScrollsToGraveyard(ps, boostScrolls, "attack_boost")
	e.advanceMasteryForUsedSkills(playerID, append([]*CardInstance{skill}, boostSkills...)...)

	if options.ResolveInstead != nil {
		e.recordSpellCast(playerID, skill)
		e.emit(GameEvent{Type: "spell_cast", Player: -1, Data: map[string]any{
			"cast_player": playerID,
			"attacker":    playerID,
			"skill":       cardToInfo(skill),
			"target":      target,
			"power":       0,
			"boost_count": len(boostSkills),
			"is_sorcery":  isSorcerySkill(skill.Card),
		}})
		options.ResolveInstead()
		return nil
	}

	if powerSacrifice != nil && powerSacrificeBonus > 0 {
		e.destroyUnitWithCause(powerSacrifice, playerID, DeathCauseSacrifice)
	}
	e.applyCoralBellyFirstSpellAttackBonus(playerID, skill)
	powerTargets := append([]SpellTarget{target}, extraTargets...)
	totalPower := e.effectiveSpellPower(playerID, skill, boostSources, powerTargets...)
	if powerSacrificeSource != nil && powerSacrificeBonus > 0 {
		totalPower += powerSacrificeBonus
	}
	powerSources := e.spellPowerSources(playerID, skill, boostSources, totalPower, powerTargets...)
	e.consumeNextSpellPowerBonuses(ps, skill)
	if len(extraTargets) > 0 {
		e.consumeNextDriveSpellExtraTarget(ps, skill)
	} else if consumeNextExtraTargetModifier {
		e.consumeNextSpellExtraTarget(ps, skill)
	}

	// Check if it's a 咒术 (sorcery - unblockable)
	isSorcery := isSorcerySkill(skill.Card)
	e.recordSpellCast(playerID, skill)
	e.triggerMagicMothAfterFocusSpellCast(playerID, skill)
	spellCastData := map[string]any{
		"cast_player": playerID,
		"attacker":    playerID,
		"purpose":     string(skillPurposeAttack),
		"skill":       cardToInfo(skill),
		"target":      target,
		"power":       totalPower,
		"boost_count": len(boostSources),
		"is_sorcery":  isSorcery,
	}
	e.emit(GameEvent{
		Type:   "spell_cast",
		Player: -1,
		Data:   spellCastData,
	})
	e.triggerEffects(TriggerOnSpellCast, skill, nil, spellCastData)
	e.dispatchPreparedSpell(&SpellCast{AttackerID: playerID, Skill: skill, Target: target,
		TotalPower: totalPower, PowerSources: powerSources, BoostSkills: boostSources, ExtraTargets: extraTargets}, spellCastData, false)
	return nil
}

func (e *Engine) startFreeSpellCastNoBoost(playerID int, skill *CardInstance, target SpellTarget, extraData map[string]any) error {
	if playerID < 0 || playerID >= len(e.State.Players) || skill == nil || skill.Card == nil {
		return fmt.Errorf("invalid free spell cast")
	}
	ps := e.State.Players[playerID]
	if e.findSkill(ps, skill.InstanceID) != skill {
		return fmt.Errorf("skill not found in skill area or bound skills")
	}
	if err := e.validateSkillForPurpose(skill, skillPurposeAttack); err != nil {
		return err
	}
	if err := e.validateSpellTargetWithPierce(playerID, skill, target, e.skillHasPierce(playerID, skill)); err != nil {
		return err
	}

	skill.IsHorizontal = true
	if !e.shouldSkipCooldown(ps, skill) {
		e.ApplyKeywordOnSkillUse(skill)
	}
	e.applySkillUseCooldownModifiers(ps, skill)
	e.notifySkillUseCommitted(playerID, skill)
	e.advanceMasteryForUsedSkills(playerID, skill)

	e.applyCoralBellyFirstSpellAttackBonus(playerID, skill)
	totalPower := e.effectiveSpellPower(playerID, skill, nil, target)
	powerSources := e.spellPowerSources(playerID, skill, nil, totalPower, target)
	e.consumeNextSpellPowerBonuses(ps, skill)

	isSorcery := isSorcerySkill(skill.Card)
	e.recordSpellCast(playerID, skill)
	e.triggerMagicMothAfterFocusSpellCast(playerID, skill)
	spellCastData := map[string]any{
		"cast_player": playerID,
		"attacker":    playerID,
		"purpose":     string(skillPurposeAttack),
		"skill":       cardToInfo(skill),
		"target":      target,
		"power":       totalPower,
		"boost_count": 0,
		"is_sorcery":  isSorcery,
		"free_cast":   true,
	}
	for key, value := range extraData {
		spellCastData[key] = value
	}
	e.emit(GameEvent{Type: "spell_cast", Player: -1, Data: spellCastData})
	e.triggerEffects(TriggerOnSpellCast, skill, nil, spellCastData)
	e.dispatchPreparedSpell(&SpellCast{AttackerID: playerID, Skill: skill, Target: target,
		TotalPower: totalPower, PowerSources: powerSources, BoostSkills: nil, ExtraTargets: nil}, spellCastData, false)
	return nil
}

func (e *Engine) promptAttackBoostSpellCastCounters(attackerID int, boostSkills []*CardInstance, afterDone func()) bool {
	if len(boostSkills) == 0 || e.State.PendingSpell == nil {
		return false
	}
	defenderID := 1 - attackerID
	var promptNext func(int, bool)
	promptNext = func(index int, continuing bool) {
		for index < len(boostSkills) {
			boost := boostSkills[index]
			index++
			if boost == nil || boost.Card == nil || boost.Statuses[iceSoulSealCancelledBoostStatus] > 0 {
				continue
			}
			cancelled := false
			data := map[string]any{
				"cast_player":  attackerID,
				"attacker":     attackerID,
				"skill":        cardToInfo(boost),
				"power":        e.effectiveSkillPowerForPurpose(attackerID, boost, skillPurposeAttackBoost),
				"purpose":      string(skillPurposeAttackBoost),
				"is_sorcery":   isSorcerySkill(boost.Card),
				"boost_use":    true,
				"attack_boost": true,
				"cancel_boost": &cancelled,
			}
			continueAfterFieldEffects := func() {
				counters := e.eligibleCounterTraps(defenderID, TriggerOnSpellCast, boost, data)
				if e.promptCounterTrapQueue(counters, TriggerOnSpellCast, boost, data, func() {
					promptNext(index, true)
				}) {
					return
				}
				promptNext(index, true)
			}
			if e.triggerSpellUseFieldEffectsWithContinuation(attackerID, boost, data, continueAfterFieldEffects) {
				return
			}
			counters := e.eligibleCounterTraps(defenderID, TriggerOnSpellCast, boost, data)
			if e.promptCounterTrapQueue(counters, TriggerOnSpellCast, boost, data, func() {
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

func (e *Engine) shouldResolveSorceryHit(skill *CardInstance) bool {
	if skill == nil || skill.Card == nil {
		return false
	}
	if b, ok := cardBehavior(skill).(SorceryHitPolicyBehavior); ok && b.HasActiveSorceryHitPolicy(skill) {
		return b.ResolvesSorceryHit(skill)
	}
	return true
}
