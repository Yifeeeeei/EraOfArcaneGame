// bot-match is a local diagnostic match between two production fixed policies.
// Its seed and transcripts are private diagnostic artifacts, never HTTP APIs.
package main

import (
	"encoding/json"
	"eraofarcane/bot"
	"eraofarcane/cards"
	"eraofarcane/game"
	"eraofarcane/model"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	seed := flag.Int64("seed", 20261117, "private replay seed")
	limit := flag.Int("max-actions", 2500, "stop and fail if the match does not finish")
	out := flag.String("out", "../agent-data/bot-matches/latest", "ignored local match directory")
	gameCommit := flag.String("game-commit", "local-worktree", "exact game commit under test (append +dirty for uncommitted changes)")
	profileA := flag.String("profile-a", bot.ProfileID, "player 0 policy profile")
	profileB := flag.String("profile-b", bot.ProfileID, "player 1 policy profile")
	flag.Parse()
	if err := os.MkdirAll(*out, 0755); err != nil {
		return err
	}
	trace, err := os.Create(filepath.Join(*out, "actions.jsonl"))
	if err != nil {
		return err
	}
	defer trace.Close()
	engineLog, err := os.Create(filepath.Join(*out, "engine.log"))
	if err != nil {
		return err
	}
	defer engineLog.Close()
	log.SetOutput(engineLog)
	if err := cards.LoadCards(); err != nil {
		return err
	}
	game.SetCardDB(cards.PlayableCardDB)
	game.RegisterAllCardEffects()
	var profiles [2]bot.Profile
	var decks [2]*model.Deck
	var policies [2]*bot.Policy
	for i, id := range []string{*profileA, *profileB} {
		profile, err := bot.LookupProfile(id)
		if err != nil {
			return err
		}
		profiles[i] = profile
		decks[i], err = model.ParseDeckCode(profile.DeckCode)
		if err != nil {
			return err
		}
		if err = decks[i].Validate(cards.PlayableCardDB); err != nil {
			return err
		}
		policies[i], err = bot.NewForProfile(id)
		if err != nil {
			return err
		}
	}
	engine := game.NewEngineWithSeed("fixed-bot-mirror", nil, *seed)
	if err = engine.SetupGameWithFirstPlayer(profiles[0].Name, decks[0], profiles[1].Name, decks[1], 0); err != nil {
		return err
	}
	rejections := 0
	start := time.Now()
	encoder := json.NewEncoder(trace)
	for step := 0; step < *limit; step++ {
		acted := false
		for player := 0; player < 2; player++ {
			obs, err := engine.ObserveDecision(player)
			if err != nil {
				return err
			}
			var view map[string]any
			if err = json.Unmarshal(obs.State, &view); err != nil {
				return err
			}
			if view["phase"] == "game_over" {
				summary := map[string]any{"profiles": []string{profiles[0].ID, profiles[1].ID}, "deck_revisions": []string{profiles[0].DeckRevision, profiles[1].DeckRevision}, "seed": *seed, "actions": step, "rejections": rejections, "elapsed_ms": time.Since(start).Milliseconds(), "state": engine.GetStateForSpectator()}
				summary["game_commit"] = *gameCommit
				return writeSummary(*out, summary)
			}
			d, ok, err := policies[player].Next(obs)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			_, _, actionErr := engine.HandleActionWithRevision(player, d.Action)
			record := map[string]any{"step": step, "player": player, "turn": view["turn_number"], "phase": view["phase"], "action": d.Action, "reason": d.Reason}
			if actionErr != nil {
				rejections++
				record["error"] = actionErr.Error()
				policies[player].Reject(d)
			} else {
				policies[player].Accept()
			}
			if err = encoder.Encode(record); err != nil {
				return err
			}
			if d.Action.Action == "surrender" {
				return fmt.Errorf("AI failed a mandatory choice: %s", d.Reason)
			}
			acted = true
			break
		}
		if !acted {
			return fmt.Errorf("no policy acted at step %d; inspect %s", step, *out)
		}
	}
	return fmt.Errorf("match exceeded %d actions; inspect %s", *limit, *out)
}
func writeSummary(dir string, summary map[string]any) error {
	f, err := os.Create(filepath.Join(dir, "summary.json"))
	if err != nil {
		return err
	}
	defer f.Close()
	// Keep stdout small; the complete final state stays in the local artifact.
	data, _ := json.MarshalIndent(summary, "", "  ")
	if _, err = f.Write(data); err != nil {
		return err
	}
	compact := map[string]any{}
	for k, v := range summary {
		if k != "state" {
			compact[k] = v
		}
	}
	state := summary["state"].(map[string]any)
	compact["winner"] = state["winner"]
	compact["turns"] = state["turn_number"]
	return json.NewEncoder(os.Stdout).Encode(compact)
}
