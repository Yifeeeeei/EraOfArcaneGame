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

	for number, name := range map[string]string{
		"1121102": "火山谷底巨兽",
		"1121105": "弗卡莱诺近卫",
		"1201101": "凛冰之龙",
		"1221101": "掠夺者海盗船",
		"1221104": "冰原猛犸",
		"1221113": "凛冬城象骑兵",
		"1401101": "普通蜥蜴",
		"1421101": "岩壁刺球",
		"1421103": "寄生虫",
		"1421109": "地穴巨蝠",
		"1521104": "旭日之龙",
		"1621105": "混沌胚胎",
		"1621107": "蔷薇死神",
		"2001102": "九霄印记",
		"2021109": "氏族战锤",
		"2121103": "浴火之翼",
		"2121112": "炎流卷轴",
		"2421101": "秋暮耳环",
		"1021112": "奥术纯净体",
		"3021102": "奥术屏障",
		"3121106": "爆炎气焰",
		"3121108": "熔岩障壁",
		"3221107": "海龙卷",
		"3321101": "急速涡旋",
		"3321103": "雷霆万钧",
		"3321106": "紫电穿空",
		"3421102": "苍岚之刃",
		"3421103": "巨岩崩落",
		"3421109": "石化死光",
		"3521103": "光铸飞弹",
		"3521105": "流光之束",
		"3521107": "虹彩之壁",
	} {
		behavior := globalRegistry.GetBehavior(number)
		if behavior == nil || behavior.ID() != number || behavior.Name() != name {
			t.Fatalf("%s should have explicit vanilla behavior, behavior=%#v", number, behavior)
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

	t.Run("rock wall guard gains shield only after enemy spell hits with no shield", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		guard := placeUnit(baseCard(t, "1021110"), 0, 0, 0, engine)
		behavior := Card1021110RockWallGuard{}

		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: guard, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 1}}); err != nil {
			t.Fatalf("1021110 enemy spell hit: %v", err)
		}
		if p0.Shield != 2 {
			t.Fatalf("1021110 should gain shield 2 after enemy spell hits while unshielded, shield=%d", p0.Shield)
		}

		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: guard, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 1}}); err != nil {
			t.Fatalf("1021110 second enemy spell hit: %v", err)
		}
		if p0.Shield != 2 {
			t.Fatalf("1021110 should not gain more shield while already shielded, shield=%d", p0.Shield)
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

	t.Run("raider gunner discards an enemy hand card after a friendly spell hits once per turn", func(t *testing.T) {
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
		if len(p1.Hand) != 1 || len(p1.Graveyard) != 1 || gunner.UsedThisTurn != 1 {
			t.Fatalf("1221111 should discard one enemy hand card and spend trigger, hand=%v grave=%v used=%d", cardsToInfo(p1.Hand), cardsToInfo(p1.Graveyard), gunner.UsedThisTurn)
		}
		if containsCardInstance(p1.Hand, p1.Graveyard[0]) {
			t.Fatalf("discarded card should leave enemy hand, hand=%v grave=%v", cardsToInfo(p1.Hand), cardsToInfo(p1.Graveyard))
		}

		if err := behavior.OnSpellHit(&EffectContext{Engine: engine, Source: gunner, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"attacker": 0}}); err != nil {
			t.Fatalf("raider gunner second friendly hit: %v", err)
		}
		if len(p1.Hand) != 1 || len(p1.Graveyard) != 1 || gunner.UsedThisTurn != 1 {
			t.Fatalf("1221111 should trigger at most once per turn, hand=%v grave=%v used=%d", cardsToInfo(p1.Hand), cardsToInfo(p1.Graveyard), gunner.UsedThisTurn)
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
		if len(p0.TempModifiers) != 0 {
			t.Fatalf("2321110 next-use power and attack should be consumed together when the skill is used, modifiers=%v", p0.TempModifiers)
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
		if got := effectiveElementsGain(seed)[model.ElementFire]; got != seed.Card.ElementsGain[model.ElementFire]+1 {
			t.Fatalf("1121111 should gain fire load when another companion takes fire damage, load=%v", effectiveElementsGain(seed))
		}

		engine.dealDamageWithExtra(enemy, 1, 1, map[string]any{"damage_source": "test", "damage_element": model.ElementFire})
		if got := effectiveElementsGain(seed)[model.ElementFire]; got != seed.Card.ElementsGain[model.ElementFire]+2 {
			t.Fatalf("1121111 should also see enemy companion fire damage, load=%v", effectiveElementsGain(seed))
		}

		engine.dealDamageWithExtra(ally, 1, 0, map[string]any{"damage_source": "test", "damage_element": model.ElementWater})
		if got := effectiveElementsGain(seed)[model.ElementFire]; got != seed.Card.ElementsGain[model.ElementFire]+2 {
			t.Fatalf("1121111 should ignore non-fire damage, load=%v", effectiveElementsGain(seed))
		}

		engine.dealDamageWithExtra(seed, 1, 0, map[string]any{"damage_source": "test", "damage_element": model.ElementFire})
		if got := effectiveElementsGain(seed)[model.ElementFire]; got != seed.Card.ElementsGain[model.ElementFire]+2 {
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
		if len(p0.TempModifiers) != 2 || p0.TempModifiers[0].TargetInstanceID != skill.InstanceID || p0.TempModifiers[1].TargetInstanceID != skill.InstanceID {
			t.Fatalf("3621101 should add two self next-use modifiers, modifiers=%+v", p0.TempModifiers)
		}
		if p0.TempModifiers[0].Type != TempModSkillPowerBonus || p0.TempModifiers[0].Amount != 2 || p0.TempModifiers[1].Type != TempModNextSkillUseAttackBonus || p0.TempModifiers[1].Amount != 1 {
			t.Fatalf("3621101 should grant +2 power and +1 attack next use, modifiers=%+v", p0.TempModifiers)
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
