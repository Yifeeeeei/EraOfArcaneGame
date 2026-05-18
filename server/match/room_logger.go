package match

import (
	"encoding/json"
	"eraofarcane/game"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var roomLogDir string

type RoomLogger struct {
	mu       sync.Mutex
	file     *os.File
	path     string
	roomID   string
	testMode bool
}

type roomLogEntry struct {
	Timestamp string              `json:"timestamp"`
	RoomID    string              `json:"room_id"`
	TestMode  bool                `json:"test_mode,omitempty"`
	Kind      string              `json:"kind"`
	Player    *int                `json:"player,omitempty"`
	PlayerID  string              `json:"player_id,omitempty"`
	Target    *int                `json:"target,omitempty"`
	Action    *game.ActionMessage `json:"action,omitempty"`
	Event     *game.GameEvent     `json:"event,omitempty"`
	Error     string              `json:"error,omitempty"`
	Data      map[string]any      `json:"data,omitempty"`
	Snapshot  map[string]any      `json:"snapshot,omitempty"`
}

func newRoomLogger(roomID string, testMode bool) (*RoomLogger, error) {
	dir := currentRoomLogDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	name := fmt.Sprintf("room-%s-%s.jsonl", roomID, time.Now().Format("20060102-150405"))
	path := filepath.Join(dir, name)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &RoomLogger{file: file, path: path, roomID: roomID, testMode: testMode}, nil
}

func currentRoomLogDir() string {
	if env := os.Getenv("ARCANE_ROOM_LOG_DIR"); env != "" {
		return env
	}
	if roomLogDir != "" {
		return roomLogDir
	}
	if dir, ok := findModuleDir(); ok {
		return filepath.Join(dir, "logs", "rooms")
	}
	return filepath.Join("logs", "rooms")
}

func findModuleDir() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func (l *RoomLogger) Path() string {
	if l == nil {
		return ""
	}
	return l.path
}

func (l *RoomLogger) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}

func (l *RoomLogger) Write(kind string, entry roomLogEntry) {
	if l == nil || l.file == nil {
		return
	}
	entry.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	entry.RoomID = l.roomID
	entry.TestMode = l.testMode
	entry.Kind = kind

	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("[RoomLog %s] marshal %s failed: %v", l.roomID, kind, err)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	if _, err := l.file.Write(append(data, '\n')); err != nil {
		log.Printf("[RoomLog %s] write %s failed: %v", l.roomID, kind, err)
	}
}
