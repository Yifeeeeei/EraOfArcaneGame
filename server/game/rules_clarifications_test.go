package game

import (
	"eraofarcane/model"
	"testing"
)

func clarifiedEngine(t *testing.T) *Engine {
	setupBaseCardSmokeSuite(t)
	return setupReportedBugEngine(t)
}

func TestReleasedHeroSetup171(t *testing.T) {
	for _, tc := range []struct {
		hero, initialSkill, expectedSkill, deckCard string
	}{
		{hero: "4111001", expectedSkill: "3101002"},
		{hero: "4211002", expectedSkill: "3201001"},
		{hero: "4611003", deckCard: "2601002"},
		{hero: "4511003", initialSkill: "3521007", expectedSkill: "3501001"},
		{hero: "4511003", initialSkill: "3021009"},
	} {
		t.Run(tc.hero+tc.initialSkill, func(t *testing.T) {
			setupBaseCardSmokeSuite(t)
			main := make([]string, 30)
			for i := range main {
				main[i] = "1021001"
			}
			deck := &model.Deck{HeroID: tc.hero, MainDeck: main}
			if tc.initialSkill != "" {
				deck.SkillPool = []string{tc.initialSkill}
			}
			var e *Engine
			checked := false
			e = NewEngine("released-hero-setup", func(event GameEvent, _ int) {
				if event.Type != "game_setup" {
					return
				}
				checked = true
				p := e.State.Players[0]
				if len(p.Hand) != 0 {
					t.Fatal("hero setup occurred after initial draw")
				}
				count := 0
				for _, skill := range p.SkillPool {
					if skill.Card.Number == tc.expectedSkill {
						count++
					}
				}
				if tc.expectedSkill != "" && count != 1 {
					t.Fatalf("missing setup skill %s", tc.expectedSkill)
				}
				if tc.expectedSkill == "" && tc.initialSkill != "" && len(p.SkillPool) != 1 {
					t.Fatal("Lexia granted a skill without replacing Hope Call")
				}
				if tc.deckCard != "" {
					count = 0
					for _, card := range p.Deck {
						if card.Card.Number == tc.deckCard {
							count++
						}
					}
					if count != 3 {
						t.Fatalf("Jieying added %d curses before draw, want 3", count)
					}
				}
			})
			if err := e.SetupGameWithFirstPlayer("P1", deck, "P2", deck, 0); err != nil {
				t.Fatal(err)
			}
			if !checked {
				t.Fatal("missing setup event")
			}
		})
	}
}

func TestTurnEndWaitsForMarkDeathChoice168(t *testing.T) {
	e := setupCoreRulesEngine(t)
	p0, p1 := e.State.Players[0], e.State.Players[1]
	p0.Hand = nil
	unit := placeUnit(getCardDB()["unit_basic"], 0, 0, 0, e)
	unit.CurrentLife = 1
	unit.Statuses[StatusBurn] = 1
	globalRegistry.Register("unit_basic", TriggerOnDeath, func(ctx *EffectContext) error {
		ctx.Engine.SetPendingAction(0, "mark_death_test", "choice", []map[string]any{{"instance_id": "continue"}}, 1, 1, func([]string) {})
		return nil
	})
	ended := false
	globalRegistry.Register("hero_air", TriggerOnTurnEnd, func(*EffectContext) error { ended = true; return nil })
	before := len(p1.Deck)
	if err := e.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction == nil || ended || len(p1.Deck) != before {
		t.Fatal("turn advanced before mark death choice")
	}
	if err := e.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{"continue"}}}); err != nil {
		t.Fatal(err)
	}
	if !ended || e.State.CurrentTurn != 1 || e.State.Phase != PhaseMain || len(p1.Deck) != before-1 {
		t.Fatal("turn did not resume after complete mark death interaction")
	}
}

func TestControlledBoundSkillDestination167(t *testing.T) {
	e := clarifiedEngine(t)
	unit := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	bound := NewCardInstance(baseCard(t, "3121001"), 1, 1)
	markTransferredBoundSkill(bound)
	unit.BoundSkills = []*CardInstance{bound}
	e.State.Players[1].Units[0][0] = nil
	transferCardControl(unit, 0)
	unit.Position = &Position{Col: 0, Row: 0}
	e.State.Players[0].Units[0][0] = unit
	e.destroyUnit(unit, 0)
	if bound.OwnerID != 0 || len(e.State.Players[0].Graveyard) != 1 || len(e.State.Players[0].Exile) != 1 || len(e.State.Players[1].Exile) != 0 {
		t.Fatal("controlled unit and bound skill did not enter new controller's zones")
	}
}

