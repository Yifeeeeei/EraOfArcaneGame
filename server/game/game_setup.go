package game

import (
	"eraofarcane/model"
)

// SetupGame initializes both players and starts the game
func (e *Engine) SetupGame(p1Name string, p1Deck *model.Deck, p2Name string, p2Deck *model.Deck) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.setupGameWithFirstPlayer(p1Name, p1Deck, p2Name, p2Deck, e.randomIntn(2))
}

// SetupGameWithFirstPlayer initializes both players with an explicit first player.
func (e *Engine) SetupGameWithFirstPlayer(p1Name string, p1Deck *model.Deck, p2Name string, p2Deck *model.Deck, firstPlayer int) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.setupGameWithFirstPlayer(p1Name, p1Deck, p2Name, p2Deck, firstPlayer)
}

func (e *Engine) setupGameWithFirstPlayer(p1Name string, p1Deck *model.Deck, p2Name string, p2Deck *model.Deck, firstPlayer int) error {
	if firstPlayer < 0 || firstPlayer > 1 {
		firstPlayer = 0
	}

	// Create player states
	e.State.Players[0] = NewPlayerState(0, p1Name, p1Deck)
	e.State.Players[1] = NewPlayerState(1, p2Name, p2Deck)
	e.State.FirstPlayer = firstPlayer

	// Initialize cards
	e.initPlayerCards(e.State.Players[0], 0)
	e.initPlayerCards(e.State.Players[1], 0)
	e.triggerInitialHeroEnterEffects()
	e.triggerGameStartEffects()
	e.emit(GameEvent{
		Type:   "game_setup",
		Player: -1,
		Data: map[string]any{
			"first_player": e.State.FirstPlayer,
			"timing":       "before_initial_draw",
		},
	})

	// Draw initial hands (4 cards each; Raven starts with one extra card)
	for i := 0; i < 2; i++ {
		drawn := e.drawCards(i, e.initialHandSizeForPlayer(e.State.Players[i]))
		e.emit(GameEvent{
			Type:   "initial_draw",
			Player: i,
			Data: map[string]any{
				"cards": cardsToInfo(drawn),
				"count": len(drawn),
			},
		})
	}

	// Enter mulligan phase
	e.State.Phase = PhaseMulligan
	e.emit(GameEvent{
		Type:   "phase_change",
		Player: -1,
		Data:   map[string]any{"phase": "mulligan"},
	})

	return nil
}

func (e *Engine) triggerGameStartEffects() {
	for playerID := 0; playerID < 2; playerID++ {
		data := map[string]any{"initial_setup": true}
		for _, card := range e.getAllFieldCards(e.State.Players[playerID]) {
			e.triggerEffects(TriggerOnGameStart, card, nil, data)
		}
	}
}

func (e *Engine) triggerInitialHeroEnterEffects() {
	for playerID := 0; playerID < 2; playerID++ {
		hero := e.State.Players[playerID].Hero
		if hero == nil {
			continue
		}
		data := map[string]any{"initial_setup": true, "entered_player": playerID}
		e.triggerEffects(TriggerOnEnter, hero, nil, data)
		e.notifyCardEntered(playerID, hero, data)
		e.triggerFieldEffectsWithData(TriggerOnUnitEnter, playerID, hero, data)
		e.triggerFieldEffectsWithData(TriggerOnUnitEnter, 1-playerID, hero, data)
	}
}

func (e *Engine) initialHandSizeForPlayer(ps *PlayerState) int {
	if ps == nil {
		return 4
	}
	size := 4
	if b, ok := cardBehavior(ps.Hero).(OpeningHandBehavior); ok && b.HasActiveOpeningHand(ps.Hero) {
		size += b.OpeningHandBonus(ps.Hero)
	}
	return size
}
