package match

import (
	"eraofarcane/game"
	"eraofarcane/model"
	"fmt"
	"math/rand"
	"sync"
)

// Room represents a game room
type Room struct {
	ID        string         `json:"id"`
	Players   [2]*RoomPlayer `json:"players"`
	Engine    *game.Engine   `json:"-"`
	IsStarted bool           `json:"is_started"`
	TestMode  bool           `json:"test_mode"`
	mu        sync.Mutex
}

// RoomPlayer represents a player in a room
type RoomPlayer struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Deck        *model.Deck                `json:"deck"`
	Ready       bool                       `json:"ready"`
	SendFn      func(event game.GameEvent) `json:"-"`
	IsConnected bool                       `json:"is_connected"`
}

// RoomManager manages all game rooms
type RoomManager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

// NewRoomManager creates a new room manager
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

// CreateRoom creates a new room
func (rm *RoomManager) CreateRoom() *Room {
	return rm.createRoom(false)
}

// CreateTestRoom creates a room whose game state may be edited through debug APIs.
func (rm *RoomManager) CreateTestRoom() *Room {
	return rm.createRoom(true)
}

func (rm *RoomManager) createRoom(testMode bool) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	id := fmt.Sprintf("%04d", rand.Intn(10000))
	// Ensure unique
	for _, exists := rm.rooms[id]; exists; _, exists = rm.rooms[id] {
		id = fmt.Sprintf("%04d", rand.Intn(10000))
	}

	room := &Room{
		ID:       id,
		TestMode: testMode,
	}
	rm.rooms[id] = room
	return room
}

// GetRoom returns a room by ID
func (rm *RoomManager) GetRoom(id string) *Room {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.rooms[id]
}

// RemoveRoom removes a room
func (rm *RoomManager) RemoveRoom(id string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.rooms, id)
}

// ListRooms returns all rooms
func (rm *RoomManager) ListRooms() []*Room {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	rooms := make([]*Room, 0, len(rm.rooms))
	for _, r := range rm.rooms {
		rooms = append(rooms, r)
	}
	return rooms
}

// JoinRoom adds a player to a room. If the game is already started and
// the same player_id reconnects, it updates the SendFn (reconnection).
func (r *Room) JoinRoom(playerID, name string, deck *model.Deck, sendFn func(game.GameEvent)) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if this player is reconnecting to a started game
	if r.IsStarted {
		for i := 0; i < 2; i++ {
			if r.Players[i] != nil && r.Players[i].ID == playerID {
				// Reconnection: update SendFn and mark connected
				r.Players[i].SendFn = sendFn
				r.Players[i].IsConnected = true
				return i, nil
			}
		}
		return -1, fmt.Errorf("game already started, cannot join as new player")
	}

	// Normal join for games not yet started
	slot := -1
	if r.Players[0] == nil {
		slot = 0
	} else if r.Players[1] == nil {
		slot = 1
	} else {
		return -1, fmt.Errorf("room is full")
	}

	r.Players[slot] = &RoomPlayer{
		ID:          playerID,
		Name:        name,
		Deck:        deck,
		SendFn:      sendFn,
		IsConnected: true,
	}

	return slot, nil
}

// DisconnectPlayer marks a player as disconnected but keeps them in the room
// if the game is started (for reconnection support).
// If the game hasn't started, removes the player entirely.
func (r *Room) DisconnectPlayer(playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := 0; i < 2; i++ {
		if r.Players[i] != nil && r.Players[i].ID == playerID {
			if r.IsStarted {
				// Game in progress: keep player data, just disconnect
				r.Players[i].SendFn = nil
				r.Players[i].IsConnected = false
			} else {
				// Game not started: remove entirely
				r.Players[i] = nil
			}
			break
		}
	}
}

// LeaveRoom removes a player from a room (kept for compatibility)
func (r *Room) LeaveRoom(playerID string) {
	r.DisconnectPlayer(playerID)
}

// GetPlayerSlot returns the slot index for a player
func (r *Room) GetPlayerSlot(playerID string) int {
	for i := 0; i < 2; i++ {
		if r.Players[i] != nil && r.Players[i].ID == playerID {
			return i
		}
	}
	return -1
}

// IsFull returns true if both player slots are occupied
func (r *Room) IsFull() bool {
	return r.Players[0] != nil && r.Players[1] != nil
}

// IsReconnection checks if a player is reconnecting to a started game
func (r *Room) IsReconnection(playerID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.IsStarted {
		return false
	}
	for i := 0; i < 2; i++ {
		if r.Players[i] != nil && r.Players[i].ID == playerID {
			return true
		}
	}
	return false
}

// StartGame initializes and starts the game
func (r *Room) StartGame() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.IsFull() {
		return fmt.Errorf("room is not full")
	}
	if r.IsStarted {
		return fmt.Errorf("game already started")
	}

	// Create engine with event callback
	r.Engine = game.NewEngine(r.ID, func(event game.GameEvent, targetPlayer int) {
		if targetPlayer == -1 {
			// Send to both players
			for i := 0; i < 2; i++ {
				if r.Players[i] != nil && r.Players[i].SendFn != nil {
					r.Players[i].SendFn(event)
				}
			}
		} else if targetPlayer >= 0 && targetPlayer < 2 {
			if r.Players[targetPlayer] != nil && r.Players[targetPlayer].SendFn != nil {
				r.Players[targetPlayer].SendFn(event)
			}
		}
	})

	// Setup game
	err := r.Engine.SetupGame(
		r.Players[0].Name, r.Players[0].Deck,
		r.Players[1].Name, r.Players[1].Deck,
	)
	if err != nil {
		return fmt.Errorf("failed to setup game: %w", err)
	}

	r.IsStarted = true
	return nil
}

// BroadcastState sends each connected player their own serialized view.
func (r *Room) BroadcastState() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Engine == nil {
		return
	}
	for i := 0; i < 2; i++ {
		if r.Players[i] == nil || r.Players[i].SendFn == nil {
			continue
		}
		r.Players[i].SendFn(game.GameEvent{
			Type:   "state_sync",
			Player: i,
			Data:   r.Engine.GetStateForPlayer(i),
		})
	}
}

// RoomInfo returns a serializable room info
func (r *Room) RoomInfo() map[string]any {
	info := map[string]any{
		"id":         r.ID,
		"is_started": r.IsStarted,
		"test_mode":  r.TestMode,
	}

	players := make([]any, 0)
	for i := 0; i < 2; i++ {
		if r.Players[i] != nil {
			players = append(players, map[string]any{
				"id":           r.Players[i].ID,
				"name":         r.Players[i].Name,
				"slot":         i,
				"is_connected": r.Players[i].IsConnected,
			})
		}
	}
	info["players"] = players

	return info
}
