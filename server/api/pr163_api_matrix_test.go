package api

import (
	"bytes"
	"encoding/json"
	"eraofarcane/game"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func pr163PostJSON(t *testing.T, serverURL, path string, body map[string]any) map[string]any {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(serverURL+path, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode %s response: %v", path, err)
	}
	if resp.StatusCode/100 != 2 {
		t.Fatalf("POST %s status=%d body=%v", path, resp.StatusCode, result)
	}
	return result
}

func pr163GetJSON(t *testing.T, serverURL, path string) map[string]any {
	t.Helper()
	resp, err := http.Get(serverURL + path)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode %s response: %v", path, err)
	}
	if resp.StatusCode/100 != 2 {
		t.Fatalf("GET %s status=%d body=%v", path, resp.StatusCode, result)
	}
	return result
}

func pr163ReadStateSync(t *testing.T, conn *websocket.Conn) map[string]any {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		conn.SetReadDeadline(deadline)
		_, raw, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read websocket event: %v", err)
		}
		var envelope map[string]any
		if err := json.Unmarshal(raw, &envelope); err != nil {
			t.Fatalf("decode websocket event: %v", err)
		}
		if envelope["type"] == "error" {
			t.Fatalf("unexpected websocket error: %v", envelope)
		}
		event, _ := envelope["event"].(map[string]any)
		if envelope["type"] == "game_event" && event["type"] == "state_sync" {
			state, ok := event["data"].(map[string]any)
			if !ok {
				t.Fatalf("state_sync data is not an object: %v", event["data"])
			}
			return state
		}
	}
}

func pr163DrainStateSync(t *testing.T, conn *websocket.Conn, count int) {
	t.Helper()
	for i := 0; i < count; i++ {
		_ = pr163ReadStateSync(t, conn)
	}
}

func pr163ReadActionResult(t *testing.T, conn *websocket.Conn, requestID string) map[string]any {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		conn.SetReadDeadline(deadline)
		_, raw, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read action result: %v", err)
		}
		var envelope map[string]any
		if err := json.Unmarshal(raw, &envelope); err != nil {
			t.Fatal(err)
		}
		if envelope["type"] == "action_result" && envelope["request_id"] == requestID {
			return envelope
		}
		if envelope["type"] == "error" {
			return envelope
		}
	}
}

func pr163CardInfos(state map[string]any, playerKey, zone string) []map[string]any {
	player, _ := state[playerKey].(map[string]any)
	items, _ := player[zone].([]any)
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if card, ok := item.(map[string]any); ok {
			result = append(result, card)
		}
	}
	return result
}

func pr163StartMatch(t *testing.T) (*httptest.Server, *websocket.Conn, *websocket.Conn, string) {
	t.Helper()
	server, rm := setupTestServer(t)
	// Mirror server/main.go: the API test helper loads card definitions but does
	// not otherwise install the lazy behavior factories used by real matches.
	game.RegisterAllCardEffects()
	room := rm.CreateTestRoom()
	c0 := connectWS(t, server, room.ID, "pr163-p0", "P0")
	if joined := readMessage(t, c0); joined["type"] != "joined" {
		t.Fatalf("p0 join: %v", joined)
	}
	c1 := connectWS(t, server, room.ID, "pr163-p1", "P1")
	if joined := readMessage(t, c1); joined["type"] != "joined" {
		t.Fatalf("p1 join: %v", joined)
	}
	// The initial event stream differs slightly by connection timing; consume one
	// authoritative state from each rather than assuming a fixed event count.
	_ = pr163ReadStateSync(t, c0)
	_ = pr163ReadStateSync(t, c1)
	sendAction(t, c0, "mulligan", map[string]any{"keep": true})
	sendAction(t, c1, "mulligan", map[string]any{"keep": true})
	_ = pr163ReadStateSync(t, c0)
	_ = pr163ReadStateSync(t, c1)
	return server, c0, c1, room.ID
}

func TestPR163APIDualPlayerStateSyncAndFixtureViews(t *testing.T) {
	server, c0, c1, roomID := pr163StartMatch(t)
	defer server.Close()
	defer c0.Close()
	defer c1.Close()

	state0 := pr163GetJSON(t, server.URL, "/api/test-room/state?room_id="+roomID+"&player=0")
	state1 := pr163GetJSON(t, server.URL, "/api/test-room/state?room_id="+roomID+"&player=1")
	if state0["phase"] != state1["phase"] {
		t.Fatalf("players saw different phases: %v vs %v", state0["phase"], state1["phase"])
	}

	// Same-name instances must remain distinguishable by instance_id in the
	// authoritative player view.
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": 0, "card_number": "1021011", "zone": "hand", "count": 2,
	})
	state0 = pr163GetJSON(t, server.URL, "/api/test-room/state?room_id="+roomID+"&player=0")
	_ = pr163ReadStateSync(t, c1)
	_ = pr163ReadStateSync(t, c0)
	var ids []string
	for _, card := range pr163CardInfos(state0, "you", "hand") {
		if card["name"] == "屠魔者杀手" {
			ids = append(ids, fmt.Sprint(card["instance_id"]))
		}
	}
	if len(ids) < 2 || ids[0] == ids[1] {
		t.Fatalf("same-name cards were not independently addressable: %v", ids)
	}

	// Rejected action from the non-active player, followed by an accepted action
	// from the active player, exercises the real request/response path.
	turn, ok := state0["current_turn"].(float64)
	if !ok {
		t.Fatalf("missing current_turn in state: %v", state0)
	}
	wrong := c1
	right := c0
	if int(turn) == 1 {
		wrong, right = c0, c1
	}
	if err := wrong.WriteJSON(map[string]any{"action": "end_turn", "data": map[string]any{}, "request_id": "pr163-wrong"}); err != nil {
		t.Fatal(err)
	}
	rejected := pr163ReadActionResult(t, wrong, "pr163-wrong")
	if rejected["ok"] != false {
		t.Fatalf("wrong-player action was not rejected: %v", rejected)
	}
	if err := right.WriteJSON(map[string]any{"action": "end_turn", "data": map[string]any{}, "request_id": "pr163-right"}); err != nil {
		t.Fatal(err)
	}
	accepted := pr163ReadActionResult(t, right, "pr163-right")
	if accepted["ok"] != true {
		t.Fatalf("active-player action was not accepted: %v", accepted)
	}
	_ = pr163ReadStateSync(t, c0)
	_ = pr163ReadStateSync(t, c1)
}

