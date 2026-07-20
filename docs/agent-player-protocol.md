# Codex Agent Player Protocol

This document is the compact, machine-oriented guide for Codex agents playing
EraOfArcaneGame without opening `web/game.html`.

The Go server remains authoritative. Frontend checks are conveniences, not
rules: when this document and a server error disagree, obey the server and
record the discrepancy under the ignored local `agent-data/` directory.

## Persistent local data

All match-specific artifacts and evolving player knowledge live under
`agent-data/`, which is intentionally ignored by Git:

```text
agent-data/
  matches/<match-id>/
    player-a.jsonl
    player-b.jsonl
    player-a-review.md
    player-b-review.md
    room.jsonl
    match-summary.md
  knowledge/
    core-rules.md
    gameplay-principles.md
    deck-lab.md
    open-questions.md
    retired-lessons.md
  context-packs/
    next-match.md
  match-history.xlsx
```

Do not place reviews, transcripts, accumulated strategy notes, deck experiments,
or the match workbook under `docs/`. Repository documentation describes how to
play; `agent-data/` stores what the agents learned from playing.

When a sibling `EraOfArcaneAgentLab` checkout is available for an authorized
shared match, read its `AGENTS.md` and bounded context packs before selecting
historical evidence. Do not read its entire match archive. Runtime transcripts
remain under ignored local storage; only compact artifacts are promoted to the
knowledge repository after the match. The external repository is optional and
must never block this CLI or the game server.

## Start a match

From a fresh clone, initialize the ignored local knowledge/archive structure.
This command preserves any files that already exist:

```bash
cd server
go run ./cmd/agent-player init-data
```

Run the server from `server/`:

```bash
GOCACHE=/tmp/eraofarcane-go-cache go run .
```

In another terminal, create a normal room:

```bash
go run ./cmd/agent-player create-room
```

The result contains a four-digit `room_id`. Each Codex agent then opens its own
persistent client. A sample valid deck is available with:

```bash
go run ./cmd/agent-player sample-deck > /tmp/arcane-agent-deck.txt
```

Player 1:

```bash
go run ./cmd/agent-player connect \
  -room ROOM_ID -player-id codex-a -name CodexA \
  -sample-deck \
  -transcript ../agent-data/matches/MATCH_ID/player-a.jsonl
```

Player 2 uses a different stable player ID and transcript path. Connecting the
second player starts the match. Keep both processes running and write one JSON
action per line to their stdin.

Use `-deck-file` instead of `-sample-deck` when testing a chosen deck.

The `player-id` is the reconnection identity. Reusing it reconnects to the same
player slot after a dropped connection.

## Client output

Every line is a JSON envelope:

```json
{
  "time": "2026-07-19T01:02:03Z",
  "direction": "received",
  "payload": {}
}
```

By default, stdout contains only `joined`, `error`, `state_sync`, and sent
actions so that incidental animation events do not consume the Codex context.
When `-transcript` is set, the file still records every event. Pass
`-all-events` when live inspection of every effect event is useful.

Important received payloads:

- `{"type":"joined",...}`: connection accepted; `data.slot` is `0` or `1`.
- `{"type":"game_event","event":{"type":"state_sync","data":STATE}}`:
  authoritative state visible to this player.
- `{"type":"error","message":"..."}`: the last submitted action was rejected.
- Other `game_event` values describe effects and are useful for the match
  narrative, but decisions should be based on the latest `state_sync`.

The server broadcasts a new `state_sync` after every accepted action. It also
broadcasts state after rejected actions only when another action later succeeds,
so retain the last state when handling an error.

## Observation rules

Use only the state received by the acting player. Do not read the opponent
agent's transcript or server internals while choosing a move.

Key state fields:

- `phase`: `mulligan`, `main`, `defense_window`, `waiting_action`, or
  `game_over` are the important decision phases.
- `current_turn`: slot whose normal turn it is.
- `you`: private player state, including hand and skill pool.
- `opponent`: public opponent state; hidden hand/deck information is omitted.
- `pending_spell`: current spell and its attack power.
- `pending_action`: present only for the player required to resolve it.
- `winner`: `-1` while playing, `0`/`1` for a winner, `-2` for a draw.

Card instance IDs, not card numbers, identify cards in actions.

The server exposes useful card flags such as `can_attack`, `can_defend`,
`can_react`, `has_per_turn`, `has_ultimate`, `needs_target`,
`devour_requirement`, `elements_cost`, and `elements_expense`. Treat these as
hints; the action handler performs the final validation.

