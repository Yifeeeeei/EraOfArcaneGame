package game

// A pending effect holds this field visit, not a later incarnation of the same
// card. Leaving clears the live instance while preserving last-known rules data.
type fieldLifetime struct {
	card     *CardInstance
	departed *CardInstance
}

func fieldVisit(card *CardInstance) *fieldLifetime {
	if card == nil {
		return nil
	}
	if card.fieldLife == nil {
		card.fieldLife = &fieldLifetime{card: card}
	}
	return card.fieldLife
}

func (visit *fieldLifetime) state() *CardInstance {
	if visit == nil {
		return nil
	}
	if visit.departed != nil {
		return visit.departed
	}
	return visit.card
}

func snapshotFieldCard(card *CardInstance) *CardInstance {
	if card == nil {
		return nil
	}
	copy := *card
	copy.Statuses = make(map[string]int, len(card.Statuses))
	for key, value := range card.Statuses {
		copy.Statuses[key] = value
	}
	copy.ElementsGainBonus = cloneElements(card.ElementsGainBonus)
	if card.ElementsGainSet != nil {
		copy.ElementsGainSet = cloneElements(card.ElementsGainSet)
	}
	copy.AttachedBehaviors = append([]AttachedBehavior(nil), card.AttachedBehaviors...)
	if card.Position != nil {
		pos := *card.Position
		copy.Position = &pos
	}
	copy.BoundSkills = make([]*CardInstance, len(card.BoundSkills))
	for i, skill := range card.BoundSkills {
		copy.BoundSkills[i] = snapshotFieldCard(skill)
	}
	return &copy
}

func (e *Engine) leaveBattlefield(card *CardInstance) *CardInstance {
	visit := fieldVisit(card)
	if card == nil {
		return nil
	}
	visit.departed = card
	// The old field object remains last-known information for in-flight damage
	// and death handlers. Only the new zone object is playable/serializable.
	clean := snapshotFieldCard(card)
	e.releaseUnderCardsToGraveyard(card.OwnerID, card)
	e.exileTransferredBoundSkills(card.OwnerID, card)
	resetCardForPublicSpecialZone(clean)
	return clean
}

func (e *Engine) effectSourceAvailable(visit *fieldLifetime) bool {
	return visit == nil || !e.hasEffectiveStatus(visit.state(), StatusPetrify)
}

func transferCardControl(card *CardInstance, playerID int) {
	if card == nil {
		return
	}
	card.OwnerID = playerID
	for _, skill := range card.BoundSkills {
		transferCardControl(skill, playerID)
	}
}

func (e *Engine) isBattlefieldInstance(card *CardInstance) bool {
	if card == nil {
		return false
	}
	for _, ps := range e.State.Players {
		if ps == nil {
			continue
		}
		for _, host := range e.getAllFieldCards(ps) {
			if host == card {
				return true
			}
			for _, skill := range host.BoundSkills {
				if skill == card {
					return true
				}
			}
		}
	}
	return false
}
