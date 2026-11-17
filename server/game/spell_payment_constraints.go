package game

import (
	"eraofarcane/model"
	"sort"
	"strconv"
)

// Resolve actual per-card payments before commit. Strict and distinct costs
// belong to each participating card, including defense and boost cards.
func (e *Engine) constrainedSpellPayments(ps *PlayerState, parts []spellPaymentPart, payment map[string]int, explicit map[string]map[string]int, allowRemainder ...bool) (map[string]map[string]int, bool) {
	type constrainedPart struct {
		part     spellPaymentPart
		required map[string]int
		distinct bool
	}
	work := []constrainedPart{}
	ordinaryCost := map[string]int{}
	known := map[string]bool{}
	for _, part := range parts {
		known[part.card.InstanceID] = true
		source := physicalSpellSource(part.card)
		required := strictPaymentRequirement(source, paymentPurposeUse, part.cost)
		distinct := requiresDistinctOwnUseCost(source, paymentPurposeUse)
		_, chosen := explicit[part.card.InstanceID]
		if chosen || totalElementCost(required) > 0 || distinct {
			work = append(work, constrainedPart{part, required, distinct})
		} else {
			ordinaryCost = mergeElementCosts(ordinaryCost, part.cost)
		}
	}
	for id := range explicit {
		if !known[id] {
			return nil, false
		}
	}
	sort.SliceStable(work, func(i, j int) bool { return totalElementCost(work[i].required) > totalElementCost(work[j].required) })
	light, lightCost := e.playerHasLightWildcard(ps), e.playerHasLightCostWildcard(ps)
	result := map[string]map[string]int{}
	failed := map[string]bool{}
	var search func(int, map[string]int) bool
	search = func(index int, available map[string]int) bool {
		if index == len(work) {
			if len(allowRemainder) > 0 && allowRemainder[0] {
				_, ok := calculateElementPaymentWithOptions(available, ordinaryCost, light, lightCost)
				return ok
			}
			return validateElementPaymentWithOptions(available, ordinaryCost, available, light, lightCost)
		}
		key := strconv.Itoa(index)
		for _, elem := range model.AllElements {
			key += "/" + strconv.Itoa(available[elem])
		}
		if failed[key] {
			return false
		}
		entry := work[index]
		part := entry.part
		try := func(own map[string]int) bool {
			if !validateElementPaymentWithOptions(available, part.cost, own, light, lightCost) || !strictPaymentSatisfied(physicalSpellSource(part.card), paymentPurposeUse, part.cost, own) {
				return false
			}
			if entry.distinct {
				for _, amount := range own {
					if amount > 1 {
						return false
					}
				}
			}
			rest := copyElements(available)
			for elem, n := range own {
				rest[elem] -= n
			}
			if !search(index+1, rest) {
				return false
			}
			result[part.card.InstanceID] = copyElements(own)
			return true
		}
		if own, chosen := explicit[part.card.InstanceID]; chosen {
			if try(own) {
				return true
			}
		} else if totalElementCost(entry.required) == totalElementCost(part.cost) && len(entry.required) > 0 {
			if try(entry.required) {
				return true
			}
		} else {
			// Enumerate the constrained card's own payment. Ordinary cards
			// are checked together after these constraints are satisfied.
			own := map[string]int{}
			var choose func(int, int) bool
			choose = func(start, remaining int) bool {
				if remaining == 0 {
					return try(own)
				}
				if start == len(model.AllElements) {
					return false
				}
				elem := model.AllElements[start]
				maximum := min(available[elem], remaining)
				if entry.distinct {
					maximum = min(maximum, 1)
				}
				for n := maximum; n >= 0; n-- {
					own[elem] = n
					if choose(start+1, remaining-n) {
						return true
					}
				}
				delete(own, elem)
				return false
			}
			if choose(0, totalElementCost(part.cost)) {
				return true
			}
		}
		failed[key] = true
		return false
	}
	if !search(0, copyElements(payment)) {
		return nil, false
	}
	// Keep records compact; zero amounts have no payment meaning.
	for _, own := range result {
		for elem, n := range own {
			if n == 0 {
				delete(own, elem)
			}
		}
	}
	return result, true
}
