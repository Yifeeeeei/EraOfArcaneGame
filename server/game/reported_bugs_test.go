package game

import (
	"testing"

	"eraofarcane/cards"
	"eraofarcane/model"
)

func setupReportedBugEngine(t *testing.T) *Engine {
	t.Helper()
	if cards.CardDB == nil {
		if err := cards.LoadCards(); err != nil {
			t.Fatalf("load cards: %v", err)
		}
	}
	SetCardDB(cards.CardDB)
	RegisterAllCardEffects()

	engine := NewEngine("reported-bugs", nil)
	engine.State.Players[0] = NewPlayerState(0, "P1", &model.Deck{})
	engine.State.Players[1] = NewPlayerState(1, "P2", &model.Deck{})
	engine.State.Phase = PhaseMain
	engine.State.CurrentTurn = 0
	engine.State.TurnNumber = 1
	return engine
}

func TestIssue31PlaytestRegressions(t *testing.T) {
	t.Run("1221014 北海飞鱼 temporary load resets without dying", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		fish := placeUnit(baseCard(t, "1221014"), 0, 1, 1, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  fish.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use fish ability: %v", err)
		}
		if got := effectiveElementsGain(fish)[model.ElementAir]; got != 1 {
			t.Fatalf("fish should temporarily become 1 air load, gain=%v", effectiveElementsGain(fish))
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
			t.Fatalf("end turn: %v", err)
		}
		if p0.Units[1][1] != fish {
			t.Fatalf("fish should remain on board after end turn")
		}
		if got := effectiveElementsGain(fish)[model.ElementWater]; got != fish.Card.ElementsGain[model.ElementWater] {
			t.Fatalf("fish temporary load should reset to printed load, gain=%v", effectiveElementsGain(fish))
		}
	})

	t.Run("1321001 渡鸦信使 is labeled as consume, not per-turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		raven := placeUnit(baseCard(t, "1321001"), 0, 1, 1, engine)
		info := cardToInfo(raven)

		if info["per_turn_label"] != "消耗" {
			t.Fatalf("raven messenger should expose consume label, info=%v", info)
		}
	})

	t.Run("3221007 水占术 searches without shuffle and orders top and bottom", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		topA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		waterCard := NewCardInstance(baseCard(t, "1221001"), 0, 1)
		topB := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		topC := NewCardInstance(baseCard(t, "1021004"), 0, 1)
		rest := NewCardInstance(baseCard(t, "1021005"), 0, 1)
		p0.Deck = []*CardInstance{topA, waterCard, topB, topC, rest}
		p0.Skills[0] = readySkill(baseCard(t, "3221007"), 0)
		p0.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast water divination: %v", err)
		}
		if engine.State.PendingAction == nil || len(engine.State.PendingAction.Candidates) != 4 {
			t.Fatalf("water divination should reveal all top four cards, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected":     []any{waterCard.InstanceID},
			"top_order":    []any{topC.InstanceID, topA.InstanceID},
			"bottom_order": []any{topB.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve water divination: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != waterCard {
			t.Fatalf("water divination should search water card to hand, hand=%v", cardsToInfo(p0.Hand))
		}
		wantDeck := []*CardInstance{topC, topA, rest, topB}
		if len(p0.Deck) != len(wantDeck) {
			t.Fatalf("deck length mismatch, deck=%v", cardsToInfo(p0.Deck))
		}
		for i, want := range wantDeck {
			if p0.Deck[i] != want {
				t.Fatalf("deck order mismatch at %d, got=%s want=%s deck=%v", i, p0.Deck[i].Card.Name, want.Card.Name, cardsToInfo(p0.Deck))
			}
		}
	})
}

func TestIssue33PlaytestRegressions(t *testing.T) {
	t.Run("ice dissolve can reduce attack power to zero and defender can confirm with no skill", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121002"), 0)
		p1.Skills[0] = readySkill(baseCard(t, "3221008"), 1)
		p0.Elements[model.ElementFire] = 10
		p1.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast arcane arrow: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "react_spell", Data: map[string]any{
			"instance_id": p1.Skills[0].InstanceID,
		}}); err != nil {
			t.Fatalf("react with ice dissolve: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != 0 {
			t.Fatalf("ice dissolve should zero pending power, pending=%+v", engine.State.PendingSpell)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{},
		}}); err != nil {
			t.Fatalf("confirm zero-power defense: %v", err)
		}
		if engine.State.Phase != PhaseMain || engine.State.PendingSpell != nil {
			t.Fatalf("zero-power defense should close defense window, phase=%v pending=%+v", engine.State.Phase, engine.State.PendingSpell)
		}
		if target.CurrentLife != target.Card.Life {
			t.Fatalf("zero-power defense should prevent spell hit, life=%d", target.CurrentLife)
		}
	})

	t.Run("weakening potion and weakening curse add weaken to enemy skills", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		enemyA := readySkill(baseCard(t, "3021005"), 1)
		enemyB := readySkill(baseCard(t, "3121002"), 1)
		p1.Skills[0] = enemyA
		p1.Skills[1] = enemyB
		potion := NewCardInstance(baseCard(t, "2621001"), 0, 1)

		engine.triggerEffects(TriggerOnUseItem, potion, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "weakening_potion" {
			t.Fatalf("weakening potion should ask for enemy skills, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{enemyA.InstanceID, enemyB.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve weakening potion: %v", err)
		}
		if enemyA.Statuses[StatusWeaken] != 2 || enemyB.Statuses[StatusWeaken] != 2 {
			t.Fatalf("weakening potion should weaken two enemy skills, a=%v b=%v", enemyA.Statuses, enemyB.Statuses)
		}

		curse := readySkill(baseCard(t, "3621009"), 0)
		p0.Skills[0] = curse
		p0.Elements[model.ElementShadow] = 10
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": curse.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast weakening curse: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "weakening_curse" {
			t.Fatalf("weakening curse should ask for an enemy skill, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{enemyA.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve weakening curse: %v", err)
		}
		if enemyA.Statuses[StatusWeaken] != 4 {
			t.Fatalf("weakening curse should add weaken 2, statuses=%v", enemyA.Statuses)
		}
	})

	t.Run("public pending spell includes attacker boost skills", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		main := readySkill(baseCard(t, "3121002"), 0)
		boost := readySkill(baseCard(t, "3121015"), 0)
		p0.Skills[0] = main
		p0.Skills[1] = boost
		p0.Elements[model.ElementFire] = 10
		p0.Elements[model.ElementArcane] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": main.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
			"boost_ids":   []any{boost.InstanceID},
		}}); err != nil {
			t.Fatalf("cast boosted spell: %v", err)
		}
		state := engine.GetStateForPlayer(1)
		pending, _ := state["pending_spell"].(map[string]any)
		boosts, _ := pending["boost_skills"].([]map[string]any)
		if len(boosts) != 1 || boosts[0]["number"] != boost.Card.Number {
			t.Fatalf("opponent state should reveal boost skills, pending=%v", pending)
		}
	})
}

func TestIssue25PlaytestRegressions(t *testing.T) {
	t.Run("square spell area affects a 2x2 block, not the whole board", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		skill := readySkill(baseCard(t, "3121005"), 0)
		inA := placeUnit(baseCard(t, "1021001"), 1, 1, 1, engine)
		inB := placeUnit(baseCard(t, "1021001"), 1, 2, 2, engine)
		out := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		units := engine.spellAffectedUnits(1, skill, SpellTarget{Type: "unit", Position: Position{Col: 1, Row: 1}})

		seen := map[*CardInstance]bool{}
		for _, unit := range units {
			seen[unit] = true
		}
		if !seen[inA] || !seen[inB] || seen[out] || len(units) != 2 {
			t.Fatalf("square should be anchored 2x2, units=%v", cardsToInfo(units))
		}
	})

	t.Run("only current player's marks settle at that player's turn end", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0Unit := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		p1Unit := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		p0Unit.Statuses[StatusBurn] = 1
		p1Unit.Statuses[StatusBurn] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
			t.Fatalf("end p0 turn: %v", err)
		}
		if p0Unit.CurrentLife != p0Unit.Card.Life-1 {
			t.Fatalf("current player's burn should settle, life=%d", p0Unit.CurrentLife)
		}
		if p1Unit.CurrentLife != p1Unit.Card.Life {
			t.Fatalf("opponent burn should not settle on p0 end, life=%d", p1Unit.CurrentLife)
		}
	})

	t.Run("burn damage at turn end still triggers fire insight", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1121012"), 0, 0, 0, engine)
		verland := placeUnit(baseCard(t, "4111002"), 0, 1, 1, engine)
		verland.Statuses[StatusBurn] = 1
		draw := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Deck = []*CardInstance{draw}

		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
			t.Fatalf("end turn with burn: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != draw {
			t.Fatalf("fire insight should draw from turn-end burn damage, hand=%v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("ice dissolve sets pending spell power to zero", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		ice := readySkill(baseCard(t, "3221008"), 1)
		spell := &SpellCast{AttackerID: 0, TotalPower: 3}
		if err := (Card3221008IceDissolve{}).OnSpellReaction(&EffectContext{
			Engine: engine, Source: ice, PlayerID: 1, OpponentID: 0,
		}, spell); err != nil {
			t.Fatalf("ice dissolve reaction: %v", err)
		}
		if spell.TotalPower != 0 {
			t.Fatalf("ice dissolve should zero power, got %d", spell.TotalPower)
		}
	})

	t.Run("firethorn deathrattle asks for a burn target", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		firethorn := placeUnit(baseCard(t, "1121014"), 0, 0, 0, engine)
		enemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		engine.destroyUnit(firethorn, 0)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "firethorn_death_burn" {
			t.Fatalf("firethorn should open target selection, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{enemy.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve firethorn: %v", err)
		}
		if enemy.Statuses[StatusBurn] != 1 {
			t.Fatalf("selected enemy should get burn, statuses=%v", enemy.Statuses)
		}
	})

	t.Run("fire arrow is an equipment sacrifice ability", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		arrow := NewCardInstance(baseCard(t, "2121004"), 0, 1)
		p0.Equipment[0] = arrow
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  arrow.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use fire arrow sacrifice: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "fire_arrow_damage" {
			t.Fatalf("fire arrow should ask target, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve fire arrow: %v", err)
		}
		if p0.Equipment[0] != nil || target.CurrentLife != target.Card.Life-1 {
			t.Fatalf("fire arrow should sacrifice and damage, equipment=%v life=%d", p0.Equipment[0], target.CurrentLife)
		}
	})

	t.Run("scorching scroll asks for a target and applies burn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2121003"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementArcane] = 1
		p0.Elements[model.ElementFire] = 1
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use scorching scroll: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "scorching_scroll_burn" {
			t.Fatalf("scorching scroll should ask target, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve scorching scroll: %v", err)
		}
		if target.Statuses[StatusBurn] != 1 {
			t.Fatalf("scorching scroll should burn selected target, statuses=%v", target.Statuses)
		}
	})

	t.Run("bound fire breath can be reset and cast from its host", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		dragon := placeUnit(baseCard(t, "1111001"), 0, 1, 1, engine)
		engine.triggerEffects(TriggerOnEnter, dragon, nil, nil)
		if len(dragon.BoundSkills) != 1 {
			t.Fatalf("fire dragon should bind one skill, bound=%v", cardsToInfo(dragon.BoundSkills))
		}
		breath := dragon.BoundSkills[0]
		engine.resetCards(p0)
		p0.Elements[model.ElementFire] = 1
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": breath.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast bound fire breath: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.Skill != breath {
			t.Fatalf("bound fire breath should enter defense window, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("learning into a full skill area replaces only a vertical skill", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		for i := 0; i < 5; i++ {
			p0.Skills[i] = readySkill(baseCard(t, "3121001"), 0)
		}
		replace := p0.Skills[2]
		replace.IsHorizontal = false
		p0.Skills[1].IsHorizontal = true
		newSkill := NewCardInstance(baseCard(t, "3221007"), 0, 1)
		p0.SkillPool = []*CardInstance{newSkill}
		p0.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "learn_skill", Data: map[string]any{
			"instance_id": newSkill.InstanceID,
			"replace_id":  replace.InstanceID,
		}}); err != nil {
			t.Fatalf("learn with replacement: %v", err)
		}
		if p0.Skills[2] != newSkill || len(p0.Graveyard) != 1 || p0.Graveyard[0] != replace {
			t.Fatalf("new skill should replace selected vertical skill, skill=%v grave=%v", p0.Skills[2], cardsToInfo(p0.Graveyard))
		}
	})
}

func TestHighRiskWaterAndAirCompanionSemantics(t *testing.T) {
	t.Run("1221008 冰域恶魔 freezes every enemy field unit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		demon := placeUnit(baseCard(t, "1221008"), 0, 1, 1, engine)
		enemyA := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		enemyB := placeUnit(baseCard(t, "1021002"), 1, 2, 1, engine)

		engine.triggerEffects(TriggerOnEnter, demon, nil, nil)

		if enemyA.Statuses[StatusFreeze] != 1 || enemyB.Statuses[StatusFreeze] != 1 {
			t.Fatalf("icefield demon should freeze all enemy units, a=%v b=%v", enemyA.Statuses, enemyB.Statuses)
		}
	})

	t.Run("1221004 寒霜傀儡 freezes selected target or enemy front row fallback", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		puppet := placeUnit(baseCard(t, "1221004"), 0, 1, 1, engine)
		selected := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		front := placeUnit(baseCard(t, "1021002"), 1, 2, 1, engine)

		engine.triggerEffects(TriggerOnEnter, puppet, selected, nil)
		if selected.Statuses[StatusFreeze] != 1 || front.Statuses[StatusFreeze] != 0 {
			t.Fatalf("selected target should be frozen exclusively, selected=%v front=%v", selected.Statuses, front.Statuses)
		}

		engine.triggerEffects(TriggerOnEnter, puppet, nil, nil)
		if selected.Statuses[StatusFreeze] != 2 {
			t.Fatalf("fallback should freeze first enemy front-row unit, selected=%v", selected.Statuses)
		}
	})

	t.Run("1221005 西境海妖 exhausts a chosen friendly companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		siren := placeUnit(baseCard(t, "1221005"), 0, 1, 1, engine)
		target := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  siren.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use western siren: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "western_siren_consume" {
			t.Fatalf("western siren should ask for a companion to exhaust, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve western siren: %v", err)
		}
		if !target.IsHorizontal {
			t.Fatalf("western siren should turn selected companion horizontal")
		}
	})

	t.Run("1211001 人鱼菲尔 searches water companion only when not adjacent to companions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		phil := placeUnit(baseCard(t, "1211001"), 0, 1, 1, engine)
		waterCompanion := NewCardInstance(baseCard(t, "1221001"), 0, 1)
		nonWater := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Deck = []*CardInstance{nonWater, waterCompanion}

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  phil.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use mermaid phil: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "mermaid_search_water_companion" {
			t.Fatalf("phil should search water companion, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{waterCompanion.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve mermaid phil: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != waterCompanion.InstanceID {
			t.Fatalf("phil should move selected water companion to hand, hand=%+v", cardsToInfo(p0.Hand))
		}

		placeUnit(baseCard(t, "1021002"), 0, 1, 0, engine)
		phil.IsHorizontal = false
		phil.UsedThisTurn = 0
		engine.State.Phase = PhaseMain
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  phil.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use mermaid phil with adjacent companion: %v", err)
		}
		if engine.State.PendingAction != nil {
			t.Fatalf("phil should not search while adjacent to a companion, pending=%+v", engine.State.PendingAction)
		}
	})

	t.Run("1221006 水栖狸猫 gains water load when adjacent to two water companions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		tanuki := placeUnit(baseCard(t, "1221006"), 0, 1, 1, engine)
		placeUnit(baseCard(t, "1221001"), 0, 0, 1, engine)
		placeUnit(baseCard(t, "1221003"), 0, 1, 0, engine)

		engine.triggerEffects(TriggerOnTurnStart, tanuki, nil, nil)

		if tanuki.ElementsGainBonus[model.ElementWater] != 1 {
			t.Fatalf("aquatic tanuki should gain +1 water load, bonuses=%v", tanuki.ElementsGainBonus)
		}
	})

	t.Run("1221015 眺望者商舰 searches water card then shuffles a hand card back", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		ship := placeUnit(baseCard(t, "1221015"), 0, 1, 1, engine)
		waterCard := NewCardInstance(baseCard(t, "1221001"), 0, 1)
		handCard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Deck = []*CardInstance{waterCard}
		p0.Hand = []*CardInstance{handCard}

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  ship.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use merchant ship: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "merchant_ship_search" {
			t.Fatalf("merchant ship should first search water card, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{waterCard.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve merchant ship search: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "merchant_ship_shuffle_hand" {
			t.Fatalf("merchant ship should then require hand card shuffle, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{handCard.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve merchant ship shuffle: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != waterCard.InstanceID || len(p0.Deck) != 1 || p0.Deck[0].InstanceID != handCard.InstanceID {
			t.Fatalf("merchant ship should keep searched card and put chosen hand card into deck, hand=%+v deck=%+v", cardsToInfo(p0.Hand), cardsToInfo(p0.Deck))
		}
	})

	t.Run("1321002 随风旅行者 gains air on enter and draws on death", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		traveler := placeUnit(baseCard(t, "1321002"), 0, 1, 1, engine)
		draw := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Deck = []*CardInstance{draw}

		engine.triggerEffects(TriggerOnEnter, traveler, nil, nil)
		if p0.Elements[model.ElementAir] != 2 {
			t.Fatalf("wind traveler should gain 2 air on enter, elements=%v", p0.Elements)
		}
		engine.destroyUnit(traveler, 0)
		if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != draw.InstanceID {
			t.Fatalf("wind traveler deathrattle should draw 1, hand=%+v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("1321004 雷电元素 stuns selected enemy on enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		elemental := placeUnit(baseCard(t, "1321004"), 0, 1, 1, engine)
		target := placeUnit(baseCard(t, "1011002"), 1, 1, 0, engine)

		engine.triggerEffects(TriggerOnEnter, elemental, target, nil)

		if target.Statuses[StatusStun] != 1 {
			t.Fatalf("lightning elemental should stun selected target, statuses=%v", target.Statuses)
		}
	})

	t.Run("1221014 北海飞鱼 changes its load to one air", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		fish := placeUnit(baseCard(t, "1221014"), 0, 1, 1, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  fish.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use north sea flying fish: %v", err)
		}
		if effectiveElementsGain(fish)[model.ElementAir] != 1 || totalLoad(fish) != 1 {
			t.Fatalf("north sea flying fish should become load 1 air, load=%v", effectiveElementsGain(fish))
		}
	})

	t.Run("1311001 大鹏 draws low-cost cards from the top eight and marks them for turn-end discard", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		roc := placeUnit(baseCard(t, "1311001"), 0, 1, 1, engine)
		lowA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		high := NewCardInstance(baseCard(t, "1021013"), 0, 1)
		lowB := NewCardInstance(baseCard(t, "1021004"), 0, 1)
		p0.Deck = []*CardInstance{lowA, high, lowB}

		engine.triggerEffects(TriggerOnEnter, roc, nil, nil)

		if len(p0.Hand) != 2 || p0.Hand[0].InstanceID != lowA.InstanceID || p0.Hand[1].InstanceID != lowB.InstanceID {
			t.Fatalf("roc should draw low-cost cards from top eight, hand=%+v", cardsToInfo(p0.Hand))
		}
		if len(p0.Deck) != 1 || p0.Deck[0].InstanceID != high.InstanceID {
			t.Fatalf("roc should leave high-cost cards in deck, deck=%+v", cardsToInfo(p0.Deck))
		}
		if !p0.DiscardAtTurnEnd[lowA.InstanceID] || !p0.DiscardAtTurnEnd[lowB.InstanceID] {
			t.Fatalf("roc drawn cards should be marked for discard, marks=%v", p0.DiscardAtTurnEnd)
		}
	})
}

