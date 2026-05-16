package game

import (
	"strings"
	"testing"

	"eraofarcane/model"
)

func setupCoreRulesEngine(t *testing.T) *Engine {
	t.Helper()
	previousDB := cardDBRef
	previousRegistry := globalRegistry
	t.Cleanup(func() {
		SetCardDB(previousDB)
		globalRegistry = previousRegistry
	})
	globalRegistry = NewEffectRegistry()
	globalRegistry.RegisterBehaviorFactory("spell_pierce", func() CardBehavior { return testPierceSkill{} })
	globalRegistry.RegisterBehaviorFactory("spell_defense", func() CardBehavior { return testDefenseSkill{} })

	db := map[string]*model.Card{
		"hero_air": {
			Number:       "hero_air",
			Type:         model.CardTypeHero,
			Name:         "Test Hero",
			Category:     model.ElementAir,
			ElementsGain: map[string]int{model.ElementAir: 4},
			Life:         6,
			Attack:       -1,
			Power:        -1,
		},
		"unit_basic": {
			Number:       "unit_basic",
			Type:         model.CardTypeCompanion,
			Name:         "Test Unit",
			Category:     model.ElementArcane,
			ElementsCost: map[string]int{model.ElementArcane: 1},
			ElementsGain: map[string]int{model.ElementArcane: 1},
			Life:         3,
			Attack:       1,
			Power:        -1,
		},
		"spell_attack": {
			Number:          "spell_attack",
			Type:            model.CardTypeSkill,
			Name:            "Front Bolt",
			Category:        model.ElementAir,
			Tag:             "法术-驱动",
			Description:     "范围:前排",
			ElementsCost:    map[string]int{model.ElementAir: 1},
			ElementsExpense: map[string]int{model.ElementAir: 1},
			Power:           3,
			Attack:          1,
		},
		"spell_pierce": {
			Number:          "spell_pierce",
			Type:            model.CardTypeSkill,
			Name:            "Piercing Bolt",
			Category:        model.ElementAir,
			Tag:             "法术-驱动",
			Description:     "穿透",
			ElementsCost:    map[string]int{model.ElementAir: 1},
			ElementsExpense: map[string]int{model.ElementAir: 1},
			Power:           3,
			Attack:          1,
		},
		"spell_defense": {
			Number:          "spell_defense",
			Type:            model.CardTypeSkill,
			Name:            "Wall",
			Category:        model.ElementEarth,
			Tag:             "法术-创造",
			Description:     "防御",
			ElementsCost:    map[string]int{model.ElementEarth: 1},
			ElementsExpense: map[string]int{model.ElementEarth: 1},
			Power:           4,
			Attack:          0,
		},
		"spell_boost": {
			Number:          "spell_boost",
			Type:            model.CardTypeSkill,
			Name:            "Boost",
			Category:        model.ElementAir,
			Tag:             "法术-幻变",
			Description:     "",
			ElementsCost:    map[string]int{model.ElementAir: 1},
			ElementsExpense: map[string]int{model.ElementAir: 1},
			Power:           2,
			Attack:          0,
		},
	}

	for _, card := range db {
		if card.ElementsCost == nil {
			card.ElementsCost = map[string]int{}
		}
		if card.ElementsGain == nil {
			card.ElementsGain = map[string]int{}
		}
		if card.ElementsExpense == nil {
			card.ElementsExpense = map[string]int{}
		}
	}

	SetCardDB(db)

	mainDeck := make([]string, 30)
	for i := range mainDeck {
		mainDeck[i] = "unit_basic"
	}
	skillPool := []string{
		"spell_attack", "spell_pierce", "spell_defense", "spell_attack", "spell_pierce",
		"spell_defense", "spell_attack", "spell_pierce", "spell_defense", "spell_boost",
	}
	deck := &model.Deck{
		HeroID:    "hero_air",
		MainDeck:  mainDeck,
		SkillPool: skillPool,
	}

	engine := NewEngine("core-rules", nil)
	if err := engine.SetupGame("P1", deck, "P2", deck); err != nil {
		t.Fatalf("setup game: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}}); err != nil {
		t.Fatalf("p0 mulligan: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}}); err != nil {
		t.Fatalf("p1 mulligan: %v", err)
	}

	return engine
}

type testPierceSkill struct{}

func (testPierceSkill) ID() string      { return "spell_pierce" }
func (testPierceSkill) Name() string    { return "Piercing Bolt" }
func (testPierceSkill) HasPierce() bool { return true }

type testDefenseSkill struct{}

func (testDefenseSkill) ID() string               { return "spell_defense" }
func (testDefenseSkill) Name() string             { return "Wall" }
func (testDefenseSkill) IsDefenseOnlySkill() bool { return true }
func (testDefenseSkill) IsSorcerySkill() bool     { return false }
func (testDefenseSkill) NeedsSpellTarget() bool   { return true }
func (testDefenseSkill) SpellArea() SpellArea     { return SpellAreaSingle }
func (testDefenseSkill) HasPierce() bool          { return false }
func (testDefenseSkill) CanUseForSkillPurpose(p skillPurpose) bool {
	return p == skillPurposeDefend || p == skillPurposeDefenseBoost
}

func readySkill(card *model.Card, ownerID int) *CardInstance {
	skill := NewCardInstance(card, ownerID, 1)
	skill.IsHorizontal = false
	return skill
}

func TestPendingActionResolvingAfterSpellHitReturnsToMain(t *testing.T) {
	engine := NewEngine("pending-after-hit", nil)
	engine.State.Players[0] = NewPlayerState(0, "p1", &model.Deck{})
	engine.State.Players[1] = NewPlayerState(1, "p2", &model.Deck{})
	engine.State.Phase = PhaseWaitingAction
	engine.State.ResumePhase = PhaseDefenseWindow
	engine.State.PendingSpell = nil
	engine.State.PendingAction = &PendingAction{
		Type:       "hit_followup",
		PlayerID:   0,
		Prompt:     "resolve follow-up",
		Candidates: []map[string]any{{"instance_id": "choice"}},
		MinSelect:  1,
		MaxSelect:  1,
	}

	err := engine.HandleAction(0, ActionMessage{
		Action: "resolve_action",
		Data:   map[string]any{"selected": []any{"choice"}},
	})
	if err != nil {
		t.Fatalf("resolve action: %v", err)
	}
	if engine.State.Phase != PhaseMain {
		t.Fatalf("expected main phase after defense follow-up resolves, got %s", engine.State.Phase)
	}
}

func TestRequiredPendingActionWithNoCandidatesIsSkipped(t *testing.T) {
	engine := NewEngine("empty-pending", nil)
	engine.State.Phase = PhaseMain

	engine.SetPendingAction(0, "empty", "cannot resolve", nil, 1, 1, nil)

	if engine.State.PendingAction != nil {
		t.Fatalf("expected empty required pending action to be skipped")
	}
	if engine.State.Phase != PhaseMain {
		t.Fatalf("expected phase to remain main, got %s", engine.State.Phase)
	}
}

func TestSpellRangeRequiresEnemyFrontRowUnlessPiercing(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]

	p0.Skills[0] = readySkill(getCardDB()["spell_attack"], 0)
	p0.Elements[model.ElementAir] = 2

	front := NewCardInstance(getCardDB()["unit_basic"], 1, 1)
	front.IsHorizontal = false
	front.Position = &Position{Col: 1, Row: 0}
	p1.Units[1][0] = front

	back := NewCardInstance(getCardDB()["unit_basic"], 1, 1)
	back.IsHorizontal = false
	back.Position = &Position{Col: 0, Row: 2}
	p1.Units[0][2] = back

	err := engine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(0),
			"target_row":  float64(2),
		},
	})
	if err == nil || !strings.Contains(err.Error(), "spell range") {
		t.Fatalf("expected spell range error, got %v", err)
	}
}

