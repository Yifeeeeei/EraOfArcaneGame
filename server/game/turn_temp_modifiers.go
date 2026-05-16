package game

func (e *Engine) applyTurnStartTemporaryModifiers(ps *PlayerState) {
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.Type != TempModDelayedElementGain {
			continue
		}
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
	}
}

func (e *Engine) applyOpponentTurnEndTemporaryModifiers(endedPlayerID int) {
	for ownerID, ps := range e.State.Players {
		if ownerID == endedPlayerID {
			continue
		}
		for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
			if modifier.Type != TempModResetSkillsOnOpponentTurnEnd {
				continue
			}
			for _, skill := range ps.Skills {
				resetCard(skill)
			}
			e.emit(GameEvent{Type: "effect_trigger", Player: ownerID, Data: map[string]any{
				"effect": "reset_skills",
			}})
			e.removeTemporaryModifier(ownerID, modifier.ID)
		}
	}
}
