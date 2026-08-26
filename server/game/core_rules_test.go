package game

import (
	"reflect"
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
		"item_weapon": {
			Number:       "item_weapon",
			Type:         model.CardTypeItem,
			Name:         "Test Weapon",
			Category:     model.ElementAir,
			Tag:          "装备-武器",
			ElementsCost: map[string]int{model.ElementAir: 1},
			ElementsGain: map[string]int{model.ElementAir: 1},
			Life:         -1,
			Attack:       2,
			Power:        -1,
		},
		"item_tool": {
			Number:       "item_tool",
			Type:         model.CardTypeItem,
			Name:         "Test Tool",
			Category:     model.ElementEarth,
			Tag:          "装备",
			ElementsCost: map[string]int{model.ElementAir: 1},
			ElementsGain: map[string]int{model.ElementEarth: 1},
			Life:         -1,
			Attack:       -1,
			Power:        -1,
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
	if err := engine.SetupGameWithFirstPlayer("P1", deck, "P2", deck, 0); err != nil {
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

type testSetupHeroEnter struct{ AlwaysActive }

func (testSetupHeroEnter) ID() string   { return "hero_enter" }
func (testSetupHeroEnter) Name() string { return "Setup Hero" }
func (testSetupHeroEnter) OnEnter(ctx *EffectContext) error {
	ctx.Source.Statuses["entered"]++
	if len(ctx.Engine.State.Players[ctx.PlayerID].Hand) == 0 {
		ctx.Source.Statuses["entered_before_initial_draw"] = 1
	}
	return nil
}

func (testSetupHeroEnter) OnGameStart(ctx *EffectContext) error {
	ctx.Source.Statuses["game_start"]++
	if ctx.Source.Statuses["entered"] == 1 && len(ctx.Engine.State.Players[ctx.PlayerID].Hand) == 0 {
		ctx.Source.Statuses["game_start_after_enter_before_initial_draw"] = 1
	}
	return nil
}

type testSetupHeroWatcher struct{ AlwaysActive }

func (testSetupHeroWatcher) ID() string   { return "hero_watcher" }
func (testSetupHeroWatcher) Name() string { return "Watcher Hero" }
func (testSetupHeroWatcher) OnUnitEnter(ctx *EffectContext) error {
	enteredPlayer, _ := ctx.ExtraData["entered_player"].(int)
	if enteredPlayer != ctx.PlayerID {
		ctx.Source.Statuses["saw_enemy_initial_enter"]++
	}
	return nil
}

type testPierceSkill struct{ AlwaysActive }

func (testPierceSkill) ID() string      { return "spell_pierce" }
func (testPierceSkill) Name() string    { return "Piercing Bolt" }
func (testPierceSkill) HasPierce() bool { return true }

type testDefenseSkill struct{ AlwaysActive }

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

func TestInitialHeroesTriggerEnterBeforeInitialDraw(t *testing.T) {
	previousDB := cardDBRef
	previousRegistry := globalRegistry
	t.Cleanup(func() {
		SetCardDB(previousDB)
		globalRegistry = previousRegistry
	})
	globalRegistry = NewEffectRegistry()
	globalRegistry.RegisterBehaviorFactory("hero_enter", func() CardBehavior { return testSetupHeroEnter{} })
	globalRegistry.RegisterBehaviorFactory("hero_watcher", func() CardBehavior { return testSetupHeroWatcher{} })

	db := map[string]*model.Card{
		"hero_enter": {
			Number:       "hero_enter",
			Type:         model.CardTypeHero,
			Name:         "Setup Hero",
			Category:     model.ElementArcane,
			ElementsGain: map[string]int{model.ElementArcane: 1},
			Life:         6,
			Attack:       -1,
			Power:        -1,
		},
		"hero_watcher": {
			Number:       "hero_watcher",
			Type:         model.CardTypeHero,
			Name:         "Watcher Hero",
			Category:     model.ElementArcane,
			ElementsGain: map[string]int{model.ElementArcane: 1},
			Life:         6,
			Attack:       -1,
			Power:        -1,
		},
		"unit_basic": {
			Number:       "unit_basic",
			Type:         model.CardTypeCompanion,
			Name:         "Test Unit",
			Category:     model.ElementArcane,
			ElementsCost: map[string]int{},
			ElementsGain: map[string]int{model.ElementArcane: 1},
			Life:         3,
			Attack:       1,
			Power:        -1,
		},
	}
	SetCardDB(db)
	mainDeck := make([]string, 30)
	for i := range mainDeck {
		mainDeck[i] = "unit_basic"
	}
	deck0 := &model.Deck{HeroID: "hero_enter", MainDeck: mainDeck, SkillPool: []string{}}
	deck1 := &model.Deck{HeroID: "hero_watcher", MainDeck: mainDeck, SkillPool: []string{}}

	events := make([]string, 0)
	engine := NewEngine("initial-hero-enter", func(event GameEvent, targetPlayer int) {
		events = append(events, event.Type)
	})
	if err := engine.SetupGameWithFirstPlayer("P1", deck0, "P2", deck1, 0); err != nil {
		t.Fatalf("setup game: %v", err)
	}

	if got := engine.State.Players[0].Hero.Statuses["entered"]; got != 1 {
		t.Fatalf("initial hero should trigger own enter once, got %d", got)
	}
	if got := engine.State.Players[0].Hero.Statuses["entered_before_initial_draw"]; got != 1 {
		t.Fatalf("initial hero enter should happen before initial draw, got %d", got)
	}
	if got := engine.State.Players[0].Hero.Statuses["game_start_after_enter_before_initial_draw"]; got != 1 {
		t.Fatalf("game start trigger should happen after hero enter and before initial draw, got %d", got)
	}
	if got := engine.State.Players[1].Hero.Statuses["saw_enemy_initial_enter"]; got != 1 {
		t.Fatalf("enemy hero should receive initial unit-enter event, got %d", got)
	}
	setupIndex := indexOfString(events, "game_setup")
	initialDrawIndex := indexOfString(events, "initial_draw")
	if setupIndex == -1 {
		t.Fatalf("setup should emit game_setup event, events=%v", events)
	}
	if initialDrawIndex == -1 {
		t.Fatalf("setup should emit initial_draw event, events=%v", events)
	}
	if setupIndex > initialDrawIndex {
		t.Fatalf("game_setup should be emitted before initial draw, events=%v", events)
	}
}

func indexOfString(values []string, target string) int {
	for i, value := range values {
		if value == target {
			return i
		}
	}
	return -1
}

func TestFirstPlayerFirstTurnDrawsAndHeroLoadIsLimited(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]

	if len(p0.Hand) != 5 {
		t.Fatalf("first player should draw on their first turn, hand=%d", len(p0.Hand))
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": p0.Hero.InstanceID}}); err != nil {
		t.Fatalf("consume first-turn hero: %v", err)
	}
	if p0.Elements[model.ElementAir] != 2 {
		t.Fatalf("first-turn hero should gain half total load rounded up, elements=%v", p0.Elements)
	}
}

func TestFirstTurnLoadLimitAppliesOnlyToFirstPlayerHero(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]

	unit := NewCardInstance(cardDBRef["unit_basic"], 0, engine.State.TurnNumber)
	unit.IsHorizontal = false
	setElementsGain(unit, map[string]int{model.ElementArcane: 3})
	p0.Units[0][0] = unit

	if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": unit.InstanceID}}); err != nil {
		t.Fatalf("consume non-hero on first turn: %v", err)
	}
	if p0.Elements[model.ElementArcane] != 3 {
		t.Fatalf("non-hero load should not be halved on first turn, elements=%v", p0.Elements)
	}
}

