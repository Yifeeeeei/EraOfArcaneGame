package match

import (
	"eraofarcane/game"
	"eraofarcane/model"
	"fmt"
	"log"
	"math/rand"
	"sync"
)

// Room represents a game room
type Room struct {
	ID         string                    `json:"id"`
	Players    [2]*RoomPlayer            `json:"players"`
	Spectators map[string]*RoomSpectator `json:"-"`
	Engine     *game.Engine              `json:"-"`
	IsStarted  bool                      `json:"is_started"`
	TestMode   bool                      `json:"test_mode"`
	Logger     *RoomLogger               `json:"-"`
	mu         sync.Mutex
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

// RoomSpectator represents a read-only observer connection in a room.
type RoomSpectator struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
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
		ID:         id,
		TestMode:   testMode,
		Spectators: make(map[string]*RoomSpectator),
	}
	logger, err := newRoomLogger(id, testMode)
	if err != nil {
		log.Printf("[RoomLog %s] disabled: %v", id, err)
	} else {
		room.Logger = logger
		room.LogRoomEvent("room_created", map[string]any{"test_mode": testMode})
	}
	rm.rooms[id] = room
	return room
}

// JoinSpectator adds or reconnects a read-only observer without occupying a
// player slot or requiring a deck.
func (r *Room) JoinSpectator(spectatorID, name string, sendFn func(game.GameEvent)) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Spectators == nil {
		r.Spectators = make(map[string]*RoomSpectator)
	}
	if existing := r.Spectators[spectatorID]; existing != nil {
		existing.Name = name
		existing.SendFn = sendFn
		existing.IsConnected = true
		r.LogPlayerEvent("spectator_reconnected", -1, spectatorID, map[string]any{"name": name})
		return nil
	}
	r.Spectators[spectatorID] = &RoomSpectator{
		ID:          spectatorID,
		Name:        name,
		SendFn:      sendFn,
		IsConnected: true,
	}
	r.LogPlayerEvent("spectator_joined", -1, spectatorID, map[string]any{"name": name})
	return nil
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
	if room := rm.rooms[id]; room != nil {
		room.LogRoomEvent("room_removed", nil)
		if err := room.Logger.Close(); err != nil {
			log.Printf("[RoomLog %s] close failed: %v", id, err)
		}
	}
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
				r.LogPlayerEvent("player_reconnected", i, playerID, map[string]any{"name": name})
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
	r.LogPlayerEvent("player_joined", slot, playerID, map[string]any{"name": name})

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
				r.LogPlayerEvent("player_disconnected", i, playerID, nil)
			} else {
				// Game not started: remove entirely
				r.LogPlayerEvent("player_left", i, playerID, nil)
				r.Players[i] = nil
			}
			break
		}
	}
}

// DisconnectSpectator marks a read-only observer as disconnected without
// touching player slots that may share the same browser-generated ID.
func (r *Room) DisconnectSpectator(spectatorID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if spectator := r.Spectators[spectatorID]; spectator != nil {
		spectator.SendFn = nil
		spectator.IsConnected = false
		r.LogPlayerEvent("spectator_disconnected", -1, spectatorID, nil)
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

	if !r.IsFull() {
		r.mu.Unlock()
		return fmt.Errorf("room is not full")
	}
	if r.IsStarted {
		r.mu.Unlock()
		return fmt.Errorf("game already started")
	}
	player0 := r.Players[0]
	player1 := r.Players[1]

	// Create engine with event callback
	engine := game.NewEngine(r.ID, func(event game.GameEvent, targetPlayer int) {
		r.sendGameEvent(event, targetPlayer)
	})
	r.mu.Unlock()

	// Setup game
	firstPlayer := rand.Intn(2)
	err := engine.SetupGameWithFirstPlayer(
		player0.Name, player0.Deck,
		player1.Name, player1.Deck,
		firstPlayer,
	)
	if err != nil {
		return fmt.Errorf("failed to setup game: %w", err)
	}

	r.mu.Lock()
	r.Engine = engine
	r.IsStarted = true
	r.mu.Unlock()

	r.LogRoomEvent("game_started", r.stateSnapshot())
	return nil
}

func (r *Room) sendGameEvent(event game.GameEvent, targetPlayer int) {
	r.mu.Lock()
	r.LogGameEvent(event, targetPlayer)

	sendFns := make([]func(game.GameEvent), 0, 3)
	if targetPlayer == -1 {
		for i := 0; i < 2; i++ {
			if r.Players[i] != nil && r.Players[i].SendFn != nil {
				sendFns = append(sendFns, r.Players[i].SendFn)
			}
		}
		if r.IsStarted && r.Engine != nil && isPublicSpectatorEvent(event) {
			for _, spectator := range r.Spectators {
				if spectator != nil && spectator.SendFn != nil {
					sendFns = append(sendFns, spectator.SendFn)
				}
			}
		}
	} else if targetPlayer >= 0 && targetPlayer < 2 {
		if r.Players[targetPlayer] != nil && r.Players[targetPlayer].SendFn != nil {
			sendFns = append(sendFns, r.Players[targetPlayer].SendFn)
		}
	}
	r.mu.Unlock()

	for _, sendFn := range sendFns {
		sendFn(event)
	}
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
		r.LogStateSync(i)
		r.Players[i].SendFn(game.GameEvent{
			Type:   "state_sync",
			Player: i,
			Data:   r.Engine.GetStateForPlayer(i),
		})
	}
	r.broadcastSpectatorStateLocked()
}

