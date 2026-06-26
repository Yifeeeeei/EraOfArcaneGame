package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

type outputEvent struct {
	Type    string         `json:"type"`
	Time    string         `json:"time,omitempty"`
	Message string         `json:"message,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
	Raw     any            `json:"raw,omitempty"`
}

type clientMessage struct {
	Type    string         `json:"type,omitempty"`
	Message string         `json:"message,omitempty"`
	Event   *gameEvent     `json:"event,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

type gameEvent struct {
	Type   string         `json:"type"`
	Player int            `json:"player"`
	Data   map[string]any `json:"data,omitempty"`
}

type stdinCommand struct {
	Command string         `json:"command,omitempty"`
	Action  string         `json:"action,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

type cliState struct {
	mu       sync.Mutex
	slot     any
	roomID   string
	last     map[string]any
	lastRaw  any
	lastSeen time.Time
}

var outputMu sync.Mutex

func main() {
	if err := run(); err != nil {
		writeJSON(os.Stdout, outputEvent{Type: "fatal", Time: now(), Message: err.Error()})
		os.Exit(1)
	}
}

func run() error {
	var (
		base       = flag.String("base", "http://127.0.0.1:9090", "game server base URL")
		room       = flag.String("room", "", "room id to join")
		createRoom = flag.Bool("create-room", false, "create a room before joining")
		playerID   = flag.String("player-id", "", "stable player id for this seat")
		playerName = flag.String("player-name", "Agent", "player display name")
		deckCode   = flag.String("deck-code", "", "deck code")
		deckFile   = flag.String("deck-file", "", "file containing deck code")
		tracePath  = flag.String("trace", "", "optional NDJSON trace file")
	)
	flag.Usage = usage
	flag.Parse()

	if *playerID == "" {
		return errors.New("missing --player-id")
	}
	deck, err := readDeckCode(*deckCode, *deckFile)
	if err != nil {
		return err
	}
	if deck == "" {
		return errors.New("missing --deck-code or --deck-file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	trace, err := openTrace(*tracePath)
	if err != nil {
		return err
	}
	if trace != nil {
		defer trace.Close()
	}

	baseURL, err := normalizeBaseURL(*base)
	if err != nil {
		return err
	}

	roomID := strings.TrimSpace(*room)
	if *createRoom {
		roomID, err = createGameRoom(ctx, baseURL)
		if err != nil {
			return err
		}
		writeBoth(trace, outputEvent{
			Type: "room_created",
			Time: now(),
			Data: map[string]any{"room_id": roomID, "base": baseURL.String()},
		})
	}
	if roomID == "" {
		return errors.New("missing --room or --create-room")
	}

	wsURL := websocketURL(baseURL, roomID, *playerID, *playerName, deck)
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return fmt.Errorf("connect websocket: %w", err)
	}
	defer conn.Close()

	state := &cliState{roomID: roomID}
	writeBoth(trace, outputEvent{
		Type: "connected",
		Time: now(),
		Data: map[string]any{"room_id": roomID, "websocket_url": redactDeck(wsURL)},
	})

	errCh := make(chan error, 2)
	done := make(chan struct{})
	go readWebSocket(conn, state, trace, errCh, done)
	go readStdin(conn, state, trace, errCh, done)

	select {
	case <-ctx.Done():
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "interrupted"))
		return nil
	case err := <-errCh:
		if err == nil || errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), `arcane-agent-cli connects one agent-controlled player seat.

Usage:
  go run ./cmd/arcane-agent-cli --create-room --player-id agent_p0 --deck-file deck.txt
  go run ./cmd/arcane-agent-cli --room 1234 --player-id agent_p1 --deck-code '...'

After connection, stdout emits one JSON object per line. Write one JSON object
per line to stdin:

  {"action":"mulligan","data":{"keep":true}}
  {"action":"end_turn","data":{}}
  {"command":"state"}
  {"command":"quit"}

Flags:
`)
	flag.PrintDefaults()
}

func readDeckCode(raw string, path string) (string, error) {
	if raw != "" {
		return strings.TrimSpace(raw), nil
	}
	if path == "" {
		return "", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read deck file: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func normalizeBaseURL(raw string) (*url.URL, error) {
	if raw == "" {
		return nil, errors.New("missing --base")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse --base: %w", err)
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	if u.Host == "" {
		return nil, fmt.Errorf("invalid --base %q", raw)
	}
	u.Path = strings.TrimRight(u.Path, "/")
	return u, nil
}

func createGameRoom(ctx context.Context, base *url.URL) (string, error) {
	u := *base
	u.Path = strings.TrimRight(u.Path, "/") + "/api/room/create"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("create room: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return "", fmt.Errorf("create room: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var payload struct {
		RoomID string `json:"room_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decode create room response: %w", err)
	}
	if payload.RoomID == "" {
		return "", errors.New("create room response missing room_id")
	}
	return payload.RoomID, nil
}