func TestHighRiskFireLightShadowCompanionSemantics(t *testing.T) {
	t.Run("1121012 火焰洞察者 draws when a unit takes fire damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1121012"), 0, 0, 0, engine)
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		draw := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		p0.Deck = []*CardInstance{draw}

		engine.dealDamageWithExtra(target, 1, 1, map[string]any{"damage_element": model.ElementAir})
		if len(p0.Hand) != 0 {
			t.Fatalf("fire insight should not draw for non-fire damage, hand=%+v", cardsToInfo(p0.Hand))
		}
		engine.dealDamageWithExtra(target, 1, 1, map[string]any{"damage_element": model.ElementFire})
		if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != draw.InstanceID {
			t.Fatalf("fire insight should draw for fire damage, hand=%+v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("1121013 纵火者 burns an enemy after friendly fire spell cast", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		arsonist := placeUnit(baseCard(t, "1121013"), 0, 0, 0, engine)
		enemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		fireSkill := readySkill(baseCard(t, "3121001"), 0)
		airSkill := readySkill(baseCard(t, "3321005"), 0)

		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, airSkill, map[string]any{"cast_player": 0})
		if enemy.Statuses[StatusBurn] != 0 {
			t.Fatalf("arsonist should ignore non-fire friendly spells, enemy=%v source=%v", enemy.Statuses, arsonist.Statuses)
		}
		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, fireSkill, map[string]any{"cast_player": 0})
		if enemy.Statuses[StatusBurn] != 1 {
			t.Fatalf("arsonist should burn an enemy after fire spell, enemy=%v", enemy.Statuses)
		}
	})

	t.Run("1121014 火荆 deathrattle burns an enemy", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		firethorn := placeUnit(baseCard(t, "1121014"), 0, 0, 0, engine)
		enemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		engine.destroyUnit(firethorn, 0)

		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{enemy.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve firethorn: %v", err)
		}
		if enemy.Statuses[StatusBurn] != 1 {
			t.Fatalf("firethorn deathrattle should burn an enemy, enemy=%v", enemy.Statuses)
		}
	})

	t.Run("1421007 高地泰坦 takes extra damage only from unboosted spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		titan := placeUnit(baseCard(t, "1421007"), 1, 1, 0, engine)
		titan.CurrentLife = 6

		engine.dealDamageWithExtra(titan, 1, 1, map[string]any{"damage_source": "attack"})
		if titan.CurrentLife != 5 {
			t.Fatalf("attack damage should not be amplified, life=%d", titan.CurrentLife)
		}
		engine.dealDamageWithExtra(titan, 1, 1, map[string]any{"damage_source": "spell", "boost_count": 1})
		if titan.CurrentLife != 4 {
			t.Fatalf("boosted spell damage should not be amplified, life=%d", titan.CurrentLife)
		}
		engine.dealDamageWithExtra(titan, 1, 1, map[string]any{"damage_source": "spell", "boost_count": 0})
		if titan.CurrentLife != 2 {
			t.Fatalf("unboosted spell damage should be amplified by 1, life=%d", titan.CurrentLife)
		}
	})

	t.Run("1511001 白袍大贤者 steals an enemy companion with ultimate", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		sage := placeUnit(baseCard(t, "1511001"), 0, 1, 1, engine)
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  sage.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use white robe sage ultimate: %v", err)
		}

		if target.OwnerID != 0 || target.Position == nil || engine.State.Players[1].Units[1][0] != nil {
			t.Fatalf("sage should take control of target, owner=%d pos=%v enemy_slot=%v", target.OwnerID, target.Position, engine.State.Players[1].Units[1][0])
		}
		if engine.State.Players[0].Units[target.Position.Col][target.Position.Row] != target {
			t.Fatalf("stolen target should be on friendly board")
		}
	})

	t.Run("1521007 虹之天使 lets light pay non-light skill costs", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1521007"), 0, 0, 0, engine)
		waterSkill := readySkill(baseCard(t, "3221001"), 0)

		cost := engine.effectiveSkillUseCost(p0, waterSkill)

		if cost[model.ElementWater] != 0 || cost[model.ElementLight] == 0 {
			t.Fatalf("rainbow angel should convert non-light skill cost to light, cost=%v", cost)
		}
	})

	t.Run("1521014 and 1521015 witches burn themselves on enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		torch := placeUnit(baseCard(t, "1521014"), 0, 0, 0, engine)
		ember := placeUnit(baseCard(t, "1521015"), 0, 1, 0, engine)

		engine.triggerEffects(TriggerOnEnter, torch, nil, nil)
		engine.triggerEffects(TriggerOnEnter, ember, nil, nil)

		if torch.Statuses[StatusBurn] != 2 || ember.Statuses[StatusBurn] != 3 {
			t.Fatalf("witches should enter burning, torch=%v ember=%v", torch.Statuses, ember.Statuses)
		}
	})

	t.Run("1611001 观察者 draws one and damages own hero", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		hero := placeUnit(baseCard(t, "4611001"), 0, 1, 1, engine)
		p0.Hero = hero
		okoru := placeUnit(baseCard(t, "1611001"), 0, 0, 0, engine)
		draw := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Deck = []*CardInstance{draw}

		engine.triggerEffects(TriggerOnEnter, okoru, nil, nil)

		if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != draw.InstanceID {
			t.Fatalf("observer should draw one, hand=%+v", cardsToInfo(p0.Hand))
		}
		if hero.CurrentLife != hero.Card.Life-1 {
			t.Fatalf("observer should damage own hero by 1, life=%d", hero.CurrentLife)
		}
	})

	t.Run("1611003 穿心人 adds phantom pain to hand", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		piercer := placeUnit(baseCard(t, "1611003"), 0, 0, 0, engine)

		engine.triggerEffects(TriggerOnEnter, piercer, nil, nil)

		if len(p0.Hand) != 1 || p0.Hand[0].Card.Number != "2601001" {
			t.Fatalf("heart piercer should add phantom pain to hand, hand=%+v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("1621001 冥界信鸽 draws on death", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		pigeon := placeUnit(baseCard(t, "1621001"), 0, 0, 0, engine)
		draw := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Deck = []*CardInstance{draw}

		engine.destroyUnit(pigeon, 0)

		if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != draw.InstanceID {
			t.Fatalf("underworld pigeon should draw on death, hand=%+v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("1621005 诅咒魔像 weakens the first enemy skill on enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p1 := engine.State.Players[1]
		p1.Skills[0] = readySkill(baseCard(t, "3021005"), 1)
		golem := placeUnit(baseCard(t, "1621005"), 0, 0, 0, engine)

		engine.triggerEffects(TriggerOnEnter, golem, nil, nil)

		if p1.Skills[0].Statuses[StatusWeaken] != 2 {
			t.Fatalf("cursed golem should weaken enemy skill by 2, skill=%v", p1.Skills[0].Statuses)
		}
	})

	t.Run("1621009 唤魔邪术士 searches shadow construct or demon after friendly death", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		summoner := placeUnit(baseCard(t, "1621009"), 0, 0, 0, engine)
		ally := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		demon := NewCardInstance(baseCard(t, "1621010"), 0, 1)
		p0.Deck = []*CardInstance{demon}

		engine.destroyUnit(ally, 0)
		if summoner.Statuses[demonSummonerDeathReady] != 1 {
			t.Fatalf("summoner should arm search after friendly death, statuses=%v", summoner.Statuses)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  summoner.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use demon summoner: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "demon_summoner_search" {
			t.Fatalf("demon summoner should open search, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{demon.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve demon summoner search: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != demon.InstanceID || summoner.Statuses[demonSummonerDeathReady] != 0 {
			t.Fatalf("demon summoner should search demon and clear ready mark, hand=%+v statuses=%v", cardsToInfo(p0.Hand), summoner.Statuses)
		}
	})
}

func TestHighRiskRemainingCompanionActivesAndAuras(t *testing.T) {
	t.Run("1211002 深渊巨口利维坦 consumes to destroy an enemy companion and sets cooldown", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		leviathan := placeUnit(baseCard(t, "1211002"), 0, 1, 1, engine)
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{
			"instance_id": leviathan.InstanceID,
		}}); err != nil {
			t.Fatalf("consume leviathan: %v", err)
		}
		if engine.State.Players[1].Units[1][0] != nil || len(engine.State.Players[1].Graveyard) != 1 || engine.State.Players[1].Graveyard[0].InstanceID != target.InstanceID {
			t.Fatalf("leviathan should destroy enemy companion, units=%v grave=%+v", engine.State.Players[1].Units[1][0], cardsToInfo(engine.State.Players[1].Graveyard))
		}
		if leviathan.Statuses["利维坦冷却"] != 1 {
			t.Fatalf("leviathan should set cooldown marker, statuses=%v", leviathan.Statuses)
		}
		engine.triggerEffects(TriggerOnTurnStart, leviathan, nil, nil)
		if leviathan.Statuses["利维坦冷却"] != 0 {
			t.Fatalf("leviathan cooldown should clear on owner turn start, statuses=%v", leviathan.Statuses)
		}
	})

	t.Run("1211003 雪女 has taunt/global range marker and freezes enemy with limited per-turn ability", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		snow := placeUnit(baseCard(t, "1211003"), 0, 1, 1, engine)
		enemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		engine.triggerEffects(TriggerOnEnter, snow, nil, nil)
		if snow.Statuses["引魔"] != 1 || !traitsForCardNumber("1211003").taunt || traitsForCardNumber("1211003").perTurnLimit != 3 {
			t.Fatalf("snow woman should expose taunt/global-range/per-turn-3 markers, traits=%+v statuses=%v", traitsForCardNumber("1211003"), snow.Statuses)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  snow.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use snow woman: %v", err)
		}
		if enemy.Statuses[StatusFreeze] != 1 {
			t.Fatalf("snow woman should freeze an enemy, statuses=%v", enemy.Statuses)
		}
	})

	t.Run("1311003 风刃 makes non-piercing air skill cost one extra air", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1311003"), 0, 0, 0, engine)
		nonPierceAir := readySkill(baseCard(t, "3321005"), 0)
		pierceAir := readySkill(baseCard(t, "3321009"), 0)

		nonPierceCost := engine.effectiveSkillUseCost(p0, nonPierceAir)
		pierceCost := engine.effectiveSkillUseCost(p0, pierceAir)

		if nonPierceCost[model.ElementAir] != skillUseCost(nonPierceAir.Card)[model.ElementAir]+1 {
			t.Fatalf("karina should add one air to non-piercing air skill, cost=%v", nonPierceCost)
		}
		if pierceCost[model.ElementAir] != skillUseCost(pierceAir.Card)[model.ElementAir] {
			t.Fatalf("karina should not tax already-piercing air skill, cost=%v base=%v", pierceCost, skillUseCost(pierceAir.Card))
		}
	})

	t.Run("1321013 传送法师 moves a friendly companion to an empty position", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		mage := placeUnit(baseCard(t, "1321013"), 0, 1, 1, engine)
		ally := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		old := *ally.Position

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  mage.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use teleport mage: %v", err)
		}

		if ally.Position == nil || (ally.Position.Col == old.Col && ally.Position.Row == old.Row) || engine.State.Players[0].Units[old.Col][old.Row] != nil {
			t.Fatalf("teleport mage should move ally from old slot, old=%+v new=%+v", old, ally.Position)
		}
	})

	t.Run("1321015 风语者 gains one air with current active implementation", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		speaker := placeUnit(baseCard(t, "1321015"), 0, 1, 1, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  speaker.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use wind speaker: %v", err)
		}
		if p0.Elements[model.ElementAir] != 1 {
			t.Fatalf("wind speaker should gain one air, elements=%v", p0.Elements)
		}
	})

	t.Run("1411003 沙之魔巫 gives single-target earth spells square area", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		placeUnit(baseCard(t, "1411003"), 0, 0, 0, engine)
		earthSkill := readySkill(baseCard(t, "3421001"), 0)
		earthSkill.OwnerID = 0
		targetA := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		targetB := placeUnit(baseCard(t, "1021002"), 1, 2, 1, engine)

		affected := engine.spellAffectedUnits(1, earthSkill, SpellTarget{Type: "unit", Position: Position{Col: 1, Row: 0}})

		if len(affected) != 2 || affected[0] != targetA || affected[1] != targetB {
			t.Fatalf("sommer should make single earth spell affect square battlefield, affected=%+v", cardsToInfo(affected))
		}
	})

	t.Run("1421012 林地飞鼠 changes its load to one air", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		squirrel := placeUnit(baseCard(t, "1421012"), 0, 1, 1, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  squirrel.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use woodland flying squirrel: %v", err)
		}
		if effectiveElementsGain(squirrel)[model.ElementAir] != 1 || totalLoad(squirrel) != 1 {
			t.Fatalf("squirrel should become load 1 air, load=%v", effectiveElementsGain(squirrel))
		}
	})

	t.Run("1421014 风息谷旅商 draws up to three for allied beast plant or spirit companions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		merchant := placeUnit(baseCard(t, "1421014"), 0, 1, 1, engine)
		placeUnit(baseCard(t, "1401002"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1421003"), 0, 0, 1, engine)
		placeUnit(baseCard(t, "1221001"), 0, 2, 0, engine)
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1021002"), 0, 1),
			NewCardInstance(baseCard(t, "1021004"), 0, 1),
			NewCardInstance(baseCard(t, "1021006"), 0, 1),
		}

		engine.triggerEffects(TriggerOnEnter, merchant, nil, nil)

		if len(p0.Hand) != 3 || len(p0.Deck) != 1 {
			t.Fatalf("merchant should draw max three, hand=%d deck=%d", len(p0.Hand), len(p0.Deck))
		}
	})

	t.Run("1611002 黑袍执行官 gains marks from friendly death and spends them to destroy enemy", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		executor := placeUnit(baseCard(t, "1611002"), 0, 0, 0, engine)
		ally := placeUnit(baseCard(t, "1021002"), 0, 1, 0, engine)
		enemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		engine.destroyUnit(ally, 0)
		if executor.Statuses["暗影标记"] != ally.Card.Life {
			t.Fatalf("executor should gain marks equal to dead companion life, statuses=%v", executor.Statuses)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  executor.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use executor ultimate: %v", err)
		}
		if engine.State.Players[1].Units[1][0] != nil || len(engine.State.Players[1].Graveyard) != 1 || engine.State.Players[1].Graveyard[0].InstanceID != enemy.InstanceID {
			t.Fatalf("executor should destroy enemy companion, units=%v grave=%+v", engine.State.Players[1].Units[1][0], cardsToInfo(engine.State.Players[1].Graveyard))
		}
	})

	t.Run("1621003 恐惧魔 requires devouring three friendly life before summoning", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		demon := NewCardInstance(baseCard(t, "1621003"), 0, 1)
		p0.Hand = []*CardInstance{demon}
		sacrifice := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		sacrifice.CurrentLife = 3

		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": demon.InstanceID,
			"col":         float64(1),
			"row":         float64(1),
			"devour_ids":  []any{sacrifice.InstanceID},
		}}); err != nil {
			t.Fatalf("summon fear demon with life devour: %v", err)
		}
		if p0.Units[1][1] != demon || p0.Units[0][0] != nil || len(p0.Graveyard) != 1 || p0.Graveyard[0].InstanceID != sacrifice.InstanceID {
			t.Fatalf("fear demon should enter after devouring sacrifice, unit=%v sacrifice_slot=%v grave=%v", p0.Units[1][1], p0.Units[0][0], cardsToInfo(p0.Graveyard))
		}
	})

	t.Run("1621004 巫术祭司 sacrifices a companion then gives its life to another character", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		priest := placeUnit(baseCard(t, "1621004"), 0, 0, 0, engine)
		sacrifice := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		target := placeUnit(baseCard(t, "1021002"), 0, 2, 0, engine)
		target.CurrentLife = 1
		life := sacrifice.CurrentLife

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  priest.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use witch priest ultimate: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "witch_priest_sacrifice" {
			t.Fatalf("witch priest should ask for sacrifice, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{sacrifice.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve witch priest sacrifice: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "witch_priest_heal" {
			t.Fatalf("witch priest should ask for recipient, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve witch priest recipient: %v", err)
		}
		if engine.State.Players[0].Units[1][0] != nil || target.CurrentLife != 1+life {
			t.Fatalf("witch priest should sacrifice and add life, sacrificed=%v target_life=%d", engine.State.Players[0].Units[1][0], target.CurrentLife)
		}
	})
}

func TestHighRiskItemSemanticsBatch(t *testing.T) {
	t.Run("2011002 统御者之冠 clears load of newly summoned friendly companions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		crown := NewCardInstance(baseCard(t, "2011002"), 0, 1)
		engine.State.Players[0].Equipment[0] = crown
		companion := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)

		engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, 0, companion, map[string]any{"entered_player": 0})

		if totalLoad(companion) != 0 {
			t.Fatalf("overlord crown should clear companion load, load=%v", effectiveElementsGain(companion))
		}
	})

	t.Run("2021002 and 2021017 expose slot expansion markers on equip", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		necklace := NewCardInstance(baseCard(t, "2021002"), 0, 1)
		pack := NewCardInstance(baseCard(t, "2021017"), 0, 1)

		engine.triggerEffects(TriggerOnEnter, necklace, nil, nil)
		engine.triggerEffects(TriggerOnEnter, pack, nil, nil)
		engine.triggerEffects(TriggerOnUseItem, necklace, nil, nil)
		engine.triggerEffects(TriggerOnUseItem, pack, nil, nil)
		engine.triggerEffects(TriggerOnUnitEnter, necklace, nil, nil)
		engine.triggerEffects(TriggerOnUnitEnter, pack, nil, nil)
		engine.triggerEffects(TriggerOnDefend, necklace, nil, nil)
		engine.triggerEffects(TriggerOnDefend, pack, nil, nil)
		engine.triggerEffects(TriggerOnSpellHit, necklace, nil, nil)
		engine.triggerEffects(TriggerOnSpellHit, pack, nil, nil)
		engine.triggerEffects(TriggerOnSpellCast, necklace, nil, nil)
		engine.triggerEffects(TriggerOnSpellCast, pack, nil, nil)
		engine.triggerEffects(TriggerOnTurnStart, necklace, nil, nil)
		engine.triggerEffects(TriggerOnTurnStart, pack, nil, nil)
		engine.triggerEffects(TriggerOnTurnEnd, necklace, nil, nil)
		engine.triggerEffects(TriggerOnTurnEnd, pack, nil, nil)
		engine.triggerEffects(TriggerOnConsume, necklace, nil, nil)
		engine.triggerEffects(TriggerOnConsume, pack, nil, nil)
		engine.triggerEffects(TriggerOnAttack, necklace, nil, nil)
		engine.triggerEffects(TriggerOnAttack, pack, nil, nil)
		engine.triggerEffects(TriggerOnDamaged, necklace, nil, nil)
		engine.triggerEffects(TriggerOnDamaged, pack, nil, nil)
		engine.triggerEffects(TriggerPerTurn, necklace, nil, nil)
		engine.triggerEffects(TriggerPerTurn, pack, nil, nil)
		engine.triggerEffects(TriggerUltimate, necklace, nil, nil)
		engine.triggerEffects(TriggerUltimate, pack, nil, nil)
		engine.triggerEffects(TriggerOnUseItem, necklace, nil, nil)
		engine.triggerEffects(TriggerOnUseItem, pack, nil, nil)
		behaviorNecklace := globalRegistry.GetBehavior("2021002")
		behaviorPack := globalRegistry.GetBehavior("2021017")
		if equip, ok := behaviorNecklace.(OnEquipBehavior); ok {
			if err := equip.OnEquip(&EffectContext{Engine: engine, Source: necklace, PlayerID: 0, OpponentID: 1}); err != nil {
				t.Fatalf("equip necklace: %v", err)
			}
		}
		if equip, ok := behaviorPack.(OnEquipBehavior); ok {
			if err := equip.OnEquip(&EffectContext{Engine: engine, Source: pack, PlayerID: 0, OpponentID: 1}); err != nil {
				t.Fatalf("equip pack: %v", err)
			}
		}

		if necklace.Statuses["技能槽位+1"] != 1 || pack.Statuses["道具槽位+3"] != 1 {
			t.Fatalf("slot expansion markers wrong, necklace=%v pack=%v", necklace.Statuses, pack.Statuses)
		}
	})

	t.Run("2021010 封印卷轴 seals one enemy skill only when enemy has at least four skills", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p1 := engine.State.Players[1]
		for i := 0; i < 4; i++ {
			p1.Skills[i] = readySkill(baseCard(t, "3021005"), 1)
		}
		scroll := NewCardInstance(baseCard(t, "2021010"), 0, 1)

		engine.triggerEffects(TriggerOnUseItem, scroll, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "sealing_scroll" {
			t.Fatalf("sealing scroll should open enemy skill selection, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{p1.Skills[2].InstanceID},
		}}); err != nil {
			t.Fatalf("resolve sealing scroll: %v", err)
		}
		if p1.Skills[2].Statuses[StatusSeal] != 2 {
			t.Fatalf("selected enemy skill should be sealed to next turn end, statuses=%v", p1.Skills[2].Statuses)
		}
	})

	t.Run("2021012 速写卷轴 casts a learned skill without tapping it", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		skill := readySkill(baseCard(t, "3121002"), 0)
		p0.Skills[0] = skill
		p0.Elements[model.ElementFire] = 10
		scroll := NewCardInstance(baseCard(t, "2021012"), 0, 1)

		engine.triggerEffects(TriggerOnUseItem, scroll, nil, nil)

		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "sketch_scroll_skill" {
			t.Fatalf("sketch scroll should ask for a learned skill, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{skill.InstanceID},
		}}); err != nil {
			t.Fatalf("choose sketch skill: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "sketch_scroll_target" {
			t.Fatalf("sketch scroll should ask for target, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("choose sketch target: %v", err)
		}
		if skill.IsHorizontal {
			t.Fatalf("sketch scroll should not tap the learned skill")
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.Skill != skill {
			t.Fatalf("sketch scroll should cast the chosen skill, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("2021015 法力增强剂C makes next skill use cost zero", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		booster := NewCardInstance(baseCard(t, "2021015"), 0, 1)
		skill := readySkill(baseCard(t, "3021005"), 0)

		engine.triggerEffects(TriggerOnUseItem, booster, nil, nil)
		cost := engine.effectiveSkillUseCost(p0, skill)

		if totalCost(cost) != 0 || len(p0.TempModifiers) == 0 {
			t.Fatalf("mana booster C should make next skill free, cost=%v modifiers=%v", cost, p0.TempModifiers)
		}
	})

	t.Run("2021018 奥术符文 gives one friendly skill +3 power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		skill := readySkill(baseCard(t, "3021005"), 0)
		engine.State.Players[0].Skills[0] = skill
		rune := NewCardInstance(baseCard(t, "2021018"), 0, 1)

		engine.triggerEffects(TriggerOnUseItem, rune, nil, nil)

		if skill.PowerBonus != 3 {
			t.Fatalf("arcane rune should add +3 power to a friendly skill, bonus=%d", skill.PowerBonus)
		}
	})

	t.Run("2111001 火龙之心 spends up to three fire for next spell power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		heart := NewCardInstance(baseCard(t, "2111001"), 0, 1)
		p0.Equipment[0] = heart
		p0.Elements[model.ElementFire] = 5

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  heart.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use fire dragon heart: %v", err)
		}
		if p0.Elements[model.ElementFire] != 2 || len(p0.TempModifiers) == 0 || p0.TempModifiers[0].Amount != 9 {
			t.Fatalf("fire dragon heart should spend 3 fire for +9 power modifier, elements=%v modifiers=%v", p0.Elements, p0.TempModifiers)
		}
	})

	t.Run("2111002 努尔之眼 counts only fire damage and converts markers", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		eye := NewCardInstance(baseCard(t, "2111002"), 0, 1)
		p0.Equipment[0] = eye

		engine.triggerEffects(TriggerOnDamaged, eye, nil, map[string]any{"damage_element": model.ElementAir})
		engine.triggerEffects(TriggerOnDamaged, eye, nil, map[string]any{"damage_element": model.ElementFire})
		if eye.Statuses["火焰标记"] != 1 {
			t.Fatalf("nur eye should count only fire damage, statuses=%v", eye.Statuses)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  eye.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use nur eye: %v", err)
		}
		if eye.Statuses["火焰标记"] != 0 || p0.Elements[model.ElementFire] != 2 {
			t.Fatalf("one marker should become 2 fire, statuses=%v elements=%v", eye.Statuses, p0.Elements)
		}
	})

	t.Run("2211002 嗜魔弓 binds winter skill and pays water for bow load on spell cast", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		bow := NewCardInstance(baseCard(t, "2211002"), 0, 1)
		p0.Equipment[0] = bow
		p0.Elements[model.ElementWater] = 1

		engine.triggerEffects(TriggerOnEnter, bow, nil, nil)
		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, readySkill(baseCard(t, "3021005"), 0), map[string]any{"cast_player": 0})

		if len(p0.SkillPool) != 0 {
			t.Fatalf("winter bow bound skill should not enter skill pool, skill_pool=%+v", cardsToInfo(p0.SkillPool))
		}
		if len(bow.BoundSkills) != 1 || bow.BoundSkills[0].Card.Number != "3201002" {
			t.Fatalf("winter bow should bind winter skill, bound=%+v", cardsToInfo(bow.BoundSkills))
		}
		if p0.Elements[model.ElementWater] != 0 || effectiveElementsGain(bow)[model.ElementWater] != bow.Card.ElementsGain[model.ElementWater]+1 {
			t.Fatalf("winter bow should pay 1 water for +1 water load, elements=%v load=%v", p0.Elements, effectiveElementsGain(bow))
		}
	})

	t.Run("2221010 and 2221011 water runes buff water companion or heal allies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		water := placeUnit(baseCard(t, "1221003"), 0, 0, 0, engine)
		ally := placeUnit(baseCard(t, "1021002"), 0, 1, 0, engine)
		water.CurrentLife = 1
		ally.CurrentLife = 1
		tide := NewCardInstance(baseCard(t, "2221010"), 0, 1)
		rain := NewCardInstance(baseCard(t, "2221011"), 0, 1)

		engine.triggerEffects(TriggerOnUseItem, tide, nil, nil)
		engine.triggerEffects(TriggerOnUseItem, rain, nil, nil)

		if effectiveElementsGain(water)[model.ElementWater] != water.Card.ElementsGain[model.ElementWater]+2 {
			t.Fatalf("tide rune should add +2 water load, load=%v", effectiveElementsGain(water))
		}
		if water.CurrentLife <= 1 || ally.CurrentLife <= 1 {
			t.Fatalf("rain of grace should heal all friendly units, water=%d ally=%d", water.CurrentLife, ally.CurrentLife)
		}
	})

	t.Run("2221012 水行之靴 gains water load near three water companions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		hero := placeUnit(baseCard(t, "4211001"), 0, 1, 1, engine)
		p0.Hero = hero
		boots := NewCardInstance(baseCard(t, "2221012"), 0, 1)
		p0.Equipment[0] = boots
		placeUnit(baseCard(t, "1221001"), 0, 0, 1, engine)
		placeUnit(baseCard(t, "1221003"), 0, 1, 0, engine)
		placeUnit(baseCard(t, "1221006"), 0, 2, 1, engine)

		engine.triggerEffects(TriggerOnTurnStart, boots, nil, nil)

		if effectiveElementsGain(boots)[model.ElementWater] != boots.Card.ElementsGain[model.ElementWater]+1 {
			t.Fatalf("water walking boots should gain +1 water load, load=%v", effectiveElementsGain(boots))
		}
	})

	t.Run("2311001 雷之源 reduces air skill use cost", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		placeUnit(baseCard(t, "2311001"), 0, 0, 0, engine)
		skill := readySkill(baseCard(t, "3321005"), 0)
		cost := engine.effectiveSkillUseCost(engine.State.Players[0], skill)

		if cost[model.ElementAir] != max(skillUseCost(skill.Card)[model.ElementAir]-1, 0) {
			t.Fatalf("thunder source should reduce air skill cost, cost=%v base=%v", cost, skillUseCost(skill.Card))
		}
	})

	t.Run("2311002 spends counters for power and 2321001 asks on draw", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		drum := NewCardInstance(baseCard(t, "2311002"), 0, 1)
		compass := NewCardInstance(baseCard(t, "2321001"), 0, 1)
		p0.Equipment[0] = drum
		p0.Equipment[1] = compass
		drum.Statuses["雷鼓标记"] = 3

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  drum.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use thunder drum: %v", err)
		}
		if drum.Statuses["雷鼓标记"] != 0 || len(p0.TempModifiers) == 0 || p0.TempModifiers[0].Amount != 3 {
			t.Fatalf("thunder drum should spend 3 marks for +3 power, statuses=%v modifiers=%v", drum.Statuses, p0.TempModifiers)
		}
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
		}
		engine.drawCards(0, 2)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "windbreath_compass" {
			t.Fatalf("windbreath compass should ask on draw, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{},
		}}); err != nil {
			t.Fatalf("decline first windbreath compass trigger: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "windbreath_compass" {
			t.Fatalf("windbreath compass should queue one prompt per drawn card, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{compass.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve windbreath compass: %v", err)
		}
		if effectiveElementsGain(compass)[model.ElementAir] != compass.Card.ElementsGain[model.ElementAir]+1 {
			t.Fatalf("windbreath compass should gain temporary air load, load=%v", effectiveElementsGain(compass))
		}
		if engine.State.PendingAction != nil {
			t.Fatalf("all queued windbreath compass prompts should resolve, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
			t.Fatalf("end turn after windbreath compass: %v", err)
		}
		if effectiveElementsGain(compass)[model.ElementAir] != compass.Card.ElementsGain[model.ElementAir] {
			t.Fatalf("windbreath compass temporary air load should expire at turn end, load=%v", effectiveElementsGain(compass))
		}
	})

	t.Run("2321012 随风斗篷 moves hero to empty position", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		hero := placeUnit(baseCard(t, "4311003"), 0, 1, 1, engine)
		p0.Hero = hero
		cloak := NewCardInstance(baseCard(t, "2321012"), 0, 1)
		p0.Equipment[0] = cloak
		old := *hero.Position

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  cloak.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use wind cloak ultimate: %v", err)
		}
		if hero.Position == nil || (hero.Position.Col == old.Col && hero.Position.Row == old.Row) || p0.Units[old.Col][old.Row] != nil {
			t.Fatalf("wind cloak should move hero, old=%+v new=%+v", old, hero.Position)
		}
	})

	t.Run("2411001, 2421001, 2421004, 2421013 earth items apply current support effects", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		unit := placeUnit(baseCard(t, "1021002"), 0, 0, 0, engine)
		mastered := placeUnit(baseCard(t, "1421003"), 0, 2, 0, engine)
		unit.CurrentLife = 1
		addElementsGainBonus(unit, model.ElementEarth, 3)
		treeHeart := NewCardInstance(baseCard(t, "2411001"), 0, 1)
		care := NewCardInstance(baseCard(t, "2421001"), 0, 1)
		care.IsHorizontal = false
		exam := NewCardInstance(baseCard(t, "2421004"), 0, 1)
		primer := NewCardInstance(baseCard(t, "2421013"), 0, 1)
		p0.Equipment[0] = treeHeart
		p0.Equipment[1] = care
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1), NewCardInstance(baseCard(t, "1021002"), 0, 1)}
		highCost := NewCardInstance(baseCard(t, "1021013"), 0, 1)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  treeHeart.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use ancient tree heart: %v", err)
		}
		engine.advanceMastery(mastered, 0, 1)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "knowledge_tree_care" {
			t.Fatalf("knowledge tree care should ask to trigger on mastery, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{care.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve knowledge tree care: %v", err)
		}
		engine.triggerEffects(TriggerOnUseItem, exam, nil, nil)
		cost := copyElementCost(highCost.Card.ElementsCost)
		if modifier, ok := globalRegistry.GetBehavior("2421013").(CardPlayCostModifier); ok {
			modifier.ModifyCardPlayCost(&EffectContext{Engine: engine, Source: primer, PlayerID: 0, OpponentID: 1}, highCost, cost)
		}

		if unit.CurrentLife <= 1 || effectiveElementsGain(unit)[model.ElementEarth] < unit.Card.ElementsGain[model.ElementEarth]+4 {
			t.Fatalf("earth support items should heal and add earth load, life=%d load=%v", unit.CurrentLife, effectiveElementsGain(unit))
		}
		if len(p0.Hand) != 1 || p0.Elements[model.ElementEarth] != 1 {
			t.Fatalf("knowledge tree care should draw and gain earth, hand=%d elements=%v", len(p0.Hand), p0.Elements)
		}
		if cost[model.ElementEarth] != max(highCost.Card.ElementsCost[model.ElementEarth]-2, 0) {
			t.Fatalf("geography primer should reduce high-cost card earth cost by 2, cost=%v base=%v", cost, highCost.Card.ElementsCost)
		}
	})

	t.Run("2511001 万灵药 supports draw and gain choices", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		panacea := NewCardInstance(baseCard(t, "2511001"), 0, 1)
		for i := 0; i < 4; i++ {
			p0.Deck = append(p0.Deck, NewCardInstance(baseCard(t, "1021001"), 0, 1))
		}

		engine.triggerEffects(TriggerOnUseItem, panacea, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "panacea_mode" {
			t.Fatalf("panacea should ask for mode, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{"draw"},
		}}); err != nil {
			t.Fatalf("resolve panacea draw: %v", err)
		}
		if len(p0.Hand) != 4 {
			t.Fatalf("panacea draw mode should draw four, hand=%d", len(p0.Hand))
		}

		engine.triggerEffects(TriggerOnUseItem, panacea, nil, nil)
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{"gain"},
		}}); err != nil {
			t.Fatalf("resolve panacea gain: %v", err)
		}
		if p0.Elements[model.ElementArcane] != 5 {
			t.Fatalf("panacea gain mode should grant 5 arcane, elements=%v", p0.Elements)
		}
	})

	t.Run("2511002 and 2601002 modify defense or weaken enemy skills", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		shield := NewCardInstance(baseCard(t, "2511002"), 0, 1)
		book := NewCardInstance(baseCard(t, "2601002"), 0, 1)
		enemyA := readySkill(baseCard(t, "3021005"), 1)
		enemyB := readySkill(baseCard(t, "3021008"), 1)
		engine.State.Players[1].Skills[0] = enemyA
		engine.State.Players[1].Skills[1] = enemyB
		stats := &SpellStats{}

		if modifier, ok := globalRegistry.GetBehavior("2511002").(SpellStatModifier); ok {
			modifier.ModifySpellStats(&EffectContext{Engine: engine, Source: shield, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"purpose": string(skillPurposeDefend)}}, stats)
		}
		engine.triggerEffects(TriggerOnEnter, book, nil, nil)

		if stats.PowerBonus != 2 {
			t.Fatalf("shining shield should add +2 defense power, stats=%+v", stats)
		}
		if enemyA.Statuses[StatusWeaken] != 1 || enemyB.Statuses[StatusWeaken] != 1 {
			t.Fatalf("spellbook should weaken all enemy skills, a=%v b=%v", enemyA.Statuses, enemyB.Statuses)
		}
	})
}

