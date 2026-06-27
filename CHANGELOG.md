# Changelog

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
