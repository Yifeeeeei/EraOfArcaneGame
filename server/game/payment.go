package game

import "eraofarcane/model"

func cloneElements(elements map[string]int) map[string]int {
	result := make(map[string]int)
	for _, elem := range model.AllElements {
		result[elem] = elements[elem]
	}
	return result
}

func calculateElementPayment(available map[string]int, cost map[string]int) (map[string]int, bool) {
	payment := make(map[string]int)
	remainingAvailable := cloneElements(available)
	remainingCost := make(map[string]int)

	for elem, amount := range cost {
		if amount <= 0 {
			continue
		}
		remainingCost[elem] += amount
	}

	for _, elem := range model.AllElements {
		if elem == model.ElementArcane {
			continue
		}
		amount := remainingCost[elem]
		if amount <= 0 {
			continue
		}

		pay := min(remainingAvailable[elem], amount)
		if pay > 0 {
			payment[elem] += pay
			remainingAvailable[elem] -= pay
			remainingCost[elem] -= pay
		}
	}

	for _, elem := range model.AllElements {
		if elem == model.ElementArcane {
			continue
		}
		amount := remainingCost[elem]
		if amount <= 0 {
			continue
		}

		pay := min(remainingAvailable[model.ElementArcane], amount)
		if pay > 0 {
			payment[model.ElementArcane] += pay
			remainingAvailable[model.ElementArcane] -= pay
			remainingCost[elem] -= pay
		}
		if remainingCost[elem] > 0 {
			return nil, false
		}
	}

	arcaneCost := remainingCost[model.ElementArcane]
	if arcaneCost > 0 {
		pay := min(remainingAvailable[model.ElementArcane], arcaneCost)
		if pay > 0 {
			payment[model.ElementArcane] += pay
			remainingAvailable[model.ElementArcane] -= pay
			arcaneCost -= pay
		}
	}

	for _, elem := range model.AllElements {
		if arcaneCost <= 0 {
			break
		}
		if elem == model.ElementArcane {
			continue
		}
		pay := min(remainingAvailable[elem], arcaneCost)
		if pay > 0 {
			payment[elem] += pay
			remainingAvailable[elem] -= pay
			arcaneCost -= pay
		}
	}
	if arcaneCost > 0 {
		return nil, false
	}

	return payment, true
}

func validateElementPayment(available map[string]int, cost map[string]int, payment map[string]int) bool {
	for elem, amount := range payment {
		if amount < 0 || amount > available[elem] {
			return false
		}
	}

	spent, ok := calculateElementPayment(payment, cost)
	if !ok {
		return false
	}
	for _, elem := range model.AllElements {
		if spent[elem] != payment[elem] {
			return false
		}
	}
	return true
}

func validateSingleElementPayment(available map[string]int, cost map[string]int, action ActionMessage) bool {
	var payment map[string]int
	if explicit := paymentFromAction(action); explicit != nil {
		payment = explicit
		if !validateElementPayment(available, cost, payment) {
			return false
		}
	} else {
		var ok bool
		payment, ok = calculateElementPayment(available, cost)
		if !ok {
			return false
		}
	}
	used := 0
	for _, amount := range payment {
		if amount > 0 {
			used++
		}
	}
	return used <= 1
}

func availableElementsWithOverexert(ps *PlayerState, units []*CardInstance) map[string]int {
	available := cloneElements(ps.Elements)
	for _, unit := range units {
		for elem, amount := range effectiveElementsGain(unit) {
			available[elem] += amount
		}
	}
	return available
}

func canPayCostWithOverexert(ps *PlayerState, cost map[string]int, units []*CardInstance) bool {
	_, ok := calculateElementPayment(availableElementsWithOverexert(ps, units), cost)
	return ok
}

func payDefenseCost(ps *PlayerState, cost map[string]int, action ActionMessage, units []*CardInstance) bool {
	available := availableElementsWithOverexert(ps, units)
	var payment map[string]int
	if explicit := paymentFromAction(action); explicit != nil {
		payment = explicit
		if !validateElementPayment(available, cost, payment) {
			return false
		}
	} else {
		var ok bool
		payment, ok = calculateElementPayment(available, cost)
		if !ok {
			return false
		}
	}

	for elem, amount := range payment {
		spendFromPool := min(ps.Elements[elem], amount)
		ps.Elements[elem] -= spendFromPool
	}
	for _, unit := range units {
		unit.IsHorizontal = true
	}
	return true
}
