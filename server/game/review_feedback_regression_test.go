package game

import (
	"eraofarcane/model"
	"testing"
)

func TestReviewFeedbackEmbersDeclarationWaitsBeforeDiscardAndReset(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	embers := readySkill(baseCard(t, "3121105"), 0)
	p.Skills[0] = embers
	p.Hero.IsHorizontal = true
	for len(p.Hand) < 6 {
		p.Hand = append(p.Hand, e.newCardInstance(baseCard(t, "1021001"), 0, 1))
	}
	if err := e.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "embers_free_cast_target" || !p.Hero.IsHorizontal || len(p.Hand) != 6 {
		t.Fatal("declaration must precede hand limit and reset")
	}
	hero := e.State.Players[1].Hero
	if !candidateContains(e.State.PendingAction.Candidates, hero.InstanceID) {
		t.Fatal("a hero in spell range must be selectable")
	}
	resolvePendingSelection(t, e, 0, hero.InstanceID)
	if e.State.PendingSpell == nil || e.State.CurrentTurn != 0 || !p.Hero.IsHorizontal {
		t.Fatal("declaration must wait for the complete spell")
	}
	if err := e.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "discard" || !p.Hero.IsHorizontal {
		t.Fatal("discard must follow the spell and precede reset")
	}
	resolvePendingSelection(t, e, 0, p.Hand[0].InstanceID)
	if e.State.CurrentTurn != 1 || p.Hero.IsHorizontal || embers.IsHorizontal {
		t.Fatal("finish cleanup and switch only after the discard")
	}
}

func TestReviewFeedbackEmbersCannotUseTheLaterReset(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	embers := readySkill(baseCard(t, "3121105"), 0)
	embers.IsHorizontal = true
	p.Skills[0] = embers
	if err := e.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction != nil || e.State.CurrentTurn != 1 || embers.IsHorizontal {
		t.Fatal("a horizontal Embers must skip the declaration, then reset normally")
	}
}

func TestReviewFeedbackPainterCopiesOntoItself(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	painter := placeUnit(baseCard(t, "1321115"), 0, 0, 0, e)
	warrior := placeUnit(baseCard(t, "1321111"), 0, 2, 0, e)
	p.Equipment[0] = e.newCardInstance(&model.Card{Number: "test_thunderlight", Type: model.CardTypeItem, ElementsGain: map[string]int{model.ElementAir: 1, model.ElementLight: 1}}, 0, 1)
	beforePainter, beforeWarrior := painter.CurrentLife, warrior.CurrentLife
	if err := (Card1321115SkyPainter{}).OnEnter(e.skillContext(0, painter)); err != nil {
		t.Fatal(err)
	}
	resolvePendingSelection(t, e, 0, warrior.InstanceID)
	resolvePendingSelection(t, e, 0, "life#0")
	if painter.CurrentLife != beforePainter+2 || warrior.CurrentLife != beforeWarrior {
		t.Fatal("the copied 'this card' must be the painter")
	}
	copy := *warrior.Card
	copy.ElementsCost = map[string]int{model.ElementAir: 6}
	copy.ElementsExpense = map[string]int{}
	warrior.Card = &copy
	if canSkyPainterCopy(e, painter, warrior) {
		t.Fatal("entry cost, rather than use cost, determines eligibility")
	}
}

