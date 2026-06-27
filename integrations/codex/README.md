# LogDiet for Codex

LogDiet is not magically built into Codex. It works through repository instructions and command hooks where supported.

## Rules fallback

Use `AGENTS.md` or `logdiet-instructions.md` to tell Codex:

- use LogDiet for noisy commands such as tests, builds, git diffs, and search;
- prefer `logdiet wrap -- <command>` when automatic hook rewrite is unavailable;
- read compact evidence first;
- expand exact output with `logdiet show`, `logdiet grep`, and `logdiet raw`;
- avoid asking the user to paste full terminal logs.

## Hook rewrite template

`hook-rewrite-template.sh` shows how a Codex command hook could call:

```sh
logdiet hook rewrite --command "$COMMAND"
```

The hook rewrite template asks LogDiet for a decision. It does not execute the command itself.

Codex may require hook review/trust. If Codex asks for review, open `/hooks` and trust the generated LogDiet hook only after comparing it with your local files.

If hooks are not enabled or trusted, Codex should still use the rules fallback and run:

```sh
logdiet wrap -- <command>
```

full raw logs stay local under `.logdiet/runs/`. Compact evidence can be expanded with `show`, `grep`, and `raw`.

Automatic command rewriting is available where Codex command hooks are supported and trusted. Other environments use rules/instructions fallback or manual `logdiet wrap`.
