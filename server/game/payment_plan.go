package game

import (
	"eraofarcane/model"
)

func (e *Engine) payCostForAction(ps *PlayerState, cost map[string]int, action ActionMessage) bool {
	payment, strictArcane, ok := e.paymentPlanForAction(ps, cost, action)
	if !ok {
		return false
	}
	e.spendPaymentPlan(ps, payment, strictArcane)
	return true
}

func (e *Engine) paymentPlanForAction(ps *PlayerState, cost map[string]int, action ActionMessage) (map[string]int, int, bool) {
	if payment := paymentFromAction(action); payment != nil {
		strictArcane, ok := e.strictArcaneUsedByExplicitPayment(ps, cost, payment)
		return payment, strictArcane, ok
	}
	return e.calculatePaymentPlan(ps, cost)
}

func (e *Engine) calculatePaymentPlan(ps *PlayerState, cost map[string]int) (map[string]int, int, bool) {
	return e.calculatePaymentPlanFromAvailable(ps, ps.Elements, ps.StrictArcane, cost, e.playerHasLightWildcard(ps))
}

func (e *Engine) paymentPlanFromAvailableForAction(ps *PlayerState, available map[string]int, strictAvailable int, cost map[string]int, action ActionMessage, lightWildcard bool) (map[string]int, int, bool) {
	if payment := paymentFromAction(action); payment != nil {
		strictArcane, ok := e.strictArcaneUsedByExplicitPaymentFromAvailable(ps, available, strictAvailable, cost, payment, lightWildcard)
		return payment, strictArcane, ok
	}
	return e.calculatePaymentPlanFromAvailable(ps, available, strictAvailable, cost, lightWildcard)
}

func (e *Engine) calculatePaymentPlanFromAvailable(ps *PlayerState, available map[string]int, strictAvailable int, cost map[string]int, lightWildcard bool) (map[string]int, int, bool) {
	maxStrict := min(strictAvailable, cost[model.ElementArcane])
	for strictArcane := maxStrict; strictArcane >= 0; strictArcane-- {
		remainingCost := cloneElements(cost)
		remainingCost[model.ElementArcane] -= strictArcane
		payment, ok := calculateElementPaymentWithOptions(available, remainingCost, lightWildcard, e.playerHasLightCostWildcard(ps))
		if !ok {
			continue
		}
		payment[model.ElementArcane] += strictArcane
		return payment, strictArcane, true
	}
	return nil, 0, false
}

func (e *Engine) strictArcaneUsedByExplicitPayment(ps *PlayerState, cost map[string]int, payment map[string]int) (int, bool) {
	return e.strictArcaneUsedByExplicitPaymentFromAvailable(ps, ps.Elements, ps.StrictArcane, cost, payment, e.playerHasLightWildcard(ps))
}

func (e *Engine) strictArcaneUsedByExplicitPaymentFromAvailable(ps *PlayerState, available map[string]int, strictAvailable int, cost map[string]int, payment map[string]int, lightWildcard bool) (int, bool) {
	maxStrict := min(strictAvailable, min(cost[model.ElementArcane], payment[model.ElementArcane]))
	for strictArcane := maxStrict; strictArcane >= 0; strictArcane-- {
		normalPayment := cloneElements(payment)
		normalPayment[model.ElementArcane] -= strictArcane
		remainingCost := cloneElements(cost)
		remainingCost[model.ElementArcane] -= strictArcane
		if validateElementPaymentWithOptions(available, remainingCost, normalPayment, lightWildcard, e.playerHasLightCostWildcard(ps)) {
			return strictArcane, true
		}
	}
	return 0, false
}

func (e *Engine) spendPaymentPlan(ps *PlayerState, payment map[string]int, strictArcane int) {
	for elem, amount := range payment {
		poolAmount := amount
		if elem == model.ElementArcane {
			poolAmount -= strictArcane
		}
		if poolAmount > 0 {
			ps.Elements[elem] -= poolAmount
		}
	}
	if strictArcane > 0 {
		ps.StrictArcane -= strictArcane
	}
}

