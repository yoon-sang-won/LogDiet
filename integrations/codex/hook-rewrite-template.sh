#!/bin/sh
set -eu

# LogDiet Codex hook rewrite template.
# Adapt this file to the exact Codex hook protocol before enabling it.
# This template asks LogDiet for a rewrite decision and does not execute commands.

: "${COMMAND:?COMMAND is required}"
logdiet hook rewrite --command "$COMMAND"
