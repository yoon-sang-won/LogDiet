# Codex Integration Verification

This document verifies the LogDiet Codex integration.

## What is verified automatically

The repository tests verify:

- `logdiet hook rewrite --command "go test ./..."` returns `wrap: true`;
- `logdiet hook rewrite --command "vim file.go"` returns `wrap: false`;
- `logdiet setup codex --mode rules` creates Codex rules;
- `logdiet setup codex --mode all` creates rules, shims, and native templates where available;
- generated rules tell Codex to use `logdiet show`, `logdiet grep`, and `logdiet raw`.

## What requires manual verification

Codex hooks may require review/trust in the Codex UI.

Manual verification:

1. Install LogDiet:

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
```

2. In a test repository:

```sh
logdiet setup codex --mode all
logdiet doctor
```

3. Start Codex from that repository:

```sh
codex
```

4. If Codex reports hook review is required, open:

```text
/hooks
```

Review and trust the LogDiet hook if it matches the local files you generated.

5. Ask Codex:

```text
Run the test suite.
```

6. Verify LogDiet captured a run:

```sh
ls .logdiet/runs
logdiet raw latest
logdiet grep latest "panic"
```

## Expected behavior

Where Codex command hooks are enabled and trusted, noisy commands such as:

```sh
go test ./...
pytest -q
npm test
git diff
rg "TODO"
```

should be routed through LogDiet or produce LogDiet compact evidence.

If hooks are unavailable or not trusted, Codex should still follow `AGENTS.md` and prefer:

```sh
logdiet wrap -- <command>
```

## Truthful limitation

Automatic command rewriting is available only where Codex command hooks are supported and trusted. Otherwise LogDiet works through rules/instructions fallback or manual wrapper mode.
