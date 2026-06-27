#!/bin/sh
set -eu

step() {
	printf '\n==> %s\n' "$*"
}

fail() {
	printf 'error: %s\n' "$*" >&2
	exit 1
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

	fail "logdiet not found after go install"
}

check_wrap() {
	cmd=$1
	expected=$2
	# Exercise the same decision path as: logdiet hook rewrite --command "$cmd".
	json=$("$logdiet_bin" hook rewrite --command "$cmd")
	printf '%s -> %s\n' "$cmd" "$json"
	printf '%s\n' "$json" | grep -F "\"wrap\":$expected" >/dev/null 2>&1 || fail "$cmd wrap decision was not $expected"
	if [ "$expected" = "true" ]; then
		printf '%s\n' "$json" | grep -F "logdiet wrap --" >/dev/null 2>&1 || fail "$cmd did not rewrite to logdiet wrap"
	fi
}

step "check repository"
root=$(git rev-parse --show-toplevel)
cd "$root"

step "install local LogDiet engine"
go install ./cmd/logdiet
logdiet_bin=$(find_logdiet)

fixture=integrations/fixtures/commands.txt
test -f "$fixture" || fail "$fixture missing"

step "exercise all fixture commands"
while IFS= read -r cmd || [ -n "$cmd" ]; do
	[ -n "$cmd" ] || continue
	"$logdiet_bin" hook rewrite --command "$cmd" >/dev/null
done <"$fixture"

step "verify wrap decisions"
check_wrap 'go test ./...' true
check_wrap 'pytest -q' true
check_wrap 'npm test' true
check_wrap 'cargo test' true
check_wrap 'git diff' true
check_wrap 'rg "TODO"' true

step "verify no-wrap decisions"
check_wrap 'vim file.go' false
check_wrap 'less README.md' false
check_wrap 'ssh example.com' false
check_wrap 'python' false
check_wrap 'node' false
check_wrap 'logdiet wrap -- go test ./...' false

printf '\n%s\n' "Adapter fixture verification passed."
