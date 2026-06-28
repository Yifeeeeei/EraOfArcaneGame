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
