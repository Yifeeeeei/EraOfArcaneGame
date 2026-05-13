package cards

import "testing"

func TestLoadCardsBuildsBasePlayablePool(t *testing.T) {
	if err := LoadCards("../../data/all_card_infos.json"); err != nil {
		t.Fatalf("LoadCards: %v", err)
	}

	if len(CardDB) <= len(PlayableCardDB) {
		t.Fatalf("expected playable base pool to be smaller than full pool: full=%d playable=%d", len(CardDB), len(PlayableCardDB))
	}
	if len(PlayableCardDB) != 378 {
		t.Fatalf("expected 378 base cards, got %d", len(PlayableCardDB))
	}

	for id, card := range PlayableCardDB {
		if card.VersionName != BaseVersionName {
			t.Fatalf("playable card %s is from %s, want %s", id, card.VersionName, BaseVersionName)
		}
	}
}
