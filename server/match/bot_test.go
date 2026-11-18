package match

import (
	"eraofarcane/bot"
	"eraofarcane/cards"
	"eraofarcane/game"
	"testing"
	"time"
)

func TestSoloBotStopsAtGameOverAndOnRoomRemoval(t *testing.T) {
	if cards.CardDB == nil {
		if err := cards.LoadCards(); err != nil {
			t.Fatal(err)
		}
		game.SetCardDB(cards.PlayableCardDB)
	}
	game.RegisterAllCardEffects()
	for _, remove := range []bool{false, true} {
		t.Run(map[bool]string{false: "game_over", true: "remove"}[remove], func(t *testing.T) {
			rm := NewRoomManager()
			r, err := rm.CreateSoloRoom()
			if err != nil {
				t.Fatal(err)
			}
			defer rm.RemoveRoom(r.ID)
			deck, _ := bot.Deck()
			if _, err = r.JoinRoom("human", "Human", deck, nil); err != nil {
				t.Fatal(err)
			}
			if err = r.StartGame(); err != nil {
				t.Fatal(err)
			}
			r.mu.Lock()
			done := r.botDone
			r.mu.Unlock()
			if done == nil {
				t.Fatal("AI worker did not start")
			}
			if remove {
				rm.RemoveRoom(r.ID)
			} else {
				if err = r.EngineSnapshot().HandleAction(0, game.ActionMessage{Action: "surrender"}); err != nil {
					t.Fatal(err)
				}
			}
			select {
			case <-done:
			case <-time.After(2 * time.Second):
				t.Fatal("AI worker outlived its match")
			}
		})
	}
}