func TestPR163APIDrawMarkersAndEffectiveLoadSerialization(t *testing.T) {
	server, c0, c1, roomID := pr163StartMatch(t)
	defer server.Close()
	defer c0.Close()
	defer c1.Close()

	state := pr163GetJSON(t, server.URL, "/api/test-room/state?room_id="+roomID+"&player=0")
	current, _ := state["current_turn"].(float64)
	actor := int(current)
	actorConn := c0
	otherConn := c1
	if actor == 1 {
		actorConn, otherConn = c1, c0
	}

	// Fixture injection is explicit; the draw and pending-action resolution are
	// still performed by the real engine through two WebSocket clients.
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": actor, "card_number": "2311002", "zone": "equipment", "count": 1,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": actor, "card_number": "2321001", "zone": "equipment", "count": 1,
	})
	// Put two air cards on top and use the real draw spell to trigger both draw
	// hooks through the WebSocket action path.
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": actor, "card_number": "1321001", "zone": "deck_top", "count": 2,
	})
	// Add a draw item and give it ample resources via the test-room fixture
	// endpoint; using it remains a real player action.
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": actor, "card_number": "2021101", "zone": "hand", "count": 1,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
		"room_id": roomID, "player": actor, "mode": "set",
		"elements": map[string]any{"无": 10, "气": 10, "火": 10, "水": 10, "地": 10, "光": 10, "暗": 10},
	})
	state = pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
	// One post-mulligan broadcast remains queued, followed by five fixture
	// broadcasts (two equipment, deck, skill, elements).
	pr163DrainStateSync(t, c0, 6)
	pr163DrainStateSync(t, c1, 6)
	var drum, compass map[string]any
	for _, card := range pr163CardInfos(state, "you", "equipment") {
		switch card["name"] {
		case "唤雷震鼓":
			drum = card
		case "风息罗盘":
			compass = card
		}
	}
	if drum == nil || compass == nil {
		t.Fatalf("fixture equipment missing from player view: drum=%v compass=%v", drum, compass)
	}
	if _, ok := drum["statuses"].(map[string]any); !ok {
		t.Fatalf("drum status payload missing: %v", drum)
	}
	if _, ok := compass["elements_gain"].(map[string]any); !ok {
		t.Fatalf("compass effective load payload missing: %v", compass)
	}

	itemID := ""
	if hand, ok := state["you"].(map[string]any)["hand"].([]any); ok {
		for _, raw := range hand {
			if item, ok := raw.(map[string]any); ok && item["name"] == "失落的银叶花" {
				itemID = fmt.Sprint(item["instance_id"])
			}
		}
	}
	if itemID == "" {
		t.Fatalf("draw item fixture missing: %v", state["you"])
	}
	if err := actorConn.WriteJSON(map[string]any{
		"action": "use_item", "request_id": "pr163-draw-use",
		"data": map[string]any{"instance_id": itemID},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-draw-use"); result["ok"] != true {
		t.Fatalf("draw item use rejected: %v", result)
	}
	state = pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, otherConn)
	afterCast := pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
	t.Logf("after draw cast state: phase=%v pending=%v pending_spell=%v deck=%v hand=%v", afterCast["phase"], afterCast["pending_action"], afterCast["pending_spell"], afterCast["you"].(map[string]any)["deck_count"], len(afterCast["you"].(map[string]any)["hand"].([]any)))
	seen := []string{}
	for step := 0; step < 8; step++ {
		pending, ok := state["pending_action"].(map[string]any)
		if !ok {
			break
		}
		pendingType := fmt.Sprint(pending["type"])
		seen = append(seen, pendingType)
		candidates, _ := pending["candidates"].([]any)
		selected := []any{}
		if len(candidates) > 0 {
			candidate, _ := candidates[0].(map[string]any)
			selected = append(selected, candidate["instance_id"])
		}
		if err := actorConn.WriteJSON(map[string]any{
			"action": "resolve_action", "request_id": fmt.Sprintf("pr163-draw-%d", step),
			"data": map[string]any{"selected": selected},
		}); err != nil {
			t.Fatal(err)
		}
		if result := pr163ReadActionResult(t, actorConn, fmt.Sprintf("pr163-draw-%d", step)); result["ok"] != true {
			t.Fatalf("pending action %s rejected: %v", pendingType, result)
		}
		state = pr163ReadStateSync(t, actorConn)
		_ = pr163ReadStateSync(t, otherConn)
	}
	if len(seen) < 2 || !strings.Contains(strings.Join(seen, ","), "thunder_drum_mark") || !strings.Contains(strings.Join(seen, ","), "windbreath_compass") {
		t.Fatalf("draw action did not exercise both queued effects: %v", seen)
	}
	var finalDrum, finalCompass map[string]any
	for _, card := range pr163CardInfos(state, "you", "equipment") {
		switch card["name"] {
		case "唤雷震鼓":
			finalDrum = card
		case "风息罗盘":
			finalCompass = card
		}
	}
	if finalDrum == nil || finalCompass == nil {
		t.Fatalf("final equipment missing after draw queue: %v", state["you"])
	}
	if marks, _ := finalDrum["statuses"].(map[string]any); fmt.Sprint(marks["雷鼓标记"]) == "0" {
		t.Fatalf("thunder drum did not gain a mark after real draw effect: %v", finalDrum)
	}
	if gains, _ := finalCompass["elements_gain"].(map[string]any); fmt.Sprint(gains["气"]) == "0" {
		t.Fatalf("compass effective air load did not change after real draw effect: %v", finalCompass)
	}
}

