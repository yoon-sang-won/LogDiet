# Release Checklist

## Required Checks

```sh
gofmt -w .
go test ./...
go install ./cmd/logdiet
logdiet --version
logdiet help
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

## GitHub

- Actions are green.
- README install command is correct.
- No `.logdiet/runs` or `.logdiet/backup` files are committed.
- Tag release:

```sh
git tag v0.1.0
git push origin v0.1.0
```

## Repository Metadata

Suggested description:

```text
Put your coding agent on a token diet. Local logs, compact evidence.
```

Suggested topics:

```text
ai coding-agent developer-tools cli go token-optimization logs terminal codex claude-code cursor antigravity
```