func TestReviewFeedbackArcaneImpactUsesActualPayments(t *testing.T) {
	for _, tc := range []struct {
		name       string
		learn, use map[string]int
		bonus      int
	}{
		{"both arcane", map[string]int{model.ElementArcane: 3}, map[string]int{model.ElementArcane: 2}, 1},
		{"colored learning", map[string]int{model.ElementFire: 3}, map[string]int{model.ElementArcane: 2}, 0},
		{"mixed use", map[string]int{model.ElementArcane: 3}, map[string]int{model.ElementArcane: 1, model.ElementFire: 1}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := setupEffectTest(t)
			p := e.State.Players[0]
			skill := readySkill(baseCard(t, "3021101"), 0)
			p.SkillPool = append(p.SkillPool, skill)
			setAllElements(p, 20)
			if !e.payCostForCardAction(p, skill, skill.Card.ElementsCost, skill.Card.ElementsCost, paymentPurposeLearn, ActionMessage{Data: map[string]any{"payment": tc.learn}}) {
				t.Fatal("learn payment")
			}
			p.Skills[0] = skill
			if !e.payCostForCardAction(p, skill, skill.Card.ElementsExpense, skill.Card.ElementsExpense, paymentPurposeUse, ActionMessage{Data: map[string]any{"payment": tc.use}}) {
				t.Fatal("use payment")
			}
			if got := e.effectiveSpellPower(0, skill, nil); got != skill.Card.Power+tc.bonus {
				t.Fatalf("power %d", got)
			}
			if got := e.effectiveSpellDamage(0, skill, skill.Card.Attack, nil); got != skill.Card.Attack+tc.bonus {
				t.Fatalf("damage %d", got)
			}
		})
	}
}

func TestReviewFeedbackBloodSandChoosesLoadAndCanPayLastLife(t *testing.T) {
	e := setupEffectTest(t)
	source := readySkill(baseCard(t, "3411102"), 0)
	unit := placeUnit(baseCard(t, "1021001"), 0, 0, 0, e)
	unit.CurrentLife = 1
	unit.ElementsGainBonus = map[string]int{model.ElementFire: 1, model.ElementWater: 1}
	paid, err := e.applyBloodSandPayments(0, []string{unit.InstanceID}, map[string]any{"payments": []any{
		map[string]any{"instance_id": unit.InstanceID, "amount": float64(1), "mode": "load", "element": model.ElementWater},
		map[string]any{"instance_id": unit.InstanceID, "amount": float64(1), "mode": "life"},
	}}, source)
	if err != nil || paid != 2 || unit.CurrentLife != 0 || e.effectiveElementsGain(unit)[model.ElementWater] != 0 || e.effectiveElementsGain(unit)[model.ElementFire] != 1 {
		t.Fatalf("payment=%d err=%v", paid, err)
	}
	e.settleTurnDeaths()
	if e.State.Players[0].Units[0][0] != nil {
		t.Fatal("a zero-life payer must die")
	}
}

func TestReviewFeedbackSilverleafChoosesOnePerDiscardAction(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	silver := readySkill(baseCard(t, "3311102"), 0)
	p.Skills[0] = silver
	a, b := e.newCardInstance(baseCard(t, "1021001"), 0, 1), e.newCardInstance(baseCard(t, "1021002"), 0, 1)
	p.Hand = []*CardInstance{a, b}
	e.discardAllHandCards(0)
	if e.State.PendingAction == nil || len(e.State.PendingAction.Candidates) != 2 {
		t.Fatal("one choice must contain the entire discard batch")
	}
	resolvePendingSelection(t, e, 0, a.InstanceID)
	if e.State.PendingAction != nil || len(silver.UnderCards) != 1 || len(p.Graveyard) != 1 || p.Graveyard[0] != b {
		t.Fatal("one batch must store exactly one card")
	}
}

func TestReviewFeedbackRainbowRingCanPayWithArcane(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	ring := e.newCardInstance(baseCard(t, "2511102"), 0, 1)
	p.Equipment[0] = ring
	p.Elements[model.ElementArcane] = 1
	if err := (Card2511102FiveRainbowRing{}).OnPerTurn(e.skillContext(0, ring)); err != nil {
		t.Fatal(err)
	}
	resolvePendingSelection(t, e, 0, model.ElementFire)
	if p.Elements[model.ElementArcane] != 0 || ring.Statuses[fiveRainbowMarkerStatus(model.ElementFire)] != 1 {
		t.Fatal("chosen requirement may be paid with arcane")
	}
}

