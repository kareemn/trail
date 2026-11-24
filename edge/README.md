# trail-edge: a TRAIL edge for nono's OS sandbox

This folder holds an implementation, in progress, of a TRAIL edge for
agents confined by [nono's OS sandbox](https://nono.sh/os-sandbox)
(Landlock on Linux, Seatbelt on macOS). The protocol is defined by the
draft in the parent folder,
[draft-nassar-agentproto-trail](../draft-nassar-agentproto-trail.md).
Section numbers below refer to it. When this code and the draft
disagree, the draft wins, or the difference goes to the draft's issue
tracker.

Licensed under [Apache-2.0](LICENSE).

## Status

| Phase | State |
|---|---|
| 0. Confinement | Passes on macOS (Seatbelt). Not yet run on Linux |
| 1. Core rule over HTTP | Works between two edges, without state, SD-CWT, or streaming |
| 2 to 6 | Not started |

## Quick start

Needs Go 1.25 or later with CGo (nono-go links nono's static library),
on macOS or Linux. No nono CLI is needed.

```sh
go build -o bin/ ./cmd/...

# Phase 0: the agent's only way out is the edge.
bin/trail-edge run -config examples/probe.json -- \
  bin/escape-probe -secret ~/.trail-edge/probe/edge.key

# Phase 1: a sandboxed agent calls another agent through two edges.
examples/two-edges.sh
```

`examples/two-edges.sh` runs curl as Agent A in the sandbox. Agent A
calls Agent B through both edges, then tries `https://example.com`
directly. The output ends with both trails:

```
== Edge A's trail ==
    0  lifecycle   completed event=start
    1  send        allowed   method=POST host=127.0.0.1:9201 path=/tasks jsonrpc-method=tasks/send aud=spiffe://example.test/edge/b
    2  recv        completed status=200 aud=spiffe://example.test/edge/b cites=e22966572224
    3  send        denied    target=CONNECT example.com:443 reason=tunnels are not allowed; use a route
    4  lifecycle   completed event=exit exit_code=0
    5  checkpoint              covers 5 records, root 3935348749bbef21

== Edge B's trail ==
    0  recv        allowed   method=POST path=/tasks jsonrpc-method=tasks/send aud=spiffe://example.test/edge/a cites=93f7666a37df
    1  send        completed status=200 aud=spiffe://example.test/edge/a cites=93f7666a37df
    2  checkpoint              covers 2 records, root d3c65fd8a189482a
```

A1 is Edge A's record 1. B1 and B2 cite it, and A2 cites B2, as in the
page's "records on both sides" block.

## How it fits together

```
+------------- Cell A (one host) -----------------+
|  OS sandbox (nono: Landlock / Seatbelt)         |
|  +-----------+   http://127.0.0.1:EDGE/<route>  |
|  |  Agent    |-------------------+              |
|  +-----------+                   |              |
|                                  v              |
|  trail-edge (outside the sandbox, holds keys) ----> HTTPS + Trail-Record
|   trail, peers, routes              <---------- inbound calls
+-------------------------------------------------+
```

`trail-edge run` starts the edge, then starts a copy of itself that
applies nono's sandbox to itself with nono-go and execs the agent. nono
cannot sandbox another process, and its sandbox cannot be undone, so
this is the only way to keep the edge outside while the agent never
runs unsandboxed. The sandbox:

* allows network access only to the edge's port (nono's
  `NetworkProxyOnly` mode);
* allows read access to system paths and what the config lists, and
  write access only where the config says; and
* is refused at startup if it would cover the edge's state directory,
  which holds its key and trail.

The agent's environment holds a few ordinary variables, the ones the
config passes through, one variable per route pointing at the edge (for
example `AGENT_B_URL=http://127.0.0.1:EDGE/a2a/b`), and proxy variables
pointing at the edge. No credential reaches it: the edge adds them
itself from its own environment, and records their hashes in
context-hashes (3.2).

The edge is a process beside the sandbox on the same host, so its
attestation level is `software` (2.1). A host compromise or a sandbox
escape goes around it, and a TPM-held key does not change that. Running
the edge and the sandbox inside a confidential VM gives `tee` with no
change to this code.

### Why the edge does not use nono's proxy

nono's proxy tunnels TLS for host allowlisting, and intercepts it only
for routes that need path-level rules. Either way, nono is the party
that sees the call. A TRAIL edge must see, sign, and bind a record to
every call, and hold the keys. So the agent talks plain HTTP to the
edge's own port, and the edge opens TLS to the real endpoint. A CONNECT
to the edge is refused and recorded.

## What works (Phase 1)

* **Outbound (4, 6.1).** Each call to a route is written to the trail as
  an Action Record before it leaves: a COSE_Sign1 with ES256, typ,
  kid (JWK Thumbprint), and CWT claims, and a deterministically encoded
  payload. The record carries call-ref (http-ref), aud, action class,
  context-hashes, and the JSON-RPC method, id, and tool name. The edge
  signs the request with an RFC 9421 signature covering `@method`,
  `@target-uri`, `content-digest`, and `trail-record`. Calls with no
  route, and CONNECT, are refused and recorded.
* **Responses.** When the route names a TRAIL counterparty (`aud`), the
  response must carry the responder's record, signed and citing the
  request's record. The edge records a recv that cites it. Otherwise it
  records a recv with outcome `unsigned`.
* **Model endpoints.** A `model-call` route gets a record and no TRAIL
  headers (4).
* **Inbound (4).** The edge refuses calls with no record, an unknown
  key, a bad signature, a call-ref or Content-Digest that does not match
  the call, an aud for someone else, a created time outside the allowed
  skew, or a replayed record. It records every refusal. It delivers
  accepted calls to the agent, records a recv citing the sender's
  record, and answers with its own signed record citing it.
* **Trail (8.1, 8.2).** An RFC 9162 tree over record identifiers,
  appended with fsync before a call leaves. A checkpoint closes each
  epoch and every run. `trail-edge verify` checks every signature, trail
  position, and checkpoint root. A changed byte fails it.
* **Lifecycle records** at agent start (command, nono version, and a
  hash of the sandbox rules) and exit (exit code).

Tests in `internal/` cover these, including tampering, replay, and
redirection. Run them with `go test ./...`.

## Not done yet

* **State and directives (9).** Peers' keys come from the config file as
  a stand-in for their state snapshots. No policy, hold, freshness, or
  require-* directives yet, so receiver step 4 (level, holds, admission,
  freshness) is skipped.
* **Witnesses and cosignatures (8.3),** publishing the trail as tlog
  tiles (6.2), and restart records after lost state (8.1).
* **Gates, admission, disclosure (10), and SD-CWT redaction (3.4).**
  Until redaction exists, records hold digests and protocol identifiers,
  not arguments.
* **Streamed responses.** SSE responses are buffered, which breaks long
  MCP streams (6.1).
* **The MOQT binding (5).**
* **WIMSE identity (11.1)** and Attestation Results (11.4). Identities
  are plain strings; there is no WIT.

## Linux

This hasn't been run on Linux yet. Two things to check there:

* Landlock's network rules (ABI v4, Linux 6.7+) match TCP ports, not
  addresses. Allowing the edge's port may allow that port on any host.
  nono's CLI adds a seccomp filter on `connect()` for this; whether
  nono-go's `NetworkProxyOnly` does too is not documented. The escape
  probe does not catch it yet, because it would need a remote host
  listening on the edge's port.
* Landlock applies per thread. `__sandbox-exec` locks its OS thread,
  applies the sandbox, and execs on that thread.

## Phases

Each phase is done when the named blocks and scenarios from the
explainer page ([TRAIL_visuals.html](../TRAIL_visuals.html), rendered
in the [repo README](../README.md)) run as tests.

| Phase | Builds | Done when these work |
|---|---|---|
| 0. Confinement | nono leaves one way out, on macOS and Linux | `escape-probe` passes under `trail-edge run` |
| 1. Core rule over HTTP | Records, trail, signatures, two edges | "Records on both sides", "tamper-evident" |
| 2. State and directives | Operator grants; policy, hold, require-record, require-level, freshness-budget, degraded mode | Grants, killswitch over HTTP, scenario 4 |
| 3. Witnesses | Witness service, cosignatures, require-witness, discovery through witness trails | Scenarios 1 and 5, discovery |
| 4. Gates, disclosure, admission | Directed records, SD-CWT, COSE-HPKE, admission tokens | Scenarios 2, 6 (without passkeys), 7 |
| 5. MOQT evidence plane | Hybrid deployment (6.6): calls on HTTP, evidence on MOQT | Killswitch and discovery over MOQT |
| 6. Measurements | The numbers listed in Open Issues | Results for Implementation Status (13) |

Phases 2 to 4 also need small mock services: an operator, an authority,
a witness, an admission issuer, a test verifier, and a test WIT issuer.

## Gaps the draft leaves to implementations

These come up as soon as the edge is built. Each should become an issue
on the draft.

* No policy bundle format. The edge has to pick one, such as Cedar or
  Rego.
* No standard way to map an HTTP call to an action class and
  counterparty. This edge uses routes in its config.
* Bootstrapping HTTP state: the state resource is "listed in its
  state", which is circular on first contact.
* The commit point over HTTP (6.4) requires an acknowledged append
  before every send, from a tile server the operator does not run. That
  adds latency to every call; batching needs guidance.
* Which inputs a send cites. The edge cannot tell which received input
  caused a later send. Scenario 3 needs a rule, such as citing every
  input in the dialog since the last send.
* Whether pausing covered agents means only refusing their traffic, or
  also suspending the process.
* How a call to a counterparty that does not run TRAIL is closed. This
  edge records a recv with outcome `unsigned`; the draft defines that
  outcome only for inbound calls.
* The sandbox's own rules have no place in the record. This edge puts a
  hash of them in a lifecycle record body.
* `@target-uri` is signed as the sender saw it. A receiver behind a
  proxy that changes Host or scheme cannot reconstruct it. The draft
  should say which URI the receiver checks.

## Layout

| Path | What it is |
|---|---|
| `cmd/trail-edge` | The edge: `run`, `serve`, `key`, `verify` |
| `cmd/escape-probe` | Phase 0 probe; with `-serve`, a stand-in edge |
| `cmd/echo-agent` | A stand-in agent or tool for demos and tests |
| `internal/edge` | Outbound and inbound pipelines, checkpoints, trail verification |
| `internal/record` | Action Records: COSE_Sign1, payload, identifiers |
| `internal/trail` | The append-only trail and RFC 9162 tree |
| `internal/httpsig` | RFC 9421 signatures and Trail-Record for TRAIL |
| `internal/sandbox` | Starting the agent in nono's OS sandbox |
| `internal/keys` | The edge key and JWK Thumbprints |
| `internal/config` | Edge configuration |
| `examples/` | `probe.json` and `two-edges.sh` |
| `spike/nono-confine.sh` | The same probe under the nono CLI instead of the library |