func TestHighRiskItemSemanticsBatchTwo(t *testing.T) {
	t.Run("2011001 大法师之杖 records stored skill marker on enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		staff := NewCardInstance(baseCard(t, "2011001"), 0, 1)
		p0.SkillPool = []*CardInstance{readySkill(baseCard(t, "3021005"), 0)}

		engine.triggerEffects(TriggerOnEnter, staff, nil, nil)

		if staff.Statuses["存储技能"] != 1 {
			t.Fatalf("archmage staff should record stored skill marker when skill pool has cards, statuses=%v", staff.Statuses)
		}
	})

	t.Run("2011003 君王法袍 lowers enemy spell damage only after ultimate is used", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		engine.State.CurrentTurn = 1
		robe := NewCardInstance(baseCard(t, "2011003"), 1, 1)
		engine.State.Players[1].Equipment[0] = robe
		loaded := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
		addElementsGainBonus(loaded, model.ElementArcane, 4)
		skill := readySkill(baseCard(t, "3021005"), 0)

		damage := engine.effectiveSpellDamage(0, skill, 3, nil)
		if damage != 3 {
			t.Fatalf("king robe should not passively reduce enemy spell damage, damage=%d", damage)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  robe.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use king robe ultimate: %v", err)
		}
		damage = engine.effectiveSpellDamage(0, skill, 3, nil)

		expected := max(3-(totalFieldLoad(engine.State.Players[1])-totalFieldLoad(engine.State.Players[0]))/2, 0)
		if damage != expected {
			t.Fatalf("king robe ultimate should reduce enemy spell damage by load difference, damage=%d expected=%d", damage, expected)
		}
	})

	t.Run("2021022, 2321010, 2521002, 2521004 reactive items create current temporary readiness markers", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		counter := NewCardInstance(baseCard(t, "2021022"), 0, 1)
		illusion := NewCardInstance(baseCard(t, "2321010"), 0, 1)
		shelter := NewCardInstance(baseCard(t, "2521002"), 0, 1)
		sanction := NewCardInstance(baseCard(t, "2521004"), 0, 1)

		engine.triggerEffects(TriggerOnUseItem, counter, nil, nil)
		engine.triggerEffects(TriggerOnUseItem, illusion, nil, nil)
		engine.triggerEffects(TriggerOnUseItem, shelter, nil, nil)
		engine.triggerEffects(TriggerOnUseItem, sanction, nil, nil)

		types := map[string]bool{}
		for _, modifier := range p0.TempModifiers {
			types[modifier.Type] = true
		}
		if !types["shelter_rune"] || !types["holy_sanction"] {
			t.Fatalf("shelter and sanction should create temporary modifiers, modifiers=%v", p0.TempModifiers)
		}
	})

	t.Run("2321011 传送符文 resets a friendly unit with current adapter effect", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		target := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		target.IsHorizontal = true
		target.Statuses[StatusCooldown] = 1
		rune := NewCardInstance(baseCard(t, "2321011"), 0, 1)

		engine.triggerEffects(TriggerOnUseItem, rune, nil, nil)

		if target.IsHorizontal || target.Statuses[StatusCooldown] != 0 {
			t.Fatalf("teleport rune should reset selected friendly unit in current adapter, horizontal=%v statuses=%v", target.IsHorizontal, target.Statuses)
		}
	})

	t.Run("2411002 裂地巨剑 consumes into next spell power modifier", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		sword := NewCardInstance(baseCard(t, "2411002"), 0, 1)

		engine.triggerEffects(TriggerOnConsume, sword, nil, nil)

		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModSkillPowerBonus || p0.TempModifiers[0].Amount != 4 {
			t.Fatalf("earthsplitter sword should add +4 next spell power in current adapter, modifiers=%v", p0.TempModifiers)
		}
	})

	t.Run("2501001 桎梏 use adapter draws one replacement card", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		shackle := NewCardInstance(baseCard(t, "2501001"), 0, 1)
		draw := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Deck = []*CardInstance{draw}

		engine.triggerEffects(TriggerOnUseItem, shackle, nil, nil)

		if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != draw.InstanceID {
			t.Fatalf("shackle should draw one replacement card in current adapter, hand=%+v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("2611002 与恶魔的契约书 sacrifices friendly unit and starts enemy-destroy selection", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		contract := NewCardInstance(baseCard(t, "2611002"), 0, 1)
		sacrifice := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		enemy := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)

		engine.triggerEffects(TriggerOnUseItem, contract, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "demon_contract_sacrifice" {
			t.Fatalf("demon contract should ask for sacrifice, pending=%+v enemy=%v", engine.State.PendingAction, enemy.InstanceID)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{sacrifice.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve demon contract sacrifice: %v", err)
		}
		if engine.State.Players[0].Units[1][0] != nil || engine.State.PendingAction == nil || engine.State.PendingAction.Type != "demon_contract_destroy" {
			t.Fatalf("demon contract should sacrifice then ask for enemy destroy, unit=%v pending=%+v", engine.State.Players[0].Units[1][0], engine.State.PendingAction)
		}
	})

	t.Run("2621002 巫毒娃娃 starts with three shadow marks and retaliates when linked side is damaged", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		doll := NewCardInstance(baseCard(t, "2621002"), 0, 1)
		engine.State.Players[0].Equipment[0] = doll
		friendly := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		enemy := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)

		engine.triggerEffects(TriggerOnEquip, doll, nil, nil)
		engine.triggerEffects(TriggerOnDamaged, doll, friendly, map[string]any{"damage": 1})

		if doll.Statuses["暗影标记"] != 2 || enemy.CurrentLife != enemy.Card.Life-1 {
			t.Fatalf("voodoo doll should spend one mark to damage enemy, statuses=%v enemy_life=%d", doll.Statuses, enemy.CurrentLife)
		}
	})

	t.Run("2621004, 2621010, 2621011, 2621013 shadow item adapters apply current defensive effects", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		hero := placeUnit(baseCard(t, "4611001"), 0, 1, 1, engine)
		p0.Hero = hero
		enemy := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
		enemy.CurrentLife = 5
		veil := NewCardInstance(baseCard(t, "2621004"), 0, 1)
		abyss := NewCardInstance(baseCard(t, "2621010"), 0, 1)
		frenzy := NewCardInstance(baseCard(t, "2621011"), 0, 1)
		ring := NewCardInstance(baseCard(t, "2621013"), 0, 1)
		p0.Equipment[0] = ring
		enemySkill := readySkill(baseCard(t, "3021005"), 1)
		enemySkill.Statuses[StatusWeaken] = 1
		engine.State.Players[1].Skills[0] = enemySkill

		engine.triggerEffects(TriggerOnSpellHit, veil, nil, map[string]any{"cast_player": 1})
		engine.triggerEffects(TriggerOnUseItem, abyss, nil, nil)
		engine.triggerEffects(TriggerOnUseItem, frenzy, nil, nil)
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  ring.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use witchcraft ring: %v", err)
		}

		if hero.Statuses["引魔"] != 1 || enemy.CurrentLife != 3 || enemy.Statuses[StatusStun] != 1 || enemySkill.Statuses[StatusWeaken] != 2 {
			t.Fatalf("shadow item adapters wrong, hero=%v enemy_life=%d enemy_status=%v skill=%v", hero.Statuses, enemy.CurrentLife, enemy.Statuses, enemySkill.Statuses)
		}
	})
}

func baseCard(t *testing.T, id string) *model.Card {
	t.Helper()
	card, ok := cards.CardDB[id]
	if !ok {
		t.Fatalf("missing card %s", id)
	}
	return card
}

func placeUnit(card *model.Card, ownerID int, col int, row int, engine *Engine) *CardInstance {
	unit := NewCardInstance(card, ownerID, engine.State.TurnNumber)
	unit.IsHorizontal = false
	unit.Position = &Position{Col: col, Row: row}
	engine.State.Players[ownerID].Units[col][row] = unit
	return unit
}

func countCardNumber(cards []*CardInstance, number string) int {
	count := 0
	for _, card := range cards {
		if card != nil && card.Card != nil && card.Card.Number == number {
			count++
		}
	}
	return count
}

func TestCounterTrapCanBeSetHiddenAndTriggeredWithOverexert(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	p0.Elements = map[string]int{}
	p1.Elements = map[string]int{}
	for _, elem := range model.AllElements {
		p0.Elements[elem] = 0
		p1.Elements[elem] = 0
	}

	counter := NewCardInstance(baseCard(t, "2121002"), 0, 1)
	p0.Hand = append(p0.Hand, counter)
	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": counter.InstanceID}}); err != nil {
		t.Fatalf("set counter trap without elements: %v", err)
	}
	if len(p0.Hand) != 0 || p0.Equipment[0] != counter || !counter.IsSetCounter {
		t.Fatalf("counter should be set in equipment, hand=%d equipment=%v set=%v", len(p0.Hand), p0.Equipment[0], counter.IsSetCounter)
	}
	opponentView := engine.playerStateToInfo(p0, false)
	hidden := opponentView["equipment"].([5]any)[0].(map[string]any)
	if hidden["is_hidden"] != true || hidden["number"] != nil {
		t.Fatalf("opponent should see only hidden counter info: %+v", hidden)
	}

	payer := placeUnit(baseCard(t, "1121001"), 0, 0, 1, engine)
	target := placeUnit(baseCard(t, "1221001"), 1, 0, 1, engine)
	engine.State.CurrentTurn = 1
	if err := engine.HandleAction(1, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": target.InstanceID}}); err != nil {
		t.Fatalf("consume to trigger counter: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_trigger" {
		t.Fatalf("expected counter trigger prompt, pending=%+v", engine.State.PendingAction)
	}
	state := engine.GetStateForPlayer(0)
	pending := state["pending_action"].(map[string]any)
	if pending["can_overexert"] != true || pending["cost"] == nil {
		t.Fatalf("counter trigger state should expose overexert payment data: %+v", pending)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected":      []any{counter.InstanceID},
		"overexert_ids": []any{payer.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve counter with overexert: %v", err)
	}
	if target.Statuses[StatusBurn] != 1 {
		t.Fatalf("fire rune counter should burn consumed unit, statuses=%v", target.Statuses)
	}
	if !payer.IsHorizontal {
		t.Fatalf("overexert payer should become horizontal")
	}
	if p0.Equipment[0] != nil || len(p0.Graveyard) != 1 || p0.Graveyard[0] != counter {
		t.Fatalf("triggered counter should go to graveyard, equipment=%v grave=%d", p0.Equipment[0], len(p0.Graveyard))
	}
}

func TestSetCounterTrapDoesNotAutoTriggerBehindExistingPrompt(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	for _, elem := range model.AllElements {
		p0.Elements[elem] = 10
	}
	first := NewCardInstance(baseCard(t, "2121002"), 0, 1)
	second := NewCardInstance(baseCard(t, "2121002"), 0, 1)
	first.IsSetCounter = true
	second.IsSetCounter = true
	p0.Equipment[0] = first
	p0.Equipment[1] = second

	target := placeUnit(baseCard(t, "1221001"), 1, 0, 1, engine)
	engine.State.CurrentTurn = 1
	if err := engine.HandleAction(1, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": target.InstanceID}}); err != nil {
		t.Fatalf("consume to trigger counters: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_trigger" {
		t.Fatalf("expected first counter prompt, pending=%+v", engine.State.PendingAction)
	}
	if target.Statuses[StatusBurn] != 0 {
		t.Fatalf("second set counter must not auto-trigger while first prompt is pending, statuses=%v", target.Statuses)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{},
	}}); err != nil {
		t.Fatalf("decline first counter: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_trigger" {
		t.Fatalf("declining first counter should ask the next eligible counter, pending=%+v", engine.State.PendingAction)
	}
	if len(engine.State.PendingAction.Candidates) != 1 || engine.State.PendingAction.Candidates[0]["instance_id"] != second.InstanceID {
		t.Fatalf("second counter should be queued after first decline, candidates=%+v", engine.State.PendingAction.Candidates)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{second.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve second counter: %v", err)
	}
	if target.Statuses[StatusBurn] != 1 {
		t.Fatalf("second counter should resolve after queue prompt, statuses=%v", target.Statuses)
	}
	if p0.Equipment[0] != first || !first.IsSetCounter || p0.Equipment[1] != nil || len(p0.Graveyard) != 1 || p0.Graveyard[0] != second {
		t.Fatalf("only the selected second counter should be spent, equipment=%v graveyard=%v", p0.Equipment, cardsToInfo(p0.Graveyard))
	}
}

func TestSpellCastCounterPromptResumesDefenseWindow(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	for _, elem := range model.AllElements {
		p0.Elements[elem] = 10
		p1.Elements[elem] = 10
	}
	p0.Skills[0] = readySkill(baseCard(t, "3121002"), 0)
	placeUnit(baseCard(t, "1221001"), 1, 0, 0, engine)
	counter := NewCardInstance(baseCard(t, "2021018"), 1, 1)
	counter.IsSetCounter = true
	p1.Equipment[0] = counter

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(0),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast spell into counter prompt: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.PlayerID != 1 || engine.State.PendingSpell == nil {
		t.Fatalf("expected counter prompt with pending spell, action=%+v spell=%+v", engine.State.PendingAction, engine.State.PendingSpell)
	}
	if engine.State.ResumePhase != PhaseDefenseWindow {
		t.Fatalf("counter prompt should resume defense window, got %s", engine.State.ResumePhase)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{}}}); err != nil {
		t.Fatalf("decline counter: %v", err)
	}
	if engine.State.Phase != PhaseDefenseWindow || engine.State.PendingSpell == nil {
		t.Fatalf("declining spell-cast counter should continue defense window, phase=%s spell=%+v", engine.State.Phase, engine.State.PendingSpell)
	}
}

func TestCounterRuneCancelsConsumableBeforeEffect(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	for _, elem := range model.AllElements {
		p0.Elements[elem] = 10
		p1.Elements[elem] = 10
	}
	scroll := NewCardInstance(baseCard(t, "2121003"), 0, 1)
	p0.Hand = append(p0.Hand, scroll)
	target := placeUnit(baseCard(t, "1221001"), 1, 0, 1, engine)
	counter := NewCardInstance(baseCard(t, "2021022"), 1, 1)
	counter.IsSetCounter = true
	p1.Equipment[0] = counter

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": scroll.InstanceID}}); err != nil {
		t.Fatalf("use scroll into counter rune: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_trigger" {
		t.Fatalf("counter rune should prompt before scroll effect, pending=%+v", engine.State.PendingAction)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{counter.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve counter rune: %v", err)
	}
	if engine.State.PendingAction != nil {
		t.Fatalf("cancelled scroll should not open its target prompt, pending=%+v", engine.State.PendingAction)
	}
	if target.Statuses[StatusBurn] != 0 {
		t.Fatalf("cancelled scroll should not burn target, statuses=%v", target.Statuses)
	}
	if p1.Equipment[0] != nil || len(p1.Graveyard) != 1 || p1.Graveyard[0] != counter {
		t.Fatalf("counter rune should go to graveyard, equipment=%v grave=%d", p1.Equipment[0], len(p1.Graveyard))
	}
}

func TestShelterRuneCancelsLowPowerSpellHitBeforeDamage(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	for _, elem := range model.AllElements {
		p0.Elements[elem] = 10
		p1.Elements[elem] = 10
	}
	skill := readySkill(baseCard(t, "3121002"), 0)
	p0.Skills[0] = skill
	target := placeUnit(baseCard(t, "1221001"), 1, 0, 1, engine)
	counter := NewCardInstance(baseCard(t, "2521002"), 1, 1)
	counter.IsSetCounter = true
	p1.Equipment[0] = counter
	beforeLife := target.CurrentLife

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": skill.InstanceID,
		"target_type": "unit",
		"target_col":  float64(0),
		"target_row":  float64(1),
	}}); err != nil {
		t.Fatalf("cast low-power spell: %v", err)
	}
	if engine.State.PendingSpell == nil || engine.State.Phase != PhaseDefenseWindow {
		t.Fatalf("spell should enter defense window before hit counter, phase=%s pending=%+v", engine.State.Phase, engine.State.PendingSpell)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("skip defense into shelter rune: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_trigger" || engine.State.PendingAction.PlayerID != 1 {
		t.Fatalf("shelter rune should prompt before hit damage, pending=%+v", engine.State.PendingAction)
	}
	if target.CurrentLife != beforeLife {
		t.Fatalf("damage must not be applied before shelter rune resolves, before=%d life=%d", beforeLife, target.CurrentLife)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{counter.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve shelter rune: %v", err)
	}
	if target.CurrentLife != beforeLife {
		t.Fatalf("shelter rune should cancel the spell hit before damage, before=%d life=%d", beforeLife, target.CurrentLife)
	}
	if engine.State.PendingSpell != nil || engine.State.Phase != PhaseMain {
		t.Fatalf("cancelled spell hit should close the defense window, phase=%s pending=%+v", engine.State.Phase, engine.State.PendingSpell)
	}
	if p1.Equipment[0] != nil || len(p1.Graveyard) != 1 || p1.Graveyard[0] != counter {
		t.Fatalf("shelter rune should be discarded after reveal, equipment=%v graveyard=%v", p1.Equipment[0], cardsToInfo(p1.Graveyard))
	}
}

func TestBoundSkillAttachesToHostInsteadOfSkillPool(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	ailaya := placeUnit(baseCard(t, "1311002"), 0, 1, 1, engine)

	engine.triggerEffects(TriggerOnEnter, ailaya, nil, nil)

	if len(p0.SkillPool) != 0 {
		t.Fatalf("bound skill should not be added to skill pool, got %d", len(p0.SkillPool))
	}
	if len(ailaya.BoundSkills) != 1 || ailaya.BoundSkills[0].Card.Number != "3301001" {
		t.Fatalf("expected Storm Fury bound to Ailaya, got %+v", ailaya.BoundSkills)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "learn_skill", Data: map[string]any{
		"instance_id": ailaya.BoundSkills[0].InstanceID,
	}}); err == nil {
		t.Fatalf("bound skill should not be learnable through the skill pool")
	}

	info := cardToInfo(ailaya)
	bound, ok := info["bound_skills"].([]map[string]any)
	if !ok || len(bound) != 1 || bound[0]["number"] != "3301001" {
		t.Fatalf("card info should expose bound skill for display, got %+v", info["bound_skills"])
	}

	engine.destroyUnit(ailaya, 0)
	if len(p0.Graveyard) != 1 || len(p0.Graveyard[0].BoundSkills) != 0 {
		t.Fatalf("bound skills should disappear when host leaves battlefield, graveyard=%+v", p0.Graveyard)
	}
}

func TestCycloneWavePaysTwoAirUsingArcaneWildcard(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Elements[model.ElementAir] = 1
	p0.Elements[model.ElementArcane] = 1
	placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast cyclone wave: %v", err)
	}
	if p0.Elements[model.ElementAir] != 0 || p0.Elements[model.ElementArcane] != 0 {
		t.Fatalf("expected air and arcane to be spent, got %v", p0.Elements)
	}
	if engine.State.Phase != PhaseDefenseWindow || p1.Units[1][0] == nil {
		t.Fatalf("expected normal spell defense window")
	}
}

func TestArcaneArrowDealsOneDamageAsSorcery(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	target.CurrentLife = 3
	p0.Skills[0] = readySkill(baseCard(t, "3021005"), 0)
	p0.Elements[model.ElementArcane] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast arcane arrow: %v", err)
	}
	if target.CurrentLife != 2 {
		t.Fatalf("arcane arrow should deal 1 damage, target life=%d", target.CurrentLife)
	}
	if engine.State.Phase != PhaseMain {
		t.Fatalf("sorcery should resolve immediately, got phase %v", engine.State.Phase)
	}
}

func TestSleepAppliesVisibleStunStatusOnHit(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3021009"), 0)
	p0.Elements[model.ElementArcane] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast sleep: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve sleep hit: %v", err)
	}
	if target.Statuses[StatusStun] != 1 {
		t.Fatalf("sleep should apply stun 1, statuses=%v", target.Statuses)
	}
}

func TestMeditationCanCastWithoutTargetAndGainArcane(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p0.Skills[0] = readySkill(baseCard(t, "3021003"), 0)
	p0.Elements[model.ElementAir] = 1
	p0.Elements[model.ElementFire] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "none",
	}}); err != nil {
		t.Fatalf("cast meditation without target: %v", err)
	}
	if p0.Elements[model.ElementArcane] != 1 {
		t.Fatalf("meditation should gain 1 arcane, elements=%v", p0.Elements)
	}
	if p0.Elements[model.ElementAir]+p0.Elements[model.ElementFire] != 0 {
		t.Fatalf("meditation should spend two non-arcane elements for arcane cost, elements=%v", p0.Elements)
	}
}

func TestSquareSpellAppliesDamageAndStatusToAllEnemyUnits(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	front := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	back := placeUnit(baseCard(t, "1021001"), 1, 2, 1, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3121005"), 0)
	p0.Elements[model.ElementFire] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast square firestorm: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve square firestorm: %v", err)
	}
	if front.CurrentLife != front.Card.Life-1 || back.CurrentLife != back.Card.Life-1 {
		t.Fatalf("square spell should damage all enemy units, front=%d back=%d", front.CurrentLife, back.CurrentLife)
	}
	if front.Statuses[StatusBurn] != 1 || back.Statuses[StatusBurn] != 1 {
		t.Fatalf("square spell should burn all enemy units, front=%v back=%v", front.Statuses, back.Statuses)
	}
}

func TestGenericFriendlySpellPowerAndAttackBonuses(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	placeUnit(baseCard(t, "1121004"), 0, 0, 0, engine) // 你的所有法术+1\威
	placeUnit(baseCard(t, "1321006"), 0, 2, 0, engine) // 你的大气法术+1\攻
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Elements[model.ElementAir] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast boosted cyclone wave: %v", err)
	}
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 2 {
		t.Fatalf("expected cyclone power 2 after generic +威, pending=%+v", engine.State.PendingSpell)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve boosted cyclone wave: %v", err)
	}
	if target.CurrentLife != target.Card.Life-2 {
		t.Fatalf("expected cyclone to deal 2 after generic +攻, life=%d", target.CurrentLife)
	}
}

func TestRavenMessengerCanTapToDrawWithoutGainingElements(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	raven := placeUnit(baseCard(t, "1321001"), 0, 0, 0, engine)
	p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 0)}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  raven.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use raven ability: %v", err)
	}
	if !raven.IsHorizontal {
		t.Fatalf("raven should tap")
	}
	if len(p0.Hand) != 1 {
		t.Fatalf("raven should draw 1 card, hand=%d", len(p0.Hand))
	}
	if p0.Elements[model.ElementAir] != 0 {
		t.Fatalf("raven draw option should not gain air, elements=%v", p0.Elements)
	}
}

func TestStormChimeraReducesAirSpellUseCost(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	placeUnit(baseCard(t, "1321010"), 0, 0, 0, engine)
	placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Elements[model.ElementAir] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("chimera should reduce cyclone wave to 1 air: %v", err)
	}
	if p0.Elements[model.ElementAir] != 0 {
		t.Fatalf("expected reduced cost to spend 1 air, elements=%v", p0.Elements)
	}
}

func TestStormChimeraDevoursFriendlyCompanion(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	chimera := NewCardInstance(baseCard(t, "1321010"), 0, engine.State.TurnNumber)
	p0.Hand = append(p0.Hand, chimera)
	food := placeUnit(baseCard(t, "1321014"), 0, 1, 0, engine)
	p0.Elements[model.ElementAir] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": chimera.InstanceID,
		"devour_id":   food.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon chimera with devour: %v", err)
	}
	if p0.Units[1][0] != nil || len(p0.Graveyard) != 1 {
		t.Fatalf("devoured unit should be destroyed")
	}
	if p0.Units[0][0] == nil || p0.Units[0][0].Card.Number != "1321010" {
		t.Fatalf("chimera should enter after devour")
	}
	if p0.Elements[model.ElementAir] != 0 {
		t.Fatalf("summon should spend 3 air, elements=%v", p0.Elements)
	}
}

func TestStormChimeraDevoursMultipleFriendlyCompanions(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	chimera := NewCardInstance(baseCard(t, "1321010"), 0, engine.State.TurnNumber)
	p0.Hand = append(p0.Hand, chimera)
	foodA := placeUnit(baseCard(t, "1321002"), 0, 0, 0, engine)
	foodB := placeUnit(baseCard(t, "1321002"), 0, 2, 0, engine)
	setElementsGain(foodA, map[string]int{model.ElementAir: 1})
	setElementsGain(foodB, map[string]int{model.ElementAir: 2})
	p0.Elements[model.ElementAir] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": chimera.InstanceID,
		"devour_ids":  []any{foodA.InstanceID, foodB.InstanceID},
		"col":         float64(1),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon chimera with multiple devours: %v", err)
	}
	if p0.Units[0][0] != nil || p0.Units[2][0] != nil || len(p0.Graveyard) != 2 {
		t.Fatalf("both devoured units should be destroyed, grave=%d", len(p0.Graveyard))
	}
	if p0.Units[1][0] == nil || p0.Units[1][0].Card.Number != "1321010" {
		t.Fatalf("chimera should enter after multiple devours")
	}
}

func TestStormChimeraRequiresDevourBeforeSummon(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	chimera := NewCardInstance(baseCard(t, "1321010"), 0, engine.State.TurnNumber)
	p0.Hand = append(p0.Hand, chimera)
	placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
	p0.Elements[model.ElementAir] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": chimera.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err == nil {
		t.Fatalf("chimera should require a valid devour target before summon")
	}
}

func TestStaticPulseRequiresEnemyTarget(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p0.Skills[0] = readySkill(baseCard(t, "3321003"), 0)
	p0.Elements[model.ElementAir] = 1
	placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "none",
	}}); err == nil {
		t.Fatalf("static pulse has power and stun text, so it should require an enemy target")
	}
}

func TestHasteSkillEntersReadyWhenLearned(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	sleep := NewCardInstance(baseCard(t, "3021009"), 0, engine.State.TurnNumber)
	p0.SkillPool = append(p0.SkillPool, sleep)
	p0.Elements[model.ElementArcane] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "learn_skill", Data: map[string]any{
		"instance_id": sleep.InstanceID,
	}}); err != nil {
		t.Fatalf("learn sleep: %v", err)
	}
	if p0.Skills[0] == nil || p0.Skills[0].IsHorizontal {
		t.Fatalf("haste skill should enter ready")
	}
}

func TestCardsResetAtEndOfOwnersTurn(t *testing.T) {
	engine := setupReportedBugEngine(t)
	unit := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	unit.IsHorizontal = true

	if err := engine.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
		t.Fatalf("end turn: %v", err)
	}
	if unit.IsHorizontal {
		t.Fatalf("owner's cards should reset at end of turn")
	}
	if engine.State.CurrentTurn != 1 {
		t.Fatalf("turn should pass to player 2")
	}
}

func TestMulingUltimateReturnsOneCompanionFromEachSide(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	muling := placeUnit(baseCard(t, "4311003"), 0, 1, 1, engine)
	own := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	enemy := placeUnit(baseCard(t, "1321010"), 1, 1, 0, engine)
	p0.Elements[model.ElementAir] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  muling.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use muling ultimate: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected muling selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{own.InstanceID, enemy.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve muling: %v", err)
	}
	if p0.Units[0][0] != nil || engine.State.Players[1].Units[1][0] != nil {
		t.Fatalf("selected companions should leave the battlefield")
	}
	if len(p0.Hand) != 1 || len(engine.State.Players[1].Hand) != 1 {
		t.Fatalf("selected companions should return to owners' hands")
	}
	if p0.Elements[model.ElementAir] != 0 {
		t.Fatalf("muling should spend cost difference 2 air, elements=%v", p0.Elements)
	}
}

