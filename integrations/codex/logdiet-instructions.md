# LogDiet Codex Instructions

Use LogDiet as an agent-native token diet layer backed by the local CLI engine.

- Prefer `logdiet wrap -- <command>` for noisy test/build/search/git commands when hooks are unavailable.
- Read compact evidence first.
- Expand exact raw output with `logdiet show latest:F1 --around 40`.
- Search raw output with `logdiet grep latest "pattern"`.
- Use `logdiet raw latest` only when compact evidence is insufficient.

Automatic command rewriting is available where the agent supports command hooks. Other agents use rules/instructions fallback or manual `logdiet wrap`.
