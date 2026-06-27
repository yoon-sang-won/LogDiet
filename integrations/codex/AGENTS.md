# LogDiet instructions for Codex

When working in this repository:

- prefer `logdiet wrap -- <command>` for noisy test/build/search/git commands if automatic hooks are unavailable;
- read compact evidence first;
- use `logdiet show latest:F1 --around 40` to expand a handle;
- use `logdiet grep latest "pattern"` to search exact raw output;
- use `logdiet raw latest` only when compact evidence is insufficient;
- do not ask the user to paste full terminal logs unless LogDiet expansion commands are insufficient;
- keep final responses concise:
  - success: changed files, tests run, risks;
  - failure: exact error, likely cause, next command.
