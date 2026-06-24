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

func TestSendGameEventSnapshotsBeforeSending(t *testing.T) {
	room := &Room{
		ID:         "event-send",
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
