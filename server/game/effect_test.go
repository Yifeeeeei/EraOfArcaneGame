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
		"2121112": SpellAreaColumn,
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
	for _, number := range []string{"2121109", "2121112", "2221110", "2521112"} {
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
	if got := engine.effectiveSpellPower(0, willErosion, nil); got != willErosion.Card.Power+5 {
		t.Fatalf("two red moon markers should add +2 to other shadow spell during red moon, got %d", got)
	}
	if got := engine.effectiveSpellPower(0, redMoon, nil); got != 2 {
		t.Fatalf("red moon markers should not buff red moon itself, got %d", got)
	}

	seviana.CurrentLife = 1
	seviana.IsHorizontal = true
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
