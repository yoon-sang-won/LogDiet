# Agent Self-Install

## Goal

A user should be able to say:

```text
Install https://github.com/yoon-sang-won/LogDiet and use it for noisy command output.
```

The agent should install LogDiet and start using it without requiring native hooks.

## Default path

1. Install the local engine.
2. Run `logdiet bootstrap --agent auto`.
3. Run `logdiet doctor`.
4. Use `logdiet wrap -- <command>` for noisy commands.
5. Use `show`, `grep`, and `raw` to expand evidence.

## Why hooks are optional

Command hooks can modify command execution. Some agents require review/trust before hooks run.

LogDiet therefore treats hooks as advanced mode.

The default path relies on rules/instructions that teach the agent to call `logdiet wrap` explicitly.

## What agents should do

Agents should use LogDiet for:

- test commands;
- build commands;
- type checks;
- git diff/status/log;
- search commands;
- any command likely to produce noisy output.

Agents should not ask users to paste full logs unless `show`, `grep`, and `raw` are insufficient.

## Manual fallback

If bootstrap fails, use:

```sh
logdiet wrap -- <command>
```
