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
	"unit":        true,
	"terrain":     true,
	"graveyard":   true,
}

// DebugAddCard creates real runtime card instances from the active card DB.
// It is intended only for explicitly marked test rooms.
func (e *Engine) DebugAddCard(playerID int, cardNumber string, zone string, count int) error {
	return e.DebugAddCardAt(playerID, cardNumber, zone, count, -1, -1)
}

func (e *Engine) DebugAddCardAt(playerID int, cardNumber string, zone string, count int, col int, row int) error {
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
	if zone == "unit" && !card.IsCompanion() && !card.IsHero() {
		return fmt.Errorf("only unit cards can be added to unit")
	}
	if zone == "terrain" && !card.IsItem() {
		return fmt.Errorf("only terrain/item cards can be added to terrain")
	}
	if (zone == "unit" || zone == "terrain") && count != 1 {
		return fmt.Errorf("%s placement supports count 1", zone)
	}

	ps := e.State.Players[playerID]
	var addedToHand []*CardInstance
	for i := 0; i < count; i++ {
		instance := e.newCardInstance(card, playerID, e.State.TurnNumber)
		switch zone {
		case "hand":
			addedToHand = append(addedToHand, instance)
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
		case "unit":
			if !validGridPosition(col, row) {
				return fmt.Errorf("unit placement requires valid col and row")
			}
			if ps.Units[col][row] != nil {
				return fmt.Errorf("unit slot is occupied")
			}
			instance.IsHorizontal = false
			instance.Position = &Position{Col: col, Row: row}
			ps.Units[col][row] = instance
		case "terrain":
			if !validGridPosition(col, row) {
				return fmt.Errorf("terrain placement requires valid col and row")
			}
			if ps.Terrain[col][row] != nil {
				return fmt.Errorf("terrain slot is occupied")
			}
			instance.IsHorizontal = false
			instance.Position = &Position{Col: col, Row: row}
			ps.Terrain[col][row] = instance
		case "graveyard":
			ps.Graveyard = append(ps.Graveyard, instance)
		}
	}
	e.addCardsToHand(playerID, addedToHand)
	return nil
}

func validGridPosition(col int, row int) bool {
	return col >= 0 && col < 3 && row >= 0 && row < 3
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

type DebugCardMutation struct {
	InstanceID string         `json:"instance_id"`
	Life       *int           `json:"life,omitempty"`
	Attack     *int           `json:"attack,omitempty"`
	Horizontal *bool          `json:"horizontal,omitempty"`
	Statuses   map[string]int `json:"statuses,omitempty"`
	Mode       string         `json:"mode,omitempty"`
}

func (e *Engine) DebugMutateCard(playerID int, mutation DebugCardMutation) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.State == nil {
		return fmt.Errorf("game is not started")
	}
	if playerID < 0 || playerID >= len(e.State.Players) {
		return fmt.Errorf("invalid player")
	}
	if mutation.InstanceID == "" {
		return fmt.Errorf("missing instance id")
	}

	card := e.findCardInstance(e.State.Players[playerID], mutation.InstanceID)
	if card == nil {
		return fmt.Errorf("card instance not found")
	}
	if mutation.Life != nil {
		if *mutation.Life < 0 {
			return fmt.Errorf("life cannot be negative")
		}
		card.CurrentLife = *mutation.Life
	}
	if mutation.Attack != nil {
		card.CurrentAttack = *mutation.Attack
	}
	if mutation.Horizontal != nil {
		card.IsHorizontal = *mutation.Horizontal
	}
	if mutation.Statuses != nil {
		if mutation.Mode == "replace" {
			card.Statuses = make(map[string]int)
		}
		if card.Statuses == nil {
			card.Statuses = make(map[string]int)
		}
		for status, amount := range mutation.Statuses {
			if amount < 0 {
				return fmt.Errorf("status %s cannot be negative", status)
			}
			if amount == 0 {
				delete(card.Statuses, status)
			} else {
				card.Statuses[status] = amount
			}
		}
	}
	return nil
}

func (e *Engine) findCardInstance(ps *PlayerState, instanceID string) *CardInstance {
	if ps == nil || instanceID == "" {
		return nil
	}
	all := make([]*CardInstance, 0, 64)
	all = append(all, ps.Hero)
	all = append(all, ps.Hand...)
	all = append(all, ps.Deck...)
	all = append(all, ps.SkillPool...)
	for i := range ps.Skills {
		all = append(all, ps.Skills[i])
	}
	for i := range ps.Equipment {
		all = append(all, ps.Equipment[i])
	}
	all = append(all, ps.Graveyard...)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			all = append(all, ps.Units[col][row], ps.Terrain[col][row])
		}
	}
	for _, card := range all {
		if card != nil && card.InstanceID == instanceID {
			return card
		}
	}
	return nil
}
