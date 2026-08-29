package game

import (
	"testing"

	"eraofarcane/model"
)

func TestTriggeredTurnReservationsCountPendingWindows(t *testing.T) {
	setupReportedBugEngine(t)
	card := NewCardInstance(baseCard(t, "2321001"), 0, 1)
	if !reserveTriggeredTurn(card) || !reserveTriggeredTurn(card) || !reserveTriggeredTurn(card) {
		t.Fatal("triggered turn 3 should reserve three simultaneous windows")
	}
	if reserveTriggeredTurn(card) {
		t.Fatal("a fourth simultaneous window must not exceed the per-turn limit")
	}
	if !resolveTriggeredTurn(card, true) {
		t.Fatal("accepted reserved window should commit")
	}
	if resolveTriggeredTurn(card, false) {
		t.Fatal("declined reserved window should release without committing")
	}
	if !resolveTriggeredTurn(card, true) {
		t.Fatal("a later accepted reserved window should still commit")
	}
	if card.UsedThisTurn != 2 || card.PendingTriggeredUses != 0 {
		t.Fatalf("unexpected trigger counters: used=%d pending=%d", card.UsedThisTurn, card.PendingTriggeredUses)
	}
	if !reserveTriggeredTurn(card) {
		t.Fatal("declining an optional window should leave one use available")
	}
}

