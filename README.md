# LogDiet

[![test](https://github.com/yoon-sang-won/LogDiet/actions/workflows/test.yml/badge.svg)](https://github.com/yoon-sang-won/LogDiet/actions/workflows/test.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

Put your coding agent on a token diet.

LogDiet keeps full command logs locally and feeds AI coding agents compact, expandable evidence instead of noisy terminal walls.

## Try It In 60 Seconds

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
logdiet install
eval "$(logdiet env)"
logdiet wrap -- go test ./...
logdiet raw latest --tail 40
logdiet grep latest "panic"
```

PowerShell activation:

```powershell
Invoke-Expression (logdiet env --shell powershell)
```

## What LogDiet Is

LogDiet is a local command I/O reduction tool for AI coding-agent sessions. It captures exact command output under `.logdiet/runs/`, prints compact evidence to the agent, and gives expansion commands for exact raw lines when more context is needed.

The core primitive is lossless local command-output capture with compact expandable evidence handles.

## What LogDiet Is Not

- Not a clone of any existing token-saving tool.
- Not a model proxy.
- Not a prompt compressor.
- Not an AI summarizer.
- Not a cloud service.
- Not a telemetry product.
- Not a replacement for provider prompt caching.
- Not a tool that discards logs.

## Core Idea: Compact Evidence, Full Raw Logs

Every wrapped command stores `stdout.txt`, `stderr.txt`, `combined.txt`, `meta.json`, and `index.json`. The terminal receives a short summary with evidence handles such as `F1`, `E1`, `D1`, or `G1`.

Use `logdiet show`, `logdiet raw`, or `logdiet grep` to expand exactly what you need.

## Agent Quickstarts

### Codex

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
logdiet setup codex
eval "$(logdiet env)"
logdiet doctor
codex
```

`logdiet setup codex` installs local shims and writes managed LogDiet rules to `AGENTS.md`. Run `logdiet doctor` in the same terminal/session where Codex will run.

### Claude Code

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
logdiet setup claude
eval "$(logdiet env)"
logdiet doctor
claude
```

`logdiet setup claude` writes managed LogDiet rules to `CLAUDE.md`. Run `logdiet doctor` in the same shell/session.

### Cursor

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
logdiet setup cursor
eval "$(logdiet env)"
logdiet doctor
```

`logdiet setup cursor` writes `.cursor/rules/logdiet.mdc`. Verify the Cursor agent terminal/session sees the same `PATH`; `logdiet doctor` should show `.logdiet/bin` in `PATH`.

### Antigravity

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
logdiet setup antigravity
eval "$(logdiet env)"
logdiet doctor
```

`logdiet setup antigravity` writes `.agents/rules/logdiet.md`. Verify with `logdiet doctor`.

### Generic Terminal Agents

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
logdiet install
eval "$(logdiet env)"
logdiet rules --print
logdiet doctor
```

## Manual Wrapper

```sh
logdiet wrap -- go test ./...
logdiet raw latest
logdiet grep latest "panic"
logdiet show latest:F1 --around 40
```

## PATH Shims

`logdiet install` creates local command shims in `.logdiet/bin`. Prepend that directory to `PATH` inside the agent session. The shims resolve the real command outside `.logdiet/bin`, set `LOGDIET_ACTIVE=1` for the child process, and preserve the command exit code.

Controls:

- `LOGDIET_BYPASS=1` runs the real command directly.
- `LOGDIET_MODE=auto` compacts known useful commands.
- `LOGDIET_MODE=force` compacts every shimmed command.
- `LOGDIET_MODE=off` bypasses compaction.

No shell profiles are modified in v0.1.

## Raw Expansion

```sh
logdiet show latest:F1 --around 40
logdiet raw latest --combined --tail 80
logdiet grep latest "AssertionError" --around 3
```

Raw output is not redacted. Compact output may shorten noise, but the raw files remain available locally.

## Instruction Lint

```sh
logdiet lint-instructions
logdiet lint-instructions --json
logdiet lint-instructions --fix
```

The linter scans common agent instruction files for token-heavy and cache-breaking patterns such as timestamps, absolute local paths, duplicate lines, large code fences, duplicate managed sections, and rules that require step-by-step narration. It also scans `.agents/rules/*.md` for Antigravity-style workspace rules.

## Optional Response Contract

```sh
logdiet rules --print
logdiet rules --install codex
logdiet rules --install claude
logdiet rules --install cursor
logdiet rules --install antigravity
logdiet rules --install gemini
```

Installed rules are wrapped in managed markers and are safe to re-run. Existing files are backed up under `.logdiet/backup/` before mutation.

## Doctor

```sh
logdiet doctor
```

`doctor` checks the current shell/session: binary path, current directory, `.logdiet` state, shim directory, `PATH`, installed shims, real command resolution, LogDiet environment variables, latest run, and installed agent rule files.

Run it in the same terminal or agent session where commands will execute.

## Privacy And Local-First Design

LogDiet makes no network calls, sends no telemetry, and stores raw logs locally. Raw logs may contain secrets, tokens, credentials, private file paths, proprietary code, or other sensitive data. Do not commit `.logdiet/runs`.

Compact output may still include snippets from raw logs.

Provider prompt caching can reduce cost or latency for repeated prefixes. LogDiet complements that by keeping local instruction files stable and by preventing huge local command outputs from entering the agent context in the first place. It does not replace provider caching.

## Design Boundaries

See [docs/design-boundaries.md](docs/design-boundaries.md). LogDiet is independently implemented Apache-2.0 code and uses synthetic fixtures.

## Fixture Benchmarks

```sh
logdiet bench-fixtures
```

Fixture benchmarks use synthetic local logs. Byte reduction is measured exactly. Token count is approximate using bytes divided by four. Benchmark numbers are not provider billing measurements.

## Limitations

- `combined.txt` appends stdout then stderr in v0.1, so cross-stream ordering is best-effort.
- Compaction is deterministic pattern extraction, not semantic understanding.
- Parsers focus on common failure shapes and may miss unusual tool output.
- Windows command lookup supports `.exe`, `.cmd`, `.bat`, and `.com`.
- No daemon, TUI, editor extension, MCP server, model proxy, or cloud dashboard is included.

## Uninstall

```sh
logdiet uninstall
logdiet uninstall --rules
```

Uninstall removes managed shims and optionally managed response-rule sections. It does not delete `.logdiet/runs` by default.

## Development

```sh
go install ./cmd/logdiet
gofmt -w .
go test ./...
```

LogDiet is standard-library-only Go. Do not add network calls, telemetry, or third-party dependencies for v0.1.

## Release

Before public announcement, maintainers should create a release tag:

```sh
git tag v0.1.0
git push origin v0.1.0
```

Suggested GitHub repository description:

```text
Put your coding agent on a token diet. Local logs, compact evidence.
```

Suggested topics:

```text
ai coding-agent developer-tools cli go token-optimization logs terminal codex claude-code cursor antigravity
```

## License

LogDiet is licensed under Apache-2.0.
