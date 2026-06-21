package game

import (
	"fmt"
	"strconv"
	"strings"

	"eraofarcane/model"
)

const (
	TempModNextSkillCostZero            = "next_skill_cost_zero"
	TempModCurrentTurnSkillCostZero     = "current_turn_skill_cost_zero"
	TempModNextLearnedSkillHaste        = "next_learned_skill_haste"
	TempModSkillPowerBonus              = "skill_power_bonus"
	TempModNextNoCooldown               = "next_skill_no_cooldown"
	TempModNextSpellHitStatus           = "next_spell_hit_status"
	TempModNextElementSpellPowerBonus   = "next_element_spell_power_bonus"
	TempModNextElementSpellDamageBonus  = "next_element_spell_damage_bonus"
	TempModDelayedElementGain           = "delayed_element_gain"
	TempModResetSkillsOnOpponentTurnEnd = "reset_skills_on_opponent_turn_end"
	TempModNextEarthSkillLearnCostMinus = "next_earth_skill_learn_cost_minus"
	TempModAllSpellDamageZero           = "all_spell_attack_zero"
	TempModFriendlySpellDamageMinus     = "friendly_spell_damage_minus"
	TempModSkillUseCooldownAdd          = "skill_use_cooldown_add"
)

type TemporaryModifier struct {
	ID               string `json:"id"`
	Type             string `json:"type"`
	SourceCardNumber string `json:"source_card_number,omitempty"`
	SourceName       string `json:"source_name,omitempty"`
	TargetInstanceID string `json:"target_instance_id,omitempty"`
	Status           string `json:"status,omitempty"`
	Amount           int    `json:"amount,omitempty"`
	RemainingUses    int    `json:"remaining_uses,omitempty"`
	ExpiresTurn      int    `json:"expires_turn,omitempty"`
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

func (e *Engine) consumeNextSkillUseModifiers(ps *PlayerState, skill *CardInstance) {
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.RemainingUses == 0 {
			continue
		}
		if modifier.TargetInstanceID != "" && modifier.TargetInstanceID != skill.InstanceID {
			continue
		}
		switch modifier.Type {
		case TempModNextSkillCostZero, TempModNextNoCooldown:
			modifier.RemainingUses--
			if modifier.RemainingUses <= 0 {
				e.removeTemporaryModifier(ps.PlayerID, modifier.ID)
			}
		}
	}
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

func (e *Engine) temporarySpellPowerBonus(playerID int, skill *CardInstance) int {
	total := 0
	for _, modifier := range e.State.Players[playerID].TempModifiers {
		switch modifier.Type {
		case TempModSkillPowerBonus:
			if modifier.TargetInstanceID == "" || modifier.TargetInstanceID == skill.InstanceID {
				total += modifier.Amount
			}
		case TempModNextElementSpellPowerBonus:
			if modifier.RemainingUses != 0 && modifier.Status == skill.Card.Category {
				total += modifier.Amount
			}
		}
	}
	return total
}

func (e *Engine) consumeNextElementSpellPowerBonus(ps *PlayerState, skill *CardInstance) {
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.Type != TempModNextElementSpellPowerBonus {
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

func (e *Engine) temporarySpellDamageBonus(playerID int, skill *CardInstance) int {
	total := 0
	for _, modifier := range e.State.Players[playerID].TempModifiers {
		if modifier.Type != TempModNextElementSpellDamageBonus {
			continue
		}
		if modifier.RemainingUses == 0 || modifier.Status != skill.Card.Category {
			continue
		}
		total += modifier.Amount
	}
	return total
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
		if modifier.Type != TempModNextSpellHitStatus {
			continue
		}
		if modifier.RemainingUses == 0 {
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
			unit.Statuses[modifier.Status] += amount
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
		modifier.RemainingUses--
		if modifier.RemainingUses <= 0 {
			e.removeTemporaryModifier(playerID, modifier.ID)
		}
	}
}
