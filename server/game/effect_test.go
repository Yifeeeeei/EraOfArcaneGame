package game

import (
	"eraofarcane/cards"
	"eraofarcane/model"
	"testing"
)

const effectTestDeck = "4311003 // 1021001 1021001 1021002 1021002 1021004 1021004 1021005 1021005 1021006 1021006 1021007 1021007 1021008 1021008 1021009 1021009 1021010 1021010 1021011 1021011 1021012 1021012 1021013 1021013 1021014 1021014 1021015 1021015 1021016 1021016 // 3321002 3001001 3001002 3021001 3021002 3021003 3021004 3021005 3021006 3021007"

func setupEffectTest(t *testing.T) *Engine {
	t.Helper()
	if cards.CardDB == nil {
		if err := cards.LoadCards(); err != nil {
			t.Fatalf("Failed to load cards: %v", err)
		}
		SetCardDB(cards.PlayableCardDB)
	}
	RegisterAllCardEffects()

	deck, err := model.ParseDeckCode(effectTestDeck)
	if err != nil {
		t.Fatalf("parse effect test deck: %v", err)
	}
	engine := NewEngine("effect-test", func(event GameEvent, targetPlayer int) {})
	if err := engine.SetupGame("P1", deck, "P2", deck); err != nil {
		t.Fatalf("setup effect test game: %v", err)
	}
	engine.HandleAction(0, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})
	engine.HandleAction(1, ActionMessage{Action: "mulligan", Data: map[string]any{"keep": true}})
	return engine
}

func TestKeywordTraits(t *testing.T) {
	setupEffectTest(t)

	rushSkill := NewCardInstance(cards.PlayableCardDB["3021009"], 0, 1)
	rushSkill.IsHorizontal = true
	engine := NewEngine("keyword-traits", func(event GameEvent, targetPlayer int) {})
	engine.ApplyKeywordOnEnter(rushSkill)
	if rushSkill.IsHorizontal {
		t.Fatal("rush skill should enter ready through explicit trait")
	}

	cooldownSkill := NewCardInstance(cards.PlayableCardDB["3421015"], 0, 1)
	engine.ApplyKeywordOnSkillUse(cooldownSkill)
	if cooldownSkill.Statuses[StatusCooldown] != 2 {
		t.Fatalf("cooldown trait = %d, want 2", cooldownSkill.Statuses[StatusCooldown])
	}

	tauntUnit := NewCardInstance(cards.PlayableCardDB["1011001"], 0, 1)
	engine.ApplyKeywordOnEnter(tauntUnit)
	if tauntUnit.Statuses["引魔"] != 1 {
		t.Fatal("taunt trait should apply 引魔 status")
	}
}

func TestEffectRegistry(t *testing.T) {
	r := NewEffectRegistry()

	// Register a test effect
	called := false
	r.Register("test-card", TriggerOnEnter, func(ctx *EffectContext) error {
		called = true
		return nil
	})

	// Verify registration
	effects := r.GetEffects("test-card", TriggerOnEnter)
	if len(effects) != 1 {
		t.Fatalf("Expected 1 effect, got %d", len(effects))
	}

	// Verify HasEffect
	if !r.HasEffect("test-card", TriggerOnEnter) {
		t.Fatal("HasEffect should return true")
	}
	if r.HasEffect("test-card", TriggerOnDeath) {
		t.Fatal("HasEffect should return false for unregistered trigger")
	}

	// Execute the effect
	ctx := &EffectContext{}
	effects[0].Handler(ctx)
	if !called {
		t.Fatal("Effect handler was not called")
	}
}