func TestBlackMarketVendorDiscardsItemAndDrawsTwo(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	vendor := placeUnit(baseCard(t, "1021012"), 0, 0, 0, engine)
	item := NewCardInstance(baseCard(t, "2021003"), 0, 1)
	p0.Hand = []*CardInstance{item}
	p0.Deck = []*CardInstance{
		NewCardInstance(baseCard(t, "1021001"), 0, 1),
		NewCardInstance(baseCard(t, "1021002"), 0, 1),
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  vendor.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use black market vendor: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected discard selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{item.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve black market vendor: %v", err)
	}
	if len(p0.Graveyard) != 1 || p0.Graveyard[0].InstanceID != item.InstanceID {
		t.Fatalf("selected item should be discarded, graveyard=%v", len(p0.Graveyard))
	}
	if len(p0.Hand) != 2 {
		t.Fatalf("vendor should draw two after discarding, hand=%d", len(p0.Hand))
	}
}

func TestFireArtistResetsAnotherFireCard(t *testing.T) {
	engine := setupReportedBugEngine(t)
	artist := placeUnit(baseCard(t, "1121010"), 0, 0, 0, engine)
	target := placeUnit(baseCard(t, "1121002"), 0, 1, 0, engine)
	target.IsHorizontal = true

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  artist.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use fire artist: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected reset selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{target.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve fire artist: %v", err)
	}
	if !artist.IsHorizontal {
		t.Fatalf("fire artist should tap itself")
	}
	if target.IsHorizontal {
		t.Fatalf("target fire card should be reset")
	}
}

func TestWindriderDiscardsAnyNumberForAir(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	windrider := placeUnit(baseCard(t, "1321005"), 0, 0, 0, engine)
	cardA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
	cardB := NewCardInstance(baseCard(t, "1021002"), 0, 1)
	cardC := NewCardInstance(baseCard(t, "1021003"), 0, 1)
	p0.Hand = []*CardInstance{cardA, cardB, cardC}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  windrider.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use windrider: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{cardA.InstanceID, cardB.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve windrider: %v", err)
	}
	if p0.Elements[model.ElementAir] != 2 {
		t.Fatalf("windrider should gain 1 air per discarded card, elements=%v", p0.Elements)
	}
	if len(p0.Hand) != 1 || len(p0.Graveyard) != 2 {
		t.Fatalf("windrider should discard exactly two selected cards, hand=%d graveyard=%d", len(p0.Hand), len(p0.Graveyard))
	}
}

func TestSunwheelMageResetsLightSkill(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	mage := placeUnit(baseCard(t, "1521011"), 0, 0, 0, engine)
	lightSkill := readySkill(baseCard(t, "3521001"), 0)
	lightSkill.IsHorizontal = true
	p0.Skills[0] = lightSkill

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  mage.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use sunwheel mage: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{lightSkill.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve sunwheel mage: %v", err)
	}
	if lightSkill.IsHorizontal {
		t.Fatalf("light skill should be reset")
	}
}

func TestSoulPriestSacrificesCompanionAndDrawsTwo(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	priest := placeUnit(baseCard(t, "1621012"), 0, 0, 0, engine)
	food := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
	p0.Deck = []*CardInstance{
		NewCardInstance(baseCard(t, "1021002"), 0, 1),
		NewCardInstance(baseCard(t, "1021003"), 0, 1),
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  priest.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use soul priest: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{food.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve soul priest: %v", err)
	}
	if p0.Units[1][0] != nil {
		t.Fatalf("sacrificed companion should leave battlefield")
	}
	if len(p0.Hand) != 2 {
		t.Fatalf("soul priest should draw two, hand=%d", len(p0.Hand))
	}
}

func TestWhimWandConsumesToResetLowCostSkillFromEquipmentSlot(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	wand := NewCardInstance(baseCard(t, "2021003"), 0, 1)
	wand.IsHorizontal = false
	p0.Equipment[0] = wand
	lowCost := readySkill(baseCard(t, "3021005"), 0)
	lowCost.IsHorizontal = true
	p0.Skills[0] = lowCost

	if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{
		"instance_id": wand.InstanceID,
	}}); err != nil {
		t.Fatalf("consume whim wand: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected low-cost skill selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{lowCost.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve whim wand: %v", err)
	}
	if !wand.IsHorizontal {
		t.Fatalf("whim wand should tap")
	}
	if lowCost.IsHorizontal {
		t.Fatalf("selected low-cost skill should be reset")
	}
}

func TestFireSpriteGainsBurnWhenConsumed(t *testing.T) {
	engine := setupReportedBugEngine(t)
	sprite := placeUnit(baseCard(t, "1121001"), 0, 0, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{
		"instance_id": sprite.InstanceID,
	}}); err != nil {
		t.Fatalf("consume fire sprite: %v", err)
	}
	if sprite.Statuses[StatusBurn] != 1 {
		t.Fatalf("fire sprite should gain burn 1 when consumed, statuses=%v", sprite.Statuses)
	}
}

func TestIssue29PlaytestRegressions(t *testing.T) {
	t.Run("explicit first player drives first turn rules", func(t *testing.T) {
		if cards.CardDB == nil {
			if err := cards.LoadCards(); err != nil {
				t.Fatalf("load cards: %v", err)
			}
		}
		SetCardDB(cards.CardDB)
		engine := NewEngine("issue29-first-player", nil)
		deck, err := model.ParseDeckCode(testDeckCode)
		if err != nil {
			t.Fatalf("parse deck: %v", err)
		}
		if err := engine.SetupGameWithFirstPlayer("P1", deck, "P2", deck, 1); err != nil {
			t.Fatalf("setup game: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}}); err != nil {
			t.Fatalf("p0 mulligan: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}}); err != nil {
			t.Fatalf("p1 mulligan: %v", err)
		}
		if engine.State.FirstPlayer != 1 || engine.State.CurrentTurn != 1 {
			t.Fatalf("expected player 1 to start, first=%d current=%d", engine.State.FirstPlayer, engine.State.CurrentTurn)
		}
		state := engine.GetStateForPlayer(0)
		if state["first_player"] != 1 {
			t.Fatalf("serialized first_player should be 1, state=%v", state["first_player"])
		}
		order := state["turn_order"].(map[string]string)
		if order["you"] != "后手" || order["opponent"] != "先手" {
			t.Fatalf("turn labels should follow randomized first player, order=%v", order)
		}
	})

	t.Run("fire sprite only burns itself when it is the consumed card", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		sprite := placeUnit(baseCard(t, "1121001"), 0, 0, 0, engine)
		other := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{
			"instance_id": other.InstanceID,
		}}); err != nil {
			t.Fatalf("consume other card: %v", err)
		}
		if sprite.Statuses[StatusBurn] != 0 {
			t.Fatalf("fire sprite should not burn when another card is consumed, statuses=%v", sprite.Statuses)
		}
	})

	t.Run("fire barrier boosts fire boost spells and grants burn from boosted hit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		barrier := readySkill(baseCard(t, "3121008"), 0)
		main := readySkill(baseCard(t, "3321005"), 0)
		boost := readySkill(baseCard(t, "3121015"), 0)
		p0.Skills[0] = barrier
		p0.Skills[1] = main
		p0.Skills[2] = boost

		withoutBarrier := main.Card.Power + boost.Card.Power
		withBarrier := engine.effectiveSpellPower(0, main, []*CardInstance{boost}, SpellTarget{Type: "unit", Position: Position{Col: 1, Row: 0}})
		if withBarrier != withoutBarrier+2 {
			t.Fatalf("fire barrier should add +2 power to fire boost spell, got=%d want=%d", withBarrier, withoutBarrier+2)
		}
		p0.Elements[model.ElementAir] = 3
		p0.Elements[model.ElementFire] = 3
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": main.InstanceID,
			"boost_ids":   []any{boost.InstanceID},
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast boosted spell: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
			t.Fatalf("no defend: %v", err)
		}
		if target.Statuses[StatusBurn] != 1 {
			t.Fatalf("fire barrier should add burn from fire boost spell, statuses=%v", target.Statuses)
		}
	})

	t.Run("Maggie and geography primer expose effective play costs without consuming discounts", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		maggie := NewCardInstance(baseCard(t, "4411003"), 0, 1)
		p0.Hero = maggie
		primer := NewCardInstance(baseCard(t, "2421013"), 0, 1)
		p0.Equipment[0] = primer
		rock := NewCardInstance(baseCard(t, "1421005"), 0, 1)
		p0.Hand = []*CardInstance{rock}
		p0.Elements[model.ElementEarth] = 2

		cost := engine.effectiveCardPlayCost(p0, rock)
		if cost[model.ElementEarth] != 2 || maggie.Statuses["麦吉折扣"] != 0 {
			t.Fatalf("effective cost should include both discounts without marking Maggie used, cost=%v statuses=%v", cost, maggie.Statuses)
		}
		state := engine.GetStateForPlayer(0)
		hand := state["you"].(map[string]any)["hand"].([]map[string]any)
		effective := hand[0]["effective_elements_cost"].(map[string]int)
		if effective[model.ElementEarth] != 2 {
			t.Fatalf("frontend state should expose effective play cost, cost=%v", effective)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": rock.InstanceID,
			"col":         float64(1),
			"row":         float64(1),
			"payment":     map[string]any{model.ElementEarth: float64(2)},
		}}); err != nil {
			t.Fatalf("summon rock with discounted cost: %v", err)
		}
		if maggie.Statuses["麦吉折扣"] != 1 {
			t.Fatalf("Maggie should be marked used only after the discounted play succeeds, statuses=%v", maggie.Statuses)
		}
	})
}

func TestBottledElementGainsOneArcaneWhenUsed(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	item := NewCardInstance(baseCard(t, "2021005"), 0, 1)
	p0.Hand = []*CardInstance{item}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": item.InstanceID,
	}}); err != nil {
		t.Fatalf("use bottled element: %v", err)
	}
	if p0.Elements[model.ElementArcane] != 1 {
		t.Fatalf("bottled element should gain 1 arcane, elements=%v", p0.Elements)
	}
	if len(p0.Hand) != 0 || len(p0.Graveyard) != 1 {
		t.Fatalf("used consumable should move from hand to graveyard, hand=%d graveyard=%d", len(p0.Hand), len(p0.Graveyard))
	}
}

func TestDefenseOverexertPaysOnlyForThisDefenseWithoutConsumeTriggers(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	fireSprite := placeUnit(baseCard(t, "1121001"), 1, 0, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Elements[model.ElementAir] = 2
	p1.Skills[0] = readySkill(baseCard(t, "3021008"), 1)
	p1.Elements[model.ElementArcane] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast cyclone wave: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
		"skill_ids":     []any{p1.Skills[0].InstanceID},
		"overexert_ids": []any{fireSprite.InstanceID},
	}}); err != nil {
		t.Fatalf("defend by overexerting fire sprite: %v", err)
	}
	if !fireSprite.IsHorizontal {
		t.Fatalf("overexerted unit should become horizontal")
	}
	if fireSprite.Statuses[StatusBurn] != 0 {
		t.Fatalf("overexertion is not consume and should not trigger fire sprite burn, statuses=%v", fireSprite.Statuses)
	}
	if p1.Elements[model.ElementArcane] != 0 || p1.Elements[model.ElementFire] != 0 {
		t.Fatalf("defense overexert should not store leftover load, elements=%v", p1.Elements)
	}
	if target.CurrentLife != target.Card.Life {
		t.Fatalf("successful defense should prevent spell hit, target life=%d", target.CurrentLife)
	}
}

func TestDefenseCanOverexertEquipment(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	equipment := NewCardInstance(baseCard(t, "2011001"), 1, 1)
	equipment.IsHorizontal = false
	p1.Equipment[0] = equipment
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Elements[model.ElementAir] = 2
	p1.Skills[0] = readySkill(baseCard(t, "3021008"), 1)

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast cyclone wave: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
		"skill_ids":     []any{p1.Skills[0].InstanceID},
		"overexert_ids": []any{equipment.InstanceID},
	}}); err != nil {
		t.Fatalf("defend by overexerting equipment: %v", err)
	}
	if !equipment.IsHorizontal {
		t.Fatalf("overexerted equipment should become horizontal")
	}
	if p1.Elements[model.ElementArcane] != 0 {
		t.Fatalf("equipment overexert should not store leftover load, elements=%v", p1.Elements)
	}
	if target.CurrentLife != target.Card.Life {
		t.Fatalf("successful defense should prevent spell hit, target life=%d", target.CurrentLife)
	}
}

func TestIssue27DarkPlaytestRegressions(t *testing.T) {
	t.Run("bloodsuck does not trigger its on-hit buff when used only as a boost", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		friend := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3021005"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3621002"), 0)
		p0.Elements[model.ElementArcane] = 2
		p0.Elements[model.ElementShadow] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"boost_ids":   []any{p0.Skills[1].InstanceID},
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast boosted arcane arrow: %v", err)
		}
		if engine.State.Phase == PhaseDefenseWindow {
			if err := engine.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
				t.Fatalf("no defend: %v", err)
			}
		}
		if engine.State.PendingAction != nil && engine.State.PendingAction.Type == "bloodsuck_buff" {
			t.Fatalf("bloodsuck boost should not open its on-hit buff pending action")
		}
		if friend.CurrentLife != friend.Card.Life {
			t.Fatalf("friendly unit should not receive bloodsuck buff from boost, life=%d", friend.CurrentLife)
		}
		if target.CurrentLife >= target.Card.Life {
			t.Fatalf("boosted spell should still hit the enemy target, life=%d", target.CurrentLife)
		}
	})

	t.Run("Alice triggers after a friendly companion dies and boosts a friendly spell", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		alice := NewCardInstance(baseCard(t, "4611001"), 0, 1)
		p0.Hero = alice
		ally := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3021005"), 0)

		engine.destroyUnit(ally, 0)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "alice_boost_spell" {
			t.Fatalf("Alice should open spell boost choice after friendly death, pending=%v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{p0.Skills[0].InstanceID},
		}}); err != nil {
			t.Fatalf("resolve Alice boost: %v", err)
		}
		if p0.Skills[0].PowerBonus != 1 {
			t.Fatalf("Alice should give selected spell +1 power, bonus=%d", p0.Skills[0].PowerBonus)
		}
	})

	t.Run("blood demon blast can be learned from the skill pool", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		blast := readySkill(baseCard(t, "3621010"), 0)
		p0.SkillPool = []*CardInstance{blast}
		p0.Elements[model.ElementShadow] = 4

		if err := engine.HandleAction(0, ActionMessage{Action: "learn_skill", Data: map[string]any{
			"instance_id": blast.InstanceID,
		}}); err != nil {
			t.Fatalf("learn blood demon blast: %v", err)
		}
		if p0.Skills[0] != blast {
			t.Fatalf("blood demon blast should move into skill slot")
		}
	})
}

func TestOpponentDeckHiddenButRevealedHandVisible(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	revealed := NewCardInstance(baseCard(t, "1321003"), 1, 1)
	hidden := NewCardInstance(baseCard(t, "1021001"), 1, 1)
	p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021002"), 0, 1)}
	p1.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021003"), 1, 1)}
	p1.Hand = []*CardInstance{revealed, hidden}
	p1.RevealedHand[revealed.InstanceID] = true

	state := engine.GetStateForPlayer(0)
	you := state["you"].(map[string]any)
	opponent := state["opponent"].(map[string]any)
	if _, ok := you["deck_summary"]; !ok {
		t.Fatalf("owner should see their own unordered deck summary")
	}
	if _, ok := opponent["deck_summary"]; ok {
		t.Fatalf("opponent deck summary should be hidden")
	}
	visible := opponent["revealed_hand"].([]map[string]any)
	if len(visible) != 1 || visible[0]["number"] != "1321003" {
		t.Fatalf("only explicitly revealed opponent hand cards should be visible, got %+v", visible)
	}
}

func TestMagicDandelionRevealsWhenDrawnAndClearsWhenLeavingHand(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	dandelion := NewCardInstance(baseCard(t, "1321003"), 0, 1)
	p0.Deck = []*CardInstance{dandelion}

	drawn := p0.DrawCards(1)
	if len(drawn) != 1 || !p0.RevealedHand[dandelion.InstanceID] {
		t.Fatalf("magic dandelion should be marked revealed when drawn, revealed=%v", p0.RevealedHand)
	}
	_, handIdx := p0.FindHandCard(dandelion.InstanceID)
	if handIdx < 0 {
		t.Fatalf("drawn dandelion should be in hand")
	}
	p0.RemoveFromHand(handIdx)
	if p0.RevealedHand[dandelion.InstanceID] {
		t.Fatalf("revealed marker should clear when card leaves hand, revealed=%v", p0.RevealedHand)
	}
}

func TestLifePotionUsesPendingSelectionToHealFriendlyUnit(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	item := NewCardInstance(baseCard(t, "2521001"), 0, 1)
	p0.Hand = []*CardInstance{item}
	p0.Elements[model.ElementArcane] = 1
	target := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
	target.CurrentLife = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": item.InstanceID,
	}}); err != nil {
		t.Fatalf("use life potion: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected life potion target selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{target.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve life potion: %v", err)
	}
	if target.CurrentLife != 3 {
		t.Fatalf("life potion should heal up to max life, life=%d", target.CurrentLife)
	}
	if p0.Elements[model.ElementArcane] != 0 {
		t.Fatalf("life potion should spend one arcane wildcard for neutral cost, elements=%v", p0.Elements)
	}
}

func TestWindcallingScrollDrawsTwoAndSkipsNextDraw(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	item := NewCardInstance(baseCard(t, "2321005"), 0, 1)
	p0.Hand = []*CardInstance{item}
	p0.Deck = []*CardInstance{
		NewCardInstance(baseCard(t, "1021001"), 0, 1),
		NewCardInstance(baseCard(t, "1021002"), 0, 1),
		NewCardInstance(baseCard(t, "1021003"), 0, 1),
	}
	p0.Elements[model.ElementAir] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": item.InstanceID,
	}}); err != nil {
		t.Fatalf("use windcalling scroll: %v", err)
	}
	if len(p0.Hand) != 2 || len(p0.Deck) != 1 {
		t.Fatalf("windcalling scroll should draw two, hand=%d deck=%d", len(p0.Hand), len(p0.Deck))
	}
	if !p0.SkipNextDraw {
		t.Fatalf("windcalling scroll should mark next turn draw to be skipped")
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
		t.Fatalf("end p0 turn: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
		t.Fatalf("end p1 turn: %v", err)
	}
	if p0.SkipNextDraw {
		t.Fatalf("skip draw flag should clear on next own turn")
	}
	if len(p0.Deck) != 1 {
		t.Fatalf("next normal draw should be skipped, deck=%d", len(p0.Deck))
	}
}

func TestArcaneArmorerSearchesEquipmentWhenNoEquipment(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	armorer := NewCardInstance(baseCard(t, "1021016"), 0, 1)
	equipment := NewCardInstance(baseCard(t, "2021003"), 0, 1)
	p0.Hand = []*CardInstance{armorer}
	p0.Deck = []*CardInstance{equipment, NewCardInstance(baseCard(t, "1021001"), 0, 1)}
	p0.Elements[model.ElementArcane] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": armorer.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon arcane armorer: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected equipment search selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{equipment.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve armorer search: %v", err)
	}
	if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != equipment.InstanceID {
		t.Fatalf("searched equipment should be added to hand, hand=%v", len(p0.Hand))
	}
}

func TestRunemasterDiscardsHandCardThenSearchesRuneOrScroll(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	runemaster := NewCardInstance(baseCard(t, "1021017"), 0, 1)
	discard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
	scroll := NewCardInstance(baseCard(t, "2321005"), 0, 1)
	p0.Hand = []*CardInstance{runemaster, discard}
	p0.Deck = []*CardInstance{scroll, NewCardInstance(baseCard(t, "1021002"), 0, 1)}
	p0.Elements[model.ElementArcane] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": runemaster.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon runemaster: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected discard selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{discard.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve runemaster discard: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected search selection after discard, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{scroll.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve runemaster search: %v", err)
	}
	if len(p0.Graveyard) != 1 || p0.Graveyard[0].InstanceID != discard.InstanceID {
		t.Fatalf("runemaster should discard selected hand card, graveyard=%d", len(p0.Graveyard))
	}
	if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != scroll.InstanceID {
		t.Fatalf("runemaster should add selected rune or scroll to hand, hand=%d", len(p0.Hand))
	}
}

func TestRedHawkSearchesHighCostFireCompanion(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	hawk := NewCardInstance(baseCard(t, "1121009"), 0, 1)
	target := NewCardInstance(baseCard(t, "1121004"), 0, 1)
	p0.Hand = []*CardInstance{hawk}
	p0.Deck = []*CardInstance{target, NewCardInstance(baseCard(t, "1121002"), 0, 1)}
	p0.Elements[model.ElementFire] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": hawk.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon red hawk: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected fire companion search selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{target.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve red hawk search: %v", err)
	}
	if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != target.InstanceID {
		t.Fatalf("red hawk should search selected high-cost fire companion, hand=%d", len(p0.Hand))
	}
}

func TestPegasusKnightSearchesUnicornPegasus(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	knight := NewCardInstance(baseCard(t, "1521009"), 0, 1)
	pegasus := NewCardInstance(baseCard(t, "1521012"), 0, 1)
	p0.Hand = []*CardInstance{knight}
	p0.Deck = []*CardInstance{pegasus, NewCardInstance(baseCard(t, "1021001"), 0, 1)}
	p0.Elements[model.ElementLight] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": knight.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon pegasus knight: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected pegasus search selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{pegasus.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve pegasus knight search: %v", err)
	}
	if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != pegasus.InstanceID {
		t.Fatalf("pegasus knight should search unicorn pegasus, hand=%d", len(p0.Hand))
	}
}

func TestManaBoosterMakesNextSkillUseFree(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	item := NewCardInstance(baseCard(t, "2021014"), 0, 1)
	p0.Hand = []*CardInstance{item}
	p0.Elements[model.ElementArcane] = 1
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": item.InstanceID,
	}}); err != nil {
		t.Fatalf("use mana booster: %v", err)
	}
	if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextSkillCostZero {
		t.Fatalf("mana booster should add next skill cost modifier, modifiers=%v", p0.TempModifiers)
	}
	if p0.Elements[model.ElementArcane] != 0 {
		t.Fatalf("mana booster item cost should be paid, elements=%v", p0.Elements)
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("next skill should be free: %v", err)
	}
	if len(p0.TempModifiers) != 0 {
		t.Fatalf("free next skill modifier should be consumed, modifiers=%v", p0.TempModifiers)
	}
}

func TestStoneforgeArtisanGivesSelectedSpellPowerThisTurn(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	artisan := placeUnit(baseCard(t, "1121003"), 0, 0, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Elements[model.ElementAir] = 2
	placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{
		"instance_id": artisan.InstanceID,
	}}); err != nil {
		t.Fatalf("consume stoneforge artisan: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected spell selection for power bonus, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{p0.Skills[0].InstanceID},
	}}); err != nil {
		t.Fatalf("resolve stoneforge artisan: %v", err)
	}
	if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModSkillPowerBonus {
		t.Fatalf("stoneforge should add spell power modifier, modifiers=%v", p0.TempModifiers)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast boosted spell: %v", err)
	}
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 3 {
		t.Fatalf("cyclone wave should have base 1 +2 power, pending=%+v", engine.State.PendingSpell)
	}
}

func TestEnergeticSeniorMakesNextSkillIgnoreCooldown(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	senior := NewCardInstance(baseCard(t, "1021015"), 0, 1)
	p0.Hand = []*CardInstance{senior}
	p0.Elements[model.ElementArcane] = 2
	p0.Elements[model.ElementLight] = 1
	p0.Skills[0] = readySkill(baseCard(t, "3021008"), 0)
	placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": senior.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon energetic senior: %v", err)
	}
	if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextNoCooldown {
		t.Fatalf("energetic senior should add no-cooldown modifier, modifiers=%v", p0.TempModifiers)
	}
	p0.Elements[model.ElementArcane] = 2
	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast cooldown spell: %v", err)
	}
	if p0.Skills[0].Statuses[StatusCooldown] > 0 {
		t.Fatalf("next skill should ignore cooldown, statuses=%v", p0.Skills[0].Statuses)
	}
	if len(p0.TempModifiers) != 0 {
		t.Fatalf("no-cooldown modifier should be consumed, modifiers=%v", p0.TempModifiers)
	}
}

func TestEndTurnResetsBeforeSettlingCooldown(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]

	normal := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	normal.IsHorizontal = true

	cooldownSkill := readySkill(baseCard(t, "3021008"), 0)
	cooldownSkill.IsHorizontal = true
	cooldownSkill.Statuses[StatusCooldown] = 1
	p0.Skills[0] = cooldownSkill

	if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
		t.Fatalf("end turn: %v", err)
	}
	if normal.IsHorizontal {
		t.Fatalf("normal horizontal card should reset before mark settlement")
	}
	if !cooldownSkill.IsHorizontal {
		t.Fatalf("cooldown skill should remain horizontal because cooldown blocks reset")
	}
	if cooldownSkill.Statuses[StatusCooldown] != 0 {
		t.Fatalf("cooldown should settle after reset, statuses=%v", cooldownSkill.Statuses)
	}
}

func TestIceDissolveReactsToEnemySpell(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]

	target := placeUnit(baseCard(t, "1221002"), 1, 1, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
	p0.Elements[model.ElementFire] = 1
	p1.Skills[0] = readySkill(baseCard(t, "3221008"), 1)

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast spell: %v", err)
	}
	if engine.State.Phase != PhaseDefenseWindow || engine.State.PendingSpell == nil {
		t.Fatalf("expected defense window, phase=%v pending=%v", engine.State.Phase, engine.State.PendingSpell)
	}
	startPower := engine.State.PendingSpell.TotalPower
	if err := engine.HandleAction(1, ActionMessage{Action: "react_spell", Data: map[string]any{
		"instance_id":   p1.Skills[0].InstanceID,
		"overexert_ids": []any{target.InstanceID},
	}}); err != nil {
		t.Fatalf("react with ice dissolve: %v", err)
	}
	if got := engine.State.PendingSpell.TotalPower; got != 0 {
		t.Fatalf("ice dissolve should zero pending spell power, got %d start %d", got, startPower)
	}
	if p1.Elements[model.ElementWater] != 0 {
		t.Fatalf("ice dissolve should not leave overexerted water in the pool, elements=%v", p1.Elements)
	}
	if !target.IsHorizontal {
		t.Fatalf("overexerted unit should become horizontal")
	}
	if !p1.Skills[0].IsHorizontal || p1.Skills[0].Statuses[StatusCooldown] != 1 {
		t.Fatalf("ice dissolve should tap and gain cooldown1, horizontal=%v statuses=%v", p1.Skills[0].IsHorizontal, p1.Skills[0].Statuses)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve spell: %v", err)
	}
	if target.CurrentLife >= target.Card.Life {
		t.Fatalf("ice dissolve changes spell power, not the spell hit damage")
	}
}

func TestIceDissolveCannotBeCastAsMainPhaseAttack(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]

	p0.Skills[0] = readySkill(baseCard(t, "3221008"), 0)
	p0.Elements[model.ElementWater] = 2

	err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
	}})
	if err == nil {
		t.Fatalf("ice dissolve should only be usable as a spell reaction")
	}
	if p0.Elements[model.ElementWater] != 2 || p0.Skills[0].IsHorizontal {
		t.Fatalf("failed main-phase cast should not spend or tap, elements=%v horizontal=%v", p0.Elements, p0.Skills[0].IsHorizontal)
	}
}

func TestWinterfellWarlockGivesNextSpellFreezeOnHit(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
	warlock := placeUnit(baseCard(t, "1221011"), 0, 0, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
	p0.Elements[model.ElementFire] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  warlock.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use winterfell warlock ultimate: %v", err)
	}
	if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextSpellHitStatus {
		t.Fatalf("warlock should add spell-hit status modifier, modifiers=%v", p0.TempModifiers)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast fireball after warlock ultimate: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve fireball hit: %v", err)
	}
	if target.Statuses[StatusFreeze] != 1 {
		t.Fatalf("next spell should apply freeze 1, statuses=%v", target.Statuses)
	}
	if len(p0.TempModifiers) != 0 {
		t.Fatalf("spell-hit status modifier should be consumed, modifiers=%v", p0.TempModifiers)
	}
}

func TestPassionOfFireDrawsWhenFriendlyFireSpellHits(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p0.Skills[0] = readySkill(baseCard(t, "3121007"), 0)
	p0.Skills[1] = readySkill(baseCard(t, "3121001"), 0)
	p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
	p0.Elements[model.ElementFire] = 1
	placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[1].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast fireball with passion of fire in play: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve fireball hit: %v", err)
	}
	if len(p0.Hand) != 1 {
		t.Fatalf("passion of fire should draw 1 after friendly fire spell hit, hand=%d deck=%d", len(p0.Hand), len(p0.Deck))
	}
}

