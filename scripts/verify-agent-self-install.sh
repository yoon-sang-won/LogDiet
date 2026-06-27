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

step "install local LogDiet engine"
go install ./cmd/logdiet
logdiet_bin=$(find_logdiet)
"$logdiet_bin" --version

step "verify link-only onboarding docs"
for file in AI_INSTALL.md README.md docs/agent-self-install.md; do
	contains "$file" "logdiet bootstrap --agent auto"
	contains "$file" "logdiet wrap --"
	contains "$file" "logdiet show latest:F1 --around 40"
	contains "$file" "logdiet grep latest"
	contains "$file" "logdiet raw latest"
done
contains AI_INSTALL.md "logdiet agent-instructions --agent auto"
contains README.md "logdiet agent-instructions --agent auto"
contains AI_INSTALL.md "Native hooks are optional advanced mode."
contains README.md "Hooks are optional advanced mode."
contains docs/agent-self-install.md "Hooks are optional advanced mode."

step "Native hooks are not required for this verification."

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT HUP INT TERM

step "create temp project"
cd "$tmpdir"
git init >/dev/null 2>&1

step "simulate agent bootstrap"
# verifies: logdiet bootstrap --agent auto
"$logdiet_bin" bootstrap --agent auto
test -f .logdiet/LOGDIET_RULES.md || fail "generic rules were not created"
# verifies: logdiet doctor
PATH="$tmpdir/.logdiet/bin:$PATH" "$logdiet_bin" doctor
# verifies: logdiet agent-instructions --agent auto
"$logdiet_bin" agent-instructions --agent auto

step "verify generated generic rules"
contains .logdiet/LOGDIET_RULES.md "logdiet wrap -- <command>"
contains .logdiet/LOGDIET_RULES.md "logdiet show latest:F1 --around 40"
contains .logdiet/LOGDIET_RULES.md "logdiet grep latest"
contains .logdiet/LOGDIET_RULES.md "logdiet raw latest"
contains .logdiet/LOGDIET_RULES.md "do not ask the user to paste full logs"

step "run wrapped command"
# verifies: logdiet wrap -- sh -c 'printf "line 1\nline 2\n"; exit 0'
"$logdiet_bin" wrap -- sh -c 'printf "line 1\nline 2\n"; exit 0'
test -d .logdiet/runs || fail ".logdiet/runs was not created"
test -s .logdiet/latest || fail "latest run pointer was not created"
"$logdiet_bin" raw latest | grep -F "line 2" >/dev/null 2>&1 || fail "raw latest did not include line 2"
# verifies: logdiet grep latest "line 2"
"$logdiet_bin" grep latest "line 2" >/dev/null 2>&1 || fail "grep latest did not find line 2"

printf '\n%s\n' "Agent self-install verification passed."