func websocketURL(base *url.URL, roomID string, playerID string, playerName string, deckCode string) string {
	u := *base
	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	default:
		u.Scheme = "ws"
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/ws"
	q := u.Query()
	q.Set("room", roomID)
	q.Set("player_id", playerID)
	q.Set("player_name", playerName)
	q.Set("deck_code", deckCode)
	u.RawQuery = q.Encode()
	return u.String()
}

func readWebSocket(conn *websocket.Conn, state *cliState, trace io.Writer, errCh chan<- error, done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
		}

		var msg clientMessage
		if err := conn.ReadJSON(&msg); err != nil {
			errCh <- err
			return
		}
		out := eventFromServerMessage(msg, state)
		writeBoth(trace, out)
	}
}

func eventFromServerMessage(msg clientMessage, state *cliState) outputEvent {
	switch msg.Type {
	case "joined":
		state.mu.Lock()
		if msg.Data != nil {
			state.slot = msg.Data["slot"]
			if room, ok := msg.Data["room_id"].(string); ok {
				state.roomID = room
			}
		}
		state.mu.Unlock()
		return outputEvent{Type: "joined", Time: now(), Data: msg.Data, Raw: msg}
	case "error":
		message := msg.Message
		if msg.Data != nil {
			if v, ok := msg.Data["message"].(string); ok {
				message = v
			}
		}
		return outputEvent{Type: "server_error", Time: now(), Message: message, Data: msg.Data, Raw: msg}
	case "game_event":
		if msg.Event != nil && msg.Event.Type == "state_sync" {
			state.mu.Lock()
			state.last = msg.Event.Data
			state.lastRaw = msg
			state.lastSeen = time.Now()
			state.mu.Unlock()
			return outputEvent{Type: "state", Time: now(), Data: msg.Event.Data, Raw: msg}
		}
		if msg.Event != nil {
			return outputEvent{
				Type: "game_event",
				Time: now(),
				Data: map[string]any{"event_type": msg.Event.Type, "player": msg.Event.Player, "data": msg.Event.Data},
				Raw:  msg,
			}
		}
	}
	return outputEvent{Type: "server_message", Time: now(), Raw: msg}
}

func readStdin(conn *websocket.Conn, state *cliState, trace io.Writer, errCh chan<- error, done chan<- struct{}) {
	defer close(done)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var cmd stdinCommand
		if err := json.Unmarshal([]byte(line), &cmd); err != nil {
			writeBoth(trace, outputEvent{Type: "client_error", Time: now(), Message: "invalid stdin JSON: " + err.Error()})
			continue
		}
		switch cmd.Command {
		case "state":
			writeBoth(trace, currentStateEvent(state))
			continue
		case "quit":
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "quit"))
			errCh <- nil
			return
		case "":
		default:
			writeBoth(trace, outputEvent{Type: "client_error", Time: now(), Message: "unknown command: " + cmd.Command})
			continue
		}
		if cmd.Action == "" {
			writeBoth(trace, outputEvent{Type: "client_error", Time: now(), Message: "missing action"})
			continue
		}
		if cmd.Data == nil {
			cmd.Data = map[string]any{}
		}
		payload := map[string]any{"action": cmd.Action, "data": cmd.Data}
		if err := conn.WriteJSON(payload); err != nil {
			errCh <- fmt.Errorf("send action: %w", err)
			return
		}
		writeBoth(trace, outputEvent{Type: "action_sent", Time: now(), Data: payload})
	}
	if err := scanner.Err(); err != nil {
		errCh <- err
		return
	}
	errCh <- nil
}

func currentStateEvent(state *cliState) outputEvent {
	state.mu.Lock()
	defer state.mu.Unlock()
	data := map[string]any{
		"room_id": state.roomID,
		"slot":    state.slot,
	}
	if state.last != nil {
		data["state"] = state.last
		data["state_seen_at"] = state.lastSeen.Format(time.RFC3339Nano)
	}
	return outputEvent{Type: "state_snapshot", Time: now(), Data: data, Raw: state.lastRaw}
}

func openTrace(path string) (io.WriteCloser, error) {
	if path == "" {
		return nil, nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open trace: %w", err)
	}
	return f, nil
}

func writeBoth(trace io.Writer, event outputEvent) {
	outputMu.Lock()
	defer outputMu.Unlock()
	writeJSON(os.Stdout, event)
	if trace != nil {
		writeJSON(trace, event)
	}
}

func writeJSON(w io.Writer, v any) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "encode output: %v\n", err)
		return
	}
	_, _ = w.Write(buf.Bytes())
}

func now() string {
	return time.Now().Format(time.RFC3339Nano)
}

func redactDeck(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	q := u.Query()
	if q.Has("deck_code") {
		q.Set("deck_code", "<redacted>")
	}
	u.RawQuery = q.Encode()
	return u.String()
}
