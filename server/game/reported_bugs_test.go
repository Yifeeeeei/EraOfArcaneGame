package game

import (
	"testing"

	"eraofarcane/cards"
	"eraofarcane/model"
)

func setupReportedBugEngine(t *testing.T) *Engine {
	t.Helper()
	if cards.CardDB == nil {
		if err := cards.LoadCards(); err != nil {
			t.Fatalf("load cards: %v", err)
		}
	}
	SetCardDB(cards.CardDB)
	RegisterAllCardEffects()

	engine := NewEngine("reported-bugs", nil)
	engine.State.Players[0] = NewPlayerState(0, "P1", &model.Deck{})
	engine.State.Players[1] = NewPlayerState(1, "P2", &model.Deck{})
	engine.State.Phase = PhaseMain
	engine.State.CurrentTurn = 0
	engine.State.TurnNumber = 1
	return engine
}

func baseCard(t *testing.T, id string) *model.Card {
	t.Helper()
	card, ok := cards.CardDB[id]
	if !ok {
		t.Fatalf("missing card %s", id)
	}
	return card
}

func placeUnit(card *model.Card, ownerID int, col int, row int, engine *Engine) *CardInstance {
	unit := NewCardInstance(card, ownerID, engine.State.TurnNumber)
	unit.IsHorizontal = false
	unit.Position = &Position{Col: col, Row: row}
	engine.State.Players[ownerID].Units[col][row] = unit
	return unit
}

func TestBoundSkillAttachesToHostInsteadOfSkillPool(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	ailaya := placeUnit(baseCard(t, "1311002"), 0, 1, 1, engine)

	engine.triggerEffects(TriggerOnEnter, ailaya, nil, nil)

	if len(p0.SkillPool) != 0 {
		t.Fatalf("bound skill should not be added to skill pool, got %d", len(p0.SkillPool))
	}
	if len(ailaya.BoundSkills) != 1 || ailaya.BoundSkills[0].Card.Number != "3301001" {
		t.Fatalf("expected Storm Fury bound to Ailaya, got %+v", ailaya.BoundSkills)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "learn_skill", Data: map[string]any{
		"instance_id": ailaya.BoundSkills[0].InstanceID,
	}}); err == nil {
		t.Fatalf("bound skill should not be learnable through the skill pool")
	}

	info := cardToInfo(ailaya)
	bound, ok := info["bound_skills"].([]map[string]any)
	if !ok || len(bound) != 1 || bound[0]["number"] != "3301001" {
		t.Fatalf("card info should expose bound skill for display, got %+v", info["bound_skills"])
	}

	engine.destroyUnit(ailaya, 0)
	if len(p0.Graveyard) != 1 || len(p0.Graveyard[0].BoundSkills) != 0 {
		t.Fatalf("bound skills should disappear when host leaves battlefield, graveyard=%+v", p0.Graveyard)
	}
}

func TestCycloneWavePaysTwoAirUsingArcaneWildcard(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Elements[model.ElementAir] = 1
	p0.Elements[model.ElementArcane] = 1
	placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast cyclone wave: %v", err)
	}
	if p0.Elements[model.ElementAir] != 0 || p0.Elements[model.ElementArcane] != 0 {
		t.Fatalf("expected air and arcane to be spent, got %v", p0.Elements)
	}
	if engine.State.Phase != PhaseDefenseWindow || p1.Units[1][0] == nil {
		t.Fatalf("expected normal spell defense window")
	}
}

func TestArcaneArrowDealsOneDamageAsSorcery(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	target.CurrentLife = 3
	p0.Skills[0] = readySkill(baseCard(t, "3021005"), 0)
	p0.Elements[model.ElementArcane] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast arcane arrow: %v", err)
	}
	if target.CurrentLife != 2 {
		t.Fatalf("arcane arrow should deal 1 damage, target life=%d", target.CurrentLife)
	}
	if engine.State.Phase != PhaseMain {
		t.Fatalf("sorcery should resolve immediately, got phase %v", engine.State.Phase)
	}
}

func TestSleepAppliesVisibleStunStatusOnHit(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3021009"), 0)
	p0.Elements[model.ElementArcane] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast sleep: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve sleep hit: %v", err)
	}
	if target.Statuses[StatusStun] != 1 {
		t.Fatalf("sleep should apply stun 1, statuses=%v", target.Statuses)
	}
}

func TestMeditationCanCastWithoutTargetAndGainArcane(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p0.Skills[0] = readySkill(baseCard(t, "3021003"), 0)
	p0.Elements[model.ElementAir] = 1
	p0.Elements[model.ElementFire] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "none",
	}}); err != nil {
		t.Fatalf("cast meditation without target: %v", err)
	}
	if p0.Elements[model.ElementArcane] != 1 {
		t.Fatalf("meditation should gain 1 arcane, elements=%v", p0.Elements)
	}
	if p0.Elements[model.ElementAir]+p0.Elements[model.ElementFire] != 0 {
		t.Fatalf("meditation should spend two non-arcane elements for arcane cost, elements=%v", p0.Elements)
	}
}

func TestSquareSpellAppliesDamageAndStatusToAllEnemyUnits(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	front := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	back := placeUnit(baseCard(t, "1021001"), 1, 0, 2, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3121005"), 0)
	p0.Elements[model.ElementFire] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast square firestorm: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve square firestorm: %v", err)
	}
	if front.CurrentLife != front.Card.Life-1 || back.CurrentLife != back.Card.Life-1 {
		t.Fatalf("square spell should damage all enemy units, front=%d back=%d", front.CurrentLife, back.CurrentLife)
	}
	if front.Statuses[StatusBurn] != 1 || back.Statuses[StatusBurn] != 1 {
		t.Fatalf("square spell should burn all enemy units, front=%v back=%v", front.Statuses, back.Statuses)
	}
}

func TestGenericFriendlySpellPowerAndAttackBonuses(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	placeUnit(baseCard(t, "1121004"), 0, 0, 0, engine) // 你的所有法术+1\威
	placeUnit(baseCard(t, "1321006"), 0, 2, 0, engine) // 你的大气法术+1\攻
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Elements[model.ElementAir] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast boosted cyclone wave: %v", err)
	}
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 2 {
		t.Fatalf("expected cyclone power 2 after generic +威, pending=%+v", engine.State.PendingSpell)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve boosted cyclone wave: %v", err)
	}
	if target.CurrentLife != target.Card.Life-2 {
		t.Fatalf("expected cyclone to deal 2 after generic +攻, life=%d", target.CurrentLife)
	}
}

func TestRavenMessengerCanTapToDrawWithoutGainingElements(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	raven := placeUnit(baseCard(t, "1321001"), 0, 0, 0, engine)
	p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 0)}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  raven.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use raven ability: %v", err)
	}
	if !raven.IsHorizontal {
		t.Fatalf("raven should tap")
	}
	if len(p0.Hand) != 1 {
		t.Fatalf("raven should draw 1 card, hand=%d", len(p0.Hand))
	}
	if p0.Elements[model.ElementAir] != 0 {
		t.Fatalf("raven draw option should not gain air, elements=%v", p0.Elements)
	}
}

func TestStormChimeraReducesAirSpellUseCost(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	placeUnit(baseCard(t, "1321010"), 0, 0, 0, engine)
	placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Elements[model.ElementAir] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("chimera should reduce cyclone wave to 1 air: %v", err)
	}
	if p0.Elements[model.ElementAir] != 0 {
		t.Fatalf("expected reduced cost to spend 1 air, elements=%v", p0.Elements)
	}
}

func TestStormChimeraDevoursFriendlyCompanion(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	chimera := NewCardInstance(baseCard(t, "1321010"), 0, engine.State.TurnNumber)
	p0.Hand = append(p0.Hand, chimera)
	food := placeUnit(baseCard(t, "1321014"), 0, 1, 0, engine)
	p0.Elements[model.ElementAir] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": chimera.InstanceID,
		"devour_id":   food.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon chimera with devour: %v", err)
	}
	if p0.Units[1][0] != nil || len(p0.Graveyard) != 1 {
		t.Fatalf("devoured unit should be destroyed")
	}
	if p0.Units[0][0] == nil || p0.Units[0][0].Card.Number != "1321010" {
		t.Fatalf("chimera should enter after devour")
	}
	if p0.Elements[model.ElementAir] != 0 {
		t.Fatalf("summon should spend 3 air, elements=%v", p0.Elements)
	}
}

func TestStormChimeraRequiresDevourBeforeSummon(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	chimera := NewCardInstance(baseCard(t, "1321010"), 0, engine.State.TurnNumber)
	p0.Hand = append(p0.Hand, chimera)
	placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
	p0.Elements[model.ElementAir] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": chimera.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err == nil {
		t.Fatalf("chimera should require a valid devour target before summon")
	}
}

func TestStaticPulseRequiresEnemyTarget(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p0.Skills[0] = readySkill(baseCard(t, "3321003"), 0)
	p0.Elements[model.ElementAir] = 1
	placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "none",
	}}); err == nil {
		t.Fatalf("static pulse has power and stun text, so it should require an enemy target")
	}
}

func TestHasteSkillEntersReadyWhenLearned(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	sleep := NewCardInstance(baseCard(t, "3021009"), 0, engine.State.TurnNumber)
	p0.SkillPool = append(p0.SkillPool, sleep)
	p0.Elements[model.ElementArcane] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "learn_skill", Data: map[string]any{
		"instance_id": sleep.InstanceID,
	}}); err != nil {
		t.Fatalf("learn sleep: %v", err)
	}
	if p0.Skills[0] == nil || p0.Skills[0].IsHorizontal {
		t.Fatalf("haste skill should enter ready")
	}
}

func TestCardsResetAtEndOfOwnersTurn(t *testing.T) {
	engine := setupReportedBugEngine(t)
	unit := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	unit.IsHorizontal = true

	if err := engine.HandleAction(0, ActionMessage{Action: "end_turn"}); err != nil {
		t.Fatalf("end turn: %v", err)
	}
	if unit.IsHorizontal {
		t.Fatalf("owner's cards should reset at end of turn")
	}
	if engine.State.CurrentTurn != 1 {
		t.Fatalf("turn should pass to player 2")
	}
}

