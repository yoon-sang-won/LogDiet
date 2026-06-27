# Changelog

## v0.2.0 - Unreleased

### Changed

- Repositioned LogDiet as an agent-native token diet layer powered by a local CLI engine.
- README now leads with the agent self-install path.
- Native hooks are documented as optional advanced mode, not the default requirement.
- README and README.ko.md now surface the agent self-install flow earlier.
- `AI_INSTALL.md`, `bootstrap`, and `agent-instructions` now more clearly tell agents to continue with `logdiet wrap` without requiring hooks.
- Integration READMEs now distinguish native adapters, rules fallback, and manual CLI usage.
- `doctor` now reports native adapter status more clearly.
- Documentation now frames hooks as best-effort native automation with wrapper fallback everywhere.

### Added

- Agent integration packages under `integrations/`.
- Command rewrite decision helper.
- `logdiet hook rewrite` for hook/plugin adapters.
- `AI_INSTALL.md` for agents installing LogDiet from a GitHub link.
- `logdiet bootstrap` for agent self-install flows.
- `logdiet agent-instructions` for current-session operating rules.
- Agent self-install documentation.
- Tests for bootstrap and agent instruction flows.
- `scripts/verify-agent-self-install.sh` for hookless self-install verification.
- `docs/first-agent-prompt.md` with a copy-paste prompt for coding agents.
- Agent-native documentation.
- Native adapter architecture documentation.
- `logdiet init` entrypoint for agent integration setup/status.
- Adapter contract documentation.
- Adapter fixture verification script.
- Agent support matrix.
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
