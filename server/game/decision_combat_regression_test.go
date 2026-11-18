package game

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDecisionDamagePreviewMatchesAreaShieldBeforeAdjustments(t *testing.T) {
	for _, number := range []string{"3321002", "3321006"} {
		for _, shield := range []int{0, 1, 2} {
			t.Run(fmt.Sprintf("%s_shield_%d", number, shield), func(t *testing.T) {
				e := setupEffectTest(t)
				titan := placeUnit(baseCard(t, "1521002"), 1, 0, 0, e)
				placeUnit(baseCard(t, "1021001"), 1, 1, 0, e)
				e.State.Players[1].Shield = shield
				skill := readySkill(baseCard(t, number), 0)
				target := SpellTarget{Type: "unit", Position: *titan.Position}
				units := e.spellAffectedUnits(1, skill, target)
				if number == "3321006" && len(units) < 2 {
					t.Fatal("fixture must use an area hit")
				}
				raw := e.effectiveSpellDamage(0, skill, skill.Card.Attack, nil, units)
				before, _ := json.Marshal(e.GetStateForPlayer(1))
				impacts := e.decisionSpellImpacts(0, skill, raw, units)
				e.decisionSpellImpacts(0, skill, raw, units)
				after, _ := json.Marshal(e.GetStateForPlayer(1))
				if string(before) != string(after) {
					t.Fatal("preview consumed shield or modifiers")
				}
				life := map[string]int{}
				for _, u := range units {
					life[u.InstanceID] = u.CurrentLife
				}
				e.resolveSpellHit(0, skill, target, nil, nil)
				for _, hit := range impacts {
					for _, u := range units {
						if u.InstanceID == hit.InstanceID && life[u.InstanceID]-u.CurrentLife != hit.Damage {
							t.Fatalf("preview %d actual %d for %s", hit.Damage, life[u.InstanceID]-u.CurrentLife, u.Card.Name)
						}
					}
				}
			})
		}
	}
}

func TestDecisionHitStatusMergesEligibleTemporaryModifiersPurely(t *testing.T) {
	e := setupEffectTest(t)
	skill := readySkill(baseCard(t, "3121003"), 0)
	hero := e.State.Players[1].Hero
	for _, m := range []TemporaryModifier{
		{Type: TempModNextSpellHitStatus, RemainingUses: 1, Status: StatusBurn, Amount: 1},
		{Type: TempModCurrentTurnElementHitStatus, Element: "火", Status: StatusBurn, Amount: 2},
		{Type: TempModCurrentTurnElementHitStatus, Element: "气", Status: StatusBurn, Amount: 20},
		{Type: TempModNextSpellHitStatus, RemainingUses: 0, Status: StatusBurn, Amount: 20},
		{Type: TempModNextSpellHitStatus, RemainingUses: 1, TargetInstanceID: "different", Status: StatusBurn, Amount: 20},
	} {
		e.addTemporaryModifier(0, m)
	}
	before, _ := json.Marshal(e.State.Players[0].TempModifiers)
	for i := 0; i < 2; i++ {
		if got := e.decisionSpellImpacts(0, skill, 1, []*CardInstance{hero})[0].Burn; got != 5 {
			t.Fatalf("burn=%d, want printed 2 + temporary 3", got)
		}
	}
	after, _ := json.Marshal(e.State.Players[0].TempModifiers)
	if string(before) != string(after) || hero.Statuses[StatusBurn] != 0 {
		t.Fatal("query spent or applied hit status")
	}
}
