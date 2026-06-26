# Agent Knowledge Map

This repository provides public game knowledge and pointers. It does not store
personal agent experience or strategy notes.

## Current Playable Card Pool

The runtime-supported card pool is the base set only.

Use one of these sources:

- Live server API: `GET /api/cards`
- Repository snapshot: `data/supported_card_infos.json`
- Compiled runtime definitions: `server/cards/definitions_gen.go`

When sources disagree, prefer the running server and compiled runtime behavior.
Card descriptions are reference text; actual effects are implemented in Go
under `server/game`.

## Rules

Use the rules documents under:

```text
docs/rules/
```

Important starting points:

- `docs/rules/core-rules.md`
- `docs/rules/glossary.md`
- `docs/rules/card-text-style-guide.md`

If a rule document and runtime behavior disagree during play, treat the server
as authoritative and report the mismatch as a bug or documentation issue.

## Deck Building

Decks used with the CLI must validate through:

```text
POST /api/deck/validate
```

Agents may use external deck-building pages or copied deck lists, but every deck
must be checked against this repository's current supported card pool before
joining a match.

## Normal Play Interface

Normal agent play should use:

- `POST /api/room/create`
- `POST /api/deck/validate`
- `/ws`
- `server/cmd/arcane-agent-cli`

Test-room APIs are for training setups and bug reproduction only, not fair
matches.