func TestMulingUltimateReturnsOneCompanionFromEachSide(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	muling := placeUnit(baseCard(t, "4311003"), 0, 1, 1, engine)
	own := placeUnit(baseCard(t, "1021001"), 0, 0, 0, engine)
	enemy := placeUnit(baseCard(t, "1321010"), 1, 1, 0, engine)
	p0.Elements[model.ElementAir] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  muling.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use muling ultimate: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected muling selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{own.InstanceID, enemy.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve muling: %v", err)
	}
	if p0.Units[0][0] != nil || engine.State.Players[1].Units[1][0] != nil {
		t.Fatalf("selected companions should leave the battlefield")
	}
	if len(p0.Hand) != 1 || len(engine.State.Players[1].Hand) != 1 {
		t.Fatalf("selected companions should return to owners' hands")
	}
	if p0.Elements[model.ElementAir] != 0 {
		t.Fatalf("muling should spend cost difference 2 air, elements=%v", p0.Elements)
	}
}

func TestBlackMarketVendorDiscardsItemAndDrawsTwo(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	vendor := placeUnit(baseCard(t, "1021012"), 0, 0, 0, engine)
	item := NewCardInstance(baseCard(t, "2021003"), 0, 1)
	p0.Hand = []*CardInstance{item}
	p0.Deck = []*CardInstance{
		NewCardInstance(baseCard(t, "1021001"), 0, 1),
		NewCardInstance(baseCard(t, "1021002"), 0, 1),
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  vendor.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use black market vendor: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected discard selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{item.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve black market vendor: %v", err)
	}
	if len(p0.Graveyard) != 1 || p0.Graveyard[0].InstanceID != item.InstanceID {
		t.Fatalf("selected item should be discarded, graveyard=%v", len(p0.Graveyard))
	}
	if len(p0.Hand) != 2 {
		t.Fatalf("vendor should draw two after discarding, hand=%d", len(p0.Hand))
	}
}

func TestFireArtistResetsAnotherFireCard(t *testing.T) {
	engine := setupReportedBugEngine(t)
	artist := placeUnit(baseCard(t, "1121010"), 0, 0, 0, engine)
	target := placeUnit(baseCard(t, "1121002"), 0, 1, 0, engine)
	target.IsHorizontal = true

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  artist.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use fire artist: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected reset selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{target.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve fire artist: %v", err)
	}
	if !artist.IsHorizontal {
		t.Fatalf("fire artist should tap itself")
	}
	if target.IsHorizontal {
		t.Fatalf("target fire card should be reset")
	}
}

func TestWindriderDiscardsAnyNumberForAir(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	windrider := placeUnit(baseCard(t, "1321005"), 0, 0, 0, engine)
	cardA := NewCardInstance(baseCard(t, "1021001"), 0, 1)
	cardB := NewCardInstance(baseCard(t, "1021002"), 0, 1)
	cardC := NewCardInstance(baseCard(t, "1021003"), 0, 1)
	p0.Hand = []*CardInstance{cardA, cardB, cardC}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  windrider.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use windrider: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{cardA.InstanceID, cardB.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve windrider: %v", err)
	}
	if p0.Elements[model.ElementAir] != 2 {
		t.Fatalf("windrider should gain 1 air per discarded card, elements=%v", p0.Elements)
	}
	if len(p0.Hand) != 1 || len(p0.Graveyard) != 2 {
		t.Fatalf("windrider should discard exactly two selected cards, hand=%d graveyard=%d", len(p0.Hand), len(p0.Graveyard))
	}
}

func TestSunwheelMageResetsLightSkill(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	mage := placeUnit(baseCard(t, "1521011"), 0, 0, 0, engine)
	lightSkill := readySkill(baseCard(t, "3521001"), 0)
	lightSkill.IsHorizontal = true
	p0.Skills[0] = lightSkill

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  mage.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use sunwheel mage: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{lightSkill.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve sunwheel mage: %v", err)
	}
	if lightSkill.IsHorizontal {
		t.Fatalf("light skill should be reset")
	}
}

func TestSoulPriestSacrificesCompanionAndDrawsTwo(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	priest := placeUnit(baseCard(t, "1621012"), 0, 0, 0, engine)
	food := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
	p0.Deck = []*CardInstance{
		NewCardInstance(baseCard(t, "1021002"), 0, 1),
		NewCardInstance(baseCard(t, "1021003"), 0, 1),
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  priest.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use soul priest: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{food.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve soul priest: %v", err)
	}
	if p0.Units[1][0] != nil {
		t.Fatalf("sacrificed companion should leave battlefield")
	}
	if len(p0.Hand) != 2 {
		t.Fatalf("soul priest should draw two, hand=%d", len(p0.Hand))
	}
}

func TestWhimWandConsumesToResetLowCostSkillFromEquipmentSlot(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	wand := NewCardInstance(baseCard(t, "2021003"), 0, 1)
	wand.IsHorizontal = false
	p0.Equipment[0] = wand
	lowCost := readySkill(baseCard(t, "3021005"), 0)
	lowCost.IsHorizontal = true
	p0.Skills[0] = lowCost

	if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{
		"instance_id": wand.InstanceID,
	}}); err != nil {
		t.Fatalf("consume whim wand: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected low-cost skill selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{lowCost.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve whim wand: %v", err)
	}
	if !wand.IsHorizontal {
		t.Fatalf("whim wand should tap")
	}
	if lowCost.IsHorizontal {
		t.Fatalf("selected low-cost skill should be reset")
	}
}

func TestFireSpriteGainsBurnWhenConsumed(t *testing.T) {
	engine := setupReportedBugEngine(t)
	sprite := placeUnit(baseCard(t, "1121001"), 0, 0, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{
		"instance_id": sprite.InstanceID,
	}}); err != nil {
		t.Fatalf("consume fire sprite: %v", err)
	}
	if sprite.Statuses[StatusBurn] != 1 {
		t.Fatalf("fire sprite should gain burn 1 when consumed, statuses=%v", sprite.Statuses)
	}
}

func TestBottledElementGainsOneArcaneWhenUsed(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	item := NewCardInstance(baseCard(t, "2021005"), 0, 1)
	p0.Hand = []*CardInstance{item}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": item.InstanceID,
	}}); err != nil {
		t.Fatalf("use bottled element: %v", err)
	}
	if p0.Elements[model.ElementArcane] != 1 {
		t.Fatalf("bottled element should gain 1 arcane, elements=%v", p0.Elements)
	}
	if len(p0.Hand) != 0 || len(p0.Graveyard) != 1 {
		t.Fatalf("used consumable should move from hand to graveyard, hand=%d graveyard=%d", len(p0.Hand), len(p0.Graveyard))
	}
}

func TestLifePotionUsesPendingSelectionToHealFriendlyUnit(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	item := NewCardInstance(baseCard(t, "2521001"), 0, 1)
	p0.Hand = []*CardInstance{item}
	p0.Elements[model.ElementArcane] = 1
	target := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
	target.CurrentLife = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": item.InstanceID,
	}}); err != nil {
		t.Fatalf("use life potion: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected life potion target selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{target.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve life potion: %v", err)
	}
	if target.CurrentLife != 3 {
		t.Fatalf("life potion should heal up to max life, life=%d", target.CurrentLife)
	}
	if p0.Elements[model.ElementArcane] != 0 {
		t.Fatalf("life potion should spend one arcane wildcard for neutral cost, elements=%v", p0.Elements)
	}
}

func TestWindcallingScrollDrawsTwoAndSkipsNextDraw(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	item := NewCardInstance(baseCard(t, "2321005"), 0, 1)
	p0.Hand = []*CardInstance{item}
	p0.Deck = []*CardInstance{
		NewCardInstance(baseCard(t, "1021001"), 0, 1),
		NewCardInstance(baseCard(t, "1021002"), 0, 1),
		NewCardInstance(baseCard(t, "1021003"), 0, 1),
	}
	p0.Elements[model.ElementAir] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": item.InstanceID,
	}}); err != nil {
		t.Fatalf("use windcalling scroll: %v", err)
	}
	if len(p0.Hand) != 2 || len(p0.Deck) != 1 {
		t.Fatalf("windcalling scroll should draw two, hand=%d deck=%d", len(p0.Hand), len(p0.Deck))
	}
	if !p0.SkipNextDraw {
		t.Fatalf("windcalling scroll should mark next turn draw to be skipped")
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
		t.Fatalf("end p0 turn: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
		t.Fatalf("end p1 turn: %v", err)
	}
	if p0.SkipNextDraw {
		t.Fatalf("skip draw flag should clear on next own turn")
	}
	if len(p0.Deck) != 1 {
		t.Fatalf("next normal draw should be skipped, deck=%d", len(p0.Deck))
	}
}

func TestArcaneArmorerSearchesEquipmentWhenNoEquipment(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	armorer := NewCardInstance(baseCard(t, "1021016"), 0, 1)
	equipment := NewCardInstance(baseCard(t, "2021003"), 0, 1)
	p0.Hand = []*CardInstance{armorer}
	p0.Deck = []*CardInstance{equipment, NewCardInstance(baseCard(t, "1021001"), 0, 1)}
	p0.Elements[model.ElementArcane] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": armorer.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon arcane armorer: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected equipment search selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{equipment.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve armorer search: %v", err)
	}
	if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != equipment.InstanceID {
		t.Fatalf("searched equipment should be added to hand, hand=%v", len(p0.Hand))
	}
}

func TestRunemasterDiscardsHandCardThenSearchesRuneOrScroll(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	runemaster := NewCardInstance(baseCard(t, "1021017"), 0, 1)
	discard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
	scroll := NewCardInstance(baseCard(t, "2321005"), 0, 1)
	p0.Hand = []*CardInstance{runemaster, discard}
	p0.Deck = []*CardInstance{scroll, NewCardInstance(baseCard(t, "1021002"), 0, 1)}
	p0.Elements[model.ElementArcane] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": runemaster.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon runemaster: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected discard selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{discard.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve runemaster discard: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected search selection after discard, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{scroll.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve runemaster search: %v", err)
	}
	if len(p0.Graveyard) != 1 || p0.Graveyard[0].InstanceID != discard.InstanceID {
		t.Fatalf("runemaster should discard selected hand card, graveyard=%d", len(p0.Graveyard))
	}
	if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != scroll.InstanceID {
		t.Fatalf("runemaster should add selected rune or scroll to hand, hand=%d", len(p0.Hand))
	}
}

func TestRedHawkSearchesHighCostFireCompanion(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	hawk := NewCardInstance(baseCard(t, "1121009"), 0, 1)
	target := NewCardInstance(baseCard(t, "1121004"), 0, 1)
	p0.Hand = []*CardInstance{hawk}
	p0.Deck = []*CardInstance{target, NewCardInstance(baseCard(t, "1121002"), 0, 1)}
	p0.Elements[model.ElementFire] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": hawk.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon red hawk: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected fire companion search selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{target.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve red hawk search: %v", err)
	}
	if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != target.InstanceID {
		t.Fatalf("red hawk should search selected high-cost fire companion, hand=%d", len(p0.Hand))
	}
}