func TestSpellStatPassivePowerAndDamageModifiers(t *testing.T) {
	t.Run("celtic lion gives all spells plus one power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1121004"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementFire] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 4 {
			t.Fatalf("fireball should have base 3 +1 power, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("raincaller gives water and air spells plus one power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1221013"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
		p0.Elements[model.ElementAir] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast cyclone wave: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 2 {
			t.Fatalf("cyclone wave should have base 1 +1 power, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("thunder beast gives air spells plus one damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		target.CurrentLife = 3
		placeUnit(baseCard(t, "1321006"), 0, 0, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
		p0.Elements[model.ElementAir] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast cyclone wave: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve cyclone wave: %v", err)
		}
		if target.CurrentLife != 1 {
			t.Fatalf("cyclone wave should deal base 1 +1 damage, target life=%d", target.CurrentLife)
		}
	})

	t.Run("divine fire beast gives attacking spells plus two power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1521013"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementFire] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 5 {
			t.Fatalf("fireball should have base 3 +2 attacking power, pending=%+v", engine.State.PendingSpell)
		}
	})
}

func TestUnitEnterListenersCanReactToEnemySummons(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	placeUnit(baseCard(t, "1111002"), 0, 0, 0, engine)
	summoned := NewCardInstance(baseCard(t, "1021001"), 1, 1)
	p1.Hand = []*CardInstance{summoned}
	p1.Elements[model.ElementArcane] = 10
	engine.State.CurrentTurn = 1

	if err := engine.HandleAction(1, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": summoned.InstanceID,
		"col":         float64(1),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon enemy unit: %v", err)
	}
	if summoned.Statuses[StatusBurn] != 1 || summoned.Statuses[StatusPetrify] != 1 {
		t.Fatalf("inferno general should burn and petrify enemy summons, statuses=%v", summoned.Statuses)
	}
	if p0.Hero != nil && p0.Hero.CurrentLife <= 0 {
		t.Fatalf("friendly-only unit-enter listeners should not misfire from enemy summon")
	}
}

func TestTwinAngelCreatesTwinInHandOnEnter(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	angel := NewCardInstance(baseCard(t, "1521005"), 0, 1)
	p0.Hand = []*CardInstance{angel}
	p0.Elements[model.ElementLight] = 10

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": angel.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon twin angel: %v", err)
	}
	if len(p0.Hand) != 1 || p0.Hand[0].Card.Number != "1501001" {
		t.Fatalf("twin angel should create card 1501001 in hand, hand=%v", cardsToInfo(p0.Hand))
	}
}

func TestRecyclingSpriteMovesSelectedGraveyardCardToDeckTop(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	sprite := NewCardInstance(baseCard(t, "1021007"), 0, 1)
	recovered := NewCardInstance(baseCard(t, "1021001"), 0, 1)
	otherDeckCard := NewCardInstance(baseCard(t, "1021002"), 0, 1)
	p0.Hand = []*CardInstance{sprite}
	p0.Graveyard = []*CardInstance{recovered}
	p0.Deck = []*CardInstance{otherDeckCard}
	p0.Elements[model.ElementArcane] = 10

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": sprite.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon recycling sprite: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction || engine.State.PendingAction == nil {
		t.Fatalf("recycling sprite should ask which graveyard card to recover")
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{recovered.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve recycling sprite action: %v", err)
	}
	if len(p0.Graveyard) != 0 {
		t.Fatalf("graveyard card should be removed, graveyard=%v", cardsToInfo(p0.Graveyard))
	}
	if len(p0.Deck) == 0 || p0.Deck[0] != recovered {
		t.Fatalf("selected graveyard card should be deck top, deck=%v", cardsToInfo(p0.Deck))
	}
}

func TestLifeFlowerBuffsAnotherFriendlyUnit(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021004"), 0, 1, 0, engine)
	target.CurrentLife = 2
	flower := NewCardInstance(baseCard(t, "1521006"), 0, 1)
	p0.Hand = []*CardInstance{flower}
	p0.Elements[model.ElementLight] = 10

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": flower.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon life flower: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{target.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve life flower: %v", err)
	}
	if target.CurrentLife != 3 {
		t.Fatalf("life flower should give +1 life, life=%d", target.CurrentLife)
	}
}

func TestDarkDeathEffects(t *testing.T) {
	t.Run("elemental husk gains none element on death", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		husk := placeUnit(baseCard(t, "1621002"), 0, 0, 0, engine)

		engine.dealDamage(husk, 99, 0)
		if p0.Elements[model.ElementArcane] != 1 {
			t.Fatalf("elemental husk should gain 1 none element, elements=%v", p0.Elements)
		}
	})

	t.Run("nightmare gains life when another friendly unit dies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		nightmare := placeUnit(baseCard(t, "1621006"), 0, 0, 0, engine)
		other := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		startLife := nightmare.CurrentLife

		engine.dealDamage(other, 99, 0)
		if nightmare.CurrentLife != startLife+1 {
			t.Fatalf("nightmare should gain +1 life, life=%d start=%d", nightmare.CurrentLife, startLife)
		}
	})

	t.Run("bone knight resummons once and loses deathrattle", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		knight := placeUnit(baseCard(t, "1621011"), 0, 0, 0, engine)

		if !cardHasActiveDeathrattle(knight) {
			t.Fatalf("bone knight should start with an active deathrattle")
		}
		engine.dealDamage(knight, 99, 0)
		if engine.State.Players[0].Units[0][0] != knight {
			t.Fatalf("bone knight should return to its position")
		}
		if knight.Statuses[boneKnightRebornStatus] != 1 {
			t.Fatalf("bone knight should lose deathrattle after return, statuses=%v", knight.Statuses)
		}
		if cardHasActiveDeathrattle(knight) {
			t.Fatalf("bone knight should no longer count as a deathrattle unit after returning")
		}
		engine.dealDamage(knight, 99, 0)
		if engine.State.Players[0].Units[0][0] != nil {
			t.Fatalf("bone knight should not return a second time")
		}
	})

	t.Run("runtime attached deathrattle counts as deathrattle and resolves", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p1 := engine.State.Players[1]
		p1.Hero = NewCardInstance(baseCard(t, "4311003"), 1, engine.State.TurnNumber)
		unit := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)

		if cardHasActiveDeathrattle(unit) {
			t.Fatalf("plain apprentice should not start as a deathrattle unit")
		}
		unit.AddAttachedBehavior(AttachedDeathrattleDamageEnemyHero{Amount: 1})
		if !cardHasActiveDeathrattle(unit) {
			t.Fatalf("attached deathrattle should make the unit count as a deathrattle unit")
		}
		startLife := p1.Hero.CurrentLife

		engine.destroyUnit(unit, 0)
		if p1.Hero.CurrentLife != startLife-1 {
			t.Fatalf("attached deathrattle should damage enemy hero by 1, got %d want %d", p1.Hero.CurrentLife, startLife-1)
		}
	})
}

func TestOnSpellCastListenersFireForFriendlySkills(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	placeUnit(baseCard(t, "1321012"), 0, 0, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
	p0.Elements[model.ElementAir] = 2
	placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast air spell: %v", err)
	}
	if len(p0.Hand) != 1 {
		t.Fatalf("wind medium should draw when friendly air skill is used, hand=%d deck=%d", len(p0.Hand), len(p0.Deck))
	}
}

func TestEquipmentEnterEffectsAndSpellStatPassives(t *testing.T) {
	t.Run("equipment on-enter triggers when equipped", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
		ring := NewCardInstance(baseCard(t, "2321007"), 0, 1)
		p0.Hand = []*CardInstance{ring}
		p0.Elements[model.ElementAir] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
			"instance_id": ring.InstanceID,
		}}); err != nil {
			t.Fatalf("equip windwhisper ring: %v", err)
		}
		if len(p0.Hand) != 1 {
			t.Fatalf("ring on-enter should draw 1, hand=%d deck=%d", len(p0.Hand), len(p0.Deck))
		}
	})

	t.Run("wizard scepter gives spell plus one power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2021004"), 0, 1)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementFire] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 4 {
			t.Fatalf("wizard scepter should give +1 power, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("manes staff only buffs water spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2221004"), 0, 1)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3221001"), 0)
		p0.Elements[model.ElementWater] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast hail: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 2 {
			t.Fatalf("manes staff should give water spell +1 power, pending=%+v", engine.State.PendingSpell)
		}
	})
}

func TestEquipmentLifeAndEnemySummonCounters(t *testing.T) {
	t.Run("life amulet buffs selected friendly role", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
		amulet := NewCardInstance(baseCard(t, "2021011"), 0, 1)
		p0.Hand = []*CardInstance{amulet}
		p0.Elements[model.ElementArcane] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
			"instance_id": amulet.InstanceID,
		}}); err != nil {
			t.Fatalf("equip life amulet: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve life amulet: %v", err)
		}
		if target.CurrentLife != target.Card.Life+1 {
			t.Fatalf("life amulet should give +1 life, life=%d", target.CurrentLife)
		}
	})

	t.Run("rattan cuirass buffs hero on enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Hero = placeUnit(baseCard(t, "4311003"), 0, 1, 1, engine)
		startLife := p0.Hero.CurrentLife
		cuirass := NewCardInstance(baseCard(t, "2421006"), 0, 1)
		p0.Hand = []*CardInstance{cuirass}
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
			"instance_id": cuirass.InstanceID,
		}}); err != nil {
			t.Fatalf("equip rattan cuirass: %v", err)
		}
		if p0.Hero.CurrentLife != startLife+2 {
			t.Fatalf("rattan cuirass should give hero +2 life, life=%d start=%d", p0.Hero.CurrentLife, startLife)
		}
	})

	t.Run("hellfire rune statuses enemy summons", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2121012"), 0, 1)
		summoned := NewCardInstance(baseCard(t, "1021001"), 1, 1)
		p1.Hand = []*CardInstance{summoned}
		p1.Elements[model.ElementArcane] = 10
		engine.State.CurrentTurn = 1

		if err := engine.HandleAction(1, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": summoned.InstanceID,
			"col":         float64(0),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("enemy summon: %v", err)
		}
		if summoned.Statuses[StatusStun] != 2 || summoned.Statuses[StatusPetrify] != 2 || summoned.Statuses[StatusBurn] != 2 {
			t.Fatalf("hellfire rune should apply three statuses, statuses=%v", summoned.Statuses)
		}
	})

	t.Run("killing instinct damages enemy summons", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p1 := engine.State.Players[1]
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2621003"), 0, 1)
		summoned := NewCardInstance(baseCard(t, "1021004"), 1, 1)
		p1.Hand = []*CardInstance{summoned}
		p1.Elements[model.ElementArcane] = 10
		engine.State.CurrentTurn = 1

		if err := engine.HandleAction(1, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": summoned.InstanceID,
			"col":         float64(0),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("enemy summon: %v", err)
		}
		if summoned.CurrentLife != summoned.Card.Life-2 {
			t.Fatalf("killing instinct should deal 2 damage, life=%d", summoned.CurrentLife)
		}
	})
}

func TestConsumableTargetedItemEffects(t *testing.T) {
	t.Run("fire arrow damages selected enemy", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		item := NewCardInstance(baseCard(t, "2121004"), 0, 1)
		p0.Hand = []*CardInstance{item}
		p0.Elements[model.ElementFire] = 10
		target := placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": item.InstanceID,
		}}); err != nil {
			t.Fatalf("use fire arrow item: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve fire arrow: %v", err)
		}
		if target.CurrentLife != target.Card.Life-1 {
			t.Fatalf("fire arrow should deal 1 damage, life=%d", target.CurrentLife)
		}
	})

	t.Run("bottled lightning gains air and stuns friendly", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		item := NewCardInstance(baseCard(t, "2321006"), 0, 1)
		p0.Hand = []*CardInstance{item}
		p0.Elements[model.ElementAir] = 10
		target := placeUnit(baseCard(t, "1021004"), 0, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": item.InstanceID,
		}}); err != nil {
			t.Fatalf("use bottled lightning: %v", err)
		}
		if p0.Elements[model.ElementAir] != 12 {
			t.Fatalf("bottled lightning should net +3 air after paying 1, elements=%v", p0.Elements)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve bottled lightning: %v", err)
		}
		if target.Statuses[StatusStun] != 2 {
			t.Fatalf("bottled lightning should stun friendly 2, statuses=%v", target.Statuses)
		}
	})

	t.Run("purification scroll removes friendly negative statuses", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		item := NewCardInstance(baseCard(t, "2521003"), 0, 1)
		p0.Hand = []*CardInstance{item}
		p0.Elements[model.ElementLight] = 10
		target := placeUnit(baseCard(t, "1021004"), 0, 1, 0, engine)
		target.Statuses[StatusBurn] = 1
		target.Statuses[StatusFreeze] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": item.InstanceID,
		}}); err != nil {
			t.Fatalf("use purification scroll: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve purification scroll: %v", err)
		}
		if target.Statuses[StatusBurn] != 0 || target.Statuses[StatusFreeze] != 0 {
			t.Fatalf("purification scroll should remove negative statuses, statuses=%v", target.Statuses)
		}
	})
}

func TestSkillPendingChoiceEffects(t *testing.T) {
	t.Run("engrave discards hand card and draws first rune or scroll", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		skill := readySkill(baseCard(t, "3021004"), 0)
		discard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		scroll := NewCardInstance(baseCard(t, "2321005"), 0, 1)
		p0.Skills[0] = skill
		p0.Hand = []*CardInstance{discard}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021002"), 0, 1), scroll}
		p0.Elements[model.ElementArcane] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast engrave: %v", err)
		}
		if engine.State.Phase != PhaseWaitingAction {
			t.Fatalf("engrave should wait for discard choice, phase=%v", engine.State.Phase)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{discard.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve engrave: %v", err)
		}
		if len(p0.Graveyard) != 1 || p0.Graveyard[0] != discard {
			t.Fatalf("engrave should discard selected card, graveyard=%v", cardsToInfo(p0.Graveyard))
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != scroll {
			t.Fatalf("engrave should draw first rune or scroll, hand=%v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("healing spell heals selected friendly unit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
		target.CurrentLife = 1
		skill := readySkill(baseCard(t, "3521001"), 0)
		p0.Skills[0] = skill
		p0.Elements[model.ElementLight] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast healing: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve healing: %v", err)
		}
		if target.CurrentLife != 3 {
			t.Fatalf("healing should restore 2 life, life=%d", target.CurrentLife)
		}
	})

	t.Run("bloodsuck keeps pending choice after spell hit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		ally := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		skill := readySkill(baseCard(t, "3621002"), 0)
		p0.Skills[0] = skill
		p0.Elements[model.ElementShadow] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast bloodsuck: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve bloodsuck hit: %v", err)
		}
		if target.CurrentLife >= target.Card.Life {
			t.Fatalf("bloodsuck should still deal spell damage, target life=%d", target.CurrentLife)
		}
		if engine.State.Phase != PhaseWaitingAction {
			t.Fatalf("bloodsuck should keep pending choice after hit, phase=%v", engine.State.Phase)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{ally.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve bloodsuck buff: %v", err)
		}
		if ally.CurrentLife != ally.Card.Life+2 {
			t.Fatalf("bloodsuck should give friendly +2 life, life=%d", ally.CurrentLife)
		}
	})
}

func TestMoreSkillChoiceEffects(t *testing.T) {
	t.Run("disarm destroys selected enemy equipment after hit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		equipment := NewCardInstance(baseCard(t, "2021004"), 1, 1)
		p1.Equipment[0] = equipment
		skill := readySkill(baseCard(t, "3021008"), 0)
		p0.Skills[0] = skill
		p0.Elements[model.ElementArcane] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast disarm: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve disarm hit: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{equipment.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve disarm equipment: %v", err)
		}
		if p1.Equipment[0] != nil || len(p1.Graveyard) != 1 {
			t.Fatalf("disarm should destroy selected equipment, equipment=%v graveyard=%v", p1.Equipment[0], cardsToInfo(p1.Graveyard))
		}
	})

	t.Run("blessing of light buffs friendly companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		ally := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
		skill := readySkill(baseCard(t, "3521014"), 0)
		p0.Skills[0] = skill
		p0.Elements[model.ElementLight] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast blessing of light: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{ally.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve blessing: %v", err)
		}
		if ally.CurrentLife != ally.Card.Life+1 || effectiveElementsGain(ally)[model.ElementLight] != 1 {
			t.Fatalf("blessing should add life and light load, life=%d gains=%v", ally.CurrentLife, effectiveElementsGain(ally))
		}
	})

	t.Run("soul recall returns up to two companions from graveyard", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		a := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		b := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		p0.Graveyard = []*CardInstance{a, b}
		skill := readySkill(baseCard(t, "3621012"), 0)
		p0.Skills[0] = skill
		p0.Elements[model.ElementShadow] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast soul recall: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{a.InstanceID, b.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve soul recall: %v", err)
		}
		if len(p0.Hand) != 2 || len(p0.Graveyard) != 0 {
			t.Fatalf("soul recall should move two companions to hand, hand=%v graveyard=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard))
		}
	})
}

func TestSpellPowerItemModifiersAndRestrictions(t *testing.T) {
	t.Run("oath ring reduces attacking spell power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2021009"), 0, 1)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementFire] = 1
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 1 {
			t.Fatalf("oath ring should reduce fireball power from 3 to 1, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("severing blade boosts attack and prevents defense", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2021013"), 0, 1)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementFire] = 1
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 5 {
			t.Fatalf("severing blade should boost fireball power from 3 to 5, pending=%+v", engine.State.PendingSpell)
		}

		engine = setupReportedBugEngine(t)
		p1 = engine.State.Players[1]
		p1.Equipment[0] = NewCardInstance(baseCard(t, "2021013"), 1, 1)
		p1.Skills[0] = readySkill(baseCard(t, "2121009"), 1)
		p1.Elements[model.ElementFire] = 10
		engine.State.CurrentTurn = 0
		engine.State.Phase = PhaseDefenseWindow
		engine.State.PendingSpell = &SpellCast{AttackerID: 0, Skill: readySkill(baseCard(t, "3121001"), 0), Target: SpellTarget{Type: "unit", Position: Position{Col: 1, Row: 0}}, TotalPower: 1}

		err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{p1.Skills[0].InstanceID},
		}})
		if err == nil {
			t.Fatalf("severing blade should prevent using spells for defense")
		}
	})

	t.Run("divine flame potion gives temporary power and burns hero", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Hero = placeUnit(baseCard(t, "4311003"), 0, 1, 1, engine)
		item := NewCardInstance(baseCard(t, "2121005"), 0, 1)
		p0.Hand = []*CardInstance{item}
		p0.Elements[model.ElementFire] = 10
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": item.InstanceID,
		}}); err != nil {
			t.Fatalf("use divine flame potion: %v", err)
		}
		if p0.Hero.Statuses[StatusBurn] != 1 {
			t.Fatalf("divine flame potion should burn hero, statuses=%v", p0.Hero.Statuses)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast boosted fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 5 {
			t.Fatalf("divine flame potion should give +2 power, pending=%+v", engine.State.PendingSpell)
		}
	})
}

func TestConsumeAndDeathRuneEffects(t *testing.T) {
	t.Run("fire rune burns any consumed unit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2121002"), 0, 1)
		unit := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		engine.State.CurrentTurn = 1

		if err := engine.HandleAction(1, ActionMessage{Action: "consume", Data: map[string]any{
			"instance_id": unit.InstanceID,
		}}); err != nil {
			t.Fatalf("consume enemy unit with fire rune watching: %v", err)
		}
		if p1.Units[1][0].Statuses[StatusBurn] != 1 {
			t.Fatalf("fire rune should burn consumed unit, statuses=%v", p1.Units[1][0].Statuses)
		}
	})

	t.Run("frost rune freezes consumed enemy partner", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2221002"), 0, 1)
		unit := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		engine.State.CurrentTurn = 1

		if err := engine.HandleAction(1, ActionMessage{Action: "consume", Data: map[string]any{
			"instance_id": unit.InstanceID,
		}}); err != nil {
			t.Fatalf("consume enemy unit with frost rune watching: %v", err)
		}
		if unit.Statuses[StatusFreeze] != 1 {
			t.Fatalf("frost rune should freeze consumed enemy unit, statuses=%v", unit.Statuses)
		}
	})

	t.Run("lightning rune stuns consumed enemy and adjacent unit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2321002"), 0, 1)
		unit := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		adjacent := placeUnit(baseCard(t, "1021004"), 1, 0, 0, engine)
		engine.State.CurrentTurn = 1

		if err := engine.HandleAction(1, ActionMessage{Action: "consume", Data: map[string]any{
			"instance_id": unit.InstanceID,
		}}); err != nil {
			t.Fatalf("consume enemy unit with lightning rune watching: %v", err)
		}
		if unit.Statuses[StatusStun] != 1 || adjacent.Statuses[StatusStun] != 1 {
			t.Fatalf("lightning rune should stun consumed and adjacent units, unit=%v adjacent=%v", unit.Statuses, adjacent.Statuses)
		}
	})

	t.Run("sacrifice rune draws when friendly partner dies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2621005"), 0, 1)
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1021004"), 0, 1),
		}
		unit := placeUnit(baseCard(t, "1021007"), 0, 0, 0, engine)

		engine.destroyUnit(unit, 0)
		if len(p0.Hand) != 2 {
			t.Fatalf("sacrifice rune should draw 2, hand=%d deck=%d", len(p0.Hand), len(p0.Deck))
		}
	})
}

func TestPhoenixFeatherCountersAndPerTurnElement(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	feather := NewCardInstance(baseCard(t, "2121001"), 0, 1)
	p0.Equipment[0] = feather

	engine.triggerEffects(TriggerOnEnter, feather, nil, nil)
	if feather.Statuses[phoenixFeatherCounter] != 3 {
		t.Fatalf("phoenix feather should enter with 3 counters, statuses=%v", feather.Statuses)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  feather.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use phoenix feather per-turn ability: %v", err)
	}
	if feather.Statuses[phoenixFeatherCounter] != 2 || p0.Elements[model.ElementFire] != 1 {
		t.Fatalf("phoenix feather should remove counter and gain fire, counters=%v elements=%v", feather.Statuses, p0.Elements)
	}
}

func TestCycloneScrollDestroysLowCostEnemyEquipment(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	scroll := NewCardInstance(baseCard(t, "2321008"), 0, 1)
	p0.Hand = []*CardInstance{scroll}
	p0.Elements[model.ElementAir] = 3
	p1.Equipment[0] = NewCardInstance(baseCard(t, "2021011"), 1, 1)

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": scroll.InstanceID,
	}}); err != nil {
		t.Fatalf("use cyclone scroll: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction || engine.State.PendingAction == nil {
		t.Fatalf("cyclone scroll should ask for enemy equipment choice")
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{p1.Equipment[0].InstanceID},
	}}); err != nil {
		t.Fatalf("resolve cyclone scroll: %v", err)
	}
	if p1.Equipment[0] != nil || len(p1.Graveyard) != 1 {
		t.Fatalf("cyclone scroll should destroy low-cost enemy equipment, equipment=%v graveyard=%d", p1.Equipment[0], len(p1.Graveyard))
	}
}

func TestFireDanceSkirtUltimatePurifiesFriendlyFireUnit(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	skirt := NewCardInstance(baseCard(t, "2121007"), 0, 1)
	p0.Equipment[0] = skirt
	unit := placeUnit(baseCard(t, "1121001"), 0, 0, 0, engine)
	unit.Statuses[StatusBurn] = 2
	unit.Statuses[StatusFreeze] = 1
	unit.Statuses[StatusStun] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  skirt.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use fire dance skirt ultimate: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction || engine.State.PendingAction == nil {
		t.Fatalf("fire dance skirt should ask for a friendly fire unit")
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{unit.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve fire dance skirt: %v", err)
	}
	if hasAnyNegativeStatus(unit) {
		t.Fatalf("fire dance skirt should clear negative statuses, statuses=%v", unit.Statuses)
	}
}

func TestShiningCrystalAddsStunToFriendlyLightSpells(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	p0.Equipment[0] = NewCardInstance(baseCard(t, "2521010"), 0, 1)
	p0.Skills[0] = readySkill(baseCard(t, "3521006"), 0)
	p0.Elements[model.ElementLight] = 4

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast light spell: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve light spell: %v", err)
	}
	if target.Statuses[StatusStun] != 1 {
		t.Fatalf("shining crystal should add stun 1, statuses=%v", target.Statuses)
	}
}

func TestFlashRuneStunsEnemyFrontRowWhenEnemyCastsSkill(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	p0.Equipment[0] = NewCardInstance(baseCard(t, "2521011"), 0, 1)
	placeUnit(baseCard(t, "1021007"), 0, 1, 1, engine)
	front := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	back := placeUnit(baseCard(t, "1021004"), 1, 0, 1, engine)
	p1.Skills[0] = readySkill(baseCard(t, "3021005"), 1)
	p1.Elements[model.ElementArcane] = 2
	engine.State.CurrentTurn = 1

	if err := engine.HandleAction(1, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p1.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(1),
	}}); err != nil {
		t.Fatalf("enemy cast skill: %v", err)
	}
	if front.Statuses[StatusStun] != 1 || back.Statuses[StatusStun] != 0 {
		t.Fatalf("flash rune should stun only enemy front row, front=%v back=%v", front.Statuses, back.Statuses)
	}
}

func TestFireBarrierBoostsFireSpellsAndAddsBurn(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3121008"), 0)
	p0.Skills[1] = readySkill(baseCard(t, "3121001"), 0)
	p0.Elements[model.ElementFire] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[1].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast fireball with fire barrier: %v", err)
	}
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 5 {
		t.Fatalf("fire barrier should boost fireball power from 3 to 5, pending=%+v", engine.State.PendingSpell)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve fireball: %v", err)
	}
	if target.Statuses[StatusBurn] != 1 {
		t.Fatalf("fire barrier should add burn 1 on hit, statuses=%v", target.Statuses)
	}
}

func TestRuntimeLoadBonusCards(t *testing.T) {
	t.Run("blue crystal lamp pays light and gains light load", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		lamp := NewCardInstance(baseCard(t, "2521007"), 0, 1)
		p0.Equipment[0] = lamp
		p0.Elements[model.ElementLight] = 5

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  lamp.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use blue crystal lamp ultimate: %v", err)
		}
		if p0.Elements[model.ElementLight] != 0 || effectiveElementsGain(lamp)[model.ElementLight] != 3 {
			t.Fatalf("lamp should pay 5 light and increase load to 3 light, elements=%v gains=%v", p0.Elements, effectiveElementsGain(lamp))
		}
	})

	t.Run("necromancy stone gains shadow load when friendly companion dies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		stone := NewCardInstance(baseCard(t, "2611001"), 0, 1)
		p0.Equipment[0] = stone
		unit := placeUnit(baseCard(t, "1021007"), 0, 0, 0, engine)

		engine.destroyUnit(unit, 0)
		if effectiveElementsGain(stone)[model.ElementShadow] != 1 {
			t.Fatalf("necromancy stone should gain shadow load, gains=%v", effectiveElementsGain(stone))
		}
	})

	t.Run("soul necklace gains shadow when friendly companion dies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		necklace := NewCardInstance(baseCard(t, "2621006"), 0, 1)
		p0.Equipment[0] = necklace
		unit := placeUnit(baseCard(t, "1021007"), 0, 0, 0, engine)

		engine.destroyUnit(unit, 0)
		if p0.Elements[model.ElementShadow] != 1 {
			t.Fatalf("soul necklace should gain 1 shadow, elements=%v", p0.Elements)
		}
	})
}

