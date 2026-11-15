package game

import (
	"eraofarcane/model"
	"testing"
)

func TestRoyalAuditLearnedSpellsDoNotGrantTheirCastOrHitEffects(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	drain := readySkill(baseCard(t, "3021103"), 0)
	flash := readySkill(baseCard(t, "3121109"), 0)
	other := readySkill(baseCard(t, "3121001"), 0)
	p.Skills[0], p.Skills[1], p.Skills[2] = drain, flash, other
	p.Hand = nil
	p.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
	target := placeUnit(baseCard(t, "1421101"), 1, 0, 0, e)
	target.CurrentLife = 20
	setAllElements(p, 20)
	if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": other.InstanceID, "target_type": "unit", "target_col": float64(0), "target_row": float64(0)}}); err != nil {
		t.Fatal(err)
	}
	if len(p.Hand) != 0 || len(p.Deck) != 1 {
		t.Fatal("learning Arcane Drain must not draw cards when another spell is cast")
	}
	before := p.Elements[model.ElementFire]
	if err := e.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatal(err)
	}
	if p.Elements[model.ElementFire] != before {
		t.Fatal("a learned Flame Flash must not reward another spell's hit")
	}
}

func TestRoyalAuditEachMagicMothKeepsItsIndependentChoice(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	p.Hand = nil
	a := NewCardInstance(baseCard(t, "1021113"), 0, 1)
	b := NewCardInstance(baseCard(t, "1021113"), 0, 1)
	p.Deck = []*CardInstance{a, b}
	focus := readySkill(baseCard(t, "3021101"), 0)
	e.triggerMagicMothAfterFocusSpellCast(0, focus)
	if e.State.PendingAction == nil {
		t.Fatal("missing first moth choice")
	}
	resolvePendingSelection(t, e, 0)
	if e.State.PendingAction == nil || !candidateContains(e.State.PendingAction.Candidates, b.InstanceID) {
		t.Fatal("declining the first moth must preserve the second trigger")
	}
	resolvePendingSelection(t, e, 0, b.InstanceID)
	if len(p.Hand) != 1 || p.Hand[0] != b || len(p.Deck) != 1 || p.Deck[0] != a {
		t.Fatal("only the accepted moth should leave the deck")
	}
}

func TestRoyalAuditBloodSandRejectsAllInvalidPaymentsBeforeSpending(t *testing.T) {
	e := setupEffectTest(t)
	a := placeUnit(baseCard(t, "1421101"), 0, 0, 0, e)
	b := placeUnit(baseCard(t, "1021001"), 0, 1, 0, e)
	a.CurrentLife = 4
	b.CurrentLife = 1
	source := readySkill(baseCard(t, "3411102"), 0)
	_, err := e.applyBloodSandPayments(0, []string{a.InstanceID, b.InstanceID}, map[string]any{"payments": []any{
		map[string]any{"instance_id": a.InstanceID, "amount": float64(1), "mode": "life"},
		map[string]any{"instance_id": b.InstanceID, "amount": float64(2), "mode": "life"},
	}}, source)
	if err == nil || a.CurrentLife != 4 || b.CurrentLife != 1 {
		t.Fatal("an invalid later payment must not consume earlier life")
	}
	if _, err := e.applyBloodSandPayments(0, nil, map[string]any{"payments": []any{map[string]any{"instance_id": a.InstanceID, "amount": float64(1), "mode": "life"}}}, source); err == nil || a.CurrentLife != 4 {
		t.Fatal("payment payload cannot spend from an unselected card")
	}
}

func TestRoyalAuditClawRequiresThreeDistinctSpellsIncludingBound(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[1]
	a, b, c := readySkill(baseCard(t, "3121001"), 1), readySkill(baseCard(t, "3221001"), 1), readySkill(baseCard(t, "3321005"), 1)
	p.Skills[0], p.Skills[1] = a, b
	p.Hero.BoundSkills = []*CardInstance{c}
	a.Statuses[StatusWeaken] = 4
	b.Statuses[StatusWeaken] = 1
	source := readySkill(baseCard(t, "3611102"), 0)
	ctx := e.skillContext(0, source)
	if err := (Card3611102ClawOfErebos{}).ValidateSkillLearn(ctx, source); err == nil {
		t.Fatal("two weakened spells are insufficient regardless of layers")
	}
	c.Statuses[StatusWeaken] = 1
	if err := (Card3611102ClawOfErebos{}).ValidateSkillLearn(ctx, source); err != nil {
		t.Fatalf("the bound third spell must count: %v", err)
	}
}

