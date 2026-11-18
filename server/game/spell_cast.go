package game

import (
	"eraofarcane/model"
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
	uses := []spellCopyUse{{instanceID, skillPurposeAttack}}
	for _, id := range stringsFromAnySlice(boostIDsRaw) {
		uses = append(uses, spellCopyUse{id, skillPurposeAttackBoost})
	}
	forms, handled, err := e.prepareSpellCopies(playerID, action, uses, func(a ActionMessage) error { return e.handleCastSpell(playerID, a) })
	if handled || err != nil {
		return err
	}

	skill := e.findSkill(ps, instanceID)
	var handSource *CardInstance
	if skill == nil {
		handSource, _ = ps.FindHandCard(instanceID)
		if handSource != nil {
			if _, ok := cardBehavior(handSource).(SpellCopyBehavior); ok {
				skill = handSource
			}
		}
	}
	if skill == nil {
		return fmt.Errorf("skill not found in skill area or bound skills")
	}
	skill = spellForm(skill, forms)
	if handled, err := e.chooseCopiedAttackTarget(playerID, skill, action, func(a ActionMessage) error { return e.handleCastSpell(playerID, a) }); handled || err != nil {
		return err
	}
	if handSource != nil && e.chooseCopiedItemBoosts(playerID, skill, action) {
		return nil
	}

	if e.chooseCastSpellArea(playerID, skill, action) {
		return nil
	}
	if err := e.validateCopiedOrLearnedSpell(skill, handSource, skillPurposeAttack); err != nil {
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
	target, err = e.spellAreaSelection(skill, target, action)
	if err != nil {
		return err
	}
	if target.Area == SpellAreaSquare && target.AreaOrigin == nil {
		handled, err := e.chooseEffectSpellSquare(playerID, skill, target, func(chosen SpellTarget) error {
			data := make(map[string]any, len(action.Data)+2)
			for k, v := range action.Data {
				data[k] = v
			}
			data["area_col"] = float64(chosen.AreaOrigin.Col)
			data["area_row"] = float64(chosen.AreaOrigin.Row)
			return e.handleCastSpell(playerID, ActionMessage{Data: data})
		})
		if handled {
			return err
		}
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
	boostSkills, _, err := e.collectSkillUses(ps, boostSkillIDs, skillPurposeAttackBoost, map[string]bool{instanceID: true}, forms)
	if err != nil {
		return err
	}
	usedBoostIDs := mergeSkillIDSet(map[string]bool{instanceID: true}, skillIDSet(boostSkills))
	boostScrolls, _, err := e.collectHandBoostScrollUses(ps, boostScrollIDs, usedBoostIDs, skillPurposeAttackBoost, forms)
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
		if err := e.validateSpellExtraTargetForSkill(playerID, skill, target, extra, hasPierce); err != nil {
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
	extraTargets = e.expandSpellTargets(playerID, skill, target, extraTargets)
	costPlan := newActionCostPlan(ps)
	cost := costPlan.skillUseCost(e, skill, skillPurposeAttack, map[string]any{
		"spell_target":      target,
		"spell_target_unit": targetUnit,
	})
	if handSource != nil {
		cost = costPlan.cardPlayCost(e, handSource)
	}
	if options.ModifyCost != nil {
		options.ModifyCost(cost)
	}
	if err := e.validateOwnCost(playerID, skill, cost, action); err != nil {
		return err
	}
	totalCost := mergeElementCosts(cost)
	paymentParts := []spellPaymentPart{{skill, cost}}
	for _, boostSkill := range boostSkills {
		partCost := costPlan.skillUseCost(e, boostSkill, skillPurposeAttackBoost, nil)
		paymentParts = append(paymentParts, spellPaymentPart{boostSkill, partCost})
		totalCost = mergeElementCosts(totalCost, partCost)
	}
	for _, boostScroll := range boostScrolls {
		partCost := costPlan.cardPlayCost(e, boostScroll)
		paymentParts = append(paymentParts, spellPaymentPart{boostScroll, partCost})
		totalCost = mergeElementCosts(totalCost, partCost)
	}
	if !e.canPayCost(ps, totalCost) {
		return fmt.Errorf("not enough elements for boost skills")
	}
	powerSacrifice, powerSacrificeSource, powerSacrificeBonus, err := e.validateSpellPowerSacrificeForSources(playerID, append([]*CardInstance{skill}, boostSources...), action)
	if err != nil {
		return err
	}
	tapCosts, handled, err := e.prepareSpellTapCosts(playerID, append([]*CardInstance{skill}, boostSources...), action, func(a ActionMessage) error { return e.handleCastSpell(playerID, a) })
	if handled || err != nil {
		return err
	}

	// Pay costs and set cards horizontal only after all validation succeeds.
	if handSource != nil && len(boostSources) > 0 && !boolFromData(action.Data, "copy_payment_done") {
		e.SetPendingActionWithError(playerID, "copied_spell_payment", "确认复制法术及强化的支付", []map[string]any{{"instance_id": "pay", "name": "支付并施放", "zone": "choice"}}, 1, 1, totalCost, false, func(_ []string, choice map[string]any) error {
			data := copyActionData(action)
			data["copy_payment_done"] = true
			data["payment"] = choice["payment"]
			return e.handleCastSpell(playerID, ActionMessage{Data: data})
		})
		return nil
	}
	actualPayment, _, _ := e.cardPaymentPlanForAction(ps, skill, cost, totalCost, paymentPurposeUse, action)
	if actualPayment == nil && hasCopiedSpellUse(skill, boostSources) {
		available := copyElements(ps.Elements)
		available[model.ElementArcane] += ps.StrictArcane
		e.retryCopiedSpellPayment(playerID, action, totalCost, available, func(a ActionMessage) error { return e.handleCastSpell(playerID, a) })
		return nil
	}
	if actualPayment != nil {
		if handled, err := e.prepareUsePaymentAttributions(ps, paymentParts, actualPayment, action, func(a ActionMessage) error { return e.handleCastSpell(playerID, a) }); handled || err != nil {
			if err != nil && hasCopiedSpellUse(skill, boostSources) {
				available := copyElements(ps.Elements)
				available[model.ElementArcane] += ps.StrictArcane
				e.retryCopiedSpellPayment(playerID, action, totalCost, available, func(a ActionMessage) error { return e.handleCastSpell(playerID, a) })
				return nil
			}
			return err
		}
	}
	assignedPayments, _ := e.constrainedSpellPayments(ps, paymentParts, actualPayment, paymentAttributions(action))
	if !e.payCostForCardAction(ps, skill, cost, totalCost, paymentPurposeUse, action) {
		return fmt.Errorf("invalid payment")
	}
	e.commitSpellTapCosts(tapCosts)
	e.recordSpellUsePayments(ps, paymentParts, actualPayment, assignedPayments)
	if options.Commit != nil {
		options.Commit()
	}
	physicalSpellSource(skill).IsHorizontal = true
	skill.IsHorizontal = true
	if handSource != nil {
		e.notifyCardPlayCostPaid(ps, handSource)
		_, idx := ps.FindHandCard(handSource.InstanceID)
		ps.RemoveFromHand(idx)
		e.addToGraveyard(playerID, handSource)
		e.emit(GameEvent{Type: "use_item", Player: -1, Data: map[string]any{"player": playerID, "card": cardToInfo(handSource), "elements": ps.Elements}})
	}
	tapSkills(boostSkills)

	// Apply cooldown from keyword
	if !e.shouldSkipCooldown(ps, skill) {
		e.ApplyKeywordOnSkillUse(physicalSpellSource(skill))
	}
	e.applySkillUseCooldownModifiers(ps, append([]*CardInstance{skill}, boostSkills...)...)
	e.notifySkillUseCommitted(playerID, skill)
	for _, boost := range boostSources {
		e.notifySkillUseCommitted(playerID, boost)
	}
	e.consumeNextSkillUseModifiers(ps, skill)
	for _, boostSkill := range boostSkills {
		e.consumeNextSkillUseModifiersForPurpose(ps, boostSkill, skillPurposeAttackBoost)
	}
	e.moveHandSpellScrollsToGraveyard(ps, boostScrolls, "attack_boost")
	e.advanceMasteryForUsedSkills(playerID, append([]*CardInstance{skill}, boostSkills...)...)

	if options.ResolveInstead != nil {
		e.recordSpellAttack(playerID, skill)
		e.emit(GameEvent{Type: "spell_cast", Player: -1, Data: map[string]any{
			"cast_player": playerID,
			"attacker":    playerID,
			"skill":       cardToInfo(skill),
			"target":      target,
			"power":       0,
			"boost_count": len(boostSkills),
			"is_sorcery":  isSorcerySkill(skill.Card),
		}})
		e.resolveSpellTapCosts(playerID, tapCosts, options.ResolveInstead)
		return nil
	}

	resolveCastAfterTap := func() {
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
		e.recordSpellAttack(playerID, skill)
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
	}
	continueCast := func() {
		if handSource != nil {
			cancelled := false
			useData := map[string]any{"used_player": playerID, "cancel_item": &cancelled, "spell_scroll": true}
			resolveItem := func() {
				if !cancelled {
					e.runResolution("copied item use", func() {
						e.triggerFieldEffectsWithData(TriggerOnUseItem, playerID, handSource, useData)
						e.triggerFieldEffectsWithData(TriggerOnUseItem, 1-playerID, handSource, useData)
					}, resolveCastAfterTap)
				} else {
					e.finishSkillUses(append([]*CardInstance{skill}, boostSources...)...)
				}
			}
			if e.promptOpponentCounterTrap(playerID, TriggerOnUseItem, handSource, useData, resolveItem) {
				return
			}
			resolveItem()
		} else {
			resolveCastAfterTap()
		}
	}
	e.resolveSpellTapCosts(playerID, tapCosts, continueCast)
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
	if handled, err := e.chooseEffectSpellSquare(playerID, skill, target, func(chosen SpellTarget) error {
		return e.startFreeSpellCastNoBoost(playerID, skill, chosen, extraData)
	}); handled {
		return err
	}

	skill.usePayment = nil
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
