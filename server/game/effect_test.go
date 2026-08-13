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

	p0.StrictArcane = 0
	p0.Elements[model.ElementArcane] = 0
	container := NewCardInstance(baseCard(t, "2021106"), 0, engine.State.TurnNumber)
	container.IsHorizontal = false
	p0.Equipment[0] = container
	if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": container.InstanceID}}); err != nil {
		t.Fatalf("consume strict arcane container: %v", err)
	}
	if p0.Elements[model.ElementArcane] != 0 || p0.StrictArcane != 2 {
		t.Fatalf("2021106 should gain strict arcane only, elements=%v strict=%d", p0.Elements, p0.StrictArcane)
	}
	if engine.canPayCost(p0, map[string]int{model.ElementFire: 1}) {
		t.Fatalf("strict arcane should not pay non-arcane costs")
	}
	if !engine.payCostForAction(p0, map[string]int{model.ElementArcane: 2}, ActionMessage{}) || p0.StrictArcane != 0 {
		t.Fatalf("strict arcane should pay arcane costs, strict=%d", p0.StrictArcane)
	}

	romulusEngine := setupReportedBugEngine(t)
	romulusP0 := romulusEngine.State.Players[0]
	romulusP0.Elements[model.ElementArcane] = 0
	romulus := NewCardInstance(baseCard(t, "4011102"), 0, romulusEngine.State.TurnNumber)
	romulus.IsHorizontal = false
	romulus.Position = &Position{Col: 1, Row: 1}
	romulusP0.Hero = romulus
	romulusP0.Units[1][1] = romulus
	if err := romulusEngine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": romulus.InstanceID}}); err != nil {
		t.Fatalf("consume Romulus: %v", err)
	}
	if romulusP0.Elements[model.ElementArcane] != 0 || romulusP0.StrictArcane != 4 {
		t.Fatalf("4011102 should gain strict arcane only, elements=%v strict=%d", romulusP0.Elements, romulusP0.StrictArcane)
	}

	defenseEngine := setupReportedBugEngine(t)
	defenseP0 := defenseEngine.State.Players[0]
	defenseP0.Elements[model.ElementArcane] = 0
	defenseContainer := NewCardInstance(baseCard(t, "2021106"), 0, defenseEngine.State.TurnNumber)
	defenseContainer.IsHorizontal = false
	defenseP0.Equipment[0] = defenseContainer
	if defenseEngine.canPayCostWithOverexertOptions(defenseP0, map[string]int{model.ElementFire: 1}, []*CardInstance{defenseContainer}, false) {
		t.Fatalf("overexerted strict arcane source should not pay fire")
	}
	if !defenseEngine.payDefenseCostWithOptions(defenseP0, map[string]int{model.ElementArcane: 2}, ActionMessage{}, []*CardInstance{defenseContainer}, false) {
		t.Fatalf("overexerted strict arcane source should pay arcane")
	}
	if !defenseContainer.IsHorizontal {
		t.Fatalf("overexerted strict arcane source should become horizontal")
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

	stealthFront.Statuses[StatusPetrify] = 1
	if engine.IsInSpellRange(0, back.Position.Col, back.Position.Row, false) {
		t.Fatal("petrified stealth front row should block default spell range again")
	}
	if !engine.IsInSpellRange(0, stealthFront.Position.Col, stealthFront.Position.Row, false) {
		t.Fatal("petrified stealth unit should be targetable as the visible front row")
	}
}

func TestRoyalConflictStealthBeatsGlobalSpellRange(t *testing.T) {
	engine := setupReportedBugEngine(t)
	placeUnit(baseCard(t, "1011002"), 0, 0, 0, engine)
	stealthBack := placeUnit(baseCard(t, "1021001"), 1, 0, 2, engine)
	stealthBack.Statuses[StatusStealth] = 1
	visibleBack := placeUnit(baseCard(t, "1021002"), 1, 1, 2, engine)

	if engine.IsInSpellRange(0, stealthBack.Position.Col, stealthBack.Position.Row, false) {
		t.Fatal("global spell range should not allow targeting opposing stealth units")
	}
	if !engine.IsInSpellRange(0, visibleBack.Position.Col, visibleBack.Position.Row, false) {
		t.Fatal("global spell range should still allow ordinary visible back-row units")
	}
}

func TestRoyalConflictAreaSpellShieldAppliesOnce(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p1 := engine.State.Players[1]
	left := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	right := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
	left.CurrentLife = 5
	right.CurrentLife = 5
	skill := readySkill(baseCard(t, "3121001"), 0)
	skill.AttackBonus = 2 - skill.Card.Attack
	p1.Shield = 1

	engine.resolveSpellHit(0, skill, SpellTarget{Type: "unit", Position: *left.Position}, nil, []SpellTarget{{Type: "unit", Position: *right.Position}})
	if p1.Shield != 0 {
		t.Fatalf("area spell should spend player shield once, got %d", p1.Shield)
	}
	if left.CurrentLife != 4 || right.CurrentLife != 4 {
		t.Fatalf("remaining area spell damage should apply equally after one shield reduction, left=%d right=%d", left.CurrentLife, right.CurrentLife)
	}
}

func TestRoyalConflictAttackPositionEffects(t *testing.T) {
	t.Run("winterfell archer can attack from non-front rows while active", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		frontAlly := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		frontAlly.CurrentAttack = 0
		archer := placeUnit(baseCard(t, "1221103"), 0, 1, 1, engine)
		archer.IsHorizontal = false
		frontEnemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		backEnemy := placeUnit(baseCard(t, "1021001"), 1, 1, 1, engine)

		if !engine.IsInAttackRange(0, archer, frontEnemy.Position.Col, frontEnemy.Position.Row) {
			t.Fatalf("1221103 should attack enemy front row from a non-front row")
		}
		if info := engine.cardToInfo(archer); info["can_attack_from_non_front"] != true {
			t.Fatalf("1221103 should serialize effective non-front attack ability, info=%v", info)
		}
		if engine.IsInAttackRange(0, archer, backEnemy.Position.Col, backEnemy.Position.Row) {
			t.Fatalf("1221103 should still require an enemy front-row target")
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "attack", Data: map[string]any{
			"attacker_id": archer.InstanceID,
			"target_col":  float64(frontEnemy.Position.Col),
			"target_row":  float64(frontEnemy.Position.Row),
		}}); err != nil {
			t.Fatalf("1221103 should attack from non-front row: %v", err)
		}
		if !archer.IsHorizontal || frontEnemy.CurrentLife != frontEnemy.Card.Life-archer.CurrentAttack {
			t.Fatalf("1221103 attack should resolve normally, horizontal=%v enemy_life=%d", archer.IsHorizontal, frontEnemy.CurrentLife)
		}
	})

	t.Run("petrified winterfell archer cannot attack from non-front rows", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		archer := placeUnit(baseCard(t, "1221103"), 0, 1, 1, engine)
		archer.IsHorizontal = false
		archer.Statuses[StatusPetrify] = 1
		frontEnemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if engine.IsInAttackRange(0, archer, frontEnemy.Position.Col, frontEnemy.Position.Row) {
			t.Fatalf("petrified 1221103 should not attack from a non-front row")
		}
		if info := engine.cardToInfo(archer); info["can_attack_from_non_front"] != false {
			t.Fatalf("petrified 1221103 should serialize inactive non-front attack ability, info=%v", info)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "attack", Data: map[string]any{
			"attacker_id": archer.InstanceID,
			"target_col":  float64(frontEnemy.Position.Col),
			"target_row":  float64(frontEnemy.Position.Row),
		}}); err == nil {
			t.Fatal("petrified 1221103 should fail the non-front attack action")
		}
	})

	t.Run("ordinary units still cannot attack from non-front rows", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		ordinary := placeUnit(baseCard(t, "1021002"), 0, 1, 1, engine)
		ordinary.CurrentAttack = 1
		ordinary.IsHorizontal = false
		frontEnemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if engine.IsInAttackRange(0, ordinary, frontEnemy.Position.Col, frontEnemy.Position.Row) {
			t.Fatalf("ordinary units should not attack from a non-front row")
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "attack", Data: map[string]any{
			"attacker_id": ordinary.InstanceID,
			"target_col":  float64(frontEnemy.Position.Col),
			"target_row":  float64(frontEnemy.Position.Row),
		}}); err == nil {
			t.Fatal("ordinary units should fail non-front attack actions")
		}
		if info := engine.cardToInfo(ordinary); info["can_attack_from_non_front"] != false {
			t.Fatalf("ordinary units should not serialize non-front attack ability, info=%v", info)
		}
	})
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

func TestRoyalConflictStaticSpellTraits(t *testing.T) {
	setupReportedBugEngine(t)

	for _, number := range []string{"3001101", "3021105", "3021107", "3021108", "3221104", "3221108", "3321101", "3321105", "3321107"} {
		card := NewCardInstance(baseCard(t, number), 0, 1)
		if !cardHasRush(card) {
			t.Fatalf("%s should have rush", number)
		}
	}

	for number, want := range map[string]int{
		"3021103": 1,
		"3021108": 1,
		"3111101": 1,
		"3221104": 1,
		"3421107": 1,
		"3011101": 2,
		"3021105": 2,
		"3211102": 2,
		"3411101": 2,
	} {
		if got := skillCooldown(NewCardInstance(baseCard(t, number), 0, 1)); got != want {
			t.Fatalf("%s cooldown=%d, want %d", number, got, want)
		}
	}

	for number, want := range map[string]SpellArea{
		"3011101": SpellAreaAll,
		"3111101": SpellAreaAll,
		"2121105": SpellAreaSplashCross,
		"2121112": SpellAreaColumn,
		"2321112": SpellAreaColumn,
		"3121104": SpellAreaColumn,
		"2521112": SpellAreaSquare,
		"3221107": SpellAreaSquare,
		"3421108": SpellAreaSquare,
		"3121107": SpellAreaFrontRow,
		"3221110": SpellAreaFrontRow,
		"3511102": SpellAreaSplashCross,
		"3621109": SpellAreaSplashCross,
	} {
		if got := spellArea(NewCardInstance(baseCard(t, number), 0, 1)); got != want {
			t.Fatalf("%s spell area=%s, want %s", number, got, want)
		}
	}

	for _, number := range []string{"3021102", "3121102", "3121108", "3221103", "3321104", "3521107"} {
		card := baseCard(t, number)
		if !isDefenseOnlySkill(card) || canUseSkillForPurpose(card, skillPurposeAttack) {
			t.Fatalf("%s should be defense-only", number)
		}
	}
	if canUseSkillForPurpose(baseCard(t, "3121105"), skillPurposeDefend) {
		t.Fatal("3121105 should not be usable for defense")
	}

	for _, number := range []string{"3121106", "3511102", "3521103"} {
		if !cardHasPierce(NewCardInstance(baseCard(t, number), 0, 1)) {
			t.Fatalf("%s should have pierce", number)
		}
	}
	for number, status := range map[string]string{
		"3211101": StatusStun,
		"3621109": StatusStun,
		"3221108": StatusFreeze,
		"2121112": StatusBurn,
	} {
		traits := traitsForCardNumber(number)
		if traits.statuses[status] != 1 {
			t.Fatalf("%s should apply %s1, statuses=%v", number, status, traits.statuses)
		}
	}

	for _, number := range []string{"3221104", "3321108", "3321110"} {
		if !skillNeedsTargetCard(baseCard(t, number)) {
			t.Fatalf("%s should require a target", number)
		}
	}
}

func TestRoyalConflictSpellScrollItemsAreSpellLike(t *testing.T) {
	setupReportedBugEngine(t)
	for _, number := range []string{"2121105", "2121109", "2121111", "2121112", "2221110", "2321106", "2521111", "2521112"} {
		card := baseCard(t, number)
		if !isSpellScrollCard(card) || !isSpellLikeCard(card) {
			t.Fatalf("%s should be a spell-like scroll item", number)
		}
		if card.Attack >= 0 || card.Power >= 0 {
			if !skillNeedsTargetCard(card) {
				t.Fatalf("%s should expose spell target requirement", number)
			}
		}
	}
}

func TestRoyalConflictAdditionalSpellBehaviors(t *testing.T) {
	t.Run("six petal snowflake freezes companions but not heroes", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		snowflake := readySkill(baseCard(t, "3221108"), 0)
		companion := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		engine.applyExplicitSpellHitStatuses(snowflake, companion)
		if companion.Statuses[StatusFreeze] != 1 {
			t.Fatalf("3221108 should freeze companion targets, statuses=%v", companion.Statuses)
		}
		hero := placeUnit(baseCard(t, "4311003"), 1, 1, 1, engine)
		engine.State.Players[1].Hero = hero
		engine.applyExplicitSpellHitStatuses(snowflake, hero)
		if hero.Statuses[StatusFreeze] != 0 {
			t.Fatalf("3221108 should not freeze hero targets, statuses=%v", hero.Statuses)
		}
	})

	t.Run("sweeping wind destroys units damaged down to one life", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3321105"), 0)
		target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		target.CurrentLife = 2
		engine.dealDamageWithExtra(target, 1, 1, map[string]any{"damage_source": "effect", "attacker": 0})
		if engine.State.Players[1].Units[0][0] != nil || len(engine.State.Players[1].Graveyard) != 1 || engine.State.Players[1].Graveyard[0] != target {
			t.Fatalf("3321105 should destroy damaged one-life unit, unit=%v grave=%v", cardToInfo(engine.State.Players[1].Units[0][0]), cardsToInfo(engine.State.Players[1].Graveyard))
		}
	})

	t.Run("war trample loses attack for every affected unit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		trample := readySkill(baseCard(t, "3121107"), 0)
		units := []*CardInstance{
			placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine),
			placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine),
			placeUnit(baseCard(t, "1021003"), 1, 2, 0, engine),
		}
		want := max(trample.Card.Attack-len(units), 0)
		if got := engine.effectiveSpellDamage(0, trample, trample.Card.Attack, nil, units); got != want {
			t.Fatalf("3121107 should reduce attack by affected unit count, got %d want %d", got, want)
		}
	})
}

func TestRoyalConflictVanillaCardsAreExplicitlyRegistered(t *testing.T) {
	setupReportedBugEngine(t)

	genericRuleVanilla := map[string]string{
		"1021112": "strict entry payment is enforced by payment rules",
		"1111102": "adjacent spell targets are added by spell target rules",
		"1311102": "empty-deck opponent draw is handled by draw rules",
		"1311103": "opponent hand limit is handled by hand-limit rules",
		"1521109": "light wildcard payment is handled by payment rules",
		"2001102": "discard damage is handled by discard pipeline",
		"2021105": "equipment slot and duplicate subtype rules are global",
		"2021106": "strict arcane source is handled by payment rules",
		"2121102": "back-row fire/air targeting is handled by range rules",
		"2121105": "scroll sacrifice and splash area are handled by spell-use rules",
		"2121111": "response suppression is handled by spell traits",
		"2121112": "column area and burn are handled by spell traits",
		"2221112": "counter behavior is handled by counter-trap rules",
		"2321111": "counter behavior is handled by counter-trap rules",
		"2521111": "support requirement and pierce are handled by spell-use rules",
		"3021102": "defense-only trait is handled by skill traits",
		"3121104": "optional sacrifice and column area are handled by spell-use rules",
		"3121106": "pierce is handled by skill traits",
		"3121108": "defense-only trait is handled by skill traits",
		"3221107": "square area is handled by skill traits",
		"3321101": "rush is handled by skill traits",
		"3321103": "front-row area is handled by skill traits",
		"3411101": "time-cycle lock and strict payment are handled by rules/payment",
		"3421102": "front-row area is handled by skill traits",
		"3421103": "no-boost trait is handled by skill traits",
		"3421109": "petrify is handled by skill traits",
		"3521103": "pierce is handled by skill traits",
		"3521107": "defense-only trait is handled by skill traits",
		"4011102": "strict arcane source is handled by payment rules",
		"4411101": "shield decay prevention is handled by shield rules",
		"4411102": "overexert remainder reward is handled by payment rules",
	}

	for number, card := range cards.PlayableCardDB {
		if card == nil || card.VersionName != "王权纷争" {
			continue
		}
		behavior := globalRegistry.GetBehavior(number)
		if behavior == nil || behavior.ID() != number {
			t.Fatalf("%s %s should have registered behavior, behavior=%#v", number, card.Name, behavior)
		}
		if _, ok := behavior.(vanillaCardBehavior); ok && card.Description != "" {
			if _, reviewed := genericRuleVanilla[number]; !reviewed {
				t.Fatalf("%s %s has effect text but only vanilla behavior; add a concrete behavior or document its generic rule support in this test", number, card.Name)
			}
		}
	}
}

func TestRoyalConflictSophiaFreezeImmunityAndUltimate(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]

	sophia := placeUnit(baseCard(t, "4211102"), 0, 1, 1, engine)
	if !cardHasActiveUltimate(sophia) {
		t.Fatal("4211102 should expose an ultimate ability")
	}
	if engine.addStatus(sophia, StatusFreeze, 1) {
		t.Fatalf("4211102 should reject freeze application, statuses=%v", sophia.Statuses)
	}
	if sophia.Statuses[StatusFreeze] != 0 || engine.hasEffectiveStatus(sophia, StatusFreeze) {
		t.Fatalf("4211102 should remain unfrozen, statuses=%v", sophia.Statuses)
	}
	for _, status := range []string{StatusBurn, StatusStun, StatusPetrify} {
		if !engine.addStatus(sophia, status, 1) {
			t.Fatalf("4211102 should not reject non-freeze negative status %s", status)
		}
		if sophia.Statuses[status] != 1 || !engine.hasEffectiveStatus(sophia, status) {
			t.Fatalf("4211102 should be affected by %s, statuses=%v", status, sophia.Statuses)
		}
		delete(sophia.Statuses, status)
	}

	friendlyFrozen := placeUnit(baseCard(t, "1021001"), 0, 0, 1, engine)
	enemyFrozen := placeUnit(baseCard(t, "1021004"), 1, 0, 0, engine)
	enemyUnfrozen := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
	friendlyFrozen.Statuses[StatusFreeze] = 2
	enemyFrozen.Statuses[StatusFreeze] = 1
	startLife := enemyFrozen.CurrentLife

	if err := (Card4211102WinterfellWarlockSophia{}).OnUltimate(&EffectContext{
		Engine:     engine,
		Source:     sophia,
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("4211102 ultimate: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "sophia_thaw_strike" {
		t.Fatalf("4211102 should prompt for a frozen unit, pending=%+v", engine.State.PendingAction)
	}
	if len(engine.State.PendingAction.Candidates) != 2 {
		t.Fatalf("4211102 should offer exactly frozen units, candidates=%+v", engine.State.PendingAction.Candidates)
	}
	resolvePendingSelection(t, engine, 0, enemyFrozen.InstanceID)
	if enemyFrozen.Statuses[StatusFreeze] != 0 {
		t.Fatalf("4211102 should remove one freeze layer, statuses=%v", enemyFrozen.Statuses)
	}
	if enemyFrozen.CurrentLife != startLife-2 {
		t.Fatalf("4211102 should deal 2 damage after removing freeze, life=%d start=%d", enemyFrozen.CurrentLife, startLife)
	}
	if friendlyFrozen.Statuses[StatusFreeze] != 2 || enemyUnfrozen.Statuses[StatusFreeze] != 0 {
		t.Fatalf("4211102 should not alter unselected units, friendly=%v enemy=%v", friendlyFrozen.Statuses, enemyUnfrozen.Statuses)
	}
	if len(p0.Graveyard) != 0 || len(p1.Graveyard) != 0 {
		t.Fatalf("4211102 test units should survive the damage, p0 grave=%v p1 grave=%v", cardsToInfo(p0.Graveyard), cardsToInfo(p1.Graveyard))
	}
}

func TestRoyalConflictGraceHealsAndRewardsFullyHealedCompanion(t *testing.T) {
	engine := setupReportedBugEngine(t)
	target := placeUnit(baseCard(t, "1021004"), 0, 0, 1, engine)
	other := placeUnit(baseCard(t, "1021003"), 0, 1, 1, engine)
	target.CurrentLife = maxLife(target) - 2
	other.CurrentLife = maxLife(other) - 1
	skill := readySkill(baseCard(t, "3521108"), 0)

	if err := (Card3521108Grace{}).OnSpellCast(&EffectContext{
		Engine:     engine,
		Source:     skill,
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("3521108 cast: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "grace_heal_companion" || len(engine.State.PendingAction.Candidates) != 2 {
		t.Fatalf("3521108 should prompt for wounded friendly companions, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, target.InstanceID)
	if target.CurrentLife != target.Card.Life+1 || maxLife(target) != target.Card.Life+1 || target.Statuses["max_life_bonus"] != 1 || effectiveElementsGain(target)[model.ElementLight] != target.Card.ElementsGain[model.ElementLight]+1 {
		t.Fatalf("3521108 should heal to full then grant +1 life/load, life=%d statuses=%v load=%v", target.CurrentLife, target.Statuses, effectiveElementsGain(target))
	}
	if other.CurrentLife != maxLife(other)-1 || other.Statuses["max_life_bonus"] != 0 {
		t.Fatalf("3521108 should not touch unselected unit, life=%d statuses=%v", other.CurrentLife, other.Statuses)
	}

	partialEngine := setupReportedBugEngine(t)
	partialTarget := placeUnit(baseCard(t, "1221113"), 0, 0, 1, partialEngine)
	partialTarget.CurrentLife = maxLife(partialTarget) - 3
	partialSkill := readySkill(baseCard(t, "3521108"), 0)
	if err := (Card3521108Grace{}).OnSpellCast(&EffectContext{
		Engine:     partialEngine,
		Source:     partialSkill,
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("3521108 partial cast: %v", err)
	}
	resolvePendingSelection(t, partialEngine, 0, partialTarget.InstanceID)
	if partialTarget.CurrentLife != partialTarget.Card.Life-1 || partialTarget.Statuses["max_life_bonus"] != 0 || effectiveElementsGain(partialTarget)[model.ElementLight] != partialTarget.Card.ElementsGain[model.ElementLight] {
		t.Fatalf("3521108 should not grant reward unless fully healed, life=%d statuses=%v load=%v", partialTarget.CurrentLife, partialTarget.Statuses, effectiveElementsGain(partialTarget))
	}
}

func TestRoyalConflictEnterGameSummonsPawnForChosenPlayer(t *testing.T) {
	engine := setupReportedBugEngine(t)
	behavior := Card3001101EnterGame{}
	source := readySkill(baseCard(t, "3001101"), 0)

	if err := behavior.OnSpellCast(&EffectContext{Engine: engine, Source: source, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("enter game spell cast: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "enter_game_player" {
		t.Fatalf("3001101 should ask for target player first, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, "player:1")
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "enter_game_position" {
		t.Fatalf("3001101 should ask for target position after player choice, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, positionSelectionID(Position{Col: 0, Row: 0}))
	summoned := engine.State.Players[1].Units[0][0]
	if summoned == nil || summoned.Card.Number != "1001101" || summoned.OwnerID != 1 {
		t.Fatalf("3001101 should summon abandoned pawn for chosen player, summoned=%v", cardToInfo(summoned))
	}

	staleEngine := setupReportedBugEngine(t)
	if err := behavior.OnSpellCast(&EffectContext{Engine: staleEngine, Source: source, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("enter game stale spell cast: %v", err)
	}
	resolvePendingSelection(t, staleEngine, 0, "player:1")
	blocker := placeUnit(baseCard(t, "1021001"), 1, 0, 0, staleEngine)
	resolvePendingSelection(t, staleEngine, 0, positionSelectionID(Position{Col: 0, Row: 0}))
	if staleEngine.State.Players[1].Units[0][0] != blocker {
		t.Fatalf("3001101 should not overwrite a stale occupied position")
	}
}

func TestRoyalConflictEmeraldBarrierScrollCountsSkillDifference(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	p0.Skills[0] = readySkill(baseCard(t, "3021005"), 0)
	p1.Skills[0] = readySkill(baseCard(t, "3121001"), 1)
	p1.Skills[1] = readySkill(baseCard(t, "3221001"), 1)
	p1.Skills[2] = readySkill(baseCard(t, "3321005"), 1)

	scroll := NewCardInstance(baseCard(t, "2421107"), 0, 1)
	if err := (Card2421107EmeraldBarrierScroll{}).OnUseItem(&EffectContext{Engine: engine, Source: scroll, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("emerald barrier scroll: %v", err)
	}
	if p0.Shield != 2 {
		t.Fatalf("2421107 should gain shield equal to enemy skill surplus, got %d", p0.Shield)
	}

	p0.Skills[1] = readySkill(baseCard(t, "3021005"), 0)
	p0.Skills[2] = readySkill(baseCard(t, "3021005"), 0)
	p0.Shield = 0
	if err := (Card2421107EmeraldBarrierScroll{}).OnUseItem(&EffectContext{Engine: engine, Source: scroll, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("emerald barrier scroll tied count: %v", err)
	}
	if p0.Shield != 0 {
		t.Fatalf("2421107 should not gain shield when enemy has no skill surplus, got %d", p0.Shield)
	}
}

func TestRoyalConflictGiftedYouthMasteryAddsChosenNonArcaneLoad(t *testing.T) {
	engine := setupReportedBugEngine(t)
	youth := placeUnit(baseCard(t, "1021107"), 0, 0, 0, engine)

	engine.advanceMastery(youth, 0, 1)
	if engine.State.PendingAction != nil {
		t.Fatalf("1021107 should not prompt before mastery 2, pending=%+v", engine.State.PendingAction)
	}
	engine.advanceMastery(youth, 0, 1)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "gifted_youth_mastery_load" {
		t.Fatalf("1021107 should prompt for non-arcane load at mastery 2, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, model.ElementAir)
	if effectiveElementsGain(youth)[model.ElementAir] != youth.Card.ElementsGain[model.ElementAir]+1 {
		t.Fatalf("1021107 should gain selected air load, load=%v", effectiveElementsGain(youth))
	}
	if youth.Statuses[StatusMastery] != 2 {
		t.Fatalf("1021107 mastery should reach 2, statuses=%v", youth.Statuses)
	}
}

func TestRoyalConflictSandDustDemonPetrifiesEnemyFrontRow(t *testing.T) {
	engine := setupReportedBugEngine(t)
	frontA := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	frontB := placeUnit(baseCard(t, "1021001"), 1, 2, 0, engine)
	back := placeUnit(baseCard(t, "1021001"), 1, 1, 1, engine)

	if err := (Card1421112SandDustDemon{}).OnPerTurn(&EffectContext{Engine: engine, Source: placeUnit(baseCard(t, "1421112"), 0, 0, 0, engine), PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("sand dust demon prayer: %v", err)
	}
	if frontA.Statuses[StatusPetrify] != 1 || frontB.Statuses[StatusPetrify] != 1 || back.Statuses[StatusPetrify] != 0 {
		t.Fatalf("1421112 should petrify only enemy front row, frontA=%v frontB=%v back=%v", frontA.Statuses, frontB.Statuses, back.Statuses)
	}
	if !(Card1421112SandDustDemon{}).IsPrayerAbility() {
		t.Fatal("1421112 should expose prayer ability")
	}
}

func TestRoyalConflictDemonChildRequiresShadowCompanionDevour(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	demon := NewCardInstance(baseCard(t, "1621108"), 0, engine.State.TurnNumber)
	p0.Hand = []*CardInstance{demon}
	p0.Elements[model.ElementShadow] = 10
	sacrifice := placeUnit(baseCard(t, "1621105"), 0, 0, 0, engine)

	info := engine.cardToInfo(demon)
	req, ok := info["devour_card_requirement"].(DevourCardRequirement)
	if !ok || req.Count != 1 || req.Category != model.ElementShadow || !req.CompanionOnly {
		t.Fatalf("1621108 should expose shadow companion devour requirement, info=%v", info["devour_card_requirement"])
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": demon.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
		"devour_ids":  []any{sacrifice.InstanceID},
	}}); err != nil {
		t.Fatalf("summon demon child with shadow companion devour: %v", err)
	}
	if p0.Units[0][0] != demon || len(p0.Graveyard) != 1 || p0.Graveyard[0] != sacrifice {
		t.Fatalf("1621108 should enter after devouring shadow companion, unit=%v grave=%v", cardToInfo(p0.Units[0][0]), cardsToInfo(p0.Graveyard))
	}

	failEngine := setupReportedBugEngine(t)
	failP0 := failEngine.State.Players[0]
	failDemon := NewCardInstance(baseCard(t, "1621108"), 0, failEngine.State.TurnNumber)
	failP0.Hand = []*CardInstance{failDemon}
	failP0.Elements[model.ElementShadow] = 10
	nonShadow := placeUnit(baseCard(t, "1021001"), 0, 0, 0, failEngine)
	if err := failEngine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": failDemon.InstanceID,
		"col":         float64(1),
		"row":         float64(1),
		"devour_ids":  []any{nonShadow.InstanceID},
	}}); err == nil {
		t.Fatal("1621108 should reject non-shadow devour target")
	}
	if failP0.Units[0][0] != nonShadow || len(failP0.Graveyard) != 0 || len(failP0.Hand) != 1 {
		t.Fatalf("failed 1621108 devour should leave state intact, units=%v grave=%v hand=%v", cardToInfo(failP0.Units[0][0]), cardsToInfo(failP0.Graveyard), cardsToInfo(failP0.Hand))
	}
}

func TestRoyalConflictDreamRippleDamagesEnemyFrontRowTotalThree(t *testing.T) {
	engine := setupReportedBugEngine(t)
	frontA := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	frontB := placeUnit(baseCard(t, "1021002"), 1, 2, 0, engine)
	back := placeUnit(baseCard(t, "1021001"), 1, 1, 1, engine)

	if err := (Card2201103DreamRipple{}).OnUseItem(&EffectContext{Engine: engine, Source: NewCardInstance(baseCard(t, "2201103"), 0, 1), PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("dream ripple: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "dream_ripple_damage" {
		t.Fatalf("2201103 should ask for front enemy damage targets, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, frontA.InstanceID, frontB.InstanceID)
	if frontA.CurrentLife != frontA.Card.Life-2 || frontB.CurrentLife != frontB.Card.Life-1 || back.CurrentLife != back.Card.Life {
		t.Fatalf("2201103 should distribute 3 damage among selected front enemies, frontA=%d frontB=%d back=%d", frontA.CurrentLife, frontB.CurrentLife, back.CurrentLife)
	}
}

func TestRoyalConflictShieldCardBehaviors(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]

	barrierBeast := NewCardInstance(baseCard(t, "1021103"), 0, 1)
	engine.triggerEffects(TriggerOnEnter, barrierBeast, nil, nil)
	if p0.Shield != 2 {
		t.Fatalf("1021103 should gain shield 2 on enter, got %d", p0.Shield)
	}

	p0.Shield = 1
	emeraldGuard := NewCardInstance(baseCard(t, "1421102"), 0, 1)
	engine.triggerEffects(TriggerOnEnter, emeraldGuard, nil, nil)
	if p0.Shield != 1 {
		t.Fatalf("1421102 should not gain shield while player already has shield, got %d", p0.Shield)
	}
	p0.Shield = 0
	engine.triggerEffects(TriggerOnEnter, emeraldGuard, nil, nil)
	if p0.Shield != 2 {
		t.Fatalf("1421102 should gain shield 2 when player has no shield, got %d", p0.Shield)
	}

	p0.Shield = 0
	skyArmor := NewCardInstance(baseCard(t, "2011101"), 0, 1)
	p0.Equipment[0] = skyArmor
	engine.triggerEffects(TriggerOnEnter, skyArmor, nil, nil)
	if p0.Shield != 2 || !p0.CannotGainShield {
		t.Fatalf("2011101 should gain initial shield then block future shield, shield=%d blocked=%v", p0.Shield, p0.CannotGainShield)
	}
	engine.gainPlayerShield(0, 2)
	if p0.Shield != 2 {
		t.Fatalf("2011101 should prevent future shield gains, got %d", p0.Shield)
	}
	engine.HandleShieldDecay(p0)
	if p0.Shield != 2 {
		t.Fatalf("2011101 should prevent shield decay while active, got %d", p0.Shield)
	}
	p0.Equipment[0] = nil
	engine.HandleShieldDecay(p0)
	if p0.Shield != 1 {
		t.Fatalf("2011101 should stop preventing shield decay after leaving field, got %d", p0.Shield)
	}
	engine.gainPlayerShield(0, 2)
	if p0.Shield != 1 {
		t.Fatalf("2011101 should keep blocking future shield gains after leaving field, got %d", p0.Shield)
	}

	p1.Shield = 4
	breakingBlade := NewCardInstance(baseCard(t, "2021102"), 0, 1)
	engine.triggerEffects(TriggerOnEnter, breakingBlade, nil, nil)
	if p1.Shield != 1 {
		t.Fatalf("2021102 should remove opponent shield 3, got %d", p1.Shield)
	}

	barrierEngine := setupReportedBugEngine(t)
	barrierP0 := barrierEngine.State.Players[0]
	barrier := NewCardInstance(baseCard(t, "2021113"), 0, 1)
	barrier.IsSetCounter = true
	barrierP0.Equipment[0] = barrier
	enemySpell := readySkill(baseCard(t, "3021005"), 1)
	if candidates := barrierEngine.eligibleCounterTraps(0, TriggerOnSpellHitBeforeDamage, enemySpell, map[string]any{"attacker": 1}); len(candidates) != 1 || candidates[0] != barrier {
		t.Fatalf("2021113 should be eligible when an enemy spell hits, candidates=%v", candidates)
	}
	barrierEngine.executeCounterTrap(barrier, TriggerOnSpellHitBeforeDamage, enemySpell, map[string]any{"attacker": 1})
	if barrierP0.Shield != 2 || barrierP0.Equipment[0] != nil || len(barrierP0.Graveyard) != 1 || barrierP0.Graveyard[0] != barrier {
		t.Fatalf("2021113 should gain shield 2 then be discarded, shield=%d equipment=%v grave=%v", barrierP0.Shield, barrierP0.Equipment, cardsToInfo(barrierP0.Graveyard))
	}

	friendlyBarrierEngine := setupReportedBugEngine(t)
	friendlyBarrier := NewCardInstance(baseCard(t, "2021113"), 0, 1)
	friendlyBarrier.IsSetCounter = true
	friendlyBarrierEngine.State.Players[0].Equipment[0] = friendlyBarrier
	friendlySpell := readySkill(baseCard(t, "3021005"), 0)
	if candidates := friendlyBarrierEngine.eligibleCounterTraps(0, TriggerOnSpellHitBeforeDamage, friendlySpell, map[string]any{"attacker": 0}); len(candidates) != 0 {
		t.Fatalf("2021113 should not be eligible for friendly spell hits, candidates=%v", candidates)
	}

	p0.CannotGainShield = false
	p0.Shield = 0
	oceanShield := NewCardInstance(baseCard(t, "2221102"), 0, 1)
	behavior, ok := globalRegistry.GetBehavior("2221102").(OnUseItemBehavior)
	if !ok {
		t.Fatal("2221102 should register an item-use behavior")
	}
	if err := behavior.OnUseItem(&EffectContext{Engine: engine, Source: oceanShield, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("use 2221102: %v", err)
	}
	if p0.Shield != 2 {
		t.Fatalf("2221102 should gain shield 2 on use, got %d", p0.Shield)
	}

	t.Run("ash kelt draws and gains light when opponent breaks your shield", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		prince := placeUnit(baseCard(t, "1511101"), 0, 0, 0, engine)
		target := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		p0.Shield = 1
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021002"), 0, 1),
			NewCardInstance(baseCard(t, "1021003"), 0, 1),
		}
		remaining := engine.applyPlayerShieldDamage(target, 2, map[string]any{"damage_source": "spell", "attacker": 1})
		if remaining != 1 || p0.Shield != 0 || len(p0.Hand) != 2 || p0.Elements[model.ElementLight] != 2 {
			t.Fatalf("1511101 should draw 2 and gain 2 light after opponent breaks shield, prince=%v remaining=%d shield=%d hand=%v elements=%v", cardToInfo(prince), remaining, p0.Shield, cardsToInfo(p0.Hand), p0.Elements)
		}
	})

	t.Run("rock wall guard gains shield only after enemy spell hits with no shield", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		guard := placeUnit(baseCard(t, "1021110"), 0, 0, 0, engine)
		behavior := Card1021110RockWallGuard{}

		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: guard, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 1}}); err != nil {
			t.Fatalf("1021110 enemy spell hit: %v", err)
		}
		if p0.Shield != 2 || !guard.UltimateUsed {
			t.Fatalf("1021110 should gain shield 2 and spend ultimate after enemy spell hits while unshielded, shield=%d used=%v", p0.Shield, guard.UltimateUsed)
		}

		p0.Shield = 0
		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: guard, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 1}}); err != nil {
			t.Fatalf("1021110 second enemy spell hit: %v", err)
		}
		if p0.Shield != 0 {
			t.Fatalf("1021110 should not gain shield after its ultimate is spent, shield=%d", p0.Shield)
		}

		friendlyEngine := setupReportedBugEngine(t)
		friendlyP0 := friendlyEngine.State.Players[0]
		friendlyGuard := placeUnit(baseCard(t, "1021110"), 0, 0, 0, friendlyEngine)
		if err := behavior.OnSpellHit(&EffectContext{Engine: friendlyEngine, Source: friendlyGuard, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("1021110 friendly spell hit: %v", err)
		}
		if friendlyP0.Shield != 0 {
			t.Fatalf("1021110 should ignore friendly spell hits, shield=%d", friendlyP0.Shield)
		}
	})
}

func TestRoyalConflictEmeraldImmortalityProtectsWhileShielded(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	emerald := NewCardInstance(baseCard(t, "2411101"), 0, 1)
	p0.Equipment[0] = emerald
	engine.triggerEffects(TriggerOnEnter, emerald, nil, nil)
	if p0.Shield != 2 {
		t.Fatalf("2411101 should gain shield 2 on enter, got %d", p0.Shield)
	}

	ally := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	ally.CurrentLife = 5
	engine.dealDamageWithExtra(ally, 2, 0, map[string]any{"damage_source": "effect", "attacker": 1})
	if ally.CurrentLife != 5 || p0.Shield != 2 {
		t.Fatalf("2411101 should prevent friendly unit damage while shielded, life=%d shield=%d", ally.CurrentLife, p0.Shield)
	}

	ally.Statuses[StatusFreeze] = 1
	if engine.hasEffectiveStatus(ally, StatusFreeze) {
		t.Fatal("2411101 should make friendly negative statuses ineffective while shielded")
	}
	skill := readySkill(baseCard(t, "3621009"), 0)
	p0.Skills[0] = skill
	skill.Statuses[StatusWeaken] = 1
	if !engine.hasEffectiveStatus(skill, StatusWeaken) {
		t.Fatal("2411101 should not protect non-unit cards from negative statuses")
	}

	p0.Shield = 0
	engine.dealDamageWithExtra(ally, 2, 0, map[string]any{"damage_source": "effect", "attacker": 1})
	if ally.CurrentLife != 3 {
		t.Fatalf("2411101 should stop protecting after shield is gone, life=%d", ally.CurrentLife)
	}
	if !engine.hasEffectiveStatus(ally, StatusFreeze) {
		t.Fatal("negative status should become effective after shield is gone")
	}
}

func TestRoyalConflictStealthCardBehaviors(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]

	mistKing := placeUnit(baseCard(t, "1211101"), 0, 0, 0, engine)
	ally := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
	alreadyStealth := placeUnit(baseCard(t, "1021002"), 0, 2, 0, engine)
	alreadyStealth.Statuses[StatusStealth] = 1
	engine.triggerEffects(TriggerOnEnter, mistKing, nil, nil)
	if mistKing.Statuses[StatusStealth] != 2 || ally.Statuses[StatusStealth] != 2 || alreadyStealth.Statuses[StatusStealth] != 1 {
		t.Fatalf("1211101 should give stealth2 only to friendly units without stealth, king=%v ally=%v existing=%v", mistKing.Statuses, ally.Statuses, alreadyStealth.Statuses)
	}

	phantom := NewCardInstance(baseCard(t, "1221109"), 0, 1)
	engine.triggerEffects(TriggerOnEnter, phantom, nil, nil)
	if phantom.Statuses[StatusStealth] != 3 || effectiveElementsGain(phantom)[model.ElementWater] != 2 {
		t.Fatalf("1221109 should enter with stealth3 and dynamic water load while stealthy, statuses=%v load=%v", phantom.Statuses, effectiveElementsGain(phantom))
	}
	phantom.Statuses[StatusStealth] = 0
	if effectiveElementsGain(phantom)[model.ElementWater] != 0 {
		t.Fatalf("1221109 dynamic water load should disappear without stealth, load=%v", effectiveElementsGain(phantom))
	}
	phantom.Statuses[StatusStealth] = 1
	phantom.Statuses[StatusPetrify] = 1
	if effectiveElementsGain(phantom)[model.ElementWater] != 0 {
		t.Fatalf("1221109 dynamic water load should disappear while petrified, load=%v", effectiveElementsGain(phantom))
	}
	shieldedEngine := setupReportedBugEngine(t)
	shieldedP0 := shieldedEngine.State.Players[0]
	shieldedP0.Shield = 1
	shieldedP0.Equipment[0] = NewCardInstance(baseCard(t, "2411101"), 0, 1)
	shieldedPhantom := placeUnit(baseCard(t, "1221109"), 0, 0, 0, shieldedEngine)
	shieldedPhantom.Statuses[StatusStealth] = 1
	shieldedPhantom.Statuses[StatusPetrify] = 1
	if shieldedEngine.effectiveElementsGain(shieldedPhantom)[model.ElementWater] != 2 {
		t.Fatalf("1221109 should keep dynamic water load when petrify is ineffective, load=%v", shieldedEngine.effectiveElementsGain(shieldedPhantom))
	}

	mage := placeUnit(baseCard(t, "1221102"), 0, 0, 1, engine)
	target := placeUnit(baseCard(t, "1021003"), 0, 1, 1, engine)
	if !cardHasActivePerTurn(mage) {
		t.Fatal("1221102 should expose a per-turn ability")
	}
	if err := globalRegistry.GetBehavior("1221102").(PerTurnAbility).OnPerTurn(&EffectContext{Engine: engine, Source: mage, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("1221102 per-turn: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "mist_mage_stealth" {
		t.Fatalf("1221102 should ask for a friendly unit, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, target.InstanceID)
	if target.Statuses[StatusStealth] != 2 {
		t.Fatalf("1221102 should grant stealth2 to selected ally, statuses=%v", target.Statuses)
	}
	p0.Units[0][1] = nil
	p0.Units[1][1] = nil

	sandworm := placeUnit(baseCard(t, "1421114"), 0, 0, 1, engine)
	sandworm.CurrentLife = 6
	engine.dealDamageWithExtra(sandworm, 1, 0, map[string]any{"damage_source": "effect", "attacker": 1})
	if sandworm.Statuses[StatusStealth] != 1 {
		t.Fatalf("1421114 should gain stealth1 after taking damage, statuses=%v", sandworm.Statuses)
	}

	promptEngine := setupReportedBugEngine(t)
	dancer := placeUnit(baseCard(t, "1221105"), 0, 0, 0, promptEngine)
	dancerTarget := placeUnit(baseCard(t, "1021004"), 0, 1, 0, promptEngine)
	promptEngine.triggerEffects(TriggerOnEnter, dancer, nil, nil)
	if promptEngine.State.PendingAction == nil || promptEngine.State.PendingAction.Type != "mist_dancer_stealth" {
		t.Fatalf("1221105 should ask for a companion target, pending=%+v", promptEngine.State.PendingAction)
	}
	resolvePendingSelection(t, promptEngine, 0, dancerTarget.InstanceID)
	if dancerTarget.Statuses[StatusStealth] != 2 {
		t.Fatalf("1221105 should grant stealth2 to selected companion, statuses=%v", dancerTarget.Statuses)
	}

	potionEngine := setupReportedBugEngine(t)
	potionTarget := placeUnit(baseCard(t, "1021005"), 0, 0, 0, potionEngine)
	potion := NewCardInstance(baseCard(t, "2021103"), 0, 1)
	if err := globalRegistry.GetBehavior("2021103").(OnUseItemBehavior).OnUseItem(&EffectContext{Engine: potionEngine, Source: potion, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("2021103 use item: %v", err)
	}
	if potionEngine.State.PendingAction == nil || potionEngine.State.PendingAction.Type != "mist_potion_stealth" {
		t.Fatalf("2021103 should ask for a companion target, pending=%+v", potionEngine.State.PendingAction)
	}
	resolvePendingSelection(t, potionEngine, 0, potionTarget.InstanceID)
	if potionTarget.Statuses[StatusStealth] != 2 {
		t.Fatalf("2021103 should grant stealth2 to selected companion, statuses=%v", potionTarget.Statuses)
	}
}

func TestRoyalConflictStealthTargetingAndDelayedSummon(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]

	frontStealth := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	frontStealth.Statuses[StatusStealth] = 1
	backStealth := placeUnit(baseCard(t, "1021002"), 1, 1, 2, engine)
	backStealth.Statuses[StatusStealth] = 1
	backVisible := placeUnit(baseCard(t, "1021003"), 1, 2, 2, engine)
	visibleFront := placeUnit(baseCard(t, "1021004"), 1, 2, 0, engine)

	undercurrent := readySkill(baseCard(t, "3221106"), 0)
	if err := engine.validateSpellTarget(0, undercurrent, SpellTarget{Type: "unit", Position: *backStealth.Position}); err != nil {
		t.Fatalf("3221106 should target stealth units regardless of row: %v", err)
	}
	if err := engine.validateSpellTarget(0, undercurrent, SpellTarget{Type: "unit", Position: *backVisible.Position}); err == nil {
		t.Fatal("3221106 should not get global range against ordinary non-stealth back-row units")
	}
	if got := engine.effectiveSpellPower(0, undercurrent, nil, SpellTarget{Type: "unit", Position: *backStealth.Position}); got != undercurrent.Card.Power+2 {
		t.Fatalf("3221106 should gain +2 power against stealth targets, got %d", got)
	}
	if got := engine.effectiveSpellPower(0, undercurrent, nil, SpellTarget{Type: "unit", Position: *backVisible.Position}); got != undercurrent.Card.Power {
		t.Fatalf("3221106 should not gain power against visible targets, got %d", got)
	}
	p1.Units[visibleFront.Position.Col][visibleFront.Position.Row] = nil

	waterEscape := readySkill(baseCard(t, "3221104"), 0)
	ally := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
	ownerID := 0
	if err := engine.validateSpellTarget(0, waterEscape, SpellTarget{Type: "unit", OwnerID: &ownerID, Position: *ally.Position}); err != nil {
		t.Fatalf("3221104 should target friendly non-stealth units: %v", err)
	}
	ally.Statuses[StatusStealth] = 1
	if err := engine.validateSpellTarget(0, waterEscape, SpellTarget{Type: "unit", OwnerID: &ownerID, Position: *ally.Position}); err == nil {
		t.Fatal("3221104 should reject units that already have stealth")
	}
	ally.Statuses[StatusPetrify] = 1
	if err := engine.validateSpellTarget(0, waterEscape, SpellTarget{Type: "unit", OwnerID: &ownerID, Position: *ally.Position}); err != nil {
		t.Fatalf("3221104 should allow targets whose stealth is disabled by petrify: %v", err)
	}
	ally.Statuses[StatusPetrify] = 0
	ally.Statuses[StatusStealth] = 0
	engine.resolveSpellHit(0, waterEscape, SpellTarget{Type: "unit", OwnerID: &ownerID, Position: *ally.Position}, nil, nil)
	if ally.Statuses[StatusStealth] != 2 {
		t.Fatalf("3221104 should grant stealth2 on hit, statuses=%v", ally.Statuses)
	}

	weaver := placeUnit(baseCard(t, "1321104"), 0, 2, 0, engine)
	engine.triggerEffects(TriggerOnEnter, weaver, nil, nil)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "mist_weaver_stealth" {
		t.Fatalf("1321104 should ask for a visible enemy target, pending=%+v", engine.State.PendingAction)
	}
	for _, candidate := range engine.State.PendingAction.Candidates {
		if candidate["instance_id"] == frontStealth.InstanceID || candidate["instance_id"] == backStealth.InstanceID {
			t.Fatalf("1321104 should not offer opposing stealth units, candidates=%+v", engine.State.PendingAction.Candidates)
		}
	}
	resolvePendingSelection(t, engine, 0, backVisible.InstanceID)
	if backVisible.Statuses[StatusStealth] != 2 {
		t.Fatalf("1321104 should grant stealth2 to selected enemy, statuses=%v", backVisible.Statuses)
	}

	if !cardHasActiveUltimate(NewCardInstance(baseCard(t, "4311102"), 0, 1)) {
		t.Fatal("4311102 should expose an ultimate ability")
	}
	fug := NewCardInstance(baseCard(t, "4311102"), 0, 1)
	if err := globalRegistry.GetBehavior("4311102").(UltimateAbility).OnUltimate(&EffectContext{Engine: engine, Source: fug, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("4311102 ultimate: %v", err)
	}
	next0 := NewCardInstance(baseCard(t, "1021005"), 0, 1)
	p0.Hand = append(p0.Hand, next0)
	for _, elem := range model.AllElements {
		p0.Elements[elem] = 9
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{"instance_id": next0.InstanceID, "col": float64(1), "row": float64(0)}}); err != nil {
		t.Fatalf("summon p0 next companion: %v", err)
	}
	if next0.Statuses[StatusStealth] != 2 || p0.NextCompanionStealth != 0 {
		t.Fatalf("4311102 should give p0 next summoned companion stealth2 once, statuses=%v pending=%d", next0.Statuses, p0.NextCompanionStealth)
	}
	next1 := NewCardInstance(baseCard(t, "1021006"), 1, 1)
	p1.Hand = append(p1.Hand, next1)
	for _, elem := range model.AllElements {
		p1.Elements[elem] = 9
	}
	engine.State.CurrentTurn = 1
	if err := engine.HandleAction(1, ActionMessage{Action: "summon", Data: map[string]any{"instance_id": next1.InstanceID, "col": float64(2), "row": float64(0)}}); err != nil {
		t.Fatalf("summon p1 next companion: %v", err)
	}
	if next1.Statuses[StatusStealth] != 2 || p1.NextCompanionStealth != 0 {
		t.Fatalf("4311102 should give p1 next summoned companion stealth2 once, statuses=%v pending=%d", next1.Statuses, p1.NextCompanionStealth)
	}

	pathEngine := setupReportedBugEngine(t)
	pathP0 := pathEngine.State.Players[0]
	pathP0.NextCompanionStealth = 2
	freeSummon := NewCardInstance(baseCard(t, "1021007"), 0, 1)
	pathP0.Hand = append(pathP0.Hand, freeSummon)
	if got := summonHandCompanionFree(&EffectContext{Engine: pathEngine, PlayerID: 0, OpponentID: 1}, func(card *CardInstance) bool {
		return card == freeSummon
	}); got != freeSummon {
		t.Fatalf("expected free hand summon path to return selected card, got %+v", got)
	}
	if freeSummon.Statuses[StatusStealth] != 2 || pathP0.NextCompanionStealth != 0 {
		t.Fatalf("free hand summon should consume next companion stealth, statuses=%v pending=%d", freeSummon.Statuses, pathP0.NextCompanionStealth)
	}

	pathP0.NextCompanionStealth = 2
	revived := NewCardInstance(baseCard(t, "1021008"), 0, 1)
	pathP0.Graveyard = append(pathP0.Graveyard, revived)
	if !pathEngine.reviveCompanionFromGraveyardWithLifeAtPosition(0, revived.InstanceID, 1, false, Position{Col: 1, Row: 0}) {
		t.Fatal("expected graveyard revive path to succeed")
	}
	if revived.Statuses[StatusStealth] != 2 || pathP0.NextCompanionStealth != 0 {
		t.Fatalf("graveyard revive should consume next companion stealth, statuses=%v pending=%d", revived.Statuses, pathP0.NextCompanionStealth)
	}
}

func TestRoyalConflictLightweightSpellAndItemEffects(t *testing.T) {
	t.Run("gospel discounts light skill use after a friendly light companion consumes", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		gospel := readySkill(baseCard(t, "3521101"), 0)
		p0.Skills[0] = gospel
		lightCompanion := placeUnit(baseCard(t, "1521104"), 0, 0, 0, engine)
		lightCompanion.IsHorizontal = false

		if got := engine.effectiveSkillUseCost(p0, gospel)[model.ElementLight]; got != gospel.Card.ElementsExpense[model.ElementLight] {
			t.Fatalf("unexpected base gospel use cost, got=%d card=%v", got, gospel.Card.ElementsExpense)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{"instance_id": lightCompanion.InstanceID}}); err != nil {
			t.Fatalf("consume light companion: %v", err)
		}
		if got := engine.effectiveSkillUseCost(p0, gospel)[model.ElementLight]; got != gospel.Card.ElementsExpense[model.ElementLight]-1 {
			t.Fatalf("gospel should reduce its light use cost this turn, got=%d modifiers=%+v", got, p0.TempModifiers)
		}
		engine.clearExpiredTemporaryModifiers(0)
		if got := engine.effectiveSkillUseCost(p0, gospel)[model.ElementLight]; got != gospel.Card.ElementsExpense[model.ElementLight] {
			t.Fatalf("gospel discount should expire at turn end, got=%d modifiers=%+v", got, p0.TempModifiers)
		}
	})

	t.Run("lingering frost scroll counts water spells already used this turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2221110"), 0, engine.State.TurnNumber)
		p0.SpellsCastThisTurn[model.ElementWater] = 2

		if got := engine.effectiveSpellPower(0, scroll, nil); got != scroll.Card.Power+6 {
			t.Fatalf("2221110 should gain +3 power per prior water spell, got=%d base=%d", got, scroll.Card.Power)
		}
	})

	t.Run("oracle scroll unity discounts itself for friendly light units", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2521112"), 0, engine.State.TurnNumber)
		placeUnit(baseCard(t, "1521104"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1521109"), 0, 1, 0, engine)

		if got := engine.effectiveCardPlayCost(p0, scroll)[model.ElementLight]; got != scroll.Card.ElementsCost[model.ElementLight]-2 {
			t.Fatalf("2521112 should cost -1 light per friendly light unit, got=%d cost=%v", got, engine.effectiveCardPlayCost(p0, scroll))
		}
	})

	t.Run("rotting erosion weakens enemy spells and advances mastery on friendly spell hit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		erosion := readySkill(baseCard(t, "3421106"), 0)
		erosion.Statuses[StatusMastery] = 2
		p0.Skills[0] = erosion
		enemyA := readySkill(baseCard(t, "3121001"), 1)
		enemyB := readySkill(baseCard(t, "3221003"), 1)
		p1.Skills[0] = enemyA
		p1.Skills[1] = enemyB
		otherSpell := readySkill(baseCard(t, "3121001"), 0)

		engine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, otherSpell, map[string]any{"attacker": 0, "spell_source": otherSpell})
		if enemyA.Statuses[StatusWeaken] != 0 || enemyB.Statuses[StatusWeaken] != 0 || erosion.Statuses[StatusMastery] != 2 {
			t.Fatalf("3421106 should ignore other friendly spell hits, a=%v b=%v mastery=%d", enemyA.Statuses, enemyB.Statuses, erosion.Statuses[StatusMastery])
		}

		engine.triggerEffects(TriggerOnSpellHit, erosion, nil, map[string]any{"attacker": 0, "spell_source": erosion})
		if enemyA.Statuses[StatusWeaken] != 1 || enemyB.Statuses[StatusWeaken] != 1 {
			t.Fatalf("3421106 should weaken all enemy spell instances by attack, a=%v b=%v", enemyA.Statuses, enemyB.Statuses)
		}
		if erosion.Statuses[StatusMastery] != 3 || erosion.PowerBonus != 1 || erosion.AttackBonus != 1 {
			t.Fatalf("3421106 should advance to mastery 3 and gain stats, mastery=%d powerBonus=%d attackBonus=%d", erosion.Statuses[StatusMastery], erosion.PowerBonus, erosion.AttackBonus)
		}
		engine.advanceMastery(erosion, 0, 3)
		if erosion.Statuses[StatusMastery] != 6 || erosion.PowerBonus != 2 || erosion.AttackBonus != 2 {
			t.Fatalf("3421106 should gain stats again at mastery 6, mastery=%d powerBonus=%d attackBonus=%d", erosion.Statuses[StatusMastery], erosion.PowerBonus, erosion.AttackBonus)
		}
	})

	t.Run("sky witch soland buffs drive and focus spells while restricting learned skill tags", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Hero = NewCardInstance(baseCard(t, "4311101"), 0, engine.State.TurnNumber)
		drive := NewCardInstance(baseCard(t, "3321101"), 0, engine.State.TurnNumber)
		creation := NewCardInstance(baseCard(t, "3121001"), 0, engine.State.TurnNumber)
		p0.SkillPool = []*CardInstance{creation, drive}
		p0.Elements = cloneElements(map[string]int{model.ElementFire: 9, model.ElementAir: 9})

		if got := engine.effectiveSpellPower(0, drive, nil); got != drive.Card.Power+1 {
			t.Fatalf("4311101 should give drive/focus spells +1 power, got=%d base=%d", got, drive.Card.Power)
		}
		if got := engine.effectiveSpellPower(0, creation, nil); got != creation.Card.Power {
			t.Fatalf("4311101 should not buff other spell tags, got=%d base=%d", got, creation.Card.Power)
		}
		if err := engine.handleLearnSkill(0, ActionMessage{Action: "learn_skill", Data: map[string]any{"instance_id": creation.InstanceID}}); err == nil {
			t.Fatal("4311101 should block learning non-drive/non-focus spells")
		}
		if err := engine.handleLearnSkill(0, ActionMessage{Action: "learn_skill", Data: map[string]any{"instance_id": drive.InstanceID}}); err != nil {
			t.Fatalf("4311101 should allow learning drive/focus spells: %v", err)
		}
	})

	t.Run("held breath buffs air spells until an extra draw happens", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3321107"), 0)
		airSpell := readySkill(baseCard(t, "3321101"), 0)
		fireSpell := readySkill(baseCard(t, "3121001"), 0)

		p0.DrawCountThisTurn = 1
		if got := engine.effectiveSpellPower(0, airSpell, nil); got != airSpell.Card.Power+1 {
			t.Fatalf("3321107 should buff air spells before extra draws, got=%d base=%d", got, airSpell.Card.Power)
		}
		if got := engine.effectiveSpellPower(0, fireSpell, nil); got != fireSpell.Card.Power {
			t.Fatalf("3321107 should not buff non-air spells, got=%d base=%d", got, fireSpell.Card.Power)
		}
		p0.DrawCountThisTurn = 2
		if got := engine.effectiveSpellPower(0, airSpell, nil); got != airSpell.Card.Power {
			t.Fatalf("3321107 should stop after extra draws, got=%d base=%d", got, airSpell.Card.Power)
		}
	})

	t.Run("devotion contract triggers once after a friendly atonement spell cast", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		contract := NewCardInstance(baseCard(t, "2621104"), 0, engine.State.TurnNumber)
		p0.Equipment[0] = contract
		p0.Hero = NewCardInstance(baseCard(t, "4611101"), 0, engine.State.TurnNumber)
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, engine.State.TurnNumber)}
		heroLife := p0.Hero.CurrentLife
		spell := readySkill(baseCard(t, "3621103"), 0)

		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, spell, map[string]any{"cast_player": 0})
		if p0.Hero.CurrentLife != heroLife-1 || len(p0.Hand) != 1 || contract.UsedThisTurn != 1 {
			t.Fatalf("2621104 should damage hero, draw, and spend trigger once, hero=%d hand=%d used=%d", p0.Hero.CurrentLife, len(p0.Hand), contract.UsedThisTurn)
		}
		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, spell, map[string]any{"cast_player": 0})
		if p0.Hero.CurrentLife != heroLife-1 || len(p0.Hand) != 1 || contract.UsedThisTurn != 1 {
			t.Fatalf("2621104 should not trigger more than once per turn, hero=%d hand=%d used=%d", p0.Hero.CurrentLife, len(p0.Hand), contract.UsedThisTurn)
		}
	})

	t.Run("sky city zenith stone damages and stuns the drawing player's front row every fifth draw", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		stone := NewCardInstance(baseCard(t, "2311101"), 0, engine.State.TurnNumber)
		p0.Equipment[0] = stone
		frontLeft := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		frontRight := placeUnit(baseCard(t, "1021002"), 1, 2, 0, engine)
		backUnit := placeUnit(baseCard(t, "1021003"), 1, 1, 1, engine)
		frontLife := frontLeft.CurrentLife
		backLife := backUnit.CurrentLife
		p1.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 1, engine.State.TurnNumber),
			NewCardInstance(baseCard(t, "1021002"), 1, engine.State.TurnNumber),
			NewCardInstance(baseCard(t, "1021003"), 1, engine.State.TurnNumber),
			NewCardInstance(baseCard(t, "1021004"), 1, engine.State.TurnNumber),
			NewCardInstance(baseCard(t, "1021005"), 1, engine.State.TurnNumber),
		}

		engine.drawCards(1, 4)
		if stone.Statuses[skyCityZenithStoneMarkerStatus] != 4 || frontLeft.CurrentLife != frontLife || frontLeft.Statuses[StatusStun] != 0 {
			t.Fatalf("2311101 should only collect markers before the fifth draw, markers=%d life=%d stun=%d", stone.Statuses[skyCityZenithStoneMarkerStatus], frontLeft.CurrentLife, frontLeft.Statuses[StatusStun])
		}
		engine.drawCards(1, 1)
		if stone.Statuses[skyCityZenithStoneMarkerStatus] != 0 {
			t.Fatalf("2311101 should remove all markers at five, statuses=%v", stone.Statuses)
		}
		if frontLeft.CurrentLife != frontLife-1 || frontRight.CurrentLife != frontRight.Card.Life-1 || frontLeft.Statuses[StatusStun] != 1 || frontRight.Statuses[StatusStun] != 1 {
			t.Fatalf("2311101 should damage and stun front row, left=%d/%v right=%d/%v", frontLeft.CurrentLife, frontLeft.Statuses, frontRight.CurrentLife, frontRight.Statuses)
		}
		if backUnit.CurrentLife != backLife || backUnit.Statuses[StatusStun] != 0 {
			t.Fatalf("2311101 should not hit back row, life=%d statuses=%v", backUnit.CurrentLife, backUnit.Statuses)
		}
	})

	t.Run("blood gu tracks hero damage then sacrifices for a current turn spell power buff", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Hero = NewCardInstance(baseCard(t, "4611101"), 0, engine.State.TurnNumber)
		p0.Hero.CurrentLife = 20
		gu := NewCardInstance(baseCard(t, "2621103"), 0, engine.State.TurnNumber)
		p0.Equipment[0] = gu
		spell := readySkill(baseCard(t, "3121001"), 0)
		basePower := engine.effectiveSpellPower(0, spell, nil)

		engine.dealDamageWithExtra(p0.Hero, 4, 0, map[string]any{"damage_source": "test", "attacker": 1})
		if gu.Statuses[bloodGuMarkerStatus] != 4 {
			t.Fatalf("2621103 should gain markers equal hero damage, statuses=%v", gu.Statuses)
		}
		other := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		engine.dealDamageWithExtra(other, 2, 0, map[string]any{"damage_source": "test", "attacker": 1})
		if gu.Statuses[bloodGuMarkerStatus] != 4 {
			t.Fatalf("2621103 should ignore non-hero damage, statuses=%v", gu.Statuses)
		}
		engine.dealDamageWithExtra(p0.Hero, 4, 0, map[string]any{"damage_source": "test", "attacker": 1})
		if gu.Statuses[bloodGuMarkerStatus] != 6 {
			t.Fatalf("2621103 markers should cap at six, statuses=%v", gu.Statuses)
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{"instance_id": gu.InstanceID}}); err != nil {
			t.Fatalf("2621103 active ability should sacrifice for spell power: %v", err)
		}
		if p0.Equipment[0] != nil || countCardsByNumber(p0.Graveyard, "2621103") != 1 {
			t.Fatalf("2621103 should be sacrificed from equipment, equipment=%v grave=%v", p0.Equipment[0], cardsToInfo(p0.Graveyard))
		}
		if got := engine.effectiveSpellPower(0, spell, nil); got != basePower+3 {
			t.Fatalf("2621103 should add +3 spell power from six markers, got=%d base=%d", got, basePower)
		}
		engine.clearExpiredTemporaryModifiers(0)
		if got := engine.effectiveSpellPower(0, spell, nil); got != basePower {
			t.Fatalf("2621103 spell power buff should expire after current turn, got=%d base=%d", got, basePower)
		}
	})

	t.Run("council judgment hammer shuffles marks after enemy spell attacks once per turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		hammer := NewCardInstance(baseCard(t, "2521108"), 0, engine.State.TurnNumber)
		p0.Equipment[0] = hammer
		enemyAttackSpell := readySkill(baseCard(t, "3121001"), 1)
		enemySorcery := readySkill(baseCard(t, "3021001"), 1)

		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, enemySorcery, map[string]any{"cast_player": 1})
		if countCardsByNumber(p1.Deck, "2001102") != 0 || hammer.UsedThisTurn != 0 {
			t.Fatalf("2521108 should ignore enemy sorceries, deck=%v used=%d", cardsToInfo(p1.Deck), hammer.UsedThisTurn)
		}
		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, enemyAttackSpell, map[string]any{"cast_player": 1})
		if countCardsByNumber(p1.Deck, "2001102") != 3 || hammer.UsedThisTurn != 1 {
			t.Fatalf("2521108 should shuffle three marks into enemy deck once, deck=%v used=%d", cardsToInfo(p1.Deck), hammer.UsedThisTurn)
		}
		engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, enemyAttackSpell, map[string]any{"cast_player": 1})
		if countCardsByNumber(p1.Deck, "2001102") != 3 || hammer.UsedThisTurn != 1 {
			t.Fatalf("2521108 should only trigger once per turn, deck=%v used=%d", cardsToInfo(p1.Deck), hammer.UsedThisTurn)
		}
	})

	t.Run("red agate chalice gives the completed light artifact set extra load", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scepter := NewCardInstance(baseCard(t, "2521006"), 0, engine.State.TurnNumber)
		lamp := NewCardInstance(baseCard(t, "2521007"), 0, engine.State.TurnNumber)
		chalice := NewCardInstance(baseCard(t, "2521103"), 0, engine.State.TurnNumber)
		p0.Equipment[0] = scepter
		p0.Equipment[1] = lamp
		p0.Equipment[2] = chalice

		if got := engine.effectiveElementsGain(scepter)[model.ElementLight]; got != scepter.Card.ElementsGain[model.ElementLight]+1 {
			t.Fatalf("2521103 should add light load to green jade scepter in the complete set, load=%v", engine.effectiveElementsGain(scepter))
		}
		if got := engine.effectiveElementsGain(lamp)[model.ElementLight]; got != lamp.Card.ElementsGain[model.ElementLight]+1 {
			t.Fatalf("2521103 should add light load to blue crystal lamp in the complete set, load=%v", engine.effectiveElementsGain(lamp))
		}
		if got := engine.effectiveElementsGain(chalice)[model.ElementLight]; got != chalice.Card.ElementsGain[model.ElementLight]+1 {
			t.Fatalf("2521103 should add light load to itself in the complete set, load=%v", engine.effectiveElementsGain(chalice))
		}

		p0.Equipment[1] = nil
		if got := engine.effectiveElementsGain(scepter)[model.ElementLight]; got != scepter.Card.ElementsGain[model.ElementLight] {
			t.Fatalf("2521103 should not add load before the set is complete, load=%v", engine.effectiveElementsGain(scepter))
		}
	})

	t.Run("quick ice bullet discounts the next consumable item or spell this turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		bullet := NewCardInstance(baseCard(t, "2221109"), 0, engine.State.TurnNumber)
		item := NewCardInstance(baseCard(t, "2221009"), 0, engine.State.TurnNumber)
		spell := readySkill(baseCard(t, "3201002"), 0)

		engine.triggerEffects(TriggerOnUseItem, bullet, nil, nil)
		if cost := engine.effectiveCardPlayCost(p0, item); cost[model.ElementWater] != item.Card.ElementsCost[model.ElementWater]-3 {
			t.Fatalf("2221109 should discount the next consumable item's water cost, cost=%v", cost)
		}
		engine.notifyCardPlayCostPaid(p0, item)
		if cost := engine.effectiveSkillUseCost(p0, spell); cost[model.ElementWater] != spell.Card.ElementsExpense[model.ElementWater] {
			t.Fatalf("2221109 item discount should be consumed after the next consumable, cost=%v", cost)
		}

		engine.triggerEffects(TriggerOnUseItem, bullet, nil, nil)
		if cost := engine.effectiveSkillUseCost(p0, spell); cost[model.ElementWater] != spell.Card.ElementsExpense[model.ElementWater]-3 {
			t.Fatalf("2221109 should discount the next spell's water use cost, cost=%v", cost)
		}
		engine.consumeNextSkillUseModifiers(p0, spell)
		if cost := engine.effectiveCardPlayCost(p0, item); cost[model.ElementWater] != item.Card.ElementsCost[model.ElementWater] {
			t.Fatalf("2221109 spell discount should be consumed after the next spell, cost=%v", cost)
		}

		engine.triggerEffects(TriggerOnUseItem, bullet, nil, nil)
		if cost := engine.effectiveSkillUseCostForPurpose(p0, spell, skillPurposeAttackBoost); cost[model.ElementWater] != spell.Card.ElementsExpense[model.ElementWater]-3 {
			t.Fatalf("2221109 should also discount a spell used as boost, cost=%v", cost)
		}
		engine.consumeNextSkillUseModifiersForPurpose(p0, spell, skillPurposeAttackBoost)
		if cost := engine.effectiveSkillUseCost(p0, spell); cost[model.ElementWater] != spell.Card.ElementsExpense[model.ElementWater] {
			t.Fatalf("2221109 boost discount should be consumed after the boosted spell use, cost=%v", cost)
		}

		learnedSpell := NewCardInstance(baseCard(t, "3221003"), 0, engine.State.TurnNumber)
		engine.triggerEffects(TriggerOnUseItem, bullet, nil, nil)
		if cost := engine.effectiveSkillLearnCost(p0, learnedSpell); cost[model.ElementWater] != learnedSpell.Card.ElementsCost[model.ElementWater]-3 {
			t.Fatalf("2221109 should discount the next spell learn cost, cost=%v", cost)
		}
		engine.notifyCardPlayCostPaid(p0, learnedSpell)
		if cost := engine.effectiveSkillUseCost(p0, spell); cost[model.ElementWater] != spell.Card.ElementsExpense[model.ElementWater] {
			t.Fatalf("2221109 learn discount should be consumed after learning a spell, cost=%v", cost)
		}
	})

	t.Run("last stand light requires fewer friendly units and scales from the best light companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		lastStand := readySkill(baseCard(t, "3511102"), 0)
		lightCompanion := placeUnit(baseCard(t, "1521001"), 0, 0, 0, engine)
		lightCompanion.CurrentLife = 5
		engine.addElementsGainBonus(lightCompanion, 0, model.ElementLight, 2, lastStand)
		placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)

		if err := engine.validateSkillUsePermissionModifiers(lastStand, skillPurposeAttack); err != nil {
			t.Fatalf("3511102 should be usable while friendly units are fewer: %v", err)
		}
		wantBonus := lightCompanion.CurrentLife + totalElementCost(engine.effectiveElementsGain(lightCompanion))
		if got := engine.effectiveSpellPower(0, lastStand, nil); got != lastStand.Card.Power+wantBonus {
			t.Fatalf("3511102 should gain power from the best light companion, got=%d base=%d bonus=%d", got, lastStand.Card.Power, wantBonus)
		}

		placeUnit(baseCard(t, "1021003"), 0, 1, 0, engine)
		if err := engine.validateSkillUsePermissionModifiers(lastStand, skillPurposeAttack); err == nil {
			t.Fatal("3511102 should require fewer friendly units than the opponent")
		}
		if err := engine.validateSkillUsePermissionModifiers(lastStand, skillPurposeDefenseBoost); err == nil {
			t.Fatal("3511102 restriction should also block boost use")
		}
	})

	t.Run("collector draws after equipping and gains arcane after using a consumable once each turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		collector := placeUnit(baseCard(t, "1011101"), 0, 0, 0, engine)
		equipmentA := NewCardInstance(baseCard(t, "2521006"), 0, engine.State.TurnNumber)
		equipmentB := NewCardInstance(baseCard(t, "2521007"), 0, engine.State.TurnNumber)
		consumableA := NewCardInstance(baseCard(t, "2221109"), 0, engine.State.TurnNumber)
		consumableB := NewCardInstance(baseCard(t, "2221109"), 0, engine.State.TurnNumber)
		p0.Hand = []*CardInstance{equipmentA, equipmentB, consumableA, consumableB}
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, engine.State.TurnNumber),
			NewCardInstance(baseCard(t, "1021002"), 0, engine.State.TurnNumber),
		}
		setAllElements(p0, 99)

		if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{"instance_id": equipmentA.InstanceID}}); err != nil {
			t.Fatalf("equip first item: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{"instance_id": equipmentB.InstanceID}}); err != nil {
			t.Fatalf("equip second item: %v", err)
		}
		if len(p0.Hand) != 3 || collector.Statuses[collectorEquipTriggeredTurnStatus] != engine.State.TurnNumber {
			t.Fatalf("1011101 should draw once after friendly equipment enters, hand=%d statuses=%v", len(p0.Hand), collector.Statuses)
		}

		beforeArcane := p0.Elements[model.ElementArcane]
		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": consumableA.InstanceID}}); err != nil {
			t.Fatalf("use first consumable: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": consumableB.InstanceID}}); err != nil {
			t.Fatalf("use second consumable: %v", err)
		}
		if p0.Elements[model.ElementArcane] != beforeArcane+1 || collector.Statuses[collectorItemTriggeredTurnStatus] != engine.State.TurnNumber {
			t.Fatalf("1011101 should gain one arcane once after friendly consumable use, elements=%v statuses=%v", p0.Elements, collector.Statuses)
		}
	})

	t.Run("council consul shuffles marks when the opponent summons companions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		placeUnit(baseCard(t, "1521111"), 0, 0, 0, engine)
		enemy := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		enemyHero := NewCardInstance(baseCard(t, "4611101"), 1, engine.State.TurnNumber)

		engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, 0, enemy, map[string]any{"entered_player": 1})
		if countCardsByNumber(p1.Deck, "2001102") != 3 {
			t.Fatalf("1521111 should shuffle three marks into the summoning opponent deck, deck=%v", cardsToInfo(p1.Deck))
		}
		engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, 0, enemyHero, map[string]any{"entered_player": 1})
		if countCardsByNumber(p1.Deck, "2001102") != 3 {
			t.Fatalf("1521111 should ignore enemy heroes, deck=%v", cardsToInfo(p1.Deck))
		}
		own := placeUnit(baseCard(t, "1021002"), 0, 1, 0, engine)
		engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, 0, own, map[string]any{"entered_player": 0})
		if countCardsByNumber(p1.Deck, "2001102") != 3 || len(p0.Deck) != 0 {
			t.Fatalf("1521111 should ignore friendly summons, ownDeck=%v enemyDeck=%v", cardsToInfo(p0.Deck), cardsToInfo(p1.Deck))
		}
	})

	t.Run("pure spirit weakens friendly spells when non-arcane cards enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Hero = NewCardInstance(baseCard(t, "4011101"), 0, engine.State.TurnNumber)
		skill := readySkill(baseCard(t, "3121001"), 0)
		p0.Skills[0] = skill
		host := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		bound := readySkill(baseCard(t, "3221003"), 0)
		host.BoundSkills = []*CardInstance{bound}
		arcaneCard := NewCardInstance(baseCard(t, "1021112"), 0, engine.State.TurnNumber)
		fireCard := NewCardInstance(baseCard(t, "1121001"), 0, engine.State.TurnNumber)
		p0.Hand = []*CardInstance{arcaneCard, fireCard}
		setAllElements(p0, 99)

		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{"instance_id": arcaneCard.InstanceID, "col": float64(1), "row": float64(0)}}); err != nil {
			t.Fatalf("summon arcane card: %v", err)
		}
		if skill.Statuses[StatusWeaken] != 0 || bound.Statuses[StatusWeaken] != 0 {
			t.Fatalf("4011101 should ignore arcane card entry, skill=%v bound=%v", skill.Statuses, bound.Statuses)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{"instance_id": fireCard.InstanceID, "col": float64(2), "row": float64(0)}}); err != nil {
			t.Fatalf("summon non-arcane card: %v", err)
		}
		if skill.Statuses[StatusWeaken] != 2 || bound.Statuses[StatusWeaken] != 2 {
			t.Fatalf("4011101 should weaken all friendly spell instances after non-arcane entry, skill=%v bound=%v", skill.Statuses, bound.Statuses)
		}
		engine.notifyCardEntered(1, fireCard, nil)
		if skill.Statuses[StatusWeaken] != 2 || bound.Statuses[StatusWeaken] != 2 {
			t.Fatalf("4011101 should ignore opponent card entry, skill=%v bound=%v", skill.Statuses, bound.Statuses)
		}
	})

	t.Run("set counter card entry triggers pure spirit but not collector equipment draw", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Hero = NewCardInstance(baseCard(t, "4011101"), 0, engine.State.TurnNumber)
		placeUnit(baseCard(t, "1011101"), 0, 0, 0, engine)
		skill := readySkill(baseCard(t, "3121001"), 0)
		p0.Skills[0] = skill
		counter := NewCardInstance(baseCard(t, "2121002"), 0, engine.State.TurnNumber)
		p0.Hand = []*CardInstance{counter}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, engine.State.TurnNumber)}

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": counter.InstanceID}}); err != nil {
			t.Fatalf("set counter: %v", err)
		}
		if p0.Equipment[0] != counter || !counter.IsSetCounter {
			t.Fatalf("counter should be set into equipment, equipment=%v set=%v", cardToInfo(p0.Equipment[0]), counter.IsSetCounter)
		}
		if skill.Statuses[StatusWeaken] != 2 {
			t.Fatalf("4011101 should see a non-arcane counter enter the equipment area, statuses=%v", skill.Statuses)
		}
		if len(p0.Hand) != 0 {
			t.Fatalf("1011101 should not treat setting a counter as equipping, hand=%v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("retribution gains attack from hero damage this turn and last turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Hero = NewCardInstance(baseCard(t, "4611101"), 0, engine.State.TurnNumber)
		p0.Hero.CurrentLife = 20
		retribution := readySkill(baseCard(t, "3621102"), 0)

		engine.dealDamageWithExtra(p0.Hero, 2, 0, map[string]any{"damage_source": "test", "attacker": 1})
		if p0.HeroDamageTakenThisTurn != 2 || p0.HeroDamageTakenLastTurn != 0 {
			t.Fatalf("hero damage should be tracked this turn, this=%d last=%d", p0.HeroDamageTakenThisTurn, p0.HeroDamageTakenLastTurn)
		}
		if got := engine.effectiveSpellDamage(0, retribution, retribution.Card.Attack, nil); got != retribution.Card.Attack+2 {
			t.Fatalf("3621102 should gain attack from this turn hero damage, got=%d base=%d", got, retribution.Card.Attack)
		}

		engine.rollFriendlyUnitDamageHistory()
		if p0.HeroDamageTakenThisTurn != 0 || p0.HeroDamageTakenLastTurn != 2 {
			t.Fatalf("hero damage should roll into last turn history, this=%d last=%d", p0.HeroDamageTakenThisTurn, p0.HeroDamageTakenLastTurn)
		}
		engine.dealDamageWithExtra(p0.Hero, 1, 0, map[string]any{"damage_source": "test", "attacker": 1})
		if got := engine.effectiveSpellDamage(0, retribution, retribution.Card.Attack, nil); got != retribution.Card.Attack+3 {
			t.Fatalf("3621102 should add this and last turn hero damage, got=%d base=%d", got, retribution.Card.Attack)
		}

		engine.rollFriendlyUnitDamageHistory()
		engine.rollFriendlyUnitDamageHistory()
		if got := engine.effectiveSpellDamage(0, retribution, retribution.Card.Attack, nil); got != retribution.Card.Attack {
			t.Fatalf("3621102 bonus should expire after the two-turn damage window, got=%d base=%d", got, retribution.Card.Attack)
		}
	})

	t.Run("panacea p damages heals and draws through the consumable item flow", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		potion := NewCardInstance(baseCard(t, "2521107"), 0, engine.State.TurnNumber)
		p0.Hand = []*CardInstance{potion}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, engine.State.TurnNumber)}
		setAllElements(p0, 99)
		friendly := placeUnit(baseCard(t, "1021002"), 0, 0, 0, engine)
		friendly.CurrentLife -= 2
		enemy := placeUnit(baseCard(t, "1021003"), 1, 0, 0, engine)
		friendlyLife := friendly.CurrentLife
		enemyLife := enemy.CurrentLife

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": potion.InstanceID}}); err != nil {
			t.Fatalf("use 2521107: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "panacea_p_damage" {
			t.Fatalf("2521107 should first prompt for a damage target, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, enemy.InstanceID)
		if enemy.CurrentLife != enemyLife-1 {
			t.Fatalf("2521107 should deal 1 damage before healing, got=%d want=%d", enemy.CurrentLife, enemyLife-1)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "panacea_p_heal" {
			t.Fatalf("2521107 should then prompt for a heal target, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, friendly.InstanceID)
		if friendly.CurrentLife != friendlyLife+1 {
			t.Fatalf("2521107 should heal 1 after damage, got=%d want=%d", friendly.CurrentLife, friendlyLife+1)
		}
		if len(p0.Hand) != 1 || p0.Hand[0].Card.Number != "1021001" || len(p0.Deck) != 0 {
			t.Fatalf("2521107 should draw one card after resolving, hand=%v deck=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Deck))
		}
		if countCardsByNumber(p0.Graveyard, "2521107") != 1 {
			t.Fatalf("2521107 should be in graveyard after use, grave=%v", cardsToInfo(p0.Graveyard))
		}
	})

	t.Run("offering torch exiles one fire spell to permanently buff another", func(t *testing.T) {
		blockedEngine := setupReportedBugEngine(t)
		blockedP0 := blockedEngine.State.Players[0]
		blockedTorch := NewCardInstance(baseCard(t, "2121110"), 0, blockedEngine.State.TurnNumber)
		blockedP0.Hand = []*CardInstance{blockedTorch}
		blockedP0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		setAllElements(blockedP0, 99)
		if err := blockedEngine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": blockedTorch.InstanceID}}); err == nil {
			t.Fatal("2121110 should require at least two friendly fire spells")
		}

		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		torch := NewCardInstance(baseCard(t, "2121110"), 0, engine.State.TurnNumber)
		p0.Hand = []*CardInstance{torch}
		exiled := readySkill(baseCard(t, "3121001"), 0)
		exiled.PowerBonus = 1
		exiled.AttackBonus = 1
		target := readySkill(baseCard(t, "3121002"), 0)
		host := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		boundTarget := readySkill(baseCard(t, "3121003"), 0)
		host.BoundSkills = []*CardInstance{boundTarget}
		p0.Skills[0] = exiled
		p0.Skills[1] = target
		setAllElements(p0, 99)
		wantPowerBonus := max(exiled.Card.Power+exiled.PowerBonus, 0)
		wantAttackBonus := max(exiled.Card.Attack+exiled.AttackBonus, 0)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": torch.InstanceID}}); err != nil {
			t.Fatalf("use 2121110: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "offering_torch_exile" {
			t.Fatalf("2121110 should ask which fire spell to exile, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, exiled.InstanceID)
		if len(p0.Exile) != 1 || p0.Exile[0] != exiled || p0.Skills[0] != nil {
			t.Fatalf("2121110 should exile the selected fire spell, exile=%v skills=%v", cardsToInfo(p0.Exile), cardsToInfo(p0.Skills[:]))
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "offering_torch_buff" {
			t.Fatalf("2121110 should ask which other fire spell to buff, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, boundTarget.InstanceID)
		if boundTarget.PowerBonus != wantPowerBonus || boundTarget.AttackBonus != wantAttackBonus {
			t.Fatalf("2121110 should permanently add exiled spell stats to bound spells, powerBonus=%d attackBonus=%d want=%d/%d", boundTarget.PowerBonus, boundTarget.AttackBonus, wantPowerBonus, wantAttackBonus)
		}
		if target.PowerBonus != 0 || target.AttackBonus != 0 {
			t.Fatalf("2121110 should only buff the selected spell, target=%d/%d", target.PowerBonus, target.AttackBonus)
		}
		if countCardsByNumber(p0.Graveyard, "2121110") != 1 {
			t.Fatalf("2121110 should go to graveyard after use, grave=%v", cardsToInfo(p0.Graveyard))
		}
	})

	t.Run("lavafort ashes exiles a fire skill to search a higher cost fire card with discount", func(t *testing.T) {
		blockedEngine := setupReportedBugEngine(t)
		blockedP0 := blockedEngine.State.Players[0]
		blockedAshes := NewCardInstance(baseCard(t, "2121101"), 0, blockedEngine.State.TurnNumber)
		blockedSkill := readySkill(baseCard(t, "3121001"), 0)
		blockedP0.Hand = []*CardInstance{blockedAshes}
		blockedP0.SkillPool = []*CardInstance{blockedSkill}
		blockedP0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1121001"), 0, blockedEngine.State.TurnNumber)}
		setAllElements(blockedP0, 99)
		if err := blockedEngine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": blockedAshes.InstanceID}}); err == nil {
			t.Fatal("2121101 should require a higher-cost fire card in deck")
		}

		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		ashes := NewCardInstance(baseCard(t, "2121101"), 0, engine.State.TurnNumber)
		source := readySkill(baseCard(t, "3121001"), 0)
		target := NewCardInstance(baseCard(t, "1121114"), 0, engine.State.TurnNumber)
		tooCheap := NewCardInstance(baseCard(t, "1121001"), 0, engine.State.TurnNumber)
		p0.Hand = []*CardInstance{ashes}
		p0.SkillPool = []*CardInstance{source}
		p0.Deck = []*CardInstance{tooCheap, target}
		setAllElements(p0, 99)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": ashes.InstanceID}}); err != nil {
			t.Fatalf("use 2121101: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "lavafort_ashes_exile_fire_skill" {
			t.Fatalf("2121101 should ask which fire skill to exile, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, source.InstanceID)
		if len(p0.Exile) != 1 || p0.Exile[0] != source || len(p0.SkillPool) != 0 {
			t.Fatalf("2121101 should exile selected fire skill from field or pool, exile=%v pool=%v", cardsToInfo(p0.Exile), cardsToInfo(p0.SkillPool))
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "lavafort_ashes_search_fire_card" || len(engine.State.PendingAction.Candidates) != 1 {
			t.Fatalf("2121101 should offer only higher-cost fire deck cards, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if len(p0.Hand) != 1 || p0.Hand[0] != target {
			t.Fatalf("2121101 should search selected fire card to hand, hand=%v", cardsToInfo(p0.Hand))
		}
		if target.Statuses["入场费用"+model.ElementFire+"-1"] != 1 {
			t.Fatalf("2121101 should give searched card fire entry discount, statuses=%v", target.Statuses)
		}
		if cost := engine.effectiveCardPlayCost(p0, target); cost[model.ElementFire] != max(target.Card.ElementsCost[model.ElementFire]-1, 0) {
			t.Fatalf("2121101 discount should affect entry cost, cost=%v base=%v", cost, target.Card.ElementsCost)
		}
		if countCardsByNumber(p0.Graveyard, "2121101") != 1 {
			t.Fatalf("2121101 should go to graveyard after use, grave=%v", cardsToInfo(p0.Graveyard))
		}
	})

	t.Run("claw of erebos requires weakened enemy spells and then weakens up to three", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		claw := NewCardInstance(baseCard(t, "3611102"), 0, engine.State.TurnNumber)
		p0.SkillPool = []*CardInstance{claw}
		p0.Elements = cloneElements(map[string]int{model.ElementShadow: 9})
		enemyA := readySkill(baseCard(t, "3121001"), 1)
		enemyB := readySkill(baseCard(t, "3221003"), 1)
		enemyC := readySkill(baseCard(t, "3321101"), 1)
		p1.Skills[0] = enemyA
		p1.Skills[1] = enemyB
		p1.Skills[2] = enemyC
		enemyA.Statuses[StatusWeaken] = 2

		if err := engine.handleLearnSkill(0, ActionMessage{Action: "learn_skill", Data: map[string]any{"instance_id": claw.InstanceID}}); err == nil {
			t.Fatal("3611102 should require at least three weakened enemy spell layers to learn")
		}
		enemyB.Statuses[StatusWeaken] = 1
		if err := engine.handleLearnSkill(0, ActionMessage{Action: "learn_skill", Data: map[string]any{"instance_id": claw.InstanceID}}); err != nil {
			t.Fatalf("3611102 should learn once enemy spell weaken layers reach three: %v", err)
		}
		if got := engine.effectiveSpellPower(0, claw, nil); got != claw.Card.Power+3 {
			t.Fatalf("3611102 should gain power for enemy weakened spell layers, got=%d base=%d", got, claw.Card.Power)
		}

		engine.triggerEffects(TriggerOnSpellCast, claw, nil, map[string]any{"cast_player": 0})
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "claw_of_erebos_weaken" {
			t.Fatalf("3611102 should prompt to weaken enemy spells after use, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, enemyA.InstanceID, enemyB.InstanceID, enemyC.InstanceID)
		if enemyA.Statuses[StatusWeaken] != 3 || enemyB.Statuses[StatusWeaken] != 2 || enemyC.Statuses[StatusWeaken] != 1 {
			t.Fatalf("3611102 should weaken up to three different enemy spells, a=%v b=%v c=%v", enemyA.Statuses, enemyB.Statuses, enemyC.Statuses)
		}
	})
}

func TestRoyalConflictStrictPaymentCards(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]

	pureBody := NewCardInstance(baseCard(t, "1021112"), 0, 1)
	playCost := engine.effectiveCardPlayCost(p0, pureBody)
	p0.Elements[model.ElementFire] = 4
	if engine.canPayCostForCardAction(p0, pureBody, playCost, playCost, paymentPurposePlay, ActionMessage{}) {
		t.Fatal("1021112 should not allow non-arcane elements for its strict entry cost")
	}
	p0.Elements[model.ElementFire] = 0
	p0.Elements[model.ElementArcane] = 4
	if !engine.canPayCostForCardAction(p0, pureBody, playCost, playCost, paymentPurposePlay, ActionMessage{}) {
		t.Fatal("1021112 should allow its strict entry cost to be paid with arcane")
	}

	absolutePure := NewCardInstance(baseCard(t, "3011101"), 0, 1)
	learnCost := engine.effectiveSkillLearnCost(p0, absolutePure)
	p0.Elements = cloneElements(map[string]int{model.ElementWater: 11})
	if engine.canPayCostForCardAction(p0, absolutePure, learnCost, learnCost, paymentPurposeLearn, ActionMessage{}) {
		t.Fatal("3011101 should require strict arcane for learning")
	}
	p0.Elements = cloneElements(map[string]int{model.ElementArcane: 11})
	if !engine.canPayCostForCardAction(p0, absolutePure, learnCost, learnCost, paymentPurposeLearn, ActionMessage{}) {
		t.Fatal("3011101 should allow learning with strict arcane")
	}
	useCost := engine.effectiveSkillUseCost(p0, absolutePure)
	totalUseCost := mergeElementCosts(useCost, map[string]int{model.ElementFire: 1})
	p0.Elements = cloneElements(map[string]int{model.ElementArcane: 7, model.ElementFire: 1})
	if !engine.canPayCostForCardAction(p0, absolutePure, useCost, totalUseCost, paymentPurposeUse, ActionMessage{Data: map[string]any{
		"payment": map[string]any{model.ElementArcane: float64(7), model.ElementFire: float64(1)},
	}}) {
		t.Fatal("3011101 strict use cost should allow separate non-strict boost payment")
	}
	p0.Elements = cloneElements(map[string]int{model.ElementWater: 7, model.ElementFire: 1})
	if engine.canPayCostForCardAction(p0, absolutePure, useCost, totalUseCost, paymentPurposeUse, ActionMessage{Data: map[string]any{
		"payment": map[string]any{model.ElementWater: float64(7), model.ElementFire: float64(1)},
	}}) {
		t.Fatal("3011101 should reject non-arcane payment for its own use cost")
	}

	timeCycle := NewCardInstance(baseCard(t, "3411101"), 0, 1)
	timeCost := engine.effectiveSkillLearnCost(p0, timeCycle)
	p0.Elements = cloneElements(map[string]int{model.ElementArcane: 2})
	if engine.canPayCostForCardAction(p0, timeCycle, timeCost, timeCost, paymentPurposeLearn, ActionMessage{}) {
		t.Fatal("3411101 should require earth for its earth component")
	}
	p0.Elements = cloneElements(map[string]int{model.ElementEarth: 1, model.ElementFire: 1})
	if engine.canPayCostForCardAction(p0, timeCycle, timeCost, timeCost, paymentPurposeLearn, ActionMessage{}) {
		t.Fatal("3411101 should require arcane for its arcane component")
	}
	p0.Elements = cloneElements(map[string]int{model.ElementEarth: 1, model.ElementArcane: 1})
	if !engine.canPayCostForCardAction(p0, timeCycle, timeCost, timeCost, paymentPurposeLearn, ActionMessage{}) {
		t.Fatal("3411101 should allow strict earth plus arcane payment")
	}
}

func TestRoyalConflictRadiantAngelLetsAnyElementPayLightCosts(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	lightCard := NewCardInstance(baseCard(t, "1521104"), 0, 1)
	cost := engine.effectiveCardPlayCost(p0, lightCard)

	p0.Elements = cloneElements(map[string]int{model.ElementFire: 5})
	if engine.canPayCostForCardAction(p0, lightCard, cost, cost, paymentPurposePlay, ActionMessage{}) {
		t.Fatal("light cost should not be payable with fire without 1521109")
	}

	angel := placeUnit(baseCard(t, "1521109"), 0, 0, 0, engine)
	if !engine.canPayCostForCardAction(p0, lightCard, cost, cost, paymentPurposePlay, ActionMessage{}) {
		t.Fatal("1521109 should let other elements pay light costs")
	}
	if !engine.payCostForCardAction(p0, lightCard, cost, cost, paymentPurposePlay, ActionMessage{}) || p0.Elements[model.ElementFire] != 0 {
		t.Fatalf("1521109 payment should spend fire as light, elements=%v", p0.Elements)
	}

	p0.Elements = cloneElements(map[string]int{model.ElementFire: 5})
	angel.Statuses[StatusPetrify] = 1
	if engine.canPayCostForCardAction(p0, lightCard, cost, cost, paymentPurposePlay, ActionMessage{}) {
		t.Fatal("petrified 1521109 should not enable other elements as light")
	}
}

func TestRoyalConflictCouncilSpokesmanReducesEnemyHandLimit(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	spokesman := placeUnit(baseCard(t, "1311103"), 0, 0, 0, engine)

	if got := engine.handLimitForPlayer(p1); got != engine.State.HandLimit-1 {
		t.Fatalf("1311103 should reduce opponent hand limit by one, got %d", got)
	}
	if got := engine.handLimitForPlayer(p0); got != engine.State.HandLimit {
		t.Fatalf("1311103 should not reduce its owner's hand limit, got %d", got)
	}
	spokesman.Statuses[StatusPetrify] = 1
	if got := engine.handLimitForPlayer(p1); got != engine.State.HandLimit {
		t.Fatalf("petrified 1311103 should not reduce hand limit, got %d", got)
	}

	spokesman.Statuses[StatusPetrify] = 0
	p1.Hand = []*CardInstance{
		NewCardInstance(baseCard(t, "1021001"), 1, 1),
		NewCardInstance(baseCard(t, "1021002"), 1, 1),
		NewCardInstance(baseCard(t, "1021003"), 1, 1),
		NewCardInstance(baseCard(t, "1021004"), 1, 1),
	}
	drawn := NewCardInstance(baseCard(t, "1021005"), 1, 1)
	p1.Deck = []*CardInstance{drawn}
	engine.drawCards(1, 1)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "discard" || engine.State.PendingAction.PlayerID != 1 || engine.State.PendingAction.MinSelect != 1 {
		t.Fatalf("1311103 should force immediate discard after exceeding reduced hand limit, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 1, drawn.InstanceID)
	if len(p1.Hand) != engine.State.HandLimit-1 || len(p1.Graveyard) != 1 || p1.Graveyard[0] != drawn {
		t.Fatalf("discard should restore reduced hand limit, hand=%d grave=%v", len(p1.Hand), cardsToInfo(p1.Graveyard))
	}

	normalEngine := setupReportedBugEngine(t)
	normalP1 := normalEngine.State.Players[1]
	for len(normalP1.Hand) < normalEngine.State.HandLimit {
		normalP1.Hand = append(normalP1.Hand, NewCardInstance(baseCard(t, "1021001"), 1, 1))
	}
	normalP1.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021002"), 1, 1)}
	normalEngine.drawCards(1, 1)
	if normalEngine.State.PendingAction != nil {
		t.Fatalf("normal draw over hand limit should not force immediate discard, pending=%+v", normalEngine.State.PendingAction)
	}

	searchEngine := setupReportedBugEngine(t)
	searchP1 := searchEngine.State.Players[1]
	placeUnit(baseCard(t, "1311103"), 0, 0, 0, searchEngine)
	for len(searchP1.Hand) < searchEngine.State.HandLimit-1 {
		searchP1.Hand = append(searchP1.Hand, NewCardInstance(baseCard(t, "1021001"), 1, 1))
	}
	searched := NewCardInstance(baseCard(t, "1021002"), 1, 1)
	searchP1.Deck = []*CardInstance{searched}
	if got := searchEngine.searchDeckCardToHand(1, searched.InstanceID); got != searched {
		t.Fatal("test setup should search card to hand")
	}
	if searchEngine.State.PendingAction == nil || searchEngine.State.PendingAction.Type != "discard" || searchEngine.State.PendingAction.PlayerID != 1 {
		t.Fatalf("1311103 should force discard after searching over hand limit, pending=%+v", searchEngine.State.PendingAction)
	}
	resolvePendingSelection(t, searchEngine, 1, searched.InstanceID)
	if len(searchP1.Hand) != searchEngine.State.HandLimit-1 || len(searchP1.Graveyard) != 1 || searchP1.Graveyard[0] != searched {
		t.Fatalf("searched card should be discarded back to reduced limit, hand=%d grave=%v", len(searchP1.Hand), cardsToInfo(searchP1.Graveyard))
	}

	graveEngine := setupReportedBugEngine(t)
	graveP1 := graveEngine.State.Players[1]
	placeUnit(baseCard(t, "1311103"), 0, 0, 0, graveEngine)
	for len(graveP1.Hand) < graveEngine.State.HandLimit-1 {
		graveP1.Hand = append(graveP1.Hand, NewCardInstance(baseCard(t, "1021001"), 1, 1))
	}
	recovered := NewCardInstance(baseCard(t, "1021003"), 1, 1)
	graveP1.Graveyard = []*CardInstance{recovered}
	if !graveEngine.moveGraveyardCardToHand(1, recovered.InstanceID) {
		t.Fatal("test setup should move graveyard card to hand")
	}
	if graveEngine.State.PendingAction == nil || graveEngine.State.PendingAction.Type != "discard" || graveEngine.State.PendingAction.PlayerID != 1 {
		t.Fatalf("1311103 should force discard after returning a card to hand, pending=%+v", graveEngine.State.PendingAction)
	}
}

func TestRoyalConflictTreasureCabinetExpandsEquipmentAndAllowsDuplicates(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	cabinet := NewCardInstance(baseCard(t, "2021105"), 0, engine.State.TurnNumber)
	cabinet.SlotIndex = 0
	p0.Equipment[0] = cabinet

	if got := equipmentSlotCapacity(p0); got != BaseEquipmentSlots+1 {
		t.Fatalf("2021105 should add one equipment slot, got %d", got)
	}

	weaponA := NewCardInstance(baseCard(t, "2121004"), 0, engine.State.TurnNumber)
	weaponA.SlotIndex = 1
	p0.Equipment[1] = weaponA
	weaponB := NewCardInstance(baseCard(t, "2121010"), 0, engine.State.TurnNumber)
	p0.Hand = []*CardInstance{weaponB}
	setAllElements(p0, 10)

	if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
		"instance_id": weaponB.InstanceID,
	}}); err != nil {
		t.Fatalf("2021105 should allow equipping duplicate weapon subtypes: %v", err)
	}
	if p0.Equipment[2] != weaponB {
		t.Fatalf("duplicate weapon should enter an empty equipment slot, equipment=%v", cardsToInfo(p0.Equipment[:]))
	}

	cabinet.Statuses[StatusPetrify] = 1
	if got := equipmentSlotCapacity(p0); got != BaseEquipmentSlots {
		t.Fatalf("petrified 2021105 should stop adding a slot, got %d", got)
	}
	weaponC := NewCardInstance(baseCard(t, "2121004"), 0, engine.State.TurnNumber)
	p0.Hand = []*CardInstance{weaponC}
	setAllElements(p0, 10)
	if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
		"instance_id": weaponC.InstanceID,
	}}); err == nil {
		t.Fatalf("petrified 2021105 should not allow another duplicate weapon")
	}
}

func TestRoyalConflictArcaneImpactGainsStatsForArcaneCosts(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	main := readySkill(baseCard(t, "3021101"), 0)
	boost := readySkill(baseCard(t, "3021101"), 0)
	p0.Skills[0] = main
	p0.Skills[1] = boost

	if got := engine.effectiveSpellPower(0, main, nil); got != main.Card.Power+1 {
		t.Fatalf("3021101 should gain +1 power as main spell, got %d", got)
	}
	if got := engine.effectiveSpellDamage(0, main, main.Card.Attack, nil); got != main.Card.Attack+1 {
		t.Fatalf("3021101 should gain +1 damage as main spell, got %d", got)
	}
	if got := engine.effectiveSpellPower(0, main, []*CardInstance{boost}); got != main.Card.Power+boost.Card.Power+2 {
		t.Fatalf("3021101 should gain +1 power for both main and boost contributions, got %d", got)
	}
	if got := engine.effectiveSpellDamage(0, main, main.Card.Attack, []*CardInstance{boost}); got != main.Card.Attack+2 {
		t.Fatalf("3021101 should gain +1 damage for both main and boost contributions, got %d", got)
	}
}

func TestRoyalConflictArcaneSealSealsEnemySkillAndRaisesOwnCost(t *testing.T) {
	noTargetEngine := setupReportedBugEngine(t)
	noTargetP0 := noTargetEngine.State.Players[0]
	noTargetSeal := readySkill(baseCard(t, "3021108"), 0)
	noTargetP0.Skills[0] = noTargetSeal
	noTargetP0.Elements[model.ElementArcane] = 10
	if err := noTargetEngine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": noTargetSeal.InstanceID,
	}}); err == nil {
		t.Fatal("3021108 should not be cast without an enemy skill target")
	}
	if noTargetP0.Elements[model.ElementArcane] != 10 || noTargetSeal.IsHorizontal || noTargetSeal.Statuses[StatusCooldown] > 0 || noTargetSeal.Statuses[arcaneSealExtraUseCostStatus] > 0 {
		t.Fatalf("failed 3021108 cast should not pay, tap, cool down, or raise cost; elements=%v horizontal=%v statuses=%v", noTargetP0.Elements, noTargetSeal.IsHorizontal, noTargetSeal.Statuses)
	}

	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	seal := readySkill(baseCard(t, "3021108"), 0)
	target := readySkill(baseCard(t, "3021005"), 1)
	p0.Skills[0] = seal
	p1.Skills[0] = target

	if cost := engine.effectiveSkillUseCost(p0, seal); cost[model.ElementArcane] != 2 {
		t.Fatalf("3021108 should start with printed use cost, cost=%v", cost)
	}
	if err := (Card3021108ArcaneSeal{}).OnSpellCast(&EffectContext{Engine: engine, Source: seal, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("3021108 spell cast: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "arcane_seal_skill" {
		t.Fatalf("3021108 should ask for an enemy skill target, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, target.InstanceID)
	if target.Statuses[StatusSeal] != 1 {
		t.Fatalf("3021108 should seal the selected enemy skill, statuses=%v", target.Statuses)
	}
	if err := engine.validateSkillForPurpose(target, skillPurposeAttack); err == nil {
		t.Fatal("sealed target skill should not be usable")
	}
	if cost := engine.effectiveSkillUseCost(p0, seal); cost[model.ElementArcane] != 4 {
		t.Fatalf("3021108 should permanently add 2 arcane to its own use cost, cost=%v", cost)
	}
}

func TestRoyalConflictArcanePurificationIgnoresFriendlyNegativeStatusesThisTurn(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	unit := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	purification := readySkill(baseCard(t, "3021105"), 0)

	if err := (Card3021105ArcanePurification{}).OnSpellCast(&EffectContext{Engine: engine, Source: purification, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("3021105 spell cast: %v", err)
	}
	if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModFriendlyNegativeStatusIgnore {
		t.Fatalf("3021105 should create friendly negative-status ignore modifier, modifiers=%+v", p0.TempModifiers)
	}
	if !engine.addStatus(unit, StatusPetrify, 2) || unit.Statuses[StatusPetrify] != 2 {
		t.Fatalf("3021105 should still allow negative statuses to be present, statuses=%v", unit.Statuses)
	}
	if engine.hasEffectiveStatus(unit, StatusPetrify) {
		t.Fatal("3021105 should make friendly negative statuses ineffective this turn")
	}

	engine.finishEndTurn(p0)
	if len(p0.TempModifiers) != 0 {
		t.Fatalf("3021105 modifier should expire at turn end, modifiers=%+v", p0.TempModifiers)
	}
	if unit.Statuses[StatusPetrify] != 1 || !engine.hasEffectiveStatus(unit, StatusPetrify) {
		t.Fatal("3021105 should stop suppressing negative statuses after turn end")
	}
}

func TestRoyalConflictArcaneDrainRequiresDistinctUseElements(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	arcaneDrain := readySkill(baseCard(t, "3021103"), 0)
	useCost := engine.effectiveSkillUseCost(p0, arcaneDrain)

	p0.Elements = cloneElements(map[string]int{model.ElementFire: 2})
	if engine.canPayCostForCardAction(p0, arcaneDrain, useCost, useCost, paymentPurposeUse, ActionMessage{}) {
		t.Fatal("3021103 should reject auto payment from only one element type")
	}
	if engine.canPayCostForCardAction(p0, arcaneDrain, useCost, useCost, paymentPurposeUse, ActionMessage{Data: map[string]any{
		"payment": map[string]any{model.ElementFire: float64(2)},
	}}) {
		t.Fatal("3021103 should reject explicit payment using the same element twice")
	}

	p0.Elements = cloneElements(map[string]int{model.ElementFire: 1, model.ElementWater: 1})
	if !engine.canPayCostForCardAction(p0, arcaneDrain, useCost, useCost, paymentPurposeUse, ActionMessage{}) {
		t.Fatal("3021103 should allow auto payment from two distinct element types")
	}
	if !engine.payCostForCardAction(p0, arcaneDrain, useCost, useCost, paymentPurposeUse, ActionMessage{}) {
		t.Fatal("3021103 should pay its use cost from distinct element types")
	}
	if p0.Elements[model.ElementFire] != 0 || p0.Elements[model.ElementWater] != 0 {
		t.Fatalf("3021103 should spend one of each distinct element, elements=%v", p0.Elements)
	}

	boostCost := map[string]int{model.ElementFire: 1}
	totalCost := mergeElementCosts(useCost, boostCost)
	p0.Elements = cloneElements(map[string]int{model.ElementFire: 2, model.ElementWater: 1})
	if !engine.canPayCostForCardAction(p0, arcaneDrain, useCost, totalCost, paymentPurposeUse, ActionMessage{Data: map[string]any{
		"payment": map[string]any{model.ElementFire: float64(2), model.ElementWater: float64(1)},
	}}) {
		t.Fatal("3021103 should allow duplicated elements when the duplicate pays a separate boost cost")
	}
}

func TestRoyalConflictArcaneDrainDrawsTwoOnCast(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	arcaneDrain := readySkill(baseCard(t, "3021103"), 0)
	drawA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
	drawB := NewCardInstance(baseCard(t, "1021002"), 0, 1)
	p0.Deck = []*CardInstance{drawA, drawB}
	p0.Hand = nil

	behavior := Card3021103ArcaneDrain{}
	if err := behavior.OnSpellCast(&EffectContext{Engine: engine, Source: arcaneDrain, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("3021103 spell cast: %v", err)
	}
	if len(p0.Hand) != 2 || p0.Hand[0] != drawA || p0.Hand[1] != drawB || len(p0.Deck) != 0 {
		t.Fatalf("3021103 should draw two cards, hand=%v deck=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Deck))
	}
}

func TestRoyalConflictFlipMechanicAndCards(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]

	unflippable := NewCardInstance(baseCard(t, "2211101"), 0, 1)
	fireOne := NewCardInstance(baseCard(t, "1121111"), 0, 1)
	fireTwo := NewCardInstance(baseCard(t, "1121101"), 0, 1)
	p0.Deck = []*CardInstance{unflippable, fireOne, fireTwo}
	rally := NewCardInstance(baseCard(t, "2121107"), 0, 1)
	if err := (Card2121107SacredFireRally{}).OnEnter(&EffectContext{Engine: engine, Source: rally, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("2121107 on enter: %v", err)
	}
	if len(p0.Hand) != 2 || !containsCardInstance(p0.Hand, fireOne) || !containsCardInstance(p0.Hand, fireTwo) || containsCardInstance(p0.Hand, unflippable) {
		t.Fatalf("2121107 should flip two fire companions and skip unflippable card, hand=%v", cardsToInfo(p0.Hand))
	}
	if p0.DrawCountThisTurn != 2 || p0.DrawnTurn[fireOne.InstanceID] == 0 || p0.DrawnTurn[fireTwo.InstanceID] == 0 {
		t.Fatalf("flipped cards should count as drawn, count=%d drawn=%v", p0.DrawCountThisTurn, p0.DrawnTurn)
	}
	if !containsCardInstance(p0.Deck, unflippable) {
		t.Fatalf("unflippable card should remain in deck, deck=%v", cardsToInfo(p0.Deck))
	}

	engine = setupReportedBugEngine(t)
	p0 = engine.State.Players[0]
	revealedOnDraw := NewCardInstance(baseCard(t, "1321003"), 0, 1)
	p0.Deck = []*CardInstance{revealedOnDraw}
	engine.flipDeckMatchesToHand(0, 1, 0, nil)
	if !p0.RevealedHand[revealedOnDraw.InstanceID] {
		t.Fatalf("flipped cards should honor reveal-on-draw, revealed=%v", p0.RevealedHand)
	}

	engine = setupReportedBugEngine(t)
	p0 = engine.State.Players[0]
	sword := NewCardInstance(baseCard(t, "2211101"), 0, 1)
	p0.Deck = []*CardInstance{sword}
	waterDivination := NewCardInstance(baseCard(t, "3221007"), 0, 1)
	if err := (Card3221007WaterDivination{}).OnSpellCast(&EffectContext{Engine: engine, Source: waterDivination, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("3221007 on spell cast: %v", err)
	}
	if engine.State.PendingAction == nil || len(engine.State.PendingAction.Candidates) != 1 || engine.State.PendingAction.Candidates[0]["can_select"] != false {
		t.Fatalf("water divination should show unsearchable sword but make it unselectable, pending=%+v", engine.State.PendingAction)
	}
	resolveWaterDivination(engine, 0, []*CardInstance{sword}, []string{sword.InstanceID}, nil)
	if len(p0.Hand) != 0 || !containsCardInstance(p0.Deck, sword) {
		t.Fatalf("forged water divination selection should not search unsearchable sword, hand=%v deck=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Deck))
	}

	engine = setupReportedBugEngine(t)
	p0 = engine.State.Players[0]
	sandworm := NewCardInstance(baseCard(t, "1421114"), 0, 1)
	p0.Deck = []*CardInstance{sandworm}
	bait := NewCardInstance(baseCard(t, "2421110"), 0, 1)
	if err := (Card2421110SandwormBait{}).OnUseItem(&EffectContext{Engine: engine, Source: bait, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("2421110 on use: %v", err)
	}
	if len(p0.Hand) != 1 || p0.Hand[0] != sandworm || sandworm.Statuses["入场费用地-2"] != 1 {
		t.Fatalf("2421110 should flip giant sandworm with earth discount, hand=%v statuses=%v", cardsToInfo(p0.Hand), sandworm.Statuses)
	}
	if cost := engine.effectiveCardPlayCost(p0, sandworm); cost[model.ElementEarth] != 4 {
		t.Fatalf("giant sandworm earth discount should reduce play cost to 4 earth, cost=%v", cost)
	}

	engine = setupReportedBugEngine(t)
	p0 = engine.State.Players[0]
	angel := NewCardInstance(baseCard(t, "1521109"), 0, 1)
	p0.Deck = []*CardInstance{angel}
	prayer := NewCardInstance(baseCard(t, "2521110"), 0, 1)
	if err := (Card2521110AngelPrayer{}).OnUseItem(&EffectContext{Engine: engine, Source: prayer, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("2521110 on use: %v", err)
	}
	if len(p0.Hand) != 1 || p0.Hand[0] != angel || angel.Statuses["入场费用光-1"] != 1 {
		t.Fatalf("2521110 should flip light spirit with light discount, hand=%v statuses=%v", cardsToInfo(p0.Hand), angel.Statuses)
	}
	if cost := engine.effectiveCardPlayCost(p0, angel); cost[model.ElementLight] != 1 {
		t.Fatalf("light spirit discount should reduce play cost to 1 light, cost=%v", cost)
	}
}

func containsCardInstance(cards []*CardInstance, target *CardInstance) bool {
	for _, card := range cards {
		if card == target {
			return true
		}
	}
	return false
}

func TestRoyalConflictSoulAndInsightUtilityCards(t *testing.T) {
	t.Run("illusionist returns a low-cost companion and gains its load", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		illusionist := placeUnit(baseCard(t, "1321105"), 0, 0, 0, engine)
		target := placeUnit(baseCard(t, "1321106"), 0, 1, 0, engine)
		target.ElementsGainBonus = map[string]int{model.ElementAir: 2}

		if err := (Card1321105Illusionist{}).OnUltimate(&EffectContext{Engine: engine, Source: illusionist, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1321105 ultimate: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "illusionist_return_companion" || !candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) {
			t.Fatalf("1321105 should ask for a low-cost friendly companion, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if p0.Units[1][0] != nil || !containsCardInstance(p0.Hand, target) {
			t.Fatalf("1321105 should return selected companion to hand, units=%v hand=%v", p0.Units[1][0], cardsToInfo(p0.Hand))
		}
		if p0.Elements[model.ElementArcane] != 0 || p0.Elements[model.ElementAir] != 3 {
			t.Fatalf("1321105 should gain returned companion load before hidden-zone reset, elements=%v", p0.Elements)
		}
		if target.ElementsGainBonus[model.ElementAir] != 0 || effectiveElementsGain(target)[model.ElementAir] != target.Card.ElementsGain[model.ElementAir] {
			t.Fatalf("1321105 should reset returned companion state in hand, bonus=%v load=%v", target.ElementsGainBonus, effectiveElementsGain(target))
		}
	})

	t.Run("soul devourer removes a soul marker to draw and gain shadow", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		p0.Hand = nil
		devourer := placeUnit(baseCard(t, "1621115"), 0, 0, 0, engine)
		skill := readySkill(baseCard(t, "3621102"), 0)
		skill.Statuses[soulMarkerStatus] = 1
		skill.PowerBonus = 2
		p0.Skills[0] = skill
		drawA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		drawB := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		p0.Deck = []*CardInstance{drawA, drawB}

		if err := (Card1621115SoulDevourer{}).OnPerTurn(&EffectContext{Engine: engine, Source: devourer, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1621115 per turn: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "soul_devourer_remove_marker" || !candidateContains(engine.State.PendingAction.Candidates, skill.InstanceID) {
			t.Fatalf("1621115 should ask for a friendly soul marker, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, skill.InstanceID)
		if skill.Statuses[soulMarkerStatus] != 0 || skill.PowerBonus != 0 {
			t.Fatalf("1621115 should remove the marker and its power bonus, statuses=%v power=%d", skill.Statuses, skill.PowerBonus)
		}
		if len(p0.Hand) != 2 || p0.Elements[model.ElementShadow] != 2 {
			t.Fatalf("1621115 should draw two and gain 2 shadow, hand=%v elements=%v", cardsToInfo(p0.Hand), p0.Elements)
		}
	})

	t.Run("soul staff exiles two shadow companions and marks a shadow spell", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		staff := NewCardInstance(baseCard(t, "2621112"), 0, 1)
		shadowA := NewCardInstance(baseCard(t, "1621101"), 0, 1)
		shadowB := NewCardInstance(baseCard(t, "1621102"), 0, 1)
		other := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		skill := readySkill(baseCard(t, "3621102"), 0)
		p0.Graveyard = []*CardInstance{shadowA, other, shadowB}
		p0.Skills[0] = skill

		if err := (Card2621112SoulStaff{}).OnPerTurn(&EffectContext{Engine: engine, Source: staff, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2621112 per turn: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "soul_staff_exile_companions" {
			t.Fatalf("2621112 should ask for shadow companion graveyard cards, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, shadowA.InstanceID, shadowB.InstanceID)
		if len(p0.Exile) != 2 || !containsCardInstance(p0.Exile, shadowA) || !containsCardInstance(p0.Exile, shadowB) || containsCardInstance(p0.Graveyard, shadowA) || containsCardInstance(p0.Graveyard, shadowB) {
			t.Fatalf("2621112 should exile selected shadow companions, exile=%v grave=%v", cardsToInfo(p0.Exile), cardsToInfo(p0.Graveyard))
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "soul_staff_mark_spell" || !candidateContains(engine.State.PendingAction.Candidates, skill.InstanceID) {
			t.Fatalf("2621112 should ask for a shadow spell after exiling, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, skill.InstanceID)
		if skill.Statuses[soulMarkerStatus] != 1 || skill.PowerBonus != 2 {
			t.Fatalf("2621112 should add one soul marker and +2 power, statuses=%v power=%d", skill.Statuses, skill.PowerBonus)
		}
	})

	t.Run("forest insight draws for earth companions then shuffles that many hand cards back", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		p0.Hand = nil
		placeUnit(baseCard(t, "1421102"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1421105"), 0, 1, 0, engine)
		skill := readySkill(baseCard(t, "3421101"), 0)
		drawA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		drawB := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		p0.Deck = []*CardInstance{drawA, drawB}

		if err := (Card3421101ForestInsight{}).OnSpellCast(&EffectContext{Engine: engine, Source: skill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0, "spell_being_cast": true}}); err != nil {
			t.Fatalf("3421101 spell cast: %v", err)
		}
		if len(p0.Hand) != 2 || engine.State.PendingAction == nil || engine.State.PendingAction.Type != "forest_insight_shuffle_hand" {
			t.Fatalf("3421101 should draw two then ask to shuffle two hand cards back, hand=%v pending=%+v", cardsToInfo(p0.Hand), engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, drawA.InstanceID, drawB.InstanceID)
		if len(p0.Hand) != 0 || len(p0.Deck) != 2 || !containsCardInstance(p0.Deck, drawA) || !containsCardInstance(p0.Deck, drawB) {
			t.Fatalf("3421101 should shuffle selected cards back into deck, hand=%v deck=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Deck))
		}

		shortDeckEngine := setupEffectTest(t)
		shortP0 := shortDeckEngine.State.Players[0]
		shortP0.Hand = nil
		for col := 0; col < 3; col++ {
			placeUnit(baseCard(t, "1421102"), 0, col, 0, shortDeckEngine)
		}
		onlyDraw := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		shortP0.Deck = []*CardInstance{onlyDraw}
		if err := (Card3421101ForestInsight{}).OnSpellCast(&EffectContext{Engine: shortDeckEngine, Source: readySkill(baseCard(t, "3421101"), 0), PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("3421101 short deck spell cast: %v", err)
		}
		if shortDeckEngine.State.PendingAction == nil || shortDeckEngine.State.PendingAction.MinSelect != 1 || shortDeckEngine.State.PendingAction.MaxSelect != 1 {
			t.Fatalf("3421101 should shuffle back actual drawn count, pending=%+v hand=%v", shortDeckEngine.State.PendingAction, cardsToInfo(shortP0.Hand))
		}
	})
}

func TestRoyalConflictSummonUtilityCards(t *testing.T) {
	t.Run("volcano salamander sacrifices itself at mastery two to summon from hand", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		target := NewCardInstance(baseCard(t, "1121105"), 0, 1)
		p0.Hand = []*CardInstance{target}
		salamander := placeUnit(baseCard(t, "1121101"), 0, 1, 1, engine)
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				if col == 1 && row == 1 {
					continue
				}
				placeUnit(baseCard(t, "1021001"), 0, col, row, engine)
			}
		}

		engine.advanceMastery(salamander, 0, 2)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "volcano_salamander_summon_card" || !candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) {
			t.Fatalf("1121101 should ask for a fire companion in hand at mastery two, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		sourcePos := Position{Col: 1, Row: 1}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "volcano_salamander_summon_position" || !candidateContains(engine.State.PendingAction.Candidates, positionSelectionID(sourcePos)) {
			t.Fatalf("1121101 should allow summoning into its sacrificed slot on a full board, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, positionSelectionID(sourcePos))
		if p0.Units[sourcePos.Col][sourcePos.Row] != target || !containsCardInstance(p0.Graveyard, salamander) || len(p0.Hand) != 0 {
			t.Fatalf("1121101 should sacrifice itself and summon target from hand, unit=%v grave=%v hand=%v", p0.Units[sourcePos.Col][sourcePos.Row], cardsToInfo(p0.Graveyard), cardsToInfo(p0.Hand))
		}
	})

	t.Run("prayer flame adds markers or spends them to summon a fire companion from hand", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		p0.Hand = nil
		flame := readySkill(baseCard(t, "3121103"), 0)
		behavior := Card3121103PrayerFlame{}

		if err := behavior.OnSpellCast(&EffectContext{Engine: engine, Source: flame, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("3121103 add markers cast: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "prayer_flame_choice" || candidateContains(engine.State.PendingAction.Candidates, "summon") {
			t.Fatalf("3121103 without target should offer only marker choice, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, "add_markers")
		if flame.Statuses[prayerFlameMarkerStatus] != 3 {
			t.Fatalf("3121103 should add three markers, statuses=%v", flame.Statuses)
		}

		fireCompanion := NewCardInstance(baseCard(t, "1121101"), 0, 1)
		p0.Hand = []*CardInstance{fireCompanion}
		if err := behavior.OnSpellCast(&EffectContext{Engine: engine, Source: flame, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("3121103 summon cast: %v", err)
		}
		if engine.State.PendingAction == nil || !candidateContains(engine.State.PendingAction.Candidates, "summon") {
			t.Fatalf("3121103 with markers and hand target should offer summon, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, "summon")
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "prayer_flame_summon_card" || !candidateContains(engine.State.PendingAction.Candidates, fireCompanion.InstanceID) {
			t.Fatalf("3121103 should ask which fire companion to summon, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, fireCompanion.InstanceID)
		pos := Position{Col: 2, Row: 2}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "prayer_flame_summon_position" {
			t.Fatalf("3121103 should ask for summon position, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, positionSelectionID(pos))
		if p0.Units[pos.Col][pos.Row] != fireCompanion || len(p0.Hand) != 0 || flame.Statuses[prayerFlameMarkerStatus] != 0 {
			t.Fatalf("3121103 should summon from hand and remove all markers, unit=%v hand=%v statuses=%v", p0.Units[pos.Col][pos.Row], cardsToInfo(p0.Hand), flame.Statuses)
		}
	})

	t.Run("phantom lizard consumes itself and splits into two normal lizards", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		p1.Equipment[0] = NewCardInstance(baseCard(t, "2121002"), 1, 1)
		lizard := placeUnit(baseCard(t, "1421106"), 0, 0, 0, engine)
		spiritSkill := readySkill(baseCard(t, "3421101"), 0)

		if err := (Card1421106PhantomLizard{}).OnSpellCast(&EffectContext{Engine: engine, Source: lizard, Target: spiritSkill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("1421106 spell cast: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "phantom_lizard_split" {
			t.Fatalf("1421106 should ask whether to split after spirit skill, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, lizard.InstanceID)
		if !containsCardInstance(p0.Graveyard, lizard) || p0.Elements[model.ElementEarth] != 1 || !lizard.UltimateUsed || lizard.Statuses[StatusBurn] != 1 {
			t.Fatalf("1421106 should consume for earth, trigger consume watchers, mark ultimate used, and move source to graveyard, grave=%v elements=%v ultimate=%v statuses=%v", cardsToInfo(p0.Graveyard), p0.Elements, lizard.UltimateUsed, lizard.Statuses)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "phantom_lizard_first_position" {
			t.Fatalf("1421106 should ask for first lizard position, pending=%+v", engine.State.PendingAction)
		}
		firstPos := Position{Col: 0, Row: 0}
		resolvePendingSelection(t, engine, 0, positionSelectionID(firstPos))
		if p0.Units[firstPos.Col][firstPos.Row] == nil || p0.Units[firstPos.Col][firstPos.Row].Card.Number != "1401101" {
			t.Fatalf("1421106 should summon first normal lizard, unit=%v", p0.Units[firstPos.Col][firstPos.Row])
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "phantom_lizard_second_position" {
			t.Fatalf("1421106 should ask for second lizard position, pending=%+v", engine.State.PendingAction)
		}
		secondPos := Position{Col: 1, Row: 0}
		resolvePendingSelection(t, engine, 0, positionSelectionID(secondPos))
		if p0.Units[secondPos.Col][secondPos.Row] == nil || p0.Units[secondPos.Col][secondPos.Row].Card.Number != "1401101" {
			t.Fatalf("1421106 should summon second normal lizard, unit=%v", p0.Units[secondPos.Col][secondPos.Row])
		}
	})
}

func TestRoyalConflictResourceConversionCards(t *testing.T) {
	t.Run("chief advisor felin sacrifices a fire companion to discount the next fire card", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		felin := placeUnit(baseCard(t, "4111101"), 0, 1, 1, engine)
		sacrifice := placeUnit(baseCard(t, "1121107"), 0, 0, 0, engine)
		target := NewCardInstance(baseCard(t, "1121107"), 0, 1)
		p0.Hand = []*CardInstance{target}

		if err := (Card4111101ChiefAdvisorFelin{}).OnUltimate(&EffectContext{Engine: engine, Source: felin, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("4111101 ultimate: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "felin_sacrifice_fire_companion" || !candidateContains(engine.State.PendingAction.Candidates, sacrifice.InstanceID) {
			t.Fatalf("4111101 should ask for a fire companion sacrifice, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, sacrifice.InstanceID)
		if !containsCardInstance(p0.Graveyard, sacrifice) || len(p0.TempModifiers) != 2 {
			t.Fatalf("4111101 should sacrifice target and create next fire discount, grave=%v modifiers=%+v", cardsToInfo(p0.Graveyard), p0.TempModifiers)
		}
		if cost := engine.effectiveCardPlayCost(p0, target); cost[model.ElementFire] != 0 || cost[model.ElementAir] != 0 {
			t.Fatalf("4111101 should discount next fire card by sacrificed card's element costs, cost=%v sacrificeCost=%v", cost, sacrifice.Card.ElementsCost)
		}
		nonFire := NewCardInstance(baseCard(t, "1221104"), 0, 1)
		if cost := engine.effectiveCardPlayCost(p0, nonFire); cost[model.ElementWater] != nonFire.Card.ElementsCost[model.ElementWater] {
			t.Fatalf("4111101 should not discount non-fire cards, cost=%v", cost)
		}
		engine.notifyCardPlayCostPaid(p0, target)
		if len(p0.TempModifiers) != 0 {
			t.Fatalf("4111101 discount should be consumed by the next fire card, modifiers=%+v", p0.TempModifiers)
		}
	})

	t.Run("beast taming collar binds and consumes the chosen fire companion for entry cost", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		collar := NewCardInstance(baseCard(t, "2121106"), 0, 1)
		p0.Equipment[0] = collar
		target := placeUnit(baseCard(t, "1121105"), 0, 0, 0, engine)
		target.IsHorizontal = false

		if err := (Card2121106BeastTamingCollar{}).OnEnter(&EffectContext{Engine: engine, Source: collar, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2121106 enter: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "beast_taming_collar_target" || !candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) {
			t.Fatalf("2121106 should ask for eligible fire companion, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if collarTarget(engine, 0, collar) != target {
			t.Fatalf("2121106 should remember selected companion, statuses=%v", collar.Statuses)
		}
		if err := (Card2121106BeastTamingCollar{}).OnPerTurn(&EffectContext{Engine: engine, Source: collar, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2121106 per-turn: %v", err)
		}
		if !target.IsHorizontal || p0.Elements[model.ElementFire] != target.Card.ElementsCost[model.ElementFire] {
			t.Fatalf("2121106 should consume target and gain its entry fire cost, horizontal=%v elements=%v cost=%v", target.IsHorizontal, p0.Elements, target.Card.ElementsCost)
		}
	})

	t.Run("lavafort archivist and legion staff officer flip cards after creation spells", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		archivist := placeUnit(baseCard(t, "1121110"), 0, 0, 0, engine)
		officer := placeUnit(baseCard(t, "1121115"), 0, 2, 0, engine)
		createSkill := readySkill(baseCard(t, "3021107"), 0)
		rune := NewCardInstance(baseCard(t, "2021111"), 0, 1)
		fireConsumable := NewCardInstance(baseCard(t, "2121108"), 0, 1)
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1), rune, fireConsumable}

		if err := (Card1121110LavafortArchivist{}).OnSpellCast(&EffectContext{Engine: engine, Source: archivist, Target: createSkill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("1121110 spell cast: %v", err)
		}
		if !containsCardInstance(p0.Hand, rune) || !archivist.UltimateUsed {
			t.Fatalf("1121110 should flip a rune or scroll and spend ultimate, hand=%v ultimate=%v", cardsToInfo(p0.Hand), archivist.UltimateUsed)
		}
		if err := (Card1121115LegionStaffOfficer{}).OnSpellCast(&EffectContext{Engine: engine, Source: officer, Target: createSkill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("1121115 spell cast: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "legion_staff_officer_flip_fire_consumable" {
			t.Fatalf("1121115 should offer optional fire consumable flip, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, officer.InstanceID)
		if !containsCardInstance(p0.Hand, fireConsumable) || !p0.DiscardAtTurnEnd[fireConsumable.InstanceID] || officer.UsedThisTurn != 1 {
			t.Fatalf("1121115 should flip a fire consumable, mark it for discard, and spend use, hand=%v discard=%v used=%d", cardsToInfo(p0.Hand), p0.DiscardAtTurnEnd, officer.UsedThisTurn)
		}
	})

	t.Run("aging potion removes earth load and advances target to next mastery", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		target := readySkill(baseCard(t, "3421001"), 0)
		target.ElementsGainBonus[model.ElementEarth] = 1
		p0.Skills[0] = target
		potion := NewCardInstance(baseCard(t, "2421106"), 0, 1)

		if err := (Card2421106AgingPotion{}).OnUseItem(&EffectContext{Engine: engine, Source: potion, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2421106 use: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "aging_potion_mastery" || !candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) {
			t.Fatalf("2421106 should ask for a friendly card with earth load and mastery, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if target.ElementsGainBonus[model.ElementEarth] != 0 || target.Statuses[StatusMastery] != 1 {
			t.Fatalf("2421106 should remove one earth load and advance mastery, load=%v mastery=%v", target.ElementsGainBonus, target.Statuses[StatusMastery])
		}
	})

	t.Run("conductor los equips finale violin after four consumed cards and resets them by consuming itself", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		los := placeUnit(baseCard(t, "1011102"), 0, 1, 1, engine)
		consumed := []*CardInstance{
			placeUnit(baseCard(t, "1001101"), 0, 0, 0, engine),
			placeUnit(baseCard(t, "1021101"), 0, 1, 0, engine),
			placeUnit(baseCard(t, "1021102"), 0, 2, 0, engine),
			placeUnit(baseCard(t, "1021103"), 0, 0, 2, engine),
		}
		for i, card := range consumed {
			card.IsHorizontal = false
			engine.consumeCardForEffectWithTriggers(0, card, card.Card.ElementsGain, "")
			if i < 3 && engine.State.PendingAction != nil {
				t.Fatalf("1011102 should wait for four consumed cards, i=%d pending=%+v", i, engine.State.PendingAction)
			}
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "conductor_los_equip_finale_violin" {
			t.Fatalf("1011102 should offer finale violin after four consumes, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0)
		if los.Statuses[conductorConsumedCountStatus] != 0 || p0.Equipment[0] != nil {
			t.Fatalf("1011102 should spend four-count even when skipped, statuses=%v equipment=%v", los.Statuses, cardsToInfo(p0.Equipment[:]))
		}
		for _, card := range consumed {
			card.IsHorizontal = false
			engine.consumeCardForEffectWithTriggers(0, card, card.Card.ElementsGain, "")
		}
		resolvePendingSelection(t, engine, 0, los.InstanceID)
		if p0.Equipment[0] == nil || p0.Equipment[0].Card.Number != "2001101" || !p0.Equipment[0].IsHorizontal {
			t.Fatalf("1011102 should equip a horizontal finale violin, equipment=%v", cardsToInfo(p0.Equipment[:]))
		}

		violin := p0.Equipment[0]
		los.IsHorizontal = false
		if err := (Card1011102ConductorLos{}).OnPerTurn(&EffectContext{Engine: engine, Source: los, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1011102 per-turn: %v", err)
		}
		if !los.IsHorizontal || violin.IsHorizontal {
			t.Fatalf("1011102 should consume itself and reset finale violins, los=%v violin=%v", los.IsHorizontal, violin.IsHorizontal)
		}

		target := placeUnit(baseCard(t, "1021101"), 1, 1, 0, engine)
		p0.Elements[model.ElementArcane] = 1
		if err := engine.handleAttack(0, ActionMessage{Data: map[string]any{
			"attacker_id": violin.InstanceID,
			"target_col":  float64(target.Position.Col),
			"target_row":  float64(target.Position.Row),
		}}); err == nil {
			t.Fatalf("2001101 should require 2 arcane to attack")
		}
		p0.Elements[model.ElementArcane] = 2
		if err := engine.handleAttack(0, ActionMessage{Data: map[string]any{
			"attacker_id": violin.InstanceID,
			"target_col":  float64(target.Position.Col),
			"target_row":  float64(target.Position.Row),
		}}); err != nil {
			t.Fatalf("2001101 attack with payment: %v", err)
		}
		if p0.Elements[model.ElementArcane] != 0 || !violin.IsHorizontal || target.CurrentLife != target.Card.Life-1 {
			t.Fatalf("2001101 should pay 2 arcane and attack for 1, elements=%v horizontal=%v targetLife=%d", p0.Elements, violin.IsHorizontal, target.CurrentLife)
		}
	})
}

func TestRoyalConflictSpellFollowupUtilityCards(t *testing.T) {
	t.Run("divine fire rider consumes another fire companion to boost the next fire spell", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		rider := placeUnit(baseCard(t, "1121109"), 0, 1, 1, engine)
		target := placeUnit(baseCard(t, "1121105"), 0, 0, 1, engine)
		target.IsHorizontal = false

		if err := (Card1121109DivineFireRider{}).OnUltimate(&EffectContext{Engine: engine, Source: rider, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1121109 ultimate: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "divine_fire_rider_consume_companion" || !candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) {
			t.Fatalf("1121109 should ask for another ready fire companion, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		bonus := totalElementCost(target.Card.ElementsCost)
		fireSpell := readySkill(baseCard(t, "3121001"), 0)
		waterSpell := readySkill(baseCard(t, "3221001"), 0)
		if !target.IsHorizontal || p0.Elements[model.ElementFire] != target.Card.ElementsGain[model.ElementFire] {
			t.Fatalf("1121109 should consume target and gain its load, horizontal=%v elements=%v", target.IsHorizontal, p0.Elements)
		}
		if got := engine.effectiveSpellPower(0, fireSpell, nil); got != fireSpell.Card.Power+bonus {
			t.Fatalf("1121109 should boost next fire spell by consumed entry cost, got=%d want=%d", got, fireSpell.Card.Power+bonus)
		}
		if got := engine.effectiveSpellPower(0, waterSpell, nil); got != waterSpell.Card.Power {
			t.Fatalf("1121109 should not boost non-fire spell, got=%d want=%d", got, waterSpell.Card.Power)
		}
		engine.consumeNextSpellPowerBonuses(p0, fireSpell)
		if len(p0.TempModifiers) != 0 {
			t.Fatalf("1121109 next fire boost should be consumed, modifiers=%+v", p0.TempModifiers)
		}
	})

	t.Run("coral wendy can pay water to reset a just-used low-cost spell", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		wendy := placeUnit(baseCard(t, "1211103"), 0, 1, 1, engine)
		spell := readySkill(baseCard(t, "3221001"), 0)
		spell.IsHorizontal = true
		p0.Skills[0] = spell
		p0.Elements[model.ElementWater] = 2

		if err := (Card1211103SeaHeroineCoralWendy{}).OnSpellCast(&EffectContext{Engine: engine, Source: wendy, Target: spell, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("1211103 spell cast: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "coral_wendy_reset_spell" {
			t.Fatalf("1211103 should offer paid reset for low-cost spell, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, spell.InstanceID)
		if spell.IsHorizontal || p0.Elements[model.ElementWater] != 0 || wendy.UsedThisTurn != 1 {
			t.Fatalf("1211103 should pay 2 water, reset spell, and spend one use, horizontal=%v elements=%v used=%d", spell.IsHorizontal, p0.Elements, wendy.UsedThisTurn)
		}

		staleEngine := setupEffectTest(t)
		staleP0 := staleEngine.State.Players[0]
		staleWendy := placeUnit(baseCard(t, "1211103"), 0, 1, 1, staleEngine)
		staleSpell := readySkill(baseCard(t, "3221001"), 0)
		staleSpell.IsHorizontal = true
		staleP0.Skills[0] = staleSpell
		staleP0.Elements[model.ElementWater] = 2
		if err := (Card1211103SeaHeroineCoralWendy{}).OnSpellCast(&EffectContext{Engine: staleEngine, Source: staleWendy, Target: staleSpell, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("1211103 stale spell cast: %v", err)
		}
		staleSpell.IsHorizontal = false
		if err := staleEngine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{staleSpell.InstanceID}}}); err == nil {
			t.Fatalf("1211103 should reject stale reset target")
		}
		if staleP0.Elements[model.ElementWater] != 2 || staleWendy.UsedThisTurn != 0 {
			t.Fatalf("1211103 should revalidate live spell state before paying, elements=%v used=%d", staleP0.Elements, staleWendy.UsedThisTurn)
		}
	})

	t.Run("heart lotus mirror mage flips a water item and sets counter traps for free", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		mage := placeUnit(baseCard(t, "1221108"), 0, 1, 1, engine)
		createSkill := readySkill(baseCard(t, "3021107"), 0)
		counter := NewCardInstance(baseCard(t, "2221002"), 0, 1)
		p0.Hand = nil
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1), counter}

		if err := (Card1221108HeartLotusMirrorMage{}).OnSpellCast(&EffectContext{Engine: engine, Source: mage, Target: createSkill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("1221108 spell cast: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "heart_lotus_mirror_mage_flip_water_item" {
			t.Fatalf("1221108 should offer optional water item flip, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, mage.InstanceID)
		if len(p0.Hand) != 0 || p0.Equipment[0] != counter || !counter.IsSetCounter || counter.Statuses[entryCostZeroStatus] != 1 || mage.UsedThisTurn != 1 {
			t.Fatalf("1221108 should flip and freely set water counter, hand=%v equipment=%v statuses=%v used=%d", cardsToInfo(p0.Hand), cardsToInfo(p0.Equipment[:]), counter.Statuses, mage.UsedThisTurn)
		}
		counter.Statuses["入场费用"+model.ElementArcane+"-1"] = 0
		if cost := engine.effectiveCardPlayCost(p0, counter); cost[model.ElementArcane] != 0 {
			t.Fatalf("1221108 free counter should keep entry cost at zero after modifiers, cost=%v statuses=%v", cost, counter.Statuses)
		}
		costEngine := setupEffectTest(t)
		costP0 := costEngine.State.Players[0]
		placeUnit(baseCard(t, "1111103"), 1, 0, 0, costEngine)
		costCounter := NewCardInstance(baseCard(t, "2221002"), 0, 1)
		costP0.Hand = []*CardInstance{costCounter}
		makeEntryCostZero(costCounter)
		if cost := costEngine.effectiveCardPlayCost(costP0, costCounter); cost[model.ElementArcane] != 0 {
			t.Fatalf("absolute free entry marker should override later global cost increases, cost=%v statuses=%v", cost, costCounter.Statuses)
		}

		limitEngine := setupEffectTest(t)
		limitP0 := limitEngine.State.Players[0]
		placeUnit(baseCard(t, "1311103"), 1, 0, 0, limitEngine)
		limitMage := placeUnit(baseCard(t, "1221108"), 0, 1, 1, limitEngine)
		limitP0.Hand = nil
		for len(limitP0.Hand) < limitEngine.State.HandLimit-1 {
			limitP0.Hand = append(limitP0.Hand, NewCardInstance(baseCard(t, "1021001"), 0, 1))
		}
		limitCounter := NewCardInstance(baseCard(t, "2221002"), 0, 1)
		limitP0.Deck = []*CardInstance{limitCounter}
		if err := (Card1221108HeartLotusMirrorMage{}).OnSpellCast(&EffectContext{Engine: limitEngine, Source: limitMage, Target: createSkill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("1221108 hand-limit spell cast: %v", err)
		}
		resolvePendingSelection(t, limitEngine, 0, limitMage.InstanceID)
		if limitEngine.State.PendingAction != nil || limitP0.Equipment[0] != limitCounter || containsCardInstance(limitP0.Hand, limitCounter) {
			t.Fatalf("1221108 should set free counter before immediate hand-limit discard, pending=%+v hand=%v equipment=%v", limitEngine.State.PendingAction, cardsToInfo(limitP0.Hand), cardsToInfo(limitP0.Equipment[:]))
		}
	})

	t.Run("king of beasts flips an earth companion and summons it for free", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		earthCompanion := NewCardInstance(baseCard(t, "1421104"), 0, 1)
		p0.Hand = nil
		p0.Deck = []*CardInstance{earthCompanion, NewCardInstance(baseCard(t, "1021001"), 0, 1)}

		if err := (Card1411103KingOfBeasts{}).OnEnter(&EffectContext{Engine: engine, Source: placeUnit(baseCard(t, "1411103"), 0, 1, 1, engine), PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1411103 enter: %v", err)
		}
		if !containsCardInstance(p0.Hand, earthCompanion) || engine.State.PendingAction == nil || engine.State.PendingAction.Type != "king_of_beasts_summon_earth_companion" {
			t.Fatalf("1411103 should flip earth companion and ask for summon position, hand=%v pending=%+v", cardsToInfo(p0.Hand), engine.State.PendingAction)
		}
		pos := Position{Col: 0, Row: 0}
		resolvePendingSelection(t, engine, 0, positionSelectionID(pos))
		if p0.Units[pos.Col][pos.Row] != earthCompanion || containsCardInstance(p0.Hand, earthCompanion) {
			t.Fatalf("1411103 should summon flipped earth companion for free, unit=%v hand=%v", cardToInfo(p0.Units[pos.Col][pos.Row]), cardsToInfo(p0.Hand))
		}

		limitEngine := setupEffectTest(t)
		limitP0 := limitEngine.State.Players[0]
		placeUnit(baseCard(t, "1311103"), 1, 0, 0, limitEngine)
		limitP0.Hand = nil
		for len(limitP0.Hand) < limitEngine.State.HandLimit-1 {
			limitP0.Hand = append(limitP0.Hand, NewCardInstance(baseCard(t, "1021001"), 0, 1))
		}
		limitEarthCompanion := NewCardInstance(baseCard(t, "1421104"), 0, 1)
		limitP0.Deck = []*CardInstance{limitEarthCompanion}
		if err := (Card1411103KingOfBeasts{}).OnEnter(&EffectContext{Engine: limitEngine, Source: placeUnit(baseCard(t, "1411103"), 0, 1, 1, limitEngine), PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1411103 hand-limit enter: %v", err)
		}
		if limitEngine.State.PendingAction == nil || limitEngine.State.PendingAction.Type != "king_of_beasts_summon_earth_companion" {
			t.Fatalf("1411103 should ask summon position before hand-limit discard, pending=%+v", limitEngine.State.PendingAction)
		}
		resolvePendingSelection(t, limitEngine, 0, positionSelectionID(Position{Col: 0, Row: 0}))
		if limitEngine.State.PendingAction != nil || limitP0.Units[0][0] != limitEarthCompanion {
			t.Fatalf("1411103 should summon before immediate hand-limit discard is considered, pending=%+v unit=%v", limitEngine.State.PendingAction, cardToInfo(limitP0.Units[0][0]))
		}
	})

	t.Run("general kelan flips a fire card after successful defense then discards", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		kelan := NewCardInstance(baseCard(t, "4111102"), 0, 1)
		p0.Hero = kelan
		discard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		flipped := NewCardInstance(baseCard(t, "1121105"), 0, 1)
		p0.Hand = []*CardInstance{discard}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1221104"), 0, 1), flipped}

		if err := (Card4111102GeneralKelan{}).OnDefend(&EffectContext{Engine: engine, Source: kelan, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": false}}); err != nil {
			t.Fatalf("4111102 failed defense: %v", err)
		}
		if engine.State.PendingAction != nil {
			t.Fatalf("4111102 should not trigger on failed defense, pending=%+v", engine.State.PendingAction)
		}
		if err := (Card4111102GeneralKelan{}).OnDefend(&EffectContext{Engine: engine, Source: kelan, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": true}}); err != nil {
			t.Fatalf("4111102 successful defense: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "general_kelan_flip_fire_card" {
			t.Fatalf("4111102 should offer fire flip after successful defense, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, kelan.InstanceID)
		if !containsCardInstance(p0.Hand, flipped) || engine.State.PendingAction == nil || engine.State.PendingAction.Type != "general_kelan_discard" {
			t.Fatalf("4111102 should flip fire card then ask to discard, hand=%v pending=%+v", cardsToInfo(p0.Hand), engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, discard.InstanceID)
		if containsCardInstance(p0.Hand, discard) || !containsCardInstance(p0.Graveyard, discard) || kelan.UsedThisTurn != 1 {
			t.Fatalf("4111102 should discard chosen hand card and spend use, hand=%v grave=%v used=%d", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard), kelan.UsedThisTurn)
		}

		limitEngine := setupEffectTest(t)
		limitP0 := limitEngine.State.Players[0]
		placeUnit(baseCard(t, "1311103"), 1, 0, 0, limitEngine)
		limitKelan := NewCardInstance(baseCard(t, "4111102"), 0, 1)
		limitP0.Hero = limitKelan
		limitP0.Hand = nil
		for len(limitP0.Hand) < limitEngine.State.HandLimit-1 {
			limitP0.Hand = append(limitP0.Hand, NewCardInstance(baseCard(t, "1021001"), 0, 1))
		}
		limitFlipped := NewCardInstance(baseCard(t, "1121105"), 0, 1)
		limitP0.Deck = []*CardInstance{limitFlipped}
		if err := (Card4111102GeneralKelan{}).OnDefend(&EffectContext{Engine: limitEngine, Source: limitKelan, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": true}}); err != nil {
			t.Fatalf("4111102 hand-limit successful defense: %v", err)
		}
		resolvePendingSelection(t, limitEngine, 0, limitKelan.InstanceID)
		if limitEngine.State.PendingAction == nil || limitEngine.State.PendingAction.Type != "general_kelan_discard" || !containsCardInstance(limitP0.Hand, limitFlipped) {
			t.Fatalf("4111102 should ask its own discard before hand-limit discard, pending=%+v hand=%v", limitEngine.State.PendingAction, cardsToInfo(limitP0.Hand))
		}
	})
}

func TestRoyalConflictDiscardDrivenCompanions(t *testing.T) {
	t.Run("sparrow silverleaf deals capped entry damage from hand cards discarded this turn", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		startLife := target.CurrentLife
		p0.DiscardedHandCountThisTurn = 5

		if err := (Card1311101SparrowSilverleaf{}).OnEnter(&EffectContext{Engine: engine, Source: placeUnit(baseCard(t, "1311101"), 0, 1, 1, engine), PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1311101 enter: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "sparrow_silverleaf_entry_damage" || !candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) {
			t.Fatalf("1311101 should ask for an in-range enemy target, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if target.CurrentLife != startLife-3 {
			t.Fatalf("1311101 should cap entry damage at 3, life=%d start=%d", target.CurrentLife, startLife)
		}

		noDiscardEngine := setupEffectTest(t)
		if err := (Card1311101SparrowSilverleaf{}).OnEnter(&EffectContext{Engine: noDiscardEngine, Source: placeUnit(baseCard(t, "1311101"), 0, 1, 1, noDiscardEngine), PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1311101 no-discard enter: %v", err)
		}
		if noDiscardEngine.State.PendingAction != nil {
			t.Fatalf("1311101 should not prompt without discarded hand cards, pending=%+v", noDiscardEngine.State.PendingAction)
		}
	})

	t.Run("speckled sparrow can pay air to summon itself after being discarded from hand", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		sparrow := NewCardInstance(baseCard(t, "1321102"), 0, 1)
		p0.Hand = []*CardInstance{sparrow}
		p0.Elements[model.ElementAir] = 1

		if discarded := engine.discardHandCardAt(0, 0); discarded != sparrow {
			t.Fatalf("test setup should discard the sparrow, got=%v", cardToInfo(discarded))
		}
		if p0.DiscardedHandCountThisTurn != 1 || engine.State.PendingAction == nil || engine.State.PendingAction.Type != "speckled_sparrow_discard_summon" {
			t.Fatalf("1321102 should count discard and offer summon, discarded=%d pending=%+v", p0.DiscardedHandCountThisTurn, engine.State.PendingAction)
		}
		pos := Position{Col: 0, Row: 0}
		resolvePendingSelection(t, engine, 0, positionSelectionID(pos))
		if p0.Elements[model.ElementAir] != 0 || p0.Units[pos.Col][pos.Row] != sparrow || containsCardInstance(p0.Graveyard, sparrow) {
			t.Fatalf("1321102 should pay 1 air and summon from graveyard, elements=%v unit=%v grave=%v", p0.Elements, cardToInfo(p0.Units[pos.Col][pos.Row]), cardsToInfo(p0.Graveyard))
		}

		failEngine := setupEffectTest(t)
		failP0 := failEngine.State.Players[0]
		failSparrow := NewCardInstance(baseCard(t, "1321102"), 0, 1)
		failP0.Hand = []*CardInstance{failSparrow}
		failEngine.discardHandCardAt(0, 0)
		if failEngine.State.PendingAction != nil || !containsCardInstance(failP0.Graveyard, failSparrow) {
			t.Fatalf("1321102 should not prompt without air payment, pending=%+v grave=%v", failEngine.State.PendingAction, cardsToInfo(failP0.Graveyard))
		}
	})

	t.Run("black pine coffin discards low-cost shadow companions and resolves their deathrattles", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		coffin := NewCardInstance(baseCard(t, "2621108"), 0, 1)
		husk := NewCardInstance(baseCard(t, "1621002"), 0, 1)
		nonShadow := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Hand = []*CardInstance{husk, nonShadow}
		startArcane := p0.Elements[model.ElementArcane]

		if err := (Card2621108BlackPineCoffin{}).OnEnter(&EffectContext{Engine: engine, Source: coffin, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2621108 enter: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "black_pine_coffin_discard_shadow_companions" {
			t.Fatalf("2621108 should ask for shadow companions, pending=%+v", engine.State.PendingAction)
		}
		if !candidateContains(engine.State.PendingAction.Candidates, husk.InstanceID) || candidateContains(engine.State.PendingAction.Candidates, nonShadow.InstanceID) {
			t.Fatalf("2621108 should only offer low-cost shadow companions, candidates=%v", engine.State.PendingAction.Candidates)
		}
		resolvePendingSelection(t, engine, 0, husk.InstanceID)
		if containsCardInstance(p0.Hand, husk) || !containsCardInstance(p0.Graveyard, husk) {
			t.Fatalf("2621108 should discard the selected hand companion, hand=%v grave=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard))
		}
		if p0.Elements[model.ElementArcane] != startArcane+1 {
			t.Fatalf("2621108 should immediately resolve the discarded card deathrattle, arcane=%d want=%d", p0.Elements[model.ElementArcane], startArcane+1)
		}
	})
}

func TestRoyalConflictDeckAndTargetKindUtilityCards(t *testing.T) {
	t.Run("sky painter copies another low-cost air card enter effect", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		painter := placeUnit(baseCard(t, "1321115"), 0, 1, 1, engine)
		hummingbird := placeUnit(baseCard(t, "1321108"), 0, 0, 1, engine)
		fireCard := placeUnit(baseCard(t, "1121001"), 0, 2, 1, engine)
		drawA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		drawB := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		p0.Hand = nil
		p0.Deck = []*CardInstance{drawA, drawB}

		if err := (Card1321115SkyPainter{}).OnEnter(&EffectContext{Engine: engine, Source: painter, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1321115 enter: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "sky_painter_copy_enter" {
			t.Fatalf("1321115 should ask which enter effect to copy, pending=%+v", engine.State.PendingAction)
		}
		if !candidateContains(engine.State.PendingAction.Candidates, hummingbird.InstanceID) ||
			candidateContains(engine.State.PendingAction.Candidates, painter.InstanceID) ||
			candidateContains(engine.State.PendingAction.Candidates, fireCard.InstanceID) {
			t.Fatalf("1321115 should only offer another low-cost air card, candidates=%v", engine.State.PendingAction.Candidates)
		}
		resolvePendingSelection(t, engine, 0, hummingbird.InstanceID)
		if !containsCardInstance(p0.Hand, drawA) || !containsCardInstance(p0.Hand, drawB) {
			t.Fatalf("1321115 should copy hummingbird enter draw, hand=%v deck=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Deck))
		}
	})

	t.Run("magic moth may be drawn from deck after casting a focus spell", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		focus := readySkill(baseCard(t, "3021101"), 0)
		p0.Skills[0] = focus
		p0.Elements = cloneElements(map[string]int{model.ElementArcane: 2})
		moth := NewCardInstance(baseCard(t, "1021113"), 0, 1)
		p0.Hand = nil
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1), moth}
		p0.DrawCountThisTurn = 0
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": focus.InstanceID,
			"target_type": "unit",
			"target_col":  float64(target.Position.Col),
			"target_row":  float64(target.Position.Row),
		}}); err != nil {
			t.Fatalf("cast focus spell: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "magic_moth_draw" {
			t.Fatalf("1021113 should offer optional deck draw, pending=%+v", engine.State.PendingAction)
		}
		if containsCardInstance(p0.Hand, moth) || !containsCardInstance(p0.Deck, moth) || p0.DrawCountThisTurn != 0 || p1.Hand == nil {
			t.Fatalf("1021113 should stay in deck until accepted, hand=%v deck=%v draw_count=%d", cardsToInfo(p0.Hand), cardsToInfo(p0.Deck), p0.DrawCountThisTurn)
		}
		resolvePendingSelection(t, engine, 0, moth.InstanceID)
		if !containsCardInstance(p0.Hand, moth) || containsCardInstance(p0.Deck, moth) || p0.DrawCountThisTurn != 1 {
			t.Fatalf("1021113 accepted draw should move card to hand and count as draw, hand=%v deck=%v draw_count=%d", cardsToInfo(p0.Hand), cardsToInfo(p0.Deck), p0.DrawCountThisTurn)
		}

		nonFocusEngine := setupEffectTest(t)
		nonFocusP0 := nonFocusEngine.State.Players[0]
		nonFocus := readySkill(baseCard(t, "3021005"), 0)
		nonFocusMoth := NewCardInstance(baseCard(t, "1021113"), 0, 1)
		nonFocusP0.Deck = []*CardInstance{nonFocusMoth}
		nonFocusEngine.triggerMagicMothAfterFocusSpellCast(0, nonFocus)
		if nonFocusEngine.State.PendingAction != nil || containsCardInstance(nonFocusP0.Hand, nonFocusMoth) || !containsCardInstance(nonFocusP0.Deck, nonFocusMoth) {
			t.Fatalf("1021113 should ignore non-focus spells, pending=%+v hand=%v deck=%v", nonFocusEngine.State.PendingAction, cardsToInfo(nonFocusP0.Hand), cardsToInfo(nonFocusP0.Deck))
		}
	})

	t.Run("absolute purity counts top arcane cards as power and shuffles after use", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		arcaneOneness := readySkill(baseCard(t, "3011101"), 0)
		p0.Skills[0] = arcaneOneness
		arcaneA := NewCardInstance(baseCard(t, "3021005"), 0, 1)
		arcaneB := NewCardInstance(baseCard(t, "3021101"), 0, 1)
		fireStop := NewCardInstance(baseCard(t, "3121001"), 0, 1)
		p0.Deck = []*CardInstance{arcaneA, arcaneB, fireStop}

		stats := engine.skillContributionStats(0, arcaneOneness, nil, skillPurposeAttack)
		if stats.PowerBonus != arcaneOneness.Card.Power+2 {
			t.Fatalf("3011101 should gain power for consecutive top arcane cards, got=%d base=%d", stats.PowerBonus, arcaneOneness.Card.Power)
		}
		if err := (Card3011101AbsolutePurityArcaneOneness{}).OnSpellCast(&EffectContext{Engine: engine, Source: arcaneOneness, Target: arcaneOneness, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("3011101 spell cast: %v", err)
		}
		if len(p0.Deck) != 3 || !containsCardInstance(p0.Deck, arcaneA) || !containsCardInstance(p0.Deck, arcaneB) || !containsCardInstance(p0.Deck, fireStop) {
			t.Fatalf("3011101 should keep all revealed cards in deck after shuffle, deck=%v", cardsToInfo(p0.Deck))
		}

		blockedEngine := setupEffectTest(t)
		blockedP0 := blockedEngine.State.Players[0]
		blockedSkill := readySkill(baseCard(t, "3011101"), 0)
		blockedP0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "3121001"), 0, 1), NewCardInstance(baseCard(t, "3021005"), 0, 1)}
		if got := blockedEngine.skillContributionStats(0, blockedSkill, nil, skillPurposeAttack).PowerBonus; got != blockedSkill.Card.Power {
			t.Fatalf("3011101 should stop counting at first non-arcane card, got=%d base=%d", got, blockedSkill.Card.Power)
		}
	})

	t.Run("starfall silverleaf stores discarded cards and recycles one after hit", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		silverleaf := readySkill(baseCard(t, "3311102"), 0)
		discarded := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Skills[0] = silverleaf
		p0.Hand = []*CardInstance{discarded}

		if got := engine.discardHandCardAt(0, 0); got != discarded {
			t.Fatalf("discard setup should discard selected card, got=%v", cardToInfo(got))
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "starfall_silverleaf_store_discard" || engine.State.PendingAction.MinSelect != 1 {
			t.Fatalf("3311102 should require storing discarded card, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, discarded.InstanceID)
		if len(silverleaf.UnderCards) != 1 || silverleaf.UnderCards[0] != discarded || containsCardInstance(p0.Graveyard, discarded) {
			t.Fatalf("3311102 should move discarded card under itself, under=%v grave=%v", cardsToInfo(silverleaf.UnderCards), cardsToInfo(p0.Graveyard))
		}
		draw := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		p0.Deck = []*CardInstance{draw}
		beforeHand := len(p0.Hand)
		hitTarget := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
		if err := (Card3311102StarfallSilverleaf{}).OnSpellHit(&EffectContext{Engine: engine, Source: silverleaf, Target: hitTarget, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("3311102 hit: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "starfall_silverleaf_recycle_under_card" {
			t.Fatalf("3311102 should ask which under card to recycle, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, discarded.InstanceID)
		if len(silverleaf.UnderCards) != 0 || containsCardInstance(p0.Graveyard, discarded) || !containsCardInstance(p0.Deck, discarded) && !containsCardInstance(p0.Hand, discarded) || len(p0.Hand) != beforeHand+1 {
			t.Fatalf("3311102 should shuffle under card into deck and draw one, under=%v deck=%v hand=%v", cardsToInfo(silverleaf.UnderCards), cardsToInfo(p0.Deck), cardsToInfo(p0.Hand))
		}
	})

	t.Run("divine radiance skyward gains power from opponent hand and resets a hand on hit", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		skill := readySkill(baseCard(t, "3511101"), 0)
		p0.Skills[0] = skill
		p1.Hand = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 1, 1),
			NewCardInstance(baseCard(t, "1021002"), 1, 1),
		}
		for len(p1.Deck) < engine.State.HandLimit {
			p1.Deck = append(p1.Deck, NewCardInstance(baseCard(t, "1021001"), 1, 1))
		}

		stats := engine.skillContributionStats(0, skill, nil, skillPurposeAttack)
		if stats.PowerBonus != skill.Card.Power+2 {
			t.Fatalf("3511101 should gain power equal to opponent hand size, got=%d base=%d", stats.PowerBonus, skill.Card.Power)
		}
		hitTarget := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		if err := (Card3511101DivineRadianceSkyward{}).OnSpellHit(&EffectContext{Engine: engine, Source: skill, Target: hitTarget, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("3511101 hit: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "divine_radiance_reset_hand" {
			t.Fatalf("3511101 should offer hand reset after hit, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, "player:opponent")
		if len(p1.Hand) != engine.handLimitForPlayer(p1) || len(p1.Graveyard) != 2 {
			t.Fatalf("3511101 should discard target hand and draw to limit, hand=%d grave=%v", len(p1.Hand), cardsToInfo(p1.Graveyard))
		}
	})

	t.Run("rending impact distributes three damage among hit column units", func(t *testing.T) {
		engine := setupEffectTest(t)
		scroll := NewCardInstance(baseCard(t, "2321112"), 0, 1)
		targetA := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
		targetB := placeUnit(baseCard(t, "1021002"), 1, 1, 1, engine)
		offColumn := placeUnit(baseCard(t, "1021002"), 1, 0, 0, engine)
		if spellArea(scroll) != SpellAreaColumn {
			t.Fatalf("2321112 should be column range, got %s", spellArea(scroll))
		}
		if err := (Card2321112RendingImpactScroll{}).OnSpellHit(&EffectContext{
			Engine: engine, Source: scroll, Target: targetA, PlayerID: 0, OpponentID: 1,
			ExtraData: map[string]any{"affected_units": []*CardInstance{targetA, targetB}},
		}); err != nil {
			t.Fatalf("2321112 hit: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "rending_impact_distribute_damage" || candidateContains(engine.State.PendingAction.Candidates, offColumn.InstanceID) {
			t.Fatalf("2321112 should offer only hit units, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, targetA.InstanceID, targetB.InstanceID)
		if targetA.CurrentLife != targetA.Card.Life-2 || targetB.CurrentLife != targetB.Card.Life-1 || offColumn.CurrentLife != offColumn.Card.Life {
			t.Fatalf("2321112 should distribute 3 damage over selected hit units, a=%d b=%d off=%d", targetA.CurrentLife, targetB.CurrentLife, offColumn.CurrentLife)
		}
	})

	t.Run("blood feast targets friendly units, rewards on hit, and can bind to hero", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		feast := readySkill(baseCard(t, "3601101"), 0)
		ally := placeUnit(baseCard(t, "1021002"), 0, 1, 0, engine)
		enemy := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
		p0.Skills[0] = feast
		ownTargetOwner := 0
		if err := engine.validateSpellTarget(0, feast, SpellTarget{Type: "unit", Position: *ally.Position, OwnerID: &ownTargetOwner}); err != nil {
			t.Fatalf("3601101 should allow friendly unit targets: %v", err)
		}
		if err := engine.validateSpellTarget(0, feast, SpellTarget{Type: "unit", Position: *enemy.Position}); err == nil {
			t.Fatalf("3601101 should reject enemy targets")
		}
		if err := (Card3601101BloodFeast{}).OnSpellHit(&EffectContext{Engine: engine, Source: feast, Target: ally, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"spell_source": feast, "attacker": 0}}); err != nil {
			t.Fatalf("3601101 hit: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_feast_reward" {
			t.Fatalf("3601101 should ask for hit reward, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, "gain_shadow")
		if p0.Elements[model.ElementShadow] != 2 {
			t.Fatalf("3601101 should gain 2 shadow, elements=%v", p0.Elements)
		}
		otherSpell := readySkill(baseCard(t, "3621107"), 0)
		engine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, otherSpell, map[string]any{"spell_source": otherSpell, "attacker": 0})
		if engine.State.PendingAction != nil {
			t.Fatalf("3601101 should not reward another spell's hit, pending=%+v", engine.State.PendingAction)
		}
		p0.Elements[model.ElementShadow] = 1
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  feast.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("3601101 bind: %v", err)
		}
		if p0.Skills[0] != nil || len(p0.Hero.BoundSkills) != 1 || p0.Hero.BoundSkills[0] != feast || p0.Elements[model.ElementShadow] != 0 {
			t.Fatalf("3601101 should move from skill slot to hero bound skill, slot=%v bound=%v elements=%v", p0.Skills[0], cardsToInfo(p0.Hero.BoundSkills), p0.Elements)
		}
	})

	t.Run("pain scream scroll weakens unweakened enemy spells after friendly damage this turn", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		scroll := NewCardInstance(baseCard(t, "2621106"), 0, 1)
		ally := placeUnit(baseCard(t, "1021002"), 0, 1, 0, engine)
		ally.CurrentLife = 8
		enemyA := readySkill(baseCard(t, "3121001"), 1)
		enemyB := readySkill(baseCard(t, "3221002"), 1)
		alreadyWeak := readySkill(baseCard(t, "3321001"), 1)
		alreadyWeak.Statuses[StatusWeaken] = 1
		p1.Skills[0] = enemyA
		p1.Skills[1] = enemyB
		p1.Skills[2] = alreadyWeak

		if err := (Card2621106PainScreamScroll{}).OnUseItem(&EffectContext{Engine: engine, Source: scroll, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2621106 use: %v", err)
		}
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModPainScreamWeakenOnDamage {
			t.Fatalf("2621106 should create current-turn damage trigger modifier, modifiers=%+v", p0.TempModifiers)
		}
		engine.dealDamageWithExtra(ally, 2, 0, map[string]any{"damage_source": "effect", "attacker": 1})
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "pain_scream_weaken_enemy_spells" ||
			!candidateContains(engine.State.PendingAction.Candidates, enemyA.InstanceID) ||
			!candidateContains(engine.State.PendingAction.Candidates, enemyB.InstanceID) ||
			candidateContains(engine.State.PendingAction.Candidates, alreadyWeak.InstanceID) {
			t.Fatalf("2621106 should offer unweakened enemy spells only, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, enemyA.InstanceID, enemyB.InstanceID)
		if enemyA.Statuses[StatusWeaken] != 2 || enemyB.Statuses[StatusWeaken] != 2 || alreadyWeak.Statuses[StatusWeaken] != 1 {
			t.Fatalf("2621106 should weaken selected enemy spells by 2, a=%v b=%v weak=%v", enemyA.Statuses, enemyB.Statuses, alreadyWeak.Statuses)
		}
	})

	t.Run("protector sival prevents friendly damage after three friendly damage this turn", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		sival := NewCardInstance(baseCard(t, "4511101"), 0, 1)
		sival.Position = &Position{Col: 1, Row: 1}
		p0.Hero = sival
		p0.Units[1][1] = sival
		ally := placeUnit(baseCard(t, "1021002"), 0, 1, 0, engine)
		ally.CurrentLife = 8

		engine.dealDamageWithExtra(ally, 2, 0, map[string]any{"damage_source": "effect", "attacker": 1})
		if engine.State.PendingAction != nil || sival.UltimateUsed {
			t.Fatalf("4511101 should wait until friendly damage reaches three, pending=%+v used=%v", engine.State.PendingAction, sival.UltimateUsed)
		}
		engine.dealDamageWithExtra(ally, 1, 0, map[string]any{"damage_source": "effect", "attacker": 1})
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "protector_sival_prevent_all_damage" {
			t.Fatalf("4511101 should offer ultimate after three friendly damage, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, sival.InstanceID)
		if !sival.UltimateUsed || sival.Statuses[protectorSivalPreventionUntilStatus] < engine.State.TurnNumber {
			t.Fatalf("4511101 should arm damage prevention, used=%v statuses=%v", sival.UltimateUsed, sival.Statuses)
		}
		allyLife := ally.CurrentLife
		sivalLife := sival.CurrentLife
		engine.dealDamageWithExtra(ally, 2, 0, map[string]any{"damage_source": "effect", "attacker": 1})
		engine.dealDamageWithExtra(sival, 2, 0, map[string]any{"damage_source": "effect", "attacker": 1})
		if ally.CurrentLife != allyLife || sival.CurrentLife != sivalLife {
			t.Fatalf("4511101 should prevent later friendly and self damage, ally=%d/%d sival=%d/%d", ally.CurrentLife, allyLife, sival.CurrentLife, sivalLife)
		}
	})

	t.Run("divine fire staff consumes itself to permanently empower a fire spell after hit", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		staff := NewCardInstance(baseCard(t, "2111101"), 0, 1)
		staff.IsHorizontal = false
		fireball := readySkill(baseCard(t, "3121001"), 0)
		p0.Equipment[0] = staff
		p0.Skills[0] = fireball
		startPower := fireball.PowerBonus
		startFire := p0.Elements[model.ElementFire]

		engine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, fireball, map[string]any{"attacker": 0, "spell_source": fireball})
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "divine_fire_staff_empower_spell" {
			t.Fatalf("2111101 should offer to consume after friendly fire spell hit, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, staff.InstanceID)
		if !staff.IsHorizontal || p0.Elements[model.ElementFire] != startFire+staff.Card.ElementsGain[model.ElementFire] {
			t.Fatalf("2111101 should consume itself and gain its load, horizontal=%v elements=%v", staff.IsHorizontal, p0.Elements)
		}
		if fireball.PowerBonus != startPower+1 || !engine.skillHasPierce(0, fireball) {
			t.Fatalf("2111101 should permanently give +1 power and pierce, power=%d statuses=%v", fireball.PowerBonus, fireball.Statuses)
		}

		enemyEngine := setupEffectTest(t)
		enemyP0 := enemyEngine.State.Players[0]
		enemyStaff := NewCardInstance(baseCard(t, "2111101"), 0, 1)
		enemyStaff.IsHorizontal = false
		enemyFireball := readySkill(baseCard(t, "3121001"), 1)
		enemyP0.Equipment[0] = enemyStaff
		enemyEngine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, enemyFireball, map[string]any{"attacker": 1, "spell_source": enemyFireball})
		if enemyEngine.State.PendingAction != nil || enemyStaff.IsHorizontal || enemyFireball.Statuses[permanentPierceStatus] > 0 {
			t.Fatalf("2111101 should ignore enemy fire spell hits, pending=%+v horizontal=%v statuses=%v", enemyEngine.State.PendingAction, enemyStaff.IsHorizontal, enemyFireball.Statuses)
		}
	})

	t.Run("mist mask discards hand cards to reduce enemy spell hit damage", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		mask := NewCardInstance(baseCard(t, "2321109"), 0, 1)
		discardA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		discardB := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		enemySpell := readySkill(baseCard(t, "3021005"), 1)
		p0.Equipment[0] = mask
		p0.Hand = []*CardInstance{discardA, discardB}
		damage := 3

		engine.triggerFieldEffectsWithData(TriggerOnSpellHitBeforeDamage, 0, enemySpell, map[string]any{
			"attacker": 1, "spell_source": enemySpell, "damage_ptr": &damage, "damage": damage,
		})
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "mist_mask_discard_reduce_spell_attack" {
			t.Fatalf("2321109 should offer to discard after enemy spell hit, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, discardA.InstanceID, discardB.InstanceID)
		if damage != 1 || !mask.UltimateUsed || !containsCardInstance(p0.Graveyard, discardA) || !containsCardInstance(p0.Graveyard, discardB) {
			t.Fatalf("2321109 should reduce damage by discarded count, damage=%d used=%v grave=%v", damage, mask.UltimateUsed, cardsToInfo(p0.Graveyard))
		}

		friendlyEngine := setupEffectTest(t)
		friendlyP0 := friendlyEngine.State.Players[0]
		friendlyMask := NewCardInstance(baseCard(t, "2321109"), 0, 1)
		friendlySpell := readySkill(baseCard(t, "3021005"), 0)
		friendlyP0.Equipment[0] = friendlyMask
		friendlyP0.Hand = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
		friendlyDamage := 2
		friendlyEngine.triggerFieldEffectsWithData(TriggerOnSpellHitBeforeDamage, 0, friendlySpell, map[string]any{
			"attacker": 0, "spell_source": friendlySpell, "damage_ptr": &friendlyDamage, "damage": friendlyDamage,
		})
		if friendlyEngine.State.PendingAction != nil || friendlyDamage != 2 || friendlyMask.UltimateUsed {
			t.Fatalf("2321109 should ignore friendly spell hits, pending=%+v damage=%d used=%v", friendlyEngine.State.PendingAction, friendlyDamage, friendlyMask.UltimateUsed)
		}

		liveEngine := setupEffectTest(t)
		liveP0 := liveEngine.State.Players[0]
		liveP1 := liveEngine.State.Players[1]
		liveSpell := readySkill(baseCard(t, "3121001"), 0)
		liveMask := NewCardInstance(baseCard(t, "2321109"), 1, 1)
		liveDiscard := NewCardInstance(baseCard(t, "1021001"), 1, 1)
		liveTarget := placeUnit(baseCard(t, "1021001"), 1, 1, 0, liveEngine)
		liveP0.Skills[0] = liveSpell
		liveP0.Elements[model.ElementFire] = 1
		liveP1.Equipment[0] = liveMask
		liveP1.Hand = []*CardInstance{liveDiscard}
		liveStartLife := liveTarget.CurrentLife
		if err := liveEngine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": liveSpell.InstanceID,
			"target_type": "unit",
			"target_col":  float64(liveTarget.Position.Col),
			"target_row":  float64(liveTarget.Position.Row),
		}}); err != nil {
			t.Fatalf("live 2321109 cast: %v", err)
		}
		if err := liveEngine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("live 2321109 no defend: %v", err)
		}
		if liveEngine.State.PendingAction == nil || liveEngine.State.PendingAction.Type != "mist_mask_discard_reduce_spell_attack" || liveTarget.CurrentLife != liveStartLife {
			t.Fatalf("2321109 should pause before damage, pending=%+v life=%d start=%d", liveEngine.State.PendingAction, liveTarget.CurrentLife, liveStartLife)
		}
		resolvePendingSelection(t, liveEngine, 1, liveDiscard.InstanceID)
		if liveTarget.CurrentLife != liveStartLife || !liveMask.UltimateUsed || liveEngine.State.PendingSpell != nil {
			t.Fatalf("2321109 live reduction should prevent damage and finish spell, life=%d start=%d used=%v pending_spell=%+v", liveTarget.CurrentLife, liveStartLife, liveMask.UltimateUsed, liveEngine.State.PendingSpell)
		}
	})

	t.Run("radiant city priest converts enemy spell damage to burn before damage resolves", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		priest := placeUnit(baseCard(t, "1521105"), 1, 0, 1, engine)
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		fireball := readySkill(baseCard(t, "3121001"), 0)
		p0.Skills[0] = fireball
		p0.Elements[model.ElementFire] = 1
		startLife := target.CurrentLife

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": fireball.InstanceID,
			"target_type": "unit",
			"target_col":  float64(target.Position.Col),
			"target_row":  float64(target.Position.Row),
		}}); err != nil {
			t.Fatalf("1521105 cast setup: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("1521105 no defend: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "radiant_city_priest_convert_damage_to_burn" || target.CurrentLife != startLife {
			t.Fatalf("1521105 should pause before damage, pending=%+v life=%d start=%d", engine.State.PendingAction, target.CurrentLife, startLife)
		}
		resolvePendingSelection(t, engine, 1, priest.InstanceID)
		if target.CurrentLife != startLife || target.Statuses[StatusBurn] != 1 || !priest.UltimateUsed || engine.State.PendingSpell != nil {
			t.Fatalf("1521105 should convert damage to burn and finish spell, life=%d burn=%d used=%v pending_spell=%+v", target.CurrentLife, target.Statuses[StatusBurn], priest.UltimateUsed, engine.State.PendingSpell)
		}

		friendlyEngine := setupEffectTest(t)
		friendlyP0 := friendlyEngine.State.Players[0]
		friendlyPriest := placeUnit(baseCard(t, "1521105"), 0, 0, 1, friendlyEngine)
		friendlyTarget := placeUnit(baseCard(t, "1021001"), 1, 1, 0, friendlyEngine)
		friendlySpell := readySkill(baseCard(t, "3121001"), 0)
		friendlyP0.Skills[0] = friendlySpell
		friendlyP0.Elements[model.ElementFire] = 1
		if err := friendlyEngine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": friendlySpell.InstanceID,
			"target_type": "unit",
			"target_col":  float64(friendlyTarget.Position.Col),
			"target_row":  float64(friendlyTarget.Position.Row),
		}}); err != nil {
			t.Fatalf("friendly 1521105 cast setup: %v", err)
		}
		if err := friendlyEngine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("friendly 1521105 no defend: %v", err)
		}
		if friendlyEngine.State.PendingAction != nil || friendlyPriest.UltimateUsed {
			t.Fatalf("1521105 should ignore friendly spell hits, pending=%+v used=%v", friendlyEngine.State.PendingAction, friendlyPriest.UltimateUsed)
		}
	})

	t.Run("frost robe freezes enemies after a friendly water unit takes enemy spell damage", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		robe := NewCardInstance(baseCard(t, "2221106"), 0, 1)
		waterAlly := placeUnit(baseCard(t, "1221001"), 0, 1, 1, engine)
		frontEnemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		backEnemy := placeUnit(baseCard(t, "1021002"), 1, 1, 1, engine)
		enemySpell := readySkill(baseCard(t, "3021005"), 1)
		p0.Equipment[0] = robe

		engine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, enemySpell, map[string]any{
			"attacker":                           1,
			"spell_source":                       enemySpell,
			"affected_units":                     []*CardInstance{waterAlly},
			"actual_friendly_damage_by_instance": map[string]int{waterAlly.InstanceID: 1},
		})
		if !robe.UltimateUsed || frontEnemy.Statuses[StatusFreeze] != 1 || backEnemy.Statuses[StatusFreeze] != 0 {
			t.Fatalf("2221106 should freeze enemies in spell range after water ally damage, used=%v front=%v back=%v", robe.UltimateUsed, frontEnemy.Statuses, backEnemy.Statuses)
		}

		noDamageEngine := setupEffectTest(t)
		noDamageP0 := noDamageEngine.State.Players[0]
		noDamageRobe := NewCardInstance(baseCard(t, "2221106"), 0, 1)
		noDamageAlly := placeUnit(baseCard(t, "1221001"), 0, 1, 1, noDamageEngine)
		noDamageEnemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, noDamageEngine)
		noDamageP0.Equipment[0] = noDamageRobe
		noDamageEngine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, readySkill(baseCard(t, "3021005"), 1), map[string]any{
			"attacker":                           1,
			"affected_units":                     []*CardInstance{noDamageAlly},
			"actual_friendly_damage_by_instance": map[string]int{},
		})
		if noDamageRobe.UltimateUsed || noDamageEnemy.Statuses[StatusFreeze] != 0 {
			t.Fatalf("2221106 should ignore hits with no actual friendly water damage, used=%v enemy=%v", noDamageRobe.UltimateUsed, noDamageEnemy.Statuses)
		}
	})

	t.Run("regroup triggers after enemy spell hit and buffs a friendly companion", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		regroup := readySkill(baseCard(t, "3521109"), 0)
		ally := placeUnit(baseCard(t, "1021001"), 0, 1, 1, engine)
		enemySpell := readySkill(baseCard(t, "3021005"), 1)
		p0.Skills[0] = regroup
		startLife := ally.CurrentLife
		startLight := effectiveElementsGain(ally)[model.ElementLight]

		engine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, enemySpell, map[string]any{"attacker": 1, "spell_source": enemySpell})
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "regroup_buff_companion" {
			t.Fatalf("3521109 should offer to buff a companion after enemy spell hit, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, ally.InstanceID)
		if !regroup.IsHorizontal || ally.CurrentLife != startLife+1 || ally.Statuses["max_life_bonus"] != 1 || effectiveElementsGain(ally)[model.ElementLight] != startLight+1 {
			t.Fatalf("3521109 should tap and grant +1 life/load, horizontal=%v life=%d statuses=%v load=%v", regroup.IsHorizontal, ally.CurrentLife, ally.Statuses, effectiveElementsGain(ally))
		}

		friendlyEngine := setupEffectTest(t)
		friendlyP0 := friendlyEngine.State.Players[0]
		friendlyRegroup := readySkill(baseCard(t, "3521109"), 0)
		friendlyP0.Skills[0] = friendlyRegroup
		friendlyEngine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, readySkill(baseCard(t, "3021005"), 0), map[string]any{"attacker": 0})
		if friendlyEngine.State.PendingAction != nil || friendlyRegroup.IsHorizontal {
			t.Fatalf("3521109 should ignore friendly spell hits, pending=%+v horizontal=%v", friendlyEngine.State.PendingAction, friendlyRegroup.IsHorizontal)
		}
	})

	t.Run("sin chooses a companion kind and gains power plus pierce against that kind", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		sin := readySkill(baseCard(t, "3521104"), 0)
		p0.Skills[0] = sin
		front := placeUnit(baseCard(t, "1121105"), 1, 1, 0, engine)
		humanBack := placeUnit(baseCard(t, "1021001"), 1, 1, 1, engine)

		if err := (Card3521104Sin{}).OnEnter(&EffectContext{Engine: engine, Source: sin, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("3521104 enter: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "sin_choose_companion_kind" || !candidateContains(engine.State.PendingAction.Candidates, "巫师") {
			t.Fatalf("3521104 should ask for a companion kind, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, "巫师")
		target := SpellTarget{Type: "unit", Position: *humanBack.Position}
		if err := engine.validateSpellTarget(0, sin, target); err != nil {
			t.Fatalf("3521104 should pierce to matching back-row target despite front blocker %s: %v", front.InstanceID, err)
		}
		if got := engine.effectiveSpellPower(0, sin, nil, target); got != sin.Card.Power+2 {
			t.Fatalf("3521104 should gain +2 power against chosen kind, got=%d want=%d", got, sin.Card.Power+2)
		}

		mismatchEngine := setupEffectTest(t)
		mismatchP0 := mismatchEngine.State.Players[0]
		mismatchSin := readySkill(baseCard(t, "3521104"), 0)
		mismatchSin.Statuses[sinTargetTagStatusPrefix+"恶魔"] = 1
		mismatchP0.Skills[0] = mismatchSin
		placeUnit(baseCard(t, "1121105"), 1, 1, 0, mismatchEngine)
		mismatchBack := placeUnit(baseCard(t, "1021001"), 1, 1, 1, mismatchEngine)
		mismatchTarget := SpellTarget{Type: "unit", Position: *mismatchBack.Position}
		if err := mismatchEngine.validateSpellTarget(0, mismatchSin, mismatchTarget); err == nil {
			t.Fatal("3521104 should not pierce to a non-matching kind")
		}
		if got := mismatchEngine.effectiveSpellPower(0, mismatchSin, nil, mismatchTarget); got != mismatchSin.Card.Power {
			t.Fatalf("3521104 should not gain power against non-matching kind, got=%d want=%d", got, mismatchSin.Card.Power)
		}
	})
}

func TestRoyalConflictPrintedBoundSkills(t *testing.T) {
	cases := []struct {
		name        string
		hostNumber  string
		boundNumber string
		equipment   bool
	}{
		{name: "1011103 弈者 binds 入局", hostNumber: "1011103", boundNumber: "3001101"},
		{name: "2511102 五虹之环 binds 五虹之束", hostNumber: "2511102", boundNumber: "3501101", equipment: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			engine := setupReportedBugEngine(t)
			p0 := engine.State.Players[0]
			host := NewCardInstance(baseCard(t, tc.hostNumber), 0, 1)
			if tc.equipment {
				p0.Equipment[0] = host
				host.SlotIndex = 0
			} else {
				p0.Units[0][0] = host
				host.Position = &Position{Col: 0, Row: 0}
			}

			engine.triggerEffects(TriggerOnEnter, host, nil, nil)

			if len(p0.SkillPool) != 0 {
				t.Fatalf("bound skill should not enter skill pool, pool=%v", cardsToInfo(p0.SkillPool))
			}
			for i, skill := range p0.Skills {
				if skill != nil {
					t.Fatalf("bound skill should not occupy skill slot %d: %v", i, cardToInfo(skill))
				}
			}
			if len(host.BoundSkills) != 1 || host.BoundSkills[0].Card.Number != tc.boundNumber {
				t.Fatalf("expected bound skill %s on host, bound=%v", tc.boundNumber, cardsToInfo(host.BoundSkills))
			}
			if host.BoundSkills[0].SlotIndex != -1 || !host.BoundSkills[0].IsHorizontal {
				t.Fatalf("bound skill should enter horizontal without a slot, bound=%v", cardToInfo(host.BoundSkills[0]))
			}
			info := cardToInfo(host)
			bound, ok := info["bound_skills"].([]map[string]any)
			if !ok || len(bound) != 1 || bound[0]["number"] != tc.boundNumber {
				t.Fatalf("card info should expose bound skill, info=%+v", info["bound_skills"])
			}
		})
	}

	t.Run("equipment replacement clears bound skills from old host", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		oldRing := NewCardInstance(baseCard(t, "2511102"), 0, 1)
		p0.Equipment[0] = oldRing
		oldRing.SlotIndex = 0
		oldRing.IsHorizontal = false
		engine.triggerEffects(TriggerOnEnter, oldRing, nil, nil)
		if len(oldRing.BoundSkills) != 1 {
			t.Fatalf("old ring should bind skill before replacement, bound=%v", cardsToInfo(oldRing.BoundSkills))
		}

		newRing := NewCardInstance(baseCard(t, "2511102"), 0, 1)
		p0.Hand = []*CardInstance{newRing}
		setAllElements(p0, 10)
		if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
			"instance_id": newRing.InstanceID,
			"replace_id":  oldRing.InstanceID,
		}}); err != nil {
			t.Fatalf("replace five rainbow ring: %v", err)
		}
		if len(p0.Graveyard) != 1 || p0.Graveyard[0] != oldRing {
			t.Fatalf("old ring should move to graveyard, grave=%v", cardsToInfo(p0.Graveyard))
		}
		if len(oldRing.BoundSkills) != 0 {
			t.Fatalf("bound skills should disappear when equipment host leaves, bound=%v", cardsToInfo(oldRing.BoundSkills))
		}
		if info := cardToInfo(oldRing); info["bound_skills"] != nil {
			t.Fatalf("graveyard equipment should not expose old bound skills, info=%+v", info["bound_skills"])
		}
		if len(newRing.BoundSkills) != 1 || newRing.BoundSkills[0].Card.Number != "3501101" {
			t.Fatalf("new ring should bind its own skill after replacement, bound=%v", cardsToInfo(newRing.BoundSkills))
		}
	})

	t.Run("enemy equipment destruction clears bound skills from old host", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p1 := engine.State.Players[1]
		ring := NewCardInstance(baseCard(t, "2511102"), 1, 1)
		p1.Equipment[0] = ring
		ring.SlotIndex = 0
		engine.triggerEffects(TriggerOnEnter, ring, nil, nil)
		if len(ring.BoundSkills) != 1 {
			t.Fatalf("ring should bind skill before destruction, bound=%v", cardsToInfo(ring.BoundSkills))
		}

		if !engine.destroyEnemyEquipment(0, ring.InstanceID) {
			t.Fatal("destroyEnemyEquipment should destroy opponent ring")
		}
		if len(p1.Graveyard) != 1 || p1.Graveyard[0] != ring {
			t.Fatalf("destroyed ring should move to opponent graveyard, grave=%v", cardsToInfo(p1.Graveyard))
		}
		if ring.SlotIndex != -1 || len(ring.BoundSkills) != 0 {
			t.Fatalf("destroyed equipment should clear slot and bound skills, slot=%d bound=%v", ring.SlotIndex, cardsToInfo(ring.BoundSkills))
		}
	})
}

func TestRoyalConflictRedMoonBasics(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	front := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	back := placeUnit(baseCard(t, "1021002"), 1, 1, 2, engine)
	redMoon := readySkill(baseCard(t, "3611101"), 0)
	willErosion := readySkill(baseCard(t, "3621107"), 0)
	p0.Skills[0] = redMoon
	p0.Skills[1] = willErosion

	if engine.redMoonActive(0) {
		t.Fatal("red moon should not be active before duration is set")
	}
	if err := engine.validateSpellTarget(0, willErosion, SpellTarget{Type: "unit", Position: *back.Position}); err == nil {
		t.Fatalf("will erosion should not pierce back-row target before red moon")
	}
	if got := engine.effectiveSpellPower(0, willErosion, nil, SpellTarget{Type: "unit", Position: *front.Position}); got != willErosion.Card.Power {
		t.Fatalf("will erosion should use base power before red moon, got %d", got)
	}

	redMoon.Statuses[StatusAbilityDuration] = 1
	if !engine.redMoonActive(0) {
		t.Fatal("red moon should be active while its ability duration mark is present")
	}
	if err := engine.validateSpellTarget(0, willErosion, SpellTarget{Type: "unit", Position: *back.Position}); err != nil {
		t.Fatalf("will erosion should pierce during red moon: %v", err)
	}
	if info := engine.cardToInfoForPlayer(p0, willErosion); info["has_pierce"] != true {
		t.Fatalf("will erosion should expose dynamic pierce during red moon, info=%+v", info)
	}
	if got := engine.effectiveSpellPower(0, willErosion, nil, SpellTarget{Type: "unit", Position: *back.Position}); got != willErosion.Card.Power+3 {
		t.Fatalf("red moon should give shadow spell +2 and will erosion +1, got %d", got)
	}

	beast := placeUnit(baseCard(t, "1621110"), 0, 0, 0, engine)
	if got := engine.effectiveSpellPower(0, willErosion, nil, SpellTarget{Type: "unit", Position: *back.Position}); got != willErosion.Card.Power+5 {
		t.Fatalf("scarlet beast should add +2 during red moon, got %d with beast %v", got, cardToInfo(beast))
	}

	delete(redMoon.Statuses, StatusAbilityDuration)
	if err := engine.validateSpellTarget(0, willErosion, SpellTarget{Type: "unit", Position: *back.Position}); err == nil {
		t.Fatalf("will erosion should stop piercing after red moon ends")
	}
}

func TestRoyalConflictRedMoonMarkersAndSevianaTransform(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	seviana := placeUnit(baseCard(t, "1611101"), 0, 0, 0, engine)
	redMoon := readySkill(baseCard(t, "3611101"), 0)
	willErosion := readySkill(baseCard(t, "3621107"), 0)
	p0.Skills[0] = redMoon
	p0.Skills[1] = willErosion

	engine.triggerEffects(TriggerOnEnter, seviana, nil, nil)
	if redMoon.Statuses[redMoonMarkerStatus] != 1 {
		t.Fatalf("Seviana should place one red moon marker on enter, statuses=%v", redMoon.Statuses)
	}
	engine.triggerPrayerAbilities(0)
	if redMoon.Statuses[redMoonMarkerStatus] != 2 {
		t.Fatalf("Seviana prayer should place another red moon marker, statuses=%v", redMoon.Statuses)
	}
	if seviana.Card.Number != "1611101" {
		t.Fatalf("Seviana should not transform before red moon is active, card=%s", seviana.Card.Number)
	}
	seviana.CurrentLife = 1
	seviana.IsHorizontal = true
	seviana.Statuses[StatusBurn] = 2

	p0.Elements[model.ElementShadow] = 1
	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": redMoon.InstanceID,
		"target_type": "none",
	}}); err != nil {
		t.Fatalf("cast red moon: %v", err)
	}
	if seviana.Card.Number != "1601101" {
		t.Fatalf("Seviana should become blood shadow body when red moon is cast, card=%s", seviana.Card.Number)
	}
	if seviana.CurrentLife != seviana.Card.Life || !seviana.IsHorizontal || seviana.Statuses[StatusBurn] != 0 {
		t.Fatalf("blood shadow body should refresh life/statuses and preserve horizontal state, life=%d horizontal=%v statuses=%v", seviana.CurrentLife, seviana.IsHorizontal, seviana.Statuses)
	}
	if got := engine.effectiveSpellPower(0, willErosion, nil); got != willErosion.Card.Power+5 {
		t.Fatalf("two red moon markers should add +2 to other shadow spell during red moon, got %d", got)
	}
	if got := engine.effectiveSpellPower(0, redMoon, nil); got != 2 {
		t.Fatalf("red moon markers should not buff red moon itself, got %d", got)
	}

	engine.processAbilityDurations(p0)
	if seviana.Card.Number != "1611101" {
		t.Fatalf("blood shadow body should revert to Seviana after red moon ends, card=%s", seviana.Card.Number)
	}
	if seviana.CurrentLife != seviana.Card.Life || seviana.IsHorizontal {
		t.Fatalf("reverted Seviana should reset, life=%d horizontal=%v", seviana.CurrentLife, seviana.IsHorizontal)
	}

	redMoon.Statuses[StatusAbilityDuration] = 1
	engine.refreshRedMoonState(0)
	if seviana.Card.Number != "1601101" {
		t.Fatalf("Seviana should transform again while red moon remains active, card=%s", seviana.Card.Number)
	}
	engine.addStatus(redMoon, StatusPetrify, 1)
	if seviana.Card.Number != "1611101" {
		t.Fatalf("petrified red moon should revert blood shadow body, card=%s", seviana.Card.Number)
	}
	engine.processEndOfTurnStatuses(p0)
	if seviana.Card.Number != "1601101" {
		t.Fatalf("red moon should transform Seviana again after petrify expires, card=%s", seviana.Card.Number)
	}

	if !engine.removeFieldCardFromGameByID(redMoon.InstanceID) {
		t.Fatal("remove red moon from field")
	}
	if seviana.Card.Number != "1611101" {
		t.Fatalf("removing red moon should revert blood shadow body, card=%s", seviana.Card.Number)
	}
}

func TestRoyalConflictBloodShadowBodySpendsRedMoonMarkerForExtraTarget(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	body := placeUnit(baseCard(t, "1601101"), 0, 0, 0, engine)
	redMoon := readySkill(baseCard(t, "3611101"), 0)
	redMoon.Statuses[redMoonMarkerStatus] = 1
	redMoon.Statuses[StatusAbilityDuration] = 1
	p0.Skills[0] = redMoon
	engine.refreshRedMoonState(0)
	spell := readySkill(baseCard(t, "3121002"), 0)
	p0.Skills[1] = spell
	p0.Elements[model.ElementFire] = 10
	front := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	back := placeUnit(baseCard(t, "1021001"), 1, 1, 2, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  body.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use blood shadow body ability: %v", err)
	}
	if redMoon.Statuses[redMoonMarkerStatus] != 0 || len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextSpellExtraTarget {
		t.Fatalf("blood shadow body should spend one marker and arm extra target, markers=%v modifiers=%v", redMoon.Statuses, p0.TempModifiers)
	}
	if !p0.TempModifiers[0].AllowSameTarget {
		t.Fatalf("blood shadow body extra target should allow the same target, modifier=%+v", p0.TempModifiers[0])
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id":      spell.InstanceID,
		"target_type":      "unit",
		"target_col":       float64(front.Position.Col),
		"target_row":       float64(front.Position.Row),
		"extra_target_col": float64(back.Position.Col),
		"extra_target_row": float64(back.Position.Row),
	}}); err != nil {
		t.Fatalf("cast with blood shadow body extra target: %v", err)
	}
	if engine.State.PendingSpell == nil || len(engine.State.PendingSpell.ExtraTargets) != 1 {
		t.Fatalf("blood shadow body should add one extra target, pending=%+v", engine.State.PendingSpell)
	}
	if engine.hasNextDriveSpellExtraTarget(p0, spell) {
		t.Fatalf("blood shadow body extra target modifier should be consumed, modifiers=%v", p0.TempModifiers)
	}

	sameTargetEngine := setupReportedBugEngine(t)
	sameP0 := sameTargetEngine.State.Players[0]
	sameBody := placeUnit(baseCard(t, "1601101"), 0, 0, 0, sameTargetEngine)
	sameRedMoon := readySkill(baseCard(t, "3611101"), 0)
	sameRedMoon.Statuses[redMoonMarkerStatus] = 1
	sameRedMoon.Statuses[StatusAbilityDuration] = 1
	sameP0.Skills[0] = sameRedMoon
	sameTargetEngine.refreshRedMoonState(0)
	sameSpell := readySkill(baseCard(t, "3121002"), 0)
	sameP0.Skills[1] = sameSpell
	sameP0.Elements[model.ElementFire] = 10
	sameFront := placeUnit(baseCard(t, "1021001"), 1, 1, 0, sameTargetEngine)
	if err := sameTargetEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  sameBody.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use blood shadow body same-target setup: %v", err)
	}
	if err := sameTargetEngine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id":      sameSpell.InstanceID,
		"target_type":      "unit",
		"target_col":       float64(sameFront.Position.Col),
		"target_row":       float64(sameFront.Position.Row),
		"extra_target_col": float64(sameFront.Position.Col),
		"extra_target_row": float64(sameFront.Position.Row),
	}}); err != nil {
		t.Fatalf("blood shadow body extra target should allow choosing the same target: %v", err)
	}
	if sameTargetEngine.State.PendingSpell == nil || len(sameTargetEngine.State.PendingSpell.ExtraTargets) != 1 {
		t.Fatalf("blood shadow body should add one same extra target, pending=%+v", sameTargetEngine.State.PendingSpell)
	}

	noExtraEngine := setupReportedBugEngine(t)
	noExtraP0 := noExtraEngine.State.Players[0]
	noExtraBody := placeUnit(baseCard(t, "1601101"), 0, 0, 0, noExtraEngine)
	noExtraRedMoon := readySkill(baseCard(t, "3611101"), 0)
	noExtraRedMoon.Statuses[redMoonMarkerStatus] = 1
	noExtraRedMoon.Statuses[StatusAbilityDuration] = 1
	noExtraP0.Skills[0] = noExtraRedMoon
	noExtraEngine.refreshRedMoonState(0)
	noExtraSpell := readySkill(baseCard(t, "3621101"), 0)
	noExtraP0.Skills[1] = noExtraSpell
	noExtraP0.Elements[model.ElementShadow] = 10
	noExtraFriend := placeUnit(baseCard(t, "1021001"), 0, 1, 0, noExtraEngine)
	ownerID := 0
	if err := noExtraEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  noExtraBody.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use blood shadow body no-extra setup: %v", err)
	}
	if err := noExtraEngine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id":  noExtraSpell.InstanceID,
		"target_type":  "unit",
		"target_owner": float64(ownerID),
		"target_col":   float64(noExtraFriend.Position.Col),
		"target_row":   float64(noExtraFriend.Position.Row),
	}}); err != nil {
		t.Fatalf("blood shadow body modifier should be consumed by the next spell even without extra target: %v", err)
	}
	if noExtraEngine.hasNextDriveSpellExtraTarget(noExtraP0, noExtraSpell) {
		t.Fatalf("blood shadow body next-spell extra target should not persist after next spell, modifiers=%v", noExtraP0.TempModifiers)
	}

	combinedEngine := setupReportedBugEngine(t)
	combinedP0 := combinedEngine.State.Players[0]
	combinedBody := placeUnit(baseCard(t, "1601101"), 0, 0, 0, combinedEngine)
	combinedChain := placeUnit(baseCard(t, "2321101"), 0, 1, 0, combinedEngine)
	combinedRedMoon := readySkill(baseCard(t, "3611101"), 0)
	combinedRedMoon.Statuses[redMoonMarkerStatus] = 1
	combinedRedMoon.Statuses[StatusAbilityDuration] = 1
	combinedP0.Skills[0] = combinedRedMoon
	combinedEngine.refreshRedMoonState(0)
	combinedDrive := readySkill(baseCard(t, "3121002"), 0)
	combinedP0.Skills[1] = combinedDrive
	combinedP0.Elements[model.ElementFire] = 10
	combinedTarget := placeUnit(baseCard(t, "1021001"), 1, 1, 0, combinedEngine)
	if err := combinedEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  combinedChain.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use thunder chain before blood shadow body: %v", err)
	}
	if err := combinedEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  combinedBody.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use blood shadow body with thunder chain armed: %v", err)
	}
	if err := combinedEngine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": combinedDrive.InstanceID,
		"target_type": "unit",
		"target_col":  float64(combinedTarget.Position.Col),
		"target_row":  float64(combinedTarget.Position.Row),
	}}); err != nil {
		t.Fatalf("cast drive spell without extra target while both modifiers exist: %v", err)
	}
	var driveRemaining, spellRemaining int
	for _, modifier := range combinedP0.TempModifiers {
		switch modifier.Type {
		case TempModNextDriveSpellExtraTarget:
			driveRemaining = modifier.RemainingUses
		case TempModNextSpellExtraTarget:
			spellRemaining = modifier.RemainingUses
		}
	}
	if driveRemaining != 1 || spellRemaining != 0 {
		t.Fatalf("no-extra drive spell should consume only blood shadow body modifier, modifiers=%v", combinedP0.TempModifiers)
	}
}

func TestRoyalConflictWillErosionHasInherentSameExtraTarget(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	willErosion := readySkill(baseCard(t, "3621107"), 0)
	p0.Skills[0] = willErosion
	p0.Elements[model.ElementShadow] = 10
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	target.CurrentLife = 20

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id":      willErosion.InstanceID,
		"target_type":      "unit",
		"target_col":       float64(target.Position.Col),
		"target_row":       float64(target.Position.Row),
		"extra_target_col": float64(target.Position.Col),
		"extra_target_row": float64(target.Position.Row),
	}}); err != nil {
		t.Fatalf("3621107 should allow an inherent same extra target: %v", err)
	}
	if engine.State.PendingSpell == nil || len(engine.State.PendingSpell.ExtraTargets) != 1 {
		t.Fatalf("3621107 should add one same extra target, pending=%+v", engine.State.PendingSpell)
	}
}

func TestRoyalConflictScarletWingsTriggersAfterRedMoon(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	wing := placeUnit(baseCard(t, "1621109"), 0, 0, 0, engine)
	redMoon := readySkill(baseCard(t, "3611101"), 0)
	p0.Skills[0] = redMoon
	p0.Elements[model.ElementShadow] = 10
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	target.CurrentLife = 5
	startWingLife := wing.CurrentLife

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": redMoon.InstanceID,
		"target_type": "none",
	}}); err != nil {
		t.Fatalf("cast red moon for scarlet wings: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "scarlet_wings_red_moon_damage" {
		t.Fatalf("scarlet wings should trigger after red moon, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, target.InstanceID)
	if target.CurrentLife != 4 || wing.CurrentLife != startWingLife+1 {
		t.Fatalf("scarlet wings should damage target and gain life, target=%d wing=%d/%d", target.CurrentLife, wing.CurrentLife, startWingLife)
	}
}

func TestRoyalConflictRedMoonPendantExtendsNextRedMoon(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	pendant := NewCardInstance(baseCard(t, "2621105"), 0, 1)
	p0.Equipment[0] = pendant
	pendant.SlotIndex = 0
	redMoon := readySkill(baseCard(t, "3611101"), 0)
	p0.Skills[0] = redMoon

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  pendant.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use red moon pendant: %v", err)
	}
	if p0.Equipment[0] != nil || len(p0.Graveyard) != 1 || p0.Graveyard[0] != pendant {
		t.Fatalf("pendant should be sacrificed to graveyard, equipment=%v grave=%v", p0.Equipment[0], cardsToInfo(p0.Graveyard))
	}
	if p0.NextRedMoonDuration != 1 {
		t.Fatalf("pendant should arm next red moon duration +1, got %d", p0.NextRedMoonDuration)
	}

	p0.Elements[model.ElementShadow] = 1
	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": redMoon.InstanceID,
		"target_type": "none",
	}}); err != nil {
		t.Fatalf("cast red moon after pendant: %v", err)
	}
	if got := redMoon.Statuses[StatusAbilityDuration]; got != 2 {
		t.Fatalf("pendant should extend next red moon duration to 2, got %d statuses=%v", got, redMoon.Statuses)
	}
	if p0.NextRedMoonDuration != 0 {
		t.Fatalf("next red moon duration should be consumed, got %d", p0.NextRedMoonDuration)
	}
}

func TestRoyalConflictRedMoonProphetReducesCurrentOrNextCooldown(t *testing.T) {
	t.Run("next red moon", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		redMoon := readySkill(baseCard(t, "3611101"), 0)
		p0.Skills[0] = redMoon
		prophet := placeUnit(baseCard(t, "1621111"), 0, 0, 0, engine)

		engine.triggerEffects(TriggerOnEnter, prophet, nil, nil)
		if p0.NextRedMoonCooldown != 1 {
			t.Fatalf("prophet should arm next red moon cooldown -1 while red moon is inactive, got %d", p0.NextRedMoonCooldown)
		}

		p0.Elements[model.ElementShadow] = 1
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": redMoon.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast red moon after prophet: %v", err)
		}
		if got := redMoon.Statuses[StatusCooldown]; got != 1 {
			t.Fatalf("prophet should reduce next red moon cooldown from 2 to 1, got %d statuses=%v", got, redMoon.Statuses)
		}
		if p0.NextRedMoonCooldown != 0 {
			t.Fatalf("next red moon cooldown should be consumed, got %d", p0.NextRedMoonCooldown)
		}
	})

	t.Run("next cooldown reduction applies after cooldown additions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		redMoon := readySkill(baseCard(t, "3611101"), 0)
		p0.Skills[0] = redMoon
		p0.NextRedMoonCooldown = 3
		p0.TempModifiers = append(p0.TempModifiers, TemporaryModifier{
			Type:        TempModSkillUseCooldownAdd,
			Amount:      2,
			ExpiresTurn: engine.State.TurnNumber + 1,
		})

		p0.Elements[model.ElementShadow] = 1
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": redMoon.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast red moon with stacked cooldown modifiers: %v", err)
		}
		if got := redMoon.Statuses[StatusCooldown]; got != 1 {
			t.Fatalf("next red moon cooldown reduction should apply after additions, got %d statuses=%v", got, redMoon.Statuses)
		}
		if p0.NextRedMoonCooldown != 0 {
			t.Fatalf("next red moon cooldown should be consumed, got %d", p0.NextRedMoonCooldown)
		}
	})

	t.Run("current red moon", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		redMoon := readySkill(baseCard(t, "3611101"), 0)
		redMoon.Statuses[StatusAbilityDuration] = 1
		redMoon.Statuses[StatusCooldown] = 2
		p0.Skills[0] = redMoon
		prophet := placeUnit(baseCard(t, "1621111"), 0, 0, 0, engine)

		engine.triggerEffects(TriggerOnEnter, prophet, nil, nil)
		if got := redMoon.Statuses[StatusCooldown]; got != 1 {
			t.Fatalf("prophet enter should reduce current red moon cooldown to 1, got %d statuses=%v", got, redMoon.Statuses)
		}
		if p0.NextRedMoonCooldown != 0 {
			t.Fatalf("current red moon reduction should not arm next cooldown, got %d", p0.NextRedMoonCooldown)
		}

		engine.triggerEffects(TriggerOnDeath, prophet, nil, nil)
		if got := redMoon.Statuses[StatusCooldown]; got != 0 {
			t.Fatalf("prophet death should remove final cooldown layer, got %d statuses=%v", got, redMoon.Statuses)
		}
	})
}

func TestRoyalConflictSimpleCardEffects(t *testing.T) {
	t.Run("enter draw and resource effects", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		drawOne := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		drawTwoA := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		drawTwoB := NewCardInstance(baseCard(t, "1021003"), 0, 1)
		p0.Deck = []*CardInstance{drawOne, drawTwoA, drawTwoB}

		geomancer := placeUnit(baseCard(t, "1421115"), 0, 0, 0, engine)
		engine.triggerEffects(TriggerOnEnter, geomancer, nil, nil)
		if len(p0.Hand) != 1 || p0.Hand[0] != drawOne {
			t.Fatalf("geomancer should draw one card, hand=%v", cardsToInfo(p0.Hand))
		}

		p0.Hand = nil
		hummingbird := placeUnit(baseCard(t, "1321108"), 0, 1, 0, engine)
		engine.triggerEffects(TriggerOnEnter, hummingbird, nil, nil)
		if len(p0.Hand) != 2 || p0.Hand[0] != drawTwoA || p0.Hand[1] != drawTwoB {
			t.Fatalf("hummingbird should draw two with fewer than two cards in hand, hand=%v", cardsToInfo(p0.Hand))
		}

		p0.Hand = []*CardInstance{
			NewCardInstance(baseCard(t, "1021004"), 0, 1),
			NewCardInstance(baseCard(t, "1021005"), 0, 1),
		}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021006"), 0, 1)}
		engine.triggerEffects(TriggerOnEnter, hummingbird, nil, nil)
		if len(p0.Hand) != 2 {
			t.Fatalf("hummingbird should not draw when hand has two cards, hand=%v", cardsToInfo(p0.Hand))
		}

		geomancer.CurrentLife = geomancer.Card.Life - 1
		ally := placeUnit(baseCard(t, "1021007"), 0, 2, 0, engine)
		ally.CurrentLife = ally.Card.Life - 1
		healthy := placeUnit(baseCard(t, "1021008"), 0, 2, 1, engine)
		prayer := placeUnit(baseCard(t, "1521114"), 0, 0, 1, engine)
		before := p0.Elements[model.ElementLight]
		engine.triggerEffects(TriggerOnEnter, prayer, nil, nil)
		if got := p0.Elements[model.ElementLight] - before; got != 2 {
			t.Fatalf("prayer should gain light for two wounded friendly units, got %d with healthy=%v", got, cardToInfo(healthy))
		}
	})

	t.Run("use item and equipment active effects", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		reshape := NewCardInstance(baseCard(t, "2021107"), 0, 1)
		discardA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		discardB := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		drawA := NewCardInstance(baseCard(t, "1021003"), 0, 1)
		drawB := NewCardInstance(baseCard(t, "1021004"), 0, 1)
		p0.Hand = []*CardInstance{reshape, discardA, discardB}
		p0.Deck = []*CardInstance{drawA, drawB}
		p0.RevealedHand[discardA.InstanceID] = true

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": reshape.InstanceID,
		}}); err != nil {
			t.Fatalf("use reshape: %v", err)
		}
		if len(p0.Hand) != 2 || p0.Hand[0] != drawA || p0.Hand[1] != drawB {
			t.Fatalf("reshape should discard hand then draw two, hand=%v", cardsToInfo(p0.Hand))
		}
		if len(p0.Graveyard) != 3 || p0.Graveyard[0] != reshape || p0.Graveyard[1] != discardA || p0.Graveyard[2] != discardB {
			t.Fatalf("reshape should place itself and discarded cards in graveyard, grave=%v", cardsToInfo(p0.Graveyard))
		}
		if p0.RevealedHand[discardA.InstanceID] {
			t.Fatalf("reshape should clear revealed flags for discarded cards")
		}

		p0.Hand = []*CardInstance{NewCardInstance(baseCard(t, "2521106"), 0, 1)}
		p0.Elements[model.ElementLight] = 1
		p0.Deck = nil
		p0.Graveyard = nil
		woundedA := placeUnit(baseCard(t, "1021004"), 0, 0, 1, engine)
		woundedA.CurrentLife = woundedA.Card.Life - 2
		ally := placeUnit(baseCard(t, "1021005"), 0, 0, 2, engine)
		ally.CurrentLife = ally.Card.Life - 2
		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": p0.Hand[0].InstanceID,
		}}); err != nil {
			t.Fatalf("use moonlight scroll: %v", err)
		}
		if woundedA.CurrentLife != woundedA.Card.Life || ally.CurrentLife != ally.Card.Life {
			t.Fatalf("moonlight scroll should heal all friendly units by 2, woundedA=%d ally=%d", woundedA.CurrentLife, ally.CurrentLife)
		}

		dragonbone := NewCardInstance(baseCard(t, "2521104"), 0, 1)
		p0.Equipment[0] = dragonbone
		dragonbone.SlotIndex = 0
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021006"), 0, 1),
			NewCardInstance(baseCard(t, "1021007"), 0, 1),
		}
		p0.Hand = nil
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  dragonbone.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use golden dragonbone: %v", err)
		}
		if p0.Equipment[0] != nil || len(p0.Hand) != 2 {
			t.Fatalf("golden dragonbone should sacrifice itself and draw two, equipment=%v hand=%v", p0.Equipment[0], cardsToInfo(p0.Hand))
		}
		if len(p0.Graveyard) == 0 || p0.Graveyard[len(p0.Graveyard)-1] != dragonbone {
			t.Fatalf("golden dragonbone should move to graveyard, grave=%v", cardsToInfo(p0.Graveyard))
		}
	})
}

func TestRoyalConflictSimpleConsumableChoiceEffects(t *testing.T) {
	t.Run("lost silverleaf draws then discards one selected hand card", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		silverleaf := NewCardInstance(baseCard(t, "2021101"), 0, 1)
		kept := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		drawA := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		drawB := NewCardInstance(baseCard(t, "1021003"), 0, 1)
		p0.Hand = []*CardInstance{silverleaf, kept}
		p0.Deck = []*CardInstance{drawA, drawB}
		p0.Elements[model.ElementArcane] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": silverleaf.InstanceID,
		}}); err != nil {
			t.Fatalf("use lost silverleaf: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "lost_silverleaf_discard" {
			t.Fatalf("lost silverleaf should prompt for discard after drawing, pending=%+v", engine.State.PendingAction)
		}
		if len(p0.Hand) != 3 || p0.Hand[0] != kept || p0.Hand[1] != drawA || p0.Hand[2] != drawB {
			t.Fatalf("lost silverleaf should draw two before discard, hand=%v", cardsToInfo(p0.Hand))
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{drawA.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve silverleaf discard: %v", err)
		}
		if len(p0.Hand) != 2 || p0.Hand[0] != kept || p0.Hand[1] != drawB {
			t.Fatalf("lost silverleaf should discard the selected card, hand=%v", cardsToInfo(p0.Hand))
		}
		if len(p0.Graveyard) != 2 || p0.Graveyard[0] != silverleaf || p0.Graveyard[1] != drawA {
			t.Fatalf("lost silverleaf graveyard order wrong, grave=%v", cardsToInfo(p0.Graveyard))
		}
	})

	t.Run("blessed lone star buffs a selected friendly companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		item := NewCardInstance(baseCard(t, "2521101"), 0, 1)
		target := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
		p0.Hand = []*CardInstance{item}
		p0.Elements[model.ElementLight] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": item.InstanceID,
		}}); err != nil {
			t.Fatalf("use blessed lone star: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blessed_lone_star_target" {
			t.Fatalf("blessed lone star should prompt for target, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve blessed lone star: %v", err)
		}
		if target.CurrentLife != target.Card.Life+1 || effectiveElementsGain(target)[model.ElementLight] != target.Card.ElementsGain[model.ElementLight]+1 {
			t.Fatalf("blessed lone star should add +1 life and +1 light load, life=%d load=%v", target.CurrentLife, effectiveElementsGain(target))
		}
	})

	t.Run("arcane bomb damages a companion in spell range", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		bomb := NewCardInstance(baseCard(t, "2021116"), 0, 1)
		front := placeUnit(baseCard(t, "1021005"), 1, 1, 0, engine)
		back := placeUnit(baseCard(t, "1021006"), 1, 1, 2, engine)
		p0.Hand = []*CardInstance{bomb}
		p0.Elements[model.ElementArcane] = 3

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": bomb.InstanceID,
		}}); err != nil {
			t.Fatalf("use arcane bomb: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "arcane_bomb_target" {
			t.Fatalf("arcane bomb should prompt for companion target, pending=%+v", engine.State.PendingAction)
		}
		for _, candidate := range engine.State.PendingAction.Candidates {
			if candidate["instance_id"] == back.InstanceID {
				t.Fatalf("arcane bomb should not offer enemy back row behind a front unit, candidates=%+v", engine.State.PendingAction.Candidates)
			}
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{front.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve arcane bomb: %v", err)
		}
		if front.CurrentLife != front.Card.Life-2 {
			t.Fatalf("arcane bomb should deal 2 damage to selected companion, life=%d", front.CurrentLife)
		}
	})
}

func TestRoyalConflictSimpleGeneratedAndPrayerEffects(t *testing.T) {
	t.Run("dream consumables draw cards or gain arcane", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		bloom := NewCardInstance(baseCard(t, "2201101"), 0, 1)
		drawA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		drawB := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		drawC := NewCardInstance(baseCard(t, "1021003"), 0, 1)
		p0.Hand = []*CardInstance{bloom}
		p0.Deck = []*CardInstance{drawA, drawB, drawC}

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": bloom.InstanceID,
		}}); err != nil {
			t.Fatalf("use dream bloom: %v", err)
		}
		if len(p0.Hand) != 3 || p0.Hand[0] != drawA || p0.Hand[1] != drawB || p0.Hand[2] != drawC {
			t.Fatalf("dream bloom should draw three cards, hand=%v", cardsToInfo(p0.Hand))
		}

		mana := NewCardInstance(baseCard(t, "2201102"), 0, 1)
		p0.Hand = []*CardInstance{mana}
		before := p0.Elements[model.ElementArcane]
		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": mana.InstanceID,
		}}); err != nil {
			t.Fatalf("use dream mana: %v", err)
		}
		if got := p0.Elements[model.ElementArcane] - before; got != 3 {
			t.Fatalf("dream mana should gain 3 arcane, got %d elements=%v", got, p0.Elements)
		}
	})

	t.Run("blood puppet damages own hero on enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		hero := placeUnit(baseCard(t, "4011001"), 0, 1, 1, engine)
		p0.Hero = hero
		puppet := placeUnit(baseCard(t, "1621103"), 0, 0, 0, engine)

		engine.triggerEffects(TriggerOnEnter, puppet, nil, nil)
		if hero.CurrentLife != hero.Card.Life-2 {
			t.Fatalf("blood puppet should deal 2 damage to own hero, life=%d", hero.CurrentLife)
		}
	})

	t.Run("prayer load effects", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		lotus := placeUnit(baseCard(t, "1221106"), 0, 0, 0, engine)
		root := placeUnit(baseCard(t, "1421105"), 0, 1, 0, engine)

		engine.triggerPrayerAbilities(0)
		if got := effectiveElementsGain(lotus)[model.ElementWater]; got != lotus.Card.ElementsGain[model.ElementWater]+1 {
			t.Fatalf("mirror lotus prayer should gain +1 water load, got %d load=%v", got, effectiveElementsGain(lotus))
		}
		if got := effectiveElementsGain(root)[model.ElementEarth]; got != 1 {
			t.Fatalf("inactive root prayer should gain 1 earth load while loadless, got %d load=%v", got, effectiveElementsGain(root))
		}

		engine.triggerPrayerAbilities(0)
		if got := effectiveElementsGain(root)[model.ElementEarth]; got != 1 {
			t.Fatalf("inactive root prayer should not add more load once it has load, got %d load=%v", got, effectiveElementsGain(root))
		}
	})
}

func TestRoyalConflictSimpleTargetedEnterAndDeathEffects(t *testing.T) {
	t.Run("swordsmanship teacher buffs adjacent friendly companion attack", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		teacher := placeUnit(baseCard(t, "1021102"), 0, 1, 1, engine)
		adjacent := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		far := placeUnit(baseCard(t, "1021002"), 0, 2, 2, engine)

		engine.triggerEffects(TriggerOnEnter, teacher, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "swordsmanship_teacher_buff" {
			t.Fatalf("swordsmanship teacher should prompt for adjacent target, pending=%+v", engine.State.PendingAction)
		}
		for _, candidate := range engine.State.PendingAction.Candidates {
			if candidate["instance_id"] == far.InstanceID {
				t.Fatalf("swordsmanship teacher should not offer non-adjacent target, candidates=%+v", engine.State.PendingAction.Candidates)
			}
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{adjacent.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve swordsmanship teacher: %v", err)
		}
		if adjacent.AttackBonus != 1 || far.AttackBonus != 0 {
			t.Fatalf("swordsmanship teacher should buff selected adjacent companion only, adjacent=%d far=%d", adjacent.AttackBonus, far.AttackBonus)
		}
	})

	t.Run("lone star guardian spirit enter and death prompts", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		guardian := placeUnit(baseCard(t, "1521103"), 0, 0, 0, engine)
		target := placeUnit(baseCard(t, "1021003"), 0, 1, 0, engine)

		engine.triggerEffects(TriggerOnEnter, guardian, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "lone_star_guardian_life" {
			t.Fatalf("guardian enter should prompt for life target, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve guardian enter: %v", err)
		}
		if target.CurrentLife != target.Card.Life+1 {
			t.Fatalf("guardian enter should give +1 life, life=%d", target.CurrentLife)
		}

		engine.triggerEffects(TriggerOnDeath, guardian, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "lone_star_guardian_load" {
			t.Fatalf("guardian death should prompt for load target, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve guardian death: %v", err)
		}
		if effectiveElementsGain(target)[model.ElementLight] != target.Card.ElementsGain[model.ElementLight]+1 {
			t.Fatalf("guardian death should add +1 light load, load=%v", effectiveElementsGain(target))
		}
	})

	t.Run("whisper elf deathrattles target enemy and friendly companions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		hunter := placeUnit(baseCard(t, "1621112"), 0, 0, 0, engine)
		enemy := placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		priest := placeUnit(baseCard(t, "1621113"), 0, 1, 0, engine)
		ally := placeUnit(baseCard(t, "1021005"), 0, 2, 0, engine)

		engine.triggerEffects(TriggerOnDeath, hunter, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "whisper_elf_hunter_damage" {
			t.Fatalf("hunter death should prompt for enemy target, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{enemy.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve hunter death: %v", err)
		}
		if enemy.CurrentLife != enemy.Card.Life-1 {
			t.Fatalf("hunter death should damage selected enemy, life=%d", enemy.CurrentLife)
		}

		engine.triggerEffects(TriggerOnDeath, priest, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "whisper_elf_priest_load" {
			t.Fatalf("priest death should prompt for friendly companion, pending=%+v", engine.State.PendingAction)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{ally.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve priest death: %v", err)
		}
		if effectiveElementsGain(ally)[model.ElementShadow] != ally.Card.ElementsGain[model.ElementShadow]+1 {
			t.Fatalf("priest death should add +1 shadow load, load=%v", effectiveElementsGain(ally))
		}
	})
}

func TestRoyalConflictSkyCityTycoonConsumesForOrderedDraw(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	tycoon := placeUnit(baseCard(t, "1021106"), 0, 0, 0, engine)
	selfDraw := NewCardInstance(baseCard(t, "1021001"), 0, 1)
	opponentDraw := NewCardInstance(baseCard(t, "1021002"), 1, 1)
	p0.Deck = []*CardInstance{selfDraw}
	p1.Deck = []*CardInstance{opponentDraw}
	beforeArcane := p0.Elements[model.ElementArcane]
	beforeAir := p0.Elements[model.ElementAir]

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  tycoon.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use sky city tycoon: %v", err)
	}
	if !tycoon.IsHorizontal {
		t.Fatalf("sky city tycoon should be horizontal after paying consume cost")
	}
	if p0.Elements[model.ElementArcane] != beforeArcane || p0.Elements[model.ElementAir] != beforeAir {
		t.Fatalf("sky city tycoon active consume should not grant its printed load, elements=%v", p0.Elements)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "sky_city_tycoon_draw_order" {
		t.Fatalf("sky city tycoon should ask for draw order, pending=%+v", engine.State.PendingAction)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{"opponent_first"},
	}}); err != nil {
		t.Fatalf("resolve sky city tycoon draw order: %v", err)
	}
	if len(p1.Hand) != 1 || p1.Hand[0] != opponentDraw || len(p0.Hand) != 1 || p0.Hand[0] != selfDraw {
		t.Fatalf("sky city tycoon should draw one for each player in chosen order, p0=%v p1=%v", cardsToInfo(p0.Hand), cardsToInfo(p1.Hand))
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  tycoon.InstanceID,
		"ability_type": "per_turn",
	}}); err == nil {
		t.Fatalf("sky city tycoon should not be usable while already horizontal")
	}
}

func TestRoyalConflictBloodNourishExilesShadowGraveyardCard(t *testing.T) {
	t.Run("exiles selected shadow card and gains shadow elements", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		nourish := readySkill(baseCard(t, "3621110"), 0)
		shadow := NewCardInstance(baseCard(t, "1621112"), 0, 1)
		other := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Skills[0] = nourish
		p0.Graveyard = []*CardInstance{other, shadow}
		p0.Elements[model.ElementShadow] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": nourish.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast blood nourish: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_nourish_exile" {
			t.Fatalf("blood nourish should prompt for shadow graveyard card, pending=%+v", engine.State.PendingAction)
		}
		if len(engine.State.PendingAction.Candidates) != 1 || engine.State.PendingAction.Candidates[0]["instance_id"] != shadow.InstanceID {
			t.Fatalf("blood nourish should only offer shadow graveyard cards, candidates=%+v", engine.State.PendingAction.Candidates)
		}
		afterPay := p0.Elements[model.ElementShadow]
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{shadow.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve blood nourish: %v", err)
		}
		if len(p0.Exile) != 1 || p0.Exile[0] != shadow {
			t.Fatalf("blood nourish should exile selected shadow card, exile=%v grave=%v", cardsToInfo(p0.Exile), cardsToInfo(p0.Graveyard))
		}
		if len(p0.Graveyard) != 1 || p0.Graveyard[0] != other {
			t.Fatalf("blood nourish should leave non-selected graveyard cards, grave=%v", cardsToInfo(p0.Graveyard))
		}
		if p0.Elements[model.ElementShadow] != afterPay+2 {
			t.Fatalf("blood nourish should gain 2 shadow after selection, before=%d elements=%v", afterPay, p0.Elements)
		}
	})

	t.Run("does nothing without shadow graveyard card", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		nourish := readySkill(baseCard(t, "3621110"), 0)
		p0.Skills[0] = nourish
		p0.Graveyard = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
		p0.Elements[model.ElementShadow] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": nourish.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast blood nourish without target: %v", err)
		}
		if engine.State.PendingAction != nil {
			t.Fatalf("blood nourish should not prompt without shadow graveyard card, pending=%+v", engine.State.PendingAction)
		}
		if len(p0.Exile) != 0 {
			t.Fatalf("blood nourish should not exile anything without shadow target, exile=%v", cardsToInfo(p0.Exile))
		}
	})
}

func TestRoyalConflictSimpleEnterBatchTwo(t *testing.T) {
	t.Run("dimensional rift beast exiles only selected enemy companion in spell range", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		beast := placeUnit(baseCard(t, "1021104"), 0, 1, 0, engine)
		inRange := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		outOfRange := placeUnit(baseCard(t, "1021002"), 1, 1, 2, engine)

		engine.triggerEffects(TriggerOnEnter, beast, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "dimensional_rift_beast_exile" {
			t.Fatalf("1021104 should ask for an enemy companion in spell range, pending=%+v", engine.State.PendingAction)
		}
		for _, candidate := range engine.State.PendingAction.Candidates {
			if candidate["instance_id"] == outOfRange.InstanceID {
				t.Fatalf("1021104 should not offer out-of-range enemy companions, candidates=%+v", engine.State.PendingAction.Candidates)
			}
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{outOfRange.InstanceID},
		}}); err == nil {
			t.Fatal("1021104 should reject forged out-of-range selections")
		}
		resolvePendingSelection(t, engine, 0, inRange.InstanceID)
		if len(engine.State.Players[1].Exile) != 1 || engine.State.Players[1].Exile[0] != inRange {
			t.Fatalf("1021104 should exile selected enemy companion, exile=%v", cardsToInfo(engine.State.Players[1].Exile))
		}
		if engine.State.Players[1].Units[outOfRange.Position.Col][outOfRange.Position.Row] != outOfRange {
			t.Fatal("1021104 should leave unselected enemy companions on the battlefield")
		}
	})

	t.Run("beacon guard gains shield only when outnumbered", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		guard := placeUnit(baseCard(t, "1121103"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
		engine.triggerEffects(TriggerOnEnter, guard, nil, nil)
		if p0.Shield != 3 {
			t.Fatalf("1121103 should gain shield 3 when friendly units are fewer, got %d", p0.Shield)
		}

		evenEngine := setupReportedBugEngine(t)
		evenP0 := evenEngine.State.Players[0]
		evenGuard := placeUnit(baseCard(t, "1121103"), 0, 0, 0, evenEngine)
		placeUnit(baseCard(t, "1021001"), 1, 0, 0, evenEngine)
		evenEngine.triggerEffects(TriggerOnEnter, evenGuard, nil, nil)
		if evenP0.Shield != 0 {
			t.Fatalf("1121103 should not gain shield when not outnumbered, got %d", evenP0.Shield)
		}
	})

	t.Run("silverleaf messenger searches lost silverleaf", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		messenger := placeUnit(baseCard(t, "1321110"), 0, 0, 0, engine)
		flower := NewCardInstance(baseCard(t, "2021101"), 0, 1)
		other := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Deck = []*CardInstance{other, flower}

		engine.triggerEffects(TriggerOnEnter, messenger, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "silverleaf_messenger_search" {
			t.Fatalf("1321110 should prompt to search lost silverleaf, pending=%+v", engine.State.PendingAction)
		}
		if len(engine.State.PendingAction.Candidates) != 1 || engine.State.PendingAction.Candidates[0]["instance_id"] != flower.InstanceID {
			t.Fatalf("1321110 should only offer lost silverleaf cards, candidates=%+v", engine.State.PendingAction.Candidates)
		}
		resolvePendingSelection(t, engine, 0, flower.InstanceID)
		if len(p0.Hand) != 1 || p0.Hand[0] != flower {
			t.Fatalf("1321110 should move selected lost silverleaf to hand, hand=%v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("council messenger gives opponent a jiuxiao mark", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		messenger := placeUnit(baseCard(t, "1321113"), 0, 0, 0, engine)
		engine.triggerEffects(TriggerOnEnter, messenger, nil, nil)
		p1 := engine.State.Players[1]
		if len(p1.Hand) != 1 || p1.Hand[0].Card.Number != "2001102" || p1.Hand[0].OwnerID != 1 {
			t.Fatalf("1321113 should add a Jiuxiao Mark to opponent hand, hand=%v", cardsToInfo(p1.Hand))
		}
	})

	t.Run("church exorcist purifies one friendly card and gains light per layer", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		exorcist := placeUnit(baseCard(t, "1521106"), 0, 0, 0, engine)
		target := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		target.Statuses[StatusBurn] = 2
		target.Statuses[StatusFreeze] = 1
		target.Statuses[StatusCooldown] = 4
		skill := readySkill(baseCard(t, "3021001"), 0)
		skill.Statuses[StatusWeaken] = 1
		p0.Skills[0] = skill

		engine.triggerEffects(TriggerOnEnter, exorcist, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "church_exorcist_purify" {
			t.Fatalf("1521106 should prompt for a friendly card with negative statuses, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if target.Statuses[StatusBurn] != 0 || target.Statuses[StatusFreeze] != 0 {
			t.Fatalf("1521106 should clear selected negative statuses, statuses=%v", target.Statuses)
		}
		if target.Statuses[StatusCooldown] != 4 {
			t.Fatalf("1521106 should not clear non-negative statuses, statuses=%v", target.Statuses)
		}
		if p0.Elements[model.ElementLight] != 3 {
			t.Fatalf("1521106 should gain 1 light per removed layer, elements=%v", p0.Elements)
		}
		if skill.Statuses[StatusWeaken] != 1 {
			t.Fatalf("1521106 should only purify the selected card, skill statuses=%v", skill.Statuses)
		}
	})

	t.Run("church exorcist can purify bound skills", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		exorcist := placeUnit(baseCard(t, "1521106"), 0, 0, 0, engine)
		host := placeUnit(baseCard(t, "1011103"), 0, 1, 0, engine)
		bound := NewCardInstance(baseCard(t, "3001101"), 0, 1)
		bound.SlotIndex = -1
		bound.Statuses[StatusWeaken] = 2
		host.BoundSkills = []*CardInstance{bound}

		engine.triggerEffects(TriggerOnEnter, exorcist, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "church_exorcist_purify" {
			t.Fatalf("1521106 should prompt for a bound skill with negative statuses, pending=%+v", engine.State.PendingAction)
		}
		if len(engine.State.PendingAction.Candidates) != 1 || engine.State.PendingAction.Candidates[0]["zone"] != "bound_skill" {
			t.Fatalf("1521106 should expose bound skill candidates, candidates=%+v", engine.State.PendingAction.Candidates)
		}
		resolvePendingSelection(t, engine, 0, bound.InstanceID)
		if bound.Statuses[StatusWeaken] != 0 {
			t.Fatalf("1521106 should clear selected bound skill negative statuses, statuses=%v", bound.Statuses)
		}
		if p0.Elements[model.ElementLight] != 2 {
			t.Fatalf("1521106 should gain light for bound skill negative layers, elements=%v", p0.Elements)
		}
	})
}

func TestRoyalConflictJiuxiaoMarkEffects(t *testing.T) {
	t.Run("jiuxiao assassin adds marks to hand and deck", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		assassin := placeUnit(baseCard(t, "1021115"), 0, 0, 0, engine)
		p1 := engine.State.Players[1]

		engine.triggerEffects(TriggerOnEnter, assassin, nil, nil)
		if len(p1.Hand) != 1 || p1.Hand[0].Card.Number != "2001102" || p1.Hand[0].OwnerID != 1 {
			t.Fatalf("1021115 enter should add a Jiuxiao Mark to opponent hand, hand=%v", cardsToInfo(p1.Hand))
		}
		engine.triggerEffects(TriggerOnDeath, assassin, nil, nil)
		if countCardsByNumber(p1.Deck, "2001102") != 4 {
			t.Fatalf("1021115 death should shuffle four Jiuxiao Marks into opponent deck, deck=%v", cardsToInfo(p1.Deck))
		}
	})

	t.Run("jiuxiao contact prayer respects opponent hand limit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p1 := engine.State.Players[1]
		contact := placeUnit(baseCard(t, "1321112"), 0, 0, 0, engine)
		if err := globalRegistry.GetBehavior("1321112").(PerTurnAbility).OnPerTurn(&EffectContext{Engine: engine, Source: contact, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1321112 prayer under hand limit: %v", err)
		}
		if len(p1.Hand) != 1 || p1.Hand[0].Card.Number != "2001102" {
			t.Fatalf("1321112 should add a Jiuxiao Mark while opponent hand is below limit, hand=%v", cardsToInfo(p1.Hand))
		}

		limitEngine := setupReportedBugEngine(t)
		limitP1 := limitEngine.State.Players[1]
		for i := 0; i < limitEngine.handLimitForPlayer(limitP1); i++ {
			limitP1.Hand = append(limitP1.Hand, NewCardInstance(baseCard(t, "1021001"), 1, 1))
		}
		limitContact := placeUnit(baseCard(t, "1321112"), 0, 0, 0, limitEngine)
		if err := globalRegistry.GetBehavior("1321112").(PerTurnAbility).OnPerTurn(&EffectContext{Engine: limitEngine, Source: limitContact, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1321112 prayer at hand limit: %v", err)
		}
		if countCardsByNumber(limitP1.Hand, "2001102") != 0 {
			t.Fatalf("1321112 should not add a mark when opponent hand reached limit, hand=%v", cardsToInfo(limitP1.Hand))
		}
	})

	t.Run("pigeon arrest order adds a Jiuxiao Mark after a friendly spell hits once per turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		order := NewCardInstance(baseCard(t, "2321107"), 0, 1)
		p0.Equipment[0] = order
		behavior := Card2321107PigeonArrestOrder{}

		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: order, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("pigeon arrest order friendly hit: %v", err)
		}
		if len(p1.Hand) != 1 || p1.Hand[0].Card.Number != "2001102" || p1.Hand[0].OwnerID != 1 || order.UsedThisTurn != 1 {
			t.Fatalf("2321107 should add one Jiuxiao Mark to opponent hand and spend trigger, hand=%v used=%d", cardsToInfo(p1.Hand), order.UsedThisTurn)
		}
		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: order, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("pigeon arrest order second friendly hit: %v", err)
		}
		if len(p1.Hand) != 1 || order.UsedThisTurn != 1 {
			t.Fatalf("2321107 should trigger at most once per turn, hand=%v used=%d", cardsToInfo(p1.Hand), order.UsedThisTurn)
		}

		enemyEngine := setupReportedBugEngine(t)
		enemyP0 := enemyEngine.State.Players[0]
		enemyP1 := enemyEngine.State.Players[1]
		enemyOrder := NewCardInstance(baseCard(t, "2321107"), 0, 1)
		enemyP0.Equipment[0] = enemyOrder
		if err := behavior.OnSpellHit(&EffectContext{Engine: enemyEngine, Source: enemyOrder, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 1}}); err != nil {
			t.Fatalf("pigeon arrest order enemy hit: %v", err)
		}
		if len(enemyP1.Hand) != 0 || enemyOrder.UsedThisTurn != 0 {
			t.Fatalf("2321107 should ignore enemy spell hits, hand=%v used=%d", cardsToInfo(enemyP1.Hand), enemyOrder.UsedThisTurn)
		}
	})

	t.Run("raider gunner discards an enemy hand card after a friendly spell hits once per game", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		gunner := placeUnit(baseCard(t, "1221111"), 0, 0, 0, engine)
		first := NewCardInstance(baseCard(t, "1021001"), 1, 1)
		second := NewCardInstance(baseCard(t, "1021002"), 1, 1)
		p1.Hand = []*CardInstance{first, second}
		behavior := Card1221111RaiderGunner{}

		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: gunner, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("raider gunner friendly hit: %v", err)
		}
		if len(p1.Hand) != 1 || len(p1.Graveyard) != 1 || !gunner.UltimateUsed {
			t.Fatalf("1221111 should discard one enemy hand card and spend ultimate, hand=%v grave=%v used=%v", cardsToInfo(p1.Hand), cardsToInfo(p1.Graveyard), gunner.UltimateUsed)
		}
		if containsCardInstance(p1.Hand, p1.Graveyard[0]) {
			t.Fatalf("discarded card should leave enemy hand, hand=%v grave=%v", cardsToInfo(p1.Hand), cardsToInfo(p1.Graveyard))
		}

		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: gunner, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("raider gunner second friendly hit: %v", err)
		}
		if len(p1.Hand) != 1 || len(p1.Graveyard) != 1 {
			t.Fatalf("1221111 should trigger at most once per game, hand=%v grave=%v", cardsToInfo(p1.Hand), cardsToInfo(p1.Graveyard))
		}

		enemyEngine := setupReportedBugEngine(t)
		enemyP1 := enemyEngine.State.Players[1]
		enemyGunner := placeUnit(baseCard(t, "1221111"), 0, 0, 0, enemyEngine)
		enemyP1.Hand = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 1, 1)}
		if err := behavior.OnSpellHit(&EffectContext{Engine: enemyEngine, Source: enemyGunner, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 1}}); err != nil {
			t.Fatalf("raider gunner enemy hit: %v", err)
		}
		if len(enemyP1.Hand) != 1 || len(enemyP1.Graveyard) != 0 || enemyGunner.UsedThisTurn != 0 {
			t.Fatalf("1221111 should ignore enemy spell hits, hand=%v grave=%v used=%d", cardsToInfo(enemyP1.Hand), cardsToInfo(enemyP1.Graveyard), enemyGunner.UsedThisTurn)
		}

		emptyEngine := setupReportedBugEngine(t)
		emptyGunner := placeUnit(baseCard(t, "1221111"), 0, 0, 0, emptyEngine)
		if err := behavior.OnSpellHit(&EffectContext{Engine: emptyEngine, Source: emptyGunner, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("raider gunner empty hand hit: %v", err)
		}
		if emptyGunner.UsedThisTurn != 0 {
			t.Fatalf("1221111 should not spend trigger when opponent has no hand cards, used=%d", emptyGunner.UsedThisTurn)
		}

		if len(p0.Hand) != 0 {
			t.Fatalf("1221111 should not touch caster hand, hand=%v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("council executor discards an extra card when it hits a mark", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p1 := engine.State.Players[1]
		p1.Hero = NewCardInstance(baseCard(t, "4311003"), 1, 1)
		p1.Hero.Position = &Position{Col: 1, Row: 1}
		p1.Units[1][1] = p1.Hero
		p1.Hand = []*CardInstance{
			NewCardInstance(baseCard(t, "2001102"), 1, 1),
			NewCardInstance(baseCard(t, "2001102"), 1, 1),
		}
		beforeLife := p1.Hero.CurrentLife
		executor := placeUnit(baseCard(t, "1321114"), 0, 0, 0, engine)
		engine.triggerEffects(TriggerOnEnter, executor, nil, nil)
		if len(p1.Hand) != 0 || countCardsByNumber(p1.Graveyard, "2001102") != 2 {
			t.Fatalf("1321114 should discard a second card after hitting a mark, hand=%v grave=%v", cardsToInfo(p1.Hand), cardsToInfo(p1.Graveyard))
		}
		if p1.Hero.CurrentLife != beforeLife-4 {
			t.Fatalf("discarded Jiuxiao Marks should deal 2 damage each to their hero, before=%d life=%d", beforeLife, p1.Hero.CurrentLife)
		}

		normalEngine := setupReportedBugEngine(t)
		normalP1 := normalEngine.State.Players[1]
		normalP1.Hero = NewCardInstance(baseCard(t, "4311003"), 1, 1)
		normalP1.Hero.Position = &Position{Col: 1, Row: 1}
		normalP1.Units[1][1] = normalP1.Hero
		normalP1.Hand = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 1, 1)}
		normalBeforeLife := normalP1.Hero.CurrentLife
		normalExecutor := placeUnit(baseCard(t, "1321114"), 0, 0, 0, normalEngine)
		normalEngine.triggerEffects(TriggerOnEnter, normalExecutor, nil, nil)
		if len(normalP1.Hand) != 0 || len(normalP1.Graveyard) != 1 || normalP1.Graveyard[0].Card.Number == "2001102" {
			t.Fatalf("1321114 should only discard once when the first discard is not a mark, hand=%v grave=%v", cardsToInfo(normalP1.Hand), cardsToInfo(normalP1.Graveyard))
		}
		if normalP1.Hero.CurrentLife != normalBeforeLife {
			t.Fatalf("non-mark discard should not damage hero, before=%d life=%d", normalBeforeLife, normalP1.Hero.CurrentLife)
		}
	})

	t.Run("discarding a mark to hand limit damages its owner hero", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		engine.State.CurrentTurn = 1
		p1 := engine.State.Players[1]
		p1.Hero = NewCardInstance(baseCard(t, "4311003"), 1, 1)
		p1.Hero.Position = &Position{Col: 1, Row: 1}
		p1.Units[1][1] = p1.Hero
		mark := NewCardInstance(baseCard(t, "2001102"), 1, 1)
		p1.Hand = []*CardInstance{
			mark,
			NewCardInstance(baseCard(t, "1021001"), 1, 1),
			NewCardInstance(baseCard(t, "1021002"), 1, 1),
			NewCardInstance(baseCard(t, "1021003"), 1, 1),
			NewCardInstance(baseCard(t, "1021004"), 1, 1),
			NewCardInstance(baseCard(t, "1021005"), 1, 1),
		}
		beforeLife := p1.Hero.CurrentLife
		engine.endTurn()
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "discard" {
			t.Fatalf("ending over hand limit should prompt discard, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 1, mark.InstanceID)
		if p1.Hero.CurrentLife != beforeLife-2 {
			t.Fatalf("discarding Jiuxiao Mark to hand limit should damage hero by 2, before=%d life=%d", beforeLife, p1.Hero.CurrentLife)
		}
		if len(p1.Graveyard) != 1 || p1.Graveyard[0] != mark {
			t.Fatalf("discarded mark should be in graveyard, grave=%v", cardsToInfo(p1.Graveyard))
		}
	})

	t.Run("council speaker shuffles marks and moves one to deck top on death", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		speaker := placeUnit(baseCard(t, "1521110"), 0, 0, 0, engine)
		p1 := engine.State.Players[1]
		engine.triggerEffects(TriggerOnEnter, speaker, nil, nil)
		if countCardsByNumber(p1.Deck, "2001102") != 4 {
			t.Fatalf("1521110 enter should shuffle four Jiuxiao Marks into opponent deck, deck=%v", cardsToInfo(p1.Deck))
		}

		deathEngine := setupReportedBugEngine(t)
		deathSpeaker := placeUnit(baseCard(t, "1521110"), 0, 0, 0, deathEngine)
		deathP1 := deathEngine.State.Players[1]
		other := NewCardInstance(baseCard(t, "1021001"), 1, 1)
		mark := NewCardInstance(baseCard(t, "2001102"), 1, 1)
		deathP1.Deck = []*CardInstance{other, mark}
		deathEngine.triggerEffects(TriggerOnDeath, deathSpeaker, nil, nil)
		if len(deathP1.Deck) == 0 || deathP1.Deck[0] != mark {
			t.Fatalf("1521110 death should move a Jiuxiao Mark from opponent deck to top, deck=%v", cardsToInfo(deathP1.Deck))
		}
	})
}

func countCardsByNumber(cards []*CardInstance, number string) int {
	count := 0
	for _, card := range cards {
		if card != nil && card.Card != nil && card.Card.Number == number {
			count++
		}
	}
	return count
}

func TestRoyalConflictLoadChoiceEffects(t *testing.T) {
	t.Run("five color coral gains two different non-arcane loads", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		coral := NewCardInstance(baseCard(t, "2021104"), 0, 1)
		engine.triggerEffects(TriggerOnEnter, coral, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "five_color_coral_load" {
			t.Fatalf("2021104 should prompt for two elements, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, model.ElementFire, model.ElementWater)
		load := effectiveElementsGain(coral)
		if load[model.ElementFire] != 1 || load[model.ElementWater] != 1 || load[model.ElementArcane] != 0 {
			t.Fatalf("2021104 should gain selected non-arcane loads, load=%v", load)
		}

		forgedEngine := setupReportedBugEngine(t)
		forgedCoral := NewCardInstance(baseCard(t, "2021104"), 0, 1)
		forgedEngine.triggerEffects(TriggerOnEnter, forgedCoral, nil, nil)
		if err := forgedEngine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{model.ElementFire, model.ElementFire},
		}}); err == nil {
			t.Fatal("2021104 should reject forged duplicate element selections")
		}
		if effectiveElementsGain(forgedCoral)[model.ElementFire] != 0 {
			t.Fatalf("2021104 should not gain load from rejected duplicate selection, load=%v", effectiveElementsGain(forgedCoral))
		}
	})

	t.Run("emerald fruit gives selected friendly companion a non-earth non-arcane load", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		target := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		fruit := NewCardInstance(baseCard(t, "2421108"), 0, 1)
		engine.triggerEffects(TriggerOnEnter, fruit, nil, nil)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "emerald_fruit_target" {
			t.Fatalf("2421108 should prompt for a friendly companion, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "emerald_fruit_element" {
			t.Fatalf("2421108 should prompt for a load element after target, pending=%+v", engine.State.PendingAction)
		}
		for _, candidate := range engine.State.PendingAction.Candidates {
			if candidate["instance_id"] == model.ElementEarth || candidate["instance_id"] == model.ElementArcane {
				t.Fatalf("2421108 should not offer earth or arcane load choices, candidates=%+v", engine.State.PendingAction.Candidates)
			}
		}
		resolvePendingSelection(t, engine, 0, model.ElementShadow)
		if effectiveElementsGain(target)[model.ElementShadow] != 1 {
			t.Fatalf("2421108 should add selected shadow load to target, load=%v", effectiveElementsGain(target))
		}
	})

	t.Run("lone star iron knight only buffs itself when isolated in front row", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		knight := placeUnit(baseCard(t, "1521115"), 0, 1, 0, engine)
		beforeLife := knight.CurrentLife
		engine.triggerEffects(TriggerOnEnter, knight, nil, nil)
		if knight.CurrentLife != beforeLife+1 || effectiveElementsGain(knight)[model.ElementLight] != knight.Card.ElementsGain[model.ElementLight]+1 {
			t.Fatalf("1521115 should gain +1 life and +1 light load when isolated front row, life=%d load=%v", knight.CurrentLife, effectiveElementsGain(knight))
		}

		blockedEngine := setupReportedBugEngine(t)
		blocked := placeUnit(baseCard(t, "1521115"), 0, 1, 0, blockedEngine)
		placeUnit(baseCard(t, "1021001"), 0, 1, 1, blockedEngine)
		blockedBeforeLife := blocked.CurrentLife
		blockedEngine.triggerEffects(TriggerOnEnter, blocked, nil, nil)
		if blocked.CurrentLife != blockedBeforeLife || effectiveElementsGain(blocked)[model.ElementLight] != blocked.Card.ElementsGain[model.ElementLight] {
			t.Fatalf("1521115 should not buff with adjacent companion, life=%d load=%v", blocked.CurrentLife, effectiveElementsGain(blocked))
		}

		dynamicFrontEngine := setupReportedBugEngine(t)
		dynamicFront := placeUnit(baseCard(t, "1521115"), 0, 1, 1, dynamicFrontEngine)
		dynamicFrontBeforeLife := dynamicFront.CurrentLife
		dynamicFrontEngine.triggerEffects(TriggerOnEnter, dynamicFront, nil, nil)
		if dynamicFront.CurrentLife != dynamicFrontBeforeLife+1 || effectiveElementsGain(dynamicFront)[model.ElementLight] != dynamicFront.Card.ElementsGain[model.ElementLight]+1 {
			t.Fatalf("1521115 should buff in row 1 when row 1 is the current front row, life=%d load=%v", dynamicFront.CurrentLife, effectiveElementsGain(dynamicFront))
		}

		backEngine := setupReportedBugEngine(t)
		placeUnit(baseCard(t, "1021001"), 0, 0, 0, backEngine)
		back := placeUnit(baseCard(t, "1521115"), 0, 1, 1, backEngine)
		backBeforeLife := back.CurrentLife
		backEngine.triggerEffects(TriggerOnEnter, back, nil, nil)
		if back.CurrentLife != backBeforeLife || effectiveElementsGain(back)[model.ElementLight] != back.Card.ElementsGain[model.ElementLight] {
			t.Fatalf("1521115 should not buff behind the current front row, life=%d load=%v", back.CurrentLife, effectiveElementsGain(back))
		}
	})

	t.Run("lone star soul gains shield and attack after enemy damage only while isolated", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		soul := placeUnit(baseCard(t, "1511102"), 0, 1, 1, engine)

		engine.dealDamageWithExtra(soul, 1, 0, map[string]any{"attacker": 1})
		if p0.Shield != 1 || soul.CurrentAttack != soul.Card.Attack+1 {
			t.Fatalf("1511102 should gain shield and attack after enemy damage while isolated, shield=%d attack=%d", p0.Shield, soul.CurrentAttack)
		}

		p0.Shield = 0
		engine.dealDamageWithExtra(soul, 1, 0, map[string]any{"attacker": 0})
		if p0.Shield != 0 || soul.CurrentAttack != soul.Card.Attack+1 {
			t.Fatalf("1511102 should ignore friendly-source damage, shield=%d attack=%d", p0.Shield, soul.CurrentAttack)
		}

		blockedEngine := setupReportedBugEngine(t)
		blockedP0 := blockedEngine.State.Players[0]
		blockedSoul := placeUnit(baseCard(t, "1511102"), 0, 1, 1, blockedEngine)
		placeUnit(baseCard(t, "1021001"), 0, 1, 0, blockedEngine)
		blockedEngine.dealDamageWithExtra(blockedSoul, 1, 0, map[string]any{"attacker": 1})
		if blockedP0.Shield != 0 || blockedSoul.CurrentAttack != blockedSoul.Card.Attack {
			t.Fatalf("1511102 should not trigger with adjacent friendly companions, shield=%d attack=%d", blockedP0.Shield, blockedSoul.CurrentAttack)
		}
	})
}

func TestRoyalConflictRuneAndDiscardEffects(t *testing.T) {
	t.Run("infusion runes permanently buff friendly spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		powerSkill := readySkill(baseCard(t, "3021001"), 0)
		p0.Skills[0] = powerSkill
		runeA := NewCardInstance(baseCard(t, "2021111"), 0, 1)
		if err := globalRegistry.GetBehavior("2021111").(OnUseItemBehavior).OnUseItem(&EffectContext{Engine: engine, Source: runeA, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("use 2021111: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "royal_infusion_rune_skill" {
			t.Fatalf("2021111 should prompt for a friendly spell, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, powerSkill.InstanceID)
		if powerSkill.PowerBonus != 2 {
			t.Fatalf("2021111 should add +2 power, bonus=%d", powerSkill.PowerBonus)
		}

		host := placeUnit(baseCard(t, "1011103"), 0, 0, 0, engine)
		bound := NewCardInstance(baseCard(t, "3001101"), 0, 1)
		bound.SlotIndex = -1
		host.BoundSkills = []*CardInstance{bound}
		runeB := NewCardInstance(baseCard(t, "2021112"), 0, 1)
		if err := globalRegistry.GetBehavior("2021112").(OnUseItemBehavior).OnUseItem(&EffectContext{Engine: engine, Source: runeB, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("use 2021112: %v", err)
		}
		resolvePendingSelection(t, engine, 0, bound.InstanceID)
		if bound.AttackBonus != 1 {
			t.Fatalf("2021112 should add +1 attack to selected bound spell, bonus=%d", bound.AttackBonus)
		}
	})

	t.Run("sky city thief discards one card from each player", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		p0.Hand = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
		p1.Hero = NewCardInstance(baseCard(t, "4311003"), 1, 1)
		p1.Hero.Position = &Position{Col: 1, Row: 1}
		p1.Units[1][1] = p1.Hero
		p1.Hand = []*CardInstance{NewCardInstance(baseCard(t, "2001102"), 1, 1)}
		beforeLife := p1.Hero.CurrentLife

		thief := placeUnit(baseCard(t, "1321107"), 0, 0, 0, engine)
		engine.triggerEffects(TriggerOnEnter, thief, nil, nil)
		if len(p0.Hand) != 0 || len(p0.Graveyard) != 1 {
			t.Fatalf("1321107 should discard one friendly hand card, hand=%v grave=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard))
		}
		if len(p1.Hand) != 0 || len(p1.Graveyard) != 1 || p1.Graveyard[0].Card.Number != "2001102" {
			t.Fatalf("1321107 should discard one enemy hand card, hand=%v grave=%v", cardsToInfo(p1.Hand), cardsToInfo(p1.Graveyard))
		}
		if p1.Hero.CurrentLife != beforeLife-2 {
			t.Fatalf("discarded Jiuxiao Mark from 1321107 should damage owner hero, before=%d life=%d", beforeLife, p1.Hero.CurrentLife)
		}
	})
}

func TestRoyalConflictPermanentSkillCostAndGraveyardBurstEffects(t *testing.T) {
	t.Run("water use cost reductions apply to selected friendly spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		waterSkill := readySkill(baseCard(t, "3221103"), 0)
		fireSkill := readySkill(baseCard(t, "3121102"), 0)
		p0.Skills[0] = waterSkill
		p0.Skills[1] = fireSkill

		spring := NewCardInstance(baseCard(t, "2221101"), 0, 1)
		if err := globalRegistry.GetBehavior("2221101").(OnUseItemBehavior).OnUseItem(&EffectContext{Engine: engine, Source: spring, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("use 2221101: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "royal_water_use_cost_reduction" {
			t.Fatalf("2221101 should prompt for a friendly spell, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, waterSkill.InstanceID)
		if got := engine.effectiveSkillUseCost(p0, waterSkill)[model.ElementWater]; got != 1 {
			t.Fatalf("2221101 should reduce water skill use cost by 1 water, cost=%v", engine.effectiveSkillUseCost(p0, waterSkill))
		}

		mirror := NewCardInstance(baseCard(t, "2221107"), 0, 1)
		mirror.SlotIndex = 0
		p0.Equipment[0] = mirror
		engine.triggerEffects(TriggerOnEnter, mirror, nil, nil)
		if engine.State.PendingAction == nil || len(engine.State.PendingAction.Candidates) != 1 || engine.State.PendingAction.Candidates[0]["instance_id"] != waterSkill.InstanceID {
			t.Fatalf("2221107 should offer only learned water spells, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, waterSkill.InstanceID)
		if got := engine.effectiveSkillUseCost(p0, waterSkill)[model.ElementWater]; got != 0 {
			t.Fatalf("2221107 should stack another -1 water use cost, cost=%v", engine.effectiveSkillUseCost(p0, waterSkill))
		}
		if got := engine.effectiveSkillUseCost(p0, fireSkill)[model.ElementFire]; got != 3 {
			t.Fatalf("water cost reducers should not affect fire spell, cost=%v", engine.effectiveSkillUseCost(p0, fireSkill))
		}
	})

	t.Run("mirrorsea spring requires a friendly spell before payment", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		spring := NewCardInstance(baseCard(t, "2221101"), 0, 1)
		p0.Hand = []*CardInstance{spring}
		p0.Elements[model.ElementWater] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": spring.InstanceID,
		}}); err == nil {
			t.Fatal("2221101 should fail before payment when there is no friendly spell")
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != spring || len(p0.Graveyard) != 0 || p0.Elements[model.ElementWater] != 1 {
			t.Fatalf("failed 2221101 should leave hand/grave/elements unchanged, hand=%v grave=%v elements=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard), p0.Elements)
		}
	})

	t.Run("dreamcatcher buffs learned spirit spells only", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		spiritSpell := readySkill(baseCard(t, "3421101"), 0)
		shadowSpiritSpell := readySkill(baseCard(t, "3621105"), 0)
		nonSpiritSpell := readySkill(baseCard(t, "3221103"), 0)
		p0.Skills[0] = spiritSpell
		p0.Skills[1] = shadowSpiritSpell
		p0.Skills[2] = nonSpiritSpell

		dreamcatcher := NewCardInstance(baseCard(t, "2421103"), 0, 1)
		dreamcatcher.SlotIndex = 0
		p0.Equipment[0] = dreamcatcher
		engine.triggerEffects(TriggerOnEnter, dreamcatcher, nil, nil)
		if spiritSpell.PowerBonus != 2 || shadowSpiritSpell.PowerBonus != 2 {
			t.Fatalf("2421103 should give learned spirit spells +2 power, spirit=%d shadow=%d", spiritSpell.PowerBonus, shadowSpiritSpell.PowerBonus)
		}
		if nonSpiritSpell.PowerBonus != 0 {
			t.Fatalf("2421103 should not buff non-spirit spells, bonus=%d", nonSpiritSpell.PowerBonus)
		}
	})

	t.Run("dark burst scroll exiles five or more shadow companions for shadow elements", func(t *testing.T) {
		failEngine := setupReportedBugEngine(t)
		failP0 := failEngine.State.Players[0]
		failScroll := NewCardInstance(baseCard(t, "2621111"), 0, 1)
		failP0.Hand = []*CardInstance{failScroll}
		failP0.Elements[model.ElementShadow] = 4
		for i := 0; i < 4; i++ {
			failP0.Graveyard = append(failP0.Graveyard, NewCardInstance(baseCard(t, "1621103"), 0, 1))
		}
		if err := failEngine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": failScroll.InstanceID,
		}}); err == nil {
			t.Fatal("2621111 should require at least five shadow companions in graveyard")
		}
		if len(failP0.Hand) != 1 || len(failP0.Graveyard) != 4 {
			t.Fatalf("failed 2621111 should leave zones unchanged, hand=%v grave=%v", cardsToInfo(failP0.Hand), cardsToInfo(failP0.Graveyard))
		}

		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2621111"), 0, 1)
		nonShadow := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementShadow] = 4
		p0.Graveyard = []*CardInstance{nonShadow}
		for i := 0; i < 5; i++ {
			p0.Graveyard = append(p0.Graveyard, NewCardInstance(baseCard(t, "1621112"), 0, 1))
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use 2621111: %v", err)
		}
		if len(p0.Exile) != 5 {
			t.Fatalf("2621111 should exile five shadow companions, exile=%v grave=%v", cardsToInfo(p0.Exile), cardsToInfo(p0.Graveyard))
		}
		if len(p0.Graveyard) != 2 || p0.Graveyard[0] != nonShadow || p0.Graveyard[1] != scroll {
			t.Fatalf("2621111 should leave non-shadow card and itself in graveyard, grave=%v", cardsToInfo(p0.Graveyard))
		}
		if p0.Elements[model.ElementShadow] != 10 {
			t.Fatalf("2621111 should spend 4 shadow then gain 10 shadow, elements=%v", p0.Elements)
		}
	})
}

func TestRoyalConflictRaiderSearchItemEffects(t *testing.T) {
	t.Run("black sail raider requires a searchable raider companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2221105"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
		p0.Elements[model.ElementWater] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err == nil {
			t.Fatal("2221105 should require a searchable raider companion")
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != scroll || len(p0.Graveyard) != 0 || p0.Elements[model.ElementWater] != 1 {
			t.Fatalf("failed 2221105 should leave zones/elements unchanged, hand=%v grave=%v elements=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard), p0.Elements)
		}
	})

	t.Run("black sail raider searches without discount when no raider is on field", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2221105"), 0, 1)
		target := NewCardInstance(baseCard(t, "1221101"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1), target}
		p0.Elements[model.ElementWater] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use 2221105: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "black_sail_raider_search" || len(engine.State.PendingAction.Candidates) != 1 {
			t.Fatalf("2221105 should ask which raider companion to search, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if engine.State.PendingAction != nil {
			t.Fatalf("2221105 should not ask for a discount without a raider on field, pending=%+v", engine.State.PendingAction)
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != target || len(p0.Graveyard) != 1 || p0.Graveyard[0] != scroll {
			t.Fatalf("2221105 should search target to hand and move itself to graveyard, hand=%v grave=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard))
		}
		if target.Statuses["入场费用"+model.ElementWater+"-1"] != 0 || target.Statuses["入场费用"+model.ElementShadow+"-1"] != 0 {
			t.Fatalf("2221105 should not discount without a raider on field, statuses=%v", target.Statuses)
		}
	})

	t.Run("black sail raider discounts searched raider when a raider is on field", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2221105"), 0, 1)
		target := NewCardInstance(baseCard(t, "1221101"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Deck = []*CardInstance{target}
		p0.Elements[model.ElementWater] = 2
		placeUnit(baseCard(t, "1221003"), 0, 0, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use 2221105 with raider on field: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "black_sail_raider_search" || len(engine.State.PendingAction.Candidates) != 1 {
			t.Fatalf("2221105 should ask which raider companion to search, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "black_sail_raider_discount" || len(engine.State.PendingAction.Candidates) != 2 {
			t.Fatalf("2221105 should ask which entry cost element to reduce, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, model.ElementShadow)
		if target.Statuses["入场费用"+model.ElementShadow+"-1"] != 1 || target.Statuses["入场费用"+model.ElementWater+"-1"] != 0 {
			t.Fatalf("2221105 should apply selected shadow discount only, statuses=%v", target.Statuses)
		}
		if cost := engine.effectiveCardPlayCost(p0, target); cost[model.ElementShadow] != 0 || cost[model.ElementWater] != 4 {
			t.Fatalf("2221105 discount should reduce target shadow entry cost only, cost=%v", cost)
		}
	})
}

func TestRoyalConflictAirAndMoonlightItemEffects(t *testing.T) {
	t.Run("burnout scroll consumes a ready fire companion for its entry cost", func(t *testing.T) {
		failEngine := setupReportedBugEngine(t)
		failP0 := failEngine.State.Players[0]
		failScroll := NewCardInstance(baseCard(t, "2121108"), 0, 1)
		failP0.Hand = []*CardInstance{failScroll}
		failP0.Elements[model.ElementFire] = 3
		tappedFire := placeUnit(baseCard(t, "1121101"), 0, 0, 0, failEngine)
		tappedFire.IsHorizontal = true
		if err := failEngine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": failScroll.InstanceID,
		}}); err == nil {
			t.Fatal("2121108 should require a ready friendly fire companion")
		}
		if len(failP0.Hand) != 1 || failP0.Hand[0] != failScroll || len(failP0.Graveyard) != 0 || failP0.Elements[model.ElementFire] != 3 {
			t.Fatalf("failed 2121108 should leave hand/grave/elements unchanged, hand=%v grave=%v elements=%v", cardsToInfo(failP0.Hand), cardsToInfo(failP0.Graveyard), failP0.Elements)
		}

		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2121108"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementFire] = 5
		target := placeUnit(baseCard(t, "1121101"), 0, 0, 0, engine)
		target.IsHorizontal = false
		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use 2121108: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "burnout_scroll_consume_fire_companion" || len(engine.State.PendingAction.Candidates) != 1 {
			t.Fatalf("2121108 should prompt with ready fire companions only, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if !target.IsHorizontal {
			t.Fatal("2121108 should consume the selected fire companion")
		}
		if p0.Elements[model.ElementFire] != 6 {
			t.Fatalf("2121108 should spend 1 fire then gain target entry cost 2 fire, elements=%v", p0.Elements)
		}
		if len(p0.Hand) != 0 || len(p0.Graveyard) != 1 || p0.Graveyard[0] != scroll {
			t.Fatalf("2121108 should move itself from hand to graveyard, hand=%v grave=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard))
		}

		burnEngine := setupReportedBugEngine(t)
		burnP0 := burnEngine.State.Players[0]
		burnScroll := NewCardInstance(baseCard(t, "2121108"), 0, 1)
		burnP0.Hand = []*CardInstance{burnScroll}
		burnP0.Elements[model.ElementFire] = 1
		fireSprite := placeUnit(baseCard(t, "1121001"), 0, 0, 0, burnEngine)
		if err := burnEngine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": burnScroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use 2121108 with fire sprite: %v", err)
		}
		resolvePendingSelection(t, burnEngine, 0, fireSprite.InstanceID)
		if fireSprite.Statuses[StatusBurn] != 1 {
			t.Fatalf("2121108 should trigger consume effects on selected companion, statuses=%v", fireSprite.Statuses)
		}

		staleEngine := setupReportedBugEngine(t)
		staleP0 := staleEngine.State.Players[0]
		staleScroll := NewCardInstance(baseCard(t, "2121108"), 0, 1)
		staleP0.Hand = []*CardInstance{staleScroll}
		staleP0.Elements[model.ElementFire] = 5
		staleTarget := placeUnit(baseCard(t, "1121101"), 0, 0, 0, staleEngine)
		if err := staleEngine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": staleScroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use 2121108 before stale target: %v", err)
		}
		staleTarget.IsHorizontal = true
		resolvePendingSelection(t, staleEngine, 0, staleTarget.InstanceID)
		if staleP0.Elements[model.ElementFire] != 4 {
			t.Fatalf("2121108 stale horizontal target should not gain entry cost, elements=%v", staleP0.Elements)
		}
	})

	t.Run("elegy scroll flips the first shadow deathrattle companion and discounts it from shadow graveyard", func(t *testing.T) {
		failEngine := setupReportedBugEngine(t)
		failP0 := failEngine.State.Players[0]
		failScroll := NewCardInstance(baseCard(t, "2621109"), 0, 1)
		failP0.Hand = []*CardInstance{failScroll}
		failP0.Elements[model.ElementShadow] = 2
		failP0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
		if err := failEngine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": failScroll.InstanceID,
		}}); err == nil {
			t.Fatal("2621109 should require a searchable shadow companion with deathrattle")
		}
		if len(failP0.Hand) != 1 || failP0.Hand[0] != failScroll || len(failP0.Graveyard) != 0 || failP0.Elements[model.ElementShadow] != 2 {
			t.Fatalf("failed 2621109 should leave hand/grave/elements unchanged, hand=%v grave=%v elements=%v", cardsToInfo(failP0.Hand), cardsToInfo(failP0.Graveyard), failP0.Elements)
		}

		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2621109"), 0, 1)
		firstTarget := NewCardInstance(baseCard(t, "1621112"), 0, 1)
		secondTarget := NewCardInstance(baseCard(t, "1621113"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1), firstTarget, secondTarget}
		p0.Graveyard = []*CardInstance{NewCardInstance(baseCard(t, "1621103"), 0, 1)}
		p0.Elements[model.ElementShadow] = 2
		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use 2621109: %v", err)
		}
		if engine.State.PendingAction != nil {
			t.Fatalf("2621109 should flip directly without exposing deck choices, pending=%+v", engine.State.PendingAction)
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != firstTarget {
			t.Fatalf("2621109 should flip the first matching companion to hand, hand=%v", cardsToInfo(p0.Hand))
		}
		if firstTarget.Statuses["入场费用"+model.ElementShadow+"-1"] != 1 {
			t.Fatalf("2621109 should discount flipped card when graveyard has shadow companion, statuses=%v", firstTarget.Statuses)
		}
		if secondTarget.Statuses["入场费用"+model.ElementShadow+"-1"] != 0 || !containsCardInstance(p0.Deck, secondTarget) {
			t.Fatalf("2621109 should not reveal or discount later matching cards, deck=%v second_statuses=%v", cardsToInfo(p0.Deck), secondTarget.Statuses)
		}
		if len(p0.Graveyard) != 2 || p0.Graveyard[1] != scroll || p0.Elements[model.ElementShadow] != 1 {
			t.Fatalf("2621109 should spend 1 shadow and move itself to graveyard, grave=%v elements=%v", cardsToInfo(p0.Graveyard), p0.Elements)
		}
	})

	t.Run("wind cycle consumes and sacrifices itself to shuffle selected air graveyard cards", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		cycle := NewCardInstance(baseCard(t, "2321102"), 0, 1)
		cycle.SlotIndex = 0
		cycle.IsHorizontal = false
		p0.Equipment[0] = cycle
		airA := NewCardInstance(baseCard(t, "1321108"), 0, 1)
		airB := NewCardInstance(baseCard(t, "2321103"), 0, 1)
		nonAir := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Graveyard = []*CardInstance{airA, nonAir, airB}

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  cycle.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use 2321102: %v", err)
		}
		if p0.Equipment[0] != nil || len(p0.Graveyard) != 4 || p0.Graveyard[3] != cycle {
			t.Fatalf("2321102 should sacrifice itself before selection, equipment=%v grave=%v", p0.Equipment[0], cardsToInfo(p0.Graveyard))
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "wind_cycle_shuffle_air" || len(engine.State.PendingAction.Candidates) != 2 {
			t.Fatalf("2321102 should prompt with air graveyard cards only, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, airA.InstanceID, airB.InstanceID)
		if countCardsByNumber(p0.Deck, airA.Card.Number) != 1 || countCardsByNumber(p0.Deck, airB.Card.Number) != 1 {
			t.Fatalf("2321102 should shuffle selected air cards into deck, deck=%v", cardsToInfo(p0.Deck))
		}
		if len(p0.Graveyard) != 2 || p0.Graveyard[0] != nonAir || p0.Graveyard[1] != cycle {
			t.Fatalf("2321102 should leave non-selected/non-air cards and itself in graveyard, grave=%v", cardsToInfo(p0.Graveyard))
		}
	})

	t.Run("thunder breath gains air when used or discarded", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		used := NewCardInstance(baseCard(t, "2321103"), 0, 1)
		p0.Hand = []*CardInstance{used}
		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": used.InstanceID,
		}}); err != nil {
			t.Fatalf("use 2321103: %v", err)
		}
		if p0.Elements[model.ElementAir] != 1 || len(p0.Graveyard) != 1 || p0.Graveyard[0] != used {
			t.Fatalf("used 2321103 should gain 1 air and enter graveyard, elements=%v grave=%v", p0.Elements, cardsToInfo(p0.Graveyard))
		}

		discarded := NewCardInstance(baseCard(t, "2321103"), 0, 1)
		p0.Hand = []*CardInstance{discarded}
		before := p0.Elements[model.ElementAir]
		engine.discardHandCardAt(0, 0)
		if p0.Elements[model.ElementAir] != before+1 || len(p0.Hand) != 0 || p0.Graveyard[len(p0.Graveyard)-1] != discarded {
			t.Fatalf("discarded 2321103 should gain 1 air and enter graveyard, elements=%v hand=%v grave=%v", p0.Elements, cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard))
		}
	})

	t.Run("moonlight dust destroys set counters or removes front stealth", func(t *testing.T) {
		failEngine := setupReportedBugEngine(t)
		failP0 := failEngine.State.Players[0]
		failDust := NewCardInstance(baseCard(t, "2521102"), 0, 1)
		failP0.Hand = []*CardInstance{failDust}
		failP0.Elements[model.ElementLight] = 1
		if err := failEngine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": failDust.InstanceID,
		}}); err == nil {
			t.Fatal("2521102 should require enemy set counters or stealthy front enemies")
		}
		if len(failP0.Hand) != 1 || len(failP0.Graveyard) != 0 || failP0.Elements[model.ElementLight] != 1 {
			t.Fatalf("failed 2521102 should leave zones and elements unchanged, hand=%v grave=%v elements=%v", cardsToInfo(failP0.Hand), cardsToInfo(failP0.Graveyard), failP0.Elements)
		}

		counterEngine := setupReportedBugEngine(t)
		counterP0 := counterEngine.State.Players[0]
		counterP1 := counterEngine.State.Players[1]
		dust := NewCardInstance(baseCard(t, "2521102"), 0, 1)
		counterP0.Hand = []*CardInstance{dust}
		counterP0.Elements[model.ElementLight] = 1
		firstCounter := NewCardInstance(baseCard(t, "2021113"), 1, 1)
		secondCounter := NewCardInstance(baseCard(t, "2021114"), 1, 1)
		firstCounter.IsSetCounter = true
		secondCounter.IsSetCounter = true
		counterP1.Equipment[0] = firstCounter
		counterP1.Equipment[1] = secondCounter
		if err := counterEngine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": dust.InstanceID,
		}}); err != nil {
			t.Fatalf("use 2521102 counters: %v", err)
		}
		resolvePendingSelection(t, counterEngine, 0, "destroy_counters")
		if counterP1.Equipment[0] != nil || counterP1.Equipment[1] != nil || len(counterP1.Graveyard) != 2 {
			t.Fatalf("2521102 should destroy all enemy set counters, equipment=%v grave=%v", counterP1.Equipment, cardsToInfo(counterP1.Graveyard))
		}

		stealthEngine := setupReportedBugEngine(t)
		stealthP0 := stealthEngine.State.Players[0]
		stealthP1 := stealthEngine.State.Players[1]
		stealthDust := NewCardInstance(baseCard(t, "2521102"), 0, 1)
		stealthP0.Hand = []*CardInstance{stealthDust}
		stealthP0.Elements[model.ElementLight] = 1
		frontRow := 0
		front := placeUnit(baseCard(t, "1021001"), 1, 0, frontRow, stealthEngine)
		front.Statuses[StatusStealth] = 2
		backRow := 1
		if backRow == stealthP1.GetFrontRow() {
			backRow = 0
		}
		back := placeUnit(baseCard(t, "1021002"), 1, 1, backRow, stealthEngine)
		back.Statuses[StatusStealth] = 2
		if err := stealthEngine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": stealthDust.InstanceID,
		}}); err != nil {
			t.Fatalf("use 2521102 stealth: %v", err)
		}
		resolvePendingSelection(t, stealthEngine, 0, "remove_front_stealth")
		if front.Statuses[StatusStealth] != 0 {
			t.Fatalf("2521102 should remove stealth from front enemy, statuses=%v", front.Statuses)
		}
		if back.Statuses[StatusStealth] != 2 {
			t.Fatalf("2521102 should not remove stealth from non-front enemy, statuses=%v", back.Statuses)
		}
	})
}

func TestRoyalConflictUtilityCompanionAndHeroEffects(t *testing.T) {
	t.Run("private teacher learns cheap skills for free and can replace a vertical skill", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		teacher := NewCardInstance(baseCard(t, "1021101"), 0, 1)
		cheap := NewCardInstance(baseCard(t, "3021005"), 0, 1)
		p0.SkillPool = []*CardInstance{cheap}
		p0.Elements = map[string]int{}

		if err := (Card1021101PrivateTeacher{}).OnEnter(&EffectContext{Engine: engine, Source: teacher, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("private teacher enter: %v", err)
		}
		resolvePendingSelection(t, engine, 0, cheap.InstanceID)
		if len(p0.SkillPool) != 0 || p0.Skills[0] != cheap || !cheap.IsHorizontal {
			t.Fatalf("private teacher should learn cheap skill without cost, pool=%v skills=%v horizontal=%v", cardsToInfo(p0.SkillPool), cardsToInfo(p0.Skills[:]), cheap.IsHorizontal)
		}
		if len(p0.Elements) != 0 {
			t.Fatalf("private teacher should not spend elements, elements=%v", p0.Elements)
		}

		fullEngine := setupReportedBugEngine(t)
		fullP0 := fullEngine.State.Players[0]
		replacement := NewCardInstance(baseCard(t, "3021005"), 0, 1)
		fullP0.SkillPool = []*CardInstance{replacement}
		oldSkill := readySkill(baseCard(t, "3021102"), 0)
		for i := 0; i < skillSlotCapacity(fullP0); i++ {
			fullP0.Skills[i] = readySkill(baseCard(t, "3021005"), 0)
		}
		fullP0.Skills[2] = oldSkill
		if err := (Card1021101PrivateTeacher{}).OnEnter(&EffectContext{Engine: fullEngine, Source: teacher, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("private teacher full enter: %v", err)
		}
		if fullEngine.State.PendingAction == nil || len(fullEngine.State.PendingAction.Candidates) == 0 {
			t.Fatalf("private teacher should offer replacement choices, pending=%+v", fullEngine.State.PendingAction)
		}
		choice := ""
		for _, candidate := range fullEngine.State.PendingAction.Candidates {
			if replaceID, _ := candidate["replace_id"].(string); replaceID == oldSkill.InstanceID {
				choice, _ = candidate["instance_id"].(string)
				break
			}
		}
		if choice == "" {
			t.Fatalf("private teacher should offer a choice replacing the selected old skill, pending=%+v", fullEngine.State.PendingAction)
		}
		resolvePendingSelection(t, fullEngine, 0, choice)
		if replacement.SlotIndex < 0 || fullP0.Skills[replacement.SlotIndex] != replacement {
			t.Fatalf("private teacher should place replacement in a skill slot, slot=%d skills=%v", replacement.SlotIndex, cardsToInfo(fullP0.Skills[:]))
		}
		if !containsCardInstance(fullP0.SkillPool, oldSkill) {
			t.Fatalf("private teacher should return replaced skill to pool, pool=%v", cardsToInfo(fullP0.SkillPool))
		}
	})

	t.Run("lone star hero entry cost increases with other hand cards", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		hero := NewCardInstance(baseCard(t, "1021111"), 0, 1)
		otherA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		otherB := NewCardInstance(baseCard(t, "3021002"), 0, 1)
		p0.Hand = []*CardInstance{hero, otherA, otherB}

		if got := engine.effectiveCardPlayCost(p0, hero)[model.ElementArcane]; got != 5 {
			t.Fatalf("1021111 should cost base 3 plus 2 other hand cards, cost=%v", engine.effectiveCardPlayCost(p0, hero))
		}
		p0.Hand = []*CardInstance{hero}
		if got := engine.effectiveCardPlayCost(p0, hero)[model.ElementArcane]; got != 3 {
			t.Fatalf("1021111 should use base cost with no other hand cards, cost=%v", engine.effectiveCardPlayCost(p0, hero))
		}
	})

	t.Run("radiant guard is free after friendly companion was damaged last turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		guard := NewCardInstance(baseCard(t, "1521107"), 0, 1)
		wounded := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)

		engine.dealDamageWithExtra(wounded, 1, 0, map[string]any{"damage_source": "test", "attacker": 1})
		if !p0.FriendlyUnitDamagedThisTurn || p0.FriendlyUnitDamagedLastTurn {
			t.Fatalf("friendly unit damage should be tracked for this turn only before rolling, this=%v last=%v", p0.FriendlyUnitDamagedThisTurn, p0.FriendlyUnitDamagedLastTurn)
		}
		if cost := engine.effectiveCardPlayCost(p0, guard); cost[model.ElementLight] != guard.Card.ElementsCost[model.ElementLight] {
			t.Fatalf("1521107 should not be free until the next turn, cost=%v", cost)
		}

		engine.rollFriendlyUnitDamageHistory()
		if p0.FriendlyUnitDamagedThisTurn || !p0.FriendlyUnitDamagedLastTurn {
			t.Fatalf("friendly unit damage history should roll to last turn, this=%v last=%v", p0.FriendlyUnitDamagedThisTurn, p0.FriendlyUnitDamagedLastTurn)
		}
		if cost := engine.effectiveCardPlayCost(p0, guard); cost[model.ElementLight] != 0 {
			t.Fatalf("1521107 should be free after a friendly companion was damaged last turn, cost=%v", cost)
		}

		placeUnit(baseCard(t, "1111103"), 1, 0, 0, engine)
		if cost := engine.effectiveCardPlayCost(p0, guard); cost[model.ElementLight] != 0 || cost[model.ElementArcane] != 0 {
			t.Fatalf("1521107 free entry should override other entry cost increases, cost=%v", cost)
		}

		nextEngine := setupReportedBugEngine(t)
		nextP0 := nextEngine.State.Players[0]
		hero := NewCardInstance(baseCard(t, "4011001"), 0, 1)
		hero.Position = &Position{Col: 1, Row: 1}
		nextP0.Hero = hero
		nextEngine.dealDamageWithExtra(hero, 1, 0, map[string]any{"damage_source": "test", "attacker": 1})
		nextEngine.rollFriendlyUnitDamageHistory()
		if nextP0.FriendlyUnitDamagedLastTurn {
			t.Fatalf("1521107 should care about friendly companions, not hero damage")
		}
	})

	t.Run("sting frog boosts later spells after friendly drive or charge skills", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1021114"), 0, 0, 0, engine)
		drive := readySkill(baseCard(t, "3321109"), 0)
		laterSpell := readySkill(baseCard(t, "3021008"), 0)
		p0.Skills[0] = drive
		p0.Skills[1] = laterSpell
		p0.Elements[model.ElementAir] = 2
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": drive.InstanceID,
			"target_type": "unit",
			"target_col":  float64(target.Position.Col),
			"target_row":  float64(target.Position.Row),
		}}); err != nil {
			t.Fatalf("cast drive spell with frog: %v", err)
		}
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModSkillPowerBonus || p0.TempModifiers[0].Amount != 1 {
			t.Fatalf("1021114 should add current-turn spell power modifier, modifiers=%+v", p0.TempModifiers)
		}
		if got := engine.effectiveSpellPower(0, laterSpell, nil); got != laterSpell.Card.Power+1 {
			t.Fatalf("1021114 should boost later spells this turn, got=%d", got)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve drive spell: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end turn after frog trigger: %v", err)
		}
		if len(p0.TempModifiers) != 0 {
			t.Fatalf("1021114 modifier should expire at turn end, modifiers=%+v", p0.TempModifiers)
		}

		nonTriggerEngine := setupReportedBugEngine(t)
		nonP0 := nonTriggerEngine.State.Players[0]
		placeUnit(baseCard(t, "1021114"), 0, 0, 0, nonTriggerEngine)
		nonTrigger := readySkill(baseCard(t, "3121001"), 0)
		checkSpell := readySkill(baseCard(t, "3021008"), 0)
		nonP0.Skills[0] = nonTrigger
		nonP0.Skills[1] = checkSpell
		nonP0.Elements[model.ElementFire] = 2
		nonTarget := placeUnit(baseCard(t, "1021001"), 1, 1, 0, nonTriggerEngine)
		if err := nonTriggerEngine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": nonTrigger.InstanceID,
			"target_type": "unit",
			"target_col":  float64(nonTarget.Position.Col),
			"target_row":  float64(nonTarget.Position.Row),
		}}); err != nil {
			t.Fatalf("cast non-drive spell with frog: %v", err)
		}
		if len(nonP0.TempModifiers) != 0 {
			t.Fatalf("1021114 should ignore non-drive/non-charge skills, modifiers=%+v", nonP0.TempModifiers)
		}

		scrollEngine := setupReportedBugEngine(t)
		scrollP0 := scrollEngine.State.Players[0]
		placeUnit(baseCard(t, "1021114"), 0, 0, 0, scrollEngine)
		scroll := NewCardInstance(baseCard(t, "2121003"), 0, 1)
		scrollEngine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, scroll, map[string]any{"cast_player": 0})
		if len(scrollP0.TempModifiers) != 0 {
			t.Fatalf("1021114 should not trigger from drive/charge spell scroll items, modifiers=%+v", scrollP0.TempModifiers)
		}
	})

	t.Run("greedy tyrant increases both players hand card entry costs", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		tyrant := placeUnit(baseCard(t, "1111103"), 0, 0, 0, engine)
		ownHandCard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		enemyHandCard := NewCardInstance(baseCard(t, "1021001"), 1, 1)
		notHandCard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Hand = []*CardInstance{ownHandCard}
		p1.Hand = []*CardInstance{enemyHandCard}

		if got := engine.effectiveCardPlayCost(p0, ownHandCard)[model.ElementArcane]; got != ownHandCard.Card.ElementsCost[model.ElementArcane]+1 {
			t.Fatalf("1111103 should increase own hand card entry cost, cost=%v", engine.effectiveCardPlayCost(p0, ownHandCard))
		}
		if got := engine.effectiveCardPlayCost(p1, enemyHandCard)[model.ElementArcane]; got != enemyHandCard.Card.ElementsCost[model.ElementArcane]+1 {
			t.Fatalf("1111103 should increase enemy hand card entry cost, cost=%v", engine.effectiveCardPlayCost(p1, enemyHandCard))
		}
		if got := engine.effectiveCardPlayCost(p0, notHandCard)[model.ElementArcane]; got != notHandCard.Card.ElementsCost[model.ElementArcane] {
			t.Fatalf("1111103 should not increase non-hand card entry cost, cost=%v", engine.effectiveCardPlayCost(p0, notHandCard))
		}
		tyrant.Statuses[StatusPetrify] = 1
		if got := engine.effectiveCardPlayCost(p0, ownHandCard)[model.ElementArcane]; got != ownHandCard.Card.ElementsCost[model.ElementArcane] {
			t.Fatalf("petrified 1111103 should not increase hand card entry cost, cost=%v", engine.effectiveCardPlayCost(p0, ownHandCard))
		}
	})

	t.Run("alchemy apprentice converts one arcane into two non-arcane elements", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		apprentice := NewCardInstance(baseCard(t, "1021108"), 0, 1)
		apprentice.IsHorizontal = false
		p0.Elements[model.ElementArcane] = 1

		if err := (Card1021108AlchemyApprentice{}).OnPerTurn(&EffectContext{Engine: engine, Source: apprentice, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("alchemy apprentice ability: %v", err)
		}
		if !apprentice.IsHorizontal || p0.Elements[model.ElementArcane] != 0 {
			t.Fatalf("alchemy apprentice should consume itself and spend one arcane, horizontal=%v elements=%v", apprentice.IsHorizontal, p0.Elements)
		}
		resolvePendingSelection(t, engine, 0, model.ElementFire+"#1", model.ElementFire+"#2")
		if p0.Elements[model.ElementFire] != 2 {
			t.Fatalf("alchemy apprentice should allow choosing the same non-arcane element twice, elements=%v", p0.Elements)
		}

		failEngine := setupReportedBugEngine(t)
		failP0 := failEngine.State.Players[0]
		failApprentice := NewCardInstance(baseCard(t, "1021108"), 0, 1)
		failApprentice.IsHorizontal = false
		if err := (Card1021108AlchemyApprentice{}).OnPerTurn(&EffectContext{Engine: failEngine, Source: failApprentice, PlayerID: 0, OpponentID: 1}); err == nil {
			t.Fatal("alchemy apprentice should require one arcane element")
		}
		if failApprentice.IsHorizontal || failP0.Elements[model.ElementArcane] != 0 {
			t.Fatalf("failed alchemy apprentice should not mutate state, horizontal=%v elements=%v", failApprentice.IsHorizontal, failP0.Elements)
		}
	})

	t.Run("thunderlight warrior chooses one reward per thunderlight item", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		warrior := NewCardInstance(baseCard(t, "1321111"), 0, 1)
		helm := NewCardInstance(baseCard(t, "2321104"), 0, 1)
		armor := NewCardInstance(baseCard(t, "2321105"), 0, 1)
		airOnly := NewCardInstance(baseCard(t, "2321107"), 0, 1)
		p0.Equipment[0] = helm
		p0.Equipment[1] = armor
		p0.Equipment[2] = airOnly

		if err := (Card1321111ThunderlightWarrior{}).OnEnter(&EffectContext{Engine: engine, Source: warrior, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("thunderlight warrior enter: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "thunderlight_warrior_rewards" || engine.State.PendingAction.MinSelect != 2 || engine.State.PendingAction.MaxSelect != 2 || len(engine.State.PendingAction.Candidates) != 8 {
			t.Fatalf("1321111 should offer four rewards per thunderlight item, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, "life#0", "air#1")
		if warrior.CurrentLife != warrior.Card.Life+2 {
			t.Fatalf("1321111 should apply selected life reward, life=%d", warrior.CurrentLife)
		}
		load := effectiveElementsGain(warrior)
		if warrior.AttackBonus != 0 || load[model.ElementAir] != warrior.Card.ElementsGain[model.ElementAir]+1 || load[model.ElementLight] != warrior.Card.ElementsGain[model.ElementLight] {
			t.Fatalf("1321111 should apply selected air load only, attack=%d load=%v", warrior.AttackBonus, load)
		}

		emptyEngine := setupReportedBugEngine(t)
		emptyWarrior := NewCardInstance(baseCard(t, "1321111"), 0, 1)
		emptyEngine.State.Players[0].Equipment[0] = NewCardInstance(baseCard(t, "2321107"), 0, 1)
		if err := (Card1321111ThunderlightWarrior{}).OnEnter(&EffectContext{Engine: emptyEngine, Source: emptyWarrior, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("thunderlight warrior no-op enter: %v", err)
		}
		if emptyEngine.State.PendingAction != nil {
			t.Fatalf("1321111 should no-op without thunderlight items, pending=%+v", emptyEngine.State.PendingAction)
		}
	})

	t.Run("thunderlight armor buffs drive and focus spells with three thunderlight items", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2321101"), 0, 1)
		p0.Equipment[1] = NewCardInstance(baseCard(t, "2321104"), 0, 1)
		p0.Equipment[2] = NewCardInstance(baseCard(t, "2321105"), 0, 1)
		drive := readySkill(baseCard(t, "3321101"), 0)
		focus := readySkill(baseCard(t, "3321103"), 0)
		mystery := readySkill(baseCard(t, "3521106"), 0)

		if got := engine.effectiveSpellPower(0, drive, nil); got != drive.Card.Power+2 {
			t.Fatalf("2321105 should buff drive spells by 2, got=%d", got)
		}
		if got := engine.effectiveSpellPower(0, focus, nil); got != focus.Card.Power+2 {
			t.Fatalf("2321105 should buff focus spells by 2, got=%d", got)
		}
		if got := engine.effectiveSpellPower(0, mystery, nil); got != mystery.Card.Power {
			t.Fatalf("2321105 should not buff mystery spells, got=%d", got)
		}

		shortEngine := setupReportedBugEngine(t)
		shortP0 := shortEngine.State.Players[0]
		shortP0.Equipment[0] = NewCardInstance(baseCard(t, "2321104"), 0, 1)
		shortP0.Equipment[1] = NewCardInstance(baseCard(t, "2321105"), 0, 1)
		shortP0.Equipment[2] = NewCardInstance(baseCard(t, "2321107"), 0, 1)
		if got := shortEngine.effectiveSpellPower(0, readySkill(baseCard(t, "3321101"), 0), nil); got != baseCard(t, "3321101").Power {
			t.Fatalf("2321105 should require three thunderlight items, got=%d", got)
		}
	})

	t.Run("thunderlight crown prayer buffs the next focus spell only", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		crown := NewCardInstance(baseCard(t, "2321104"), 0, 1)
		p0.Equipment[0] = crown
		focus := readySkill(baseCard(t, "3321103"), 0)
		drive := readySkill(baseCard(t, "3321101"), 0)

		if !cardHasActivePrayer(crown) {
			t.Fatal("2321104 should expose a prayer ability")
		}
		if err := (Card2321104ThunderlightCrown{}).OnPerTurn(&EffectContext{Engine: engine, Source: crown, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("thunderlight crown prayer: %v", err)
		}
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextTaggedSpellPowerBonus || p0.TempModifiers[0].Status != "聚能" || p0.TempModifiers[0].Amount != 1 {
			t.Fatalf("2321104 should create a next focus spell power modifier, modifiers=%v", p0.TempModifiers)
		}
		if got := engine.effectiveSpellPower(0, drive, nil); got != drive.Card.Power {
			t.Fatalf("2321104 should not buff drive spells, got=%d", got)
		}
		engine.consumeNextSpellPowerBonuses(p0, drive)
		if len(p0.TempModifiers) != 1 {
			t.Fatalf("2321104 modifier should not be consumed by drive spells, modifiers=%v", p0.TempModifiers)
		}
		if got := engine.effectiveSpellPower(0, focus, nil); got != focus.Card.Power+1 {
			t.Fatalf("2321104 should buff the next focus spell, got=%d", got)
		}
		engine.consumeNextSpellPowerBonuses(p0, focus)
		if len(p0.TempModifiers) != 0 {
			t.Fatalf("2321104 modifier should be consumed by the focus spell, modifiers=%v", p0.TempModifiers)
		}
	})

	t.Run("pigeon raid order buffs a rush spell learned this turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		order := NewCardInstance(baseCard(t, "2321110"), 0, 1)
		rushThisTurn := readySkill(baseCard(t, "3321101"), 0)
		rushThisTurn.EnterTurn = engine.State.TurnNumber
		oldRush := readySkill(baseCard(t, "3321101"), 0)
		oldRush.EnterTurn = engine.State.TurnNumber - 1
		nonRushThisTurn := readySkill(baseCard(t, "3521106"), 0)
		nonRushThisTurn.EnterTurn = engine.State.TurnNumber
		p0.Skills[0] = rushThisTurn
		p0.Skills[1] = oldRush
		p0.Skills[2] = nonRushThisTurn

		if err := (Card2321110PigeonRaidOrder{}).OnUseItem(&EffectContext{Engine: engine, Source: order, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("pigeon raid order: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "pigeon_raid_order_skill" || len(engine.State.PendingAction.Candidates) != 1 {
			t.Fatalf("2321110 should only offer rush spells learned this turn, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, rushThisTurn.InstanceID)
		if len(p0.TempModifiers) != 2 {
			t.Fatalf("2321110 should create one power and one attack modifier, modifiers=%v", p0.TempModifiers)
		}
		if p0.TempModifiers[0].Type != TempModSkillPowerBonus || p0.TempModifiers[1].Type != TempModNextSkillUseAttackBonus {
			t.Fatalf("2321110 should use next-use power and attack modifiers, modifiers=%v", p0.TempModifiers)
		}
		if got := engine.effectiveSpellPower(0, rushThisTurn, nil); got != rushThisTurn.Card.Power+1 {
			t.Fatalf("2321110 should buff selected skill power, got=%d", got)
		}
		if got := engine.effectiveSpellDamage(0, rushThisTurn, rushThisTurn.Card.Attack, nil); got != rushThisTurn.Card.Attack+1 {
			t.Fatalf("2321110 should buff selected skill attack, got=%d", got)
		}
		if got := engine.effectiveSpellPower(0, oldRush, nil); got != oldRush.Card.Power {
			t.Fatalf("2321110 should not buff unselected rush skill, got=%d", got)
		}
		engine.consumeNextSpellPowerBonuses(p0, rushThisTurn)
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextSkillUseAttackBonus {
			t.Fatalf("2321110 power consumption should leave attack bonus for damage timing, modifiers=%v", p0.TempModifiers)
		}
		engine.consumeNextSpellAttackBonuses(p0, rushThisTurn)
		if len(p0.TempModifiers) != 0 {
			t.Fatalf("2321110 next-use attack should be consumed at damage timing, modifiers=%v", p0.TempModifiers)
		}
	})

	t.Run("uncontrolled divine fire beast buffs both players attacking spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		attackerSkill := readySkill(baseCard(t, "3121001"), 0)
		defenseSkill := readySkill(baseCard(t, "3121102"), 0)

		placeUnit(baseCard(t, "1121107"), 0, 0, 0, engine)
		if got := engine.effectiveSpellPower(0, attackerSkill, nil); got != attackerSkill.Card.Power+2 {
			t.Fatalf("1121107 should buff friendly attacking spells, got=%d", got)
		}
		if got := engine.effectiveSkillPowerForPurposeWithData(0, defenseSkill, nil, skillPurposeDefend, nil); got != defenseSkill.Card.Power {
			t.Fatalf("1121107 should not buff defensive spells, got=%d", got)
		}

		enemyEngine := setupReportedBugEngine(t)
		enemyAttackerSkill := readySkill(baseCard(t, "3121001"), 0)
		placeUnit(baseCard(t, "1121107"), 1, 0, 0, enemyEngine)
		if got := enemyEngine.effectiveSpellPower(0, enemyAttackerSkill, nil); got != enemyAttackerSkill.Card.Power+2 {
			t.Fatalf("1121107 should buff enemy attacking spells too, got=%d", got)
		}
	})

	t.Run("killing wind gains power from absolute hand size difference", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		killingWind := readySkill(baseCard(t, "3321102"), 0)

		p0.Hand = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
		p1.Hand = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 1, 1),
			NewCardInstance(baseCard(t, "1021002"), 1, 2),
			NewCardInstance(baseCard(t, "1021004"), 1, 3),
			NewCardInstance(baseCard(t, "1021005"), 1, 4),
		}
		if got := engine.effectiveSpellPower(0, killingWind, nil); got != killingWind.Card.Power+3 {
			t.Fatalf("3321102 should gain power from opponent's larger hand difference, got=%d", got)
		}

		p0.Hand = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1021002"), 0, 2),
			NewCardInstance(baseCard(t, "1021004"), 0, 3),
			NewCardInstance(baseCard(t, "1021005"), 0, 4),
			NewCardInstance(baseCard(t, "1021006"), 0, 5),
		}
		p1.Hand = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 1, 1),
			NewCardInstance(baseCard(t, "1021002"), 1, 2),
		}
		if got := engine.effectiveSpellPower(0, killingWind, nil); got != killingWind.Card.Power+3 {
			t.Fatalf("3321102 should gain power from caster's larger hand difference, got=%d", got)
		}

		p0.Hand = p0.Hand[:2]
		if got := engine.effectiveSpellPower(0, killingWind, nil); got != killingWind.Card.Power {
			t.Fatalf("3321102 should not gain power when hand sizes are equal, got=%d", got)
		}
	})

	t.Run("divine help gains power only when boosting mystery spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		divineHelp := readySkill(baseCard(t, "3521102"), 0)
		mysterySpell := readySkill(baseCard(t, "3521106"), 0)
		nonMysterySpell := readySkill(baseCard(t, "3321101"), 0)

		if got := engine.effectiveSkillPowerForPurposeWithData(0, divineHelp, mysterySpell, skillPurposeAttackBoost, nil); got != divineHelp.Card.Power+2 {
			t.Fatalf("3521102 should gain power when boosting a mystery spell, got=%d", got)
		}
		if got := engine.effectiveSkillPowerForPurposeWithData(0, divineHelp, nonMysterySpell, skillPurposeAttackBoost, nil); got != divineHelp.Card.Power {
			t.Fatalf("3521102 should not gain power when boosting non-mystery spells, got=%d", got)
		}
		if got := engine.effectiveSkillPowerForPurposeWithData(0, divineHelp, divineHelp, skillPurposeAttack, nil); got != divineHelp.Card.Power {
			t.Fatalf("3521102 should not gain power when used as the main attack spell, got=%d", got)
		}
	})

	t.Run("returning heart counts only friendly light and non-light companions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		returningHeart := readySkill(baseCard(t, "3521106"), 0)

		placeUnit(baseCard(t, "1521001"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1521103"), 0, 1, 0, engine)
		placeUnit(baseCard(t, "1021001"), 0, 2, 0, engine)
		placeUnit(baseCard(t, "1421001"), 1, 0, 0, engine)
		if got := engine.effectiveSpellPower(0, returningHeart, nil); got != returningHeart.Card.Power+1 {
			t.Fatalf("3521106 should count friendly light companions minus friendly non-light companions only, got=%d", got)
		}

		placeUnit(baseCard(t, "1421001"), 0, 1, 1, engine)
		if got := engine.effectiveSpellPower(0, returningHeart, nil); got != returningHeart.Card.Power {
			t.Fatalf("3521106 should reduce power for each friendly non-light companion, got=%d", got)
		}
	})

	t.Run("intimidation gains power and attack from weakened enemy spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p1 := engine.State.Players[1]
		intimidation := readySkill(baseCard(t, "3621105"), 0)
		enemyA := readySkill(baseCard(t, "3121001"), 1)
		enemyB := readySkill(baseCard(t, "3121002"), 1)
		enemyC := readySkill(baseCard(t, "3221003"), 1)
		p1.Skills[0] = enemyA
		p1.Skills[1] = enemyB
		p1.Skills[2] = enemyC

		if got := engine.effectiveSpellPower(0, intimidation, nil); got != intimidation.Card.Power {
			t.Fatalf("3621105 should not gain power without weakened enemy spells, got=%d", got)
		}
		enemyA.Statuses[StatusWeaken] = 1
		if got := engine.effectiveSpellPower(0, intimidation, nil); got != intimidation.Card.Power+1 {
			t.Fatalf("3621105 should gain +1 power from one weakened enemy spell, got=%d", got)
		}
		if got := engine.effectiveSpellDamage(0, intimidation, intimidation.Card.Attack, nil); got != intimidation.Card.Attack+1 {
			t.Fatalf("3621105 should gain +1 attack from one weakened enemy spell, got=%d", got)
		}

		enemyB.Statuses[StatusWeaken] = 2
		enemyC.Statuses[StatusWeaken] = 1
		if got := engine.effectiveSpellPower(0, intimidation, nil); got != intimidation.Card.Power+2 {
			t.Fatalf("3621105 should cap power bonus at 2, got=%d", got)
		}
		if got := engine.effectiveSpellDamage(0, intimidation, intimidation.Card.Attack, nil); got != intimidation.Card.Attack+2 {
			t.Fatalf("3621105 should cap attack bonus at 2, got=%d", got)
		}
	})

	t.Run("crushing stone gains power against high life targets", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		crushingStone := readySkill(baseCard(t, "3421110"), 0)
		lowLifeTarget := placeUnit(baseCard(t, "1021002"), 1, 0, 0, engine)
		highLifeTarget := placeUnit(baseCard(t, "1121104"), 1, 1, 0, engine)

		if got := engine.effectiveSpellPower(0, crushingStone, nil); got != crushingStone.Card.Power {
			t.Fatalf("3421110 should not gain power without a target, got=%d", got)
		}
		if got := engine.effectiveSpellPower(0, crushingStone, nil, SpellTarget{Type: "unit", Position: *lowLifeTarget.Position}); got != crushingStone.Card.Power {
			t.Fatalf("3421110 should not gain power against targets with life 2, got=%d", got)
		}
		highLifeTarget.CurrentLife = 3
		if got := engine.effectiveSpellPower(0, crushingStone, nil, SpellTarget{Type: "unit", Position: *highLifeTarget.Position}); got != crushingStone.Card.Power+1 {
			t.Fatalf("3421110 should gain power against targets with life above 2, got=%d", got)
		}
	})

	t.Run("giant rock collapse cannot be used as a boost but can be boosted", func(t *testing.T) {
		setupReportedBugEngine(t)
		collapse := readySkill(baseCard(t, "3421103"), 0)

		if canUseSkillForPurpose(collapse.Card, skillPurposeAttackBoost) {
			t.Fatalf("3421103 should not be usable as an attack boost")
		}
		if canUseSkillForPurpose(collapse.Card, skillPurposeDefenseBoost) {
			t.Fatalf("3421103 should not be usable as a defense boost")
		}
		if !canUseSkillForPurpose(collapse.Card, skillPurposeAttack) {
			t.Fatalf("3421103 should still be usable as a main attack spell")
		}
		if !canSkillBeBoosted(collapse) {
			t.Fatalf("3421103 should still be boostable as a main spell")
		}

		info := CardRuleInfo(collapse.Card)
		if info["can_attack_boost"] != false || info["can_defense_boost"] != false || info["can_be_boosted"] != true {
			t.Fatalf("3421103 should expose no-boost but boostable rule info, info=%v", info)
		}
	})

	t.Run("guarding stone array reduces its main defense use cost only", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		stoneArray := readySkill(baseCard(t, "3421108"), 0)
		otherDefense := readySkill(baseCard(t, "3221103"), 0)
		p0.Skills[0] = stoneArray
		p0.Skills[1] = otherDefense

		if got := engine.effectiveSkillUseCostForPurpose(p0, stoneArray, skillPurposeAttack)[model.ElementEarth]; got != 3 {
			t.Fatalf("3421108 should not reduce attack use cost, cost=%v", engine.effectiveSkillUseCostForPurpose(p0, stoneArray, skillPurposeAttack))
		}
		if got := engine.effectiveSkillUseCostForPurpose(p0, stoneArray, skillPurposeDefend)[model.ElementEarth]; got != 2 {
			t.Fatalf("3421108 should reduce its main defense use cost by 1 earth, cost=%v", engine.effectiveSkillUseCostForPurpose(p0, stoneArray, skillPurposeDefend))
		}
		if got := engine.effectiveSkillUseCostForPurpose(p0, stoneArray, skillPurposeDefenseBoost)[model.ElementEarth]; got != 3 {
			t.Fatalf("3421108 should not reduce defense boost cost, cost=%v", engine.effectiveSkillUseCostForPurpose(p0, stoneArray, skillPurposeDefenseBoost))
		}
		if got := engine.effectiveSkillUseCostForPurpose(p0, otherDefense, skillPurposeDefend)[model.ElementWater]; got != 2 {
			t.Fatalf("3421108 should not reduce another defense skill cost, cost=%v", engine.effectiveSkillUseCostForPurpose(p0, otherDefense, skillPurposeDefend))
		}
	})

	t.Run("magma fortress chariot burns the target when attacking", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		chariot := placeUnit(baseCard(t, "1121104"), 0, 1, 0, engine)
		target := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
		beforeLife := target.CurrentLife

		if err := engine.HandleAction(0, ActionMessage{Action: "attack", Data: map[string]any{
			"attacker_id": chariot.InstanceID,
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("chariot attack: %v", err)
		}
		if target.Statuses[StatusBurn] != 1 {
			t.Fatalf("1121104 should burn its attack target, statuses=%v", target.Statuses)
		}
		if target.CurrentLife != beforeLife-chariot.CurrentAttack {
			t.Fatalf("chariot attack should still deal normal attack damage, life=%d want=%d", target.CurrentLife, beforeLife-chariot.CurrentAttack)
		}
	})

	t.Run("psychic disk reduces medium skill use cost", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		disk := NewCardInstance(baseCard(t, "2021108"), 0, 1)
		p0.Equipment[0] = disk

		medium := readySkill(baseCard(t, "3621105"), 0)
		nonMedium := readySkill(baseCard(t, "3121105"), 0)
		arcaneMedium := readySkill(baseCard(t, "3021002"), 0)
		if got := engine.effectiveSkillUseCost(p0, medium)[model.ElementShadow]; got != 1 {
			t.Fatalf("2021108 should reduce dark medium skill cost from 2 to 1, cost=%v", engine.effectiveSkillUseCost(p0, medium))
		}
		if got := engine.effectiveSkillUseCost(p0, nonMedium)[model.ElementFire]; got != 2 {
			t.Fatalf("2021108 should not reduce non-medium skill cost, cost=%v", engine.effectiveSkillUseCost(p0, nonMedium))
		}
		if got := engine.effectiveSkillUseCost(p0, arcaneMedium)[model.ElementArcane]; got != 0 {
			t.Fatalf("2021108 should reduce one-cost arcane medium skill to 0, cost=%v", engine.effectiveSkillUseCost(p0, arcaneMedium))
		}
	})

	t.Run("spirit guard amulet gains arcane load only as sole equipment", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		amulet := NewCardInstance(baseCard(t, "2021110"), 0, 1)
		p0.Equipment[0] = amulet

		if got := engine.effectiveElementsGain(amulet)[model.ElementArcane]; got != amulet.Card.ElementsGain[model.ElementArcane]+1 {
			t.Fatalf("2021110 should gain +1 arcane load as sole equipment, load=%v", engine.effectiveElementsGain(amulet))
		}

		p0.Equipment[1] = NewCardInstance(baseCard(t, "2021108"), 0, 1)
		if got := engine.effectiveElementsGain(amulet)[model.ElementArcane]; got != amulet.Card.ElementsGain[model.ElementArcane] {
			t.Fatalf("2021110 should lose bonus with another equipment, load=%v", engine.effectiveElementsGain(amulet))
		}

		p0.Equipment[1] = nil
		amulet.Statuses[StatusPetrify] = 1
		if got := engine.effectiveElementsGain(amulet)[model.ElementArcane]; got != amulet.Card.ElementsGain[model.ElementArcane] {
			t.Fatalf("2021110 aura should be inactive while petrified, load=%v statuses=%v", engine.effectiveElementsGain(amulet), amulet.Statuses)
		}
	})

	t.Run("raider ghost captain gives other friendly raiders water load", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		captain := NewCardInstance(baseCard(t, "1221110"), 0, 1)
		friendlyRaider := NewCardInstance(baseCard(t, "1221111"), 0, 1)
		friendlyNonRaider := NewCardInstance(baseCard(t, "1221106"), 0, 1)
		enemyRaider := NewCardInstance(baseCard(t, "1221111"), 1, 1)
		p0.Units[1][0] = captain
		p0.Units[1][1] = friendlyRaider
		p0.Units[1][2] = friendlyNonRaider
		p1.Units[1][1] = enemyRaider

		if got := engine.effectiveElementsGain(friendlyRaider)[model.ElementWater]; got != friendlyRaider.Card.ElementsGain[model.ElementWater]+1 {
			t.Fatalf("1221110 should add water load to other friendly raiders, load=%v", engine.effectiveElementsGain(friendlyRaider))
		}
		if got := engine.effectiveElementsGain(captain)[model.ElementWater]; got != captain.Card.ElementsGain[model.ElementWater] {
			t.Fatalf("1221110 should not buff itself, load=%v", engine.effectiveElementsGain(captain))
		}
		if got := engine.effectiveElementsGain(friendlyNonRaider)[model.ElementWater]; got != friendlyNonRaider.Card.ElementsGain[model.ElementWater] {
			t.Fatalf("1221110 should not buff non-raider companions, load=%v", engine.effectiveElementsGain(friendlyNonRaider))
		}
		if got := engine.effectiveElementsGain(enemyRaider)[model.ElementWater]; got != enemyRaider.Card.ElementsGain[model.ElementWater] {
			t.Fatalf("1221110 should not buff enemy raiders, load=%v", engine.effectiveElementsGain(enemyRaider))
		}

		captain.Statuses[StatusPetrify] = 1
		if got := engine.effectiveElementsGain(friendlyRaider)[model.ElementWater]; got != friendlyRaider.Card.ElementsGain[model.ElementWater] {
			t.Fatalf("1221110 aura should be inactive while petrified, load=%v statuses=%v", engine.effectiveElementsGain(friendlyRaider), captain.Statuses)
		}
	})

	t.Run("seven gods blessing rewards distinct skill elements", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		blessing := readySkill(baseCard(t, "3021104"), 0)
		fireSpell := readySkill(baseCard(t, "3121001"), 0)
		waterSorcery := readySkill(baseCard(t, "3221007"), 0)
		p0.Skills[0] = blessing
		p0.Skills[1] = fireSpell
		p0.Skills[2] = waterSorcery

		if got := engine.effectiveSkillUseCost(p0, fireSpell)[model.ElementFire]; got != max(fireSpell.Card.ElementsExpense[model.ElementFire]-1, 0) {
			t.Fatalf("3021104 should reduce distinct fire skill use cost, cost=%v", engine.effectiveSkillUseCost(p0, fireSpell))
		}
		if got := engine.effectiveSkillUseCost(p0, waterSorcery)[model.ElementWater]; got != max(waterSorcery.Card.ElementsExpense[model.ElementWater]-1, 0) {
			t.Fatalf("3021104 should reduce distinct water skill use cost, cost=%v", engine.effectiveSkillUseCost(p0, waterSorcery))
		}
		if got := engine.effectiveSpellPower(0, fireSpell, nil); got != fireSpell.Card.Power+2 {
			t.Fatalf("3021104 should give spell skills +2 power when elements are distinct, got=%d", got)
		}

		duplicateFire := readySkill(baseCard(t, "3121002"), 0)
		p0.Skills[3] = duplicateFire
		if got := engine.effectiveSkillUseCost(p0, fireSpell)[model.ElementFire]; got != fireSpell.Card.ElementsExpense[model.ElementFire] {
			t.Fatalf("3021104 should not reduce costs when skill elements repeat, cost=%v", engine.effectiveSkillUseCost(p0, fireSpell))
		}
		if got := engine.effectiveSpellPower(0, fireSpell, nil); got != fireSpell.Card.Power {
			t.Fatalf("3021104 should not add power when skill elements repeat, got=%d", got)
		}
	})

	t.Run("arcane flow doubles only itself while friendly field is all arcane", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		hero := NewCardInstance(baseCard(t, "4011102"), 0, 1)
		arcaneFlow := readySkill(baseCard(t, "3021106"), 0)
		otherSpell := readySkill(baseCard(t, "3021008"), 0)
		p0.Hero = hero
		p0.Units[1][1] = hero
		p0.Skills[0] = arcaneFlow
		p0.Skills[1] = otherSpell

		if got := engine.effectiveSpellPower(0, arcaneFlow, nil); got != arcaneFlow.Card.Power*2 {
			t.Fatalf("3021106 should double its own power on all-arcane field, got=%d", got)
		}
		if got := engine.effectiveSpellPower(0, otherSpell, nil); got != otherSpell.Card.Power {
			t.Fatalf("3021106 should not double other spells, got=%d", got)
		}

		p0.Units[1][0] = NewCardInstance(baseCard(t, "1121001"), 0, 1)
		if got := engine.effectiveSpellPower(0, arcaneFlow, nil); got != arcaneFlow.Card.Power {
			t.Fatalf("3021106 should not double with a non-arcane friendly field card, got=%d", got)
		}

		p0.Units[1][0] = nil
		arcaneFlow.Statuses[StatusPetrify] = 1
		if got := engine.effectiveSpellPower(0, arcaneFlow, nil); got != arcaneFlow.Card.Power {
			t.Fatalf("3021106 aura should be inactive while petrified, got=%d statuses=%v", got, arcaneFlow.Statuses)
		}
	})

	t.Run("moonlit spirit boosts the next spell attack then loses its aura", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		spirit := placeUnit(baseCard(t, "1521101"), 0, 0, 0, engine)
		firstSpell := readySkill(baseCard(t, "3021008"), 0)
		secondSpell := readySkill(baseCard(t, "3021008"), 0)
		p0.Skills[0] = firstSpell
		p0.Skills[1] = secondSpell
		p0.Elements[model.ElementArcane] = 2
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if got := engine.effectiveSpellPower(0, firstSpell, nil); got != firstSpell.Card.Power+2 {
			t.Fatalf("1521101 should boost friendly spells before its aura is spent, got=%d", got)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": firstSpell.InstanceID,
			"target_type": "unit",
			"target_col":  float64(target.Position.Col),
			"target_row":  float64(target.Position.Row),
		}}); err != nil {
			t.Fatalf("cast boosted spell with moonlit spirit: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != firstSpell.Card.Power+2 {
			t.Fatalf("1521101 should boost the spell attack that spends it, pending=%+v", engine.State.PendingSpell)
		}
		if spirit.Statuses[moonlitSpiritAuraSpentStatus] != 1 {
			t.Fatalf("1521101 should lose aura after friendly spell attack, statuses=%v", spirit.Statuses)
		}
		if got := engine.effectiveSpellPower(0, secondSpell, nil); got != secondSpell.Card.Power {
			t.Fatalf("1521101 should not boost later spells after aura is spent, got=%d", got)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve moonlit boosted spell: %v", err)
		}

		enemyEngine := setupReportedBugEngine(t)
		enemyP1 := enemyEngine.State.Players[1]
		enemySpirit := placeUnit(baseCard(t, "1521101"), 0, 0, 0, enemyEngine)
		enemySpell := readySkill(baseCard(t, "3021008"), 1)
		enemyP1.Skills[0] = enemySpell
		enemyEngine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, enemySpell, map[string]any{"cast_player": 1, "attacker": 1})
		if enemySpirit.Statuses[moonlitSpiritAuraSpentStatus] != 0 {
			t.Fatalf("1521101 should not lose aura from enemy spell attacks, statuses=%v", enemySpirit.Statuses)
		}
	})

	t.Run("ripple slash improves after one copy was cast this turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		first := readySkill(baseCard(t, "3221109"), 0)
		second := readySkill(baseCard(t, "3221109"), 0)
		p0.Skills[0] = first
		p0.Skills[1] = second
		p0.Elements[model.ElementWater] = 2
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if got := engine.effectiveSpellPower(0, second, nil); got != second.Card.Power {
			t.Fatalf("3221109 should start at printed power before any copy is cast, got=%d", got)
		}
		if got := engine.effectiveSpellArea(second); got != SpellAreaSingle {
			t.Fatalf("3221109 should start as single target, got=%s", got)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": first.InstanceID,
			"target_type": "unit",
			"target_col":  float64(target.Position.Col),
			"target_row":  float64(target.Position.Row),
		}}); err != nil {
			t.Fatalf("cast first ripple slash: %v", err)
		}
		if spellCastByNumberThisTurn(p0, "3221109") != 1 {
			t.Fatalf("first 3221109 cast should be recorded by number, casts=%v", p0.SpellsCastByNumberThisTurn)
		}
		if got := engine.effectiveSpellPower(0, second, nil); got != second.Card.Power+2 {
			t.Fatalf("3221109 should gain +2 power after a copy was cast, got=%d", got)
		}
		if got := engine.effectiveSpellArea(second); got != SpellAreaFrontRow {
			t.Fatalf("3221109 should become front-row area after a copy was cast, got=%s", got)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve first ripple slash: %v", err)
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end p0 turn: %v", err)
		}
		if spellCastByNumberThisTurn(p0, "3221109") != 0 {
			t.Fatalf("3221109 cast count should reset at turn end, casts=%v", p0.SpellsCastByNumberThisTurn)
		}
		if got := engine.effectiveSpellPower(0, second, nil); got != second.Card.Power {
			t.Fatalf("3221109 bonus should not persist after turn end, got=%d", got)
		}
	})

	t.Run("silverleaf cyclone power becomes six after a card enters graveyard this turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		cyclone := readySkill(baseCard(t, "3321109"), 0)
		discard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Skills[0] = cyclone
		p0.Hand = []*CardInstance{discard}

		if got := engine.effectiveSpellPower(0, cyclone, nil); got != cyclone.Card.Power {
			t.Fatalf("3321109 should start at printed power before any card enters graveyard, got=%d", got)
		}
		if engine.discardHandCardAt(0, 0) != discard {
			t.Fatalf("discard setup card")
		}
		if !engine.State.CardEnteredGraveyardThisTurn {
			t.Fatalf("discarding a hand card should mark graveyard entry this turn")
		}
		if got := engine.effectiveSpellPower(0, cyclone, nil); got != 6 {
			t.Fatalf("3321109 power should become 6 after hand discard, got=%d", got)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end turn after discard: %v", err)
		}
		if engine.State.CardEnteredGraveyardThisTurn {
			t.Fatalf("graveyard entry marker should clear at turn end")
		}
		if got := engine.effectiveSpellPower(0, cyclone, nil); got != cyclone.Card.Power {
			t.Fatalf("3321109 power should reset after turn end, got=%d", got)
		}

		deathEngine := setupReportedBugEngine(t)
		deathCyclone := readySkill(baseCard(t, "3321109"), 0)
		deathEngine.State.Players[0].Skills[0] = deathCyclone
		unit := placeUnit(baseCard(t, "1021001"), 1, 1, 0, deathEngine)
		deathEngine.destroyUnit(unit, 1)
		if got := deathEngine.effectiveSpellPower(0, deathCyclone, nil); got != 6 {
			t.Fatalf("3321109 power should become 6 after unit death, got=%d", got)
		}
	})

	t.Run("rock wall monster limits damage while its owner has no learned spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		monster := placeUnit(baseCard(t, "1421111"), 0, 0, 0, engine)
		before := monster.CurrentLife

		engine.dealDamageWithExtra(monster, 3, 0, map[string]any{"damage_source": "test"})
		if monster.CurrentLife != before-1 {
			t.Fatalf("1421111 should take at most 1 damage with no learned spells, life=%d want=%d", monster.CurrentLife, before-1)
		}

		sorceryEngine := setupReportedBugEngine(t)
		sorceryP0 := sorceryEngine.State.Players[0]
		sorceryP0.Skills[0] = readySkill(baseCard(t, "3221007"), 0)
		sorceryMonster := placeUnit(baseCard(t, "1421111"), 0, 0, 0, sorceryEngine)
		sorceryBefore := sorceryMonster.CurrentLife
		sorceryEngine.dealDamageWithExtra(sorceryMonster, 3, 0, map[string]any{"damage_source": "test"})
		if sorceryMonster.CurrentLife != sorceryBefore-1 {
			t.Fatalf("1421111 should still limit damage when owner learned only sorceries, life=%d want=%d", sorceryMonster.CurrentLife, sorceryBefore-1)
		}

		spellEngine := setupReportedBugEngine(t)
		spellP0 := spellEngine.State.Players[0]
		spellP0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		spellMonster := placeUnit(baseCard(t, "1421111"), 0, 0, 0, spellEngine)
		spellBefore := spellMonster.CurrentLife
		spellEngine.dealDamageWithExtra(spellMonster, 3, 0, map[string]any{"damage_source": "test"})
		if spellMonster.CurrentLife != spellBefore-3 {
			t.Fatalf("1421111 should take full damage after owner learns a spell, life=%d want=%d", spellMonster.CurrentLife, spellBefore-3)
		}
	})

	t.Run("rock wall colossus gives summoned earth companions life before learned spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1421110"), 0, 0, 0, engine)
		p0.Elements[model.ElementEarth] = 10

		earthCompanion := NewCardInstance(baseCard(t, "1421001"), 0, 1)
		p0.Hand = []*CardInstance{earthCompanion}
		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": earthCompanion.InstanceID,
			"col":         float64(1),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("summon earth companion: %v", err)
		}
		if maxLife(earthCompanion) != earthCompanion.Card.Life+1 || earthCompanion.CurrentLife != earthCompanion.Card.Life+1 {
			t.Fatalf("1421110 should buff later earth companions, max=%d current=%d", maxLife(earthCompanion), earthCompanion.CurrentLife)
		}

		spellEngine := setupReportedBugEngine(t)
		spellP0 := spellEngine.State.Players[0]
		placeUnit(baseCard(t, "1421110"), 0, 0, 0, spellEngine)
		spellP0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		blockedEarth := NewCardInstance(baseCard(t, "1421001"), 0, 1)
		spellP0.Hand = []*CardInstance{blockedEarth}
		spellP0.Elements[model.ElementEarth] = 10
		if err := spellEngine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": blockedEarth.InstanceID,
			"col":         float64(1),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("summon earth companion after learned spell: %v", err)
		}
		if maxLife(blockedEarth) != blockedEarth.Card.Life || blockedEarth.CurrentLife != blockedEarth.Card.Life {
			t.Fatalf("1421110 should not buff earth companions after owner learns a spell, max=%d current=%d", maxLife(blockedEarth), blockedEarth.CurrentLife)
		}
	})

	t.Run("church envoy removes negative statuses from friendly cards", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		envoy := NewCardInstance(baseCard(t, "1021109"), 0, 1)
		target := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		target.Statuses[StatusBurn] = 2
		target.Statuses[StatusFreeze] = 1
		target.Statuses["mastery"] = 1
		p0.Units[0][0] = target

		if err := (Card1021109ChurchEnvoy{}).OnUltimate(&EffectContext{Engine: engine, Source: envoy, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("church envoy ultimate: %v", err)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if target.Statuses[StatusBurn] != 0 || target.Statuses[StatusFreeze] != 0 || target.Statuses["mastery"] != 1 {
			t.Fatalf("church envoy should clear only negative statuses, statuses=%v", target.Statuses)
		}
	})

	t.Run("shadow heroes add blood feast or mill both decks", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		hubert := NewCardInstance(baseCard(t, "4611101"), 0, 1)
		if err := (Card4611101BloodCountHubert{}).OnEnter(&EffectContext{Engine: engine, Source: hubert, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("hubert enter: %v", err)
		}
		if len(p0.SkillPool) != 1 || p0.SkillPool[0].Card.Number != "3601101" {
			t.Fatalf("hubert should add blood feast to skill pool, pool=%v", cardsToInfo(p0.SkillPool))
		}

		millEngine := setupReportedBugEngine(t)
		millP0 := millEngine.State.Players[0]
		millP1 := millEngine.State.Players[1]
		for i := 0; i < 5; i++ {
			millP0.Deck = append(millP0.Deck, NewCardInstance(baseCard(t, "1021001"), 0, 1))
			millP1.Deck = append(millP1.Deck, NewCardInstance(baseCard(t, "1021002"), 1, 1))
		}
		firstP0 := millP0.Deck[0]
		firstP1 := millP1.Deck[0]
		dom := NewCardInstance(baseCard(t, "4611102"), 0, 1)
		if err := (Card4611102CalamityRoseDom{}).OnEnter(&EffectContext{Engine: millEngine, Source: dom, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("dom enter: %v", err)
		}
		if len(millP0.Deck) != 1 || len(millP1.Deck) != 1 || len(millP0.Graveyard) != 4 || len(millP1.Graveyard) != 4 {
			t.Fatalf("dom should mill top four from both decks, p0 deck/grave=%d/%d p1 deck/grave=%d/%d", len(millP0.Deck), len(millP0.Graveyard), len(millP1.Deck), len(millP1.Graveyard))
		}
		if millP0.Graveyard[0] != firstP0 || millP1.Graveyard[0] != firstP1 {
			t.Fatalf("dom should preserve top-deck mill order, p0 grave=%v p1 grave=%v", cardsToInfo(millP0.Graveyard), cardsToInfo(millP1.Graveyard))
		}
	})
}

func TestRoyalConflictResetAndTemporaryAbilityEffects(t *testing.T) {
	t.Run("fire butterfly temporarily changes its load to one air", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		butterfly := NewCardInstance(baseCard(t, "1121108"), 0, 1)
		if err := (Card1121108FireButterfly{}).OnPerTurn(&EffectContext{Engine: engine, Source: butterfly, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("fire butterfly ability: %v", err)
		}
		load := effectiveElementsGain(butterfly)
		if len(load) != 1 || load[model.ElementAir] != 1 {
			t.Fatalf("fire butterfly should temporarily set load to exactly 1 air, load=%v", load)
		}
		if err := (Card1121108FireButterfly{}).OnTurnEnd(&EffectContext{Engine: engine, Source: butterfly, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("fire butterfly turn end: %v", err)
		}
		if effectiveElementsGain(butterfly)[model.ElementAir] != butterfly.Card.ElementsGain[model.ElementAir] || butterfly.Statuses[fireButterflyTemporaryLoadStatus] != 0 {
			t.Fatalf("fire butterfly temporary load should expire, load=%v statuses=%v", effectiveElementsGain(butterfly), butterfly.Statuses)
		}

		preset := NewCardInstance(baseCard(t, "1121108"), 0, 1)
		setElementsGain(preset, map[string]int{model.ElementFire: 2})
		if err := (Card1121108FireButterfly{}).OnPerTurn(&EffectContext{Engine: engine, Source: preset, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("fire butterfly preset ability: %v", err)
		}
		if err := (Card1121108FireButterfly{}).OnTurnEnd(&EffectContext{Engine: engine, Source: preset, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("fire butterfly preset turn end: %v", err)
		}
		presetLoad := effectiveElementsGain(preset)
		if presetLoad[model.ElementFire] != 2 || presetLoad[model.ElementAir] != 0 {
			t.Fatalf("fire butterfly should restore an earlier load override, load=%v", presetLoad)
		}

		overridden := NewCardInstance(baseCard(t, "1121108"), 0, 1)
		if err := (Card1121108FireButterfly{}).OnPerTurn(&EffectContext{Engine: engine, Source: overridden, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("fire butterfly second ability: %v", err)
		}
		setElementsGain(overridden, map[string]int{model.ElementFire: 2})
		if err := (Card1121108FireButterfly{}).OnTurnEnd(&EffectContext{Engine: engine, Source: overridden, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("fire butterfly overridden turn end: %v", err)
		}
		overriddenLoad := effectiveElementsGain(overridden)
		if overriddenLoad[model.ElementFire] != 2 || overriddenLoad[model.ElementAir] != 0 {
			t.Fatalf("fire butterfly should not overwrite a later load override, load=%v", overriddenLoad)
		}

		sameValueOverride := NewCardInstance(baseCard(t, "1121108"), 0, 1)
		setElementsGain(sameValueOverride, map[string]int{model.ElementFire: 2})
		if err := (Card1121108FireButterfly{}).OnPerTurn(&EffectContext{Engine: engine, Source: sameValueOverride, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("fire butterfly same-value ability: %v", err)
		}
		setElementsGain(sameValueOverride, map[string]int{model.ElementAir: 1})
		if err := (Card1121108FireButterfly{}).OnTurnEnd(&EffectContext{Engine: engine, Source: sameValueOverride, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("fire butterfly same-value turn end: %v", err)
		}
		sameValueLoad := effectiveElementsGain(sameValueOverride)
		if sameValueLoad[model.ElementAir] != 1 || sameValueLoad[model.ElementFire] != 0 {
			t.Fatalf("fire butterfly should not restore over a later equal-value load override, load=%v", sameValueLoad)
		}
	})

	t.Run("water mage resets a low-cost water spell", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		mage := NewCardInstance(baseCard(t, "1221112"), 0, 1)
		waterSkill := readySkill(baseCard(t, "3221103"), 0)
		waterSkill.IsHorizontal = true
		p0.Skills[0] = waterSkill
		if err := (Card1221112WaterMage{}).OnUltimate(&EffectContext{Engine: engine, Source: mage, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("water mage ultimate: %v", err)
		}
		resolvePendingSelection(t, engine, 0, waterSkill.InstanceID)
		if waterSkill.IsHorizontal {
			t.Fatalf("water mage should reset selected water skill")
		}

		failEngine := setupReportedBugEngine(t)
		failMage := placeUnit(baseCard(t, "1221112"), 0, 0, 0, failEngine)
		if err := (Card1221112WaterMage{}).OnUltimate(&EffectContext{Engine: failEngine, Source: failMage, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("water mage should no-op without a resettable water spell: %v", err)
		}
		if failEngine.State.PendingAction != nil {
			t.Fatalf("water mage should not open a prompt without legal targets, pending=%+v", failEngine.State.PendingAction)
		}
		err := failEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  failMage.InstanceID,
			"ability_type": "ultimate",
		}})
		if err == nil || failMage.UltimateUsed {
			t.Fatalf("water mage action should fail without burning ultimate, err=%v ultimateUsed=%v", err, failMage.UltimateUsed)
		}

		boundEngine := setupReportedBugEngine(t)
		boundMage := placeUnit(baseCard(t, "1221112"), 0, 0, 0, boundEngine)
		host := placeUnit(baseCard(t, "1021001"), 0, 1, 0, boundEngine)
		boundWaterSkill := readySkill(baseCard(t, "3221103"), 0)
		boundWaterSkill.IsHorizontal = true
		boundWaterSkill.SlotIndex = -1
		host.BoundSkills = []*CardInstance{boundWaterSkill}
		if err := boundEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  boundMage.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("water mage should accept bound water spell target: %v", err)
		}
		resolvePendingSelection(t, boundEngine, 0, boundWaterSkill.InstanceID)
		if boundWaterSkill.IsHorizontal {
			t.Fatalf("water mage should reset selected bound water skill")
		}
	})

	t.Run("winterfell anti mage discounts each learned skill's next use by one water", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		mage := placeUnit(baseCard(t, "1221115"), 0, 0, 0, engine)
		first := readySkill(baseCard(t, "3221106"), 0)
		second := readySkill(baseCard(t, "3221107"), 0)
		p0.Skills[0] = first
		p0.Skills[1] = second
		p0.Elements = map[string]int{model.ElementWater: 10, model.ElementAir: 10}

		if err := (Card1221115WinterfellAntiMage{}).OnPrayer(&EffectContext{Engine: engine, Source: mage, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1221115 prayer: %v", err)
		}
		if len(p0.TempModifiers) != 2 {
			t.Fatalf("1221115 should add one next-use cost modifier per learned skill, modifiers=%+v", p0.TempModifiers)
		}
		if got := engine.effectiveSkillUseCost(p0, first)[model.ElementWater]; got != 1 {
			t.Fatalf("1221115 should reduce first skill water use cost by 1, cost=%v", engine.effectiveSkillUseCost(p0, first))
		}
		secondCost := engine.effectiveSkillUseCost(p0, second)
		if secondCost[model.ElementWater] != 2 || secondCost[model.ElementAir] != 1 {
			t.Fatalf("1221115 should reduce only water component of second skill, cost=%v", secondCost)
		}

		target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id":  first.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(target.Position.Col),
			"target_row":   float64(target.Position.Row),
			"target_owner": float64(1),
		}})
		if err != nil {
			t.Fatalf("use discounted first skill: %v", err)
		}
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].TargetInstanceID != second.InstanceID {
			t.Fatalf("using first skill should consume only its next-use modifier, modifiers=%+v", p0.TempModifiers)
		}
		if got := engine.effectiveSkillUseCost(p0, first)[model.ElementWater]; got != 2 {
			t.Fatalf("first skill discount should be gone after use, cost=%v", engine.effectiveSkillUseCost(p0, first))
		}
		if got := engine.effectiveSkillUseCost(p0, second)[model.ElementWater]; got != 2 {
			t.Fatalf("second skill discount should remain after first skill use, cost=%v", engine.effectiveSkillUseCost(p0, second))
		}

		boostEngine := setupEffectTest(t)
		boostP0 := boostEngine.State.Players[0]
		boostMage := placeUnit(baseCard(t, "1221115"), 0, 0, 0, boostEngine)
		mainSkill := readySkill(baseCard(t, "3221106"), 0)
		boostSkill := readySkill(baseCard(t, "3221003"), 0)
		boostP0.Skills[0] = mainSkill
		boostP0.Skills[1] = boostSkill
		boostP0.Elements = map[string]int{model.ElementWater: 10}
		if err := (Card1221115WinterfellAntiMage{}).OnPrayer(&EffectContext{Engine: boostEngine, Source: boostMage, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1221115 boost prayer: %v", err)
		}
		boostTarget := placeUnit(baseCard(t, "1021001"), 1, 0, 0, boostEngine)
		if err := boostEngine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id":  mainSkill.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(boostTarget.Position.Col),
			"target_row":   float64(boostTarget.Position.Row),
			"target_owner": float64(1),
			"boost_ids":    []any{boostSkill.InstanceID},
		}}); err != nil {
			t.Fatalf("use discounted boost skill: %v", err)
		}
		if len(boostP0.TempModifiers) != 0 {
			t.Fatalf("main attack should consume next-use modifiers for both main and boost skills, modifiers=%+v", boostP0.TempModifiers)
		}

		zeroEngine := setupEffectTest(t)
		zeroP0 := zeroEngine.State.Players[0]
		zeroMage := placeUnit(baseCard(t, "1221115"), 0, 0, 0, zeroEngine)
		zeroMain := readySkill(baseCard(t, "3221106"), 0)
		zeroBoost := readySkill(baseCard(t, "3221003"), 0)
		zeroP0.Skills[0] = zeroMain
		zeroP0.Skills[1] = zeroBoost
		zeroP0.Elements = map[string]int{model.ElementWater: 10}
		zeroP0.TempModifiers = append(zeroP0.TempModifiers, TemporaryModifier{
			ID:               "boost-zero-should-remain",
			Type:             TempModNextSkillCostZero,
			TargetInstanceID: zeroBoost.InstanceID,
			RemainingUses:    1,
		})
		if err := (Card1221115WinterfellAntiMage{}).OnPrayer(&EffectContext{Engine: zeroEngine, Source: zeroMage, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1221115 zero boost prayer: %v", err)
		}
		zeroTarget := placeUnit(baseCard(t, "1021001"), 1, 0, 0, zeroEngine)
		if err := zeroEngine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id":  zeroMain.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(zeroTarget.Position.Col),
			"target_row":   float64(zeroTarget.Position.Row),
			"target_owner": float64(1),
			"boost_ids":    []any{zeroBoost.InstanceID},
		}}); err != nil {
			t.Fatalf("use boost skill with unrelated zero-cost modifier: %v", err)
		}
		if len(zeroP0.TempModifiers) != 1 || zeroP0.TempModifiers[0].Type != TempModNextSkillCostZero || zeroP0.TempModifiers[0].TargetInstanceID != zeroBoost.InstanceID {
			t.Fatalf("boost use should consume only applied -1 water modifier and keep unapplied zero-cost modifier, modifiers=%+v", zeroP0.TempModifiers)
		}
	})

	t.Run("coral belly permanently empowers the first spell attack this game", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		p0.Hero = NewCardInstance(baseCard(t, "4211101"), 0, engine.State.TurnNumber)
		sorcery := readySkill(baseCard(t, "3021003"), 0)
		first := readySkill(baseCard(t, "3221106"), 0)
		second := readySkill(baseCard(t, "3121109"), 0)
		p0.Skills[2] = sorcery
		p0.Skills[0] = first
		p0.Skills[1] = second
		p0.Elements = map[string]int{model.ElementWater: 10, model.ElementFire: 10}
		target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": sorcery.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast sorcery before coral belly trigger: %v", err)
		}
		if sorcery.PowerBonus != 0 || p0.Hero.Statuses[coralBellyFirstSpellAttackUsedStatus] != 0 {
			t.Fatalf("4211101 should ignore sorceries before the first spell attack, sorcery_power=%d statuses=%v", sorcery.PowerBonus, p0.Hero.Statuses)
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id":  first.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(target.Position.Col),
			"target_row":   float64(target.Position.Row),
			"target_owner": float64(1),
		}}); err != nil {
			t.Fatalf("cast first spell with coral belly: %v", err)
		}
		if first.PowerBonus != 3 || p0.Hero.Statuses[coralBellyFirstSpellAttackUsedStatus] != 1 {
			t.Fatalf("4211101 should permanently give the first spell +3 power and mark itself, power=%d statuses=%v", first.PowerBonus, p0.Hero.Statuses)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != first.Card.Power+3 {
			t.Fatalf("4211101 bonus should affect the current spell power, pending=%+v base=%d", engine.State.PendingSpell, first.Card.Power)
		}

		engine.State.PendingSpell = nil
		engine.State.Phase = PhaseMain
		secondTarget := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id":  second.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(secondTarget.Position.Col),
			"target_row":   float64(secondTarget.Position.Row),
			"target_owner": float64(1),
		}}); err != nil {
			t.Fatalf("cast second spell with coral belly: %v", err)
		}
		if second.PowerBonus != 0 {
			t.Fatalf("4211101 should trigger only once per game, second power bonus=%d", second.PowerBonus)
		}
	})

	t.Run("silverleaf ranger consumes for the next spell attack bonus", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		ranger := NewCardInstance(baseCard(t, "1321106"), 0, 1)
		ranger.IsHorizontal = false
		if err := (Card1321106SilverleafRanger{}).OnPerTurn(&EffectContext{Engine: engine, Source: ranger, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("silverleaf ranger ability: %v", err)
		}
		if !ranger.IsHorizontal {
			t.Fatal("silverleaf ranger should consume itself")
		}
		spell := readySkill(baseCard(t, "3021005"), 0)
		damage := engine.effectiveSpellDamage(0, spell, max(spell.Card.Attack+spell.AttackBonus, 0), nil)
		if damage != 1 {
			t.Fatalf("silverleaf ranger should add +1 attack to next spell, damage=%d modifiers=%v", damage, p0.TempModifiers)
		}
		engine.consumeNextSpellAttackBonuses(p0, spell)
		if len(p0.TempModifiers) != 0 {
			t.Fatalf("silverleaf ranger attack bonus should be consumed once, modifiers=%v", p0.TempModifiers)
		}
	})

	t.Run("cave elf pickaxe consumes to flip a chosen card kind from top five", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		pickaxe := NewCardInstance(baseCard(t, "2421109"), 0, 1)
		pickaxe.SlotIndex = 0
		pickaxe.IsHorizontal = false
		p0.Equipment[0] = pickaxe
		unflippable := NewCardInstance(baseCard(t, "2211101"), 0, 1)
		item := NewCardInstance(baseCard(t, "2021101"), 0, 1)
		companion := NewCardInstance(baseCard(t, "1421101"), 0, 1)
		tooDeep := NewCardInstance(baseCard(t, "1421102"), 0, 1)
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "3021005"), 0, 1),
			unflippable,
			item,
			companion,
			NewCardInstance(baseCard(t, "3021006"), 0, 1),
			tooDeep,
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  pickaxe.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use 2421109: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "cave_elf_pickaxe_kind" || len(engine.State.PendingAction.Candidates) != 2 {
			t.Fatalf("2421109 should ask for companion/item kind, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, "companion")
		if !pickaxe.IsHorizontal || pickaxe.UsedThisTurn != 1 {
			t.Fatalf("2421109 should consume itself and spend ability use, horizontal=%v used=%d", pickaxe.IsHorizontal, pickaxe.UsedThisTurn)
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != companion || containsCardInstance(p0.Hand, tooDeep) {
			t.Fatalf("2421109 should flip the first companion within top five only, hand=%v", cardsToInfo(p0.Hand))
		}
		if !containsCardInstance(p0.Deck, item) || !containsCardInstance(p0.Deck, unflippable) || !containsCardInstance(p0.Deck, tooDeep) {
			t.Fatalf("2421109 should leave nonmatching/unflippable/out-of-range cards in deck, deck=%v", cardsToInfo(p0.Deck))
		}

		itemEngine := setupReportedBugEngine(t)
		itemP0 := itemEngine.State.Players[0]
		itemPickaxe := NewCardInstance(baseCard(t, "2421109"), 0, 1)
		itemPickaxe.SlotIndex = 0
		itemPickaxe.IsHorizontal = false
		itemP0.Equipment[0] = itemPickaxe
		itemTarget := NewCardInstance(baseCard(t, "2021101"), 0, 1)
		itemP0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1421101"), 0, 1), itemTarget}
		if err := itemEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  itemPickaxe.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use 2421109 for item: %v", err)
		}
		resolvePendingSelection(t, itemEngine, 0, "item")
		if len(itemP0.Hand) != 1 || itemP0.Hand[0] != itemTarget {
			t.Fatalf("2421109 should flip an item when item kind is selected, hand=%v", cardsToInfo(itemP0.Hand))
		}

		failEngine := setupReportedBugEngine(t)
		failP0 := failEngine.State.Players[0]
		failPickaxe := NewCardInstance(baseCard(t, "2421109"), 0, 1)
		failPickaxe.SlotIndex = 0
		failPickaxe.IsHorizontal = true
		failP0.Equipment[0] = failPickaxe
		err := failEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  failPickaxe.InstanceID,
			"ability_type": "per_turn",
		}})
		if err == nil || failPickaxe.UsedThisTurn != 0 {
			t.Fatalf("2421109 should require a ready source without burning use, err=%v used=%d", err, failPickaxe.UsedThisTurn)
		}
	})

	t.Run("autumn maple gem marks itself and resets a horizontal earth companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		gem := NewCardInstance(baseCard(t, "2421112"), 0, 1)
		if err := (Card2421112AutumnMapleGem{}).OnEnter(&EffectContext{Engine: engine, Source: gem, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("autumn maple gem enter: %v", err)
		}
		if gem.Statuses[autumnMapleGemCounter] != 2 {
			t.Fatalf("autumn maple gem should enter with two counters, statuses=%v", gem.Statuses)
		}
		earth := placeUnit(baseCard(t, "1421101"), 0, 0, 0, engine)
		earth.IsHorizontal = true
		if err := (Card2421112AutumnMapleGem{}).OnPerTurn(&EffectContext{Engine: engine, Source: gem, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("autumn maple gem ability: %v", err)
		}
		resolvePendingSelection(t, engine, 0, earth.InstanceID)
		if earth.IsHorizontal || gem.Statuses[autumnMapleGemCounter] != 1 {
			t.Fatalf("autumn maple gem should spend one counter to reset earth companion, horizontal=%v statuses=%v", earth.IsHorizontal, gem.Statuses)
		}

		failEngine := setupReportedBugEngine(t)
		failGem := NewCardInstance(baseCard(t, "2421112"), 0, 1)
		failEngine.State.Players[0].Equipment[0] = failGem
		if err := (Card2421112AutumnMapleGem{}).OnPerTurn(&EffectContext{Engine: failEngine, Source: failGem, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("autumn maple gem should no-op without counters: %v", err)
		}
		if failEngine.State.PendingAction != nil {
			t.Fatalf("autumn maple gem should not open a prompt without counters, pending=%+v", failEngine.State.PendingAction)
		}
		err := failEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  failGem.InstanceID,
			"ability_type": "per_turn",
		}})
		if err == nil || failGem.UsedThisTurn != 0 {
			t.Fatalf("autumn maple gem action should fail without burning use, err=%v used=%d", err, failGem.UsedThisTurn)
		}
	})
}

func TestRoyalConflictDeathrattleEffects(t *testing.T) {
	t.Run("abandoned pawn damages adjacent units and replaces killed companions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		pawn := NewCardInstance(baseCard(t, "1001101"), 0, 1)
		pawn.Position = &Position{Col: 1, Row: 1}
		killed := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		killed.CurrentLife = 1
		survivor := placeUnit(baseCard(t, "1021001"), 0, 0, 1, engine)
		survivor.CurrentLife = 3
		if err := (Card1001101AbandonedPawn{}).OnDeath(&EffectContext{Engine: engine, Source: pawn, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("abandoned pawn deathrattle: %v", err)
		}
		if engine.State.Players[1].Units[1][0] == nil || engine.State.Players[1].Units[1][0].Card.Number != "1001101" {
			t.Fatalf("abandoned pawn should replace killed adjacent companion, unit=%v", cardToInfo(engine.State.Players[1].Units[1][0]))
		}
		if len(engine.State.Players[1].Graveyard) != 1 || engine.State.Players[1].Graveyard[0] != killed {
			t.Fatalf("killed adjacent companion should enter graveyard, grave=%v", cardsToInfo(engine.State.Players[1].Graveyard))
		}
		if survivor.CurrentLife != 2 || engine.State.Players[0].Units[0][1] != survivor {
			t.Fatalf("surviving adjacent unit should only take damage, life=%d", survivor.CurrentLife)
		}
	})

	t.Run("contradictory knight is summoned for opponent with reduced max life", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		knight := NewCardInstance(baseCard(t, "1521108"), 0, 1)
		engine.State.Players[0].Graveyard = append(engine.State.Players[0].Graveyard, knight)
		if err := (Card1521108ContradictoryKnight{}).OnDeath(&EffectContext{Engine: engine, Source: knight, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("contradictory knight deathrattle: %v", err)
		}
		posID := positionSelectionID(Position{Col: 0, Row: 0})
		resolvePendingSelection(t, engine, 1, posID)
		if len(engine.State.Players[0].Graveyard) != 0 {
			t.Fatalf("contradictory knight should leave original graveyard, grave=%v", cardsToInfo(engine.State.Players[0].Graveyard))
		}
		summoned := engine.State.Players[1].Units[0][0]
		if summoned != knight || summoned.OwnerID != 1 || summoned.Card.Life != 3 || summoned.CurrentLife != 3 {
			t.Fatalf("contradictory knight should switch sides with max life -1, summoned=%+v", cardToInfo(summoned))
		}

		staleEngine := setupReportedBugEngine(t)
		staleKnight := NewCardInstance(baseCard(t, "1521108"), 0, 1)
		staleEngine.State.Players[0].Graveyard = append(staleEngine.State.Players[0].Graveyard, staleKnight)
		if err := (Card1521108ContradictoryKnight{}).OnDeath(&EffectContext{Engine: staleEngine, Source: staleKnight, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("contradictory knight stale deathrattle: %v", err)
		}
		blocker := placeUnit(baseCard(t, "1021001"), 1, 0, 0, staleEngine)
		resolvePendingSelection(t, staleEngine, 1, positionSelectionID(Position{Col: 0, Row: 0}))
		if len(staleEngine.State.Players[0].Graveyard) != 1 || staleEngine.State.Players[0].Graveyard[0] != staleKnight || staleKnight.OwnerID != 0 || staleKnight.Card.Life != 4 {
			t.Fatalf("stale contradictory knight position should leave card in original graveyard, grave=%v owner=%d life=%d", cardsToInfo(staleEngine.State.Players[0].Graveyard), staleKnight.OwnerID, staleKnight.Card.Life)
		}
		if staleEngine.State.Players[1].Units[0][0] != blocker {
			t.Fatalf("stale contradictory knight should not overwrite occupied position")
		}
	})

	t.Run("radiant watchdog searches a discounted companion only when enemy killed it", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		watchdog := NewCardInstance(baseCard(t, "1521113"), 0, 1)
		target := NewCardInstance(baseCard(t, "1521101"), 0, 1)
		p0.Deck = []*CardInstance{target}
		if err := (Card1521113RadiantWatchdog{}).OnDeath(&EffectContext{Engine: engine, Source: watchdog, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 1}}); err != nil {
			t.Fatalf("radiant watchdog deathrattle: %v", err)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if len(p0.Hand) != 1 || p0.Hand[0] != target {
			t.Fatalf("radiant watchdog should search companion to hand, hand=%v deck=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Deck))
		}
		if cost := engine.effectiveCardPlayCost(p0, target); cost[model.ElementLight] != max(target.Card.ElementsCost[model.ElementLight]-1, 0) {
			t.Fatalf("searched companion should have -1 light entry cost, cost=%v base=%v", cost, target.Card.ElementsCost)
		}

		friendlyEngine := setupReportedBugEngine(t)
		friendlyEngine.State.Players[0].Deck = []*CardInstance{NewCardInstance(baseCard(t, "1521101"), 0, 1)}
		if err := (Card1521113RadiantWatchdog{}).OnDeath(&EffectContext{Engine: friendlyEngine, Source: watchdog, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("friendly killed watchdog deathrattle: %v", err)
		}
		if friendlyEngine.State.PendingAction != nil {
			t.Fatalf("radiant watchdog should not trigger when killed by friendly source, pending=%+v", friendlyEngine.State.PendingAction)
		}
	})

	t.Run("soul symbiote marks up to two friendly spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		skill := readySkill(baseCard(t, "3021005"), 0)
		p0.Skills[0] = skill
		host := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		bound := readySkill(baseCard(t, "3221103"), 0)
		bound.SlotIndex = -1
		host.BoundSkills = []*CardInstance{bound}
		symbiote := NewCardInstance(baseCard(t, "1621114"), 0, 1)
		if err := (Card1621114SoulSymbiote{}).OnDeath(&EffectContext{Engine: engine, Source: symbiote, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("soul symbiote deathrattle: %v", err)
		}
		resolvePendingSelection(t, engine, 0, skill.InstanceID, bound.InstanceID)
		if skill.Statuses[soulMarkerStatus] != 1 || skill.PowerBonus != 2 || bound.Statuses[soulMarkerStatus] != 1 || bound.PowerBonus != 2 {
			t.Fatalf("soul symbiote should mark selected spells, skill status/power=%v/%d bound=%v/%d", skill.Statuses, skill.PowerBonus, bound.Statuses, bound.PowerBonus)
		}
	})
}

func TestRoyalConflictSimpleActiveAbilityEffects(t *testing.T) {
	t.Run("lone star tower watcher discards up to three hand cards for shield", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		watcher := placeUnit(baseCard(t, "1321103"), 0, 0, 0, engine)
		cardA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		cardB := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		cardC := NewCardInstance(baseCard(t, "1021003"), 0, 1)
		p0.Hand = []*CardInstance{cardA, cardB, cardC}
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  watcher.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("watcher ultimate: %v", err)
		}
		resolvePendingSelection(t, engine, 0, cardA.InstanceID, cardC.InstanceID)
		if p0.Shield != 2 || len(p0.Hand) != 1 || p0.Hand[0] != cardB || len(p0.Graveyard) != 2 {
			t.Fatalf("watcher should discard selected cards for shield, shield=%d hand=%v grave=%v", p0.Shield, cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard))
		}
	})

	t.Run("storm horn discards a hand card to search air equipment", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		horn := placeUnit(baseCard(t, "1321109"), 0, 0, 0, engine)
		discard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		equipment := NewCardInstance(baseCard(t, "2321101"), 0, 1)
		p0.Hand = []*CardInstance{discard}
		p0.Deck = []*CardInstance{equipment, NewCardInstance(baseCard(t, "1021002"), 0, 1)}
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  horn.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("storm horn ultimate: %v", err)
		}
		resolvePendingSelection(t, engine, 0, discard.InstanceID)
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "storm_horn_search_air_equipment" {
			t.Fatalf("storm horn should prompt to search air equipment, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, equipment.InstanceID)
		if len(p0.Graveyard) != 1 || p0.Graveyard[0] != discard || len(p0.Hand) != 1 || p0.Hand[0] != equipment {
			t.Fatalf("storm horn should discard cost and search equipment, hand=%v grave=%v deck=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard), cardsToInfo(p0.Deck))
		}

		failEngine := setupReportedBugEngine(t)
		failHorn := placeUnit(baseCard(t, "1321109"), 0, 0, 0, failEngine)
		failEngine.State.Players[0].Deck = []*CardInstance{NewCardInstance(baseCard(t, "2321101"), 0, 1)}
		err := failEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  failHorn.InstanceID,
			"ability_type": "ultimate",
		}})
		if err == nil || failHorn.UltimateUsed {
			t.Fatalf("storm horn should fail without a hand card and not burn ultimate, err=%v ultimate=%v", err, failHorn.UltimateUsed)
		}
	})

	t.Run("jiuxiao radiance discards both hands then draws the same counts", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		radiance := NewCardInstance(baseCard(t, "2511101"), 0, 1)
		p0.Equipment[0] = radiance
		oldP0 := []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1), NewCardInstance(baseCard(t, "1021002"), 0, 1)}
		oldP1 := []*CardInstance{NewCardInstance(baseCard(t, "1021003"), 1, 1)}
		p0.Hand = oldP0
		p1.Hand = oldP1
		draw0A := NewCardInstance(baseCard(t, "1021004"), 0, 1)
		draw0B := NewCardInstance(baseCard(t, "1021005"), 0, 1)
		draw1 := NewCardInstance(baseCard(t, "1021006"), 1, 1)
		p0.Deck = []*CardInstance{draw0A, draw0B}
		p1.Deck = []*CardInstance{draw1}
		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  radiance.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("jiuxiao radiance ultimate: %v", err)
		}
		if len(p0.Graveyard) != 2 || len(p1.Graveyard) != 1 || len(p0.Hand) != 2 || len(p1.Hand) != 1 {
			t.Fatalf("jiuxiao radiance should discard all hands then draw same counts, p0 hand/grave=%v/%v p1 hand/grave=%v/%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard), cardsToInfo(p1.Hand), cardsToInfo(p1.Graveyard))
		}
		if p0.Hand[0] != draw0A || p0.Hand[1] != draw0B || p1.Hand[0] != draw1 {
			t.Fatalf("jiuxiao radiance should draw replacement cards from deck, p0=%v p1=%v", cardsToInfo(p0.Hand), cardsToInfo(p1.Hand))
		}
	})
}

func TestRoyalConflictTriggeredPerTurnEffects(t *testing.T) {
	t.Run("lava fort hellhound damages two units after effect consume", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		hellhound := placeUnit(baseCard(t, "1121113"), 0, 0, 0, engine)
		ally := placeUnit(baseCard(t, "1121101"), 0, 1, 0, engine)
		enemy := placeUnit(baseCard(t, "1121101"), 1, 0, 0, engine)
		outOfRangeEnemy := placeUnit(baseCard(t, "1121101"), 1, 1, 1, engine)
		behavior := Card1121113LavaFortHellhound{}

		if err := behavior.OnConsume(&EffectContext{
			Engine:     engine,
			Source:     hellhound,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"consumed_player": 0, "gained": map[string]int{model.ElementFire: 3}},
		}); err != nil {
			t.Fatalf("1121113 normal consume: %v", err)
		}
		if engine.State.PendingAction != nil || hellhound.UsedThisTurn != 0 {
			t.Fatalf("1121113 should ignore ordinary consumes, pending=%+v used=%d", engine.State.PendingAction, hellhound.UsedThisTurn)
		}

		if err := behavior.OnConsume(&EffectContext{
			Engine:     engine,
			Source:     hellhound,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"consumed_player": 0, "consume_source": "2121108", "gained": map[string]int{model.ElementFire: 3}},
		}); err != nil {
			t.Fatalf("1121113 effect consume: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "lava_fort_hellhound_damage" || hellhound.UsedThisTurn != 1 {
			t.Fatalf("1121113 should ask for two damage targets after effect consume, pending=%+v used=%d", engine.State.PendingAction, hellhound.UsedThisTurn)
		}
		for _, candidate := range engine.State.PendingAction.Candidates {
			if candidate["instance_id"] == outOfRangeEnemy.InstanceID {
				t.Fatalf("1121113 should not offer out-of-range enemy units, candidates=%+v", engine.State.PendingAction.Candidates)
			}
		}
		resolvePendingSelection(t, engine, 0, ally.InstanceID, enemy.InstanceID)
		if ally.CurrentLife != ally.Card.Life-1 || enemy.CurrentLife != enemy.Card.Life-1 || outOfRangeEnemy.CurrentLife != outOfRangeEnemy.Card.Life {
			t.Fatalf("1121113 should damage selected units only, ally=%d enemy=%d out=%d", ally.CurrentLife, enemy.CurrentLife, outOfRangeEnemy.CurrentLife)
		}

		if err := behavior.OnConsume(&EffectContext{
			Engine:     engine,
			Source:     hellhound,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"consumed_player": 0, "consume_source": "2121108"},
		}); err != nil {
			t.Fatalf("1121113 second effect consume: %v", err)
		}
		if engine.State.PendingAction != nil || hellhound.UsedThisTurn != 1 {
			t.Fatalf("1121113 should trigger at most once per turn, pending=%+v used=%d", engine.State.PendingAction, hellhound.UsedThisTurn)
		}

		otherConsumeEngine := setupReportedBugEngine(t)
		observer := placeUnit(baseCard(t, "1121113"), 0, 0, 0, otherConsumeEngine)
		consumedOther := placeUnit(baseCard(t, "1121101"), 0, 1, 0, otherConsumeEngine)
		if err := behavior.OnConsume(&EffectContext{
			Engine:     otherConsumeEngine,
			Source:     observer,
			Target:     consumedOther,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"consumed_player": 0, "consume_source": "2121108"},
		}); err != nil {
			t.Fatalf("1121113 observing other effect consume: %v", err)
		}
		if otherConsumeEngine.State.PendingAction != nil || observer.UsedThisTurn != 0 {
			t.Fatalf("1121113 should not trigger when another unit is consumed, pending=%+v used=%d", otherConsumeEngine.State.PendingAction, observer.UsedThisTurn)
		}
	})

	t.Run("curse box marks deaths and spends markers to weaken enemy spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p1 := engine.State.Players[1]
		box := NewCardInstance(baseCard(t, "2621107"), 0, 1)
		p1.Skills[0] = readySkill(baseCard(t, "3321005"), 1)
		p1.Skills[1] = readySkill(baseCard(t, "3221001"), 1)
		p1.Skills[2] = readySkill(baseCard(t, "3621006"), 1)
		behavior := Card2621107CurseBox{}

		if err := behavior.OnFriendlyDeath(&EffectContext{Engine: engine, Source: box, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2621107 friendly death: %v", err)
		}
		if err := behavior.OnEnemyDeath(&EffectContext{Engine: engine, Source: box, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2621107 enemy death: %v", err)
		}
		box.Statuses[curseBoxMarkerStatus] += 2
		if box.Statuses[curseBoxMarkerStatus] != 4 {
			t.Fatalf("2621107 should mark every unit death, statuses=%v", box.Statuses)
		}

		if err := behavior.OnPerTurn(&EffectContext{Engine: engine, Source: box, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2621107 per-turn: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "curse_box_weaken" || engine.State.PendingAction.MaxSelect != 3 {
			t.Fatalf("2621107 should ask for up to 3 enemy spells, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, p1.Skills[0].InstanceID, p1.Skills[1].InstanceID, p1.Skills[2].InstanceID)
		if box.Statuses[curseBoxMarkerStatus] != 1 {
			t.Fatalf("2621107 should remove one marker per selected spell, statuses=%v", box.Statuses)
		}
		for i := 0; i < 3; i++ {
			if p1.Skills[i].Statuses[StatusWeaken] != 1 {
				t.Fatalf("2621107 should weaken selected enemy spells by 1, skill %d statuses=%v", i, p1.Skills[i].Statuses)
			}
		}

		emptyBox := NewCardInstance(baseCard(t, "2621107"), 0, 1)
		if err := behavior.OnPerTurn(&EffectContext{Engine: engine, Source: emptyBox, PlayerID: 0, OpponentID: 1}); err == nil {
			t.Fatal("2621107 should reject active ability with no markers")
		}

		deathEngine := setupReportedBugEngine(t)
		deathBox := NewCardInstance(baseCard(t, "2621107"), 0, 1)
		deathEngine.State.Players[0].Equipment[0] = deathBox
		friendlyDead := placeUnit(baseCard(t, "1021001"), 0, 0, 0, deathEngine)
		enemyDead := placeUnit(baseCard(t, "1021001"), 1, 0, 0, deathEngine)
		deathEngine.destroyUnit(friendlyDead, 0)
		deathEngine.destroyUnit(enemyDead, 1)
		if deathBox.Statuses[curseBoxMarkerStatus] != 2 {
			t.Fatalf("2621107 should mark real friendly and enemy unit deaths once each, statuses=%v", deathBox.Statuses)
		}
	})

	t.Run("soul hunter marks friendly spell once after it hits", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		hunter := placeUnit(baseCard(t, "1621106"), 0, 0, 0, engine)
		skill := readySkill(baseCard(t, "3021005"), 0)
		behavior := Card1621106SoulHunter{}

		if err := behavior.OnSpellHit(&EffectContext{
			Engine:     engine,
			Source:     hunter,
			Target:     skill,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"attacker": 0, "spell_source": skill},
		}); err != nil {
			t.Fatalf("1621106 friendly spell hit: %v", err)
		}
		if skill.Statuses[soulMarkerStatus] != 1 || skill.PowerBonus != 2 || hunter.UsedThisTurn != 1 {
			t.Fatalf("1621106 should mark the hit spell once, statuses=%v power=%d used=%d", skill.Statuses, skill.PowerBonus, hunter.UsedThisTurn)
		}

		if err := behavior.OnSpellHit(&EffectContext{
			Engine:     engine,
			Source:     hunter,
			Target:     skill,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"attacker": 0, "spell_source": skill},
		}); err != nil {
			t.Fatalf("1621106 second friendly spell hit: %v", err)
		}
		if skill.Statuses[soulMarkerStatus] != 1 || skill.PowerBonus != 2 || hunter.UsedThisTurn != 1 {
			t.Fatalf("1621106 should trigger at most once per turn, statuses=%v power=%d used=%d", skill.Statuses, skill.PowerBonus, hunter.UsedThisTurn)
		}

		enemyEngine := setupReportedBugEngine(t)
		enemyHunter := placeUnit(baseCard(t, "1621106"), 0, 0, 0, enemyEngine)
		enemySkill := readySkill(baseCard(t, "3021005"), 1)
		if err := behavior.OnSpellHit(&EffectContext{
			Engine:     enemyEngine,
			Source:     enemyHunter,
			Target:     enemySkill,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"attacker": 1, "spell_source": enemySkill},
		}); err != nil {
			t.Fatalf("1621106 enemy spell hit: %v", err)
		}
		if enemySkill.Statuses[soulMarkerStatus] != 0 || enemySkill.PowerBonus != 0 || enemyHunter.UsedThisTurn != 0 {
			t.Fatalf("1621106 should ignore enemy spell hits, statuses=%v power=%d used=%d", enemySkill.Statuses, enemySkill.PowerBonus, enemyHunter.UsedThisTurn)
		}

		missingSourceEngine := setupReportedBugEngine(t)
		missingSourceHunter := placeUnit(baseCard(t, "1621106"), 0, 0, 0, missingSourceEngine)
		if err := behavior.OnSpellHit(&EffectContext{
			Engine:     missingSourceEngine,
			Source:     missingSourceHunter,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"attacker": 0},
		}); err != nil {
			t.Fatalf("1621106 missing source spell hit: %v", err)
		}
		if missingSourceHunter.UsedThisTurn != 0 {
			t.Fatalf("1621106 should not spend trigger without a skill source, used=%d", missingSourceHunter.UsedThisTurn)
		}

		fieldEngine := setupReportedBugEngine(t)
		fieldHunter := placeUnit(baseCard(t, "1621106"), 0, 0, 0, fieldEngine)
		fieldSkill := readySkill(baseCard(t, "3021005"), 0)
		fieldEngine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, fieldSkill, map[string]any{"attacker": 0, "spell_source": fieldSkill})
		if fieldSkill.Statuses[soulMarkerStatus] != 1 || fieldSkill.PowerBonus != 2 || fieldHunter.UsedThisTurn != 1 {
			t.Fatalf("1621106 should trigger through field spell-hit plumbing, statuses=%v power=%d used=%d", fieldSkill.Statuses, fieldSkill.PowerBonus, fieldHunter.UsedThisTurn)
		}

		scrollEngine := setupReportedBugEngine(t)
		scrollHunter := placeUnit(baseCard(t, "1621106"), 0, 0, 0, scrollEngine)
		scroll := NewCardInstance(baseCard(t, "2121003"), 0, 1)
		scrollEngine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, scroll, map[string]any{"attacker": 0, "spell_source": scroll})
		if scroll.Statuses[soulMarkerStatus] != 1 || scroll.PowerBonus != 2 || scrollHunter.UsedThisTurn != 1 {
			t.Fatalf("1621106 should mark spell scroll hits too, statuses=%v power=%d used=%d", scroll.Statuses, scroll.PowerBonus, scrollHunter.UsedThisTurn)
		}
	})

	t.Run("rock wall monk zeros the first enemy spell hit while no skills are learned", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		monk := placeUnit(baseCard(t, "1421113"), 0, 0, 0, engine)
		enemySkill := readySkill(baseCard(t, "3021005"), 1)
		behavior := Card1421113RockWallMonk{}

		damage := 3
		if err := behavior.OnSpellHitBeforeDamage(&EffectContext{
			Engine:     engine,
			Source:     monk,
			Target:     monk,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"attacker": 1, "spell_source": enemySkill, "damage_ptr": &damage, "damage": damage},
		}); err != nil {
			t.Fatalf("1421113 enemy spell hit: %v", err)
		}
		if damage != 0 || monk.UsedThisTurn != 1 {
			t.Fatalf("1421113 should zero the first enemy spell hit and spend trigger, damage=%d used=%d", damage, monk.UsedThisTurn)
		}

		secondDamage := 4
		if err := behavior.OnSpellHitBeforeDamage(&EffectContext{
			Engine:     engine,
			Source:     monk,
			Target:     monk,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"attacker": 1, "spell_source": enemySkill, "damage_ptr": &secondDamage, "damage": secondDamage},
		}); err != nil {
			t.Fatalf("1421113 second enemy spell hit: %v", err)
		}
		if secondDamage != 4 || monk.UsedThisTurn != 1 {
			t.Fatalf("1421113 should trigger at most once per turn, damage=%d used=%d", secondDamage, monk.UsedThisTurn)
		}

		friendlyEngine := setupReportedBugEngine(t)
		friendlyMonk := placeUnit(baseCard(t, "1421113"), 0, 0, 0, friendlyEngine)
		friendlyDamage := 3
		if err := behavior.OnSpellHitBeforeDamage(&EffectContext{
			Engine:     friendlyEngine,
			Source:     friendlyMonk,
			Target:     friendlyMonk,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"attacker": 0, "spell_source": readySkill(baseCard(t, "3021005"), 0), "damage_ptr": &friendlyDamage, "damage": friendlyDamage},
		}); err != nil {
			t.Fatalf("1421113 friendly spell hit: %v", err)
		}
		if friendlyDamage != 3 || friendlyMonk.UsedThisTurn != 0 {
			t.Fatalf("1421113 should ignore friendly spell hits, damage=%d used=%d", friendlyDamage, friendlyMonk.UsedThisTurn)
		}

		learnedEngine := setupReportedBugEngine(t)
		learnedMonk := placeUnit(baseCard(t, "1421113"), 0, 0, 0, learnedEngine)
		learnedEngine.State.Players[0].Skills[0] = readySkill(baseCard(t, "3021005"), 0)
		learnedDamage := 3
		if err := behavior.OnSpellHitBeforeDamage(&EffectContext{
			Engine:     learnedEngine,
			Source:     learnedMonk,
			Target:     learnedMonk,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"attacker": 1, "spell_source": enemySkill, "damage_ptr": &learnedDamage, "damage": learnedDamage},
		}); err != nil {
			t.Fatalf("1421113 learned skill enemy spell hit: %v", err)
		}
		if learnedDamage != 3 || learnedMonk.UsedThisTurn != 0 {
			t.Fatalf("1421113 should not trigger after learning a skill, damage=%d used=%d", learnedDamage, learnedMonk.UsedThisTurn)
		}

		fieldEngine := setupReportedBugEngine(t)
		fieldMonk := placeUnit(baseCard(t, "1421113"), 0, 0, 0, fieldEngine)
		fieldDamage := 5
		fieldEngine.triggerFieldEffectsWithData(TriggerOnSpellHitBeforeDamage, 0, enemySkill, map[string]any{
			"attacker": 1, "spell_source": enemySkill, "damage_ptr": &fieldDamage, "damage": fieldDamage,
		})
		if fieldDamage != 0 || fieldMonk.UsedThisTurn != 1 {
			t.Fatalf("1421113 should trigger through field before-damage plumbing, damage=%d used=%d", fieldDamage, fieldMonk.UsedThisTurn)
		}
	})

	t.Run("spark moth reveals from hand after fire spell hits for entry discount", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		mothA := NewCardInstance(baseCard(t, "1121112"), 0, 1)
		mothB := NewCardInstance(baseCard(t, "1121112"), 0, 1)
		otherFire := NewCardInstance(baseCard(t, "1121101"), 0, 1)
		enemyMoth := NewCardInstance(baseCard(t, "1121112"), 1, 1)
		p0.Hand = []*CardInstance{mothA, mothB, otherFire}
		p1.Hand = []*CardInstance{enemyMoth}

		engine.triggerSparkMothAfterSpellHit(readySkill(baseCard(t, "3021005"), 0))
		if engine.State.PendingAction != nil {
			t.Fatalf("1121112 should ignore non-fire spell hits, pending=%+v", engine.State.PendingAction)
		}

		engine.triggerSparkMothAfterSpellHit(readySkill(baseCard(t, "3121001"), 0))
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "spark_moth_reveal" || engine.State.PendingAction.PlayerID != 0 || engine.State.PendingAction.MaxSelect != 2 {
			t.Fatalf("1121112 should prompt the first player with moths in hand, pending=%+v", engine.State.PendingAction)
		}
		if len(engine.State.PendingActionQueue) != 1 || engine.State.PendingActionQueue[0].PlayerID != 1 {
			t.Fatalf("1121112 should queue the other player's reveal prompt, queue=%+v", engine.State.PendingActionQueue)
		}

		resolvePendingSelection(t, engine, 0)
		if p0.RevealedHand[mothA.InstanceID] || p0.RevealedHand[mothB.InstanceID] || mothA.Statuses["入场费用"+model.ElementFire+"-1"] != 0 || mothB.Statuses["入场费用"+model.ElementFire+"-1"] != 0 {
			t.Fatalf("1121112 skipped reveal should not discount or reveal, revealed=%v statuses=%v/%v", p0.RevealedHand, mothA.Statuses, mothB.Statuses)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.PlayerID != 1 {
			t.Fatalf("1121112 should advance to queued enemy reveal prompt, pending=%+v", engine.State.PendingAction)
		}

		resolvePendingSelection(t, engine, 1, enemyMoth.InstanceID)
		if !p1.RevealedHand[enemyMoth.InstanceID] || enemyMoth.Statuses["入场费用"+model.ElementFire+"-1"] != 1 {
			t.Fatalf("1121112 selected enemy moth should reveal and discount, revealed=%v statuses=%v", p1.RevealedHand, enemyMoth.Statuses)
		}
		if cost := engine.effectiveCardPlayCost(p1, enemyMoth); cost[model.ElementFire] != enemyMoth.Card.ElementsCost[model.ElementFire]-1 {
			t.Fatalf("1121112 discount should reduce effective entry fire cost, cost=%v", cost)
		}

		engine.triggerSparkMothAfterSpellHit(readySkill(baseCard(t, "3121001"), 0))
		if engine.State.PendingAction == nil || engine.State.PendingAction.PlayerID != 0 {
			t.Fatalf("1121112 should prompt again on later fire spell hits, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, mothA.InstanceID, mothB.InstanceID)
		if !p0.RevealedHand[mothA.InstanceID] || !p0.RevealedHand[mothB.InstanceID] || mothA.Statuses["入场费用"+model.ElementFire+"-1"] != 1 || mothB.Statuses["入场费用"+model.ElementFire+"-1"] != 1 {
			t.Fatalf("1121112 selected moths should reveal and discount, revealed=%v statuses=%v/%v", p0.RevealedHand, mothA.Statuses, mothB.Statuses)
		}
		if cost := engine.effectiveCardPlayCost(p0, mothA); cost[model.ElementFire] != mothA.Card.ElementsCost[model.ElementFire]-1 {
			t.Fatalf("1121112 discount should apply to selected own moth, cost=%v", cost)
		}
		resolvePendingSelection(t, engine, 1)
		if enemyMoth.Statuses["入场费用"+model.ElementFire+"-1"] != 1 {
			t.Fatalf("1121112 skipped later reveal should not add another discount, statuses=%v", enemyMoth.Statuses)
		}

		hitEngine := setupReportedBugEngine(t)
		hitP0 := hitEngine.State.Players[0]
		hitMoth := NewCardInstance(baseCard(t, "1121112"), 0, 1)
		hitP0.Hand = []*CardInstance{hitMoth}
		target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, hitEngine)
		fireSpell := readySkill(baseCard(t, "3121001"), 0)
		hitEngine.resolveSpellHit(0, fireSpell, SpellTarget{Type: "unit", Position: *target.Position}, nil, nil)
		if hitEngine.State.PendingAction == nil || hitEngine.State.PendingAction.Type != "spark_moth_reveal" || hitEngine.State.PendingAction.PlayerID != 0 {
			t.Fatalf("1121112 should trigger through real fire spell hit resolution, pending=%+v", hitEngine.State.PendingAction)
		}
	})

	t.Run("celtic deer resets once after any medium skill is used", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		deer := placeUnit(baseCard(t, "1421108"), 0, 1, 1, engine)
		deer.IsHorizontal = true
		medium := readySkill(baseCard(t, "3421104"), 0)
		nonMedium := readySkill(baseCard(t, "3121001"), 0)
		behavior := Card1421108CelticDeer{}

		if err := behavior.OnSpellCast(&EffectContext{Engine: engine, Source: deer, Target: medium, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("1421108 friendly medium spell cast: %v", err)
		}
		if deer.IsHorizontal || deer.UsedThisTurn != 1 {
			t.Fatalf("1421108 should reset once after a medium skill is used, horizontal=%v used=%d", deer.IsHorizontal, deer.UsedThisTurn)
		}

		deer.IsHorizontal = true
		if err := behavior.OnSpellCast(&EffectContext{Engine: engine, Source: deer, Target: medium, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 1}}); err != nil {
			t.Fatalf("1421108 second medium spell cast: %v", err)
		}
		if !deer.IsHorizontal || deer.UsedThisTurn != 1 {
			t.Fatalf("1421108 should trigger at most once per turn, horizontal=%v used=%d", deer.IsHorizontal, deer.UsedThisTurn)
		}

		nextEngine := setupReportedBugEngine(t)
		nextDeer := placeUnit(baseCard(t, "1421108"), 0, 1, 1, nextEngine)
		nextDeer.IsHorizontal = true
		if err := behavior.OnSpellCast(&EffectContext{Engine: nextEngine, Source: nextDeer, Target: medium, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 1}}); err != nil {
			t.Fatalf("1421108 enemy medium spell cast: %v", err)
		}
		if nextDeer.IsHorizontal || nextDeer.UsedThisTurn != 1 {
			t.Fatalf("1421108 should reset after an enemy medium skill is used, horizontal=%v used=%d", nextDeer.IsHorizontal, nextDeer.UsedThisTurn)
		}

		failEngine := setupReportedBugEngine(t)
		failDeer := placeUnit(baseCard(t, "1421108"), 0, 1, 1, failEngine)
		failDeer.IsHorizontal = true
		if err := behavior.OnSpellCast(&EffectContext{Engine: failEngine, Source: failDeer, Target: nonMedium, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("1421108 non-medium spell cast: %v", err)
		}
		if !failDeer.IsHorizontal || failDeer.UsedThisTurn != 0 {
			t.Fatalf("1421108 should ignore non-medium skills, horizontal=%v used=%d", failDeer.IsHorizontal, failDeer.UsedThisTurn)
		}
	})

	t.Run("lone star fire seed gains fire load after other companions take fire damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		seed := placeUnit(baseCard(t, "1121111"), 0, 0, 0, engine)
		ally := placeUnit(baseCard(t, "1121101"), 0, 1, 0, engine)
		enemy := placeUnit(baseCard(t, "1121101"), 1, 0, 0, engine)

		engine.dealDamageWithExtra(ally, 1, 0, map[string]any{"damage_source": "test", "damage_element": model.ElementFire})
		if got := effectiveElementsGain(seed)[model.ElementFire]; got != seed.Card.ElementsGain[model.ElementFire]+1 || !seed.UltimateUsed {
			t.Fatalf("1121111 should gain fire load once and spend ultimate when another companion takes fire damage, load=%v used=%v", effectiveElementsGain(seed), seed.UltimateUsed)
		}

		engine.dealDamageWithExtra(enemy, 1, 1, map[string]any{"damage_source": "test", "damage_element": model.ElementFire})
		if got := effectiveElementsGain(seed)[model.ElementFire]; got != seed.Card.ElementsGain[model.ElementFire]+1 {
			t.Fatalf("1121111 should not trigger again after its ultimate is spent, load=%v", effectiveElementsGain(seed))
		}

		engine.dealDamageWithExtra(ally, 1, 0, map[string]any{"damage_source": "test", "damage_element": model.ElementWater})
		if got := effectiveElementsGain(seed)[model.ElementFire]; got != seed.Card.ElementsGain[model.ElementFire]+1 {
			t.Fatalf("1121111 should ignore non-fire damage, load=%v", effectiveElementsGain(seed))
		}

		engine.dealDamageWithExtra(seed, 1, 0, map[string]any{"damage_source": "test", "damage_element": model.ElementFire})
		if got := effectiveElementsGain(seed)[model.ElementFire]; got != seed.Card.ElementsGain[model.ElementFire]+1 {
			t.Fatalf("1121111 should ignore damage to itself, load=%v", effectiveElementsGain(seed))
		}
	})

	t.Run("pain soul gains shadow load once after being damaged", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		soul := placeUnit(baseCard(t, "1621101"), 0, 0, 0, engine)
		engine.dealDamageWithExtra(soul, 1, 0, map[string]any{"attacker": 1})
		if effectiveElementsGain(soul)[model.ElementShadow] != soul.Card.ElementsGain[model.ElementShadow]+1 || soul.UsedThisTurn != 1 {
			t.Fatalf("pain soul should gain one shadow load after damage, load=%v used=%d", effectiveElementsGain(soul), soul.UsedThisTurn)
		}
		engine.dealDamageWithExtra(soul, 1, 0, map[string]any{"attacker": 1})
		if effectiveElementsGain(soul)[model.ElementShadow] != soul.Card.ElementsGain[model.ElementShadow]+1 {
			t.Fatalf("pain soul should trigger at most once per turn, load=%v used=%d", effectiveElementsGain(soul), soul.UsedThisTurn)
		}

		watchEngine := setupReportedBugEngine(t)
		watcher := placeUnit(baseCard(t, "1621101"), 0, 0, 0, watchEngine)
		other := placeUnit(baseCard(t, "1021001"), 0, 1, 0, watchEngine)
		watchEngine.dealDamageWithExtra(other, 1, 0, map[string]any{"attacker": 1})
		if effectiveElementsGain(watcher)[model.ElementShadow] != watcher.Card.ElementsGain[model.ElementShadow] || watcher.UsedThisTurn != 0 {
			t.Fatalf("pain soul should not trigger when another unit is damaged, load=%v used=%d", effectiveElementsGain(watcher), watcher.UsedThisTurn)
		}
	})

	t.Run("pain avenger gains attack once after being damaged", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		avenger := placeUnit(baseCard(t, "1621102"), 0, 0, 0, engine)
		engine.dealDamageWithExtra(avenger, 1, 0, map[string]any{"attacker": 1})
		if avenger.CurrentAttack != avenger.Card.Attack+1 || avenger.UsedThisTurn != 1 {
			t.Fatalf("pain avenger should gain one attack after damage, attack=%d used=%d", avenger.CurrentAttack, avenger.UsedThisTurn)
		}
		engine.dealDamageWithExtra(avenger, 1, 0, map[string]any{"attacker": 1})
		if avenger.CurrentAttack != avenger.Card.Attack+1 {
			t.Fatalf("pain avenger should trigger at most once per turn, attack=%d used=%d", avenger.CurrentAttack, avenger.UsedThisTurn)
		}

		watchEngine := setupReportedBugEngine(t)
		watcher := placeUnit(baseCard(t, "1621102"), 0, 0, 0, watchEngine)
		other := placeUnit(baseCard(t, "1021001"), 0, 1, 0, watchEngine)
		watchEngine.dealDamageWithExtra(other, 1, 0, map[string]any{"attacker": 1})
		if watcher.CurrentAttack != watcher.Card.Attack || watcher.UsedThisTurn != 0 {
			t.Fatalf("pain avenger should not trigger when another unit is damaged, attack=%d used=%d", watcher.CurrentAttack, watcher.UsedThisTurn)
		}
	})

	t.Run("rose garden gardener heals a friendly unit once after a unit dies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		gardener := placeUnit(baseCard(t, "1621104"), 0, 0, 0, engine)
		wounded := placeUnit(baseCard(t, "1021002"), 0, 1, 0, engine)
		wounded.CurrentLife = maxLife(wounded) - 2
		dead := placeUnit(baseCard(t, "1021001"), 0, 2, 0, engine)
		engine.destroyUnitWithData(dead, 0, map[string]any{"attacker": 1})
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "rose_garden_gardener_heal" || gardener.UsedThisTurn != 0 {
			t.Fatalf("gardener should prompt once after friendly death, pending=%+v used=%d", engine.State.PendingAction, gardener.UsedThisTurn)
		}
		resolvePendingSelection(t, engine, 0, wounded.InstanceID)
		if wounded.CurrentLife != maxLife(wounded) || gardener.UsedThisTurn != 1 {
			t.Fatalf("gardener should heal selected friendly unit up to max and spend trigger, life=%d max=%d used=%d", wounded.CurrentLife, maxLife(wounded), gardener.UsedThisTurn)
		}

		anotherDead := placeUnit(baseCard(t, "1021001"), 0, 2, 1, engine)
		engine.destroyUnitWithData(anotherDead, 0, map[string]any{"attacker": 1})
		if engine.State.PendingAction != nil {
			t.Fatalf("gardener should trigger at most once per turn, pending=%+v used=%d", engine.State.PendingAction, gardener.UsedThisTurn)
		}

		staleEngine := setupReportedBugEngine(t)
		staleGardener := placeUnit(baseCard(t, "1621104"), 0, 0, 0, staleEngine)
		staleWounded := placeUnit(baseCard(t, "1021002"), 0, 1, 0, staleEngine)
		staleWounded.CurrentLife = maxLife(staleWounded) - 1
		staleDead := placeUnit(baseCard(t, "1021001"), 0, 2, 0, staleEngine)
		staleEngine.destroyUnitWithData(staleDead, 0, map[string]any{"attacker": 1})
		healUnit(staleWounded, 99)
		resolvePendingSelection(t, staleEngine, 0, staleWounded.InstanceID)
		if staleGardener.UsedThisTurn != 0 || staleWounded.CurrentLife != maxLife(staleWounded) {
			t.Fatalf("gardener should not spend trigger on stale full-health target, used=%d life=%d", staleGardener.UsedThisTurn, staleWounded.CurrentLife)
		}
	})
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

func TestRoyalConflictSimpleSkillEffects(t *testing.T) {
	t.Run("dragon blood treant removes a friendly load and gains shadow load", func(t *testing.T) {
		selfEngine := setupEffectTest(t)
		setElementsGain(selfEngine.State.Players[0].Hero, map[string]int{})
		selfTreant := placeUnit(baseCard(t, "1421107"), 0, 0, 0, selfEngine)
		behavior := Card1421107DragonBloodTreant{}

		if err := behavior.OnEnter(&EffectContext{Engine: selfEngine, Source: selfTreant, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1421107 self load enter: %v", err)
		}
		selfLoad := effectiveElementsGain(selfTreant)
		if selfLoad[model.ElementEarth] != 1 || selfLoad[model.ElementShadow] != 1 {
			t.Fatalf("1421107 should be able to remove its own load and gain shadow, load=%v", selfLoad)
		}

		engine := setupEffectTest(t)
		setElementsGain(engine.State.Players[0].Hero, map[string]int{})
		treant := placeUnit(baseCard(t, "1421107"), 0, 0, 0, engine)
		setElementsGain(treant, map[string]int{})
		target := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		setElementsGain(target, map[string]int{model.ElementFire: 1})

		if err := behavior.OnEnter(&EffectContext{Engine: engine, Source: treant, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1421107 single load enter: %v", err)
		}
		if effectiveElementsGain(target)[model.ElementFire] != 0 || effectiveElementsGain(treant)[model.ElementShadow] != 1 {
			t.Fatalf("1421107 should auto-remove the sole friendly load and gain shadow, target=%v treant=%v", effectiveElementsGain(target), effectiveElementsGain(treant))
		}

		multiEngine := setupEffectTest(t)
		setElementsGain(multiEngine.State.Players[0].Hero, map[string]int{})
		multiTreant := placeUnit(baseCard(t, "1421107"), 0, 0, 0, multiEngine)
		setElementsGain(multiTreant, map[string]int{})
		multiTarget := placeUnit(baseCard(t, "1021001"), 0, 1, 0, multiEngine)
		setElementsGain(multiTarget, map[string]int{model.ElementFire: 1, model.ElementWater: 1})
		if err := behavior.OnEnter(&EffectContext{Engine: multiEngine, Source: multiTreant, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1421107 multi load enter: %v", err)
		}
		if multiEngine.State.PendingAction == nil || multiEngine.State.PendingAction.Type != "dragon_blood_treant_remove_load" || len(multiEngine.State.PendingAction.Candidates) != 2 {
			t.Fatalf("1421107 should ask which load to remove, pending=%+v", multiEngine.State.PendingAction)
		}
		resolvePendingSelection(t, multiEngine, 0, multiTarget.InstanceID+"|"+model.ElementWater)
		load := effectiveElementsGain(multiTarget)
		if load[model.ElementFire] != 1 || load[model.ElementWater] != 0 || effectiveElementsGain(multiTreant)[model.ElementShadow] != 1 {
			t.Fatalf("1421107 should remove selected load only and gain shadow, target=%v treant=%v", load, effectiveElementsGain(multiTreant))
		}

		bonusEngine := setupEffectTest(t)
		setElementsGain(bonusEngine.State.Players[0].Hero, map[string]int{})
		bonusTreant := placeUnit(baseCard(t, "1421107"), 0, 0, 0, bonusEngine)
		setElementsGain(bonusTreant, map[string]int{})
		bonusTarget := placeUnit(baseCard(t, "1021001"), 0, 1, 0, bonusEngine)
		setElementsGain(bonusTarget, map[string]int{model.ElementFire: 1})
		bonusTarget.ElementsGainBonus = map[string]int{model.ElementWater: 1}
		if err := behavior.OnEnter(&EffectContext{Engine: bonusEngine, Source: bonusTreant, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1421107 bonus load enter: %v", err)
		}
		resolvePendingSelection(t, bonusEngine, 0, bonusTarget.InstanceID+"|"+model.ElementWater)
		bonusLoad := effectiveElementsGain(bonusTarget)
		if bonusLoad[model.ElementFire] != 1 || bonusLoad[model.ElementWater] != 0 || bonusTarget.ElementsGainBonus[model.ElementWater] != 0 || bonusTarget.ElementsGainSet[model.ElementFire] != 1 {
			t.Fatalf("1421107 should remove selected bonus load without changing base load, target=%v bonus=%v set=%v", bonusLoad, bonusTarget.ElementsGainBonus, bonusTarget.ElementsGainSet)
		}
	})

	t.Run("royal tax collector gains arcane when opponent draws until their next turn ends", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		collector := placeUnit(baseCard(t, "1021105"), 0, 0, 0, engine)
		p1.Equipment[0] = NewCardInstance(baseCard(t, "2311002"), 1, engine.State.TurnNumber)
		p1.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 1, engine.State.TurnNumber),
			NewCardInstance(baseCard(t, "1021002"), 1, engine.State.TurnNumber),
		}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021004"), 0, engine.State.TurnNumber)}
		behavior := Card1021105RoyalTaxCollector{}

		if err := behavior.OnEnter(&EffectContext{Engine: engine, Source: collector, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1021105 enter: %v", err)
		}
		if collector.Statuses[royalTaxCollectorUntilOpponentTurnEndStatus] != engine.State.TurnNumber {
			t.Fatalf("1021105 should arm tax window, statuses=%v", collector.Statuses)
		}

		engine.drawCards(0, 1)
		if p0.Elements[model.ElementArcane] != 0 || engine.State.PendingAction != nil {
			t.Fatalf("1021105 should ignore own draw and opponent draw listeners should not trigger for own-only cards, elements=%v pending=%+v", p0.Elements, engine.State.PendingAction)
		}
		engine.drawCards(1, 2)
		if p0.Elements[model.ElementArcane] != 2 {
			t.Fatalf("1021105 should gain one arcane per opponent draw, elements=%v", p0.Elements)
		}

		if err := behavior.OnTurnEnd(&EffectContext{
			Engine:     engine,
			Source:     collector,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"ended_player": 0},
		}); err != nil {
			t.Fatalf("1021105 own turn end: %v", err)
		}
		if collector.Statuses[royalTaxCollectorUntilOpponentTurnEndStatus] == 0 {
			t.Fatalf("1021105 should stay active through own turn end, statuses=%v", collector.Statuses)
		}
		if err := behavior.OnTurnEnd(&EffectContext{
			Engine:     engine,
			Source:     collector,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"ended_player": 1},
		}); err != nil {
			t.Fatalf("1021105 opponent turn end: %v", err)
		}
		if collector.Statuses[royalTaxCollectorUntilOpponentTurnEndStatus] != 0 {
			t.Fatalf("1021105 should expire at opponent turn end, statuses=%v", collector.Statuses)
		}
		p1.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021006"), 1, engine.State.TurnNumber)}
		engine.drawCards(1, 1)
		if p0.Elements[model.ElementArcane] != 2 {
			t.Fatalf("1021105 should stop after opponent turn end, elements=%v", p0.Elements)
		}
	})

	t.Run("fire beast trainer discounts the next fire beast or monster companion", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		trainer := placeUnit(baseCard(t, "1121106"), 0, 0, 0, engine)
		beast := NewCardInstance(baseCard(t, "1121102"), 0, engine.State.TurnNumber)
		machine := NewCardInstance(baseCard(t, "1121104"), 0, engine.State.TurnNumber)
		behavior := Card1121106FireBeastTrainer{}

		if err := behavior.OnEnter(&EffectContext{Engine: engine, Source: trainer, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1121106 enter: %v", err)
		}
		if trainer.Statuses[fireBeastTrainerDiscountStatus] != 1 {
			t.Fatalf("1121106 should arm one discount, statuses=%v", trainer.Statuses)
		}
		if cost := engine.effectiveCardPlayCost(p0, machine); cost[model.ElementFire] != machine.Card.ElementsCost[model.ElementFire] {
			t.Fatalf("1121106 should not discount fire machines, cost=%v", cost)
		}
		if cost := engine.effectiveCardPlayCost(p0, beast); cost[model.ElementFire] != 4 {
			t.Fatalf("1121106 should discount fire beast/monster companion by two, cost=%v", cost)
		}

		p0.Hand = append(p0.Hand, beast)
		p0.Elements[model.ElementFire] = 4
		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": beast.InstanceID,
			"col":         float64(1),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("1121106 discounted summon: %v", err)
		}
		if trainer.Statuses[fireBeastTrainerDiscountStatus] != 0 {
			t.Fatalf("1121106 discount should be consumed after matching summon, statuses=%v", trainer.Statuses)
		}
		nextBeast := NewCardInstance(baseCard(t, "1121101"), 0, engine.State.TurnNumber)
		if cost := engine.effectiveCardPlayCost(p0, nextBeast); cost[model.ElementFire] != nextBeast.Card.ElementsCost[model.ElementFire] {
			t.Fatalf("1121106 should discount only one matching companion, cost=%v", cost)
		}
	})

	t.Run("legion general prayer buffs fire spells until next turn end", func(t *testing.T) {
		engine := setupEffectTest(t)
		general := placeUnit(baseCard(t, "1121114"), 0, 0, 0, engine)
		fireSkill := readySkill(baseCard(t, "3121106"), 0)
		waterSkill := readySkill(baseCard(t, "3221107"), 0)
		behavior := Card1121114LegionGeneral{}

		if !cardHasActivePrayer(general) {
			t.Fatal("1121114 should expose prayer ability")
		}
		if err := behavior.OnPerTurn(&EffectContext{Engine: engine, Source: general, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1121114 prayer: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "legion_general_prayer" {
			t.Fatalf("1121114 should ask which fire spell buff to apply, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, "power")
		if got := engine.effectiveSpellPower(0, fireSkill, nil); got != fireSkill.Card.Power+2 {
			t.Fatalf("1121114 should give fire spells +2 power, got=%d", got)
		}
		if got := engine.effectiveSpellPower(0, waterSkill, nil); got != waterSkill.Card.Power {
			t.Fatalf("1121114 should not buff non-fire spell power, got=%d", got)
		}
		engine.State.TurnNumber += 2
		engine.clearExpiredTemporaryModifiers(0)
		if got := engine.effectiveSpellPower(0, fireSkill, nil); got != fireSkill.Card.Power {
			t.Fatalf("1121114 power buff should expire by next turn end, got=%d modifiers=%v", got, engine.State.Players[0].TempModifiers)
		}

		attackEngine := setupEffectTest(t)
		attackGeneral := placeUnit(baseCard(t, "1121114"), 0, 0, 0, attackEngine)
		attackSkill := readySkill(baseCard(t, "3121106"), 0)
		if err := behavior.OnPerTurn(&EffectContext{Engine: attackEngine, Source: attackGeneral, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1121114 attack prayer: %v", err)
		}
		resolvePendingSelection(t, attackEngine, 0, "attack")
		if got := attackEngine.effectiveSpellDamage(0, attackSkill, attackSkill.Card.Attack, nil); got != attackSkill.Card.Attack+1 {
			t.Fatalf("1121114 should give fire spells +1 attack, got=%d", got)
		}
	})

	t.Run("andis gift grants shadow load then kills the unit at turn end", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		item := NewCardInstance(baseCard(t, "2621110"), 0, engine.State.TurnNumber)
		target := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		behavior := Card2621110AndisGift{}

		if err := behavior.OnUseItem(&EffectContext{Engine: engine, Source: item, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2621110 use: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "andis_gift_target" ||
			!candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) ||
			!candidateContains(engine.State.PendingAction.Candidates, p0.Hero.InstanceID) {
			t.Fatalf("2621110 should ask for a friendly unit including hero, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if effectiveElementsGain(target)[model.ElementShadow] != 2 || target.Statuses[andisGiftDoomedStatus] != engine.State.TurnNumber {
			t.Fatalf("2621110 should grant +2 shadow load and mark target, load=%v statuses=%v", effectiveElementsGain(target), target.Statuses)
		}

		engine.finishEndTurn(p0)
		if p0.Units[1][0] != nil || !containsCardInstance(p0.Graveyard, target) {
			t.Fatalf("2621110 target should die at turn end, unit=%v grave=%v", p0.Units[1][0], cardsToInfo(p0.Graveyard))
		}

		heroEngine := setupEffectTest(t)
		heroP0 := heroEngine.State.Players[0]
		heroItem := NewCardInstance(baseCard(t, "2621110"), 0, heroEngine.State.TurnNumber)
		if err := behavior.OnUseItem(&EffectContext{Engine: heroEngine, Source: heroItem, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2621110 hero use: %v", err)
		}
		resolvePendingSelection(t, heroEngine, 0, heroP0.Hero.InstanceID)
		heroEngine.finishEndTurn(heroP0)
		if heroP0.Hero.CurrentLife > 0 || heroEngine.State.Phase != PhaseGameOver || heroEngine.State.Winner != 1 {
			t.Fatalf("2621110 should kill selected hero and end the game, life=%d phase=%s winner=%d", heroP0.Hero.CurrentLife, heroEngine.State.Phase, heroEngine.State.Winner)
		}
	})

	t.Run("western chart grants pierce to the selected water spell while equipped", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		chart := NewCardInstance(baseCard(t, "2221108"), 0, engine.State.TurnNumber)
		waterSkill := readySkill(baseCard(t, "3221106"), 0)
		fireSkill := readySkill(baseCard(t, "3121109"), 0)
		p0.Equipment[0] = chart
		p0.Skills[0] = waterSkill
		p0.Skills[1] = fireSkill
		placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		backEnemy := placeUnit(baseCard(t, "1021002"), 1, 0, 2, engine)

		if err := (Card2221108WesternChart{}).OnEnter(&EffectContext{Engine: engine, Source: chart, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2221108 enter: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "western_chart_pierce_target" ||
			!candidateContains(engine.State.PendingAction.Candidates, waterSkill.InstanceID) ||
			candidateContains(engine.State.PendingAction.Candidates, fireSkill.InstanceID) {
			t.Fatalf("2221108 should ask for a water spell only, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, waterSkill.InstanceID)
		if !engine.skillHasPierce(0, waterSkill) {
			t.Fatalf("2221108 should grant pierce to selected water spell")
		}
		if info := engine.cardToInfoForPlayer(p0, waterSkill); info["has_pierce"] != true {
			t.Fatalf("2221108-granted pierce should be serialized, info=%v", info)
		}
		if err := engine.validateSpellTarget(0, waterSkill, SpellTarget{Type: "unit", Position: *backEnemy.Position}); err != nil {
			t.Fatalf("2221108 should let selected water spell target back row: %v", err)
		}

		p0.Equipment[0] = nil
		if engine.skillHasPierce(0, waterSkill) {
			t.Fatalf("2221108 pierce should stop when the chart leaves equipment")
		}
	})

	t.Run("skycarrier e2 prayer draws or recycles two air graveyard cards", func(t *testing.T) {
		drawEngine := setupEffectTest(t)
		drawP0 := drawEngine.State.Players[0]
		drawP0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1321101"), 0, drawEngine.State.TurnNumber),
			NewCardInstance(baseCard(t, "1321102"), 0, drawEngine.State.TurnNumber),
		}
		startHand := len(drawP0.Hand)
		carrier := placeUnit(baseCard(t, "1321101"), 0, 0, 0, drawEngine)
		behavior := Card1321101SkycarrierE2{}

		if !cardHasActivePrayer(carrier) {
			t.Fatal("1321101 should expose prayer ability")
		}
		if err := behavior.OnPerTurn(&EffectContext{Engine: drawEngine, Source: carrier, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1321101 draw prayer: %v", err)
		}
		if drawEngine.State.PendingAction == nil || drawEngine.State.PendingAction.Type != "skycarrier_e2_prayer" || !candidateContains(drawEngine.State.PendingAction.Candidates, "draw") {
			t.Fatalf("1321101 should offer draw prayer, pending=%+v", drawEngine.State.PendingAction)
		}
		resolvePendingSelection(t, drawEngine, 0, "draw")
		if len(drawP0.Hand) != startHand+2 || len(drawP0.Deck) != 0 {
			t.Fatalf("1321101 draw prayer should draw two cards, hand=%v deck=%v", cardsToInfo(drawP0.Hand), cardsToInfo(drawP0.Deck))
		}

		recycleEngine := setupEffectTest(t)
		recycleP0 := recycleEngine.State.Players[0]
		recycleCarrier := placeUnit(baseCard(t, "1321101"), 0, 0, 0, recycleEngine)
		airA := NewCardInstance(baseCard(t, "1321101"), 0, recycleEngine.State.TurnNumber)
		airB := NewCardInstance(baseCard(t, "1321102"), 0, recycleEngine.State.TurnNumber)
		nonAir := NewCardInstance(baseCard(t, "1021001"), 0, recycleEngine.State.TurnNumber)
		recycleP0.Graveyard = []*CardInstance{airA, nonAir, airB}
		if err := behavior.OnPerTurn(&EffectContext{Engine: recycleEngine, Source: recycleCarrier, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("1321101 recycle prayer: %v", err)
		}
		if !candidateContains(recycleEngine.State.PendingAction.Candidates, "recycle") {
			t.Fatalf("1321101 should offer recycle when two air graveyard cards exist, pending=%+v", recycleEngine.State.PendingAction)
		}
		resolvePendingSelection(t, recycleEngine, 0, "recycle")
		if recycleEngine.State.PendingAction == nil || recycleEngine.State.PendingAction.Type != "skycarrier_e2_recycle" {
			t.Fatalf("1321101 should ask which air graveyard cards to recycle, pending=%+v", recycleEngine.State.PendingAction)
		}
		resolvePendingSelection(t, recycleEngine, 0, airA.InstanceID, airB.InstanceID)
		if len(recycleP0.Graveyard) != 1 || recycleP0.Graveyard[0] != nonAir || !containsCardInstance(recycleP0.Deck, airA) || !containsCardInstance(recycleP0.Deck, airB) {
			t.Fatalf("1321101 should move selected air graveyard cards to deck, grave=%v deck=%v", cardsToInfo(recycleP0.Graveyard), cardsToInfo(recycleP0.Deck))
		}
	})

	t.Run("arcane shield grants shield at next turn start", func(t *testing.T) {
		engine := setupEffectTest(t)
		skill := readySkill(baseCard(t, "3021107"), 0)
		behavior := Card3021107ArcaneShield{}

		if err := behavior.OnSpellCast(&EffectContext{Engine: engine, Source: skill, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("3021107 cast: %v", err)
		}
		if engine.State.Players[0].Shield != 0 || len(engine.State.Players[0].TempModifiers) != 1 {
			t.Fatalf("3021107 should defer shield gain, shield=%d modifiers=%+v", engine.State.Players[0].Shield, engine.State.Players[0].TempModifiers)
		}

		engine.applyTurnStartTemporaryModifiers(engine.State.Players[0])
		if engine.State.Players[0].Shield != 1 || len(engine.State.Players[0].TempModifiers) != 0 {
			t.Fatalf("3021107 should grant shield once at turn start, shield=%d modifiers=%+v", engine.State.Players[0].Shield, engine.State.Players[0].TempModifiers)
		}
	})

	t.Run("flame flash gains fire on spell hit", func(t *testing.T) {
		engine := setupEffectTest(t)
		skill := readySkill(baseCard(t, "3121109"), 0)
		behavior := Card3121109FlameFlash{}

		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: skill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("3121109 hit: %v", err)
		}
		if got := engine.State.Players[0].Elements[model.ElementFire]; got != 3 {
			t.Fatalf("3121109 should gain 3 fire, got %d elements=%v", got, engine.State.Players[0].Elements)
		}
	})

	t.Run("water mirror wall gains shield only on successful defense", func(t *testing.T) {
		engine := setupEffectTest(t)
		skill := readySkill(baseCard(t, "3221103"), 0)
		behavior := Card3221103WaterMirrorWall{}

		if err := behavior.OnDefend(&EffectContext{Engine: engine, Source: skill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": false}}); err != nil {
			t.Fatalf("3221103 failed defense: %v", err)
		}
		if engine.State.Players[0].Shield != 0 {
			t.Fatalf("3221103 should not shield failed defense, shield=%d", engine.State.Players[0].Shield)
		}

		if err := behavior.OnDefend(&EffectContext{Engine: engine, Source: skill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": true}}); err != nil {
			t.Fatalf("3221103 successful defense: %v", err)
		}
		if engine.State.Players[0].Shield != 1 {
			t.Fatalf("3221103 should gain 1 shield after successful defense, shield=%d", engine.State.Players[0].Shield)
		}
	})

	t.Run("summon defense spells damage enemies only after successful defense", func(t *testing.T) {
		engine := setupEffectTest(t)
		fireSnake := readySkill(baseCard(t, "3121101"), 0)
		target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		startLife := target.CurrentLife
		behavior := Card3121101SummonFireSnake{}

		if err := behavior.OnDefend(&EffectContext{Engine: engine, Source: fireSnake, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": false}}); err != nil {
			t.Fatalf("3121101 failed defense: %v", err)
		}
		if engine.State.PendingAction != nil || target.CurrentLife != startLife {
			t.Fatalf("3121101 should do nothing on failed defense, pending=%+v life=%d", engine.State.PendingAction, target.CurrentLife)
		}
		if err := behavior.OnDefend(&EffectContext{Engine: engine, Source: fireSnake, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": true}}); err != nil {
			t.Fatalf("3121101 successful defense: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "summon_fire_snake_defense_damage" || !candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) {
			t.Fatalf("3121101 should ask for an in-range enemy target, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, target.InstanceID)
		if target.CurrentLife != startLife-1 {
			t.Fatalf("3121101 should deal 1 damage to selected target, life=%d start=%d", target.CurrentLife, startLife)
		}

		houndEngine := setupEffectTest(t)
		hound := NewCardInstance(baseCard(t, "2121109"), 0, houndEngine.State.TurnNumber)
		houndTarget := placeUnit(baseCard(t, "1021002"), 1, 0, 0, houndEngine)
		houndStart := houndTarget.CurrentLife
		if err := (Card2121109SummonBlazingHoundScroll{}).OnDefend(&EffectContext{Engine: houndEngine, Source: hound, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": true}}); err != nil {
			t.Fatalf("2121109 successful defense: %v", err)
		}
		resolvePendingSelection(t, houndEngine, 0, houndTarget.InstanceID)
		if houndTarget.CurrentLife != houndStart-2 {
			t.Fatalf("2121109 should deal 2 damage to selected target, life=%d start=%d", houndTarget.CurrentLife, houndStart)
		}

		dragonEngine := setupEffectTest(t)
		dragon := readySkill(baseCard(t, "3221102"), 0)
		first := placeUnit(baseCard(t, "1021001"), 1, 0, 0, dragonEngine)
		second := placeUnit(baseCard(t, "1021002"), 1, 1, 0, dragonEngine)
		firstStart := first.CurrentLife
		secondStart := second.CurrentLife
		if err := (Card3221102SummonFloodDragon{}).OnDefend(&EffectContext{Engine: dragonEngine, Source: dragon, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": true}}); err != nil {
			t.Fatalf("3221102 successful defense: %v", err)
		}
		if first.CurrentLife != firstStart-1 || second.CurrentLife != secondStart-1 {
			t.Fatalf("3221102 should deal 1 damage to each in-range enemy, first=%d/%d second=%d/%d", first.CurrentLife, firstStart, second.CurrentLife, secondStart)
		}
	})

	t.Run("lion guardian permanently buffs other fire spells after successful defense", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		guardian := readySkill(baseCard(t, "3121102"), 0)
		fireSkill := readySkill(baseCard(t, "3121109"), 0)
		otherFireSkill := readySkill(baseCard(t, "3121007"), 0)
		nonFireSkill := readySkill(baseCard(t, "3321106"), 0)
		p0.Skills[0] = guardian
		p0.Skills[1] = fireSkill
		p0.Skills[2] = otherFireSkill
		p0.Skills[3] = nonFireSkill
		behavior := Card3121102LionGuardian{}

		if err := behavior.OnDefend(&EffectContext{Engine: engine, Source: guardian, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": false}}); err != nil {
			t.Fatalf("3121102 failed defense: %v", err)
		}
		if fireSkill.PowerBonus != 0 || otherFireSkill.PowerBonus != 0 || guardian.PowerBonus != 0 {
			t.Fatalf("3121102 should not buff on failed defense, guardian=%d fire=%d other=%d", guardian.PowerBonus, fireSkill.PowerBonus, otherFireSkill.PowerBonus)
		}

		if err := behavior.OnDefend(&EffectContext{Engine: engine, Source: guardian, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": true}}); err != nil {
			t.Fatalf("3121102 successful defense: %v", err)
		}
		if fireSkill.PowerBonus != 1 || otherFireSkill.PowerBonus != 1 || guardian.PowerBonus != 0 || nonFireSkill.PowerBonus != 0 {
			t.Fatalf("3121102 should buff other fire spells only, guardian=%d fire=%d other=%d nonfire=%d", guardian.PowerBonus, fireSkill.PowerBonus, otherFireSkill.PowerBonus, nonFireSkill.PowerBonus)
		}
	})

	t.Run("gather momentum buffs next attacking spell after successful defense", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		defense := readySkill(baseCard(t, "3321104"), 0)
		attackSpell := readySkill(baseCard(t, "3321106"), 0)
		behavior := Card3321104GatherMomentum{}

		if err := behavior.OnDefend(&EffectContext{Engine: engine, Source: defense, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": false}}); err != nil {
			t.Fatalf("3321104 failed defense: %v", err)
		}
		if len(p0.TempModifiers) != 0 {
			t.Fatalf("3321104 should not add modifier on failed defense, modifiers=%+v", p0.TempModifiers)
		}

		if err := behavior.OnDefend(&EffectContext{Engine: engine, Source: defense, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"defense_success": true}}); err != nil {
			t.Fatalf("3321104 successful defense: %v", err)
		}
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextAttackSpellPowerBonus || p0.TempModifiers[0].Amount != 3 || p0.TempModifiers[0].RemainingUses != 1 {
			t.Fatalf("3321104 should add one +3 power next-use modifier, modifiers=%+v", p0.TempModifiers)
		}
		if got := engine.temporarySpellPowerBonusForPurpose(0, attackSpell, skillPurposeDefend); got != 0 {
			t.Fatalf("3321104 temporary power bonus should not apply to defense, got %d", got)
		}
		if got := engine.temporarySpellPowerBonusForPurpose(0, attackSpell, skillPurposeAttack); got != 3 {
			t.Fatalf("3321104 temporary power bonus = %d, want 3", got)
		}
		engine.consumeNextSpellPowerBonuses(p0, attackSpell)
		if len(p0.TempModifiers) != 0 {
			t.Fatalf("3321104 modifier should be consumed after next attacking spell, modifiers=%+v", p0.TempModifiers)
		}
	})

	t.Run("corrosive flow discards a random enemy hand card on hit", func(t *testing.T) {
		engine := setupEffectTest(t)
		p1 := engine.State.Players[1]
		p1.Hand = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 1, engine.State.TurnNumber),
			NewCardInstance(baseCard(t, "1021002"), 1, engine.State.TurnNumber),
		}
		skill := readySkill(baseCard(t, "3221105"), 0)
		behavior := Card3221105CorrosiveFlow{}

		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: skill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("3221105 hit: %v", err)
		}
		if len(p1.Hand) != 1 || len(p1.Graveyard) != 1 {
			t.Fatalf("3221105 should discard one enemy hand card, hand=%d grave=%d", len(p1.Hand), len(p1.Graveyard))
		}
	})

	t.Run("plundering tide discards and draws for each hit unit", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		startHand := len(p0.Hand)
		p1.Hand = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 1, engine.State.TurnNumber),
			NewCardInstance(baseCard(t, "1021002"), 1, engine.State.TurnNumber),
			NewCardInstance(baseCard(t, "1021004"), 1, engine.State.TurnNumber),
		}
		unitA := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		unitB := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
		skill := readySkill(baseCard(t, "3221110"), 0)
		behavior := Card3221110PlunderingTide{}

		if err := behavior.OnSpellHitBeforeDamage(&EffectContext{
			Engine:     engine,
			Source:     skill,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"affected_units": []*CardInstance{unitA, unitB}, "attacker": 0, "spell_source": skill},
		}); err != nil {
			t.Fatalf("3221110 hit before damage: %v", err)
		}
		if len(p1.Hand) != 1 || len(p1.Graveyard) != 2 {
			t.Fatalf("3221110 should discard one enemy hand card per hit unit, hand=%d grave=%d", len(p1.Hand), len(p1.Graveyard))
		}
		if len(p0.Hand) != startHand+2 {
			t.Fatalf("3221110 should draw one card per hit unit, hand=%d start=%d", len(p0.Hand), startHand)
		}
	})

	t.Run("petrifying death ray applies petrify three", func(t *testing.T) {
		engine := setupEffectTest(t)
		skill := readySkill(baseCard(t, "3421109"), 0)
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if !skillNeedsTargetInstance(skill) || traitsForCardNumber("3421109").statuses[StatusPetrify] != 3 {
			t.Fatalf("3421109 should target and carry petrify 3 traits")
		}
		engine.applyExplicitSpellHitStatuses(skill, target)
		if target.Statuses[StatusPetrify] != 3 {
			t.Fatalf("3421109 should apply petrify 3, statuses=%v", target.Statuses)
		}
	})

	t.Run("goshawk buffs a friendly air spell next use", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		goshawk := readySkill(baseCard(t, "3321108"), 0)
		airSkill := readySkill(baseCard(t, "3321106"), 0)
		p0.Skills[0] = goshawk
		p0.Skills[1] = airSkill
		behavior := Card3321108CallSpiritGoshawk{}

		if err := behavior.OnEnter(&EffectContext{Engine: engine, Source: goshawk, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("3321108 enter: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "goshawk_air_skill_buff" {
			t.Fatalf("3321108 should ask for an air skill target, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, airSkill.InstanceID)
		if len(p0.TempModifiers) != 2 {
			t.Fatalf("3321108 should add two next-use modifiers, modifiers=%+v", p0.TempModifiers)
		}
		if p0.TempModifiers[0].TargetInstanceID != airSkill.InstanceID || p0.TempModifiers[1].TargetInstanceID != airSkill.InstanceID {
			t.Fatalf("3321108 modifiers should target selected skill, modifiers=%+v", p0.TempModifiers)
		}
		if p0.TempModifiers[0].Type != TempModSkillPowerBonus || p0.TempModifiers[1].Type != TempModNextSkillUseAttackBonus {
			t.Fatalf("3321108 should grant +1 power and +1 attack next use, modifiers=%+v", p0.TempModifiers)
		}
	})

	t.Run("air flow triggers on learned skill enter and hastes only next air spell", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		for i := range p0.Skills {
			p0.Skills[i] = nil
		}
		airFlow := NewCardInstance(baseCard(t, "3321110"), 0, engine.State.TurnNumber)
		fireSkill := NewCardInstance(baseCard(t, "3121109"), 0, engine.State.TurnNumber)
		airSkill := NewCardInstance(baseCard(t, "3321106"), 0, engine.State.TurnNumber)
		p0.SkillPool = []*CardInstance{airFlow, fireSkill, airSkill}
		p0.Elements = map[string]int{model.ElementAir: 20, model.ElementFire: 20}

		if err := engine.handleLearnSkill(0, ActionMessage{Action: "learn_skill", Data: map[string]any{"instance_id": airFlow.InstanceID}}); err != nil {
			t.Fatalf("learn 3321110: %v", err)
		}
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextLearnedSkillHaste || p0.TempModifiers[0].Element != model.ElementAir {
			t.Fatalf("3321110 should add air-only next learned haste, modifiers=%+v", p0.TempModifiers)
		}

		if err := engine.handleLearnSkill(0, ActionMessage{Action: "learn_skill", Data: map[string]any{"instance_id": fireSkill.InstanceID}}); err != nil {
			t.Fatalf("learn non-air skill: %v", err)
		}
		if !fireSkill.IsHorizontal || len(p0.TempModifiers) != 1 {
			t.Fatalf("3321110 should not haste or consume on non-air skill, horizontal=%v modifiers=%+v", fireSkill.IsHorizontal, p0.TempModifiers)
		}

		if err := engine.handleLearnSkill(0, ActionMessage{Action: "learn_skill", Data: map[string]any{"instance_id": airSkill.InstanceID}}); err != nil {
			t.Fatalf("learn air skill: %v", err)
		}
		if airSkill.IsHorizontal || len(p0.TempModifiers) != 0 {
			t.Fatalf("3321110 should haste next learned air skill once, horizontal=%v modifiers=%+v", airSkill.IsHorizontal, p0.TempModifiers)
		}
	})

	t.Run("aging touch removes all companion load on hit", func(t *testing.T) {
		engine := setupEffectTest(t)
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		setElementsGain(target, map[string]int{model.ElementFire: 2, model.ElementEarth: 1})
		target.ElementsGainBonus = map[string]int{model.ElementWater: 1}
		skill := readySkill(baseCard(t, "3421105"), 0)
		behavior := Card3421105AgingTouch{}

		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: skill, Target: target, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("3421105 hit: %v", err)
		}
		if totalLoad(target) != 0 {
			t.Fatalf("3421105 should remove all target load, load=%v bonus=%v", effectiveElementsGain(target), target.ElementsGainBonus)
		}

		hero := engine.State.Players[1].Hero
		setElementsGain(hero, map[string]int{model.ElementFire: 2})
		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: skill, Target: hero, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("3421105 hero hit: %v", err)
		}
		if totalLoad(hero) != 2 {
			t.Fatalf("3421105 should not remove hero load, load=%v", effectiveElementsGain(hero))
		}
	})

	t.Run("light spirit drain grants light load to a friendly light companion on hit", func(t *testing.T) {
		engine := setupEffectTest(t)
		lightA := placeUnit(baseCard(t, "1521104"), 0, 0, 0, engine)
		skill := readySkill(baseCard(t, "3521110"), 0)
		behavior := Card3521110LightSpiritDrain{}

		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: skill, Target: engine.State.Players[1].Hero, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("3521110 single target hit: %v", err)
		}
		if effectiveElementsGain(lightA)[model.ElementLight] != lightA.Card.ElementsGain[model.ElementLight]+1 {
			t.Fatalf("3521110 should auto-load sole friendly light companion, load=%v", effectiveElementsGain(lightA))
		}

		multiEngine := setupEffectTest(t)
		lightB := placeUnit(baseCard(t, "1521104"), 0, 0, 0, multiEngine)
		lightC := placeUnit(baseCard(t, "1521104"), 0, 1, 0, multiEngine)
		nonLight := placeUnit(baseCard(t, "1021001"), 0, 2, 0, multiEngine)
		multiSkill := readySkill(baseCard(t, "3521110"), 0)
		if err := behavior.OnSpellHit(&EffectContext{Engine: multiEngine, Source: multiSkill, Target: multiEngine.State.Players[1].Hero, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("3521110 multi target hit: %v", err)
		}
		if multiEngine.State.PendingAction == nil || multiEngine.State.PendingAction.Type != "light_spirit_drain_load" || len(multiEngine.State.PendingAction.Candidates) != 2 {
			t.Fatalf("3521110 should ask between light companion targets, pending=%+v", multiEngine.State.PendingAction)
		}
		resolvePendingSelection(t, multiEngine, 0, lightC.InstanceID)
		if effectiveElementsGain(lightC)[model.ElementLight] != lightC.Card.ElementsGain[model.ElementLight]+1 {
			t.Fatalf("3521110 should load selected light companion, load=%v", effectiveElementsGain(lightC))
		}
		if effectiveElementsGain(lightB)[model.ElementLight] != lightB.Card.ElementsGain[model.ElementLight] || effectiveElementsGain(nonLight)[model.ElementLight] != nonLight.Card.ElementsGain[model.ElementLight] {
			t.Fatalf("3521110 should not load unselected or non-light companions, lightB=%v nonLight=%v", effectiveElementsGain(lightB), effectiveElementsGain(nonLight))
		}
	})

	t.Run("blood soul slash hurts and heals own hero", func(t *testing.T) {
		engine := setupEffectTest(t)
		hero := engine.State.Players[0].Hero
		hero.CurrentLife = maxLife(hero) - 1
		skill := readySkill(baseCard(t, "3621103"), 0)
		behavior := Card3621103BloodSoulSlash{}

		if err := behavior.OnSpellCast(&EffectContext{Engine: engine, Source: skill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
			t.Fatalf("3621103 cast: %v", err)
		}
		if hero.CurrentLife != maxLife(hero)-2 {
			t.Fatalf("3621103 should deal 1 damage to own hero on attack cast, life=%d max=%d", hero.CurrentLife, maxLife(hero))
		}
		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: skill, Target: engine.State.Players[1].Hero, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("3621103 hit: %v", err)
		}
		if hero.CurrentLife != maxLife(hero) {
			t.Fatalf("3621103 should heal own hero by 2 on hit, life=%d max=%d", hero.CurrentLife, maxLife(hero))
		}
	})

	t.Run("blood pledge rewards damaging a friendly unit", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		friend := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		enemy := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		skill := readySkill(baseCard(t, "3621101"), 0)
		behavior := Card3621101BloodPledge{}

		if err := behavior.OnSpellHit(&EffectContext{
			Engine:     engine,
			Source:     skill,
			Target:     enemy,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"actual_friendly_damage_by_instance": map[string]int{}, "attacker": 0, "spell_source": skill},
		}); err != nil {
			t.Fatalf("3621101 enemy hit: %v", err)
		}
		if p0.Elements[model.ElementShadow] != 0 || len(p0.TempModifiers) != 0 {
			t.Fatalf("3621101 should ignore enemy targets, elements=%v modifiers=%+v", p0.Elements, p0.TempModifiers)
		}

		if err := behavior.OnSpellHit(&EffectContext{
			Engine:     engine,
			Source:     skill,
			Target:     friend,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"actual_friendly_damage_by_instance": map[string]int{friend.InstanceID: 1}, "attacker": 0, "spell_source": skill},
		}); err != nil {
			t.Fatalf("3621101 friendly hit: %v", err)
		}
		if p0.Elements[model.ElementShadow] != 2 {
			t.Fatalf("3621101 should gain 2 shadow after damaging friendly unit, elements=%v", p0.Elements)
		}
		if len(p0.TempModifiers) != 2 || p0.TempModifiers[0].TargetInstanceID != "" || p0.TempModifiers[1].TargetInstanceID != "" {
			t.Fatalf("3621101 should add two next spell modifiers, modifiers=%+v", p0.TempModifiers)
		}
		if p0.TempModifiers[0].Type != TempModNextAttackSpellPowerBonus || p0.TempModifiers[0].Amount != 2 || p0.TempModifiers[1].Type != TempModNextSkillUseAttackBonus || p0.TempModifiers[1].Amount != 1 {
			t.Fatalf("3621101 should grant +2 power and +1 attack to the next spell, modifiers=%+v", p0.TempModifiers)
		}
		nextSpell := readySkill(baseCard(t, "3621107"), 0)
		if power := engine.effectiveSpellPower(0, nextSpell, nil); power != nextSpell.Card.Power+2 {
			t.Fatalf("3621101 should grant +2 power to next spell, power=%d want=%d", power, nextSpell.Card.Power+2)
		}
		if damage := engine.effectiveSpellDamage(0, nextSpell, nextSpell.Card.Attack, nil); damage != nextSpell.Card.Attack+1 {
			t.Fatalf("3621101 should grant +1 attack to next spell, damage=%d want=%d", damage, nextSpell.Card.Attack+1)
		}
		engine.consumeNextSpellPowerBonuses(p0, nextSpell)
		engine.consumeNextSpellAttackBonuses(p0, nextSpell)
		if len(p0.TempModifiers) != 0 {
			t.Fatalf("3621101 next spell modifiers should be consumed together, modifiers=%+v", p0.TempModifiers)
		}

		defendedEngine := setupEffectTest(t)
		defendedP0 := defendedEngine.State.Players[0]
		defendedP1 := defendedEngine.State.Players[1]
		defendedFriend := placeUnit(baseCard(t, "1021001"), 0, 0, 0, defendedEngine)
		defendedEnemy := placeUnit(baseCard(t, "1021001"), 1, 1, 0, defendedEngine)
		defendedBloodPledge := readySkill(baseCard(t, "3621101"), 0)
		if err := behavior.OnSpellHit(&EffectContext{
			Engine:     defendedEngine,
			Source:     defendedBloodPledge,
			Target:     defendedFriend,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"actual_friendly_damage_by_instance": map[string]int{defendedFriend.InstanceID: 1}, "attacker": 0, "spell_source": defendedBloodPledge},
		}); err != nil {
			t.Fatalf("3621101 defended setup: %v", err)
		}
		defendedSpell := readySkill(baseCard(t, "3121002"), 0)
		defenseSpell := readySkill(baseCard(t, "3521013"), 1)
		defendedP0.Skills[0] = defendedSpell
		defendedP1.Skills[0] = defenseSpell
		defendedP0.Elements[model.ElementFire] = 10
		defendedP1.Elements[model.ElementLight] = 10
		if err := defendedEngine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": defendedSpell.InstanceID,
			"target_type": "unit",
			"target_col":  float64(defendedEnemy.Position.Col),
			"target_row":  float64(defendedEnemy.Position.Row),
		}}); err != nil {
			t.Fatalf("cast next spell after 3621101: %v", err)
		}
		if err := defendedEngine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{defenseSpell.InstanceID},
		}}); err != nil {
			t.Fatalf("defend next spell after 3621101: %v", err)
		}
		if len(defendedP0.TempModifiers) != 0 || defendedEnemy.CurrentLife != defendedEnemy.Card.Life {
			t.Fatalf("3621101 next spell modifiers should be consumed when the spell is defended, modifiers=%+v enemyLife=%d", defendedP0.TempModifiers, defendedEnemy.CurrentLife)
		}

		killedEngine := setupEffectTest(t)
		killedP0 := killedEngine.State.Players[0]
		killedSkill := readySkill(baseCard(t, "3621101"), 0)
		deadFriendlyID := "dead-friendly-unit"
		if err := behavior.OnSpellHit(&EffectContext{
			Engine:     killedEngine,
			Source:     killedSkill,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"actual_friendly_damage_by_instance": map[string]int{deadFriendlyID: 1}, "attacker": 0, "spell_source": killedSkill},
		}); err != nil {
			t.Fatalf("3621101 lethal friendly hit: %v", err)
		}
		if killedP0.Elements[model.ElementShadow] != 2 || len(killedP0.TempModifiers) != 2 {
			t.Fatalf("3621101 should reward actual friendly damage even after target leaves, elements=%v modifiers=%+v", killedP0.Elements, killedP0.TempModifiers)
		}

		lethalEngine := setupEffectTest(t)
		lethalP0 := lethalEngine.State.Players[0]
		for i := range lethalP0.Skills {
			lethalP0.Skills[i] = nil
		}
		lethalSkill := readySkill(baseCard(t, "3621101"), 0)
		lethalP0.Skills[0] = lethalSkill
		lethalFriend := placeUnit(baseCard(t, "1021001"), 0, 0, 0, lethalEngine)
		lethalFriend.CurrentLife = 1
		ownerID := 0
		lethalEngine.resolveSpellHit(0, lethalSkill, SpellTarget{Type: "unit", Position: *lethalFriend.Position, OwnerID: &ownerID}, nil, nil)
		if lethalP0.Elements[model.ElementShadow] != 2 || len(lethalP0.TempModifiers) != 2 {
			t.Fatalf("3621101 should reward real lethal friendly spell damage, elements=%v modifiers=%+v", lethalP0.Elements, lethalP0.TempModifiers)
		}
		if lethalEngine.State.Players[0].Units[0][0] != nil {
			t.Fatalf("3621101 test target should have died from real spell damage")
		}

		preventEngine := setupEffectTest(t)
		preventP0 := preventEngine.State.Players[0]
		for i := range preventP0.Skills {
			preventP0.Skills[i] = nil
		}
		preventSkill := readySkill(baseCard(t, "3621101"), 0)
		preventP0.Skills[0] = preventSkill
		preventFriend := placeUnit(baseCard(t, "1021001"), 0, 0, 0, preventEngine)
		preventFriend.Statuses[sturdyScrollShieldStatus] = 1
		preventFriend.Statuses[sturdyScrollShieldUntilStatus] = preventEngine.State.TurnNumber
		preventEngine.resolveSpellHit(0, preventSkill, SpellTarget{Type: "unit", Position: *preventFriend.Position, OwnerID: &ownerID}, nil, nil)
		if preventP0.Elements[model.ElementShadow] != 0 || len(preventP0.TempModifiers) != 0 {
			t.Fatalf("3621101 should not reward prevented friendly spell damage, elements=%v modifiers=%+v", preventP0.Elements, preventP0.TempModifiers)
		}
	})
}

func TestRoyalConflictEndlessWindTideAndDesertLeggings(t *testing.T) {
	t.Run("endless wind tide returns to hand and grows only this instance", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2321106"), 0, engine.State.TurnNumber)
		p0.Hand = append(p0.Hand, scroll)
		setAllElements(p0, 9)
		baseCost := engine.effectiveCardPlayCost(p0, scroll)[model.ElementAir]
		target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id":  scroll.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(target.Position.Col),
			"target_row":   float64(target.Position.Row),
			"target_owner": float64(1),
		}}); err != nil {
			t.Fatalf("2321106 use item: %v", err)
		}
		if !containsCardInstance(p0.Graveyard, scroll) || containsCardInstance(p0.Hand, scroll) {
			t.Fatalf("2321106 should wait in graveyard before hit resolves, hand=%v grave=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard))
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
			t.Fatalf("2321106 no defend: %v", err)
		}
		if containsCardInstance(p0.Graveyard, scroll) || len(p0.Hand) == 0 || p0.Hand[len(p0.Hand)-1] != scroll {
			t.Fatalf("2321106 should move from graveyard to hand after hit, hand=%v grave=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard))
		}
		if scroll.PowerBonus != 2 {
			t.Fatalf("2321106 should permanently gain +2 power, power_bonus=%d", scroll.PowerBonus)
		}
		if got := engine.effectiveCardPlayCost(p0, scroll)[model.ElementAir]; got != baseCost+1 {
			t.Fatalf("2321106 returned instance should cost +1 air, got=%d base=%d cost=%v", got, baseCost, engine.effectiveCardPlayCost(p0, scroll))
		}
		fresh := NewCardInstance(baseCard(t, "2321106"), 0, engine.State.TurnNumber)
		if got := engine.effectiveCardPlayCost(p0, fresh)[model.ElementAir]; got != baseCost {
			t.Fatalf("2321106 should not mutate global card cost, fresh=%d base=%d", got, baseCost)
		}
		scroll.IsHorizontal = false
		if err := engine.validateSkillForPurpose(scroll, skillPurposeAttack); err == nil {
			t.Fatalf("2321106 should not be usable again in the same turn")
		}
	})

	t.Run("desert leggings reduces large friendly companion damage", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		leggings := NewCardInstance(baseCard(t, "2421111"), 0, engine.State.TurnNumber)
		p0.Equipment[0] = leggings
		ally := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
		ally.CurrentLife = 5
		engine.dealDamageWithExtra(ally, 3, 0, map[string]any{"attacker": 1})
		if ally.CurrentLife != 4 || !leggings.UltimateUsed {
			t.Fatalf("2421111 should reduce first 3 friendly companion damage to 1 and spend ultimate, life=%d used=%v", ally.CurrentLife, leggings.UltimateUsed)
		}

		ally.CurrentLife = 5
		engine.dealDamageWithExtra(ally, 3, 0, map[string]any{"attacker": 1})
		if ally.CurrentLife != 2 {
			t.Fatalf("2421111 should not reduce damage again after ultimate is spent, life=%d", ally.CurrentLife)
		}

		small := placeUnit(baseCard(t, "1021002"), 0, 1, 0, engine)
		small.CurrentLife = 5
		engine.dealDamageWithExtra(small, 1, 0, map[string]any{"attacker": 1})
		if small.CurrentLife != 4 {
			t.Fatalf("2421111 should not reduce 1 damage, life=%d", small.CurrentLife)
		}

		enemy := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		enemy.CurrentLife = 5
		engine.dealDamageWithExtra(enemy, 3, 1, map[string]any{"attacker": 0})
		if enemy.CurrentLife != 2 {
			t.Fatalf("2421111 should not reduce enemy unit damage, life=%d", enemy.CurrentLife)
		}
	})
}

func TestRoyalConflictPrimalDivineFlameLopsius(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	lopsius := readySkill(baseCard(t, "3111102"), 0)
	fireSpell := readySkill(baseCard(t, "3121001"), 0)
	waterSpell := readySkill(baseCard(t, "3221001"), 0)
	p0.Skills[0] = lopsius
	p0.Skills[1] = fireSpell
	p0.Skills[2] = waterSpell

	behavior := Card3111102PrimalDivineFlameLopsius{}
	if err := behavior.ValidateSkillUse(&EffectContext{Engine: engine, Source: lopsius, PlayerID: 0, OpponentID: 1}, lopsius, skillPurposeAttackBoost); err == nil {
		t.Fatalf("3111102 should not be usable as a boost")
	}
	if err := behavior.OnPerTurn(&EffectContext{Engine: engine, Source: lopsius, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("3111102 per-turn failed: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "primal_divine_flame_exile" || !candidateContains(engine.State.PendingAction.Candidates, fireSpell.InstanceID) || candidateContains(engine.State.PendingAction.Candidates, waterSpell.InstanceID) {
		t.Fatalf("3111102 should offer other fire skills only, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, fireSpell.InstanceID)
	if p0.Skills[1] != nil || len(p0.Exile) != 1 || p0.Exile[0] != fireSpell {
		t.Fatalf("3111102 should exile selected fire skill, skills=%v exile=%v", cardsToInfo(p0.Skills[:]), cardsToInfo(p0.Exile))
	}
	if lopsius.AttackBonus != 1 || lopsius.PowerBonus != 2 || lopsius.UsedThisTurn != 1 {
		t.Fatalf("3111102 should gain +1 attack +2 power and spend use, attack=%d power=%d used=%d", lopsius.AttackBonus, lopsius.PowerBonus, lopsius.UsedThisTurn)
	}
	lopsius.UsedThisTurn = 0
	secondFireSpell := readySkill(baseCard(t, "3121002"), 0)
	p0.Skills[1] = secondFireSpell
	if err := behavior.OnPerTurn(&EffectContext{Engine: engine, Source: lopsius, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("3111102 second per-turn failed: %v", err)
	}
	resolvePendingSelection(t, engine, 0, secondFireSpell.InstanceID)
	info := engine.cardToInfoForPlayer(p0, lopsius)
	if info["attack"] != 4 || info["current_attack"] != 4 || info["power"] != 10 {
		t.Fatalf("3111102 serialized stats should include repeated growth, info=%v", info)
	}

	actionEngine := setupEffectTest(t)
	actionP0 := actionEngine.State.Players[0]
	actionLopsius := readySkill(baseCard(t, "3111102"), 0)
	actionFireSpell := readySkill(baseCard(t, "3121103"), 0)
	actionP0.Skills[0] = actionLopsius
	actionP0.Skills[1] = actionFireSpell
	if err := actionEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  actionLopsius.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("3111102 handle action per-turn failed: %v", err)
	}
	if actionLopsius.UsedThisTurn != 1 {
		t.Fatalf("3111102 should be marked used after opening pending action, used=%d", actionLopsius.UsedThisTurn)
	}
	resolvePendingSelection(t, actionEngine, 0, actionFireSpell.InstanceID)
	if actionP0.Skills[1] != nil || len(actionP0.Exile) != 1 || actionP0.Exile[0] != actionFireSpell {
		t.Fatalf("3111102 handle action path should exile selected fire skill, skills=%v exile=%v", cardsToInfo(actionP0.Skills[:]), cardsToInfo(actionP0.Exile))
	}
	if actionLopsius.AttackBonus != 1 || actionLopsius.PowerBonus != 2 || actionLopsius.UsedThisTurn != 1 {
		t.Fatalf("3111102 handle action path should grow once without double-counting use, attack=%d power=%d used=%d", actionLopsius.AttackBonus, actionLopsius.PowerBonus, actionLopsius.UsedThisTurn)
	}
}

func TestRoyalConflictMindSeaMazeStacksAfterUse(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	skill := readySkill(baseCard(t, "3211101"), 0)
	p0.Skills[0] = skill
	behavior := Card3211101MindSeaMaze{}
	ctx := &EffectContext{
		Engine:     engine,
		Source:     skill,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"cast_player": 0},
	}
	if got := engine.effectiveSpellArea(skill); got != SpellAreaSingle {
		t.Fatalf("3211101 should start single-target, area=%s", got)
	}
	if err := behavior.OnSpellCast(ctx); err != nil {
		t.Fatalf("3211101 first cast: %v", err)
	}
	if err := behavior.OnSpellCast(ctx); err != nil {
		t.Fatalf("3211101 second cast: %v", err)
	}
	if got := engine.temporarySpellPowerBonus(0, skill); got != 2 {
		t.Fatalf("3211101 should gain +1 power per use this turn, got=%d modifiers=%+v", got, p0.TempModifiers)
	}
	if got := engine.effectiveSpellArea(skill); got != SpellAreaAll {
		t.Fatalf("3211101 should become any/all range this turn, area=%s", got)
	}
	skill.Statuses[mindSeaMazeAnyRangeUntilStatus] = engine.State.TurnNumber - 1
	if got := engine.effectiveSpellArea(skill); got != SpellAreaSingle {
		t.Fatalf("3211101 range status should expire by turn, area=%s", got)
	}
}

func TestRoyalConflictTreadingWaveResetsOtherWaterSpellWithRisingCost(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	treading := readySkill(baseCard(t, "3221101"), 0)
	usedWater := readySkill(baseCard(t, "3221001"), 0)
	resetTarget := readySkill(baseCard(t, "3221104"), 0)
	resetTarget.IsHorizontal = true
	p0.Skills[0] = treading
	p0.Skills[1] = usedWater
	p0.Skills[2] = resetTarget

	if err := (Card3221101TreadingWave{}).OnSpellCast(&EffectContext{
		Engine:     engine,
		Source:     treading,
		Target:     usedWater,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"cast_player": 0},
	}); err != nil {
		t.Fatalf("3221101 trigger failed: %v", err)
	}
	if resetTarget.IsHorizontal {
		t.Fatalf("3221101 should reset another horizontal water spell")
	}
	if got := engine.effectiveSkillUseCost(p0, resetTarget)[model.ElementWater]; got != resetTarget.Card.ElementsExpense[model.ElementWater]+2 {
		t.Fatalf("3221101 first trigger should add 2 water to next use, cost=%v", engine.effectiveSkillUseCost(p0, resetTarget))
	}
	engine.consumeNextSkillUseModifiers(p0, resetTarget)
	if got := engine.effectiveSkillUseCost(p0, resetTarget)[model.ElementWater]; got != resetTarget.Card.ElementsExpense[model.ElementWater] {
		t.Fatalf("3221101 extra cost should be consumed after next use, cost=%v statuses=%v", engine.effectiveSkillUseCost(p0, resetTarget), resetTarget.Statuses)
	}

	resetTarget.IsHorizontal = true
	if err := (Card3221101TreadingWave{}).OnSpellCast(&EffectContext{
		Engine:     engine,
		Source:     treading,
		Target:     usedWater,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"cast_player": 0},
	}); err != nil {
		t.Fatalf("3221101 second trigger failed: %v", err)
	}
	if got := engine.effectiveSkillUseCost(p0, resetTarget)[model.ElementWater]; got != resetTarget.Card.ElementsExpense[model.ElementWater]+3 {
		t.Fatalf("3221101 second same-turn trigger should add 3 water, cost=%v", engine.effectiveSkillUseCost(p0, resetTarget))
	}

	enemyEngine := setupEffectTest(t)
	enemyTreading := readySkill(baseCard(t, "3221101"), 0)
	enemyTarget := readySkill(baseCard(t, "3221104"), 0)
	enemyTarget.IsHorizontal = true
	enemyEngine.State.Players[0].Skills[0] = enemyTreading
	enemyEngine.State.Players[0].Skills[1] = enemyTarget
	if err := (Card3221101TreadingWave{}).OnSpellCast(&EffectContext{
		Engine:     enemyEngine,
		Source:     enemyTreading,
		Target:     readySkill(baseCard(t, "3221001"), 1),
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"cast_player": 1},
	}); err != nil {
		t.Fatalf("3221101 enemy spell trigger failed: %v", err)
	}
	if !enemyTarget.IsHorizontal || enemyTarget.Statuses[skillUseExtraCostStatus(model.ElementWater, 2)] > 0 {
		t.Fatalf("3221101 should ignore enemy water spells, horizontal=%v statuses=%v", enemyTarget.IsHorizontal, enemyTarget.Statuses)
	}
}

func TestRoyalConflictHolyChildBonusOnSelfLoadGain(t *testing.T) {
	engine := setupEffectTest(t)
	child := placeUnit(baseCard(t, "1521102"), 0, 0, 0, engine)
	engine.addElementsGainBonus(child, 0, model.ElementLight, 1, child)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "holy_child_bonus" {
		t.Fatalf("1521102 should prompt after gaining load, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, "gain_light_load")
	if got := effectiveElementsGain(child)[model.ElementLight]; got != child.Card.ElementsGain[model.ElementLight]+2 {
		t.Fatalf("1521102 should gain original light load plus one bonus light load, got=%d load=%v", got, effectiveElementsGain(child))
	}
	if !child.UltimateUsed {
		t.Fatalf("1521102 should spend its triggered ultimate after bonus load")
	}
	if engine.State.PendingAction != nil {
		t.Fatalf("1521102 bonus load should not recursively reopen prompt, pending=%+v", engine.State.PendingAction)
	}

	lifeEngine := setupEffectTest(t)
	lifeChild := placeUnit(baseCard(t, "1521102"), 0, 0, 0, lifeEngine)
	startLife := lifeChild.CurrentLife
	lifeEngine.addElementsGainBonus(lifeChild, 0, model.ElementLight, 1, lifeChild)
	resolvePendingSelection(t, lifeEngine, 0, "gain_life")
	if lifeChild.CurrentLife != startLife+1 {
		t.Fatalf("1521102 should be able to choose +1 life, life=%d start=%d", lifeChild.CurrentLife, startLife)
	}

	lifeGainEngine := setupEffectTest(t)
	lifeGainChild := placeUnit(baseCard(t, "1521102"), 0, 0, 0, lifeGainEngine)
	lifeGainChild.CurrentLife++
	lifeGainEngine.triggerHolyChildAfterLifeGain(0, lifeGainChild)
	if lifeGainEngine.State.PendingAction == nil || lifeGainEngine.State.PendingAction.Type != "holy_child_bonus" {
		t.Fatalf("1521102 should prompt after gaining life, pending=%+v", lifeGainEngine.State.PendingAction)
	}
}

func TestRoyalConflictThunderChainGrantsNextDriveExtraTarget(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	chain := NewCardInstance(baseCard(t, "2321101"), 0, engine.State.TurnNumber)
	chain.IsHorizontal = false
	p0.Equipment[0] = chain
	drive := readySkill(baseCard(t, "3321002"), 0)
	p0.Skills[0] = drive
	setAllElements(p0, 9)
	mainTarget := placeUnit(baseCard(t, "1021001"), 1, 1, 1, engine)
	extraTarget := placeUnit(baseCard(t, "1021002"), 1, 2, 1, engine)
	mainLife := mainTarget.CurrentLife
	extraLife := extraTarget.CurrentLife

	if err := engine.HandleAction(0, ActionMessage{
		Action: "use_ability",
		Data: map[string]any{
			"instance_id":  chain.InstanceID,
			"ability_type": "per_turn",
		},
	}); err != nil {
		t.Fatalf("use 2321101 ability: %v", err)
	}
	if !chain.IsHorizontal || len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextDriveSpellExtraTarget {
		t.Fatalf("2321101 should tap and grant next drive extra target, horizontal=%v modifiers=%+v", chain.IsHorizontal, p0.TempModifiers)
	}

	if err := engine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id":      drive.InstanceID,
			"target_type":      "unit",
			"target_col":       float64(1),
			"target_row":       float64(1),
			"extra_target_col": float64(2),
			"extra_target_row": float64(1),
		},
	}); err != nil {
		t.Fatalf("cast drive spell with extra target: %v", err)
	}
	if engine.State.Phase == PhaseDefenseWindow {
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve defense: %v", err)
		}
	}
	if mainTarget.CurrentLife >= mainLife || extraTarget.CurrentLife >= extraLife {
		t.Fatalf("drive spell should damage main and extra target, main=%d/%d extra=%d/%d p1=%v", mainTarget.CurrentLife, mainLife, extraTarget.CurrentLife, extraLife, cardsToInfo(p1.Units[0][:]))
	}
	if p0.TempModifiers[0].RemainingUses != 0 {
		t.Fatalf("2321101 extra target modifier should be consumed, modifiers=%+v", p0.TempModifiers)
	}
}

func TestRoyalConflictSoulRendingScreamWeakensDefenseSpellsAfterBattle(t *testing.T) {
	engine := setupEffectTest(t)
	attacker := readySkill(baseCard(t, "3621109"), 0)
	defense := readySkill(baseCard(t, "3221001"), 1)
	boost := readySkill(baseCard(t, "3221002"), 1)
	if err := (Card3621109SoulRendingScream{}).OnDefend(&EffectContext{
		Engine:     engine,
		Source:     attacker,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData: map[string]any{
			"attack_skill":    attacker,
			"defense_skills":  []*CardInstance{defense},
			"defense_boosts":  []*CardInstance{boost},
			"defense_success": true,
		},
	}); err != nil {
		t.Fatalf("3621109 defend trigger failed: %v", err)
	}
	if defense.Statuses[StatusWeaken] != 1 || boost.Statuses[StatusWeaken] != 1 {
		t.Fatalf("3621109 should weaken enemy defense spell and boost, defense=%v boost=%v", defense.Statuses, boost.Statuses)
	}

	other := readySkill(baseCard(t, "3321002"), 0)
	if err := (Card3621109SoulRendingScream{}).OnDefend(&EffectContext{
		Engine:     engine,
		Source:     attacker,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData: map[string]any{
			"attack_skill":   other,
			"defense_skills": []*CardInstance{defense},
		},
	}); err != nil {
		t.Fatalf("3621109 unrelated defend trigger failed: %v", err)
	}
	if defense.Statuses[StatusWeaken] != 1 {
		t.Fatalf("3621109 should only trigger when it is the main attack spell, defense=%v", defense.Statuses)
	}
}

func TestRoyalConflictRedMoonDevourDestroysAndFeedsShadowUnit(t *testing.T) {
	engine := setupEffectTest(t)
	skill := readySkill(baseCard(t, "3621106"), 0)
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	target.CurrentLife = 3
	if err := (Card3621106RedMoonDevour{}).OnSpellHit(&EffectContext{
		Engine:     engine,
		Source:     skill,
		Target:     target,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"attacker": 0, "spell_source": skill},
	}); err != nil {
		t.Fatalf("3621106 hit without red moon: %v", err)
	}
	if engine.State.Players[1].Units[0][0] != nil || len(engine.State.Players[1].Graveyard) != 1 || engine.State.PendingAction != nil {
		t.Fatalf("3621106 should destroy target without prompt when red moon inactive, unit=%v grave=%v pending=%+v", engine.State.Players[1].Units[0][0], cardsToInfo(engine.State.Players[1].Graveyard), engine.State.PendingAction)
	}

	redEngine := setupEffectTest(t)
	redP0 := redEngine.State.Players[0]
	redMoon := readySkill(baseCard(t, "3611101"), 0)
	redMoon.Statuses[StatusAbilityDuration] = 1
	redP0.Skills[0] = redMoon
	redEngine.refreshRedMoonState(0)
	shadowAlly := placeUnit(baseCard(t, "1621001"), 0, 1, 0, redEngine)
	startLife := shadowAlly.CurrentLife
	redSkill := readySkill(baseCard(t, "3621106"), 0)
	redTarget := placeUnit(baseCard(t, "1021002"), 1, 0, 0, redEngine)
	redTarget.CurrentLife = 2
	if err := (Card3621106RedMoonDevour{}).OnSpellHit(&EffectContext{
		Engine:     redEngine,
		Source:     redSkill,
		Target:     redTarget,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"attacker": 0, "spell_source": redSkill},
	}); err != nil {
		t.Fatalf("3621106 hit with red moon: %v", err)
	}
	if redEngine.State.PendingAction == nil || redEngine.State.PendingAction.Type != "red_moon_devour_life" {
		t.Fatalf("3621106 should prompt shadow life gain when red moon active, pending=%+v", redEngine.State.PendingAction)
	}
	resolvePendingSelection(t, redEngine, 0, shadowAlly.InstanceID)
	if shadowAlly.CurrentLife != startLife+2 {
		t.Fatalf("3621106 should grant remaining life to selected shadow ally, life=%d start=%d", shadowAlly.CurrentLife, startLife)
	}

	multiEngine := setupEffectTest(t)
	multiP0 := multiEngine.State.Players[0]
	multiRedMoon := readySkill(baseCard(t, "3611101"), 0)
	multiRedMoon.Statuses[StatusAbilityDuration] = 1
	multiP0.Skills[0] = multiRedMoon
	multiEngine.refreshRedMoonState(0)
	multiShadowAlly := placeUnit(baseCard(t, "1621001"), 0, 1, 0, multiEngine)
	multiStartLife := multiShadowAlly.CurrentLife
	multiSkill := readySkill(baseCard(t, "3621106"), 0)
	mainTarget := placeUnit(baseCard(t, "1021002"), 1, 0, 0, multiEngine)
	extraTarget := placeUnit(baseCard(t, "1021001"), 1, 1, 0, multiEngine)
	mainTarget.CurrentLife = 2
	extraTarget.CurrentLife = 3
	if err := (Card3621106RedMoonDevour{}).OnSpellHit(&EffectContext{
		Engine:     multiEngine,
		Source:     multiSkill,
		Target:     mainTarget,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData: map[string]any{
			"attacker":       0,
			"spell_source":   multiSkill,
			"affected_units": []*CardInstance{mainTarget, extraTarget},
		},
	}); err != nil {
		t.Fatalf("3621106 multi hit with red moon: %v", err)
	}
	if multiEngine.State.Players[1].Units[0][0] != nil || multiEngine.State.Players[1].Units[1][0] != nil {
		t.Fatalf("3621106 should destroy all hit enemy companions, units=%v", multiEngine.State.Players[1].Units)
	}
	if multiEngine.State.PendingAction == nil || multiEngine.State.PendingAction.Type != "red_moon_devour_life" {
		t.Fatalf("3621106 multi hit should prompt shadow life gain, pending=%+v", multiEngine.State.PendingAction)
	}
	resolvePendingSelection(t, multiEngine, 0, multiShadowAlly.InstanceID)
	if multiShadowAlly.CurrentLife != multiStartLife+5 {
		t.Fatalf("3621106 should grant total remaining life from destroyed targets, life=%d start=%d", multiShadowAlly.CurrentLife, multiStartLife)
	}
}

func TestRoyalConflictMoonshadowResetsAfterWeakDefenseBattle(t *testing.T) {
	engine := setupEffectTest(t)
	moonshadow := readySkill(baseCard(t, "3621108"), 0)
	moonshadow.IsHorizontal = true
	weakDefense := readySkill(baseCard(t, "3221001"), 1)
	weakDefense.Statuses[StatusWeaken] = 1
	boost := readySkill(baseCard(t, "3221002"), 1)

	if err := (Card3621108Moonshadow{}).OnDefend(&EffectContext{
		Engine:     engine,
		Source:     moonshadow,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData: map[string]any{
			"attack_skill":    moonshadow,
			"defense_skills":  []*CardInstance{weakDefense},
			"defense_boosts":  []*CardInstance{boost},
			"defense_success": true,
		},
	}); err != nil {
		t.Fatalf("3621108 defend trigger failed: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "moonshadow_reset" {
		t.Fatalf("3621108 should offer optional reset after battling a weakened defense spell, pending=%+v", engine.State.PendingAction)
	}
	if !moonshadow.IsHorizontal {
		t.Fatalf("3621108 should not reset before the optional window resolves")
	}
	resolvePendingSelection(t, engine, 0, moonshadow.InstanceID)
	if moonshadow.IsHorizontal {
		t.Fatalf("3621108 should reset when accepted")
	}

	moonshadow.IsHorizontal = true
	otherAttack := readySkill(baseCard(t, "3621109"), 0)
	if err := (Card3621108Moonshadow{}).OnDefend(&EffectContext{
		Engine:     engine,
		Source:     moonshadow,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData: map[string]any{
			"attack_skill":   otherAttack,
			"defense_skills": []*CardInstance{weakDefense},
		},
	}); err != nil {
		t.Fatalf("3621108 unrelated defend trigger failed: %v", err)
	}
	if !moonshadow.IsHorizontal {
		t.Fatalf("3621108 should only reset when it is the attacking spell")
	}
}

func TestRoyalConflictRedeemerEveAutumnMapleUltimate(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	eve := NewCardInstance(baseCard(t, "4511102"), 0, engine.State.TurnNumber)
	eve.IsHorizontal = false
	p0.Hero = eve
	p0.Units[1][1] = eve
	woundedCompanion := placeUnit(baseCard(t, "1521001"), 0, 0, 1, engine)
	woundedCompanion.CurrentLife = maxLife(woundedCompanion) - 1
	woundedHero := p0.Hero
	woundedHero.CurrentLife = maxLife(woundedHero) - 1
	placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
	placeUnit(baseCard(t, "1021003"), 1, 2, 0, engine)
	p0.Elements[model.ElementLight] = 2
	startLoad := effectiveElementsGain(woundedCompanion)[model.ElementLight]

	if err := engine.HandleAction(0, ActionMessage{
		Action: "use_ability",
		Data: map[string]any{
			"instance_id":  eve.InstanceID,
			"ability_type": "ultimate",
		},
	}); err != nil {
		t.Fatalf("use 4511102 ultimate: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "redeemer_eve_autumn_maple_target" || !candidateContains(engine.State.PendingAction.Candidates, woundedCompanion.InstanceID) || candidateContains(engine.State.PendingAction.Candidates, woundedHero.InstanceID) {
		t.Fatalf("4511102 should ask for a wounded friendly companion only, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, woundedCompanion.InstanceID)
	if !eve.UltimateUsed || p0.Elements[model.ElementLight] != 0 {
		t.Fatalf("4511102 should spend ultimate and 2 light, ultimate=%v elements=%v", eve.UltimateUsed, p0.Elements)
	}
	if maxLife(woundedCompanion) != woundedCompanion.Card.Life+2 || woundedCompanion.CurrentLife != woundedCompanion.Card.Life+1 {
		t.Fatalf("4511102 should grant +2 max/current life, life=%d max=%d", woundedCompanion.CurrentLife, maxLife(woundedCompanion))
	}
	if effectiveElementsGain(woundedCompanion)[model.ElementLight] != startLoad+2 {
		t.Fatalf("4511102 should grant +2 light load, load=%v start=%d", effectiveElementsGain(woundedCompanion), startLoad)
	}

	p1.Units[0][0] = nil
	failEngine := setupEffectTest(t)
	failP0 := failEngine.State.Players[0]
	failEve := NewCardInstance(baseCard(t, "4511102"), 0, failEngine.State.TurnNumber)
	failP0.Hero = failEve
	failP0.Units[1][1] = failEve
	failTarget := placeUnit(baseCard(t, "1521001"), 0, 0, 0, failEngine)
	failTarget.CurrentLife = maxLife(failTarget) - 1
	placeUnit(baseCard(t, "1021001"), 1, 0, 0, failEngine)
	failP0.Elements[model.ElementLight] = 1
	if err := (Card4511102RedeemerEveAutumnMaple{}).OnUltimate(&EffectContext{Engine: failEngine, Source: failEve, PlayerID: 0, OpponentID: 1}); err == nil {
		t.Fatalf("4511102 should reject when enemy unit count is not higher")
	}
	if failEve.UltimateUsed || failEngine.State.PendingAction != nil {
		t.Fatalf("4511102 failed ultimate should not spend state, ultimate=%v pending=%+v", failEve.UltimateUsed, failEngine.State.PendingAction)
	}
}

func TestRoyalConflictFlowerSeaDreamWhaleShufflesAndSearchesDreams(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	whale := NewCardInstance(baseCard(t, "1211102"), 0, engine.State.TurnNumber)
	p0.Deck = nil
	if err := (Card1211102FlowerSeaDreamWhale{}).OnEnter(&EffectContext{Engine: engine, Source: whale, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("1211102 enter failed: %v", err)
	}
	if len(p0.Deck) != 3 {
		t.Fatalf("1211102 should shuffle three dreams into deck, deck=%v", cardsToInfo(p0.Deck))
	}
	seenDreams := map[string]bool{}
	for _, card := range p0.Deck {
		if !isDreamCreationCardInstance(card) {
			t.Fatalf("1211102 should only add dream cards, deck=%v", cardsToInfo(p0.Deck))
		}
		seenDreams[card.Card.Number] = true
	}
	for _, number := range []string{"2201101", "2201102", "2201103"} {
		if !seenDreams[number] {
			t.Fatalf("1211102 missing dream %s in deck=%v", number, cardsToInfo(p0.Deck))
		}
	}

	creation := readySkill(baseCard(t, "3121105"), 0)
	if err := (Card1211102FlowerSeaDreamWhale{}).OnSpellCast(&EffectContext{
		Engine:     engine,
		Source:     whale,
		Target:     creation,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"cast_player": 0},
	}); err != nil {
		t.Fatalf("1211102 first creation trigger failed: %v", err)
	}
	if engine.State.PendingAction != nil || whale.Statuses[flowerSeaDreamWhaleCreationCountStatus] != 1 {
		t.Fatalf("1211102 should wait for the second creation spell, pending=%+v statuses=%v", engine.State.PendingAction, whale.Statuses)
	}
	if err := (Card1211102FlowerSeaDreamWhale{}).OnSpellCast(&EffectContext{
		Engine:     engine,
		Source:     whale,
		Target:     creation,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"cast_player": 0},
	}); err != nil {
		t.Fatalf("1211102 second creation trigger failed: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "flower_sea_dream_whale_search" {
		t.Fatalf("1211102 should prompt to search a dream after two creation spells, pending=%+v", engine.State.PendingAction)
	}
	handBeforeSearch := len(p0.Hand)
	selected := engine.State.PendingAction.Candidates[0]["instance_id"].(string)
	resolvePendingSelection(t, engine, 0, selected)
	if len(p0.Hand) != handBeforeSearch+1 || !containsDreamCreationCard(p0.Hand) || whale.Statuses[flowerSeaDreamWhaleCreationCountStatus] != 0 {
		t.Fatalf("1211102 should search a dream and spend two-count, hand=%v statuses=%v", cardsToInfo(p0.Hand), whale.Statuses)
	}

	enemyCastEngine := setupEffectTest(t)
	enemyWhale := NewCardInstance(baseCard(t, "1211102"), 0, enemyCastEngine.State.TurnNumber)
	if err := (Card1211102FlowerSeaDreamWhale{}).OnSpellCast(&EffectContext{
		Engine:     enemyCastEngine,
		Source:     enemyWhale,
		Target:     readySkill(baseCard(t, "3121105"), 1),
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"cast_player": 1},
	}); err != nil {
		t.Fatalf("1211102 enemy creation trigger failed: %v", err)
	}
	if enemyWhale.Statuses[flowerSeaDreamWhaleCreationCountStatus] != 0 {
		t.Fatalf("1211102 should ignore enemy creation spells, statuses=%v", enemyWhale.Statuses)
	}
}

func containsDreamCreationCard(cards []*CardInstance) bool {
	for _, card := range cards {
		if isDreamCreationCardInstance(card) {
			return true
		}
	}
	return false
}

func TestRoyalConflictRottenAncientTreeHeartRemovesLoadEveryTwoSpells(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	heart := NewCardInstance(baseCard(t, "2411102"), 0, engine.State.TurnNumber)
	p0.Equipment[0] = heart
	heart.SlotIndex = 0
	ally := placeUnit(baseCard(t, "1421101"), 0, 0, 0, engine)
	spell := readySkill(baseCard(t, "3121105"), 0)

	trigger := func() {
		t.Helper()
		if err := (Card2411102RottenAncientTreeHeart{}).OnSpellCast(&EffectContext{
			Engine:     engine,
			Source:     heart,
			Target:     spell,
			PlayerID:   0,
			OpponentID: 1,
			ExtraData:  map[string]any{"cast_player": 0},
		}); err != nil {
			t.Fatalf("2411102 spell trigger failed: %v", err)
		}
	}
	trigger()
	if engine.State.PendingAction != nil || heart.Statuses[rottenAncientTreeHeartSpellCountPrefix+"0"] != 1 {
		t.Fatalf("2411102 should count first spell only, pending=%+v statuses=%v", engine.State.PendingAction, heart.Statuses)
	}
	trigger()
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "rotten_ancient_tree_heart_remove_load" {
		t.Fatalf("2411102 should ask current caster to remove load on second spell, pending=%+v", engine.State.PendingAction)
	}
	allyEarthBefore := reducibleElementLoad(ally, model.ElementEarth)
	resolvePendingSelection(t, engine, 0, ally.InstanceID+"|"+model.ElementEarth)
	if reducibleElementLoad(ally, model.ElementEarth) != allyEarthBefore-1 || heart.Statuses[rottenAncientTreeHeartSpellCountPrefix+"0"] != 0 {
		t.Fatalf("2411102 should remove selected earth load and spend count, ally_load=%v statuses=%v", effectiveElementsGain(ally), heart.Statuses)
	}

	p0.Equipment[0] = heart
	heart.SlotIndex = 0
	setElementsGain(heart, map[string]int{model.ElementEarth: 1})
	trigger()
	trigger()
	if engine.State.PendingAction != nil {
		resolvePendingSelection(t, engine, 0, heart.InstanceID+"|"+model.ElementEarth)
	}
	if p0.Equipment[0] != nil || len(p0.Graveyard) == 0 || p0.Graveyard[len(p0.Graveyard)-1] != heart {
		t.Fatalf("2411102 should destroy itself after losing all load, equipment=%v grave=%v", p0.Equipment[0], cardsToInfo(p0.Graveyard))
	}

	enemyEngine := setupEffectTest(t)
	enemyHeart := NewCardInstance(baseCard(t, "2411102"), 0, enemyEngine.State.TurnNumber)
	enemyEngine.State.Players[0].Equipment[0] = enemyHeart
	enemySpell := readySkill(baseCard(t, "3121105"), 1)
	if err := (Card2411102RottenAncientTreeHeart{}).OnSpellCast(&EffectContext{
		Engine:     enemyEngine,
		Source:     enemyHeart,
		Target:     enemySpell,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"cast_player": 1},
	}); err != nil {
		t.Fatalf("2411102 enemy first spell trigger failed: %v", err)
	}
	if enemyHeart.Statuses[rottenAncientTreeHeartSpellCountPrefix+"1"] != 1 || enemyHeart.Statuses[rottenAncientTreeHeartSpellCountPrefix+"0"] != 0 {
		t.Fatalf("2411102 should track spell counts by casting player, statuses=%v", enemyHeart.Statuses)
	}
}

func TestRoyalConflictRoseWhipGainsShadowLoadAfterFriendlyLoadLoss(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	whip := NewCardInstance(baseCard(t, "2421102"), 0, engine.State.TurnNumber)
	p0.Equipment[0] = whip
	whip.SlotIndex = 0
	ally := placeUnit(baseCard(t, "1421101"), 0, 0, 0, engine)
	setElementsGain(ally, map[string]int{model.ElementEarth: 3})
	enemy := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	setElementsGain(enemy, map[string]int{model.ElementEarth: 1})

	if removed := engine.reduceCardElementLoadWithTriggers(0, ally, model.ElementEarth, 1, nil); removed != 1 {
		t.Fatalf("expected one friendly load removed, got=%d", removed)
	}
	if got := whip.ElementsGainBonus[model.ElementShadow]; got != 1 {
		t.Fatalf("2421102 should gain one shadow load after friendly load loss, got=%d load=%v", got, effectiveElementsGain(whip))
	}

	engine.reduceCardElementLoadWithTriggers(0, ally, model.ElementEarth, 1, nil)
	engine.reduceCardElementLoadWithTriggers(0, ally, model.ElementEarth, 1, nil)
	if got := whip.ElementsGainBonus[model.ElementShadow]; got != 2 {
		t.Fatalf("2421102 gained shadow load should cap at two, got=%d load=%v", got, effectiveElementsGain(whip))
	}

	engine.reduceCardElementLoadWithTriggers(1, enemy, model.ElementEarth, 1, nil)
	if got := whip.ElementsGainBonus[model.ElementShadow]; got != 2 {
		t.Fatalf("2421102 should ignore enemy load loss, got=%d", got)
	}
}

func TestRoyalConflictAgedFrankenPrayerAndLoadLossWeakensEnemySpell(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	franken := placeUnit(baseCard(t, "1411101"), 0, 0, 0, engine)
	enemySpell := readySkill(baseCard(t, "3221001"), 1)
	p1.Skills[0] = enemySpell
	startEarth := reducibleElementLoad(franken, model.ElementEarth)

	if !(Card1411101AgedFrankenBaililan{}).IsPrayerAbility() {
		t.Fatal("1411101 should expose prayer timing")
	}
	if err := (Card1411101AgedFrankenBaililan{}).OnPerTurn(&EffectContext{
		Engine:     engine,
		Source:     franken,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"prayer": true},
	}); err != nil {
		t.Fatalf("1411101 prayer failed: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "aged_franken_prayer_load" {
		t.Fatalf("1411101 should ask which load to lose when multiple loads exist, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, model.ElementEarth)
	if reducibleElementLoad(franken, model.ElementEarth) != startEarth-1 || enemySpell.PowerBonus != -2 {
		t.Fatalf("1411101 prayer should lose earth load and weaken enemy spell, load=%v power=%d", effectiveElementsGain(franken), enemySpell.PowerBonus)
	}

	externalEngine := setupEffectTest(t)
	externalFranken := placeUnit(baseCard(t, "1411101"), 0, 0, 0, externalEngine)
	externalSpell := readySkill(baseCard(t, "3221001"), 1)
	externalEngine.State.Players[1].Skills[0] = externalSpell
	externalEngine.reduceCardElementLoadWithTriggers(0, externalFranken, model.ElementEarth, 1, nil)
	if externalSpell.PowerBonus != -2 {
		t.Fatalf("1411101 should also trigger when another effect removes its load, power=%d", externalSpell.PowerBonus)
	}

	if cardHasActivePerTurn(franken) || !cardHasActivePrayer(franken) || !isPrayerAbilityNumber("1411101") {
		t.Fatalf("1411101 should be exposed as prayer, not normal per-turn, perTurn=%v prayer=%v trait=%v p0=%v", cardHasActivePerTurn(franken), cardHasActivePrayer(franken), isPrayerAbilityNumber("1411101"), p0.PlayerID)
	}
}

func TestRoyalConflictBlackPineWandReducesFriendlyTargetSpellCost(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	wand := NewCardInstance(baseCard(t, "2621101"), 0, engine.State.TurnNumber)
	p0.Equipment[0] = wand
	grace := readySkill(baseCard(t, "3521108"), 0)
	p0.Skills[0] = grace
	friendly := placeUnit(baseCard(t, "1521001"), 0, 0, 0, engine)
	friendly.CurrentLife = maxLife(friendly) - 1
	target := SpellTarget{Type: "unit", Position: *friendly.Position}
	ownerID := 0
	target.OwnerID = &ownerID

	cost := engine.effectiveSkillUseCostForPurposeWithData(p0, grace, skillPurposeAttack, map[string]any{
		"spell_target":      target,
		"spell_target_unit": friendly,
	})
	if got := cost[model.ElementLight]; got != grace.Card.ElementsExpense[model.ElementLight]-1 {
		t.Fatalf("2621101 should reduce friendly-target spell by 1 light, cost=%v", cost)
	}
	enemy := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	enemyTarget := SpellTarget{Type: "unit", Position: *enemy.Position}
	enemyCost := engine.effectiveSkillUseCostForPurposeWithData(p0, grace, skillPurposeAttack, map[string]any{
		"spell_target":      enemyTarget,
		"spell_target_unit": enemy,
	})
	if got := enemyCost[model.ElementLight]; got != grace.Card.ElementsExpense[model.ElementLight] {
		t.Fatalf("2621101 should not reduce enemy-target spell, cost=%v", enemyCost)
	}

	p0.Elements[model.ElementLight] = 1
	if err := engine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id":  grace.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(friendly.Position.Col),
			"target_row":   float64(friendly.Position.Row),
			"target_owner": float64(0),
		},
	}); err != nil {
		t.Fatalf("cast 3521108 with 2621101 reduced cost: %v", err)
	}
	if p0.Elements[model.ElementLight] != 0 {
		t.Fatalf("2621101 reduced cast should spend exactly 1 light, elements=%v", p0.Elements)
	}
}

func TestRoyalConflictBloodRoseContractBindsOwnSpellToEarthOrShadowCompanion(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	contract := NewCardInstance(baseCard(t, "2421104"), 0, engine.State.TurnNumber)
	skill := readySkill(baseCard(t, "3421001"), 0)
	p0.Skills[0] = skill
	host := placeUnit(baseCard(t, "1421101"), 0, 0, 0, engine)
	hostLoad := engine.totalLoad(host)

	if err := (Card2421104BloodRoseContract{}).OnUseItem(&EffectContext{
		Engine:     engine,
		Source:     contract,
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("2421104 use failed: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_rose_contract_spell" {
		t.Fatalf("2421104 should ask for spell first, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, skill.InstanceID)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_rose_contract_host" {
		t.Fatalf("2421104 should ask for host second, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, host.InstanceID)
	if p0.Skills[0] != nil || len(host.BoundSkills) != 1 || host.BoundSkills[0] != skill {
		t.Fatalf("2421104 should move selected spell from slot to host, skills=%v bound=%v", cardsToInfo(p0.Skills[:]), cardsToInfo(host.BoundSkills))
	}
	if skill.SlotIndex != -1 || skill.IsHorizontal || skill.PowerBonus != hostLoad || !isTransferredBoundSkill(skill) {
		t.Fatalf("2421104 bound spell should preserve ready state, become slotless, gain host load power, and be marked transferred; slot=%d horizontal=%v power=%d hostLoad=%d statuses=%v", skill.SlotIndex, skill.IsHorizontal, skill.PowerBonus, hostLoad, skill.Statuses)
	}
	engine.destroyUnit(host, 0)
	if len(p0.Graveyard) == 0 || p0.Graveyard[len(p0.Graveyard)-1] != host || len(host.BoundSkills) != 0 || !containsCardInstance(p0.Exile, skill) {
		t.Fatalf("2421104 transferred bound spell should be exiled when host dies, grave=%v bound=%v exile=%v", cardsToInfo(p0.Graveyard), cardsToInfo(host.BoundSkills), cardsToInfo(p0.Exile))
	}
}

func TestRoyalConflictBloodRoseCurseBindsEnemySpellToEnemyChosenCompanion(t *testing.T) {
	engine := setupEffectTest(t)
	p1 := engine.State.Players[1]
	curse := NewCardInstance(baseCard(t, "2621102"), 0, engine.State.TurnNumber)
	enemySpell := readySkill(baseCard(t, "3221001"), 1)
	p1.Skills[0] = enemySpell
	host := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)

	if err := (Card2621102BloodRoseCurse{}).OnUseItem(&EffectContext{
		Engine:     engine,
		Source:     curse,
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("2621102 use failed: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_rose_curse_spell" || engine.State.PendingAction.PlayerID != 0 {
		t.Fatalf("2621102 should first ask caster to choose enemy spell, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, enemySpell.InstanceID)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_rose_curse_host" || engine.State.PendingAction.PlayerID != 1 {
		t.Fatalf("2621102 should then ask opponent to choose host, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 1, host.InstanceID)
	if p1.Skills[0] != nil || len(host.BoundSkills) != 1 || host.BoundSkills[0] != enemySpell {
		t.Fatalf("2621102 should move enemy spell from slot to chosen enemy host, skills=%v bound=%v", cardsToInfo(p1.Skills[:]), cardsToInfo(host.BoundSkills))
	}
	if enemySpell.SlotIndex != -1 || enemySpell.IsHorizontal || !isTransferredBoundSkill(enemySpell) {
		t.Fatalf("2621102 bound spell should preserve ready state, become slotless, and be marked transferred, slot=%d horizontal=%v statuses=%v", enemySpell.SlotIndex, enemySpell.IsHorizontal, enemySpell.Statuses)
	}
	engine.destroyUnit(host, 1)
	if len(host.BoundSkills) != 0 || !containsCardInstance(p1.Exile, enemySpell) {
		t.Fatalf("2621102 transferred bound spell should be exiled for its owner when host dies, bound=%v exile=%v", cardsToInfo(host.BoundSkills), cardsToInfo(p1.Exile))
	}
}

func TestRoyalConflictExileZoneIsPrivateToOwner(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	exiled := NewCardInstance(baseCard(t, "1021104"), 0, engine.State.TurnNumber)
	p0.Exile = append(p0.Exile, exiled)

	ownerView := engine.playerStateToInfo(p0, true)
	if exile, ok := ownerView["exile"].([]map[string]any); !ok || len(exile) != 1 || exile[0]["number"] != "1021104" {
		t.Fatalf("owner should see full exile cards, exile=%v", ownerView["exile"])
	}
	if ownerView["exile_count"] != 1 {
		t.Fatalf("owner should see exile count, state=%v", ownerView)
	}

	opponentView := engine.playerStateToInfo(p0, false)
	if _, ok := opponentView["exile"]; ok {
		t.Fatalf("opponent should not see full exile cards, state=%v", opponentView)
	}
	if opponentView["exile_count"] != 1 {
		t.Fatalf("opponent should see only exile count, state=%v", opponentView)
	}
}

func TestRoyalConflictGamblerUltimateDamagesAllPawns(t *testing.T) {
	engine := setupEffectTest(t)
	gambler := placeUnit(baseCard(t, "1011103"), 0, 0, 0, engine)
	ownPawn := placeUnit(baseCard(t, "1001101"), 0, 1, 0, engine)
	enemyPawn := placeUnit(baseCard(t, "1001101"), 1, 1, 0, engine)
	nonPawn := placeUnit(baseCard(t, "1021001"), 1, 2, 0, engine)
	ownPawn.CurrentLife = 2
	enemyPawn.CurrentLife = 2
	nonPawn.CurrentLife = 2

	if err := (Card1011103Gambler{}).OnUltimate(&EffectContext{
		Engine:     engine,
		Source:     gambler,
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("1011103 ultimate: %v", err)
	}
	if ownPawn.CurrentLife != 1 || enemyPawn.CurrentLife != 1 || nonPawn.CurrentLife != 2 {
		t.Fatalf("1011103 should damage only pawns on both fields, own=%d enemy=%d other=%d", ownPawn.CurrentLife, enemyPawn.CurrentLife, nonPawn.CurrentLife)
	}
}

func TestRoyalConflictNaturalCommunionRedistributesTwoEarthCompanionLoads(t *testing.T) {
	engine := setupEffectTest(t)
	a := placeUnit(baseCard(t, "1421101"), 0, 0, 0, engine)
	b := placeUnit(baseCard(t, "1421102"), 0, 1, 0, engine)
	setElementsGain(a, map[string]int{model.ElementEarth: 2})
	a.ElementsGainBonus = map[string]int{model.ElementWater: 1}
	setElementsGain(b, map[string]int{model.ElementAir: 1})

	if err := (Card2421105NaturalCommunion{}).OnUseItem(&EffectContext{
		Engine:     engine,
		Source:     NewCardInstance(baseCard(t, "2421105"), 0, engine.State.TurnNumber),
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("2421105 use failed: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "natural_communion_distribute" {
		t.Fatalf("2421105 should ask for two earth companions, pending=%+v", engine.State.PendingAction)
	}
	err := resolvePendingSelectionWithData(engine, 0, []string{a.InstanceID, b.InstanceID}, map[string]any{
		"load_distribution": map[string]any{
			a.InstanceID: map[string]any{model.ElementAir: float64(1), model.ElementWater: float64(1)},
			b.InstanceID: map[string]any{model.ElementEarth: float64(2)},
		},
	})
	if err != nil {
		t.Fatalf("resolve 2421105 distribution: %v", err)
	}
	if got := effectiveElementsGain(a); got[model.ElementAir] != 1 || got[model.ElementWater] != 1 || got[model.ElementEarth] != 0 {
		t.Fatalf("2421105 should apply redistributed load to first companion, load=%v", got)
	}
	if got := effectiveElementsGain(b); got[model.ElementEarth] != 2 || got[model.ElementAir] != 0 || got[model.ElementWater] != 0 {
		t.Fatalf("2421105 should apply redistributed load to second companion, load=%v", got)
	}
	if len(a.ElementsGainBonus) != 0 || len(b.ElementsGainBonus) != 0 {
		t.Fatalf("2421105 should make redistributed load the actual load, bonus a=%v b=%v", a.ElementsGainBonus, b.ElementsGainBonus)
	}

	invalidEngine := setupEffectTest(t)
	invalidA := placeUnit(baseCard(t, "1421101"), 0, 0, 0, invalidEngine)
	invalidB := placeUnit(baseCard(t, "1421102"), 0, 1, 0, invalidEngine)
	setElementsGain(invalidA, map[string]int{model.ElementEarth: 1})
	setElementsGain(invalidB, map[string]int{model.ElementEarth: 1})
	if err := (Card2421105NaturalCommunion{}).OnUseItem(&EffectContext{Engine: invalidEngine, Source: NewCardInstance(baseCard(t, "2421105"), 0, 1), PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("2421105 invalid setup use failed: %v", err)
	}
	err = resolvePendingSelectionWithData(invalidEngine, 0, []string{invalidA.InstanceID, invalidB.InstanceID}, map[string]any{
		"load_distribution": map[string]any{
			invalidA.InstanceID: map[string]any{model.ElementEarth: float64(2)},
			invalidB.InstanceID: map[string]any{model.ElementEarth: float64(2)},
		},
	})
	if err == nil {
		t.Fatalf("2421105 should reject non-conserving load distribution")
	}
	if invalidEngine.State.PendingAction == nil || invalidEngine.State.PendingAction.Type != "natural_communion_distribute" {
		t.Fatalf("2421105 invalid distribution should keep pending action, pending=%+v", invalidEngine.State.PendingAction)
	}
}

func TestRoyalConflictLopsiusRageSuppressesOpponentResponses(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	scroll := NewCardInstance(baseCard(t, "2121111"), 0, engine.State.TurnNumber)
	p0.Hand = []*CardInstance{scroll}
	setAllElements(p0, 9)
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)

	castCounter := NewCardInstance(baseCard(t, "2021018"), 1, engine.State.TurnNumber)
	castCounter.IsSetCounter = true
	p1.Equipment[0] = castCounter
	hitCounter := NewCardInstance(baseCard(t, "2021113"), 1, engine.State.TurnNumber)
	hitCounter.IsSetCounter = true
	p1.Equipment[1] = hitCounter

	if err := engine.HandleAction(0, ActionMessage{
		Action: "use_item",
		Data: map[string]any{
			"instance_id": scroll.InstanceID,
			"target_type": "unit",
			"target_col":  float64(target.Position.Col),
			"target_row":  float64(target.Position.Row),
		},
	}); err != nil {
		t.Fatalf("use 2121111: %v", err)
	}
	if engine.State.PendingAction != nil || engine.State.Phase != PhaseDefenseWindow {
		t.Fatalf("2121111 should skip spell-cast counters and open defense, phase=%s pending=%+v", engine.State.Phase, engine.State.PendingAction)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "react_spell", Data: map[string]any{"instance_id": "any"}}); err == nil {
		t.Fatalf("2121111 should reject opponent spell reactions")
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve 2121111 without defense: %v", err)
	}
	if engine.State.PendingAction != nil || p1.Equipment[0] != castCounter || p1.Equipment[1] != hitCounter {
		t.Fatalf("2121111 should skip hit counters too, pending=%+v equipment=%v", engine.State.PendingAction, cardsToInfo(p1.Equipment[:]))
	}
	if target.CurrentLife >= target.Card.Life {
		t.Fatalf("2121111 should still resolve damage, target life=%d start=%d", target.CurrentLife, target.Card.Life)
	}
}

func TestRoyalConflictIceLockRuneLocksLearnedEnemySkill(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	skill := NewCardInstance(baseCard(t, "3121105"), 0, engine.State.TurnNumber)
	p0.SkillPool = []*CardInstance{skill}
	setAllElements(p0, 9)
	counter := NewCardInstance(baseCard(t, "2221103"), 1, engine.State.TurnNumber)
	counter.IsSetCounter = true
	p1.Equipment[0] = counter
	p1.Elements[model.ElementWater] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "learn_skill", Data: map[string]any{"instance_id": skill.InstanceID}}); err != nil {
		t.Fatalf("learn skill into ice lock: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_trigger" {
		t.Fatalf("2221103 should prompt when opponent learns a skill, pending=%+v", engine.State.PendingAction)
	}
	if err := resolvePendingSelectionWithData(engine, 1, []string{counter.InstanceID}, nil); err != nil {
		t.Fatalf("resolve 2221103 counter: %v", err)
	}
	if skill.Statuses[StatusCannotUseSkillUntilTurn] != engine.State.TurnNumber+1 {
		t.Fatalf("2221103 should lock learned skill until next turn ends, statuses=%v turn=%d", skill.Statuses, engine.State.TurnNumber)
	}
	if p1.Equipment[0] != nil || countCardNumber(p1.Graveyard, "2221103") != 1 {
		t.Fatalf("2221103 should move to graveyard after reveal, equipment=%v grave=%v", cardToInfo(p1.Equipment[0]), cardsToInfo(p1.Graveyard))
	}
}

func TestRoyalConflictPunishmentRuneTriggersAfterThirdEnemySpellAttack(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	counter := NewCardInstance(baseCard(t, "2521109"), 1, engine.State.TurnNumber)
	counter.IsSetCounter = true
	p1.Equipment[0] = counter
	p1.Elements[model.ElementLight] = 2
	spell := readySkill(baseCard(t, "3121105"), 0)
	target := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	startLife := target.CurrentLife
	engine.recordSpellCast(0, spell)
	engine.recordSpellCast(0, spell)
	engine.recordSpellCast(0, spell)
	if totalSpellsCastThisTurn(p0) != 3 {
		t.Fatalf("test setup should record three spell attacks, got %d", totalSpellsCastThisTurn(p0))
	}

	prompted := engine.promptOpponentCounterTrap(0, TriggerOnSpellCast, spell, map[string]any{"cast_player": 0, "attacker": 0}, nil)
	if !prompted || engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_trigger" {
		t.Fatalf("2521109 should prompt after the third enemy spell attack, prompted=%v pending=%+v", prompted, engine.State.PendingAction)
	}
	if err := resolvePendingSelectionWithData(engine, 1, []string{counter.InstanceID}, nil); err != nil {
		t.Fatalf("resolve 2521109 counter: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "punishment_rune_damage" {
		t.Fatalf("2521109 should ask for an enemy companion target, pending=%+v", engine.State.PendingAction)
	}
	if err := resolvePendingSelectionWithData(engine, 1, []string{target.InstanceID}, nil); err != nil {
		t.Fatalf("resolve 2521109 damage target: %v", err)
	}
	if target.CurrentLife != startLife-2 {
		t.Fatalf("2521109 should deal 2 damage, life=%d start=%d", target.CurrentLife, startLife)
	}
	if p1.Equipment[0] != nil || countCardNumber(p1.Graveyard, "2521109") != 1 {
		t.Fatalf("2521109 should move to graveyard after reveal, equipment=%v grave=%v", cardToInfo(p1.Equipment[0]), cardsToInfo(p1.Graveyard))
	}
}

func TestRoyalConflictAutumnMapleLordRewardsEarthOverexertRemainder(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p0.Hero = NewCardInstance(baseCard(t, "4411102"), 0, engine.State.TurnNumber)
	earth := placeUnit(baseCard(t, "1421101"), 0, 0, 0, engine)
	setElementsGain(earth, map[string]int{model.ElementEarth: 3})
	p0.Elements = map[string]int{model.ElementEarth: 1}

	if !engine.payDefenseCostWithOptions(p0, map[string]int{model.ElementEarth: 2}, ActionMessage{}, []*CardInstance{earth}, false) {
		t.Fatal("defense payment with earth overexert should succeed")
	}
	if !earth.IsHorizontal {
		t.Fatal("overexerted earth unit should become horizontal")
	}
	if got := p0.Elements[model.ElementEarth]; got != 4 {
		t.Fatalf("4411102 should return twice the unused earth load: got %d elements=%v", got, p0.Elements)
	}

	blockedEngine := setupEffectTest(t)
	blockedP0 := blockedEngine.State.Players[0]
	blockedP0.Hero = NewCardInstance(baseCard(t, "4311003"), 0, blockedEngine.State.TurnNumber)
	blockedEarth := placeUnit(baseCard(t, "1421101"), 0, 0, 0, blockedEngine)
	setElementsGain(blockedEarth, map[string]int{model.ElementEarth: 3})
	blockedP0.Elements = map[string]int{model.ElementEarth: 1}
	if !blockedEngine.payDefenseCostWithOptions(blockedP0, map[string]int{model.ElementEarth: 2}, ActionMessage{}, []*CardInstance{blockedEarth}, false) {
		t.Fatal("baseline defense payment should succeed")
	}
	if got := blockedP0.Elements[model.ElementEarth]; got != 0 {
		t.Fatalf("without 4411102, overexert remainder should still be lost, got %d elements=%v", got, blockedP0.Elements)
	}
}

func TestRoyalConflictGuardianRunePreventsLethalDamage(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	target.CurrentLife = 1
	counter := NewCardInstance(baseCard(t, "2021114"), 0, engine.State.TurnNumber)
	counter.IsSetCounter = true
	p0.Equipment[0] = counter
	p0.Elements[model.ElementArcane] = 1

	engine.dealDamageWithExtra(target, 3, 0, map[string]any{"damage_source": "test", "attacker": 1})
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_trigger" {
		t.Fatalf("2021114 should prompt before lethal friendly damage, pending=%+v", engine.State.PendingAction)
	}
	if err := resolvePendingSelectionWithData(engine, 0, []string{counter.InstanceID}, nil); err != nil {
		t.Fatalf("resolve 2021114 counter: %v", err)
	}
	if target.CurrentLife != 1 || engine.State.Players[0].Units[0][0] != target {
		t.Fatalf("2021114 should prevent the whole lethal damage event, life=%d unit=%v", target.CurrentLife, cardToInfo(engine.State.Players[0].Units[0][0]))
	}
	if p0.Equipment[0] != nil || countCardNumber(p0.Graveyard, "2021114") != 1 {
		t.Fatalf("2021114 should move to graveyard after reveal, equipment=%v grave=%v", cardToInfo(p0.Equipment[0]), cardsToInfo(p0.Graveyard))
	}

	declineEngine := setupEffectTest(t)
	declineP0 := declineEngine.State.Players[0]
	declineTarget := placeUnit(baseCard(t, "1021001"), 0, 0, 0, declineEngine)
	declineTarget.CurrentLife = 1
	declineCounter := NewCardInstance(baseCard(t, "2021114"), 0, declineEngine.State.TurnNumber)
	declineCounter.IsSetCounter = true
	declineP0.Equipment[0] = declineCounter
	declineP0.Elements[model.ElementArcane] = 1
	declineEngine.dealDamageWithExtra(declineTarget, 3, 0, map[string]any{"damage_source": "test", "attacker": 1})
	if declineEngine.State.PendingAction == nil {
		t.Fatal("2021114 decline setup should prompt")
	}
	if err := resolvePendingSelectionWithData(declineEngine, 0, nil, nil); err != nil {
		t.Fatalf("decline 2021114: %v", err)
	}
	if declineEngine.State.Players[0].Units[0][0] != nil {
		t.Fatalf("declining 2021114 should let lethal damage resolve, unit=%v", cardToInfo(declineEngine.State.Players[0].Units[0][0]))
	}
}

func TestRoyalConflictIceSoulSealForgeHalvesHighPowerSpell(t *testing.T) {
	engine := setupEffectTest(t)
	p1 := engine.State.Players[1]
	spell := readySkill(baseCard(t, "3111101"), 0)
	counter := NewCardInstance(baseCard(t, "2221111"), 1, engine.State.TurnNumber)
	counter.IsSetCounter = true
	p1.Equipment[0] = counter
	p1.Elements[model.ElementWater] = 2
	engine.State.PendingSpell = &SpellCast{AttackerID: 0, Skill: spell, TotalPower: 13}
	extraData := map[string]any{"cast_player": 0, "attacker": 0, "power": 13}

	if !engine.promptOpponentCounterTrap(0, TriggerOnSpellCast, spell, extraData, nil) {
		t.Fatal("2221111 should prompt for enemy spell total power above 10")
	}
	if err := resolvePendingSelectionWithData(engine, 1, []string{counter.InstanceID}, nil); err != nil {
		t.Fatalf("resolve 2221111 counter: %v", err)
	}
	if engine.State.PendingSpell.TotalPower != 7 || extraData["power"].(int) != 7 {
		t.Fatalf("2221111 should halve power upward, pending=%d data=%v", engine.State.PendingSpell.TotalPower, extraData["power"])
	}
	if p1.Equipment[0] != nil || countCardNumber(p1.Graveyard, "2221111") != 1 {
		t.Fatalf("2221111 should move to graveyard after reveal, equipment=%v grave=%v", cardToInfo(p1.Equipment[0]), cardsToInfo(p1.Graveyard))
	}

	blockedEngine := setupEffectTest(t)
	blockedP1 := blockedEngine.State.Players[1]
	blockedCounter := NewCardInstance(baseCard(t, "2221111"), 1, blockedEngine.State.TurnNumber)
	blockedCounter.IsSetCounter = true
	blockedP1.Equipment[0] = blockedCounter
	if candidates := blockedEngine.eligibleCounterTraps(1, TriggerOnSpellCast, spell, map[string]any{"cast_player": 0, "attacker": 0, "power": 10}); len(candidates) != 0 {
		t.Fatalf("2221111 should require power greater than 10, candidates=%v", cardsToInfo(candidates))
	}
}

func TestRoyalConflictNaturalEchoResetsAndGrantsSameExtraTarget(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	echo := readySkill(baseCard(t, "3421104"), 0)
	echo.IsHorizontal = true
	p0.Skills[0] = echo
	support := placeUnit(baseCard(t, "1421101"), 0, 0, 0, engine)
	setElementsGain(support, map[string]int{model.ElementEarth: 2})

	if err := (Card3421104NaturalEcho{}).OnPerTurn(&EffectContext{Engine: engine, Source: echo, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("3421104 per-turn: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "natural_echo_remove_load" {
		t.Fatalf("3421104 should ask which earth load to remove, pending=%+v", engine.State.PendingAction)
	}
	if err := resolvePendingSelectionWithData(engine, 0, []string{support.InstanceID}, nil); err != nil {
		t.Fatalf("resolve 3421104 load removal: %v", err)
	}
	if effectiveElementsGain(support)[model.ElementEarth] != 1 || echo.IsHorizontal || echo.UsedThisTurn != 1 {
		t.Fatalf("3421104 should remove one earth load and reset itself, support_load=%v horizontal=%v used=%d", effectiveElementsGain(support), echo.IsHorizontal, echo.UsedThisTurn)
	}
	if got := engine.effectiveSpellPower(0, echo, nil); got != echo.Card.Power+2 {
		t.Fatalf("3421104 next cast should gain +2 power, got=%d base=%d", got, echo.Card.Power)
	}

	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	target.CurrentLife = 20
	startLife := target.CurrentLife
	p0.Elements[model.ElementEarth] = 10
	if err := engine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id":      echo.InstanceID,
			"target_type":      "unit",
			"target_col":       float64(target.Position.Col),
			"target_row":       float64(target.Position.Row),
			"extra_target_col": float64(target.Position.Col),
			"extra_target_row": float64(target.Position.Row),
		},
	}); err != nil {
		t.Fatalf("cast 3421104 with same extra target: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve 3421104 hit: %v", err)
	}
	wantDamage := max(echo.Card.Attack+echo.AttackBonus, 0) * 2
	if target.CurrentLife != startLife-wantDamage {
		t.Fatalf("3421104 same extra target should apply damage twice, life=%d start=%d wantDamage=%d", target.CurrentLife, startLife, wantDamage)
	}
	if engine.hasNextDriveSpellExtraTarget(p0, echo) {
		t.Fatalf("3421104 extra target modifier should be consumed, modifiers=%+v", p0.TempModifiers)
	}
}

func TestRoyalConflictBloodRoseSealBindsWhenMarkedEnemyDies(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	p0.Hero = NewCardInstance(baseCard(t, "1011001"), 0, 1)
	seal := readySkill(baseCard(t, "3621104"), 0)
	p0.Skills[0] = seal
	seal.SlotIndex = 0
	target := NewCardInstance(baseCard(t, "1021001"), 1, 1)
	target.Position = &Position{Col: 1, Row: 1}
	p1.Units[1][1] = target

	if err := (Card3621104BloodRoseSeal{}).OnEnter(&EffectContext{Engine: engine, Source: seal, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("blood rose seal enter: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_rose_seal_mark" {
		t.Fatalf("blood rose seal should ask for an enemy unit, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, target.InstanceID)
	if target.Statuses[bloodRoseSealMarkerStatus(seal)] != 1 {
		t.Fatalf("blood rose seal should mark the selected enemy, statuses=%v", target.Statuses)
	}

	engine.destroyUnitWithData(target, 1, map[string]any{"death_cause": "test"})
	if p0.Skills[0] != nil || len(p0.Hero.BoundSkills) != 1 || p0.Hero.BoundSkills[0] != seal {
		t.Fatalf("marked enemy death should bind seal to hero, skills=%v bound=%v", cardsToInfo(p0.Skills[:]), cardsToInfo(p0.Hero.BoundSkills))
	}
	if seal.SlotIndex != -1 || seal.IsHorizontal || !isTransferredBoundSkill(seal) {
		t.Fatalf("bound seal should preserve ready state, become slotless, and be marked transferred, slot=%d horizontal=%v statuses=%v", seal.SlotIndex, seal.IsHorizontal, seal.Statuses)
	}
	if cost := engine.effectiveSkillUseCost(p0, seal); cost[model.ElementShadow] != 1 {
		t.Fatalf("bound seal should cost one less shadow, cost=%v", cost)
	}
}

func TestRoyalConflictBloodRoseSealExpiresAtNextOwnTurnEnd(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	p0.Hero = NewCardInstance(baseCard(t, "1011001"), 0, 1)
	seal := readySkill(baseCard(t, "3621104"), 0)
	p0.Skills[0] = seal
	seal.SlotIndex = 0
	target := NewCardInstance(baseCard(t, "1021001"), 1, 1)
	target.Position = &Position{Col: 1, Row: 1}
	p1.Units[1][1] = target

	if err := (Card3621104BloodRoseSeal{}).OnEnter(&EffectContext{Engine: engine, Source: seal, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("blood rose seal enter: %v", err)
	}
	resolvePendingSelection(t, engine, 0, target.InstanceID)

	if err := (Card3621104BloodRoseSeal{}).OnTurnEnd(&EffectContext{
		Engine:     engine,
		Source:     seal,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"ended_player": 0},
	}); err != nil {
		t.Fatalf("blood rose seal current turn end: %v", err)
	}
	if target.Statuses[bloodRoseSealMarkerStatus(seal)] != 1 {
		t.Fatalf("blood rose seal should last past the current turn end, statuses=%v", target.Statuses)
	}

	engine.State.TurnNumber = 3
	if err := (Card3621104BloodRoseSeal{}).OnTurnEnd(&EffectContext{
		Engine:     engine,
		Source:     seal,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"ended_player": 0},
	}); err != nil {
		t.Fatalf("blood rose seal next turn end: %v", err)
	}
	if target.Statuses[bloodRoseSealMarkerStatus(seal)] != 0 {
		t.Fatalf("blood rose seal marker should expire at next own turn end, statuses=%v", target.Statuses)
	}

	engine.destroyUnitWithData(target, 1, map[string]any{"death_cause": "test"})
	if p0.Skills[0] != seal || len(p0.Hero.BoundSkills) != 0 {
		t.Fatalf("expired blood rose seal should not bind, skills=%v bound=%v", cardsToInfo(p0.Skills[:]), cardsToInfo(p0.Hero.BoundSkills))
	}
}

func TestRoyalConflictExileSotorAddsAdjacentSpellTargets(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	placeUnit(baseCard(t, "1111102"), 0, 0, 2, engine)
	spell := readySkill(baseCard(t, "3121002"), 0)
	p0.Skills[0] = spell
	spell.SlotIndex = 0
	p0.Elements[model.ElementFire] = 10

	main := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	left := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	right := placeUnit(baseCard(t, "1021001"), 1, 2, 0, engine)
	for _, unit := range []*CardInstance{main, left, right} {
		unit.CurrentLife = 10
	}

	if err := engine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id": spell.InstanceID,
			"target_type": "unit",
			"target_col":  float64(main.Position.Col),
			"target_row":  float64(main.Position.Row),
		},
	}); err != nil {
		t.Fatalf("cast spell with Sotor aura: %v", err)
	}
	if engine.State.PendingSpell == nil || len(engine.State.PendingSpell.ExtraTargets) != 3 {
		t.Fatalf("Sotor should add three adjacent extra targets, pending=%+v", engine.State.PendingSpell)
	}
	heroStartLife := p1.Hero.CurrentLife
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve Sotor spell hit: %v", err)
	}
	for _, unit := range []*CardInstance{main, left, right} {
		if unit.CurrentLife != 8 {
			t.Fatalf("Sotor spell should damage main and adjacent targets, unit=%s life=%d", unit.InstanceID, unit.CurrentLife)
		}
	}
	if p1.Hero.CurrentLife != heroStartLife-2 {
		t.Fatalf("Sotor spell should also damage adjacent hero, life=%d start=%d", p1.Hero.CurrentLife, heroStartLife)
	}
}

func TestRoyalConflictBloodThornGardenResummonsAfterFriendlyKill(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	thorn := placeUnit(baseCard(t, "1611102"), 0, 0, 0, engine)
	thorn.CurrentLife = 0
	p0.Elements[model.ElementShadow] = 1
	ctx := &EffectContext{
		Engine:     engine,
		Source:     thorn,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"attacker": 0},
	}

	engine.destroyUnitWithData(thorn, 0, ctx.ExtraData)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_thorn_resummon" {
		t.Fatalf("blood thorn should ask to resummon after friendly kill, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, thorn.InstanceID)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_thorn_resummon_position" {
		t.Fatalf("blood thorn should ask for resummon position, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, positionSelectionID(Position{Col: 0, Row: 0}))
	if p0.Units[0][0] != thorn || thorn.CurrentLife != thorn.Card.Life || !thorn.IsHorizontal {
		t.Fatalf("blood thorn should return fresh at its old position, unit=%v card=%v", p0.Units[0][0], cardToInfo(thorn))
	}
	if p0.Elements[model.ElementShadow] != 0 || containsCardInstance(p0.Graveyard, thorn) {
		t.Fatalf("blood thorn should spend one shadow and leave graveyard, elements=%v grave=%v", p0.Elements, cardsToInfo(p0.Graveyard))
	}
}

func TestRoyalConflictBloodThornGardenIgnoresNonFriendlyKill(t *testing.T) {
	engine := setupEffectTest(t)
	thorn := placeUnit(baseCard(t, "1611102"), 0, 0, 0, engine)
	thorn.CurrentLife = 0
	engine.State.Players[0].Elements[model.ElementShadow] = 1

	engine.destroyUnitWithData(thorn, 0, map[string]any{"attacker": 1})
	if engine.State.PendingAction != nil {
		t.Fatalf("blood thorn should not trigger after enemy kill, pending=%+v", engine.State.PendingAction)
	}
	if engine.State.Players[0].Units[0][0] != nil || !containsCardInstance(engine.State.Players[0].Graveyard, thorn) {
		t.Fatalf("blood thorn should stay dead after enemy kill, unit=%v grave=%v", engine.State.Players[0].Units[0][0], cardsToInfo(engine.State.Players[0].Graveyard))
	}
}

func TestRoyalConflictRebornUnitsCanChooseCurrentEmptyPositions(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	bone := placeUnit(baseCard(t, "1621011"), 0, 0, 0, engine)
	devourFuel := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
	devourFuel.CurrentLife = 3
	fearDemon := NewCardInstance(baseCard(t, "1621003"), 0, engine.State.TurnNumber)
	p0.Hand = append(p0.Hand, fearDemon)

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": fearDemon.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
		"devour_ids":  []any{bone.InstanceID, devourFuel.InstanceID},
	}}); err != nil {
		t.Fatalf("summon fear demon with devour: %v", err)
	}
	if p0.Units[0][0] != fearDemon {
		t.Fatalf("fear demon should occupy bone knight's old position, unit=%v", cardToInfo(p0.Units[0][0]))
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "bone_knight_reborn" {
		t.Fatalf("bone knight should ask to reborn after devour, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, bone.InstanceID)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "bone_knight_reborn_position" ||
		candidateContains(engine.State.PendingAction.Candidates, positionSelectionID(Position{Col: 0, Row: 0})) {
		t.Fatalf("bone knight should ask for a current empty position, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, positionSelectionID(Position{Col: 2, Row: 2}))
	if p0.Units[2][2] != bone || bone.Statuses[boneKnightRebornStatus] != 1 || containsCardInstance(p0.Graveyard, bone) {
		t.Fatalf("bone knight should reborn at chosen current empty position, pos=%v grave=%v statuses=%v", bone.Position, cardsToInfo(p0.Graveyard), bone.Statuses)
	}

	thorn := placeUnit(baseCard(t, "1611102"), 0, 1, 0, engine)
	thorn.CurrentLife = 0
	p0.Elements[model.ElementShadow] = 1
	engine.destroyUnitWithData(thorn, 0, map[string]any{"attacker": 0})
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_thorn_resummon" {
		t.Fatalf("blood thorn should ask to resummon after friendly kill, pending=%+v", engine.State.PendingAction)
	}
	blocker := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
	resolvePendingSelection(t, engine, 0, thorn.InstanceID)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_thorn_resummon_position" ||
		candidateContains(engine.State.PendingAction.Candidates, positionSelectionID(Position{Col: 1, Row: 0})) {
		t.Fatalf("blood thorn should ask for a current empty position, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, positionSelectionID(Position{Col: 2, Row: 1}))
	if p0.Units[2][1] != thorn || containsCardInstance(p0.Graveyard, thorn) || p0.Elements[model.ElementShadow] != 0 {
		t.Fatalf("blood thorn should spend one shadow and resummon at chosen current empty position, unit=%v blocker=%v grave=%v elements=%v", cardToInfo(p0.Units[2][1]), cardToInfo(blocker), cardsToInfo(p0.Graveyard), p0.Elements)
	}
}

func TestRoyalConflictRobertBlackPineMarksFriendlyDamageAndDeath(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	robert := placeUnit(baseCard(t, "1611103"), 0, 0, 0, engine)
	ally := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
	ally.CurrentLife = 5

	engine.dealDamageWithExtra(ally, 1, 0, map[string]any{"damage_source": "effect", "attacker": 0})
	if robert.Statuses[robertBlackPineMarkerStatus] != 1 {
		t.Fatalf("Robert should gain one marker from friendly damage, statuses=%v", robert.Statuses)
	}
	ally.CurrentLife = 0
	engine.destroyUnitWithData(ally, 0, map[string]any{"death_cause": "test", "attacker": 0})
	if robert.Statuses[robertBlackPineMarkerStatus] != 3 {
		t.Fatalf("Robert should gain two markers from friendly-caused death, statuses=%v", robert.Statuses)
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  robert.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use Robert ability: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "robert_black_pine_reward" {
		t.Fatalf("Robert should ask for reward, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, "load")
	if robert.Statuses[robertBlackPineMarkerStatus] != 0 || effectiveElementsGain(robert)[model.ElementShadow] != robert.Card.ElementsGain[model.ElementShadow]+1 {
		t.Fatalf("Robert should spend markers for shadow load, statuses=%v load=%v", robert.Statuses, effectiveElementsGain(robert))
	}
	if p0.Units[0][0] != robert {
		t.Fatalf("Robert should remain on field")
	}
}

func TestRoyalConflictRobertBlackPineIgnoresEnemyDamage(t *testing.T) {
	engine := setupEffectTest(t)
	robert := placeUnit(baseCard(t, "1611103"), 0, 0, 0, engine)
	ally := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)

	engine.dealDamageWithExtra(ally, 1, 0, map[string]any{"damage_source": "spell", "attacker": 1})
	if robert.Statuses[robertBlackPineMarkerStatus] != 0 {
		t.Fatalf("Robert should ignore enemy damage, statuses=%v", robert.Statuses)
	}
}

func TestRoyalConflictSpellHitStatsRollToLastTurn(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]

	engine.recordSpellHitStats(0, 2, 3)
	engine.recordSpellHitStats(0, 1, 0)
	rollSpellHitTracking(p0)
	if p0.SpellHitsLastTurn != 2 || p0.SpellHitTargetsLastTurn != 3 || p0.SpellDamageLastTurn != 3 {
		t.Fatalf("spell hit stats should roll to last turn, player=%+v", p0)
	}
	if p0.SpellHitsThisTurn != 0 || p0.SpellHitTargetsThisTurn != 0 || p0.SpellDamageThisTurn != 0 {
		t.Fatalf("spell hit stats should clear current turn after roll, player=%+v", p0)
	}
}

func TestRoyalConflictSindarielDamagesForEnemyLastTurnSpellMilestones(t *testing.T) {
	engine := setupEffectTest(t)
	p1 := engine.State.Players[1]
	sindariel := placeUnit(baseCard(t, "1411102"), 0, 0, 0, engine)
	targetA := placeUnit(baseCard(t, "1421101"), 1, 0, 0, engine)
	targetB := placeUnit(baseCard(t, "1421101"), 1, 1, 0, engine)
	targetC := placeUnit(baseCard(t, "1421101"), 1, 2, 0, engine)
	for _, target := range []*CardInstance{targetA, targetB, targetC} {
		target.CurrentLife = 5
	}
	p1.SpellHitsLastTurn = 3
	p1.SpellHitTargetsLastTurn = 3
	p1.SpellDamageLastTurn = 3

	if err := (Card1411102WhisperElfKingSindariel{}).OnEnter(&EffectContext{
		Engine:     engine,
		Source:     sindariel,
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("Sindariel enter failed: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "sindariel_entry_damage" ||
		engine.State.PendingAction.MinSelect != 0 || engine.State.PendingAction.MaxSelect != 3 {
		t.Fatalf("Sindariel should allow up to three optional targets, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, targetA.InstanceID, targetB.InstanceID, targetC.InstanceID)
	for _, target := range []*CardInstance{targetA, targetB, targetC} {
		if target.CurrentLife != 3 {
			t.Fatalf("Sindariel should deal two damage to each selected target, target=%s life=%d", target.InstanceID, target.CurrentLife)
		}
	}
}

func TestRoyalConflictSupremeQueenSummonsFireCompanionsWithTemporaryImmunity(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	queen := placeUnit(baseCard(t, "1111101"), 0, 1, 0, engine)
	fireA := NewCardInstance(baseCard(t, "1121001"), 0, engine.State.TurnNumber)
	fireB := NewCardInstance(baseCard(t, "1121002"), 0, engine.State.TurnNumber)
	water := NewCardInstance(baseCard(t, "1221001"), 0, engine.State.TurnNumber)
	p0.Hand = []*CardInstance{fireA, fireB, water}

	if err := (Card1111101SupremeQueenDailinCeltic{}).OnEnter(&EffectContext{
		Engine:     engine,
		Source:     queen,
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("queen enter failed: %v", err)
	}
	if queen.Statuses[temporaryDamageAndNegativeImmunityUntilStatus] < engine.State.TurnNumber {
		t.Fatalf("queen should mark herself immune, statuses=%v", queen.Statuses)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "supreme_queen_summon_cards" ||
		candidateContains(engine.State.PendingAction.Candidates, water.InstanceID) ||
		!candidateContains(engine.State.PendingAction.Candidates, fireA.InstanceID) ||
		!candidateContains(engine.State.PendingAction.Candidates, fireB.InstanceID) {
		t.Fatalf("queen should ask for fire companions only, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, fireA.InstanceID, fireB.InstanceID)
	firstPos := Position{Col: 0, Row: 0}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "supreme_queen_summon_position" ||
		!candidateContains(engine.State.PendingAction.Candidates, positionSelectionID(firstPos)) {
		t.Fatalf("queen should ask for adjacent position, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, positionSelectionID(firstPos))
	secondPos := Position{Col: 2, Row: 0}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "supreme_queen_summon_position" ||
		!candidateContains(engine.State.PendingAction.Candidates, positionSelectionID(secondPos)) {
		t.Fatalf("queen should ask for next adjacent position, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, positionSelectionID(secondPos))
	if p0.Units[firstPos.Col][firstPos.Row] != fireA || p0.Units[secondPos.Col][secondPos.Row] != fireB ||
		containsCardInstance(p0.Hand, fireA) || containsCardInstance(p0.Hand, fireB) || !containsCardInstance(p0.Hand, water) {
		t.Fatalf("queen should summon selected fire companions from hand, units=%v/%v hand=%v", p0.Units[firstPos.Col][firstPos.Row], p0.Units[secondPos.Col][secondPos.Row], cardsToInfo(p0.Hand))
	}

	queenLife := queen.CurrentLife
	fireLife := fireA.CurrentLife
	engine.dealDamageWithExtra(queen, 3, 0, map[string]any{"damage_source": "test", "attacker": 1})
	engine.dealDamageWithExtra(fireA, 3, 0, map[string]any{"damage_source": "test", "attacker": 1})
	if queen.CurrentLife != queenLife || fireA.CurrentLife != fireLife {
		t.Fatalf("queen immunity should prevent damage, queen=%d/%d fire=%d/%d", queen.CurrentLife, queenLife, fireA.CurrentLife, fireLife)
	}
	if engine.addStatus(fireA, StatusFreeze, 1) || fireA.Statuses[StatusFreeze] != 0 {
		t.Fatalf("queen immunity should reject negative statuses, statuses=%v", fireA.Statuses)
	}

	engine.State.TurnNumber += 2
	if !engine.addStatus(fireA, StatusFreeze, 1) {
		t.Fatalf("queen immunity should expire for negative statuses, statuses=%v", fireA.Statuses)
	}
	engine.dealDamageWithExtra(fireA, 1, 0, map[string]any{"damage_source": "test", "attacker": 1})
	if fireA.CurrentLife != fireLife-1 {
		t.Fatalf("queen immunity should expire for damage, life=%d start=%d", fireA.CurrentLife, fireLife)
	}
}

func TestRoyalConflictManesArbitrationPermanentWaterPowerChoice(t *testing.T) {
	engine := setupEffectTest(t)
	waterSkill := readySkill(baseCard(t, "3221001"), 0)
	item := NewCardInstance(baseCard(t, "2211102"), 0, engine.State.TurnNumber)

	if err := (Card2211102ManesArbitration{}).OnEquip(&EffectContext{
		Engine:     engine,
		Source:     item,
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("Manes equip failed: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "manes_arbitration_choice" ||
		!candidateContains(engine.State.PendingAction.Candidates, "all_water_power") {
		t.Fatalf("Manes should ask for permanent mode, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, "all_water_power")
	if power := engine.effectiveSkillPowerForPurpose(0, waterSkill, skillPurposeAttack); power != waterSkill.Card.Power+2 {
		t.Fatalf("Manes all-water mode should add +2 power, got=%d base=%d modifiers=%+v", power, waterSkill.Card.Power, engine.State.Players[0].TempModifiers)
	}
}

func TestRoyalConflictManesArbitrationEmpowersOneWaterSkillAndLocksLearning(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	waterSkill := readySkill(baseCard(t, "3221001"), 0)
	p0.Skills[0] = waterSkill
	item := NewCardInstance(baseCard(t, "2211102"), 0, engine.State.TurnNumber)

	if err := (Card2211102ManesArbitration{}).OnEquip(&EffectContext{
		Engine:     engine,
		Source:     item,
		PlayerID:   0,
		OpponentID: 1,
	}); err != nil {
		t.Fatalf("Manes equip failed: %v", err)
	}
	resolvePendingSelection(t, engine, 0, "one_water_power_attack")
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "manes_arbitration_skill" ||
		!candidateContains(engine.State.PendingAction.Candidates, waterSkill.InstanceID) {
		t.Fatalf("Manes should ask for a learned water skill, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, waterSkill.InstanceID)
	if waterSkill.PowerBonus != 3 || waterSkill.AttackBonus != 1 {
		t.Fatalf("Manes should empower one water skill, power=%d attack=%d", waterSkill.PowerBonus, waterSkill.AttackBonus)
	}
	if err := engine.validateSkillLearnPermissionModifiers(0, NewCardInstance(baseCard(t, "3221002"), 0, engine.State.TurnNumber)); err == nil {
		t.Fatal("Manes one-skill mode should prevent learning water skills")
	}
	if err := engine.validateSkillLearnPermissionModifiers(0, NewCardInstance(baseCard(t, "3121001"), 0, engine.State.TurnNumber)); err != nil {
		t.Fatalf("Manes one-skill mode should not prevent non-water skills: %v", err)
	}
}

func TestRoyalConflictDeepSwordRevealsOnDrawWhenEnemySpellPowerIsHigher(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	sword := NewCardInstance(baseCard(t, "2211101"), 0, engine.State.TurnNumber)
	ownSkill := readySkill(baseCard(t, "3121105"), 0)
	enemySkill := readySkill(baseCard(t, "3421107"), 1)
	p0.Skills[0] = ownSkill
	p1.Skills[0] = enemySkill
	front := placeUnit(baseCard(t, "1421101"), 1, 0, 0, engine)
	back := placeUnit(baseCard(t, "1421101"), 1, 0, 2, engine)
	front.CurrentLife = 5
	back.CurrentLife = 5
	p0.Hand = []*CardInstance{sword}

	if err := (Card2211101DeepSword{}).OnDraw(&EffectContext{
		Engine:     engine,
		Source:     sword,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"drawn_player": 0},
	}); err != nil {
		t.Fatalf("deep sword draw failed: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "deep_sword_reveal_damage" {
		t.Fatalf("deep sword should ask to reveal, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, sword.InstanceID)
	if !p0.RevealedHand[sword.InstanceID] {
		t.Fatalf("deep sword should reveal itself, revealed=%v", p0.RevealedHand)
	}
	if front.CurrentLife != 3 {
		t.Fatalf("deep sword should damage enemies in spell range, front life=%d", front.CurrentLife)
	}
	if back.CurrentLife != 5 {
		t.Fatalf("deep sword should not damage enemies outside spell range, back life=%d", back.CurrentLife)
	}
}

func TestRoyalConflictDeepSwordDoesNotTriggerWithoutEnemySpellPowerLead(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	sword := NewCardInstance(baseCard(t, "2211101"), 0, engine.State.TurnNumber)
	p0.Skills[0] = readySkill(baseCard(t, "3421107"), 0)
	p1.Skills[0] = readySkill(baseCard(t, "3121105"), 1)
	p0.Hand = []*CardInstance{sword}
	placeUnit(baseCard(t, "1421101"), 1, 0, 0, engine)

	if err := (Card2211101DeepSword{}).OnDraw(&EffectContext{
		Engine:     engine,
		Source:     sword,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"drawn_player": 0},
	}); err != nil {
		t.Fatalf("deep sword draw failed: %v", err)
	}
	if engine.State.PendingAction != nil {
		t.Fatalf("deep sword should not trigger without enemy power lead, pending=%+v", engine.State.PendingAction)
	}
}

func TestRoyalConflictTimeCycleLocksPlayerActionsWhileActive(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	timeCycle := readySkill(baseCard(t, "3411101"), 0)
	p0.Skills[0] = timeCycle
	companion := NewCardInstance(baseCard(t, "1021001"), 0, engine.State.TurnNumber)
	equipment := NewCardInstance(baseCard(t, "2211101"), 0, engine.State.TurnNumber)
	item := NewCardInstance(baseCard(t, "2121002"), 0, engine.State.TurnNumber)
	p0.Hand = []*CardInstance{companion, equipment, item}
	p0.SkillPool = []*CardInstance{NewCardInstance(baseCard(t, "3121105"), 0, engine.State.TurnNumber)}
	attacker := placeUnit(baseCard(t, "1121001"), 0, 1, 0, engine)
	attacker.IsHorizontal = false
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	_ = target
	castSkill := readySkill(baseCard(t, "3121105"), 0)
	p0.Skills[1] = castSkill

	if !engine.timeCycleLockActive() {
		t.Fatal("time cycle should be active while 3411101 is a non-petrified learned skill")
	}
	if err := engine.handleSummon(0, ActionMessage{Data: map[string]any{
		"instance_id": companion.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err == nil {
		t.Fatal("time cycle should prevent summoning companions")
	}
	if err := engine.handleEquip(0, ActionMessage{Data: map[string]any{"instance_id": equipment.InstanceID}}); err == nil {
		t.Fatal("time cycle should prevent equipping cards")
	}
	if err := engine.handleUseItem(0, ActionMessage{Data: map[string]any{"instance_id": item.InstanceID}}); err == nil {
		t.Fatal("time cycle should prevent using items")
	}
	if err := engine.handleLearnSkill(0, ActionMessage{Data: map[string]any{"instance_id": p0.SkillPool[0].InstanceID}}); err == nil {
		t.Fatal("time cycle should prevent learning skills")
	}
	if err := engine.handleCastSpell(0, ActionMessage{Data: map[string]any{
		"instance_id": castSkill.InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err == nil {
		t.Fatal("time cycle should prevent using skills")
	}
	if err := engine.handleAttack(0, ActionMessage{Data: map[string]any{
		"attacker_id": attacker.InstanceID,
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err == nil {
		t.Fatal("time cycle should prevent card attacks")
	}

	timeCycle.Statuses[StatusPetrify] = 1
	if engine.timeCycleLockActive() {
		t.Fatal("petrified time cycle should stop locking actions")
	}
}

func TestRoyalConflictBurrowAddsMasteryExtraSplashTargets(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	burrow := readySkill(baseCard(t, "3421107"), 0)
	burrow.Statuses[StatusMastery] = 2
	p0.Skills[0] = burrow
	p0.Elements[model.ElementEarth] = 10
	main := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	extraA := placeUnit(baseCard(t, "1021002"), 1, 0, 2, engine)
	extraB := placeUnit(baseCard(t, "1021004"), 1, 2, 2, engine)

	if err := engine.handleCastSpell(0, ActionMessage{Data: map[string]any{
		"instance_id": burrow.InstanceID,
		"target_type": "unit",
		"target_col":  float64(main.Position.Col),
		"target_row":  float64(main.Position.Row),
		"extra_targets": []any{
			map[string]any{"col": float64(extraA.Position.Col), "row": float64(extraA.Position.Row)},
			map[string]any{"col": float64(extraB.Position.Col), "row": float64(extraB.Position.Row)},
		},
	}}); err != nil {
		t.Fatalf("cast 3421107 with mastery extra targets: %v", err)
	}
	if engine.State.PendingSpell == nil || len(engine.State.PendingSpell.ExtraTargets) != 2 {
		t.Fatalf("3421107 should carry two mastery extra targets, pending=%+v", engine.State.PendingSpell)
	}
	if got := engine.State.PendingSpell.ExtraTargets[0].Position; got != *extraA.Position {
		t.Fatalf("first extra target = %+v, want %+v", got, *extraA.Position)
	}
	if got := engine.State.PendingSpell.ExtraTargets[1].Position; got != *extraB.Position {
		t.Fatalf("second extra target = %+v, want %+v", got, *extraB.Position)
	}

	failEngine := setupEffectTest(t)
	failP0 := failEngine.State.Players[0]
	failBurrow := readySkill(baseCard(t, "3421107"), 0)
	failP0.Skills[0] = failBurrow
	failP0.Elements[model.ElementEarth] = 10
	failMain := placeUnit(baseCard(t, "1021001"), 1, 1, 0, failEngine)
	failExtra := placeUnit(baseCard(t, "1021002"), 1, 0, 2, failEngine)
	if err := failEngine.handleCastSpell(0, ActionMessage{Data: map[string]any{
		"instance_id": failBurrow.InstanceID,
		"target_type": "unit",
		"target_col":  float64(failMain.Position.Col),
		"target_row":  float64(failMain.Position.Row),
		"extra_targets": []any{
			map[string]any{"col": float64(failExtra.Position.Col), "row": float64(failExtra.Position.Row)},
		},
	}}); err == nil {
		t.Fatal("3421107 should reject mastery extra targets before reaching mastery")
	}
	if failEngine.State.PendingSpell != nil || failBurrow.IsHorizontal {
		t.Fatalf("invalid 3421107 extra target should not spend state, pending=%+v horizontal=%v", failEngine.State.PendingSpell, failBurrow.IsHorizontal)
	}
}

func TestRoyalConflictCursedFireFlipsCheapFireSpellScrollForFree(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	cursedFire := readySkill(baseCard(t, "3121110"), 0)
	p0.Skills[0] = cursedFire
	p0.Elements[model.ElementFire] = 10
	nonScroll := NewCardInstance(baseCard(t, "3121105"), 0, 1)
	tooExpensive := NewCardInstance(baseCard(t, "2121104"), 0, 1)
	scroll := NewCardInstance(baseCard(t, "2121112"), 0, 1)
	p0.Deck = []*CardInstance{nonScroll, tooExpensive, scroll}
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

	if err := engine.handleCastSpell(0, ActionMessage{Data: map[string]any{
		"instance_id": cursedFire.InstanceID,
		"target_type": "none",
	}}); err != nil {
		t.Fatalf("cast 3121110: %v", err)
	}
	if !containsCardInstance(p0.Hand, scroll) || containsCardInstance(p0.Hand, nonScroll) || containsCardInstance(p0.Hand, tooExpensive) {
		t.Fatalf("3121110 should flip first cheap fire spell scroll only, hand=%v deck=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Deck))
	}
	if cost := engine.effectiveCardPlayCost(p0, scroll); totalElementCost(cost) != 0 {
		t.Fatalf("3121110 flipped scroll should be free to use, cost=%v statuses=%v", cost, scroll.Statuses)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "cursed_fire_use_scroll_target" {
		t.Fatalf("3121110 should offer to immediately use the flipped scroll, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, target.InstanceID)
	if containsCardInstance(p0.Hand, scroll) || engine.State.PendingSpell == nil || engine.State.PendingSpell.Skill != scroll {
		t.Fatalf("3121110 accepted immediate use should cast the flipped scroll, hand=%v pending=%+v", cardsToInfo(p0.Hand), engine.State.PendingSpell)
	}
}

func TestRoyalConflictJadeFacedSnowFoxMovesAndForcesRetarget(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	spell := readySkill(baseCard(t, "3121105"), 0)
	p0.Skills[0] = spell
	p0.Elements[model.ElementFire] = 3
	fox := placeUnit(baseCard(t, "1221114"), 1, 1, 0, engine)
	other := placeUnit(baseCard(t, "1021001"), 1, 2, 0, engine)

	if err := engine.handleCastSpell(0, ActionMessage{Data: map[string]any{
		"instance_id": spell.InstanceID,
		"target_type": "unit",
		"target_col":  float64(fox.Position.Col),
		"target_row":  float64(fox.Position.Row),
	}}); err != nil {
		t.Fatalf("cast at snow fox: %v", err)
	}
	if engine.State.Phase != PhaseDefenseWindow || engine.State.PendingSpell == nil {
		t.Fatalf("expected defense window pending spell, phase=%s spell=%+v", engine.State.Phase, engine.State.PendingSpell)
	}
	if err := engine.handleReactSpell(1, ActionMessage{Data: map[string]any{"instance_id": fox.InstanceID}}); err != nil {
		t.Fatalf("snow fox reaction: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "jade_faced_snow_fox_move" {
		t.Fatalf("snow fox should ask for movement, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 1, positionSelectionID(Position{Col: 0, Row: 2}))
	if p1.Elements[model.ElementWater] != 2 || !fox.UltimateUsed || fox.Position == nil || *fox.Position != (Position{Col: 0, Row: 2}) {
		t.Fatalf("snow fox should move and gain water, elements=%v used=%v pos=%+v", p1.Elements, fox.UltimateUsed, fox.Position)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "jade_faced_snow_fox_retarget" || engine.State.PendingAction.PlayerID != 0 {
		t.Fatalf("snow fox should ask attacker to retarget, pending=%+v", engine.State.PendingAction)
	}
	if engine.State.PendingAction.Prompt != "玉面雪狐:重新选择法术攻击目标" {
		t.Fatalf("snow fox retarget prompt should name snow fox, prompt=%q", engine.State.PendingAction.Prompt)
	}
	resolvePendingSelection(t, engine, 0, other.InstanceID)
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.Target.Position != *other.Position {
		t.Fatalf("snow fox should update pending spell target, spell=%+v other=%+v", engine.State.PendingSpell, other.Position)
	}
}

func TestRoyalConflictScatterAwayProtectsDamagedAirCompanion(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	counter := NewCardInstance(baseCard(t, "2321108"), 0, engine.State.TurnNumber)
	counter.IsSetCounter = true
	counter.IsHorizontal = true
	counter.SlotIndex = 0
	p0.Equipment[0] = counter
	p0.Elements[model.ElementAir] = 2
	ally := placeUnit(baseCard(t, "1321004"), 0, 1, 0, engine)
	startLife := ally.CurrentLife

	engine.dealDamageWithExtra(ally, 1, 0, map[string]any{"damage_source": "test", "attacker": 1})
	if ally.CurrentLife != startLife-1 || engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_trigger" {
		t.Fatalf("2321108 should trigger after friendly air damage, life=%d/%d pending=%+v", ally.CurrentLife, startLife, engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, counter.InstanceID)
	if ally.Statuses[temporaryDamageAndNegativeImmunityUntilStatus] == 0 || !containsCardInstance(p0.Graveyard, counter) {
		t.Fatalf("2321108 should grant temporary immunity and go to graveyard, statuses=%v grave=%v", ally.Statuses, cardsToInfo(p0.Graveyard))
	}
	protectedLife := ally.CurrentLife
	engine.dealDamageWithExtra(ally, 2, 0, map[string]any{"damage_source": "test", "attacker": 1})
	if ally.CurrentLife != protectedLife {
		t.Fatalf("2321108 immunity should prevent later damage, life=%d want=%d", ally.CurrentLife, protectedLife)
	}
}

func TestRoyalConflictInfusionRuneEResetsDefenseSpells(t *testing.T) {
	engine := setupEffectTest(t)
	defense := readySkill(baseCard(t, "3221004"), 0)
	boost := readySkill(baseCard(t, "3221003"), 0)
	defense.IsHorizontal = true
	boost.IsHorizontal = true
	defense.Statuses[StatusCooldown] = 1
	boost.Statuses[StatusCooldown] = 1
	rune := NewCardInstance(baseCard(t, "2021115"), 0, engine.State.TurnNumber)

	if err := (Card2021115InfusionRuneE{}).OnDefend(&EffectContext{
		Engine:   engine,
		Source:   rune,
		PlayerID: 0,
		ExtraData: map[string]any{
			"defender":       0,
			"defense_skills": []*CardInstance{defense},
			"defense_boosts": []*CardInstance{boost},
		},
	}); err != nil {
		t.Fatalf("2021115 defend: %v", err)
	}
	if defense.IsHorizontal || boost.IsHorizontal {
		t.Fatalf("2021115 should reset defense and defense-boost spells through cooldown, defense=%v boost=%v", defense.IsHorizontal, boost.IsHorizontal)
	}
}

func TestRoyalConflictFireRebirthScrollRevivesFireCompanionsThatDiedThisTurn(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	engine.State.CurrentTurn = 1
	setAllElements(p0, 10)

	scroll := NewCardInstance(baseCard(t, "2121104"), 0, engine.State.TurnNumber)
	scroll.IsSetCounter = true
	scroll.IsHorizontal = true
	scroll.SlotIndex = 0
	p0.Equipment[0] = scroll

	fire := placeUnit(baseCard(t, "1121102"), 0, 1, 0, engine)
	fire2 := placeUnit(baseCard(t, "1121113"), 0, 2, 0, engine)
	engine.destroyUnitWithData(fire, 0, map[string]any{"attacker": 1, "death_cause": "test"})
	engine.destroyUnitWithData(fire2, 0, map[string]any{"attacker": 1, "death_cause": "test"})
	if len(p0.Graveyard) < 2 || p0.Graveyard[0].InstanceID != fire.InstanceID || p0.Graveyard[1].InstanceID != fire2.InstanceID {
		t.Fatalf("fire companion should be in graveyard before rebirth, grave=%v", cardsToInfo(p0.Graveyard))
	}

	triggered := engine.triggerFieldEffectsWithData(TriggerOnTurnEnd, 0, nil, map[string]any{"ended_player": 1})
	if !triggered || engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_trigger" {
		t.Fatalf("2121104 should prompt at enemy turn end, triggered=%v pending=%+v", triggered, engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, scroll.InstanceID)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "fire_rebirth_scroll" || len(engine.State.PendingAction.Candidates) != 2 {
		t.Fatalf("2121104 should ask which recent fire companion to revive, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, fire2.InstanceID)
	revivePos := Position{Col: 0, Row: 1}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "fire_rebirth_scroll_position" || !candidateContains(engine.State.PendingAction.Candidates, positionSelectionID(revivePos)) {
		t.Fatalf("2121104 should ask which empty unit position to revive into, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, positionSelectionID(revivePos))

	if p0.Units[revivePos.Col][revivePos.Row] != fire2 || fire2.IsHorizontal || fire2.CurrentLife != fire2.Card.Life {
		t.Fatalf("2121104 should revive selected fire companion vertical at chosen position, pos=%+v horizontal=%v life=%d/%d", fire2.Position, fire2.IsHorizontal, fire2.CurrentLife, fire2.Card.Life)
	}
	if !containsCardInstance(p0.Graveyard, fire) || containsCardInstance(p0.Graveyard, fire2) || !containsCardInstance(p0.Graveyard, scroll) {
		t.Fatalf("2121104 should revive only selected fire and discard scroll, grave=%v", cardsToInfo(p0.Graveyard))
	}
	if fire2.Statuses[enteredGraveyardTurnStatus] != 0 {
		t.Fatalf("revived fire should not keep internal graveyard turn status, statuses=%v", fire2.Statuses)
	}

	staleEngine := setupEffectTest(t)
	staleP0 := staleEngine.State.Players[0]
	staleScroll := NewCardInstance(baseCard(t, "2121104"), 0, staleEngine.State.TurnNumber)
	staleScroll.IsSetCounter = true
	staleScroll.IsHorizontal = true
	staleScroll.SlotIndex = 0
	staleP0.Equipment[0] = staleScroll
	setAllElements(staleP0, 10)
	staleFire := placeUnit(baseCard(t, "1121102"), 0, 1, 0, staleEngine)
	staleEngine.destroyUnitWithData(staleFire, 0, map[string]any{"attacker": 1, "death_cause": "test"})
	staleEngine.State.TurnNumber++
	if triggered := staleEngine.triggerFieldEffectsWithData(TriggerOnTurnEnd, 0, nil, map[string]any{"ended_player": 1}); triggered || staleEngine.State.PendingAction != nil {
		t.Fatalf("2121104 should ignore fire companions from older turns, triggered=%v pending=%+v", triggered, staleEngine.State.PendingAction)
	}

	fullEngine := setupEffectTest(t)
	fullP0 := fullEngine.State.Players[0]
	fullEngine.State.CurrentTurn = 1
	fullScroll := NewCardInstance(baseCard(t, "2121104"), 0, fullEngine.State.TurnNumber)
	fullScroll.IsSetCounter = true
	fullScroll.IsHorizontal = true
	fullScroll.SlotIndex = 0
	fullP0.Equipment[0] = fullScroll
	setAllElements(fullP0, 10)
	fullFire := placeUnit(baseCard(t, "1121102"), 0, 1, 0, fullEngine)
	fullEngine.destroyUnitWithData(fullFire, 0, map[string]any{"attacker": 1, "death_cause": "test"})
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if fullP0.Units[col][row] == nil {
				placeUnit(baseCard(t, "1021001"), 0, col, row, fullEngine)
			}
		}
	}
	if triggered := fullEngine.triggerFieldEffectsWithData(TriggerOnTurnEnd, 0, nil, map[string]any{"ended_player": 1}); triggered || fullEngine.State.PendingAction != nil {
		t.Fatalf("2121104 should not trigger when no friendly unit slot is empty, triggered=%v pending=%+v", triggered, fullEngine.State.PendingAction)
	}
	if !containsCardInstance(fullP0.Graveyard, fullFire) || fullP0.Equipment[0] != fullScroll {
		t.Fatalf("2121104 should not consume scroll or revive when board is full, equipment=%v grave=%v", cardToInfo(fullP0.Equipment[0]), cardsToInfo(fullP0.Graveyard))
	}
}

func TestRoyalConflictLavaArmorYeYanSacrificesAndEquipsMoltenArmorAfterShieldBreak(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	armor := NewCardInstance(baseCard(t, "2111102"), 0, engine.State.TurnNumber)
	armor.SlotIndex = 0
	p0.Equipment[0] = armor
	molten := NewCardInstance(baseCard(t, "2121013"), 0, engine.State.TurnNumber)
	p0.Hand = []*CardInstance{molten}
	p0.Shield = 1

	remaining := engine.applyPlayerShieldDamage(target, 2, map[string]any{"damage_source": "spell", "attacker": 1})
	if remaining != 1 || p0.Shield != 0 || !p0.ShieldBrokenThisTurn {
		t.Fatalf("spell damage should break player shield, remaining=%d shield=%d broken=%v", remaining, p0.Shield, p0.ShieldBrokenThisTurn)
	}
	if err := (Card2111102LavaArmorYeYan{}).OnSpellHit(&EffectContext{
		Engine:     engine,
		Source:     armor,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"attacker": 1, "spell_source": readySkill(baseCard(t, "3021005"), 1)},
	}); err != nil {
		t.Fatalf("2111102 spell hit: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "lava_armor_yeyan_sacrifice" {
		t.Fatalf("2111102 should ask whether to sacrifice after enemy spell hit, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, armor.InstanceID)
	if p0.Equipment[0] != molten || !molten.IsHorizontal || len(p0.Hand) != 0 || !containsCardInstance(p0.Graveyard, armor) {
		t.Fatalf("2111102 should equip molten armor from hand after shield break, equipment=%v hand=%v grave=%v", cardToInfo(p0.Equipment[0]), cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard))
	}
	if p0.Shield != 2 {
		t.Fatalf("2111102 should gain shield 2 after sacrifice, shield=%d", p0.Shield)
	}

	noBreakEngine := setupEffectTest(t)
	noBreakP0 := noBreakEngine.State.Players[0]
	noBreakTarget := placeUnit(baseCard(t, "1021001"), 0, 0, 0, noBreakEngine)
	noBreakArmor := NewCardInstance(baseCard(t, "2111102"), 0, noBreakEngine.State.TurnNumber)
	noBreakArmor.SlotIndex = 0
	noBreakP0.Equipment[0] = noBreakArmor
	deckMolten := NewCardInstance(baseCard(t, "2121013"), 0, noBreakEngine.State.TurnNumber)
	noBreakP0.Deck = []*CardInstance{deckMolten}
	if err := (Card2111102LavaArmorYeYan{}).OnSpellHit(&EffectContext{
		Engine:     noBreakEngine,
		Source:     noBreakArmor,
		PlayerID:   0,
		OpponentID: 1,
		ExtraData:  map[string]any{"attacker": 1, "spell_source": readySkill(baseCard(t, "3021005"), 1)},
	}); err != nil {
		t.Fatalf("2111102 no-break spell hit: %v", err)
	}
	resolvePendingSelection(t, noBreakEngine, 0, noBreakArmor.InstanceID)
	if noBreakP0.Equipment[0] != nil || len(noBreakP0.Deck) != 1 || noBreakP0.Shield != 2 || len(noBreakP0.TempModifiers) != 1 || noBreakP0.TempModifiers[0].Type != TempModLavaArmorYeYanShieldBreak {
		t.Fatalf("2111102 should wait for later shield break after sacrifice, equipment=%v deck=%v shield=%d modifiers=%+v", cardToInfo(noBreakP0.Equipment[0]), cardsToInfo(noBreakP0.Deck), noBreakP0.Shield, noBreakP0.TempModifiers)
	}
	laterRemaining := noBreakEngine.applyPlayerShieldDamage(noBreakTarget, 3, map[string]any{"damage_source": "spell", "attacker": 1})
	if laterRemaining != 1 || !noBreakP0.ShieldBrokenThisTurn || noBreakP0.Equipment[0] != deckMolten || len(noBreakP0.Deck) != 0 || len(noBreakP0.TempModifiers) != 0 {
		t.Fatalf("2111102 should equip molten armor from deck after later shield break, remaining=%d broken=%v equipment=%v deck=%v modifiers=%+v", laterRemaining, noBreakP0.ShieldBrokenThisTurn, cardToInfo(noBreakP0.Equipment[0]), cardsToInfo(noBreakP0.Deck), noBreakP0.TempModifiers)
	}
}

func TestRoyalConflictErebosSoulChainMarksOverexertSpellAndWeakensOnConsumeOrOverexert(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	chain := NewCardInstance(baseCard(t, "2611101"), 0, engine.State.TurnNumber)
	p0.Equipment[0] = chain
	payer := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	spell := readySkill(baseCard(t, "3121001"), 1)
	boost := readySkill(baseCard(t, "3321001"), 1)
	p1.Skills[0] = spell
	p1.Skills[1] = boost

	engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, spell, map[string]any{
		"cast_player":     1,
		"attacker":        1,
		"overexert_units": []*CardInstance{payer},
		"boost_skills":    []*CardInstance{boost},
	})
	if payer.Statuses[erebosSoulChainMarkedUnitStatus] == 0 || spell.Statuses[erebosSoulChainMarkedSpellStatus] == 0 || boost.Statuses[erebosSoulChainMarkedSpellStatus] == 0 {
		t.Fatalf("2611101 should mark overexerted companion and used spells, payer=%v spell=%v boost=%v", payer.Statuses, spell.Statuses, boost.Statuses)
	}
	if !chain.UltimateUsed {
		t.Fatalf("2611101 should spend its triggered ultimate after marking, used=%v", chain.UltimateUsed)
	}
	secondPayer := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
	secondSpell := readySkill(baseCard(t, "3121002"), 1)
	engine.triggerFieldEffectsWithData(TriggerOnSpellCast, 0, secondSpell, map[string]any{
		"cast_player":     1,
		"attacker":        1,
		"overexert_units": []*CardInstance{secondPayer},
	})
	if secondPayer.Statuses[erebosSoulChainMarkedUnitStatus] != 0 || secondSpell.Statuses[erebosSoulChainMarkedSpellStatus] != 0 {
		t.Fatalf("2611101 should not mark another spell after its ultimate is spent, payer=%v spell=%v", secondPayer.Statuses, secondSpell.Statuses)
	}

	engine.triggerFieldEffectsWithData(TriggerOnConsume, 0, payer, map[string]any{"consumed_player": 1})
	if spell.Statuses[StatusWeaken] != 1 || boost.Statuses[StatusWeaken] != 1 {
		t.Fatalf("2611101 should weaken marked spells when marked companion is consumed, spell=%v boost=%v", spell.Statuses, boost.Statuses)
	}
	engine.triggerErebosSoulChainMarkedOverexert(1, []*CardInstance{payer})
	if spell.Statuses[StatusWeaken] != 2 || boost.Statuses[StatusWeaken] != 2 {
		t.Fatalf("2611101 should weaken marked spells when marked companion is overexerted, spell=%v boost=%v", spell.Statuses, boost.Statuses)
	}

	unmarked := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
	engine.triggerErebosSoulChainMarkedOverexert(1, []*CardInstance{unmarked})
	if spell.Statuses[StatusWeaken] != 2 || boost.Statuses[StatusWeaken] != 2 {
		t.Fatalf("2611101 should ignore unmarked overexerted companions, spell=%v boost=%v", spell.Statuses, boost.Statuses)
	}
}

func TestRoyalConflictLampusSwordDelaysAndDistributesAirDiscardDamage(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	sword := NewCardInstance(baseCard(t, "2311102"), 0, engine.State.TurnNumber)
	sword.SlotIndex = 0
	p0.Equipment[0] = sword
	airA := NewCardInstance(baseCard(t, "3321001"), 0, engine.State.TurnNumber)
	airB := NewCardInstance(baseCard(t, "3321002"), 0, engine.State.TurnNumber)
	fire := NewCardInstance(baseCard(t, "3121001"), 0, engine.State.TurnNumber)
	p0.Hand = []*CardInstance{airA, airB, fire}
	front := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	back := placeUnit(baseCard(t, "1021002"), 1, 1, 2, engine)
	frontStart := front.CurrentLife
	backStart := back.CurrentLife

	if err := (Card2311102LampusSword{}).OnPerTurn(&EffectContext{Engine: engine, Source: sword, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("2311102 ability: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "lampus_sword_discard_air" {
		t.Fatalf("2311102 should ask for air hand cards, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, airA.InstanceID, airB.InstanceID)
	if p0.Equipment[0] != nil || !containsCardInstance(p0.Graveyard, sword) || len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModLampusSwordDelayedDamage || p0.TempModifiers[0].Amount != 2 {
		t.Fatalf("2311102 should sacrifice and store delayed damage, equipment=%v grave=%v mods=%+v", p0.Equipment[0], cardsToInfo(p0.Graveyard), p0.TempModifiers)
	}
	if len(p0.Hand) != 1 || p0.Hand[0] != fire {
		t.Fatalf("2311102 should discard selected air cards only, hand=%v", cardsToInfo(p0.Hand))
	}

	engine.applyOpponentTurnEndTemporaryModifiers(1)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "lampus_sword_distribute_damage" {
		t.Fatalf("2311102 should ask to distribute damage at opponent turn end, pending=%+v", engine.State.PendingAction)
	}
	for _, candidate := range engine.State.PendingAction.Candidates {
		if candidate["instance_id"] == back.InstanceID {
			t.Fatalf("2311102 should only offer enemies in spell range, candidates=%+v", engine.State.PendingAction.Candidates)
		}
	}
	resolvePendingSelection(t, engine, 0, front.InstanceID)
	if front.CurrentLife != frontStart-2 || back.CurrentLife != backStart || len(p0.TempModifiers) != 0 {
		t.Fatalf("2311102 should deal stored damage and clear modifier, front=%d/%d back=%d/%d mods=%+v", front.CurrentLife, frontStart, back.CurrentLife, backStart, p0.TempModifiers)
	}
}

func TestRoyalConflictBloodSandArrayPaymentsCreateMarkersAndModifyStats(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	array := readySkill(baseCard(t, "3411102"), 0)
	p0.Skills[0] = array
	ally := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	setElementsGain(ally, map[string]int{model.ElementEarth: 2})
	enemy := placeUnit(baseCard(t, "1021002"), 1, 0, 0, engine)
	enemyStartLife := enemy.CurrentLife

	if err := (Card3411102BloodSandArray{}).OnPerTurn(&EffectContext{Engine: engine, Source: array, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("3411102 ability: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_sand_array_pay" {
		t.Fatalf("3411102 should ask owner payment first, pending=%+v", engine.State.PendingAction)
	}
	if err := resolvePendingSelectionWithData(engine, 0, nil, map[string]any{
		"payments": []any{map[string]any{"instance_id": ally.InstanceID, "amount": float64(2), "mode": "load"}},
	}); err != nil {
		t.Fatalf("resolve 3411102 owner payment: %v", err)
	}
	if effectiveElementsGain(ally)[model.ElementEarth] != 0 {
		t.Fatalf("3411102 should remove selected load from owner unit, load=%v", effectiveElementsGain(ally))
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "blood_sand_array_pay_opponent" || engine.State.PendingAction.PlayerID != 1 {
		t.Fatalf("3411102 should ask opponent payment second, pending=%+v", engine.State.PendingAction)
	}
	if err := resolvePendingSelectionWithData(engine, 1, nil, map[string]any{
		"payments": []any{map[string]any{"instance_id": enemy.InstanceID, "amount": float64(1), "mode": "life"}},
	}); err != nil {
		t.Fatalf("resolve 3411102 opponent payment: %v", err)
	}
	if enemy.CurrentLife != enemyStartLife-1 || array.Statuses[bloodSandArrayMarkerStatus] != 1 {
		t.Fatalf("3411102 should pay opponent life and gain difference markers, life=%d/%d statuses=%v", enemy.CurrentLife, enemyStartLife, array.Statuses)
	}
	stats := engine.skillContributionStats(0, array, nil, skillPurposeAttack)
	if stats.PowerBonus != array.Card.Power+3 || stats.DamageBonus != 1 {
		t.Fatalf("3411102 markers should add +3 power and +1 damage, stats=%+v basePower=%d", stats, array.Card.Power)
	}
}

func TestRoyalConflictFlameArrayScrollRequiresFireSacrificeAndAddsPower(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	scroll := NewCardInstance(baseCard(t, "2121105"), 0, engine.State.TurnNumber)
	p0.Hand = []*CardInstance{scroll}
	setAllElements(p0, 9)
	sacrifice := placeUnit(baseCard(t, "1121102"), 0, 0, 0, engine)
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	bonus := totalElementCost(sacrifice.Card.ElementsCost)

	if err := engine.HandleAction(0, ActionMessage{
		Action: "use_item",
		Data: map[string]any{
			"instance_id":  scroll.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(target.Position.Col),
			"target_row":   float64(target.Position.Row),
			"sacrifice_id": sacrifice.InstanceID,
		},
	}); err != nil {
		t.Fatalf("use 2121105: %v", err)
	}
	if engine.State.PendingAction != nil {
		resolvePendingSelection(t, engine, engine.State.PendingAction.PlayerID)
	}
	power, ok := lastSpellCastPower(engine, scroll.InstanceID)
	if !ok {
		t.Fatalf("2121105 should emit spell_cast, events=%v", engine.log)
	}
	if power != scroll.Card.Power+bonus {
		t.Fatalf("2121105 should add sacrificed entry cost to power, power=%d want=%d", power, scroll.Card.Power+bonus)
	}
	if p0.Units[0][0] != nil || !containsCardInstance(p0.Graveyard, sacrifice) {
		t.Fatalf("2121105 should sacrifice selected fire companion, unit=%v grave=%v", p0.Units[0][0], cardsToInfo(p0.Graveyard))
	}

	failEngine := setupEffectTest(t)
	failP0 := failEngine.State.Players[0]
	failScroll := NewCardInstance(baseCard(t, "2121105"), 0, failEngine.State.TurnNumber)
	failP0.Hand = []*CardInstance{failScroll}
	setAllElements(failP0, 9)
	nonFire := placeUnit(baseCard(t, "1021001"), 0, 0, 0, failEngine)
	failTarget := placeUnit(baseCard(t, "1021001"), 1, 0, 0, failEngine)
	if err := failEngine.HandleAction(0, ActionMessage{
		Action: "use_item",
		Data: map[string]any{
			"instance_id":  failScroll.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(failTarget.Position.Col),
			"target_row":   float64(failTarget.Position.Row),
			"sacrifice_id": nonFire.InstanceID,
		},
	}); err == nil {
		t.Fatalf("2121105 should reject non-fire sacrifice")
	}
	if len(failP0.Hand) != 1 || failP0.Units[0][0] != nonFire || len(failP0.Graveyard) != 0 {
		t.Fatalf("2121105 invalid sacrifice should not spend state, hand=%v unit=%v grave=%v", cardsToInfo(failP0.Hand), cardToInfo(failP0.Units[0][0]), cardsToInfo(failP0.Graveyard))
	}
}

func TestRoyalConflictForesightOrbPreviewsAndReordersDeckTop(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	orb := NewCardInstance(baseCard(t, "2011102"), 0, engine.State.TurnNumber)
	orb.IsHorizontal = false
	p0.Equipment[0] = orb
	a := NewCardInstance(baseCard(t, "1021001"), 0, 1)
	b := NewCardInstance(baseCard(t, "1021002"), 0, 1)
	c := NewCardInstance(baseCard(t, "1021003"), 0, 1)
	d := NewCardInstance(baseCard(t, "1021004"), 0, 1)
	p0.Deck = []*CardInstance{a, b, c, d}

	state := engine.playerStateToInfo(p0, true)
	preview, ok := state["top_deck_preview"].([]map[string]any)
	if !ok || len(preview) != 3 || preview[0]["number"] != a.Card.Number || preview[2]["number"] != c.Card.Number {
		t.Fatalf("2011102 should expose owner top three preview, preview=%v", state["top_deck_preview"])
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  orb.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use 2011102 ability: %v", err)
	}
	if !orb.IsHorizontal || engine.State.PendingAction == nil || engine.State.PendingAction.Type != "foresight_orb_reorder" {
		t.Fatalf("2011102 should consume itself and ask for deck reorder, horizontal=%v pending=%+v", orb.IsHorizontal, engine.State.PendingAction)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected":     []any{},
		"top_order":    []any{c.InstanceID, a.InstanceID},
		"bottom_order": []any{b.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve 2011102 reorder: %v", err)
	}
	if p0.Deck[0] != c || p0.Deck[1] != a || p0.Deck[2] != d || p0.Deck[3] != b {
		t.Fatalf("2011102 deck order wrong, deck=%v", cardsToInfo(p0.Deck))
	}
}

func TestRoyalConflictCouncilGuardReordersOpponentTopFiveOnlyWithMark(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	guard := placeUnit(baseCard(t, "1521112"), 0, 0, 0, engine)
	a := NewCardInstance(baseCard(t, "1021001"), 1, 1)
	mark := NewCardInstance(baseCard(t, "2001102"), 1, 1)
	c := NewCardInstance(baseCard(t, "1021002"), 1, 1)
	d := NewCardInstance(baseCard(t, "1021003"), 1, 1)
	e := NewCardInstance(baseCard(t, "1021004"), 1, 1)
	rest := NewCardInstance(baseCard(t, "1021005"), 1, 1)
	p1.Deck = []*CardInstance{a, mark, c, d, e, rest}
	startHandCount := len(p0.Hand)

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  guard.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use 1521112 ability: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "council_guard_reorder" || engine.State.PendingAction.PlayerID != 0 {
		t.Fatalf("1521112 should ask controller to reorder opponent deck, pending=%+v", engine.State.PendingAction)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected":     []any{},
		"top_order":    []any{e.InstanceID, mark.InstanceID},
		"bottom_order": []any{a.InstanceID, d.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve 1521112 reorder: %v", err)
	}
	if p1.Deck[0] != e || p1.Deck[1] != mark || p1.Deck[2] != c || p1.Deck[3] != rest || p1.Deck[4] != a || p1.Deck[5] != d {
		t.Fatalf("1521112 opponent deck order wrong, deck=%v", cardsToInfo(p1.Deck))
	}
	if len(p0.Hand) != startHandCount {
		t.Fatalf("1521112 should only reorder, not draw revealed cards")
	}

	noMarkEngine := setupEffectTest(t)
	noMarkP1 := noMarkEngine.State.Players[1]
	noMarkGuard := placeUnit(baseCard(t, "1521112"), 0, 0, 0, noMarkEngine)
	noMarkP1.Deck = []*CardInstance{
		NewCardInstance(baseCard(t, "1021001"), 1, 1),
		NewCardInstance(baseCard(t, "1021002"), 1, 1),
		NewCardInstance(baseCard(t, "1021003"), 1, 1),
	}
	if err := noMarkEngine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  noMarkGuard.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use 1521112 no mark ability: %v", err)
	}
	if noMarkEngine.State.PendingAction != nil {
		t.Fatalf("1521112 should not open reorder prompt without a mark, pending=%+v", noMarkEngine.State.PendingAction)
	}
}

func TestRoyalConflictRoseProphetReordersOpponentTopThreeAfterShuffle(t *testing.T) {
	engine := setupEffectTest(t)
	p1 := engine.State.Players[1]
	placeUnit(baseCard(t, "1511103"), 0, 0, 0, engine)
	p1.Deck = []*CardInstance{
		NewCardInstance(baseCard(t, "1021001"), 1, 1),
		NewCardInstance(baseCard(t, "1021002"), 1, 1),
		NewCardInstance(baseCard(t, "1021003"), 1, 1),
		NewCardInstance(baseCard(t, "1021004"), 1, 1),
	}

	engine.shuffleDeck(1)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "rose_prophet_reorder" || engine.State.PendingAction.PlayerID != 0 {
		t.Fatalf("1511103 should ask controller to reorder opponent deck after shuffle, pending=%+v", engine.State.PendingAction)
	}
	candidates := engine.State.PendingAction.Candidates
	if len(candidates) != 3 {
		t.Fatalf("1511103 should reveal top three candidates, candidates=%+v", candidates)
	}
	firstID, _ := candidates[0]["instance_id"].(string)
	secondID, _ := candidates[1]["instance_id"].(string)
	thirdID, _ := candidates[2]["instance_id"].(string)
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected":     []any{},
		"top_order":    []any{thirdID, firstID},
		"bottom_order": []any{secondID},
	}}); err != nil {
		t.Fatalf("resolve 1511103 reorder: %v", err)
	}
	if p1.Deck[0].InstanceID != thirdID || p1.Deck[1].InstanceID != firstID || p1.Deck[len(p1.Deck)-1].InstanceID != secondID {
		t.Fatalf("1511103 deck order wrong, deck=%v", cardsToInfo(p1.Deck))
	}
}

func TestRoyalConflictOracleScrollGloryRequiresSupportAndAddsPower(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	scroll := NewCardInstance(baseCard(t, "2521111"), 0, engine.State.TurnNumber)
	p0.Hand = []*CardInstance{scroll}
	setAllElements(p0, 9)
	support := placeUnit(baseCard(t, "1521001"), 0, 0, 0, engine)
	support.CurrentLife = 4
	setElementsGain(support, map[string]int{model.ElementLight: 2})
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	bonus := support.CurrentLife + engine.totalLoad(support)

	if err := engine.HandleAction(0, ActionMessage{
		Action: "use_item",
		Data: map[string]any{
			"instance_id": scroll.InstanceID,
			"target_type": "unit",
			"target_col":  float64(target.Position.Col),
			"target_row":  float64(target.Position.Row),
			"support_id":  support.InstanceID,
		},
	}); err != nil {
		t.Fatalf("use 2521111: %v", err)
	}
	if engine.State.PendingAction != nil {
		resolvePendingSelection(t, engine, engine.State.PendingAction.PlayerID)
	}
	power, ok := lastSpellCastPower(engine, scroll.InstanceID)
	if !ok {
		t.Fatalf("2521111 should emit spell_cast, events=%v", engine.log)
	}
	if power != scroll.Card.Power+bonus {
		t.Fatalf("2521111 should add support life+load to power, power=%d want=%d", power, scroll.Card.Power+bonus)
	}

	failEngine := setupEffectTest(t)
	failP0 := failEngine.State.Players[0]
	failScroll := NewCardInstance(baseCard(t, "2521111"), 0, failEngine.State.TurnNumber)
	failP0.Hand = []*CardInstance{failScroll}
	setAllElements(failP0, 9)
	failTarget := placeUnit(baseCard(t, "1021001"), 1, 0, 0, failEngine)
	if err := failEngine.HandleAction(0, ActionMessage{
		Action: "use_item",
		Data: map[string]any{
			"instance_id": failScroll.InstanceID,
			"target_type": "unit",
			"target_col":  float64(failTarget.Position.Col),
			"target_row":  float64(failTarget.Position.Row),
		},
	}); err == nil {
		t.Fatalf("2521111 should require support_id")
	}
	if len(failP0.Hand) != 1 || len(failP0.Graveyard) != 0 {
		t.Fatalf("2521111 failed use should not leave hand, hand=%v grave=%v", cardsToInfo(failP0.Hand), cardsToInfo(failP0.Graveyard))
	}
}

func TestRoyalConflictCloudTopTradingHouseDrawsOpponentDeckAsNeutral(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	placeUnit(baseCard(t, "1311102"), 0, 0, 0, engine)
	stolen := NewCardInstance(baseCard(t, "1021002"), 1, engine.State.TurnNumber)
	p0.Deck = nil
	p1.Deck = []*CardInstance{stolen}

	drawn := engine.drawCards(0, 1)
	if len(drawn) != 1 || drawn[0] != stolen || len(p0.Hand) == 0 || p0.Hand[len(p0.Hand)-1] != stolen {
		t.Fatalf("1311102 should draw from opponent deck when own deck is empty, drawn=%v hand=%v", cardsToInfo(drawn), cardsToInfo(p0.Hand))
	}
	if len(p1.Deck) != 0 || stolen.OwnerID != 0 {
		t.Fatalf("1311102 should remove drawn card from opponent deck and transfer owner, p1deck=%v owner=%d", cardsToInfo(p1.Deck), stolen.OwnerID)
	}
	if got := engine.effectiveCardPlayCost(p0, stolen); got[model.ElementArcane] != totalElementCost(stolen.Card.ElementsCost) || len(got) != 1 {
		t.Fatalf("1311102 should convert entry cost to equivalent neutral/arcane, cost=%v original=%v", got, stolen.Card.ElementsCost)
	}
	if got := effectiveElementsGain(stolen); got[model.ElementArcane] != totalElementCost(stolen.Card.ElementsGain) || len(got) != 1 {
		t.Fatalf("1311102 should convert load to equivalent neutral/arcane, load=%v original=%v", got, stolen.Card.ElementsGain)
	}

	blockedEngine := setupEffectTest(t)
	blockedP0 := blockedEngine.State.Players[0]
	blockedP1 := blockedEngine.State.Players[1]
	blockedP0.Deck = nil
	blockedP1.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021002"), 1, blockedEngine.State.TurnNumber)}
	if drawn := blockedEngine.drawCards(0, 1); len(drawn) != 0 || len(blockedP1.Deck) != 1 {
		t.Fatalf("without 1311102, empty own deck should not draw opponent deck, drawn=%v p1deck=%v", cardsToInfo(drawn), cardsToInfo(blockedP1.Deck))
	}
}

func resolvePendingSelectionWithData(engine *Engine, playerID int, selected []string, extra map[string]any) error {
	values := make([]any, 0, len(selected))
	for _, id := range selected {
		values = append(values, id)
	}
	data := map[string]any{"selected": values}
	for key, value := range extra {
		data[key] = value
	}
	return engine.HandleAction(playerID, ActionMessage{Action: "resolve_action", Data: data})
}

func lastSpellCastPower(engine *Engine, instanceID string) (int, bool) {
	for i := len(engine.log) - 1; i >= 0; i-- {
		event := engine.log[i]
		if event.Type != "spell_cast" {
			continue
		}
		if skill, ok := event.Data["skill"].(map[string]any); ok {
			id, _ := skill["instance_id"].(string)
			if id != instanceID {
				continue
			}
		}
		switch v := event.Data["power"].(type) {
		case int:
			return v, true
		case float64:
			return int(v), true
		}
	}
	return 0, false
}

func TestRoyalConflictHellRoarCanSacrificeFireCompanionForPower(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	roar := readySkill(baseCard(t, "3121104"), 0)
	p0.Skills[0] = roar
	setAllElements(p0, 9)
	sacrifice := placeUnit(baseCard(t, "1121101"), 0, 0, 0, engine)
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	bonus := totalElementCost(sacrifice.Card.ElementsCost)

	if err := engine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id":  roar.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(target.Position.Col),
			"target_row":   float64(target.Position.Row),
			"sacrifice_id": sacrifice.InstanceID,
		},
	}); err != nil {
		t.Fatalf("cast 3121104 with sacrifice: %v", err)
	}
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != roar.Card.Power+bonus {
		t.Fatalf("3121104 should add sacrificed companion entry cost to pending power, pending=%+v bonus=%d", engine.State.PendingSpell, bonus)
	}
	if p0.Units[0][0] != nil || len(p0.Graveyard) == 0 || p0.Graveyard[len(p0.Graveyard)-1] != sacrifice {
		t.Fatalf("3121104 should sacrifice selected fire companion, units=%v grave=%v", p0.Units[0][0], cardsToInfo(p0.Graveyard))
	}

	boostEngine := setupEffectTest(t)
	boostP0 := boostEngine.State.Players[0]
	mainAttack := readySkill(baseCard(t, "3121106"), 0)
	boostRoar := readySkill(baseCard(t, "3121104"), 0)
	boostP0.Skills[0] = mainAttack
	boostP0.Skills[1] = boostRoar
	setAllElements(boostP0, 9)
	boostSacrifice := placeUnit(baseCard(t, "1121101"), 0, 0, 0, boostEngine)
	boostTarget := placeUnit(baseCard(t, "1021001"), 1, 0, 0, boostEngine)
	boostBonus := totalElementCost(boostSacrifice.Card.ElementsCost)
	if err := boostEngine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id":  mainAttack.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(boostTarget.Position.Col),
			"target_row":   float64(boostTarget.Position.Row),
			"boost_ids":    []any{boostRoar.InstanceID},
			"sacrifice_id": boostSacrifice.InstanceID,
		},
	}); err != nil {
		t.Fatalf("attack boost 3121104 with sacrifice: %v", err)
	}
	wantBoostPower := mainAttack.Card.Power + boostRoar.Card.Power + boostBonus
	if boostEngine.State.PendingSpell == nil || boostEngine.State.PendingSpell.TotalPower != wantBoostPower {
		t.Fatalf("3121104 attack boost should add sacrificed companion entry cost to pending power, pending=%+v want=%d", boostEngine.State.PendingSpell, wantBoostPower)
	}
	if boostP0.Units[0][0] != nil || !containsCardInstance(boostP0.Graveyard, boostSacrifice) {
		t.Fatalf("3121104 attack boost should sacrifice selected fire companion, unit=%v grave=%v", boostP0.Units[0][0], cardsToInfo(boostP0.Graveyard))
	}

	invalidEngine := setupEffectTest(t)
	invalidP0 := invalidEngine.State.Players[0]
	invalidRoar := readySkill(baseCard(t, "3121104"), 0)
	invalidP0.Skills[0] = invalidRoar
	setAllElements(invalidP0, 9)
	nonFire := placeUnit(baseCard(t, "1021001"), 0, 0, 0, invalidEngine)
	invalidTarget := placeUnit(baseCard(t, "1021002"), 1, 0, 0, invalidEngine)
	if err := invalidEngine.HandleAction(0, ActionMessage{
		Action: "cast_spell",
		Data: map[string]any{
			"instance_id":  invalidRoar.InstanceID,
			"target_type":  "unit",
			"target_col":   float64(invalidTarget.Position.Col),
			"target_row":   float64(invalidTarget.Position.Row),
			"sacrifice_id": nonFire.InstanceID,
		},
	}); err == nil {
		t.Fatalf("3121104 should reject non-fire sacrifice")
	}
	if invalidP0.Units[0][0] != nonFire || len(invalidP0.Graveyard) != 0 {
		t.Fatalf("3121104 invalid sacrifice should not spend board state, unit=%v grave=%v", cardToInfo(invalidP0.Units[0][0]), cardsToInfo(invalidP0.Graveyard))
	}

	defenseEngine := setupEffectTest(t)
	defender := defenseEngine.State.Players[1]
	defenseEngine.State.Phase = PhaseDefenseWindow
	defenseEngine.State.CurrentTurn = 0
	defenseRoar := readySkill(baseCard(t, "3121104"), 1)
	defender.Skills[0] = defenseRoar
	setAllElements(defender, 9)
	defenseSacrifice := placeUnit(baseCard(t, "1121101"), 1, 0, 0, defenseEngine)
	defenseBonus := totalElementCost(defenseSacrifice.Card.ElementsCost)
	attackSpell := readySkill(baseCard(t, "3021005"), 0)
	defenseEngine.State.PendingSpell = &SpellCast{
		AttackerID: 0,
		Skill:      attackSpell,
		Target:     SpellTarget{Type: "hero"},
		TotalPower: defenseRoar.Card.Power + defenseBonus,
	}
	if err := defenseEngine.HandleAction(1, ActionMessage{
		Action: "defend",
		Data: map[string]any{
			"skill_ids":    []any{defenseRoar.InstanceID},
			"sacrifice_id": defenseSacrifice.InstanceID,
		},
	}); err != nil {
		t.Fatalf("defend 3121104 with sacrifice: %v", err)
	}
	if defender.Units[0][0] != nil || !containsCardInstance(defender.Graveyard, defenseSacrifice) {
		t.Fatalf("3121104 defense should sacrifice selected fire companion, unit=%v grave=%v", defender.Units[0][0], cardsToInfo(defender.Graveyard))
	}
	if power, ok := lastDefenseAttemptPower(defenseEngine); !ok || power != defenseRoar.Card.Power+defenseBonus {
		t.Fatalf("3121104 defense should add sacrifice cost to defense power, power=%d ok=%v want=%d", power, ok, defenseRoar.Card.Power+defenseBonus)
	}

	defenseBoostEngine := setupEffectTest(t)
	defenseBoostP1 := defenseBoostEngine.State.Players[1]
	defenseBoostEngine.State.Phase = PhaseDefenseWindow
	defenseBoostEngine.State.CurrentTurn = 0
	mainDefense := readySkill(baseCard(t, "3221004"), 1)
	defenseBoostRoar := readySkill(baseCard(t, "3121104"), 1)
	defenseBoostP1.Skills[0] = mainDefense
	defenseBoostP1.Skills[1] = defenseBoostRoar
	setAllElements(defenseBoostP1, 9)
	defenseBoostSacrifice := placeUnit(baseCard(t, "1121101"), 1, 0, 0, defenseBoostEngine)
	defenseBoostBonus := totalElementCost(defenseBoostSacrifice.Card.ElementsCost)
	defenseBoostEngine.State.PendingSpell = &SpellCast{
		AttackerID: 0,
		Skill:      readySkill(baseCard(t, "3021005"), 0),
		Target:     SpellTarget{Type: "hero"},
		TotalPower: 999,
	}
	if err := defenseBoostEngine.HandleAction(1, ActionMessage{
		Action: "defend",
		Data: map[string]any{
			"skill_ids":    []any{mainDefense.InstanceID},
			"boost_ids":    []any{defenseBoostRoar.InstanceID},
			"sacrifice_id": defenseBoostSacrifice.InstanceID,
		},
	}); err != nil {
		t.Fatalf("defense boost 3121104 with sacrifice: %v", err)
	}
	wantDefenseBoostPower := mainDefense.Card.Power + defenseBoostRoar.Card.Power + defenseBoostBonus
	if defenseBoostP1.Units[0][0] != nil || !containsCardInstance(defenseBoostP1.Graveyard, defenseBoostSacrifice) {
		t.Fatalf("3121104 defense boost should sacrifice selected fire companion, unit=%v grave=%v", defenseBoostP1.Units[0][0], cardsToInfo(defenseBoostP1.Graveyard))
	}
	if power, ok := lastDefenseAttemptPower(defenseBoostEngine); !ok || power != wantDefenseBoostPower {
		t.Fatalf("3121104 defense boost should add sacrifice cost to defense power, power=%d ok=%v want=%d", power, ok, wantDefenseBoostPower)
	}
}

func lastDefenseAttemptPower(engine *Engine) (int, bool) {
	for i := len(engine.log) - 1; i >= 0; i-- {
		event := engine.log[i]
		if event.Type != "defense_attempt" {
			continue
		}
		switch v := event.Data["defense_power"].(type) {
		case int:
			return v, true
		case float64:
			return int(v), true
		}
	}
	return 0, false
}

func TestRoyalConflictFireCloudFanTargetsBackEnemyWithEmptyFront(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	fan := NewCardInstance(baseCard(t, "2121102"), 0, engine.State.TurnNumber)
	p0.Equipment[0] = fan
	fireSpell := readySkill(baseCard(t, "3121105"), 0)
	airSpell := readySkill(baseCard(t, "3321002"), 0)
	waterSpell := readySkill(baseCard(t, "3221001"), 0)
	placeUnit(baseCard(t, "1021002"), 1, 0, 0, engine)
	backEnemy := placeUnit(baseCard(t, "1021001"), 1, 1, 1, engine)
	target := SpellTarget{Type: "unit", Position: *backEnemy.Position}

	if err := engine.validateSpellTarget(0, fireSpell, target); err != nil {
		t.Fatalf("2121102 should let fire spells target a back enemy with empty front, err=%v", err)
	}
	if err := engine.validateSpellTarget(0, airSpell, target); err != nil {
		t.Fatalf("2121102 should let air spells target a back enemy with empty front, err=%v", err)
	}
	if err := engine.validateSpellTarget(0, waterSpell, target); err == nil {
		t.Fatalf("2121102 should not help non-fire/non-air spells")
	}

	blockedEngine := setupEffectTest(t)
	blockedP0 := blockedEngine.State.Players[0]
	blockedP0.Equipment[0] = NewCardInstance(baseCard(t, "2121102"), 0, blockedEngine.State.TurnNumber)
	blockedBack := placeUnit(baseCard(t, "1021001"), 1, 1, 1, blockedEngine)
	placeUnit(baseCard(t, "1021002"), 1, 1, 0, blockedEngine)
	if err := blockedEngine.validateSpellTarget(0, readySkill(baseCard(t, "3121105"), 0), SpellTarget{Type: "unit", Position: *blockedBack.Position}); err == nil {
		t.Fatalf("2121102 should not help when the target has a unit directly in front")
	}
}

func TestRoyalConflictDragonSnowfieldFreezesAfterSpellHitAndSummonsFrostDragon(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	snowfield := readySkill(baseCard(t, "3211102"), 0)
	p0.Skills[0] = snowfield
	frontTarget := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
	backTarget := placeUnit(baseCard(t, "1021002"), 0, 1, 1, engine)
	enemySpell := readySkill(baseCard(t, "3121001"), 1)

	engine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, enemySpell, map[string]any{"attacker": 1, "spell_source": enemySpell})
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "dragon_snowfield_freeze" ||
		!candidateContains(engine.State.PendingAction.Candidates, frontTarget.InstanceID) ||
		candidateContains(engine.State.PendingAction.Candidates, backTarget.InstanceID) {
		t.Fatalf("3211102 should prompt the caster to freeze an enemy in spell range, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 1, frontTarget.InstanceID)
	if frontTarget.Statuses[StatusFreeze] != 1 || snowfield.Statuses[dragonSnowfieldTriggerStatus] != 1 || engine.State.PendingAction != nil {
		t.Fatalf("3211102 should freeze and count its trigger, freeze=%d count=%d pending=%+v", frontTarget.Statuses[StatusFreeze], snowfield.Statuses[dragonSnowfieldTriggerStatus], engine.State.PendingAction)
	}

	snowfield.Statuses[dragonSnowfieldTriggerStatus] = 4
	engine.triggerFieldEffectsWithData(TriggerOnSpellHit, 0, enemySpell, map[string]any{"attacker": 1, "spell_source": enemySpell})
	resolvePendingSelection(t, engine, 1, frontTarget.InstanceID)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "dragon_snowfield_summon_frost_dragon" {
		t.Fatalf("3211102 should offer frost dragon summon after fifth trigger, pending=%+v", engine.State.PendingAction)
	}
	pos := Position{Col: 0, Row: 0}
	if !candidateContains(engine.State.PendingAction.Candidates, positionSelectionID(pos)) {
		t.Fatalf("3211102 summon prompt should include empty own position, candidates=%+v", engine.State.PendingAction.Candidates)
	}
	resolvePendingSelection(t, engine, 0, positionSelectionID(pos))
	summoned := p0.Units[pos.Col][pos.Row]
	if summoned == nil || summoned.Card == nil || summoned.Card.Number != "1201101" || snowfield.Statuses[dragonSnowfieldTriggerStatus] != 0 {
		t.Fatalf("3211102 should summon frost dragon and spend five counts, summoned=%v count=%d", cardToInfo(summoned), snowfield.Statuses[dragonSnowfieldTriggerStatus])
	}
}

func TestRoyalConflictIceSoulSealCancelCountersLowPowerBoosts(t *testing.T) {
	t.Run("attack boost", func(t *testing.T) {
		engine := setupEffectTest(t)
		p1 := engine.State.Players[1]
		main := readySkill(baseCard(t, "3121001"), 0)
		boost := readySkill(baseCard(t, "3221001"), 0)
		counter := NewCardInstance(baseCard(t, "2221112"), 1, engine.State.TurnNumber)
		counter.IsSetCounter = true
		p1.Equipment[0] = counter
		p1.Elements[model.ElementWater] = 1
		engine.State.PendingSpell = &SpellCast{
			AttackerID:   0,
			Skill:        main,
			Target:       SpellTarget{Type: "unit", Position: Position{Col: 0, Row: 0}},
			TotalPower:   main.Card.Power + boost.Card.Power,
			BoostSkills:  []*CardInstance{boost},
			PowerSources: []SpellPowerSource{spellPowerSourceForCard(main, main.Card.Power, true), spellPowerSourceForCard(boost, boost.Card.Power, false)},
		}
		done := false
		if !engine.promptAttackBoostSpellCastCounters(0, []*CardInstance{boost}, func() { done = true }) {
			t.Fatalf("2221112 should open a counter window for low-power attack boost")
		}
		resolvePendingSelection(t, engine, 1, counter.InstanceID)
		if !done || len(engine.State.PendingSpell.BoostSkills) != 0 || engine.State.PendingSpell.TotalPower != main.Card.Power {
			t.Fatalf("2221112 should cancel attack boost and refresh power, done=%v boosts=%d power=%d", done, len(engine.State.PendingSpell.BoostSkills), engine.State.PendingSpell.TotalPower)
		}
	})

	t.Run("defense boost", func(t *testing.T) {
		engine := setupEffectTest(t)
		p1 := engine.State.Players[1]
		defense := readySkill(baseCard(t, "3221001"), 0)
		boost := readySkill(baseCard(t, "3221001"), 0)
		counter := NewCardInstance(baseCard(t, "2221112"), 1, engine.State.TurnNumber)
		counter.IsSetCounter = true
		p1.Equipment[0] = counter
		p1.Elements[model.ElementWater] = 1
		engine.State.PendingSpell = &SpellCast{
			AttackerID: 1,
			Skill:      readySkill(baseCard(t, "3121001"), 1),
			Target:     SpellTarget{Type: "unit", Position: Position{Col: 0, Row: 0}},
			TotalPower: 99,
		}
		done := false
		if !engine.promptDefenseSpellCastCounters(1, 0, []*CardInstance{defense}, []*CardInstance{boost}, nil, func() { done = true }) {
			t.Fatalf("2221112 should open a counter window for low-power defense boost")
		}
		resolvePendingSelection(t, engine, 1, counter.InstanceID)
		filtered := engine.filterIceSoulSealCancelledBoosts([]*CardInstance{defense, boost})
		if !done || len(filtered) != 1 || filtered[0] != defense {
			t.Fatalf("2221112 should cancel defense boost before defense power, done=%v filtered=%v", done, cardsToInfo(filtered))
		}
	})
}

func TestRoyalConflictEmbersCastsForFreeAtTurnEndWithNoFire(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	embers := readySkill(baseCard(t, "3121105"), 0)
	p0.Skills[0] = embers
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	p0.Elements[model.ElementFire] = 0

	if err := (Card3121105Embers{}).OnTurnEnd(&EffectContext{Engine: engine, Source: embers, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("3121105 turn end: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "embers_free_cast_target" ||
		!candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) {
		t.Fatalf("3121105 should prompt for a free target at turn end, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, target.InstanceID)
	if !embers.IsHorizontal || engine.State.PendingSpell == nil || engine.State.PendingSpell.Skill != embers ||
		engine.State.PendingSpell.TotalPower != embers.Card.Power {
		t.Fatalf("3121105 should cast itself free and enter defense window, horizontal=%v pending=%+v", embers.IsHorizontal, engine.State.PendingSpell)
	}
	if p0.Elements[model.ElementFire] != 0 {
		t.Fatalf("3121105 free cast should not spend fire, elements=%v", p0.Elements)
	}

	noTriggerEngine := setupEffectTest(t)
	noTriggerP0 := noTriggerEngine.State.Players[0]
	noTrigger := readySkill(baseCard(t, "3121105"), 0)
	noTriggerP0.Skills[0] = noTrigger
	noTriggerP0.Elements[model.ElementFire] = 1
	placeUnit(baseCard(t, "1021001"), 1, 0, 0, noTriggerEngine)
	if err := (Card3121105Embers{}).OnTurnEnd(&EffectContext{Engine: noTriggerEngine, Source: noTrigger, PlayerID: 0, OpponentID: 1}); err != nil {
		t.Fatalf("3121105 nonzero fire turn end: %v", err)
	}
	if noTriggerEngine.State.PendingAction != nil || noTrigger.IsHorizontal {
		t.Fatalf("3121105 should not trigger with remaining fire, pending=%+v horizontal=%v", noTriggerEngine.State.PendingAction, noTrigger.IsHorizontal)
	}
}

func TestRoyalConflictSkyPhantasmCastsAnotherLearnedDriveOrFocusSpell(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	sky := readySkill(baseCard(t, "3311101"), 0)
	drive := readySkill(baseCard(t, "3321002"), 0)
	p0.Skills[0] = sky
	p0.Skills[1] = drive
	p0.Elements[model.ElementLight] = 2
	p0.Elements[model.ElementAir] = 3
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": sky.InstanceID,
		"target_type": "none",
	}}); err != nil {
		t.Fatalf("3311101 cast: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "sky_phantasm_spell_choice" ||
		!candidateContains(engine.State.PendingAction.Candidates, drive.InstanceID) {
		t.Fatalf("3311101 should ask for another learned drive/focus spell, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, drive.InstanceID)
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "sky_phantasm_target" ||
		!candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) {
		t.Fatalf("3311101 should ask for copied spell target, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, target.InstanceID)
	if !sky.IsHorizontal || drive.IsHorizontal || engine.State.PendingSpell == nil ||
		engine.State.PendingSpell.Skill == drive ||
		engine.State.PendingSpell.Skill.Card.Number != drive.Card.Number ||
		engine.State.PendingSpell.TotalPower != drive.Card.Power {
		t.Fatalf("3311101 should cast a virtual copy without tapping the chosen spell, sky=%v drive=%v pending=%+v", sky.IsHorizontal, drive.IsHorizontal, engine.State.PendingSpell)
	}

	blockedEngine := setupEffectTest(t)
	blockedP0 := blockedEngine.State.Players[0]
	blockedSky := readySkill(baseCard(t, "3311101"), 0)
	blockedBoost := readySkill(baseCard(t, "3321002"), 0)
	blockedP0.Skills[0] = blockedSky
	blockedP0.Skills[1] = blockedBoost
	blockedP0.Elements[model.ElementLight] = 2
	blockedP0.Elements[model.ElementAir] = 3
	if err := blockedEngine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": blockedSky.InstanceID,
		"target_type": "none",
		"boost_ids":   []any{blockedBoost.InstanceID},
	}}); err == nil {
		t.Fatalf("3311101 should reject direct boost ids instead of charging and discarding them")
	}
	if blockedSky.IsHorizontal || blockedBoost.IsHorizontal || blockedP0.Elements[model.ElementLight] != 2 || blockedP0.Elements[model.ElementAir] != 3 {
		t.Fatalf("3311101 rejected boost should not spend or tap, sky=%v boost=%v elements=%v", blockedSky.IsHorizontal, blockedBoost.IsHorizontal, blockedP0.Elements)
	}
}

func TestRoyalConflictWaterMirrorScrollCopiesPreviousLowCostWaterSpell(t *testing.T) {
	engine := setupEffectTest(t)
	p0 := engine.State.Players[0]
	waterSpell := readySkill(baseCard(t, "3221001"), 0)
	waterSpell.PowerBonus = 2
	engine.recordSpellCast(0, waterSpell)
	scroll := NewCardInstance(baseCard(t, "2221104"), 0, engine.State.TurnNumber)
	p0.Hand = append(p0.Hand, scroll)
	p0.Elements[model.ElementWater] = 2
	target := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": scroll.InstanceID,
	}}); err != nil {
		t.Fatalf("2221104 use: %v", err)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "water_mirror_scroll_target" ||
		!candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) {
		t.Fatalf("2221104 should ask for copied spell target, pending=%+v", engine.State.PendingAction)
	}
	resolvePendingSelection(t, engine, 0, target.InstanceID)
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.Skill.Card.Number != waterSpell.Card.Number ||
		engine.State.PendingSpell.Skill == waterSpell ||
		engine.State.PendingSpell.TotalPower != waterSpell.Card.Power+2 {
		t.Fatalf("2221104 should cast a virtual copy of previous low-cost water spell, pending=%+v", engine.State.PendingSpell)
	}

	failEngine := setupEffectTest(t)
	failP0 := failEngine.State.Players[0]
	failScroll := NewCardInstance(baseCard(t, "2221104"), 0, failEngine.State.TurnNumber)
	failP0.Hand = append(failP0.Hand, failScroll)
	failP0.Elements[model.ElementWater] = 2
	if err := failEngine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": failScroll.InstanceID}}); err == nil {
		t.Fatalf("2221104 should reject use without previous low-cost water spell")
	}
}

func TestRoyalConflictFlameInfernoBurnsAndRaisesDefenseThreshold(t *testing.T) {
	engine := setupEffectTest(t)
	targetA := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	targetB := placeUnit(baseCard(t, "1021002"), 1, 1, 0, engine)
	inferno := readySkill(baseCard(t, "3111101"), 0)
	if err := (Card3111101FlameInferno{}).OnSpellHit(&EffectContext{
		Engine:    engine,
		Source:    inferno,
		PlayerID:  0,
		ExtraData: map[string]any{"power": 17, "affected_units": []*CardInstance{targetA, targetB}},
	}); err != nil {
		t.Fatalf("3111101 hit: %v", err)
	}
	if targetA.Statuses[StatusBurn] != 2 || targetB.Statuses[StatusBurn] != 2 {
		t.Fatalf("3111101 should add one burn per 8 power, a=%v b=%v", targetA.Statuses, targetB.Statuses)
	}
	if got := engine.requiredDefensePowerForSpell(inferno, 16); got != 24 {
		t.Fatalf("3111101 should require +4 defense per burn layer, got %d", got)
	}

	failEngine := setupEffectTest(t)
	failTarget := placeUnit(baseCard(t, "1021001"), 1, 0, 0, failEngine)
	failTarget.CurrentLife = 10
	failInferno := readySkill(baseCard(t, "3111101"), 0)
	failDefense := readySkill(baseCard(t, "3221001"), 1)
	failDefense.PowerBonus = 15
	failEngine.State.PendingSpell = &SpellCast{
		AttackerID: 0,
		Skill:      failInferno,
		Target:     SpellTarget{Type: "unit", Position: *failTarget.Position},
		TotalPower: 16,
	}
	startLife := failTarget.CurrentLife
	failEngine.finishDefenseResolution(1, []*CardInstance{failDefense}, nil, 0, nil, 0)
	if failTarget.CurrentLife >= startLife || failTarget.Statuses[StatusBurn] != 2 {
		t.Fatalf("3111101 should still hit when defense meets base power but misses extra burn threshold, life=%d start=%d statuses=%v", failTarget.CurrentLife, startLife, failTarget.Statuses)
	}

	successEngine := setupEffectTest(t)
	successTarget := placeUnit(baseCard(t, "1021001"), 1, 0, 0, successEngine)
	successTarget.CurrentLife = 10
	successInferno := readySkill(baseCard(t, "3111101"), 0)
	successDefense := readySkill(baseCard(t, "3221001"), 1)
	successDefense.PowerBonus = 23
	successEngine.State.PendingSpell = &SpellCast{
		AttackerID: 0,
		Skill:      successInferno,
		Target:     SpellTarget{Type: "unit", Position: *successTarget.Position},
		TotalPower: 16,
	}
	successStartLife := successTarget.CurrentLife
	successEngine.finishDefenseResolution(1, []*CardInstance{successDefense}, nil, 0, nil, 0)
	if successTarget.CurrentLife != successStartLife || successTarget.Statuses[StatusBurn] != 0 {
		t.Fatalf("3111101 should be fully defended when extra burn threshold is met, life=%d start=%d statuses=%v", successTarget.CurrentLife, successStartLife, successTarget.Statuses)
	}
}

func TestRoyalConflictFiveRainbowRingMarkersAndBeamModes(t *testing.T) {
	t.Run("ring marker", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		ring := NewCardInstance(baseCard(t, "2511102"), 0, engine.State.TurnNumber)
		p0.Equipment[0] = ring
		p0.Elements[model.ElementFire] = 1
		if err := (Card2511102FiveRainbowRing{}).OnPerTurn(&EffectContext{Engine: engine, Source: ring, PlayerID: 0, OpponentID: 1}); err != nil {
			t.Fatalf("2511102 per-turn: %v", err)
		}
		if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "five_rainbow_ring_marker" ||
			!candidateContains(engine.State.PendingAction.Candidates, model.ElementFire) {
			t.Fatalf("2511102 should prompt payable marker elements, pending=%+v", engine.State.PendingAction)
		}
		resolvePendingSelection(t, engine, 0, model.ElementFire)
		if p0.Elements[model.ElementFire] != 0 || ring.Statuses[fiveRainbowMarkerStatus(model.ElementFire)] != 1 {
			t.Fatalf("2511102 should pay fire and add fire marker, elements=%v statuses=%v", p0.Elements, ring.Statuses)
		}
	})

	t.Run("beam consumes markers", func(t *testing.T) {
		engine := setupEffectTest(t)
		p0 := engine.State.Players[0]
		ring := NewCardInstance(baseCard(t, "2511102"), 0, engine.State.TurnNumber)
		beam := readySkill(baseCard(t, "3501101"), 0)
		ring.BoundSkills = []*CardInstance{beam}
		p0.Equipment[0] = ring
		for _, elem := range fiveRainbowElements() {
			ring.Statuses[fiveRainbowMarkerStatus(elem)] = 1
		}
		front := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		back := placeUnit(baseCard(t, "1021002"), 1, 1, 1, engine)
		extra := placeUnit(baseCard(t, "1021003"), 1, 0, 0, engine)
		back.CurrentLife = 10
		extra.CurrentLife = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": beam.InstanceID,
			"target_type": "unit",
			"target_col":  float64(back.Position.Col),
			"target_row":  float64(back.Position.Row),
			"rainbow_markers": map[string]any{
				model.ElementFire:  1.0,
				model.ElementWater: 1.0,
				model.ElementEarth: 1.0,
				model.ElementAir:   1.0,
				model.ElementLight: 1.0,
			},
			"extra_targets": []any{map[string]any{"col": float64(extra.Position.Col), "row": float64(extra.Position.Row)}},
		}}); err != nil {
			t.Fatalf("3501101 cast: %v", err)
		}
		if front == nil || engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 14 || len(engine.State.PendingSpell.ExtraTargets) != 1 {
			t.Fatalf("3501101 should pierce, double power, and carry extra target, pending=%+v front=%v", engine.State.PendingSpell, front)
		}
		for _, elem := range fiveRainbowElements() {
			if ring.Statuses[fiveRainbowMarkerStatus(elem)] != 0 {
				t.Fatalf("3501101 should remove consumed marker %s, statuses=%v", elem, ring.Statuses)
			}
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("3501101 no defend: %v", err)
		}
		if back.CurrentLife != 7 || extra.CurrentLife != 7 {
			t.Fatalf("3501101 fire marker should make attack 3 to both targets, back=%d extra=%d", back.CurrentLife, extra.CurrentLife)
		}
		if beam.Statuses[fiveRainbowBeamSelectedStatus(model.ElementAir)] != 0 || beam.Statuses[fiveRainbowBeamAllStatus] != 0 {
			t.Fatalf("3501101 should clear temporary marker selection after resolution, statuses=%v", beam.Statuses)
		}
	})
}

func TestRoyalConflictCounterWindHoleScrollReflectsMissedOrCancelledSingleSpell(t *testing.T) {
	engine := setupEffectTest(t)
	p1 := engine.State.Players[1]
	counter := NewCardInstance(baseCard(t, "2321111"), 1, engine.State.TurnNumber)
	counter.IsSetCounter = true
	p1.Equipment[0] = counter
	p1.Elements[model.ElementAir] = 2
	original := readySkill(baseCard(t, "3121001"), 0)
	boost := readySkill(baseCard(t, "3321002"), 0)
	target := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	target.CurrentLife = 10

	prompted := engine.promptSpellMissOrCancelledCounters(0, original, []*CardInstance{boost}, nil, "test_cancelled")
	if !prompted || engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_trigger" {
		t.Fatalf("2321111 should prompt after a non-area enemy spell is cancelled, prompted=%v pending=%+v", prompted, engine.State.PendingAction)
	}
	if err := resolvePendingSelectionWithData(engine, 1, []string{counter.InstanceID}, nil); err != nil {
		t.Fatalf("resolve 2321111 reveal: %v", err)
	}
	if p1.Equipment[0] != nil || countCardNumber(p1.Graveyard, "2321111") != 1 || p1.Elements[model.ElementAir] != 0 {
		t.Fatalf("2321111 should pay and discard after reveal, equipment=%v grave=%v elements=%v", cardToInfo(p1.Equipment[0]), cardsToInfo(p1.Graveyard), p1.Elements)
	}
	if engine.State.PendingAction == nil || engine.State.PendingAction.Type != "counter_wind_hole_scroll_target" ||
		!candidateContains(engine.State.PendingAction.Candidates, target.InstanceID) {
		t.Fatalf("2321111 should ask for an enemy companion target, pending=%+v target=%s", engine.State.PendingAction, target.InstanceID)
	}
	if err := resolvePendingSelectionWithData(engine, 1, []string{target.InstanceID}, nil); err != nil {
		t.Fatalf("resolve 2321111 target: %v", err)
	}
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.AttackerID != 1 ||
		engine.State.PendingSpell.Skill.Card.Number != original.Card.Number || len(engine.State.PendingSpell.BoostSkills) != 1 ||
		engine.State.PendingSpell.BoostSkills[0].Card.Number != boost.Card.Number {
		t.Fatalf("2321111 should create a reflected copy with boosts, pending=%+v", engine.State.PendingSpell)
	}
	if engine.State.PendingSpell.TotalPower != 6 || engine.State.Phase != PhaseDefenseWindow {
		t.Fatalf("2321111 reflected spell should carry boosted power and open defense, phase=%s pending=%+v", engine.State.Phase, engine.State.PendingSpell)
	}

	lostEngine := setupEffectTest(t)
	lostP1 := lostEngine.State.Players[1]
	lostCounter := NewCardInstance(baseCard(t, "2321111"), 1, lostEngine.State.TurnNumber)
	lostCounter.IsSetCounter = true
	lostP1.Equipment[0] = lostCounter
	lostP1.Elements[model.ElementAir] = 2
	lostSkill := readySkill(baseCard(t, "3121001"), 0)
	lostTarget := placeUnit(baseCard(t, "1021001"), 1, 0, 0, lostEngine)
	reflectedTarget := placeUnit(baseCard(t, "1021002"), 0, 0, 0, lostEngine)
	lostEngine.State.PendingSpell = &SpellCast{
		AttackerID: 0,
		Skill:      lostSkill,
		Target:     SpellTarget{Type: "unit", Position: *lostTarget.Position},
		TotalPower: lostSkill.Card.Power,
	}
	lostEngine.State.Players[1].Units[lostTarget.Position.Col][lostTarget.Position.Row] = nil
	lostEngine.resolvePendingSpellHit()
	if lostEngine.State.PendingAction == nil || lostEngine.State.PendingAction.Type != "counter_trigger" {
		t.Fatalf("2321111 should trigger from real target-lost spell miss, pending=%+v", lostEngine.State.PendingAction)
	}
	if err := resolvePendingSelectionWithData(lostEngine, 1, []string{lostCounter.InstanceID}, nil); err != nil {
		t.Fatalf("resolve target-lost 2321111 reveal: %v", err)
	}
	if lostEngine.State.PendingAction == nil || lostEngine.State.PendingAction.Type != "counter_wind_hole_scroll_target" {
		t.Fatalf("2321111 target-lost path should ask for reflected target, pending=%+v", lostEngine.State.PendingAction)
	}
	if err := resolvePendingSelectionWithData(lostEngine, 1, []string{reflectedTarget.InstanceID}, nil); err != nil {
		t.Fatalf("resolve target-lost 2321111 target: %v", err)
	}
	if lostEngine.State.PendingSpell == nil || lostEngine.State.PendingSpell.AttackerID != 1 ||
		lostEngine.State.PendingSpell.Skill.Card.Number != lostSkill.Card.Number {
		t.Fatalf("2321111 target-lost continuation should preserve reflected pending spell, pending=%+v", lostEngine.State.PendingSpell)
	}

	blockedEngine := setupEffectTest(t)
	blockedP1 := blockedEngine.State.Players[1]
	blockedCounter := NewCardInstance(baseCard(t, "2321111"), 1, blockedEngine.State.TurnNumber)
	blockedCounter.IsSetCounter = true
	blockedP1.Equipment[0] = blockedCounter
	blockedP1.Elements[model.ElementAir] = 2
	placeUnit(baseCard(t, "1021001"), 0, 0, 0, blockedEngine)
	areaSpell := readySkill(baseCard(t, "3221001"), 0)
	if prompted := blockedEngine.promptSpellMissOrCancelledCounters(0, areaSpell, nil, nil, "test_cancelled"); prompted {
		t.Fatalf("2321111 should not trigger for area spells, pending=%+v", blockedEngine.State.PendingAction)
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
