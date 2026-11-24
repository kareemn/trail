#!/bin/sh
# Phase 0: check that nono leaves the agent exactly one way out, the
# edge's localhost port. See ../README.md, "Phase 0".
#
# Usage:
#   spike/nono-confine.sh               run the probe inside nono
#   spike/nono-confine.sh --control     run it without nono (most checks
#                                       should FAIL: proves the probe works)
#
# The nono flags below are a first guess from nono's docs and are not yet
# verified. Override them with NONO_FLAGS to try others, for example:
#   NONO_FLAGS="--block-net --open-port 7777" spike/nono-confine.sh
set -eu

cd "$(dirname "$0")/.."
EDGE_ADDR=127.0.0.1:7777
OTHER_ADDR=127.0.0.1:7778
WORK=$(mktemp -d)
BIN="$WORK/bin"
mkdir -p "$BIN" "$WORK/edge-secrets"

go build -o "$BIN/escape-probe" ./cmd/escape-probe
echo "not for the agent" > "$WORK/edge-secrets/edge.key"
chmod 600 "$WORK/edge-secrets/edge.key"

"$BIN/escape-probe" -serve "$EDGE_ADDR" &
EDGE_PID=$!
"$BIN/escape-probe" -serve "$OTHER_ADDR" &
OTHER_PID=$!
trap 'kill $EDGE_PID $OTHER_PID 2>/dev/null; rm -rf "$WORK"' EXIT INT TERM
sleep 1

PROBE_ARGS="-edge $EDGE_ADDR -other-local $OTHER_ADDR -secret $WORK/edge-secrets/edge.key -edge-pid $EDGE_PID"

if [ "${1:-}" = "--control" ]; then
  echo "== control run, no sandbox =="
  # shellcheck disable=SC2086
  "$BIN/escape-probe" $PROBE_ARGS || true
  exit 0
fi

command -v nono >/dev/null || { echo "nono not found: brew install nono" >&2; exit 1; }
NONO_FLAGS=${NONO_FLAGS:-"--block-net --open-port 7777 --allow $BIN"}
echo "== nono $(nono --version 2>/dev/null || echo '?') with: $NONO_FLAGS =="
# shellcheck disable=SC2086
nono run $NONO_FLAGS -- "$BIN/escape-probe" $PROBE_ARGS
