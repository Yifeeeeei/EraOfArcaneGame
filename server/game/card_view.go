package game

import (
	"eraofarcane/cards"
	"eraofarcane/model"
	"sort"
)

func cardToInfo(ci *CardInstance) map[string]any {
	if ci == nil {
		return nil
	}
	info := map[string]any{
		"instance_id":               ci.InstanceID,
		"owner":                     ci.OwnerID,
		"number":                    ci.Card.Number,
		"name":                      ci.Card.Name,
		"type":                      ci.Card.Type,
		"category":                  ci.Card.Category,
		"tag":                       ci.Card.Tag,
		"description":               ci.Card.Description,
		"attack":                    ci.Card.Attack + ci.AttackBonus,
		"life":                      maxLife(ci),
		"power":                     ci.Card.Power + ci.PowerBonus,
		"duration":                  ci.Card.Duration,
		"elements_cost":             ci.Card.ElementsCost,
		"elements_gain":             effectiveElementsGain(ci),
		"elements_expense":          ci.Card.ElementsExpense,
		"current_life":              ci.CurrentLife,
		"current_attack":            effectiveCurrentAttack(ci),
		"is_horizontal":             ci.IsHorizontal,
		"is_terrain":                cards.IsTerrain(ci.Card.Number),
		"is_companion":              ci.Card.IsCompanion(),
		"is_consumable":             cards.IsConsumable(ci.Card.Number),
		"is_equipment":              cards.IsEquipment(ci.Card.Number),
		"is_weapon":                 cards.IsWeapon(ci.Card.Number),
		"can_attack_from_non_front": cardCanAttackFromNonFront(ci),
		"has_taunt":                 cardHasTaunt(ci),
		"has_global_spell_range":    cardHasActiveGlobalSpellRange(ci),
		"is_counter_trap":           isCounterTrapCard(ci.Card.Number),
		"is_set_counter":            ci.IsSetCounter,
		"statuses":                  ci.Statuses,
		"position":                  ci.Position,
		"output_path":               ci.Card.OutputPath,
		"used_this_turn":            ci.UsedThisTurn,
		"ultimate_used":             ci.UltimateUsed,
		"uses_remaining":            ci.UsesRemaining,
	}
	addCardEffectMetadata(info, ci.Card)
	if len(ci.BoundSkills) > 0 {
		info["bound_skills"] = cardsToInfo(ci.BoundSkills)
	}
	if len(ci.UnderCards) > 0 {
		info["under_cards"] = cardsToInfo(ci.UnderCards)
	}
	if attached := attachedBehaviorsInfo(ci); len(attached) > 0 {
		info["attached_behaviors"] = attached
	}

	hasPerTurn := cardHasActivePerTurn(ci)
	hasUltimate := cardHasActiveUltimate(ci)
	info["has_per_turn"] = hasPerTurn
	info["has_prayer"] = cardHasActivePrayer(ci)
	info["has_ultimate"] = hasUltimate
	if cardHasActiveSpellReaction(ci) {
		info["can_react"] = true
	}

	if hasPerTurn {
		info["per_turn_limit"] = perTurnLimit(ci)
		info["per_turn_label"] = "回合技"
		behavior := behaviorForNumber(ci.Card.Number)
		if labeler, ok := behavior.(PerTurnLabelBehavior); ok {
			if label := labeler.PerTurnLabel(ci); label != "" {
				info["per_turn_label"] = label
			}
		}
	}
	if requirement := summonDevourRequirement(ci); len(requirement) > 0 {
		info["devour_requirement"] = requirement
	}
	if requirement := summonDevourCardRequirement(ci); requirement.Count > 0 {
		info["devour_card_requirement"] = requirement
	}

	// Mark spell-like skills and spell scrolls.
	if isSpellLikeCard(ci.Card) {
		info["is_defense_only"] = isDefenseOnlySkill(ci.Card)
		info["is_sorcery"] = isSorcerySkill(ci.Card)
		info["needs_target"] = skillNeedsTargetInstance(ci)
		info["has_pierce"] = cardHasPierce(ci)
		if friendly, ok := behaviorForNumber(ci.Card.Number).(FriendlySpellTargetBehavior); ok && friendly.HasActiveFriendlySpellTarget(ci) && friendly.AllowsFriendlySpellTarget() {
			info["allows_friendly_target"] = true
		}
		info["can_attack"] = canUseSkillForPurpose(ci.Card, skillPurposeAttack)
		info["can_defend"] = canUseSkillForPurpose(ci.Card, skillPurposeDefend)
		info["can_attack_boost"] = canUseSkillForPurpose(ci.Card, skillPurposeAttackBoost)
		info["can_defense_boost"] = canUseSkillForPurpose(ci.Card, skillPurposeDefenseBoost)
		info["can_react"] = cardHasActiveSpellReaction(ci)
		info["can_boost"] = info["can_attack_boost"]
		info["spell_area"] = spellArea(ci)
	}

	return info
}

