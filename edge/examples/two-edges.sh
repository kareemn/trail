#!/bin/sh
# Two edges on one machine, as in the page's "records on both sides" block.
#
#   Agent A (curl, in nono's OS sandbox) -> Edge A -> Edge B -> Agent B (echo)
#
# Agent A calls B once through its route, then tries to reach the internet
# directly. The first call carries records on both sides; the second is
# refused and recorded. Both trails are then verified.
set -eu
cd "$(dirname "$0")/.."
go build -o bin/ ./cmd/...
WORK=$(mktemp -d)
trap 'kill $(jobs -p) 2>/dev/null; rm -rf "$WORK"' EXIT INT TERM

mkdir -p "$WORK/a" "$WORK/b"
cat > "$WORK/a/edge.json" <<JSON
{ "identity": "spiffe://example.test/edge/a", "agent": "spiffe://example.test/agent/a",
  "state_dir": "$WORK/a/state", "routes": [], "sandbox": {} }
JSON
cat > "$WORK/b/edge.json" <<JSON
{ "identity": "spiffe://example.test/edge/b", "agent": "spiffe://example.test/agent/b",
  "state_dir": "$WORK/b/state", "routes": [], "sandbox": {} }
JSON
PEER_A=$(bin/trail-edge key -config "$WORK/a/edge.json")
PEER_B=$(bin/trail-edge key -config "$WORK/b/edge.json")

cat > "$WORK/a/edge.json" <<JSON
{
  "identity": "spiffe://example.test/edge/a",
  "agent": "spiffe://example.test/agent/a",
  "state_dir": "$WORK/a/state",
  "epoch_records": 8,
  "routes": [
    { "prefix": "/a2a/b", "upstream": "http://127.0.0.1:9201", "kind": "send",
      "aud": "spiffe://example.test/edge/b", "env": "AGENT_B_URL" }
  ],
  "peers": [ $PEER_B ],
  "sandbox": { "workdir": false }
}
JSON
cat > "$WORK/b/edge.json" <<JSON
{
  "identity": "spiffe://example.test/edge/b",
  "agent": "spiffe://example.test/agent/b",
  "state_dir": "$WORK/b/state",
  "inbound": { "listen": "127.0.0.1:9201", "agent": "127.0.0.1:9202" },
  "routes": [],
  "peers": [ $PEER_A ]
}
JSON

bin/echo-agent -listen 127.0.0.1:9202 -name agent-b 2>/dev/null &
bin/trail-edge serve -config "$WORK/b/edge.json" 2>"$WORK/b/log" &
sleep 1

echo "== Agent A, in the sandbox, calls Agent B and then tries the internet =="
bin/trail-edge run -config "$WORK/a/edge.json" -- /bin/sh -c '
  echo "-> call B through the edge:"
  curl -s -X POST "$AGENT_B_URL/tasks" -H "Content-Type: application/json" \
    -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tasks/send\",\"params\":{\"text\":\"hello B\"}}"
  echo
  echo "-> try https://example.com directly:"
  curl -s -o /dev/null -w "%{http_code}\n" --max-time 5 https://example.com || echo "refused"
' 2>"$WORK/a/log"

kill %2 2>/dev/null; sleep 1
echo
echo "== Edge A's trail =="
bin/trail-edge verify -config "$WORK/a/edge.json"
echo
echo "== Edge B's trail =="
bin/trail-edge verify -config "$WORK/b/edge.json"
