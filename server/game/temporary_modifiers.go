package game

import (
	"fmt"
	"strconv"
	"strings"

	"eraofarcane/model"
)

const (
	TempModNextSkillCostZero            = "next_skill_cost_zero"
	TempModNextSkillUseCostMinus        = "next_skill_use_cost_minus"
	TempModNextItemOrSkillCostMinus     = "next_item_or_skill_cost_minus"
	TempModNextFireCardPlayCostMinus    = "next_fire_card_play_cost_minus"
	TempModCurrentTurnSkillUseCostMinus = "current_turn_skill_use_cost_minus"
	TempModCurrentTurnSkillCostZero     = "current_turn_skill_cost_zero"
	TempModNextLearnedSkillHaste        = "next_learned_skill_haste"
	TempModSkillPowerBonus              = "skill_power_bonus"
	TempModNextAttackSpellPowerBonus    = "next_attack_spell_power_bonus"
	TempModSkillAttackBonus             = "skill_attack_bonus"
	TempModNextSkillUseAttackBonus      = "next_skill_use_attack_bonus"
	TempModNextNoCooldown               = "next_skill_no_cooldown"
	TempModNextSpellHitStatus           = "next_spell_hit_status"
	TempModNextElementSpellPowerBonus   = "next_element_spell_power_bonus"
	TempModNextTaggedSpellPowerBonus    = "next_tagged_spell_power_bonus"
	TempModNextElementSpellDamageBonus  = "next_element_spell_damage_bonus"
	TempModCurrentTurnElementDamage     = "current_turn_element_damage"
	TempModCurrentTurnElementHitStatus  = "current_turn_element_hit_status"
	TempModDelayedElementGain           = "delayed_element_gain"
	TempModDelayedShieldGain            = "delayed_shield_gain"
	TempModResetSkillsOnOpponentTurnEnd = "reset_skills_on_opponent_turn_end"
	TempModLampusSwordDelayedDamage     = "lampus_sword_delayed_damage"
	TempModNextEarthSkillLearnCostMinus = "next_earth_skill_learn_cost_minus"
	TempModAllSpellDamageZero           = "all_spell_attack_zero"
	TempModFriendlySpellDamageMinus     = "friendly_spell_damage_minus"
	TempModSkillUseCooldownAdd          = "skill_use_cooldown_add"
	TempModFriendlyNegativeStatusIgnore = "friendly_negative_status_ignore"
	TempModPainScreamWeakenOnDamage     = "pain_scream_weaken_on_damage"
	TempModNextDriveSpellExtraTarget    = "next_drive_spell_extra_target"
	TempModNextSpellExtraTarget         = "next_spell_extra_target"
	TempModCannotLearnElementSkill      = "cannot_learn_element_skill"
)

const skillUseExtraCostStatusPrefix = "使用费用额外"

type TemporaryModifier struct {
	ID               string `json:"id"`
	Type             string `json:"type"`
	SourceCardNumber string `json:"source_card_number,omitempty"`
	SourceName       string `json:"source_name,omitempty"`
	TargetInstanceID string `json:"target_instance_id,omitempty"`
	Element          string `json:"element,omitempty"`
	Status           string `json:"status,omitempty"`
	Amount           int    `json:"amount,omitempty"`
	RemainingUses    int    `json:"remaining_uses,omitempty"`
	ExpiresTurn      int    `json:"expires_turn,omitempty"`
	AllowSameTarget  bool   `json:"allow_same_target,omitempty"`
}

func (e *Engine) addTemporaryModifier(playerID int, modifier TemporaryModifier) {
	if modifier.ID == "" {
		modifier.ID = generateID()
	}
	e.State.Players[playerID].TempModifiers = append(e.State.Players[playerID].TempModifiers, modifier)
	e.emit(GameEvent{
		Type:   "effect_trigger",
		Player: playerID,
		Data: map[string]any{
			"effect":   "temp_modifier",
			"modifier": modifier,
		},
	})
}

