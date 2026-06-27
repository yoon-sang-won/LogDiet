# Changelog

## v0.2.0 - Unreleased

### Changed

- Repositioned LogDiet as an agent-native token diet layer powered by a local CLI engine.

### Added

- Agent integration packages under `integrations/`.
- Command rewrite decision helper.
- `logdiet hook rewrite` for hook/plugin adapters.
- Agent-native documentation.
- Setup modes for rules, shims, and native templates.
- Doctor output for agent integration status.

### Notes

Automatic command rewriting is available where an agent supports command hooks. Other agents use rules/instructions fallback or manual `logdiet wrap`.

## v0.1.0 - Initial public release

LogDiet v0.1.0 is the first public release of LogDiet, a local token-diet layer for AI coding agents.

### Added

- Local command output capture under `.logdiet/runs/`
- Compact, expandable evidence output for AI coding agents
- `logdiet wrap -- <cmd>`
- `logdiet raw`, `logdiet show`, and `logdiet grep`
- Project-local PATH shims via `logdiet install`
- Agent setup flows:
  - Codex via `AGENTS.md`
  - Claude Code via `CLAUDE.md`
  - Cursor via `.cursor/rules/logdiet.mdc`
  - Antigravity via `.agents/rules/logdiet.md`
  - Gemini via `GEMINI.md`
- `logdiet doctor`
- `logdiet lint-instructions`
- Optional LogDiet response contract rules
- Synthetic fixture benchmarks via `logdiet bench-fixtures`
- Release verification docs and script

### Privacy

- No network calls
- No telemetry
- No model/API calls
- Raw logs stay local

### Notes

Token estimates are approximate and based on local byte counts. They are not provider billing measurements.
