package game

func masteryBehavior(card *CardInstance) (MasteryBehavior, bool) {
	if card == nil || card.Card == nil {
		return nil, false
	}
	behavior, ok := globalRegistry.GetBehavior(card.Card.Number).(MasteryBehavior)
	if !ok || !behavior.HasActiveMastery(card) {
		return nil, false
	}
	return behavior, true
}

func (e *Engine) advanceMastery(card *CardInstance, playerID int, amount int) {
	if amount <= 0 {
		return
	}
	behavior, ok := masteryBehavior(card)
	if !ok {
		return
	}
	maxLevel := behavior.MasteryMax()
	if maxLevel <= 0 {
		return
	}
	oldLevel := card.Statuses[StatusMastery]
	newLevel := oldLevel + amount
	if newLevel > maxLevel {
		newLevel = maxLevel
	}
	if newLevel <= oldLevel {
		return
	}
	card.Statuses[StatusMastery] = newLevel
	for level := oldLevel + 1; level <= newLevel; level++ {
		ctx := &EffectContext{
			Engine:     e,
			Source:     card,
			PlayerID:   playerID,
			OpponentID: 1 - playerID,
			ExtraData:  map[string]any{"mastery": level},
		}
		_ = behavior.OnMastery(ctx, level)
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"source":  cardToInfo(card),
			"effect":  "mastery",
			"mastery": level,
		}})
		e.triggerFieldEffectsWithData(TriggerOnMastery, playerID, card, map[string]any{
			"mastery":         level,
			"mastered_card":   card,
			"mastered_player": playerID,
		})
	}
}

func (e *Engine) advanceAllMasteryToMax(ps *PlayerState) {
	if ps == nil {
		return
	}
	for _, card := range e.getAllFieldCards(ps) {
		behavior, ok := masteryBehavior(card)
		if !ok {
			continue
		}
		e.advanceMastery(card, ps.PlayerID, behavior.MasteryMax()-card.Statuses[StatusMastery])
	}
}

func (e *Engine) settleMastery(ps *PlayerState) {
	// 精通 does not advance during mark settlement. It advances when the card
	// itself is consumed, then threshold effects resolve immediately.
}