// BroadcastSpectatorState sends the public serialized view to connected spectators.
func (r *Room) BroadcastSpectatorState() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Engine == nil {
		return
	}
	r.broadcastSpectatorStateLocked()
}

func (r *Room) broadcastSpectatorStateLocked() {
	for _, spectator := range r.Spectators {
		if spectator == nil || spectator.SendFn == nil {
			continue
		}
		r.LogStateSync(-1)
		spectator.SendFn(game.GameEvent{
			Type:   "state_sync",
			Player: -1,
			Data:   r.Engine.GetStateForSpectator(),
		})
	}
}

// RoomInfo returns a serializable room info
func (r *Room) RoomInfo() map[string]any {
	r.mu.Lock()
	defer r.mu.Unlock()

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
	spectators := make([]any, 0)
	for _, spectator := range r.Spectators {
		if spectator != nil && spectator.IsConnected {
			spectators = append(spectators, map[string]any{
				"id":           spectator.ID,
				"name":         spectator.Name,
				"is_connected": spectator.IsConnected,
			})
		}
	}
	info["spectators"] = spectators
	info["spectator_count"] = len(spectators)
	if r.Logger != nil {
		info["log_path"] = r.Logger.Path()
	}

	return info
}

func isPublicSpectatorEvent(event game.GameEvent) bool {
	switch event.Type {
	case "game_setup", "game_start", "phase_change", "turn_start", "spell_cast", "spell_reaction",
		"defense_success", "spell_hit", "unit_attack", "unit_moved", "unit_destroyed",
		"unit_returned", "game_over", "ability_used", "defense_attempt", "discard",
		"use_item", "place_terrain", "consume", "counter_revealed", "item_cancelled",
		"spell_hit_cancelled", "card_removed_from_game", "summon":
		return true
	default:
		return false
	}
}

func (r *Room) LogRoomEvent(kind string, data map[string]any) {
	r.Logger.Write(kind, roomLogEntry{
		Data:     data,
		Snapshot: r.stateSnapshot(),
	})
}

func (r *Room) LogPlayerEvent(kind string, player int, playerID string, data map[string]any) {
	r.Logger.Write(kind, roomLogEntry{
		Player:   intPtr(player),
		PlayerID: playerID,
		Data:     data,
		Snapshot: r.stateSnapshot(),
	})
}

func (r *Room) LogClientAction(player int, playerID string, action game.ActionMessage) {
	r.Logger.Write("client_action", roomLogEntry{
		Player:   intPtr(player),
		PlayerID: playerID,
		Action:   &action,
		Snapshot: r.stateSnapshot(),
	})
}

func (r *Room) LogClientLog(player int, playerID string, data map[string]any) {
	r.Logger.Write("client_log", roomLogEntry{
		Player:   intPtr(player),
		PlayerID: playerID,
		Data:     data,
		Snapshot: r.stateSnapshot(),
	})
}

func (r *Room) LogActionError(player int, playerID string, action game.ActionMessage, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	r.Logger.Write("action_error", roomLogEntry{
		Player:   intPtr(player),
		PlayerID: playerID,
		Action:   &action,
		Error:    msg,
		Snapshot: r.stateSnapshot(),
	})
}

