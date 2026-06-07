package cards

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"eraofarcane/model"
)

func TestSupportedCardInfoSnapshotMatchesPlayablePool(t *testing.T) {
	if err := LoadCards(); err != nil {
		t.Fatalf("LoadCards: %v", err)
	}

	data, err := os.ReadFile("../../data/supported_card_infos.json")
	if err != nil {
		t.Fatalf("read supported card snapshot: %v", err)
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	var snapshot []model.Card
	if err := json.Unmarshal(data, &snapshot); err != nil {
		t.Fatalf("parse supported card snapshot: %v", err)
	}

	current := PlayableCardsSorted()
	if len(snapshot) != len(current) {
		t.Fatalf("supported card snapshot has %d cards, current playable pool has %d; regenerate data/supported_card_infos.json", len(snapshot), len(current))
	}

	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	currentJSON, err := json.Marshal(current)
	if err != nil {
		t.Fatalf("marshal current playable pool: %v", err)
	}
	if string(snapshotJSON) != string(currentJSON) {
		t.Fatalf("supported card snapshot is out of date; regenerate data/supported_card_infos.json from the current playable pool")
	}
}
