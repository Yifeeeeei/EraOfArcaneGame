# Agent Playbook

This document is for developer agents such as Codex. It explains how to play
Era of Arcane through the existing human client protocol. Do not add or expect a
user-facing "connect agent" button. An agent should act like a normal player:
create or join a room, open a websocket seat, read `state_sync`, and submit the
same actions a browser client would submit.

## Scope

- Use the existing HTTP and websocket APIs.
- Do not require new agent-specific endpoints.
- Do not expose agent controls in the ordinary player UI.
- Let the server validate all moves through the normal game engine.

## Required Server

Run the game server from the repository root:

```bash
cd server
go run .
```

The local server listens on:

```text
http://127.0.0.1:9090
```

## Deck Code

A player websocket must provide a valid deck code. This sample deck is already
used by the lobby and tests:

```text
4311003 // 1021001 1021001 1021004 1021004 1021011 1021011 1321002 1321002 1321001 1321001 1321003 1321003 1321007 1321007 1321013 1321013 1321011 1321011 1321010 1321010 1321008 1321008 1321004 1321004 1021006 1021006 1021007 1021007 1021012 1021012 // 3321002 3321003 3321005 3321010 3321013 3321015 3021008 3021009 3021005 3021003
```

Validate a deck before joining:

```bash
curl -s http://127.0.0.1:9090/api/deck/validate \
  -H 'Content-Type: application/json' \
  -d '{"deck_code":"4311003 // 1021001 1021001 1021004 1021004 1021011 1021011 1321002 1321002 1321001 1321001 1321003 1321003 1321007 1321007 1321013 1321013 1321011 1321011 1321010 1321010 1321008 1321008 1321004 1321004 1021006 1021006 1021007 1021007 1021012 1021012 // 3321002 3321003 3321005 3321010 3321013 3321015 3021008 3021009 3021005 3021003"}'
```

## Create A Room

Create a normal room:

```bash
curl -s -X POST http://127.0.0.1:9090/api/room/create
```

The response is:

```json
{"room_id":"1234"}
```

List rooms:

```bash
curl -s http://127.0.0.1:9090/api/room/list
```

Inspect one room:

```bash
curl -s 'http://127.0.0.1:9090/api/room/info?id=1234'
```

## Join As A Player

Open a websocket connection to occupy a player seat:

```text
ws://127.0.0.1:9090/ws?room=1234&player_id=agent_p0&player_name=Codex&deck_code=<url-encoded-deck-code>
```

Important rules:

- Use a stable unique `player_id` per controlled seat.
- Use two websocket connections for an agent-vs-agent game.
- The game starts automatically when both player seats are occupied.
- To reconnect to a started game, reuse the same `player_id`.
- A spectator connection uses `role=spectator` and does not include `deck_code`.

## Message Shapes

The server first sends a `joined` message:

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

Most gameplay messages are wrapped game events:

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

When `event.type == "state_sync"`, `event.data` is the current state visible to
that seat. Keep the latest state per websocket and decide actions from it.

Submit actions as plain websocket JSON:

```json
{"action":"mulligan","data":{"keep":true}}
```

Errors arrive as:

```json
{"type":"error","message":"not your turn"}
```

On an error, read the next `state_sync` or request a new plan from the current
state. Do not assume a failed action changed the game.

## State Reading

The player state has these top-level fields:

- `phase`: current phase, such as `mulligan`, `main`, `defense_window`,
  `waiting_action`, or `game_over`.
- `current_turn`: player slot whose main turn it is.
- `turn_number`: current turn number.
- `winner`: winner slot when the game is over.
- `you`: private state for this websocket seat, including hand, skill pool,
  deck summary, field, elements, and graveyard.
- `opponent`: public opponent state.
- `pending_spell`: spell being defended or reacted to.
- `pending_action`: private pending choice for this seat, when any.

For action planning, prefer IDs from the state over card names. Cards and
candidates expose `instance_id`; actions usually refer to those IDs.

## Basic Action Loop

An agent can make useful progress with this conservative loop:

1. Wait for `state_sync`.
2. If `phase == "game_over"`, stop.
3. If `pending_action` is present, resolve it.
4. If `phase == "mulligan"`, send `mulligan`.
5. If `phase == "defense_window"` and the pending spell is aimed at this
   player, send `no_defend` unless intentionally testing defense.
6. If `phase == "main"` and `current_turn` is this seat, choose a legal main
   action or send `end_turn`.
7. Otherwise wait for another event.

Minimal safe actions:

```json
{"action":"mulligan","data":{"keep":true}}
{"action":"no_defend","data":{}}
{"action":"end_turn","data":{}}
```

Pending action resolution:

```json
{"action":"resolve_action","data":{"selected":["candidate_instance_id"]}}
```

For optional pending actions with `min_select == 0`, an agent may choose:

```json
{"action":"resolve_action","data":{"selected":[]}}
```

If `pending_action.cost` is present, include a valid `payment` map in `data`.
If `pending_action.can_overexert` is true, include `overexert_ids` only when the
selected payment plan needs them.

## Main Phase Actions

The full action surface is the same as `game.Engine.HandleAction`:

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

Useful examples:

```json
{"action":"consume","data":{"instance_id":"card_instance_id"}}
```

```json
{"action":"learn_skill","data":{"instance_id":"skill_pool_instance_id"}}
```

If the skill area is full, replace a vertical learned skill:

```json
{"action":"learn_skill","data":{"instance_id":"skill_pool_instance_id","replace_id":"learned_skill_instance_id"}}
```

```json
{"action":"summon","data":{"instance_id":"hand_instance_id","col":1,"row":1}}
```

```json
{"action":"attack","data":{"attacker_id":"unit_instance_id","target_col":1,"target_row":0}}
```

Many actions require costs, targets, replacement choices, or card-specific
fields. If an action fails, use the error plus the latest state to choose a
simpler legal action. For broad smoke testing, it is acceptable to fall back to
`end_turn` when no confident legal main action is available.

## Human Vs Agent

To play with a human:

1. Create a room through `/api/room/create` or let the human create one in the
   lobby.
2. The human joins one seat through the browser.
3. The agent opens a websocket for the other seat with its deck code.
4. The agent follows the basic action loop from its private `state_sync`.

The human UI remains unchanged. The agent is just another websocket client.

## Agent Vs Agent

For a pure agent game:

1. Create a room.
2. Open websocket A with `player_id=agent_p0`.
3. Open websocket B with `player_id=agent_p1`.
4. Maintain separate latest states for A and B.
5. Drive each seat only from its own private `state_sync`.

For subagent testing, give each subagent only its own websocket URL, deck code,
seat identity, and this playbook. Do not share the other player hand or private
state unless the test intentionally requires perfect information.

## Logging

Room logs are written under:

```text
server/logs/rooms/*.jsonl
```

They include room events, client actions, action errors, game events, and state
snapshots. When reporting an agent-found bug, include:

- commit hash
- room id
- log path
- controlled player id
- latest action sent
- server error, if any

## Expectations For Agents

- Treat the server as authoritative.
- Never mutate game state outside normal actions.
- Never infer hidden opponent information from local files or logs during a
  normal playtest.
- Prefer deterministic, simple policies for reproduction.
- Preserve the exact websocket messages that led to a bug.
