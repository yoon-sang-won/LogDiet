# LogDiet

<p align="center">
  <a href="./README.md">English</a> |
  <a href="./README.ko.md">한국어</a>
</p>

<p align="center">
  <strong>Agent-native token diet for coding agents.</strong>
</p>

<p align="center">
  LogDiet rewrites noisy terminal commands into compact, expandable evidence while keeping full raw logs local.
</p>

<p align="center">
  Agent-first. CLI-powered. No network. No telemetry.
</p>

<p align="center">
  <a href="https://github.com/yoon-sang-won/LogDiet/actions/workflows/test.yml"><img alt="test" src="https://github.com/yoon-sang-won/LogDiet/actions/workflows/test.yml/badge.svg"></a>
  <a href="./LICENSE"><img alt="License" src="https://img.shields.io/badge/license-Apache--2.0-blue.svg"></a>
  <img alt="Go" src="https://img.shields.io/badge/Go-1.22+-00ADD8">
  <img alt="No Network" src="https://img.shields.io/badge/network-none-brightgreen">
  <img alt="No Telemetry" src="https://img.shields.io/badge/telemetry-none-brightgreen">
</p>

No network. No telemetry. No model/API calls.

## Why

Coding agents need command evidence, not terminal walls. Long test logs, diffs, search output, and stack traces consume context while hiding the lines that matter.

LogDiet keeps the complete raw output on disk and gives the agent a compact report with handles for exact expansion.

## Before / After

### Before

```text
pytest -q
... thousands of lines of traceback, warnings, retries, and progress output ...
... repeated stack frames ...
... unrelated warnings ...
... the actual failure is buried somewhere above ...
```

### After

```text
logdiet run 20260627T120000Z-12345-a1b2 exit=1 raw=.logdiet/runs/20260627T120000Z-12345-a1b2
cmd: pytest -q
summary: 2 failed, 31 passed
F1 tests/test_api.py:42 AssertionError: expected 200, got 500
F2 tests/test_auth.py:17 ValueError: missing token
show: logdiet show latest:F1 --around 40
raw:  logdiet raw latest
grep: logdiet grep latest "pattern"
stats: raw=18420B compact=610B approx_saved=96.7%
```

This example is synthetic. `approx_saved` is a byte-based reduction estimate, not a provider billing measurement.

## How LogDiet works

LogDiet has two layers:

1. Agent integration layer:
   - plugin / skill / rules / hook packages for coding agents;
   - teaches agents not to paste log walls;
   - rewrites noisy commands where hooks are supported.

2. Local CLI engine:
   - runs `logdiet wrap -- <cmd>`;
   - stores raw logs under `.logdiet/runs/`;
   - prints compact evidence;
   - expands exact output with `show`, `grep`, and `raw`.

```mermaid
flowchart LR
    A[Coding agent] --> B[Rules / skill / hook template]
    B --> C[logdiet local engine]
    C --> D[Compact evidence]
    C --> E[Raw logs in .logdiet/runs]
    D --> F[show / grep / raw]
```

Automatic command rewriting is available where the agent supports command hooks. Other agents use rules/instructions fallback or manual `logdiet wrap`.

## Quickstart: install for your agent

### Codex

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
logdiet setup codex --mode all
logdiet doctor
codex
```

### Claude Code

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
logdiet setup claude --mode all
logdiet doctor
claude
```

### Other agents

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
logdiet setup cursor --mode rules
logdiet setup gemini --mode rules
logdiet setup antigravity --mode rules
logdiet doctor
```

### Manual engine mode

```sh
logdiet wrap -- go test ./...
logdiet show latest:F1 --around 40
logdiet grep latest "panic"
logdiet raw latest
```

`@latest` works best after a release tag exists.

## Hook rewrite bridge

Agents with trusted command hooks can call:

```sh
logdiet hook rewrite --command "go test ./..."
```

Example output:

```json
{"wrap":true,"command":"logdiet wrap -- go test ./...","reason":"known noisy developer command"}
```

The bridge only returns a decision. It does not execute commands.

## Works with

Integration packages live under `integrations/`:

- Codex: `integrations/codex/`
- Claude Code: `integrations/claude-code/`
- Cursor: `integrations/cursor/`
- Gemini: `integrations/gemini/`
- Antigravity: `integrations/antigravity/`
- Generic terminal agents: `integrations/generic/`

See [docs/agent-native.md](docs/agent-native.md) for the v0.2 architecture.

## Core commands

```sh
logdiet install
logdiet setup codex --mode rules
logdiet setup codex --mode shim
logdiet setup codex --mode native
logdiet setup codex --mode all
logdiet doctor
logdiet wrap -- pytest -q
logdiet show latest:F1 --around 40
logdiet raw latest
logdiet grep latest "pattern"
logdiet hook rewrite --command "go test ./..."
logdiet bench-fixtures
```

## Setup modes

| Mode | Behavior |
| ---- | -------- |
| `rules` | Installs agent rules/instructions only |
| `shim` | Installs rules plus local `.logdiet/bin` PATH shims |
| `native` | Installs rules plus local native hook/plugin templates |
| `all` | Installs rules, shims, and native templates |

Native templates are reviewable files. LogDiet does not silently enable risky hooks.

## Privacy and security

- Raw logs stay local under `.logdiet/runs/`.
- Hooks can change command execution, so review generated hook templates before enabling them.
- Raw logs may contain secrets, tokens, private paths, or proprietary output.
- Do not commit `.logdiet/runs/` or `.logdiet/backup/`.

## What LogDiet is not

LogDiet is not:

- a model proxy;
- a prompt compressor;
- a cloud service;
- a telemetry collector;
- a daemon;
- a web UI;
- a benchmark claiming exact provider-token savings.

## Verification

```sh
gofmt -w .
go test ./...
go install ./cmd/logdiet
./scripts/verify-release.sh
```

For v0.2 checks, see [docs/v0.2-verification.md](docs/v0.2-verification.md).

## License

LogDiet is licensed under Apache-2.0.