func TestRoyalAuditSoulChainsDoNotShareMarkedGroups(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	chainA := NewCardInstance(baseCard(t, "2611101"), 0, 1)
	chainB := NewCardInstance(baseCard(t, "2611101"), 0, 1)
	p.Equipment[0], p.Equipment[1] = chainA, chainB
	u := placeUnit(baseCard(t, "1021001"), 1, 0, 0, e)
	v := placeUnit(baseCard(t, "1021002"), 1, 1, 0, e)
	a, b := readySkill(baseCard(t, "3121001"), 1), readySkill(baseCard(t, "3221001"), 1)
	e.State.Players[1].Skills[0], e.State.Players[1].Skills[1] = a, b
	u.Statuses[soulChainUnitKey(chainA)] = 1
	a.Statuses[soulChainSpellKey(chainA)] = 1
	v.Statuses[soulChainUnitKey(chainB)] = 1
	b.Statuses[soulChainSpellKey(chainB)] = 1
	e.triggerErebosSoulChainMarkedOverexert(1, []*CardInstance{u})
	if a.Statuses[StatusWeaken] != 1 || b.Statuses[StatusWeaken] != 0 {
		t.Fatal("only the chain linked to the paid unit may weaken its own spell group")
	}
}

func TestRoyalAuditLopsiusRejectsOtherEffectsButKeepsOwnGrowth(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	source := readySkill(baseCard(t, "3111102"), 0)
	fuel := readySkill(baseCard(t, "3121001"), 0)
	p.Skills[0], p.Skills[1] = source, fuel
	if canSkillBeBoosted(source) || e.addStatus(source, StatusWeaken, 2) || e.addStatus(source, StatusPetrify, 1) {
		t.Fatal("printed boost and negative-effect immunity must apply")
	}
	e.addTemporaryModifier(0, TemporaryModifier{Type: TempModSkillPowerBonus, TargetInstanceID: source.InstanceID, Amount: 10})
	if got := e.effectiveSpellPower(0, source, nil); got != source.Card.Power {
		t.Fatal("another card's temporary bonus must not affect Lopsius")
	}
	if err := e.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{"instance_id": source.InstanceID, "ability_type": "per_turn"}}); err != nil {
		t.Fatal(err)
	}
	resolvePendingSelection(t, e, 0, fuel.InstanceID)
	if e.effectiveSpellPower(0, source, nil) != source.Card.Power+2 || source.AttackBonus != 1 || len(p.Exile) != 1 {
		t.Fatal("its own sacrifice ability must still grow the spell")
	}
}

func TestRoyalAuditDivineRadianceUsesTheChosenGeometry(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	skill := readySkill(baseCard(t, "3511101"), 0)
	p.Skills[0] = skill
	setAllElements(p, 20)
	placeUnit(baseCard(t, "1421101"), 1, 0, 0, e)
	placeUnit(baseCard(t, "1421101"), 1, 0, 2, e)
	placeUnit(baseCard(t, "1421101"), 1, 2, 0, e)
	before := p.Elements[model.ElementLight]
	if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": skill.InstanceID, "target_type": "unit", "target_col": float64(0), "target_row": float64(0)}}); err != nil {
		t.Fatal(err)
	}
	if p.Elements[model.ElementLight] != before || skill.IsHorizontal || e.State.PendingAction == nil {
		t.Fatal("area must be chosen before payment and tap")
	}
	resolvePendingSelection(t, e, 0, string(SpellAreaColumn))
	if e.State.PendingSpell == nil || e.State.PendingSpell.Target.Area != SpellAreaColumn {
		t.Fatal("the pending cast must retain its selected area")
	}
	if got := e.pendingSpellAffectedPositions(e.State.PendingSpell); len(got) != 3 {
		t.Fatalf("column preview must include three cells, got %v", got)
	}
	if got := e.spellAffectedUnits(1, skill, e.State.PendingSpell.Target); len(got) != 2 {
		t.Fatalf("column must hit both units in that column, got %d", len(got))
	}
}

