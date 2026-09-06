package game

import "fmt"

// dealDamage deals damage to a card instance
func (e *Engine) dealDamage(target *CardInstance, amount int, ownerID int) {
	e.ApplyDamage(DamageRequest{Target: target, Amount: amount})
}

func (e *Engine) fieldDamagePreventionSource(target *CardInstance, ownerID int, damageData map[string]any) *CardInstance {
	if target == nil || ownerID < 0 || ownerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[ownerID]
	for _, source := range e.getAllFieldCards(ps) {
		if source == nil || source.Card == nil || source == target {
			continue
		}
		behavior, ok := globalRegistry.GetBehavior(source.Card.Number).(FieldDamagePreventionBehavior)
		if !ok || !behavior.HasActiveFieldDamagePrevention(source) {
			continue
		}
		ctx := &EffectContext{
			Engine:     e,
			Source:     source,
			Target:     target,
			PlayerID:   ownerID,
			OpponentID: 1 - ownerID,
			ExtraData:  damageData,
		}
		if behavior.PreventsFieldDamage(ctx) {
			return source
		}
	}
	return nil
}

func (e *Engine) resolveDamage(target *CardInstance, amount int, ownerID int, extraData map[string]any) {
	if target == nil || target.Card == nil || amount <= 0 {
		return
	}
	damageData := map[string]any{
		"damage": amount,
	}
	for key, value := range extraData {
		damageData[key] = value
	}
	if e.temporaryDamageAndNegativeImmunityActive(target) {
		e.emit(GameEvent{
			Type:   "damage_prevented",
			Player: -1,
			Data: map[string]any{
				"target": cardToInfo(target),
				"amount": amount,
				"reason": "temporary_immunity",
			},
		})
		return
	}
	if behavior, ok := globalRegistry.GetBehavior(target.Card.Number).(DamagePreventionBehavior); ok && behavior.HasActiveDamagePrevention(target) {
		ctx := &EffectContext{
			Engine:     e,
			Source:     target,
			Target:     target,
			PlayerID:   ownerID,
			OpponentID: 1 - ownerID,
			ExtraData:  damageData,
		}
		if behavior.PreventsDamage(ctx) {
			e.emit(GameEvent{
				Type:   "damage_prevented",
				Player: -1,
				Data: map[string]any{
					"target": cardToInfo(target),
					"amount": amount,
				},
			})
			return
		}
	}
	if source := e.fieldDamagePreventionSource(target, ownerID, damageData); source != nil {
		e.emit(GameEvent{
			Type:   "damage_prevented",
			Player: -1,
			Data: map[string]any{
				"source": cardToInfo(source),
				"target": cardToInfo(target),
				"amount": amount,
				"reason": "field_prevention",
			},
		})
		return
	}

	amount = e.adjustedDamage(damageEventFromContext(&EffectContext{Source: target, Target: target, ExtraData: damageData}))
	damageData["damage"] = amount
	if amount <= 0 {
		return
	}

	if target.Statuses[sturdyScrollShieldStatus] > 0 && target.Statuses[sturdyScrollShieldUntilStatus] >= e.State.TurnNumber {
		prevented := min(amount, target.Statuses[sturdyScrollShieldStatus])
		target.Statuses[sturdyScrollShieldStatus] -= prevented
		if target.Statuses[sturdyScrollShieldStatus] <= 0 {
			delete(target.Statuses, sturdyScrollShieldStatus)
			delete(target.Statuses, sturdyScrollShieldUntilStatus)
		}
		amount -= prevented
		e.emit(GameEvent{
			Type:   "damage_prevented",
			Player: -1,
			Data: map[string]any{
				"target": cardToInfo(target),
				"amount": prevented,
				"reason": "sturdy_scroll",
			},
		})
		if amount <= 0 {
			return
		}
	}

	amount = e.applyPlayerShieldDamage(target, amount, damageData)
	if amount <= 0 {
		return
	}
	modifierPlan := e.planDamageModifiers(target, amount, ownerID, damageData)
	amount = modifierPlan.Amount
	modifierPlan.commit()
	if amount <= 0 {
		return
	}
	if target.Statuses["防止致命"] > 0 && target.CurrentLife-amount <= 0 {
		target.Statuses["防止致命"]--
		if target.Statuses["防止致命"] <= 0 {
			delete(target.Statuses, "防止致命")
		}
		e.emit(GameEvent{
			Type:   "damage_prevented",
			Player: -1,
			Data: map[string]any{
				"target": cardToInfo(target),
				"amount": damageData["damage"],
				"reason": "prevent_lethal",
			},
		})
		return
	}

	damageData["damage"] = amount
	commit := func() { e.commitDamage(target, amount, ownerID, damageData, extraData) }
	continueAfterCounters := func() {
		if !e.promptLethalSacrificePrevention(target, amount, ownerID, damageData, commit) {
			commit()
		}
	}
	if !e.promptDamagePreventionCounters(target, amount, ownerID, damageData, continueAfterCounters) {
		continueAfterCounters()
	}
}