func TestPiercingSpellCanTargetBackRow(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]

	p0.Skills[0] = readySkill(getCardDB()["spell_pierce"], 0)
	p0.Elements[model.ElementAir] = 2

	front := NewCardInstance(getCardDB()["unit_basic"], 1, 1)
	front.IsHorizontal = false
	front.Position = &Position{Col: 1, Row: 0}
	p1.Units[1][0] = front

	back := NewCardInstance(getCardDB()["unit_basic"], 1, 1)
	back.IsHorizontal = false
	back.Position = &Position{Col: 0, Row: 2}
	p1.Units[0][2] = back

	err := engine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(0),
			"target_row":  float64(2),
		},
	})
	if err != nil {
		t.Fatalf("piercing spell should target back row: %v", err)
	}
	if engine.State.Phase != PhaseDefenseWindow {
		t.Fatalf("expected defense window, got %v", engine.State.Phase)
	}
}

func TestDefenseRequiresAndPaysElements(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]

	p0.Skills[0] = readySkill(getCardDB()["spell_attack"], 0)
	p0.Elements[model.ElementAir] = 1
	p1.Skills[0] = readySkill(getCardDB()["spell_defense"], 1)

	err := engine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(1),
		},
	})
	if err != nil {
		t.Fatalf("cast spell: %v", err)
	}

	err = engine.HandleAction(1, ActionMessage{
		Action: "defend",
		Data:   map[string]any{"skill_ids": []any{p1.Skills[0].InstanceID}},
	})
	if err == nil || !strings.Contains(err.Error(), "not enough elements") {
		t.Fatalf("expected defense cost error, got %v", err)
	}

	p1.Elements[model.ElementEarth] = 1
	err = engine.HandleAction(1, ActionMessage{
		Action: "defend",
		Data:   map[string]any{"skill_ids": []any{p1.Skills[0].InstanceID}},
	})
	if err != nil {
		t.Fatalf("defend with elements: %v", err)
	}
	if p1.Elements[model.ElementEarth] != 0 {
		t.Fatalf("defense should spend earth element, got %v", p1.Elements)
	}
	if engine.State.Phase != PhaseMain {
		t.Fatalf("expected main phase after defense, got %v", engine.State.Phase)
	}
}

