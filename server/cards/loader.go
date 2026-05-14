package cards

import (
	"fmt"

	"eraofarcane/model"
)

// CardDB holds all card definitions, keyed by card number
var CardDB map[string]*model.Card

const BaseVersionName = "基础包"

// BaseCardDB holds only cards from the base set. This is the current playable
// card pool while the rest of the game is being finished.
var BaseCardDB map[string]*model.Card

// PlayableCardDB is the active card pool used by deck validation and games.
var PlayableCardDB map[string]*model.Card

// LoadCards loads all compiled card definitions.
//
// The path argument is kept for compatibility with older call sites and tools,
// but the server no longer parses JSON at runtime. Card definitions are Go
// implementations in definitions_gen.go.
func LoadCards(path string) error {
	CardDB = make(map[string]*model.Card, len(compiledCardDefinitions))
	BaseCardDB = make(map[string]*model.Card)
	for _, definition := range compiledCardDefinitions {
		card := definition.Card()
		c := normalizeCard(card)
		CardDB[c.Number] = c
		if c.VersionName == BaseVersionName {
			BaseCardDB[c.Number] = c
		}
	}
	PlayableCardDB = BaseCardDB

	fmt.Printf("Loaded %d cards (%d playable base cards)\n", len(CardDB), len(PlayableCardDB))
	return nil
}

func normalizeCard(card model.Card) *model.Card {
	if card.ElementsCost == nil {
		card.ElementsCost = make(map[string]int)
	}
	if card.ElementsGain == nil {
		card.ElementsGain = make(map[string]int)
	}
	if card.ElementsExpense == nil {
		card.ElementsExpense = make(map[string]int)
	}
	if card.Spawns == nil {
		card.Spawns = []string{}
	}
	return &card
}

// GetCard returns a card by number
func GetCard(number string) (*model.Card, bool) {
	c, ok := CardDB[number]
	return c, ok
}

// GetPlayableCard returns a card from the active playable card pool.
func GetPlayableCard(number string) (*model.Card, bool) {
	c, ok := PlayableCardDB[number]
	return c, ok
}
