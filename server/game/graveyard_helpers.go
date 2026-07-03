package game

func (e *Engine) reviveCompanionFromGraveyard(playerID int, instanceID string) bool {
	return e.reviveCompanionFromGraveyardWithLife(playerID, instanceID, 0, true)
}

func (e *Engine) reviveCompanionFromGraveyardAtOneLife(playerID int, instanceID string) bool {
	return e.reviveCompanionFromGraveyardWithLife(playerID, instanceID, 1, false)
}

func (e *Engine) reviveCompanionFromGraveyardWithLife(playerID int, instanceID string, life int, payCost bool) bool {
	ps := e.State.Players[playerID]
	pos := ps.FindEmptyPosition()
	if pos == nil {
		return false
	}
	return e.reviveCompanionFromGraveyardWithLifeAtPosition(playerID, instanceID, life, payCost, *pos)
}

func (e *Engine) reviveCompanionFromGraveyardWithLifeAtPosition(playerID int, instanceID string, life int, payCost bool, pos Position) bool {
	ps := e.State.Players[playerID]
	if !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil {
		return false
	}
	for i, card := range ps.Graveyard {
		if card == nil || card.InstanceID != instanceID || !card.Card.IsCompanion() {
			continue
		}
		if payCost {
			if !e.canPayCost(ps, card.Card.ElementsCost) {
				return false
			}
			e.payCostForAction(ps, card.Card.ElementsCost, ActionMessage{})
		}
		ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
		if life > 0 {
			card.CurrentLife = life
		} else {
			card.CurrentLife = max(card.CurrentLife, card.Card.Life)
		}
		card.IsHorizontal = true
		card.Position = &pos
		ps.Units[pos.Col][pos.Row] = card
		e.emit(GameEvent{Type: "summon", Player: -1, Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"position": &pos,
			"elements": ps.Elements,
		}})
		e.triggerEffects(TriggerOnEnter, card, nil, nil)
		return true
	}
	return false
}

func (e *Engine) removeEquipmentFromGame(playerID int, instanceID string) bool {
	ps := e.State.Players[playerID]
	for i, card := range ps.Equipment {
		if card == nil || card.InstanceID != instanceID {
			continue
		}
		ps.Equipment[i] = nil
		e.emit(GameEvent{Type: "card_removed_from_game", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
		return true
	}
	return false
}

func (e *Engine) removeStoredArchmageStaffSkillAfterUse(playerID int, skill *CardInstance) bool {
	if skill == nil || skill.Statuses[archmageStaffStoredSkillStatus] <= 0 {
		return false
	}
	ps := e.State.Players[playerID]
	for _, host := range e.getAllFieldCards(ps) {
		if host == nil {
			continue
		}
		for i, bound := range host.BoundSkills {
			if bound == nil || bound.InstanceID != skill.InstanceID {
				continue
			}
			host.BoundSkills = append(host.BoundSkills[:i], host.BoundSkills[i+1:]...)
			e.emit(GameEvent{Type: "card_removed_from_game", Player: playerID, Data: map[string]any{"card": cardToInfo(skill)}})
			return true
		}
	}
	return false
}

func (e *Engine) nonHeroFieldCardCandidates(playerID int) []map[string]any {
	candidates := make([]map[string]any, 0)
	for ownerID, ps := range e.State.Players {
		side := "enemy"
		if ownerID == playerID {
			side = "own"
		}
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				if unit := ps.Units[col][row]; unit != nil && !unit.Card.IsHero() {
					candidates = append(candidates, candidateInfo(unit, "unit", side))
				}
				if terrain := ps.Terrain[col][row]; terrain != nil {
					candidates = append(candidates, candidateInfo(terrain, "terrain", side))
				}
			}
		}
		for _, skill := range ps.Skills {
			if skill != nil {
				candidates = append(candidates, candidateInfo(skill, "skill", side))
			}
		}
		for _, equipment := range ps.Equipment {
			if equipment != nil {
				candidates = append(candidates, candidateInfo(equipment, "equipment", side))
			}
		}
	}
	return candidates
}

func (e *Engine) removeFieldCardFromGameByID(instanceID string) bool {
	for playerID, ps := range e.State.Players {
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				if unit := ps.Units[col][row]; unit != nil && unit.InstanceID == instanceID && !unit.Card.IsHero() {
					ps.Units[col][row] = nil
					unit.Position = nil
					unit.BoundSkills = nil
					e.emit(GameEvent{Type: "card_removed_from_game", Player: playerID, Data: map[string]any{"card": cardToInfo(unit)}})
					return true
				}
				if terrain := ps.Terrain[col][row]; terrain != nil && terrain.InstanceID == instanceID {
					ps.Terrain[col][row] = nil
					terrain.Position = nil
					e.emit(GameEvent{Type: "card_removed_from_game", Player: playerID, Data: map[string]any{"card": cardToInfo(terrain)}})
					return true
				}
			}
		}
		for i, skill := range ps.Skills {
			if skill != nil && skill.InstanceID == instanceID {
				ps.Skills[i] = nil
				e.emit(GameEvent{Type: "card_removed_from_game", Player: playerID, Data: map[string]any{"card": cardToInfo(skill)}})
				return true
			}
		}
		for i, equipment := range ps.Equipment {
			if equipment != nil && equipment.InstanceID == instanceID {
				ps.Equipment[i] = nil
				e.emit(GameEvent{Type: "card_removed_from_game", Player: playerID, Data: map[string]any{"card": cardToInfo(equipment)}})
				return true
			}
		}
	}
	return false
}
