package game

import "testing"

type countingDamageReduction struct {
	AlwaysActive
	calls *int
}

func (countingDamageReduction) ID() string   { return "1021001" }
func (countingDamageReduction) Name() string { return "damage reduction fixture" }
func (b countingDamageReduction) ModifyDamageAmount(_ *EffectContext, amount int) int {
	*b.calls++
	return amount - 1
}

func TestDamageDeclinedReplacementDoesNotRecalculate(t *testing.T) {
	e := setupReportedBugEngine(t)
	calls := 0
	globalRegistry.RegisterBehaviorFactory("1021001", func() CardBehavior { return countingDamageReduction{calls: &calls} })
	placeUnit(baseCard(t, "1221001"), 0, 0, 0, e)
	target := placeUnit(baseCard(t, "1021001"), 0, 1, 0, e)
	target.CurrentLife = 2
	e.dealDamage(target, 5, 0)
	if calls != 1 || target.CurrentLife != 2 || e.State.PendingAction == nil {
		t.Fatalf("damage must be calculated once, then wait: calls=%d life=%d pending=%v", calls, target.CurrentLife, e.State.PendingAction)
	}
	resolveEmptyChoice(t, e)
	if calls != 1 || target.CurrentLife != -2 || target.DamageTakenThisTurn != 4 {
		t.Fatalf("resume recalculated damage: calls=%d life=%d damage=%d", calls, target.CurrentLife, target.DamageTakenThisTurn)
	}
}

func TestDamageAdjustmentUsesOneReceiptAndPreventionWindow(t *testing.T) {
	t.Run("burn aura produces one damage event", func(t *testing.T) {
		e := setupReportedBugEngine(t)
		placeUnit(baseCard(t, "1111003"), 0, 0, 0, e)
		target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
		target.CurrentLife = 4
		e.dealDamageWithExtra(target, 1, 1, map[string]any{"status_damage": StatusBurn})
		count := 0
		for _, event := range e.log {
			if event.Type == "damage" {
				count++
				if event.Data["amount"] != 2 {
					t.Fatalf("wrong receipt: %v", event.Data)
				}
			}
		}
		if count != 1 || target.DamageTakenThisTurn != 2 {
			t.Fatalf("expected one amplified event, count=%d damage=%d", count, target.DamageTakenThisTurn)
		}
	})
	t.Run("amplified damage is preventable before becoming lethal", func(t *testing.T) {
		e := setupReportedBugEngine(t)
		dolphin := placeUnit(baseCard(t, "1221001"), 0, 0, 0, e)
		titan := placeUnit(baseCard(t, "1421007"), 0, 1, 0, e)
		titan.CurrentLife = 2
		e.dealDamageWithExtra(titan, 1, 0, map[string]any{"damage_source": "spell", "boost_count": 0})
		if e.State.PendingAction == nil || titan.CurrentLife != 2 {
			t.Fatal("amplification bypassed lethal prevention")
		}
		if err := resolvePendingSelectionWithData(e, 0, []string{dolphin.InstanceID}, nil); err != nil {
			t.Fatal(err)
		}
		if titan.CurrentLife != 2 || titan.DamageTakenThisTurn != 0 {
			t.Fatal("prevented damage was committed")
		}
	})
}

func TestSpellHitWaitsForDamageReplacementReceipt(t *testing.T) {
	for _, prevent := range []bool{false, true} {
		t.Run(map[bool]string{false: "decline", true: "prevent"}[prevent], func(t *testing.T) {
			e := setupReportedBugEngine(t)
			dolphin := placeUnit(baseCard(t, "1221001"), 1, 0, 0, e)
			target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, e)
			target.CurrentLife = 1
			skill := readySkill(baseCard(t, "3121001"), 0)
			e.State.PendingSpell = &SpellCast{AttackerID: 0, Skill: skill, Target: SpellTarget{Type: "unit", Position: *target.Position}, TotalPower: 5}
			e.State.Phase = PhaseDefenseWindow
			e.resolvePendingSpellHit()
			if e.State.PendingAction == nil || e.State.PendingSpell == nil {
				t.Fatal("spell must wait for damage replacement")
			}
			if e.State.Players[0].SpellHitsThisTurn != 0 || e.State.Players[0].SpellDamageThisTurn != 0 {
				t.Fatal("hit effects ran before the damage decision")
			}
			var selected []string
			if prevent {
				selected = []string{dolphin.InstanceID}
			}
			if err := resolvePendingSelectionWithData(e, 1, selected, nil); err != nil {
				t.Fatal(err)
			}
			if e.State.PendingSpell != nil || e.State.Players[0].SpellHitsThisTurn != 1 {
				t.Fatal("spell did not complete exactly once")
			}
			if got := e.State.Players[0].SpellDamageThisTurn; got != target.DamageTakenThisTurn {
				t.Fatalf("hit receipt=%d actual=%d", got, target.DamageTakenThisTurn)
			}
			if prevent && target.DamageTakenThisTurn != 0 {
				t.Fatal("prevented damage entered the receipt")
			}
			if !prevent && target.DamageTakenThisTurn == 0 {
				t.Fatal("declined damage missing from the receipt")
			}
		})
	}
}

func TestHiddenDamageTriggerQueuesBehindUnrelatedChoice(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	p.Hand = []*CardInstance{NewCardInstance(baseCard(t, "1401002"), 0, 1)}
	p.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1401002"), 0, 1)}
	target := placeUnit(baseCard(t, "1021001"), 0, 0, 0, e)
	target.CurrentLife = 4
	source := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	e.dealDamage(target, 1, 0)
	if e.State.PendingAction != nil {
		t.Fatal("unknown damage source must not be treated as an enemy")
	}
	e.SetPendingAction(0, "earlier_choice", "earlier choice", nil, 0, 0, func([]string) {})
	ctx := &EffectContext{Engine: e, Source: source, PlayerID: 1, OpponentID: 0}
	ctx.DealDamage(target, 1)
	if e.State.PendingAction.Type != "earlier_choice" || len(e.State.PendingActionQueue) != 1 {
		t.Fatalf("hidden trigger was dropped or duplicated: pending=%v queue=%d", e.State.PendingAction, len(e.State.PendingActionQueue))
	}
	resolveEmptyChoice(t, e)
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "xinke_summon" || len(e.State.PendingAction.Candidates) != 2 {
		t.Fatal("queued hidden-zone choice was not activated")
	}
	resolveEmptyChoice(t, e)
	if e.State.PendingAction != nil {
		t.Fatal("identical hidden copies created duplicate choices for one event")
	}
}

func TestDamageModifierQueriesDoNotSpendReduction(t *testing.T) {
	e := setupReportedBugEngine(t)
	leggings := NewCardInstance(baseCard(t, "2421111"), 0, 1)
	e.State.Players[0].Equipment[0] = leggings
	target := placeUnit(baseCard(t, "1021001"), 0, 0, 0, e)
	target.CurrentLife = 5
	for i := 0; i < 3; i++ {
		plan := e.planDamageModifiers(target, 3, 0, nil)
		if plan.Amount != 1 || leggings.UltimateUsed || target.CurrentLife != 5 {
			t.Fatal("query spent the damage reduction")
		}
	}
	e.dealDamage(target, 3, 0)
	if target.CurrentLife != 4 || !leggings.UltimateUsed {
		t.Fatal("actual damage did not commit the queried reduction")
	}
	if next := e.planDamageModifiers(target, 3, 0, nil); next.Amount != 3 {
		t.Fatal("spent reduction applied again")
	}
}
