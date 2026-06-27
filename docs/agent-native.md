# Agent-Native LogDiet

## Goal

LogDiet v0.2 moves from CLI-first usage to agent-first usage.

The CLI remains the local engine, but users should install LogDiet through their coding agent whenever possible.

## Architecture

```text
coding agent
    |
    | command hook / plugin / skill / rules
    v
logdiet local engine
    |
    +-- captures stdout/stderr
    +-- preserves exit code
    +-- stores raw logs in .logdiet/runs/
    +-- prints compact evidence
```

## Integration levels

| Level | Name | What it does |
| ----- | ---- | ------------ |
| 1 | Rules only | Teaches the agent to prefer `logdiet wrap`, `show`, `grep`, and `raw` |
| 2 | PATH shim | Captures supported commands when `.logdiet/bin` is first in PATH |
| 3 | Hook rewrite | Rewrites supported shell commands to `logdiet wrap -- <cmd>` before execution |
| 4 | Native plugin package | Installs rules, hooks, and setup docs as one agent package |

## Agent support matrix

| Agent | v0.2 target | Integration type | Automatic command rewrite |
| ----- | ----------- | ---------------- | ------------------------- |
| Codex | yes | plugin/rules/hook package where supported | where hooks are trusted |
| Claude Code | yes | plugin/skill/hook package where supported | where hooks are supported |
| Cursor | yes | rules + hook config where supported | where hooks are supported |
| Gemini | yes | instructions + hook config where supported | where hooks are supported |
| Antigravity | yes | rules fallback first | not guaranteed |
| Generic terminal agents | yes | PATH shim + manual wrapper | via PATH shim only |

## CLI engine

The CLI remains the source of truth:

- `logdiet wrap -- <cmd>`
- `logdiet raw latest`
- `logdiet show latest:F1 --around 40`
- `logdiet grep latest "pattern"`
- `logdiet doctor`

## Truthful wording

Do not claim universal automatic behavior. Use:

> Automatic command rewriting is available where the agent supports command hooks. Other agents use rules/instructions fallback or manual `logdiet wrap`.

## Security

Hooks can change command execution. Generated hook files must be readable, local, and explicit. Users should review and trust hooks according to their agent's normal security flow.
