package api

import (
	"bytes"
	"encoding/json"
	"eraofarcane/bot"
	"eraofarcane/game"
	"net/http"
	"testing"
)

func TestSoloCreateReservesOneBotAndForbidsSeatImpersonation(t *testing.T) {
	server, rm := setupTestServer(t)
	defer server.Close()
	game.RegisterAllCardEffects()
	response, err := http.Post(server.URL+"/api/solo/create", "application/json", bytes.NewBufferString(`{"profile":"royal-fire-aggro-v1"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		t.Fatalf("status %d", response.StatusCode)
	}
	var result map[string]any
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	room := rm.GetRoom(result["room_id"].(string))
	defer rm.RemoveRoom(room.ID)
	info := room.RoomInfo()
	players := info["players"].([]any)
	if len(players) != 1 || players[0].(map[string]any)["is_bot"] != true {
		t.Fatal("bot seat not reserved")
	}
	botID := players[0].(map[string]any)["id"].(string)
	if _, err = room.JoinRoom(botID, "impostor", nil, nil); err == nil {
		t.Fatal("could take over bot seat")
	}
	if room.TestMode {
		t.Fatal("single player used a test room")
	}
	if room.ShouldStart() {
		t.Fatal("started without a human")
	}
}

func TestSoloCreateSelectsThunderlightDeck(t *testing.T) {
	server, rm := setupTestServer(t)
	defer server.Close()
	response, err := http.Post(server.URL+"/api/solo/create", "application/json", bytes.NewBufferString(`{"profile":"royal-thunderlight-midrange-v1"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var result struct {
		RoomID   string `json:"room_id"`
		Profile  string `json:"profile"`
		Opponent string `json:"opponent"`
	}
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 || result.Profile != bot.ThunderlightProfileID || result.Opponent != "雷光气 · 索兰德" {
		t.Fatalf("wrong profile response: %+v", result)
	}
	room := rm.GetRoom(result.RoomID)
	defer rm.RemoveRoom(room.ID)
	if room.Players[1].Deck.HeroID != "4311101" {
		t.Fatal("thunderlight profile silently used fire deck")
	}
	for _, n := range room.Players[1].Deck.SkillPool {
		if n == "3511010" {
			t.Fatal("still uses obsolete Dawn precon")
		}
	}
}

func TestSoloCreateValidatesMethodAndProfile(t *testing.T) {
	server, rm := setupTestServer(t)
	defer server.Close()
	response, err := http.Get(server.URL + "/api/solo/create")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 405 {
		t.Fatal("GET created a room")
	}
	response, err = http.Post(server.URL+"/api/solo/create", "application/json", bytes.NewBufferString(`{"profile":"missing"}`))
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 400 {
		t.Fatal("unknown profile accepted")
	}
	if len(rm.ListRooms()) != 0 {
		t.Fatal("invalid requests left rooms behind")
	}
}