func (e *Engine) removeTemporaryModifier(playerID int, modifierID string) {
	ps := e.State.Players[playerID]
	kept := ps.TempModifiers[:0]
	for _, modifier := range ps.TempModifiers {
		if modifier.ID != modifierID {
			kept = append(kept, modifier)
		}
	}
	ps.TempModifiers = kept
}

func (e *Engine) clearExpiredTemporaryModifiers(playerID int) {
	ps := e.State.Players[playerID]
	kept := ps.TempModifiers[:0]
	for _, modifier := range ps.TempModifiers {
		if modifier.ExpiresTurn > 0 && modifier.ExpiresTurn <= e.State.TurnNumber {
			continue
		}
		kept = append(kept, modifier)
	}
	ps.TempModifiers = kept
}

func (e *Engine) nextSkillCostZeroModifier(ps *PlayerState, skill *CardInstance) *TemporaryModifier {
	for i := range ps.TempModifiers {
		modifier := &ps.TempModifiers[i]
		if modifier.Type != TempModNextSkillCostZero && modifier.Type != TempModCurrentTurnSkillCostZero {
			continue
		}
		if modifier.RemainingUses == 0 {
			continue
		}
		if modifier.TargetInstanceID == "" || modifier.TargetInstanceID == skill.InstanceID {
			return modifier
		}
	}
	return nil
}

func (e *Engine) temporaryNextSkillUseCostMinus(ps *PlayerState, skill *CardInstance, elem string) int {
	if ps == nil || skill == nil || elem == "" {
		return 0
	}
	total := 0
	for _, modifier := range ps.TempModifiers {
		if modifier.Type != TempModNextSkillUseCostMinus && modifier.Type != TempModCurrentTurnSkillUseCostMinus && modifier.Type != TempModNextItemOrSkillCostMinus {
			continue
		}
		if (modifier.Type == TempModNextSkillUseCostMinus || modifier.Type == TempModNextItemOrSkillCostMinus) && modifier.RemainingUses == 0 {
			continue
		}
		if modifier.Element != "" && modifier.Element != elem {
			continue
		}
		if modifier.TargetInstanceID != "" && modifier.TargetInstanceID != skill.InstanceID {
			continue
		}
		amount := modifier.Amount
		if amount <= 0 {
			amount = 1
		}
		total += amount
	}
	return total
}

func (e *Engine) temporaryNextCardPlayCostMinus(ps *PlayerState, card *CardInstance, elem string) int {
	if ps == nil || card == nil || card.Card == nil || elem == "" {
		return 0
	}
	total := 0
	for _, modifier := range ps.TempModifiers {
		switch modifier.Type {
		case TempModNextItemOrSkillCostMinus:
			if !isConsumableOrSkillCardInstance(card) {
				continue
			}
		case TempModNextFireCardPlayCostMinus:
			if card.Card.Category != model.ElementFire {
				continue
			}
		default:
			continue
		}
		if modifier.RemainingUses == 0 || modifier.Element != "" && modifier.Element != elem || modifier.TargetInstanceID != "" && modifier.TargetInstanceID != card.InstanceID {
			continue
		}
		amount := modifier.Amount
		if amount <= 0 {
			amount = 1
		}
		total += amount
	}
	return total
}

func (e *Engine) consumeNextSkillUseModifiers(ps *PlayerState, skill *CardInstance) {
	e.consumeNextSkillUseModifiersForPurpose(ps, skill, skillPurposeAttack)
}

func (e *Engine) consumeNextSkillUseModifiersForPurpose(ps *PlayerState, skill *CardInstance, purpose skillPurpose) {
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.RemainingUses == 0 {
			continue
		}
		if modifier.TargetInstanceID != "" && modifier.TargetInstanceID != skill.InstanceID {
			continue
		}
		if isBoostPurpose(purpose) && modifier.Type != TempModNextSkillUseCostMinus && modifier.Type != TempModNextItemOrSkillCostMinus {
			continue
		}
		switch modifier.Type {
		case TempModNextSkillCostZero, TempModNextSkillUseCostMinus, TempModNextItemOrSkillCostMinus, TempModNextNoCooldown:
			modifier.RemainingUses--
			if modifier.RemainingUses <= 0 {
				e.removeTemporaryModifier(ps.PlayerID, modifier.ID)
			}
		}
	}
	consumeSkillUseExtraCostStatuses(skill)
}