func TestRoyalAuditTurnExpiryAdvancesOnEitherPlayer(t *testing.T) {
	for _, first := range []int{0, 1} {
		for _, current := range []int{0, 1} {
			e := setupEffectTest(t)
			e.State.FirstPlayer, e.State.CurrentTurn, e.State.TurnNumber = first, current, 3
			unit := e.State.Players[0].Hero
			unit.Statuses[temporaryDamageAndNegativeImmunityUntilStatus] = e.turnOrdinal() + 1
			if !e.temporaryDamageAndNegativeImmunityActive(unit) {
				t.Fatal("protection must apply in the granting turn")
			}
			for step := 1; step <= 2; step++ {
				e.State.CurrentTurn = 1 - e.State.CurrentTurn
				if e.State.CurrentTurn == first {
					e.State.TurnNumber++
				}
				if got := e.temporaryDamageAndNegativeImmunityActive(unit); got != (step == 1) {
					t.Fatal("unqualified next-turn protection must cover exactly the following player's turn")
				}
			}
		}
	}
}

func TestRoyalAuditFrozenResetDoesNotSpendWendyPayment(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	wendy := placeUnit(baseCard(t, "1211103"), 0, 0, 0, e)
	skill := readySkill(baseCard(t, "3221001"), 0)
	skill.IsHorizontal = true
	skill.Statuses[StatusFreeze] = 1
	p.Skills[0] = skill
	p.Elements[model.ElementWater] = 2
	if err := (Card1211103SeaHeroineCoralWendy{}).OnSpellCast(&EffectContext{Engine: e, Source: wendy, Target: skill, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"cast_player": 0}}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction != nil || p.Elements[model.ElementWater] != 2 || wendy.UsedThisTurn != 0 {
		t.Fatal("a frozen spell must not offer or spend a reset payment")
	}
}

func TestRoyalAuditInfernoSeparatesDefenseFromBurnPrevention(t *testing.T) {
	for _, defensePower := range []int{9, 12} {
		e := setupEffectTest(t)
		attacker, defender := e.State.Players[0], e.State.Players[1]
		inferno := readySkill(baseCard(t, "3111101"), 0)
		inferno.PowerBonus = 8 - inferno.Card.Power
		block := readySkill(baseCard(t, "3221103"), 1)
		block.PowerBonus = defensePower - block.Card.Power
		attacker.Skills[0], defender.Skills[0] = inferno, block
		setAllElements(attacker, 20)
		setAllElements(defender, 20)
		target := placeUnit(baseCard(t, "1421101"), 1, 0, 0, e)
		target.CurrentLife = 10
		if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": inferno.InstanceID, "target_type": "unit", "target_col": float64(0), "target_row": float64(0)}}); err != nil {
			t.Fatal(err)
		}
		if err := e.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{"skill_ids": []any{block.InstanceID}}}); err != nil {
			t.Fatal(err)
		}
		if target.CurrentLife != 10 {
			t.Fatal("ordinary attack damage must be blocked once defense reaches attack power")
		}
		wantBurn := 1
		if defensePower == 12 {
			wantBurn = 0
		}
		if target.Statuses[StatusBurn] != wantBurn {
			t.Fatalf("defense %d should leave %d burn, got %v", defensePower, wantBurn, target.Statuses)
		}
	}
}

