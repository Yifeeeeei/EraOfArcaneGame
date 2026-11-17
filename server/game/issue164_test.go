package game

import (
	"eraofarcane/model"
	"testing"
)

func TestIssue164DrawEffectsRevealAndSerialize(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	drum := NewCardInstance(baseCard(t, "2311002"), 0, 1)
	compass := NewCardInstance(baseCard(t, "2321001"), 0, 1)
	p.Equipment[0], p.Equipment[1] = drum, compass
	for i := 0; i < 3; i++ {
		drawn := NewCardInstance(baseCard(t, "1321001"), 0, 1)
		p.Deck = []*CardInstance{drawn}
		e.drawCards(0, 1)
		if len(e.State.PendingAction.Candidates) != 1 || e.State.PendingAction.Candidates[0]["instance_id"] != drawn.InstanceID {
			t.Fatal("must offer only the drawn card")
		}
		resolvePendingSelection(t, e, 0, drawn.InstanceID)
		if e.State.PendingAction != nil {
			resolvePendingSelection(t, e, 0, drawn.InstanceID)
		}
		if !p.RevealedHand[drawn.InstanceID] {
			t.Fatal("accepted draw must reveal hand card")
		}
	}
	if got := cardToInfo(drum)["statuses"].(map[string]int)["雷鼓标记"]; got != 3 {
		t.Fatalf("serialized marks=%d", got)
	}
	if got := cardToInfo(compass)["elements_gain"].(map[string]int)[model.ElementAir]; got != compass.Card.ElementsGain[model.ElementAir]+3 {
		t.Fatalf("serialized load=%d", got)
	}
	if err := e.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{"instance_id": drum.InstanceID, "ability_type": "per_turn"}}); err != nil {
		t.Fatal(err)
	}
	resolvePendingSelection(t, e, 0, "attack")
	if thunderDrumMarks(drum) != 0 {
		t.Fatal("three marks must be spent")
	}
	if err := e.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
		t.Fatal(err)
	}
	if got := cardToInfo(compass)["elements_gain"].(map[string]int)[model.ElementAir]; got != compass.Card.ElementsGain[model.ElementAir] {
		t.Fatalf("load did not expire: %d", got)
	}
}