func cardsToInfo(cards []*CardInstance) []map[string]any {
	result := make([]map[string]any, len(cards))
	for i, c := range cards {
		result[i] = cardToInfo(c)
	}
	return result
}

func effectiveCurrentAttack(card *CardInstance) int {
	if card == nil {
		return 0
	}
	return max(card.CurrentAttack+card.AttackBonus, 0)
}

func (e *Engine) cardToInfo(ci *CardInstance) map[string]any {
	info := cardToInfo(ci)
	if info == nil {
		return nil
	}
	info["elements_gain"] = e.effectiveElementsGain(ci)
	info["can_attack_from_non_front"] = !e.hasEffectiveStatus(ci, StatusPetrify) && cardCanAttackFromNonFront(ci)
	if len(ci.BoundSkills) > 0 {
		info["bound_skills"] = e.cardsToInfo(ci.BoundSkills)
	}
	if len(ci.UnderCards) > 0 {
		info["under_cards"] = e.cardsToInfo(ci.UnderCards)
	}
	return info
}

func (e *Engine) cardsToInfo(cards []*CardInstance) []map[string]any {
	result := make([]map[string]any, len(cards))
	for i, c := range cards {
		result[i] = e.cardToInfo(c)
	}
	return result
}

func deckSummaryToInfo(deck []*CardInstance) []map[string]any {
	type summary struct {
		card  *model.Card
		count int
	}

	byNumber := map[string]*summary{}
	for _, ci := range deck {
		if ci == nil || ci.Card == nil {
			continue
		}
		number := ci.Card.Number
		entry := byNumber[number]
		if entry == nil {
			entry = &summary{card: ci.Card}
			byNumber[number] = entry
		}
		entry.count++
	}

	numbers := make([]string, 0, len(byNumber))
	for number := range byNumber {
		numbers = append(numbers, number)
	}
	sort.Strings(numbers)

	result := make([]map[string]any, 0, len(numbers))
	for _, number := range numbers {
		entry := byNumber[number]
		card := entry.card
		info := map[string]any{
			"number":           card.Number,
			"name":             card.Name,
			"type":             card.Type,
			"category":         card.Category,
			"tag":              card.Tag,
			"description":      card.Description,
			"attack":           card.Attack,
			"life":             card.Life,
			"power":            card.Power,
			"duration":         card.Duration,
			"elements_cost":    card.ElementsCost,
			"elements_gain":    card.ElementsGain,
			"elements_expense": card.ElementsExpense,
			"output_path":      card.OutputPath,
			"count":            entry.count,
			"is_terrain":       cards.IsTerrain(card.Number),
			"is_consumable":    cards.IsConsumable(card.Number),
			"is_equipment":     cards.IsEquipment(card.Number),
			"is_weapon":        cards.IsWeapon(card.Number),
		}
		addCardEffectMetadata(info, card)
		result = append(result, info)
	}

	return result
}

