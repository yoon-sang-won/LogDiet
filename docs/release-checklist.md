# Release Checklist

## Required Checks

```sh
./scripts/verify-release.sh
gofmt -w .
go test ./...
go install ./cmd/logdiet
logdiet --version
logdiet help
logdiet bench-fixtures
```

## Smoke Test

```sh
tmpdir="$(mktemp -d)"
cd "$tmpdir"

logdiet install
eval "$(logdiet env)"
logdiet doctor

logdiet wrap -- sh -c 'echo ok; echo "panic: synthetic failure" >&2; exit 7'
test "$?" = "7"

logdiet raw latest
logdiet grep latest "panic"
```

## Fresh Clone Verification

```sh
tmpdir="$(mktemp -d)"
git clone https://github.com/yoon-sang-won/LogDiet "$tmpdir/LogDiet"
cd "$tmpdir/LogDiet"
./scripts/verify-release.sh
```

## GitHub

- Actions are green.
- `.github/workflows/test.yml` runs direct `go test ./...`.
- Fresh clone verification passes.
- README install command is correct.
- No `.logdiet/runs` or `.logdiet/backup` files are committed.
- Run a real dogfood test in at least one Go, Python, or Node repo.
- Run at least one agent-specific setup flow.
- Set repository description.
- Set repository topics.
- Tag release:

```sh
./scripts/verify-release.sh
git tag v0.1.0
git push origin v0.1.0
```

- Create GitHub Release from tag `v0.1.0`.
- Paste `docs/release-notes-v0.1.0.md` into the release body.

## Repository Metadata

Suggested description:

```text
Put your coding agent on a token diet. Local logs, compact evidence.
```

Suggested topics:

```text
ai
coding-agent
developer-tools
cli
go
token-optimization
logs
terminal
codex
claude-code
cursor
antigravity
```
