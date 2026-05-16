package game

func (e *Engine) friendlyTopDeckCards(playerID int, limit int, predicate func(*CardInstance) bool) []map[string]any {
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0)
	for i, card := range ps.Deck {
		if i >= limit {
			break
		}
		if card == nil {
			continue
		}
		if predicate != nil && !predicate(card) {
			continue
		}
		candidates = append(candidates, candidateInfo(card, "deck", "own"))
	}
	return candidates
}

func (e *Engine) moveDeckCardToBottom(playerID int, instanceID string) bool {
	ps := e.State.Players[playerID]
	for i, card := range ps.Deck {
		if card == nil || card.InstanceID != instanceID {
			continue
		}
		ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
		ps.Deck = append(ps.Deck, card)
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"effect": "deck_card_to_bottom",
			"card":   cardToInfo(card),
		}})
		return true
	}
	return false
}

func (e *Engine) moveHandCardToDeckBottom(playerID int, instanceID string) bool {
	ps := e.State.Players[playerID]
	for i, card := range ps.Hand {
		if card == nil || card.InstanceID != instanceID {
			continue
		}
		ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
		ps.Deck = append(ps.Deck, card)
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"effect": "hand_card_to_deck",
			"card":   cardToInfo(card),
		}})
		return true
	}
	return false
}
