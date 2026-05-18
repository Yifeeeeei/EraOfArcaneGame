package match

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRoomLoggerWritesJSONL(t *testing.T) {
	oldDir := roomLogDir
	roomLogDir = t.TempDir()
	t.Cleanup(func() { roomLogDir = oldDir })

	room := (&RoomManager{rooms: make(map[string]*Room)}).createRoom(true)
	room.LogRoomEvent("test_event", map[string]any{"ok": true})

	files, err := filepath.Glob(filepath.Join(roomLogDir, "room-"+room.ID+"-*.jsonl"))
	if err != nil {
		t.Fatalf("glob room log: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected one room log file, got %d", len(files))
	}

	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("read room log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected room_created and test_event lines, got %q", string(data))
	}

	var entry roomLogEntry
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &entry); err != nil {
		t.Fatalf("last log line should be valid json: %v", err)
	}
	if entry.Kind != "test_event" || entry.RoomID != room.ID || !entry.TestMode {
		t.Fatalf("unexpected entry: %+v", entry)
	}
}
