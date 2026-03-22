package cards

import (
	"encoding/json"
	"fmt"
	"os"

	"eraofarcane/model"
)

// CardDB holds all card definitions, keyed by card number
var CardDB map[string]*model.Card

// LoadCards loads all card data from the JSON file
func LoadCards(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read card data: %w", err)
	}

	var cardList []model.Card
	if err := json.Unmarshal(data, &cardList); err != nil {
		return fmt.Errorf("failed to parse card data: %w", err)
	}

	CardDB = make(map[string]*model.Card, len(cardList))
	for i := range cardList {
		c := &cardList[i]
		if c.ElementsCost == nil {
			c.ElementsCost = make(map[string]int)
		}
		if c.ElementsGain == nil {
			c.ElementsGain = make(map[string]int)
		}
		if c.ElementsExpense == nil {
			c.ElementsExpense = make(map[string]int)
		}
		if c.Spawns == nil {
			c.Spawns = []string{}
		}
		CardDB[c.Number] = c
	}

	fmt.Printf("Loaded %d cards\n", len(CardDB))
	return nil
}

// GetCard returns a card by number
func GetCard(number string) (*model.Card, bool) {
	c, ok := CardDB[number]
	return c, ok
}
