package match

import (
	"eraofarcane/game"
	"testing"
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
