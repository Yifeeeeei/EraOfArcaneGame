package game

import (
	"fmt"

	"eraofarcane/model"
)

var debugCardZones = map[string]bool{
	"hand":        true,
	"deck_top":    true,
	"deck_bottom": true,
	"skill_pool":  true,
	"skill_slot":  true,
	"equipment":   true,
	"graveyard":   true,
}

// DebugAddCard creates real runtime card instances from the active card DB.
// It is intended only for explicitly marked test rooms.
func (e *Engine) DebugAddCard(playerID int, cardNumber string, zone string, count int) error {
	if count <= 0 {
		count = 1
	}
	if count > 20 {
		return fmt.Errorf("count must be at most 20")
	}
	if zone == "" {
		zone = "hand"
	}
	if !debugCardZones[zone] {
		return fmt.Errorf("unsupported debug zone %q", zone)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.State == nil {
		return fmt.Errorf("game is not started")
	}
	if playerID < 0 || playerID >= len(e.State.Players) {
		return fmt.Errorf("invalid player")
	}
	card := getCardDB()[cardNumber]
	if card == nil {
		return fmt.Errorf("card %s is not in the supported card pool", cardNumber)
	}
	if zone == "skill_pool" && !card.IsSkill() {
		return fmt.Errorf("only skill cards can be added to skill_pool")
	}
	if zone == "skill_slot" && !card.IsSkill() {
		return fmt.Errorf("only skill cards can be added to skill_slot")
	}
	if zone == "equipment" && !card.IsItem() {
		return fmt.Errorf("only item cards can be added to equipment")
	}

	ps := e.State.Players[playerID]
	for i := 0; i < count; i++ {
		instance := NewCardInstance(card, playerID, e.State.TurnNumber)
		switch zone {
		case "hand":
			ps.Hand = append(ps.Hand, instance)
		case "deck_top":
			ps.Deck = append([]*CardInstance{instance}, ps.Deck...)
		case "deck_bottom":
			ps.Deck = append(ps.Deck, instance)
		case "skill_pool":
			ps.SkillPool = append(ps.SkillPool, instance)
		case "skill_slot":
			slot := firstEmptySkillSlot(ps)
			if slot < 0 {
				return fmt.Errorf("no empty skill slot")
			}
			instance.IsHorizontal = false
			instance.SlotIndex = slot
			ps.Skills[slot] = instance
		case "equipment":
			slot := firstEmptyEquipmentSlot(ps)
			if slot < 0 {
				return fmt.Errorf("no empty equipment slot")
			}
			instance.IsHorizontal = false
			instance.SlotIndex = slot
			ps.Equipment[slot] = instance
		case "graveyard":
			ps.Graveyard = append(ps.Graveyard, instance)
		}
	}
	return nil
}

func firstEmptySkillSlot(ps *PlayerState) int {
	for i, card := range ps.Skills {
		if card == nil {
			return i
		}
	}
	return -1
}

func firstEmptyEquipmentSlot(ps *PlayerState) int {
	for i, card := range ps.Equipment {
		if card == nil {
			return i
		}
	}
	return -1
}

// DebugUpdateElements changes a player's element pool in a test room.
// mode="set" replaces values; all other modes add deltas.
func (e *Engine) DebugUpdateElements(playerID int, elements map[string]int, mode string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.State == nil {
		return fmt.Errorf("game is not started")
	}
	if playerID < 0 || playerID >= len(e.State.Players) {
		return fmt.Errorf("invalid player")
	}
	if len(elements) == 0 {
		return fmt.Errorf("missing elements")
	}

	allowed := make(map[string]bool, len(model.AllElements))
	for _, elem := range model.AllElements {
		allowed[elem] = true
	}

	ps := e.State.Players[playerID]
	if ps.Elements == nil {
		ps.Elements = make(map[string]int)
	}
	for elem, amount := range elements {
		if !allowed[elem] {
			return fmt.Errorf("unknown element %q", elem)
		}
		if mode == "set" {
			if amount < 0 {
				return fmt.Errorf("element amount cannot be negative")
			}
			ps.Elements[elem] = amount
			continue
		}
		next := ps.Elements[elem] + amount
		if next < 0 {
			return fmt.Errorf("%s would become negative", elem)
		}
		ps.Elements[elem] = next
	}
	return nil
}
