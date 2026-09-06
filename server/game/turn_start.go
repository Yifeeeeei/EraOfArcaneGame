package game

import (
	"fmt"
)

// handleMulligan handles the mulligan (redraw) action
func (e *Engine) handleMulligan(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMulligan {
		return fmt.Errorf("not in mulligan phase")
	}
	if e.State.MulliganDone[playerID] {
		return fmt.Errorf("already submitted mulligan")
	}

	// Check if player wants to redraw
	keepStr, _ := action.Data["keep"].(bool)

	if !keepStr {
		// Redraw to the same starting hand size. Raven starts with one extra card,
		// so its mulligan should also redraw that extra card.
		ps := e.State.Players[playerID]
		ps.Deck = append(ps.Deck, ps.Hand...)
		ps.Hand = make([]*CardInstance, 0)
		e.shuffleCards(ps.Deck)
		drawn := e.drawCards(playerID, e.initialHandSizeForPlayer(ps))
		e.emit(GameEvent{
			Type:   "mulligan_redraw",
			Player: playerID,
			Data: map[string]any{
				"cards": cardsToInfo(drawn),
				"count": len(drawn),
			},
		})
	}

	e.State.MulliganDone[playerID] = true

	e.emit(GameEvent{
		Type:   "mulligan_done",
		Player: playerID,
		Data:   map[string]any{"player": playerID},
	})

	// If both players are done, start the game
	if e.State.MulliganDone[0] && e.State.MulliganDone[1] {
		e.startGame()
	}

	return nil
}

// startGame begins the actual game
func (e *Engine) startGame() {
	e.State.CurrentTurn = e.State.FirstPlayer
	e.State.TurnNumber = 1
	e.State.IsFirstTurn = true

	e.emit(GameEvent{
		Type:   "game_start",
		Player: -1,
		Data: map[string]any{
			"first_player": e.State.FirstPlayer,
		},
	})

	e.startTurn()
}

// startTurn begins a new turn for the current player
func (e *Engine) startTurn() {
	e.clearDamageTakenThisTurn()

	ps := e.State.Players[e.State.CurrentTurn]
	clearSpellCastTracking(ps)
	e.clearGraveyardTurnTracking()
	ps.DrawCountThisTurn = 0
	ps.DiscardedHandCountThisTurn = 0

	// Elements are cleared at the end of their owner's turn. Start turn should
	// not be the rule point that removes remaining elements.
	e.applyTurnStartTemporaryModifiers(ps)

	e.continuePreDrawTurnStartEffects(ps, append([]*CardInstance(nil), e.getAllFieldCards(ps)...), 0)
}

func (e *Engine) continueStartTurnAfterPreDraw(ps *PlayerState) {
	if ps == nil {
		return
	}
	if ps.SkipNextDraw {
		ps.SkipNextDraw = false
		e.emit(GameEvent{
			Type:   "effect_trigger",
			Player: ps.PlayerID,
			Data: map[string]any{
				"effect": "skip_draw",
			},
		})
	} else {
		drawn := e.drawCards(ps.PlayerID, 1)
		if len(drawn) > 0 {
			// Notify opponent about the draw (without card info)
			e.emit(GameEvent{
				Type:   "opponent_draw",
				Player: 1 - ps.PlayerID,
				Data: map[string]any{
					"count": 1,
				},
			})
		}
		if e.State.PendingAction != nil {
			e.continueAfterPendingAction(func() {
				e.continueStartTurnAfterDraw(ps)
			})
			return
		}
	}

	e.continueStartTurnAfterDraw(ps)
}

func (e *Engine) continueStartTurnAfterDraw(ps *PlayerState) {
	if ps == nil {
		return
	}
	e.State.Phase = PhaseMain

	e.emit(GameEvent{
		Type:   "turn_start",
		Player: -1,
		Data: map[string]any{
			"current_player": ps.PlayerID,
			"turn_number":    e.State.TurnNumber,
			"elements":       ps.Elements,
		},
	})

	// Trigger 回合开始 effects for all cards on the current player's field
	allCards := e.getAllFieldCards(ps)
	for _, card := range allCards {
		e.triggerEffects(TriggerOnTurnStart, card, nil, nil)
	}
	if e.State.PendingAction == nil {
		e.triggerPrayerAbilities(ps.PlayerID)
	}
}

