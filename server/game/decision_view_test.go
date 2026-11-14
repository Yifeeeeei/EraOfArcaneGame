package game

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDecisionRevisionRejectsOldDraftAndDuplicate(t *testing.T) {
	e := setupTestEngine(t)
	e.State.Phase = PhaseMain
	e.State.CurrentTurn = 0
	revision := e.revision
	action := ActionMessage{Action: "offer_draw", ExpectedRevision: &revision}
	if err := e.HandleAction(0, action); err != nil {
		t.Fatal(err)
	}
	if e.revision != revision+1 {
		t.Fatal("accepted action must advance revision")
	}
	if err := e.HandleAction(0, action); err == nil || !strings.Contains(err.Error(), "stale revision") {
		t.Fatalf("duplicate accepted: %v", err)
	}
	if e.revision != revision+1 {
		t.Fatal("rejected action advanced revision")
	}
}

func TestConsecutiveIdenticalChoicesHaveDistinctIDsAndPrivacy(t *testing.T) {
	e := setupTestEngine(t)
	e.State.Phase = PhaseMain
	candidates := []map[string]any{{"instance_id": "private-candidate", "name": "private-counter"}}
	e.SetPendingAction(0, "counter_trigger", "private-prompt", candidates, 0, 1, func([]string) {})
	first := e.State.PendingAction.DecisionID
	e.SetPendingAction(0, "counter_trigger", "private-prompt", candidates, 0, 1, func([]string) {})
	if err := e.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{}}}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction == nil || e.State.PendingAction.DecisionID <= first {
		t.Fatal("consecutive decision identity reused")
	}
	for _, view := range []map[string]any{e.GetStateForPlayer(1), e.GetStateForSpectator()} {
		b, _ := json.Marshal(view)
		for _, secret := range []string{"private-candidate", "private-counter", "private-prompt", "counter_trigger"} {
			if strings.Contains(string(b), secret) {
				t.Fatalf("private decision leaked: %s", secret)
			}
		}
		if view["decision_player"] != 0 {
			t.Fatal("decision owner missing")
		}
	}
	own := e.GetStateForPlayer(0)["pending_action"].(map[string]any)
	if own["decision_id"] != e.State.PendingAction.DecisionID {
		t.Fatal("owner identity missing")
	}
}

func TestDecisionPrimitiveContracts(t *testing.T) {
	for typ, kind := range map[string]string{"water_divination_search": "reorder", "natural_communion_distribute": "load_distribution", "new_mechanic": "selection"} {
		if pendingActionUI(&PendingAction{Type: typ})["kind"] != kind {
			t.Fatal(typ)
		}
	}
}

func TestAllocationPrimitiveValidatesBeforeCallback(t *testing.T) {
	e := setupTestEngine(t)
	e.State.Phase = PhaseMain
	var received []string
	e.SetAllocatedPendingAction(0, "allocation", "allocate", []map[string]any{{"instance_id": "a"}, {"instance_id": "b"}}, 1, 2, 4, func(ids []string) { received = ids })
	invalid := ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{"a", "b"}, "allocations": map[string]any{"a": float64(2), "b": float64(1)}}}
	if err := e.HandleAction(0, invalid); err == nil {
		t.Fatal("invalid total accepted")
	}
	if len(received) != 0 || e.State.PendingAction == nil {
		t.Fatal("invalid allocation consumed decision")
	}
	invalid.Data["allocations"] = map[string]any{"a": float64(3.5), "b": float64(1)}
	if err := e.HandleAction(0, invalid); err == nil {
		t.Fatal("fractional allocation accepted")
	}
	invalid.Data["allocations"] = map[string]any{"a": float64(2), "b": float64(2)}
	if err := e.HandleAction(0, invalid); err != nil {
		t.Fatal(err)
	}
	if strings.Join(received, ",") != "a,a,b,b" {
		t.Fatalf("wrong allocation: %v", received)
	}
}

