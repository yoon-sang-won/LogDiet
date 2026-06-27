# First Agent Prompt

Copy this into your coding agent:

```text
Install https://github.com/yoon-sang-won/LogDiet and use it for noisy test/build/git/search output.

After installing, use LogDiet for the rest of this session:
- run noisy commands with `logdiet wrap -- <command>`;
- expand compact evidence with `logdiet show latest:F1 --around 40`;
- search raw output with `logdiet grep latest "pattern"`;
- print full raw output with `logdiet raw latest` only when needed;
- do not ask me to paste full terminal logs unless LogDiet expansion commands are insufficient.

Hooks are optional. Do not block on hook setup.
```