func skillUseExtraCostStatus(elem string, amount int) string {
	if amount <= 1 {
		return skillUseExtraCostStatusPrefix + elem
	}
	return fmt.Sprintf("%s%s%d", skillUseExtraCostStatusPrefix, elem, amount)
}

func consumeSkillUseExtraCostStatuses(skill *CardInstance) {
	if skill == nil {
		return
	}
	for status, amount := range skill.Statuses {
		if amount <= 0 || !strings.HasPrefix(status, skillUseExtraCostStatusPrefix) {
			continue
		}
		if amount == 1 {
			delete(skill.Statuses, status)
			continue
		}
		skill.Statuses[status] = amount - 1
	}
}

func (e *Engine) consumeNextCardPlayCostModifiers(ps *PlayerState, card *CardInstance) {
	if ps == nil || card == nil || card.Card == nil {
		return
	}
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.RemainingUses == 0 || modifier.TargetInstanceID != "" && modifier.TargetInstanceID != card.InstanceID {
			continue
		}
		switch modifier.Type {
		case TempModNextItemOrSkillCostMinus:
			if !isConsumableOrSkillCardInstance(card) {
				continue
			}
		case TempModNextFireCardPlayCostMinus:
			if card.Card.Category != model.ElementFire {
				continue
			}
		default:
			continue
		}
		modifier.RemainingUses--
		if modifier.RemainingUses <= 0 {
			e.removeTemporaryModifier(ps.PlayerID, modifier.ID)
		}
	}
}

func isConsumableOrSkillCardInstance(card *CardInstance) bool {
	return isConsumableCardInstance(card) || (card != nil && card.Card != nil && card.Card.IsSkill())
}

func (e *Engine) applySkillUseCooldownModifiers(ps *PlayerState, skills ...*CardInstance) {
	for _, skill := range skills {
		if skill == nil || skill.Card == nil || !skill.Card.IsSkill() {
			continue
		}
		for i := range ps.TempModifiers {
			modifier := &ps.TempModifiers[i]
			if modifier.Type != TempModSkillUseCooldownAdd {
				continue
			}
			if modifier.TargetInstanceID != "" && modifier.TargetInstanceID != skill.InstanceID {
				continue
			}
			amount := modifier.Amount
			if amount <= 0 {
				amount = 1
			}
			skill.Statuses[StatusCooldown] += amount
		}
	}
}

func (e *Engine) shouldSkipCooldown(ps *PlayerState, skill *CardInstance) bool {
	for i := range ps.TempModifiers {
		modifier := &ps.TempModifiers[i]
		if modifier.Type != TempModNextNoCooldown {
			continue
		}
		if modifier.RemainingUses == 0 {
			continue
		}
		if modifier.TargetInstanceID == "" || modifier.TargetInstanceID == skill.InstanceID {
			return true
		}
	}
	return false
}

func (e *Engine) hasNextDriveSpellExtraTarget(ps *PlayerState, skill *CardInstance) bool {
	if ps == nil || skill == nil || skill.Card == nil {
		return false
	}
	for _, modifier := range ps.TempModifiers {
		if modifier.Type == TempModNextDriveSpellExtraTarget && modifier.RemainingUses != 0 && hasCardTag(skill.Card, "驱动") {
			return true
		}
		if modifier.Type == TempModNextSpellExtraTarget && modifier.RemainingUses != 0 &&
			(modifier.TargetInstanceID == "" || modifier.TargetInstanceID == skill.InstanceID) {
			return true
		}
	}
	return false
}