func addCardEffectMetadata(info map[string]any, card *model.Card) {
	if card == nil {
		return
	}
	if len(card.EffectCategories) > 0 {
		info["effect_categories"] = card.EffectCategories
	}
	if len(card.EffectOptionality) > 0 {
		info["effect_optionality"] = card.EffectOptionality
	}
}

func turnOrderLabel(playerID int, firstPlayer int) string {
	if playerID == firstPlayer {
		return "先手"
	}
	return "后手"
}

func (e *Engine) playerStateToInfo(ps *PlayerState, isOwner bool) map[string]any {
	info := map[string]any{
		"player_id":                      ps.PlayerID,
		"player_name":                    ps.PlayerName,
		"hero":                           e.cardToInfoForPlayer(ps, ps.Hero),
		"elements":                       ps.Elements,
		"strict_arcane":                  ps.StrictArcane,
		"shield":                         ps.Shield,
		"cannot_gain_shield":             ps.CannotGainShield,
		"next_red_moon_duration":         ps.NextRedMoonDuration,
		"next_red_moon_cooldown":         ps.NextRedMoonCooldown,
		"charge":                         ps.Charge,
		"temp_modifiers":                 ps.TempModifiers,
		"deck_count":                     len(ps.Deck),
		"graveyard":                      e.cardsToInfo(ps.Graveyard),
		"exile_count":                    len(ps.Exile),
		"discarded_hand_count_this_turn": ps.DiscardedHandCountThisTurn,
	}

	// Units grid
	units := [3][3]any{}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			units[col][row] = e.cardToInfoForPlayer(ps, ps.Units[col][row])
		}
	}
	info["units"] = units

	// Terrain grid
	terrain := [3][3]any{}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			terrain[col][row] = e.cardToInfoForPlayer(ps, ps.Terrain[col][row])
		}
	}
	info["terrain"] = terrain

	// Skills
	skills := make([]any, skillSlotCapacity(ps))
	for i := range skills {
		skills[i] = e.cardToInfoForPlayer(ps, ps.Skills[i])
	}
	info["skills"] = skills
	info["skill_slot_capacity"] = len(skills)

	// Equipment
	equipment := make([]any, equipmentSlotCapacity(ps))
	for i := range equipment {
		if !isOwner && ps.Equipment[i] != nil && ps.Equipment[i].IsSetCounter {
			equipment[i] = hiddenCounterInfo(ps.Equipment[i])
		} else {
			equipment[i] = e.cardToInfoForPlayer(ps, ps.Equipment[i])
		}
	}
	info["equipment"] = equipment
	info["equipment_slot_capacity"] = len(equipment)

	if isOwner {
		// Show full hand
		info["hand"] = e.cardsToInfoWithEffectiveCosts(ps, ps.Hand, false)
		info["deck_summary"] = deckSummaryToInfo(ps.Deck)
		info["exile"] = e.cardsToInfo(ps.Exile)
		if e.hasForesightOrbActive(ps.PlayerID) {
			info["top_deck_preview"] = e.cardsToInfo(ps.Deck[:min(3, len(ps.Deck))])
		}
		info["skill_pool"] = e.cardsToInfoWithEffectiveCosts(ps, ps.SkillPool, true)
	} else {
		// Only show count
		info["hand_count"] = len(ps.Hand)
		revealed := make([]*CardInstance, 0)
		for _, card := range ps.Hand {
			if ps.RevealedHand[card.InstanceID] {
				revealed = append(revealed, card)
			}
		}
		info["revealed_hand"] = e.cardsToInfo(revealed)
		info["skill_pool_count"] = len(ps.SkillPool)
	}

	return info
}

