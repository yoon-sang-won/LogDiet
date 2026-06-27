#!/bin/sh
set -eu

step() {
	printf '\n==> %s\n' "$*"
}

fail() {
	printf 'error: %s\n' "$*" >&2
	exit 1
}

contains() {
	file=$1
	text=$2
	grep -F "$text" "$file" >/dev/null 2>&1 || fail "$file missing: $text"
}

find_logdiet() {
	if command -v logdiet >/dev/null 2>&1; then
		command -v logdiet
		return 0
	fi

	step "install local logdiet"
	go install ./cmd/logdiet

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

step "check repository"
root=$(git rev-parse --show-toplevel)
cd "$root"

logdiet_bin=$(find_logdiet)
step "using $logdiet_bin"
"$logdiet_bin" --version

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT HUP INT TERM

step "create temp repository"
cd "$tmpdir"
git init >/dev/null 2>&1

step "setup codex rules"
# verifies: logdiet setup codex --mode rules
"$logdiet_bin" setup codex --mode rules
test -f AGENTS.md || fail "AGENTS.md was not created"

step "setup codex all"
# verifies: logdiet setup codex --mode all
"$logdiet_bin" setup codex --mode all
test -d .logdiet || fail ".logdiet was not created"
test -d .logdiet/bin || fail ".logdiet/bin was not created"
test -f .logdiet/integrations/codex/hook-rewrite-template.sh || fail "Codex hook template was not created"

step "verify AGENTS.md"
contains AGENTS.md "logdiet wrap"
contains AGENTS.md "logdiet show latest:F1 --around 40"
contains AGENTS.md "logdiet grep latest"
contains AGENTS.md "logdiet raw latest"
contains AGENTS.md "do not ask the user to paste full terminal logs"

step "verify hook rewrite JSON"
# verifies: logdiet hook rewrite --command "go test ./..."
wrap_json=$("$logdiet_bin" hook rewrite --command "go test ./...")
printf '%s\n' "$wrap_json"
printf '%s\n' "$wrap_json" | grep -F '"wrap":true' >/dev/null 2>&1 || fail "go test was not selected for wrapping"

nowrap_json=$("$logdiet_bin" hook rewrite --command "vim file.go")
printf '%s\n' "$nowrap_json"
printf '%s\n' "$nowrap_json" | grep -F '"wrap":false' >/dev/null 2>&1 || fail "vim should not be wrapped"

step "doctor"
PATH="$tmpdir/.logdiet/bin:$PATH" "$logdiet_bin" doctor

step "codex integration verification passed"
printf '%s\n' "Codex native hook template verified. Runtime trust must be verified manually in Codex with /hooks."
