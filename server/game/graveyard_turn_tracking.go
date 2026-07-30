package game

func (e *Engine) addToGraveyard(playerID int, card *CardInstance) {
	if e == nil || card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	ps.Graveyard = append(ps.Graveyard, card)
	e.State.CardEnteredGraveyardThisTurn = true
}

func (e *Engine) clearGraveyardTurnTracking() {
	if e == nil {
		return
	}
	e.State.CardEnteredGraveyardThisTurn = false
}
