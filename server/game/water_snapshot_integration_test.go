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
				wantPower, wantPositions := 4, 1
				if i > 0 {
					wantPower, wantPositions = 6, 3
				}
				if spell.TotalPower != wantPower || len(e.pendingSpellAffectedPositions(spell)) != wantPositions {
					t.Fatalf("%s cast %d power/area mismatch", mode, i+1)
				}
				beforeA, beforeB := a.CurrentLife, b.CurrentLife
				waterResolve(t, e)
				if a.CurrentLife >= beforeA || (i == 0 && b.CurrentLife != beforeB) || (i > 0 && b.CurrentLife >= beforeB) {
					t.Fatal("virtual/sketch damage hit wrong set")
				}
				if e.State.PendingSpell != nil {
					t.Fatal("resolved virtual source remained pending")
				}
				for _, source := range e.getAllFieldCards(e.State.Players[0]) {
					if source == skill {
						t.Fatal("virtual source became a permanent aura")
					}
				}
			}
		})
	}
}

func TestWaterSourceModifiersThroughScrollActionsAndNestedCounters(t *testing.T) {
	for _, mode := range []string{"mirror", "sketch"} {
		t.Run(mode, func(t *testing.T) {
			e := setupReportedBugEngine(t)
			p := e.State.Players[0]
			learned := readySkill(baseCard(t, "3221109"), 0)
			if mode == "sketch" {
				p.Skills[0] = learned
			} else {
				// A prior-turn recording survives while this turn's use count starts at zero.
				p.LastLowCostWaterSpell = e.cloneVirtualSpell(learned, 0, e.State.TurnNumber)
			}
			for _, ps := range e.State.Players {
				for _, el := range []string{"火", "水", "地", "气", "光", "暗", "无"} {
					ps.Elements[el] = 100
				}
			}
			a := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
			a.CurrentLife = 100
			b := placeUnit(baseCard(t, "1021001"), 1, 1, 0, e)
			b.CurrentLife = 100
			defense := readySkill(baseCard(t, "3321002"), 1)
			e.State.Players[1].Skills[0] = defense
			fieldCount := len(e.getAllFieldCards(p))
			choose := func(player int, id string) {
				t.Helper()
				if err := e.HandleAction(player, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{id}}}); err != nil {
					t.Fatal(err)
				}
			}
			for i := 0; i < 2; i++ {
				counter := NewCardInstance(baseCard(t, "2021018"), 1, e.State.TurnNumber)
				counter.IsSetCounter = true
				e.State.Players[1].Equipment[0] = counter
				number := "2221104"
				if mode == "sketch" {
					number = "2021012"
				}
				scroll := NewCardInstance(baseCard(t, number), 0, e.State.TurnNumber)
				p.Hand = append(p.Hand, scroll)
				if err := e.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": scroll.InstanceID}}); err != nil {
					t.Fatal(err)
				}
				if mode == "sketch" {
					if e.State.PendingAction == nil || e.State.PendingAction.Type != "sketch_scroll_skill" {
						t.Fatal("missing skill choice")
					}
					choose(0, learned.InstanceID)
				}
				if e.State.PendingAction == nil {
					t.Fatal("missing target choice")
				}
				choose(0, a.InstanceID)
				spell := e.State.PendingSpell
				wantPower, wantPositions := 4, 1
				if i > 0 {
					wantPower, wantPositions = 6, 3
				}
				check := func() {
					t.Helper()
					if spell == nil || e.State.PendingSpell != spell || spell.TotalPower != wantPower || len(e.pendingSpellAffectedPositions(spell)) != wantPositions || spell.priorSameNumberCasts == nil || *spell.priorSameNumberCasts != i {
						t.Fatalf("%s cast %d lost power/area/snapshot", mode, i+1)
					}
				}
				check()
				if mode == "mirror" && spell.Skill == learned {
					t.Fatal("mirror did not create a virtual source")
				}
				if e.State.PendingAction == nil || e.State.PendingAction.Type != "counter_trigger" {
					t.Fatal("missing counter window")
				}
				choose(1, counter.InstanceID)
				if e.State.PendingAction == nil || e.State.PendingAction.Type != "arcane_rune_skill" {
					t.Fatal("missing nested choice")
				}
				choose(1, defense.InstanceID)
				check()
				beforeA, beforeB := a.CurrentLife, b.CurrentLife
				waterResolve(t, e)
				if a.CurrentLife >= beforeA || (i == 0 && b.CurrentLife != beforeB) || (i > 0 && b.CurrentLife >= beforeB) {
					t.Fatal("wrong damage set after nested continuation")
				}
				if e.State.PendingSpell != nil || e.State.PendingAction != nil || len(e.getAllFieldCards(p)) != fieldCount {
					t.Fatal("source lifetime escaped cast")
				}
				unrelated := readySkill(baseCard(t, "3321002"), 0)
				if e.spellStatBonuses(0, unrelated, skillPurposeAttack).PowerBonus != 0 || e.effectiveSpellArea(unrelated) != SpellAreaSingle {
					t.Fatal("resolved ripple leaked aura to unrelated spell")
				}
			}
		})
	}
}

func TestWaterSourceModifierOptInAndIdentity(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	ripple := readySkill(baseCard(t, "3221109"), 0)
	p.Skills[0] = ripple
	count := 0
	for _, c := range e.spellModifierSources(p, ripple) {
		if c == ripple {
			count++
		}
	}
	if count != 1 {
		t.Fatal("in-field source duplicated")
	}
	offField := readySkill(baseCard(t, "3021106"), 0)
	for _, c := range e.spellModifierSources(p, offField) {
		if c == offField {
			t.Fatal("non-opted-in off-field aura was admitted")
		}
	}
	enemy := readySkill(baseCard(t, "3221109"), 1)
	for _, c := range e.spellModifierSources(p, enemy) {
		if c == enemy {
			t.Fatal("enemy source admitted")
		}
	}
}