func TestLifeAndZoneClarification167(t *testing.T) {
	e := clarifiedEngine(t)
	p := e.State.Players[0]
	card := NewCardInstance(baseCard(t, "1021001"), 0, 1)
	p.Hand = []*CardInstance{card}
	p.Elements[model.ElementArcane] = 10
	e.gainLife(card, 2, card)
	if err := e.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{"instance_id": card.InstanceID, "col": float64(0), "row": float64(0)}}); err != nil {
		t.Fatal(err)
	}
	if card.CurrentLife != 3 || maxLife(card) != 3 {
		t.Fatalf("off-field bonus lost on entry: %v", cardToInfo(card))
	}
	e.ApplyDamage(DamageRequest{Target: card, Amount: 1})
	e.healUnit(card, 1, card)
	if card.CurrentLife != 3 {
		t.Fatal("healing did not restore the increased maximum")
	}
	e.returnUnitToHand(card, 0)
	returned, _ := p.FindHandCard(card.InstanceID)
	if returned == nil || returned.CurrentLife != 1 || maxLife(returned) != 1 || len(returned.Statuses) != 0 {
		t.Fatalf("return did not clear field state: %v", cardToInfo(returned))
	}
	// A later off-field bonus survives unrelated hidden-zone transfers.
	e.gainLife(returned, 1, returned)
	p.Hand = nil
	p.Graveyard = []*CardInstance{returned}
	e.moveGraveyardCardToHand(0, returned.InstanceID)
	if returned.CurrentLife != 2 || maxLife(returned) != 2 {
		t.Fatal("off-field transfer cleared a new bonus")
	}
	host := placeUnit(baseCard(t, "1021001"), 0, 0, 0, e)
	if !e.placeCardUnder(host, returned) {
		t.Fatal("could not place card under host")
	}
	e.destroyUnit(host, 0)
	if returned.CurrentLife != 2 || maxLife(returned) != 2 {
		t.Fatal("off-field under-card transfers cleared a bonus")
	}
}

func TestDeathZoneUsesCleanCardAndLastKnownSource167(t *testing.T) {
	e := clarifiedEngine(t)
	p := e.State.Players[0]
	card := placeUnit(baseCard(t, "1021001"), 0, 0, 0, e)
	e.gainLife(card, 3, card)
	card.PowerBonus = 4
	e.destroyUnit(card, 0)
	if len(p.Graveyard) != 1 {
		t.Fatal("missing graveyard card")
	}
	dead := p.Graveyard[0]
	if dead.InstanceID != card.InstanceID || dead.CurrentLife != dead.Card.Life || dead.PowerBonus != 0 || dead.Statuses["max_life_bonus"] != 0 {
		t.Fatalf("graveyard retained field state: %v", cardToInfo(dead))
	}
	if card.CurrentLife != 4 {
		t.Fatal("last-known source was overwritten")
	}
}

func TestLethalMarkStopsNextTurn168(t *testing.T) {
	e := setupCoreRulesEngine(t)
	p0, p1 := e.State.Players[0], e.State.Players[1]
	p0.Hand = nil
	p0.Hero.CurrentLife = 1
	p0.Hero.Statuses[StatusBurn] = 1
	before := len(p1.Deck)
	if err := e.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
		t.Fatal(err)
	}
	if e.State.Phase != PhaseGameOver || len(p1.Deck) != before {
		t.Fatalf("next turn advanced after lethal mark: phase=%v deck=%d/%d", e.State.Phase, len(p1.Deck), before)
	}
}

func TestEndTriggersAfterMarksAndUntilEndProtection168(t *testing.T) {
	e := setupCoreRulesEngine(t)
	p := e.State.Players[0]
	p.Hand = nil
	p.Hero.Statuses[StatusBurn] = 1
	p.Hero.Statuses[temporaryDamageAndNegativeImmunityUntilStatus] = e.State.TurnNumber
	before := p.Hero.CurrentLife
	seen := false
	globalRegistry.Register("hero_air", TriggerOnTurnEnd, func(ctx *EffectContext) error {
		if ctx.PlayerID == 0 {
			seen = ctx.Source.Statuses[StatusBurn] == 0 && !ctx.Source.IsHorizontal && ctx.Source.CurrentLife == before
		}
		return nil
	})
	if err := e.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
		t.Fatal(err)
	}
	if !seen || p.Hero.CurrentLife != before {
		t.Fatal("end trigger ran before reset/marks or protection expired early")
	}
}

