package api

import (
	"encoding/json"
	"eraofarcane/game"
	"net/http"
	"strings"
	"testing"
)

func TestTutorialRoomIsPrivateAndUsesFixedDeck(t *testing.T) {
	server, rm := setupTestServer(t)
	defer server.Close()
	game.RegisterAllCardEffects()
	resp, err := http.Post(server.URL+"/api/tutorial/create", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var result struct {
		RoomID   string `json:"room_id"`
		PlayerID string `json:"player_id"`
		DeckCode string `json:"deck_code"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 || result.PlayerID == "" || result.DeckCode == "" {
		t.Fatal("missing tutorial identity")
	}
	room := rm.GetRoom(result.RoomID)
	defer rm.RemoveRoom(room.ID)
	if !room.IsTutorial() || room.TestMode {
		t.Fatal("tutorial must be private and not editable by test APIs")
	}
	if _, err = room.JoinRoom("other-player", "Other", nil, nil); err == nil {
		t.Fatal("another human joined tutorial")
	}
	if err = room.JoinSpectator("watcher", "Watcher", nil); err == nil {
		t.Fatal("tutorial allowed spectator")
	}
	list, err := http.Get(server.URL + "/api/room/list")
	if err != nil {
		t.Fatal(err)
	}
	defer list.Body.Close()
	var rooms []any
	if err = json.NewDecoder(list.Body).Decode(&rooms); err != nil || len(rooms) != 0 {
		t.Fatal("tutorial appeared in public lobby")
	}
	deck, _ := game.TutorialDecks()
	if _, err = room.JoinRoom(result.PlayerID, "Learner", deck, nil); err != nil {
		t.Fatal(err)
	}
	if err = room.StartGame(); err != nil {
		t.Fatal(err)
	}
	view := room.PlayerStateView(0)
	if view["tutorial"] == nil || view["phase"] != "mulligan" {
		t.Fatal("lesson wasn't attached to real game")
	}
	raw, _ := json.Marshal(room.RoomInfo())
	if strings.Contains(string(raw), result.PlayerID) {
		t.Fatal("room info exposed reserved identity")
	}
	if _, err = room.JoinRoom(result.PlayerID, "Learner", deck, nil); err != nil {
		t.Fatal("reserved owner could not reconnect", err)
	}
}

func TestTutorialEndpointRequiresPost(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.Close()
	resp, err := http.Get(server.URL + "/api/tutorial/create")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 405 {
		t.Fatal("expected POST-only endpoint")
	}
}
