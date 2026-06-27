# LogDiet for Codex

This directory contains original LogDiet templates for using LogDiet with Codex.

## Layer A: rules fallback

Use `AGENTS.md` or `logdiet-instructions.md` to teach Codex to prefer compact LogDiet evidence and avoid pasted log walls.

## Layer B: hook rewrite template

`hook-rewrite-template.sh` shows how a Codex command hook could call:

```sh
logdiet hook rewrite --command "$COMMAND"
```

This is a hook rewrite template. It is not automatically installed or enabled by this repository.

Automatic command rewriting is available where the agent supports command hooks. Other agents use rules/instructions fallback or manual `logdiet wrap`.
