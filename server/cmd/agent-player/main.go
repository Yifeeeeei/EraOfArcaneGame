// Command agent-player is a headless client for Codex agents playing
// EraOfArcaneGame without using the browser frontend.
package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

const sampleDeckCode = "4311003 // 1021001 1021001 1021002 1021002 1021004 1021004 1021005 1021005 1021006 1021006 1021007 1021007 1021008 1021008 1021009 1021009 1021010 1021010 1021011 1021011 1021012 1021012 1021013 1021013 1021014 1021014 1021015 1021015 1021016 1021016 // 3321002 3001001 3001002 3021001 3021002 3021003 3021004 3021005 3021006 3021007"

//go:embed templates/*.md
var dataTemplates embed.FS

type outputWriter struct {
	mu         sync.Mutex
	transcript io.Writer
}

func (w *outputWriter) write(direction string, payload []byte, show bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	envelope := struct {
		Time      string          `json:"time"`
		Direction string          `json:"direction"`
		Payload   json.RawMessage `json:"payload"`
	}{
		Time:      time.Now().UTC().Format(time.RFC3339Nano),
		Direction: direction,
		Payload:   append(json.RawMessage(nil), payload...),
	}
	line, err := json.Marshal(envelope)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode output: %v\n", err)
		return
	}
	if show {
		os.Stdout.Write(line)
		os.Stdout.Write([]byte("\n"))
	}
	if w.transcript != nil {
		w.transcript.Write(line)
		w.transcript.Write([]byte("\n"))
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "init-data":
		err = initData(os.Args[2:])
	case "create-room":
		err = createRoom(os.Args[2:])
	case "validate-deck":
		err = validateDeck(os.Args[2:])
	case "connect":
		err = connect(os.Args[2:])
	case "sample-deck":
		fmt.Println(sampleDeckCode)
	case "help", "-h", "--help":
		usage()
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "agent-player:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `Usage:
  agent-player init-data [-root ../agent-data]
  agent-player create-room [-server http://127.0.0.1:9090] [-test]
  agent-player validate-deck [-server URL] (-deck CODE | -deck-file PATH)
  agent-player connect -room ID -player-id ID -name NAME (-deck CODE | -deck-file PATH | -sample-deck)
                       [-server URL] [-transcript PATH] [-all-events]
  agent-player sample-deck

The connect command prints JSONL envelopes for decision-relevant messages and
all outgoing actions. With -transcript, every event is recorded. Use -all-events
to also print every event. Write one action per line to stdin:
  {"action":"mulligan","data":{"keep":true}}
  {"action":"end_turn","data":{}}

Lines beginning with # and blank lines are ignored. Ctrl-C closes the client.
`)
}

func initData(args []string) error {
	flags := flag.NewFlagSet("init-data", flag.ContinueOnError)
	root := flags.String("root", "../agent-data", "local ignored agent data directory")
	if err := flags.Parse(args); err != nil {
		return err
	}
	targets := map[string]string{
		"templates/core-rules.md":          "knowledge/core-rules.md",
		"templates/gameplay-principles.md": "knowledge/gameplay-principles.md",
		"templates/deck-lab.md":            "knowledge/deck-lab.md",
		"templates/open-questions.md":      "knowledge/open-questions.md",
		"templates/retired-lessons.md":     "knowledge/retired-lessons.md",
		"templates/next-match.md":          "context-packs/next-match.md",
	}
	for source, relativeTarget := range targets {
		content, err := fs.ReadFile(dataTemplates, source)
		if err != nil {
			return err
		}
		target := filepath.Join(*root, relativeTarget)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := writeFileIfMissing(target, content); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Join(*root, "matches"), 0o755); err != nil {
		return err
	}
	fmt.Printf("Agent data initialized at %s (existing files preserved)\n", *root)
	return nil
}

func writeFileIfMissing(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, os.ErrExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(content)
	return err
}

func createRoom(args []string) error {
	fs := flag.NewFlagSet("create-room", flag.ContinueOnError)
	server := fs.String("server", "http://127.0.0.1:9090", "server base URL")
	testMode := fs.Bool("test", false, "create a mutable card-test room")
	if err := fs.Parse(args); err != nil {
		return err
	}
	path := "/api/room/create"
	if *testMode {
		path = "/api/test-room/create"
	}
	resp, err := http.Post(strings.TrimRight(*server, "/")+path, "application/json", bytes.NewReader([]byte("{}")))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("server returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	fmt.Println(strings.TrimSpace(string(body)))
	return nil
}

func validateDeck(args []string) error {
	fs := flag.NewFlagSet("validate-deck", flag.ContinueOnError)
	server := fs.String("server", "http://127.0.0.1:9090", "server base URL")
	deck := fs.String("deck", "", "deck code")
	deckFile := fs.String("deck-file", "", "file containing a deck code")
	if err := fs.Parse(args); err != nil {
		return err
	}
	code, err := loadDeck(*deck, *deckFile)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"deck_code": code})
	resp, err := http.Post(strings.TrimRight(*server, "/")+"/api/deck/validate", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(strings.TrimSpace(string(responseBody)))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("deck rejected with %s", resp.Status)
	}
	return nil
}

func connect(args []string) error {
	fs := flag.NewFlagSet("connect", flag.ContinueOnError)
	server := fs.String("server", "http://127.0.0.1:9090", "server base URL")
	room := fs.String("room", "", "room ID")
	playerID := fs.String("player-id", "", "stable player ID, also used for reconnection")
	name := fs.String("name", "", "display name")
	deck := fs.String("deck", "", "deck code")
	deckFile := fs.String("deck-file", "", "file containing a deck code")
	useSampleDeck := fs.Bool("sample-deck", false, "use the built-in valid sample deck")
	transcriptPath := fs.String("transcript", "", "append JSONL messages to this file")
	allEvents := fs.Bool("all-events", false, "print every game event, not only decision-relevant messages")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *room == "" || *playerID == "" || *name == "" {
		return errors.New("-room, -player-id, and -name are required")
	}
	if *useSampleDeck {
		if *deck != "" || *deckFile != "" {
			return errors.New("-sample-deck cannot be combined with -deck or -deck-file")
		}
		*deck = sampleDeckCode
	}
	code, err := loadDeck(*deck, *deckFile)
	if err != nil {
		return err
	}
	wsURL, err := websocketURL(*server, *room, *playerID, *name, code)
	if err != nil {
		return err
	}

	var transcript *os.File
	if *transcriptPath != "" {
		transcript, err = os.OpenFile(*transcriptPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return fmt.Errorf("open transcript: %w", err)
		}
		defer transcript.Close()
	}
	out := &outputWriter{transcript: transcript}

	conn, response, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		if response != nil {
			body, _ := io.ReadAll(response.Body)
			return fmt.Errorf("connect: %w: %s", err, strings.TrimSpace(string(body)))
		}
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()

	done := make(chan error, 1)
	go func() {
		for {
			_, payload, readErr := conn.ReadMessage()
			if readErr != nil {
				done <- readErr
				return
			}
			if !json.Valid(payload) {
				payload, _ = json.Marshal(map[string]string{"unparsed": string(payload)})
			}
			out.write("received", payload, *allEvents || decisionRelevant(payload))
		}
	}()

	input := make(chan []byte)
	inputErr := make(chan error, 1)
	go scanInput(os.Stdin, input, inputErr)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	for {
		select {
		case payload, ok := <-input:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "stdin closed"))
				return nil
			}
			if err := validateAction(payload); err != nil {
				fmt.Fprintf(os.Stderr, "ignored invalid action: %v\n", err)
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return fmt.Errorf("send action: %w", err)
			}
			out.write("sent", payload, true)
		case err := <-inputErr:
			return fmt.Errorf("read stdin: %w", err)
		case err := <-done:
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return nil
			}
			return fmt.Errorf("connection closed: %w", err)
		case <-signals:
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "interrupted"))
			return nil
		}
	}
}

func decisionRelevant(payload []byte) bool {
	var message struct {
		Type  string `json:"type"`
		Event *struct {
			Type string `json:"type"`
		} `json:"event"`
	}
	if json.Unmarshal(payload, &message) != nil {
		return true
	}
	if message.Type == "joined" || message.Type == "error" {
		return true
	}
	return message.Type == "game_event" && message.Event != nil && message.Event.Type == "state_sync"
}

func scanInput(r io.Reader, output chan<- []byte, failures chan<- error) {
	defer close(output)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		output <- append([]byte(nil), line...)
	}
	if err := scanner.Err(); err != nil {
		failures <- err
	}
}

func validateAction(payload []byte) error {
	var action struct {
		Action string         `json:"action"`
		Data   map[string]any `json:"data"`
	}
	if err := json.Unmarshal(payload, &action); err != nil {
		return err
	}
	if action.Action == "" {
		return errors.New(`missing string field "action"`)
	}
	if action.Data == nil {
		return errors.New(`missing object field "data"`)
	}
	return nil
}

func loadDeck(inline, path string) (string, error) {
	if inline != "" && path != "" {
		return "", errors.New("use only one of -deck and -deck-file")
	}
	if path != "" {
		content, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		inline = strings.TrimSpace(string(content))
	}
	if strings.TrimSpace(inline) == "" {
		return "", errors.New("a deck is required via -deck or -deck-file")
	}
	return strings.TrimSpace(inline), nil
}

func websocketURL(server, room, playerID, name, deck string) (string, error) {
	parsed, err := url.Parse(server)
	if err != nil {
		return "", err
	}
	switch parsed.Scheme {
	case "http":
		parsed.Scheme = "ws"
	case "https":
		parsed.Scheme = "wss"
	case "ws", "wss":
	default:
		return "", fmt.Errorf("unsupported server URL scheme %q", parsed.Scheme)
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/ws"
	query := parsed.Query()
	query.Set("room", room)
	query.Set("player_id", playerID)
	query.Set("player_name", name)
	query.Set("deck_code", deck)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}