func (e *Engine) payCostForCardAction(ps *PlayerState, card *CardInstance, strictCost map[string]int, totalCost map[string]int, purpose paymentPurpose, action ActionMessage) bool {
	payment, strictArcane, ok := e.cardPaymentPlanForAction(ps, card, strictCost, totalCost, purpose, action)
	if !ok {
		return false
	}
	e.spendPaymentPlan(ps, payment, strictArcane)
	return true
}

func (e *Engine) cardPaymentPlanForAction(ps *PlayerState, card *CardInstance, strictCost map[string]int, totalCost map[string]int, purpose paymentPurpose, action ActionMessage) (map[string]int, int, bool) {
	if payment := paymentFromAction(action); payment != nil {
		strictArcane, ok := e.strictArcaneUsedByExplicitCardPayment(ps, card, strictCost, totalCost, purpose, payment)
		return payment, strictArcane, ok
	}
	maxStrict := min(ps.StrictArcane, totalCost[model.ElementArcane])
	for strictArcane := maxStrict; strictArcane >= 0; strictArcane-- {
		remainingTotalCost := cloneElements(totalCost)
		remainingTotalCost[model.ElementArcane] -= strictArcane
		payment, ok := calculateCardActionPaymentWithOptions(ps.Elements, card, strictCost, remainingTotalCost, purpose, e.playerHasLightWildcard(ps), e.playerHasLightCostWildcard(ps))
		if !ok {
			continue
		}
		payment[model.ElementArcane] += strictArcane
		if strictPaymentSatisfied(card, purpose, strictCost, payment) {
			return payment, strictArcane, true
		}
	}
	return nil, 0, false
}

func (e *Engine) strictArcaneUsedByExplicitCardPayment(ps *PlayerState, card *CardInstance, strictCost map[string]int, totalCost map[string]int, purpose paymentPurpose, payment map[string]int) (int, bool) {
	maxStrict := min(ps.StrictArcane, min(totalCost[model.ElementArcane], payment[model.ElementArcane]))
	for strictArcane := maxStrict; strictArcane >= 0; strictArcane-- {
		normalPayment := cloneElements(payment)
		normalPayment[model.ElementArcane] -= strictArcane
		remainingTotalCost := cloneElements(totalCost)
		remainingTotalCost[model.ElementArcane] -= strictArcane
		if !validateCardActionPaymentWithOptions(ps.Elements, card, strictCost, remainingTotalCost, purpose, normalPayment, e.playerHasLightWildcard(ps), e.playerHasLightCostWildcard(ps)) {
			continue
		}
		if strictPaymentSatisfied(card, purpose, strictCost, payment) {
			return strictArcane, true
		}
	}
	return 0, false
}

func (e *Engine) canPayCostForAction(ps *PlayerState, cost map[string]int, action ActionMessage) bool {
	_, _, ok := e.paymentPlanForAction(ps, cost, action)
	return ok
}

func (e *Engine) canPayCost(ps *PlayerState, cost map[string]int) bool {
	_, _, ok := e.calculatePaymentPlan(ps, cost)
	return ok
}

func (e *Engine) canPayCostForCardAction(ps *PlayerState, card *CardInstance, strictCost map[string]int, totalCost map[string]int, purpose paymentPurpose, action ActionMessage) bool {
	_, _, ok := e.cardPaymentPlanForAction(ps, card, strictCost, totalCost, purpose, action)
	return ok
}

func (e *Engine) playerHasLightWildcard(ps *PlayerState) bool {
	return e.playerPaymentRules(ps).LightPaysAny
}

func (e *Engine) playerHasLightCostWildcard(ps *PlayerState) bool {
	return e.playerPaymentRules(ps).AnyPaysLight
}

func paymentFromAction(action ActionMessage) map[string]int {
	raw, ok := action.Data["payment"]
	if !ok || raw == nil {
		return nil
	}
	payment := make(map[string]int)
	switch values := raw.(type) {
	case map[string]any:
		for elem, value := range values {
			switch amount := value.(type) {
			case float64:
				payment[elem] = int(amount)
			case int:
				payment[elem] = amount
			}
		}
	case map[string]int:
		for elem, amount := range values {
			payment[elem] = amount
		}
	}
	return payment
}
