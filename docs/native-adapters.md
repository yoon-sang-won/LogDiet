# Native Agent Adapters

## Goal

LogDiet should provide the best available integration for each coding agent:

1. native hook/plugin adapter where supported and trusted;
2. rules plus explicit `logdiet wrap -- COMMAND` fallback everywhere else;
3. manual CLI usage as the final fallback.

Core principle:

```text
Native where possible. Fallback everywhere. Raw logs always local.
```

## Architecture

```text
Agent shell command
  -> native adapter hook/plugin, if available and trusted
  -> `logdiet hook rewrite --command "COMMAND"`
  -> rewritten command, usually `logdiet wrap -- COMMAND`
  -> compact evidence to the agent
  -> full raw logs stored locally under `.logdiet/runs`
```

## Thin adapter rule

Agent-specific adapters should be thin delegates.

They may:

- parse the agent's hook payload;
- extract the shell command;
- call `logdiet hook rewrite`;
- return the modified command in the agent's expected format.

They should not:

- duplicate rewrite policy;
- summarize logs themselves;
- store logs themselves;
- make network calls;
- contain agent-independent decision logic.

## Central rewrite engine

The Go binary owns command selection through `internal/agentrewrite` and `logdiet hook rewrite`.

## Integration tiers

| Tier | Meaning | Example behavior |
| --- | --- | --- |
| Native adapter | Hook/plugin can rewrite shell commands before execution | Agent runs `go test ./...`; adapter rewrites to `logdiet wrap -- go test ./...` |
| Rules fallback | Agent reads rules and explicitly calls LogDiet | Agent runs `logdiet wrap -- go test ./...` |
| Manual CLI | Human or agent calls LogDiet directly | `logdiet wrap -- COMMAND` |

## Supported agents matrix

| Agent | Rules fallback | Native adapter files | Transparent command rewrite | Requires trust/review | Verification status | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| Codex | yes | template | partial | yes | LogDiet-side verified; runtime not verified | Generated hook template delegates to `logdiet hook rewrite`. |
| Claude Code | yes | template | template | yes | not verified | Hook protocol must be adapted and trusted in Claude Code. |
| Cursor | yes | template | template | yes | not verified | Rules fallback is available; native hook behavior is not claimed. |
| Gemini CLI | yes | template | template | yes | not verified | Instructions and hook template are local scaffolding. |
| Antigravity | yes | template | unknown | unknown | not verified | Rules fallback is the reliable path. |
| Generic | yes | not applicable | not supported yet | no | supported fallback | Use rules or explicit `logdiet wrap -- COMMAND`. |

## Scope of transparent rewrite

Transparent rewrite only applies to shell/terminal command execution paths that the agent exposes to hooks/plugins.

Built-in file tools, search tools, editor tools, or non-shell actions may bypass LogDiet.

When in doubt, agents should use:

```sh
logdiet wrap -- COMMAND
```

## Safety

LogDiet never silently enables command-rewriting hooks.

If an agent requires hook review/trust, users must complete that agent's trust flow.

## Fallback behavior

If native rewrite is unavailable, agents should continue through rules and explicit wrapper mode.

## Verification

Native adapter verification has two levels:

1. LogDiet-side verification:
   - adapter files exist;
   - adapter calls `logdiet hook rewrite`;
   - fixtures prove command rewrite behavior;
   - `logdiet doctor` reports status.

2. Agent-runtime verification:
   - run inside the actual agent;
   - enable/trust hook if required;
   - observe a normal shell command being rewritten;
   - confirm `.logdiet/runs` contains the captured run.

Do not claim level 2 unless performed.