func TestRoyalAuditPrayerNormalUseAndFreeCounterHaveDifferentRewards(t *testing.T) {
	for _, mode := range []string{"normal", "set"} {
		e := setupEffectTest(t)
		p := e.State.Players[0]
		prayer := NewCardInstance(baseCard(t, "2521110"), 0, 1)
		spirit := NewCardInstance(baseCard(t, "1521109"), 0, 1)
		p.Hand = []*CardInstance{prayer}
		p.Deck = []*CardInstance{spirit}
		p.Elements[model.ElementLight] = 1
		if err := e.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{"instance_id": prayer.InstanceID}}); err != nil {
			t.Fatal(err)
		}
		resolvePendingSelection(t, e, 0, mode)
		if mode == "normal" {
			resolvePendingSelection(t, e, 0, prayer.InstanceID)
			if p.Elements[model.ElementLight] != 0 || len(p.Hand) != 1 || p.Hand[0] != spirit || spirit.Statuses["入场费用光-1"] != 0 {
				t.Fatal("normal use pays one light and grants no spirit discount")
			}
			continue
		}
		resolvePendingSelection(t, e, 0, "set")
		p.Elements[model.ElementLight] = 0
		enemy := readySkill(baseCard(t, "3121001"), 1)
		if !e.promptCounterTrapIfEligible(prayer, TriggerOnSpellCast, enemy, map[string]any{"cast_player": 1, "purpose": string(skillPurposeAttack)}, nil) {
			t.Fatal("the set prayer must trigger without any elements")
		}
		resolvePendingSelection(t, e, 0, prayer.InstanceID)
		if len(p.Hand) != 1 || p.Hand[0] != spirit || spirit.Statuses["入场费用光-1"] != 1 {
			t.Fatal("counter use must flip and discount the spirit")
		}
	}
}

func TestRoyalAuditDecliningCursedFireDoesNotDiscountTheScroll(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	source := readySkill(baseCard(t, "3121110"), 0)
	scroll := NewCardInstance(baseCard(t, "2121109"), 0, 1)
	p.Hand = nil
	p.Deck = []*CardInstance{scroll}
	placeUnit(baseCard(t, "1421101"), 1, 0, 0, e)
	if err := (Card3121110CursedFire{}).OnSpellHit(&EffectContext{Engine: e, Source: source, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"spell_source": source, "attacker": 0}}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction == nil {
		t.Fatal("the immediate use must be offered")
	}
	resolvePendingSelection(t, e, 0)
	if cost := e.effectiveCardPlayCost(p, scroll); !sameElementCost(cost, scroll.Card.ElementsCost) {
		t.Fatalf("declining immediate use must preserve the scroll's printed cost, got %v", cost)
	}
}

func TestRoyalAuditPlainActiveCanBeUsedAgainAfterResetWithoutGrantingLoad(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	conductor := placeUnit(baseCard(t, "1011102"), 0, 0, 0, e)
	p.Elements = map[string]int{}
	for use := 1; use <= 2; use++ {
		if err := e.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{"instance_id": conductor.InstanceID, "ability_type": "per_turn"}}); err != nil {
			t.Fatalf("plain active use %d failed after reset: %v", use, err)
		}
		if !conductor.IsHorizontal || totalElementCost(p.Elements) != 0 {
			t.Fatal("the active's tap payment must not grant the printed load")
		}
		e.resetCard(conductor)
	}
	if conductor.Statuses[conductorConsumedCountStatus] != 2 {
		t.Fatal("the conductor must include itself in its global consume count")
	}
	if info := e.cardToInfo(conductor); info["per_turn_unlimited"] != true {
		t.Fatal("the frontend must see the repeated-use permission")
	}
}

func TestRoyalAuditAgingPotionJumpsToTheNextPrintedMilestone(t *testing.T) {
	e := setupEffectTest(t)
	treant := placeUnit(baseCard(t, "1421003"), 0, 0, 0, e)
	potion := NewCardInstance(baseCard(t, "2421106"), 0, 1)
	if err := (Card2421106AgingPotion{}).OnUseItem(e.skillContext(0, potion)); err != nil {
		t.Fatal(err)
	}
	resolvePendingSelection(t, e, 0, treant.InstanceID)
	if treant.Statuses[StatusMastery] != 2 {
		t.Fatal("the next printed mastery milestone is 2, not one generic mark")
	}
}

