package game

import (
	"fmt"
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

func (e *Engine) effectiveSkillUseCostForPurpose(ps *PlayerState, skill *CardInstance, purpose skillPurpose) map[string]int {
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
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		if modifier, ok := behavior.(SkillUseCostModifier); ok && modifier.HasActiveSkillUseCostModifier(fieldCard) {
			ctx.Target = fieldCard
			modifier.ModifySkillUseCost(ctx, cost)
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
		if modifier, ok := behavior.(CardPlayCostModifier); ok && modifier.HasActiveCardPlayCostModifier(fieldCard) {
			ctx.Source = fieldCard
			modifier.ModifyCardPlayCost(ctx, card, cost)
		}
	}
	if card.Statuses["入场费用水-1"] > 0 {
		reduceCost(cost, model.ElementWater, card.Statuses["入场费用水-1"])
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
	if capacity > MaxEquipmentSlots {
		return MaxEquipmentSlots
	}
	return capacity
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
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
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

func (e *Engine) effectiveSpellDamage(playerID int, skill *CardInstance, baseDamage int, boostSkills []*CardInstance) int {
	damage := baseDamage + e.genericSpellBonus(playerID, skill, "攻")
	extra := map[string]any{"stat": "damage"}
	if e.State.PendingSpell != nil && e.State.PendingSpell.Skill == skill {
		extra["final_power"] = e.State.PendingSpell.TotalPower
	} else {
		extra["final_power"] = e.effectiveSpellPower(playerID, skill, boostSkills)
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
	stats.DamageBonus = max(stats.DamageBonus, 0)
	return stats
}

func (e *Engine) effectiveSkillPowerForPurposeWithData(playerID int, skill *CardInstance, target *CardInstance, purpose skillPurpose, extra map[string]any) int {
	if skill == nil || skill.Card == nil {
		return 0
	}
	power := e.skillContributionStatsWithData(playerID, skill, target, purpose, extra).PowerBonus
	power += e.spellStatBonusesWithData(playerID, skill, purpose, extra).PowerBonus
	power += e.genericSpellBonus(playerID, skill, "威")
	power += e.temporarySpellPowerBonus(playerID, skill)
	if weak := skill.Statuses[StatusWeaken]; weak > 0 && e.hasEffectiveStatus(skill, StatusWeaken) {
		power -= weak
	}
	return max(power, 0)
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
	return e.validateSpellTargetWithPierce(playerID, skill, target, cardHasPierce(skill))
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

	if !e.IsInSpellRange(playerID, target.Position.Col, target.Position.Row, hasPierce) {
		return fmt.Errorf("target is not in spell range")
	}

	return nil
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

func (e *Engine) spellHasPierceWithBoosts(playerID int, skill *CardInstance, boostSkills []*CardInstance) bool {
	if cardHasPierce(skill) {
		return true
	}
	for _, boostSkill := range boostSkills {
		if e.skillContributionStats(playerID, boostSkill, skill, skillPurposeAttackBoost).Pierce {
			return true
		}
	}
	return false
}
