package cards

import "testing"

func TestLoadCardsBuildsBasePlayablePool(t *testing.T) {
	if err := LoadCards(); err != nil {
		t.Fatalf("LoadCards: %v", err)
	}

	if len(CardDB) != 379 {
		t.Fatalf("expected 379 compiled base cards, got %d", len(CardDB))
	}
	if len(PlayableCardDB) != 379 {
		t.Fatalf("expected 379 base cards, got %d", len(PlayableCardDB))
	}

	for id, card := range PlayableCardDB {
		if card.VersionName != BaseVersionName {
			t.Fatalf("playable card %s is from %s, want %s", id, card.VersionName, BaseVersionName)
		}
	}
}
