package game

func (e *Engine) discardMarkedEndOfTurnCards(ps *PlayerState) {
	if len(ps.DiscardAtTurnEnd) == 0 {
		return
	}
	kept := ps.Hand[:0]
	for _, card := range ps.Hand {
		if card != nil && ps.DiscardAtTurnEnd[card.InstanceID] {
			ps.Graveyard = append(ps.Graveyard, card)
			e.emit(GameEvent{Type: "discard", Player: ps.PlayerID, Data: map[string]any{"card": cardToInfo(card)}})
			continue
		}
		kept = append(kept, card)
	}
	ps.Hand = kept
	ps.DiscardAtTurnEnd = make(map[string]bool)
}
