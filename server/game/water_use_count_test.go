package game

import (
	"eraofarcane/model"
	"testing"
)

func waterCast(t *testing.T, e *Engine, skill *CardInstance) {
	t.Helper()
	target := "unit"
	if isSorcerySkill(skill.Card) {
		target = "none"
	}
	if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": skill.InstanceID, "target_type": target, "target_col": float64(0), "target_row": float64(0)}}); err != nil {
		t.Fatal(err)
	}
}

func TestWaterCountRippleFirstResolution(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	s := readySkill(baseCard(t, "3221109"), 0)
	p.Skills[0] = s
	p.Elements[model.ElementWater] = 100
	a := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	b := placeUnit(baseCard(t, "1021001"), 1, 1, 0, e)
	a.CurrentLife = 100
	b.CurrentLife = 100
	waterCast(t, e, s)
	if got := len(e.spellAffectedUnits(1, s, e.State.PendingSpell.Target)); got != 1 {
		t.Fatalf("first cast affects %d units, want 1 (power=%d)", got, e.State.PendingSpell.TotalPower)
	}
}

func TestWaterCountTreadingActivation(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	wave := readySkill(baseCard(t, "3221101"), 0)
	s := readySkill(baseCard(t, "3221109"), 0)
	s.IsHorizontal = true
	p.Skills[0], p.Skills[1] = wave, s
	p.Elements[model.ElementWater] = 100
	waterCast(t, e, wave)
	if got := wave.Statuses[treadingWaveTriggerCountStatus]; got != 0 || !s.IsHorizontal {
		t.Fatalf("activating sorcery consumed trigger: count=%d target=%v", got, s.Statuses)
	}
}

func waterResolve(t *testing.T, e *Engine) {
	t.Helper()
	if err := e.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
		t.Fatal(err)
	}
}

func TestWaterCountRippleRepeatedAndReset(t *testing.T) {
	for _, free := range []bool{false, true} {
		t.Run(map[bool]string{false: "paid", true: "free"}[free], func(t *testing.T) {
			e := setupReportedBugEngine(t)
			p := e.State.Players[0]
			a := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
			b := placeUnit(baseCard(t, "1021001"), 1, 1, 0, e)
			a.CurrentLife, b.CurrentLife = 100, 100
			p.Elements[model.ElementWater] = 100
			first := readySkill(baseCard(t, "3221109"), 0)
			second := readySkill(baseCard(t, "3221109"), 0)
			p.Skills[0], p.Skills[1] = first, second
			// Invalid payment must not grant the aura to the next valid attempt.
			p.Elements[model.ElementWater] = 0
			if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": first.InstanceID, "target_type": "unit", "target_col": float64(0), "target_row": float64(0)}}); err == nil {
				t.Fatal("unpaid cast succeeded")
			}
			if spellCastByNumberThisTurn(p, "3221109") != 0 || first.IsHorizontal {
				t.Fatal("failed cast mutated use state")
			}
			p.Elements[model.ElementWater] = 100
			for i, s := range []*CardInstance{first, second, first} {
				s.IsHorizontal = false // fixture: a reset allows reusing the same instance.
				if free {
					if err := e.startFreeSpellCastNoBoost(0, s, SpellTarget{Type: "unit", Position: Position{0, 0}}, nil); err != nil {
						t.Fatal(err)
					}
				} else {
					waterCast(t, e, s)
				}
				wantPower, wantTargets := s.Card.Power, 1
				if i > 0 {
					wantPower += 2
					wantTargets = 2
				}
				if e.State.PendingSpell.TotalPower != wantPower {
					t.Fatalf("cast %d power=%d want=%d", i+1, e.State.PendingSpell.TotalPower, wantPower)
				}
				if got := len(e.spellAffectedUnits(1, s, e.State.PendingSpell.Target)); got != wantTargets {
					t.Fatalf("cast %d targets=%d want=%d", i+1, got, wantTargets)
				}
				wantPositions := 1
				if i > 0 {
					wantPositions = 3
				}
				if got := len(e.pendingSpellAffectedPositions(e.State.PendingSpell)); got != wantPositions {
					t.Fatalf("preview=%d want=%d", got, wantPositions)
				}
				beforeA, beforeB := a.CurrentLife, b.CurrentLife
				waterResolve(t, e)
				if a.CurrentLife >= beforeA || (i == 0 && b.CurrentLife != beforeB) || (i > 0 && b.CurrentLife >= beforeB) {
					t.Fatalf("wrong hit set on cast %d", i+1)
				}
			}
			if err := e.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
				t.Fatal(err)
			}
			if spellCastByNumberThisTurn(p, "3221109") != 0 || e.effectiveSpellArea(first) != SpellAreaSingle {
				t.Fatal("turn count/aura did not reset")
			}
		})
	}
}

