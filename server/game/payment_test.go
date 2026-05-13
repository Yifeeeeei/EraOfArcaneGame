package game

import (
	"testing"

	"eraofarcane/model"
)

func newPaymentPlayer(elements map[string]int) *PlayerState {
	ps := NewPlayerState(0, "P", &model.Deck{})
	for elem, amount := range elements {
		ps.Elements[elem] = amount
	}
	return ps
}

func TestSpecificCostUsesMatchingElementBeforeArcane(t *testing.T) {
	ps := newPaymentPlayer(map[string]int{
		model.ElementFire:   1,
		model.ElementArcane: 1,
	})

	if !ps.PayCost(map[string]int{model.ElementFire: 1}) {
		t.Fatalf("expected fire cost to be payable")
	}
	if ps.Elements[model.ElementFire] != 0 {
		t.Fatalf("expected matching fire to be spent first, got %v", ps.Elements)
	}
	if ps.Elements[model.ElementArcane] != 1 {
		t.Fatalf("expected arcane to remain, got %v", ps.Elements)
	}
}

func TestSpecificCostCanUseArcaneAsWildcard(t *testing.T) {
	ps := newPaymentPlayer(map[string]int{model.ElementArcane: 1})

	if !ps.PayCost(map[string]int{model.ElementFire: 1}) {
		t.Fatalf("expected arcane to pay fire cost")
	}
	if ps.Elements[model.ElementArcane] != 0 {
		t.Fatalf("expected arcane to be spent, got %v", ps.Elements)
	}
}

func TestArcaneCostCanUseAnyElement(t *testing.T) {
	ps := newPaymentPlayer(map[string]int{
		model.ElementFire:  1,
		model.ElementWater: 1,
	})

	if !ps.PayCost(map[string]int{model.ElementArcane: 2}) {
		t.Fatalf("expected any elements to pay arcane cost")
	}
	if ps.Elements[model.ElementFire]+ps.Elements[model.ElementWater] != 0 {
		t.Fatalf("expected both elements to be spent, got %v", ps.Elements)
	}
}

func TestCanPayCostDoesNotMutateElements(t *testing.T) {
	ps := newPaymentPlayer(map[string]int{
		model.ElementFire:   1,
		model.ElementArcane: 1,
	})

	if !ps.CanPayCost(map[string]int{model.ElementFire: 2}) {
		t.Fatalf("expected cost to be payable")
	}
	if ps.Elements[model.ElementFire] != 1 || ps.Elements[model.ElementArcane] != 1 {
		t.Fatalf("CanPayCost should not mutate elements, got %v", ps.Elements)
	}
}

func TestCannotUseSpecificElementAsDifferentSpecificElement(t *testing.T) {
	ps := newPaymentPlayer(map[string]int{model.ElementWater: 1})

	if ps.CanPayCost(map[string]int{model.ElementFire: 1}) {
		t.Fatalf("water should not pay fire-specific cost")
	}
	if ps.PayCost(map[string]int{model.ElementFire: 1}) {
		t.Fatalf("water should not be spent for fire-specific cost")
	}
	if ps.Elements[model.ElementWater] != 1 {
		t.Fatalf("failed payment should not mutate elements, got %v", ps.Elements)
	}
}
