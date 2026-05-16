package game

func (e *Engine) sacrificeEquipment(playerID int, instanceID string) bool {
	ps := e.State.Players[playerID]
	for i, card := range ps.Equipment {
		if card == nil || card.InstanceID != instanceID {
			continue
		}
		ps.Equipment[i] = nil
		ps.Graveyard = append(ps.Graveyard, card)
		e.emit(GameEvent{Type: "discard", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
		return true
	}
	return false
}