func TestIssue164ExtraTargetRangeAndRedMoon(t *testing.T) {
	for _, active := range []bool{false, true} {
		t.Run(map[bool]string{false: "inactive", true: "active"}[active], func(t *testing.T) {
			e := setupReportedBugEngine(t)
			p := e.State.Players[0]
			skill := readySkill(baseCard(t, "3621107"), 0)
			p.Skills[0] = skill
			p.Elements[model.ElementShadow] = 10
			moon := readySkill(baseCard(t, "3611101"), 0)
			p.Skills[1] = moon
			if active {
				moon.Statuses[redMoonMarkerStatus] = 1
				moon.Statuses[StatusAbilityDuration] = 1
			}
			front := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
			back := placeUnit(baseCard(t, "1021001"), 1, 2, 2, e)
			action := ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": skill.InstanceID, "target_type": "unit", "target_col": float64(front.Position.Col), "target_row": float64(front.Position.Row), "extra_target_col": float64(back.Position.Col), "extra_target_row": float64(back.Position.Row)}}
			err := e.HandleAction(0, action)
			if !active {
				if err == nil {
					t.Fatal("out-of-range extra target accepted")
				}
				if skill.IsHorizontal || p.Elements[model.ElementShadow] != 10 || e.State.PendingSpell != nil {
					t.Fatal("rejected target consumed state")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if e.State.PendingSpell == nil || !e.skillHasPierce(0, skill) || e.State.PendingSpell.TotalPower != skill.Card.Power+4 {
					t.Fatalf("red moon stats: %+v", e.State.PendingSpell)
				}
				delete(moon.Statuses, StatusAbilityDuration)
				if e.skillHasPierce(0, skill) {
					t.Fatal("pierce survived red moon")
				}
			}
		})
	}
}

func TestIssue164MultipleWingsRefreshTargets(t *testing.T) {
	e := setupReportedBugEngine(t)
	a := placeUnit(baseCard(t, "1621109"), 0, 0, 0, e)
	b := placeUnit(baseCard(t, "1621109"), 0, 2, 0, e)
	front := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	front.CurrentLife = 1
	back := placeUnit(baseCard(t, "1021001"), 1, 2, 2, e)
	back.CurrentLife = 5
	e.triggerScarletWingsAfterRedMoon(0)
	if len(e.State.PendingActionQueue) != 1 {
		t.Fatal("both sources must queue")
	}
	resolvePendingSelection(t, e, 0, front.InstanceID)
	if e.State.PendingAction == nil {
		t.Fatal("second wing trigger missing")
	}
	for _, c := range e.State.PendingAction.Candidates {
		if c["instance_id"] == front.InstanceID {
			t.Fatal("dead candidate retained")
		}
	}
	resolvePendingSelection(t, e, 0, back.InstanceID)
	if a.CurrentLife != 2 || b.CurrentLife != 2 || back.CurrentLife != 4 || e.State.PendingAction != nil {
		t.Fatal("both wings must finish independently")
	}
}

func TestIssue164PainScreamQueuesAndRefreshes(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	ally := placeUnit(baseCard(t, "1021002"), 0, 0, 0, e)
	p.TempModifiers = append(p.TempModifiers, TemporaryModifier{Type: TempModPainScreamWeakenOnDamage, RemainingUses: -1})
	a := readySkill(baseCard(t, "3121001"), 1)
	b := readySkill(baseCard(t, "3121002"), 1)
	e.State.Players[1].Skills[0], e.State.Players[1].Skills[1] = a, b
	e.promptPainScreamWeakenAfterFriendlyDamage(0, ally, 1)
	e.promptPainScreamWeakenAfterFriendlyDamage(0, ally, 1)
	if len(e.State.PendingActionQueue) != 1 {
		t.Fatal("second damage trigger dropped")
	}
	resolvePendingSelection(t, e, 0, a.InstanceID)
	if e.State.PendingAction == nil || len(e.State.PendingAction.Candidates) != 1 || e.State.PendingAction.Candidates[0]["instance_id"] != b.InstanceID {
		t.Fatal("queued candidates not refreshed")
	}
	resolvePendingSelection(t, e, 0, b.InstanceID)
	if a.Statuses[StatusWeaken] != 2 || b.Statuses[StatusWeaken] != 2 {
		t.Fatal("both damage events must resolve")
	}
}

func TestIssue164BoostPierceExtraTargets(t *testing.T) {
	for _, stealth := range []bool{false, true} {
		t.Run(map[bool]string{false: "boost", true: "stealth"}[stealth], func(t *testing.T) {
			e := setupReportedBugEngine(t)
			p := e.State.Players[0]
			skill := readySkill(baseCard(t, "3621107"), 0)
			boost := readySkill(baseCard(t, "3121015"), 0)
			p.Skills[0], p.Skills[1] = skill, boost
			p.Elements = map[string]int{model.ElementShadow: 10, model.ElementFire: 10, model.ElementAir: 10}
			placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
			back := placeUnit(baseCard(t, "1021001"), 1, 2, 2, e)
			if stealth {
				back.Statuses[StatusStealth] = 1
			}
			err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": skill.InstanceID, "boost_ids": []any{boost.InstanceID}, "target_type": "unit", "target_col": float64(0), "target_row": float64(0), "extra_target_col": float64(2), "extra_target_row": float64(2)}})
			if stealth && err == nil {
				t.Fatal("pierce must not bypass stealth")
			}
			if !stealth && err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIssue164WingsSelfDamageWithPrevention(t *testing.T) {
	for _, prevention := range []bool{false, true} {
		t.Run(map[bool]string{false: "without", true: "with"}[prevention], func(t *testing.T) {
			e := setupReportedBugEngine(t)
			a := placeUnit(baseCard(t, "1621109"), 0, 0, 0, e)
			b := placeUnit(baseCard(t, "1621109"), 0, 2, 0, e)
			if prevention {
				placeUnit(baseCard(t, "1221001"), 0, 1, 0, e)
			}
			target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
			target.CurrentLife = 5
			e.triggerScarletWingsAfterRedMoon(0)
			resolvePendingSelection(t, e, 0, a.InstanceID)
			if prevention {
				if e.State.PendingAction.Type != "prevent_lethal_sacrifice" {
					t.Fatal("missing damage prevention")
				}
				resolvePendingSelection(t, e, 0)
			}
			if e.State.PendingAction == nil || e.State.PendingAction.Type != "scarlet_wings_red_moon_damage" {
				t.Fatal("second wing missing")
			}
			resolvePendingSelection(t, e, 0, target.InstanceID)
			if e.State.Players[0].Units[0][0] != a || a.CurrentLife != 1 || b.CurrentLife != 2 {
				t.Fatalf("inconsistent damage/life: a=%d b=%d", a.CurrentLife, b.CurrentLife)
			}
		})
	}
}

func TestIssue164PainScreamRequiresAllDamageChoices(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	p.TempModifiers = append(p.TempModifiers, TemporaryModifier{Type: TempModPainScreamWeakenOnDamage, RemainingUses: -1})
	ally := placeUnit(baseCard(t, "1021002"), 0, 0, 0, e)
	a := readySkill(baseCard(t, "3121001"), 1)
	b := readySkill(baseCard(t, "3121002"), 1)
	e.State.Players[1].Skills[0], e.State.Players[1].Skills[1] = a, b
	e.promptPainScreamWeakenAfterFriendlyDamage(0, ally, 2)
	e.promptPainScreamWeakenAfterFriendlyDamage(0, ally, 2)
	if err := e.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{a.InstanceID}}}); err == nil {
		t.Fatal("must select both legal spells")
	}
	resolvePendingSelection(t, e, 0, a.InstanceID, b.InstanceID)
	if e.State.PendingAction != nil || len(e.State.PendingActionQueue) != 0 {
		t.Fatal("empty queued choice must be skipped")
	}
}
