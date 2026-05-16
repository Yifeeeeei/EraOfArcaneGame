package game

func scheduleLoadGainAtTurnEnd(ps *PlayerState, instanceID string, elem string, amount int) {
	if ps == nil || instanceID == "" || elem == "" || amount == 0 {
		return
	}
	if ps.LoadGainAtTurnEnd == nil {
		ps.LoadGainAtTurnEnd = make(map[string]map[string]int)
	}
	if ps.LoadGainAtTurnEnd[instanceID] == nil {
		ps.LoadGainAtTurnEnd[instanceID] = make(map[string]int)
	}
	ps.LoadGainAtTurnEnd[instanceID][elem] += amount
}

func (e *Engine) applyLoadGainAtTurnEnd(ps *PlayerState) {
	if ps == nil || len(ps.LoadGainAtTurnEnd) == 0 {
		return
	}
	for instanceID, gains := range ps.LoadGainAtTurnEnd {
		card := e.findFieldCardByInstance(ps, instanceID)
		if card == nil {
			continue
		}
		for elem, amount := range gains {
			addElementsGainBonus(card, elem, amount)
			e.emit(GameEvent{Type: "effect_trigger", Player: ps.PlayerID, Data: map[string]any{
				"source":  cardToInfo(card),
				"effect":  "load_gain",
				"element": elem,
				"amount":  amount,
			}})
		}
	}
	ps.LoadGainAtTurnEnd = make(map[string]map[string]int)
}

func (e *Engine) findFieldCardByInstance(ps *PlayerState, instanceID string) *CardInstance {
	for _, card := range e.getAllFieldCards(ps) {
		if card != nil && card.InstanceID == instanceID {
			return card
		}
	}
	return nil
}
