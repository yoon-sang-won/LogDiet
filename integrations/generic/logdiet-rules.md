# LogDiet generic rules

- Prefer compact LogDiet evidence for noisy command output.
- Do not paste full logs.
- Use `logdiet wrap -- <cmd>` when automatic rewrite is unavailable.
- Expand with `logdiet show latest:F1 --around 40`.
- Search exact raw output with `logdiet grep latest "pattern"`.
- Use `logdiet raw latest` only when compact evidence is insufficient.
