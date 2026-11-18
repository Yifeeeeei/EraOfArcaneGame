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

func TestReview177TapPaymentsBeforeTriggers(t *testing.T) {
	e := setupEffectTest(t)
	attacker, defender := e.State.Players[0], e.State.Players[1]
	rune := NewCardInstance(baseCard(t, "2321002"), 0, e.State.TurnNumber)
	attacker.Hand = []*CardInstance{rune}
	attacker.Elements["气"], attacker.Elements["火"] = 9, 9
	if err := e.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": rune.InstanceID}}); err != nil {
		t.Fatal(err)
	}
	first := placeUnit(baseCard(t, "1121102"), 1, 0, 0, e)
	second := placeUnit(baseCard(t, "1121102"), 1, 1, 0, e)
	scrolls := []*CardInstance{
		NewCardInstance(baseCard(t, "2121105"), 1, e.State.TurnNumber),
		NewCardInstance(baseCard(t, "2121105"), 1, e.State.TurnNumber),
	}
	defender.Hand = scrolls
	defender.Elements["火"] = 4
	spell := readySkill(baseCard(t, "3121001"), 0)
	attacker.Skills[0] = spell
	if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": spell.InstanceID, "target_type": "unit", "target_col": float64(0), "target_row": float64(0),
	}}); err != nil {
		t.Fatal(err)
	}
	if err := e.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
		"scroll_ids": []any{scrolls[0].InstanceID, scrolls[1].InstanceID},
	}}); err != nil {
		t.Fatal(err)
	}
	for _, unit := range []*CardInstance{first, second} {
		if e.State.PendingAction == nil || e.State.PendingAction.Type != "spell_tap_cost" {
			t.Fatal("missing scroll tap-cost choice")
		}
		resolvePendingSelection(t, e, 1, unit.InstanceID)
	}
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "counter_trigger" || e.State.PendingAction.PlayerID != 0 {
		t.Fatal("missing lightning rune response to the first consumption")
	}
	allPaidBeforeCounter := first.IsHorizontal && second.IsHorizontal
	resolvePendingSelection(t, e, 0, rune.InstanceID)
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "lightning_rune_adjacent" {
		t.Fatal("missing adjacent stun choice")
	}
	resolvePendingSelection(t, e, 0, second.InstanceID)
	if !allPaidBeforeCounter || !first.IsHorizontal || !second.IsHorizontal {
		t.Fatalf("counter interrupted tap payment: paid before counter=%v first=%v second=%v", allPaidBeforeCounter, first.IsHorizontal, second.IsHorizontal)
	}
	if first.Statuses[StatusStun] != 1 || second.Statuses[StatusStun] != 1 {
		t.Fatal("rune must still stun both paid companions")
	}
	consumed := map[string]int{}
	defensePower := 0
	for _, event := range e.log {
		if event.Type == "consume" {
			id, _ := event.Data["instance_id"].(string)
			consumed[id]++
		}
		if event.Type == "defense_attempt" {
			defensePower, _ = event.Data["defense_power"].(int)
		}
	}
	if consumed[first.InstanceID] != 1 || consumed[second.InstanceID] != 1 || defensePower != 16 || e.State.Phase != PhaseMain || e.State.PendingSpell != nil {
		t.Fatalf("paid consume events or defense continuation incorrect: consumed=%v power=%d phase=%s", consumed, defensePower, e.State.Phase)
	}
	if defender.Elements["火"] != 0 || len(defender.Hand) != 0 || !containsCardInstance(attacker.Graveyard, rune) {
		t.Fatal("scroll/counter payment did not finish exactly once")
	}
	for _, scroll := range scrolls {
		if !containsCardInstance(defender.Graveyard, scroll) || scroll.spellTapPowerBonus != 0 {
			t.Fatal("used scroll or its bonus was not cleaned up")
		}
	}
}

