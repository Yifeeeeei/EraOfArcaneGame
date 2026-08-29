package game

func triggeredTurnAvailable(card *CardInstance) bool {
	return card != nil && card.UsedThisTurn+card.PendingTriggeredUses < perTurnLimit(card)
}

func useTriggeredTurn(card *CardInstance) bool {
	if !triggeredTurnAvailable(card) {
		return false
	}
	card.UsedThisTurn++
	return true
}

func reserveTriggeredTurn(card *CardInstance) bool {
	if !triggeredTurnAvailable(card) {
		return false
	}
	card.PendingTriggeredUses++
	return true
}

func resolveTriggeredTurn(card *CardInstance, use bool) bool {
	if card == nil || card.PendingTriggeredUses <= 0 {
		return false
	}
	card.PendingTriggeredUses--
	if !use || card.UsedThisTurn >= perTurnLimit(card) {
		return false
	}
	card.UsedThisTurn++
	return true
}
