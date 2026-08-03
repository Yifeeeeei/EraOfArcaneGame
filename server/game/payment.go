package game

import "eraofarcane/model"

type paymentPurpose string

const (
	paymentPurposePlay  paymentPurpose = "play"
	paymentPurposeLearn paymentPurpose = "learn"
	paymentPurposeUse   paymentPurpose = "use"
)

func cloneElements(elements map[string]int) map[string]int {
	result := make(map[string]int)
	for _, elem := range model.AllElements {
		result[elem] = elements[elem]
	}
	return result
}

func calculateElementPayment(available map[string]int, cost map[string]int) (map[string]int, bool) {
	return calculateElementPaymentWithOptions(available, cost, false)
}

func calculateElementPaymentWithOptions(available map[string]int, cost map[string]int, lightWildcard bool, lightCostWildcardOptions ...bool) (map[string]int, bool) {
	lightCostWildcard := len(lightCostWildcardOptions) > 0 && lightCostWildcardOptions[0]
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
		if remainingCost[elem] > 0 && lightWildcard && elem != model.ElementLight {
			pay := min(remainingAvailable[model.ElementLight], remainingCost[elem])
			if pay > 0 {
				payment[model.ElementLight] += pay
				remainingAvailable[model.ElementLight] -= pay
				remainingCost[elem] -= pay
			}
		}
		if remainingCost[elem] > 0 && !(lightCostWildcard && elem == model.ElementLight) {
			return nil, false
		}
	}

	if lightCostWildcard && remainingCost[model.ElementLight] > 0 {
		for _, elem := range model.AllElements {
			if elem == model.ElementArcane || elem == model.ElementLight {
				continue
			}
			pay := min(remainingAvailable[elem], remainingCost[model.ElementLight])
			if pay > 0 {
				payment[elem] += pay
				remainingAvailable[elem] -= pay
				remainingCost[model.ElementLight] -= pay
			}
			if remainingCost[model.ElementLight] <= 0 {
				break
			}
		}
		if remainingCost[model.ElementLight] > 0 {
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

func calculateCardActionPaymentWithOptions(available map[string]int, card *CardInstance, ownCost map[string]int, totalCost map[string]int, purpose paymentPurpose, lightWildcard bool, lightCostWildcardOptions ...bool) (map[string]int, bool) {
	lightCostWildcard := len(lightCostWildcardOptions) > 0 && lightCostWildcardOptions[0]
	if !requiresDistinctOwnUseCost(card, purpose) {
		return calculateElementPaymentWithOptions(available, totalCost, lightWildcard, lightCostWildcard)
	}
	return calculateDistinctOwnCostPayment(available, ownCost, totalCost, lightWildcard, lightCostWildcard)
}

func calculateDistinctOwnCostPayment(available map[string]int, ownCost map[string]int, totalCost map[string]int, lightWildcard bool, lightCostWildcardOptions ...bool) (map[string]int, bool) {
	lightCostWildcard := len(lightCostWildcardOptions) > 0 && lightCostWildcardOptions[0]
	distinctCount := totalElementCost(ownCost)
	if distinctCount <= 0 {
		return calculateElementPaymentWithOptions(available, totalCost, lightWildcard, lightCostWildcard)
	}
	extraCost := subtractElementCosts(totalCost, ownCost)
	elements := model.AllElements
	ownPayment := make(map[string]int)
	var search func(int, int) (map[string]int, bool)
	search = func(start int, remaining int) (map[string]int, bool) {
		if remaining == 0 {
			if !validateElementPaymentWithOptions(available, ownCost, ownPayment, lightWildcard, lightCostWildcard) {
				return nil, false
			}
			remainingAvailable := cloneElements(available)
			for elem, amount := range ownPayment {
				remainingAvailable[elem] -= amount
			}
			extraPayment, ok := calculateElementPaymentWithOptions(remainingAvailable, extraCost, lightWildcard, lightCostWildcard)
			if !ok {
				return nil, false
			}
			payment := cloneElements(extraPayment)
			for elem, amount := range ownPayment {
				payment[elem] += amount
			}
			return payment, true
		}
		for i := start; i < len(elements); i++ {
			elem := elements[i]
			if available[elem] <= ownPayment[elem] {
				continue
			}
			ownPayment[elem]++
			if payment, ok := search(i+1, remaining-1); ok {
				return payment, true
			}
			ownPayment[elem]--
		}
		return nil, false
	}
	return search(0, distinctCount)
}

func validateElementPayment(available map[string]int, cost map[string]int, payment map[string]int) bool {
	return validateElementPaymentWithOptions(available, cost, payment, false)
}

func validateElementPaymentWithOptions(available map[string]int, cost map[string]int, payment map[string]int, lightWildcard bool, lightCostWildcardOptions ...bool) bool {
	lightCostWildcard := len(lightCostWildcardOptions) > 0 && lightCostWildcardOptions[0]
	for elem, amount := range payment {
		if amount < 0 || amount > available[elem] {
			return false
		}
	}

	spent, ok := calculateElementPaymentWithOptions(payment, cost, lightWildcard, lightCostWildcard)
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

func validateCardActionPaymentWithOptions(available map[string]int, card *CardInstance, ownCost map[string]int, totalCost map[string]int, purpose paymentPurpose, payment map[string]int, lightWildcard bool, lightCostWildcardOptions ...bool) bool {
	lightCostWildcard := len(lightCostWildcardOptions) > 0 && lightCostWildcardOptions[0]
	if !validateElementPaymentWithOptions(available, totalCost, payment, lightWildcard, lightCostWildcard) {
		return false
	}
	if !requiresDistinctOwnUseCost(card, purpose) {
		return true
	}
	return distinctOwnCostPaymentSatisfied(card, purpose, ownCost, totalCost, payment, lightWildcard, lightCostWildcard)
}

func strictPaymentRequirement(card *CardInstance, purpose paymentPurpose, cost map[string]int) map[string]int {
	if card == nil || card.Card == nil {
		return nil
	}
	switch card.Card.Number {
	case "1021112":
		if purpose == paymentPurposePlay {
			return map[string]int{model.ElementArcane: totalElementCost(cost)}
		}
	case "3011101":
		if purpose == paymentPurposeLearn || purpose == paymentPurposeUse {
			return map[string]int{model.ElementArcane: totalElementCost(cost)}
		}
	case "3411101":
		if purpose == paymentPurposeLearn || purpose == paymentPurposeUse {
			earth := cost[model.ElementEarth]
			arcane := totalElementCost(cost) - earth
			requirement := map[string]int{}
			if earth > 0 {
				requirement[model.ElementEarth] = earth
			}
			if arcane > 0 {
				requirement[model.ElementArcane] = arcane
			}
			return requirement
		}
	}
	return nil
}

func requiresDistinctOwnUseCost(card *CardInstance, purpose paymentPurpose) bool {
	return card != nil && card.Card != nil && card.Card.Number == "3021103" && purpose == paymentPurposeUse
}

func distinctOwnCostPaymentSatisfied(card *CardInstance, purpose paymentPurpose, ownCost map[string]int, totalCost map[string]int, payment map[string]int, lightWildcard bool, lightCostWildcardOptions ...bool) bool {
	lightCostWildcard := len(lightCostWildcardOptions) > 0 && lightCostWildcardOptions[0]
	if !requiresDistinctOwnUseCost(card, purpose) {
		return true
	}
	distinctCount := totalElementCost(ownCost)
	if distinctCount <= 0 {
		return true
	}
	extraCost := subtractElementCosts(totalCost, ownCost)
	elements := model.AllElements
	ownPayment := make(map[string]int)
	var search func(int, int) bool
	search = func(start int, remaining int) bool {
		if remaining == 0 {
			if !validateElementPaymentWithOptions(payment, ownCost, ownPayment, lightWildcard, lightCostWildcard) {
				return false
			}
			remainingPayment := cloneElements(payment)
			for elem, amount := range ownPayment {
				remainingPayment[elem] -= amount
			}
			return validateElementPaymentWithOptions(remainingPayment, extraCost, remainingPayment, lightWildcard, lightCostWildcard)
		}
		for i := start; i < len(elements); i++ {
			elem := elements[i]
			if payment[elem] <= ownPayment[elem] {
				continue
			}
			ownPayment[elem]++
			if search(i+1, remaining-1) {
				return true
			}
			ownPayment[elem]--
		}
		return false
	}
	return search(0, distinctCount)
}

func subtractElementCosts(total map[string]int, own map[string]int) map[string]int {
	result := make(map[string]int)
	for _, elem := range model.AllElements {
		if amount := total[elem] - own[elem]; amount > 0 {
			result[elem] = amount
		}
	}
	return result
}

func strictPaymentSatisfied(card *CardInstance, purpose paymentPurpose, strictCost map[string]int, payment map[string]int) bool {
	requirement := strictPaymentRequirement(card, purpose, strictCost)
	if len(requirement) == 0 {
		return true
	}
	for elem, amount := range requirement {
		if payment[elem] < amount {
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

func (e *Engine) availableElementsWithOverexert(ps *PlayerState, units []*CardInstance) map[string]int {
	available := cloneElements(ps.Elements)
	for _, unit := range units {
		for elem, amount := range e.effectiveElementsGain(unit) {
			available[elem] += amount
		}
	}
	return available
}

func canPayCostWithOverexert(ps *PlayerState, cost map[string]int, units []*CardInstance) bool {
	return canPayCostWithOverexertOptions(ps, cost, units, false)
}

func canPayCostWithOverexertOptions(ps *PlayerState, cost map[string]int, units []*CardInstance, lightWildcard bool) bool {
	_, ok := calculateElementPaymentWithOptions(availableElementsWithOverexert(ps, units), cost, lightWildcard)
	return ok
}

func (e *Engine) canPayCostWithOverexertOptions(ps *PlayerState, cost map[string]int, units []*CardInstance, lightWildcard bool) bool {
	_, ok := calculateElementPaymentWithOptions(e.availableElementsWithOverexert(ps, units), cost, lightWildcard, e.playerHasLightCostWildcard(ps))
	return ok
}

func payDefenseCost(ps *PlayerState, cost map[string]int, action ActionMessage, units []*CardInstance) bool {
	return payDefenseCostWithOptions(ps, cost, action, units, false)
}

func payDefenseCostWithOptions(ps *PlayerState, cost map[string]int, action ActionMessage, units []*CardInstance, lightWildcard bool) bool {
	available := availableElementsWithOverexert(ps, units)
	var payment map[string]int
	if explicit := paymentFromAction(action); explicit != nil {
		payment = explicit
		if !validateElementPaymentWithOptions(available, cost, payment, lightWildcard) {
			return false
		}
	} else {
		var ok bool
		payment, ok = calculateElementPaymentWithOptions(available, cost, lightWildcard)
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

func (e *Engine) payDefenseCostWithOptions(ps *PlayerState, cost map[string]int, action ActionMessage, units []*CardInstance, lightWildcard bool) bool {
	available := e.availableElementsWithOverexert(ps, units)
	payment := paymentFromAction(action)
	if payment == nil {
		var ok bool
		payment, ok = calculateElementPaymentWithOptions(available, cost, lightWildcard, e.playerHasLightCostWildcard(ps))
		if !ok {
			return false
		}
	} else if !validateElementPaymentWithOptions(available, cost, payment, lightWildcard, e.playerHasLightCostWildcard(ps)) {
		return false
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
