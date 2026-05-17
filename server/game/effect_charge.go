package game

// Legacy charge/overload helpers.
// Do not use this for 精通. Mastery is a per-card instance mark handled in
// mastery.go through CardInstance.Statuses[StatusMastery].

// addCharge adds charge counters to a player's charge pool
func (e *Engine) addCharge(playerID int, amount int) {
	ps := e.State.Players[playerID]
	ps.Charge += amount

	e.emit(GameEvent{
		Type:   "charge_change",
		Player: -1,
		Data: map[string]any{
			"player": playerID,
			"amount": amount,
			"total":  ps.Charge,
		},
	})

	// Check for overload triggers on all field cards
	e.checkOverload(playerID)
}

// removeCharge removes charge counters from a player's charge pool
func (e *Engine) removeCharge(playerID int, amount int) bool {
	ps := e.State.Players[playerID]
	if ps.Charge < amount {
		return false
	}
	ps.Charge -= amount

	e.emit(GameEvent{
		Type:   "charge_change",
		Player: -1,
		Data: map[string]any{
			"player": playerID,
			"amount": -amount,
			"total":  ps.Charge,
		},
	})

	return true
}

// checkOverload checks if any overload (过载) thresholds are met
func (e *Engine) checkOverload(playerID int) {
	ps := e.State.Players[playerID]
	allCards := e.getAllFieldCards(ps)

	for _, card := range allCards {
		behavior, ok := globalRegistry.GetBehavior(card.Card.Number).(OverloadBehavior)
		if !ok || !behavior.HasActiveOverload(card) {
			continue
		}
		threshold := behavior.OverloadThreshold()
		if threshold > 0 && ps.Charge >= threshold {
			// Trigger overload effect
			e.triggerEffects(TriggerOnEnter, card, nil, map[string]any{
				"overload": true,
			})
		}
	}
}
