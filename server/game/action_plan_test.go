package game

import (
	"eraofarcane/model"
	"reflect"
	"testing"
)

func TestRejectedSummonKeepsDevourCostsAndPayment(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	card := NewCardInstance(baseCard(t, "1621010"), 0, e.State.TurnNumber)
	p.Hand = []*CardInstance{card}
	p.Elements[model.ElementShadow] = 10
	food := placeUnit(baseCard(t, "1021001"), 0, 0, 0, e)
	setElementsGain(food, map[string]int{model.ElementShadow: 4})
	occupant := placeUnit(baseCard(t, "1021002"), 0, 1, 0, e)
	before := e.GetStateForPlayer(0)
	err := e.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": card.InstanceID, "col": float64(1), "row": float64(0), "devour_ids": []any{food.InstanceID},
	}})
	if err == nil {
		t.Fatal("occupied destination should be rejected")
	}
	if p.Units[0][0] != food || p.Units[1][0] != occupant || len(p.Graveyard) != 0 || p.Elements[model.ElementShadow] != 10 {
		t.Fatalf("rejected summon paid a cost: food=%v occupant=%v graveyard=%v elements=%v", p.Units[0][0], p.Units[1][0], p.Graveyard, p.Elements)
	}
	if !reflect.DeepEqual(before, e.GetStateForPlayer(0)) {
		t.Fatal("rejected summon changed the player view")
	}
}
