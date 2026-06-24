package match

import (
	"eraofarcane/game"
	"testing"
	"time"
)

func TestDisconnectPlayerAndSpectatorWithSameIDAreRoleSpecific(t *testing.T) {
	room := &Room{
		ID:         "same-id",
		IsStarted:  true,
		Spectators: make(map[string]*RoomSpectator),
	}
	room.Players[0] = &RoomPlayer{
		ID:          "shared",
		Name:        "Player",
		SendFn:      func(game.GameEvent) {},
		IsConnected: true,
	}
	room.Spectators["shared"] = &RoomSpectator{
		ID:          "shared",
		Name:        "Watcher",
		SendFn:      func(game.GameEvent) {},
		IsConnected: true,
	}

	room.DisconnectPlayer("shared")

	if room.Players[0] == nil || room.Players[0].IsConnected {
		t.Fatalf("player connection should be marked disconnected")
	}
	if spectator := room.Spectators["shared"]; spectator == nil || !spectator.IsConnected || spectator.SendFn == nil {
		t.Fatalf("spectator should remain connected after player disconnect, got %+v", spectator)
	}

	room.DisconnectSpectator("shared")

	if room.Players[0] == nil || room.Players[0].SendFn != nil || room.Players[0].IsConnected {
		t.Fatalf("spectator disconnect should not reconnect or remove player, got %+v", room.Players[0])
	}
	if spectator := room.Spectators["shared"]; spectator == nil || spectator.IsConnected || spectator.SendFn != nil {
		t.Fatalf("spectator should be marked disconnected, got %+v", spectator)
	}
}

func TestDisconnectPlayerDuringStartKeepsSeat(t *testing.T) {
	room := &Room{
		ID:         "starting",
		isStarting: true,
	}
	room.Players[0] = &RoomPlayer{
		ID:          "p1",
		Name:        "Player",
		SendFn:      func(game.GameEvent) {},
		IsConnected: true,
	}

	room.DisconnectPlayer("p1")

	if room.Players[0] == nil {
		t.Fatal("player seat should be kept while game start is in progress")
	}
	if room.Players[0].IsConnected || room.Players[0].SendFn != nil {
		t.Fatalf("player should be marked disconnected during start, got %+v", room.Players[0])
	}
}

func TestJoinRoomDuringStartDoesNotReplacePlayers(t *testing.T) {
	room := &Room{
		ID:         "starting-join",
		isStarting: true,
	}
	room.Players[0] = &RoomPlayer{ID: "p1", Name: "Player1"}
	room.Players[1] = &RoomPlayer{ID: "p2", Name: "Player2"}

	if _, err := room.JoinRoom("p3", "Player3", nil, nil); err == nil {
		t.Fatal("new players should not join while game start is in progress")
	}
}

func TestSendGameEventSnapshotsBeforeSending(t *testing.T) {
	room := &Room{
		ID:         "event-send",
		IsStarted:  true,
		Engine:     game.NewEngine("event-send", nil),
		Spectators: make(map[string]*RoomSpectator),
	}
	done := make(chan struct{}, 1)
	room.Spectators["watcher"] = &RoomSpectator{
		ID:          "watcher",
		Name:        "Watcher",
		IsConnected: true,
		SendFn: func(game.GameEvent) {
			_ = room.RoomInfo()
			done <- struct{}{}
		},
	}

	room.sendGameEvent(game.GameEvent{Type: "game_start", Player: -1}, -1)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("sendGameEvent should not hold room lock while invoking SendFn")
	}
}

func TestSendGameEventSkipsSpectatorsBeforeStartPublished(t *testing.T) {
	room := &Room{
		ID:         "setup-events",
		Spectators: make(map[string]*RoomSpectator),
	}
	called := false
	room.Spectators["watcher"] = &RoomSpectator{
		ID:          "watcher",
		Name:        "Watcher",
		IsConnected: true,
		SendFn: func(game.GameEvent) {
			called = true
		},
	}

	room.sendGameEvent(game.GameEvent{Type: "phase_change", Player: -1}, -1)

	if called {
		t.Fatal("setup-time public events should not reach spectators before state_sync baseline is available")
	}
}

func TestBroadcastSpectatorStateDoesNotHoldRoomLock(t *testing.T) {
	engine := game.NewEngine("broadcast-state", nil)
	engine.State.Players[0] = &game.PlayerState{PlayerID: 0, PlayerName: "P1"}
	engine.State.Players[1] = &game.PlayerState{PlayerID: 1, PlayerName: "P2"}
	room := &Room{
		ID:         "broadcast-state",
		IsStarted:  true,
		Engine:     engine,
		Spectators: make(map[string]*RoomSpectator),
	}
	done := make(chan struct{}, 1)
	room.Spectators["watcher"] = &RoomSpectator{
		ID:          "watcher",
		Name:        "Watcher",
		IsConnected: true,
		SendFn: func(game.GameEvent) {
			_ = room.RoomInfo()
			done <- struct{}{}
		},
	}

	room.BroadcastSpectatorState()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("BroadcastSpectatorState should not hold room lock while serializing or sending")
	}
}
