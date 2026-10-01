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
	e.runResolution("end turn",
		func() {
			// Collect the same batch as before; a choice pauses the next step,
			// not collection of other triggers from this turn-end event.
			for _, card := range e.getAllFieldCards(ps) {
				e.triggerEffects(TriggerOnTurnEnd, card, nil, nil)
			}
			e.triggerFieldEffectsWithData(TriggerOnTurnEnd, 1-ps.PlayerID, nil, map[string]any{"ended_player": ps.PlayerID})
		},
		func() { e.applyOpponentTurnEndTemporaryModifiers(ps.PlayerID) },
		func() { e.finishEndTurnAfterOpponentTemp(ps) },
	)
}

func (e *Engine) finishEndTurnAfterOpponentTemp(ps *PlayerState) {
	e.discardMarkedEndOfTurnCards(ps)
	e.applyLoadGainAtTurnEnd(ps)

	e.clearTemporaryModifiersAtTurnEnd(ps.PlayerID)
	e.clearExpiredTemporaryModifiers(ps.PlayerID)
	e.processAbilityDurations(ps)

	// Remove 临时 (temporary) units before the cleanup/reset steps.
	e.destroyAndisGiftDoomedUnits(ps)
	e.HandleTemporaryUnits(ps)

	// Discard phase has already happened above. The cleanup order is:
	// reset cards first, then settle marks. A skill with 冷却1 therefore remains
	// horizontal through the next turn, because 冷却 blocks this reset before it
	// is removed by mark settlement.
	e.resetCards(ps)

	// Process status marks (点燃, 冻结, 冷却, etc.) after reset.
	e.processEndOfTurnStatuses(ps)

	// Decay 护盾 and 隐蔽 as part of mark settlement.
	e.HandleShieldDecay(ps)

	// 精通 is a card-instance mark. It advances during the unified mark
	// settlement step, after reset has already considered 冷却/冻结/etc.
	e.settleMastery(ps)

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
	allCards := e.getAllFieldCards(ps)
	changes := []CardStateChange{}

	for _, card := range allCards {
		// 点燃: remove 1 stack, deal 1 fire damage
		if card.Statuses[StatusBurn] > 0 {
			effective := e.hasEffectiveStatus(card, StatusBurn)
			card.Statuses[StatusBurn]--
			if effective {
				e.ApplyDamage(DamageRequest{Target: card, Amount: 1, Status: StatusBurn})
			}
		}
		// 冻结: remove 1 stack
		if card.Statuses[StatusFreeze] > 0 {
			card.Statuses[StatusFreeze]--
		}
		// 眩晕: remove 1 stack
		if card.Statuses[StatusStun] > 0 {
			card.Statuses[StatusStun]--
		}
		// 石化: remove 1 stack
		if card.Statuses[StatusPetrify] > 0 {
			before := card.Statuses[StatusPetrify]
			card.Statuses[StatusPetrify]--
			if card.Statuses[StatusPetrify] <= 0 {
				changes = append(changes, CardStateChange{Card: card, Status: StatusPetrify, Before: before, After: 0})
			}
		}
		// 冷却: remove 1 stack
		if card.Statuses[StatusCooldown] > 0 {
			card.Statuses[StatusCooldown]--
		}
	}

	// 虚弱 is on skills, handled separately
	for i := range ps.Skills {
		if ps.Skills[i] != nil && ps.Skills[i].Statuses[StatusWeaken] > 0 {
			ps.Skills[i].Statuses[StatusWeaken]--
		}
		if ps.Skills[i] != nil && ps.Skills[i].Statuses[StatusSeal] > 0 {
			ps.Skills[i].Statuses[StatusSeal]--
		}
	}
	e.notifyCardStateChanges(changes...)
}