func TestNextTurnShield168(t *testing.T) {
	e := setupCoreRulesEngine(t)
	p := e.State.Players[0]
	p.Hand = nil
	e.addTemporaryModifier(0, TemporaryModifier{Type: TempModDelayedShieldGain, Amount: 1})
	if err := e.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
		t.Fatal(err)
	}
	if e.State.CurrentTurn != 1 || p.Shield != 1 {
		t.Fatalf("next-turn shield was not active for opponent turn: %d", p.Shield)
	}
}

func TestScarletWingsQueueSurvivesDeparture169(t *testing.T) {
	e := clarifiedEngine(t)
	a := placeUnit(baseCard(t, "1621109"), 0, 0, 0, e)
	b := placeUnit(baseCard(t, "1621109"), 0, 1, 0, e)
	enemy := placeUnit(baseCard(t, "1021002"), 1, 0, 0, e)
	enemy.CurrentLife = 8
	e.triggerScarletWingsAfterRedMoon(0)
	resolve := func(id string) {
		t.Helper()
		if err := e.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{id}}}); err != nil {
			t.Fatal(err)
		}
	}
	resolve(b.InstanceID)
	if e.State.PendingAction == nil {
		t.Fatal("second wing trigger was lost")
	}
	resolve(enemy.InstanceID)
	if enemy.CurrentLife != 7 || a.CurrentLife != 2 || len(e.State.Players[0].Graveyard) != 1 {
		t.Fatal("remaining results did not resolve independently")
	}
	if dead := e.State.Players[0].Graveyard[0]; dead.Statuses["max_life_bonus"] != 0 {
		t.Fatal("departed source received a new field bonus")
	}
}

func TestPetrificationCancelsQueuedWingAndDeathrattle169(t *testing.T) {
	e := clarifiedEngine(t)
	placeUnit(baseCard(t, "1621109"), 0, 0, 0, e)
	b := placeUnit(baseCard(t, "1621109"), 0, 1, 0, e)
	enemy := placeUnit(baseCard(t, "1021002"), 1, 0, 0, e)
	enemy.CurrentLife = 8
	e.triggerScarletWingsAfterRedMoon(0)
	b.Statuses[StatusPetrify] = 1
	if err := e.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{enemy.InstanceID}}}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction != nil || enemy.CurrentLife != 7 {
		t.Fatal("petrified wing still resolved")
	}
	prophet := placeUnit(baseCard(t, "1621111"), 0, 2, 0, e)
	prophet.Statuses[StatusPetrify] = 1
	e.destroyUnit(prophet, 0)
	if e.State.Players[0].NextRedMoonCooldown != 0 {
		t.Fatal("petrified deathrattle resolved")
	}
}

func TestSquareAndStealthRules170(t *testing.T) {
	e := clarifiedEngine(t)
	center := placeUnit(baseCard(t, "1021001"), 1, 1, 1, e)
	center.Statuses[StatusStealth] = 1
	placeUnit(baseCard(t, "1021001"), 1, 0, 2, e)
	skill := readySkill(baseCard(t, "3121005"), 0)
	target := SpellTarget{Type: "unit", Position: Position{Col: 1, Row: 1}, AreaOrigin: &Position{Col: 0, Row: 0}}
	if err := e.validateSpellTarget(0, skill, target); err != nil {
		t.Fatal(err)
	}
	if units := e.spellAffectedUnits(1, skill, target); len(units) != 1 || units[0] != center {
		t.Fatalf("selected square ignored: %v", units)
	}
	if front := e.State.Players[1].GetFrontRow(); front != 2 {
		t.Fatalf("stealth formed front row: %d", front)
	}
	single := readySkill(baseCard(t, "3121006"), 0)
	if err := e.validateSpellTarget(0, single, SpellTarget{Type: "unit", Position: *center.Position}); err == nil {
		t.Fatal("piercing single spell bypassed stealth")
	}
	chain := readySkill(baseCard(t, "3321001"), 0)
	if err := e.validateSpellExtraTargetForSkill(0, chain, SpellTarget{Type: "unit", Position: Position{Col: 0, Row: 2}}, SpellTarget{Type: "unit", Position: *center.Position}, true); err == nil {
		t.Fatal("ignore-range extra target bypassed stealth")
	}
	target.AreaOrigin = &Position{Col: 2, Row: 0}
	if err := e.validateSpellTarget(0, skill, target); err == nil {
		t.Fatal("invalid square origin accepted")
	}
}

