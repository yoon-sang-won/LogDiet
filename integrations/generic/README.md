# LogDiet for generic terminal agents

Use this package when an agent has no native LogDiet support.

## manual wrapper mode

Run noisy commands through `logdiet wrap -- <cmd>`.

## PATH shim mode

Run `logdiet install`, put `.logdiet/bin` first in PATH, then use `logdiet doctor`.

## rules-only fallback

Use `logdiet-rules.md` as the agent instruction file.

## hook adapter

`logdiet-hook-rewrite.sh` is a template for agents that support custom command-hook scripts.
