#!/bin/sh
set -eu

# Generic LogDiet hook rewrite adapter.
# Agent hook input/output contracts differ, so adapt this template before use.
# It must not execute the command itself; it only asks LogDiet for a decision.

: "${COMMAND:?COMMAND is required}"
if ! command -v logdiet >/dev/null 2>&1; then
	printf '%s\n' "logdiet not found; install LogDiet or use logdiet wrap -- COMMAND manually" >&2
	exit 127
fi
logdiet hook rewrite --command "$COMMAND"