func TestEffectSpellSquareSelection170(t *testing.T) {
	for _, path := range []string{"copy", "reflected", "sketch"} {
		t.Run(path, func(t *testing.T) {
			e := clarifiedEngine(t)
			placeUnit(baseCard(t, "1021001"), 1, 1, 1, e)
			skill := readySkill(baseCard(t, "3221001"), 0)
			e.State.Players[0].Elements[model.ElementWater] = 5
			target := SpellTarget{Type: "unit", Position: Position{Col: 1, Row: 1}}
			var err error
			switch path {
			case "copy":
				err = e.startVirtualSpellCastNoBoost(0, skill, target, nil)
			case "reflected":
				err = e.startVirtualSpellCastWithBoosts(0, skill, target, nil, nil)
			case "sketch":
				err = e.castSkillFromSketchScroll(0, skill, target)
			}
			if err != nil {
				t.Fatal(err)
			}
			if e.State.PendingAction == nil || len(e.State.PendingAction.Candidates) != 4 || e.State.PendingSpell != nil || e.State.Players[0].Elements[model.ElementWater] != 5 {
				t.Fatal("effect cast committed before square choice")
			}
			if err = e.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{"square:0:1"}}}); err != nil {
				t.Fatal(err)
			}
			if e.State.PendingSpell == nil || e.State.PendingSpell.Target.AreaOrigin == nil || *e.State.PendingSpell.Target.AreaOrigin != (Position{Col: 0, Row: 1}) {
				t.Fatal("selected effect square lost before defense")
			}
		})
	}
}

func TestSlotStackingAndCounterReplacement172(t *testing.T) {
	e := clarifiedEngine(t)
	p := e.State.Players[0]
	for i, number := range []string{"2021017", "2021105", "2021002", "2021002"} {
		p.Equipment[i] = NewCardInstance(baseCard(t, number), 0, 1)
	}
	if equipmentSlotCapacity(p) != 9 || skillSlotCapacity(p) != 7 {
		t.Fatalf("grants did not stack: %d/%d", equipmentSlotCapacity(p), skillSlotCapacity(p))
	}
	for i := range p.Equipment {
		p.Equipment[i] = nil
	}
	for i := 0; i < BaseEquipmentSlots; i++ {
		p.Equipment[i] = NewCardInstance(baseCard(t, "2021006"), 0, 1)
		p.Equipment[i].IsHorizontal = false
	}
	trap := NewCardInstance(baseCard(t, "2121002"), 0, 1)
	p.Hand = []*CardInstance{trap}
	old := p.Equipment[0]
	if err := e.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": trap.InstanceID, "replace_id": old.InstanceID}}); err != nil {
		t.Fatal(err)
	}
	if p.Equipment[0] != trap || !trap.IsSetCounter || len(p.Graveyard) != 1 {
		t.Fatal("full-zone counter replacement failed")
	}
}

func TestReplacementCooldownAndSubtype172(t *testing.T) {
	e := setupCoreRulesEngine(t)
	p := e.State.Players[0]
	for i := 0; i < BaseEquipmentSlots; i++ {
		p.Equipment[i] = NewCardInstance(getCardDB()["item_tool"], 0, 1)
		p.Equipment[i].IsHorizontal = false
	}
	weapon := NewCardInstance(getCardDB()["item_weapon"], 0, 1)
	p.Hand = []*CardInstance{weapon}
	p.Elements[model.ElementAir] = 10
	if err := e.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{"instance_id": weapon.InstanceID, "replace_id": p.Equipment[0].InstanceID}}); err != nil {
		t.Fatal(err)
	}
	old := NewCardInstance(getCardDB()["spell_attack"], 0, 1)
	old.IsHorizontal = false
	old.Statuses[StatusCooldown] = 2
	p.Skills[0] = old
	skill := NewCardInstance(getCardDB()["spell_boost"], 0, 1)
	p.SkillPool = []*CardInstance{skill}
	before := p.Elements[model.ElementAir]
	if err := e.HandleAction(0, ActionMessage{Action: "learn_skill", Data: map[string]any{"instance_id": skill.InstanceID, "replace_id": old.InstanceID}}); err == nil || p.Elements[model.ElementAir] != before {
		t.Fatal("cooldown replacement accepted or spent resources")
	}
}

