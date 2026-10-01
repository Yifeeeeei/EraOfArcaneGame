package game

import (
	"fmt"
)

// handleLearnSkill handles learning a skill from the skill pool
func (e *Engine) handleLearnSkill(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}
	instanceID, _ := action.Data["instance_id"].(string)
	replaceID, _ := action.Data["replace_id"].(string) // optional: which skill to replace

	ps := e.State.Players[playerID]

	// Find skill in skill pool
	var skill *CardInstance
	var poolIdx int
	for i, s := range ps.SkillPool {
		if s.InstanceID == instanceID {
			skill = s
			poolIdx = i
			break
		}
	}
	if skill == nil {
		return fmt.Errorf("skill not found in skill pool")
	}
	if err := e.validateSkillLearnPermissionModifiers(playerID, skill); err != nil {
		return err
	}

	// Check cost
	cost := e.effectiveSkillLearnCost(ps, skill)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}
	if err := e.validateOwnCost(playerID, skill, cost, action); err != nil {
		return err
	}

	// Find slot
	slotIdx := -1
	var replacedSkill *CardInstance
	if replaceID != "" {
		// Replace existing skill
		for i := 0; i < skillSlotCapacity(ps); i++ {
			if ps.Skills[i] != nil && ps.Skills[i].InstanceID == replaceID {
				if ps.Skills[i].IsHorizontal {
					return fmt.Errorf("can only replace vertical skills")
				}
				if !skillAllowedInSlot(ps, skill, i) {
					return fmt.Errorf("skill cannot be learned into this slot")
				}
				replacedSkill = ps.Skills[i]
				slotIdx = i
				break
			}
		}
		if slotIdx == -1 {
			return fmt.Errorf("replacement skill not found")
		}
	} else {
		// Find empty slot
		for i := 0; i < skillSlotCapacity(ps); i++ {
			if ps.Skills[i] == nil && skillAllowedInSlot(ps, skill, i) {
				slotIdx = i
				break
			}
		}
		if slotIdx == -1 {
			return fmt.Errorf("skill area is full, must replace an existing skill")
		}
	}

	// Pay cost and place
	if !e.payCostForCardAction(ps, skill, cost, cost, paymentPurposeLearn, action) {
		return fmt.Errorf("invalid payment")
	}
	e.notifyCardPlayCostPaid(ps, skill)
	e.consumeEarthSkillLearnCostModifier(ps, skill)
	ps.SkillPool = append(ps.SkillPool[:poolIdx], ps.SkillPool[poolIdx+1:]...)
	if replacedSkill != nil {
		ps.Skills[slotIdx] = nil
		returnSkillToPool(replacedSkill)
		ps.SkillPool = append(ps.SkillPool, replacedSkill)
		e.notifyCardStateChanges(CardStateChange{Card: replacedSkill, LeftField: true})
	}
	skill.IsHorizontal = true
	skill.SlotIndex = slotIdx
	skill.EnterTurn = e.State.TurnNumber
	e.ApplyKeywordOnEnter(skill)
	ps.Skills[slotIdx] = skill
	e.applyNextLearnedSkillHasteModifier(playerID, skill)

	e.emit(GameEvent{
		Type:   "learn_skill",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(skill),
			"slot":     slotIdx,
			"elements": ps.Elements,
		},
	})
	e.triggerEffects(TriggerOnEnter, skill, nil, nil)
	learnData := map[string]any{"entered_player": playerID, "learned_skill": true}
	e.notifyCardEntered(playerID, skill, learnData)
	e.promptOpponentCounterTrap(playerID, TriggerOnCardEnter, skill, learnData, nil)

	return nil
}

func (e *Engine) learnSkillFromPoolWithoutCost(playerID int, instanceID string, replaceID string) bool {
	ps := e.State.Players[playerID]
	if ps == nil {
		return false
	}
	var skill *CardInstance
	poolIdx := -1
	for i, s := range ps.SkillPool {
		if s != nil && s.InstanceID == instanceID {
			skill = s
			poolIdx = i
			break
		}
	}
	if skill == nil || skill.Card == nil || !skill.Card.IsSkill() {
		return false
	}
	if err := e.validateSkillLearnPermissionModifiers(playerID, skill); err != nil {
		return false
	}

	slotIdx := -1
	var replacedSkill *CardInstance
	if replaceID != "" {
		for i := 0; i < skillSlotCapacity(ps); i++ {
			if ps.Skills[i] != nil && ps.Skills[i].InstanceID == replaceID && !ps.Skills[i].IsHorizontal {
				if !skillAllowedInSlot(ps, skill, i) {
					return false
				}
				replacedSkill = ps.Skills[i]
				slotIdx = i
				break
			}
		}
	} else {
		for i := 0; i < skillSlotCapacity(ps); i++ {
			if ps.Skills[i] == nil && skillAllowedInSlot(ps, skill, i) {
				slotIdx = i
				break
			}
		}
	}
	if slotIdx == -1 {
		return false
	}

	ps.SkillPool = append(ps.SkillPool[:poolIdx], ps.SkillPool[poolIdx+1:]...)
	if replacedSkill != nil {
		ps.Skills[slotIdx] = nil
		returnSkillToPool(replacedSkill)
		ps.SkillPool = append(ps.SkillPool, replacedSkill)
		e.notifyCardStateChanges(CardStateChange{Card: replacedSkill, LeftField: true})
	}
	skill.IsHorizontal = true
	skill.SlotIndex = slotIdx
	skill.EnterTurn = e.State.TurnNumber
	e.ApplyKeywordOnEnter(skill)
	ps.Skills[slotIdx] = skill
	e.applyNextLearnedSkillHasteModifier(playerID, skill)

	e.emit(GameEvent{
		Type:   "learn_skill",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(skill),
			"slot":     slotIdx,
			"elements": ps.Elements,
		},
	})
	e.triggerEffects(TriggerOnEnter, skill, nil, nil)
	learnData := map[string]any{"entered_player": playerID, "learned_skill": true}
	e.notifyCardEntered(playerID, skill, learnData)
	e.promptOpponentCounterTrap(playerID, TriggerOnCardEnter, skill, learnData, nil)
	return true
}

func (e *Engine) applyNextLearnedSkillHasteModifier(playerID int, skill *CardInstance) {
	ps := e.State.Players[playerID]
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.Type != TempModNextLearnedSkillHaste || modifier.RemainingUses == 0 {
			continue
		}
		if modifier.Element != "" && (skill == nil || skill.Card == nil || skill.Card.Category != modifier.Element) {
			continue
		}
		skill.IsHorizontal = false
		e.removeTemporaryModifier(playerID, modifier.ID)
		break
	}
}

func returnSkillToPool(skill *CardInstance) {
	if skill == nil {
		return
	}
	skill.IsHorizontal = true
	skill.Position = nil
	skill.SlotIndex = -1
	skill.EnterTurn = 0
	skill.UsedThisTurn = 0
	skill.UltimateUsed = false
	skill.Statuses = make(map[string]int)
	skill.ElementsGainBonus = make(map[string]int)
	skill.ElementsGainSet = nil
	skill.PowerBonus = 0
	skill.AttackBonus = 0
	skill.AttachedBehaviors = nil
}
