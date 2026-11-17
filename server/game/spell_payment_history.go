package game

import "eraofarcane/model"

type spellPaymentPart struct {
	card *CardInstance
	cost map[string]int
}

// Attribute a combined payment once. Exact colored requirements are paid
// first; remaining substitutions and wildcard costs follow main/boost order.
func (e *Engine) recordSpellUsePayments(ps *PlayerState, parts []spellPaymentPart, payment map[string]int, attribution ...map[string]map[string]int) {
	remaining := copyElements(payment)
	costs := make([]map[string]int, len(parts))
	assigned := map[string]map[string]int{}
	if len(attribution) > 0 {
		assigned = attribution[0]
	}
	for _, part := range parts {
		if own, ok := assigned[part.card.InstanceID]; ok {
			for elem, n := range own {
				remaining[elem] -= n
			}
		}
	}
	for i, part := range parts {
		if own, ok := assigned[part.card.InstanceID]; ok {
			part.card.usePayment = copyElements(own)
			continue
		}
		part.card.usePayment = map[string]int{}
		costs[i] = copyElements(part.cost)
		for _, elem := range model.AllElements {
			if elem == model.ElementArcane {
				continue
			}
			n := min(costs[i][elem], remaining[elem])
			if n > 0 {
				part.card.usePayment[elem] += n
				remaining[elem] -= n
				costs[i][elem] -= n
			}
		}
	}
	for _, colored := range []bool{true, false} {
		for i, part := range parts {
			if _, ok := assigned[part.card.InstanceID]; ok {
				continue
			}
			if (totalElementCost(costs[i]) > costs[i][model.ElementArcane]) != colored {
				continue
			}
			paid, ok := calculateElementPaymentWithOptions(remaining, costs[i], e.playerHasLightWildcard(ps), e.playerHasLightCostWildcard(ps))
			if !ok {
				continue
			} // Payment was validated as a whole before commit.
			for elem, n := range paid {
				if n > 0 {
					part.card.usePayment[elem] += n
					remaining[elem] -= n
				}
			}
		}
	}
}