func (r *Room) LogMalformedClientMessage(player int, playerID string, raw string, err error) {
	data := map[string]any{"raw": raw}
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	r.Logger.Write("malformed_client_message", roomLogEntry{
		Player:   intPtr(player),
		PlayerID: playerID,
		Data:     data,
		Error:    msg,
		Snapshot: r.stateSnapshot(),
	})
}

func (r *Room) LogGameEvent(event game.GameEvent, targetPlayer int) {
	r.Logger.Write("game_event", roomLogEntry{
		Player:   intPtr(event.Player),
		Target:   intPtr(targetPlayer),
		Event:    &event,
		Snapshot: r.stateSnapshot(),
	})
}

func (r *Room) LogStateSync(player int) {
	r.Logger.Write("state_sync", roomLogEntry{
		Player:   intPtr(player),
		Snapshot: r.stateSnapshot(),
	})
}

func (r *Room) stateSnapshot() map[string]any {
	if r == nil || r.Engine == nil || r.Engine.State == nil {
		return nil
	}
	state := r.Engine.State
	snapshot := map[string]any{
		"phase":        state.Phase.String(),
		"current_turn": state.CurrentTurn,
		"first_player": state.FirstPlayer,
		"turn_number":  state.TurnNumber,
		"winner":       state.Winner,
		"mulligan":     state.MulliganDone,
	}
	if state.PendingAction != nil {
		snapshot["pending_action"] = map[string]any{
			"type":             state.PendingAction.Type,
			"player_id":        state.PendingAction.PlayerID,
			"prompt":           state.PendingAction.Prompt,
			"candidate_count":  len(state.PendingAction.Candidates),
			"min_select":       state.PendingAction.MinSelect,
			"max_select":       state.PendingAction.MaxSelect,
			"candidate_sample": state.PendingAction.Candidates,
		}
	}
	if state.PendingSpell != nil {
		snapshot["pending_spell"] = map[string]any{
			"attacker_id": state.PendingSpell.AttackerID,
			"skill":       cardBrief(state.PendingSpell.Skill),
			"target":      state.PendingSpell.Target,
			"total_power": state.PendingSpell.TotalPower,
		}
	}
	players := make([]map[string]any, 0, 2)
	for i := 0; i < 2; i++ {
		players = append(players, playerSnapshot(state.Players[i]))
	}
	snapshot["players"] = players
	return snapshot
}

func playerSnapshot(ps *game.PlayerState) map[string]any {
	if ps == nil {
		return nil
	}
	units := make([]map[string]any, 0)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if card := ps.Units[col][row]; card != nil {
				units = append(units, cardBrief(card))
			}
		}
	}
	skills := make([]map[string]any, 0, len(ps.Skills))
	for _, card := range ps.Skills {
		if card != nil {
			skills = append(skills, cardBrief(card))
		}
	}
	equipment := make([]map[string]any, 0, len(ps.Equipment))
	for _, card := range ps.Equipment {
		if card != nil {
			equipment = append(equipment, cardBrief(card))
		}
	}
	return map[string]any{
		"player_id":  ps.PlayerID,
		"name":       ps.PlayerName,
		"hero":       cardBrief(ps.Hero),
		"life":       cardLife(ps.Hero),
		"elements":   ps.Elements,
		"hand_count": len(ps.Hand),
		"deck_count": len(ps.Deck),
		"graveyard":  cardListBrief(ps.Graveyard),
		"units":      units,
		"skills":     skills,
		"equipment":  equipment,
	}
}

func cardListBrief(cards []*game.CardInstance) []map[string]any {
	items := make([]map[string]any, 0, len(cards))
	for _, card := range cards {
		items = append(items, cardBrief(card))
	}
	return items
}

func cardBrief(card *game.CardInstance) map[string]any {
	if card == nil || card.Card == nil {
		return nil
	}
	info := map[string]any{
		"instance_id":   card.InstanceID,
		"number":        card.Card.Number,
		"name":          card.Card.Name,
		"owner":         card.OwnerID,
		"life":          card.CurrentLife,
		"attack":        card.CurrentAttack,
		"horizontal":    card.IsHorizontal,
		"statuses":      card.Statuses,
		"slot":          card.SlotIndex,
		"elements_gain": card.ElementsGainBonus,
	}
	if card.Position != nil {
		info["position"] = card.Position
	}
	return info
}

func cardLife(card *game.CardInstance) int {
	if card == nil {
		return 0
	}
	return card.CurrentLife
}

func intPtr(v int) *int {
	return &v
}