func TestPR163APIRedMoonExtraTargetAndDefenseHandshake(t *testing.T) {
	for _, active := range []bool{false, true} {
		t.Run(map[bool]string{false: "inactive_rejects", true: "active_accepts"}[active], func(t *testing.T) {
			server, c0, c1, roomID := pr163StartMatch(t)
			defer server.Close()
			defer c0.Close()
			defer c1.Close()

			state := pr163GetJSON(t, server.URL, "/api/test-room/state?room_id="+roomID+"&player=0")
			current, _ := state["current_turn"].(float64)
			actor := int(current)
			actorConn, defenderConn := c0, c1
			if actor == 1 {
				actorConn, defenderConn = c1, c0
			}
			for _, cardNumber := range []string{"3621107", "3611101"} {
				pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
					"room_id": roomID, "player": actor, "card_number": cardNumber, "zone": "skill_slot", "count": 1,
				})
			}
			for _, pos := range [][2]int{{0, 0}, {2, 1}} {
				pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
					"room_id": roomID, "player": 1 - actor, "card_number": "1021001", "zone": "unit", "count": 1,
					"col": pos[0], "row": pos[1],
				})
			}
			pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
				"room_id": roomID, "player": actor, "mode": "set",
				"elements": map[string]any{"暗": 10, "气": 10, "无": 10, "火": 10, "水": 10, "地": 10, "光": 10},
			})

			state = pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
			// Start/mulligan leaves one state broadcast; five fixture calls follow.
			pr163DrainStateSync(t, c0, 6)
			pr163DrainStateSync(t, c1, 6)
			findSkill := func(name string) string {
				player := state["you"].(map[string]any)
				for _, raw := range player["skills"].([]any) {
					if card, ok := raw.(map[string]any); ok && card["name"] == name {
						return fmt.Sprint(card["instance_id"])
					}
				}
				return ""
			}
			willID, moonID := findSkill("意志侵蚀"), findSkill("红月")
			if willID == "" || moonID == "" {
				t.Fatalf("skills not present in fixture state: will=%q moon=%q", willID, moonID)
			}
			if active {
				pr163PostJSON(t, server.URL, "/api/test-room/card-state", map[string]any{
					"room_id": roomID, "player": actor, "instance_id": moonID,
					"statuses": map[string]any{"红月标记": 1, "异能持续": 1}, "mode": "replace",
				})
				pr163DrainStateSync(t, c0, 1)
				pr163DrainStateSync(t, c1, 1)
			}
			before := pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
			if err := actorConn.WriteJSON(map[string]any{
				"action": "cast_spell", "request_id": "pr163-red-moon-cast",
				"data": map[string]any{
					"instance_id": willID, "target_type": "unit", "target_owner": float64(1 - actor),
					"target_col": float64(0), "target_row": float64(0),
					"extra_target_col": float64(2), "extra_target_row": float64(1), "boost_ids": []any{},
				},
			}); err != nil {
				t.Fatal(err)
			}
			result := pr163ReadActionResult(t, actorConn, "pr163-red-moon-cast")
			if !active {
				if result["ok"] != false {
					t.Fatalf("inactive red moon accepted out-of-range extra target: %v", result)
				}
				after := pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
				if before["pending_spell"] != nil || after["pending_spell"] != nil {
					t.Fatalf("rejected red moon cast left pending spell: before=%v after=%v", before["pending_spell"], after["pending_spell"])
				}
				return
			}
			if result["ok"] != true {
				t.Fatalf("active red moon cast rejected: %v", result)
			}
			state = pr163ReadStateSync(t, actorConn)
			_ = pr163ReadStateSync(t, defenderConn)
			spell, ok := state["pending_spell"].(map[string]any)
			if !ok {
				t.Fatalf("accepted cast did not expose pending spell to actor: %v", state)
			}
			if spell["attacker_id"] == nil && spell["attacker"] == nil {
				t.Fatalf("pending spell lacks attacker identity: %v", spell)
			}

			// Toggle the red moon off through the fixture API while the spell is in
			// the defense window, then complete the real defense handoff.
			pr163PostJSON(t, server.URL, "/api/test-room/card-state", map[string]any{
				"room_id": roomID, "player": actor, "instance_id": moonID,
				"statuses": map[string]any{"红月标记": 1}, "mode": "replace",
			})
			pr163DrainStateSync(t, c0, 1)
			pr163DrainStateSync(t, c1, 1)
			if err := defenderConn.WriteJSON(map[string]any{"action": "no_defend", "request_id": "pr163-red-moon-no-defend", "data": map[string]any{}}); err != nil {
				t.Fatal(err)
			}
			if result := pr163ReadActionResult(t, defenderConn, "pr163-red-moon-no-defend"); result["ok"] != true {
				t.Fatalf("defense handoff rejected: %v", result)
			}
			_ = pr163ReadStateSync(t, actorConn)
			final := pr163ReadStateSync(t, defenderConn)
			if final["pending_spell"] != nil || final["phase"] == "defense_window" {
				t.Fatalf("defense handoff did not complete spell: %v", final)
			}
		})
	}
}

