package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"eraofarcane/model"
)

const (
	sourcePath      = "../data/supported_card_infos.json"
	definitionsPath = "cards/definitions_gen.go"
	markersPath     = "cards/category_markers_gen.go"
)

func main() {
	cards, err := readCards(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := writeFormattedGo(definitionsPath, generateDefinitions(cards)); err != nil {
		log.Fatalf("write definitions: %v", err)
	}
	if err := writeFormattedGo(markersPath, generateMarkers(cards)); err != nil {
		log.Fatalf("write markers: %v", err)
	}
	log.Printf("generated %d card definitions", len(cards))
}

func readCards(path string) ([]model.Card, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	var cards []model.Card
	if err := json.Unmarshal(data, &cards); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	sort.Slice(cards, func(i, j int) bool {
		left, leftErr := strconv.Atoi(cards[i].Number)
		right, rightErr := strconv.Atoi(cards[j].Number)
		if leftErr == nil && rightErr == nil {
			return left < right
		}
		return cards[i].Number < cards[j].Number
	})
	return cards, nil
}

func writeFormattedGo(path string, src []byte) error {
	formatted, err := format.Source(src)
	if err != nil {
		return fmt.Errorf("format %s: %w\n%s", path, err, src)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, formatted, 0644)
}

func generateDefinitions(cards []model.Card) []byte {
	var b bytes.Buffer
	b.WriteString("// Code generated from data/supported_card_infos.json; DO NOT EDIT BY HAND.\n")
	b.WriteString("package cards\n\n")
	b.WriteString("import \"eraofarcane/model\"\n\n")
	b.WriteString("// CardDefinition is the compiled, code-owned definition of a playable card.\n")
	b.WriteString("// Runtime card instances point at the model.Card returned by these definitions.\n")
	b.WriteString("type CardDefinition interface {\n")
	b.WriteString("\tID() string\n\tName() string\n\tKind() string\n\tElement() string\n\tCard() model.Card\n")
	b.WriteString("}\n\n")

	for _, card := range cards {
		typeName := cardTypeName(card.Number)
		fmt.Fprintf(&b, "type %s struct{}\n\n", typeName)
		fmt.Fprintf(&b, "func (%s) ID() string { return %q }\n", typeName, card.Number)
		fmt.Fprintf(&b, "func (%s) Name() string { return %q }\n", typeName, card.Name)
		fmt.Fprintf(&b, "func (%s) Kind() string { return %q }\n", typeName, card.Type)
		fmt.Fprintf(&b, "func (%s) Element() string { return %q }\n\n", typeName, card.Category)
		fmt.Fprintf(&b, "func (%s) Card() model.Card {\n", typeName)
		b.WriteString("\treturn model.Card{\n")
		writeStringField(&b, "Number", card.Number)
		writeStringField(&b, "Type", card.Type)
		writeStringField(&b, "Name", card.Name)
		writeStringField(&b, "Category", card.Category)
		writeStringField(&b, "Tag", card.Tag)
		writeStringField(&b, "Description", card.Description)
		writeOptionalSliceField(&b, "EffectCategories", card.EffectCategories)
		writeOptionalSliceField(&b, "EffectOptionality", card.EffectOptionality)
		writeStringField(&b, "Quote", card.Quote)
		writeMapField(&b, "ElementsCost", card.ElementsCost)
		writeMapField(&b, "ElementsGain", card.ElementsGain)
		writeMapField(&b, "ElementsExpense", card.ElementsExpense)
		writeStringField(&b, "VersionNum", card.VersionNum)
		writeStringField(&b, "VersionName", card.VersionName)
		writeIntField(&b, "Attack", card.Attack)
		writeIntField(&b, "Life", card.Life)
		writeIntField(&b, "Duration", card.Duration)
		writeIntField(&b, "Power", card.Power)
		writeSliceField(&b, "Spawns", card.Spawns)
		writeStringField(&b, "OutputPath", card.OutputPath)
		b.WriteString("\t}\n")
		b.WriteString("}\n\n")
	}

	b.WriteString("var compiledCardDefinitions = []CardDefinition{\n")
	for _, card := range cards {
		fmt.Fprintf(&b, "\t%s{},\n", cardTypeName(card.Number))
	}
	b.WriteString("}\n")
	return b.Bytes()
}

func generateMarkers(cards []model.Card) []byte {
	var b bytes.Buffer
	b.WriteString("// Code generated from data/supported_card_infos.json; DO NOT EDIT BY HAND.\n")
	b.WriteString("package cards\n\n")
	for _, card := range cards {
		typeName := cardTypeName(card.Number)
		switch card.Type {
		case model.CardTypeHero:
			fmt.Fprintf(&b, "func (%s) isHeroCard() {}\n\n", typeName)
		case model.CardTypeCompanion:
			fmt.Fprintf(&b, "func (%s) isCompanionCard() {}\n\n", typeName)
		case model.CardTypeSkill:
			fmt.Fprintf(&b, "func (%s) isSkillCard() {}\n\n", typeName)
		case model.CardTypeItem:
			fmt.Fprintf(&b, "func (%s) isItemCard() {}\n", typeName)
			if isEquipment(card) {
				fmt.Fprintf(&b, "func (%s) isEquipmentCard() {}\n", typeName)
			}
			if isWeapon(card) {
				fmt.Fprintf(&b, "func (%s) isWeaponCard() {}\n", typeName)
			}
			if isConsumable(card) {
				fmt.Fprintf(&b, "func (%s) isConsumableCard() {}\n", typeName)
			}
			if isTerrain(card) {
				fmt.Fprintf(&b, "func (%s) isTerrainCard() {}\n", typeName)
			}
			b.WriteString("\n")
		}
	}
	return b.Bytes()
}

func cardTypeName(number string) string {
	return "CardDef" + number
}

func writeStringField(b *bytes.Buffer, name, value string) {
	fmt.Fprintf(b, "\t\t%s: %q,\n", name, value)
}

func writeIntField(b *bytes.Buffer, name string, value int) {
	fmt.Fprintf(b, "\t\t%s: %d,\n", name, value)
}

func writeMapField(b *bytes.Buffer, name string, value map[string]int) {
	if len(value) == 0 {
		fmt.Fprintf(b, "\t\t%s: map[string]int{},\n", name)
		return
	}
	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	fmt.Fprintf(b, "\t\t%s: map[string]int{", name)
	for i, key := range keys {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(b, "%q: %d", key, value[key])
	}
	b.WriteString("},\n")
}

func writeSliceField(b *bytes.Buffer, name string, value []string) {
	if len(value) == 0 {
		fmt.Fprintf(b, "\t\t%s: []string{},\n", name)
		return
	}
	fmt.Fprintf(b, "\t\t%s: []string{", name)
	for i, item := range value {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(b, "%q", item)
	}
	b.WriteString("},\n")
}

func writeOptionalSliceField(b *bytes.Buffer, name string, value []string) {
	if len(value) == 0 {
		return
	}
	writeSliceField(b, name, value)
}

func isEquipment(card model.Card) bool {
	return card.Type == model.CardTypeItem && strings.Contains(card.Tag, "装备")
}

func isWeapon(card model.Card) bool {
	return isEquipment(card) && strings.Contains(card.Tag, "武器")
}

func isConsumable(card model.Card) bool {
	return card.Type == model.CardTypeItem && strings.Contains(card.Tag, "消耗品")
}

func isTerrain(card model.Card) bool {
	return card.Type == model.CardTypeItem && strings.Contains(card.Tag, "地形")
}
