package game

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"eraofarcane/cards"
	"eraofarcane/model"
)

func setupBaseCardSmokeSuite(t *testing.T) {
	t.Helper()
	if cards.CardDB == nil {
		if err := cards.LoadCards("../../data/all_card_infos.json"); err != nil {
			t.Fatalf("load cards: %v", err)
		}
	}
	if cards.PlayableCardDB == nil {
		t.Fatal("playable card DB is nil")
	}

	previousDB := cardDBRef
	previousRegistry := globalRegistry
	SetCardDB(cards.PlayableCardDB)
	globalRegistry = NewEffectRegistry()
	RegisterAllCardEffects()
	AutoParseAndRegister()

	t.Cleanup(func() {
		if previousDB != nil {
			SetCardDB(previousDB)
		} else {
			SetCardDB(cards.CardDB)
		}
		globalRegistry = previousRegistry
	})
}

func sortedBaseCards(t *testing.T) []*model.Card {
	t.Helper()
	result := make([]*model.Card, 0, len(cards.PlayableCardDB))
	for _, card := range cards.PlayableCardDB {
		result = append(result, card)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Number < result[j].Number })
	return result
}

func baseSmokeEngine(t *testing.T) *Engine {
	t.Helper()
	hero := cards.PlayableCardDB["4311003"]
	if hero == nil {
		t.Fatal("missing base smoke hero 4311003")
	}

	engine := NewEngine("base-card-smoke", nil)
	engine.State.Phase = PhaseMain
	engine.State.CurrentTurn = 0
	engine.State.TurnNumber = 2
	engine.State.Winner = -1
	engine.State.HandLimit = 10

	for playerID := 0; playerID < 2; playerID++ {
		ps := &PlayerState{
			PlayerID:   playerID,
			PlayerName: fmt.Sprintf("P%d", playerID+1),
			Elements:   make(map[string]int),
			DeckDef:    &model.Deck{},
		}
		for _, elem := range model.AllElements {
			ps.Elements[elem] = 99
		}
		ps.Hero = NewCardInstance(hero, playerID, 1)
		ps.Hero.IsHorizontal = false
		ps.Hero.Position = &Position{Col: 1, Row: 1}
		ps.Units[1][1] = ps.Hero
		ps.Deck = baseDrawPile(playerID)
		engine.State.Players[playerID] = ps
	}

	return engine
}

func baseDrawPile(ownerID int) []*CardInstance {
	var ids []string
	for _, card := range cards.PlayableCardDB {
		if !card.IsHero() && !card.IsSkill() {
			ids = append(ids, card.Number)
		}
	}
	sort.Strings(ids)
	pile := make([]*CardInstance, 0, min(len(ids), 12))
	for _, id := range ids {
		pile = append(pile, NewCardInstance(cards.PlayableCardDB[id], ownerID, 0))
		if len(pile) >= 12 {
			break
		}
	}
	return pile
}

func setAllElements(ps *PlayerState, amount int) {
	for _, elem := range model.AllElements {
		ps.Elements[elem] = amount
	}
}

func attachSourceForEffect(engine *Engine, source *CardInstance) {
	ps := engine.State.Players[source.OwnerID]
	source.IsHorizontal = false
	switch {
	case source.Card.IsHero() || source.Card.IsCompanion():
		if source.Card.IsHero() {
			ps.Units[1][1] = source
			ps.Hero = source
			source.Position = &Position{Col: 1, Row: 1}
			return
		}
		ps.Units[0][1] = source
		source.Position = &Position{Col: 0, Row: 1}
	case source.Card.IsSkill():
		ps.Skills[0] = source
		source.SlotIndex = 0
	case source.Card.IsItem():
		if source.Card.IsTerrain() {
			ps.Terrain[0][1] = source
			source.Position = &Position{Col: 0, Row: 1}
		} else {
			ps.Equipment[0] = source
			source.SlotIndex = 0
		}
	}
}

func isConsumableItem(card *model.Card) bool {
	return card.IsItem() && strings.Contains(card.Tag, "消耗品")
}

func TestBaseCardPoolRejectsEveryNonBaseCard(t *testing.T) {
	setupBaseCardSmokeSuite(t)

	for id, card := range cards.CardDB {
		if card.VersionName == cards.BaseVersionName {
			continue
		}
		if _, ok := cards.PlayableCardDB[id]; ok {
			t.Fatalf("non-base card %s (%s/%s) is present in playable DB", id, card.VersionName, card.Name)
		}
	}
}