func TestPR163APIScarletWingsQueueRefreshOverWebSocket(t *testing.T) {
	server, c0, c1, roomID := pr163StartMatch(t)
	defer server.Close()
	defer c0.Close()
	defer c1.Close()

	state := pr163GetJSON(t, server.URL, "/api/test-room/state?room_id="+roomID+"&player=0")
	current, _ := state["current_turn"].(float64)
	actor := int(current)
	actorConn, otherConn := c0, c1
	if actor == 1 {
		actorConn, otherConn = c1, c0
	}
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": actor, "card_number": "3611101", "zone": "skill_slot", "count": 1,
	})
	for _, pos := range [][2]int{{0, 0}, {2, 0}} {
		pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
			"room_id": roomID, "player": actor, "card_number": "1621109", "zone": "unit", "count": 1,
			"col": pos[0], "row": pos[1],
		})
	}
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": 1 - actor, "card_number": "1021001", "zone": "unit", "count": 1,
		"col": 0, "row": 0,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": 1 - actor, "card_number": "1021002", "zone": "unit", "count": 1,
		"col": 2, "row": 2,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
		"room_id": roomID, "player": actor, "mode": "set",
		"elements": map[string]any{"暗": 10, "无": 10, "气": 10, "火": 10, "水": 10, "地": 10, "光": 10},
	})
	state = pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
	pr163DrainStateSync(t, c0, 7)
	pr163DrainStateSync(t, c1, 7)
	moonID := ""
	for _, raw := range state["you"].(map[string]any)["skills"].([]any) {
		if card, ok := raw.(map[string]any); ok && card["name"] == "红月" {
			moonID = fmt.Sprint(card["instance_id"])
		}
	}
	if moonID == "" {
		t.Fatalf("red moon skill missing: %v", state["you"])
	}
	if err := actorConn.WriteJSON(map[string]any{
		"action": "cast_spell", "request_id": "pr163-wings-cast",
		"data": map[string]any{"instance_id": moonID, "target_type": "none", "boost_ids": []any{}},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-wings-cast"); result["ok"] != true {
		t.Fatalf("red moon cast rejected: %v", result)
	}
	state = pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, otherConn)
	pending, ok := state["pending_action"].(map[string]any)
	if !ok || pending["type"] != "scarlet_wings_red_moon_damage" {
		t.Fatalf("red moon did not open Scarlet Wings prompt: %v", state)
	}
	var firstTarget string
	for _, raw := range pending["candidates"].([]any) {
		candidate, _ := raw.(map[string]any)
		if candidate["side"] == "enemy" && candidate["name"] == "巫师的学徒" {
			firstTarget = fmt.Sprint(candidate["instance_id"])
			break
		}
	}
	if firstTarget == "" {
		t.Fatalf("first opponent candidate missing: %v", pending["candidates"])
	}
	if err := actorConn.WriteJSON(map[string]any{
		"action": "resolve_action", "request_id": "pr163-wings-invalid",
		"data": map[string]any{"selected": []any{"not-a-candidate"}},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-wings-invalid"); result["ok"] != false {
		t.Fatalf("invalid pending selection was not rejected: %v", result)
	}
	unchanged := pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
	if unchanged["pending_action"] == nil {
		t.Fatalf("invalid pending selection consumed the prompt: %v", unchanged)
	}
	if err := actorConn.WriteJSON(map[string]any{
		"action": "resolve_action", "request_id": "pr163-wings-first",
		"data": map[string]any{"selected": []any{firstTarget}},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-wings-first"); result["ok"] != true {
		t.Fatalf("first wing target rejected: %v", result)
	}
	state = pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, otherConn)
	pending, ok = state["pending_action"].(map[string]any)
	if !ok || pending["type"] != "scarlet_wings_red_moon_damage" {
		t.Fatalf("second wing prompt missing after first target died: %v", state)
	}
	for _, raw := range pending["candidates"].([]any) {
		candidate, _ := raw.(map[string]any)
		if fmt.Sprint(candidate["instance_id"]) == firstTarget {
			t.Fatalf("dead first target remained in refreshed candidates: %v", pending["candidates"])
		}
	}
	var secondTarget string
	for _, raw := range pending["candidates"].([]any) {
		candidate, _ := raw.(map[string]any)
		if candidate["side"] == "enemy" {
			secondTarget = fmt.Sprint(candidate["instance_id"])
			break
		}
	}
	if secondTarget == "" {
		t.Fatalf("refreshed second opponent candidate missing: %v", pending["candidates"])
	}
	if err := actorConn.WriteJSON(map[string]any{
		"action": "resolve_action", "request_id": "pr163-wings-second",
		"data": map[string]any{"selected": []any{secondTarget}},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-wings-second"); result["ok"] != true {
		t.Fatalf("second wing target rejected: %v", result)
	}
	final := pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, otherConn)
	if final["pending_action"] != nil {
		t.Fatalf("Scarlet Wings queue not empty after second selection: %v", final["pending_action"])
	}
}

