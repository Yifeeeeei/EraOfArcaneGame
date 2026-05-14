# AGENTS.md

This file is for future coding agents working on EraOfArcaneGame. Treat it as the current project map and read it before changing code.

## What This Project Is

EraOfArcaneGame is a browser-playable prototype for the tabletop/card game "奥术纪元 / Era of Arcane".

The current product goal is narrow: make the base set playable and testable in a real two-player frontend match. Do not assume any expansion cards are supported yet.

## Current Scope

- Only the base set (`version_name == "基础包"`) is present in this repository and supported for live games.
- `server/cards/definitions_gen.go` is the compiled Go definition file for the 378 playable base cards.
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
- `server/cmd/snapshot-supported-cards/main.go`: regenerates `data/supported_card_infos.json`.
- `server/game/card_behavior.go`: card behavior interfaces such as `OnEnterBehavior`, `OnDeathBehavior`, `PerTurnAbility`, and `UltimateAbility`.
- `server/game/card_effects_base_*.go`: concrete base-set card structs that own their custom behavior.
- `server/game/card_effects_catalog.go`: registers the current base-set behavior objects with the engine adapter.
- `server/game/engine.go`: main game engine and action handling.
- `server/game/rules.go`: focused rules helpers.
- `server/game/payment.go`: element payment logic.
- `server/game/base_cards_smoke_test.go`: smoke coverage for all currently supported base cards.
- `web/index.html`: lobby.
- `web/game.html`: actual match UI.
- `web/card-test.html`: base-card testing workbench.
- `web/css/main.css`, `web/css/lobby.css`, `web/css/game.css`: current visual language.

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

## Card Behavior Architecture

The game should be understandable from Go code alone. JSON snapshots are reference material only; they are not runtime truth.

- Card metadata lives in compiled Go definitions under `server/cards`.
- Card categories are Go interfaces (`HeroCard`, `CompanionCard`, `SkillCard`, `ItemCard`, plus item subtypes) rather than runtime-only string checks.
- Custom rules live on concrete structs under `server/game`, for example `Card1021006Grocer`.
- Category and trigger behavior is expressed through Go interfaces. A card gets an enter effect by implementing `OnEnter(*EffectContext) error`; an ultimate by implementing `OnUltimate(*EffectContext) error`; and so on.
- `EffectRegistry` still exists as an engine adapter, but new work should add or change card behavior structs, not description parsers or string-inferred effects.
- Expansion cards should not be added until the base-set scope changes.

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

After changing base card data or balance values, regenerate the compiled definitions and snapshot:

```bash
cd server
go run ./cmd/snapshot-supported-cards
go test ./...
```

Then inspect:

```bash
git diff -- data/supported_card_infos.json server/cards/definitions_gen.go
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

- Keep commits small and descriptive.
- Do not commit generated binaries.
- Do not commit stale progress notes unless they are updated to current facts.
- If adding a new playable card behavior, add or update tests.