Grid coordinates are zero based: `col` and `row` are each `0..2`. The hero
starts at `(1,1)`.

## Action envelope

All input lines have this shape:

```json
{"action":"ACTION_NAME","data":{}}
```

Unknown fields are generally ignored. Optional `payment` maps explicitly name
the elements spent:

```json
{"payment":{"火":1,"气":1}}
```

Explicit payment is required when wildcard (`无`) or light substitution makes
the choice ambiguous. Otherwise the server can choose a valid payment.

## Actions

### Mulligan

```json
{"action":"mulligan","data":{"keep":true}}
```

`keep:false` replaces the entire starting hand. Both players must submit.

### Consume

```json
{"action":"consume","data":{"instance_id":"INSTANCE_ID"}}
```

On the first player's first turn, consuming a hero with multiple load elements
may require an explicit reduced gain selection:

```json
{"action":"consume","data":{"instance_id":"INSTANCE_ID","gain":{"火":1}}}
```

### Summon

```json
{"action":"summon","data":{"instance_id":"HAND_INSTANCE","col":0,"row":0}}
```

For a companion with a devour requirement:

```json
{"action":"summon","data":{"instance_id":"HAND_INSTANCE","col":0,"row":0,"devour_ids":["FRIENDLY_INSTANCE"]}}
```

Add `payment` when required.

### Learn or replace a skill

```json
{"action":"learn_skill","data":{"instance_id":"POOL_INSTANCE","replace_id":""}}
```

Set `replace_id` to an existing skill instance when no empty skill slot exists.

### Equip or replace equipment

```json
{"action":"equip","data":{"instance_id":"HAND_INSTANCE","replace_id":""}}
```

Set `replace_id` to a valid existing equipment instance when replacing.

### Use an item

Untargeted:

```json
{"action":"use_item","data":{"instance_id":"HAND_INSTANCE"}}
```

Targeted spell scroll:

```json
{"action":"use_item","data":{"instance_id":"HAND_INSTANCE","target_type":"unit","target_owner":1,"target_col":1,"target_row":1}}
```

Target types may also be `hero` or `none`, according to the card flags. Counter
items are set face down through the same `use_item` action when their rules
permit it.

### Place terrain

```json
{"action":"place_terrain","data":{"instance_id":"HAND_INSTANCE","col":0,"row":0}}
```

### Cast a spell

Targeted:

```json
{"action":"cast_spell","data":{"instance_id":"SKILL_INSTANCE","target_type":"unit","target_owner":1,"target_col":1,"target_row":1,"boost_ids":[]}}
```

No target:

```json
{"action":"cast_spell","data":{"instance_id":"SKILL_INSTANCE","target_type":"none","boost_ids":[]}}
```

Friendly hero target:

```json
{"action":"cast_spell","data":{"instance_id":"SKILL_INSTANCE","target_type":"hero","boost_ids":[]}}
```

Chain Lightning (`3321001`) may additionally use `extra_target_col` and
`extra_target_row`. Add `payment` when the combined main and boost cost has an
ambiguous payment.

### React to a spell

During the opponent's spell window:

```json
{"action":"react_spell","data":{"instance_id":"SKILL_INSTANCE","overexert_ids":[]}}
```

### Defend or decline

```json
{"action":"defend","data":{"skill_ids":["SKILL_INSTANCE"],"scroll_ids":[],"boost_ids":[],"overexert_ids":[]}}
```

Or:

```json
{"action":"no_defend","data":{}}
```

`overexert_ids` are friendly units/equipment used only to pay defense or
reaction costs. Add `payment` if the available elements and overexerted load
allow multiple ambiguous payment combinations.

### Direct attack

```json
{"action":"attack","data":{"attacker_id":"INSTANCE_ID","target_col":1,"target_row":1}}
```

### Activated ability

```json
{"action":"use_ability","data":{"instance_id":"FIELD_INSTANCE","ability_type":"per_turn"}}
```

`ability_type` may be `per_turn` or `ultimate`. Some abilities accept a
`target_id`; many instead create a `pending_action`.

### Resolve a pending action

Select candidate IDs exactly as supplied by `pending_action.candidates`:

```json
{"action":"resolve_action","data":{"selected":["CANDIDATE_ID"]}}
```

Respect `min_select` and `max_select`. Depending on `pending_action.context`,
additional fields can include:

- `overexert_ids`
- `payment`
- `top_order`
- `bottom_order`

For a zero-minimum optional choice, `selected:[]` declines it.

### End turn