func TestPR163APIPaymentFailureLeavesSpellRetryable(t *testing.T) {
	server, c0, c1, roomID := pr163StartMatch(t)
	defer server.Close()
	defer c0.Close()
	defer c1.Close()

	state := pr163GetJSON(t, server.URL, "/api/test-room/state?room_id="+roomID+"&player=0")
	current, _ := state["current_turn"].(float64)
	actor := int(current)
	actorConn, defenderConn := c0, c1
	if actor == 1 {
		actorConn, defenderConn = c1, c0
	}
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": actor, "card_number": "3321001", "zone": "skill_slot", "count": 1,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": 1 - actor, "card_number": "1021001", "zone": "unit", "count": 1,
		"col": 0, "row": 0,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
		"room_id": roomID, "player": actor, "mode": "set",
		"elements": map[string]any{"气": 0, "无": 0, "暗": 0, "火": 0, "水": 0, "地": 0, "光": 0},
	})
	state = pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
	pr163DrainStateSync(t, c0, 4)
	pr163DrainStateSync(t, c1, 4)
	skillID := ""
	for _, raw := range state["you"].(map[string]any)["skills"].([]any) {
		if card, ok := raw.(map[string]any); ok && card["name"] == "闪电链" {
			skillID = fmt.Sprint(card["instance_id"])
		}
	}
	if skillID == "" {
		t.Fatalf("lightning chain fixture missing: %v", state["you"])
	}
	castData := map[string]any{
		"instance_id": skillID, "target_type": "unit", "target_owner": float64(1 - actor),
		"target_col": float64(0), "target_row": float64(0), "boost_ids": []any{},
	}
	if err := actorConn.WriteJSON(map[string]any{"action": "cast_spell", "request_id": "pr163-payment-fail", "data": castData}); err != nil {
		t.Fatal(err)
	}
	failed := pr163ReadActionResult(t, actorConn, "pr163-payment-fail")
	if failed["ok"] != false {
		t.Fatalf("insufficient-payment cast was not rejected: %v", failed)
	}
	afterFailure := pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
	if afterFailure["pending_spell"] != nil {
		t.Fatalf("failed payment left a pending spell: %v", afterFailure["pending_spell"])
	}
	for _, raw := range afterFailure["you"].(map[string]any)["skills"].([]any) {
		if card, ok := raw.(map[string]any); ok && fmt.Sprint(card["instance_id"]) == skillID && card["is_horizontal"] == true {
			t.Fatalf("failed payment consumed/horizontalized skill: %v", card)
		}
	}
	pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
		"room_id": roomID, "player": actor, "mode": "set",
		"elements": map[string]any{"气": 10, "无": 10, "暗": 10, "火": 10, "水": 10, "地": 10, "光": 10},
	})
	pr163DrainStateSync(t, c0, 1)
	pr163DrainStateSync(t, c1, 1)
	if err := actorConn.WriteJSON(map[string]any{"action": "cast_spell", "request_id": "pr163-payment-retry", "data": castData}); err != nil {
		t.Fatal(err)
	}
	if retried := pr163ReadActionResult(t, actorConn, "pr163-payment-retry"); retried["ok"] != true {
		t.Fatalf("valid payment retry rejected: %v", retried)
	}
	_ = pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, defenderConn)
	if err := defenderConn.WriteJSON(map[string]any{"action": "no_defend", "request_id": "pr163-payment-no-defend", "data": map[string]any{}}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, defenderConn, "pr163-payment-no-defend"); result["ok"] != true {
		t.Fatalf("defense after retry rejected: %v", result)
	}
	final := pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, defenderConn)
	if final["pending_spell"] != nil || final["phase"] == "defense_window" {
		t.Fatalf("successful retry did not complete defense flow: %v", final)
	}
}