func TestItemSpellScrollEffects(t *testing.T) {
	t.Run("waterform bind consumes an enemy companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021007"), 1, 1, 0, engine)
		scroll := NewCardInstance(baseCard(t, "2221008"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementWater] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use waterform bind: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve waterform bind: %v", err)
		}
		if !target.IsHorizontal {
			t.Fatalf("waterform bind should consume target")
		}
	})

	t.Run("thunderstorm damages and stuns hit companions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		a := placeUnit(baseCard(t, "1021007"), 1, 0, 0, engine)
		b := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		scroll := NewCardInstance(baseCard(t, "2321003"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementAir] = 4

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use thunderstorm: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{a.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve thunderstorm: %v", err)
		}
		if a.Statuses[StatusStun] != 1 || b.Statuses[StatusStun] != 1 {
			t.Fatalf("thunderstorm should stun all hit companions, a=%v b=%v", a.Statuses, b.Statuses)
		}
	})

	t.Run("chain lightning damages and draws", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021007"), 1, 1, 0, engine)
		target.CurrentLife = 3
		scroll := NewCardInstance(baseCard(t, "2321009"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
		p0.Elements[model.ElementAir] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use chain lightning: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve chain lightning: %v", err)
		}
		if target.CurrentLife != 2 || len(p0.Hand) != 1 {
			t.Fatalf("chain lightning should damage and draw, life=%d hand=%d", target.CurrentLife, len(p0.Hand))
		}
	})
}

func TestLightSearchSpells(t *testing.T) {
	t.Run("united hope searches a light companion among top five", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3501001"), 0)
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1521005"), 0, 1),
			NewCardInstance(baseCard(t, "1021004"), 0, 1),
		}
		p0.Elements[model.ElementLight] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast united hope: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0].Card.Number != "1521005" {
			t.Fatalf("united hope should find light companion, hand=%v", p0.Hand)
		}
	})

	t.Run("call of hope searches first light companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3521007"), 0)
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1021004"), 0, 1),
			NewCardInstance(baseCard(t, "1521005"), 0, 1),
		}
		p0.Elements[model.ElementLight] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast call of hope: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0].Card.Number != "1521005" {
			t.Fatalf("call of hope should find first light companion, hand=%v", p0.Hand)
		}
	})
}

func TestDelayedAndTurnCountEffects(t *testing.T) {
	t.Run("focus scroll gains arcane next turn start", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2021021"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementArcane] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use focus scroll: %v", err)
		}
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModDelayedElementGain {
			t.Fatalf("focus scroll should add delayed element modifier, modifiers=%v", p0.TempModifiers)
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end p0 turn: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end p1 turn: %v", err)
		}
		if p0.Elements[model.ElementArcane] != 3 || len(p0.TempModifiers) != 0 {
			t.Fatalf("focus scroll should gain 3 arcane next turn, elements=%v modifiers=%v", p0.Elements, p0.TempModifiers)
		}
	})

	t.Run("energy potion modifier expires when opponent turn ends", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3021005"), 0)
		p0.Skills[0].IsHorizontal = true
		potion := NewCardInstance(baseCard(t, "2221005"), 0, 1)
		p0.Hand = []*CardInstance{potion}
		p0.Elements[model.ElementWater] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": potion.InstanceID,
		}}); err != nil {
			t.Fatalf("use energy potion: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end p0 turn: %v", err)
		}
		if p0.Skills[0].IsHorizontal {
			t.Fatalf("skills should reset at own turn end under current turn rules")
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end p1 turn: %v", err)
		}
		if p0.Skills[0].IsHorizontal || len(p0.TempModifiers) != 0 {
			t.Fatalf("energy potion should reset skills after opponent turn, horizontal=%v modifiers=%v", p0.Skills[0].IsHorizontal, p0.TempModifiers)
		}
	})

	t.Run("flame rekindle gains fire for fire spells cast this turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3121014"), 0)
		p0.Elements[model.ElementFire] = 3
		placeUnit(baseCard(t, "1021007"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve fireball: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[1].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast flame rekindle: %v", err)
		}
		if p0.Elements[model.ElementFire] != 3 {
			t.Fatalf("flame rekindle should pay 1 and gain 2 fire for two fire casts, elements=%v", p0.Elements)
		}
	})
}

func TestBagOfTricksSacrificesToSearchConsumable(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	bag := NewCardInstance(baseCard(t, "2021006"), 0, 1)
	p0.Equipment[0] = bag
	target := NewCardInstance(baseCard(t, "2021021"), 0, 1)
	p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1), target}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  bag.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use bag of tricks ultimate: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction || engine.State.PendingAction == nil {
		t.Fatalf("bag of tricks should ask which consumable to search")
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{target.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve bag of tricks: %v", err)
	}
	if p0.Equipment[0] != nil || len(p0.Graveyard) != 1 || len(p0.Hand) != 1 || p0.Hand[0] != target {
		t.Fatalf("bag should sacrifice and search target, equipment=%v grave=%d hand=%d", p0.Equipment[0], len(p0.Graveyard), len(p0.Hand))
	}
}

func TestUtilityScrollAndForesightEffects(t *testing.T) {
	t.Run("cursed scroll draws two and discards them at turn end", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2021019"), 0, 1)
		draw1 := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		draw2 := NewCardInstance(baseCard(t, "1021004"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Deck = []*CardInstance{draw1, draw2}
		p0.Elements[model.ElementArcane] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use cursed scroll: %v", err)
		}
		if len(p0.Hand) != 2 || len(p0.DiscardAtTurnEnd) != 2 {
			t.Fatalf("cursed scroll should draw and mark cards, hand=%d markers=%v", len(p0.Hand), p0.DiscardAtTurnEnd)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end turn: %v", err)
		}
		if len(p0.Hand) != 0 || len(p0.Graveyard) != 3 {
			t.Fatalf("cursed scroll should discard drawn cards at turn end, hand=%d graveyard=%d", len(p0.Hand), len(p0.Graveyard))
		}
	})

	t.Run("foresight can move selected top deck cards to bottom", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3021002"), 0)
		a := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		b := NewCardInstance(baseCard(t, "1021004"), 0, 1)
		c := NewCardInstance(baseCard(t, "1021007"), 0, 1)
		p0.Deck = []*CardInstance{a, b, c}
		p0.Elements[model.ElementArcane] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast foresight: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{a.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve foresight: %v", err)
		}
		if p0.Deck[0] != b || p0.Deck[2] != a {
			t.Fatalf("foresight should move selected card to bottom, deck=%s,%s,%s", p0.Deck[0].Card.Number, p0.Deck[1].Card.Number, p0.Deck[2].Card.Number)
		}
	})

	t.Run("deep frost curse permanently freezes enemy companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021007"), 1, 1, 0, engine)
		scroll := NewCardInstance(baseCard(t, "2221013"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementWater] = 3

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use deep frost curse: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve deep frost curse: %v", err)
		}
		if target.Statuses[StatusFreeze] < 50 {
			t.Fatalf("deep frost curse should apply long freeze, statuses=%v", target.Statuses)
		}
	})

	t.Run("soul devour weakens selected enemy skills", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		p1.Skills[0] = readySkill(baseCard(t, "3021005"), 1)
		p1.Skills[1] = readySkill(baseCard(t, "3121001"), 1)
		scroll := NewCardInstance(baseCard(t, "2621008"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementShadow] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use soul devour: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{p1.Skills[0].InstanceID, p1.Skills[1].InstanceID},
		}}); err != nil {
			t.Fatalf("resolve soul devour: %v", err)
		}
		if p1.Skills[0].Statuses[StatusWeaken] != 1 || p1.Skills[1].Statuses[StatusWeaken] != 1 {
			t.Fatalf("soul devour should weaken selected skills, s0=%v s1=%v", p1.Skills[0].Statuses, p1.Skills[1].Statuses)
		}
	})
}

func TestRebirthScrollRevivesLightCompanion(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	dead := NewCardInstance(baseCard(t, "1521005"), 0, 1)
	dead.CurrentLife = 0
	p0.Graveyard = []*CardInstance{dead}
	scroll := NewCardInstance(baseCard(t, "2521005"), 0, 1)
	p0.Hand = []*CardInstance{scroll}
	p0.Elements[model.ElementLight] = 10

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": scroll.InstanceID,
	}}); err != nil {
		t.Fatalf("use rebirth scroll: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{dead.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve rebirth scroll: %v", err)
	}
	if len(p0.Graveyard) != 1 || dead.Position == nil || p0.Units[dead.Position.Col][dead.Position.Row] != dead {
		t.Fatalf("rebirth should revive companion onto board, graveyard=%d position=%v", len(p0.Graveyard), dead.Position)
	}
	if dead.CurrentLife != dead.Card.Life {
		t.Fatalf("rebirth should restore life, life=%d", dead.CurrentLife)
	}
}

func TestSkillContributionModifiers(t *testing.T) {
	t.Run("water and frost blade boost spell power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3221002"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3221003"), 0)
		p0.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
			"boost_ids":   []any{p0.Skills[1].InstanceID},
		}}); err != nil {
			t.Fatalf("cast boosted water spell: %v", err)
		}
		want := baseCard(t, "3221002").Power + baseCard(t, "3221003").Power + 2
		if engine.State.PendingSpell.TotalPower != want {
			t.Fatalf("freezing current should add +2 power while boosting water, got %d want %d", engine.State.PendingSpell.TotalPower, want)
		}
	})

	t.Run("frost blade adds power while attacking", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3221009"), 0)
		p0.Elements[model.ElementWater] = 10
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast frost blade: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != baseCard(t, "3221009").Power+2 {
			t.Fatalf("frost blade should add +2 power while attacking, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("joint casting boosts spell damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3421008"), 0)
		p0.Elements[model.ElementFire] = 10
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
			"boost_ids":   []any{p0.Skills[1].InstanceID},
		}}); err != nil {
			t.Fatalf("cast spell with joint casting: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve joint casting hit: %v", err)
		}
		if target.CurrentLife != target.Card.Life-2 {
			t.Fatalf("joint casting should add +1 attack damage, life=%d", target.CurrentLife)
		}
	})

	t.Run("moonlight can defend above its printed power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3101002"), 0)
		p1.Skills[0] = readySkill(baseCard(t, "3521013"), 1)
		p0.Elements[model.ElementFire] = 10
		p1.Elements[model.ElementLight] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast all fires as one: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != 5 {
			t.Fatalf("expected 5 attack power, got %+v", engine.State.PendingSpell)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{p1.Skills[0].InstanceID},
		}}); err != nil {
			t.Fatalf("defend with moonlight: %v", err)
		}
		if target.CurrentLife != target.Card.Life {
			t.Fatalf("moonlight +2 defense power should stop the spell, life=%d", target.CurrentLife)
		}
	})

	t.Run("all fires as one gains attack from power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3101002"), 0)
		p0.Elements[model.ElementFire] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast all fires as one: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve all fires as one: %v", err)
		}
		if target.CurrentLife != target.Card.Life-2 {
			t.Fatalf("all fires as one should deal 2 at 5 power, life=%d", target.CurrentLife)
		}
	})
}

func TestSplashBlizzardAndSoulBiteEffects(t *testing.T) {
	t.Run("splash range damages and freezes adjacent units", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		center := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		left := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		down := placeUnit(baseCard(t, "1021001"), 1, 1, 1, engine)
		diagonal := placeUnit(baseCard(t, "1021001"), 1, 0, 1, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3221005"), 0)
		p0.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast ice field: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve ice field: %v", err)
		}
		for _, unit := range []*CardInstance{center, left, down} {
			if unit.CurrentLife != unit.Card.Life-1 || unit.Statuses[StatusFreeze] != 1 {
				t.Fatalf("splash unit should be damaged and frozen, unit=%s life=%d statuses=%v", unit.Card.Name, unit.CurrentLife, unit.Statuses)
			}
		}
		if diagonal.CurrentLife != diagonal.Card.Life || diagonal.Statuses[StatusFreeze] != 0 {
			t.Fatalf("diagonal unit should not be hit by splash, life=%d statuses=%v", diagonal.CurrentLife, diagonal.Statuses)
		}
	})

	t.Run("blizzard gives water and air spells power and freeze on hit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3221002"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3221015"), 0)
		p0.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast water spell with blizzard: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != baseCard(t, "3221002").Power+1 {
			t.Fatalf("blizzard should add +1 power to water spells, pending=%+v", engine.State.PendingSpell)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve water spell with blizzard: %v", err)
		}
		if target.Statuses[StatusFreeze] != 1 {
			t.Fatalf("blizzard should freeze hit target, statuses=%v", target.Statuses)
		}
	})

	t.Run("dead soul bite weakens enemy skills and weaken lowers power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3621006"), 0)
		p1.Skills[0] = readySkill(baseCard(t, "3121001"), 1)
		p0.Elements[model.ElementShadow] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast dead soul bite: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve dead soul bite: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{p1.Skills[0].InstanceID},
		}}); err != nil {
			t.Fatalf("resolve dead soul bite weaken: %v", err)
		}
		if p1.Skills[0].Statuses[StatusWeaken] != 1 {
			t.Fatalf("dead soul bite should add weaken, statuses=%v", p1.Skills[0].Statuses)
		}
		if power := engine.effectiveSpellPower(1, p1.Skills[0], nil); power != baseCard(t, "3121001").Power-1 {
			t.Fatalf("weaken should lower spell power by 1, got %d", power)
		}
	})
}

func TestStormEarthAndDeadFuryEffects(t *testing.T) {
	t.Run("storm fury adds air spell power for each card in hand", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3301001"), 0)
		p0.Hand = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1021002"), 0, 1),
		}
		p0.Elements[model.ElementAir] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast air spell with storm fury: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != baseCard(t, "3321005").Power+2 {
			t.Fatalf("storm fury should add hand size to air spell power, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("earth resonance hits all enemies and gains attack from large allies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		front := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		back := placeUnit(baseCard(t, "1021001"), 1, 0, 1, engine)
		placeUnit(baseCard(t, "1421001"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1521002"), 0, 2, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3421013"), 0)
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast earth resonance: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve earth resonance: %v", err)
		}
		if front.CurrentLife != front.Card.Life-2 || back.CurrentLife != back.Card.Life-2 {
			t.Fatalf("earth resonance should deal 2 to all enemies, front=%d back=%d", front.CurrentLife, back.CurrentLife)
		}
	})

	t.Run("dead fury gains permanent instance power when companions die", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3621008"), 0)
		dead := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		engine.destroyUnit(dead, 1)
		if p0.Skills[0].PowerBonus != 1 {
			t.Fatalf("dead fury should gain +1 instance power, bonus=%d", p0.Skills[0].PowerBonus)
		}
		if power := engine.effectiveSpellPower(0, p0.Skills[0], nil); power != baseCard(t, "3621008").Power+1 {
			t.Fatalf("dead fury power should include death bonus, got %d", power)
		}
	})
}

func TestBurnSpellGetsPowerAgainstBurningTarget(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	target.Statuses[StatusBurn] = 1
	p0.Skills[0] = readySkill(baseCard(t, "3121002"), 0)
	p0.Elements[model.ElementFire] = 10

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast burn against burning target: %v", err)
	}
	if engine.State.PendingSpell.TotalPower != baseCard(t, "3121002").Power+2 {
		t.Fatalf("burn should gain +2 power against burning target, pending=%+v", engine.State.PendingSpell)
	}
}

func TestDefenseAndPositionSkillEffects(t *testing.T) {
	t.Run("rivers to sea gains water after successful defense", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p1.Skills[0] = readySkill(baseCard(t, "3201001"), 1)
		p0.Elements[model.ElementFire] = 10
		p1.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		before := p1.Elements[model.ElementWater]
		if err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{p1.Skills[0].InstanceID},
		}}); err != nil {
			t.Fatalf("defend with rivers to sea: %v", err)
		}
		if p1.Elements[model.ElementWater] != before-baseCard(t, "3201001").ElementsExpense[model.ElementWater]+1 {
			t.Fatalf("rivers to sea should gain water equal attack spell attack, elements=%v before=%d", p1.Elements, before)
		}
	})

	t.Run("static barrier stuns attacker front row after failed defense", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		attackerFront := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3421013"), 0)
		p1.Skills[0] = readySkill(baseCard(t, "3321015"), 1)
		p0.Elements[model.ElementEarth] = 10
		p1.Elements[model.ElementAir] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast earth resonance: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{p1.Skills[0].InstanceID},
		}}); err != nil {
			t.Fatalf("defend with static barrier: %v", err)
		}
		if attackerFront.Statuses[StatusStun] != 1 {
			t.Fatalf("static barrier should stun attacker front row, statuses=%v", attackerFront.Statuses)
		}
		if target.CurrentLife == target.Card.Life {
			t.Fatalf("failed defense should still let the spell hit")
		}
	})

	t.Run("sky sense adds power when affected area includes back row", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		front := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		back := placeUnit(baseCard(t, "1021001"), 1, 1, 1, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3421013"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3321012"), 0)
		p0.Elements[model.ElementEarth] = 10
		p0.Elements[model.ElementAir] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast earth resonance with sky sense: %v", err)
		}
		want := baseCard(t, "3421013").Power + 2
		if engine.State.PendingSpell.TotalPower != want {
			t.Fatalf("sky sense should add +2 power when area includes back row, got %d want %d front=%v back=%v", engine.State.PendingSpell.TotalPower, want, front.Position, back.Position)
		}
	})

	t.Run("earthshaker gains power from friendly earth load", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1421008"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1421010"), 0, 2, 0, engine)
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3421012"), 0)
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast earthshaker: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != baseCard(t, "3421012").Power+2 {
			t.Fatalf("earthshaker should add friendly earth load to power, pending=%+v", engine.State.PendingSpell)
		}
	})
}

func TestSelectionSorcerySkillEffects(t *testing.T) {
	t.Run("call lightning discards then stuns selected enemy", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		discard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Hand = []*CardInstance{discard}
		p0.Skills[0] = readySkill(baseCard(t, "3321014"), 0)
		p0.Elements[model.ElementAir] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast call lightning: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{discard.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve call lightning discard: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve call lightning stun: %v", err)
		}
		if len(p0.Graveyard) != 1 || target.Statuses[StatusStun] != 1 {
			t.Fatalf("call lightning should discard and stun, graveyard=%d statuses=%v", len(p0.Graveyard), target.Statuses)
		}
	})

	t.Run("natural growth gives selected earth companion load at turn end", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1421010"), 0, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3421011"), 0)
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast natural growth: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve natural growth: %v", err)
		}
		if effectiveElementsGain(target)[model.ElementEarth] != 1 {
			t.Fatalf("natural growth should not apply before end turn, load=%v", effectiveElementsGain(target))
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end turn after natural growth: %v", err)
		}
		if effectiveElementsGain(target)[model.ElementEarth] != 2 {
			t.Fatalf("natural growth should add earth load at end turn, load=%v", effectiveElementsGain(target))
		}
	})

	t.Run("natural growth filters load four targets and keeps pending action recoverable", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		valid := placeUnit(baseCard(t, "1421010"), 0, 0, 0, engine)
		invalid := placeUnit(baseCard(t, "1421010"), 0, 2, 0, engine)
		addElementsGainBonus(invalid, model.ElementEarth, 3)
		p0.Skills[0] = readySkill(baseCard(t, "3421011"), 0)
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast natural growth: %v", err)
		}
		if engine.State.PendingAction == nil {
			t.Fatalf("natural growth should create pending action")
		}
		for _, candidate := range engine.State.PendingAction.Candidates {
			if candidate["instance_id"] == invalid.InstanceID {
				t.Fatalf("load four target should not be selectable, candidates=%v", engine.State.PendingAction.Candidates)
			}
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{invalid.InstanceID},
		}}); err == nil {
			t.Fatalf("invalid target should be rejected")
		}
		if engine.State.PendingAction == nil || engine.State.Phase != PhaseWaitingAction {
			t.Fatalf("invalid selection should keep pending action recoverable, phase=%v pending=%+v", engine.State.Phase, engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{valid.InstanceID},
		}}); err != nil {
			t.Fatalf("valid target should still resolve after invalid attempt: %v", err)
		}
	})

	t.Run("blood demon blast sacrifices companion and damages enemy front row", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		sacrifice := placeUnit(baseCard(t, "1521002"), 0, 1, 0, engine)
		frontA := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		frontB := placeUnit(baseCard(t, "1021001"), 1, 2, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3621010"), 0)
		p0.Elements[model.ElementShadow] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast blood demon blast: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{sacrifice.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve blood demon blast: %v", err)
		}
		if p0.Units[1][0] != nil || len(p0.Graveyard) == 0 {
			t.Fatalf("blood demon blast should sacrifice the selected companion")
		}
		if frontA.CurrentLife != frontA.Card.Life-sacrifice.Card.Life || frontB.CurrentLife != frontB.Card.Life-sacrifice.Card.Life {
			t.Fatalf("blood demon blast should damage enemy front row, a=%d b=%d", frontA.CurrentLife, frontB.CurrentLife)
		}
	})
}

func TestRemainingPassiveSkillEffects(t *testing.T) {
	t.Run("rapid sandstorm weakens low original power spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3421015"), 0)
		p0.Elements[model.ElementFire] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball under rapid sandstorm: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != 1 {
			t.Fatalf("rapid sandstorm should reduce fireball power to 1, pending=%+v", engine.State.PendingSpell)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve fireball under rapid sandstorm: %v", err)
		}
		if target.CurrentLife != target.Card.Life {
			t.Fatalf("rapid sandstorm should reduce fireball damage to 0, life=%d", target.CurrentLife)
		}
	})

	t.Run("undead wall gains defense power after friendly death", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		dead := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3101002"), 0)
		p1.Skills[0] = readySkill(baseCard(t, "3621013"), 1)
		p0.Elements[model.ElementFire] = 10
		p1.Elements[model.ElementShadow] = 10

		engine.destroyUnit(dead, 1)
		if p1.Skills[0].Statuses[recentFriendlyDeathStatus] <= 0 {
			t.Fatalf("undead wall should remember recent friendly death, statuses=%v", p1.Skills[0].Statuses)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast all fires as one: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{p1.Skills[0].InstanceID},
		}}); err != nil {
			t.Fatalf("defend with undead wall: %v", err)
		}
		if target.CurrentLife != target.Card.Life {
			t.Fatalf("undead wall should defend successfully with recent death bonus, life=%d", target.CurrentLife)
		}
	})
}

func TestChoiceUtilitySkillEffects(t *testing.T) {
	t.Run("elemental enchant lets next spell apply chosen status", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3021007"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementArcane] = 10
		p0.Elements[model.ElementFire] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast elemental enchant: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{StatusBurn},
		}}); err != nil {
			t.Fatalf("choose elemental enchant status: %v", err)
		}
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Status != StatusBurn {
			t.Fatalf("elemental enchant should add next spell status modifier, modifiers=%v", p0.TempModifiers)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[1].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast enchanted fireball: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve enchanted fireball: %v", err)
		}
		if target.Statuses[StatusBurn] != 1 {
			t.Fatalf("enchanted fireball should burn target, statuses=%v", target.Statuses)
		}
	})

	t.Run("mind tempering permanently buffs selected skill", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		targetSkill := readySkill(baseCard(t, "3121001"), 0)
		p0.Skills[0] = readySkill(baseCard(t, "3021012"), 0)
		p0.Skills[1] = targetSkill
		p0.Elements[model.ElementArcane] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast mind tempering: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{targetSkill.InstanceID},
		}}); err != nil {
			t.Fatalf("choose mind tempering target: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{"power"},
		}}); err != nil {
			t.Fatalf("choose mind tempering mode: %v", err)
		}
		if targetSkill.PowerBonus != 3 {
			t.Fatalf("mind tempering should add +3 power, bonus=%d", targetSkill.PowerBonus)
		}
	})

	t.Run("water divination searches water card among top four", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		waterCard := NewCardInstance(baseCard(t, "1221001"), 0, 1)
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			waterCard,
			NewCardInstance(baseCard(t, "1021002"), 0, 1),
		}
		p0.Skills[0] = readySkill(baseCard(t, "3221007"), 0)
		p0.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast water divination: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{waterCard.InstanceID},
		}}); err != nil {
			t.Fatalf("choose water divination card: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != waterCard {
			t.Fatalf("water divination should search water card to hand, hand=%v", cardsToInfo(p0.Hand))
		}
	})
}

func TestCompanionUtilityEffects(t *testing.T) {
	t.Run("farolank gains adjacent companion load on enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		left := placeUnit(baseCard(t, "1421008"), 0, 0, 1, engine)
		right := placeUnit(baseCard(t, "1521002"), 0, 2, 1, engine)
		farolank := NewCardInstance(baseCard(t, "1011003"), 0, 1)
		p0.Hand = []*CardInstance{farolank}
		p0.Elements[model.ElementArcane] = 20

		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": farolank.InstanceID,
			"col":         float64(1),
			"row":         float64(1),
		}}); err != nil {
			t.Fatalf("summon farolank: %v", err)
		}
		load := effectiveElementsGain(farolank)
		if load[model.ElementEarth] != effectiveElementsGain(left)[model.ElementEarth] ||
			load[model.ElementAir] != effectiveElementsGain(left)[model.ElementAir] ||
			load[model.ElementLight] != effectiveElementsGain(right)[model.ElementLight]+farolank.Card.ElementsGain[model.ElementLight] {
			t.Fatalf("farolank should gain adjacent companion load, load=%v", load)
		}
	})

	t.Run("specialist mage changes load to chosen element", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		mage := NewCardInstance(baseCard(t, "1021010"), 0, 1)
		p0.Hand = []*CardInstance{mage}
		p0.Elements[model.ElementArcane] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": mage.InstanceID,
			"col":         float64(1),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("summon specialist mage: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{model.ElementFire},
		}}); err != nil {
			t.Fatalf("choose specialist mage element: %v", err)
		}
		load := effectiveElementsGain(mage)
		if load[model.ElementFire] != 2 || len(load) != 1 {
			t.Fatalf("specialist mage should convert load to chosen element, load=%v", load)
		}
	})

	t.Run("blessed girl gives adjacent earth companion load", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		girl := placeUnit(baseCard(t, "1421009"), 0, 1, 1, engine)
		target := placeUnit(baseCard(t, "1421010"), 0, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  girl.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use blessed girl ability: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve blessed girl load: %v", err)
		}
		if effectiveElementsGain(target)[model.ElementEarth] != target.Card.ElementsGain[model.ElementEarth]+1 {
			t.Fatalf("blessed girl should add earth load, load=%v", effectiveElementsGain(target))
		}
	})

	t.Run("lundesal buffs a selected skill on enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		lundesal := NewCardInstance(baseCard(t, "1511002"), 0, 1)
		targetSkill := readySkill(baseCard(t, "3121001"), 0)
		p0.Hand = []*CardInstance{lundesal}
		p0.Skills[0] = targetSkill
		p0.Elements[model.ElementLight] = 20

		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": lundesal.InstanceID,
			"col":         float64(1),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("summon lundesal: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{targetSkill.InstanceID},
		}}); err != nil {
			t.Fatalf("choose lundesal skill: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{"attack"},
		}}); err != nil {
			t.Fatalf("choose lundesal buff: %v", err)
		}
		if targetSkill.AttackBonus != 1 {
			t.Fatalf("lundesal should add +1 attack, bonus=%d", targetSkill.AttackBonus)
		}
	})
}

