package game

import (
	"eraofarcane/model"
	"fmt"
)

// handleConsume handles consuming a card (横置 to gain elements)
func (e *Engine) elementsGainedFromConsume(playerID int, card *CardInstance, action ActionMessage) (map[string]int, error) {
	gains := e.effectiveElementsGain(card)
	if !e.isFirstPlayerFirstTurnHeroConsume(playerID, card) {
		return gains, nil
	}

	total := 0
	positiveElements := 0
	lastPositiveElement := ""
	for elem, amount := range gains {
		if amount <= 0 {
			continue
		}
		total += amount
		positiveElements++
		lastPositiveElement = elem
	}
	if total <= 1 {
		return gains, nil
	}

	limit := (total + 1) / 2
	if positiveElements == 1 {
		return map[string]int{lastPositiveElement: limit}, nil
	}

	selected := elementMapFromAction(action, "gain")
	if selected == nil {
		selected = elementMapFromAction(action, "gained")
	}
	if selected == nil {
		return nil, fmt.Errorf("first turn hero load requires choosing %d elements to gain", limit)
	}

	selectedTotal := 0
	for elem, amount := range selected {
		if amount < 0 || amount > gains[elem] {
			return nil, fmt.Errorf("invalid first turn hero load choice")
		}
		selectedTotal += amount
	}
	if selectedTotal != limit {
		return nil, fmt.Errorf("first turn hero load must gain exactly %d elements", limit)
	}
	return selected, nil
}

func (e *Engine) isFirstPlayerFirstTurnHeroConsume(playerID int, card *CardInstance) bool {
	if card == nil {
		return false
	}
	ps := e.State.Players[playerID]
	return e.State.IsFirstTurn && playerID == e.State.FirstPlayer && ps != nil && card == ps.Hero
}

func elementMapFromAction(action ActionMessage, key string) map[string]int {
	raw, ok := action.Data[key]
	if !ok || raw == nil {
		return nil
	}
	result := make(map[string]int)
	switch values := raw.(type) {
	case map[string]any:
		for elem, value := range values {
			switch amount := value.(type) {
			case float64:
				result[elem] = int(amount)
			case int:
				result[elem] = amount
			}
		}
	case map[string]int:
		for elem, amount := range values {
			result[elem] = amount
		}
	default:
		return nil
	}
	return result
}

func (e *Engine) handleConsume(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain && e.State.Phase != PhaseDefenseWindow {
		return fmt.Errorf("cannot consume now")
	}
	// During defense window, only the defending player can consume (透支)
	if e.State.Phase == PhaseDefenseWindow {
		if e.State.PendingSpell == nil || playerID == e.State.PendingSpell.AttackerID {
			return fmt.Errorf("only defender can overdraft during defense")
		}
	}
	if e.State.Phase == PhaseMain && e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	ps := e.State.Players[playerID]

	// Find the card on the field
	card := e.findCardOnField(ps, instanceID)
	if card == nil {
		return fmt.Errorf("card not found on field")
	}
	if !e.canConsumeCard(card) {
		return fmt.Errorf("card cannot be consumed")
	}

	gains, err := e.elementsGainedFromConsume(playerID, card, action)
	if err != nil {
		return err
	}

	// Consume: set horizontal and gain elements
	card.IsHorizontal = true
	strictArcaneGained := 0
	if providesStrictArcaneOnly(card) {
		strictArcaneGained = gains[model.ElementArcane]
		gains = cloneElements(gains)
		gains[model.ElementArcane] = 0
	}
	ps.GainElements(gains)
	if strictArcaneGained > 0 {
		e.gainStrictArcane(playerID, strictArcaneGained)
	}

	e.emit(GameEvent{
		Type:   "consume",
		Player: -1,
		Data: map[string]any{
			"player":               playerID,
			"instance_id":          instanceID,
			"elements":             ps.Elements,
			"gained":               gains,
			"strict_arcane_gained": strictArcaneGained,
		},
	})

	// Trigger 消耗 effects
	e.triggerEffects(TriggerOnConsume, card, nil, map[string]any{
		"gained": gains,
	})
	e.triggerFieldEffectsWithData(TriggerOnConsume, playerID, card, map[string]any{
		"consumed_player": playerID,
		"gained":          gains,
	})
	e.triggerFieldEffectsWithData(TriggerOnConsume, 1-playerID, card, map[string]any{
		"consumed_player": playerID,
		"gained":          gains,
	})
	e.advanceMastery(card, playerID, 1)
	e.destroyFuyeDoomedCardAfterExert(card)

	return nil
}