func TestPR163APIReactionResponseContinuesIntoDefense(t *testing.T) {
	server, c0, c1, roomID := pr163StartMatch(t)
	defer server.Close()
	defer c0.Close()
	defer c1.Close()

	state := pr163GetJSON(t, server.URL, "/api/test-room/state?room_id="+roomID+"&player=0")
	current, _ := state["current_turn"].(float64)
	caster := int(current)
	defender := 1 - caster
	casterConn, defenderConn := c0, c1
	if caster == 1 {
		casterConn, defenderConn = c1, c0
	}

	// The fixture only prepares cards. The counter is actually set and the
	// spell is actually cast over WebSocket below.
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": caster, "card_number": "3321001", "zone": "skill_slot", "count": 1,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": defender, "card_number": "3321002", "zone": "skill_slot", "count": 1,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": defender, "card_number": "2021018", "zone": "hand", "count": 2,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": defender, "card_number": "1021001", "zone": "unit", "count": 1, "col": 0, "row": 0,
	})
	allElements := map[string]any{"气": 10, "无": 10, "暗": 10, "火": 10, "水": 10, "地": 10, "光": 10}
	pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
		"room_id": roomID, "player": caster, "mode": "set", "elements": allElements,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
		"room_id": roomID, "player": defender, "mode": "set", "elements": allElements,
	})
	pr163DrainStateSync(t, c0, 7)
	pr163DrainStateSync(t, c1, 7)

	defenderState := pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, defender))
	counterID := ""
	for _, card := range pr163CardInfos(defenderState, "you", "hand") {
		if card["name"] == "奥术符文" {
			counterID = fmt.Sprint(card["instance_id"])
			break
		}
	}
	if counterID == "" {
		t.Fatalf("counter fixture missing: %v", defenderState["you"])
	}
	if err := casterConn.WriteJSON(map[string]any{"action": "end_turn", "request_id": "pr163-counter-pass", "data": map[string]any{}}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, casterConn, "pr163-counter-pass"); result["ok"] != true {
		t.Fatalf("caster could not pass to defender: %v", result)
	}
	_ = pr163ReadStateSync(t, casterConn)
	_ = pr163ReadStateSync(t, defenderConn)
	if err := defenderConn.WriteJSON(map[string]any{
		"action": "use_item", "request_id": "pr163-counter-set", "data": map[string]any{"instance_id": counterID},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, defenderConn, "pr163-counter-set"); result["ok"] != true {
		t.Fatalf("defender could not set counter: %v", result)
	}
	_ = pr163ReadStateSync(t, casterConn)
	_ = pr163ReadStateSync(t, defenderConn)
	defenderState = pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, defender))
	secondCounterID := ""
	for _, card := range pr163CardInfos(defenderState, "you", "hand") {
		if card["name"] == "奥术符文" {
			secondCounterID = fmt.Sprint(card["instance_id"])
			break
		}
	}
	if secondCounterID == "" {
		t.Fatalf("second counter fixture missing after first set: %v", defenderState["you"])
	}
	if err := defenderConn.WriteJSON(map[string]any{
		"action": "use_item", "request_id": "pr163-counter-set-second", "data": map[string]any{"instance_id": secondCounterID},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, defenderConn, "pr163-counter-set-second"); result["ok"] != true {
		t.Fatalf("defender could not set second counter: %v", result)
	}
	_ = pr163ReadStateSync(t, casterConn)
	_ = pr163ReadStateSync(t, defenderConn)
	if err := defenderConn.WriteJSON(map[string]any{"action": "end_turn", "request_id": "pr163-counter-return", "data": map[string]any{}}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, defenderConn, "pr163-counter-return"); result["ok"] != true {
		t.Fatalf("defender could not return the turn: %v", result)
	}
	_ = pr163ReadStateSync(t, casterConn)
	_ = pr163ReadStateSync(t, defenderConn)
	pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
		"room_id": roomID, "player": caster, "mode": "set", "elements": allElements,
	})
	pr163DrainStateSync(t, c0, 1)
	pr163DrainStateSync(t, c1, 1)
	pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
		"room_id": roomID, "player": defender, "mode": "set", "elements": allElements,
	})
	pr163DrainStateSync(t, c0, 1)
	pr163DrainStateSync(t, c1, 1)

	casterState := pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, caster))
	spellID := ""
	for _, card := range pr163CardInfos(casterState, "you", "skills") {
		if card["name"] == "闪电链" {
			spellID = fmt.Sprint(card["instance_id"])
			break
		}
	}
	if spellID == "" {
		t.Fatalf("caster spell fixture missing: %v", casterState["you"])
	}
	cast := map[string]any{
		"instance_id": spellID, "target_type": "unit", "target_owner": float64(defender),
		"target_col": float64(0), "target_row": float64(0), "boost_ids": []any{},
	}
	if err := casterConn.WriteJSON(map[string]any{"action": "cast_spell", "request_id": "pr163-counter-cast", "data": cast}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, casterConn, "pr163-counter-cast"); result["ok"] != true {
		t.Fatalf("spell with counter response was rejected: %v", result)
	}
	_ = pr163ReadStateSync(t, casterConn)
	defenderState = pr163ReadStateSync(t, defenderConn)
	pending, ok := defenderState["pending_action"].(map[string]any)
	if !ok || pending["type"] != "counter_trigger" {
		t.Fatalf("counter response window missing: %v", defenderState)
	}
	if err := defenderConn.WriteJSON(map[string]any{
		"action": "resolve_action", "request_id": "pr163-counter-resolve", "data": map[string]any{"selected": []any{counterID}},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, defenderConn, "pr163-counter-resolve"); result["ok"] != true {
		t.Fatalf("counter response was rejected: %v", result)
	}
	_ = pr163ReadStateSync(t, casterConn)
	defenderState = pr163ReadStateSync(t, defenderConn)
	if pending, ok := defenderState["pending_action"].(map[string]any); ok && pending["type"] == "arcane_rune_skill" {
		candidates, _ := pending["candidates"].([]any)
		if len(candidates) == 0 {
			t.Fatalf("accepted counter did not expose its response choice: %v", pending)
		}
		responseSkill := fmt.Sprint(candidates[0].(map[string]any)["instance_id"])
		if err := defenderConn.WriteJSON(map[string]any{
			"action": "resolve_action", "request_id": "pr163-counter-response", "data": map[string]any{"selected": []any{responseSkill}},
		}); err != nil {
			t.Fatal(err)
		}
		if result := pr163ReadActionResult(t, defenderConn, "pr163-counter-response"); result["ok"] != true {
			t.Fatalf("counter response effect was rejected: %v", result)
		}
		_ = pr163ReadStateSync(t, casterConn)
		defenderState = pr163ReadStateSync(t, defenderConn)
	}
	if pending, ok := defenderState["pending_action"].(map[string]any); ok && pending["type"] == "counter_trigger" {
		if pending["min_select"] != float64(0) {
			t.Fatalf("optional counter prompt did not allow empty selection: %v", pending)
		}
		if err := defenderConn.WriteJSON(map[string]any{
			"action": "resolve_action", "request_id": "pr163-counter-decline-second", "data": map[string]any{"selected": []any{}},
		}); err != nil {
			t.Fatal(err)
		}
		if result := pr163ReadActionResult(t, defenderConn, "pr163-counter-decline-second"); result["ok"] != true {
			t.Fatalf("optional empty counter response was rejected: %v", result)
		}
		_ = pr163ReadStateSync(t, casterConn)
		defenderState = pr163ReadStateSync(t, defenderConn)
	}
	if defenderState["pending_spell"] == nil || defenderState["phase"] != "defense_window" {
		t.Fatalf("counter response did not continue into defense: %v", defenderState)
	}
	if err := defenderConn.WriteJSON(map[string]any{"action": "no_defend", "request_id": "pr163-counter-no-defend", "data": map[string]any{}}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, defenderConn, "pr163-counter-no-defend"); result["ok"] != true {
		t.Fatalf("defense after counter response was rejected: %v", result)
	}
	final := pr163ReadStateSync(t, casterConn)
	_ = pr163ReadStateSync(t, defenderConn)
	if final["pending_spell"] != nil || final["phase"] == "defense_window" {
		t.Fatalf("counter -> response -> defense flow did not settle: %v", final)
	}
}