func TestReviewFeedbackMindSeaChoosesAfterUseAndKeepsFirstSingle(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	skill := readySkill(baseCard(t, "3211101"), 0)
	p.Skills[0] = skill
	setAllElements(p, 20)
	a := placeUnit(baseCard(t, "1421101"), 1, 0, 0, e)
	b := placeUnit(baseCard(t, "1421101"), 1, 2, 0, e)
	a.CurrentLife = 10
	b.CurrentLife = 10
	cast := func() {
		t.Helper()
		if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": skill.InstanceID, "target_type": "unit", "target_col": float64(0), "target_row": float64(0)}}); err != nil {
			t.Fatal(err)
		}
	}
	cast()
	if e.State.PendingAction != nil || e.State.PendingSpell == nil || e.State.PendingSpell.TotalPower != skill.Card.Power {
		t.Fatal("first cast should be single with printed power")
	}
	if err := e.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
		t.Fatal(err)
	}
	if b.CurrentLife != 10 || e.State.PendingAction == nil || e.State.PendingAction.Type != "mind_sea_maze_next_area" {
		t.Fatal("choose the new area only after finishing the first use")
	}
	resolvePendingSelection(t, e, 0, string(SpellAreaFrontRow))
	skill.IsHorizontal = false
	delete(skill.Statuses, StatusStun)
	cast()
	if e.State.PendingSpell.TotalPower != skill.Card.Power+1 || e.effectiveSpellArea(skill) != SpellAreaFrontRow {
		t.Fatal("next cast must use the stored area and +1 power")
	}
}

func TestReviewFeedbackArcaneOnenessSnapshotsOnlyCurrentUse(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	skill := readySkill(baseCard(t, "3011101"), 0)
	p.Skills[0] = skill
	p.Elements[model.ElementArcane] = 7
	p.Deck = []*CardInstance{e.newCardInstance(baseCard(t, "3021005"), 0, 1), e.newCardInstance(baseCard(t, "3021101"), 0, 1), e.newCardInstance(baseCard(t, "3121001"), 0, 1)}
	if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": skill.InstanceID, "target_type": "unit", "target_col": float64(1), "target_row": float64(1)}}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingSpell == nil || e.State.PendingSpell.TotalPower != 2 {
		t.Fatal("reveal should provide two power")
	}
	p.Deck = []*CardInstance{e.newCardInstance(baseCard(t, "3121001"), 0, 1)}
	e.refreshPendingSpellPowerForModifiedSkill(0, skill)
	if e.State.PendingSpell.TotalPower != 2 {
		t.Fatal("counter recalculation must retain the revealed count")
	}
	if err := e.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
		t.Fatal(err)
	}
	if e.effectiveSpellPower(0, skill, nil) != skill.Card.Power {
		t.Fatal("reveal bonus must disappear after the use")
	}
}