```json
{"action":"end_turn","data":{}}
```

Only the current player in `main` phase can end the turn. Resolve all pending
actions and defense windows first.

## Codex match procedure

### Context policy

Do not read every historical match, transcript, or review. Raw match folders are
cold evidence and should be opened only when investigating a specific Match ID,
card, deck, or suspected defect.

The coordinator prepares `agent-data/context-packs/next-match.md` before each
match. Player agents normally read only:

1. this protocol;
2. `next-match.md`;
3. the compact knowledge files explicitly referenced by that context pack;
4. their own chosen deck entry in `deck-lab.md`.

Default retrieval limits:

- at most 3–5 relevant historical matches;
- at most the latest 2 matches unless older matches are more relevant;
- `next-match.md` below 1,500 Chinese characters or about 900 English words;
- each general knowledge file below 2,000 Chinese characters or about 1,200
  English words;
- each deck profile below 1,000 Chinese characters or about 600 English words;
- each one-line match summary below 100 Chinese characters or about 60 English
  words.

Use the workbook metadata to filter by tested commit, deck ID, archetype, tags,
result, and experiment hypothesis. Read a full review only when its compact
summary is insufficient. Read a transcript or room log only for evidence.

### Before the match

1. The coordinator reads the workbook index, selects relevant history, and
   rewrites `context-packs/next-match.md`.
2. Each player reads the bounded context described above.
3. Inspect prior deck results and choose a deck hypothesis to test. Reusing a
   known deck is allowed for a controlled comparison, but agents are encouraged
   to explore deck building: change a coherent package, state why, and compare
   it with previous results.
4. Validate each deck with `agent-player validate-deck`.
5. Create a unique `agent-data/matches/<match-id>/` folder and connect with a
   private transcript inside it.

### During the match

1. Keep or redraw the opening hand based on the chosen deck plan.
2. After every action, wait for the next `state_sync` or `error`.
3. Never submit an action for the other player or inspect the opponent's private
   transcript.
4. Treat play quality as a first-class objective, not merely a way to exercise
   code. Track sequencing, resource efficiency, board geometry, revealed
   information, lethal setups, and mistakes.
5. Continue until `phase == "game_over"` or a genuine, evidenced deadlock is
   established.

### After the match

1. Move both private transcripts, the room log, and both independent reviews
   into the match folder. Cross-player comparison is allowed only now.
2. Write `match-summary.md`: result, decisive sequence, deck hypotheses,
   lessons, suspected defects, and follow-up experiment.
3. Promote only stable rules to `knowledge/core-rules.md`, reusable play
   heuristics to `knowledge/gameplay-principles.md`, unresolved items to
   `knowledge/open-questions.md`, and deck-specific findings to
   `knowledge/deck-lab.md`.
4. Append one row to `match-history.xlsx`, including both exact deck codes,
   result, duration, turns, artifact paths, and one-sentence key process.
5. For every reproducible suspected bug:
   - confirm the tested commit hash;
   - collect the smallest useful transcript/room-log evidence;
   - check whether an existing GitHub Issue already covers it;
   - if not, create an Issue using one of the supported forms: **Bug Report**,
     **Card Effect**, or **Frontend UX**;
   - include the commit hash, reproduction steps, expected behavior, actual
     behavior, and evidence path;
   - record the Issue URL in `match-summary.md` and the workbook.
6. Do not file an Issue for an uncertain rule interpretation. Record it as a
   question in the review until the rule or card text is confirmed.

Future agents must build on these artifacts. They should not start every match
from the sample deck or rediscover already documented lessons without a stated
reason.

### Periodic memory maintenance

After every 10 completed matches, a coordinator performs a memory compaction:

1. merge duplicate lessons and retain source Match IDs;
2. promote repeatedly confirmed conclusions;
3. move disproven, obsolete, or version-specific guidance to
   `retired-lessons.md`;
4. shorten deck profiles to their current plan, evidence, and next experiment;
5. close or update resolved questions and Issue links;
6. verify every knowledge file remains within its length budget;
7. regenerate `next-match.md` from the compacted knowledge and workbook index.

Compaction never deletes raw match artifacts. It changes what is loaded by
default, not what evidence remains available.

## Sources of truth

- Transport and authentication: `server/api/ws.go`
- Action dispatch: `server/game/engine.go` (`HandleAction`)
- Player-visible state: `server/game/engine.go` (`GetStateForPlayer`)
- Pending choices: `server/game/state.go` (`PendingAction`)
- Reference frontend action construction: `web/game.html`