func TestRoyalAuditBoundArcaneFlowKeepsItsSelfAura(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	p.Hero = NewCardInstance(baseCard(t, "4011101"), 0, 1)
	p.Hero.Position = &Position{Col: 1, Row: 1}
	p.Units[1][1] = p.Hero
	flow := readySkill(baseCard(t, "3021106"), 0)
	p.Hero.BoundSkills = []*CardInstance{flow}
	if power := e.effectiveSpellPower(0, flow, nil); power != 2*flow.Card.Power {
		t.Fatalf("the actual attached spell must retain its self aura, got power %d", power)
	}
}

func TestRoyalAuditDeepSwordDrawAndAreaSweepIncludeStealth(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	sword := NewCardInstance(baseCard(t, "2211101"), 0, 1)
	p.Hand = nil
	p.Deck = []*CardInstance{sword}
	e.State.Players[1].Skills[0] = readySkill(baseCard(t, "3421107"), 1)
	hidden := placeUnit(baseCard(t, "1421101"), 1, 0, 0, e)
	hidden.Statuses[StatusStealth] = 1
	hidden.CurrentLife = 5
	areaSpell := readySkill(baseCard(t, "3511101"), 0)
	hero := e.State.Players[1].Hero
	if err := e.validateSpellTarget(0, areaSpell, SpellTarget{Type: "unit", Position: *hero.Position, Area: SpellAreaColumn}); err != nil {
		t.Fatalf("a hidden-only front row must not block an area center in the visible spell range: %v", err)
	}
	e.drawCards(0, 1)
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "deep_sword_reveal_damage" {
		t.Fatal("drawing the sword itself must expose its reveal choice")
	}
	resolvePendingSelection(t, e, 0, sword.InstanceID)
	if hidden.CurrentLife != 3 {
		t.Fatal("an untargeted range sweep must include a stealthed unit before the visible front row")
	}
	p.Hand = nil
	p.Equipment[0] = sword
	p.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
	e.drawCards(0, 1)
	if e.State.PendingAction != nil {
		t.Fatal("an equipped sword must not react to drawing a different card")
	}
}

func TestRoyalAuditCloudTopCanDeclineAndPreservesDrawContinuations(t *testing.T) {
	for _, accept := range []bool{false, true} {
		e := setupEffectTest(t)
		p, opponent := e.State.Players[0], e.State.Players[1]
		placeUnit(baseCard(t, "1311102"), 0, 0, 0, e)
		own := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		foreign := NewCardInstance(baseCard(t, "1021002"), 1, 1)
		p.Hand = nil
		p.Deck = []*CardInstance{own}
		opponent.Deck = []*CardInstance{foreign}
		scroll := NewCardInstance(baseCard(t, "2021019"), 0, 1)
		if err := (Card2021019CursedScroll{}).OnUseItem(e.skillContext(0, scroll)); err != nil {
			t.Fatal(err)
		}
		if e.State.PendingAction == nil || len(opponent.Deck) != 1 {
			t.Fatal("the replacement draw must wait for consent")
		}
		if accept {
			resolvePendingSelection(t, e, 0, "draw:1")
		} else {
			resolvePendingSelection(t, e, 0)
		}
		if !p.DiscardAtTurnEnd[own.InstanceID] || p.DiscardAtTurnEnd[foreign.InstanceID] != accept {
			t.Fatal("the original effect must track precisely the cards drawn after the choice")
		}
		if (foreign.OwnerID == 0) != accept || (len(opponent.Deck) == 0) != accept {
			t.Fatal("declining must preserve the opponent's card and its ownership")
		}
	}
}