func TestDefenseAndStunClarifications173(t *testing.T) {
	t.Run("window consume rejected", func(t *testing.T) {
		e := setupCoreRulesEngine(t)
		p := e.State.Players[1]
		p.Hero.IsHorizontal = false
		e.State.Phase = PhaseDefenseWindow
		e.State.PendingSpell = &SpellCast{AttackerID: 0, Skill: NewCardInstance(getCardDB()["spell_attack"], 0, 1)}
		if err := e.HandleAction(1, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": p.Hero.InstanceID}}); err == nil || p.Hero.IsHorizontal {
			t.Fatal("window consume accepted")
		}
	})
	t.Run("zero power still hits without spell", func(t *testing.T) {
		e := setupCoreRulesEngine(t)
		p := e.State.Players[1]
		before := p.Hero.CurrentLife
		e.State.Phase = PhaseDefenseWindow
		e.State.PendingSpell = &SpellCast{AttackerID: 0, Skill: NewCardInstance(getCardDB()["spell_attack"], 0, 1), TotalPower: 0, Target: SpellTarget{Type: "unit", Position: Position{Col: 1, Row: 1}}}
		if err := e.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{"skill_ids": []any{}}}); err != nil {
			t.Fatal(err)
		}
		if p.Hero.CurrentLife != before-1 {
			t.Fatal("empty defense prevented zero-power hit")
		}
	})
	t.Run("no tap ultimate survives stun", func(t *testing.T) {
		e := clarifiedEngine(t)
		card := placeUnit(baseCard(t, "1021012"), 0, 0, 0, e)
		card.IsHorizontal = true
		card.Statuses[StatusStun] = 1
		e.State.Players[0].Hand = []*CardInstance{NewCardInstance(baseCard(t, "2021002"), 0, 1)}
		if err := e.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{"instance_id": card.InstanceID, "ability_type": "ultimate"}}); err != nil || e.State.PendingAction == nil {
			t.Fatalf("non-tap ultimate rejected: %v", err)
		}
	})
	t.Run("tap ability blocked by stun", func(t *testing.T) {
		e := clarifiedEngine(t)
		card := placeUnit(baseCard(t, "1321001"), 0, 0, 0, e)
		card.Statuses[StatusStun] = 1
		if err := e.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{"instance_id": card.InstanceID, "ability_type": "per_turn"}}); err == nil || card.IsHorizontal {
			t.Fatal("stunned source paid a tap cost")
		}
	})
}

func TestDevourAcceptsEquipmentAndLearnedSkills173(t *testing.T) {
	for _, zone := range []string{"equipment", "skill"} {
		t.Run(zone, func(t *testing.T) {
			e := clarifiedEngine(t)
			p := e.State.Players[0]
			number := "2321013"
			if zone == "skill" {
				number = "3121001"
			}
			target := NewCardInstance(baseCard(t, number), 0, 1)
			target.ElementsGainSet = map[string]int{model.ElementAir: 3}
			if zone == "skill" {
				p.Skills[0] = target
			} else {
				p.Equipment[0] = target
			}
			summoned := NewCardInstance(baseCard(t, "1321010"), 0, 1)
			p.Hand = []*CardInstance{summoned}
			setAllElements(p, 30)
			if err := e.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{"instance_id": summoned.InstanceID, "col": float64(0), "row": float64(0), "devour_ids": []any{target.InstanceID}}}); err != nil {
				t.Fatal(err)
			}
			if p.Units[0][0] != summoned || len(p.Graveyard) != 1 || p.Graveyard[0].InstanceID != target.InstanceID || p.Graveyard[0].ElementsGainSet != nil {
				t.Fatal("devour did not discard and clean the eligible battlefield card")
			}
		})
	}
}

// Field visits can have distinct Go objects while retaining the same game ID.
func sameCardInstance(a, b *CardInstance) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.InstanceID == "" || b.InstanceID == "" {
		return a == b
	}
	return a.InstanceID == b.InstanceID
}
