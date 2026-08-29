package game

func triggeredTurnAvailable(card *CardInstance) bool {
	return card != nil && card.UsedThisTurn < perTurnLimit(card)
}

func useTriggeredTurn(card *CardInstance) bool {
	if !triggeredTurnAvailable(card) {
		return false
	}
	card.UsedThisTurn++
	return true
}

func (e *Engine) SetTriggeredTurnAction(source *CardInstance, playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, callback func([]string)) {
	if !triggeredTurnAvailable(source) {
		return
	}
	e.setPendingActionWithOptions(playerID, actionType, prompt, candidates, minSelect, maxSelect, nil, false, callback, nil, nil, nil, func() bool {
		return triggeredTurnAvailable(source)
	})
}