func TestDebugMutationsAreAtomicAndVersioned(t *testing.T) {
	e := setupTestEngine(t)
	rev := e.revision
	if rev == 0 {
		t.Fatal("setup must version its initial state")
	}
	p := e.State.Players[0]
	before := p.Elements["火"]
	if err := e.DebugUpdateElements(0, map[string]int{"火": 4, "invalid": 1}, "set"); err == nil {
		t.Fatal("invalid element accepted")
	}
	if p.Elements["火"] != before || e.revision != rev {
		t.Fatal("failed element update mutated state")
	}
	life := p.Hero.CurrentLife + 1
	if err := e.DebugMutateCard(0, DebugCardMutation{InstanceID: p.Hero.InstanceID, Life: &life, Statuses: map[string]int{"invalid": -1}}); err == nil {
		t.Fatal("invalid status accepted")
	}
	if p.Hero.CurrentLife == life || e.revision != rev {
		t.Fatal("failed card update mutated state")
	}
	for i := range p.Equipment {
		p.Equipment[i] = nil
	}
	if err := e.DebugAddCard(0, "2011001", "equipment", len(p.Equipment)+1); err == nil {
		t.Fatal("overflow accepted")
	}
	for _, card := range p.Equipment {
		if card != nil {
			t.Fatal("failed batch partially added equipment")
		}
	}
	if e.revision != rev {
		t.Fatal("failed batch changed revision")
	}
	if err := e.DebugUpdateElements(0, map[string]int{"火": 4}, "set"); err != nil {
		t.Fatal(err)
	}
	if e.revision != rev+1 {
		t.Fatal("successful fixture not versioned")
	}
	if err := e.HandleAction(0, ActionMessage{Action: "mulligan", ExpectedRevision: &rev, Data: map[string]any{"keep": true}}); err == nil {
		t.Fatal("draft from before fixture accepted")
	}
}

func TestExecutedErrorVersionsRealPartialPaymentWithoutRollback(t *testing.T) {
	e := setupReportedBugEngine(t)
	e.State.Phase = PhaseMain
	e.State.CurrentTurn = 0
	source := readySkill(baseCard(t, "3411102"), 0)
	e.State.Players[0].Skills[0] = source
	unit := placeUnit(baseCard(t, "1021002"), 0, 0, 0, e)
	unit.CurrentLife = 2
	if err := e.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{"instance_id": source.InstanceID, "ability": "per_turn"}}); err != nil {
		t.Fatal(err)
	}
	decision := e.State.PendingAction.DecisionID
	before := e.revision
	action := ActionMessage{Action: "resolve_action", ExpectedRevision: &before, Data: map[string]any{"selected": []any{unit.InstanceID}, "payments": []any{map[string]any{"instance_id": unit.InstanceID, "amount": float64(2), "mode": "life"}}}}
	revision, executed, err := e.HandleActionWithRevision(0, action)
	if err == nil || err.Error() != "cannot pay life from selected unit" {
		t.Fatalf("original error changed: %v", err)
	}
	if !executed || revision != before+1 || unit.CurrentLife != 1 {
		t.Fatalf("partial mutation not versioned: revision=%d executed=%v life=%d", revision, executed, unit.CurrentLife)
	}
	if e.State.PendingAction == nil || e.State.PendingAction.DecisionID != decision {
		t.Fatal("failed decision identity changed")
	}
	staleRev, staleExecuted, staleErr := e.HandleActionWithRevision(0, action)
	if staleErr == nil || staleExecuted || staleRev != revision || unit.CurrentLife != 1 {
		t.Fatal("old action repeated partial mutation")
	}
	// Invalid selection enters the rule handler but does not modify the unit.
	action.ExpectedRevision = &revision
	action.Data = map[string]any{"selected": []any{"not-a-candidate"}}
	next, executed, err := e.HandleActionWithRevision(0, action)
	if err == nil || err.Error() != "invalid selection" || !executed || next != revision+1 || unit.CurrentLife != 1 {
		t.Fatal("unmodified rule rejection must still conservatively advance")
	}
	for _, view := range []map[string]any{e.GetStateForPlayer(1), e.GetStateForSpectator()} {
		raw, _ := json.Marshal(view)
		if strings.Contains(string(raw), "blood_sand_array_pay") || strings.Contains(string(raw), "payments") {
			t.Fatal("private choice leaked after error")
		}
		if view["revision"] != next {
			t.Fatal("filtered view missed rejected attempt revision")
		}
	}
	for _, bad := range []ActionMessage{{Action: "unknown"}, {Action: "resolve_action", ExpectedRevision: &before}} {
		got, ran, err := e.HandleActionWithRevision(0, bad)
		if err == nil || ran || got != next {
			t.Fatal("preflight failure entered rules or advanced revision")
		}
	}
}
