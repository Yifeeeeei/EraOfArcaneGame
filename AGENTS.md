# AGENTS.md

This file is for future coding agents working on EraOfArcaneGame. Treat it as the current project map and read it before changing code.

## What This Project Is

EraOfArcaneGame is a browser-playable prototype for the tabletop/card game "奥术纪元 / Era of Arcane".

The current product goal is narrow: make the base set playable and testable in a real two-player frontend match. Do not assume any expansion cards are supported yet.

## Current Scope

- Only the base set (`version_name == "基础包"`) is present in this repository and supported for live games.
- `server/cards/definitions_gen.go` is the compiled Go definition file for the 379 playable base cards.
- `data/supported_card_infos.json` is the base-card snapshot used for balance review and regeneration of compiled definitions.
- Non-base cards are intentionally absent from the runtime card pool.
- Runtime card behavior must be explicit Go code. Do not add text parsers that infer effects from card descriptions.

## Tech Stack

Backend:
- Go.
- Standard `net/http` server.
- `github.com/gorilla/websocket` for realtime matches.
- In-memory room and game state. There is no database.

Frontend:
- Static HTML/CSS/JS.
- Vue 3 loaded from CDN in each page.
- No npm frontend build step.
- Card images are loaded from `https://yifeeeeei.github.io/ArcaneImages/`.

Serving:
- The Go server serves API routes, WebSocket, and static frontend files from one process.
- Main routes:
  - `/`
  - `/game.html`
  - `/card-test.html`
  - `/api/cards`
  - `/api/deck/validate`
  - `/api/room/create`
  - `/api/room/list`
  - `/api/room/info`
  - `/ws`

## Important Files

- `server/main.go`: process entrypoint. Loads compiled base cards, sets playable card DB, registers behavior objects, serves routes.
- `server/cards/definitions_gen.go`: generated Go definitions for all currently supported base cards.
- `server/cards/interfaces.go` and `server/cards/category_markers_gen.go`: card category interfaces and generated marker methods for hero/companion/skill/item subtypes.
- `server/cards/loader.go`: loads compiled base cards and builds `BaseCardDB` / `PlayableCardDB`.
- `server/cards/snapshot.go`: exports the playable card pool as a stable JSON snapshot.
- `server/cmd/extract-supported-cards/main.go`: extracts `version_name == "基础包"` cards from `data/all_card_infos.json` into `data/supported_card_infos.json`.
- `server/cmd/check-card-metadata/main.go`: audits structured card metadata such as `effect_categories` and `effect_optionality`.
- `server/cmd/generate-card-definitions/main.go`: regenerates compiled Go card definitions from `data/supported_card_infos.json`.
- `server/cmd/snapshot-supported-cards/main.go`: regenerates `data/supported_card_infos.json`.
- `server/cmd/agent-player/main.go`: headless CLI used by Codex agents to initialize local match data, validate decks, create rooms, and play through the backend WebSocket without opening the frontend.
- `server/game/card_behavior.go`: card behavior interfaces such as `OnEnterBehavior`, `OnDeathBehavior`, `PerTurnAbility`, and `UltimateAbility`.
- `server/game/card_<number>_<name>.go`: one file per concrete base-set card with custom behavior.
- `server/game/card_effects_catalog.go`: registers lazy behavior factories with the engine adapter; it should not instantiate every behavior at startup.
- `server/game/engine.go`: main game engine and action handling.
- `server/game/rules.go`: focused rules helpers.
- `server/game/payment.go`: element payment logic.
- `server/game/base_cards_smoke_test.go`: smoke coverage for all currently supported base cards.
- `server/logs/rooms/*.jsonl`: ignored local room replay/debug logs written at runtime; inspect these first when reproducing a live-game bug.
- `web/index.html`: lobby.
- `web/game.html`: actual match UI.
- `web/card-test.html`: base-card testing workbench.
- `web/css/main.css`, `web/css/lobby.css`, `web/css/game.css`: current visual language.
- `docs/agent-player-protocol.md`: required machine-oriented workflow and action payload reference for headless Codex-agent matches.
- `docs/agent-player-data-layout.md`: explains the ignored local archive, bounded context packs, and long-term knowledge layout.

## Running Locally

Run from the `server` directory because paths are currently relative:

```bash
cd server
go run .
```

Then open:

```text
http://localhost:9090/
```

The server currently assumes `../web` is available relative to the `server` directory. Card data is compiled into Go and is not read from JSON at startup.

