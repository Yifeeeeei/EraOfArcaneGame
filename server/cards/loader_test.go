package cards

import "testing"

func TestLoadCardsBuildsSupportedPlayablePool(t *testing.T) {
	if err := LoadCards(); err != nil {
		t.Fatalf("LoadCards: %v", err)
	}

	const expectedBaseCards = 393
	const expectedRoyalConflictCards = 335
	expectedSupportedCards := expectedBaseCards + expectedRoyalConflictCards

	if len(CardDB) != expectedSupportedCards {
		t.Fatalf("expected %d compiled supported cards, got %d", expectedSupportedCards, len(CardDB))
	}
	if len(BaseCardDB) != expectedBaseCards {
		t.Fatalf("expected %d base cards, got %d", expectedBaseCards, len(BaseCardDB))
	}
	if len(PlayableCardDB) != expectedSupportedCards {
		t.Fatalf("expected %d playable supported cards, got %d", expectedSupportedCards, len(PlayableCardDB))
	}

	for id, card := range PlayableCardDB {
		if !IsSupportedVersion(card.VersionName) {
			t.Fatalf("playable card %s is from unsupported version %s", id, card.VersionName)
		}
	}
}
