package game

import (
	"eraofarcane/cards"
	"eraofarcane/model"
	"fmt"
	"testing"
)

const testDeckCode = "4311003 // 1021001 1021001 1211203 1311202 1321002 1321002 1321010 1321012 1321013 1321013 1321105 1321204 1321204 1321205 1321205 1321207 1321207 1321210 1321210 1321212 2221205 2221205 2321006 2321006 2321201 2321202 2321207 2321207 2321208 2321211 // 3221209 3311201 3321002 3321007 3321015 3321106 3321206 3321207 3321208 3321209 // 2001201 2001202 2001203 2001204 2001205 2001206 2001207 2001208"

func setupTestEngine(t *testing.T) *Engine {
	t.Helper()

	// Load cards
	if cards.CardDB == nil {
		if err := cards.LoadCards("../../data/all_card_infos.json"); err != nil {
			t.Fatalf("Failed to load cards: %v", err)
		}
		SetCardDB(cards.CardDB)
	}

	// Parse deck
	deck, err := model.ParseDeckCode(testDeckCode)
	if err != nil {
		t.Fatalf("Failed to parse deck: %v", err)
	}
	if err := deck.Validate(cards.CardDB); err != nil {
		t.Fatalf("Deck validation failed: %v", err)
	}

	// Create engine
	events := make([]GameEvent, 0)
	engine := NewEngine("test-game", func(event GameEvent, targetPlayer int) {
		events = append(events, event)
	})

	// Setup game
	if err := engine.SetupGame("Player1", deck, "Player2", deck); err != nil {
		t.Fatalf("Failed to setup game: %v", err)
	}

	return engine
}

func TestGameSetup(t *testing.T) {
	engine := setupTestEngine(t)

	// Check initial state
	if engine.State.Phase != PhaseMulligan {
		t.Errorf("Expected mulligan phase, got %v", engine.State.Phase)
	}

	// Both players should have 4 cards in hand
	for i := 0; i < 2; i++ {
		if len(engine.State.Players[i].Hand) != 4 {
			t.Errorf("Player %d should have 4 cards, got %d", i, len(engine.State.Players[i].Hand))
		}
		if engine.State.Players[i].Hero == nil {
			t.Errorf("Player %d hero is nil", i)
		}
		if engine.State.Players[i].Hero.CurrentLife != 6 {
			t.Errorf("Player %d hero should have 6 HP, got %d", i, engine.State.Players[i].Hero.CurrentLife)
		}
		if len(engine.State.Players[i].Deck) != 26 {
			t.Errorf("Player %d deck should have 26 cards, got %d", i, len(engine.State.Players[i].Deck))
		}
		if len(engine.State.Players[i].SkillPool) != 10 {
			t.Errorf("Player %d skill pool should have 10 cards, got %d", i, len(engine.State.Players[i].SkillPool))
		}
	}

	// Hero should be at position (1,1)
	hero := engine.State.Players[0].Units[1][1]
	if hero == nil {
		t.Errorf("Hero should be at position (1,1)")
	}
}

func TestMulligan(t *testing.T) {
	engine := setupTestEngine(t)

	// Player 0 keeps
	err := engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})
	if err != nil {
		t.Fatalf("Mulligan failed: %v", err)
	}
	if !engine.State.MulliganDone[0] {
		t.Errorf("Player 0 mulligan should be done")
	}

	// Player 1 redraws
	err = engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": false}})
	if err != nil {
		t.Fatalf("Mulligan failed: %v", err)
	}
	if !engine.State.MulliganDone[1] {
		t.Errorf("Player 1 mulligan should be done")
	}

	// Game should now be in main phase
	if engine.State.Phase != PhaseMain {
		t.Errorf("Expected main phase after mulligan, got %v", engine.State.Phase)
	}
	if engine.State.CurrentTurn != 0 {
		t.Errorf("Expected player 0 to start, got %d", engine.State.CurrentTurn)
	}
}

func TestConsumeAndSummon(t *testing.T) {
	engine := setupTestEngine(t)

	// Complete mulligan
	engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})
	engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})

	ps := engine.State.Players[0]

	// Consume hero to get elements
	heroID := ps.Hero.InstanceID
	err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": heroID}})
	if err != nil {
		t.Fatalf("Consume hero failed: %v", err)
	}

	// Check hero is now horizontal
	if !ps.Hero.IsHorizontal {
		t.Errorf("Hero should be horizontal after consume")
	}

	// Check elements gained
	totalElements := 0
	for _, v := range ps.Elements {
		totalElements += v
	}
	// First turn, hero load is halved (4/2=2)
	if totalElements != 2 {
		t.Errorf("Expected 2 elements on first turn (half load), got %d (elements: %v)", totalElements, ps.Elements)
	}

	// Try to summon a companion from hand
	var companion *CardInstance
	for _, card := range ps.Hand {
		if card.Card.IsCompanion() && card.Card.TotalCost() <= totalElements {
			companion = card
			break
		}
	}

	if companion != nil {
		err = engine.HandleAction(0, ActionMessage{
			Action: "summon",
			Data: map[string]any{
				"instance_id": companion.InstanceID,
				"col":         float64(0),
				"row":         float64(0),
			},
		})
		if err != nil {
			t.Logf("Summon failed (may not have enough elements for this card): %v", err)
		} else {
			// Check unit placed
			if ps.Units[0][0] == nil {
				t.Errorf("Unit should be at position (0,0)")
			}
			t.Logf("Successfully summoned %s", companion.Card.Name)
		}
	} else {
		t.Log("No affordable companion in hand, skipping summon test")
	}
}

