package game

import (
	"testing"

	"eraofarcane/model"
)

func TestDrawPendingActionsQueueByEquipmentOrder(t *testing.T) {
	t.Run("thunder drum before windbreath compass", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		drum := NewCardInstance(baseCard(t, "2311002"), 0, 1)
		compass := NewCardInstance(baseCard(t, "2321001"), 0, 1)
		p0.Equipment[0] = drum
		p0.Equipment[1] = compass
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}

		engine.drawCards(0, 1)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "thunder_drum_mark" {
			t.Fatalf("thunder drum should ask before compass, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{drum.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve thunder drum trigger: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "windbreath_compass" {
			t.Fatalf("windbreath compass should remain queued after thunder drum, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{compass.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve windbreath compass trigger: %v", err)
		}
		if thunderDrumMarks(drum) != 1 || effectiveElementsGain(compass)[model.ElementAir] != compass.Card.ElementsGain[model.ElementAir]+1 {
			t.Fatalf("both draw prompts should resolve, drum=%v compass_load=%v", drum.Statuses, effectiveElementsGain(compass))
		}
		if engine.State.PendingAction != nil || len(engine.State.PendingActionQueue) != 0 {
			t.Fatalf("draw prompt queue should be empty, pending=%+v queue=%d", engine.State.PendingAction, len(engine.State.PendingActionQueue))
		}
	})

	t.Run("windbreath compass before thunder drum", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		compass := NewCardInstance(baseCard(t, "2321001"), 0, 1)
		drum := NewCardInstance(baseCard(t, "2311002"), 0, 1)
		p0.Equipment[0] = compass
		p0.Equipment[1] = drum
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}

		engine.drawCards(0, 1)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "windbreath_compass" {
			t.Fatalf("windbreath compass should ask before thunder drum, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{},
		}}); err != nil {
			t.Fatalf("decline windbreath compass trigger: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "thunder_drum_mark" {
			t.Fatalf("thunder drum should remain queued after compass, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{drum.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve thunder drum trigger: %v", err)
		}
		if thunderDrumMarks(drum) != 1 {
			t.Fatalf("thunder drum should gain a mark from queued trigger, statuses=%v", drum.Statuses)
		}
		if engine.State.PendingAction != nil || len(engine.State.PendingActionQueue) != 0 {
			t.Fatalf("draw prompt queue should be empty, pending=%+v queue=%d", engine.State.PendingAction, len(engine.State.PendingActionQueue))
		}
	})
}

