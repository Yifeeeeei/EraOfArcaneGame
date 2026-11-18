package game

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFlameArrayConsumeRejectsIneligibleWithoutSpending(t *testing.T) {
	for _, condition := range []string{"horizontal", "stun", "cooldown", "cheap", "enemy"} {
		t.Run(condition, func(t *testing.T) {
			e := setupEffectTest(t)
			p := e.State.Players[0]
			scroll := NewCardInstance(baseCard(t, "2121105"), 0, 1)
			p.Hand = []*CardInstance{scroll}
			p.Elements["火"] = 2
			number, owner := "1121102", 0
			if condition == "cheap" {
				number = "1121001"
			}
			if condition == "enemy" {
				owner = 1
			}
			u := placeUnit(baseCard(t, number), owner, 0, 0, e)
			if condition == "horizontal" {
				u.IsHorizontal = true
			}
			if condition == "stun" {
				u.Statuses[StatusStun] = 1
			}
			if condition == "cooldown" {
				u.Statuses[StatusCooldown] = 1
			}
			before, _ := json.Marshal(e.State.Players)
			err := e.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
				"instance_id": scroll.InstanceID, "consume_id": u.InstanceID,
				"target_type": "unit", "target_col": float64(1), "target_row": float64(1),
			}})
			if err == nil {
				t.Fatal("invalid tap cost accepted")
			}
			after, _ := json.Marshal(e.State.Players)
			if string(before) != string(after) {
				t.Fatal("invalid tap cost spent state")
			}
		})
	}
}

func TestRebirthCounterTurnEndPaymentAndSingleUse(t *testing.T) {
	for _, ownTurn := range []bool{true, false} {
		t.Run(fmt.Sprintf("own_turn_%v", ownTurn), func(t *testing.T) {
			e := setupEffectTest(t)
			p := e.State.Players[0]
			if !ownTurn {
				e.State.CurrentTurn = 1
			}
			for _, player := range e.State.Players {
				player.Hand = nil
			}
			scroll := NewCardInstance(baseCard(t, "2121104"), 0, 1)
			scroll.IsSetCounter, scroll.IsHorizontal, scroll.SlotIndex = true, true, 0
			p.Equipment[0] = scroll
			fire := placeUnit(baseCard(t, "1121102"), 0, 0, 0, e)
			e.destroyUnitWithCause(fire, 0, DeathCauseSacrifice)
			if err := e.HandleAction(e.State.CurrentTurn, ActionMessage{Action: "end_turn"}); err != nil {
				t.Fatal(err)
			}
			if ownTurn {
				if e.State.PendingAction != nil || p.Equipment[0] != scroll || !scroll.IsSetCounter || !containsCardInstance(p.Graveyard, fire) {
					t.Fatal("a set counter resolved for free at its owner's turn end")
				}
				return
			}
			if e.State.PendingAction == nil || e.State.PendingAction.Type != "counter_trigger" {
				t.Fatal("missing enemy-turn counter window")
			}
			if err := resolvePendingSelectionWithData(e, 0, []string{scroll.InstanceID}, map[string]any{"payment": map[string]any{"火": float64(2)}}); err == nil {
				t.Fatal("revealed without paying")
			}
			if !scroll.IsSetCounter || p.Equipment[0] != scroll || !containsCardInstance(p.Graveyard, fire) {
				t.Fatal("failed payment changed cards")
			}
			p.Elements["火"] = 2
			if err := resolvePendingSelectionWithData(e, 0, []string{scroll.InstanceID}, map[string]any{"payment": map[string]any{"火": float64(2)}}); err != nil {
				t.Fatal(err)
			}
			if p.Elements["火"] != 0 || !containsCardInstance(p.Graveyard, scroll) {
				t.Fatal("counter did not pay and discard exactly once")
			}
			if e.State.PendingAction == nil || e.State.PendingAction.Type != "fire_rebirth_scroll_position" {
				t.Fatal("missing revival placement")
			}
			resolvePendingSelection(t, e, 0, "pos:0:0")
			if !sameCardInstance(p.Units[0][0], fire) || p.Units[0][0].IsHorizontal {
				t.Fatal("paid counter failed to revive/reset companion")
			}
			e.destroyUnitWithCause(p.Units[0][0], 0, DeathCauseSacrifice)
			if err := e.HandleAction(e.State.CurrentTurn, ActionMessage{Action: "end_turn"}); err != nil {
				t.Fatal(err)
			}
			if e.State.PendingAction != nil || !containsCardInstance(p.Graveyard, fire) {
				t.Fatal("discarded scroll revived a second time")
			}
		})
	}
}

func TestRebirthCandidatesRequireDeathRatherThanHandDiscard(t *testing.T) {
	e := setupEffectTest(t)
	dead := placeUnit(baseCard(t, "1121102"), 0, 0, 0, e)
	e.destroyUnitWithCause(dead, 0, DeathCauseSacrifice)
	discarded := NewCardInstance(baseCard(t, "1121102"), 0, 1)
	e.addToGraveyard(0, discarded)
	candidates := e.rebirthScrollReviveCandidates(0)
	if len(candidates) != 1 || candidates[0].InstanceID != dead.InstanceID {
		t.Fatal("hand discard was treated as a companion death")
	}
	// A card that died, returned to hand, and was discarded again is a fresh
	// graveyard entry, rather than another death at the old battlefield visit.
	card := candidates[0]
	if !e.moveGraveyardCardToHand(0, card.InstanceID) {
		t.Fatal("return failed")
	}
	for i, c := range e.State.Players[0].Hand {
		if c == card {
			e.State.Players[0].RemoveFromHand(i)
			break
		}
	}
	e.addToGraveyard(0, card)
	if len(e.rebirthScrollReviveCandidates(0)) != 0 {
		t.Fatal("stale death survived a later hand discard")
	}
}

