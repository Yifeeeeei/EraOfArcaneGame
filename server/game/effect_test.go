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
	if err := engine.SetupGameWithFirstPlayer("P1", deck, "P2", deck, 0); err != nil {
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
	engine := setupReportedBugEngine(t)
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	target.CurrentLife = 5
	engine.State.Players[1].Shield = 3

	engine.dealDamageWithExtra(target, 2, 1, map[string]any{"damage_source": "spell", "attacker": 0})
	if target.CurrentLife != 5 || engine.State.Players[1].Shield != 1 {
		t.Fatalf("enemy spell damage should hit player shield first, life=%d shield=%d", target.CurrentLife, engine.State.Players[1].Shield)
	}

	engine.dealDamageWithExtra(target, 3, 1, map[string]any{"damage_source": "spell", "attacker": 0})
	if target.CurrentLife != 3 || engine.State.Players[1].Shield != 0 {
		t.Fatalf("spell damage should overflow after shield breaks, life=%d shield=%d", target.CurrentLife, engine.State.Players[1].Shield)
	}

	engine.State.Players[1].Shield = 2
	engine.dealDamageWithExtra(target, 1, 1, map[string]any{"damage_source": "attack", "attacker": 0})
	if target.CurrentLife != 2 || engine.State.Players[1].Shield != 2 {
		t.Fatalf("non-spell damage should ignore player shield, life=%d shield=%d", target.CurrentLife, engine.State.Players[1].Shield)
	}

	engine.dealDamageWithExtra(target, 1, 1, map[string]any{"damage_source": "spell", "attacker": 1})
	if target.CurrentLife != 1 || engine.State.Players[1].Shield != 2 {
		t.Fatalf("friendly spell damage should ignore player shield, life=%d shield=%d", target.CurrentLife, engine.State.Players[1].Shield)
	}
}

func TestRoyalConflictShieldDecayAndStrictArcane(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]

	p0.Shield = 2
	engine.HandleShieldDecay(p0)
	if p0.Shield != 1 {
		t.Fatalf("shield should decay by one without support, got %d", p0.Shield)
	}

	jadeBaron := placeUnit(baseCard(t, "4411101"), 0, 0, 0, engine)
	p0.Shield = 2
	engine.HandleShieldDecay(p0)
	if p0.Shield != 2 {
		t.Fatalf("翡翠男爵 should keep shield below 3 from decaying, got %d", p0.Shield)
	}

	p0.Shield = 3
	engine.HandleShieldDecay(p0)
	if p0.Shield != 2 {
		t.Fatalf("翡翠男爵 should not prevent decay at 3 shield, got %d", p0.Shield)
	}

	jadeBaron.Statuses[StatusPetrify] = 1
	p0.Shield = 2
	engine.HandleShieldDecay(p0)
	if p0.Shield != 1 {
		t.Fatalf("petrified 翡翠男爵 should not prevent shield decay, got %d", p0.Shield)
	}

	engine.gainStrictArcane(0, 3)
	if p0.StrictArcane != 3 {
		t.Fatalf("strict arcane should be tracked separately, got %d", p0.StrictArcane)
	}
	if !engine.spendStrictArcane(0, 2) || p0.StrictArcane != 1 {
		t.Fatalf("strict arcane spend failed, got %d", p0.StrictArcane)
	}
	if engine.spendStrictArcane(0, 2) || p0.StrictArcane != 1 {
		t.Fatalf("strict arcane should not overspend, got %d", p0.StrictArcane)
	}
}

func TestRoyalConflictStealthDoesNotBlockSpellRange(t *testing.T) {
	engine := setupReportedBugEngine(t)
	stealthFront := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	stealthFront.Statuses[StatusStealth] = 1
	back := placeUnit(baseCard(t, "1021002"), 1, 0, 1, engine)

	if engine.IsInSpellRange(0, stealthFront.Position.Col, stealthFront.Position.Row, true) {
		t.Fatal("opposing stealth unit should not be targetable even with pierce")
	}
	if !engine.IsInSpellRange(0, back.Position.Col, back.Position.Row, false) {
		t.Fatal("stealth front row should not block spell range to the next visible row")
	}

	attacker := placeUnit(baseCard(t, "1021003"), 0, 0, 0, engine)
	attacker.IsHorizontal = false
	if engine.IsInAttackRange(0, attacker, stealthFront.Position.Col, stealthFront.Position.Row) {
		t.Fatal("opposing stealth unit should not be targetable by direct attack")
	}
}

func TestRoyalConflictPublicSpecialZones(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	host := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	under := NewCardInstance(baseCard(t, "1021002"), 0, 1)
	p0.Hand = []*CardInstance{under}

	if !engine.placeCardUnder(host, under) {
		t.Fatal("expected hand card to be placed under host")
	}
	if len(p0.Hand) != 0 || len(host.UnderCards) != 1 || host.UnderCards[0] != under {
		t.Fatalf("under card should move from hand to host, hand=%v under=%v", cardsToInfo(p0.Hand), cardsToInfo(host.UnderCards))
	}
	if info := cardToInfo(host); len(info["under_cards"].([]map[string]any)) != 1 {
		t.Fatalf("card info should expose public under cards, info=%v", info)
	}

	engine.destroyUnit(host, 0)
	if len(p0.Graveyard) != 2 || p0.Graveyard[0] != under || p0.Graveyard[1] != host || len(host.UnderCards) != 0 {
		t.Fatalf("destroying host should release under cards before host, graveyard=%v hostUnder=%v", cardsToInfo(p0.Graveyard), cardsToInfo(host.UnderCards))
	}

	exiled := NewCardInstance(baseCard(t, "1021003"), 0, 1)
	p0.Graveyard = append(p0.Graveyard, exiled)
	if !engine.exileCard(0, exiled) {
		t.Fatal("expected graveyard card to be exiled")
	}
	if len(p0.Exile) != 1 || p0.Exile[0] != exiled {
		t.Fatalf("exile zone should contain moved card, exile=%v", cardsToInfo(p0.Exile))
	}
	for _, card := range p0.Graveyard {
		if card == exiled {
			t.Fatalf("exiled card should leave graveyard, graveyard=%v", cardsToInfo(p0.Graveyard))
		}
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
