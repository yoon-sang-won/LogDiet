# LogDiet for generic terminal agents

Use this package when an agent has no native LogDiet support.

## Status

- rules fallback: supported through `logdiet-rules.md`
- native adapter: not applicable by default
- transparent rewrite: no
- trust required: no

## Recommended setup

```sh
logdiet init --agent generic
logdiet doctor
```

## Default fallback behavior

Generic agents should always use:

```sh
logdiet wrap -- <command>
```

## manual wrapper mode

Run noisy commands through `logdiet wrap -- <cmd>`.

## PATH shim mode

Run `logdiet install`, put `.logdiet/bin` first in PATH, then use `logdiet doctor`.

## hook adapter

`logdiet-hook-rewrite.sh` is a template for agents that support custom command-hook scripts. It is not enabled automatically.

## How it works

`logdiet-rules.md` teaches the agent to prefer compact evidence and avoid full log pastes.

## Verification

Fixture tests verify the central rewrite policy. Generic hook runtime behavior depends on the host agent.

## Limitations

- No native automatic rewrite is assumed.
- Agent-specific tools may bypass LogDiet.
- Raw logs stay local under `.logdiet/runs/`.
