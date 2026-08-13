package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"eraofarcane/cards"
	"eraofarcane/model"
)

const (
	sourcePath = "../data/all_card_infos.json"
	targetPath = "../data/supported_card_infos.json"
)

func main() {
	allCards, err := readCards(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	supported := make([]model.Card, 0)
	for _, card := range allCards {
		if cards.IsSupportedVersion(card.VersionName) {
			supported = append(supported, card)
		}
	}
	sort.Slice(supported, func(i, j int) bool {
		left, leftErr := strconv.Atoi(supported[i].Number)
		right, rightErr := strconv.Atoi(supported[j].Number)
		if leftErr == nil && rightErr == nil {
			return left < right
		}
		return supported[i].Number < supported[j].Number
	})

	if len(supported) == 0 {
		log.Fatalf("no cards found with supported version_name in %q", strings.Join(cards.SupportedVersionNames, ", "))
	}
	if err := writeCards(targetPath, supported); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %d supported cards (%s) to %s", len(supported), strings.Join(cards.SupportedVersionNames, ", "), targetPath)
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
	return cards, nil
}

func writeCards(path string, cards []model.Card) error {
	data, err := json.MarshalIndent(cards, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
