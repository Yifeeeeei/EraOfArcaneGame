package game

import (
	"fmt"
	"strconv"
	"strings"

	"eraofarcane/model"
)

type skillPurpose string

const (
	skillPurposeAttack       skillPurpose = "attack"
	skillPurposeDefend       skillPurpose = "defend"
	skillPurposeBoost        skillPurpose = "boost"
	skillPurposeAttackBoost  skillPurpose = "attack_boost"
	skillPurposeDefenseBoost skillPurpose = "defense_boost"
	skillPurposeReaction     skillPurpose = "reaction"
)

func isBoostPurpose(purpose skillPurpose) bool {
	return purpose == skillPurposeBoost || purpose == skillPurposeAttackBoost || purpose == skillPurposeDefenseBoost
}

func parsePositiveInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0
	}
	return n
}

func skillUseCost(card *model.Card) map[string]int {
	if card == nil {
		return map[string]int{}
	}
	return card.ElementsExpense
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

func restrictedEquipmentSubtype(card *model.Card) string {
	if card == nil || card.Type != model.CardTypeItem || !strings.Contains(card.Tag, "装备") {
		return ""
	}
	switch {
	case strings.Contains(card.Tag, "武器"):
		return "武器"
	case strings.Contains(card.Tag, "防具"):
		return "防具"
	case strings.Contains(card.Tag, "饰物"):
		return "饰物"
	case strings.Contains(card.Tag, "神器"):
		return "神器"
	default:
		return ""
	}
}

func isEquipmentCard(card *model.Card) bool {
	return card != nil && card.Type == model.CardTypeItem && strings.Contains(card.Tag, "装备")
}

func (e *Engine) effectiveSkillUseCost(ps *PlayerState, skill *CardInstance) map[string]int {
	return e.effectiveSkillUseCostForPurpose(ps, skill, skillPurposeAttack)
}

func (e *Engine) effectiveAttackCost(ps *PlayerState, attacker *CardInstance) map[string]int {
	if e == nil || ps == nil || attacker == nil || attacker.Card == nil {
		return map[string]int{}
	}
	behavior := globalRegistry.GetBehavior(attacker.Card.Number)
	costBehavior, ok := behavior.(AttackCostBehavior)
	if !ok || !costBehavior.HasActiveAttackCost(attacker) {
		return map[string]int{}
	}
	ctx := &EffectContext{
		Engine:     e,
		Source:     attacker,
		PlayerID:   ps.PlayerID,
		OpponentID: 1 - ps.PlayerID,
	}
	return copyElementCost(costBehavior.AttackCost(ctx))
}

func (e *Engine) effectiveSkillUseCostForPurpose(ps *PlayerState, skill *CardInstance, purpose skillPurpose) map[string]int {
	return e.effectiveSkillUseCostForPurposeWithData(ps, skill, purpose, nil)
}

func (e *Engine) effectiveSkillUseCostForPurposeWithData(ps *PlayerState, skill *CardInstance, purpose skillPurpose, extra map[string]any) map[string]int {
	if skill == nil || skill.Card == nil {
		return map[string]int{}
	}
	cost := copyElementCost(skillUseCost(skill.Card))
	data := map[string]any{"purpose": string(purpose)}
	for key, value := range extra {
		data[key] = value
	}
	ctx := &EffectContext{
		Engine:     e,
		Source:     skill,
		PlayerID:   ps.PlayerID,
		OpponentID: 1 - ps.PlayerID,
		ExtraData:  data,
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		if modifier, ok := behavior.(SkillUseCostModifier); ok && modifier.HasActiveSkillUseCostModifier(fieldCard) {
			ctx.Target = fieldCard
			modifier.ModifySkillUseCost(ctx, cost)
		}
	}
	for _, elem := range model.AllElements {
		key := "使用费用" + elem + "-"
		for status, amount := range skill.Statuses {
			if amount <= 0 || !strings.HasPrefix(status, key) {
				continue
			}
			reduceCost(cost, elem, parsePositiveInt(status[len(key):])*amount)
		}
		extraKey := skillUseExtraCostStatusPrefix + elem
		for status, amount := range skill.Statuses {
			if amount <= 0 || !strings.HasPrefix(status, extraKey) {
				continue
			}
			extra := parsePositiveInt(status[len(extraKey):])
			if extra <= 0 {
				extra = 1
			}
			cost[elem] += extra * amount
		}
		if reduction := e.temporaryNextSkillUseCostMinus(ps, skill, elem); reduction > 0 {
			reduceCost(cost, elem, reduction)
		}
	}
	if !isBoostPurpose(purpose) && e.nextSkillCostZeroModifier(ps, skill) != nil {
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
	if neutral := card.Statuses[StatusEntryCostNeutralAmount]; neutral > 0 {
		cost = map[string]int{model.ElementArcane: neutral}
	}
	ctx := &EffectContext{
		Engine:     e,
		Target:     card,
		PlayerID:   ps.PlayerID,
		OpponentID: 1 - ps.PlayerID,
	}
	for _, player := range e.State.Players {
		if player == nil {
			continue
		}
		for _, fieldCard := range e.getAllFieldCards(player) {
			if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
				continue
			}
			behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
			if modifier, ok := behavior.(GlobalCardPlayCostModifier); ok && modifier.HasActiveGlobalCardPlayCostModifier(fieldCard) {
				ctx.Source = fieldCard
				modifier.ModifyGlobalCardPlayCost(ctx, card, cost)
			}
		}
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		if modifier, ok := behavior.(CardPlayCostModifier); ok && modifier.HasActiveCardPlayCostModifier(fieldCard) {
			ctx.Source = fieldCard
			modifier.ModifyCardPlayCost(ctx, card, cost)
		}
	}
	for _, elem := range model.AllElements {
		key := "入场费用" + elem + "-"
		for status, amount := range card.Statuses {
			if amount <= 0 || !strings.HasPrefix(status, key) {
				continue
			}
			reduceCost(cost, elem, parsePositiveInt(status[len(key):])*amount)
		}
		if reduction := e.temporaryNextCardPlayCostMinus(ps, card, elem); reduction > 0 {
			reduceCost(cost, elem, reduction)
		}
	}
	if modifier, ok := globalRegistry.GetBehavior(card.Card.Number).(SelfCardPlayCostModifier); ok && modifier.HasActiveSelfCardPlayCostModifier(card) {
		ctx.Source = card
		modifier.ModifySelfCardPlayCost(ctx, cost)
	}
	if card.Statuses[entryCostZeroStatus] > 0 {
		for elem := range cost {
			cost[elem] = 0
		}
	}
	return cost
}

func (e *Engine) effectiveSkillLearnCost(ps *PlayerState, skill *CardInstance) map[string]int {
	cost := e.effectiveCardPlayCost(ps, skill)
	if modifier := e.nextEarthSkillLearnCostMinus(ps, skill); modifier != nil {
		amount := modifier.Amount
		if amount <= 0 {
			amount = 2
		}
		reduceCost(cost, model.ElementEarth, amount)
	}
	return cost
}

func (e *Engine) validateSkillLearnPermissionModifiers(playerID int, skill *CardInstance) error {
	if skill == nil || skill.Card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	if e.timeCycleLockActive() {
		return fmt.Errorf("time cycle prevents learning skills")
	}
	if e.playerCannotLearnElementSkill(ps, skill.Card.Category) {
		return fmt.Errorf("cannot learn %s skills", skill.Card.Category)
	}
	ctx := &EffectContext{
		Engine:     e,
		Target:     skill,
		PlayerID:   playerID,
		OpponentID: 1 - playerID,
	}
	if skillBehavior, ok := behaviorForNumber(skill.Card.Number).(SkillLearnPermissionModifier); ok && skillBehavior.HasActiveSkillLearnPermission(skill) {
		ctx.Source = skill
		if err := skillBehavior.ValidateSkillLearn(ctx, skill); err != nil {
			return err
		}
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			continue
		}
		if fieldCard == skill {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		modifier, ok := behavior.(SkillLearnPermissionModifier)
		if !ok || !modifier.HasActiveSkillLearnPermission(fieldCard) {
			continue
		}
		ctx.Source = fieldCard
		if err := modifier.ValidateSkillLearn(ctx, skill); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) playerCannotLearnElementSkill(ps *PlayerState, element string) bool {
	if ps == nil || element == "" {
		return false
	}
	for _, modifier := range ps.TempModifiers {
		if modifier.Type == TempModCannotLearnElementSkill && modifier.Element == element {
			return true
		}
	}
	return false
}

func (e *Engine) timeCycleLockActive() bool {
	if e == nil {
		return false
	}
	for _, ps := range e.State.Players {
		if ps == nil {
			continue
		}
		for _, card := range e.getAllFieldCards(ps) {
			if card != nil && card.Card != nil && card.Card.Number == "3411101" && !e.hasEffectiveStatus(card, StatusPetrify) {
				return true
			}
		}
	}
	return false
}

func (e *Engine) notifyCardPlayCostPaid(ps *PlayerState, card *CardInstance) {
	if ps == nil || card == nil || card.Card == nil {
		return
	}
	ctx := &EffectContext{
		Engine:     e,
		Target:     card,
		PlayerID:   ps.PlayerID,
		OpponentID: 1 - ps.PlayerID,
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		paid, ok := behavior.(CardPlayCostPaidBehavior)
		if !ok || !paid.HasActiveCardPlayCostPaid(fieldCard) {
			continue
		}
		ctx.Source = fieldCard
		paid.OnCardPlayCostPaid(ctx, card)
	}
	e.consumeNextCardPlayCostModifiers(ps, card)
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

func (e *Engine) handLimitForPlayer(ps *PlayerState) int {
	if ps == nil {
		return e.State.HandLimit
	}
	limit := e.State.HandLimit
	if ps.Hero != nil && ps.Hero.Card != nil && ps.Hero.Card.Number == "4311002" {
		limit++
	}
	opponentID := 1 - ps.PlayerID
	if opponentID >= 0 && opponentID < len(e.State.Players) && e.playerHasActiveCard(e.State.Players[opponentID], "1311103") {
		limit--
	}
	if limit < 0 {
		return 0
	}
	return limit
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

func skillSlotCapacity(ps *PlayerState) int {
	capacity := baseSkillSlotCapacity(ps)
	if hasActiveSpiritCandle(ps) {
		capacity += 2
	}
	if capacity > MaxSkillSlots {
		return MaxSkillSlots
	}
	return capacity
}

func baseSkillSlotCapacity(ps *PlayerState) int {
	if ps == nil {
		return BaseSkillSlots
	}
	capacity := BaseSkillSlots
	for _, equipment := range ps.Equipment {
		if equipment != nil && equipment.Card != nil && equipment.Card.Number == "2021002" && equipment.Statuses[StatusPetrify] <= 0 {
			capacity++
			break
		}
	}
	if capacity > MaxSkillSlots {
		return MaxSkillSlots
	}
	return capacity
}

func hasActiveSpiritCandle(ps *PlayerState) bool {
	if ps == nil {
		return false
	}
	for _, equipment := range ps.Equipment {
		if equipment != nil && equipment.Card != nil && equipment.Card.Number == "2611102" && equipment.Statuses[StatusPetrify] <= 0 {
			return true
		}
	}
	return false
}

func skillAllowedInSlot(ps *PlayerState, skill *CardInstance, slotIdx int) bool {
	if slotIdx < 0 || slotIdx >= skillSlotCapacity(ps) {
		return false
	}
	if slotIdx < baseSkillSlotCapacity(ps) {
		return true
	}
	if !hasActiveSpiritCandle(ps) || skill == nil || skill.Card == nil {
		return false
	}
	return hasCardTag(skill.Card, "灵媒") || hasCardTag(skill.Card, "神秘")
}

func equipmentSlotCapacity(ps *PlayerState) int {
	if ps == nil {
		return BaseEquipmentSlots
	}
	capacity := BaseEquipmentSlots
	for _, equipment := range ps.Equipment {
		if equipment != nil && equipment.Card != nil && equipment.Card.Number == "2021017" && equipment.Statuses[StatusPetrify] <= 0 {
			capacity += 3
			break
		}
	}
	for _, equipment := range ps.Equipment {
		if equipment != nil && equipment.Card != nil && equipment.Card.Number == "2021105" && equipment.Statuses[StatusPetrify] <= 0 {
			capacity++
			break
		}
	}
	if capacity > MaxEquipmentSlots {
		return MaxEquipmentSlots
	}
	return capacity
}

func playerCanEquipDuplicateSubtypes(ps *PlayerState) bool {
	if ps == nil {
		return false
	}
	for _, equipment := range ps.Equipment {
		if equipment != nil && equipment.Card != nil && equipment.Card.Number == "2021105" && equipment.Statuses[StatusPetrify] <= 0 {
			return true
		}
	}
	return false
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
		for elem, amount := range e.effectiveSkillUseCostForPurpose(ps, skill, purpose) {
			totalCost[elem] += amount
		}
	}

	return skills, totalCost, nil
}

func (e *Engine) validateSkillForPurpose(skill *CardInstance, purpose skillPurpose) error {
	if e.timeCycleLockActive() && (skill == nil || skill.Card == nil || skill.Card.Number != "3411101") {
		return fmt.Errorf("time cycle prevents skill use")
	}
	if skill.Statuses[StatusCannotUseSkillUntilTurn] >= e.State.TurnNumber {
		return fmt.Errorf("skill cannot be used this turn")
	}
	if purpose != skillPurposeReaction || skill.Card.IsSkill() {
		if err := e.validateReadySkill(skill); err != nil {
			return err
		}
	} else if e.hasEffectiveStatus(skill, StatusPetrify) {
		return fmt.Errorf("card is petrified")
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
	case skillPurposeReaction:
		behavior, ok := behaviorForNumber(skill.Card.Number).(SpellReactionBehavior)
		if !ok || !behavior.HasActiveSpellReaction(skill) {
			return fmt.Errorf("skill cannot react to spells")
		}
		ctx := &EffectContext{
			Engine:     e,
			Source:     skill,
			PlayerID:   skill.OwnerID,
			OpponentID: 1 - skill.OwnerID,
		}
		if !behavior.CanReactToSpell(ctx, e.State.PendingSpell) {
			return fmt.Errorf("skill cannot react to this spell")
		}
	default:
		return fmt.Errorf("unknown skill purpose: %s", purpose)
	}

	if err := e.validateSkillUsePermissionModifiers(skill, purpose); err != nil {
		return err
	}

	return nil
}

func (e *Engine) validateSkillUsePermissionModifiers(skill *CardInstance, purpose skillPurpose) error {
	ps := e.State.Players[skill.OwnerID]
	ctx := &EffectContext{
		Engine:     e,
		Target:     skill,
		PlayerID:   skill.OwnerID,
		OpponentID: 1 - skill.OwnerID,
	}
	if skillBehavior, ok := behaviorForNumber(skill.Card.Number).(SkillUsePermissionModifier); ok && skillBehavior.HasActiveSkillUsePermissionModifier(skill) {
		ctx.Source = skill
		if err := skillBehavior.ValidateSkillUse(ctx, skill, purpose); err != nil {
			return err
		}
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			continue
		}
		if fieldCard == skill {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		modifier, ok := behavior.(SkillUsePermissionModifier)
		if !ok || !modifier.HasActiveSkillUsePermissionModifier(fieldCard) {
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
	if e.hasEffectiveStatus(skill, StatusPetrify) {
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
	extra := map[string]any{"stat": "power"}
	if len(targets) > 0 {
		extra["spell_target"] = targets[0]
		extra["spell_targets"] = targets
		if unit := e.spellTargetUnit(1-playerID, targets[0]); unit != nil {
			extra["spell_target_unit"] = unit
		}
	}
	power := e.effectiveSkillPowerForPurposeWithData(playerID, skill, skill, skillPurposeAttack, extra)
	for _, boostSkill := range boostSkills {
		power += e.effectiveSkillPowerForPurposeWithData(playerID, boostSkill, skill, skillPurposeAttackBoost, extra)
	}
	return max(power, 0)
}

func (e *Engine) effectiveSpellDamage(playerID int, skill *CardInstance, baseDamage int, boostSkills []*CardInstance, affectedUnitGroups ...[]*CardInstance) int {
	damage := baseDamage + e.genericSpellBonus(playerID, skill, "攻")
	extra := map[string]any{"stat": "damage"}
	if e.State.PendingSpell != nil && e.State.PendingSpell.Skill == skill {
		extra["final_power"] = e.State.PendingSpell.TotalPower
	} else {
		extra["final_power"] = e.effectiveSpellPower(playerID, skill, boostSkills)
	}
	if len(affectedUnitGroups) > 0 {
		extra["affected_units"] = affectedUnitGroups[0]
	}
	damage += e.skillContributionStatsWithData(playerID, skill, skill, skillPurposeAttack, extra).DamageBonus
	for _, boostSkill := range boostSkills {
		damage += e.skillContributionStatsWithData(playerID, boostSkill, skill, skillPurposeAttackBoost, extra).DamageBonus
	}
	damage += e.spellStatBonusesWithData(playerID, skill, skillPurposeAttack, extra).DamageBonus
	damage += e.temporarySpellDamageBonus(playerID, skill)
	return max(damage, 0)
}

func (e *Engine) spellPowerSources(playerID int, skill *CardInstance, boostSkills []*CardInstance, totalPower int, targets ...SpellTarget) []SpellPowerSource {
	extra := map[string]any{"stat": "power"}
	if len(targets) > 0 {
		extra["spell_target"] = targets[0]
		extra["spell_targets"] = targets
		if unit := e.spellTargetUnit(1-playerID, targets[0]); unit != nil {
			extra["spell_target_unit"] = unit
		}
	}

	mainPower := e.effectiveSkillPowerForPurposeWithData(playerID, skill, skill, skillPurposeAttack, extra)
	sources := []SpellPowerSource{spellPowerSourceForCard(skill, max(mainPower, 0), true)}
	sum := sources[0].Power
	for _, boostSkill := range boostSkills {
		boostPower := e.effectiveSkillPowerForPurposeWithData(playerID, boostSkill, skill, skillPurposeAttackBoost, extra)
		source := spellPowerSourceForCard(boostSkill, max(boostPower, 0), false)
		sources = append(sources, source)
		sum += source.Power
	}
	if len(sources) > 0 && totalPower != sum {
		sources[0].Power = max(sources[0].Power+totalPower-sum, 0)
	}
	return sources
}

func spellPowerSourceForCard(card *CardInstance, power int, isMain bool) SpellPowerSource {
	source := SpellPowerSource{Power: power, IsMain: isMain}
	if card != nil {
		source.InstanceID = card.InstanceID
		if card.Card != nil {
			source.CardName = card.Card.Name
		}
	}
	return source
}

func (e *Engine) totalEffectiveSkillPower(playerID int, skills []*CardInstance, purpose skillPurpose) int {
	total := 0
	for _, skill := range skills {
		total += e.effectiveSkillPowerForPurposeWithData(playerID, skill, nil, purpose, map[string]any{"stat": "power"})
	}
	return max(total, 0)
}

func (e *Engine) skillContributionStats(playerID int, skill *CardInstance, target *CardInstance, purpose skillPurpose) SpellStats {
	return e.skillContributionStatsWithData(playerID, skill, target, purpose, nil)
}

func (e *Engine) skillContributionStatsWithData(playerID int, skill *CardInstance, target *CardInstance, purpose skillPurpose, extra map[string]any) SpellStats {
	stats := SpellStats{PowerBonus: max(skill.Card.Power+skill.PowerBonus, 0)}
	behavior := globalRegistry.GetBehavior(skill.Card.Number)
	if modifier, ok := behavior.(SkillContributionModifier); ok && modifier.HasActiveSkillContributionModifier(skill) {
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
	return stats
}

func (e *Engine) effectiveSkillPowerForPurposeWithData(playerID int, skill *CardInstance, target *CardInstance, purpose skillPurpose, extra map[string]any) int {
	if skill == nil || skill.Card == nil {
		return 0
	}
	power := e.skillContributionStatsWithData(playerID, skill, target, purpose, extra).PowerBonus
	power += e.spellStatBonusesWithData(playerID, skill, purpose, extra).PowerBonus
	power += e.genericSpellBonus(playerID, skill, "威")
	power += e.temporarySpellPowerBonusForPurpose(playerID, skill, purpose)
	if weak := skill.Statuses[StatusWeaken]; weak > 0 && e.hasEffectiveStatus(skill, StatusWeaken) {
		power -= weak
	}
	if hasActiveSpiritCandle(e.State.Players[playerID]) && !hasCardTag(skill.Card, "灵媒") && !hasCardTag(skill.Card, "神秘") {
		power = (power + 1) / 2
	}
	return max(power, 0)
}

func (e *Engine) spellTargetUnit(defenderID int, target SpellTarget) *CardInstance {
	if target.Type != "unit" || !target.Position.Valid() {
		return nil
	}
	return e.State.Players[defenderID].Units[target.Position.Col][target.Position.Row]
}

func (e *Engine) spellTargetUnitForCaster(playerID int, target SpellTarget) *CardInstance {
	if target.Type != "unit" || !target.Position.Valid() {
		return nil
	}
	targetOwnerID := 1 - playerID
	if target.OwnerID != nil {
		targetOwnerID = *target.OwnerID
	}
	if targetOwnerID < 0 || targetOwnerID >= len(e.State.Players) {
		return nil
	}
	return e.State.Players[targetOwnerID].Units[target.Position.Col][target.Position.Row]
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
	if _, ok := data["cast_player"]; !ok {
		data["cast_player"] = playerID
	}
	ctx := &EffectContext{
		Engine:     e,
		Target:     skill,
		PlayerID:   playerID,
		OpponentID: 1 - playerID,
		ExtraData:  data,
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		modifier, ok := behavior.(SpellStatModifier)
		if !ok || !modifier.HasActiveSpellStatModifier(fieldCard) {
			continue
		}
		ctx.Source = fieldCard
		modifier.ModifySpellStats(ctx, &stats)
	}
	opponentID := 1 - playerID
	opponent := e.State.Players[opponentID]
	enemyCtx := &EffectContext{
		Engine:     e,
		Target:     skill,
		PlayerID:   opponentID,
		OpponentID: playerID,
		ExtraData:  data,
	}
	for _, fieldCard := range e.getAllFieldCards(opponent) {
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		modifier, ok := behavior.(EnemySpellStatModifier)
		if !ok || !modifier.HasActiveEnemySpellStatModifier(fieldCard) {
			continue
		}
		enemyCtx.Source = fieldCard
		modifier.ModifyEnemySpellStats(enemyCtx, &stats)
	}
	for _, owner := range []*PlayerState{ps, opponent} {
		for _, modifier := range owner.TempModifiers {
			if modifier.Type == TempModAllSpellDamageZero && modifier.RemainingUses != 0 && data["stat"] == "damage" {
				if modifier.TargetInstanceID != "" && (skill == nil || modifier.TargetInstanceID != skill.InstanceID) {
					continue
				}
				stats.DamageBonus -= 99
			}
		}
	}
	for _, modifier := range opponent.TempModifiers {
		if modifier.Type == TempModFriendlySpellDamageMinus && modifier.RemainingUses != 0 && data["stat"] == "damage" {
			if modifier.TargetInstanceID != "" && modifier.TargetInstanceID != skill.InstanceID {
				continue
			}
			stats.DamageBonus -= modifier.Amount
		}
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
	return e.validateSpellTargetWithPierce(playerID, skill, target, e.skillHasPierce(playerID, skill))
}

func (e *Engine) validateSpellTargetWithPierce(playerID int, skill *CardInstance, target SpellTarget, hasPierce bool) error {
	if !skillNeedsTargetInstance(skill) {
		if target.Type == "" || target.Type == "none" {
			return nil
		}
	}
	if target.Type == "hero" {
		if friendly, ok := behaviorForNumber(skill.Card.Number).(FriendlySpellTargetBehavior); ok && friendly.HasActiveFriendlySpellTarget(skill) && friendly.AllowsFriendlySpellTarget() {
			return nil
		}
		return fmt.Errorf("spell cannot target hero")
	}
	if target.Type != "unit" {
		return fmt.Errorf("unsupported spell target type: %s", target.Type)
	}
	if !target.Position.Valid() {
		return fmt.Errorf("invalid target position")
	}

	targetOwnerID := 1 - playerID
	if target.OwnerID != nil {
		targetOwnerID = *target.OwnerID
	}
	if targetOwnerID != playerID && targetOwnerID != 1-playerID {
		return fmt.Errorf("invalid spell target owner")
	}
	if targetOwnerID == playerID {
		targetUnit := e.State.Players[playerID].Units[target.Position.Col][target.Position.Row]
		if targetUnit != nil {
			if err := e.validateCardSpecificSpellTarget(playerID, skill, target, targetUnit); err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("no friendly unit at target position")
	}

	opponent := e.State.Players[1-playerID]
	if opponent.Units[target.Position.Col][target.Position.Row] == nil {
		if e.effectiveSpellArea(skill) == SpellAreaColumn {
			for row := 0; row < 3; row++ {
				if opponent.Units[target.Position.Col][row] == nil {
					continue
				}
				frontRow := opponent.GetFrontRow()
				if e.IsInSpellRange(playerID, target.Position.Col, target.Position.Row, hasPierce) || frontRow == -1 || target.Position.Row == frontRow {
					return nil
				}
			}
		}
		if friendly, ok := behaviorForNumber(skill.Card.Number).(FriendlySpellTargetBehavior); ok && friendly.HasActiveFriendlySpellTarget(skill) && friendly.AllowsFriendlySpellTarget() {
			own := e.State.Players[playerID]
			if own.Units[target.Position.Col][target.Position.Row] != nil {
				return nil
			}
		}
		return fmt.Errorf("no enemy unit at target position")
	}
	targetUnit := opponent.Units[target.Position.Col][target.Position.Row]
	if err := e.validateCardSpecificSpellTarget(playerID, skill, target, targetUnit); err != nil {
		return err
	}
	if sinMatchesTargetTag(skill, targetUnit) {
		hasPierce = true
	}

	if e.effectiveSpellArea(skill) == SpellAreaFrontRow {
		frontRow := opponent.GetFrontRow()
		if frontRow >= 0 && target.Position.Row != frontRow {
			return fmt.Errorf("front-row spell target must be in the actual front row")
		}
	}

	if e.hasStealthFromOpponent(playerID, targetUnit) && e.spellAllowsStealthTarget(skill) {
		return nil
	}
	if e.fireCloudFanGrantsTarget(playerID, skill, target) {
		return nil
	}
	if !e.IsInSpellRange(playerID, target.Position.Col, target.Position.Row, hasPierce) {
		return fmt.Errorf("target is not in spell range")
	}

	return nil
}

func (e *Engine) fireCloudFanGrantsTarget(playerID int, skill *CardInstance, target SpellTarget) bool {
	if e == nil || skill == nil || skill.Card == nil || target.Type != "unit" || !target.Position.Valid() {
		return false
	}
	if skill.Card.Category != model.ElementFire && skill.Card.Category != model.ElementAir {
		return false
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return false
	}
	hasFan := false
	for _, card := range e.getAllFieldCards(ps) {
		if card != nil && card.Card != nil && card.Card.Number == "2121102" && !e.hasEffectiveStatus(card, StatusPetrify) {
			hasFan = true
			break
		}
	}
	if !hasFan {
		return false
	}
	opponent := e.State.Players[1-playerID]
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 || target.Position.Row <= frontRow {
		return false
	}
	frontOfTarget := target.Position.Row - 1
	return frontOfTarget >= 0 && opponent.Units[target.Position.Col][frontOfTarget] == nil
}

func (e *Engine) validateCardSpecificSpellTarget(playerID int, skill *CardInstance, target SpellTarget, targetUnit *CardInstance) error {
	if skill == nil || skill.Card == nil {
		return nil
	}
	validator, ok := behaviorForNumber(skill.Card.Number).(SpellTargetValidationBehavior)
	if !ok || !validator.HasActiveSpellTargeting(skill) {
		return nil
	}
	return validator.ValidateSpellTarget(&EffectContext{
		Engine:     e,
		Source:     skill,
		Target:     targetUnit,
		PlayerID:   playerID,
		OpponentID: 1 - playerID,
	}, target, targetUnit)
}

func (e *Engine) spellAllowsStealthTarget(skill *CardInstance) bool {
	if skill == nil || skill.Card == nil {
		return false
	}
	behavior, ok := behaviorForNumber(skill.Card.Number).(StealthSpellTargetBehavior)
	return ok && behavior.HasActiveSpellTargeting(skill) && behavior.AllowsStealthSpellTarget()
}

func (e *Engine) validateSpellExtraTarget(playerID int, target SpellTarget) error {
	if target.Type != "unit" {
		return fmt.Errorf("unsupported spell target type: %s", target.Type)
	}
	if !target.Position.Valid() {
		return fmt.Errorf("invalid target position")
	}
	if e.State.Players[1-playerID].Units[target.Position.Col][target.Position.Row] == nil {
		return fmt.Errorf("no enemy unit at extra target position")
	}
	return nil
}

func (e *Engine) validateSpellExtraTargetForSkill(playerID int, skill *CardInstance, mainTarget SpellTarget, extra SpellTarget) error {
	if extra.Type == mainTarget.Type && extra.Position == mainTarget.Position {
		if e.allowsSameSpellExtraTarget(e.State.Players[playerID], skill) {
			return e.validateSpellTarget(playerID, skill, extra)
		}
		return fmt.Errorf("extra target cannot be the same as the main target")
	}
	return e.validateSpellExtraTarget(playerID, extra)
}

func (e *Engine) burrowExtraTargetsFromAction(playerID int, skill *CardInstance, mainTarget SpellTarget, action ActionMessage) ([]SpellTarget, error) {
	if skill == nil || skill.Card == nil || skill.Card.Number != "3421107" {
		return nil, nil
	}
	raw, ok := action.Data["extra_targets"].([]any)
	if !ok || len(raw) == 0 {
		return nil, nil
	}
	maxTargets := skill.Statuses[StatusMastery]
	if maxTargets > 2 {
		maxTargets = 2
	}
	if maxTargets <= 0 {
		return nil, fmt.Errorf("burrow extra targets require mastery")
	}
	if len(raw) > maxTargets {
		return nil, fmt.Errorf("too many extra targets for burrow")
	}
	result := make([]SpellTarget, 0, len(raw))
	seen := map[Position]bool{mainTarget.Position: true}
	for _, entry := range raw {
		data, ok := entry.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid burrow extra target")
		}
		pos, err := requiredBoardPosition(data, "col", "row")
		if err != nil {
			return nil, fmt.Errorf("invalid burrow extra target")
		}
		extra := SpellTarget{Type: "unit", Position: pos}
		if err := e.validateSpellExtraTarget(playerID, extra); err != nil {
			return nil, err
		}
		if seen[extra.Position] {
			return nil, fmt.Errorf("duplicate burrow extra target")
		}
		seen[extra.Position] = true
		result = append(result, extra)
	}
	return result, nil
}

func (e *Engine) allowsSameSpellExtraTarget(ps *PlayerState, skill *CardInstance) bool {
	if ps == nil || skill == nil || skill.Card == nil {
		return false
	}
	if skill.Card.Number == "3621107" {
		return true
	}
	for _, modifier := range ps.TempModifiers {
		if modifier.Type != TempModNextSpellExtraTarget || modifier.RemainingUses == 0 {
			continue
		}
		if (modifier.TargetInstanceID == "" || modifier.TargetInstanceID == skill.InstanceID) && modifier.AllowSameTarget {
			return true
		}
	}
	return false
}

func (e *Engine) spellHasPierceWithBoosts(playerID int, skill *CardInstance, boostSkills []*CardInstance) bool {
	if e.skillHasPierce(playerID, skill) {
		return true
	}
	for _, boostSkill := range boostSkills {
		if e.skillContributionStats(playerID, boostSkill, skill, skillPurposeAttackBoost).Pierce {
			return true
		}
	}
	return false
}

func (e *Engine) skillHasPierce(playerID int, skill *CardInstance) bool {
	if skill != nil && skill.Statuses[permanentPierceStatus] > 0 {
		return true
	}
	if cardHasPierce(skill) || e.windBladeGrantsPierce(playerID, skill) || e.westernChartGrantsPierce(playerID, skill) {
		return true
	}
	return skill != nil && skill.Card != nil && skill.Card.Number == "3621107" && e.redMoonActive(playerID)
}

func (e *Engine) westernChartGrantsPierce(playerID int, skill *CardInstance) bool {
	if skill == nil || skill.Card == nil || skill.InstanceID == "" || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	key := westernChartPierceTargetPrefix + skill.InstanceID
	for _, equipment := range e.State.Players[playerID].Equipment {
		if equipment == nil || equipment.Card == nil || equipment.Card.Number != "2221108" || e.hasEffectiveStatus(equipment, StatusPetrify) {
			continue
		}
		if equipment.Statuses[key] > 0 {
			return true
		}
	}
	return false
}

func (e *Engine) redMoonActive(playerID int) bool {
	if playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	for _, card := range e.getAllFieldCards(e.State.Players[playerID]) {
		if card != nil && card.Card != nil && card.Card.Number == "3611101" && abilityDurationActive(card) && !e.hasEffectiveStatus(card, StatusPetrify) {
			return true
		}
	}
	return false
}

func (e *Engine) windBladeGrantsPierce(playerID int, skill *CardInstance) bool {
	if skill == nil || skill.Card == nil || skill.Card.Category != model.ElementAir || !skillNeedsTargetInstance(skill) {
		return false
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return false
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard != nil && fieldCard.Card != nil && fieldCard.Card.Number == "1311003" && !e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			return true
		}
	}
	return false
}
