package game

import (
	"fmt"
	"strings"

	"eraofarcane/model"
)

type skillPurpose string

const (
	skillPurposeAttack skillPurpose = "attack"
	skillPurposeDefend skillPurpose = "defend"
	skillPurposeBoost  skillPurpose = "boost"
)

func skillUseCost(card *model.Card) map[string]int {
	if len(card.ElementsExpense) > 0 {
		return card.ElementsExpense
	}
	return card.ElementsCost
}

func mergeElementCosts(costs ...map[string]int) map[string]int {
	merged := make(map[string]int)
	for _, cost := range costs {
		for elem, amount := range cost {
			merged[elem] += amount
		}
	}
	return merged
}

func stringsFromAnySlice(values []any) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if str, ok := value.(string); ok && str != "" {
			result = append(result, str)
		}
	}
	return result
}

func (e *Engine) collectSkillUses(ps *PlayerState, ids []string, purpose skillPurpose, reserved map[string]bool) ([]*CardInstance, map[string]int, error) {
	skills := make([]*CardInstance, 0, len(ids))
	totalCost := make(map[string]int)
	seen := make(map[string]bool)
	for id := range reserved {
		seen[id] = true
	}

	for _, id := range ids {
		if seen[id] {
			return nil, nil, fmt.Errorf("skill %s selected more than once", id)
		}
		seen[id] = true

		skill := e.findSkill(ps, id)
		if skill == nil {
			return nil, nil, fmt.Errorf("skill not found: %s", id)
		}
		if err := e.validateSkillForPurpose(skill, purpose); err != nil {
			return nil, nil, fmt.Errorf("skill %s cannot be used for %s: %w", id, purpose, err)
		}

		skills = append(skills, skill)
		for elem, amount := range skillUseCost(skill.Card) {
			totalCost[elem] += amount
		}
	}

	return skills, totalCost, nil
}

func (e *Engine) validateSkillForPurpose(skill *CardInstance, purpose skillPurpose) error {
	if err := e.validateReadySkill(skill); err != nil {
		return err
	}

	switch purpose {
	case skillPurposeAttack:
		if !canUseSkillToAttack(skill.Card) {
			return fmt.Errorf("skill cannot be used to attack")
		}
	case skillPurposeDefend:
		if !canUseSkillToDefend(skill.Card) {
			return fmt.Errorf("skill cannot be used to defend")
		}
	case skillPurposeBoost:
		if !canUseSkillToBoost(skill.Card) {
			return fmt.Errorf("skill cannot be used to boost")
		}
	default:
		return fmt.Errorf("unknown skill purpose: %s", purpose)
	}

	return nil
}

func (e *Engine) validateReadySkill(skill *CardInstance) error {
	if skill.IsHorizontal {
		return fmt.Errorf("skill is horizontal (already used)")
	}
	if skill.Statuses[StatusCooldown] > 0 {
		return fmt.Errorf("skill is on cooldown")
	}
	if skill.Statuses[StatusPetrify] > 0 {
		return fmt.Errorf("skill is petrified")
	}
	return nil
}

func tapSkills(skills []*CardInstance) {
	for _, skill := range skills {
		skill.IsHorizontal = true
	}
}

func totalSkillPower(skills []*CardInstance) int {
	total := 0
	for _, skill := range skills {
		total += max(skill.Card.Power, 0)
	}
	return total
}

func skillIDSet(skills []*CardInstance) map[string]bool {
	ids := make(map[string]bool, len(skills))
	for _, skill := range skills {
		ids[skill.InstanceID] = true
	}
	return ids
}

func (e *Engine) validateSpellTarget(playerID int, skill *CardInstance, target SpellTarget) error {
	if target.Type != "unit" {
		return fmt.Errorf("unsupported spell target type: %s", target.Type)
	}
	if !target.Position.Valid() {
		return fmt.Errorf("invalid target position")
	}

	opponent := e.State.Players[1-playerID]
	if opponent.Units[target.Position.Col][target.Position.Row] == nil {
		return fmt.Errorf("no enemy unit at target position")
	}

	hasPierce := HasKeyword(skill.Card.Description, KW_Pierce)
	if !e.IsInSpellRange(playerID, target.Position.Col, target.Position.Row, hasPierce) {
		return fmt.Errorf("target is not in spell range")
	}

	return nil
}

func canUseSkillToAttack(card *model.Card) bool {
	if !card.IsSkill() {
		return false
	}
	if isDefenseOnlySkill(card) {
		return false
	}
	if strings.Contains(card.Description, "不可用于攻击") {
		return false
	}
	return true
}

func canUseSkillToDefend(card *model.Card) bool {
	if !card.IsSkill() {
		return false
	}
	if isSorcerySkill(card) {
		return false
	}
	if strings.Contains(card.Description, "不可用于防御") {
		return false
	}
	return card.Power > 0
}

func canUseSkillToBoost(card *model.Card) bool {
	if !card.IsSkill() {
		return false
	}
	if isSorcerySkill(card) {
		return false
	}
	desc := card.Description
	if strings.Contains(desc, "无法用于强化") || strings.Contains(desc, "不可用于强化") {
		return false
	}
	if strings.Contains(desc, "无法强化或被强化") {
		return false
	}
	return card.Power > 0
}
