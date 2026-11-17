package main

import (
	"encoding/json"
	"os"
	"reflect"
	"regexp"
	"testing"

	"eraofarcane/cards"
	"eraofarcane/model"
)

func TestPackagedRecipesMatchCurrentPlayablePool(t *testing.T) {
	if err := cards.LoadCards(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile("../../../web/assets/preconstructed-decks.json")
	if err != nil {
		t.Fatal(err)
	}
	var snapshot catalog
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		t.Fatal(err)
	}
	if snapshot.Source != repository || !regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(snapshot.Revision) {
		t.Fatal("catalog must identify its exact official source commit")
	}
	if len(snapshot.Decks) == 0 {
		t.Fatal("empty player catalog")
	}
	seen := map[string]bool{}
	for _, r := range append(append([]recipe{}, snapshot.Decks...), snapshot.Unavailable...) {
		if r.ID == "" || r.Name == "" || seen[r.ID] || !supported(r) {
			t.Fatalf("invalid recipe identity or unsupported version: %+v", r)
		}
		seen[r.ID] = true
		deck, err := model.ParseDeckCode(r.DeckCode)
		if err == nil {
			err = deck.Validate(cards.PlayableCardDB)
		}
		if r.Reason == "" && err != nil {
			t.Errorf("player recipe %s is no longer legal: %v", r.ID, err)
		}
		if r.Reason != "" && err == nil {
			t.Errorf("recipe %s is now legal; re-sync to make it available", r.ID)
		}
	}
	// Unsupported packs are excluded by metadata, even when they reuse legal IDs.
	other := snapshot.Decks[0]
	other.ID = "future/test"
	other.Version = []string{"future"}
	input := append(append([]recipe{}, snapshot.Decks...), snapshot.Unavailable...)
	input = append(input, other)
	encoded, err := json.Marshal(struct {
		Decks []recipe `json:"decks"`
	}{input})
	if err != nil {
		t.Fatal(err)
	}
	regenerated, err := selectRecipes(encoded, snapshot.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(regenerated, snapshot) {
		t.Fatal("sync must preserve official recipes and keep rejected decks out of the playable list")
	}
}