## Headless Codex-Agent Matches

When asked to run, coordinate, debug, or review a Codex-agent match without the
frontend, read these files before acting:

1. `docs/agent-player-protocol.md`
2. `docs/agent-player-data-layout.md`

Discover the CLI and initialize its ignored local data from the `server`
directory:

```bash
cd server
go run ./cmd/agent-player --help
go run ./cmd/agent-player init-data
```

Start exactly one game server from `server/`, then create a room in another
terminal:

```bash
GOCACHE=/tmp/eraofarcane-go-cache go run .
go run ./cmd/agent-player create-room
```

Use `agent-player validate-deck` before connecting. Each player agent must run a
separate persistent `agent-player connect` process with a distinct stable
`player-id` and transcript path. The second connection starts the match; actions
are newline-delimited JSON written to each client's stdin. `sample-deck` is
available for a first smoke match.

Do not commit match transcripts, reviews, evolving strategy, deck experiments,
or the match workbook. They belong under the ignored `agent-data/` directory,
not under `docs/`. Repository documentation should contain only the reusable
instructions required for another agent to operate the workflow.

Shared cross-machine match knowledge lives in the separate repository:

```text
https://github.com/Yifeeeeei/EraOfArcaneAgentLab
```

For agent-match tasks, prefer a sibling checkout at
`../EraOfArcaneAgentLab`. If it exists, read its `AGENTS.md`, then its bounded
`context-packs/bootstrap.md` and `context-packs/next-match.md`; do not preload
its complete match history. Record the exact EraOfArcaneGame commit in every
shared match or deck artifact.

The knowledge repository is optional infrastructure, not a game dependency.
Normal coding, tests, builds, and human frontend matches must work without it.
If it is unavailable, continue with ignored local `agent-data/` and report that
shared knowledge was not loaded. Do not silently clone, pull, push, or modify
the external repository unless the current task authorizes those network or
repository changes.

Routine raw transcripts and server room logs stay ignored locally. When a task
does authorize sharing new experience, promote only compact match metadata,
summaries, independent reviews, exact deck codes, and repeatedly supported
knowledge to EraOfArcaneAgentLab; game bugs themselves still belong in this
repository's GitHub Issues.

## Card Behavior Architecture

The game should be understandable from Go code alone. JSON snapshots are reference material only; they are not runtime truth.

- A horizontal card may still use `回合技` / `绝技` unless that specific card text or implementation says otherwise.
- A `消耗:` effect is different from a normal `回合技`: it pays by turning a vertical card horizontal, so it cannot be used when the source is already horizontal. Most `消耗:` cards should be implemented through `TriggerOnConsume` so `handleConsume` enforces this. Cards such as `渡鸦信使` that use an active button for `消耗:` must explicitly reject horizontal sources and set themselves horizontal without granting their printed load.
- Defensive payment uses `透支`, not `消耗`. During a defense window the player selects defense skills, optional boost skills, and units to overexert, then confirms together. Overexerting a unit turns it horizontal only for paying that defense cost; it does not trigger consume effects, does not grant elements to the pool, and any excess load is lost.
- End-of-turn cleanup follows the actual phase order: discard down to hand limit, reset cards, then settle marks. This matters for `冷却`: a horizontal skill with `冷却1` must fail to reset while the mark is still present, then the mark is removed, leaving the skill unavailable for the next turn as intended.
- Reactive sorceries such as `冰封消解` are not normal main-phase attack spells. Model them with explicit reaction behavior and a frontend action during the defense/spell window.
- If a card implements `SpellReactionBehavior`, make sure it is exposed to the frontend as `can_react`. Forgetting this makes the backend effect technically present but impossible for a human to use. Current base-set reaction examples include `冰封消解`, `风洞`, and `虹吸`.
- A player may inspect their own remaining deck only as an unordered summary. Opponent deck contents are hidden; opponent graveyards are public. If a hand card is revealed by a card effect or keyword, expose it through an explicit revealed-hand zone rather than by making the whole hand visible.
- Skill card UI must distinguish learn/entry cost (`elements_cost`) from cast/use cost (`elements_expense`). Do not collapse them into a single generic "费用" label when both exist.
- Card metadata lives in compiled Go definitions under `server/cards`.
- Card text metadata may include `effect_categories` (`主动`, `条件`, `光环`, `入场`, `遗言`, `反制`, `响应`, etc.) and `effect_optionality` (`强制`, `可选`). These fields are for review and frontend display only; runtime effects still require explicit Go behavior.
- Card categories are Go interfaces (`HeroCard`, `CompanionCard`, `SkillCard`, `ItemCard`, plus item subtypes) rather than runtime-only string checks.
- Custom rules live on concrete structs under `server/game`, one file per card, for example `card_1021006_grocer.go` containing `Card1021006Grocer`.
- Category and trigger behavior is expressed through Go interfaces. A card gets an enter effect by implementing `OnEnter(*EffectContext) error`; an ultimate by implementing `OnUltimate(*EffectContext) error`; and so on.
- Behavior interfaces are instance-aware. Do not decide that a card "has deathrattle", "has a per-turn ability", "can react", or "has a modifier" from interface conformance alone. The concrete behavior must implement the matching `HasActive...(*CardInstance) bool`, usually by embedding `AlwaysActive`. Rule checks should use active helpers such as `cardHasActiveDeathrattle`, `cardHasActivePerTurn`, and `cardHasActiveSpellReaction` when current card state matters. Example: `白骨骑士` still has the same card number after it returns, but its `HasActiveDeathrattle` becomes false.
- Runtime-granted abilities should use attached behavior objects on `CardInstance`, not ad hoc `Statuses` string checks. For example, "give a unit deathrattle" should attach an `AttachedDeathrattleBehavior`; then `cardHasActiveDeathrattle` and death resolution will see it like any printed deathrattle.
- `EffectRegistry` still exists as an engine adapter, but new work should add or change card behavior structs, not description parsers or string-inferred effects.
- Behavior registration is lazy. `RegisterAllCardEffects` registers factories only; concrete behavior objects should be constructed only when a card number is queried during play or serialization.
- Expansion cards should not be added until the base-set scope changes.

