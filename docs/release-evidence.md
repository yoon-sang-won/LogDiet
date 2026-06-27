# Release Evidence

This document records the checks expected before publishing a LogDiet release.

## Required local checks

```sh
gofmt -w .
go test ./...
go install ./cmd/logdiet
logdiet --version
logdiet help
logdiet bench-fixtures
```

## Smoke test

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

## Source integrity

```sh
git show HEAD:go.mod | cat -n
git show HEAD:cmd/logdiet/main.go | cat -n
git show HEAD:.github/workflows/test.yml | cat -n
git show HEAD:README.md | wc -l

git grep -n "package .* import" -- '*.go' || true
git grep -n "logdiet/internal" -- '*.go' || true
```

## Expected result

- Go tests pass.
- CLI installs.
- Smoke test preserves wrapped exit code.
- Raw output is stored locally.
- `raw`, `grep`, and `doctor` work.
- Source files are multiline and valid.
- GitHub Actions runs direct `go test ./...`.