func TestReviewFeedbackSkyCopiesCanAttackDefendAndBoost(t *testing.T) {
	for _, purpose := range []skillPurpose{skillPurposeAttack, skillPurposeAttackBoost, skillPurposeDefend, skillPurposeDefenseBoost} {
		t.Run(string(purpose), func(t *testing.T) {
			e := setupEffectTest(t)
			p := e.State.Players[0]
			setAllElements(p, 20)
			sky := readySkill(baseCard(t, "3311101"), 0)
			base := readySkill(baseCard(t, "3321002"), 0)
			main := readySkill(baseCard(t, "3321002"), 0)
			base.IsHorizontal = true
			base.Statuses[StatusCooldown] = 2
			p.Skills[0], p.Skills[1], p.Skills[2] = sky, base, main
			data := map[string]any{"spell_copies": map[string]any{sky.InstanceID: base.InstanceID}, "target_type": "unit", "target_col": float64(1), "target_row": float64(1)}
			action := "cast_spell"
			if purpose == skillPurposeAttack {
				data["instance_id"] = sky.InstanceID
				data["boost_ids"] = []any{main.InstanceID}
			}
			if purpose == skillPurposeAttackBoost {
				data["instance_id"] = main.InstanceID
				data["boost_ids"] = []any{sky.InstanceID}
			}
			if purpose == skillPurposeDefend || purpose == skillPurposeDefenseBoost {
				action = "defend"
				e.State.Phase = PhaseDefenseWindow
				e.State.CurrentTurn = 1
				e.State.PendingSpell = &SpellCast{AttackerID: 1, Skill: readySkill(baseCard(t, "3321002"), 1), Target: SpellTarget{Type: "unit", Position: *p.Hero.Position}, TotalPower: 5}
				if purpose == skillPurposeDefend {
					data["skill_ids"] = []any{sky.InstanceID}
					data["boost_ids"] = []any{main.InstanceID}
				} else {
					data["skill_ids"] = []any{main.InstanceID}
					data["boost_ids"] = []any{sky.InstanceID}
				}
			}
			if err := e.HandleAction(0, ActionMessage{Action: action, Data: data}); err != nil {
				t.Fatal(err)
			}
			if !sky.IsHorizontal || !main.IsHorizontal || base.Statuses[StatusCooldown] != 2 {
				t.Fatal("tap physical sources, leave copied source untouched")
			}
			if action == "cast_spell" && (e.State.PendingSpell == nil || e.State.PendingSpell.TotalPower != 6) {
				t.Fatal("copied spell can receive and provide boosts")
			}
			if action == "defend" && e.State.PendingSpell != nil {
				t.Fatal("copied defense plus boost should stop the attack")
			}
		})
	}
}

func TestReviewFeedbackWaterMirrorDefendsAndBoosts(t *testing.T) {
	for _, defend := range []bool{true, false} {
		t.Run(map[bool]string{true: "defend", false: "attack boost"}[defend], func(t *testing.T) {
			e := setupEffectTest(t)
			p := e.State.Players[0]
			setAllElements(p, 20)
			recorded := readySkill(baseCard(t, "3221001"), 0)
			e.recordSpellCast(0, recorded)
			mirror := e.newCardInstance(baseCard(t, "2221104"), 0, 1)
			p.Hand = append(p.Hand, mirror)
			main := readySkill(baseCard(t, "3321002"), 0)
			p.Skills[0] = main
			if defend {
				e.State.Phase = PhaseDefenseWindow
				e.State.CurrentTurn = 1
				e.State.PendingSpell = &SpellCast{AttackerID: 1, Skill: readySkill(baseCard(t, "3321002"), 1), Target: SpellTarget{Type: "unit", Position: *p.Hero.Position}, TotalPower: 1}
				if err := e.HandleAction(0, ActionMessage{Action: "defend", Data: map[string]any{"scroll_ids": []any{mirror.InstanceID}, "boost_ids": []any{main.InstanceID}}}); err != nil {
					t.Fatal(err)
				}
				if e.State.PendingSpell != nil {
					t.Fatal("mirror must defend and accept boosts")
				}
			} else {
				if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": main.InstanceID, "target_type": "unit", "target_col": float64(1), "target_row": float64(1), "boost_ids": []any{mirror.InstanceID}}}); err != nil {
					t.Fatal(err)
				}
				if e.State.PendingSpell == nil || e.State.PendingSpell.TotalPower != main.Card.Power+recorded.Card.Power {
					t.Fatal("mirror must supply copied boost power")
				}
			}
			if len(p.Graveyard) != 1 || p.Graveyard[0].Card.Number != "2221104" {
				t.Fatal("the physical scroll, rather than the copied skill, must enter graveyard")
			}
		})
	}
}

