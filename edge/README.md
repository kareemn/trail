# trail-edge: a TRAIL edge for nono sandboxes

This folder holds an implementation, in progress, of a TRAIL edge for
agents that run in a [nono](https://nono.sh/) sandbox. The protocol is
defined by the draft in the parent folder,
[draft-nassar-agentproto-trail](../draft-nassar-agentproto-trail.md).
Section numbers below refer to it. When this code and the draft
disagree, the draft wins, or the difference goes to the draft's issue
tracker.

Status: Phase 0. Nothing in the edge works yet. Licensed under
[Apache-2.0](LICENSE).

## How it fits together

nono confines the agent so that its only way out is a localhost port
owned by the edge. The edge runs outside the sandbox, as its own
process and ideally its own user. It holds every key, signs every call,
and keeps the trail. That makes every path out of the cell pass through
the edge, as the draft's sandbox definition (1.6) and core rule (4)
require.

```
+------------- Cell A (one host) -----------------+
|  nono run (Landlock / Seatbelt)                 |
|  +-----------+   http://127.0.0.1:EDGE/<route>  |
|  |  Agent    |-------------------+              |
|  +-----------+   listen-port <---+--+           |
|                                  v  |           |
|  Edge (own uid, holds the keys) -------> HTTPS + Trail-Record
|   trail db, state cache, policy     <------- inbound calls
+------------------------|------------------------+
                         v
        evidence plane: tlog tiles + state + inbox (HTTP),
        MOQT later
```

The edge is a process beside the sandbox on the same host, so its
attestation level is `software` (2.1). A host compromise or a sandbox
escape goes around it, and a TPM-held key does not change that. Running
nono and the edge inside a confidential VM gives `tee` with no change
to this code.

### Why the edge does not use nono's proxy

nono's network proxy allows or denies CONNECT by domain. TLS passes
through it unopened, so an edge behind it could not see, sign, or bind
a record to the call. The edge instead copies nono's credential-route
idea: the agent gets plain-HTTP base URLs that point at the edge, such
as `MCP_URL=http://127.0.0.1:EDGE/mcp/github` or
`OPENAI_BASE_URL=http://127.0.0.1:EDGE/model/openai`. The edge parses
the call, signs it, opens TLS to the real endpoint, and adds the
credentials itself. It must add them itself, because the record carries
their hashes in context-hashes (3.2).

| nono feature | Use in the edge |
|---|---|
| Kernel filesystem deny-by-default | The edge's keys, trail database, and config stay outside every allowed path |
| Network limited to one localhost port | The only way out. Not yet verified: see Phase 0 |
| `--listen-port` | Inbound calls: the edge verifies Trail-Record, then forwards to the agent |
| Credential proxy | Not used. The edge adds credentials |
| Supervisor approvals, audit Merkle log, undo | Not in the draft. At most a lifecycle record could cite nono's session root, as an extension |
| Process exit | The edge writes a lifecycle record, and admission-release if the agent held a token |

## What the edge contains

1. **Launcher.** Builds the nono profile, starts `nono run` with the
   route environment variables and no secrets, and watches the process.
2. **Outbound pipeline (4, 6.1).** Classifies each call (kind, action
   class, counterparty, MCP or A2A method from JSON-RPC), checks it
   against directives, builds the Action Record (COSE_Sign1, SD-CWT
   payload), appends it to the trail and waits for the acknowledgment,
   signs the request with RFC 9421 and adds Trail-Record, sends it, then
   verifies the responder's record and appends a record citing it.
   Model endpoints get model-call records only. Any other destination
   is denied.
3. **Inbound pipeline.** The six receiver steps of section 4: verify the
   record and call-ref, check the sender's level, holds, admission, and
   freshness, reject replays, deliver, and append a recv record. SSE
   responses carry the record in the event that completes the exchange.
4. **Trail (8).** An RFC 9162 tree in SQLite, published as tlog tiles.
   A checkpoint closes each epoch and goes to the witnesses for
   cosignature. A restart record starts a new trail after lost state.
5. **State (9).** Publishes the edge's snapshot: its WIT, the agent's
   WIT, the Attestation Result, and the latest cosigned checkpoint.
   Subscribes to the operator's, the authorities', and peers' snapshots,
   ignores directives outside a grant, and tracks freshness per
   authority and action class.
6. **Directive engine (9.2 to 9.5).** All eight directive types, degraded
   mode with grace periods from the last fresh snapshot, and a
   hold-applied record the first time a hold stops an action.
7. **Directed records (5.3, 10).** Gate requests and approvals,
   disclosures encrypted with COSE-HPKE, disclosure requests (refused
   and recorded outside a grant), and admission tokens.

The edge is written in Go. Go has tlog-tiles (Tessera), C2SP checkpoint
and witness tooling, go-cose, RFC 9421 libraries, a nono SDK, and MOQT
libraries. Whether a Go MOQT library supports draft-22 is still to be
checked.

## Phases

Each phase is done when the named blocks and scenarios from the
explainer page ([TRAIL_visuals.html](../TRAIL_visuals.html), rendered
in the [repo README](../README.md)) run as tests.

| Phase | Builds | Done when these work |
|---|---|---|
| 0. Confinement | Confirm nono leaves one way out, on macOS and Linux (Landlock ABI v4+) | `spike/nono-confine.sh` passes |
| 1. Core rule over HTTP | Records, trail, signatures, two edges, a mock MCP or A2A peer | "Records on both sides", "tamper-evident" |
| 2. State and directives | Operator grants; policy, hold, require-record, require-level, freshness-budget, degraded mode | Grants, killswitch over HTTP, scenario 4 |
| 3. Witnesses | Witness service, cosignatures, require-witness, discovery through witness trails | Scenarios 1 and 5, discovery |
| 4. Gates, disclosure, admission | Directed records, SD-CWT, COSE-HPKE, admission tokens | Scenarios 2, 6 (without passkeys), 7 |
| 5. MOQT evidence plane | Hybrid deployment (6.6): calls on HTTP, evidence on MOQT | Killswitch and discovery over MOQT |
| 6. Measurements | The numbers listed in Open Issues | Results for Implementation Status (13) |

Phases 1 to 4 also need small mock services: an operator, an authority,
a witness, an admission issuer, a test verifier, and a test WIT issuer.

## Phase 0: confinement

`cmd/escape-probe` runs inside the sandbox and tries every way out it
knows: the edge port (must work), another local port, IPv4 and IPv6
TCP, UDP DNS, the system resolver, the cloud metadata address, reading
a file the edge owns, and signalling the edge's process. Everything but
the edge port must fail.

```sh
spike/nono-confine.sh --control   # no sandbox: most checks FAIL, proving the probe works
spike/nono-confine.sh             # inside nono
NONO_FLAGS="--block-net --open-port 7777" spike/nono-confine.sh
```

The default nono flags are a guess from nono's documentation and have
not been run yet. If no combination passes, the gap goes to nolabs as a
feature request: an option that limits a child to one localhost port.

On some networks UDP DNS and the metadata address are already blocked,
so those two checks pass even in the control run. They prove something
only on a network where the control run reaches them.

## Gaps the draft leaves to implementations

These come up as soon as the edge is built. Each should become an issue
on the draft.

- No policy bundle format. The edge has to pick one, such as Cedar or
  Rego.
- No standard way to map an HTTP call to an action class and
  counterparty. The edge needs a configuration schema for it.
- Bootstrapping HTTP state: the state resource is "listed in its
  state", which is circular on first contact.
- The commit point over HTTP (6.4) requires an acknowledged append
  before every send, from a tile server the operator does not run. That
  adds latency to every call; batching needs guidance.
- Which inputs a send cites. The edge cannot tell which received input
  caused a later send. Scenario 3 needs a rule, such as citing every
  input in the dialog since the last send.
- Whether pausing covered agents means only refusing their traffic, or
  also suspending the process.
- nono's own evidence (supervisor approvals, audit log) has no place in
  the record.

## Layout

| Path | What it is |
|---|---|
| `cmd/trail-edge` | The edge. A stub for now |
| `cmd/escape-probe` | Phase 0 probe, also a stand-in edge with `-serve` |
| `spike/nono-confine.sh` | Runs the Phase 0 check |

Build with Go 1.25 or later: `go build ./...` from this folder.
