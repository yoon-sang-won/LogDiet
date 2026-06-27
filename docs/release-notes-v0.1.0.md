# LogDiet v0.1.0

First public release of LogDiet.

LogDiet is a local token-diet layer for AI coding agents. It keeps full command logs locally and feeds agents compact, expandable evidence instead of noisy terminal walls.

## Highlights

- Captures full command output under `.logdiet/runs/`
- Prints compact evidence with expansion handles
- Expands exact raw output with `logdiet show`, `logdiet raw`, and `logdiet grep`
- Adds project-local PATH shims via `logdiet install`
- Adds setup flows for:
  - Codex
  - Claude Code
  - Cursor
  - Antigravity
  - Gemini
- Adds `logdiet doctor` for setup verification
- Adds instruction lint and optional response contract rules
- Adds synthetic fixture benchmarks
- Adds release verification docs and script

## Install

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
```

## Quickstart

```sh
logdiet install
eval "$(logdiet env)"
logdiet doctor
logdiet wrap -- go test ./...
```

PowerShell:

```powershell
logdiet install
Invoke-Expression (logdiet env --shell powershell)
logdiet doctor
logdiet wrap -- go test ./...
```

## Agent setup

```sh
logdiet setup codex
logdiet setup claude
logdiet setup cursor
logdiet setup antigravity
logdiet setup gemini
```

## Privacy

- No network calls
- No telemetry
- No model/API calls
- Raw logs stay local

## Verification

```sh
./scripts/verify-release.sh
```

## Notes

Token estimates are approximate and based on local byte counts. They are not provider billing measurements.