func (e *Engine) commitDamage(target *CardInstance, amount int, ownerID int, damageData, extraData map[string]any) {
	target.CurrentLife -= amount
	target.DamageTakenThisTurn += amount
	if target.Card != nil && target.Card.IsHero() && ownerID >= 0 && ownerID < len(e.State.Players) {
		e.State.Players[ownerID].HeroDamageTakenThisTurn += amount
	}
	if actualDamage, ok := damageData["actual_damage_by_instance"].(map[string]int); ok && target.InstanceID != "" {
		actualDamage[target.InstanceID] += amount
	}
	if actualFriendlyDamage, ok := damageData["actual_friendly_damage_by_instance"].(map[string]int); ok && target.InstanceID != "" {
		if attacker, ok := damageData["attacker"].(int); ok && attacker == ownerID {
			actualFriendlyDamage[target.InstanceID] += amount
		}
	}
	if target.Card != nil && target.Card.IsCompanion() && ownerID >= 0 && ownerID < len(e.State.Players) {
		e.State.Players[ownerID].FriendlyUnitDamagedThisTurn = true
	}

	e.emit(GameEvent{
		Type:   "damage",
		Player: -1,
		Data: map[string]any{
			"target":    cardToInfo(target),
			"amount":    amount,
			"remaining": target.CurrentLife,
		},
	})

	damageData["damage"] = amount
	damageData["damage_taken_this_turn"] = target.DamageTakenThisTurn

	// Trigger 受伤 effects
	e.triggerEffects(TriggerOnDamaged, target, nil, damageData)
	fieldDamageData := map[string]any{
		"damaged_player": ownerID,
		"damage":         amount,
	}
	for key, value := range extraData {
		fieldDamageData[key] = value
	}
	e.triggerFieldEffectsWithData(TriggerOnDamaged, ownerID, target, fieldDamageData)
	enemyDamageData := map[string]any{"damaged_player": ownerID, "damage": amount}
	for key, value := range extraData {
		enemyDamageData[key] = value
	}
	e.triggerFieldEffectsWithData(TriggerOnDamaged, 1-ownerID, target, enemyDamageData)
	e.triggerHiddenFriendlyDamaged(ownerID, target, fieldDamageData)
	e.promptPainScreamWeakenAfterFriendlyDamage(ownerID, target, amount)

	if target.CurrentLife <= 0 {
		if attacker, ok := damageData["attacker"].(int); ok {
			target.Statuses["lethal_source_player"] = attacker + 1
		} else {
			delete(target.Statuses, "lethal_source_player")
		}
		e.queueDeathWithData(target, ownerID, extraData)
		if e.resolutionDepth == 0 && !e.resolvingDeaths {
			e.resolvePendingDeaths()
		}
	}
}

func (e *Engine) promptLethalSacrificePrevention(target *CardInstance, amount int, ownerID int, damageData map[string]any, afterDecline func()) bool {
	if e == nil || target == nil || amount <= 0 || target.CurrentLife-amount > 0 {
		return false
	}
	ps := e.State.Players[ownerID]
	if ps == nil {
		return false
	}
	candidates := make([]map[string]any, 0)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit == nil || unit == target || !e.canPreventLethalBySacrifice(unit, target, amount, damageData) {
				continue
			}
			candidates = append(candidates, candidateInfo(unit, "unit", "own"))
		}
	}
	if len(candidates) == 0 {
		return false
	}
	prompt := fmt.Sprintf("是否献祭1个伙伴，防止%s受到的%d点伤害", target.Card.Name, amount)
	e.SetPendingAction(ownerID, "prevent_lethal_sacrifice", prompt, candidates, 0, 1, func(selected []string) {
		dolphin := selectedUnitFromCandidates(e, selected, candidates)
		if dolphin == nil || dolphin.OwnerID != ownerID || dolphin == target || !e.canPreventLethalBySacrifice(dolphin, target, amount, damageData) {
			afterDecline()
			return
		}
		e.destroyUnitWithCause(dolphin, ownerID, DeathCauseSacrifice)
		e.emit(GameEvent{
			Type:   "damage_prevented",
			Player: -1,
			Data: map[string]any{
				"source": cardToInfo(dolphin),
				"target": cardToInfo(target),
				"amount": amount,
				"reason": "sacrifice_prevention",
			},
		})
	})
	return true
}

func (e *Engine) promptDamagePreventionCounters(target *CardInstance, amount int, ownerID int, damageData map[string]any, afterDecline func()) bool {
	if e == nil || target == nil || amount <= 0 || target.CurrentLife-amount > 0 {
		return false
	}
	if ownerID < 0 || ownerID >= len(e.State.Players) {
		return false
	}
	prevented := false
	counterData := map[string]any{
		"damaged_player": ownerID,
		"damage":         amount,
		"prevent_damage": &prevented,
	}
	for key, value := range damageData {
		counterData[key] = value
	}
	if e.promptCounterTrapQueue(e.eligibleCounterTraps(ownerID, TriggerBeforeDamage, target, counterData), TriggerBeforeDamage, target, counterData, func() {
		if prevented {
			return
		}
		afterDecline()
	}) {
		return true
	}
	return false
}

func (e *Engine) triggerHiddenFriendlyDamaged(playerID int, target *CardInstance, extraData map[string]any) {
	ps := e.State.Players[playerID]
	hidden := append([]*CardInstance{}, ps.Hand...)
	hidden = append(hidden, ps.Deck...)
	seenGroups := map[string]bool{}
	for _, card := range hidden {
		if card == nil || card.Card == nil {
			continue
		}
		behavior, ok := globalRegistry.GetBehavior(card.Card.Number).(OnFriendlyDamagedFromHiddenBehavior)
		if !ok || !behavior.HasActiveFriendlyDamagedFromHidden(card) {
			continue
		}
		ctx := &EffectContext{
			Engine:     e,
			Source:     card,
			Target:     target,
			PlayerID:   playerID,
			OpponentID: 1 - playerID,
			ExtraData:  extraData,
		}
		if group, ok := behavior.(HiddenDamageGroupBehavior); ok {
			key := group.HiddenDamageGroup(card)
			if key != "" {
				if seenGroups[key] {
					continue
				}
				seenGroups[key] = true
			}
		}
		_ = behavior.OnFriendlyDamagedFromHidden(ctx, damageEventFromContext(ctx))
	}
}
