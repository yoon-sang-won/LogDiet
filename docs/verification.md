# Verification

Use this document when local results, GitHub Actions, or raw GitHub URLs appear to disagree.

Release verification should check all of these:

1. Working tree.
2. Committed `HEAD`.
3. `origin/main`.
4. Fresh clone.
5. GitHub Actions.
6. Optional commit-specific raw URLs.

GitHub Actions is a useful signal, but `git show origin/main:<file>` and a fresh clone are the best ways to verify committed source content.

## 1. Working tree

```sh
git status --short
gofmt -w .
go test ./...
go install ./cmd/logdiet
```

## 2. Committed HEAD

```sh
git show HEAD:go.mod | cat -n
git show HEAD:cmd/logdiet/main.go | cat -n
git show HEAD:.github/workflows/test.yml | cat -n
git show HEAD:README.md | wc -l
git grep -n "package .* import" -- '*.go' || true
git grep -n "logdiet/internal" -- '*.go' || true
```

Expected:

- `go.mod` has separate `module` and `go` lines.
- `cmd/logdiet/main.go` is multiline Go source.
- workflow YAML is multiline.
- README has real Markdown line breaks.
- malformed `package ... import` matches are absent.
- old `logdiet/internal` imports are absent.

## 3. origin/main

```sh
git fetch origin
git show origin/main:go.mod | cat -n
git show origin/main:cmd/logdiet/main.go | cat -n
git show origin/main:.github/workflows/test.yml | cat -n
git show origin/main:README.md | wc -l
```

## 4. Fresh clone

```sh
tmpdir="$(mktemp -d)"
git clone https://github.com/yoon-sang-won/LogDiet "$tmpdir/LogDiet"
cd "$tmpdir/LogDiet"

gofmt -w .
git diff --exit-code
go test ./...
go install ./cmd/logdiet
```

## 5. Smoke test

```sh
tmpdir="$(mktemp -d)"
cd "$tmpdir"

logdiet install
eval "$(logdiet env)"
logdiet doctor

logdiet wrap -- sh -c 'echo ok; echo "panic: synthetic failure" >&2; exit 7'
code="$?"
test "$code" = "7"

logdiet raw latest
logdiet grep latest "panic"
```

## 6. GitHub Actions

Check the workflow page:

```text
https://github.com/yoon-sang-won/LogDiet/actions/workflows/test.yml
```

The workflow should run direct `go test ./...`.

## 7. GitHub raw URLs

Prefer `git show origin/main:<file>` as the source of truth. Raw GitHub URLs can be cached, but commit-specific raw URLs are useful for spot checks.