func (e *Engine) hasNextSpellExtraTarget(ps *PlayerState, skill *CardInstance) bool {
	if ps == nil || skill == nil || skill.Card == nil {
		return false
	}
	for _, modifier := range ps.TempModifiers {
		if modifier.Type != TempModNextSpellExtraTarget || modifier.RemainingUses == 0 {
			continue
		}
		if modifier.TargetInstanceID == "" || modifier.TargetInstanceID == skill.InstanceID {
			return true
		}
	}
	return false
}

func (e *Engine) consumeNextDriveSpellExtraTarget(ps *PlayerState, skill *CardInstance) {
	if ps == nil || skill == nil || skill.Card == nil {
		return
	}
	for i := range ps.TempModifiers {
		modifier := &ps.TempModifiers[i]
		switch modifier.Type {
		case TempModNextDriveSpellExtraTarget:
			if !hasCardTag(skill.Card, "驱动") || modifier.RemainingUses == 0 {
				continue
			}
		case TempModNextSpellExtraTarget:
			if modifier.RemainingUses == 0 || (modifier.TargetInstanceID != "" && modifier.TargetInstanceID != skill.InstanceID) {
				continue
			}
		default:
			continue
		}
		modifier.RemainingUses--
		return
	}
}

func (e *Engine) temporarySpellPowerBonus(playerID int, skill *CardInstance) int {
	return e.temporarySpellPowerBonusForPurpose(playerID, skill, skillPurposeAttack)
}

func (e *Engine) temporarySpellPowerBonusForPurpose(playerID int, skill *CardInstance, purpose skillPurpose) int {
	total := 0
	for _, modifier := range e.State.Players[playerID].TempModifiers {
		switch modifier.Type {
		case TempModSkillPowerBonus:
			if modifier.Element != "" && modifier.Element != skill.Card.Category {
				continue
			}
			if modifier.TargetInstanceID == "" || modifier.TargetInstanceID == skill.InstanceID {
				total += modifier.Amount
			}
		case TempModNextAttackSpellPowerBonus:
			if purpose == skillPurposeAttack && modifier.RemainingUses != 0 && (modifier.TargetInstanceID == "" || modifier.TargetInstanceID == skill.InstanceID) {
				total += modifier.Amount
			}
		case TempModNextElementSpellPowerBonus:
			if modifier.RemainingUses != 0 && modifier.Status == skill.Card.Category {
				total += modifier.Amount
			}
		case TempModNextTaggedSpellPowerBonus:
			if modifier.RemainingUses != 0 && hasCardTag(skill.Card, modifier.Status) {
				total += modifier.Amount
			}
		}
	}
	return total
}

func (e *Engine) consumeNextSpellPowerBonuses(ps *PlayerState, skill *CardInstance) {
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		switch modifier.Type {
		case TempModSkillPowerBonus:
			if modifier.RemainingUses == 0 {
				continue
			}
			if modifier.Element != "" && modifier.Element != skill.Card.Category {
				continue
			}
			if modifier.TargetInstanceID != "" && modifier.TargetInstanceID != skill.InstanceID {
				continue
			}
		case TempModNextAttackSpellPowerBonus:
			if modifier.RemainingUses == 0 {
				continue
			}
			if modifier.TargetInstanceID != "" && modifier.TargetInstanceID != skill.InstanceID {
				continue
			}
		case TempModNextElementSpellPowerBonus:
			if modifier.RemainingUses == 0 || modifier.Status != skill.Card.Category {
				continue
			}
		case TempModNextTaggedSpellPowerBonus:
			if modifier.RemainingUses == 0 || !hasCardTag(skill.Card, modifier.Status) {
				continue
			}
		default:
			continue
		}
		modifier.RemainingUses--
		if modifier.RemainingUses <= 0 {
			e.removeTemporaryModifier(ps.PlayerID, modifier.ID)
		}
	}
}