func TestFirstTurnMultiElementHeroLoadRequiresExplicitChoice(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	setElementsGain(p0.Hero, map[string]int{model.ElementFire: 1, model.ElementWater: 2})

	err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": p0.Hero.InstanceID}})
	if err == nil {
		t.Fatalf("multi-element first-turn hero load should require an explicit gain choice")
	}
	if p0.Hero.IsHorizontal {
		t.Fatalf("failed first-turn gain choice should not consume the hero")
	}

	err = engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{
		"instance_id": p0.Hero.InstanceID,
		"gain":        map[string]any{model.ElementWater: float64(2)},
	}})
	if err != nil {
		t.Fatalf("consume first-turn multi-element hero with choice: %v", err)
	}
	if p0.Elements[model.ElementWater] != 2 || p0.Elements[model.ElementFire] != 0 {
		t.Fatalf("first-turn hero should gain the explicitly chosen half-load, elements=%v", p0.Elements)
	}
}

func TestElementsClearAtEndOfOwnersTurnNotStartTurn(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	p0.Elements[model.ElementAir] = 3
	p1.Elements[model.ElementWater] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
		t.Fatalf("end turn: %v", err)
	}

	if p0.Elements[model.ElementAir] != 0 {
		t.Fatalf("active player's elements should clear at end of turn, elements=%v", p0.Elements)
	}
	if p1.Elements[model.ElementWater] != 2 {
		t.Fatalf("next player's start turn should not clear their existing elements, elements=%v", p1.Elements)
	}
}
func TestReplacingSkillReturnsOldSkillToPool(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	p0.Elements[model.ElementAir] = 10

	oldSkill := p0.SkillPool[0]
	p0.SkillPool = p0.SkillPool[1:]
	oldSkill.IsHorizontal = false
	oldSkill.SlotIndex = 0
	oldSkill.Statuses[StatusCooldown] = 2
	oldSkill.PowerBonus = 3
	p0.Skills[0] = oldSkill

	newSkill := p0.SkillPool[0]
	if err := engine.HandleAction(0, ActionMessage{Action: "learn_skill", Data: map[string]any{
		"instance_id": newSkill.InstanceID,
		"replace_id":  oldSkill.InstanceID,
	}}); err != nil {
		t.Fatalf("replace skill: %v", err)
	}

	if p0.Skills[0] != newSkill {
		t.Fatalf("new skill should occupy replaced slot")
	}
	if len(p0.Graveyard) != 0 {
		t.Fatalf("replaced skill should not go to graveyard, graveyard=%v", cardsToInfo(p0.Graveyard))
	}
	foundInPool := false
	for _, skill := range p0.SkillPool {
		if skill == oldSkill {
			foundInPool = true
			break
		}
	}
	if !foundInPool || oldSkill.SlotIndex != -1 || len(oldSkill.Statuses) != 0 || oldSkill.PowerBonus != 0 {
		t.Fatalf("replaced skill should return to pool with battlefield state cleared, found=%v slot=%d statuses=%v power=%d", foundInPool, oldSkill.SlotIndex, oldSkill.Statuses, oldSkill.PowerBonus)
	}
}