func TestReviewFeedbackWeakenAppliesAfterPowerScaling(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	skill := readySkill(baseCard(t, "3321002"), 0)
	p.Skills[0] = skill
	p.Equipment[0] = e.newCardInstance(baseCard(t, "2611102"), 0, 1)
	skill.Statuses[StatusWeaken] = 1
	if got := e.effectiveSpellPower(0, skill, nil); got != 1 {
		t.Fatalf("ceil(3/2)-1 must be 1, got %d", got)
	}
}

func TestReviewFeedbackKnightStopsAtZeroAndSnowfieldDoesNotCycle(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	knight := e.newCardInstance(baseCard(t, "1521108"), 0, 1)
	copy := *knight.Card
	copy.Life = 1
	knight.Card = &copy
	p.Graveyard = append(p.Graveyard, knight)
	if err := (Card1521108ContradictoryKnight{}).OnDeath(e.skillContext(0, knight)); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction != nil || len(p.Graveyard) != 1 || knight.Card.Life != 0 {
		t.Fatal("zero maximum life must remain in graveyard without a resurrection choice")
	}
	snow := readySkill(baseCard(t, "3211102"), 0)
	p.Skills[0] = snow
	for _, count := range []int{6, 10, 15} {
		snow.Statuses[dragonSnowfieldTriggerStatus] = count
		e.promptDragonSnowfieldSummon(0, snow)
		if e.State.PendingAction != nil {
			t.Fatalf("must not repeat at count %d", count)
		}
	}
}

func TestReviewFeedbackConsumeCostsAdvanceMasteryWithoutLoad(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	unit := placeUnit(baseCard(t, "1421003"), 0, 0, 0, e)
	unit.IsHorizontal = false
	conductor := placeUnit(baseCard(t, "1011102"), 0, 2, 0, e)
	e.consumeCardForEffectWithTriggers(0, unit, nil, "1121109")
	if unit.Statuses[StatusMastery] != 1 || !unit.IsHorizontal || totalElementCost(p.Elements) != 0 {
		t.Fatal("a consume cost advances mastery but grants no printed load")
	}
	if conductor.Statuses[conductorConsumedCountStatus] != 1 {
		t.Fatal("a consume cost must notify the conductor's field-wide observer once")
	}
}

func TestReviewFeedbackCombinedPaymentChoosesImpactAttributionBeforeCommit(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	main := readySkill(baseCard(t, "3021101"), 0)
	impact := readySkill(baseCard(t, "3021101"), 0)
	impact.learnPayment = map[string]int{model.ElementArcane: 3}
	p.Skills[0], p.Skills[1] = main, impact
	mainCard := *main.Card
	mainCard.ElementsExpense = map[string]int{model.ElementArcane: 2}
	main.Card = &mainCard
	p.Elements[model.ElementArcane] = 2
	p.Elements[model.ElementFire] = 2
	if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": main.InstanceID, "target_type": "unit", "target_col": float64(1), "target_row": float64(1), "boost_ids": []any{impact.InstanceID}, "payment": map[string]int{model.ElementArcane: 2, model.ElementFire: 2}}}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "spell_payment_attribution" || main.IsHorizontal || impact.IsHorizontal || p.Elements[model.ElementArcane] != 2 {
		t.Fatal("ambiguous attribution must be chosen before spending or tapping")
	}
	if err := e.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{impact.InstanceID}, "payment": map[string]int{model.ElementArcane: 2}}}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingSpell == nil || e.State.PendingSpell.TotalPower != main.Card.Power+impact.Card.Power+1 || p.Elements[model.ElementArcane] != 0 || p.Elements[model.ElementFire] != 0 {
		t.Fatal("the chosen arcane payment must grant the impact bonus exactly once")
	}
}

