# Agent Match Protocol Playbook

This playbook explains how an agent can operate Era of Arcane as a normal
player. Prefer `server/cmd/arcane-agent-cli` for day-to-day use; it wraps the
same protocol described here. Agents may also implement this protocol directly.

Do not add user-facing agent controls to the browser UI. Do not use debug APIs
for a fair match.

## Run The Server

Run from the `server` directory:

```bash
go run .
```

The local server listens on:

```text
http://127.0.0.1:9090
```

## Validate A Deck

Use the normal deck validation API before joining:

```text
POST /api/deck/validate
Content-Type: application/json

{"deck_code":"4311003 // ... // ..."}
```

The response confirms the card pool and deck counts. A valid deck is still
checked again when joining `/ws`.

## Room APIs

Create a room:

```text
POST /api/room/create
```

Response:

```json
{"room_id":"1234"}
```

List rooms:

```text
GET /api/room/list
```

Inspect a room:

```text
GET /api/room/info?id=1234
```

## Join A Player Seat

Connect to the normal player WebSocket endpoint:

```text
ws://127.0.0.1:9090/ws?room=1234&player_id=agent_p0&player_name=Codex&deck_code=<url-encoded-deck>
```

Rules:

- Use one stable `player_id` per controlled seat.
- Reuse the same `player_id` to reconnect to a started game.
- Use two WebSocket connections for agent-vs-agent play.
- A spectator connection uses `role=spectator`, but spectators cannot act and
  must not be used as hidden-information sources for a player agent.

## Server Messages

The first message is normally `joined`:

```json
{
  "type": "joined",
  "data": {
    "room_id": "1234",
    "slot": 0,
    "role": "",
    "reconnected": false
  }
}
```

Game state arrives inside `game_event`:

```json
{
  "type": "game_event",
  "event": {
    "type": "state_sync",
    "player": 0,
    "data": {}
  }
}
```

When `event.type` is `state_sync`, `event.data` is the current state visible to
that seat. Keep the latest state for each WebSocket and decide only from that
seat's visible data.

Errors arrive as:

```json
{"type":"error","message":"not your turn"}
```

Treat failed actions as no-ops. Wait for the next state or choose again from
the latest known state.

## Submit Actions

Send one JSON object through the WebSocket:

```json
{"action":"mulligan","data":{"keep":true}}
```

Common safe actions:

```json
{"action":"mulligan","data":{"keep":true}}
{"action":"no_defend","data":{}}
{"action":"resolve_action","data":{"selected":[]}}
{"action":"end_turn","data":{}}
```

For required pending choices, select candidate `instance_id` values from the
current `pending_action`:

```json
{"action":"resolve_action","data":{"selected":["ci_12"]}}
```

If `pending_action.cost` is present, include a `payment` map. If
`pending_action.can_overexert` is true, include `overexert_ids` only when using
overexertion to pay.

## State Reading

Important top-level fields in player state include:

- `phase`: `mulligan`, `main`, `defense_window`, `waiting_action`, or
  `game_over`
- `current_turn`: active main-phase player
- `turn_number`: current turn number
- `winner`: winner slot when the game is over
- `you`: private state for this player seat
- `opponent`: public opponent state
- `pending_spell`: spell currently being defended or reacted to
- `pending_action`: private pending choice for this player seat, when present

Prefer `instance_id` values from state over card names.

## Basic Play Loop

A conservative agent loop can be:

1. Wait for `state_sync`.
2. If `phase == "game_over"`, stop.
3. If `pending_action` is present for this seat, resolve it.
4. If `phase == "mulligan"`, send `mulligan`.
5. If `phase == "defense_window"` and this seat is defending, send
   `no_defend` unless intentionally testing defense.
6. If `phase == "main"` and `current_turn` is this seat, choose one legal main
   action or send `end_turn`.
7. Otherwise wait for another event.

Main-phase action names include:

- `summon`
- `consume`
- `cast_spell`
- `react_spell`
- `defend`
- `no_defend`
- `attack`
- `equip`
- `learn_skill`
- `use_item`
- `place_terrain`
- `use_ability`
- `resolve_action`
- `end_turn`
- `mulligan`

Let the server reject invalid actions. Do not mutate game state directly.

## Human Vs Agent

1. Human creates or joins a normal room in the browser.
2. Agent joins the other seat with `/ws` or `arcane-agent-cli`.
3. Agent reads only its own `state_sync`.
4. Human and agent act through the same server validation path.

## Agent Vs Agent

1. Create a normal room.
2. Open one WebSocket for `agent_p0`.
3. Open another WebSocket for `agent_p1`.
4. Maintain separate latest states and traces for each seat.
5. Do not share private hand or deck information between seats unless the test
   is explicitly a perfect-information experiment.

## Logs And Bug Reports

Room logs are written under:

```text
server/logs/rooms/*.jsonl
```

Use logs after a game or when reproducing a bug. Do not read room logs during a
fair match to choose moves.

When reporting an agent-found bug, include:

- commit hash
- room id
- controlled player id
- latest action sent
- server error, if any
- relevant log path or trace file