func TestReplacingEquipmentSendsOldEquipmentToGraveyard(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	for i := 0; i < 5; i++ {
		card := cardDBRef["item_tool"]
		if i == 2 {
			card = cardDBRef["item_weapon"]
		}
		equipment := NewCardInstance(card, 0, engine.State.TurnNumber)
		equipment.IsHorizontal = false
		equipment.SlotIndex = i
		p0.Equipment[i] = equipment
	}
	replaced := p0.Equipment[2]
	p0.Equipment[1].IsHorizontal = true
	newEquipment := NewCardInstance(cardDBRef["item_weapon"], 0, engine.State.TurnNumber)
	p0.Hand = append(p0.Hand, newEquipment)
	p0.Elements[model.ElementAir] = 10

	err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
		"instance_id": newEquipment.InstanceID,
		"replace_id":  p0.Equipment[1].InstanceID,
	}})
	if err == nil {
		t.Fatalf("should not replace horizontal equipment")
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
		"instance_id": newEquipment.InstanceID,
		"replace_id":  replaced.InstanceID,
	}}); err != nil {
		t.Fatalf("replace equipment: %v", err)
	}
	if p0.Equipment[2] != newEquipment || len(p0.Graveyard) != 1 || p0.Graveyard[0] != replaced {
		t.Fatalf("new equipment should replace selected vertical equipment and send old equipment to graveyard, equipment=%v grave=%v", p0.Equipment[2], cardsToInfo(p0.Graveyard))
	}
}