func TestMasteryIsPerCardAndSettlesAsMark(t *testing.T) {
	t.Run("knowledge tree maxes friendly mastery cards without changing player charge", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		tree := placeUnit(baseCard(t, "1421003"), 0, 0, 0, engine)
		guard := placeUnit(baseCard(t, "1421004"), 0, 1, 0, engine)
		knowledgeTree := placeUnit(baseCard(t, "1411002"), 0, 2, 0, engine)

		engine.triggerEffects(TriggerOnEnter, knowledgeTree, nil, nil)

		if p0.Charge != 0 {
			t.Fatalf("精通 must not be stored as player charge, charge=%d", p0.Charge)
		}
		if tree.Statuses[StatusMastery] != 4 || guard.Statuses[StatusMastery] != 5 {
			t.Fatalf("knowledge tree should max per-card mastery, tree=%v guard=%v", tree.Statuses, guard.Statuses)
		}
		if effectiveElementsGain(tree)[model.ElementEarth] != tree.Card.ElementsGain[model.ElementEarth]+2 {
			t.Fatalf("growing treant should receive both mastery load bonuses, load=%v", effectiveElementsGain(tree))
		}
		if guard.CurrentLife != guard.Card.Life+1 || effectiveElementsGain(guard)[model.ElementEarth] != guard.Card.ElementsGain[model.ElementEarth]+1 || guard.AttackBonus != 2 {
			t.Fatalf("forest guard mastery bonuses wrong, life=%d load=%v attack_bonus=%d", guard.CurrentLife, effectiveElementsGain(guard), guard.AttackBonus)
		}
	})

	t.Run("consume advances mastery immediately and mark settlement does not", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		tree := placeUnit(baseCard(t, "1421003"), 0, 0, 0, engine)
		tree.IsHorizontal = false

		engine.processEndOfTurnStatuses(engine.State.Players[0])
		engine.settleMastery(engine.State.Players[0])
		if tree.Statuses[StatusMastery] != 0 || effectiveElementsGain(tree)[model.ElementEarth] != tree.Card.ElementsGain[model.ElementEarth] {
			t.Fatalf("mark settlement should not advance mastery, statuses=%v load=%v", tree.Statuses, effectiveElementsGain(tree))
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{
			"instance_id": tree.InstanceID,
		}}); err != nil {
			t.Fatalf("consume treant once: %v", err)
		}
		if tree.Statuses[StatusMastery] != 1 || effectiveElementsGain(tree)[model.ElementEarth] != tree.Card.ElementsGain[model.ElementEarth] {
			t.Fatalf("mastery 1 should advance on consume without threshold load, statuses=%v load=%v", tree.Statuses, effectiveElementsGain(tree))
		}
		tree.IsHorizontal = false
		if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{
			"instance_id": tree.InstanceID,
		}}); err != nil {
			t.Fatalf("consume treant twice: %v", err)
		}
		if tree.Statuses[StatusMastery] != 2 || effectiveElementsGain(tree)[model.ElementEarth] != tree.Card.ElementsGain[model.ElementEarth]+1 {
			t.Fatalf("mastery 2 should trigger exactly one treant load bonus on consume, statuses=%v load=%v", tree.Statuses, effectiveElementsGain(tree))
		}
	})

	t.Run("dragon prince searches water companion with play cost discount at mastery 2", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		dragon := placeUnit(baseCard(t, "1221012"), 0, 0, 0, engine)
		target := NewCardInstance(baseCard(t, "1221010"), 0, 1)
		p0.Deck = []*CardInstance{target}

		engine.advanceMastery(dragon, 0, 2)

		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "dragon_prince_search" {
			t.Fatalf("dragon prince mastery should open search action, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("choose dragon prince search: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != target || target.Statuses["入场费用水-1"] != 1 {
			t.Fatalf("searched water companion should enter hand with discount, hand=%v statuses=%v", cardsToInfo(p0.Hand), target.Statuses)
		}
		cost := engine.effectiveCardPlayCost(p0, target)
		if cost[model.ElementWater] != max(target.Card.ElementsCost[model.ElementWater]-1, 0) {
			t.Fatalf("discounted water play cost wrong, cost=%v base=%v", cost, target.Card.ElementsCost)
		}
	})

	t.Run("earth mastery skills read their own mastery markers", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		earthCrack := readySkill(baseCard(t, "3421003"), 0)
		forestShelter := readySkill(baseCard(t, "3421001"), 0)

		stats := engine.skillContributionStats(0, earthCrack, earthCrack, skillPurposeAttack)
		if stats.PowerBonus != 4 || stats.DamageBonus != 0 {
			t.Fatalf("earth crack without mastery should be base only, stats=%+v", stats)
		}
		earthCrack.Statuses[StatusMastery] = 3
		stats = engine.skillContributionStats(0, earthCrack, earthCrack, skillPurposeAttack)
		if stats.PowerBonus != 6 || stats.DamageBonus != 2 {
			t.Fatalf("earth crack mastery 3 should add +2 power/+2 damage, stats=%+v", stats)
		}

		stats = engine.skillContributionStats(0, forestShelter, nil, skillPurposeDefend)
		if stats.PowerBonus != 2 {
			t.Fatalf("forest shelter without mastery should be base 2 power, stats=%+v", stats)
		}
		forestShelter.Statuses[StatusMastery] = 1
		stats = engine.skillContributionStats(0, forestShelter, nil, skillPurposeDefend)
		if stats.PowerBonus != 4 {
			t.Fatalf("forest shelter mastery 1 should become 4 power, stats=%+v", stats)
		}
		forestShelter.Statuses[StatusMastery] = 3
		stats = engine.skillContributionStats(0, forestShelter, nil, skillPurposeDefend)
		if stats.PowerBonus != 6 {
			t.Fatalf("forest shelter mastery 3 should become 6 power, stats=%+v", stats)
		}
	})

	t.Run("great elder mastery discounts only the next earth skill learn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		elder := placeUnit(baseCard(t, "1421011"), 0, 0, 0, engine)
		earthSkill := NewCardInstance(baseCard(t, "3421003"), 0, 1)
		p0.SkillPool = []*CardInstance{earthSkill}
		p0.Elements[model.ElementEarth] = 3

		engine.advanceMastery(elder, 0, 1)
		cost := engine.effectiveSkillLearnCost(p0, earthSkill)
		if cost[model.ElementEarth] != 1 {
			t.Fatalf("great elder mastery should reduce next earth skill learn by 2, cost=%v", cost)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "learn_skill", Data: map[string]any{
			"instance_id": earthSkill.InstanceID,
		}}); err != nil {
			t.Fatalf("learn discounted earth skill: %v", err)
		}
		if p0.Elements[model.ElementEarth] != 2 || len(p0.TempModifiers) != 0 {
			t.Fatalf("discount should be consumed after learning once, elements=%v modifiers=%v", p0.Elements, p0.TempModifiers)
		}
	})

	t.Run("parasitic touch gains load from its own mastery marker", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		item := NewCardInstance(baseCard(t, "2421007"), 0, 1)
		engine.State.Players[0].Equipment[0] = item

		if effectiveElementsGain(item)[model.ElementEarth] != item.Card.ElementsGain[model.ElementEarth] {
			t.Fatalf("parasitic touch should not gain mastery load before mark settlement")
		}
		engine.advanceMastery(item, 0, 1)
		if item.Statuses[StatusMastery] != 1 || effectiveElementsGain(item)[model.ElementEarth] != item.Card.ElementsGain[model.ElementEarth]+1 {
			t.Fatalf("parasitic touch mastery should add earth load, statuses=%v load=%v", item.Statuses, effectiveElementsGain(item))
		}
	})
}

func TestEnemyAndGlobalSpellDamageReducers(t *testing.T) {
	t.Run("wall keeper entry makes every player's spell damage zero", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		wallKeeper := placeUnit(baseCard(t, "1221010"), 0, 0, 0, engine)
		engine.triggerEffects(TriggerOnEnter, wallKeeper, nil, nil)

		p0Skill := readySkill(baseCard(t, "3021005"), 0)
		p1Skill := readySkill(baseCard(t, "3021005"), 1)

		if damage := engine.effectiveSpellDamage(0, p0Skill, 1, nil); damage != 0 {
			t.Fatalf("wall keeper should reduce friendly spell damage to 0 too, damage=%d", damage)
		}
		if damage := engine.effectiveSpellDamage(1, p1Skill, 1, nil); damage != 0 {
			t.Fatalf("wall keeper should reduce enemy spell damage to 0, damage=%d", damage)
		}
	})

	t.Run("frost heart optionally reacts and sacrifices itself to zero an enemy spell hit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		target.CurrentLife = 3
		frostHeart := NewCardInstance(baseCard(t, "2221001"), 1, 1)
		p1.Equipment[0] = frostHeart
		p0.Skills[0] = readySkill(baseCard(t, "3121002"), 0)
		p0.Elements[model.ElementFire] = 10
		if _, ok := globalRegistry.GetBehavior("2221001").(SpellReactionBehavior); !ok {
			t.Fatalf("frost heart should expose spell reaction behavior")
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast arcane arrow into frost heart: %v", err)
		}
		if target.CurrentLife != 3 || p1.Equipment[0] != frostHeart {
			t.Fatalf("frost heart should not auto-trigger before reaction, life=%d equipment=%v", target.CurrentLife, p1.Equipment[0])
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "react_spell", Data: map[string]any{
			"instance_id": frostHeart.InstanceID,
		}}); err != nil {
			t.Fatalf("react with frost heart: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve spell after frost heart: %v", err)
		}
		if target.CurrentLife != 3 {
			t.Fatalf("frost heart should zero the incoming spell damage, life=%d", target.CurrentLife)
		}
		if p1.Equipment[0] != nil || len(p1.Graveyard) != 1 || p1.Graveyard[0] != frostHeart {
			t.Fatalf("frost heart should sacrifice to graveyard, equipment=%v graveyard=%v", p1.Equipment[0], cardsToInfo(p1.Graveyard))
		}
	})

	t.Run("shadow cloak prevents only the first enemy spell hit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		target.CurrentLife = 3
		cloak := NewCardInstance(baseCard(t, "2621012"), 1, 1)
		p1.Equipment[0] = cloak
		p0.Skills[0] = readySkill(baseCard(t, "3021005"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3021005"), 0)
		p0.Elements[model.ElementArcane] = 4

		for i := 0; i < 2; i++ {
			if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
				"instance_id": p0.Skills[i].InstanceID,
				"target_type": "unit",
				"target_col":  float64(1),
				"target_row":  float64(0),
			}}); err != nil {
				t.Fatalf("cast arcane arrow %d into shadow cloak: %v", i+1, err)
			}
		}
		if target.CurrentLife != 2 || cloak.Statuses["已防护"] != 1 {
			t.Fatalf("shadow cloak should block first hit only, life=%d statuses=%v", target.CurrentLife, cloak.Statuses)
		}
	})
}

func TestReviewCardsDamagePreventionAndEarthLightActives(t *testing.T) {
	t.Run("arena phantom only takes spell damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		phantom := placeUnit(baseCard(t, "1021009"), 1, 1, 0, engine)
		attacker := placeUnit(baseCard(t, "1021013"), 0, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "attack", Data: map[string]any{
			"attacker_id": attacker.InstanceID,
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("attack arena phantom: %v", err)
		}
		if phantom.CurrentLife != phantom.Card.Life {
			t.Fatalf("arena phantom should ignore attack damage, life=%d", phantom.CurrentLife)
		}

		p0.Skills[0] = readySkill(baseCard(t, "3021005"), 0)
		p0.Elements[model.ElementArcane] = 2
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast arcane arrow at arena phantom: %v", err)
		}
		if phantom.CurrentLife != phantom.Card.Life-1 {
			t.Fatalf("arena phantom should take spell damage, life=%d", phantom.CurrentLife)
		}
	})

	t.Run("healing warlock heals a selected friendly unit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		warlock := placeUnit(baseCard(t, "1521001"), 0, 0, 0, engine)
		target := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		target.CurrentLife = target.Card.Life - 1

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  warlock.InstanceID,
			"ability_type": "per_turn",
			"target_id":    target.InstanceID,
		}}); err != nil {
			t.Fatalf("use healing warlock: %v", err)
		}
		if target.CurrentLife != target.Card.Life {
			t.Fatalf("healing warlock should heal target to max, life=%d", target.CurrentLife)
		}
	})

	t.Run("growth potion resets an earth companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1421003"), 0, 1, 0, engine)
		target.IsHorizontal = true
		target.Statuses[StatusCooldown] = 1
		potion := NewCardInstance(baseCard(t, "2421002"), 0, 1)
		p0.Hand = append(p0.Hand, potion)
		p0.Elements[model.ElementEarth] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": potion.InstanceID,
		}}); err != nil {
			t.Fatalf("use growth potion: %v", err)
		}
		if target.IsHorizontal || target.Statuses[StatusCooldown] != 0 {
			t.Fatalf("growth potion should reset earth companion, horizontal=%v statuses=%v", target.IsHorizontal, target.Statuses)
		}
	})

	t.Run("sturdy and nature seal scrolls reduce spell damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		sturdy := NewCardInstance(baseCard(t, "2421003"), 1, 1)
		p1.Hand = append(p1.Hand, sturdy)
		p1.Elements[model.ElementEarth] = 2
		engine.State.CurrentTurn = 1
		if err := engine.HandleAction(1, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": sturdy.InstanceID,
		}}); err != nil {
			t.Fatalf("use sturdy scroll: %v", err)
		}
		engine.State.CurrentTurn = 0

		damageSkill := baseCard(t, "3121002")
		p0.Skills[0] = readySkill(damageSkill, 0)
		p0.Elements[model.ElementFire] = 10
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball into sturdy scroll: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve fireball into sturdy scroll: %v", err)
		}
		wantLife := target.Card.Life - max(damageSkill.Attack-1, 0)
		if target.CurrentLife != wantLife {
			t.Fatalf("sturdy scroll should reduce spell damage by 1, life=%d want=%d", target.CurrentLife, wantLife)
		}

		engine = setupReportedBugEngine(t)
		p0 = engine.State.Players[0]
		p1 = engine.State.Players[1]
		target = placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		seal := NewCardInstance(baseCard(t, "2421010"), 1, 1)
		p1.Hand = append(p1.Hand, seal)
		p1.Elements[model.ElementEarth] = 2
		engine.State.CurrentTurn = 1
		if err := engine.HandleAction(1, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": seal.InstanceID,
		}}); err != nil {
			t.Fatalf("use nature seal scroll: %v", err)
		}
		engine.State.CurrentTurn = 0
		p0.Skills[0] = readySkill(damageSkill, 0)
		p0.Elements[model.ElementFire] = 10
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball into nature seal: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve fireball into nature seal: %v", err)
		}
		if target.CurrentLife != target.Card.Life {
			t.Fatalf("nature seal should reduce spell damage to zero, life=%d", target.CurrentLife)
		}
	})

	t.Run("elf armor per-turn ability heals hero", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Hero = placeUnit(baseCard(t, "4311003"), 0, 1, 1, engine)
		p0.Hero.CurrentLife = p0.Hero.Card.Life - 1
		armor := NewCardInstance(baseCard(t, "2421011"), 0, 1)
		p0.Equipment[0] = armor

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  armor.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use elf armor ability: %v", err)
		}
		if p0.Hero.CurrentLife != p0.Hero.Card.Life {
			t.Fatalf("elf armor should heal hero, life=%d", p0.Hero.CurrentLife)
		}
	})
}

func TestReviewCardsLoadReviveAndFriendlyTargetSpells(t *testing.T) {
	t.Run("mask changes hero load to same total arcane load", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		hero := placeUnit(baseCard(t, "4311003"), 0, 1, 1, engine)
		p0.Hero = hero
		setElementsGain(hero, map[string]int{model.ElementFire: 1, model.ElementWater: 2})
		mask := NewCardInstance(baseCard(t, "2021020"), 0, 1)
		p0.Hand = append(p0.Hand, mask)
		p0.Elements[model.ElementArcane] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
			"instance_id": mask.InstanceID,
		}}); err != nil {
			t.Fatalf("equip mask: %v", err)
		}
		load := effectiveElementsGain(hero)
		if load[model.ElementArcane] != 3 || load[model.ElementFire] != 0 || load[model.ElementWater] != 0 {
			t.Fatalf("mask should convert hero load to arcane with same total, load=%v", load)
		}
	})

	t.Run("mermaid tear removes itself and revives one companion at one life", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		tear := NewCardInstance(baseCard(t, "2211001"), 0, 1)
		p0.Equipment[0] = tear
		dead := NewCardInstance(baseCard(t, "1021013"), 0, 1)
		dead.CurrentLife = 0
		p0.Graveyard = []*CardInstance{dead}

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  tear.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use mermaid tear: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{dead.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve mermaid tear: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "mermaid_tear_position" {
			t.Fatalf("mermaid tear should ask for revive position, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{positionSelectionID(Position{Col: 2, Row: 2})},
		}}); err != nil {
			t.Fatalf("resolve mermaid tear position: %v", err)
		}
		if p0.Equipment[0] != nil || len(p0.Graveyard) != 0 || dead.CurrentLife != 1 || dead.Position == nil || *dead.Position != (Position{Col: 2, Row: 2}) {
			t.Fatalf("mermaid tear should leave play and revive at one life, equipment=%v grave=%d life=%d pos=%v", p0.Equipment[0], len(p0.Graveyard), dead.CurrentLife, dead.Position)
		}
	})

	t.Run("iridescent paint converts up to four light load to arcane", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		a := placeUnit(baseCard(t, "1521002"), 0, 0, 0, engine)
		b := placeUnit(baseCard(t, "1521005"), 0, 1, 0, engine)
		setElementsGain(a, map[string]int{model.ElementLight: 3})
		setElementsGain(b, map[string]int{model.ElementLight: 3})
		paint := NewCardInstance(baseCard(t, "2521012"), 0, 1)
		p0.Hand = append(p0.Hand, paint)
		p0.Elements[model.ElementLight] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": paint.InstanceID,
		}}); err != nil {
			t.Fatalf("use iridescent paint: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{a.InstanceID, b.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve iridescent paint: %v", err)
		}
		totalLight := effectiveElementsGain(a)[model.ElementLight] + effectiveElementsGain(b)[model.ElementLight]
		totalArcane := effectiveElementsGain(a)[model.ElementArcane] + effectiveElementsGain(b)[model.ElementArcane]
		if totalLight != 2 || totalArcane != 4 {
			t.Fatalf("iridescent paint should convert exactly four light load, light=%d arcane=%d", totalLight, totalArcane)
		}
	})

	t.Run("regeneration resets a selected earth companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1421003"), 0, 1, 0, engine)
		target.IsHorizontal = true
		target.Statuses[StatusCooldown] = 1
		p0.Skills[0] = readySkill(baseCard(t, "3421004"), 0)
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast regeneration: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve regeneration: %v", err)
		}
		if target.IsHorizontal || target.Statuses[StatusCooldown] != 0 {
			t.Fatalf("regeneration should reset earth companion, horizontal=%v statuses=%v", target.IsHorizontal, target.Statuses)
		}
	})

	t.Run("holy fire cleanses friendly target and damages enemy target", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		ally := placeUnit(baseCard(t, "1021013"), 0, 1, 0, engine)
		ally.Statuses[StatusBurn] = 1
		ally.Statuses[StatusStun] = 1
		p0.Skills[0] = readySkill(baseCard(t, "3521002"), 0)
		p0.Elements[model.ElementLight] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast holy fire on ally: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve holy fire on ally: %v", err)
		}
		if ally.CurrentLife != ally.Card.Life || ally.Statuses[StatusBurn] != 0 || ally.Statuses[StatusStun] != 0 {
			t.Fatalf("holy fire should cleanse friendly target without damage, life=%d statuses=%v", ally.CurrentLife, ally.Statuses)
		}

		engine = setupReportedBugEngine(t)
		p0 = engine.State.Players[0]
		enemy := placeUnit(baseCard(t, "1021013"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3521002"), 0)
		p0.Elements[model.ElementLight] = 10
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast holy fire on enemy: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve holy fire on enemy: %v", err)
		}
		if enemy.CurrentLife != enemy.Card.Life-1 {
			t.Fatalf("holy fire should damage enemy target, life=%d", enemy.CurrentLife)
		}
	})
}

func TestReviewCardsNegativeStatusImmunity(t *testing.T) {
	t.Run("divine guardian ignores negative status effects", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		guardian := placeUnit(baseCard(t, "1521010"), 0, 1, 0, engine)
		guardian.Statuses[StatusStun] = 1
		guardian.Statuses[StatusFreeze] = 1

		if !engine.canConsumeCard(guardian) {
			t.Fatalf("divine guardian should ignore stun for consume")
		}
		guardian.IsHorizontal = true
		engine.resetCards(engine.State.Players[0])
		if guardian.IsHorizontal {
			t.Fatalf("divine guardian should ignore freeze while resetting")
		}
	})

	t.Run("blessing priest protects itself and adjacent units while marks remain visible", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		priest := placeUnit(baseCard(t, "1421002"), 0, 1, 1, engine)
		adjacent := placeUnit(baseCard(t, "1021013"), 0, 1, 0, engine)
		far := placeUnit(baseCard(t, "1021013"), 0, 0, 0, engine)
		for _, unit := range []*CardInstance{priest, adjacent, far} {
			unit.Statuses[StatusStun] = 1
			unit.Statuses[StatusFreeze] = 1
		}

		if !engine.canConsumeCard(priest) || !engine.canConsumeCard(adjacent) {
			t.Fatalf("priest and adjacent unit should ignore stun")
		}
		if engine.canConsumeCard(far) {
			t.Fatalf("non-adjacent unit should still be affected by stun")
		}
		for _, unit := range []*CardInstance{priest, adjacent, far} {
			unit.IsHorizontal = true
		}
		engine.resetCards(engine.State.Players[0])
		if priest.IsHorizontal || adjacent.IsHorizontal {
			t.Fatalf("priest protection should let protected cards reset through freeze")
		}
		if !far.IsHorizontal {
			t.Fatalf("non-adjacent unit should stay frozen")
		}
		if priest.Statuses[StatusStun] != 1 || adjacent.Statuses[StatusFreeze] != 1 {
			t.Fatalf("protection should not erase visible marks, priest=%v adjacent=%v", priest.Statuses, adjacent.Statuses)
		}
	})
}

func TestWizardVolleyLineResetsSkillAndMakesNextCastFrontRow(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	frontA := placeUnit(baseCard(t, "1021013"), 1, 0, 0, engine)
	frontB := placeUnit(baseCard(t, "1021013"), 1, 1, 0, engine)
	back := placeUnit(baseCard(t, "1021013"), 1, 1, 1, engine)
	skill := readySkill(baseCard(t, "3021005"), 0)
	skill.IsHorizontal = true
	p0.Skills[0] = skill
	item := NewCardInstance(baseCard(t, "2021007"), 0, 1)
	p0.Hand = append(p0.Hand, item)
	p0.Elements[model.ElementArcane] = 10

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": item.InstanceID,
	}}); err != nil {
		t.Fatalf("use wizard volley line: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{skill.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve wizard volley line: %v", err)
	}
	if skill.IsHorizontal || skill.Statuses["下一次范围前排"] != 1 {
		t.Fatalf("wizard volley line should reset skill and mark next front-row area, horizontal=%v statuses=%v", skill.IsHorizontal, skill.Statuses)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": skill.InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast marked spell: %v", err)
	}
	if frontA.CurrentLife != frontA.Card.Life-1 || frontB.CurrentLife != frontB.Card.Life-1 || back.CurrentLife != back.Card.Life {
		t.Fatalf("marked spell should hit enemy front row only, frontA=%d frontB=%d back=%d", frontA.CurrentLife, frontB.CurrentLife, back.CurrentLife)
	}
	if skill.Statuses["下一次范围前排"] != 0 {
		t.Fatalf("front-row range mark should be consumed, statuses=%v", skill.Statuses)
	}
}

func TestPureArcaneLetsPlayerChooseElementAndAmountForNextMatchingSpell(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	placeUnit(baseCard(t, "1021013"), 1, 1, 0, engine)
	pure := readySkill(baseCard(t, "3001002"), 0)
	p0.Skills[0] = pure
	p0.Elements[model.ElementFire] = 4
	p0.Elements[model.ElementWater] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": pure.InstanceID,
		"target_type": "none",
	}}); err != nil {
		t.Fatalf("cast pure arcane: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "pure_arcane_spend" {
		t.Fatalf("pure arcane should ask for element/amount choice, pending=%+v", engine.State.PendingAction)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{model.ElementFire + ":3"},
	}}); err != nil {
		t.Fatalf("resolve pure arcane: %v", err)
	}
	if p0.Elements[model.ElementFire] != 1 || p0.Elements[model.ElementWater] != 2 {
		t.Fatalf("pure arcane should spend selected fire only, elements=%v", p0.Elements)
	}
	if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextElementSpellPowerBonus || p0.TempModifiers[0].Status != model.ElementFire || p0.TempModifiers[0].Amount != 3 {
		t.Fatalf("pure arcane should create next fire spell power modifier, modifiers=%v", p0.TempModifiers)
	}

	water := readySkill(baseCard(t, "3221009"), 0)
	p0.Skills[1] = water
	p0.Elements[model.ElementWater] = 10
	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": water.InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast water spell: %v", err)
	}
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != water.Card.Power+2 {
		t.Fatalf("pure arcane fire modifier should not boost water spell, pending=%+v", engine.State.PendingSpell)
	}
	if len(p0.TempModifiers) != 1 {
		t.Fatalf("fire modifier should remain after nonmatching spell, modifiers=%v", p0.TempModifiers)
	}
	engine.State.PendingSpell = nil
	engine.State.Phase = PhaseMain

	fire := readySkill(baseCard(t, "3121001"), 0)
	p0.Skills[2] = fire
	p0.Elements[model.ElementFire] = 10
	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": fire.InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast fire spell: %v", err)
	}
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != fire.Card.Power+3 {
		t.Fatalf("pure arcane should boost next matching fire spell by 3, pending=%+v", engine.State.PendingSpell)
	}
	if len(p0.TempModifiers) != 0 {
		t.Fatalf("pure arcane modifier should be consumed by matching spell, modifiers=%v", p0.TempModifiers)
	}
}

func TestHighRiskSkillReactionAndBoostSemantics(t *testing.T) {
	t.Run("3121015 焚风 grants pierce while boosting instead of adding damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		front := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		back := placeUnit(baseCard(t, "1021002"), 1, 1, 2, engine)
		main := readySkill(baseCard(t, "3321005"), 0)
		boost := readySkill(baseCard(t, "3121015"), 0)
		p0.Skills[0] = main
		p0.Skills[1] = boost
		for _, element := range model.AllElements {
			p0.Elements[element] = 10
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": main.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(2),
			"boost_ids":   []any{boost.InstanceID},
		}}); err != nil {
			t.Fatalf("burning wind boost should let spell pierce to back row: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.Target.Position != *back.Position {
			t.Fatalf("boosted spell should target back row, pending=%+v back=%v front=%v", engine.State.PendingSpell, back.Position, front.Position)
		}
		if got := engine.effectiveSpellDamage(0, main, main.Card.Attack, []*CardInstance{boost}); got != main.Card.Attack {
			t.Fatalf("burning wind should not add damage, got %d want %d", got, main.Card.Attack)
		}
	})

	t.Run("3321008 风洞 reacts to a non-area enemy spell and cancels it", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		attacker := readySkill(baseCard(t, "3121003"), 0)
		windHole := readySkill(baseCard(t, "3321008"), 1)
		p0.Skills[0] = attacker
		p1.Skills[0] = windHole
		for _, element := range model.AllElements {
			p0.Elements[element] = 10
			p1.Elements[element] = 10
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": attacker.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast spell for wind hole: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "react_spell", Data: map[string]any{
			"instance_id": windHole.InstanceID,
		}}); err != nil {
			t.Fatalf("wind hole should react: %v", err)
		}
		if engine.State.PendingSpell != nil || engine.State.Phase != PhaseMain {
			t.Fatalf("wind hole should cancel pending spell, pending=%+v phase=%s", engine.State.PendingSpell, engine.State.Phase)
		}
		if target.CurrentLife != target.Card.Life {
			t.Fatalf("cancelled spell should not damage target, life=%d", target.CurrentLife)
		}
		if !windHole.IsHorizontal || windHole.Statuses[StatusCooldown] != 1 {
			t.Fatalf("wind hole should be used and cooled down, horizontal=%v statuses=%v", windHole.IsHorizontal, windHole.Statuses)
		}
	})

	t.Run("3621015 虹吸 converts enemy spell damage into healing for affected units", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		target := placeUnit(baseCard(t, "1011002"), 1, 1, 0, engine)
		target.CurrentLife = 1
		attacker := readySkill(baseCard(t, "3121003"), 0)
		siphon := readySkill(baseCard(t, "3621015"), 1)
		p0.Skills[0] = attacker
		p1.Skills[0] = siphon
		for _, element := range model.AllElements {
			p0.Elements[element] = 10
			p1.Elements[element] = 10
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": attacker.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast spell for siphon: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "react_spell", Data: map[string]any{
			"instance_id": siphon.InstanceID,
		}}); err != nil {
			t.Fatalf("siphon should react: %v", err)
		}
		if engine.State.PendingSpell != nil || target.CurrentLife <= 1 {
			t.Fatalf("siphon should cancel spell and heal target, pending=%+v life=%d", engine.State.PendingSpell, target.CurrentLife)
		}
		if siphon.Statuses[StatusCooldown] != 2 {
			t.Fatalf("siphon should take cooldown 2, statuses=%v", siphon.Statuses)
		}
	})

	t.Run("3221010 水幻影 creates exactly one next-water-copy modifier on hit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		phantom := readySkill(baseCard(t, "3221010"), 0)
		p0.Skills[0] = phantom
		for _, element := range model.AllElements {
			p0.Elements[element] = 10
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": phantom.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast water phantom: %v", err)
		}
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != "next_water_copy" || p0.TempModifiers[0].RemainingUses != 1 {
			t.Fatalf("water phantom should arm one next water copy modifier, modifiers=%v", p0.TempModifiers)
		}
	})

	t.Run("3621007 安迪斯的惩罚 gains power from friendly damage only", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		punishment := readySkill(baseCard(t, "3621007"), 0)
		engine.State.Players[0].Skills[0] = punishment
		friend := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		enemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		engine.dealDamage(friend, 2, 0)
		if punishment.PowerBonus != 2 {
			t.Fatalf("friendly damage should increase andis punishment by 2, bonus=%d", punishment.PowerBonus)
		}
		engine.dealDamage(enemy, 2, 1)
		if punishment.PowerBonus != 2 {
			t.Fatalf("enemy damage should not increase andis punishment, bonus=%d", punishment.PowerBonus)
		}
	})
}