## Rule Clarifications From Prior Corrections

These are rules that previous agents have misunderstood. Treat them as hard constraints unless the user explicitly changes the game rules.

- Before fixing a bug or implementing a card-data update, compare the current card/rule text with the implementation and any test notes. If a test note conflicts with current effect text, prefer the effect text unless the user explicitly says that text is obsolete.
- Do not infer mechanics from `Description` text. Code such as `strings.Contains(card.Description, "精通")`, `includes("穿透")`, or similar text parsing is tech debt and should be removed, not expanded. Keywords and categories must be represented through Go interfaces or explicit card behavior methods.
- A keyword belongs to a specific card instance, not to a player globally. For example, `精通` is not a player-wide value; cards that care about mastered cards must inspect cards that implement the relevant interface/state.
- Spell stat modifiers have different owners and timing. Friendly spell bonuses are collected from the caster's field, but cards that say `敌方法术...` must implement the explicit enemy-spell modifier surface and be evaluated from the non-caster's field. Also keep power calculation and hit-damage calculation separate: damage prevention cards such as `冰霜之心` / `暗影披风` must not consume themselves while the engine is only calculating spell `威`.
- Some triggers originate from zones other than the battlefield. Example: `灵兽 辛柯` responds from hand/deck after a friendly unit takes enemy damage and should be summoned for free, not searched to hand, and not implemented as a battlefield-only passive.
- If a card cares about a specific damage source such as `点燃伤害`, pass explicit source metadata through the damage pipeline. Do not infer it from the target merely having a mark; a burning unit can take normal damage too.
- `绑定技能` is not a learnable skill. It does not enter the skill pool, does not occupy one of the 5 skill slots, and cannot be learned with `learn_skill`. It is attached to its host card as runtime state such as `BoundSkills`; when the host leaves the battlefield, the bound skill disappears with it and does not go to the graveyard as an independent card. Example: `"风暴之女" 艾拉雅` binds `风暴之怒` on herself.
- A derived/generated skill may still have a concrete Go behavior file. Being generated or bound does not justify special parser logic.
- `反制` consumable items are set into the item/equipment area for free and face-down. They pay their printed entry cost only when their trigger window opens and the owner chooses to reveal/resolve them; that reveal payment may use overexertion during the opponent's turn.
- Cards reset from horizontal to vertical at the end of their owner's turn, not at the beginning of that player's next turn.
- `消耗:` is a tap payment. It requires the source card to be vertical first. A horizontal card can still use `回合技` or `绝技` if rules allow, but it cannot pay a `消耗:` cost because that cost is turning vertical to horizontal.
- `奥术`/`无` in a play cost is a wildcard requirement whose exact paid element can be ambiguous. If the player has multiple possible elements that could satisfy an arcane cost, the frontend should ask the player to choose the payment rather than silently picking one.
- `吞噬:3\气` and similar summon requirements happen before choosing/placing the summoned card on the battlefield. The player must select a valid friendly companion with sufficient load and destroy it as part of summoning. Do not model this as a later optional ability.
- `速攻` means the card can be vertical and usable immediately after entering, according to the card's type/implementation. Apply the keyword generally, not only for one named card.
- Bound, generated, or temporary cards still need to be visible enough in the frontend for a human to understand why an effect exists. Showing bound skills in the inspector/detail view is appropriate; showing them in the skill pool is wrong.
- Remaining deck visibility must not reveal order. If exposing deck contents, expose only an unordered summary by card number/name/count, never the ordered deck slice or instance IDs.