func (e *Engine) continuePreDrawTurnStartEffects(ps *PlayerState, cards []*CardInstance, start int) {
	if ps == nil {
		return
	}
	for i := start; i < len(cards); i++ {
		card := cards[i]
		if card == nil || card.Card == nil || e.hasEffectiveStatus(card, StatusPetrify) {
			continue
		}
		if !e.cardStillOnField(card) {
			continue
		}
		behavior, ok := globalRegistry.GetBehavior(card.Card.Number).(BeforeDrawBehavior)
		if !ok || !behavior.HasActiveBeforeDraw(card) {
			continue
		}
		ctx := &EffectContext{
			Engine:     e,
			Source:     card,
			PlayerID:   ps.PlayerID,
			OpponentID: 1 - ps.PlayerID,
			ExtraData:  map[string]any{"timing": "pre_draw"},
		}
		_ = behavior.OnBeforeDraw(ctx)
		if e.State.PendingAction != nil {
			next := i + 1
			e.continueAfterPendingAction(func() {
				e.continuePreDrawTurnStartEffects(ps, cards, next)
			})
			return
		}
	}
	e.continueStartTurnAfterPreDraw(ps)
}

func (e *Engine) triggerPrayerAbilities(playerID int) {
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	e.continuePrayerAbilities(playerID, append([]*CardInstance(nil), e.getAllFieldCards(ps)...), 0)
}

func (e *Engine) continuePrayerAbilities(playerID int, cards []*CardInstance, start int) {
	for i := start; i < len(cards); i++ {
		card := cards[i]
		if !cardHasActivePrayer(card) || !e.cardStillOnField(card) {
			continue
		}
		if e.executePrayerAbility(playerID, card); e.State.PendingAction != nil {
			next := i + 1
			e.continueAfterPendingAction(func() {
				e.continuePrayerAbilities(playerID, cards, next)
			})
			return
		}
	}
}

func (e *Engine) executePrayerAbility(playerID int, card *CardInstance) {
	behavior := cardBehavior(card)
	perTurn, ok := behavior.(PerTurnAbility)
	if !ok || !perTurn.HasActivePerTurn(card) {
		return
	}
	ctx := &EffectContext{
		Engine:     e,
		Source:     card,
		PlayerID:   playerID,
		OpponentID: 1 - playerID,
		ExtraData:  map[string]any{"prayer": true},
	}
	run := func() {
		e.emit(GameEvent{
			Type:   "effect_trigger",
			Player: -1,
			Data: map[string]any{
				"effect": "prayer",
				"card":   cardToInfo(card),
				"player": playerID,
			},
		})
		if err := perTurn.OnPerTurn(ctx); err != nil {
			e.emit(GameEvent{
				Type:   "effect_trigger",
				Player: playerID,
				Data: map[string]any{
					"effect": "prayer_error",
					"card":   cardToInfo(card),
					"error":  err.Error(),
				},
			})
		}
	}
	if optional, ok := behavior.(OptionalPrayerAbility); ok && optional.IsPrayerOptional(card) {
		e.SetPendingAction(playerID, "optional_prayer", "是否发动祈咒: "+card.Card.Name, []map[string]any{candidateInfo(card, "card", "own")}, 0, 1, func(selected []string) {
			if len(selected) > 0 {
				run()
			}
		})
		return
	}
	run()
}

func (e *Engine) cardStillOnField(card *CardInstance) bool {
	if card == nil || card.OwnerID < 0 || card.OwnerID >= len(e.State.Players) {
		return false
	}
	for _, fieldCard := range e.getAllFieldCards(e.State.Players[card.OwnerID]) {
		if fieldCard == card {
			return true
		}
	}
	return false
}

func (e *Engine) clearDamageTakenThisTurn() {
	for _, ps := range e.State.Players {
		if ps == nil {
			continue
		}
		for _, card := range e.getAllFieldCards(ps) {
			if card != nil {
				card.DamageTakenThisTurn = 0
			}
		}
		for _, card := range ps.Graveyard {
			if card != nil {
				card.DamageTakenThisTurn = 0
			}
		}
	}
}

func (e *Engine) rollFriendlyUnitDamageHistory() {
	for _, ps := range e.State.Players {
		if ps == nil {
			continue
		}
		ps.FriendlyUnitDamagedLastTurn = ps.FriendlyUnitDamagedThisTurn
		ps.FriendlyUnitDamagedThisTurn = false
		ps.HeroDamageTakenLastTurn = ps.HeroDamageTakenThisTurn
		ps.HeroDamageTakenThisTurn = 0
	}
}