func TestWaterCountTreadingLifecycleAndPayments(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	wave := readySkill(baseCard(t, "3221101"), 0)
	a := readySkill(baseCard(t, "3221109"), 0)
	b := readySkill(baseCard(t, "3221109"), 0)
	p.Skills[0], p.Skills[1], p.Skills[2] = wave, a, b
	p.Elements[model.ElementWater] = 100
	unit := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	unit.CurrentLife = 100
	b.IsHorizontal = true
	waterCast(t, e, a)
	waterResolve(t, e)
	if !b.IsHorizontal || wave.Statuses[treadingWaveTriggerCountStatus] != 0 {
		t.Fatal("inactive ability triggered")
	}
	waterCast(t, e, wave)
	if wave.Statuses[treadingWaveTriggerCountStatus] != 0 {
		t.Fatal("activation counted as water spell")
	}
	a.IsHorizontal = false
	for i := 1; i <= 3; i++ {
		used, reset := a, b
		if i%2 == 0 {
			used, reset = b, a
		}
		before := p.Elements[model.ElementWater]
		expectedCost := used.Card.ElementsExpense[model.ElementWater]
		if i > 1 {
			expectedCost += i
		}
		waterCast(t, e, used)
		if got := before - p.Elements[model.ElementWater]; got != expectedCost {
			t.Fatalf("cast %d paid %d want %d", i, got, expectedCost)
		}
		if got := wave.Statuses[treadingWaveTriggerCountStatus]; got != i {
			t.Fatalf("count=%d want=%d", got, i)
		}
		if got := e.effectiveSkillUseCost(p, reset)[model.ElementWater]; got != reset.Card.ElementsExpense[model.ElementWater]+i+1 {
			t.Fatalf("trigger %d next cost=%d", i, got)
		}
		if reset.IsHorizontal {
			t.Fatal("target not reset")
		}
		waterResolve(t, e)
	}
	// Failed next payment preserves both the surcharge and trigger count.
	p.Elements[model.ElementWater] = 0
	if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": b.InstanceID, "target_type": "unit", "target_col": float64(0), "target_row": float64(0)}}); err == nil {
		t.Fatal("unpaid reuse succeeded")
	}
	if wave.Statuses[treadingWaveTriggerCountStatus] != 3 || e.effectiveSkillUseCost(p, b)[model.ElementWater] != b.Card.ElementsExpense[model.ElementWater]+4 {
		t.Fatal("failed use consumed count or surcharge")
	}
	if err := e.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
		t.Fatal(err)
	}
	if abilityDurationActive(wave) {
		t.Fatal("ability survived owner turn end")
	}
	if err := e.HandleAction(1, ActionMessage{Action: "end_turn"}); err != nil {
		t.Fatal(err)
	}
	// Fixture removes printed cooldown to test reactivation on a later turn.
	wave.IsHorizontal = false
	delete(wave.Statuses, StatusCooldown)
	p.Elements[model.ElementWater] = 100
	waterCast(t, e, wave)
	b.IsHorizontal = true
	a.IsHorizontal = false
	waterCast(t, e, a)
	if wave.Statuses[treadingWaveTriggerCountStatus] != 1 {
		t.Fatal("new turn did not restart trigger count")
	}
	// Unused prior surcharge survives turn end and adds to the new +2.
	if got := e.effectiveSkillUseCost(p, b)[model.ElementWater]; got != b.Card.ElementsExpense[model.ElementWater]+6 {
		t.Fatalf("next-use surcharge lifetime wrong: %d", got)
	}
}

