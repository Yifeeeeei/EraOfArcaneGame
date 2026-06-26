# Agent CLI Usage

`arcane-agent-cli` is a thin player-client wrapper for agents. It does not
choose moves, evaluate positions, build decks, or store strategy. It only keeps
one normal player WebSocket connection open and forwards actions to the game
server.

Run the game server first:

```bash
cd server
go run .
```

Run the CLI from the `server` directory:

```bash
go run ./cmd/arcane-agent-cli \
  --create-room \
  --player-id agent_p0 \
  --player-name Codex \
  --deck-file ../tmp/agent-deck.txt
```

To join an existing room:

```bash
go run ./cmd/arcane-agent-cli \
  --room 1234 \
  --player-id agent_p1 \
  --player-name Codex \
  --deck-code '4311003 // ... // ...'
```

The process writes newline-delimited JSON events to stdout:

```json
{"type":"connected","data":{"room_id":"1234"}}
{"type":"joined","data":{"slot":0,"room_id":"1234"}}
{"type":"state","data":{"phase":"mulligan"}}
{"type":"server_error","message":"not your turn"}
```

Write one JSON command per line to stdin:

```json
{"action":"mulligan","data":{"keep":true}}
{"action":"end_turn","data":{}}
{"command":"state"}
{"command":"quit"}
```

Actions are sent to the normal `/ws` player endpoint. The server remains the
only authority for legality and game state.

Optional trace logging:

```bash
go run ./cmd/arcane-agent-cli ... --trace ../tmp/agent-trace.ndjson
```

Trace files are local debugging artifacts. Do not commit traces that contain
private hand states, experimental strategy, or personal agent notes.

## Scope

The CLI intentionally handles only protocol mechanics:

- creating a normal room when requested
- connecting as one player seat
- printing `joined`, `state_sync`, game events, and server errors
- sending action JSON from stdin
- returning the latest cached state on `{"command":"state"}`

It intentionally does not:

- inspect hidden opponent information
- call test-room debug APIs
- pick moves
- build decks
- record strategic lessons
- maintain win-rate or experience databases
