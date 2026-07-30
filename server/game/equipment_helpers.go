package game

func (e *Engine) sacrificeEquipment(playerID int, instanceID string) bool {
	ps := e.State.Players[playerID]
	for i, card := range ps.Equipment {
		if card == nil || card.InstanceID != instanceID {
			continue
		}
		e.moveEquipmentToGraveyard(playerID, i, card)
		return true
	}
	return false
}

func (e *Engine) moveEquipmentToGraveyard(playerID int, slot int, card *CardInstance) {
	if card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	if slot >= 0 && slot < len(ps.Equipment) && ps.Equipment[slot] == card {
		ps.Equipment[slot] = nil
	}
	card.SlotIndex = -1
	e.releaseUnderCardsToGraveyard(playerID, card)
	card.BoundSkills = nil
	e.addToGraveyard(playerID, card)
	e.emit(GameEvent{Type: "discard", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
}