func TestRegisterAllCardEffectsIsLazy(t *testing.T) {
	previousRegistry := globalRegistry
	t.Cleanup(func() { globalRegistry = previousRegistry })

	RegisterAllCardEffects()

	if got := len(globalRegistry.effects); got != 0 {
		t.Fatalf("RegisterAllCardEffects should not instantiate behavior effects, got %d effect entries", got)
	}

	if !globalRegistry.HasEffect("1021006", TriggerOnEnter) {
		t.Fatal("expected lazy lookup to materialize 1021006 enter behavior")
	}
	if got := len(globalRegistry.effects); got != 1 {
		t.Fatalf("expected only queried card behavior to be materialized, got %d effect entries", got)
	}
}

func TestCardRuleInfoDoesNotMaterializeLazyBehaviors(t *testing.T) {
	previousRegistry := globalRegistry
	t.Cleanup(func() { globalRegistry = previousRegistry })

	if cards.CardDB == nil {
		if err := cards.LoadCards(); err != nil {
			t.Fatalf("load cards: %v", err)
		}
	}
	RegisterAllCardEffects()

	for _, card := range cards.PlayableCardDB {
		_ = CardRuleInfo(card)
	}
	if got := len(globalRegistry.effects); got != 0 {
		t.Fatalf("CardRuleInfo should not materialize behavior effects, got %d effect entries", got)
	}
}

func TestShieldMechanic(t *testing.T) {
	// Test shield damage reduction
	card := &CardInstance{
		CurrentLife: 10,
		Statuses:    map[string]int{"护盾": 3},
	}

	// Shield should block damage
	remaining := ApplyShieldDamage(card, 2)
	if remaining != 0 {
		t.Errorf("Expected 0 remaining damage with shield, got %d", remaining)
	}
	if card.Statuses["护盾"] != 2 {
		t.Errorf("Expected shield 2 after block, got %d", card.Statuses["护盾"])
	}

	// Damage exceeding shield
	remaining = ApplyShieldDamage(card, 5)
	if remaining != 3 {
		t.Errorf("Expected 3 remaining damage, got %d", remaining)
	}
	if card.Statuses["护盾"] != 0 {
		t.Errorf("Expected shield 0 after break, got %d", card.Statuses["护盾"])
	}
}

func TestChargeSystem(t *testing.T) {
	engine := setupEffectTest(t)

	// Add charge
	engine.addCharge(0, 3)
	if engine.State.Players[0].Charge != 3 {
		t.Fatalf("Expected 3 charge, got %d", engine.State.Players[0].Charge)
	}

	// Remove charge
	ok := engine.removeCharge(0, 2)
	if !ok {
		t.Fatal("removeCharge should succeed")
	}
	if engine.State.Players[0].Charge != 1 {
		t.Fatalf("Expected 1 charge, got %d", engine.State.Players[0].Charge)
	}

	// Remove too much
	ok = engine.removeCharge(0, 5)
	if ok {
		t.Fatal("removeCharge should fail when insufficient")
	}
}

func TestEffectSystemIntegration(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]

	// Consume hero to gain elements
	engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": p0.Hero.InstanceID}})

	// Find a summonable companion
	var summonCard *CardInstance
	for _, c := range p0.Hand {
		if c.Card.IsCompanion() && p0.CanPayCost(c.Card.ElementsCost) {
			summonCard = c
			break
		}
	}

	if summonCard != nil {
		desc := summonCard.Card.Description
		t.Logf("Summoning %s (desc: %s)", summonCard.Card.Name, desc)

		err := engine.HandleAction(0, ActionMessage{
			Action: "summon",
			Data:   map[string]any{"instance_id": summonCard.InstanceID, "col": float64(0), "row": float64(0)},
		})
		if err != nil {
			t.Logf("Summon failed: %v", err)
		} else {
			t.Logf("Summoned successfully at (0,0)")

			placed := p0.Units[0][0]
			if placed != nil {
				if cardHasRush(placed) && placed.IsHorizontal {
					t.Error("Rush card should be vertical on enter")
				}
				t.Logf("Card is_horizontal: %v", placed.IsHorizontal)
			}
		}
	}

	t.Log("Effect system integration test completed")
}