func TestWaterCountTreadingChoiceFailureAndCancelledSpell(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	wave := readySkill(baseCard(t, "3221101"), 0)
	a := readySkill(baseCard(t, "3221109"), 0)
	b := readySkill(baseCard(t, "3221109"), 0)
	c := readySkill(baseCard(t, "3221109"), 0)
	p.Skills[0], p.Skills[1], p.Skills[2], p.Skills[3] = wave, a, b, c
	p.Elements[model.ElementWater] = 100
	unit := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	unit.CurrentLife = 100
	waterCast(t, e, wave)
	b.IsHorizontal, c.IsHorizontal = true, true
	waterCast(t, e, a)
	pending := e.State.PendingAction
	if pending == nil || pending.Type != "treading_wave_reset" {
		t.Fatal("missing choice")
	}
	for _, ids := range [][]any{{}, {"invalid"}, {b.InstanceID, c.InstanceID}} {
		if err := e.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": ids}}); err == nil {
			t.Fatal("invalid/empty mandatory choice accepted")
		}
		if e.State.PendingAction != pending || wave.Statuses[treadingWaveTriggerCountStatus] != 0 || !b.IsHorizontal {
			t.Fatal("rejected choice changed count/reset")
		}
	}
	if err := e.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{b.InstanceID}}}); err != nil {
		t.Fatal(err)
	}
	if wave.Statuses[treadingWaveTriggerCountStatus] != 1 || b.IsHorizontal || !c.IsHorizontal {
		t.Fatal("valid selection not counted exactly once")
	}
	// Cancellation is after commitment: it must keep the cast and reset counts.
	e.cancelPendingSpell(1, nil, "test cancellation")
	if e.State.PendingSpell != nil || spellCastByNumberThisTurn(p, "3221109") != 1 || wave.Statuses[treadingWaveTriggerCountStatus] != 1 {
		t.Fatal("cancelled committed spell lost count")
	}
	// Remove extra candidate to test the next real cast without another choice.
	c.IsHorizontal = false
	waterCast(t, e, b)
	if e.State.PendingSpell.TotalPower != b.Card.Power+2 || wave.Statuses[treadingWaveTriggerCountStatus] != 2 {
		t.Fatal("next cast ignored committed cancelled use")
	}
	waterResolve(t, e)
}

func TestWaterCountTreadingNoEligibleResetDoesNotCount(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	wave := readySkill(baseCard(t, "3221101"), 0)
	otherSorcery := readySkill(baseCard(t, "3221101"), 0)
	a := readySkill(baseCard(t, "3221109"), 0)
	b := readySkill(baseCard(t, "3221109"), 0)
	p.Skills[0], p.Skills[1], p.Skills[2], p.Skills[3] = wave, a, b, otherSorcery
	p.Elements[model.ElementWater] = 100
	unit := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	unit.CurrentLife = 100
	waterCast(t, e, wave)
	otherSorcery.IsHorizontal = true
	waterCast(t, e, a)
	if wave.Statuses[treadingWaveTriggerCountStatus] != 0 || !otherSorcery.IsHorizontal || e.State.PendingAction != nil {
		t.Fatal("ineligible sorcery reset or empty trigger counted")
	}
	waterResolve(t, e)
	waterCast(t, e, b)
	if wave.Statuses[treadingWaveTriggerCountStatus] != 1 || e.effectiveSkillUseCost(p, a)[model.ElementWater] != a.Card.ElementsExpense[model.ElementWater]+2 {
		t.Fatal("first actual reset must add exactly two")
	}
	waterResolve(t, e)
}
