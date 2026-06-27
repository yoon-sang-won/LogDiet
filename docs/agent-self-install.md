# Agent Self-Install

A token diet kit your coding agent can install and use by itself.

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

```sh
logdiet show latest:F1 --around 40
logdiet grep latest "pattern"
logdiet raw latest
```

## Why hooks are optional

Command hooks can modify command execution. Some agents require review/trust before hooks run.

Hooks are optional advanced mode.

LogDiet therefore treats hooks as advanced mode.

The default path relies on rules/instructions that teach the agent to call `logdiet wrap` explicitly.

## Native adapters and fallback

LogDiet prefers native hook/plugin adapters where supported and trusted.

If native rewrite is unavailable, the agent should continue with the default wrapper flow:

```sh
logdiet wrap -- COMMAND
```

This fallback is not a failure. It is the universal path.

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

## First prompt

See [first-agent-prompt.md](first-agent-prompt.md) for a copy-paste prompt to give a coding agent.
