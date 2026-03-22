package model

import (
	"fmt"
	"strings"
)

// Deck represents a player's deck
type Deck struct {
	HeroID    string   `json:"hero_id"`
	MainDeck  []string `json:"main_deck"`  // 30 cards
	SkillPool []string `json:"skill_pool"` // 10 cards
	ExtraDeck []string `json:"extra_deck"` // spawned cards
}

// ParseDeckCode parses a deck code string into a Deck struct
// Format: heroID // card1 card2 ... // skill1 skill2 ... // extra1 extra2 ...
func ParseDeckCode(code string) (*Deck, error) {
	parts := strings.Split(code, "//")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid deck code: need at least 3 sections (hero // main // skills), got %d", len(parts))
	}

	deck := &Deck{}

	// Parse hero
	heroStr := strings.TrimSpace(parts[0])
	if heroStr == "" {
		return nil, fmt.Errorf("hero ID is empty")
	}
	deck.HeroID = heroStr

	// Parse main deck
	mainStr := strings.TrimSpace(parts[1])
	if mainStr != "" {
		deck.MainDeck = strings.Fields(mainStr)
	}

	// Parse skill pool
	skillStr := strings.TrimSpace(parts[2])
	if skillStr != "" {
		deck.SkillPool = strings.Fields(skillStr)
	}

	// Parse extra deck (optional)
	if len(parts) >= 4 {
		extraStr := strings.TrimSpace(parts[3])
		if extraStr != "" {
			deck.ExtraDeck = strings.Fields(extraStr)
		}
	}

	return deck, nil
}

// Validate checks if the deck is valid against card data
func (d *Deck) Validate(cardDB map[string]*Card) error {
	// Check hero
	hero, ok := cardDB[d.HeroID]
	if !ok {
		return fmt.Errorf("unknown hero card: %s", d.HeroID)
	}
	if !hero.IsHero() {
		return fmt.Errorf("card %s is not a hero", d.HeroID)
	}

	// Check main deck count
	if len(d.MainDeck) != 30 {
		return fmt.Errorf("main deck must have 30 cards, got %d", len(d.MainDeck))
	}

	// Check main deck cards
	mainCounts := make(map[string]int)
	for _, id := range d.MainDeck {
		card, ok := cardDB[id]
		if !ok {
			return fmt.Errorf("unknown card in main deck: %s", id)
		}
		if card.IsHero() || card.IsSkill() {
			return fmt.Errorf("card %s (%s) cannot be in main deck", id, card.Type)
		}
		mainCounts[id]++
	}

	// Check duplicates (legendary max 1, normal max 2)
	for id, count := range mainCounts {
		if count > 2 {
			return fmt.Errorf("card %s appears %d times in main deck (max 2)", id, count)
		}
	}

	// Check skill pool count
	if len(d.SkillPool) != 10 {
		return fmt.Errorf("skill pool must have 10 cards, got %d", len(d.SkillPool))
	}

	// Check skill pool cards
	skillCounts := make(map[string]int)
	for _, id := range d.SkillPool {
		card, ok := cardDB[id]
		if !ok {
			return fmt.Errorf("unknown card in skill pool: %s", id)
		}
		if !card.IsSkill() {
			return fmt.Errorf("card %s (%s) is not a skill", id, card.Type)
		}
		skillCounts[id]++
		if skillCounts[id] > 1 {
			return fmt.Errorf("skill %s appears more than once in skill pool", id)
		}
	}

	// Check extra deck cards
	for _, id := range d.ExtraDeck {
		if _, ok := cardDB[id]; !ok {
			return fmt.Errorf("unknown card in extra deck: %s", id)
		}
	}

	return nil
}

// ToDeckCode converts a Deck back to deck code string
func (d *Deck) ToDeckCode() string {
	parts := []string{
		d.HeroID,
		strings.Join(d.MainDeck, " "),
		strings.Join(d.SkillPool, " "),
	}
	if len(d.ExtraDeck) > 0 {
		parts = append(parts, strings.Join(d.ExtraDeck, " "))
	}
	return strings.Join(parts, " // ")
}
