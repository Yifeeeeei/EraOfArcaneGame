package game

import (
	"eraofarcane/cards"
	"eraofarcane/model"
	"testing"
)

// TestFullCombatFlow tests a complete game with summon, consume, attack
func TestFullCombatFlow(t *testing.T) {
	if cards.CardDB == nil {
		if err := cards.LoadCards("../../data/all_card_infos.json"); err != nil {
			t.Fatalf("Failed to load cards: %v", err)
		}
		SetCardDB(cards.CardDB)
	}

	deck, _ := model.ParseDeckCode(testDeckCode)

	events := make([]GameEvent, 0)
	engine := NewEngine("combat-test", func(event GameEvent, targetPlayer int) {
		events = append(events, event)
	})
	engine.SetupGame("P1", deck, "P2", deck)

	// Mulligan
	engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})
	engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})

	if engine.State.Phase != PhaseMain {
		t.Fatalf("Expected main phase, got %v", engine.State.Phase)
	}

	// === Turn 1 (P0): Consume hero, summon a creature ===
	p0 := engine.State.Players[0]
	engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": p0.Hero.InstanceID}})

	t.Logf("P0 elements after consume: %v", p0.Elements)

	// Find a summonable companion
	var summonCard *CardInstance
	for _, c := range p0.Hand {
		if c.Card.IsCompanion() && p0.CanPayCost(c.Card.ElementsCost) {
			summonCard = c
			break
		}
	}
	if summonCard != nil {
		err := engine.HandleAction(0, ActionMessage{
			Action: "summon",
			Data:   map[string]any{"instance_id": summonCard.InstanceID, "col": float64(1), "row": float64(0)},
		})
		if err != nil {
			t.Logf("Summon failed: %v", err)
		} else {
			t.Logf("P0 summoned %s (ATK:%d, HP:%d) at front row", summonCard.Card.Name, summonCard.Card.Attack, summonCard.Card.Life)
		}
	}

	// End turn
	engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}})

	// === Turn 2 (P1): Consume hero, summon a creature ===
	p1 := engine.State.Players[1]
	engine.HandleAction(1, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": p1.Hero.InstanceID}})

	var summonCard2 *CardInstance
	for _, c := range p1.Hand {
		if c.Card.IsCompanion() && p1.CanPayCost(c.Card.ElementsCost) {
			summonCard2 = c
			break
		}
	}
	if summonCard2 != nil {
		engine.HandleAction(1, ActionMessage{
			Action: "summon",
			Data:   map[string]any{"instance_id": summonCard2.InstanceID, "col": float64(1), "row": float64(0)},
		})
		t.Logf("P1 summoned %s (ATK:%d, HP:%d) at front row", summonCard2.Card.Name, summonCard2.Card.Attack, summonCard2.Card.Life)
	}

	engine.HandleAction(1, ActionMessage{Action: "end_turn", Data: map[string]any{}})

	// === Turn 3 (P0): Consume creatures and attack ===
	// Now P0's units should be reset (vertical)
	// Consume all units for elements, then try to attack with front row units
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := p0.Units[col][row]
			if unit != nil && unit.CanConsume() && unit != p0.Hero {
				engine.HandleAction(0, ActionMessage{
					Action: "consume",
					Data:   map[string]any{"instance_id": unit.InstanceID},
				})
				t.Logf("P0 consumed %s, gained %v", unit.Card.Name, unit.Card.ElementsGain)
			}
		}
	}

	// Try to attack with hero if in front row
	p0Front := p0.GetFrontRow()
	p1Front := p1.GetFrontRow()
	t.Logf("P0 front row: %d, P1 front row: %d", p0Front, p1Front)

	// Find an attacker in front row
	for col := 0; col < 3; col++ {
		unit := p0.Units[col][p0Front]
		if unit != nil && !unit.IsHorizontal && unit.CurrentAttack > 0 {
			// Find target in P1 front row
			for tcol := 0; tcol < 3; tcol++ {
				target := p1.Units[tcol][p1Front]
				if target != nil {
					err := engine.HandleAction(0, ActionMessage{
						Action: "attack",
						Data: map[string]any{
							"attacker_id": unit.InstanceID,
							"target_col":  float64(tcol),
							"target_row":  float64(p1Front),
						},
					})
					if err != nil {
						t.Logf("Attack failed: %v", err)
					} else {
						t.Logf("P0 %s attacked P1 %s!", unit.Card.Name, target.Card.Name)
						t.Logf("P1 target HP: %d", target.CurrentLife)
					}
					break
				}
			}
			break
		}
	}

	// End turn
	engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}})

	// Log final state
	t.Logf("=== After 3 turns ===")
	t.Logf("P0 Hero HP: %d, P1 Hero HP: %d", p0.Hero.CurrentLife, p1.Hero.CurrentLife)
	t.Logf("P0 units: %d, P1 units: %d", p0.CountUnits(), p1.CountUnits())
	t.Logf("Game phase: %v, Winner: %d", engine.State.Phase, engine.State.Winner)
}