func TestHighRiskCompanionAndHeroSemantics(t *testing.T) {
	t.Run("1111001 火龙辉煌 requires fire devour and binds fire breath to itself", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		dragon := NewCardInstance(baseCard(t, "1111001"), 0, 1)
		food := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		addElementsGainBonus(food, model.ElementFire, 3)
		p0.Hand = []*CardInstance{dragon}
		for _, element := range model.AllElements {
			p0.Elements[element] = 10
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": dragon.InstanceID,
			"col":         float64(1),
			"row":         float64(1),
			"devour_id":   food.InstanceID,
		}}); err != nil {
			t.Fatalf("summon fire dragon with devour: %v", err)
		}
		if p0.Units[0][0] != nil || len(p0.Graveyard) != 1 || p0.Graveyard[0] != food {
			t.Fatalf("fire dragon should devour the selected companion, unit=%v graveyard=%v", p0.Units[0][0], cardsToInfo(p0.Graveyard))
		}
		if len(p0.SkillPool) != 0 {
			t.Fatalf("bound fire breath should not enter skill pool, pool=%v", cardsToInfo(p0.SkillPool))
		}
		if len(dragon.BoundSkills) != 1 || dragon.BoundSkills[0].Card.Number != "3101001" {
			t.Fatalf("fire breath should be bound to dragon, bound=%v", cardsToInfo(dragon.BoundSkills))
		}
	})

	t.Run("1621013 言灵 weakens only from the opponent third spell onward", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		wordSpirit := placeUnit(baseCard(t, "1621013"), 0, 1, 1, engine)
		first := readySkill(baseCard(t, "3121001"), 1)
		second := readySkill(baseCard(t, "3221009"), 1)
		third := readySkill(baseCard(t, "3321005"), 1)
		p1 := engine.State.Players[1]
		p1.SpellsCastThisTurn = map[string]int{model.ElementFire: 1}
		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, first, map[string]any{"cast_player": 1})
		if first.Statuses[StatusWeaken] != 0 {
			t.Fatalf("first enemy spell should not be weakened, source=%v first=%v", wordSpirit.Statuses, first.Statuses)
		}
		p1.SpellsCastThisTurn[model.ElementWater] = 1
		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, second, map[string]any{"cast_player": 1})
		if second.Statuses[StatusWeaken] != 0 {
			t.Fatalf("second enemy spell should not be weakened, second=%v", second.Statuses)
		}
		p1.SpellsCastThisTurn[model.ElementAir] = 1
		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, third, map[string]any{"cast_player": 1})
		if third.Statuses[StatusWeaken] != 1 {
			t.Fatalf("third enemy spell should be weakened, third=%v", third.Statuses)
		}
	})

	t.Run("hero passive cards participate in field triggers", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		noFace := NewCardInstance(baseCard(t, "4011002"), 0, 1)
		p0.Hero = noFace
		noFace.CurrentLife = noFace.Card.Life
		placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		entering := placeUnit(baseCard(t, "1021002"), 0, 1, 0, engine)

		engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, 0, entering, map[string]any{"entered_player": 0})

		if noFace.CurrentLife != noFace.Card.Life-1 {
			t.Fatalf("No Face should damage its hero on same-element friendly entry, life=%d", noFace.CurrentLife)
		}
	})

	t.Run("4111002 维兰德 converts one fire load to arcane for the turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		verland := NewCardInstance(baseCard(t, "4111002"), 0, 1)
		engine.State.Players[0].Hero = verland

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  verland.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use Verland per-turn: %v", err)
		}
		gain := effectiveElementsGain(verland)
		if verland.Statuses[StatusBurn] != 1 || gain[model.ElementFire] != 3 || gain[model.ElementArcane] != 1 {
			t.Fatalf("Verland should burn herself and convert fire load to arcane, statuses=%v gain=%v", verland.Statuses, gain)
		}
		engine.triggerEffects(TriggerOnTurnEnd, verland, nil, nil)
		gain = effectiveElementsGain(verland)
		if gain[model.ElementFire] != 4 || gain[model.ElementArcane] != 0 {
			t.Fatalf("Verland load conversion should end at turn end, gain=%v statuses=%v", gain, verland.Statuses)
		}
	})

	t.Run("startup heroes add their promised cards exactly once", func(t *testing.T) {
		cases := []struct {
			id       string
			location string
			card     string
			count    int
		}{
			{"4111001", "pool", "3101002", 1},
			{"4211002", "pool", "3201001", 1},
			{"4411002", "deck", "1401002", 1},
			{"4511002", "opponent_deck", "2501001", 5},
			{"4511003", "pool", "3501001", 1},
			{"4611003", "deck", "2601002", 3},
		}
		for _, tc := range cases {
			t.Run(tc.id, func(t *testing.T) {
				engine := setupReportedBugEngine(t)
				hero := NewCardInstance(baseCard(t, tc.id), 0, 1)
				engine.State.Players[0].Hero = hero
				engine.triggerEffects(TriggerOnTurnStart, hero, nil, nil)
				engine.triggerEffects(TriggerOnTurnStart, hero, nil, nil)

				switch tc.location {
				case "pool":
					if got := countCardNumber(engine.State.Players[0].SkillPool, tc.card); got != tc.count {
						t.Fatalf("skill pool count for %s got %d want %d", tc.card, got, tc.count)
					}
				case "deck":
					if got := countCardNumber(engine.State.Players[0].Deck, tc.card); got != tc.count {
						t.Fatalf("deck count for %s got %d want %d", tc.card, got, tc.count)
					}
				case "opponent_deck":
					if got := countCardNumber(engine.State.Players[1].Deck, tc.card); got != tc.count {
						t.Fatalf("opponent deck count for %s got %d want %d", tc.card, got, tc.count)
					}
				}
			})
		}
	})

	t.Run("4011001 斯卡尔蒂 discards for elements and locks non-arcane elements", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		skadi := NewCardInstance(baseCard(t, "4011001"), 0, 1)
		p0.Hero = skadi
		fire := NewCardInstance(baseCard(t, "3121001"), 0, 1)
		arcane := NewCardInstance(baseCard(t, "3021001"), 0, 1)
		p0.Hand = []*CardInstance{fire, arcane}

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  skadi.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use Skadi per-turn: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "skadi_discard" {
			t.Fatalf("Skadi should ask which hand card to discard, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{fire.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve Skadi discard: %v", err)
		}
		if p0.Elements[model.ElementFire] != 2 || len(p0.Graveyard) != 1 || skadi.Statuses["斯卡蒂已用:"+model.ElementFire] != 1 {
			t.Fatalf("Skadi should gain fire and mark it used, elements=%v graveyard=%v statuses=%v", p0.Elements, cardsToInfo(p0.Graveyard), skadi.Statuses)
		}
		skadi.UsedThisTurn = 0
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  skadi.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use Skadi again for arcane: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{arcane.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve Skadi arcane discard: %v", err)
		}
		if p0.Elements[model.ElementArcane] != 2 {
			t.Fatalf("Skadi should allow arcane gain, elements=%v", p0.Elements)
		}
	})

	t.Run("4111003 梵天 and 4211003 水晶心 arm spell-hit effects for the turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		brahma := NewCardInstance(baseCard(t, "4111003"), 0, 1)
		crystal := NewCardInstance(baseCard(t, "4211003"), 0, 1)
		p0.Hero = brahma
		fireSkill := readySkill(baseCard(t, "3121003"), 0)
		target := placeUnit(baseCard(t, "1011002"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  brahma.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use Brahma ultimate: %v", err)
		}
		engine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, fireSkill, map[string]any{"attacker": 0})
		if effectiveElementsGain(brahma)[model.ElementFire] != brahma.Card.ElementsGain[model.ElementFire]+1 {
			t.Fatalf("Brahma should gain one fire load after friendly fire hit, gain=%v", effectiveElementsGain(brahma))
		}

		p0.Hero = crystal
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  crystal.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use Crystal Heart ultimate: %v", err)
		}
		engine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, fireSkill, map[string]any{"attacker": 0})
		if target.Statuses[StatusFreeze] != 0 {
			t.Fatalf("field trigger without target should not freeze arbitrary units, target=%v", target.Statuses)
		}
		engine.triggerEffects(TriggerOnSpellHit, crystal, target, map[string]any{"attacker": 0})
		if target.Statuses[StatusFreeze] != 1 {
			t.Fatalf("Crystal Heart armed spells should freeze hit target, target=%v", target.Statuses)
		}
	})

	t.Run("4311001 肃 discards two air cards to damage an enemy", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		su := NewCardInstance(baseCard(t, "4311001"), 0, 1)
		p0.Hero = su
		p0.Hand = []*CardInstance{
			NewCardInstance(baseCard(t, "3321005"), 0, 1),
			NewCardInstance(baseCard(t, "3321006"), 0, 1),
		}
		enemy := placeUnit(baseCard(t, "1011002"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  su.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use Su ultimate: %v", err)
		}
		if len(p0.Hand) != 0 || len(p0.Graveyard) != 2 || enemy.CurrentLife != enemy.Card.Life-1 {
			t.Fatalf("Su should discard two air cards and deal one damage, hand=%d graveyard=%d life=%d", len(p0.Hand), len(p0.Graveyard), enemy.CurrentLife)
		}
	})

	t.Run("4411001 白须 searches an earth beast plant or spirit on first turn only", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		whitebeard := NewCardInstance(baseCard(t, "4411001"), 0, 1)
		earth := NewCardInstance(baseCard(t, "1421003"), 0, 1)
		other := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Hero = whitebeard
		p0.Deck = []*CardInstance{other, earth}
		engine.State.TurnNumber = 1

		engine.triggerEffects(TriggerOnTurnStart, whitebeard, nil, nil)
		if len(p0.Hand) != 1 || p0.Hand[0] != earth || len(p0.Deck) != 1 {
			t.Fatalf("Whitebeard should search the earth creature, hand=%v deck=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Deck))
		}
		engine.State.TurnNumber = 2
		engine.triggerEffects(TriggerOnTurnStart, whitebeard, nil, nil)
		if len(p0.Hand) != 1 {
			t.Fatalf("Whitebeard should not search after turn one, hand=%v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("4411003 麦吉 discounts first original high-cost play or learn by earth", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		maggie := NewCardInstance(baseCard(t, "4411003"), 0, 1)
		p0.Hero = maggie
		highCost := NewCardInstance(baseCard(t, "1411003"), 0, 1)
		cost := engine.effectiveCardPlayCost(p0, highCost)
		if cost[model.ElementEarth] != max(highCost.Card.ElementsCost[model.ElementEarth]-2, 0) || maggie.Statuses["麦吉折扣"] != 0 {
			t.Fatalf("Maggie should discount first high-cost card by earth, cost=%v statuses=%v", cost, maggie.Statuses)
		}
		skill := NewCardInstance(baseCard(t, "3421014"), 0, 1)
		learnCost := engine.effectiveSkillLearnCost(p0, skill)
		if learnCost[model.ElementEarth] != max(skill.Card.ElementsCost[model.ElementEarth]-2, 0) {
			t.Fatalf("Maggie should also discount learning high-cost skills, cost=%v", learnCost)
		}
	})

	t.Run("4511001 玛丽斯 offers triggered choice when friendly units take enemy damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		maris := NewCardInstance(baseCard(t, "4511001"), 0, 1)
		p0.Hero = maris
		friend := placeUnit(baseCard(t, "1011002"), 0, 1, 0, engine)

		engine.dealDamageWithExtra(friend, 1, 0, map[string]any{"attacker": 1})
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "maris_gain_light" {
			t.Fatalf("Maris should open a triggered choice after enemy damage, pending=%v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{maris.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve Maris triggered choice: %v", err)
		}
		if p0.Elements[model.ElementLight] != 2 {
			t.Fatalf("Maris should grant 2 light after enemy damage, elements=%v", p0.Elements)
		}
	})

	t.Run("4611002 芙雅 doubles companion attack and load then marks it temporary", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		fuye := NewCardInstance(baseCard(t, "4611002"), 0, 1)
		p0.Hero = fuye
		target := placeUnit(baseCard(t, "1011002"), 0, 1, 0, engine)
		attack := target.CurrentAttack
		load := effectiveElementsGain(target)[model.ElementArcane]

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  fuye.InstanceID,
			"ability_type": "ultimate",
			"target_id":    target.InstanceID,
		}}); err != nil {
			t.Fatalf("use Fuye ultimate: %v", err)
		}
		if target.CurrentAttack != attack*2 || effectiveElementsGain(target)[model.ElementArcane] != load*2 || target.Statuses["临时"] != 1 {
			t.Fatalf("Fuye should double attack/load and mark temporary, attack=%d gain=%v statuses=%v", target.CurrentAttack, effectiveElementsGain(target), target.Statuses)
		}
	})
}

func TestRemainingHighRiskBaseCardSemantics(t *testing.T) {
	t.Run("3021011 统御者的制裁 must be paid with a single element", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1011002"), 1, 1, 0, engine)
		sanction := readySkill(baseCard(t, "3021011"), 0)
		p0.Skills[0] = sanction
		p0.Elements[model.ElementFire] = 2
		p0.Elements[model.ElementWater] = 2

		err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": sanction.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}})
		if err == nil {
			t.Fatalf("mixed payment should be rejected for Overlord Sanction")
		}

		p0.Elements[model.ElementFire] = 4
		p0.Elements[model.ElementWater] = 0
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": sanction.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
			"payment":     map[string]any{model.ElementFire: float64(4)},
		}}); err != nil {
			t.Fatalf("single-element payment should cast Overlord Sanction: %v target=%v", err, target.InstanceID)
		}
	})

	t.Run("3321001 闪电链 can damage one extra enemy target outside normal range", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		front := placeUnit(baseCard(t, "1011002"), 1, 1, 0, engine)
		back := placeUnit(baseCard(t, "1011002"), 1, 1, 2, engine)
		chain := readySkill(baseCard(t, "3321001"), 0)
		p0.Skills[0] = chain
		p0.Elements[model.ElementAir] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id":      chain.InstanceID,
			"target_type":      "unit",
			"target_col":       float64(1),
			"target_row":       float64(0),
			"extra_target_col": float64(1),
			"extra_target_row": float64(2),
		}}); err != nil {
			t.Fatalf("cast lightning chain: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve lightning chain: %v", err)
		}
		if front.CurrentLife != front.Card.Life-1 || back.CurrentLife != back.Card.Life-1 {
			t.Fatalf("lightning chain should damage main and extra target, front=%d back=%d", front.CurrentLife, back.CurrentLife)
		}
	})

	t.Run("3521011 光之庇护 protects the chosen friendly companion from lethal damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1011002"), 0, 1, 0, engine)
		shelter := readySkill(baseCard(t, "3521011"), 0)
		p0.Skills[0] = shelter
		p0.Elements[model.ElementLight] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": shelter.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast light shelter: %v", err)
		}
		if target.Statuses["防止致命"] == 0 {
			t.Fatalf("light shelter should mark chosen target, statuses=%v", target.Statuses)
		}
		engine.dealDamageWithExtra(target, 99, 0, map[string]any{"attacker": 1})
		if target.CurrentLife != 1 {
			t.Fatalf("light shelter should prevent lethal damage and leave target at 1, life=%d", target.CurrentLife)
		}
	})

	t.Run("4311002 渡鸦 starts with one extra hand card and should not draw again on turn start", func(t *testing.T) {
		if cards.CardDB == nil {
			if err := cards.LoadCards(); err != nil {
				t.Fatalf("load cards: %v", err)
			}
		}
		SetCardDB(cards.CardDB)
		RegisterAllCardEffects()
		deck := &model.Deck{
			HeroID:   "4311002",
			MainDeck: []string{"1021001", "1021001", "1021002", "1021002", "1021003", "1021003"},
		}
		other := &model.Deck{
			HeroID:   "4011001",
			MainDeck: []string{"1021001", "1021001", "1021002", "1021002", "1021003", "1021003"},
		}
		engine := NewEngine("raven-start", nil)
		if err := engine.SetupGame("Raven", deck, "Other", other); err != nil {
			t.Fatalf("setup game: %v", err)
		}
		if len(engine.State.Players[0].Hand) != 5 || len(engine.State.Players[1].Hand) != 4 {
			t.Fatalf("Raven should start with 5 cards while other hero starts with 4, p0=%d p1=%d", len(engine.State.Players[0].Hand), len(engine.State.Players[1].Hand))
		}
		before := len(engine.State.Players[0].Hand)
		engine.triggerEffects(TriggerOnTurnStart, engine.State.Players[0].Hero, nil, nil)
		if len(engine.State.Players[0].Hand) != before {
			t.Fatalf("Raven should not draw an extra card at turn start, before=%d after=%d", before, len(engine.State.Players[0].Hand))
		}
	})

	t.Run("3021010 解咒 remains a cooldown reaction card placeholder for defense-spell countering", func(t *testing.T) {
		dispel := readySkill(baseCard(t, "3021010"), 0)
		if skillCooldown(dispel) != 1 || skillNeedsTargetInstance(dispel) {
			t.Fatalf("dispel should be cooldown 1 and targetless, cooldown=%d needsTarget=%v", skillCooldown(dispel), skillNeedsTargetInstance(dispel))
		}
	})
}

func TestRemainingMediumRiskGenericBaseCards(t *testing.T) {
	t.Run("1511003 2121008 2221003 2221009 2421005 2421008 2621001 2621009 generic companion/item traits", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		pegasus := NewCardInstance(baseCard(t, "1511003"), 0, 1)
		if !cardHasTaunt(pegasus) {
			t.Fatalf("1511003 should expose taunt from its 引魔 keyword")
		}

		if spellArea(readySkill(baseCard(t, "2121008"), 0)) != SpellAreaSquare {
			t.Fatalf("2121008 should be square area")
		}
		if traitsForCardNumber("2221003").statuses[StatusFreeze] != 1 {
			t.Fatalf("2221003 should carry freeze 1 generic status")
		}
		if spellArea(readySkill(baseCard(t, "2221009"), 0)) != SpellAreaSplashCross || traitsForCardNumber("2221009").statuses[StatusFreeze] != 1 {
			t.Fatalf("2221009 should be splash cross freeze")
		}
		petrify := readySkill(baseCard(t, "2421005"), 0)
		target := placeUnit(baseCard(t, "1011002"), 1, 1, 0, engine)
		engine.applyGenericStatusFromDescription(petrify, target)
		if target.Statuses[StatusPetrify] != 2 {
			t.Fatalf("2421005 should apply petrify 2, statuses=%v", target.Statuses)
		}
		if spellArea(readySkill(baseCard(t, "2421008"), 0)) != SpellAreaSquare {
			t.Fatalf("2421008 should be square area")
		}
		if traitsForCardNumber("2621001").statuses[StatusWeaken] != 2 {
			t.Fatalf("2621001 should carry weaken 2 generic status")
		}
		if spellArea(readySkill(baseCard(t, "2621009"), 0)) != SpellAreaSplashCross {
			t.Fatalf("2621009 should be splash cross area")
		}
	})

	t.Run("3121010 3121011 3121013 3221014 3321007 3321011 generic skill traits", func(t *testing.T) {
		magma := readySkill(baseCard(t, "3121010"), 0)
		if spellArea(magma) != SpellAreaSquare || !cardHasPierce(magma) || traitsForCardNumber("3121010").statuses[StatusBurn] != 1 {
			t.Fatalf("3121010 should be square pierce burn")
		}
		ignite := readySkill(baseCard(t, "3121011"), 0)
		if !cardHasRush(ignite) || traitsForCardNumber("3121011").statuses[StatusBurn] != 1 {
			t.Fatalf("3121011 should be rush burn")
		}
		backlash := readySkill(baseCard(t, "3121013"), 0)
		if !isDefenseOnlySkill(backlash.Card) || traitsForCardNumber("3121013").statuses[StatusBurn] != 1 {
			t.Fatalf("3121013 should be defense-only with burn marker")
		}
		iceField := readySkill(baseCard(t, "3221014"), 0)
		if !isDefenseOnlySkill(iceField.Card) {
			t.Fatalf("3221014 should be defense-only")
		}
		sourceWind := readySkill(baseCard(t, "3321007"), 0)
		if skillCooldown(sourceWind) != 1 || skillNeedsTargetInstance(sourceWind) {
			t.Fatalf("3321007 should be cooldown 1 targetless hand refill")
		}
		if spellArea(readySkill(baseCard(t, "3321011"), 0)) != SpellAreaColumn {
			t.Fatalf("3321011 should be column area")
		}
	})

	t.Run("3421007 3421009 3521012 3621009 3621011 3621014 generic skill traits", func(t *testing.T) {
		quake := readySkill(baseCard(t, "3421007"), 0)
		if spellArea(quake) != SpellAreaSquare || traitsForCardNumber("3421007").statuses[StatusStun] != 1 {
			t.Fatalf("3421007 should be square stun")
		}
		fear := readySkill(baseCard(t, "3421009"), 0)
		if !cardHasPierce(fear) || skillCooldown(fear) != 1 || traitsForCardNumber("3421009").statuses[StatusPetrify] != 2 {
			t.Fatalf("3421009 should be pierce cooldown petrify")
		}
		if spellArea(readySkill(baseCard(t, "3521012"), 0)) != SpellAreaColumn {
			t.Fatalf("3521012 should be column area")
		}
		curse := readySkill(baseCard(t, "3621009"), 0)
		if !cardHasRush(curse) || traitsForCardNumber("3621009").statuses[StatusWeaken] != 2 {
			t.Fatalf("3621009 should be rush weaken")
		}
		dimensional := readySkill(baseCard(t, "3621011"), 0)
		if !cardHasPierce(dimensional) || spellArea(dimensional) != SpellAreaSquare {
			t.Fatalf("3621011 should be pierce square")
		}
		karma := readySkill(baseCard(t, "3621014"), 0)
		if !isDefenseOnlySkill(karma.Card) {
			t.Fatalf("3621014 should be defense-only")
		}
	})
}

func TestDamagedAndDeathTriggeredCardEffects(t *testing.T) {
	t.Run("dolphin sacrifices to prevent lethal damage to another friendly unit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		dolphin := placeUnit(baseCard(t, "1221001"), 0, 0, 0, engine)
		ally := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		ally.CurrentLife = 1

		engine.dealDamage(ally, 5, 0)

		if ally.CurrentLife != 1 || engine.State.Players[0].Units[1][0] != ally {
			t.Fatalf("dolphin should leave ally alive at 1, ally=%+v", ally)
		}
		if engine.State.Players[0].Units[0][0] != nil || len(engine.State.Players[0].Graveyard) != 1 || engine.State.Players[0].Graveyard[0] != dolphin {
			t.Fatalf("dolphin should be sacrificed, units=%v graveyard=%v", engine.State.Players[0].Units[0][0], cardsToInfo(engine.State.Players[0].Graveyard))
		}
	})

	t.Run("bifang increases only burn damage, not all damage to burning units", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		bifang := placeUnit(baseCard(t, "1111003"), 0, 0, 0, engine)
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		target.CurrentLife = 4
		target.Statuses[StatusBurn] = 1

		engine.dealDamage(target, 1, 1)
		if target.CurrentLife != 3 {
			t.Fatalf("normal damage to a burning unit should not be increased by bifang, life=%d source=%+v", target.CurrentLife, bifang)
		}

		engine.dealDamageWithExtra(target, 1, 1, map[string]any{"status_damage": StatusBurn})
		if target.CurrentLife != 1 {
			t.Fatalf("burn damage should be increased by bifang to 2 total, life=%d", target.CurrentLife)
		}
	})

	t.Run("xinke can be summoned free from hand or deck after friendly damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		ally := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		handXinke := NewCardInstance(baseCard(t, "1401002"), 0, 1)
		deckXinke := NewCardInstance(baseCard(t, "1401002"), 0, 1)
		p0.Hand = []*CardInstance{handXinke}
		p0.Deck = []*CardInstance{deckXinke}

		engine.dealDamage(ally, 1, 0)

		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "xinke_summon" || len(engine.State.PendingAction.Candidates) != 2 {
			t.Fatalf("xinke should offer hand and deck summon candidates, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{handXinke.InstanceID},
		}}); err != nil {
			t.Fatalf("summon xinke from hand: %v", err)
		}
		if len(p0.Hand) != 0 || handXinke.Position == nil || p0.Units[handXinke.Position.Col][handXinke.Position.Row] != handXinke {
			t.Fatalf("hand xinke should be summoned free, hand=%v position=%v", cardsToInfo(p0.Hand), handXinke.Position)
		}
		if len(p0.Deck) != 1 || p0.Deck[0] != deckXinke {
			t.Fatalf("deck xinke should remain when hand copy chosen, deck=%v", cardsToInfo(p0.Deck))
		}
	})

	t.Run("great druid ultimate arms the next friendly companion death into a life seed summon", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		druid := placeUnit(baseCard(t, "1411001"), 0, 0, 0, engine)
		ally := placeUnit(baseCard(t, "1421003"), 0, 1, 0, engine)
		addElementsGainBonus(ally, model.ElementEarth, 2)
		ally.CurrentLife = ally.Card.Life + 2

		engine.destroyUnit(ally, 0)
		if len(p0.Graveyard) != 1 || len(p0.Hand) != 0 {
			t.Fatalf("druid should not trigger before ultimate, graveyard=%v hand=%v", cardsToInfo(p0.Graveyard), cardsToInfo(p0.Hand))
		}

		ally2 := placeUnit(baseCard(t, "1421003"), 0, 1, 0, engine)
		addElementsGainBonus(ally2, model.ElementEarth, 2)
		ally2.CurrentLife = ally2.Card.Life + 2
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  druid.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use druid ultimate: %v", err)
		}
		engine.destroyUnit(ally2, 0)

		var seed *CardInstance
		for _, unit := range engine.getAllFieldCards(p0) {
			if unit.Card.Number == "1401001" {
				seed = unit
				break
			}
		}
		if seed == nil {
			t.Fatalf("druid ultimate should summon a life seed")
		}
		if seed.CurrentLife != seed.Card.Life+2 || effectiveElementsGain(seed)[model.ElementEarth] != seed.Card.ElementsGain[model.ElementEarth]+2 {
			t.Fatalf("life seed should inherit life/load bonuses, life=%d load=%v", seed.CurrentLife, effectiveElementsGain(seed))
		}
	})
}
