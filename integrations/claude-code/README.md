# LogDiet for Claude Code

This directory contains LogDiet rule, skill, and hook templates for Claude Code.

It is not an official Claude plugin and does not claim automatic install behavior.

## Status

- rules fallback: supported through `CLAUDE.md` or `skill.md`
- native adapter: hook rewrite template
- transparent rewrite: template / not verified
- trust required: yes, if your Claude Code setup supports command hooks

## Recommended setup

```sh
logdiet init --agent claude --mode native
logdiet doctor
```

Review all generated files before enabling any hook.

## Default fallback behavior

If a Claude Code hook is unavailable or untrusted, use:

```sh
logdiet wrap -- <command>
```

## How it works

`hook-rewrite-template.sh` is a thin delegate template. It reads a command from `COMMAND`, calls `logdiet hook rewrite --command "$COMMAND"`, and returns LogDiet's decision.

## Verification

LogDiet verifies the local rewrite engine and template files. Runtime Claude Code hook behavior is not verified by this repository.

## Limitations

- The hook protocol must be adapted to the exact Claude Code environment.
- Built-in non-shell tools may bypass LogDiet.
- Raw logs remain local under `.logdiet/runs/`.
