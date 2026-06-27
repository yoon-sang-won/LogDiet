#!/bin/sh
set -eu

# LogDiet Claude Code hook rewrite template.
# Adapt this file to the exact Claude Code hook protocol before enabling it.
# This template asks LogDiet for a rewrite decision and does not execute commands.

: "${COMMAND:?COMMAND is required}"
if ! command -v logdiet >/dev/null 2>&1; then
	printf '%s\n' "logdiet not found; install LogDiet or use logdiet wrap -- COMMAND manually" >&2
	exit 127
fi
logdiet hook rewrite --command "$COMMAND"