func TestFlameArrayTapCostAppliesToAttackBoost(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	base := readySkill(baseCard(t, "3121001"), 0)
	p.Skills[0] = base
	scroll := NewCardInstance(baseCard(t, "2121105"), 0, 1)
	p.Hand = []*CardInstance{scroll}
	p.Elements["火"] = 10
	unit := placeUnit(baseCard(t, "1121102"), 0, 0, 0, e)
	move := ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": base.InstanceID, "target_type": "unit", "target_col": float64(1), "target_row": float64(1),
		"boost_ids": []any{scroll.InstanceID},
	}}
	if err := e.HandleAction(0, move); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "spell_tap_cost" || p.Elements["火"] != 10 || unit.IsHorizontal {
		t.Fatal("boost skipped additional tap selection or paid early")
	}
	resolvePendingSelection(t, e, 0, unit.InstanceID)
	want := base.Card.Power + scroll.Card.Power + totalElementCost(unit.Card.ElementsCost)
	if e.State.PendingSpell == nil || e.State.PendingSpell.TotalPower != want || !unit.IsHorizontal || p.Units[0][0] != unit {
		t.Fatalf("boost tap/bonus wrong, pending=%+v want=%d", e.State.PendingSpell, want)
	}
	if err := e.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
		t.Fatal(err)
	}
	if scroll.spellTapPowerBonus != 0 {
		t.Fatal("boost power survived spell cleanup")
	}
}

func TestFlameArrayTapCostAppliesToDefenseAndCannotDoublePay(t *testing.T) {
	for _, kind := range []string{"main", "boost", "two_scrolls", "duplicate_payment"} {
		t.Run(kind, func(t *testing.T) {
			e := setupEffectTest(t)
			p := e.State.Players[0]
			e.State.CurrentTurn = 1
			e.State.Phase = PhaseDefenseWindow
			attack := readySkill(baseCard(t, "3121001"), 1)
			e.State.PendingSpell = &SpellCast{AttackerID: 1, Skill: attack, Target: SpellTarget{Type: "unit", Position: Position{Col: 1, Row: 1}}, TotalPower: 5}
			scroll := NewCardInstance(baseCard(t, "2121105"), 0, 1)
			p.Hand = []*CardInstance{scroll}
			p.Elements["火"] = 10
			unit := placeUnit(baseCard(t, "1121102"), 0, 0, 0, e)
			data := map[string]any{"scroll_ids": []any{scroll.InstanceID}}
			want := scroll.Card.Power + totalElementCost(unit.Card.ElementsCost)
			if kind == "boost" {
				base := readySkill(baseCard(t, "3121001"), 0)
				p.Skills[0] = base
				data["skill_ids"], data["scroll_ids"], data["boost_ids"] = []any{base.InstanceID}, []any{}, []any{scroll.InstanceID}
				want += base.Card.Power
			}
			var second *CardInstance
			if kind == "two_scrolls" {
				second = placeUnit(baseCard(t, "1121102"), 0, 2, 0, e)
				other := NewCardInstance(baseCard(t, "2121105"), 0, 1)
				p.Hand = append(p.Hand, other)
				data["scroll_ids"] = []any{scroll.InstanceID, other.InstanceID}
				want *= 2
			}
			if kind == "duplicate_payment" {
				data["consume_id"], data["overexert_ids"] = unit.InstanceID, []any{unit.InstanceID}
				if err := e.HandleAction(0, ActionMessage{Action: "defend", Data: data}); err == nil {
					t.Fatal("same unit paid both tap and overexertion")
				}
				if p.Elements["火"] != 10 || unit.IsHorizontal || len(p.Hand) != 1 {
					t.Fatal("invalid cost changed state")
				}
				return
			}
			life := p.Hero.CurrentLife
			if err := e.HandleAction(0, ActionMessage{Action: "defend", Data: data}); err != nil {
				t.Fatal(err)
			}
			if e.State.PendingAction == nil || e.State.PendingAction.Type != "spell_tap_cost" {
				t.Fatal("missing defense tap prompt")
			}
			resolvePendingSelection(t, e, 0, unit.InstanceID)
			if second != nil {
				if e.State.PendingAction == nil || e.State.PendingAction.Type != "spell_tap_cost" {
					t.Fatal("second scroll reused the first tap")
				}
				resolvePendingSelection(t, e, 0, second.InstanceID)
			}
			power := 0
			for _, event := range e.log {
				if event.Type == "defense_attempt" {
					power, _ = event.Data["defense_power"].(int)
				}
			}
			if power != want || p.Hero.CurrentLife != life || !unit.IsHorizontal || p.Units[0][0] != unit || scroll.spellTapPowerBonus != 0 {
				t.Fatalf("defense tap/bonus/cleanup wrong: power=%d want=%d", power, want)
			}
		})
	}
}