func TestPR163APIPainScreamPromptOverWebSocket(t *testing.T) {
	server, c0, c1, roomID := pr163StartMatch(t)
	defer server.Close()
	defer c0.Close()
	defer c1.Close()

	state := pr163GetJSON(t, server.URL, "/api/test-room/state?room_id="+roomID+"&player=0")
	current, _ := state["current_turn"].(float64)
	actor := int(current)
	opponent := 1 - actor
	actorConn, opponentConn := c0, c1
	if actor == 1 {
		actorConn, opponentConn = c1, c0
	}
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": actor, "card_number": "2621106", "zone": "hand", "count": 1,
	})
	for _, slot := range []int{0, 1} {
		card := "3321001"
		if slot == 1 {
			card = "3321003"
		}
		pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
			"room_id": roomID, "player": opponent, "card_number": card, "zone": "skill_slot", "count": 1,
		})
	}
	for _, fixture := range []struct {
		card     string
		col, row int
	}{
		{card: "1121113", col: 0, row: 0},
		{card: "1121109", col: 1, row: 0},
		{card: "1021002", col: 2, row: 0},
	} {
		pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
			"room_id": roomID, "player": actor, "card_number": fixture.card, "zone": "unit", "count": 1,
			"col": fixture.col, "row": fixture.row,
		})
	}
	allElements := map[string]any{"气": 10, "无": 10, "暗": 10, "火": 10, "水": 10, "地": 10, "光": 10}
	pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
		"room_id": roomID, "player": actor, "mode": "set", "elements": allElements,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
		"room_id": roomID, "player": opponent, "mode": "set", "elements": allElements,
	})
	pr163DrainStateSync(t, c0, 9)
	pr163DrainStateSync(t, c1, 9)

	actorState := pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
	painID := ""
	for _, card := range pr163CardInfos(actorState, "you", "hand") {
		if card["name"] == "苦痛尖啸卷轴" {
			painID = fmt.Sprint(card["instance_id"])
			break
		}
	}
	if painID == "" {
		t.Fatalf("pain scream fixture missing: %v", actorState["you"])
	}
	if err := actorConn.WriteJSON(map[string]any{"action": "use_item", "request_id": "pr163-pain-use", "data": map[string]any{"instance_id": painID}}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-pain-use"); result["ok"] != true {
		t.Fatalf("pain scream use was rejected: %v", result)
	}
	_ = pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, opponentConn)
	actorState = pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
	units, _ := actorState["you"].(map[string]any)["units"].([]any)
	if len(units) < 3 {
		t.Fatalf("pain damage unit fixtures missing: %v", actorState["you"])
	}
	riderUnit, _ := units[1].([]any)
	riderCard, _ := riderUnit[0].(map[string]any)
	if riderCard["name"] != "神火兽骑手" {
		t.Fatalf("unexpected rider fixture: %v", riderCard)
	}
	if err := actorConn.WriteJSON(map[string]any{
		"action": "use_ability", "request_id": "pr163-pain-trigger", "data": map[string]any{
			"instance_id": fmt.Sprint(riderCard["instance_id"]), "ability_type": "ultimate",
		},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-pain-trigger"); result["ok"] != true {
		t.Fatalf("Divine Fire Rider trigger was rejected: %v", result)
	}
	actorState = pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, opponentConn)
	consumePrompt, ok := actorState["pending_action"].(map[string]any)
	if !ok || consumePrompt["type"] != "divine_fire_rider_consume_companion" {
		t.Fatalf("Divine Fire Rider consume prompt missing: %v", actorState)
	}
	units, _ = actorState["you"].(map[string]any)["units"].([]any)
	hellhoundUnit, _ := units[0].([]any)
	hellhoundCard, _ := hellhoundUnit[0].(map[string]any)
	hellhoundID := fmt.Sprint(hellhoundCard["instance_id"])
	if err := actorConn.WriteJSON(map[string]any{
		"action": "resolve_action", "request_id": "pr163-pain-consume-hellhound", "data": map[string]any{"selected": []any{hellhoundID}},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-pain-consume-hellhound"); result["ok"] != true {
		t.Fatalf("Hellhound consume selection was rejected: %v", result)
	}
	actorState = pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, opponentConn)
	damagePrompt, ok := actorState["pending_action"].(map[string]any)
	if !ok || damagePrompt["type"] != "lava_fort_hellhound_damage" {
		t.Fatalf("Hellhound damage prompt missing: %v", actorState)
	}
	var damagedUnits []string
	for _, raw := range damagePrompt["candidates"].([]any) {
		candidate, _ := raw.(map[string]any)
		if candidate["side"] == "own" && candidate["name"] != "熔岩堡地狱犬" {
			damagedUnits = append(damagedUnits, fmt.Sprint(candidate["instance_id"]))
			if len(damagedUnits) == 2 {
				break
			}
		}
	}
	if len(damagedUnits) != 2 {
		t.Fatalf("Hellhound did not offer two friendly targets: %v", damagePrompt)
	}
	if err := actorConn.WriteJSON(map[string]any{
		"action": "resolve_action", "request_id": "pr163-pain-damage-targets", "data": map[string]any{"selected": []any{damagedUnits[0], damagedUnits[1]}},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-pain-damage-targets"); result["ok"] != true {
		t.Fatalf("Hellhound damage targets were rejected: %v", result)
	}
	actorState = pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, opponentConn)
	pending, ok := actorState["pending_action"].(map[string]any)
	if !ok || pending["type"] != "pain_scream_weaken_enemy_spells" {
		t.Fatalf("pain scream prompt missing after friendly damage: %v", actorState)
	}
	candidates, _ := pending["candidates"].([]any)
	if len(candidates) != 2 {
		t.Fatalf("pain scream prompt should see both enemy spells: %v", pending)
	}
	if pending["min_select"] != float64(1) || pending["max_select"] != float64(1) {
		t.Fatalf("pain scream prompt should require exactly one selection: %v", pending)
	}
	if err := actorConn.WriteJSON(map[string]any{
		"action": "resolve_action", "request_id": "pr163-pain-empty", "data": map[string]any{"selected": []any{}},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-pain-empty"); result["ok"] != false {
		t.Fatalf("empty forced selection was accepted: %v", result)
	}
	unchanged := pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
	if unchanged["pending_action"] == nil {
		t.Fatalf("empty forced selection consumed the prompt: %v", unchanged)
	}
	first := fmt.Sprint(candidates[0].(map[string]any)["instance_id"])
	if err := actorConn.WriteJSON(map[string]any{
		"action": "resolve_action", "request_id": "pr163-pain-first", "data": map[string]any{"selected": []any{first}},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-pain-first"); result["ok"] != true {
		t.Fatalf("first pain scream selection was rejected: %v", result)
	}
	actorState = pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, opponentConn)
	pending, ok = actorState["pending_action"].(map[string]any)
	if !ok || pending["type"] != "pain_scream_weaken_enemy_spells" {
		t.Fatalf("queued pain scream prompt missing after first damage: %v", actorState)
	}
	candidates, _ = pending["candidates"].([]any)
	if len(candidates) != 1 || fmt.Sprint(candidates[0].(map[string]any)["instance_id"]) == first {
		t.Fatalf("queued pain scream candidates were not refreshed: %v", pending)
	}
	second := fmt.Sprint(candidates[0].(map[string]any)["instance_id"])
	if err := actorConn.WriteJSON(map[string]any{
		"action": "resolve_action", "request_id": "pr163-pain-second", "data": map[string]any{"selected": []any{second}},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-pain-second"); result["ok"] != true {
		t.Fatalf("second pain scream selection was rejected: %v", result)
	}
	final := pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, opponentConn)
	if final["pending_action"] != nil {
		t.Fatalf("pain scream prompt did not settle: %v", final["pending_action"])
	}
}