func (e *Engine) temporarySpellDamageBonus(playerID int, skill *CardInstance) int {
	total := 0
	for _, modifier := range e.State.Players[playerID].TempModifiers {
		if modifier.Type == TempModSkillAttackBonus {
			if modifier.RemainingUses != 0 && (modifier.TargetInstanceID == "" || modifier.TargetInstanceID == skill.InstanceID) {
				total += modifier.Amount
			}
			continue
		}
		if modifier.Type == TempModNextSkillUseAttackBonus {
			if modifier.RemainingUses != 0 && (modifier.TargetInstanceID == "" || modifier.TargetInstanceID == skill.InstanceID) {
				total += modifier.Amount
			}
			continue
		}
		if modifier.Type != TempModNextElementSpellDamageBonus {
			if modifier.Type != TempModCurrentTurnElementDamage {
				continue
			}
		}
		if modifier.Type == TempModNextElementSpellDamageBonus && modifier.Status != skill.Card.Category {
			continue
		}
		if modifier.Type == TempModCurrentTurnElementDamage && modifier.Element != skill.Card.Category {
			continue
		}
		if modifier.Type == TempModNextElementSpellDamageBonus && modifier.RemainingUses == 0 {
			continue
		}
		total += modifier.Amount
	}
	return total
}

func (e *Engine) consumeNextSpellAttackBonuses(ps *PlayerState, skill *CardInstance) {
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.Type != TempModSkillAttackBonus && modifier.Type != TempModNextSkillUseAttackBonus {
			continue
		}
		if modifier.RemainingUses == 0 {
			continue
		}
		if modifier.TargetInstanceID != "" && modifier.TargetInstanceID != skill.InstanceID {
			continue
		}
		modifier.RemainingUses--
		if modifier.RemainingUses <= 0 {
			e.removeTemporaryModifier(ps.PlayerID, modifier.ID)
		}
	}
}

func (e *Engine) consumeNextElementSpellDamageBonus(ps *PlayerState, skill *CardInstance) {
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.Type != TempModNextElementSpellDamageBonus {
			continue
		}
		if modifier.RemainingUses == 0 || modifier.Status != skill.Card.Category {
			continue
		}
		modifier.RemainingUses--
		if modifier.RemainingUses <= 0 {
			e.removeTemporaryModifier(ps.PlayerID, modifier.ID)
		}
	}
}

func (e *Engine) consumeFriendlySpellDamageMinus(ps *PlayerState, skill *CardInstance) {
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.Type != TempModFriendlySpellDamageMinus {
			continue
		}
		if modifier.RemainingUses == 0 {
			continue
		}
		if modifier.TargetInstanceID != "" && (skill == nil || modifier.TargetInstanceID != skill.InstanceID) {
			continue
		}
		modifier.RemainingUses--
		if modifier.RemainingUses <= 0 {
			e.removeTemporaryModifier(ps.PlayerID, modifier.ID)
		}
	}
}

func (e *Engine) consumeAllSpellDamageZero(ps *PlayerState, skill *CardInstance) {
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.Type != TempModAllSpellDamageZero {
			continue
		}
		if modifier.RemainingUses == 0 {
			continue
		}
		if modifier.TargetInstanceID != "" && (skill == nil || modifier.TargetInstanceID != skill.InstanceID) {
			continue
		}
		modifier.RemainingUses--
		if modifier.RemainingUses <= 0 {
			e.removeTemporaryModifier(ps.PlayerID, modifier.ID)
		}
	}
}

func (e *Engine) addNextElementSpellPowerBonus(playerID int, elem string, amount int) {
	if amount <= 0 {
		return
	}
	e.addTemporaryModifier(playerID, TemporaryModifier{
		Type:          TempModNextElementSpellPowerBonus,
		Status:        elem,
		Amount:        amount,
		RemainingUses: 1,
		ExpiresTurn:   e.State.TurnNumber + 2,
	})
}

func (e *Engine) addNextTaggedSpellPowerBonus(playerID int, tag string, amount int) {
	if amount <= 0 || tag == "" {
		return
	}
	e.addTemporaryModifier(playerID, TemporaryModifier{
		Type:          TempModNextTaggedSpellPowerBonus,
		Status:        tag,
		Amount:        amount,
		RemainingUses: 1,
		ExpiresTurn:   e.State.TurnNumber + 2,
	})
}

