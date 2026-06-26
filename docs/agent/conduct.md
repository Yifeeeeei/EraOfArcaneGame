# Agent Conduct And Local Experience

This repository provides shared tooling, public rules, and public knowledge
entry points. It must not become a store for one agent's private experience,
strategy, prompts, or metagame conclusions.

## Fair Play

During a normal match, an agent must:

- act only through the normal player interface or `arcane-agent-cli`
- use only the state visible to its own player seat
- let the server validate every action
- avoid spectator state, room logs, ordered deck data, and test-room debug APIs
- preserve enough action/error context to report bugs honestly

During a normal match, an agent must not:

- inspect the opponent's hidden hand or ordered deck
- mutate game state through debug APIs
- use server logs to choose moves
- treat a test-room scenario as a fair match result

## Training And Reproduction

Debug APIs and test rooms may be used for:

- reproducing a suspected bug
- building a fixed tactical scenario
- checking whether an action shape is accepted
- studying one card interaction in isolation

Results from debug or test-room sessions should be marked as non-match data.

## Local Experience

Agents may keep local experience outside the repository. Recommended locations:

```text
.agent-memory/
local/agent-experience/
```

These paths are ignored by git. Users may copy or replace them to transfer an
agent's experience between machines or agents.

Recommended local structure:

```text
.agent-memory/
  profile.md
  match-notes/
  deck-notes/
  policy-notes/
  traces/
```

Repository commits should not include:

- personal strategy notes
- deck preference notes
- learned matchup evaluations
- private prompts
- win-rate tables
- raw match traces containing private state

Reusable, non-personal improvements belong in code or public docs. Personal
experience belongs outside version control.
