# LogDiet for Codex

LogDiet is not magically built into Codex. It works through repository instructions and command hooks where supported.

## Status

- rules fallback: supported through `AGENTS.md` or `logdiet-instructions.md`
- native adapter: hook rewrite template
- transparent rewrite: partial, requires supported and trusted Codex hooks
- trust required: yes, review hooks in `/hooks` before enabling them

## Recommended setup

```sh
logdiet init --agent codex --mode all
logdiet doctor
```

If Codex asks for hook review, open `/hooks` and trust the generated LogDiet hook only after comparing it with your local files.

## Default fallback behavior

Without trusted hooks, Codex should still use the rules fallback and run:

```sh
logdiet wrap -- <command>
```

This is the reliable path for noisy tests, builds, git diffs, and search.

## How it works

`hook-rewrite-template.sh` shows how a Codex command hook could call:

```sh
logdiet hook rewrite --command "$COMMAND"
```

The hook rewrite template asks LogDiet for a decision. It does not execute the command itself.

## Verification

```sh
./scripts/verify-codex-integration.sh
```

This verifies LogDiet-side files and rewrite decisions. Runtime hook trust must be verified manually inside Codex with `/hooks`.

## Limitations

- Built-in Codex file/search/editor tools may bypass shell hooks.
- Automatic rewrite applies only to command execution paths Codex exposes to trusted hooks.
- full raw logs stay local under `.logdiet/runs/`; compact evidence can be expanded with `show`, `grep`, and `raw`.
