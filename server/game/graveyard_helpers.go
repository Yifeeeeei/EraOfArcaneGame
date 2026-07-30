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
		e.ApplySummonModifiersOnEnter(card)
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
		e.releaseUnderCardsToGraveyard(playerID, card)
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
					e.releaseUnderCardsToGraveyard(playerID, unit)
					unit.BoundSkills = nil
					e.emit(GameEvent{Type: "card_removed_from_game", Player: playerID, Data: map[string]any{"card": cardToInfo(unit)}})
					return true
				}
				if terrain := ps.Terrain[col][row]; terrain != nil && terrain.InstanceID == instanceID {
					ps.Terrain[col][row] = nil
					terrain.Position = nil
					e.releaseUnderCardsToGraveyard(playerID, terrain)
					e.emit(GameEvent{Type: "card_removed_from_game", Player: playerID, Data: map[string]any{"card": cardToInfo(terrain)}})
					return true
				}
			}
		}
		for i, skill := range ps.Skills {
			if skill != nil && skill.InstanceID == instanceID {
				ps.Skills[i] = nil
				e.releaseUnderCardsToGraveyard(playerID, skill)
				e.emit(GameEvent{Type: "card_removed_from_game", Player: playerID, Data: map[string]any{"card": cardToInfo(skill)}})
				if skill.Card.Number == "3611101" {
					e.refreshRedMoonState(playerID)
				}
				return true
			}
		}
		for i, equipment := range ps.Equipment {
			if equipment != nil && equipment.InstanceID == instanceID {
				ps.Equipment[i] = nil
				e.releaseUnderCardsToGraveyard(playerID, equipment)
				e.emit(GameEvent{Type: "card_removed_from_game", Player: playerID, Data: map[string]any{"card": cardToInfo(equipment)}})
				return true
			}
		}
	}
	return false
}

func (e *Engine) exileCard(playerID int, card *CardInstance) bool {
	if card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	owner := e.State.Players[playerID]
	for zoneOwnerID, ps := range e.State.Players {
		if ps == nil {
			continue
		}
		if removed := e.detachCardFromKnownZones(ps, card.InstanceID); removed != nil {
			e.releaseUnderCardsToGraveyard(zoneOwnerID, removed)
			resetCardForPublicSpecialZone(removed)
			owner.Exile = append(owner.Exile, removed)
			e.emit(GameEvent{Type: "card_exiled", Player: playerID, Data: map[string]any{"card": cardToInfo(removed)}})
			return true
		}
	}
	return false
}

func (e *Engine) placeCardUnder(host *CardInstance, card *CardInstance) bool {
	if host == nil || card == nil || host.InstanceID == card.InstanceID {
		return false
	}
	for _, ps := range e.State.Players {
		if ps == nil {
			continue
		}
		if removed := e.detachCardFromKnownZones(ps, card.InstanceID); removed != nil {
			resetCardForPublicSpecialZone(removed)
			host.UnderCards = append(host.UnderCards, removed)
			e.emit(GameEvent{Type: "card_placed_under", Player: host.OwnerID, Data: map[string]any{
				"host": cardToInfo(host),
				"card": cardToInfo(removed),
			}})
			return true
		}
	}
	return false
}

func (e *Engine) detachCardFromKnownZones(ps *PlayerState, instanceID string) *CardInstance {
	if ps == nil || instanceID == "" {
		return nil
	}
	for i, card := range ps.Hand {
		if card != nil && card.InstanceID == instanceID {
			ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
			if ps.RevealedHand != nil {
				delete(ps.RevealedHand, instanceID)
			}
			return card
		}
	}
	for i, card := range ps.Deck {
		if card != nil && card.InstanceID == instanceID {
			ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
			return card
		}
	}
	for i, card := range ps.SkillPool {
		if card != nil && card.InstanceID == instanceID {
			ps.SkillPool = append(ps.SkillPool[:i], ps.SkillPool[i+1:]...)
			return card
		}
	}
	for i, card := range ps.Graveyard {
		if card != nil && card.InstanceID == instanceID {
			ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
			return card
		}
	}
	for i, card := range ps.Exile {
		if card != nil && card.InstanceID == instanceID {
			ps.Exile = append(ps.Exile[:i], ps.Exile[i+1:]...)
			return card
		}
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if unit := ps.Units[col][row]; unit != nil && unit.InstanceID == instanceID && !unit.Card.IsHero() {
				ps.Units[col][row] = nil
				return unit
			}
			if terrain := ps.Terrain[col][row]; terrain != nil && terrain.InstanceID == instanceID {
				ps.Terrain[col][row] = nil
				return terrain
			}
		}
	}
	for i, card := range ps.Skills {
		if card != nil && card.InstanceID == instanceID {
			ps.Skills[i] = nil
			return card
		}
	}
	for i, card := range ps.Equipment {
		if card != nil && card.InstanceID == instanceID {
			ps.Equipment[i] = nil
			return card
		}
	}
	for _, host := range e.getAllFieldCards(ps) {
		if host == nil {
			continue
		}
		for i, under := range host.UnderCards {
			if under != nil && under.InstanceID == instanceID {
				host.UnderCards = append(host.UnderCards[:i], host.UnderCards[i+1:]...)
				return under
			}
		}
		for i, bound := range host.BoundSkills {
			if bound != nil && bound.InstanceID == instanceID {
				host.BoundSkills = append(host.BoundSkills[:i], host.BoundSkills[i+1:]...)
				return bound
			}
		}
	}
	return nil
}

func (e *Engine) releaseUnderCardsToGraveyard(ownerID int, host *CardInstance) {
	if host == nil || len(host.UnderCards) == 0 || ownerID < 0 || ownerID >= len(e.State.Players) {
		return
	}
	for _, card := range host.UnderCards {
		if card == nil {
			continue
		}
		graveyardOwner := card.OwnerID
		if graveyardOwner < 0 || graveyardOwner >= len(e.State.Players) || e.State.Players[graveyardOwner] == nil {
			graveyardOwner = ownerID
		}
		resetCardForPublicSpecialZone(card)
		e.addToGraveyard(graveyardOwner, card)
		e.emit(GameEvent{Type: "discard", Player: graveyardOwner, Data: map[string]any{
			"card":   cardToInfo(card),
			"reason": "under_card_released",
		}})
	}
	host.UnderCards = nil
}