func TestReviewFeedbackUseSnapshotClearsAfterPetrifyOrReplacement(t *testing.T) {
	for _, replaced := range []bool{false, true} {
		t.Run(map[bool]string{false: "petrified", true: "replaced"}[replaced], func(t *testing.T) {
			e := setupEffectTest(t)
			p := e.State.Players[0]
			source := readySkill(baseCard(t, "3011101"), 0)
			p.Skills[0] = source
			count := 2
			source.arcaneOnenessReveal = &count
			original := &SpellCast{AttackerID: 0, Skill: source, Target: SpellTarget{Type: "unit", Position: *e.State.Players[1].Hero.Position}, TotalPower: 2}
			e.State.PendingSpell = original
			e.State.Phase = PhaseDefenseWindow
			if replaced {
				if err := e.replacePendingSpell(func() error {
					e.State.PendingSpell = &SpellCast{AttackerID: 0, Skill: readySkill(baseCard(t, "3321002"), 0), Target: original.Target, TotalPower: 1}
					return nil
				}); err != nil {
					t.Fatal(err)
				}
				if source.arcaneOnenessReveal == nil {
					t.Fatal("replacement must finish before the original use is cleaned up")
				}
				if err := e.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
					t.Fatal(err)
				}
			} else {
				source.Statuses[StatusPetrify] = 1
				e.completePendingSpell(original)
			}
			if source.arcaneOnenessReveal != nil {
				t.Fatal("per-use snapshot must not survive the completed use")
			}
		})
	}
}

func TestReviewFeedbackPawnNeverDamagesTheOtherGrid(t *testing.T) {
	e := setupEffectTest(t)
	pawn := e.newCardInstance(baseCard(t, "1001101"), 0, 1)
	pawn.Position = &Position{Col: 1, Row: 1}
	own := placeUnit(baseCard(t, "1021001"), 0, 0, 1, e)
	own.CurrentLife = 3
	enemy := placeUnit(baseCard(t, "1021001"), 1, 0, 1, e)
	enemy.CurrentLife = 3
	if err := (Card1001101AbandonedPawn{}).OnDeath(e.skillContext(0, pawn)); err != nil {
		t.Fatal(err)
	}
	if own.CurrentLife != 2 || enemy.CurrentLife != 3 {
		t.Fatal("adjacency must remain inside the pawn owner's grid")
	}
}

func TestReviewFeedbackStackedSilverleafRefreshesRemainingDiscards(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	first, second := readySkill(baseCard(t, "3311102"), 0), readySkill(baseCard(t, "3311102"), 0)
	p.Skills[0], p.Skills[1] = first, second
	a, b := e.newCardInstance(baseCard(t, "1021001"), 0, 1), e.newCardInstance(baseCard(t, "1021002"), 0, 1)
	p.Hand = []*CardInstance{a, b}
	e.discardAllHandCards(0)
	resolvePendingSelection(t, e, 0, a.InstanceID)
	if e.State.PendingAction == nil || len(e.State.PendingAction.Candidates) != 1 || !candidateContains(e.State.PendingAction.Candidates, b.InstanceID) {
		t.Fatal("the second copy must only offer the card still in graveyard")
	}
	resolvePendingSelection(t, e, 0, b.InstanceID)
	if len(first.UnderCards) != 1 || len(second.UnderCards) != 1 || len(p.Graveyard) != 0 {
		t.Fatal("both legal discard triggers must resolve")
	}
}

func TestReviewFeedbackCombinedCostsKeepEachCardsPaymentConstraint(t *testing.T) {
	for _, number := range []string{"3011101", "3021103"} {
		t.Run(number, func(t *testing.T) {
			e := setupEffectTest(t)
			p := e.State.Players[0]
			source := readySkill(baseCard(t, number), 0)
			cost := e.effectiveSkillUseCost(p, source)
			parts := []spellPaymentPart{{source, cost}}
			if _, ok := e.constrainedSpellPayments(p, parts, map[string]int{model.ElementFire: totalElementCost(cost)}, nil); ok {
				t.Fatal("strict or distinct costs must reject an all-fire payment")
			}
			good := map[string]int{}
			if number == "3011101" {
				good[model.ElementArcane] = totalElementCost(cost)
			} else {
				for _, elem := range model.AllElements[:totalElementCost(cost)] {
					good[elem] = 1
				}
			}
			if _, ok := e.constrainedSpellPayments(p, parts, good, nil); !ok {
				t.Fatal("valid constrained payment rejected")
			}
		})
	}
}