## Tests

Always run:

```bash
cd server
go test ./...
```

For frontend/gameplay changes, backend tests are not enough. The project has already found issues that only appear when operating the UI. Do a browser-level check when changing `web/game.html`, `web/css/game.css`, WebSocket state shape, targeting, pending actions, or turn flow.

Minimum frontend sanity check:
- Start server.
- Create a room from `card-test.html`.
- Join with two real frontend clients.
- Keep mulligan on both.
- Learn a skill.
- Consume for elements.
- Cast a spell.
- Resolve defense window.
- Resolve discard pending action.
- Reach game-over on both clients.

## Supported Card Snapshot Workflow

After changing base card data or balance values from the spreadsheet export, regenerate the supported snapshot and compiled definitions:

```bash
cd server
go run ./cmd/extract-supported-cards
go run ./cmd/generate-card-definitions
go test ./...
```

For metadata review, run this separately:

```bash
cd server
go run ./cmd/check-card-metadata
```

`check-card-metadata` may fail while `effect_categories` / `effect_optionality` are still being filled in; treat that as a data-review checklist, not as runtime behavior or a required test gate until the metadata is complete.

Then inspect:

```bash
git diff -- data/supported_card_infos.json server/cards/definitions_gen.go server/cards/category_markers_gen.go
```

This diff is intended to show exactly which base cards changed, were added, or were removed.

## Design Language

The frontend should visually match:

- `https://yifeeeeei.github.io/EraOfArcane/`
- `https://yifeeeeei.github.io/ArcaneComposer/`

Current direction:
- dark black/brown tabletop surface
- thin gold borders
- square, not pill-shaped, controls
- serif Chinese fantasy typography
- restrained tool-like layout for lobby/test bench
- dense but readable game board

Avoid introducing generic modern SaaS styling, bright blue panels, purple gradients, large rounded cards, or marketing landing-page sections.

## Deployment Notes

The current app is easiest to deploy as one Go service behind a reverse proxy:

```text
Nginx/Caddy HTTPS
  -> http://127.0.0.1:9090
      Go server
        /api/*
        /ws
        static web/*
```

Do not hardcode frontend API calls to `localhost`; in the browser, localhost means the player's machine. Use relative paths and `location.host` for WebSocket, as `web/game.html` currently does.

## Known Constraints And Risks

- Room/game state is in memory. Restarting the server drops rooms.
- No authentication or rate limiting exists yet.
- Not ready for multi-instance deployment.
- Many older root docs were written by a previous agent and may be stale. Prefer code, tests, this file, and recent git history over old status docs.
- Be careful with existing dirty worktrees. Do not revert user changes unless explicitly asked.

## Git Hygiene

- This project is changing quickly. Before investigating a bug or starting new implementation work, sync with the latest `main` when possible; the issue may already be fixed. If the worktree is dirty, inspect the changes first and do not overwrite user work.
- Keep commits small and descriptive.
- Do not commit generated binaries.
- Do not commit stale progress notes unless they are updated to current facts.
- If adding a new playable card behavior, add or update tests.

## Issue Workflow

- Long-term bug and feature tracking should go through GitHub Issues, not only chat history.
- Only use the supported issue forms: Bug Report, Card Effect, and Frontend UX. Blank issues are disabled intentionally.
- Every issue must include a version or commit hash. If there is no formal version number yet, use the exact commit tested.
- Do not require every PR to close an issue; normal iteration PRs are allowed. When a PR does fix a tracked issue, mention it in the PR description.