func (e *Engine) cardsToInfoWithEffectiveCosts(ps *PlayerState, cards []*CardInstance, learn bool) []map[string]any {
	result := make([]map[string]any, len(cards))
	for i, c := range cards {
		info := e.cardToInfoForPlayer(ps, c)
		if c != nil && c.Card != nil {
			if learn {
				info["effective_learn_cost"] = e.effectiveSkillLearnCost(ps, c)
			} else {
				info["effective_elements_cost"] = e.effectiveCardPlayCost(ps, c)
				info["action_base_play_cost"] = e.actionBaseCardPlayCost(ps, c)
			}
		}
		result[i] = info
	}
	return result
}

func (e *Engine) cardToInfoForPlayer(ps *PlayerState, card *CardInstance) map[string]any {
	info := e.cardToInfo(card)
	if ps == nil || card == nil || card.Card == nil {
		return info
	}
	if isSpellLikeCard(card.Card) {
		info["has_pierce"] = e.skillHasPierce(ps.PlayerID, card)
		info["spell_area"] = e.effectiveSpellArea(card)
		info["effective_defense_power"] = e.effectiveSkillPowerForPurpose(ps.PlayerID, card, skillPurposeDefend)
		info["effective_defense_boost_power"] = e.effectiveSkillPowerForPurpose(ps.PlayerID, card, skillPurposeDefenseBoost)
		info["effective_attack_power"] = e.effectiveSkillPowerForPurpose(ps.PlayerID, card, skillPurposeAttack)
		info["effective_attack_boost_power"] = e.effectiveSkillPowerForPurpose(ps.PlayerID, card, skillPurposeAttackBoost)
		info["action_base_attack_cost"] = e.actionBaseSkillUseCost(ps, card, skillPurposeAttack)
		info["action_base_attack_boost_cost"] = e.actionBaseSkillUseCost(ps, card, skillPurposeAttackBoost)
		info["action_base_defense_cost"] = e.actionBaseSkillUseCost(ps, card, skillPurposeDefend)
		info["action_base_defense_boost_cost"] = e.actionBaseSkillUseCost(ps, card, skillPurposeDefenseBoost)
		info["action_base_reaction_cost"] = e.actionBaseSkillUseCost(ps, card, skillPurposeReaction)
	}
	if len(card.BoundSkills) > 0 {
		boundSkills := make([]map[string]any, len(card.BoundSkills))
		for i, bound := range card.BoundSkills {
			boundSkills[i] = e.cardToInfoForPlayer(ps, bound)
		}
		info["bound_skills"] = boundSkills
	}
	if len(card.UnderCards) > 0 {
		underCards := make([]map[string]any, len(card.UnderCards))
		for i, under := range card.UnderCards {
			underCards[i] = e.cardToInfoForPlayer(ps, under)
		}
		info["under_cards"] = underCards
	}
	return info
}

func (e *Engine) effectiveSkillPowerForPurpose(playerID int, skill *CardInstance, purpose skillPurpose) int {
	if skill == nil || skill.Card == nil {
		return 0
	}
	return e.effectiveSkillPowerForPurposeWithData(playerID, skill, nil, purpose, map[string]any{"stat": "power"})
}

func (e *Engine) refreshPendingSpellPowerForModifiedSkill(playerID int, skill *CardInstance) {
	if e == nil || e.State.PendingSpell == nil || skill == nil {
		return
	}
	spell := e.State.PendingSpell
	if spell.AttackerID != playerID {
		return
	}
	if spell.Skill != skill {
		foundBoost := false
		for _, boost := range spell.BoostSkills {
			if boost == skill {
				foundBoost = true
				break
			}
		}
		if !foundBoost {
			return
		}
	}
	powerTargets := append([]SpellTarget{spell.Target}, spell.ExtraTargets...)
	spell.TotalPower = e.effectiveSpellPower(playerID, spell.Skill, spell.BoostSkills, powerTargets...)
	spell.PowerSources = e.spellPowerSources(playerID, spell.Skill, spell.BoostSkills, spell.TotalPower, powerTargets...)
}