// TestSkillCastAndDefense tests spell casting and defense mechanics
func TestSkillCastAndDefense(t *testing.T) {
	if cards.CardDB == nil {
		if err := cards.LoadCards("../../data/all_card_infos.json"); err != nil {
			t.Fatalf("Failed to load cards: %v", err)
		}
		SetCardDB(cards.CardDB)
	}

	deck, _ := model.ParseDeckCode(testDeckCode)

	engine := NewEngine("spell-test", func(event GameEvent, targetPlayer int) {})
	engine.SetupGame("P1", deck, "P2", deck)

	// Mulligan
	engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})
	engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})

	p0 := engine.State.Players[0]

	// Consume hero
	engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": p0.Hero.InstanceID}})

	// Learn a skill
	for _, skill := range p0.SkillPool {
		if p0.CanPayCost(skill.Card.ElementsCost) {
			err := engine.HandleAction(0, ActionMessage{
				Action: "learn_skill",
				Data:   map[string]any{"instance_id": skill.InstanceID, "replace_id": ""},
			})
			if err == nil {
				t.Logf("Learned skill: %s (Power:%d, ATK:%d)", skill.Card.Name, skill.Card.Power, skill.Card.Attack)
				break
			}
		}
	}

	// Consume more units for elements if needed
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := p0.Units[col][row]
			if unit != nil && unit.CanConsume() {
				engine.HandleAction(0, ActionMessage{
					Action: "consume",
					Data:   map[string]any{"instance_id": unit.InstanceID},
				})
			}
		}
	}

	// Try to cast a skill
	for i := 0; i < 5; i++ {
		skill := p0.Skills[i]
		if skill != nil && !skill.IsHorizontal {
			cost := skill.Card.ElementsExpense
			if len(cost) == 0 {
				cost = skill.Card.ElementsCost
			}
			if p0.CanPayCost(cost) {
				p1 := engine.State.Players[1]
				frontRow := p1.GetFrontRow()
				// Target enemy hero or front row unit
				var targetCol, targetRow int
				found := false
				for c := 0; c < 3; c++ {
					if p1.Units[c][frontRow] != nil {
						targetCol = c
						targetRow = frontRow
						found = true
						break
					}
				}
				if !found {
					targetCol = 1
					targetRow = 1 // hero position
				}

				err := engine.HandleAction(0, ActionMessage{
					Action: "cast_spell",
					Data: map[string]any{
						"instance_id": skill.InstanceID,
						"target_type": "unit",
						"target_col":  float64(targetCol),
						"target_row":  float64(targetRow),
					},
				})
				if err != nil {
					t.Logf("Cast failed: %v", err)
				} else {
					t.Logf("Cast %s at (%d,%d)", skill.Card.Name, targetCol, targetRow)

					// If it's a sorcery (咒术), it resolves immediately
					// If it's a regular spell, defense window opens
					if engine.State.Phase == PhaseDefenseWindow {
						t.Logf("Defense window opened, power: %d", engine.State.PendingSpell.TotalPower)

						// P1 chooses not to defend
						err = engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}})
						if err != nil {
							t.Logf("No defend failed: %v", err)
						} else {
							t.Logf("Spell hit!")
						}
					} else {
						t.Logf("Sorcery resolved immediately")
					}
					break
				}
			}
		}
	}

	t.Logf("P1 Hero HP after spell: %d", engine.State.Players[1].Hero.CurrentLife)
}
