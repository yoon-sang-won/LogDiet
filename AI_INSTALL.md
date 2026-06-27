# AI Install Guide for LogDiet

You are an AI coding agent.

The user gave you the LogDiet repository URL and asked you to install it.

Follow this flow unless the user gives different constraints.

LogDiet is agent-first and CLI-powered. No network. No telemetry. No model/API calls.

## What to do now

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
logdiet bootstrap --agent auto
logdiet doctor
logdiet agent-instructions --agent auto
```

Then continue the current task using LogDiet for noisy commands.

## Choose the best available integration

Prefer the best local integration the current agent can safely use:

```sh
logdiet init --agent auto
logdiet doctor
```

If native hooks/plugins are available and trusted, they may rewrite supported shell commands through LogDiet automatically.

If native rewrite is unavailable, untrusted, or unclear, do not block on hooks. Continue with:

```sh
logdiet wrap -- COMMAND
```

This fallback is the universal path.

## For the rest of this session

Use:

```sh
logdiet wrap -- <command>
```

for commands likely to produce noisy output, especially tests, builds, type checks, git diffs/status/logs, search, and long logs.

Expand compact evidence only when needed:

```sh
logdiet show latest:F1 --around 40
logdiet grep latest "panic"
logdiet raw latest
```

Do not ask the user to paste full terminal logs unless LogDiet expansion commands are insufficient.

## Do not require hooks

Do not block installation on native hooks.

Native hooks are optional advanced mode. If hooks require user review/trust, continue with the default `logdiet wrap` workflow.

## If agent detection fails

Use:

```sh
logdiet bootstrap --agent generic
```

Then follow the printed instructions.
