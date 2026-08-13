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

func equipmentSlotOf(ps *PlayerState, card *CardInstance) int {
	if ps == nil || card == nil {
		return -1
	}
	for i, equipment := range ps.Equipment {
		if equipment != nil && equipment.InstanceID == card.InstanceID {
			return i
		}
	}
	return -1
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
	e.exileTransferredBoundSkills(playerID, card)
	card.BoundSkills = nil
	e.addToGraveyard(playerID, card)
	e.emit(GameEvent{Type: "discard", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
}
