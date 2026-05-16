package game

const (
	TempModNextSkillCostZero            = "next_skill_cost_zero"
	TempModNextLearnedSkillHaste        = "next_learned_skill_haste"
	TempModSkillPowerBonus              = "skill_power_bonus"
	TempModNextNoCooldown               = "next_skill_no_cooldown"
	TempModNextSpellHitStatus           = "next_spell_hit_status"
	TempModDelayedElementGain           = "delayed_element_gain"
	TempModResetSkillsOnOpponentTurnEnd = "reset_skills_on_opponent_turn_end"
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
		if modifier.Type != TempModNextSkillCostZero {
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
		if modifier.Type != TempModSkillPowerBonus {
			continue
		}
		if modifier.TargetInstanceID == "" || modifier.TargetInstanceID == skill.InstanceID {
			total += modifier.Amount
		}
	}
	return total
}

func (e *Engine) applyTemporarySpellHitStatus(playerID int, skill *CardInstance, affectedUnits []*CardInstance) {
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
