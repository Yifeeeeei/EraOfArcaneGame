package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"eraofarcane/model"
)

const sourcePath = "../data/supported_card_infos.json"

var allowedCategories = map[string]bool{
	"主动":  true,
	"条件":  true,
	"光环":  true,
	"入场":  true,
	"遗言":  true,
	"反制":  true,
	"响应":  true,
	"祈咒":  true,
	"回合技": true,
	"绝技":  true,
}

var allowedOptionality = map[string]bool{
	"强制": true,
	"可选": true,
}

func main() {
	cards, err := readCards(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	issues := make([]string, 0)
	for _, card := range cards {
		issues = append(issues, checkCard(card)...)
	}

	if len(issues) == 0 {
		fmt.Println("card metadata check passed")
		return
	}

	for _, issue := range issues {
		fmt.Println(issue)
	}
	fmt.Printf("card metadata check found %d issue(s)\n", len(issues))
	os.Exit(1)
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
		return cards[i].Number < cards[j].Number
	})
	return cards, nil
}

func checkCard(card model.Card) []string {
	issues := make([]string, 0)
	label := fmt.Sprintf("%s %s", card.Number, card.Name)

	for _, category := range card.EffectCategories {
		if !allowedCategories[category] {
			issues = append(issues, fmt.Sprintf("%s: unknown effect category %q", label, category))
		}
	}
	for _, optionality := range card.EffectOptionality {
		if !allowedOptionality[optionality] {
			issues = append(issues, fmt.Sprintf("%s: unknown effect optionality %q", label, optionality))
		}
	}

	if strings.TrimSpace(card.Description) == "" {
		return issues
	}

	expectedCategories := expectedEffectCategories(card.Description)
	if len(expectedCategories) > 0 && len(card.EffectCategories) == 0 {
		issues = append(issues, fmt.Sprintf("%s: description suggests %s but effect_categories is empty", label, strings.Join(expectedCategories, "/")))
	}
	if len(expectedCategories) > 0 && len(card.EffectOptionality) == 0 {
		issues = append(issues, fmt.Sprintf("%s: description has effect text but effect_optionality is empty", label))
	}

	return issues
}

func expectedEffectCategories(description string) []string {
	checks := []struct {
		label   string
		needles []string
	}{
		{label: "主动", needles: []string{"主动", "回合技", "绝技", "消耗:"}},
		{label: "条件", needles: []string{"诱发", "当", "每当", "入场", "遗言", "祈咒", "游戏开始"}},
		{label: "光环", needles: []string{"光环"}},
		{label: "反制", needles: []string{"反制"}},
		{label: "响应", needles: []string{"反应"}},
	}

	found := make([]string, 0)
	for _, check := range checks {
		for _, needle := range check.needles {
			if strings.Contains(description, needle) {
				found = append(found, check.label)
				break
			}
		}
	}
	return found
}
