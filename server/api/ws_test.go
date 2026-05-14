package api

import (
	"bytes"
	"encoding/json"
	"eraofarcane/cards"
	"eraofarcane/game"
	"eraofarcane/match"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

const testDeckCode = "4311003 // 1021001 1021001 1021002 1021002 1021004 1021004 1021005 1021005 1021006 1021006 1021007 1021007 1021008 1021008 1021009 1021009 1021010 1021010 1021011 1021011 1021012 1021012 1021013 1021013 1021014 1021014 1021015 1021015 1021016 1021016 // 3321002 3001001 3001002 3021001 3021002 3021003 3021004 3021005 3021006 3021007"

func setupTestServer(t *testing.T) (*httptest.Server, *match.RoomManager) {
	t.Helper()

	if cards.CardDB == nil {
		if err := cards.LoadCards(); err != nil {
			t.Fatalf("Failed to load cards: %v", err)
		}
		game.SetCardDB(cards.PlayableCardDB)
	}

	rm := match.NewRoomManager()
	mux := http.NewServeMux()
	SetupRoutes(mux, rm)

	server := httptest.NewServer(mux)
	return server, rm
}

func connectWS(t *testing.T, server *httptest.Server, roomID, playerID, playerName string) *websocket.Conn {
	t.Helper()

	import_net_url := func(s string) string {
		return strings.ReplaceAll(strings.ReplaceAll(s, " ", "%20"), "/", "%2F")
	}
	_ = import_net_url

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") +
		"/ws?room=" + roomID + "&player_id=" + playerID + "&player_name=" + playerName + "&deck_code=" + strings.ReplaceAll(testDeckCode, " ", "%20")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("WebSocket dial failed: %v", err)
	}
	return conn
}

func readMessage(t *testing.T, conn *websocket.Conn) map[string]any {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}
	var result map[string]any
	json.Unmarshal(msg, &result)
	return result
}

func sendAction(t *testing.T, conn *websocket.Conn, action string, data map[string]any) {
	t.Helper()
	msg := map[string]any{"action": action, "data": data}
	if err := conn.WriteJSON(msg); err != nil {
		t.Fatalf("Failed to send action: %v", err)
	}
}

func drainMessages(conn *websocket.Conn, count int) {
	for i := 0; i < count; i++ {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		conn.ReadMessage()
	}
}

func TestWebSocketGameFlow(t *testing.T) {
	server, rm := setupTestServer(t)
	defer server.Close()

	// Create a room
	room := rm.CreateRoom()
	roomID := room.ID
	t.Logf("Created room: %s", roomID)

	// Player 1 connects
	conn1 := connectWS(t, server, roomID, "p1", "Player1")
	defer conn1.Close()

	msg1 := readMessage(t, conn1) // "joined" message
	t.Logf("P1 joined: %v", msg1["type"])
	if msg1["type"] != "joined" {
		t.Fatalf("Expected 'joined', got %v", msg1["type"])
	}

	// Player 2 connects (triggers game start)
	conn2 := connectWS(t, server, roomID, "p2", "Player2")
	defer conn2.Close()

	msg2 := readMessage(t, conn2) // "joined" message
	t.Logf("P2 joined: %v", msg2["type"])

	// Both should receive game events (initial_draw, phase_change, state_sync)
	// Drain the initial events
	drainMessages(conn1, 10)
	drainMessages(conn2, 10)

	// Both players mulligan (keep)
	sendAction(t, conn1, "mulligan", map[string]any{"keep": true})
	time.Sleep(100 * time.Millisecond)
	drainMessages(conn1, 5)

	sendAction(t, conn2, "mulligan", map[string]any{"keep": true})
	time.Sleep(100 * time.Millisecond)

	// Drain game_start and turn_start events
	drainMessages(conn1, 10)
	drainMessages(conn2, 10)

	// Player 0 ends turn
	sendAction(t, conn1, "end_turn", map[string]any{})
	time.Sleep(100 * time.Millisecond)
	drainMessages(conn1, 10)
	drainMessages(conn2, 10)

	// Player 1 ends turn
	sendAction(t, conn2, "end_turn", map[string]any{})
	time.Sleep(100 * time.Millisecond)

	t.Log("Full WebSocket game flow test passed!")
}

func TestDeckValidationRejectsNonBaseCards(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.Close()

	reqBody := []byte(`{"deck_code":"4311003 // 1021001 1021001 1211203 1211203 1321002 1321002 1321010 1321010 1321012 1321012 1321013 1321013 1321105 1321105 1321204 1321204 1321205 1321205 1321207 1321207 1321210 1321210 2221205 2221205 2321006 2321006 2321201 2321201 2321202 2321202 // 3221209 3311201 3321002 3321007 3321015 3321106 3321206 3321207 3321208 3321209"}`)
	resp, err := http.Post(server.URL+"/api/deck/validate", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("validate request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected non-base deck to be rejected, got status %d", resp.StatusCode)
	}
}
