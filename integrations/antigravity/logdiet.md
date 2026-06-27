# LogDiet rules for Antigravity

- rules fallback first;
- automatic command rewrite is not guaranteed;
- prefer `logdiet wrap -- <cmd>` for noisy test/build/search/git commands;
- use `logdiet show latest:F1 --around 40`, `logdiet grep latest "pattern"`, and `logdiet raw latest` for expansion;
- do not paste full log walls.

Install target: `.agents/rules/logdiet.md`.
