# LogDiet Skill

- do not paste log walls;
- prefer LogDiet compact evidence;
- use `logdiet show latest:F1 --around 40`, `logdiet grep latest "pattern"`, and `logdiet raw latest` to expand evidence;
- keep success/failure reports short;
- use `logdiet wrap -- <cmd>` when hook rewrite is unavailable.

Automatic command rewriting is available where the agent supports command hooks. Other agents use rules/instructions fallback or manual `logdiet wrap`.