func TestEquipmentCanDirectAttackEnemyFrontRowWithoutHitTrigger(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	weapon := NewCardInstance(cardDBRef["item_weapon"], 0, engine.State.TurnNumber)
	weapon.IsHorizontal = false
	p0.Equipment[0] = weapon
	target := NewCardInstance(getCardDB()["unit_basic"], 1, engine.State.TurnNumber)
	target.Position = &Position{Col: 0, Row: 0}
	target.CurrentLife = 5
	p1.Units[0][0] = target

	hitTriggered := false
	globalRegistry.Register("item_weapon", TriggerOnHit, func(ctx *EffectContext) error {
		hitTriggered = true
		return nil
	})

	if err := engine.HandleAction(0, ActionMessage{Action: "attack", Data: map[string]any{
		"attacker_id": weapon.InstanceID,
		"target_col":  float64(0),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("equipment attack: %v", err)
	}
	if !weapon.IsHorizontal || target.CurrentLife != 3 {
		t.Fatalf("equipment attack should tap source and deal attack damage, horizontal=%v target_life=%d", weapon.IsHorizontal, target.CurrentLife)
	}
	if hitTriggered {
		t.Fatalf("direct attacks should not trigger hit effects")
	}
}

func TestDirectAttackTriggersAttackedBeforeDamage(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	weapon := NewCardInstance(cardDBRef["item_weapon"], 0, engine.State.TurnNumber)
	weapon.IsHorizontal = false
	p0.Equipment[0] = weapon
	target := NewCardInstance(cardDBRef["unit_basic"], 1, engine.State.TurnNumber)
	target.Position = &Position{Col: 0, Row: 0}
	target.CurrentLife = 5
	p1.Units[0][0] = target

	order := []string{}
	globalRegistry.Register("unit_basic", TriggerOnAttacked, func(ctx *EffectContext) error {
		if ctx.Source == target {
			order = append(order, "attacked")
		}
		return nil
	})
	globalRegistry.Register("unit_basic", TriggerOnDamaged, func(ctx *EffectContext) error {
		if ctx.Source == target {
			order = append(order, "damaged")
		}
		return nil
	})

	if err := engine.HandleAction(0, ActionMessage{Action: "attack", Data: map[string]any{
		"attacker_id": weapon.InstanceID,
		"target_col":  float64(0),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("equipment attack: %v", err)
	}
	if strings.Join(order, ",") != "attacked,damaged" {
		t.Fatalf("direct attack should trigger attacked before damaged, order=%v", order)
	}
}

func TestCoordinateActionsRejectMissingOrInvalidCoordinates(t *testing.T) {
	t.Run("summon does not default missing col row to zero", func(t *testing.T) {
		engine := setupCoreRulesEngine(t)
		p0 := engine.State.Players[0]
		unit := NewCardInstance(getCardDB()["unit_basic"], 0, engine.State.TurnNumber)
		p0.Hand = []*CardInstance{unit}
		p0.Elements[model.ElementArcane] = 1

		err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": unit.InstanceID,
			"position":    map[string]any{"col": float64(1), "row": float64(0)},
		}})
		if err == nil {
			t.Fatalf("summon with nested position but missing top-level col/row should fail")
		}
		if p0.Units[0][0] != nil || len(p0.Hand) != 1 || p0.Hand[0] != unit {
			t.Fatalf("invalid summon should not place or remove card, unit00=%v hand=%v", p0.Units[0][0], cardsToInfo(p0.Hand))
		}
	})

	t.Run("attack does not default target id payload to zero position", func(t *testing.T) {
		engine := setupCoreRulesEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		weapon := NewCardInstance(getCardDB()["item_weapon"], 0, engine.State.TurnNumber)
		weapon.IsHorizontal = false
		p0.Equipment[0] = weapon
		target := NewCardInstance(getCardDB()["unit_basic"], 1, engine.State.TurnNumber)
		target.Position = &Position{Col: 0, Row: 0}
		target.CurrentLife = 5
		p1.Units[0][0] = target

		err := engine.HandleAction(0, ActionMessage{Action: "attack", Data: map[string]any{
			"attacker_id": weapon.InstanceID,
			"target_id":   target.InstanceID,
		}})
		if err == nil {
			t.Fatalf("attack with target_id but missing target_col/target_row should fail")
		}
		if weapon.IsHorizontal || target.CurrentLife != 5 {
			t.Fatalf("invalid attack should not tap attacker or damage target, horizontal=%v target_life=%d", weapon.IsHorizontal, target.CurrentLife)
		}
	})

	t.Run("spell target coordinates must be numeric integers in range", func(t *testing.T) {
		cases := []struct {
			name string
			data map[string]any
		}{
			{
				name: "missing target row",
				data: map[string]any{"target_type": "unit", "target_col": float64(0)},
			},
			{
				name: "string target col",
				data: map[string]any{"target_type": "unit", "target_col": "0", "target_row": float64(0)},
			},
			{
				name: "fractional target row",
				data: map[string]any{"target_type": "unit", "target_col": float64(0), "target_row": 0.5},
			},
			{
				name: "out of range target col",
				data: map[string]any{"target_type": "unit", "target_col": float64(3), "target_row": float64(0)},
			},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				engine := setupCoreRulesEngine(t)
				p0 := engine.State.Players[0]
				p1 := engine.State.Players[1]
				skill := readySkill(getCardDB()["spell_attack"], 0)
				p0.Skills[0] = skill
				p0.Elements[model.ElementAir] = 1
				target := NewCardInstance(getCardDB()["unit_basic"], 1, engine.State.TurnNumber)
				target.Position = &Position{Col: 0, Row: 0}
				target.CurrentLife = 5
				p1.Units[0][0] = target

				data := map[string]any{"instance_id": skill.InstanceID}
				for key, value := range tc.data {
					data[key] = value
				}
				err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: data})
				if err == nil {
					t.Fatalf("cast_spell with %s should fail", tc.name)
				}
				if skill.IsHorizontal || target.CurrentLife != 5 || engine.State.PendingSpell != nil {
					t.Fatalf("invalid cast should not mutate state, horizontal=%v target_life=%d pending=%v", skill.IsHorizontal, target.CurrentLife, engine.State.PendingSpell)
				}
			})
		}
	})
}