func TestEndTurnAndSwitch(t *testing.T) {
	engine := setupTestEngine(t)

	// Complete mulligan
	engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})
	engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})

	// Player 0 ends turn
	err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}})
	if err != nil {
		t.Fatalf("End turn failed: %v", err)
	}

	// Should be player 1's turn now
	if engine.State.CurrentTurn != 1 {
		t.Errorf("Expected player 1's turn, got %d", engine.State.CurrentTurn)
	}

	// Player 1 should have drawn a card (5 cards now)
	if len(engine.State.Players[1].Hand) != 5 {
		t.Errorf("Player 1 should have 5 cards after drawing, got %d", len(engine.State.Players[1].Hand))
	}
}

func TestTurnNotYours(t *testing.T) {
	engine := setupTestEngine(t)

	engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})
	engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})

	// Player 1 tries to act during player 0's turn
	err := engine.HandleAction(1, ActionMessage{Action: "end_turn", Data: map[string]any{}})
	if err == nil {
		t.Errorf("Should error when non-active player acts")
	}
}

func TestLearnSkill(t *testing.T) {
	engine := setupTestEngine(t)

	engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})
	engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})

	ps := engine.State.Players[0]

	// Consume hero
	engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": ps.Hero.InstanceID}})

	// Try to learn a skill from pool
	if len(ps.SkillPool) > 0 {
		skill := ps.SkillPool[0]
		err := engine.HandleAction(0, ActionMessage{
			Action: "learn_skill",
			Data:   map[string]any{"instance_id": skill.InstanceID, "replace_id": ""},
		})
		if err != nil {
			t.Logf("Learn skill failed (may not have enough elements): %v", err)
		} else {
			// Check skill is in skill area
			found := false
			for i := 0; i < 5; i++ {
				if ps.Skills[i] != nil && ps.Skills[i].Card.Number == skill.Card.Number {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Skill should be in skill area")
			}
			t.Logf("Successfully learned skill: %s", skill.Card.Name)
		}
	}
}

func TestFullGameFlow(t *testing.T) {
	engine := setupTestEngine(t)

	// Mulligan
	engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})
	engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})

	// Play several turns
	for turn := 0; turn < 6; turn++ {
		currentPlayer := engine.State.CurrentTurn
		ps := engine.State.Players[currentPlayer]

		// Consume hero for elements
		if !ps.Hero.IsHorizontal {
			engine.HandleAction(currentPlayer, ActionMessage{
				Action: "consume",
				Data:   map[string]any{"instance_id": ps.Hero.InstanceID},
			})
		}

		// Try to summon companions
		for _, card := range ps.Hand {
			if card.Card.IsCompanion() && ps.CanPayCost(card.Card.ElementsCost) {
				pos := ps.FindEmptyPosition()
				if pos != nil {
					err := engine.HandleAction(currentPlayer, ActionMessage{
						Action: "summon",
						Data: map[string]any{
							"instance_id": card.InstanceID,
							"col":         float64(pos.Col),
							"row":         float64(pos.Row),
						},
					})
					if err == nil {
						t.Logf("Turn %d P%d: Summoned %s at (%d,%d)", turn+1, currentPlayer, card.Card.Name, pos.Col, pos.Row)
						break
					}
				}
			}
		}

		// End turn
		err := engine.HandleAction(currentPlayer, ActionMessage{Action: "end_turn", Data: map[string]any{}})
		if err != nil {
			t.Fatalf("Turn %d: End turn failed: %v", turn+1, err)
		}

		if engine.State.Phase == PhaseGameOver {
			t.Logf("Game ended at turn %d, winner: %d", turn+1, engine.State.Winner)
			return
		}
	}

	t.Logf("Game still running after 6 turns")
	t.Logf("P0 Hero HP: %d, P1 Hero HP: %d", engine.State.Players[0].Hero.CurrentLife, engine.State.Players[1].Hero.CurrentLife)
	fmt.Printf("P0 units: %d, P1 units: %d\n", engine.State.Players[0].CountUnits(), engine.State.Players[1].CountUnits())
}

func TestGetStateForPlayer(t *testing.T) {
	engine := setupTestEngine(t)

	engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})
	engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})

	// Get state for player 0
	state0 := engine.GetStateForPlayer(0)
	you := state0["you"].(map[string]any)
	opp := state0["opponent"].(map[string]any)

	// Player 0 should see their hand
	if _, ok := you["hand"]; !ok {
		t.Errorf("Player should see their own hand")
	}

	// Should not see opponent's hand
	if _, ok := opp["hand"]; ok {
		t.Errorf("Should not see opponent's hand details")
	}
	if _, ok := opp["hand_count"]; !ok {
		t.Errorf("Should see opponent's hand count")
	}
}