func TestPegasusKnightSearchesUnicornPegasus(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	knight := NewCardInstance(baseCard(t, "1521009"), 0, 1)
	pegasus := NewCardInstance(baseCard(t, "1521012"), 0, 1)
	p0.Hand = []*CardInstance{knight}
	p0.Deck = []*CardInstance{pegasus, NewCardInstance(baseCard(t, "1021001"), 0, 1)}
	p0.Elements[model.ElementLight] = 2

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": knight.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon pegasus knight: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected pegasus search selection, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{pegasus.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve pegasus knight search: %v", err)
	}
	if len(p0.Hand) != 1 || p0.Hand[0].InstanceID != pegasus.InstanceID {
		t.Fatalf("pegasus knight should search unicorn pegasus, hand=%d", len(p0.Hand))
	}
}

func TestManaBoosterMakesNextSkillUseFree(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	item := NewCardInstance(baseCard(t, "2021014"), 0, 1)
	p0.Hand = []*CardInstance{item}
	p0.Elements[model.ElementArcane] = 1
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": item.InstanceID,
	}}); err != nil {
		t.Fatalf("use mana booster: %v", err)
	}
	if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextSkillCostZero {
		t.Fatalf("mana booster should add next skill cost modifier, modifiers=%v", p0.TempModifiers)
	}
	if p0.Elements[model.ElementArcane] != 0 {
		t.Fatalf("mana booster item cost should be paid, elements=%v", p0.Elements)
	}

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("next skill should be free: %v", err)
	}
	if len(p0.TempModifiers) != 0 {
		t.Fatalf("free next skill modifier should be consumed, modifiers=%v", p0.TempModifiers)
	}
}

func TestStoneforgeArtisanGivesSelectedSpellPowerThisTurn(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	artisan := placeUnit(baseCard(t, "1121003"), 0, 0, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Elements[model.ElementAir] = 2
	placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "consume", Data: map[string]any{
		"instance_id": artisan.InstanceID,
	}}); err != nil {
		t.Fatalf("consume stoneforge artisan: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction {
		t.Fatalf("expected spell selection for power bonus, phase=%v", engine.State.Phase)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{p0.Skills[0].InstanceID},
	}}); err != nil {
		t.Fatalf("resolve stoneforge artisan: %v", err)
	}
	if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModSkillPowerBonus {
		t.Fatalf("stoneforge should add spell power modifier, modifiers=%v", p0.TempModifiers)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast boosted spell: %v", err)
	}
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 3 {
		t.Fatalf("cyclone wave should have base 1 +2 power, pending=%+v", engine.State.PendingSpell)
	}
}

func TestEnergeticSeniorMakesNextSkillIgnoreCooldown(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	senior := NewCardInstance(baseCard(t, "1021015"), 0, 1)
	p0.Hand = []*CardInstance{senior}
	p0.Elements[model.ElementArcane] = 2
	p0.Elements[model.ElementLight] = 1
	p0.Skills[0] = readySkill(baseCard(t, "3021008"), 0)
	placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": senior.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon energetic senior: %v", err)
	}
	if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextNoCooldown {
		t.Fatalf("energetic senior should add no-cooldown modifier, modifiers=%v", p0.TempModifiers)
	}
	p0.Elements[model.ElementArcane] = 2
	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast cooldown spell: %v", err)
	}
	if p0.Skills[0].Statuses[StatusCooldown] > 0 {
		t.Fatalf("next skill should ignore cooldown, statuses=%v", p0.Skills[0].Statuses)
	}
	if len(p0.TempModifiers) != 0 {
		t.Fatalf("no-cooldown modifier should be consumed, modifiers=%v", p0.TempModifiers)
	}
}

func TestWinterfellWarlockGivesNextSpellFreezeOnHit(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
	warlock := placeUnit(baseCard(t, "1221011"), 0, 0, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
	p0.Elements[model.ElementFire] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  warlock.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use winterfell warlock ultimate: %v", err)
	}
	if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModNextSpellHitStatus {
		t.Fatalf("warlock should add spell-hit status modifier, modifiers=%v", p0.TempModifiers)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast fireball after warlock ultimate: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve fireball hit: %v", err)
	}
	if target.Statuses[StatusFreeze] != 1 {
		t.Fatalf("next spell should apply freeze 1, statuses=%v", target.Statuses)
	}
	if len(p0.TempModifiers) != 0 {
		t.Fatalf("spell-hit status modifier should be consumed, modifiers=%v", p0.TempModifiers)
	}
}

func TestPassionOfFireDrawsWhenFriendlyFireSpellHits(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p0.Skills[0] = readySkill(baseCard(t, "3121007"), 0)
	p0.Skills[1] = readySkill(baseCard(t, "3121001"), 0)
	p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
	p0.Elements[model.ElementFire] = 1
	placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[1].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast fireball with passion of fire in play: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve fireball hit: %v", err)
	}
	if len(p0.Hand) != 1 {
		t.Fatalf("passion of fire should draw 1 after friendly fire spell hit, hand=%d deck=%d", len(p0.Hand), len(p0.Deck))
	}
}

func TestSpellStatPassivePowerAndDamageModifiers(t *testing.T) {
	t.Run("celtic lion gives all spells plus one power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1121004"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementFire] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 4 {
			t.Fatalf("fireball should have base 3 +1 power, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("raincaller gives water and air spells plus one power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1221013"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
		p0.Elements[model.ElementAir] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast cyclone wave: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 2 {
			t.Fatalf("cyclone wave should have base 1 +1 power, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("thunder beast gives air spells plus one damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		target.CurrentLife = 3
		placeUnit(baseCard(t, "1321006"), 0, 0, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
		p0.Elements[model.ElementAir] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast cyclone wave: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve cyclone wave: %v", err)
		}
		if target.CurrentLife != 1 {
			t.Fatalf("cyclone wave should deal base 1 +1 damage, target life=%d", target.CurrentLife)
		}
	})

	t.Run("divine fire beast gives attacking spells plus two power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1521013"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementFire] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 5 {
			t.Fatalf("fireball should have base 3 +2 attacking power, pending=%+v", engine.State.PendingSpell)
		}
	})
}

func TestUnitEnterListenersCanReactToEnemySummons(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	placeUnit(baseCard(t, "1111002"), 0, 0, 0, engine)
	summoned := NewCardInstance(baseCard(t, "1021001"), 1, 1)
	p1.Hand = []*CardInstance{summoned}
	p1.Elements[model.ElementArcane] = 10
	engine.State.CurrentTurn = 1

	if err := engine.HandleAction(1, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": summoned.InstanceID,
		"col":         float64(1),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon enemy unit: %v", err)
	}
	if summoned.Statuses[StatusBurn] != 1 || summoned.Statuses[StatusPetrify] != 1 {
		t.Fatalf("inferno general should burn and petrify enemy summons, statuses=%v", summoned.Statuses)
	}
	if p0.Hero != nil && p0.Hero.CurrentLife <= 0 {
		t.Fatalf("friendly-only unit-enter listeners should not misfire from enemy summon")
	}
}

func TestTwinAngelCreatesTwinInHandOnEnter(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	angel := NewCardInstance(baseCard(t, "1521005"), 0, 1)
	p0.Hand = []*CardInstance{angel}
	p0.Elements[model.ElementLight] = 10

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": angel.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon twin angel: %v", err)
	}
	if len(p0.Hand) != 1 || p0.Hand[0].Card.Number != "1501001" {
		t.Fatalf("twin angel should create card 1501001 in hand, hand=%v", cardsToInfo(p0.Hand))
	}
}

func TestRecyclingSpriteMovesSelectedGraveyardCardToDeckTop(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	sprite := NewCardInstance(baseCard(t, "1021007"), 0, 1)
	recovered := NewCardInstance(baseCard(t, "1021001"), 0, 1)
	otherDeckCard := NewCardInstance(baseCard(t, "1021002"), 0, 1)
	p0.Hand = []*CardInstance{sprite}
	p0.Graveyard = []*CardInstance{recovered}
	p0.Deck = []*CardInstance{otherDeckCard}
	p0.Elements[model.ElementArcane] = 10

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": sprite.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon recycling sprite: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction || engine.State.PendingAction == nil {
		t.Fatalf("recycling sprite should ask which graveyard card to recover")
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{recovered.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve recycling sprite action: %v", err)
	}
	if len(p0.Graveyard) != 0 {
		t.Fatalf("graveyard card should be removed, graveyard=%v", cardsToInfo(p0.Graveyard))
	}
	if len(p0.Deck) == 0 || p0.Deck[0] != recovered {
		t.Fatalf("selected graveyard card should be deck top, deck=%v", cardsToInfo(p0.Deck))
	}
}

func TestLifeFlowerBuffsAnotherFriendlyUnit(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021004"), 0, 1, 0, engine)
	target.CurrentLife = 2
	flower := NewCardInstance(baseCard(t, "1521006"), 0, 1)
	p0.Hand = []*CardInstance{flower}
	p0.Elements[model.ElementLight] = 10

	if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
		"instance_id": flower.InstanceID,
		"col":         float64(0),
		"row":         float64(0),
	}}); err != nil {
		t.Fatalf("summon life flower: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{target.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve life flower: %v", err)
	}
	if target.CurrentLife != 3 {
		t.Fatalf("life flower should give +1 life, life=%d", target.CurrentLife)
	}
}

func TestDarkDeathEffects(t *testing.T) {
	t.Run("elemental husk gains none element on death", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		husk := placeUnit(baseCard(t, "1621002"), 0, 0, 0, engine)

		engine.dealDamage(husk, 99, 0)
		if p0.Elements[model.ElementArcane] != 1 {
			t.Fatalf("elemental husk should gain 1 none element, elements=%v", p0.Elements)
		}
	})

	t.Run("nightmare gains life when another friendly unit dies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		nightmare := placeUnit(baseCard(t, "1621006"), 0, 0, 0, engine)
		other := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		startLife := nightmare.CurrentLife

		engine.dealDamage(other, 99, 0)
		if nightmare.CurrentLife != startLife+1 {
			t.Fatalf("nightmare should gain +1 life, life=%d start=%d", nightmare.CurrentLife, startLife)
		}
	})

	t.Run("bone knight resummons once and loses deathrattle", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		knight := placeUnit(baseCard(t, "1621011"), 0, 0, 0, engine)

		engine.dealDamage(knight, 99, 0)
		if engine.State.Players[0].Units[0][0] != knight {
			t.Fatalf("bone knight should return to its position")
		}
		if knight.Statuses[boneKnightRebornStatus] != 1 {
			t.Fatalf("bone knight should lose deathrattle after return, statuses=%v", knight.Statuses)
		}
		engine.dealDamage(knight, 99, 0)
		if engine.State.Players[0].Units[0][0] != nil {
			t.Fatalf("bone knight should not return a second time")
		}
	})
}

