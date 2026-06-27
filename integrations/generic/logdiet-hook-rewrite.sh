#!/bin/sh
set -eu

# Generic LogDiet hook rewrite adapter.
# Agent hook input/output contracts differ, so adapt this template before use.
# It must not execute the command itself; it only asks LogDiet for a decision.

: "${COMMAND:?COMMAND is required}"
logdiet hook rewrite --command "$COMMAND"
