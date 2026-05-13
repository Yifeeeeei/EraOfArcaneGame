package cards

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"eraofarcane/model"
)

const SupportedCardInfoSnapshotPath = "../data/supported_card_infos.json"

// PlayableCardsSorted returns the active playable card pool as a stable,
// number-sorted card list for snapshots and comparisons.
func PlayableCardsSorted() []model.Card {
	cards := make([]model.Card, 0, len(PlayableCardDB))
	for _, card := range PlayableCardDB {
		cards = append(cards, *card)
	}
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Number < cards[j].Number
	})
	return cards
}

// WriteSupportedCardInfoSnapshot writes the current playable card pool in the
// same schema as all_card_infos.json. It is intended as a balance/version
// snapshot for future diffs against newer source card data.
func WriteSupportedCardInfoSnapshot(path string) error {
	if PlayableCardDB == nil {
		return fmt.Errorf("playable card database is not loaded")
	}

	data, err := json.MarshalIndent(PlayableCardsSorted(), "", "  ")
	if err != nil {
		return fmt.Errorf("marshal supported card snapshot: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write supported card snapshot: %w", err)
	}
	return nil
}