func TestOnSpellCastListenersFireForFriendlySkills(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	placeUnit(baseCard(t, "1321012"), 0, 0, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
	p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
	p0.Elements[model.ElementAir] = 2
	placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast air spell: %v", err)
	}
	if len(p0.Hand) != 1 {
		t.Fatalf("wind medium should draw when friendly air skill is used, hand=%d deck=%d", len(p0.Hand), len(p0.Deck))
	}
}

func TestEquipmentEnterEffectsAndSpellStatPassives(t *testing.T) {
	t.Run("equipment on-enter triggers when equipped", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
		ring := NewCardInstance(baseCard(t, "2321007"), 0, 1)
		p0.Hand = []*CardInstance{ring}
		p0.Elements[model.ElementAir] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
			"instance_id": ring.InstanceID,
		}}); err != nil {
			t.Fatalf("equip windwhisper ring: %v", err)
		}
		if len(p0.Hand) != 1 {
			t.Fatalf("ring on-enter should draw 1, hand=%d deck=%d", len(p0.Hand), len(p0.Deck))
		}
	})

	t.Run("wizard scepter gives spell plus one power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2021004"), 0, 1)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementFire] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 4 {
			t.Fatalf("wizard scepter should give +1 power, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("manes staff only buffs water spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2221004"), 0, 1)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3221001"), 0)
		p0.Elements[model.ElementWater] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast hail: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 2 {
			t.Fatalf("manes staff should give water spell +1 power, pending=%+v", engine.State.PendingSpell)
		}
	})
}

func TestEquipmentLifeAndEnemySummonCounters(t *testing.T) {
	t.Run("life amulet buffs selected friendly role", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
		amulet := NewCardInstance(baseCard(t, "2021011"), 0, 1)
		p0.Hand = []*CardInstance{amulet}
		p0.Elements[model.ElementArcane] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
			"instance_id": amulet.InstanceID,
		}}); err != nil {
			t.Fatalf("equip life amulet: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve life amulet: %v", err)
		}
		if target.CurrentLife != target.Card.Life+1 {
			t.Fatalf("life amulet should give +1 life, life=%d", target.CurrentLife)
		}
	})

	t.Run("rattan cuirass buffs hero on enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Hero = placeUnit(baseCard(t, "4311003"), 0, 1, 1, engine)
		startLife := p0.Hero.CurrentLife
		cuirass := NewCardInstance(baseCard(t, "2421006"), 0, 1)
		p0.Hand = []*CardInstance{cuirass}
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "equip", Data: map[string]any{
			"instance_id": cuirass.InstanceID,
		}}); err != nil {
			t.Fatalf("equip rattan cuirass: %v", err)
		}
		if p0.Hero.CurrentLife != startLife+2 {
			t.Fatalf("rattan cuirass should give hero +2 life, life=%d start=%d", p0.Hero.CurrentLife, startLife)
		}
	})

	t.Run("hellfire rune statuses enemy summons", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2121012"), 0, 1)
		summoned := NewCardInstance(baseCard(t, "1021001"), 1, 1)
		p1.Hand = []*CardInstance{summoned}
		p1.Elements[model.ElementArcane] = 10
		engine.State.CurrentTurn = 1

		if err := engine.HandleAction(1, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": summoned.InstanceID,
			"col":         float64(0),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("enemy summon: %v", err)
		}
		if summoned.Statuses[StatusStun] != 2 || summoned.Statuses[StatusPetrify] != 2 || summoned.Statuses[StatusBurn] != 2 {
			t.Fatalf("hellfire rune should apply three statuses, statuses=%v", summoned.Statuses)
		}
	})

	t.Run("killing instinct damages enemy summons", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p1 := engine.State.Players[1]
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2621003"), 0, 1)
		summoned := NewCardInstance(baseCard(t, "1021004"), 1, 1)
		p1.Hand = []*CardInstance{summoned}
		p1.Elements[model.ElementArcane] = 10
		engine.State.CurrentTurn = 1

		if err := engine.HandleAction(1, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": summoned.InstanceID,
			"col":         float64(0),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("enemy summon: %v", err)
		}
		if summoned.CurrentLife != summoned.Card.Life-2 {
			t.Fatalf("killing instinct should deal 2 damage, life=%d", summoned.CurrentLife)
		}
	})
}

func TestConsumableTargetedItemEffects(t *testing.T) {
	t.Run("fire arrow damages selected enemy", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		item := NewCardInstance(baseCard(t, "2121004"), 0, 1)
		p0.Hand = []*CardInstance{item}
		p0.Elements[model.ElementFire] = 10
		target := placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": item.InstanceID,
		}}); err != nil {
			t.Fatalf("use fire arrow item: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve fire arrow: %v", err)
		}
		if target.CurrentLife != target.Card.Life-1 {
			t.Fatalf("fire arrow should deal 1 damage, life=%d", target.CurrentLife)
		}
	})

	t.Run("bottled lightning gains air and stuns friendly", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		item := NewCardInstance(baseCard(t, "2321006"), 0, 1)
		p0.Hand = []*CardInstance{item}
		p0.Elements[model.ElementAir] = 10
		target := placeUnit(baseCard(t, "1021004"), 0, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": item.InstanceID,
		}}); err != nil {
			t.Fatalf("use bottled lightning: %v", err)
		}
		if p0.Elements[model.ElementAir] != 12 {
			t.Fatalf("bottled lightning should net +3 air after paying 1, elements=%v", p0.Elements)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve bottled lightning: %v", err)
		}
		if target.Statuses[StatusStun] != 2 {
			t.Fatalf("bottled lightning should stun friendly 2, statuses=%v", target.Statuses)
		}
	})

	t.Run("purification scroll removes friendly negative statuses", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		item := NewCardInstance(baseCard(t, "2521003"), 0, 1)
		p0.Hand = []*CardInstance{item}
		p0.Elements[model.ElementLight] = 10
		target := placeUnit(baseCard(t, "1021004"), 0, 1, 0, engine)
		target.Statuses[StatusBurn] = 1
		target.Statuses[StatusFreeze] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": item.InstanceID,
		}}); err != nil {
			t.Fatalf("use purification scroll: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve purification scroll: %v", err)
		}
		if target.Statuses[StatusBurn] != 0 || target.Statuses[StatusFreeze] != 0 {
			t.Fatalf("purification scroll should remove negative statuses, statuses=%v", target.Statuses)
		}
	})
}

func TestSkillPendingChoiceEffects(t *testing.T) {
	t.Run("engrave discards hand card and draws first rune or scroll", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		skill := readySkill(baseCard(t, "3021004"), 0)
		discard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		scroll := NewCardInstance(baseCard(t, "2321005"), 0, 1)
		p0.Skills[0] = skill
		p0.Hand = []*CardInstance{discard}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021002"), 0, 1), scroll}
		p0.Elements[model.ElementArcane] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast engrave: %v", err)
		}
		if engine.State.Phase != PhaseWaitingAction {
			t.Fatalf("engrave should wait for discard choice, phase=%v", engine.State.Phase)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{discard.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve engrave: %v", err)
		}
		if len(p0.Graveyard) != 1 || p0.Graveyard[0] != discard {
			t.Fatalf("engrave should discard selected card, graveyard=%v", cardsToInfo(p0.Graveyard))
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != scroll {
			t.Fatalf("engrave should draw first rune or scroll, hand=%v", cardsToInfo(p0.Hand))
		}
	})

	t.Run("healing spell heals selected friendly unit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
		target.CurrentLife = 1
		skill := readySkill(baseCard(t, "3521001"), 0)
		p0.Skills[0] = skill
		p0.Elements[model.ElementLight] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast healing: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve healing: %v", err)
		}
		if target.CurrentLife != 3 {
			t.Fatalf("healing should restore 2 life, life=%d", target.CurrentLife)
		}
	})

	t.Run("bloodsuck keeps pending choice after spell hit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		ally := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		skill := readySkill(baseCard(t, "3621002"), 0)
		p0.Skills[0] = skill
		p0.Elements[model.ElementShadow] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast bloodsuck: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve bloodsuck hit: %v", err)
		}
		if target.CurrentLife >= target.Card.Life {
			t.Fatalf("bloodsuck should still deal spell damage, target life=%d", target.CurrentLife)
		}
		if engine.State.Phase != PhaseWaitingAction {
			t.Fatalf("bloodsuck should keep pending choice after hit, phase=%v", engine.State.Phase)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{ally.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve bloodsuck buff: %v", err)
		}
		if ally.CurrentLife != ally.Card.Life+2 {
			t.Fatalf("bloodsuck should give friendly +2 life, life=%d", ally.CurrentLife)
		}
	})
}

func TestMoreSkillChoiceEffects(t *testing.T) {
	t.Run("disarm destroys selected enemy equipment after hit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
		equipment := NewCardInstance(baseCard(t, "2021004"), 1, 1)
		p1.Equipment[0] = equipment
		skill := readySkill(baseCard(t, "3021008"), 0)
		p0.Skills[0] = skill
		p0.Elements[model.ElementArcane] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast disarm: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve disarm hit: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{equipment.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve disarm equipment: %v", err)
		}
		if p1.Equipment[0] != nil || len(p1.Graveyard) != 1 {
			t.Fatalf("disarm should destroy selected equipment, equipment=%v graveyard=%v", p1.Equipment[0], cardsToInfo(p1.Graveyard))
		}
	})

	t.Run("blessing of light buffs friendly companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		ally := placeUnit(baseCard(t, "1021004"), 0, 0, 0, engine)
		skill := readySkill(baseCard(t, "3521014"), 0)
		p0.Skills[0] = skill
		p0.Elements[model.ElementLight] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast blessing of light: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{ally.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve blessing: %v", err)
		}
		if ally.CurrentLife != ally.Card.Life+1 || effectiveElementsGain(ally)[model.ElementLight] != 1 {
			t.Fatalf("blessing should add life and light load, life=%d gains=%v", ally.CurrentLife, effectiveElementsGain(ally))
		}
	})

	t.Run("soul recall returns up to two companions from graveyard", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		a := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		b := NewCardInstance(baseCard(t, "1021002"), 0, 1)
		p0.Graveyard = []*CardInstance{a, b}
		skill := readySkill(baseCard(t, "3621012"), 0)
		p0.Skills[0] = skill
		p0.Elements[model.ElementShadow] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": skill.InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast soul recall: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{a.InstanceID, b.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve soul recall: %v", err)
		}
		if len(p0.Hand) != 2 || len(p0.Graveyard) != 0 {
			t.Fatalf("soul recall should move two companions to hand, hand=%v graveyard=%v", cardsToInfo(p0.Hand), cardsToInfo(p0.Graveyard))
		}
	})
}

