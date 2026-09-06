package game

import (
	"fmt"
	"reflect"
	"testing"
)

func resolveEmptyChoice(t *testing.T, e *Engine) {
	t.Helper()
	if e.State.PendingAction == nil {
		t.Fatal("expected a pending choice")
	}
	if err := e.HandleAction(e.State.PendingAction.PlayerID, ActionMessage{
		Action: "resolve_action", Data: map[string]any{"selected": []any{}},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestResolutionWaitsForNestedChoiceSpellAndQueuedChoices(t *testing.T) {
	e := setupReportedBugEngine(t)
	var trace []string
	var spell *SpellCast
	frame := e.runResolution("nested sequence", func() {
		e.SetPendingAction(0, "outer", "outer", nil, 0, 0, func([]string) {
			trace = append(trace, "outer")
			e.SetPendingActionWithData(0, "inner", "inner", nil, 0, 0, func([]string, map[string]any) {
				trace = append(trace, "inner")
				spell = &SpellCast{Skill: readySkill(baseCard(t, "3121002"), 0)}
				e.State.PendingSpell = spell
				e.State.Phase = PhaseDefenseWindow
			})
		})
		e.SetPendingAction(1, "sibling", "sibling", nil, 0, 0, func([]string) {
			trace = append(trace, "sibling")
		})
	}, func() { trace = append(trace, "parent") })
	resolveEmptyChoice(t, e)
	resolveEmptyChoice(t, e)
	resolveEmptyChoice(t, e)
	if !reflect.DeepEqual(trace, []string{"outer", "inner", "sibling"}) || frame.waitingSpell != spell {
		t.Fatalf("parent must wait for nested spell: trace=%v frame=%+v", trace, frame)
	}
	e.completePendingSpell(spell)
	e.completePendingSpell(spell)
	e.resumeResolution(frame)
	if !reflect.DeepEqual(trace, []string{"outer", "inner", "sibling", "parent"}) || !frame.done {
		t.Fatalf("parent must resume exactly once: %v, %+v", trace, frame)
	}
}

func TestResolutionFailedChoiceRetainsContinuation(t *testing.T) {
	e := setupReportedBugEngine(t)
	attempts, resumed := 0, 0
	e.SetPendingActionWithError(0, "retry", "retry", nil, 0, 0, nil, false,
		func([]string, map[string]any) error {
			attempts++
			if attempts == 1 {
				return fmt.Errorf("payment not valid")
			}
			return nil
		})
	choice := e.State.PendingAction
	e.continueAfterPendingAction(func() { resumed++ })
	if err := e.HandleAction(0, ActionMessage{Action: "resolve_action"}); err == nil {
		t.Fatal("expected retryable error")
	}
	if resumed != 0 || e.State.PendingAction != choice || len(choice.resolutions) != 1 {
		t.Fatal("failed choice consumed its continuation")
	}
	resolveEmptyChoice(t, e)
	if resumed != 1 || attempts != 2 {
		t.Fatalf("attempts=%d resumed=%d", attempts, resumed)
	}
}

func TestResolutionResumesBaselineSpellBeforeParent(t *testing.T) {
	e := setupReportedBugEngine(t)
	var trace []string
	spell := &SpellCast{Skill: readySkill(baseCard(t, "3121002"), 0)}
	e.SetPendingAction(0, "cast", "cast", nil, 0, 0, func([]string) {
		e.State.PendingSpell = spell
		e.SetPendingAction(0, "on cast", "on cast", nil, 0, 0, nil)
		e.continueAfterPendingAction(func() { trace = append(trace, "defense"); e.State.Phase = PhaseDefenseWindow })
	})
	e.continueAfterPendingAction(func() { trace = append(trace, "end turn") })
	resolveEmptyChoice(t, e)
	resolveEmptyChoice(t, e)
	if !reflect.DeepEqual(trace, []string{"defense"}) {
		t.Fatalf("baseline spell deadlocked or parent resumed early: %v", trace)
	}
	e.completePendingSpell(spell)
	if !reflect.DeepEqual(trace, []string{"defense", "end turn"}) {
		t.Fatalf("wrong continuation order: %v", trace)
	}
}

func TestResolutionSkippedChoiceContinuesWithoutCallingEffect(t *testing.T) {
	e := setupReportedBugEngine(t)
	available := true
	called, resumed := false, 0
	e.SetPendingAction(0, "first", "first", nil, 0, 0, nil)
	queued := e.setPendingActionWithOptions(1, "stale", "stale", nil, 0, 0, nil, false,
		func([]string) { called = true }, nil, nil, nil, func() bool { return available })
	e.addActionContinuation(queued, "after stale", func() { resumed++ })
	e.SetPendingAction(1, "last", "last", nil, 0, 0, nil)
	available = false
	resolveEmptyChoice(t, e)
	if resumed != 0 || e.State.PendingAction == nil || e.State.PendingAction.Type != "last" {
		t.Fatal("skipped choice resumed before a remaining sibling")
	}
	resolveEmptyChoice(t, e)
	if called || resumed != 1 || e.State.PendingAction != nil {
		t.Fatalf("called=%v resumed=%d pending=%v", called, resumed, e.State.PendingAction)
	}
}

func TestResolutionSpellReplacementTransfersWaitersAndCanRetry(t *testing.T) {
	e := setupReportedBugEngine(t)
	original := &SpellCast{Skill: readySkill(baseCard(t, "3121002"), 0)}
	replacement := &SpellCast{Skill: readySkill(baseCard(t, "3121002"), 1)}
	resumed := 0
	frame := e.runResolution("reflection", func() { e.State.PendingSpell = original }, func() { resumed++ })
	err := e.replacePendingSpell(func() error { return fmt.Errorf("invalid target") })
	if err == nil || e.State.PendingSpell != original || resumed != 0 || frame.waitingSpell != original {
		t.Fatal("failed replacement lost original spell or its waiters")
	}
	if err := e.replacePendingSpell(func() error { e.State.PendingSpell = replacement; return nil }); err != nil {
		t.Fatal(err)
	}
	if resumed != 0 || frame.waitingSpell != replacement || len(original.resolutions) != 0 {
		t.Fatal("replacement did not inherit waiting parent")
	}
	e.completePendingSpell(original)
	e.completePendingSpell(replacement)
	if resumed != 1 {
		t.Fatalf("resumed=%d", resumed)
	}
}

func TestResolutionStopsAtGameOver(t *testing.T) {
	e := setupReportedBugEngine(t)
	resumed := false
	e.runResolution("stop", func() {
		e.SetPendingAction(0, "finish", "finish", nil, 0, 0, func([]string) { e.finishGame(0, "test", 0) })
	}, func() { resumed = true; e.State.Phase = PhaseMain })
	resolveEmptyChoice(t, e)
	if resumed || e.State.Phase != PhaseGameOver {
		t.Fatal("continuation reopened a finished game")
	}
}

func TestResolutionFinalSpellRestoresMainPhase(t *testing.T) {
	e := setupReportedBugEngine(t)
	skill := readySkill(baseCard(t, "3121002"), 0)
	e.State.Players[0].Skills[0] = skill
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	frame := e.runResolution("last step is a spell", func() {
		if err := e.startVirtualSpellCastNoBoost(0, skill, SpellTarget{Type: "unit", Position: *target.Position}, nil); err != nil {
			t.Fatal(err)
		}
	})
	if frame.done || frame.waitingSpell == nil {
		t.Fatal("final spell must remain observable as waiting")
	}
	if err := e.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
		t.Fatal(err)
	}
	if !frame.done || e.State.PendingSpell != nil || e.State.Phase != PhaseMain {
		t.Fatalf("final spell left a phantom defense window: frame=%+v phase=%s", frame, e.State.Phase)
	}
}

func TestResolutionEmbersReflectionFinishesBeforeTurnAdvances(t *testing.T) {
	e := setupReportedBugEngine(t)
	embers := readySkill(baseCard(t, "3121105"), 0)
	e.State.Players[0].Skills[0] = embers
	originalTarget := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	reflectedTarget := placeUnit(baseCard(t, "1021001"), 0, 0, 0, e)
	counter := NewCardInstance(baseCard(t, "2321111"), 1, e.State.TurnNumber)
	counter.IsSetCounter = true
	e.State.Players[1].Equipment[0] = counter
	e.State.Players[1].Elements["气"] = 2
	act := func(player int, action string, selected ...string) {
		t.Helper()
		ids := make([]any, len(selected))
		for i, id := range selected {
			ids[i] = id
		}
		if err := e.HandleAction(player, ActionMessage{Action: action, Data: map[string]any{"selected": ids}}); err != nil {
			t.Fatal(err)
		}
	}
	act(0, "end_turn")
	act(0, "resolve_action", originalTarget.InstanceID)
	e.cancelPendingSpell(1, counter, "test reflection")
	act(1, "resolve_action", counter.InstanceID)
	act(1, "resolve_action", reflectedTarget.InstanceID)
	if e.State.CurrentTurn != 0 || e.State.PendingSpell == nil || e.State.PendingSpell.AttackerID != 1 {
		t.Fatalf("turn advanced before reflection: turn=%d spell=%+v", e.State.CurrentTurn, e.State.PendingSpell)
	}
	act(0, "no_defend")
	if e.State.CurrentTurn != 1 || e.State.PendingAction != nil || e.State.PendingSpell != nil {
		t.Fatalf("reflection did not finish turn: turn=%d action=%+v spell=%+v", e.State.CurrentTurn, e.State.PendingAction, e.State.PendingSpell)
	}
}
