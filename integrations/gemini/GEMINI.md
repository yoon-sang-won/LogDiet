# LogDiet instructions

- use compact evidence from LogDiet before expanding raw output;
- prefer `logdiet wrap -- <cmd>` when automatic hooks are unavailable;
- use `logdiet show latest:F1 --around 40` for focused expansion;
- use `logdiet grep latest "pattern"` for exact raw search;
- use raw expansion only when needed;
- avoid full log pastes.

Automatic command rewriting is available where the agent supports command hooks. Other agents use rules/instructions fallback or manual `logdiet wrap`.
