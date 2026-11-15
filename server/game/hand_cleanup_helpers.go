package game

func (e *Engine) discardMarkedEndOfTurnCards(ps *PlayerState) {
	if len(ps.DiscardAtTurnEnd) == 0 {
		return
	}
	kept := ps.Hand[:0]
	var discarded []*CardInstance
	for _, card := range ps.Hand {
		if card != nil && ps.DiscardAtTurnEnd[card.InstanceID] {
			discarded = append(discarded, card)
			continue
		}
		kept = append(kept, card)
	}
	ps.Hand = kept
	e.discardHandCardsToGraveyard(ps.PlayerID, discarded)
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
	e.discardHandCardsToGraveyard(playerID, []*CardInstance{card})
}

type DiscardBatchBehavior interface {
	HasActiveDiscardBatch(*CardInstance) bool
	OnDiscardBatch(*EffectContext, []*CardInstance) error
}

func (AlwaysActive) HasActiveDiscardBatch(*CardInstance) bool { return true }

// Per-card triggers retain their semantics; batch observers receive one event
// after all cards in this discard action have reached the graveyard.
func (e *Engine) discardHandCardsToGraveyard(playerID int, cards []*CardInstance) {
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	var discarded []*CardInstance
	for _, card := range cards {
		if card == nil {
			continue
		}
		e.addToGraveyard(playerID, card)
		ps.DiscardedHandCountThisTurn++
		delete(ps.RevealedHand, card.InstanceID)
		discarded = append(discarded, card)
		e.emit(GameEvent{Type: "discard", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
		e.triggerDiscardEffects(playerID, card)
		e.resolveDiscardedCardEffects(playerID, card)
	}
	if len(discarded) == 0 {
		return
	}
	for _, source := range e.ruleSources(ps) {
		if e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		if b, ok := cardBehavior(source).(DiscardBatchBehavior); ok && b.HasActiveDiscardBatch(source) {
			_ = b.OnDiscardBatch(e.skillContext(playerID, source), discarded)
		}
	}
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
