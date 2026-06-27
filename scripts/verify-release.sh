#!/bin/sh
set -eu

step() {
	printf '\n==> %s\n' "$*"
}

run() {
	step "$*"
	"$@"
}

check_no_matches() {
	pattern=$1
	step "git grep -n \"$pattern\" -- '*.go'"
	set +e
	matches=$(git grep -n "$pattern" -- '*.go')
	status=$?
	set -e
	if [ "$status" -eq 0 ]; then
		printf '%s\n' "$matches"
		return 1
	fi
	if [ "$status" -eq 1 ]; then
		return 0
	fi
	return "$status"
}

find_logdiet() {
	if command -v logdiet >/dev/null 2>&1; then
		command -v logdiet
		return 0
	fi

	goexe=$(go env GOEXE)
	gobin=$(go env GOBIN)
	if [ -n "$gobin" ] && [ -x "$gobin/logdiet$goexe" ]; then
		printf '%s\n' "$gobin/logdiet$goexe"
		return 0
	fi

	gopath=$(go env GOPATH)
	if [ -n "$gopath" ] && [ -x "$gopath/bin/logdiet$goexe" ]; then
		printf '%s\n' "$gopath/bin/logdiet$goexe"
		return 0
	fi

	printf '%s\n' "logdiet"
}

step "check git repository"
root=$(git rev-parse --show-toplevel)
cd "$root"

step "current commit"
git rev-parse HEAD

run git status --short
run gofmt -w .
run go test ./...
run go install ./cmd/logdiet

check_no_matches "package .* import"
check_no_matches "logdiet/internal"

step "line counts"
wc -l \
	go.mod \
	cmd/logdiet/main.go \
	internal/cli/cli.go \
	README.md \
	.github/workflows/test.yml \
	docs/release-checklist.md \
	docs/verification.md

logdiet_bin=$(find_logdiet)
run "$logdiet_bin" --version
run "$logdiet_bin" help
run "$logdiet_bin" bench-fixtures

step "temp-dir smoke test"
tmpdir=$(mktemp -d)
(
	cd "$tmpdir"
	run "$logdiet_bin" install
	eval "$("$logdiet_bin" env --shell sh)"
	run "$logdiet_bin" doctor

	step "$logdiet_bin wrap -- sh -c 'echo ok; echo \"panic: synthetic failure\" >&2; exit 7'"
	set +e
	"$logdiet_bin" wrap -- sh -c 'echo ok; echo "panic: synthetic failure" >&2; exit 7'
	code=$?
	set -e
	test "$code" = "7"

	run "$logdiet_bin" raw latest
	run "$logdiet_bin" grep latest "panic"
)

step "release verification passed"