func TestSpellPowerItemModifiersAndRestrictions(t *testing.T) {
	t.Run("oath ring reduces attacking spell power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2021009"), 0, 1)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementFire] = 1
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 1 {
			t.Fatalf("oath ring should reduce fireball power from 3 to 1, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("severing blade boosts attack and prevents defense", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2021013"), 0, 1)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementFire] = 1
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 5 {
			t.Fatalf("severing blade should boost fireball power from 3 to 5, pending=%+v", engine.State.PendingSpell)
		}

		engine = setupReportedBugEngine(t)
		p1 = engine.State.Players[1]
		p1.Equipment[0] = NewCardInstance(baseCard(t, "2021013"), 1, 1)
		p1.Skills[0] = readySkill(baseCard(t, "2121009"), 1)
		p1.Elements[model.ElementFire] = 10
		engine.State.CurrentTurn = 0
		engine.State.Phase = PhaseDefenseWindow
		engine.State.PendingSpell = &SpellCast{AttackerID: 0, Skill: readySkill(baseCard(t, "3121001"), 0), Target: SpellTarget{Type: "unit", Position: Position{Col: 1, Row: 0}}, TotalPower: 1}

		err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{p1.Skills[0].InstanceID},
		}})
		if err == nil {
			t.Fatalf("severing blade should prevent using spells for defense")
		}
	})

	t.Run("divine flame potion gives temporary power and burns hero", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Hero = placeUnit(baseCard(t, "4311003"), 0, 1, 1, engine)
		item := NewCardInstance(baseCard(t, "2121005"), 0, 1)
		p0.Hand = []*CardInstance{item}
		p0.Elements[model.ElementFire] = 10
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": item.InstanceID,
		}}); err != nil {
			t.Fatalf("use divine flame potion: %v", err)
		}
		if p0.Hero.Statuses[StatusBurn] != 1 {
			t.Fatalf("divine flame potion should burn hero, statuses=%v", p0.Hero.Statuses)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast boosted fireball: %v", err)
		}
		if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 5 {
			t.Fatalf("divine flame potion should give +2 power, pending=%+v", engine.State.PendingSpell)
		}
	})
}

func TestConsumeAndDeathRuneEffects(t *testing.T) {
	t.Run("fire rune burns any consumed unit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2121002"), 0, 1)
		unit := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		engine.State.CurrentTurn = 1

		if err := engine.HandleAction(1, ActionMessage{Action: "consume", Data: map[string]any{
			"instance_id": unit.InstanceID,
		}}); err != nil {
			t.Fatalf("consume enemy unit with fire rune watching: %v", err)
		}
		if p1.Units[1][0].Statuses[StatusBurn] != 1 {
			t.Fatalf("fire rune should burn consumed unit, statuses=%v", p1.Units[1][0].Statuses)
		}
	})

	t.Run("frost rune freezes consumed enemy partner", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2221002"), 0, 1)
		unit := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		engine.State.CurrentTurn = 1

		if err := engine.HandleAction(1, ActionMessage{Action: "consume", Data: map[string]any{
			"instance_id": unit.InstanceID,
		}}); err != nil {
			t.Fatalf("consume enemy unit with frost rune watching: %v", err)
		}
		if unit.Statuses[StatusFreeze] != 1 {
			t.Fatalf("frost rune should freeze consumed enemy unit, statuses=%v", unit.Statuses)
		}
	})

	t.Run("lightning rune stuns consumed enemy and adjacent unit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2321002"), 0, 1)
		unit := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		adjacent := placeUnit(baseCard(t, "1021004"), 1, 0, 0, engine)
		engine.State.CurrentTurn = 1

		if err := engine.HandleAction(1, ActionMessage{Action: "consume", Data: map[string]any{
			"instance_id": unit.InstanceID,
		}}); err != nil {
			t.Fatalf("consume enemy unit with lightning rune watching: %v", err)
		}
		if unit.Statuses[StatusStun] != 1 || adjacent.Statuses[StatusStun] != 1 {
			t.Fatalf("lightning rune should stun consumed and adjacent units, unit=%v adjacent=%v", unit.Statuses, adjacent.Statuses)
		}
	})

	t.Run("sacrifice rune draws when friendly partner dies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Equipment[0] = NewCardInstance(baseCard(t, "2621005"), 0, 1)
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1021004"), 0, 1),
		}
		unit := placeUnit(baseCard(t, "1021007"), 0, 0, 0, engine)

		engine.destroyUnit(unit, 0)
		if len(p0.Hand) != 2 {
			t.Fatalf("sacrifice rune should draw 2, hand=%d deck=%d", len(p0.Hand), len(p0.Deck))
		}
	})
}

func TestPhoenixFeatherCountersAndPerTurnElement(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	feather := NewCardInstance(baseCard(t, "2121001"), 0, 1)
	p0.Equipment[0] = feather

	engine.triggerEffects(TriggerOnEnter, feather, nil, nil)
	if feather.Statuses[phoenixFeatherCounter] != 3 {
		t.Fatalf("phoenix feather should enter with 3 counters, statuses=%v", feather.Statuses)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  feather.InstanceID,
		"ability_type": "per_turn",
	}}); err != nil {
		t.Fatalf("use phoenix feather per-turn ability: %v", err)
	}
	if feather.Statuses[phoenixFeatherCounter] != 2 || p0.Elements[model.ElementFire] != 1 {
		t.Fatalf("phoenix feather should remove counter and gain fire, counters=%v elements=%v", feather.Statuses, p0.Elements)
	}
}

func TestCycloneScrollDestroysLowCostEnemyEquipment(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	scroll := NewCardInstance(baseCard(t, "2321008"), 0, 1)
	p0.Hand = []*CardInstance{scroll}
	p0.Elements[model.ElementAir] = 3
	p1.Equipment[0] = NewCardInstance(baseCard(t, "2021011"), 1, 1)

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": scroll.InstanceID,
	}}); err != nil {
		t.Fatalf("use cyclone scroll: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction || engine.State.PendingAction == nil {
		t.Fatalf("cyclone scroll should ask for enemy equipment choice")
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{p1.Equipment[0].InstanceID},
	}}); err != nil {
		t.Fatalf("resolve cyclone scroll: %v", err)
	}
	if p1.Equipment[0] != nil || len(p1.Graveyard) != 1 {
		t.Fatalf("cyclone scroll should destroy low-cost enemy equipment, equipment=%v graveyard=%d", p1.Equipment[0], len(p1.Graveyard))
	}
}

func TestFireDanceSkirtUltimatePurifiesFriendlyFireUnit(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	skirt := NewCardInstance(baseCard(t, "2121007"), 0, 1)
	p0.Equipment[0] = skirt
	unit := placeUnit(baseCard(t, "1121001"), 0, 0, 0, engine)
	unit.Statuses[StatusBurn] = 2
	unit.Statuses[StatusFreeze] = 1
	unit.Statuses[StatusStun] = 1

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  skirt.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use fire dance skirt ultimate: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction || engine.State.PendingAction == nil {
		t.Fatalf("fire dance skirt should ask for a friendly fire unit")
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{unit.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve fire dance skirt: %v", err)
	}
	if hasAnyNegativeStatus(unit) {
		t.Fatalf("fire dance skirt should clear negative statuses, statuses=%v", unit.Statuses)
	}
}

func TestShiningCrystalAddsStunToFriendlyLightSpells(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	p0.Equipment[0] = NewCardInstance(baseCard(t, "2521010"), 0, 1)
	p0.Skills[0] = readySkill(baseCard(t, "3521006"), 0)
	p0.Elements[model.ElementLight] = 4

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast light spell: %v", err)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve light spell: %v", err)
	}
	if target.Statuses[StatusStun] != 1 {
		t.Fatalf("shining crystal should add stun 1, statuses=%v", target.Statuses)
	}
}

func TestFlashRuneStunsEnemyFrontRowWhenEnemyCastsSkill(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	p1 := engine.State.Players[1]
	p0.Equipment[0] = NewCardInstance(baseCard(t, "2521011"), 0, 1)
	placeUnit(baseCard(t, "1021007"), 0, 1, 1, engine)
	front := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
	back := placeUnit(baseCard(t, "1021004"), 1, 0, 1, engine)
	p1.Skills[0] = readySkill(baseCard(t, "3021005"), 1)
	p1.Elements[model.ElementArcane] = 2
	engine.State.CurrentTurn = 1

	if err := engine.HandleAction(1, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p1.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(1),
	}}); err != nil {
		t.Fatalf("enemy cast skill: %v", err)
	}
	if front.Statuses[StatusStun] != 1 || back.Statuses[StatusStun] != 0 {
		t.Fatalf("flash rune should stun only enemy front row, front=%v back=%v", front.Statuses, back.Statuses)
	}
}

func TestFireBarrierBoostsFireSpellsAndAddsBurn(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021004"), 1, 1, 0, engine)
	p0.Skills[0] = readySkill(baseCard(t, "3121008"), 0)
	p0.Skills[1] = readySkill(baseCard(t, "3121001"), 0)
	p0.Elements[model.ElementFire] = 3

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[1].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast fireball with fire barrier: %v", err)
	}
	if engine.State.PendingSpell == nil || engine.State.PendingSpell.TotalPower != 5 {
		t.Fatalf("fire barrier should boost fireball power from 3 to 5, pending=%+v", engine.State.PendingSpell)
	}
	if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
		t.Fatalf("resolve fireball: %v", err)
	}
	if target.Statuses[StatusBurn] != 1 {
		t.Fatalf("fire barrier should add burn 1 on hit, statuses=%v", target.Statuses)
	}
}