func TestEveryBaseCardHasRunnablePrimaryAction(t *testing.T) {
	setupBaseCardSmokeSuite(t)

	for _, card := range sortedBaseCards(t) {
		card := card
		t.Run(card.Number+"_"+card.Name, func(t *testing.T) {
			engine := baseSmokeEngine(t)
			ps := engine.State.Players[0]
			setAllElements(ps, 99)

			switch {
			case card.IsHero():
				hero := NewCardInstance(card, 0, engine.State.TurnNumber)
				hero.IsHorizontal = false
				hero.Position = &Position{Col: 1, Row: 1}
				ps.Hero = hero
				ps.Units[1][1] = hero
				if len(card.ElementsGain) > 0 {
					if err := engine.HandleAction(0, ActionMessage{
						Action: "consume",
						Data:   map[string]any{"instance_id": hero.InstanceID},
					}); err != nil {
						t.Fatalf("consume hero load failed: %v", err)
					}
				}
				for _, trigger := range []struct {
					name string
					typ  EffectTrigger
				}{
					{name: "per_turn", typ: TriggerPerTurn},
					{name: "ultimate", typ: TriggerUltimate},
				} {
					if !globalRegistry.HasEffect(card.Number, trigger.typ) {
						continue
					}
					hero.IsHorizontal = false
					if err := engine.HandleAction(0, ActionMessage{
						Action: "use_ability",
						Data: map[string]any{
							"instance_id":  hero.InstanceID,
							"ability_type": trigger.name,
							"target_id":    engine.State.Players[1].Hero.InstanceID,
						},
					}); err != nil {
						t.Fatalf("%s ability failed: %v", trigger.name, err)
					}
				}

			case card.IsCompanion():
				instance := NewCardInstance(card, 0, engine.State.TurnNumber)
				ps.Hand = append(ps.Hand, instance)
				if err := engine.HandleAction(0, ActionMessage{
					Action: "summon",
					Data: map[string]any{
						"instance_id": instance.InstanceID,
						"col":         float64(0),
						"row":         float64(0),
					},
				}); err != nil {
					t.Fatalf("summon failed: %v", err)
				}
				if ps.Units[0][0] == nil || ps.Units[0][0].InstanceID != instance.InstanceID {
					t.Fatalf("summon did not place card on board")
				}

			case card.IsSkill():
				instance := NewCardInstance(card, 0, engine.State.TurnNumber)
				ps.SkillPool = append(ps.SkillPool, instance)
				if err := engine.HandleAction(0, ActionMessage{
					Action: "learn_skill",
					Data:   map[string]any{"instance_id": instance.InstanceID},
				}); err != nil {
					t.Fatalf("learn skill failed: %v", err)
				}
				if ps.Skills[0] == nil || ps.Skills[0].InstanceID != instance.InstanceID {
					t.Fatalf("learn skill did not place card in skill slot")
				}
				if canUseSkillToAttack(card) {
					ps.Skills[0].IsHorizontal = false
					setAllElements(ps, 99)
					err := engine.HandleAction(0, ActionMessage{
						Action: "cast_spell",
						Data: map[string]any{
							"instance_id": instance.InstanceID,
							"target_type": "unit",
							"target_col":  float64(1),
							"target_row":  float64(1),
						},
					})
					if err != nil {
						t.Fatalf("cast skill failed: %v", err)
					}
					if engine.State.Phase == PhaseDefenseWindow {
						if err := engine.HandleAction(1, ActionMessage{Action: "no_defend", Data: map[string]any{}}); err != nil {
							t.Fatalf("resolve spell without defense failed: %v", err)
						}
					}
				}

			case card.IsItem():
				instance := NewCardInstance(card, 0, engine.State.TurnNumber)
				ps.Hand = append(ps.Hand, instance)
				action := "equip"
				data := map[string]any{"instance_id": instance.InstanceID}
				if card.IsTerrain() {
					action = "place_terrain"
					data["col"] = float64(0)
					data["row"] = float64(0)
				} else if isConsumableItem(card) {
					action = "use_item"
				}
				if err := engine.HandleAction(0, ActionMessage{Action: action, Data: data}); err != nil {
					t.Fatalf("%s item failed: %v", action, err)
				}
			}
		})
	}
}

func TestEveryRegisteredBaseCardEffectHandlerRuns(t *testing.T) {
	setupBaseCardSmokeSuite(t)

	for _, card := range sortedBaseCards(t) {
		effects := globalRegistry.GetAllEffects(card.Number)
		if len(effects) == 0 {
			continue
		}
		card := card
		t.Run(card.Number+"_"+card.Name, func(t *testing.T) {
			for i, effect := range effects {
				effect := effect
				t.Run(fmt.Sprintf("%s_%d", triggerName(effect.Trigger), i), func(t *testing.T) {
					engine := baseSmokeEngine(t)
					source := NewCardInstance(card, 0, engine.State.TurnNumber)
					attachSourceForEffect(engine, source)
					target := engine.State.Players[1].Hero
					ctx := &EffectContext{
						Engine:       engine,
						Source:       source,
						Target:       target,
						TargetPos:    target.Position,
						PlayerID:     0,
						OpponentID:   1,
						DamageAmount: 1,
						ExtraData: map[string]any{
							"damage": 1,
							"target": SpellTarget{Type: "unit", Position: *target.Position},
							"gained": map[string]int{model.ElementAir: 1},
						},
					}
					defer func() {
						if r := recover(); r != nil {
							t.Fatalf("effect panicked: %v", r)
						}
					}()
					if err := effect.Handler(ctx); err != nil {
						t.Fatalf("effect returned error: %v", err)
					}
				})
			}
		})
	}
}