func TestReviewFeedbackCopiedUseCanCorrectAnInvalidCombinedPayment(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	setAllElements(p, 20)
	sky, base, drain := readySkill(baseCard(t, "3311101"), 0), readySkill(baseCard(t, "3021101"), 0), readySkill(baseCard(t, "3321002"), 0)
	p.Skills[0], p.Skills[1], p.Skills[2] = sky, base, drain
	skyCost := e.effectiveSkillUseCost(p, sky)
	drainCost := e.effectiveSkillUseCostForPurpose(p, drain, skillPurposeAttackBoost)
	bad := copyElements(skyCost)
	bad[model.ElementFire] += totalElementCost(drainCost)
	action := ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": sky.InstanceID, "spell_copies": map[string]any{sky.InstanceID: base.InstanceID}, "boost_ids": []any{drain.InstanceID}, "target_type": "unit", "target_col": float64(1), "target_row": float64(1), "payment": bad}}
	if err := e.HandleAction(0, action); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "copied_spell_payment" || sky.IsHorizontal || drain.IsHorizontal {
		t.Fatal("an invalid copied-use payment must remain editable before commit")
	}
	good := mergeElementCosts(skyCost, drainCost)
	if err := e.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{"selected": []any{"pay"}, "payment": good}}); err != nil {
		t.Fatal(err)
	}
	if !sky.IsHorizontal || !drain.IsHorizontal || base.IsHorizontal {
		t.Fatal("corrected payment should use the physical sources only")
	}
}

func TestReviewFeedbackStrictMainPaymentCannotBeReassignedToImpact(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	main, impact := readySkill(baseCard(t, "3011101"), 0), readySkill(baseCard(t, "3021101"), 0)
	p.Skills[0], p.Skills[1] = main, impact
	impact.learnPayment = map[string]int{model.ElementArcane: 3}
	p.Elements[model.ElementArcane] = totalElementCost(e.effectiveSkillUseCost(p, main))
	p.Elements[model.ElementFire] = totalElementCost(e.effectiveSkillUseCostForPurpose(p, impact, skillPurposeAttackBoost))
	p.Deck = []*CardInstance{e.newCardInstance(baseCard(t, "3021005"), 0, 1), e.newCardInstance(baseCard(t, "3021101"), 0, 1), e.newCardInstance(baseCard(t, "3121002"), 0, 1)}
	if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": main.InstanceID, "boost_ids": []any{impact.InstanceID}, "target_type": "unit", "target_col": float64(1), "target_row": float64(1)}}); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction != nil || e.State.PendingSpell == nil || e.State.PendingSpell.TotalPower != 2+impact.Card.Power {
		t.Fatal("strict main payment must stay reserved; colored Impact receives no bonus")
	}
}

func TestReviewFeedbackSorceriesCannotReceiveBoosts(t *testing.T) {
	e := setupEffectTest(t)
	p := e.State.Players[0]
	sorcery := readySkill(baseCard(t, "3021103"), 0)
	boost := readySkill(baseCard(t, "3321002"), 0)
	p.Skills[0], p.Skills[1] = sorcery, boost
	setAllElements(p, 20)
	if canSkillBeBoosted(sorcery) || e.cardToInfoForPlayer(p, sorcery)["can_be_boosted"] != false {
		t.Fatal("a sorcery must not offer boosting")
	}
	if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": sorcery.InstanceID, "target_type": "none", "boost_ids": []any{boost.InstanceID}}}); err == nil {
		t.Fatal("server must reject a boosted sorcery")
	}
	if sorcery.IsHorizontal || boost.IsHorizontal {
		t.Fatal("invalid boosted sorcery must not tap either source")
	}
}