func TestPR163APIZeroCandidateItemSkipsPrompt(t *testing.T) {
	server, c0, c1, roomID := pr163StartMatch(t)
	defer server.Close()
	defer c0.Close()
	defer c1.Close()

	state := pr163GetJSON(t, server.URL, "/api/test-room/state?room_id="+roomID+"&player=0")
	current, _ := state["current_turn"].(float64)
	actor := int(current)
	actorConn, otherConn := c0, c1
	if actor == 1 {
		actorConn, otherConn = c1, c0
	}
	pr163PostJSON(t, server.URL, "/api/test-room/add-card", map[string]any{
		"room_id": roomID, "player": actor, "card_number": "2421002", "zone": "hand", "count": 1,
	})
	pr163PostJSON(t, server.URL, "/api/test-room/elements", map[string]any{
		"room_id": roomID, "player": actor, "mode": "set",
		"elements": map[string]any{"气": 10, "无": 10, "暗": 10, "火": 10, "水": 10, "地": 10, "光": 10},
	})
	pr163DrainStateSync(t, c0, 3)
	pr163DrainStateSync(t, c1, 3)

	state = pr163GetJSON(t, server.URL, fmt.Sprintf("/api/test-room/state?room_id=%s&player=%d", roomID, actor))
	scrollID := ""
	for _, card := range pr163CardInfos(state, "you", "hand") {
		if card["name"] == "生长药水" {
			scrollID = fmt.Sprint(card["instance_id"])
			break
		}
	}
	if scrollID == "" {
		t.Fatalf("zero-candidate item fixture missing: %v", state["you"])
	}
	if err := actorConn.WriteJSON(map[string]any{
		"action": "use_item", "request_id": "pr163-zero-candidate", "data": map[string]any{"instance_id": scrollID},
	}); err != nil {
		t.Fatal(err)
	}
	if result := pr163ReadActionResult(t, actorConn, "pr163-zero-candidate"); result["ok"] != true {
		t.Fatalf("zero-candidate item action was rejected: %v", result)
	}
	final := pr163ReadStateSync(t, actorConn)
	_ = pr163ReadStateSync(t, otherConn)
	if final["pending_action"] != nil {
		t.Fatalf("zero-candidate item unexpectedly opened a prompt: %v", final["pending_action"])
	}
}