func TestTriggeredTurnOncePerTurnRegressions(t *testing.T) {
	t.Run("fire insight draws once for multiple fire damage events", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		insight := placeUnit(baseCard(t, "1121012"), 0, 0, 0, engine)
		targetA := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		targetB := placeUnit(baseCard(t, "1021002"), 0, 2, 0, engine)
		engine.State.Players[0].Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1021002"), 0, 1),
		}

		data := map[string]any{"damage_element": model.ElementFire, "attacker": 1}
		engine.triggerEffects(TriggerOnDamaged, insight, targetA, data)
		engine.triggerEffects(TriggerOnDamaged, insight, targetB, data)

		if len(engine.State.Players[0].Hand) != 1 || insight.UsedThisTurn != 1 {
			t.Fatalf("fire insight should draw once, hand=%d used=%d", len(engine.State.Players[0].Hand), insight.UsedThisTurn)
		}
	})

	t.Run("soul necklace gains one shadow for simultaneous deaths", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		necklace := NewCardInstance(baseCard(t, "2621006"), 0, 1)
		p0.Equipment[0] = necklace
		first := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		second := placeUnit(baseCard(t, "1021002"), 0, 1, 0, engine)

		engine.destroyUnitWithData(first, 0, map[string]any{"attacker": 1})
		engine.destroyUnitWithData(second, 0, map[string]any{"attacker": 1})

		if p0.Elements[model.ElementShadow] != 1 || necklace.UsedThisTurn != 1 {
			t.Fatalf("soul necklace should trigger once, shadow=%d used=%d", p0.Elements[model.ElementShadow], necklace.UsedThisTurn)
		}
	})

	t.Run("word spirit releases a declined optional window", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		wordSpirit := placeUnit(baseCard(t, "1621013"), 0, 0, 0, engine)
		enemySpell := readySkill(baseCard(t, "3121001"), 1)
		enemySpell.IsHorizontal = true
		engine.State.Players[1].Skills[0] = enemySpell

		data := map[string]any{"cast_player": 1}
		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, enemySpell, data)
		resolvePendingSelection(t, engine, 0)
		if wordSpirit.UsedThisTurn != 0 || wordSpirit.PendingTriggeredUses != 0 {
			t.Fatalf("declining should release the window, used=%d pending=%d", wordSpirit.UsedThisTurn, wordSpirit.PendingTriggeredUses)
		}

		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, enemySpell, data)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "word_spirit_weaken" {
			t.Fatalf("a later eligible cast should offer the trigger again, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, wordSpirit.InstanceID)
		if enemySpell.Statuses[StatusWeaken] != 1 || wordSpirit.UsedThisTurn != 1 {
			t.Fatalf("accepted word spirit trigger should resolve once, statuses=%v used=%d", enemySpell.Statuses, wordSpirit.UsedThisTurn)
		}
	})

	t.Run("windbreath compass caps simultaneous draw windows at three", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		compass := NewCardInstance(baseCard(t, "2321001"), 0, 1)
		p0.Equipment[0] = compass
		drawn := make([]*CardInstance, 4)
		for i := range drawn {
			drawn[i] = NewCardInstance(baseCard(t, "1321001"), 0, 1)
		}
		p0.Deck = append([]*CardInstance(nil), drawn...)

		engine.drawCards(0, 4)
		if compass.PendingTriggeredUses != 3 || len(engine.State.PendingActionQueue) != 2 {
			t.Fatalf("compass should reserve exactly three prompts, pending_uses=%d queue=%d", compass.PendingTriggeredUses, len(engine.State.PendingActionQueue))
		}
		resolvePendingSelection(t, engine, 0)
		resolvePendingSelection(t, engine, 0, drawn[1].InstanceID)
		resolvePendingSelection(t, engine, 0, drawn[2].InstanceID)
		if compass.UsedThisTurn != 2 || compass.PendingTriggeredUses != 0 || engine.State.PendingAction != nil {
			t.Fatalf("compass should commit only accepted prompts, used=%d pending=%d action=%+v", compass.UsedThisTurn, compass.PendingTriggeredUses, engine.State.PendingAction)
		}
		if p0.RevealedHand[drawn[0].InstanceID] || !p0.RevealedHand[drawn[1].InstanceID] || !p0.RevealedHand[drawn[2].InstanceID] || p0.RevealedHand[drawn[3].InstanceID] {
			t.Fatalf("only accepted eligible draws should be revealed, revealed=%v", p0.RevealedHand)
		}
	})

	t.Run("ancient tree heart cannot recurse between life and load", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		tree := NewCardInstance(baseCard(t, "2411001"), 0, 1)
		p0.Equipment[0] = tree
		unit := placeUnit(baseCard(t, "1021002"), 0, 0, 0, engine)
		unit.CurrentLife = maxLife(unit) - 2
		startLife := unit.CurrentLife
		startEarth := effectiveElementsGain(unit)[model.ElementEarth]

		engine.addElementsGainBonus(unit, 0, model.ElementEarth, 1, tree)
		resolvePendingSelection(t, engine, 0, unit.InstanceID)

		if unit.CurrentLife != startLife+1 || effectiveElementsGain(unit)[model.ElementEarth] != startEarth+1 || tree.UsedThisTurn != 1 {
			t.Fatalf("tree heart should resolve only the chosen side once, life=%d earth=%d used=%d", unit.CurrentLife, effectiveElementsGain(unit)[model.ElementEarth], tree.UsedThisTurn)
		}
		if engine.State.PendingAction != nil || len(engine.State.PendingActionQueue) != 0 {
			t.Fatalf("tree heart life gain must not recursively queue its other mode, pending=%+v queue=%d", engine.State.PendingAction, len(engine.State.PendingActionQueue))
		}
	})

	t.Run("witchcraft ring augments only the first weaken event", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		ring := NewCardInstance(baseCard(t, "2621013"), 0, 1)
		engine.State.Players[0].Equipment[0] = ring
		first := readySkill(baseCard(t, "3121001"), 1)
		second := readySkill(baseCard(t, "3221009"), 1)
		engine.State.Players[1].Skills[0] = first
		engine.State.Players[1].Skills[1] = second

		engine.addStatus(first, StatusWeaken, 1)
		engine.addStatus(second, StatusWeaken, 1)

		if first.Statuses[StatusWeaken] != 2 || second.Statuses[StatusWeaken] != 1 || ring.UsedThisTurn != 1 {
			t.Fatalf("witchcraft ring should augment once, first=%d second=%d used=%d", first.Statuses[StatusWeaken], second.Statuses[StatusWeaken], ring.UsedThisTurn)
		}
	})

	t.Run("defense triggers apply at most once", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		shield := NewCardInstance(baseCard(t, "2511002"), 0, 1)
		phantomPain := readySkill(baseCard(t, "2601001"), 0)
		defenseSpell := readySkill(baseCard(t, "3121001"), 1)
		engine.State.Players[0].Equipment[0] = shield
		engine.State.Players[0].Skills[0] = phantomPain
		engine.State.Players[1].Skills[0] = defenseSpell

		ownDefense := map[string]any{"defense_success": true, "defender": 0}
		engine.triggerEffects(TriggerOnDefend, shield, nil, ownDefense)
		engine.triggerEffects(TriggerOnDefend, shield, nil, ownDefense)
		enemyDefense := map[string]any{"defense_success": true, "defender": 1, "defense_skills": []*CardInstance{defenseSpell}}
		engine.triggerEffects(TriggerOnDefend, phantomPain, nil, enemyDefense)
		engine.triggerEffects(TriggerOnDefend, phantomPain, nil, enemyDefense)

		if shield.UsedThisTurn != 1 || phantomPain.UsedThisTurn != 1 || defenseSpell.Statuses[StatusWeaken] != 2 {
			t.Fatalf("defense triggers should resolve once, shield=%d pain=%d weaken=%d", shield.UsedThisTurn, phantomPain.UsedThisTurn, defenseSpell.Statuses[StatusWeaken])
		}
	})
}

func TestTriggeredOnlyCardsAreNotSerializedAsActiveAbilities(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	for _, number := range []string{"1211003", "1321015", "1621009", "2411001", "2621013"} {
		card := NewCardInstance(baseCard(t, number), 0, 1)
		if info := engine.cardToInfoForPlayer(p0, card); info["has_per_turn"] == true {
			t.Errorf("%s is triggered-only and must not expose an active per-turn button: %v", number, info)
		}
	}
}
