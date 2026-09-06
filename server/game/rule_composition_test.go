package game

import "testing"

func TestCounterRuneRecognizesBothPacksAndRejectsOtherCategories(t *testing.T) {
	setupReportedBugEngine(t)
	source := NewCardInstance(baseCard(t, "2021022"), 0, 1)
	for _, tc := range []struct {
		number string
		want   bool
	}{
		{"2021010", true}, {"2221104", true}, {"2021114", true},
		{"2121005", false}, {"2021002", false}, {"3121001", false},
	} {
		t.Run(tc.number, func(t *testing.T) {
			played := NewCardInstance(baseCard(t, tc.number), 1, 1)
			ctx := &CounterContext{Source: source, Event: CounterEvent{Trigger: TriggerOnUseItem, PlayerID: 1, Card: played}}
			if got := (Card2021022CounterRune{}).CanTriggerCounter(ctx); got != tc.want {
				t.Fatalf("%s (%s): got %v want %v", played.Card.Name, played.Card.Tag, got, tc.want)
			}
			ctx.Event.PlayerID = 0
			if (Card2021022CounterRune{}).CanTriggerCounter(ctx) {
				t.Fatal("must not cancel own item")
			}
		})
	}
}

func TestSlotGrantsComposeAndWithdraw(t *testing.T) {
	e := setupReportedBugEngine(t)
	ps := e.State.Players[0]
	for i, number := range []string{"2021002", "2611102", "2021002", "2021017", "2021105"} {
		ps.Equipment[i] = NewCardInstance(baseCard(t, number), 0, 1)
	}
	if got := baseSkillSlotCapacity(ps); got != BaseSkillSlots+1 {
		t.Fatalf("duplicate necklace stacked: %d", got)
	}
	if got := skillSlotCapacity(ps); got != min(MaxSkillSlots, BaseSkillSlots+3) {
		t.Fatalf("restricted grant lost: %d", got)
	}
	ordinary := NewCardInstance(baseCard(t, "3121001"), 0, 1)
	restrictedSlot := baseSkillSlotCapacity(ps)
	if skillAllowedInSlot(ps, ordinary, restrictedSlot) {
		t.Fatal("unrelated spell entered a restricted slot")
	}
	if !playerCanEquipDuplicateSubtypes(ps) {
		t.Fatal("cabinet grant missing")
	}
	ps.Equipment[4].Statuses[StatusPetrify] = 1
	if playerCanEquipDuplicateSubtypes(ps) {
		t.Fatal("petrification did not withdraw grant")
	}
	ps.Equipment[1] = nil
	if skillSlotCapacity(ps) != baseSkillSlotCapacity(ps) {
		t.Fatal("removed candle kept its slots")
	}
}

func TestSearchTriggerQueuesAndResumesParentAfterChoice(t *testing.T) {
	e := setupReportedBugEngine(t)
	placeUnit(baseCard(t, "1211003"), 0, 0, 0, e)
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	e.SetPendingAction(0, "earlier", "earlier", nil, 0, 0, func([]string) {})
	resumed := false
	e.notifyCardSearchedThen(0, NewCardInstance(baseCard(t, "2221002"), 0, 1), func() { resumed = true })
	if resumed || len(e.State.PendingActionQueue) != 1 {
		t.Fatalf("search trigger lost behind earlier choice: resumed=%v queued=%d", resumed, len(e.State.PendingActionQueue))
	}
	resolveEmptyChoice(t, e)
	if resumed || e.State.PendingAction == nil || e.State.PendingAction.Type != "snow_woman_freeze_after_search" {
		t.Fatal("parent resumed before search trigger")
	}
	if err := resolvePendingSelectionWithData(e, 0, []string{target.InstanceID}, nil); err != nil {
		t.Fatal(err)
	}
	if !resumed || target.Statuses[StatusFreeze] != 1 {
		t.Fatal("search effect or continuation missing")
	}
}

func TestDamageRequestUnknownSourceDoesNotBecomePlayerZero(t *testing.T) {
	e := setupReportedBugEngine(t)
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	request := DamageRequest{Target: target, Amount: 1}
	event := damageEventFromContext(&EffectContext{Target: target, ExtraData: request.triggerData()})
	if event.SourcePlayer != -1 || event.IsEnemyDamage(1) {
		t.Fatal("unknown request acquired an attacker")
	}
	request.SourceKnown, request.SourcePlayer = true, 0
	event = damageEventFromContext(&EffectContext{Target: target, ExtraData: request.triggerData()})
	if !event.IsEnemyDamage(1) {
		t.Fatal("explicit player zero source was lost")
	}
}
