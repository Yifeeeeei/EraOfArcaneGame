package game

import (
	"eraofarcane/cards"
	"reflect"
	"testing"
)

func newLesson(t *testing.T) *Engine {
	t.Helper()
	if err := cards.LoadCards(); err != nil {
		t.Fatal(err)
	}
	SetCardDB(cards.PlayableCardDB)
	RegisterAllCardEffects()
	e := NewEngine("lesson-regression", nil)
	if err := e.SetupTutorialGame("新玩家"); err != nil {
		t.Fatal(err)
	}
	return e
}

func TestFirstLessonCompletesRealSpellBattleAndHandLimit(t *testing.T) {
	e := newLesson(t)
	seenDefense, seenDamage, seenDiscard := false, false, false
	for i := 0; i < 150 && e.tutorial.index < len(e.tutorial.steps); i++ {
		step := e.tutorial.steps[e.tutorial.index]
		action := e.tutorialStepAction(step)
		before := e.tutorial.index
		if err := e.HandleAction(step.player, action); err != nil {
			t.Fatalf("step %d %s (%s): %v", before, step.title, action.Action, err)
		}
		if e.tutorial.index != before+1 {
			t.Fatalf("step failed to advance exactly once: %d -> %d", before, e.tutorial.index)
		}
		if action.Action == "defend" {
			seenDefense = true
			if e.State.Players[0].Hero.CurrentLife != 6 && e.State.Players[0].Hero.CurrentLife != 5 {
				t.Fatal("defense did not preserve life")
			}
		}
		if action.Action == "resolve_action" && step.player == 0 {
			seenDiscard = true
			if len(e.State.Players[0].Hand) != 5 {
				t.Fatal("real hand limit discard did not finish")
			}
			if totalElementCost(e.State.Players[0].Elements) != 0 {
				t.Fatal("end-turn elements did not clear")
			}
			if e.State.CurrentTurn != 1 {
				t.Fatal("turn didn't pass after discard")
			}
		}
		if e.State.Players[1].Hero.CurrentLife == 4 {
			seenDamage = true
		}
	}
	if !seenDefense || !seenDamage || !seenDiscard {
		t.Fatal("lesson skipped a required combat/turn phase")
	}
	if e.State.Phase != PhaseGameOver || e.State.Winner != 0 || e.State.TurnNumber != 5 {
		t.Fatalf("lesson didn't reach real final victory: phase=%v winner=%d turn=%d", e.State.Phase, e.State.Winner, e.State.TurnNumber)
	}
	if e.State.Players[0].Hero.CurrentLife != 5 || e.State.Players[1].Hero.CurrentLife != 0 {
		t.Fatal("unexpected final life totals")
	}
	if e.tutorial.failure != "" || e.tutorial.index != len(e.tutorial.steps) {
		t.Fatalf("incomplete lesson: %+v", e.tutorial)
	}
}

func TestTutorialCannotSkipOrReplaySteps(t *testing.T) {
	e := newLesson(t)
	before := e.GetStateForPlayer(0)
	rev := e.revision
	for _, action := range []ActionMessage{
		{Action: "end_turn", ExpectedRevision: &rev},
		{Action: "mulligan", Data: map[string]any{"keep": true}, ExpectedRevision: &rev},
		{Action: "tutorial_continue"},
	} {
		if _, executed, err := e.HandleActionWithRevision(0, action); err == nil || executed {
			t.Fatal("tutorial bypass was executed")
		}
	}
	if !reflect.DeepEqual(before, e.GetStateForPlayer(0)) {
		t.Fatal("rejected tutorial action changed state")
	}
	action := e.tutorialStepAction(e.tutorial.steps[0])
	if err := e.HandleAction(0, action); err != nil {
		t.Fatal(err)
	}
	if _, executed, err := e.HandleActionWithRevision(0, action); err == nil || executed {
		t.Fatal("stale reading acknowledgement advanced twice")
	}
	if e.NextTutorialOpponentAction() != nil {
		t.Fatal("opponent skipped a human reading pause")
	}
}

func TestTutorialPinsOtherwiseLegalPaymentAndPrivateViews(t *testing.T) {
	e := newLesson(t)
	for e.tutorial.index < len(e.tutorial.steps) {
		s := e.tutorial.steps[e.tutorial.index]
		if s.action == "learn_skill" && s.player == 0 {
			break
		}
		if err := e.HandleAction(s.player, e.tutorialStepAction(s)); err != nil {
			t.Fatal(err)
		}
	}
	s := e.tutorial.steps[e.tutorial.index]
	a := e.tutorialStepAction(s)
	a.Data["payment"] = map[string]int{"火": 4, "无": 1} // Legal base payment, wrong lesson budget.
	before := e.GetStateForPlayer(0)
	if _, executed, err := e.HandleActionWithRevision(0, a); err == nil || executed {
		t.Fatal("wrong payment accepted")
	}
	if !reflect.DeepEqual(before, e.GetStateForPlayer(0)) {
		t.Fatal("wrong payment spent resources")
	}
	if e.GetStateForPlayer(1)["tutorial"] != nil {
		t.Fatal("opponent received lesson script")
	}
	if _, ok := e.GetStateForSpectator()["tutorial"]; ok {
		t.Fatal("spectator received lesson script")
	}
}

func TestTutorialWaitingViewDoesNotExposeOpponentPrivateSkillIDs(t *testing.T) {
	e := newLesson(t)
	for e.tutorial.index < len(e.tutorial.steps) {
		step := e.tutorial.steps[e.tutorial.index]
		if step.player == 1 && step.action == "learn_skill" {
			view := e.GetStateForPlayer(0)["tutorial"].(map[string]any)
			if view["card_id"] != nil || view["allowed"] != nil {
				t.Fatal("waiting view exposed opponent action identity")
			}
			return
		}
		if err := e.HandleAction(step.player, e.tutorialStepAction(step)); err != nil {
			t.Fatal(err)
		}
	}
	t.Fatal("did not reach opponent learning window")
}
