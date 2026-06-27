# Announcement Drafts

## Short post

```md
I just released LogDiet v0.1.0.

It puts AI coding agents on a token diet: full command logs stay local, while the agent sees compact, expandable evidence.

Works with Codex, Claude Code, Cursor, Antigravity, and terminal-based agents.

No network. No telemetry. No model/API calls.

https://github.com/yoon-sang-won/LogDiet
```

## Developer-focused post

```md
LogDiet v0.1.0 is out.

It is a local command I/O layer for AI coding agents.

Problem: agents often eat huge terminal outputs: test logs, build logs, diffs, grep results, tracebacks.

LogDiet keeps the full logs under `.logdiet/runs/` and feeds the agent compact evidence with handles:

- `logdiet show latest:F1 --around 40`
- `logdiet raw latest`
- `logdiet grep latest "panic"`

It also includes project-local PATH shims and setup flows for Codex, Claude Code, Cursor, Antigravity, and Gemini.

No network. No telemetry. No model/API calls.

https://github.com/yoon-sang-won/LogDiet
```

## Very short tagline

```md
LogDiet: Put your coding agent on a token diet.

Keep full logs locally. Feed agents compact evidence.
```
