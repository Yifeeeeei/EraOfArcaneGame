# Agent Player Local Data

Codex-vs-Codex match artifacts are intentionally excluded from Git under:

```text
agent-data/
```

Initialize this directory after cloning:

```bash
cd server
go run ./cmd/agent-player init-data
```

The command creates missing templates and directories without overwriting
existing local knowledge.

The live local directory contains:

- `matches/<match-id>/`: transcripts, room log, independent reviews, and match
  summary.
- `knowledge/core-rules.md`: compact, repeatedly confirmed rules.
- `knowledge/gameplay-principles.md`: compact reusable play heuristics.
- `knowledge/deck-lab.md`: deck hypotheses, exact codes, matchup findings, and
  next controlled experiments.
- `knowledge/open-questions.md`: unresolved rules and active Issue links.
- `knowledge/retired-lessons.md`: disproven or obsolete conclusions, excluded
  from normal context.
- `context-packs/next-match.md`: bounded, coordinator-generated context for the
  next match; this is the default history entry point for player agents.
- `match-history.xlsx`: the structured match ledger and deck performance
  summary, used as a retrieval index rather than read wholesale.

See `docs/agent-player-protocol.md` for the required pre-match, in-match, and
post-match workflow.

This file documents the layout only. Do not commit the contents of
`agent-data/`.

Raw transcripts, room logs, and full reviews are cold evidence. Agents should
open them only when a relevant workbook row or context pack points to them.
Every 10 matches, compact the knowledge layer while preserving the raw archive.