func (e *Engine) addNextElementSpellDamageBonus(playerID int, elem string, amount int) {
	if amount <= 0 {
		return
	}
	e.addTemporaryModifier(playerID, TemporaryModifier{
		Type:          TempModNextElementSpellDamageBonus,
		Status:        elem,
		Amount:        amount,
		RemainingUses: 1,
		ExpiresTurn:   e.State.TurnNumber + 2,
	})
}

func (e *Engine) pureArcaneChoices(playerID int) []map[string]any {
	ps := e.State.Players[playerID]
	choices := make([]map[string]any, 0)
	for _, elem := range []string{model.ElementFire, model.ElementWater, model.ElementAir, model.ElementEarth, model.ElementLight, model.ElementShadow, model.ElementArcane} {
		available := ps.Elements[elem]
		if available <= 0 {
			continue
		}
		for amount := 1; amount <= min(available, 10); amount++ {
			id := elem + ":" + fmt.Sprintf("%d", amount)
			choices = append(choices, map[string]any{
				"instance_id": id,
				"number":      "3001002",
				"name":        elem + " " + fmt.Sprintf("%d", amount),
				"type":        "元素选择",
				"zone":        "choice",
				"side":        "own",
				"element":     elem,
				"amount":      amount,
			})
		}
	}
	return choices
}

func parsePureArcaneChoice(choice string) (string, int, bool) {
	parts := strings.Split(choice, ":")
	if len(parts) != 2 {
		return "", 0, false
	}
	amount, err := strconv.Atoi(parts[1])
	if err != nil || amount <= 0 || amount > 10 {
		return "", 0, false
	}
	return parts[0], amount, true
}

func (e *Engine) nextEarthSkillLearnCostMinus(ps *PlayerState, skill *CardInstance) *TemporaryModifier {
	if skill == nil || skill.Card == nil || !skill.Card.IsSkill() || skill.Card.Category != model.ElementEarth {
		return nil
	}
	for i := range ps.TempModifiers {
		modifier := &ps.TempModifiers[i]
		if modifier.Type != TempModNextEarthSkillLearnCostMinus || modifier.RemainingUses == 0 {
			continue
		}
		return modifier
	}
	return nil
}

func (e *Engine) consumeEarthSkillLearnCostModifier(ps *PlayerState, skill *CardInstance) {
	modifier := e.nextEarthSkillLearnCostMinus(ps, skill)
	if modifier == nil {
		return
	}
	modifier.RemainingUses--
	if modifier.RemainingUses <= 0 {
		e.removeTemporaryModifier(ps.PlayerID, modifier.ID)
	}
}

func (e *Engine) applyTemporarySpellHitStatus(playerID int, skill *CardInstance, affectedUnits []*CardInstance) {
	if len(affectedUnits) == 0 {
		return
	}
	ps := e.State.Players[playerID]
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.Type != TempModNextSpellHitStatus && modifier.Type != TempModCurrentTurnElementHitStatus {
			continue
		}
		if modifier.Type == TempModNextSpellHitStatus && modifier.RemainingUses == 0 {
			continue
		}
		if modifier.Type == TempModCurrentTurnElementHitStatus && modifier.Element != skill.Card.Category {
			continue
		}
		if modifier.TargetInstanceID != "" && modifier.TargetInstanceID != skill.InstanceID {
			continue
		}
		amount := modifier.Amount
		if amount <= 0 {
			amount = 1
		}
		for _, unit := range affectedUnits {
			if unit == nil {
				continue
			}
			if !e.addStatus(unit, modifier.Status, amount) {
				continue
			}
			e.emit(GameEvent{
				Type:   "effect_trigger",
				Player: playerID,
				Data: map[string]any{
					"source": cardToInfo(skill),
					"target": cardToInfo(unit),
					"effect": "apply_status",
					"status": modifier.Status,
					"amount": amount,
				},
			})
		}
		if modifier.Type == TempModNextSpellHitStatus {
			modifier.RemainingUses--
			if modifier.RemainingUses <= 0 {
				e.removeTemporaryModifier(playerID, modifier.ID)
			}
		}
	}
}