func TestCoordinateActionsAllowExplicitZeroCoordinates(t *testing.T) {
	t.Run("summon at zero zero", func(t *testing.T) {
		engine := setupCoreRulesEngine(t)
		p0 := engine.State.Players[0]
		unit := NewCardInstance(getCardDB()["unit_basic"], 0, engine.State.TurnNumber)
		p0.Hand = []*CardInstance{unit}
		p0.Elements[model.ElementArcane] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": unit.InstanceID,
			"col":         float64(0),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("summon at explicit zero coordinates should work: %v", err)
		}
		if p0.Units[0][0] != unit || len(p0.Hand) != 0 {
			t.Fatalf("unit should be summoned at (0,0), unit00=%v hand=%v", p0.Units[0][0], cardsToInfo(p0.Hand))
		}
	})

	t.Run("attack at zero zero", func(t *testing.T) {
		engine := setupCoreRulesEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		weapon := NewCardInstance(getCardDB()["item_weapon"], 0, engine.State.TurnNumber)
		weapon.IsHorizontal = false
		p0.Equipment[0] = weapon
		target := NewCardInstance(getCardDB()["unit_basic"], 1, engine.State.TurnNumber)
		target.Position = &Position{Col: 0, Row: 0}
		target.CurrentLife = 5
		p1.Units[0][0] = target

		if err := engine.HandleAction(0, ActionMessage{Action: "attack", Data: map[string]any{
			"attacker_id": weapon.InstanceID,
			"target_col":  float64(0),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("attack at explicit zero coordinates should work: %v", err)
		}
		if !weapon.IsHorizontal || target.CurrentLife != 3 {
			t.Fatalf("attack should tap source and damage target, horizontal=%v target_life=%d", weapon.IsHorizontal, target.CurrentLife)
		}
	})
}

func TestSpellHitBeforeAndAfterDamageTiming(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	skill := readySkill(getCardDB()["spell_attack"], 0)
	p0.Skills[0] = skill
	target := NewCardInstance(getCardDB()["unit_basic"], 1, engine.State.TurnNumber)
	target.Position = &Position{Col: 0, Row: 0}
	target.CurrentLife = 5
	p1.Units[0][0] = target

	order := []string{}
	globalRegistry.Register("spell_attack", TriggerOnSpellHitBeforeDamage, func(ctx *EffectContext) error {
		order = append(order, "before_damage")
		if ctx.ExtraData["timing"] != "before_damage" {
			t.Fatalf("before damage trigger should expose timing, data=%v", ctx.ExtraData)
		}
		return nil
	})
	globalRegistry.Register("unit_basic", TriggerOnDamaged, func(ctx *EffectContext) error {
		if ctx.Source == target {
			order = append(order, "damaged")
		}
		return nil
	})
	globalRegistry.Register("spell_attack", TriggerOnSpellHit, func(ctx *EffectContext) error {
		order = append(order, "after_damage")
		if ctx.ExtraData["timing"] != "after_damage" {
			t.Fatalf("after damage trigger should expose timing, data=%v", ctx.ExtraData)
		}
		return nil
	})

	engine.resolveSpellHit(0, skill, SpellTarget{Type: "unit", Position: Position{Col: 0, Row: 0}}, nil, nil)
	if strings.Join(order, ",") != "before_damage,damaged,after_damage" {
		t.Fatalf("spell hit timing should be before, damage, after; order=%v", order)
	}
}
func TestSpellDoesNotHitWhenAllTargetsAreGone(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	skill := readySkill(getCardDB()["spell_attack"], 0)
	p0.Skills[0] = skill
	lostTarget := p1.Units[0][0]
	if lostTarget == nil {
		lostTarget = NewCardInstance(getCardDB()["unit_basic"], 1, engine.State.TurnNumber)
		lostTarget.Position = &Position{Col: 0, Row: 0}
		p1.Units[0][0] = lostTarget
	}
	p1.Units[0][0] = nil

	hitTriggered := false
	globalRegistry.Register("spell_attack", TriggerOnSpellHit, func(ctx *EffectContext) error {
		hitTriggered = true
		return nil
	})

	delayed := engine.resolveSpellHit(0, skill, SpellTarget{Type: "unit", Position: Position{Col: 0, Row: 0}}, nil, nil)
	if delayed {
		t.Fatalf("lost target should not open a delayed hit window")
	}
	if hitTriggered {
		t.Fatalf("spell should not trigger hit effects when all targets are gone")
	}
	if len(engine.log) == 0 || engine.log[len(engine.log)-1].Type != "spell_miss" {
		t.Fatalf("lost target should emit spell_miss, last=%+v", engine.log)
	}
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

func TestLethalDamageDestroysAfterResolutionScope(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	unit := NewCardInstance(getCardDB()["unit_basic"], 0, engine.State.TurnNumber)
	unit.Position = &Position{Col: 0, Row: 0}
	unit.CurrentLife = 2
	p0.Units[0][0] = unit

	engine.beginResolution()
	engine.dealDamage(unit, 2, 0)
	if p0.Units[0][0] != unit {
		t.Fatalf("lethal damage should not remove unit before resolution ends")
	}
	if len(p0.Graveyard) != 0 {
		t.Fatalf("lethal damage should not send unit to graveyard before resolution ends, grave=%d", len(p0.Graveyard))
	}
	engine.endResolution()

	if p0.Units[0][0] != nil {
		t.Fatalf("unit should be destroyed after resolution ends")
	}
	if len(p0.Graveyard) != 1 || p0.Graveyard[0] != unit {
		t.Fatalf("destroyed unit should enter graveyard after resolution, grave=%v", p0.Graveyard)
	}
}

func TestQueuedDeathIsSkippedIfUnitRecoversBeforeResolutionEnds(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	unit := NewCardInstance(getCardDB()["unit_basic"], 0, engine.State.TurnNumber)
	unit.Position = &Position{Col: 0, Row: 0}
	unit.CurrentLife = 2
	p0.Units[0][0] = unit

	engine.beginResolution()
	engine.dealDamage(unit, 2, 0)
	unit.CurrentLife = 1
	engine.endResolution()

	if p0.Units[0][0] != unit {
		t.Fatalf("unit that recovered to positive life should remain on the field")
	}
	if len(p0.Graveyard) != 0 {
		t.Fatalf("recovered unit should not enter graveyard, grave=%d", len(p0.Graveyard))
	}
}

func TestDeathTriggersResolveAfterAllQueuedDeathsLeaveField(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	first := NewCardInstance(getCardDB()["unit_basic"], 0, engine.State.TurnNumber)
	first.Position = &Position{Col: 0, Row: 0}
	first.CurrentLife = 2
	second := NewCardInstance(getCardDB()["unit_basic"], 0, engine.State.TurnNumber)
	second.Position = &Position{Col: 1, Row: 0}
	second.CurrentLife = 2
	p0.Units[0][0] = first
	p0.Units[1][0] = second

	firstDeathSawSecondAlreadyDead := false
	globalRegistry.Register("unit_basic", TriggerOnDeath, func(ctx *EffectContext) error {
		if ctx.Source == first {
			firstDeathSawSecondAlreadyDead = p0.Units[1][0] == nil && len(p0.Graveyard) == 2
		}
		return nil
	})

	engine.beginResolution()
	engine.dealDamage(first, 2, 0)
	engine.dealDamage(second, 2, 0)
	if firstDeathSawSecondAlreadyDead {
		t.Fatalf("death trigger should not resolve before the death queue is flushed")
	}
	engine.endResolution()

	if !firstDeathSawSecondAlreadyDead {
		t.Fatalf("death trigger should resolve only after all queued deaths leave the field, units=%v grave=%d", p0.Units[1][0], len(p0.Graveyard))
	}
}

func TestDeathTriggersPreferCurrentTurnPlayer(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	p0Unit := NewCardInstance(getCardDB()["unit_basic"], 0, engine.State.TurnNumber)
	p0Unit.Position = &Position{Col: 0, Row: 0}
	p0Unit.CurrentLife = 2
	p1Unit := NewCardInstance(getCardDB()["unit_basic"], 1, engine.State.TurnNumber)
	p1Unit.Position = &Position{Col: 0, Row: 0}
	p1Unit.CurrentLife = 2
	p0.Units[0][0] = p0Unit
	p1.Units[0][0] = p1Unit

	order := []int{}
	globalRegistry.Register("unit_basic", TriggerOnDeath, func(ctx *EffectContext) error {
		order = append(order, ctx.Source.OwnerID)
		return nil
	})

	engine.beginResolution()
	engine.dealDamage(p1Unit, 2, 1)
	engine.dealDamage(p0Unit, 2, 0)
	engine.endResolution()

	if len(order) < 2 || order[0] != 0 || order[1] != 1 {
		t.Fatalf("current turn player's death triggers should resolve first, order=%v", order)
	}
}

func TestDeathsCausedByDeathTriggersWaitForCurrentTriggerGroup(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	first := NewCardInstance(getCardDB()["unit_basic"], 0, engine.State.TurnNumber)
	first.Position = &Position{Col: 0, Row: 0}
	first.CurrentLife = 2
	second := NewCardInstance(getCardDB()["unit_basic"], 0, engine.State.TurnNumber)
	second.Position = &Position{Col: 1, Row: 0}
	second.CurrentLife = 2
	third := NewCardInstance(getCardDB()["unit_basic"], 0, engine.State.TurnNumber)
	third.Position = &Position{Col: 2, Row: 0}
	third.CurrentLife = 2
	p0.Units[0][0] = first
	p0.Units[1][0] = second
	p0.Units[2][0] = third

	order := []string{}
	globalRegistry.Register("unit_basic", TriggerOnDeath, func(ctx *EffectContext) error {
		switch ctx.Source {
		case first:
			order = append(order, "first")
			engine.dealDamage(third, 2, 0)
		case second:
			order = append(order, "second")
		case third:
			order = append(order, "third")
		}
		return nil
	})

	engine.beginResolution()
	engine.dealDamage(first, 2, 0)
	engine.dealDamage(second, 2, 0)
	engine.endResolution()

	want := []string{"first", "second", "third"}
	if !reflect.DeepEqual(order, want) {
		t.Fatalf("deaths caused by death triggers should wait for the current trigger group, got=%v want=%v", order, want)
	}
}

func TestDirectDestroyQueuesDeathTriggersUntilResolutionEnds(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0 := engine.State.Players[0]
	unit := NewCardInstance(getCardDB()["unit_basic"], 0, engine.State.TurnNumber)
	unit.Position = &Position{Col: 0, Row: 0}
	unit.CurrentLife = 2
	p0.Units[0][0] = unit

	triggered := false
	globalRegistry.Register("unit_basic", TriggerOnDeath, func(ctx *EffectContext) error {
		triggered = true
		return nil
	})

	engine.beginResolution()
	engine.destroyUnit(unit, 0)
	if p0.Units[0][0] != nil {
		t.Fatalf("directly destroyed unit should leave the battlefield immediately")
	}
	if triggered {
		t.Fatalf("death trigger should wait until the current resolution ends")
	}
	engine.endResolution()

	if !triggered {
		t.Fatalf("death trigger should resolve after the current resolution ends")
	}
}

func TestSimultaneousHeroDeathInSameResolutionIsDraw(t *testing.T) {
	engine := setupCoreRulesEngine(t)
	p0Hero := engine.State.Players[0].Hero
	p1Hero := engine.State.Players[1].Hero

	engine.beginResolution()
	engine.dealDamage(p0Hero, p0Hero.CurrentLife, 0)
	engine.dealDamage(p1Hero, p1Hero.CurrentLife, 1)
	if engine.State.Phase == PhaseGameOver {
		t.Fatalf("game should not end before simultaneous death resolution completes")
	}
	engine.endResolution()

	if engine.State.Phase != PhaseGameOver {
		t.Fatalf("game should end after both heroes die")
	}
	if engine.State.Winner != -2 {
		t.Fatalf("simultaneous hero death should be a draw, winner=%d", engine.State.Winner)
	}
	last := engine.log[len(engine.log)-1]
	if last.Type != "game_over" || last.Data["reason"] != "both_heroes_killed" {
		t.Fatalf("draw should emit both_heroes_killed game_over, last=%+v", last)
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

func TestSurrenderEndsGameForOpponent(t *testing.T) {
	engine := setupCoreRulesEngine(t)

	if err := engine.HandleAction(0, ActionMessage{Action: "surrender", Data: map[string]any{}}); err != nil {
		t.Fatalf("surrender: %v", err)
	}
	if engine.State.Phase != PhaseGameOver || engine.State.Winner != 1 {
		t.Fatalf("surrender should make opponent win, phase=%s winner=%d", engine.State.Phase, engine.State.Winner)
	}
	last := engine.log[len(engine.log)-1]
	if last.Type != "game_over" || last.Data["reason"] != "surrender" || last.Data["actor"] != 0 {
		t.Fatalf("surrender should emit game_over with actor, last=%+v", last)
	}
}

func TestDrawOfferAcceptEndsGameAsDraw(t *testing.T) {
	engine := setupCoreRulesEngine(t)

	if err := engine.HandleAction(0, ActionMessage{Action: "offer_draw", Data: map[string]any{}}); err != nil {
		t.Fatalf("offer draw: %v", err)
	}
	if engine.State.DrawOfferBy != 0 || engine.State.Phase == PhaseGameOver {
		t.Fatalf("draw offer should remain pending without ending the game, offer=%d phase=%s", engine.State.DrawOfferBy, engine.State.Phase)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "respond_draw_offer", Data: map[string]any{"accept": true}}); err != nil {
		t.Fatalf("accept draw offer: %v", err)
	}
	if engine.State.Phase != PhaseGameOver || engine.State.Winner != -2 || engine.State.DrawOfferBy != -1 {
		t.Fatalf("accepted draw should end as draw and clear offer, phase=%s winner=%d offer=%d", engine.State.Phase, engine.State.Winner, engine.State.DrawOfferBy)
	}
	last := engine.log[len(engine.log)-1]
	if last.Type != "game_over" || last.Data["reason"] != "draw_agreement" || last.Data["winner"] != -2 {
		t.Fatalf("accepted draw should emit draw game_over, last=%+v", last)
	}
}

func TestDrawOfferRejectClearsOffer(t *testing.T) {
	engine := setupCoreRulesEngine(t)

	if err := engine.HandleAction(0, ActionMessage{Action: "offer_draw", Data: map[string]any{}}); err != nil {
		t.Fatalf("offer draw: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "respond_draw_offer", Data: map[string]any{"accept": false}}); err != nil {
		t.Fatalf("reject draw offer: %v", err)
	}
	if engine.State.Phase == PhaseGameOver || engine.State.DrawOfferBy != -1 {
		t.Fatalf("rejected draw should keep game active and clear offer, phase=%s offer=%d", engine.State.Phase, engine.State.DrawOfferBy)
	}
	last := engine.log[len(engine.log)-1]
	if last.Type != "draw_offer_declined" || last.Data["player"] != 1 || last.Data["offer_by"] != 0 {
		t.Fatalf("reject should emit declined event, last=%+v", last)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "offer_draw", Data: map[string]any{}}); err != nil {
		t.Fatalf("opponent should be able to offer draw after rejection: %v", err)
	}
	if engine.State.DrawOfferBy != 1 {
		t.Fatalf("new draw offer should be tracked, offer=%d", engine.State.DrawOfferBy)
	}
}