func TestRoyalAuditFreeArmorStillUsesNormalEquipmentReplacementRules(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	old := NewCardInstance(baseCard(t, "2121013"), 0, 1)
	old.IsHorizontal = true
	old.SlotIndex = 0
	p.Equipment[0] = old
	fresh := NewCardInstance(baseCard(t, "2121013"), 0, 1)
	p.Hand = []*CardInstance{fresh}
	if got := e.equipCardFromHandOrDeckFree(0, "2121013"); got != nil || len(p.Hand) != 1 || p.Equipment[1] != nil || e.State.PendingAction != nil {
		t.Fatal("free entry cannot bypass the subtype limit or replace a horizontal armor")
	}
	old.IsHorizontal = false
	e.equipCardFromHandOrDeckFree(0, "2121013")
	if e.State.PendingAction == nil {
		t.Fatal("a legal vertical replacement must be offered")
	}
	resolvePendingSelection(t, e, 0, old.InstanceID)
	if p.Equipment[0] != fresh || len(p.Hand) != 0 || !containsCardInstance(p.Graveyard, old) {
		t.Fatal("free entry must commit the selected legal replacement")
	}
}

func TestRoyalAuditGameLongManesRewardSurvivesLeavingAndControlChange(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	skill := readySkill(baseCard(t, "3221001"), 0)
	p.Skills[0] = skill
	item := NewCardInstance(baseCard(t, "2211102"), 0, 1)
	if err := (Card2211102ManesArbitration{}).OnEquip(e.skillContext(0, item)); err != nil {
		t.Fatal(err)
	}
	resolvePendingSelection(t, e, 0, "one_water_power_attack")
	resolvePendingSelection(t, e, 0, skill.InstanceID)
	p.Skills[0] = nil
	clean := e.leaveBattlefield(skill)
	clean.OwnerID = 1
	e.State.Players[1].Hero.BoundSkills = []*CardInstance{clean}
	if e.effectiveSpellPower(1, clean, nil) != clean.Card.Power+3 || e.temporarySpellDamageBonus(1, clean) != 1 {
		t.Fatal("a game-long targeted reward belongs to the game and survives clean-zone movement and transfer")
	}
}

func TestRoyalAuditBoundObserversAndAreaEffectsResolveOnce(t *testing.T) {
	e := setupEffectTest(t)
	moon := readySkill(baseCard(t, "3621108"), 0)
	moon.IsHorizontal = true
	e.State.Players[0].Hero.BoundSkills = []*CardInstance{moon}
	attack := readySkill(baseCard(t, "3121001"), 0)
	defense := readySkill(baseCard(t, "3221001"), 1)
	defense.Statuses[StatusWeaken] = 1
	e.triggerFieldEffectsWithData(TriggerOnDefend, 0, attack, map[string]any{"attack_skill": attack, "defense_skills": []*CardInstance{defense}})
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "moonshadow_reset" {
		t.Fatal("a bound Moonshadow must still observe the battle")
	}
	resolvePendingSelection(t, e, 0, moon.InstanceID)
	if moon.IsHorizontal || moon.UsedThisTurn != 1 {
		t.Fatal("the bound observer must resolve exactly once")
	}
	bound := readySkill(baseCard(t, "3121001"), 1)
	e.State.Players[1].Hero.BoundSkills = []*CardInstance{bound}
	erosion := readySkill(baseCard(t, "3421106"), 0)
	if err := (Card3421106RottingErosion{}).OnSpellHit(&EffectContext{Engine: e, Source: erosion, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"spell_source": erosion}}); err != nil {
		t.Fatal(err)
	}
	if bound.Statuses[StatusWeaken] != max(0, erosion.Card.Attack) {
		t.Fatal("a bound spell must not be counted twice by an all-spells effect")
	}
}

func TestRoyalAuditNextSpellBonusWaitsUntilAnActualMatchingUse(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	e.addNextElementSpellPowerBonus(0, model.ElementFire, 3)
	e.State.TurnNumber += 4
	e.clearExpiredTemporaryModifiers(0)
	fire := readySkill(baseCard(t, "3121001"), 0)
	water := readySkill(baseCard(t, "3221001"), 0)
	e.consumeNextSpellPowerBonuses(p, water)
	if e.temporarySpellPowerBonus(0, fire) != 3 {
		t.Fatal("the unqualified next-spell bonus must survive time and unrelated elements")
	}
	e.consumeNextSpellPowerBonuses(p, fire)
	if len(p.TempModifiers) != 0 {
		t.Fatal("the matching use must spend the bonus exactly once")
	}
}