func TestDefenseOnlySkillCannotAttack(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]

	p0.Skills[0] = readySkill(getCardDB()["spell_defense"], 0)
	p0.Elements[model.ElementEarth] = 1

	err := engine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(1),
		},
	})
	if err == nil || !strings.Contains(err.Error(), "cannot be used to attack") {
		t.Fatalf("expected defense-only attack error, got %v", err)
	}
}

func TestDefenseValidationDoesNotPartiallySpendOrTap(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]

	p0.Skills[0] = readySkill(getCardDB()["spell_attack"], 0)
	p0.Elements[model.ElementAir] = 1
	p1.Skills[0] = readySkill(getCardDB()["spell_defense"], 1)
	p1.Skills[1] = readySkill(getCardDB()["spell_boost"], 1)
	p1.Elements[model.ElementEarth] = 1

	err := engine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(1),
		},
	})
	if err != nil {
		t.Fatalf("cast spell: %v", err)
	}

	err = engine.HandleAction(1, ActionMessage{
		Action: "defend",
		Data: map[string]any{
			"skill_ids": []any{p1.Skills[0].InstanceID},
			"boost_ids": []any{p1.Skills[1].InstanceID},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "not enough elements") {
		t.Fatalf("expected combined defense cost error, got %v", err)
	}
	if p1.Elements[model.ElementEarth] != 1 {
		t.Fatalf("failed defense validation should not spend earth, got %v", p1.Elements)
	}
	if p1.Skills[0].IsHorizontal {
		t.Fatalf("failed defense validation should not tap defense skill")
	}
	if p1.Skills[1].IsHorizontal {
		t.Fatalf("failed defense validation should not tap boost skill")
	}
	if engine.State.Phase != PhaseDefenseWindow {
		t.Fatalf("failed defense validation should keep defense window open, got %v", engine.State.Phase)
	}
}

func TestSpellBoostValidationDoesNotPartiallySpendOrTap(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]

	p0.Skills[0] = readySkill(getCardDB()["spell_attack"], 0)
	p0.Skills[1] = readySkill(getCardDB()["spell_boost"], 0)
	p0.Elements[model.ElementAir] = 1

	err := engine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(1),
			"boost_ids":   []any{p0.Skills[1].InstanceID},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "not enough elements") {
		t.Fatalf("expected combined spell cost error, got %v", err)
	}
	if p0.Elements[model.ElementAir] != 1 {
		t.Fatalf("failed spell validation should not spend air, got %v", p0.Elements)
	}
	if p0.Skills[0].IsHorizontal {
		t.Fatalf("failed spell validation should not tap main skill")
	}
	if p0.Skills[1].IsHorizontal {
		t.Fatalf("failed spell validation should not tap boost skill")
	}
	if engine.State.Phase != PhaseMain {
		t.Fatalf("failed spell validation should keep main phase, got %v", engine.State.Phase)
	}
}

func TestHeroDeathEmitsGameOverOnce(t *testing.T) {
	events := make([]GameEvent, 0)
	engine := setupCoreRulesEngine(t)
	engine.callback = func(event GameEvent, targetPlayer int) {
		events = append(events, event)
	}

	p1Hero := engine.State.Players[1].Hero
	engine.dealDamage(p1Hero, p1Hero.CurrentLife, 1)
	engine.checkWinCondition()

	gameOverCount := 0
	for _, event := range events {
		if event.Type == "game_over" {
			gameOverCount++
		}
	}
	if gameOverCount != 1 {
		t.Fatalf("expected one game_over event, got %d", gameOverCount)
	}
}
