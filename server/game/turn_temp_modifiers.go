package game

func (e *Engine) applyTurnStartTemporaryModifiers(ps *PlayerState) {
	// Unqualified "next turn" means the next player's turn, regardless of owner.
	for ownerID, owner := range e.State.Players {
		if owner == nil {
			continue
		}
		for _, modifier := range append([]TemporaryModifier(nil), owner.TempModifiers...) {
			if modifier.Type == TempModDelayedShieldGain {
				e.gainPlayerShield(ownerID, max(1, modifier.Amount))
				e.removeTemporaryModifier(ownerID, modifier.ID)
			}
		}
	}
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		switch modifier.Type {
		case TempModDelayedElementGain:
			amount := modifier.Amount
			if amount <= 0 {
				amount = 1
			}
			ps.GainElements(map[string]int{modifier.Status: amount})
			e.emit(GameEvent{Type: "effect_trigger", Player: ps.PlayerID, Data: map[string]any{
				"effect":   "gain_element",
				"element":  modifier.Status,
				"amount":   amount,
				"elements": ps.Elements,
			}})
			e.removeTemporaryModifier(ps.PlayerID, modifier.ID)
		default:
			continue
		}
	}
}

func (e *Engine) applyOpponentTurnEndTemporaryModifiers(endedPlayerID int) {
	for ownerID, ps := range e.State.Players {
		if ownerID == endedPlayerID {
			continue
		}
		for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
			switch modifier.Type {
			case TempModResetSkillsOnOpponentTurnEnd:
				for _, skill := range ps.Skills {
					resetCard(skill)
				}
				e.emit(GameEvent{Type: "effect_trigger", Player: ownerID, Data: map[string]any{
					"effect": "reset_skills",
				}})
				e.removeTemporaryModifier(ownerID, modifier.ID)
			case TempModLampusSwordDelayedDamage:
				e.promptLampusSwordDelayedDamage(ownerID, modifier)
			default:
				continue
			}
		}
	}
}
