package game

import ()

func (e *Engine) drawCards(playerID int, n int) []*CardInstance {
	if n <= 0 {
		return nil
	}
	ps := e.State.Players[playerID]
	drawn := ps.DrawCards(n)
	if len(drawn) < n && len(ps.Deck) == 0 {
		for _, source := range e.getAllFieldCards(ps) {
			if source == nil || e.hasEffectiveStatus(source, StatusPetrify) {
				continue
			}
			behavior, ok := cardBehavior(source).(EmptyDeckDrawBehavior)
			if !ok || !behavior.HasActiveEmptyDeckDraw(source) {
				continue
			}
			ctx := &EffectContext{Engine: e, Source: source, PlayerID: playerID, OpponentID: 1 - playerID}
			drawn = append(drawn, behavior.DrawFromEmptyDeck(ctx, n-len(drawn))...)
			break
		}
	}
	for _, card := range drawn {
		e.notifyCardDrawn(playerID, card)
	}
	if e.shouldImmediatelyEnforceHandLimit(playerID) {
		e.promptDiscardToHandLimit(playerID, nil)
	}
	return drawn
}

func (e *Engine) addCardToHand(playerID int, card *CardInstance) bool {
	return e.addCardsToHand(playerID, []*CardInstance{card}) > 0
}

func (e *Engine) addCardsToHand(playerID int, cards []*CardInstance) int {
	added := e.appendCardsToHand(playerID, cards)
	if added > 0 {
		e.enforceImmediateHandLimitAfterHandGain(playerID)
	}
	return added
}

func (e *Engine) appendCardsToHand(playerID int, cards []*CardInstance) int {
	if playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return 0
	}
	added := 0
	for _, card := range cards {
		if card == nil {
			continue
		}
		ps.Hand = append(ps.Hand, card)
		added++
	}
	return added
}

func (e *Engine) enforceImmediateHandLimitAfterHandGain(playerID int) bool {
	if !e.shouldImmediatelyEnforceHandLimit(playerID) {
		return false
	}
	return e.promptDiscardToHandLimit(playerID, nil)
}

func (e *Engine) millTopDeckCards(playerID int, n int) []*CardInstance {
	if n <= 0 {
		return nil
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	count := min(n, len(ps.Deck))
	milled := make([]*CardInstance, 0, count)
	for i := 0; i < count; i++ {
		card := ps.Deck[0]
		ps.Deck = ps.Deck[1:]
		e.addToGraveyard(playerID, card)
		milled = append(milled, card)
		e.emit(GameEvent{Type: "discard", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
	}
	return milled
}

func (e *Engine) notifyCardDrawn(playerID int, card *CardInstance) {
	if card == nil {
		return
	}
	ps := e.State.Players[playerID]
	if ps.DrawnTurn == nil {
		ps.DrawnTurn = make(map[string]int)
	}
	if cardRevealsOnDraw(card) {
		ps.RevealedHand[card.InstanceID] = true
	}
	ps.DrawnTurn[card.InstanceID] = e.State.TurnNumber
	ps.DrawCountThisTurn++
	e.emit(GameEvent{
		Type:   "draw_card",
		Player: playerID,
		Data:   map[string]any{"card": cardToInfo(card)},
	})
	data := map[string]any{
		"drawn_card":           card,
		"drawn_player":         playerID,
		"draw_count_this_turn": ps.DrawCountThisTurn,
		"initial_hand":         e.State.Phase == PhaseWaitingPlayers || e.State.Phase == PhaseMulligan,
	}
	e.triggerFieldEffectsWithData(TriggerOnDraw, playerID, card, data)
	e.triggerFieldEffectsWithData(TriggerOnDraw, 1-playerID, card, data)
	if h, ok := cardBehavior(card).(OnSelfDrawBehavior); ok && h.HasActiveDraw(card) {
		_ = h.OnSelfDraw(&EffectContext{
			Engine:     e,
			Source:     card,
			PlayerID:   playerID,
			OpponentID: 1 - playerID,
			ExtraData:  data,
		})
	}
}
