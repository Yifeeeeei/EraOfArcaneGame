package game

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWaterSnapshotSurvivesNestedCounterChoice(t *testing.T) {
	e := setupReportedBugEngine(t)
	ripple := readySkill(baseCard(t, "3221109"), 0)
	e.State.Players[0].Skills[0] = ripple
	e.State.Players[0].Elements["水"] = 100
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	target.CurrentLife = 100
	other := placeUnit(baseCard(t, "1021001"), 1, 1, 0, e)
	other.CurrentLife = 100
	defense := readySkill(baseCard(t, "3321002"), 1)
	e.State.Players[1].Skills[0] = defense
	counter := NewCardInstance(baseCard(t, "2021018"), 1, e.State.TurnNumber)
	counter.IsSetCounter = true
	e.State.Players[1].Equipment[0] = counter
	for _, el := range []string{"火", "水", "气", "地", "光", "暗", "无"} {
		e.State.Players[1].Elements[el] = 30
	}
	waterCast(t, e, ripple)
	original := e.State.PendingSpell
	if original == nil || original.priorSameNumberCasts == nil || *original.priorSameNumberCasts != 0 {
		t.Fatal("missing first-cast snapshot")
	}
	choose := func(id string) {
		t.Helper()
		if err := e.HandleAction(1, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{id}}}); err != nil {
			t.Fatal(err)
		}
	}
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "counter_trigger" {
		t.Fatal("missing counter")
	}
	choose(counter.InstanceID)
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "arcane_rune_skill" {
		t.Fatal("missing nested choice")
	}
	choose(defense.InstanceID)
	if e.State.PendingSpell != original || *original.priorSameNumberCasts != 0 || len(e.pendingSpellAffectedPositions(original)) != 1 {
		t.Fatal("nested choice changed first-cast snapshot")
	}
	raw, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "priorSameNumberCasts") || strings.Contains(string(raw), "resolutions") {
		t.Fatal("private snapshot/frames serialized")
	}
	waterResolve(t, e)
	if other.CurrentLife != 100 || target.CurrentLife >= 100 {
		t.Fatal("first cast hit wrong units after nested choice")
	}
}

func TestWaterSnapshotVirtualAndSketchEntryPoints(t *testing.T) {
	for _, mode := range []string{"virtual", "virtual_boosted", "sketch"} {
		t.Run(mode, func(t *testing.T) {
			e := setupReportedBugEngine(t)
			e.State.Players[0].Elements["水"] = 100
			a := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
			a.CurrentLife = 100
			b := placeUnit(baseCard(t, "1021001"), 1, 1, 0, e)
			b.CurrentLife = 100
			for i := 0; i < 2; i++ {
				skill := readySkill(baseCard(t, "3221109"), 0)
				target := SpellTarget{Type: "unit", Position: Position{0, 0}}
				var err error
				switch mode {
				case "virtual":
					err = e.startVirtualSpellCastNoBoost(0, skill, target, nil)
				case "virtual_boosted":
					err = e.startVirtualSpellCastWithBoosts(0, skill, target, nil, nil)
				case "sketch":
					err = e.castSkillFromSketchScroll(0, skill, target)
				}
				if err != nil {
					t.Fatal(err)
				}
				spell := e.State.PendingSpell
				if spell == nil || spell.priorSameNumberCasts == nil || *spell.priorSameNumberCasts != i {
					t.Fatalf("%s cast %d lost snapshot", mode, i)
				}
				t.Logf("%s cast=%d prior=%d power=%d affected_positions=%d", mode, i+1, *spell.priorSameNumberCasts, spell.TotalPower, len(e.pendingSpellAffectedPositions(spell)))
				// This verifies the snapshot, not virtual-card self-aura semantics.
				// Existing modifier enumeration only visits cards on the field.
				waterResolve(t, e)
			}
		})
	}
}
