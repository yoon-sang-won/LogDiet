# LogDiet

Put your coding agent on a token diet.

LogDiet keeps full command logs locally and feeds AI coding agents compact, expandable evidence instead of noisy terminal walls.

## What LogDiet Is

LogDiet is a local command I/O layer for AI coding-agent sessions. It captures exact command output under `.logdiet/runs/`, prints compact evidence to the agent, and gives expansion commands for raw lines when more context is needed.

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

Every wrapped command stores `stdout.txt`, `stderr.txt`, `combined.txt`, `meta.json`, and `index.json`. The terminal receives a short summary with evidence handles such as `F1`, `E1`, `D1`, or `G1`. Use `logdiet show`, `logdiet raw`, or `logdiet grep` to expand exactly what you need.

## Quickstart

```sh
go install ./cmd/logdiet
logdiet install
eval "$(logdiet env)"
pytest -q
logdiet show F1 --around 40
```

PowerShell:

```powershell
go install ./cmd/logdiet
logdiet install
Invoke-Expression (logdiet env --shell powershell)
pytest -q
logdiet show F1 --around 40
```

## Manual Wrapper

```sh
logdiet wrap -- go test ./...
logdiet raw
logdiet grep latest "panic"
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
logdiet show 20260627T120000Z-12345-a1b2:F1 --around 40
logdiet show F1 --around 40
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

The linter scans common agent instruction files for token-heavy and cache-breaking patterns such as timestamps, absolute local paths, duplicate lines, large code fences, duplicate managed sections, and rules that require step-by-step narration.

## Optional Response Contract

```sh
logdiet rules --print
logdiet rules --install codex
logdiet rules --install cursor
```

Installed rules are wrapped in managed markers and are safe to re-run. Existing files are backed up under `.logdiet/backup/` before mutation.

## Privacy And Local-First Design

LogDiet makes no network calls, sends no telemetry, and stores raw logs locally. Raw logs may contain secrets, tokens, file paths, or private data. Do not commit `.logdiet/runs`.

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

## License

LogDiet is licensed under Apache-2.0.