func TestRuntimeLoadBonusCards(t *testing.T) {
	t.Run("blue crystal lamp pays light and gains light load", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		lamp := NewCardInstance(baseCard(t, "2521007"), 0, 1)
		p0.Equipment[0] = lamp
		p0.Elements[model.ElementLight] = 5

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  lamp.InstanceID,
			"ability_type": "ultimate",
		}}); err != nil {
			t.Fatalf("use blue crystal lamp ultimate: %v", err)
		}
		if p0.Elements[model.ElementLight] != 0 || effectiveElementsGain(lamp)[model.ElementLight] != 3 {
			t.Fatalf("lamp should pay 5 light and increase load to 3 light, elements=%v gains=%v", p0.Elements, effectiveElementsGain(lamp))
		}
	})

	t.Run("necromancy stone gains shadow load when friendly companion dies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		stone := NewCardInstance(baseCard(t, "2611001"), 0, 1)
		p0.Equipment[0] = stone
		unit := placeUnit(baseCard(t, "1021007"), 0, 0, 0, engine)

		engine.destroyUnit(unit, 0)
		if effectiveElementsGain(stone)[model.ElementShadow] != 1 {
			t.Fatalf("necromancy stone should gain shadow load, gains=%v", effectiveElementsGain(stone))
		}
	})

	t.Run("soul necklace gains shadow when friendly companion dies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		necklace := NewCardInstance(baseCard(t, "2621006"), 0, 1)
		p0.Equipment[0] = necklace
		unit := placeUnit(baseCard(t, "1021007"), 0, 0, 0, engine)

		engine.destroyUnit(unit, 0)
		if p0.Elements[model.ElementShadow] != 1 {
			t.Fatalf("soul necklace should gain 1 shadow, elements=%v", p0.Elements)
		}
	})
}

func TestItemSpellScrollEffects(t *testing.T) {
	t.Run("waterform bind consumes an enemy companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021007"), 1, 1, 0, engine)
		scroll := NewCardInstance(baseCard(t, "2221008"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementWater] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use waterform bind: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve waterform bind: %v", err)
		}
		if !target.IsHorizontal {
			t.Fatalf("waterform bind should consume target")
		}
	})

	t.Run("thunderstorm damages and stuns hit companions", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		a := placeUnit(baseCard(t, "1021007"), 1, 0, 0, engine)
		b := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		scroll := NewCardInstance(baseCard(t, "2321003"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementAir] = 4

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use thunderstorm: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{a.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve thunderstorm: %v", err)
		}
		if a.Statuses[StatusStun] != 1 || b.Statuses[StatusStun] != 1 {
			t.Fatalf("thunderstorm should stun all hit companions, a=%v b=%v", a.Statuses, b.Statuses)
		}
	})

	t.Run("chain lightning damages and draws", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021007"), 1, 1, 0, engine)
		target.CurrentLife = 3
		scroll := NewCardInstance(baseCard(t, "2321009"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1)}
		p0.Elements[model.ElementAir] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use chain lightning: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve chain lightning: %v", err)
		}
		if target.CurrentLife != 2 || len(p0.Hand) != 1 {
			t.Fatalf("chain lightning should damage and draw, life=%d hand=%d", target.CurrentLife, len(p0.Hand))
		}
	})
}

func TestLightSearchSpells(t *testing.T) {
	t.Run("united hope searches a light companion among top five", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3501001"), 0)
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1521005"), 0, 1),
			NewCardInstance(baseCard(t, "1021004"), 0, 1),
		}
		p0.Elements[model.ElementLight] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast united hope: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0].Card.Number != "1521005" {
			t.Fatalf("united hope should find light companion, hand=%v", p0.Hand)
		}
	})

	t.Run("call of hope searches first light companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3521007"), 0)
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1021004"), 0, 1),
			NewCardInstance(baseCard(t, "1521005"), 0, 1),
		}
		p0.Elements[model.ElementLight] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast call of hope: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0].Card.Number != "1521005" {
			t.Fatalf("call of hope should find first light companion, hand=%v", p0.Hand)
		}
	})
}

func TestDelayedAndTurnCountEffects(t *testing.T) {
	t.Run("focus scroll gains arcane next turn start", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2021021"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementArcane] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use focus scroll: %v", err)
		}
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Type != TempModDelayedElementGain {
			t.Fatalf("focus scroll should add delayed element modifier, modifiers=%v", p0.TempModifiers)
		}

		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end p0 turn: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end p1 turn: %v", err)
		}
		if p0.Elements[model.ElementArcane] != 3 || len(p0.TempModifiers) != 0 {
			t.Fatalf("focus scroll should gain 3 arcane next turn, elements=%v modifiers=%v", p0.Elements, p0.TempModifiers)
		}
	})

	t.Run("energy potion modifier expires when opponent turn ends", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3021005"), 0)
		p0.Skills[0].IsHorizontal = true
		potion := NewCardInstance(baseCard(t, "2221005"), 0, 1)
		p0.Hand = []*CardInstance{potion}
		p0.Elements[model.ElementWater] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": potion.InstanceID,
		}}); err != nil {
			t.Fatalf("use energy potion: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end p0 turn: %v", err)
		}
		if p0.Skills[0].IsHorizontal {
			t.Fatalf("skills should reset at own turn end under current turn rules")
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end p1 turn: %v", err)
		}
		if p0.Skills[0].IsHorizontal || len(p0.TempModifiers) != 0 {
			t.Fatalf("energy potion should reset skills after opponent turn, horizontal=%v modifiers=%v", p0.Skills[0].IsHorizontal, p0.TempModifiers)
		}
	})

	t.Run("flame rekindle gains fire for fire spells cast this turn", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3121014"), 0)
		p0.Elements[model.ElementFire] = 3
		placeUnit(baseCard(t, "1021007"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve fireball: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[1].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast flame rekindle: %v", err)
		}
		if p0.Elements[model.ElementFire] != 3 {
			t.Fatalf("flame rekindle should pay 1 and gain 2 fire for two fire casts, elements=%v", p0.Elements)
		}
	})
}

func TestBagOfTricksSacrificesToSearchConsumable(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	bag := NewCardInstance(baseCard(t, "2021006"), 0, 1)
	p0.Equipment[0] = bag
	target := NewCardInstance(baseCard(t, "2021021"), 0, 1)
	p0.Deck = []*CardInstance{NewCardInstance(baseCard(t, "1021001"), 0, 1), target}

	if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
		"instance_id":  bag.InstanceID,
		"ability_type": "ultimate",
	}}); err != nil {
		t.Fatalf("use bag of tricks ultimate: %v", err)
	}
	if engine.State.Phase != PhaseWaitingAction || engine.State.PendingAction == nil {
		t.Fatalf("bag of tricks should ask which consumable to search")
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{target.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve bag of tricks: %v", err)
	}
	if p0.Equipment[0] != nil || len(p0.Graveyard) != 1 || len(p0.Hand) != 1 || p0.Hand[0] != target {
		t.Fatalf("bag should sacrifice and search target, equipment=%v grave=%d hand=%d", p0.Equipment[0], len(p0.Graveyard), len(p0.Hand))
	}
}

func TestUtilityScrollAndForesightEffects(t *testing.T) {
	t.Run("cursed scroll draws two and discards them at turn end", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		scroll := NewCardInstance(baseCard(t, "2021019"), 0, 1)
		draw1 := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		draw2 := NewCardInstance(baseCard(t, "1021004"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Deck = []*CardInstance{draw1, draw2}
		p0.Elements[model.ElementArcane] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use cursed scroll: %v", err)
		}
		if len(p0.Hand) != 2 || len(p0.DiscardAtTurnEnd) != 2 {
			t.Fatalf("cursed scroll should draw and mark cards, hand=%d markers=%v", len(p0.Hand), p0.DiscardAtTurnEnd)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end turn: %v", err)
		}
		if len(p0.Hand) != 0 || len(p0.Graveyard) != 3 {
			t.Fatalf("cursed scroll should discard drawn cards at turn end, hand=%d graveyard=%d", len(p0.Hand), len(p0.Graveyard))
		}
	})

	t.Run("foresight can move selected top deck cards to bottom", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3021002"), 0)
		a := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		b := NewCardInstance(baseCard(t, "1021004"), 0, 1)
		c := NewCardInstance(baseCard(t, "1021007"), 0, 1)
		p0.Deck = []*CardInstance{a, b, c}
		p0.Elements[model.ElementArcane] = 1

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast foresight: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{a.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve foresight: %v", err)
		}
		if p0.Deck[0] != b || p0.Deck[2] != a {
			t.Fatalf("foresight should move selected card to bottom, deck=%s,%s,%s", p0.Deck[0].Card.Number, p0.Deck[1].Card.Number, p0.Deck[2].Card.Number)
		}
	})

	t.Run("deep frost curse permanently freezes enemy companion", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021007"), 1, 1, 0, engine)
		scroll := NewCardInstance(baseCard(t, "2221013"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementWater] = 3

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use deep frost curse: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve deep frost curse: %v", err)
		}
		if target.Statuses[StatusFreeze] < 50 {
			t.Fatalf("deep frost curse should apply long freeze, statuses=%v", target.Statuses)
		}
	})

	t.Run("soul devour weakens selected enemy skills", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		p1.Skills[0] = readySkill(baseCard(t, "3021005"), 1)
		p1.Skills[1] = readySkill(baseCard(t, "3121001"), 1)
		scroll := NewCardInstance(baseCard(t, "2621008"), 0, 1)
		p0.Hand = []*CardInstance{scroll}
		p0.Elements[model.ElementShadow] = 2

		if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
			"instance_id": scroll.InstanceID,
		}}); err != nil {
			t.Fatalf("use soul devour: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{p1.Skills[0].InstanceID, p1.Skills[1].InstanceID},
		}}); err != nil {
			t.Fatalf("resolve soul devour: %v", err)
		}
		if p1.Skills[0].Statuses[StatusWeaken] != 1 || p1.Skills[1].Statuses[StatusWeaken] != 1 {
			t.Fatalf("soul devour should weaken selected skills, s0=%v s1=%v", p1.Skills[0].Statuses, p1.Skills[1].Statuses)
		}
	})
}

func TestRebirthScrollRevivesLightCompanion(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	dead := NewCardInstance(baseCard(t, "1521005"), 0, 1)
	dead.CurrentLife = 0
	p0.Graveyard = []*CardInstance{dead}
	scroll := NewCardInstance(baseCard(t, "2521005"), 0, 1)
	p0.Hand = []*CardInstance{scroll}
	p0.Elements[model.ElementLight] = 10

	if err := engine.HandleAction(0, ActionMessage{Action: "use_item", Data: map[string]any{
		"instance_id": scroll.InstanceID,
	}}); err != nil {
		t.Fatalf("use rebirth scroll: %v", err)
	}
	if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
		"selected": []any{dead.InstanceID},
	}}); err != nil {
		t.Fatalf("resolve rebirth scroll: %v", err)
	}
	if len(p0.Graveyard) != 1 || dead.Position == nil || p0.Units[dead.Position.Col][dead.Position.Row] != dead {
		t.Fatalf("rebirth should revive companion onto board, graveyard=%d position=%v", len(p0.Graveyard), dead.Position)
	}
	if dead.CurrentLife != dead.Card.Life {
		t.Fatalf("rebirth should restore life, life=%d", dead.CurrentLife)
	}
}

