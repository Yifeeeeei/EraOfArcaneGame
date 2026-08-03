package game

import "eraofarcane/model"

func isNegativeStatus(status string) bool {
	for _, candidate := range negativeStatuses {
		if candidate == status {
			return true
		}
	}
	return false
}

func (e *Engine) hasEffectiveStatus(card *CardInstance, status string) bool {
	if card == nil || card.Statuses[status] <= 0 {
		return false
	}
	return !e.negativeStatusIneffective(card, status)
}

func (e *Engine) addStatus(card *CardInstance, status string, amount int) bool {
	if card == nil || amount <= 0 {
		return false
	}
	if status == StatusWeaken && !canCardBeWeakened(card.Card) {
		return false
	}
	if isNegativeStatus(status) && e.rejectsNegativeStatusApplication(card, status) {
		return false
	}
	card.Statuses[status] += amount
	if status == StatusPetrify && card.Card != nil && card.Card.Number == "3611101" {
		e.refreshRedMoonState(card.OwnerID)
	}
	return true
}

func (e *Engine) rejectsNegativeStatusApplication(card *CardInstance, status string) bool {
	if card == nil || card.Card == nil {
		return false
	}
	return e.cardHasNegativeStatusImmunity(card, status)
}

func (e *Engine) negativeStatusIneffective(card *CardInstance, status string) bool {
	if card == nil || card.Card == nil || !isNegativeStatus(status) {
		return false
	}
	if card.Statuses[fireNegativeStatusImmunityUntil] >= e.State.TurnNumber && card.Card.Category == model.ElementFire {
		return true
	}
	if e.cardHasNegativeStatusImmunity(card, status) {
		return true
	}
	ps := e.State.Players[card.OwnerID]
	if ps != nil && card.Position != nil && ps.Shield > 0 && e.playerHasActiveCard(ps, "2411101") {
		return true
	}
	if ps == nil || card.Position == nil {
		return false
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			protector := ps.Units[col][row]
			if protector == nil || protector.Card == nil || protector.Position == nil {
				continue
			}
			behavior, ok := behaviorForNumber(protector.Card.Number).(AdjacentNegativeStatusProtectionBehavior)
			if !ok || !behavior.HasActiveAdjacentNegativeStatusProtection(protector) || !behavior.ProtectsAdjacentFromNegativeStatus() {
				continue
			}
			if abs(protector.Position.Col-card.Position.Col)+abs(protector.Position.Row-card.Position.Row) <= 1 {
				return true
			}
		}
	}
	return false
}

func (e *Engine) cardHasNegativeStatusImmunity(card *CardInstance, status string) bool {
	if card == nil || card.Card == nil {
		return false
	}
	behavior := behaviorForNumber(card.Card.Number)
	if immune, ok := behavior.(SpecificNegativeStatusImmunityBehavior); ok && immune.HasActiveNegativeStatusImmunity(card) && immune.HasNegativeStatusImmunity(status) {
		return true
	}
	if immune, ok := behavior.(NegativeStatusImmunityBehavior); ok && immune.HasActiveNegativeStatusImmunity(card) && immune.HasNegativeStatusImmunity() {
		return true
	}
	return false
}

func (e *Engine) resetCard(card *CardInstance) {
	if card == nil {
		return
	}
	if e.hasEffectiveStatus(card, StatusFreeze) {
		return
	}
	if card.Statuses[StatusCooldown] > 0 {
		return
	}
	card.IsHorizontal = false
}

func (e *Engine) canConsumeCard(card *CardInstance) bool {
	if card == nil || card.IsHorizontal {
		return false
	}
	if e.hasEffectiveStatus(card, StatusStun) {
		return false
	}
	if card.Statuses[StatusCooldown] > 0 {
		return false
	}
	return true
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
