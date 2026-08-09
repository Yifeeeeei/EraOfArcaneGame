package game

const enteredGraveyardTurnStatus = "entered_graveyard_turn"

func (e *Engine) addToGraveyard(playerID int, card *CardInstance) {
	if e == nil || card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	if card.Statuses == nil {
		card.Statuses = make(map[string]int)
	}
	card.Statuses[enteredGraveyardTurnStatus] = e.State.TurnNumber
	ps.Graveyard = append(ps.Graveyard, card)
	e.State.CardEnteredGraveyardThisTurn = true
}

func (e *Engine) clearGraveyardTurnTracking() {
	if e == nil {
		return
	}
	e.State.CardEnteredGraveyardThisTurn = false
	for _, ps := range e.State.Players {
		if ps == nil {
			continue
		}
		for _, card := range ps.Graveyard {
			if card != nil && card.Statuses != nil {
				delete(card.Statuses, enteredGraveyardTurnStatus)
			}
		}
	}
}
