package game

func (e *Engine) discardMarkedEndOfTurnCards(ps *PlayerState) {
	if len(ps.DiscardAtTurnEnd) == 0 {
		return
	}
	kept := ps.Hand[:0]
	for _, card := range ps.Hand {
		if card != nil && ps.DiscardAtTurnEnd[card.InstanceID] {
			e.discardHandCardToGraveyard(ps.PlayerID, card)
			continue
		}
		kept = append(kept, card)
	}
	ps.Hand = kept
	ps.DiscardAtTurnEnd = make(map[string]bool)
}

func (e *Engine) discardHandCardAt(playerID int, index int) *CardInstance {
	ps := e.State.Players[playerID]
	if ps == nil || index < 0 || index >= len(ps.Hand) {
		return nil
	}
	card := ps.Hand[index]
	ps.Hand = append(ps.Hand[:index], ps.Hand[index+1:]...)
	e.discardHandCardToGraveyard(playerID, card)
	return card
}

func (e *Engine) discardHandCardToGraveyard(playerID int, card *CardInstance) {
	if card == nil {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	e.addToGraveyard(playerID, card)
	ps.DiscardedHandCountThisTurn++
	delete(ps.RevealedHand, card.InstanceID)
	e.emit(GameEvent{Type: "discard", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
	e.triggerDiscardEffects(playerID, card)
	e.resolveDiscardedCardEffects(playerID, card)
}

func (e *Engine) triggerDiscardEffects(playerID int, card *CardInstance) {
	if e == nil || card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	data := map[string]any{
		"discarded_player": playerID,
		"discarded_card":   card,
	}
	e.triggerFieldEffectsWithData(TriggerOnDiscard, playerID, card, data)
	e.triggerFieldEffectsWithData(TriggerOnDiscard, 1-playerID, card, data)
}