func TestReview177AttackBoostTapPaymentsBeforeTriggers(t *testing.T) {
	e := setupEffectTest(t)
	attacker, opponent := e.State.Players[0], e.State.Players[1]
	rune := NewCardInstance(baseCard(t, "2321002"), 1, e.State.TurnNumber)
	opponent.Hand = []*CardInstance{rune}
	opponent.Elements["气"] = 9
	e.State.CurrentTurn = 1
	if err := e.HandleAction(1, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": rune.InstanceID}}); err != nil {
		t.Fatal(err)
	}
	e.State.CurrentTurn = 0
	first := placeUnit(baseCard(t, "1121102"), 0, 0, 0, e)
	second := placeUnit(baseCard(t, "1121102"), 0, 1, 0, e)
	scrolls := []*CardInstance{
		NewCardInstance(baseCard(t, "2121105"), 0, e.State.TurnNumber),
		NewCardInstance(baseCard(t, "2121105"), 0, e.State.TurnNumber),
	}
	attacker.Hand = scrolls
	attacker.Elements["火"] = 20
	spell := readySkill(baseCard(t, "3121001"), 0)
	attacker.Skills[0] = spell
	if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": spell.InstanceID, "target_type": "unit", "target_col": float64(1), "target_row": float64(1),
		"boost_ids":   []any{scrolls[0].InstanceID, scrolls[1].InstanceID},
		"consume_ids": map[string]any{scrolls[0].InstanceID: first.InstanceID, scrolls[1].InstanceID: second.InstanceID},
	}}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "counter_trigger" || !first.IsHorizontal || !second.IsHorizontal {
		t.Fatal("attack boost did not commit both tap costs before the counter")
	}
	resolvePendingSelection(t, e, 1, rune.InstanceID)
	resolvePendingSelection(t, e, 1, second.InstanceID)
	if e.State.Phase != PhaseDefenseWindow || e.State.PendingSpell == nil || e.State.PendingSpell.TotalPower != spell.Card.Power+16 || !second.IsHorizontal || second.Statuses[StatusStun] != 1 {
		t.Fatal("attack boost did not resume with both paid bonuses")
	}
	if err := e.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
		t.Fatal(err)
	}
	for _, scroll := range scrolls {
		if !containsCardInstance(attacker.Graveyard, scroll) || scroll.spellTapPowerBonus != 0 {
			t.Fatal("boost scroll cleanup failed")
		}
	}
}

func TestReview177FailedElementPaymentDoesNotCommitTapCosts(t *testing.T) {
	for _, defense := range []bool{false, true} {
		t.Run(fmt.Sprintf("defense_%v", defense), func(t *testing.T) {
			e := setupEffectTest(t)
			p := e.State.Players[0]
			first := placeUnit(baseCard(t, "1121102"), 0, 0, 0, e)
			second := placeUnit(baseCard(t, "1121102"), 0, 1, 0, e)
			scrolls := []*CardInstance{NewCardInstance(baseCard(t, "2121105"), 0, 1), NewCardInstance(baseCard(t, "2121105"), 0, 1)}
			p.Hand, p.Elements["火"] = scrolls, 10
			move := ActionMessage{Action: "cast_spell", Data: map[string]any{
				"payment":     map[string]any{"火": float64(1)},
				"consume_ids": map[string]any{scrolls[0].InstanceID: first.InstanceID, scrolls[1].InstanceID: second.InstanceID},
			}}
			spell := readySkill(baseCard(t, "3121001"), 0)
			if defense {
				e.State.CurrentTurn, e.State.Phase = 1, PhaseDefenseWindow
				e.State.PendingSpell = &SpellCast{AttackerID: 1, Skill: readySkill(baseCard(t, "3121001"), 1), Target: SpellTarget{Type: "unit", Position: Position{Col: 0, Row: 0}}, TotalPower: 3}
				move.Action, move.Data["scroll_ids"] = "defend", []any{scrolls[0].InstanceID, scrolls[1].InstanceID}
			} else {
				p.Skills[0] = spell
				move.Data["instance_id"], move.Data["target_type"], move.Data["target_col"], move.Data["target_row"] = spell.InstanceID, "unit", float64(1), float64(1)
				move.Data["boost_ids"] = []any{scrolls[0].InstanceID, scrolls[1].InstanceID}
			}
			before, _ := json.Marshal(e.State.Players)
			if err := e.HandleAction(0, move); err == nil {
				t.Fatal("invalid element payment accepted")
			}
			after, _ := json.Marshal(e.State.Players)
			if string(before) != string(after) || scrolls[0].spellTapPowerBonus != 0 || scrolls[1].spellTapPowerBonus != 0 {
				t.Fatal("failed element payment spent tap costs or granted bonuses")
			}
		})
	}
}
