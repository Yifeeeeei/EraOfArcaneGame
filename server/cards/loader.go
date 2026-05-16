package cards

import (
	"fmt"

	"eraofarcane/model"
)

// CardDB holds all card definitions, keyed by card number
var CardDB map[string]*model.Card

// DefinitionDB holds compiled card definition objects keyed by card number.
// It is the source for interface-based category checks outside this package.
var DefinitionDB map[string]CardDefinition

const BaseVersionName = "基础包"

// BaseCardDB holds only cards from the base set. This is the current playable
// card pool while the rest of the game is being finished.
var BaseCardDB map[string]*model.Card

// PlayableCardDB is the active card pool used by deck validation and games.
var PlayableCardDB map[string]*model.Card

// LoadCards loads all compiled base-set card definitions.
func LoadCards() error {
	CardDB = make(map[string]*model.Card, len(compiledCardDefinitions))
	DefinitionDB = make(map[string]CardDefinition, len(compiledCardDefinitions))
	BaseCardDB = make(map[string]*model.Card)
	for _, definition := range compiledCardDefinitions {
		card := definition.Card()
		c := normalizeCard(card)
		CardDB[c.Number] = c
		DefinitionDB[c.Number] = definition
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

// GetDefinition returns the compiled Go definition for a card number.
func GetDefinition(number string) (CardDefinition, bool) {
	definition, ok := DefinitionDB[number]
	return definition, ok
}

func IsHero(number string) bool {
	definition, ok := GetDefinition(number)
	if !ok {
		return false
	}
	_, ok = definition.(HeroCard)
	return ok
}

func IsCompanion(number string) bool {
	definition, ok := GetDefinition(number)
	if !ok {
		return false
	}
	_, ok = definition.(CompanionCard)
	return ok
}

func IsSkill(number string) bool {
	definition, ok := GetDefinition(number)
	if !ok {
		return false
	}
	_, ok = definition.(SkillCard)
	return ok
}

func IsItem(number string) bool {
	definition, ok := GetDefinition(number)
	if !ok {
		return false
	}
	_, ok = definition.(ItemCard)
	return ok
}

func IsEquipment(number string) bool {
	definition, ok := GetDefinition(number)
	if !ok {
		return false
	}
	_, ok = definition.(EquipmentCard)
	return ok
}

func IsWeapon(number string) bool {
	definition, ok := GetDefinition(number)
	if !ok {
		return false
	}
	_, ok = definition.(WeaponCard)
	return ok
}

func IsConsumable(number string) bool {
	definition, ok := GetDefinition(number)
	if !ok {
		return false
	}
	_, ok = definition.(ConsumableCard)
	return ok
}

func IsTerrain(number string) bool {
	definition, ok := GetDefinition(number)
	if !ok {
		return false
	}
	_, ok = definition.(TerrainCard)
	return ok
}
