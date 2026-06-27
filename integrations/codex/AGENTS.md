# LogDiet for Codex

When working in this repository:

- use LogDiet for noisy commands such as tests, builds, git diffs, and search;
- prefer `logdiet wrap -- <command>` when automatic hook rewrite is unavailable;
- read compact evidence first;
- expand one handle with `logdiet show latest:F1 --around 40`;
- search raw output with `logdiet grep latest "pattern"`;
- print full raw output with `logdiet raw latest` only when compact evidence is insufficient;
- do not ask the user to paste full terminal logs unless LogDiet expansion commands are insufficient;
- keep final responses concise:
  - success: changed files, tests run, risks;
  - failure: exact error, likely cause, next command.
