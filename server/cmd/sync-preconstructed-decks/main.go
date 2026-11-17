// Command sync-preconstructed-decks imports official recipes without changing them.
// Run from server/:
// go run ./cmd/sync-preconstructed-decks -revision <ArcaneDecks commit>
// An existing checkout can be read with -source ../../ArcaneDecks/index.json.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"eraofarcane/cards"
	"eraofarcane/model"
)

const repository = "https://github.com/Yifeeeeei/ArcaneDecks"

type recipe struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Version  []string `json:"version"`
	DeckCode string   `json:"deckCode"`
	Reason   string   `json:"reason,omitempty"`
}

type catalog struct {
	Source      string   `json:"source"`
	Revision    string   `json:"revision"`
	Decks       []recipe `json:"decks"`
	Unavailable []recipe `json:"unavailable"`
}

func supported(r recipe) bool {
	if len(r.Version) == 0 {
		return false
	}
	for _, version := range r.Version {
		if !cards.IsSupportedVersion(version) {
			return false
		}
	}
	return true
}

func selectRecipes(raw []byte, revision string) (catalog, error) {
	var index struct {
		Decks []recipe `json:"decks"`
	}
	if err := json.Unmarshal(raw, &index); err != nil {
		return catalog{}, err
	}
	result := catalog{Source: repository, Revision: revision, Decks: []recipe{}, Unavailable: []recipe{}}
	seen := map[string]bool{}
	for _, r := range index.Decks {
		if !supported(r) {
			continue
		}
		if r.ID == "" || r.Name == "" || seen[r.ID] {
			return catalog{}, fmt.Errorf("missing or duplicate recipe identity: %q", r.ID)
		}
		seen[r.ID] = true
		deck, err := model.ParseDeckCode(r.DeckCode)
		if err == nil {
			err = deck.Validate(cards.PlayableCardDB)
		}
		if err != nil {
			r.Reason = err.Error()
			result.Unavailable = append(result.Unavailable, r)
		} else {
			r.Reason = ""
			result.Decks = append(result.Decks, r)
		}
	}
	if len(result.Decks) == 0 {
		return catalog{}, fmt.Errorf("no supported, legal recipes found")
	}
	return result, nil
}

func run() error {
	revision := flag.String("revision", "", "exact 40-character ArcaneDecks commit (required)")
	source := flag.String("source", "", "optional local index.json; otherwise download the pinned commit")
	output := flag.String("output", "../web/assets/preconstructed-decks.json", "output snapshot")
	flag.Parse()
	if !regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(*revision) {
		return fmt.Errorf("-revision must be a full commit hash")
	}
	var raw []byte
	var err error
	if *source != "" {
		raw, err = os.ReadFile(*source)
	} else {
		client := &http.Client{Timeout: 30 * time.Second}
		var response *http.Response
		response, err = client.Get("https://raw.githubusercontent.com/Yifeeeeei/ArcaneDecks/" + *revision + "/index.json")
		if err == nil {
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("download index: %s", response.Status)
			}
			raw, err = io.ReadAll(io.LimitReader(response.Body, 4<<20))
		}
	}
	if err != nil {
		return err
	}
	if err := cards.LoadCards(); err != nil {
		return err
	}
	result, err := selectRecipes(raw, *revision)
	if err != nil {
		return err
	}
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(*output, append(encoded, '\n'), 0644); err != nil {
		return err
	}
	fmt.Printf("Imported %d legal recipes; %d need upstream corrections.\n", len(result.Decks), len(result.Unavailable))
	for _, r := range result.Unavailable {
		fmt.Printf("  %s: %s\n", r.ID, r.Reason)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
