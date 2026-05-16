package game

import (
	"fmt"

	"eraofarcane/model"
)

type skillPurpose string

const (
	skillPurposeAttack       skillPurpose = "attack"
	skillPurposeDefend       skillPurpose = "defend"
	skillPurposeBoost        skillPurpose = "boost"
	skillPurposeAttackBoost  skillPurpose = "attack_boost"
	skillPurposeDefenseBoost skillPurpose = "defense_boost"
)

func isBoostPurpose(purpose skillPurpose) bool {
	return purpose == skillPurposeBoost || purpose == skillPurposeAttackBoost || purpose == skillPurposeDefenseBoost
}

func skillUseCost(card *model.Card) map[string]int {
	if len(card.ElementsExpense) > 0 {
		return card.ElementsExpense
	}
	return card.ElementsCost
}

func copyElementCost(cost map[string]int) map[string]int {
	copied := make(map[string]int)
	for elem, amount := range cost {
		if amount > 0 {
			copied[elem] = amount
		}
	}
	return copied
}

func (e *Engine) effectiveSkillUseCost(ps *PlayerState, skill *CardInstance) map[string]int {
	if skill == nil || skill.Card == nil {
		return map[string]int{}
	}
	cost := copyElementCost(skillUseCost(skill.Card))
	ctx := &EffectContext{
		Engine:     e,
		Source:     skill,
		PlayerID:   ps.PlayerID,
		OpponentID: 1 - ps.PlayerID,
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || fieldCard.Statuses[StatusPetrify] > 0 {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		if modifier, ok := behavior.(SkillUseCostModifier); ok {
			ctx.Target = fieldCard
			modifier.ModifySkillUseCost(ctx, cost)
		}
	}
	if e.nextSkillCostZeroModifier(ps, skill) != nil {
		for elem := range cost {
			delete(cost, elem)
		}
	}
	return cost
}

func (e *Engine) effectiveCardPlayCost(ps *PlayerState, card *CardInstance) map[string]int {
	if card == nil || card.Card == nil {
		return map[string]int{}
	}
	cost := copyElementCost(card.Card.ElementsCost)
	ctx := &EffectContext{
		Engine:     e,
		Target:     card,
		PlayerID:   ps.PlayerID,
		OpponentID: 1 - ps.PlayerID,
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || fieldCard.Statuses[StatusPetrify] > 0 {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		if modifier, ok := behavior.(CardPlayCostModifier); ok {
			ctx.Source = fieldCard
			modifier.ModifyCardPlayCost(ctx, card, cost)
		}
	}
	return cost
}

func (e *Engine) effectiveSkillLearnCost(ps *PlayerState, skill *CardInstance) map[string]int {
	return e.effectiveCardPlayCost(ps, skill)
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
		for elem, amount := range e.effectiveSkillUseCost(ps, skill) {
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
		if !canUseSkillForPurpose(skill.Card, skillPurposeAttack) {
			return fmt.Errorf("skill cannot be used to attack")
		}
	case skillPurposeDefend:
		if !canUseSkillForPurpose(skill.Card, skillPurposeDefend) {
			return fmt.Errorf("skill cannot be used to defend")
		}
	case skillPurposeBoost, skillPurposeAttackBoost, skillPurposeDefenseBoost:
		if !canUseSkillForPurpose(skill.Card, purpose) {
			return fmt.Errorf("skill cannot be used to boost")
		}
	default:
		return fmt.Errorf("unknown skill purpose: %s", purpose)
	}

	ps := e.State.Players[skill.OwnerID]
	ctx := &EffectContext{
		Engine:     e,
		Target:     skill,
		PlayerID:   skill.OwnerID,
		OpponentID: 1 - skill.OwnerID,
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || fieldCard.Statuses[StatusPetrify] > 0 {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		modifier, ok := behavior.(SkillUsePermissionModifier)
		if !ok {
			continue
		}
		ctx.Source = fieldCard
		if err := modifier.ValidateSkillUse(ctx, skill, purpose); err != nil {
			return err
		}
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
	if skill.Statuses[StatusSeal] > 0 {
		return fmt.Errorf("skill is sealed")
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
		total += max(skill.Card.Power+skill.PowerBonus, 0)
	}
	return total
}

func (e *Engine) effectiveSpellPower(playerID int, skill *CardInstance, boostSkills []*CardInstance, targets ...SpellTarget) int {
	extra := map[string]any{}
	if len(targets) > 0 {
		extra["spell_target"] = targets[0]
		if unit := e.spellTargetUnit(1-playerID, targets[0]); unit != nil {
			extra["spell_target_unit"] = unit
		}
	}
	power := e.skillContributionStatsWithData(playerID, skill, skill, skillPurposeAttack, extra).PowerBonus
	for _, boostSkill := range boostSkills {
		power += e.skillContributionStatsWithData(playerID, boostSkill, skill, skillPurposeAttackBoost, extra).PowerBonus
	}
	power += e.spellStatBonusesWithData(playerID, skill, skillPurposeAttack, extra).PowerBonus
	power += e.genericSpellBonus(playerID, skill, "威")
	power += e.temporarySpellPowerBonus(playerID, skill)
	return max(power, 0)
}

func (e *Engine) effectiveSpellDamage(playerID int, skill *CardInstance, baseDamage int, boostSkills []*CardInstance) int {
	damage := baseDamage + e.genericSpellBonus(playerID, skill, "攻")
	damage += e.skillContributionStats(playerID, skill, skill, skillPurposeAttack).DamageBonus
	for _, boostSkill := range boostSkills {
		damage += e.skillContributionStats(playerID, boostSkill, skill, skillPurposeAttackBoost).DamageBonus
	}
	damage += e.spellStatBonuses(playerID, skill, skillPurposeAttack).DamageBonus
	return max(damage, 0)
}

func (e *Engine) totalEffectiveSkillPower(playerID int, skills []*CardInstance, purpose skillPurpose) int {
	total := 0
	for _, skill := range skills {
		total += e.skillContributionStats(playerID, skill, nil, purpose).PowerBonus
	}
	return max(total, 0)
}

func (e *Engine) skillContributionStats(playerID int, skill *CardInstance, target *CardInstance, purpose skillPurpose) SpellStats {
	return e.skillContributionStatsWithData(playerID, skill, target, purpose, nil)
}

func (e *Engine) skillContributionStatsWithData(playerID int, skill *CardInstance, target *CardInstance, purpose skillPurpose, extra map[string]any) SpellStats {
	stats := SpellStats{PowerBonus: max(skill.Card.Power+skill.PowerBonus, 0)}
	if weak := skill.Statuses[StatusWeaken]; weak > 0 {
		stats.PowerBonus = max(stats.PowerBonus-weak, 0)
	}
	behavior := globalRegistry.GetBehavior(skill.Card.Number)
	if modifier, ok := behavior.(SkillContributionModifier); ok {
		data := map[string]any{"purpose": string(purpose)}
		for key, value := range extra {
			data[key] = value
		}
		ctx := &EffectContext{
			Engine:     e,
			Source:     skill,
			Target:     target,
			PlayerID:   playerID,
			OpponentID: 1 - playerID,
			ExtraData:  data,
		}
		modifier.ModifySkillContribution(ctx, &stats)
	}
	stats.PowerBonus = max(stats.PowerBonus, 0)
	stats.DamageBonus = max(stats.DamageBonus, 0)
	return stats
}

func (e *Engine) spellTargetUnit(defenderID int, target SpellTarget) *CardInstance {
	if target.Type != "unit" || !target.Position.Valid() {
		return nil
	}
	return e.State.Players[defenderID].Units[target.Position.Col][target.Position.Row]
}

func (e *Engine) spellStatBonuses(playerID int, skill *CardInstance, purpose skillPurpose) SpellStats {
	return e.spellStatBonusesWithData(playerID, skill, purpose, nil)
}

func (e *Engine) spellStatBonusesWithData(playerID int, skill *CardInstance, purpose skillPurpose, extra map[string]any) SpellStats {
	stats := SpellStats{}
	ps := e.State.Players[playerID]
	data := map[string]any{"purpose": string(purpose)}
	for key, value := range extra {
		data[key] = value
	}
	ctx := &EffectContext{
		Engine:     e,
		Target:     skill,
		PlayerID:   playerID,
		OpponentID: 1 - playerID,
		ExtraData:  data,
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || fieldCard.Statuses[StatusPetrify] > 0 {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		modifier, ok := behavior.(SpellStatModifier)
		if !ok {
			continue
		}
		ctx.Source = fieldCard
		modifier.ModifySpellStats(ctx, &stats)
	}
	return stats
}

func skillIDSet(skills []*CardInstance) map[string]bool {
	ids := make(map[string]bool, len(skills))
	for _, skill := range skills {
		ids[skill.InstanceID] = true
	}
	return ids
}

func (e *Engine) validateSpellTarget(playerID int, skill *CardInstance, target SpellTarget) error {
	if !skillNeedsTargetInstance(skill) {
		if target.Type == "" || target.Type == "none" {
			return nil
		}
	}
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

	hasPierce := cardHasPierce(skill)
	if !e.IsInSpellRange(playerID, target.Position.Col, target.Position.Row, hasPierce) {
		return fmt.Errorf("target is not in spell range")
	}

	return nil
}