func TestSkillContributionModifiers(t *testing.T) {
	t.Run("water and frost blade boost spell power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3221002"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3221003"), 0)
		p0.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
			"boost_ids":   []any{p0.Skills[1].InstanceID},
		}}); err != nil {
			t.Fatalf("cast boosted water spell: %v", err)
		}
		want := baseCard(t, "3221002").Power + baseCard(t, "3221003").Power + 2
		if engine.State.PendingSpell.TotalPower != want {
			t.Fatalf("freezing current should add +2 power while boosting water, got %d want %d", engine.State.PendingSpell.TotalPower, want)
		}
	})

	t.Run("frost blade adds power while attacking", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3221009"), 0)
		p0.Elements[model.ElementWater] = 10
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast frost blade: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != baseCard(t, "3221009").Power+2 {
			t.Fatalf("frost blade should add +2 power while attacking, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("joint casting boosts spell damage", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3421008"), 0)
		p0.Elements[model.ElementFire] = 10
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
			"boost_ids":   []any{p0.Skills[1].InstanceID},
		}}); err != nil {
			t.Fatalf("cast spell with joint casting: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve joint casting hit: %v", err)
		}
		if target.CurrentLife != target.Card.Life-2 {
			t.Fatalf("joint casting should add +1 attack damage, life=%d", target.CurrentLife)
		}
	})

	t.Run("moonlight can defend above its printed power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3101002"), 0)
		p1.Skills[0] = readySkill(baseCard(t, "3521013"), 1)
		p0.Elements[model.ElementFire] = 10
		p1.Elements[model.ElementLight] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast all fires as one: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != 5 {
			t.Fatalf("expected 5 attack power, got %+v", engine.State.PendingSpell)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{p1.Skills[0].InstanceID},
		}}); err != nil {
			t.Fatalf("defend with moonlight: %v", err)
		}
		if target.CurrentLife != target.Card.Life {
			t.Fatalf("moonlight +2 defense power should stop the spell, life=%d", target.CurrentLife)
		}
	})

	t.Run("all fires as one gains attack from power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3101002"), 0)
		p0.Elements[model.ElementFire] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast all fires as one: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve all fires as one: %v", err)
		}
		if target.CurrentLife != target.Card.Life-2 {
			t.Fatalf("all fires as one should deal 2 at 5 power, life=%d", target.CurrentLife)
		}
	})
}

func TestSplashBlizzardAndSoulBiteEffects(t *testing.T) {
	t.Run("splash range damages and freezes adjacent units", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		center := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		left := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		down := placeUnit(baseCard(t, "1021001"), 1, 1, 1, engine)
		diagonal := placeUnit(baseCard(t, "1021001"), 1, 0, 1, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3221005"), 0)
		p0.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast ice field: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve ice field: %v", err)
		}
		for _, unit := range []*CardInstance{center, left, down} {
			if unit.CurrentLife != unit.Card.Life-1 || unit.Statuses[StatusFreeze] != 1 {
				t.Fatalf("splash unit should be damaged and frozen, unit=%s life=%d statuses=%v", unit.Card.Name, unit.CurrentLife, unit.Statuses)
			}
		}
		if diagonal.CurrentLife != diagonal.Card.Life || diagonal.Statuses[StatusFreeze] != 0 {
			t.Fatalf("diagonal unit should not be hit by splash, life=%d statuses=%v", diagonal.CurrentLife, diagonal.Statuses)
		}
	})

	t.Run("blizzard gives water and air spells power and freeze on hit", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3221002"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3221015"), 0)
		p0.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast water spell with blizzard: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != baseCard(t, "3221002").Power+1 {
			t.Fatalf("blizzard should add +1 power to water spells, pending=%+v", engine.State.PendingSpell)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve water spell with blizzard: %v", err)
		}
		if target.Statuses[StatusFreeze] != 1 {
			t.Fatalf("blizzard should freeze hit target, statuses=%v", target.Statuses)
		}
	})

	t.Run("dead soul bite weakens enemy skills and weaken lowers power", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3621006"), 0)
		p1.Skills[0] = readySkill(baseCard(t, "3121001"), 1)
		p0.Elements[model.ElementShadow] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast dead soul bite: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve dead soul bite: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{p1.Skills[0].InstanceID},
		}}); err != nil {
			t.Fatalf("resolve dead soul bite weaken: %v", err)
		}
		if p1.Skills[0].Statuses[StatusWeaken] != 1 {
			t.Fatalf("dead soul bite should add weaken, statuses=%v", p1.Skills[0].Statuses)
		}
		if power := engine.effectiveSpellPower(1, p1.Skills[0], nil); power != baseCard(t, "3121001").Power-1 {
			t.Fatalf("weaken should lower spell power by 1, got %d", power)
		}
	})
}

func TestStormEarthAndDeadFuryEffects(t *testing.T) {
	t.Run("storm fury adds air spell power for each card in hand", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3321005"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3301001"), 0)
		p0.Hand = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			NewCardInstance(baseCard(t, "1021002"), 0, 1),
		}
		p0.Elements[model.ElementAir] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast air spell with storm fury: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != baseCard(t, "3321005").Power+2 {
			t.Fatalf("storm fury should add hand size to air spell power, pending=%+v", engine.State.PendingSpell)
		}
	})

	t.Run("earth resonance hits all enemies and gains attack from large allies", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		front := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		back := placeUnit(baseCard(t, "1021001"), 1, 0, 1, engine)
		placeUnit(baseCard(t, "1421001"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1521002"), 0, 2, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3421013"), 0)
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast earth resonance: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve earth resonance: %v", err)
		}
		if front.CurrentLife != front.Card.Life-2 || back.CurrentLife != back.Card.Life-2 {
			t.Fatalf("earth resonance should deal 2 to all enemies, front=%d back=%d", front.CurrentLife, back.CurrentLife)
		}
	})

	t.Run("dead fury gains permanent instance power when companions die", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p0.Skills[0] = readySkill(baseCard(t, "3621008"), 0)
		dead := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)

		engine.destroyUnit(dead, 1)
		if p0.Skills[0].PowerBonus != 1 {
			t.Fatalf("dead fury should gain +1 instance power, bonus=%d", p0.Skills[0].PowerBonus)
		}
		if power := engine.effectiveSpellPower(0, p0.Skills[0], nil); power != baseCard(t, "3621008").Power+1 {
			t.Fatalf("dead fury power should include death bonus, got %d", power)
		}
	})
}

func TestBurnSpellGetsPowerAgainstBurningTarget(t *testing.T) {
	engine := setupReportedBugEngine(t)
	p0 := engine.State.Players[0]
	target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
	target.Statuses[StatusBurn] = 1
	p0.Skills[0] = readySkill(baseCard(t, "3121002"), 0)
	p0.Elements[model.ElementFire] = 10

	if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
		"instance_id": p0.Skills[0].InstanceID,
		"target_type": "unit",
		"target_col":  float64(1),
		"target_row":  float64(0),
	}}); err != nil {
		t.Fatalf("cast burn against burning target: %v", err)
	}
	if engine.State.PendingSpell.TotalPower != baseCard(t, "3121002").Power+2 {
		t.Fatalf("burn should gain +2 power against burning target, pending=%+v", engine.State.PendingSpell)
	}
}

func TestDefenseAndPositionSkillEffects(t *testing.T) {
	t.Run("rivers to sea gains water after successful defense", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p1.Skills[0] = readySkill(baseCard(t, "3201001"), 1)
		p0.Elements[model.ElementFire] = 10
		p1.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball: %v", err)
		}
		before := p1.Elements[model.ElementWater]
		if err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{p1.Skills[0].InstanceID},
		}}); err != nil {
			t.Fatalf("defend with rivers to sea: %v", err)
		}
		if p1.Elements[model.ElementWater] != before-baseCard(t, "3201001").ElementsExpense[model.ElementWater]+1 {
			t.Fatalf("rivers to sea should gain water equal attack spell attack, elements=%v before=%d", p1.Elements, before)
		}
	})

	t.Run("static barrier stuns attacker front row after failed defense", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		attackerFront := placeUnit(baseCard(t, "1021001"), 0, 1, 0, engine)
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3421013"), 0)
		p1.Skills[0] = readySkill(baseCard(t, "3321015"), 1)
		p0.Elements[model.ElementEarth] = 10
		p1.Elements[model.ElementAir] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast earth resonance: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{p1.Skills[0].InstanceID},
		}}); err != nil {
			t.Fatalf("defend with static barrier: %v", err)
		}
		if attackerFront.Statuses[StatusStun] != 1 {
			t.Fatalf("static barrier should stun attacker front row, statuses=%v", attackerFront.Statuses)
		}
		if target.CurrentLife == target.Card.Life {
			t.Fatalf("failed defense should still let the spell hit")
		}
	})

	t.Run("sky sense adds power when affected area includes back row", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		front := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		back := placeUnit(baseCard(t, "1021001"), 1, 1, 1, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3421013"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3321012"), 0)
		p0.Elements[model.ElementEarth] = 10
		p0.Elements[model.ElementAir] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast earth resonance with sky sense: %v", err)
		}
		want := baseCard(t, "3421013").Power + 2
		if engine.State.PendingSpell.TotalPower != want {
			t.Fatalf("sky sense should add +2 power when area includes back row, got %d want %d front=%v back=%v", engine.State.PendingSpell.TotalPower, want, front.Position, back.Position)
		}
	})

	t.Run("earthshaker gains power from friendly earth load", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		placeUnit(baseCard(t, "1421008"), 0, 0, 0, engine)
		placeUnit(baseCard(t, "1421010"), 0, 2, 0, engine)
		placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3421012"), 0)
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast earthshaker: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != baseCard(t, "3421012").Power+2 {
			t.Fatalf("earthshaker should add friendly earth load to power, pending=%+v", engine.State.PendingSpell)
		}
	})
}

