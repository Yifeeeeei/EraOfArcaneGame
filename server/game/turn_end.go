package game

import (
	"fmt"
)

// handleEndTurn handles ending the current turn
func (e *Engine) handleEndTurn(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	e.endTurn()
	return nil
}

// endTurn processes end-of-turn effects and switches to next player
func (e *Engine) endTurn() {
	ps := e.State.Players[e.State.CurrentTurn]

	// Discard to hand limit
	if e.promptDiscardToHandLimit(e.State.CurrentTurn, func() {
		e.finishEndTurn(ps)
	}) {
		return // Wait for player to choose
	}

	e.finishEndTurn(ps)
}

func (e *Engine) promptDiscardToHandLimit(playerID int, afterDiscard func()) bool {
	if playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return false
	}
	handLimit := e.handLimitForPlayer(ps)
	if len(ps.Hand) <= handLimit {
		return false
	}

	discardCount := len(ps.Hand) - handLimit
	candidates := make([]map[string]any, len(ps.Hand))
	for i, c := range ps.Hand {
		candidates[i] = cardToInfo(c)
	}
	e.SetPendingAction(playerID, "discard",
		fmt.Sprintf("弃牌至手牌上限（需弃%d张）", discardCount),
		candidates, discardCount, discardCount,
		func(selected []string) {
			toDiscard := make(map[string]bool)
			for _, id := range selected {
				toDiscard[id] = true
			}
			remaining := make([]*CardInstance, 0, len(ps.Hand)-len(selected))
			for _, c := range ps.Hand {
				if toDiscard[c.InstanceID] {
					e.discardHandCardToGraveyard(playerID, c)
				} else {
					remaining = append(remaining, c)
				}
			}
			ps.Hand = remaining
			if afterDiscard != nil {
				afterDiscard()
			}
		})
	return true
}

func (e *Engine) shouldImmediatelyEnforceHandLimit(playerID int) bool {
	if playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	opponentID := 1 - playerID
	if opponentID < 0 || opponentID >= len(e.State.Players) {
		return false
	}
	for _, source := range e.getAllFieldCards(e.State.Players[opponentID]) {
		if source == nil || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		if b, ok := cardBehavior(source).(HandLimitEnforcementBehavior); ok && b.EnforcesHandLimit(source, playerID) {
			return true
		}
	}
	return false
}

// finishEndTurn completes end-of-turn processing (after optional discard)
func (e *Engine) finishEndTurn(ps *PlayerState) {
	e.State.Phase = PhaseTurnEnd
	e.runResolution("end turn",
		func() { e.resetCardOrientations(ps) },
		func() { e.processEndOfTurnStatuses(ps) },
		func() { e.HandleShieldDecay(ps); e.settleTurnDeaths() },
		func() {
			// Collect the same batch as before; a choice pauses the next step,
			// not collection of other triggers from this turn-end event.
			for _, card := range e.getAllFieldCards(ps) {
				e.triggerEffects(TriggerOnTurnEnd, card, nil, nil)
			}
			e.triggerFieldEffectsWithData(TriggerOnTurnEnd, 1-ps.PlayerID, nil, map[string]any{"ended_player": ps.PlayerID})
		},
		func() { e.settleTurnDeaths() },
		func() { e.applyOpponentTurnEndTemporaryModifiers(ps.PlayerID) },
		func() { e.settleTurnDeaths() },
		func() { e.finishEndTurnAfterOpponentTemp(ps) },
		func() { e.settleTurnDeaths() },
		func() { e.advanceTurnAfterCleanup(ps) },
	)
}

func (e *Engine) finishEndTurnAfterOpponentTemp(ps *PlayerState) {
	e.discardMarkedEndOfTurnCards(ps)
	e.applyLoadGainAtTurnEnd(ps)

	e.clearTemporaryModifiersAtTurnEnd(ps.PlayerID)
	e.clearExpiredTemporaryModifiers(ps.PlayerID)
	e.processAbilityDurations(ps)

	// Until-end effects last through reset, marks, and turn-end triggers.
	e.destroyAndisGiftDoomedUnits(ps)
	e.HandleTemporaryUnits(ps)
}

func (e *Engine) advanceTurnAfterCleanup(ps *PlayerState) {
	e.refreshTurnAbilityUses(ps)

	// Clear elements
	for elem := range ps.Elements {
		ps.Elements[elem] = 0
	}
	clearSpellCastTracking(ps)
	rollSpellHitTracking(ps)
	e.clearGraveyardTurnTracking()
	e.rollFriendlyUnitDamageHistory()
	for _, player := range e.State.Players {
		if player != nil {
			kept := player.TempModifiers[:0]
			for _, modifier := range player.TempModifiers {
				if modifier.Type != TempModLavaArmorYeYanShieldBreak {
					kept = append(kept, modifier)
				}
			}
			player.TempModifiers = kept
			player.ShieldBrokenThisTurn = false
		}
	}

	e.emit(GameEvent{
		Type:   "turn_end",
		Player: -1,
		Data: map[string]any{
			"player": e.State.CurrentTurn,
		},
	})

	// Switch turns
	if e.State.IsFirstTurn && e.State.CurrentTurn == e.State.FirstPlayer {
		e.State.IsFirstTurn = false
	}
	e.State.CurrentTurn = 1 - e.State.CurrentTurn
	if e.State.CurrentTurn == e.State.FirstPlayer {
		e.State.TurnNumber++
	}

	if e.State.Phase != PhaseGameOver {
		e.startTurn()
	}
}

// Finish a phase's complete death batch before another phase or player draws.
func (e *Engine) settleTurnDeaths() {
	depth := e.resolutionDepth
	e.resolutionDepth = 0
	defer func() { e.resolutionDepth = depth }()
	e.resolvePendingDeaths()
	e.checkWinCondition()
}

func (e *Engine) processAbilityDurations(ps *PlayerState) {
	changes := []CardStateChange{}
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil || card.Statuses[StatusAbilityDuration] <= 0 {
			continue
		}
		before := card.Statuses[StatusAbilityDuration]
		card.Statuses[StatusAbilityDuration]--
		if card.Statuses[StatusAbilityDuration] <= 0 {
			delete(card.Statuses, StatusAbilityDuration)
			changes = append(changes, CardStateChange{Card: card, Status: StatusAbilityDuration, Before: before, After: 0})
		}
	}
	e.notifyCardStateChanges(changes...)
}

// processEndOfTurnStatuses processes status marks at end of turn
func (e *Engine) processEndOfTurnStatuses(ps *PlayerState) {
	steps := []func(){}
	for _, card := range e.getAllFieldCards(ps) {
		card, visit := card, fieldVisit(card)
		present := func() bool { return e.isBattlefieldInstance(card) && card.fieldLife == visit }
		steps = append(steps, func() {
			if present() && card.Statuses[StatusBurn] > 0 {
				effective := e.hasEffectiveStatus(card, StatusBurn)
				card.Statuses[StatusBurn]--
				if effective {
					e.ApplyDamage(DamageRequest{Target: card, Amount: 1, Status: StatusBurn})
				}
			}
		}, e.settleTurnDeaths, func() {
			if !present() {
				return
			}
			var changes []CardStateChange
			for _, status := range []string{StatusFreeze, StatusStun, StatusPetrify, StatusCooldown, StatusWeaken, StatusSeal} {
				if before := card.Statuses[status]; before > 0 {
					card.Statuses[status]--
					changes = append(changes, CardStateChange{Card: card, Status: status, Before: before, After: before - 1})
				}
			}
			e.notifyCardStateChanges(changes...)
		})
	}
	e.runResolution("mark settlement", steps...)
}
