package api

import (
	"encoding/json"
	"eraofarcane/cards"
	"eraofarcane/game"
	"eraofarcane/match"
	"eraofarcane/model"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WSConnection represents a WebSocket connection for a player in a room
type WSConnection struct {
	conn     *websocket.Conn
	playerID string
	roomID   string
	slot     int
	mu       sync.Mutex
}

// SendJSON sends a JSON message through the WebSocket
func (wsc *WSConnection) SendJSON(v any) error {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()
	return wsc.conn.WriteJSON(v)
}

// HandleWebSocket handles WebSocket connections for game communication
func HandleWebSocket(rm *match.RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("room")
		playerID := r.URL.Query().Get("player_id")
		playerName := r.URL.Query().Get("player_name")
		deckCode := r.URL.Query().Get("deck_code")
		role := r.URL.Query().Get("role")
		isSpectator := role == "spectator"

		if roomID == "" || playerID == "" || playerName == "" || (!isSpectator && deckCode == "") {
			http.Error(w, "missing parameters", http.StatusBadRequest)
			return
		}

		var deck *model.Deck
		if !isSpectator {
			// Parse and validate deck
			var err error
			deck, err = model.ParseDeckCode(deckCode)
			if err != nil {
				http.Error(w, "invalid deck: "+err.Error(), http.StatusBadRequest)
				return
			}
			if err := deck.Validate(cards.PlayableCardDB); err != nil {
				http.Error(w, "invalid deck: "+err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Get or verify room
		room := rm.GetRoom(roomID)
		if room == nil {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}

		// Upgrade to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}

		wsc := &WSConnection{
			conn:     conn,
			playerID: playerID,
			roomID:   roomID,
			slot:     -1,
		}

		// Check if this is a reconnection
		isReconnect := !isSpectator && room.IsReconnection(playerID)

		defer func() {
			conn.Close()
			if isSpectator {
				room.DisconnectSpectator(playerID)
				log.Printf("Spectator %s disconnected from room %s", playerID, roomID)
			} else {
				room.DisconnectPlayer(playerID)
				log.Printf("Player %s disconnected from room %s", playerID, roomID)
			}
		}()

		// Join room (handles both new join and reconnection)
		sendFn := func(event game.GameEvent) {
			msg := map[string]any{
				"type":  "game_event",
				"event": event,
			}
			if err := wsc.SendJSON(msg); err != nil {
				log.Printf("Failed to send event to player %s: %v", playerID, err)
			}
		}

		slot := -1
		if isSpectator {
			if err := room.JoinSpectator(playerID, playerName, sendFn); err != nil {
				wsc.SendJSON(map[string]any{"type": "error", "message": err.Error()})
				return
			}
			log.Printf("Spectator %s (%s) joined room %s", playerName, playerID, roomID)
		} else {
			var err error
			slot, err = room.JoinRoom(playerID, playerName, deck, sendFn)
			if err != nil {
				wsc.SendJSON(map[string]any{"type": "error", "message": err.Error()})
				return
			}
			wsc.slot = slot

			if isReconnect {
				log.Printf("Player %s (%s) reconnected to room %s as slot %d", playerName, playerID, roomID, slot)
			} else {
				log.Printf("Player %s (%s) joined room %s as slot %d", playerName, playerID, roomID, slot)
			}
		}

		// Notify about joining
		wsc.SendJSON(map[string]any{
			"type": "joined",
			"data": map[string]any{
				"room_id":     roomID,
				"slot":        slot,
				"role":        role,
				"room":        room.RoomInfo(),
				"reconnected": isReconnect,
			},
		})

		if isSpectator {
			if room.Engine != nil {
				stateEvent := game.GameEvent{
					Type:   "state_sync",
					Player: -1,
					Data:   room.Engine.GetStateForSpectator(),
				}
				room.LogStateSync(-1)
				sendFn(stateEvent)
			}
		} else if isReconnect {
			// Send current game state immediately on reconnection
			if room.Engine != nil {
				stateEvent := game.GameEvent{
					Type:   "state_sync",
					Player: slot,
					Data:   room.Engine.GetStateForPlayer(slot),
				}
				room.LogStateSync(slot)
				sendFn(stateEvent)
			}
		} else {
			// If room is full, start the game (first time only)
			if room.IsFull() && !room.IsStarted {
				if err := room.StartGame(); err != nil {
					log.Printf("Failed to start game: %v", err)
					wsc.SendJSON(map[string]any{"type": "error", "message": err.Error()})
				} else {
					// Send initial state to connected players and spectators.
					room.BroadcastState()
				}
			}
		}

		// Read loop
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}

			var action game.ActionMessage
			if err := json.Unmarshal(message, &action); err != nil {
				room.LogMalformedClientMessage(slot, playerID, string(message), err)
				wsc.SendJSON(map[string]any{"type": "error", "message": "invalid action format"})
				continue
			}
			if action.Action == "client_log" {
				room.LogClientLog(slot, playerID, action.Data)
				continue
			}
			if isSpectator {
				room.LogActionError(slot, playerID, action, fmt.Errorf("spectators cannot act"))
				wsc.SendJSON(map[string]any{"type": "error", "message": "spectators cannot act"})
				continue
			}
			room.LogClientAction(slot, playerID, action)

			if room.Engine == nil {
				room.LogActionError(slot, playerID, action, fmt.Errorf("game not started"))
				wsc.SendJSON(map[string]any{"type": "error", "message": "game not started"})
				continue
			}

			if err := room.Engine.HandleAction(slot, action); err != nil {
				room.LogActionError(slot, playerID, action, err)
				wsc.SendJSON(map[string]any{"type": "error", "message": err.Error()})
				continue
			}

			// Send updated state after each action
			room.BroadcastState()
		}
	}
}
