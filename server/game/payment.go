package game

import "eraofarcane/model"

func calculateElementPayment(available map[string]int, cost map[string]int) (map[string]int, bool) {
	payment := make(map[string]int)
	remainingAvailable := make(map[string]int)
	remainingCost := make(map[string]int)

	for _, elem := range model.AllElements {
		remainingAvailable[elem] = available[elem]
	}
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