func TestSelectionSorcerySkillEffects(t *testing.T) {
	t.Run("call lightning discards then stuns selected enemy", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		discard := NewCardInstance(baseCard(t, "1021001"), 0, 1)
		p0.Hand = []*CardInstance{discard}
		p0.Skills[0] = readySkill(baseCard(t, "3321014"), 0)
		p0.Elements[model.ElementAir] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast call lightning: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{discard.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve call lightning discard: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve call lightning stun: %v", err)
		}
		if len(p0.Graveyard) != 1 || target.Statuses[StatusStun] != 1 {
			t.Fatalf("call lightning should discard and stun, graveyard=%d statuses=%v", len(p0.Graveyard), target.Statuses)
		}
	})

	t.Run("natural growth gives selected earth companion load at turn end", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1421010"), 0, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3421011"), 0)
		p0.Elements[model.ElementEarth] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast natural growth: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve natural growth: %v", err)
		}
		if effectiveElementsGain(target)[model.ElementEarth] != 1 {
			t.Fatalf("natural growth should not apply before end turn, load=%v", effectiveElementsGain(target))
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "end_turn", Data: map[string]any{}}); err != nil {
			t.Fatalf("end turn after natural growth: %v", err)
		}
		if effectiveElementsGain(target)[model.ElementEarth] != 2 {
			t.Fatalf("natural growth should add earth load at end turn, load=%v", effectiveElementsGain(target))
		}
	})

	t.Run("blood demon blast sacrifices companion and damages enemy front row", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		sacrifice := placeUnit(baseCard(t, "1521002"), 0, 1, 0, engine)
		frontA := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		frontB := placeUnit(baseCard(t, "1021001"), 1, 2, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3621010"), 0)
		p0.Elements[model.ElementShadow] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast blood demon blast: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{sacrifice.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve blood demon blast: %v", err)
		}
		if p0.Units[1][0] != nil || len(p0.Graveyard) == 0 {
			t.Fatalf("blood demon blast should sacrifice the selected companion")
		}
		if frontA.CurrentLife != frontA.Card.Life-sacrifice.Card.Life || frontB.CurrentLife != frontB.Card.Life-sacrifice.Card.Life {
			t.Fatalf("blood demon blast should damage enemy front row, a=%d b=%d", frontA.CurrentLife, frontB.CurrentLife)
		}
	})
}

func TestRemainingPassiveSkillEffects(t *testing.T) {
	t.Run("rapid sandstorm weakens low original power spells", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3121001"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3421015"), 0)
		p0.Elements[model.ElementFire] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast fireball under rapid sandstorm: %v", err)
		}
		if engine.State.PendingSpell.TotalPower != 1 {
			t.Fatalf("rapid sandstorm should reduce fireball power to 1, pending=%+v", engine.State.PendingSpell)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve fireball under rapid sandstorm: %v", err)
		}
		if target.CurrentLife != target.Card.Life {
			t.Fatalf("rapid sandstorm should reduce fireball damage to 0, life=%d", target.CurrentLife)
		}
	})

	t.Run("undead wall gains defense power after friendly death", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		p1 := engine.State.Players[1]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		dead := placeUnit(baseCard(t, "1021001"), 1, 0, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3101002"), 0)
		p1.Skills[0] = readySkill(baseCard(t, "3621013"), 1)
		p0.Elements[model.ElementFire] = 10
		p1.Elements[model.ElementShadow] = 10

		engine.destroyUnit(dead, 1)
		if p1.Skills[0].Statuses[recentFriendlyDeathStatus] <= 0 {
			t.Fatalf("undead wall should remember recent friendly death, statuses=%v", p1.Skills[0].Statuses)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast all fires as one: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "defend", Data: map[string]any{
			"skill_ids": []any{p1.Skills[0].InstanceID},
		}}); err != nil {
			t.Fatalf("defend with undead wall: %v", err)
		}
		if target.CurrentLife != target.Card.Life {
			t.Fatalf("undead wall should defend successfully with recent death bonus, life=%d", target.CurrentLife)
		}
	})
}

func TestChoiceUtilitySkillEffects(t *testing.T) {
	t.Run("elemental enchant lets next spell apply chosen status", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		target := placeUnit(baseCard(t, "1021001"), 1, 1, 0, engine)
		p0.Skills[0] = readySkill(baseCard(t, "3021007"), 0)
		p0.Skills[1] = readySkill(baseCard(t, "3121001"), 0)
		p0.Elements[model.ElementArcane] = 10
		p0.Elements[model.ElementFire] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast elemental enchant: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{StatusBurn},
		}}); err != nil {
			t.Fatalf("choose elemental enchant status: %v", err)
		}
		if len(p0.TempModifiers) != 1 || p0.TempModifiers[0].Status != StatusBurn {
			t.Fatalf("elemental enchant should add next spell status modifier, modifiers=%v", p0.TempModifiers)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[1].InstanceID,
			"target_type": "unit",
			"target_col":  float64(1),
			"target_row":  float64(0),
		}}); err != nil {
			t.Fatalf("cast enchanted fireball: %v", err)
		}
		if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
			t.Fatalf("resolve enchanted fireball: %v", err)
		}
		if target.Statuses[StatusBurn] != 1 {
			t.Fatalf("enchanted fireball should burn target, statuses=%v", target.Statuses)
		}
	})

	t.Run("mind tempering permanently buffs selected skill", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		targetSkill := readySkill(baseCard(t, "3121001"), 0)
		p0.Skills[0] = readySkill(baseCard(t, "3021012"), 0)
		p0.Skills[1] = targetSkill
		p0.Elements[model.ElementArcane] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast mind tempering: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{targetSkill.InstanceID},
		}}); err != nil {
			t.Fatalf("choose mind tempering target: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{"power"},
		}}); err != nil {
			t.Fatalf("choose mind tempering mode: %v", err)
		}
		if targetSkill.PowerBonus != 3 {
			t.Fatalf("mind tempering should add +3 power, bonus=%d", targetSkill.PowerBonus)
		}
	})

	t.Run("water divination searches water card among top four", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		waterCard := NewCardInstance(baseCard(t, "1221001"), 0, 1)
		p0.Deck = []*CardInstance{
			NewCardInstance(baseCard(t, "1021001"), 0, 1),
			waterCard,
			NewCardInstance(baseCard(t, "1021002"), 0, 1),
		}
		p0.Skills[0] = readySkill(baseCard(t, "3221007"), 0)
		p0.Elements[model.ElementWater] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{
			"instance_id": p0.Skills[0].InstanceID,
			"target_type": "none",
		}}); err != nil {
			t.Fatalf("cast water divination: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{waterCard.InstanceID},
		}}); err != nil {
			t.Fatalf("choose water divination card: %v", err)
		}
		if len(p0.Hand) != 1 || p0.Hand[0] != waterCard {
			t.Fatalf("water divination should search water card to hand, hand=%v", cardsToInfo(p0.Hand))
		}
	})
}

func TestCompanionUtilityEffects(t *testing.T) {
	t.Run("farolank gains adjacent companion load on enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		left := placeUnit(baseCard(t, "1421008"), 0, 0, 1, engine)
		right := placeUnit(baseCard(t, "1521002"), 0, 2, 1, engine)
		farolank := NewCardInstance(baseCard(t, "1011003"), 0, 1)
		p0.Hand = []*CardInstance{farolank}
		p0.Elements[model.ElementArcane] = 20

		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": farolank.InstanceID,
			"col":         float64(1),
			"row":         float64(1),
		}}); err != nil {
			t.Fatalf("summon farolank: %v", err)
		}
		load := effectiveElementsGain(farolank)
		if load[model.ElementEarth] != effectiveElementsGain(left)[model.ElementEarth] ||
			load[model.ElementAir] != effectiveElementsGain(left)[model.ElementAir] ||
			load[model.ElementLight] != effectiveElementsGain(right)[model.ElementLight]+farolank.Card.ElementsGain[model.ElementLight] {
			t.Fatalf("farolank should gain adjacent companion load, load=%v", load)
		}
	})

	t.Run("specialist mage changes load to chosen element", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		mage := NewCardInstance(baseCard(t, "1021010"), 0, 1)
		p0.Hand = []*CardInstance{mage}
		p0.Elements[model.ElementArcane] = 10

		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": mage.InstanceID,
			"col":         float64(1),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("summon specialist mage: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{model.ElementFire},
		}}); err != nil {
			t.Fatalf("choose specialist mage element: %v", err)
		}
		load := effectiveElementsGain(mage)
		if load[model.ElementFire] != 2 || len(load) != 1 {
			t.Fatalf("specialist mage should convert load to chosen element, load=%v", load)
		}
	})

	t.Run("blessed girl gives adjacent earth companion load", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		girl := placeUnit(baseCard(t, "1421009"), 0, 1, 1, engine)
		target := placeUnit(baseCard(t, "1421010"), 0, 1, 0, engine)

		if err := engine.HandleAction(0, ActionMessage{Action: "use_ability", Data: map[string]any{
			"instance_id":  girl.InstanceID,
			"ability_type": "per_turn",
		}}); err != nil {
			t.Fatalf("use blessed girl ability: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{target.InstanceID},
		}}); err != nil {
			t.Fatalf("resolve blessed girl load: %v", err)
		}
		if effectiveElementsGain(target)[model.ElementEarth] != target.Card.ElementsGain[model.ElementEarth]+1 {
			t.Fatalf("blessed girl should add earth load, load=%v", effectiveElementsGain(target))
		}
	})

	t.Run("lundesal buffs a selected skill on enter", func(t *testing.T) {
		engine := setupReportedBugEngine(t)
		p0 := engine.State.Players[0]
		lundesal := NewCardInstance(baseCard(t, "1511002"), 0, 1)
		targetSkill := readySkill(baseCard(t, "3121001"), 0)
		p0.Hand = []*CardInstance{lundesal}
		p0.Skills[0] = targetSkill
		p0.Elements[model.ElementLight] = 20

		if err := engine.HandleAction(0, ActionMessage{Action: "summon", Data: map[string]any{
			"instance_id": lundesal.InstanceID,
			"col":         float64(1),
			"row":         float64(0),
		}}); err != nil {
			t.Fatalf("summon lundesal: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{targetSkill.InstanceID},
		}}); err != nil {
			t.Fatalf("choose lundesal skill: %v", err)
		}
		if err := engine.HandleAction(0, ActionMessage{Action: "resolve_action", Data: map[string]any{
			"selected": []any{"attack"},
		}}); err != nil {
			t.Fatalf("choose lundesal buff: %v", err)
		}
		if targetSkill.AttackBonus != 1 {
			t.Fatalf("lundesal should add +1 attack, bonus=%d", targetSkill.AttackBonus)
		}
	})
}