func TestFriendlyDeathPendingActionsAndStateTriggersAreNotSuppressed(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]

	alice := placeUnit(baseCard(t, "4611001"), 0, 0, 0, engine)
	demonSummoner := placeUnit(baseCard(t, "1621009"), 0, 2, 0, engine)
	dead := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
	targetSkill := readySkill(baseCard(t, "3121001"), 0)
	p0.Skills[0] = targetSkill
	searchTarget := NewCardInstance(baseCard(t, "1621010"), 0, 1)
	p0.Deck = []*CardInstance{searchTarget}

	engine.destroyUnit(dead, 0)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "alice_boost_spell" {
		t.Fatalf("Alice should open the first friendly-death prompt, pending=%+v", engine.State.PendingAction)
	}
	if demonSummoner.Statuses[demonSummonerDeathReady] != 1 {
		t.Fatalf("demon summoner should still arm its death trigger while Alice prompt is pending, statuses=%v", demonSummoner.Statuses)
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{targetSkill.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve Alice prompt: %v", err)
	}
	if targetSkill.PowerBonus != 1 || alice.UsedThisTurn != 1 {
		t.Fatalf("Alice should boost the selected skill, power_bonus=%d used=%d", targetSkill.PowerBonus, alice.UsedThisTurn)
	}
	if !hasEvent(engine.log, "pending_action_cleared", "alice_boost_spell") {
		t.Fatalf("resolving Alice prompt should emit pending_action_cleared, events=%v", eventTypes(engine.log))
	}
	if engine.State.PendingAction != nil || len(engine.State.PendingActionQueue) != 0 || engine.State.Phase != PhaseMain {
		t.Fatalf("Alice prompt should clear immediately after resolution, phase=%s pending=%+v queue=%d", engine.State.Phase, engine.State.PendingAction, len(engine.State.PendingActionQueue))
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  demonSummoner.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use armed demon summoner: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "demon_summoner_search" {
		t.Fatalf("demon summoner should open its search after the queued death state was preserved, pending=%+v", engine.State.PendingAction)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{searchTarget.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve demon summoner search: %v", err)
	}
	if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != searchTarget.InstanceID || demonSummoner.Statuses[demonSummonerDeathReady] != 0 {
		t.Fatalf("demon summoner should search and clear ready mark, hand=%v statuses=%v", cardsToInfo(p0.Hand), demonSummoner.Statuses)
	}
}

func TestFriendlyDeathPendingActionsQueueWhenMultiplePromptsTrigger(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]

	placeUnit(baseCard(t, "4611001"), 0, 0, 0, engine)
	greatDruid := placeUnit(baseCard(t, "1411001"), 0, 2, 0, engine)
	dead := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
	targetSkill := readySkill(baseCard(t, "3121001"), 0)
	p0.Skills[0] = targetSkill
	if err := (Card1411001GreatDruidCycle{}).OnUltimate(&EffectContext{
		Engine:     engine,
		Source:     greatDruid,
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("arm great druid ultimate: %v", err)
	}

	engine.destroyUnit(dead, 0)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "alice_boost_spell" {
		t.Fatalf("Alice should open first by field order, pending=%+v", engine.State.PendingAction)
	}
	if len(engine.State.PendingActionQueue) != 1 || engine.State.PendingActionQueue[0].Type != "great_druid_life_seed" {
		t.Fatalf("great druid prompt should be queued behind Alice, queue=%+v", engine.State.PendingActionQueue)
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{targetSkill.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve Alice prompt: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "great_druid_life_seed" {
		t.Fatalf("great druid prompt should activate immediately after Alice clears, pending=%+v", engine.State.PendingAction)
	}
	if !hasEvent(engine.log, "pending_action_cleared", "alice_boost_spell") {
		t.Fatalf("Alice prompt should emit pending_action_cleared before next queued prompt, events=%v", eventTypes(engine.log))
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{greatDruid.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve great druid prompt: %v", err)
	}
	if !hasEvent(engine.log, "pending_action_cleared", "great_druid_life_seed") {
		t.Fatalf("great druid prompt should emit pending_action_cleared, events=%v", eventTypes(engine.log))
	}
	if engine.State.PendingAction != nil || len(engine.State.PendingActionQueue) != 0 || engine.State.Phase != PhaseMain {
		t.Fatalf("all friendly death prompts should resolve cleanly, phase=%s pending=%+v queue=%d", engine.State.Phase, engine.State.PendingAction, len(engine.State.PendingActionQueue))
	}
	seedCount := 0
	for _, unit := range engine.getAllFieldCards(p0) {
		if unit.Card.Number == "1401001" {
			seedCount++
		}
	}
	if targetSkill.PowerBonus != 1 || seedCount != 1 {
		t.Fatalf("both queued effects should resolve, power_bonus=%d life_seed_count=%d", targetSkill.PowerBonus, seedCount)
	}
}

func hasEvent(events []GameEvent, eventType string, pendingType string) bool {
	for _, event := range events {
		if event.Type != eventType {
			continue
		}
		if pendingType == "" || event.Data["type"] == pendingType {
			return true
		}
	}
	return false
}

func eventTypes(events []GameEvent) []string {
	types := make([]string, 0, len(events))
	for _, event := range events {
		types = append(types, event.Type)
	}
	return types
}
