---
title: "TRAIL: Tamper-evident Records for Agent Interactions and Liveness over Media over QUIC"
abbrev: "TRAIL"
category: std

docname: draft-nandakumar-agentproto-trail-latest
submissiontype: IETF
number:
date: 2026-10-29
consensus: false
v: 3
area: "Applications and Real-Time"
workgroup: "Agent Communication Protocols"
keyword:
 - agent evidence
 - oversight
 - attestation
 - transparency
 - MOQT
 - WIMSE

venue:
  group: "Agent Communication Protocols"
  type: "Working Group"
  mail: "agentproto@ietf.org"
  arch: "https://mailarchive.ietf.org/arch/browse/agentproto/"

author:
 -
    name: "Suhas Nandakumar"
    organization: "Cisco"
    email: "snandaku@cisco.com"

normative:
  RFC8392:
  RFC8610:
  RFC8949:
  RFC9052:
  RFC9162:
  RFC9338:
  RFC9449:
  RFC9596:
  RFC9597:
  RFC9942:
  RFC9943:
  MoQ-TRANSPORT:
    title: "Media over QUIC Transport"
    date: 2026-10
    author:
    - ins: S. Nandakumar
    - ins: V. Vasiliev
    - ins: I. Swett
      role: editor
    - ins: A. Frindell
      role: editor
    seriesinfo:
      Internet-Draft: draft-ietf-moq-transport-22
    target: https://datatracker.ietf.org/doc/draft-ietf-moq-transport/
  MoQ-C4M:
    title: "Authorization scheme for MOQT using Common Access Tokens"
    date: 2026
    author:
    - ins: W. Law
    - ins: C. Lemmons
    - ins: G. Simon
    - ins: S. Nandakumar
    seriesinfo:
      Internet-Draft: draft-ietf-moq-c4m-01
    target: https://datatracker.ietf.org/doc/draft-ietf-moq-c4m/
  MoQ-SECOBJ:
    title: "End-to-End Secure Objects for Media over QUIC Transport"
    date: 2026
    author:
    - org: IETF MoQ Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-moq-secure-objects-01
    target: https://datatracker.ietf.org/doc/draft-ietf-moq-secure-objects/
  WIMSE-CREDS:
    title: "WIMSE Workload Credentials"
    date: 2026
    author:
    - org: IETF WIMSE Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-wimse-workload-creds-02
    target: https://datatracker.ietf.org/doc/draft-ietf-wimse-workload-creds/
  OAUTH-STATUS-LIST:
    title: "Token Status List (TSL)"
    date: 2026
    author:
    - org: IETF OAuth Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-oauth-status-list-21
    target: https://datatracker.ietf.org/doc/draft-ietf-oauth-status-list/
  SD-CWT:
    title: "SPICE SD-CWT"
    date: 2026
    author:
    - org: IETF SPICE Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-spice-sd-cwt
    target: https://datatracker.ietf.org/doc/draft-ietf-spice-sd-cwt/
  COSE-HPKE:
    title: "Use of Hybrid Public-Key Encryption (HPKE) with CBOR Object Signing and Encryption (COSE)"
    date: 2026
    author:
    - org: IETF COSE Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-cose-hpke
    target: https://datatracker.ietf.org/doc/draft-ietf-cose-hpke/

informative:
  RFC3552:
  RFC6973:
  RFC8693:
  RFC9334:
  RFC9711:
  RFC9999:
  RFC9396:
  RFC9421:
  RFC9562:
  RFC9901:
  PACE:
    title: "PACE: Protocol for Agent Communication Exchange"
    date: 2026
    author:
    - ins: S. Nandakumar
    - ins: C. Jennings
    seriesinfo:
      Internet-Draft: draft-nandakumar-agentproto-moq-pace
    target: https://datatracker.ietf.org/doc/draft-nandakumar-agentproto-moq-pace/
  MCP-MOQT:
    title: "Model Context Protocol and Agent Skills over Media over QUIC Transport"
    date: 2026
    author:
    - ins: C. Jennings
    - ins: I. Swett
    - ins: J. Rosenberg
    - ins: S. Nandakumar
    seriesinfo:
      Internet-Draft: draft-jennings-agentproto-mcp-over-moqt
    target: https://datatracker.ietf.org/doc/draft-jennings-agentproto-mcp-over-moqt/
  SD-CARD:
    title: "Selective Disclosure for Agent Discovery and Identity Management"
    date: 2026-02
    author:
    - ins: S. Nandakumar
    - ins: C. Jennings
    seriesinfo:
      Internet-Draft: draft-nandakumar-agent-sd-jwt
    target: https://datatracker.ietf.org/doc/draft-nandakumar-agent-sd-jwt/
  WIMSE-ARCH:
    title: "Workload Identity in a Multi System Environment (WIMSE) Architecture"
    date: 2026
    author:
    - org: IETF WIMSE Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-wimse-arch-08
    target: https://datatracker.ietf.org/doc/draft-ietf-wimse-arch/
  WIMSE-WPT:
    title: "WIMSE Workload Proof Token"
    date: 2026
    author:
    - org: IETF WIMSE Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-wimse-wpt-02
    target: https://datatracker.ietf.org/doc/draft-ietf-wimse-wpt/
  WIMSE-HTTPSIG:
    title: "WIMSE Workload-to-Workload Authentication with HTTP Signatures"
    date: 2026
    author:
    - org: IETF WIMSE Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-wimse-http-signature-08
    target: https://datatracker.ietf.org/doc/draft-ietf-wimse-http-signature/
  WIMSE-AIMS:
    title: "AI Identity Management System"
    date: 2026
    author:
    - org: IETF WIMSE Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-wimse-aims-00
    target: https://datatracker.ietf.org/doc/draft-ietf-wimse-aims/
  WIMSE-ATTEST:
    title: "WIMSE Workload Attestation"
    date: 2026
    author:
    - ins: T. Reddy
    seriesinfo:
      Internet-Draft: draft-reddy-wimse-workload-attestation-00
    target: https://datatracker.ietf.org/doc/draft-reddy-wimse-workload-attestation/
  TXN-TOKENS:
    title: "Transaction Tokens"
    date: 2026
    author:
    - org: IETF OAuth Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-oauth-transaction-tokens-11
    target: https://datatracker.ietf.org/doc/draft-ietf-oauth-transaction-tokens/
  RATS-AR4SI:
    title: "Attestation Results for Secure Interactions"
    date: 2026
    author:
    - org: IETF RATS Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-rats-ar4si-10
    target: https://datatracker.ietf.org/doc/draft-ietf-rats-ar4si/
  RATS-EAR:
    title: "EAT Attestation Results"
    date: 2026
    author:
    - org: IETF RATS Working Group
    seriesinfo:
      Internet-Draft: draft-ietf-rats-ear-04
    target: https://datatracker.ietf.org/doc/draft-ietf-rats-ear/
  WIMSE-AUTHZ-EVIDENCE:
    title: "Authorization Evidence for Workloads"
    date: 2026
    author:
    - ins: Munoz
    seriesinfo:
      Internet-Draft: draft-munoz-wimse-authorization-evidence-01
    target: https://datatracker.ietf.org/doc/draft-munoz-wimse-authorization-evidence/
  WIMSE-AGENT-EVIDENCE:
    title: "Agent Evidence Binding for WIMSE"
    date: 2026
    author:
    - ins: Thaha
    seriesinfo:
      Internet-Draft: draft-thaha-wimse-agent-evidence-binding-01
    target: https://datatracker.ietf.org/doc/draft-thaha-wimse-agent-evidence-binding/
  SCITT-LOCAL-LOG:
    title: "Checkpointed Local Logs for SCITT"
    date: 2026
    author:
    - ins: Mih
    seriesinfo:
      Internet-Draft: draft-mih-scitt-checkpointed-local-log-02
    target: https://datatracker.ietf.org/doc/draft-mih-scitt-checkpointed-local-log/
  AGENTPROTO-VERIFIER:
    title: "Verifier Evidence for Agent Protocols"
    date: 2026
    author:
    - ins: He
    seriesinfo:
      Internet-Draft: draft-he-agentproto-verifier-evidence-01
    target: https://datatracker.ietf.org/doc/draft-he-agentproto-verifier-evidence/
  AAUTH:
    title: "AAuth Protocol"
    date: 2026
    author:
    - ins: D. Hardt
    seriesinfo:
      Internet-Draft: draft-hardt-oauth-aauth-protocol-11
    target: https://datatracker.ietf.org/doc/draft-hardt-oauth-aauth-protocol/
  A2A:
    title: "Agent2Agent Protocol Specification"
    author:
    - org: A2A Project
    target: https://a2a-protocol.org/latest/specification/
  MCP:
    title: "Model Context Protocol Specification"
    author:
    - org: Anthropic
    target: https://modelcontextprotocol.io/specification/

--- abstract

TRAIL produces verifiable evidence of AI agent actions and supports
live oversight of them, using Media over QUIC Transport (MOQT) to carry
both the agent calls and the oversight. An enforcement point outside the
agent signs each call as a COSE Action Statement. That one statement
authenticates the call, serves as the receiver's receipt, and is
appended to a tamper-evident log.

--- middle

# Introduction

An agent is a model running in a loop with tools, credentials, and a
goal. It reads untrusted input, acts with real authority, and works
faster than people can watch. Two things are needed outside the agent:

1. An enforcement point that the agent cannot see, reach, or turn off.
2. A record of what the agent did that the operator cannot quietly
   rewrite, that parties outside the operator can check, and that
   those parties can act on while the agent is still running.

Single-node sandboxes and gateways provide the first. The second is
missing as soon as one agent's call reaches an agent on another node.

## Problem Statement

Agent calls cross operator boundaries, but the controls and the records
for those calls do not. Today:

* A receiver cannot verify, per call, which workload sent it, on whose
  behalf, under which policy, and whether it was altered in transit by
  an intermediary such as a relay.
* Records of agent actions are written by the same software and
  operator whose conduct they describe, so they can be dropped,
  reordered, or rewritten without detection.
* Parties with oversight rights cannot observe activity while it
  happens, pause or stop it, or require approval before an action
  takes effect.
* Limits that apply across many nodes, such as a cap on concurrent
  agents in a swarm, cannot be enforced because no node sees the total.
* Sharing records with one party exposes them to every party, because
  there is no way to disclose part of a record to only those entitled
  to it.

## Goals

TRAIL aims to provide the following properties for agent calls that
cross node boundaries:

Integrity:
: Any change to a call, its response, or its record after signing is
  detectable.

Origin:
: Every call and record is attributable to an enforcement point and,
  through it, to a workload identity and a principal.

Completeness:
: Omitted, reordered, or rewritten records are detectable by parties
  other than the operator.

Policy binding:
: Each record names the policy and the delegated authority under which
  the action was allowed.

Causality:
: A response, and any action it triggers, can be tied to the call that
  caused it, across nodes.

Liveness:
: A node that stops publishing evidence, or loses contact with the
  authorities it answers to, stops acting.

Live oversight:
: Authorities, including ones the operator did not choose, can watch,
  gate, hold, and stop activity while it runs.

Bounded scope:
: Limits that span many nodes can be counted and enforced.

Scoped disclosure:
: Each party sees only the parts of a record it is entitled to.

Media evidence:
: Media produced or consumed by agents can be bound to the same
  records.

## Design Principles

One artifact:
: A single signed object, the Action Statement, serves as per-message
  proof of possession for the receiver, as evidence for auditors, and
  (when a response cites it) as the counterparty's receipt.

One transport:
: Agent calls travel over MOQT through the agentproto MOQT binding
  {{PACE}} and protocol mappings such as MCP over MOQT {{MCP-MOQT}}. The
  evidence and control tracks use the same relays, the same
  authorization, and the same naming. A fallback for HTTP calls is in
  {{http-fallback}}.

Three tracks per party:
: Every party (cell, authority, witness, operator) publishes a `log`
  track, a `state` track, and directed `to` tracks. There are no
  per-function namespaces.

Two freshness rules:
: A cell must hold a fresh state snapshot from each authority it
  answers to, and its log must hold a fresh witness cosignature. These
  two rules provide kill-switch reachability, monitoring leases,
  heartbeats, and checkpoint liveness.

## Scope

This document specifies:

* The roles and architecture ({{architecture}}).
* The Action Statement, its claims, and its signature ({{statement}}).
* The core rule: no statement, no action ({{core-rule}}).
* Why MOQT is the primary transport ({{why-moqt}}).
* The MOQT binding: namespaces, the three track kinds, the
  TRAIL_STATEMENT Property, and delivery rules ({{binding}}).
* How TRAIL attaches to PACE and to MCP over MOQT ({{protocols}}).
* Logs, checkpoints, and witnesses ({{log}}).
* State snapshots, directives, and the freshness rules ({{state}}).
* Gates, admission, scoped disclosure, and media evidence
  ({{functions}}).
* Identity, delegation, and attestation profiles ({{identity}}).
* Relay authorization with C4M ({{authz}}).

This document does not specify:

* Agent application protocols. A2A {{A2A}}, MCP {{MCP}}, and others
  keep their message formats.
* Dialog management. TRAIL consumes the agentproto dialog context; it
  does not define it.
* Policy languages, inspector models, or the detection of harmful
  content. A signed statement shows what happened and who did it, not
  whether it was safe.
* New MOQT messages.

## Requirements Language

{::boilerplate bcp14-tagged}

## Terminology

This document uses the terms Track, Track Namespace, Full Track Name,
Group, Object, Subgroup, Location, Property, Immutable Properties,
Original Publisher, Relay, and Fetch from {{MoQ-TRANSPORT}}; Signed
Statement, Transparency Service, and Receipt from {{RFC9943}};
Workload Identity Token (WIT) from {{WIMSE-CREDS}}; and Attestation
Result from {{RFC9334}}.

Cell:
: A node running one or more agent harnesses behind an edge. The cell
  is the unit at which TRAIL enforces.

Edge:
: The component that every call into or out of a cell passes through.
  It holds the keys, evaluates policy, signs Action Statements, and
  applies holds. The agent cannot reach it. An edge can be a sandbox
  egress layer, an agent gateway, or a component on separate hardware
  in the network path. When it is split across components, they act
  as one edge.

Party:
: Any entity that publishes TRAIL tracks: a cell edge, an authority,
  a witness, an operator, or an admission issuer.

Operator:
: The organization that runs cells and agents.

Principal:
: The person or organization an agent acts for.

Authority:
: A party with oversight rights over some scope, such as a regulator,
  a counterparty's inspector, or the operator's own safety team.

Witness:
: A party that verifies a cell's log and countersigns its checkpoints.
  An authority can be its own witness.

Action Statement:
: A signed statement by an edge about one action ({{statement}}).

Directive:
: A signed instruction carried in a party's state snapshot: a policy
  reference, a hold, a requirement, or a grant ({{state}}).

Hold:
: A directive that tells edges to stop actions that match a selector.

Gate:
: An action that waits for a signed approval of its request.

Swarm:
: A set of cells and parties that share one TRAIL deployment and one
  top-level namespace.

# Architecture {#architecture}

~~~
 +---------------- Cell A -----------------+
 |  +---------+       +------------------+  |
 |  | Agent   |------>|      Edge A      |--+--> call Object on PACE track
 |  | harness | local |  policy, keys,   |  |    + TRAIL_STATEMENT
 |  +---------+ MOQT  |  holds, signing  |  |
 |                    +------------------+  |
 +-----------------------|---------^--------+
            log, state   |         |  authority state (policy, holds,
            (Action      v         |  status lists), approvals
             Statements)           |
                   +-----------------------+
                   |      MOQT relays      |
                   +-----------------------+
                     /         |         \
                    v          v          v
               Authority    Witness     Edge B
               (watch,      (verify,    (verify statement,
                hold, gate)  cosign)     deliver, cite)
~~~
{: #fig-arch title="TRAIL architecture"}

The agent never holds a signing key. The edge terminates the agent's
local MOQT session, or a local API, and is the Original Publisher of
the agent's call Objects towards the network. For each consequential
action the edge:

1. checks the action against cached state from its authorities;
2. signs an Action Statement bound to the call Object;
3. appends the statement to its log; and
4. publishes the call Object with the statement attached.

The receiving edge verifies the attached statement without any network
round trip, delivers the call, and appends a statement that cites it.
Authorities and witnesses subscribe to logs by namespace prefix.
Directives flow back on the authorities' state tracks.

## Attestation Levels

The value of a statement depends on where the edge's key lives. TRAIL
names three levels. The level is never asserted by the edge itself: it
is derived from an Attestation Result {{RFC9334}} about the edge's key
({{attestation}}).

| Level | Key held in | Example |
|---|---|---|
| `software` | A process on the host | Sandbox or gateway with no hardware root |
| `tee` | A TPM or trusted execution environment | Confidential VM, enclave-backed gateway |
| `isolated` | Hardware in the network path, separate from the host | DPU or SmartNIC enforcement point |
{: #tab-levels title="Attestation levels"}

At the `software` level a statement still shows integrity and binds the
operator to what it signed, but a compromised host can sign anything.
Authorities can require a higher level for higher-impact actions.

## Dialogs {#dialogs}

Agent interactions are often grouped into dialogs: a shared context,
named by an identifier and a generation, that several agents open,
join, update, and leave {{PACE}}. TRAIL does not manage dialogs. It
records the actions taken inside them, so that those records can be
correlated and the dialog's own changes are witnessed:

* Every Action Statement for an action inside a dialog carries the
  dialog identifier and generation ({{claims}}). Records of all
  participants can be correlated across cells and operators.
* Dialog lifecycle operations (OPEN, ACCEPT, UPDATE, REVOKE,
  TERMINATE) are consequential actions and SHOULD be recorded. Changes
  to dialog membership are then witnessed.
* A hold is applied by the edge regardless of dialog state. When a hold
  stops a participant, the dialog controller SHOULD follow with REVOKE,
  and the edge MUST keep applying the hold until it is lifted.
* Where a dialog requires confirmation or evidence of approval for a
  high-risk operation, a TRAIL gate approval ({{gates}}) is a suitable
  format.
* TRAIL tracks share relays and authorization with the dialog's own
  tracks. Dialog identifiers appear only inside signed payloads, never
  in TRAIL namespaces, so relays do not learn dialog membership from
  track names.

## Why MOQT {#why-moqt}

TRAIL carries agent calls on MOQT whenever it can, and always carries
its own `log`, `state`, and `to` tracks on MOQT. Several of its
simplifications depend on properties that MOQT has and HTTP does not:

Exact binding:
: A MOQT Object has a stable name and address (Full Track Name, Group,
  Object) and does not change once published {{MoQ-TRANSPORT}}. An
  Action Statement can therefore bind to exactly one call
  (`object-ref`). An HTTP request has no such name, and a statement
  can only cover a hash of selected components.

Evidence that cannot be stripped:
: Only the Original Publisher can add Immutable Properties, relays
  forward them unchanged, and Secure Objects {{MoQ-SECOBJ}}
  authenticates them end to end. An HTTP intermediary can drop a
  header field.

One commit point:
: The edge writes the statement to its log and publishes the call on
  the same session, through the same first-hop relay, so "appended
  before the call leaves" ({{commit}}) is a local ordering rule. When
  the call travels on another transport, the ordering spans two
  systems.

Fan-out to parties the publisher does not choose:
: SUBSCRIBE_TRACKS on a namespace prefix lets authorities, witnesses,
  and archives receive every log in a swarm, including logs of cells
  that appear later, without the publisher knowing who they are.

Pushed state:
: Inline checks need holds, status lists, and cosigned checkpoints at
  every edge as soon as they change. A subscription with
  FILL_PARAMETERS delivers the latest snapshot and every later one, and
  the freshness rules ({{freshness}}) depend on that push.

History and live data under one name:
: The same `log` track serves live statements by SUBSCRIBE and past
  ones by FETCH, from relays or archives.

One authorization system and one key:
: C4M with DPoP {{MoQ-C4M}} authorizes call tracks and TRAIL tracks
  with the same WIT-bound key and the same prefix rule.

Media:
: Agents that see and hear receive audio and video as MOQT Objects in
  the same dialog, so media seals and key release ({{media}}) use the
  same tracks and relays as the rest of the evidence.

The agentproto MOQT binding {{PACE}} and MCP over MOQT {{MCP-MOQT}}
already carry A2A and MCP calls as MOQT Objects. Calls
that cannot use MOQT use the HTTP fallback ({{http-fallback}}), which
loses exact binding, protection against stripping, and the single
commit point, but keeps the log, state, and freshness rules.

# The Action Statement {#statement}

An Action Statement is a COSE_Sign1 {{RFC9052}} structure. It is a
SCITT Signed Statement {{RFC9943}}, so it can also be registered
with a Transparency Service without conversion.

## Protected Header

The protected header MUST contain:

* `alg` (1): an asymmetric signature algorithm. ES256 is mandatory to
  implement.
* `typ` (16) {{RFC9596}}: `application/trail-statement+cbor`.
* `kid` (4): the JWK Thumbprint of the edge key, which is the `cnf` key
  of the edge's WIT ({{identity}}).
* CWT Claims (15) {{RFC9597}} with at least:
  * `iss` (1): the edge's WIMSE identifier.
  * `sub` (2): the identifier of the agent that acted.
  * `iat` (6): the time of signing.

## Payload {#claims}

The payload is a CBOR {{RFC8949}} map, described in CDDL {{RFC8610}}:

~~~ cddl
trail-statement = {
  1 => uint,                  ; version (1)
  2 => kind,
  3 => log-ref,               ; issuer's log and leaf index
  ? 4 => object-ref,          ; the MOQT Object this action produced
  5 => bstr,                  ; policy: hash of the directive set in force
  ? 6 => outcome,
  ? 7 => dialog-ref,          ; agentproto dialog context
  ? 8 => tstr,                ; aud: counterparty identifier
  ? 9 => context-hashes,      ; credentials relied on, by hash
  ? 10 => [+ cite],           ; what this action depended on
  ? 11 => any,                ; kind-specific body
  * int => any                ; SD-CWT redacted claims and extensions
}

kind = &(
  send: 1, recv: 2, tool-call: 3, model-call: 4,
  state-write: 5, state-read: 6,
  gate-request: 7, approval: 8,
  hold-applied: 9, disclosure: 10,
  checkpoint: 11, media-seal: 12,
  lifecycle: 13, admission-release: 14, restart: 15
)

log-ref = [ log-id: bstr .size 16, leaf: uint ]

object-ref = [
  ftn-hash: bstr,             ; SHA-256 of the serialized Full Track Name
  group: uint,
  object: uint,
  payload-hash: bstr          ; SHA-256 of the application payload
]

outcome = &( allowed: 0, denied: 1, held: 2, gated: 3,
             completed: 4, failed: 5 )

dialog-ref = [ dialog-id: bstr .size 16, generation: uint ]

context-hashes = {
  ? 1 => bstr,                ; wth: hash of the agent's WIT
  ? 2 => bstr,                ; tth: hash of the Transaction Token
  ? 3 => { + tstr => bstr },  ; oth: other context tokens by name
  ? 4 => bstr,                ; sdh: sd_hash of the capability card
  ? 5 => bstr,                ; adm: hash of the admission token
}

cite = [ cite-type, ref: (log-ref / object-ref), hash: bstr ]
cite-type = &( input: 1, peer-statement: 2, media-group: 3,
               state: 4, directive: 5, approval: 6 )
~~~
{: #fig-statement title="Action Statement payload"}

The `context-hashes` member follows the WIMSE Workload Proof Token
{{WIMSE-WPT}}, which binds a proof to the tokens presented with it by
hash: `wth` for the WIT, `tth` for a Transaction Token, and `oth` for
other context tokens such as an access token. An Action Statement is, in effect, a Workload Proof Token for
one MOQT Object that is also kept as evidence.

The statement identifier is the SHA-256 hash of the COSE_Sign1 with an
empty unprotected header. Receipts and other material MAY be added to
the unprotected header, as in SCITT, without changing the identifier.

## Signature Binding

`log-ref` binds the statement to one position in the issuer's log, and
`object-ref` binds it to one Object at one Location. A relay or archive
therefore cannot serve the statement at another position, and the
statement cannot be attached to another Object. Verifiers MUST check
both against where they received the statement and the Object.

## Selective Disclosure {#sd}

Sensitive claims, such as amounts, principal identifiers, or tool
arguments, SHOULD be redacted with SD-CWT {{SD-CWT}}. The signed
statement then holds salted digests, which every verifier can check.
Each disclosure is labeled with a data class. The edge encrypts the
disclosures of a class to each party entitled to that class with
COSE-HPKE {{COSE-HPKE}}, and sends them on that party's directed track
({{tracks}}). Every authority verifies the same signed statement, and
each authority opens only the claims of its classes.

This construction is the CBOR equivalent of the SD-JWT construction
used by SD-Card {{SD-CARD}}, so capability cards and evidence share one
disclosure model.

# The Core Rule {#core-rule}

Agent-to-agent and agent-to-tool traffic MUST NOT leave a cell unless
an Action Statement for it has been appended to the edge's log. An edge
MUST drop agent traffic that does not carry a valid Action Statement,
both outbound and inbound.

| Destination | Edge behavior |
|---|---|
| Agent or tool reached over MOQT | Check, sign, append to log, publish the call Object with TRAIL_STATEMENT |
| Agent or tool reached over HTTP | As above, with the statement carried as in {{http-fallback}} |
| Model endpoint | Record a `model-call` statement; the payload is not forwarded to the plane |
| Anything else | Deny by default |
{: #tab-core title="Core rule at the edge"}

Agent logic does not change. Routine checks run against state already
cached at the edge ({{state}}) and do not wait on the network. Only
gated actions wait for a round trip.

The receiving edge MUST:

1. Drop the Object if TRAIL_STATEMENT is missing.
2. Verify the statement signature with the key named by `kid`, obtained
   from the sender's state ({{state}}).
3. Check that `object-ref` matches the Object's Full Track Name,
   Location, and payload hash, and that `aud`, if present, names the
   receiver.
4. Check the sender against the receiver's requirements: attestation
   level, upstream freshness ({{freshness}}), holds, and the status of
   the tokens named in `context-hashes`.
5. Reject a statement identifier it has already accepted, unless the
   action is idempotent.
6. Deliver the Object, and append a `recv` statement that cites the
   sender's statement.

A response Object carries the responder's statement, which cites the
request's statement. The requester keeps it as the receipt. Each side
can then prove what the other said.

# MOQT Binding {#binding}

## Requirements on MOQT

Implementations MUST support {{MoQ-TRANSPORT}} draft-22 or later,
including PUBLISH_NAMESPACE, SUBSCRIBE_TRACKS, PUBLISH, SUBSCRIBE with
FILL_PARAMETERS, FETCH, TRACK_STATUS, Immutable Properties, and the
AUTHORIZATION TOKEN parameter with C4M tokens {{MoQ-C4M}}.

`log`, `state`, and `to` tracks MUST use Subgroup Delivery Mode, MUST
NOT use Datagram Delivery Mode, and MUST NOT set a delivery timeout.

Nothing that matters for safety may be carried in a mutable Property.

## Namespaces {#namespaces}

Every TRAIL track lives under the publishing party's prefix:

~~~
("trail", S, P)        -- "log"     append-only statements of party P
("trail", S, P)        -- "state"   current snapshot of party P
("trail", S, P, "to")  -- Q         directed statements from P to Q
~~~
{: #fig-ns title="TRAIL namespace layout"}

S is the swarm pseudonym. P and Q are party pseudonyms. All three MUST
be 16-byte values that do not encode identities, locations, or tasks.
The mapping from a pseudonym to a WIMSE identifier is in that party's
state snapshot, which relays cannot read without authorization
({{privacy}}).

Call Objects stay on their dialog tracks, for example
`("pace", "agents", {agent-id}, {dialog-id}, "a2a", "task") -- {task}`
{{PACE}}. Media stays on its dialog media tracks.

An authority sees the whole swarm with one SUBSCRIBE_TRACKS on
`("trail", S)`, or one party with SUBSCRIBE_TRACKS on `("trail", S, P)`.

## The Three Track Kinds {#tracks}

`log`:
: The party's append-only sequence of statements. Each Object is one
  statement. A Group is one epoch and ends with a checkpoint ({{log}}).
  Object IDs within a Group start at 0 and are consecutive. The
  publisher MUST NOT use the Prior Object ID Gap Property on a log
  track, so any gap is evidence of loss or withholding.

`state`:
: The party's current snapshot. Each Group holds exactly one Object, a
  complete signed snapshot ({{state}}). A receiver never merges deltas.
  The track SHOULD use Descending group order. Subscribers join with
  FILL_PARAMETERS selecting the latest Group, so they get current
  state and live updates in one step.

`to`:
: Statements addressed to one other party: gate approvals, admission
  tokens, encrypted disclosures, and requests for disclosure. Payloads
  that are not public MUST be encrypted to the recipient with
  COSE-HPKE. Every directed statement that matters as evidence is also
  appended, in redacted form, to the sender's log.

## The TRAIL_STATEMENT Property {#property}

TRAIL defines one Object Property, TRAIL_STATEMENT (odd type, byte
value), whose value is the Action Statement for that Object.

* The edge MUST place it inside Immutable Properties. If the call track
  is protected with Secure Objects {{MoQ-SECOBJ}}, the edge SHOULD
  place it in the Encrypted Properties List instead, so that relays
  cannot see it.
* Relays forward it unchanged, as they do any unknown Property. A relay
  cannot add it, because MOQT does not let relays add Immutable
  Properties. This is why the edge must be the Original Publisher.
* Because the statement is carried with the Object, the receiver needs
  no FETCH and no second subscription to verify a call.

## Commit Point {#commit}

MOQT does not acknowledge individual Objects. TRAIL therefore defines
the commit point as follows:

1. The edge MUST write the statement to its `log` track before it
   releases the call Object. Deployments SHOULD make the edge its own
   first-hop relay, so the statement is cached before the call leaves.
2. The receiver verifies the statement inline (it is in the Object).
3. The checkpoint for the epoch MUST include the statement. A receiver
   that holds a statement whose leaf is covered by a cosigned checkpoint
   that does not include it has proof of equivocation.
4. For actions above a policy threshold, the edge MAY wait until a
   witness has cosigned a checkpoint that includes the statement before
   it releases the call (synchronous mode).

## Delivery Parameters

| Track | Publisher priority | Group order | Joining |
|---|---|---|---|
| `state` | 0 | Descending | Latest Group with FILL_PARAMETERS |
| `to` | 0 | Ascending | Subscribe; FETCH for history |
| `log` | 128 | Ascending | FETCH from an archive, then subscribe |
| media | 32 to 63 | Descending | Live; FETCH complete Groups for evidence |
{: #tab-delivery title="Recommended delivery parameters"}

MOQT priority applies only within a session and only under congestion.
An edge SHOULD carry its `state` and `to` subscriptions on a session
separate from bulk log and media traffic.

# Carrying TRAIL in Agent Protocols over MOQT {#protocols}

TRAIL does not change agent protocols. It attaches to the MOQT
Objects that their MOQT mappings already publish. This section states
how for the agentproto MOQT binding {{PACE}} and for MCP over MOQT
{{MCP-MOQT}}. Other protocols carried by PACE bindings, such as A2A,
follow the PACE rules.

## The General Rule

In every mapping, the edge acts for the agent and is the Original
Publisher of the agent's Objects. For each Object that carries a
request, a response, a result, or a control operation, the edge
attaches TRAIL_STATEMENT ({{property}}). Objects that carry only
progress, keep-alive, or streaming fragments MAY omit it, provided the
Object that completes the exchange carries a statement whose
`object-ref` covers it and whose body lists the hashes of the
fragments it completes.

The protocol's own identifiers (PACE `message_id`, A2A task
identifiers, MCP JSON-RPC `id`) SHOULD be recorded in the statement
body, so that the evidence can be matched to the protocol's own logs.
They are not used for any TRAIL check.

## PACE {#pace}

Dialog tracks:
: Call Objects stay on the participant's dialog tracks, for example
  `("pace", "agents", {agent-id}, {dialog-id}, "a2a", "task") --
  {task-id}`. TRAIL tracks are separate, under `("trail", S, P)`. The
  PACE agent-id and dialog-id never appear in TRAIL namespaces; they
  appear only inside signed statements.

Dialog context:
: Every statement for an Object on a dialog track carries `dialog-ref`
  with the PACE dialog-id and the generation the sender last knew.

Control operations:
: Each object on a participant's `ctrl` track (OPEN, ACCEPT, REJECT,
  UPDATE, JOIN_REQUEST, LEAVE, REVOKE, TERMINATE) SHOULD carry a
  `lifecycle` statement. PACE requires that deployments that do not
  trust every relay provide object-level sender authentication for
  `ctrl` objects, and notes that Secure Objects does not provide it.
  TRAIL_STATEMENT provides it, and the log makes changes to dialog
  membership witnessed.

Envelope:
: The PACE envelope wraps the application payload. TRAIL_STATEMENT is
  an Object Property outside the payload. `payload-hash` covers the
  full Object payload, envelope included.

Credentials:
: Identity assertions, Transaction Tokens, and SD-Card presentations
  travel in the PACE `credentials` member, as PACE specifies. The
  statement carries their hashes in `context-hashes`. An SD-Card key
  binding JWT is per-verifier, so it belongs in `credentials` on OPEN,
  ACCEPT, or JOIN_REQUEST, never on the cached `card` track.

Confirmation and evidence:
: Where PACE lets a peer require evidence of approval for a high-risk
  operation, a TRAIL `approval` statement ({{gates}}) is the evidence.
  The receiving edge finds it as a `cite` of type `approval` in the
  sender's statement.

Holds and REVOKE:
: A hold stops matching actions at the edge regardless of dialog
  state. When a hold stops a participant, the controller SHOULD follow
  with REVOKE so that the dialog state matches, and the edge keeps
  applying the hold until it is lifted.

Discovery:
: The PACE `card` track and TRAIL `state` track are separate. A cell's
  `state` MAY include the issuer-signed SD-Card that the agent also
  publishes on its `card` track, so that a receiver can check
  capability claims and key binding from the same subscription it uses
  for the edge's credentials.

## MCP over MOQT {#mcp}

Request and response tracks:
: MCP over MOQT maps each JSON-RPC exchange onto a request track
  published by the requester and a response track published by the
  responder, with one Group per request. The requester's edge attaches
  a statement of kind `tool-call` (or `send` for other requests) to
  the request Object. The responder's edge attaches a statement that
  cites the request's statement to the Object that carries the result
  or error. That statement is the requester's receipt.

Multi-object responses:
: A response Group can hold progress notifications followed by a
  result. The Object that carries the result MUST carry a statement,
  and the general rule above applies to the notifications.

Server-initiated requests:
: Sampling and elicitation requests come from the MCP server. The
  server's edge treats them as requests and the client's edge as
  responses, so model use and user input requested by a tool are
  recorded with the same rules. A sampling response SHOULD also be
  recorded as a `model-call` statement by the client's edge.

Origin of tool results:
: MCP over MOQT requires Secure Objects on tool, sampling, and
  elicitation tracks that cross relays, and notes that shared keys do
  not show which endpoint produced an Object, so the application needs
  a per-publisher signature. TRAIL_STATEMENT is that signature, and the
  statement in the Encrypted Properties List keeps it hidden from
  relays.

Standalone and PACE modes:
: In the standalone mode, on `(mcp, <session-id>, ...)` tracks, there
  is no dialog; statements omit `dialog-ref` and record the MCP session
  identifier in the body. In the PACE mode, on
  `("pace", "agents", {agent-id}, {dialog-id}, "mcp", ...)` tracks,
  statements carry `dialog-ref` as in {{pace}}.

Authorization:
: MCP over MOQT maps the MCP OAuth framework onto MOQT by minting C4M
  tokens from the OAuth grant and binding them with DPoP. TRAIL uses
  the same tokens ({{authz}}). The OAuth access token presented to the
  MCP server, or a Transaction Token derived from it, is recorded by
  hash in `oth` or `tth`.

Resources and skills:
: Resource, prompt, and skill tracks carry content the server
  publishes, not actions. An edge SHOULD record a `state-read`
  statement when an agent uses such content in a later action, and the
  later action cites it, so context provenance follows data from a
  tool into a decision.

# Logs, Checkpoints, and Witnesses {#log}

## The Log

Each party's log is a Merkle tree as defined in {{RFC9162}}. The leaves
are statement identifiers, in log order. The leaf index is the `leaf`
value in the statement's `log-ref`. One log per cell is the default. An
edge MAY run more than one log, each with its own `log-id` and `log`
track, for example one per tenant.

An edge that loses its log state MUST start a new log with a new
`log-id`, MUST record a `restart` statement in the new log that names
the old `log-id` and last known checkpoint, and MUST NOT reuse an old
`log-id`.

## Checkpoints

The last Object of every Group on a `log` track is a `checkpoint`
statement. Its body is a COSE Receipt {{RFC9942}} for the
RFC 9162 verifiable data structure: a signed tree head that states the
tree size and root after the epoch. A checkpoint also lists the media
seals of the epoch ({{media}}), because those are statements in the log.

The epoch length is set by policy. A witness MAY ask for an early
checkpoint with REQUEST_UPDATE and NEW_GROUP_REQUEST, when the track
sets DYNAMIC_GROUPS.

## Witnesses

A witness subscribes to a party's `log` track. For each checkpoint it:

1. checks that it holds every statement up to the checkpoint's tree
   size, and fetches any it lacks;
2. checks that the new tree is consistent with the last tree head it
   cosigned for that log;
3. countersigns the checkpoint with a COSE countersignature
   {{RFC9338}}; and
4. publishes the countersignature on its own `log` track and on its
   `to` track for the party.

A witness MUST NOT countersign two different roots for the same log and
tree size. Two such countersignatures are proof that the witness
misbehaved, and two conflicting checkpoints with countersignatures from
different witnesses are proof that the log operator equivocated.

A cosignature says more than "the root is consistent". It says that the
witness has received every statement in the tree. A cosignature from an
authority is therefore also proof that the authority is watching that
log. TRAIL uses this in place of a separate monitoring lease.

The party copies its latest cosigned checkpoint into its `state`
snapshot, so that any peer can check it with one subscription.

A witness MAY be a SCITT Transparency Service. Deployments that want
long-term public transparency MAY also register checkpoints with one.

## Archives

MOQT relays are not required to cache. Each deployment MUST name at
least one archive for every witness. An archive subscribes to `log`
tracks and serves them by FETCH as a publisher. Witnesses and
authorities that miss Objects fetch them from an archive. A witness
MAY be its own archive.

# State, Directives, and Freshness {#state}

## State Snapshots

A party's `state` track carries its current self-description and its
directives. Each snapshot is a COSE_Sign1 signed by the party's key,
whose payload is:

~~~ cddl
trail-state = {
  1 => uint,                    ; version (1)
  2 => uint,                    ; sequence, increases with every snapshot
  3 => uint,                    ; iat
  4 => uint,                    ; exp
  5 => tstr,                    ; WIMSE identifier of the party
  ? 6 => [+ bstr],              ; credentials: WIT, Attestation Result,
                                ;   capability cards (cells)
  ? 7 => bstr,                  ; latest cosigned checkpoint (cells)
  ? 8 => [+ directive],         ; directives (authorities)
  ? 9 => [+ bstr],              ; Token Status Lists (issuers)
  ? 10 => [+ grant],            ; grants to authorities (operators)
  * int => any
}

directive = {
  1 => directive-type,
  2 => bstr,                    ; directive id
  3 => selector,                ; what it applies to
  ? 4 => any,                   ; type-specific value
  ? 5 => uint,                  ; not-before
  ? 6 => uint,                  ; not-after
}

directive-type = &(
  policy: 1,          ; value: hash and location of a policy bundle
  hold: 2,            ; stop matching actions
  require-level: 3,   ; minimum attestation level for matching actions
  require-witness: 4, ; value: [witness id, max age]
  require-gate: 5,    ; matching actions need an approval from this party
  require-admission: 6  ; matching agents need an admission token
)

selector = {
  ? 1 => [+ bstr],              ; party pseudonyms
  ? 2 => [+ tstr],              ; agent identifiers or prefixes
  ? 3 => [+ kind],              ; statement kinds
  ? 4 => [+ uint],              ; action classes (deployment-defined)
  ? 5 => [+ tstr],              ; counterparty identifiers or prefixes
}

grant = {
  1 => tstr,                    ; authority's WIMSE identifier
  2 => [+ directive-type],      ; directive types it may issue
  3 => selector,                ; the scope it may issue them for
  4 => [+ uint],                ; data classes it may read
}
~~~
{: #fig-state title="State snapshot"}

Policy bundles, which can be large, are referenced from directives by
hash. The bundle itself is a statement in the authority's log, and
edges fetch it once. The `policy` claim of every Action Statement is
the hash of the set of directives the edge applied, so "it was allowed"
can be checked later.

An authority's directives apply only within the grant its operator
gives it. An edge MUST ignore any directive that falls outside the
issuing authority's grant in the operator's current state snapshot.
One authority therefore cannot stop a whole swarm unless the operator
agreed to that.

Holds that apply to a single credential, such as one agent's WIT or
one admission token, SHOULD be expressed as a `suspended` entry in the
issuer's Token Status List {{OAUTH-STATUS-LIST}}, carried in the
issuer's state. Revocation is an `invalid` entry. This uses the same
mechanism that OAuth and WIMSE deployments use for credential status,
delivered by push instead of by polling.

An edge MUST record a `hold-applied` statement the first time it stops
an action because of a hold.

## The Two Freshness Rules {#freshness}

Downstream freshness (authority to cell):
: An edge MUST hold, for every authority that has a grant in its
  operator's state, a snapshot whose `exp` has not passed. If it does
  not, it MUST pause the agents covered by that authority's grant.

Upstream freshness (cell to witness):
: For every `require-witness` directive that applies, a cell's latest
  checkpoint MUST carry a countersignature from the named witness that
  is no older than the stated maximum age. If it does not, the cell's
  edge MUST pause the covered agents, and receivers MUST reject calls
  from that cell that are covered by the directive.

An authority keeps its snapshot fresh by republishing it before `exp`,
even when nothing changes. This is the heartbeat, and it proves that
the authority can still reach the edge. A witness keeps a cell fresh by
cosigning its checkpoints. This proves that the witness is still
receiving the cell's log.

These two rules replace heartbeats, monitoring leases, checkpoint
deadlines, and maximum staleness bounds. For example, the rule "any
agent calling service Z must be monitored by regulator R" is a
`require-witness` directive from R, with Z as the counterparty in its
selector. If R stops watching, its cosignatures stop, the callers'
checkpoints age out, and Z's edge rejects their calls.

Pausing MUST also be triggered when a subscription to a required
`state` track ends with PUBLISH_DONE, a stream reset, or REQUEST_ERROR
and cannot be restored before the snapshot's `exp`.

# Gates, Admission, Disclosure, and Media {#functions}

## Gates {#gates}

For an action matched by a `require-gate` directive, the edge holds
the action and:

1. appends a `gate-request` statement to its log, whose body is the
   hash of the pending call Object and any redacted context the
   approver needs; and
2. sends the same statement on its `to` track for the approver, with
   disclosures for the approver's classes.

The approver publishes an `approval` statement on its `to` track for
the cell, and appends it to its own log. The approval cites the request
by statement identifier and states the decision. The edge continues
only after a valid approval, and its `send` statement cites it. When
several approvals are required, the edge waits for each of them.

The approver can be a person. In that case the approval is signed with
a key bound to that person, for example through an OAuth client or a
platform authenticator, and the approval statement names that person
in a redacted claim.

## Admission and Counting {#admission}

Limits that span many cells, such as "no more than 1,000 live agents in
swarm S", use admission tokens:

* The admission issuer is a party. It issues each agent an admission
  token: a CWT {{RFC8392}} with `cnf` set to the agent's WIT key, a
  unique `cti`, an expiry, and a status claim that points into the
  issuer's Token Status List {{OAUTH-STATUS-LIST}}. The token is sent on
  the issuer's `to` track for the cell.
* The issuer never has more than the cap of valid tokens at once.
* A `require-admission` directive tells receivers to accept calls only
  from agents whose statements carry the hash of a valid admission
  token in `context-hashes`.
* When an agent retires, its edge appends an `admission-release`
  statement. The issuer then marks the token `invalid` and can issue a
  new one.
* Any witness can count live tokens from the issuer's log, and the same
  token cited by two agents shows up in the logs.

Runtime limits use the token's expiry, and the C4M `moqt-reval` claim
({{authz}}) for relay access.

## Disclosure on Request

An authority that needs a claim it was not sent proactively sends a
`disclosure` request on its `to` track for the cell. The edge answers
on its `to` track for the authority, with disclosures encrypted to the
authority, and appends a `disclosure` statement to its log. Every look
at protected data is therefore itself a witnessed statement.

## Media Evidence {#media}

Agents that see and hear act on media that is already carried on MOQT
tracks. For each media Group the edge publishes, or receives and
forwards to the agent, it computes a Merkle root over the (Object ID,
payload hash) pairs of the Group, and appends a `media-seal` statement
that names the track, the Group, the Object count, and the root.

An Action Statement that depended on media cites it with a
`media-group` cite (track, Group, root). Changing one frame breaks the
link. Inclusion proofs let a verifier check single frames.

Live media MAY use Datagram Delivery Mode and delivery timeouts.
Evidence checks MUST use complete Groups fetched from an archive.

For approval before playback, each media Group uses its own Secure
Objects key. The publisher sends the key, encrypted with COSE-HPKE, to
the approver on its `to` track. When the approver approves, it sends
the key to the entitled viewers on its own `to` tracks. Nothing plays
before approval, and each release is a statement in the approver's log.

# Identity, Delegation, and Attestation {#identity}

## Identity: WIMSE

WIMSE workload credentials {{WIMSE-CREDS}} are mandatory to implement.

* Each edge, authority, witness, operator, and admission issuer has a
  WIMSE identifier and a WIT. The `kid` of every TRAIL signature is the
  thumbprint of the signer's WIT `cnf` key.
* Each agent has its own WIMSE identifier and WIT. The agent's WIT key
  is held by the edge, never by the harness, consistent with
  {{WIMSE-AIMS}}. The statement's `sub` is the agent's identifier and
  `wth` is the hash of the agent's WIT.
* WITs are published in the party's `state` snapshot. A receiver
  subscribes to the sender's `state` track once and then verifies every
  call from that sender without further lookups.

The WIMSE architecture {{WIMSE-ARCH}} requires an intermediary that
changes identity context to make the change explicit and auditable. The
edge acts for the agent, and the Action Statement is that explicit
statement.

Other identity schemes, such as X.509 certificates or SPIFFE SVIDs, can
be used through a profile that provides the same four things: a key
bound to an operator and an agent, a link to the principal, an
attestation result for the key, and a status mechanism.

## Principal and Delegation: OAuth

TRAIL does not define a delegation format. The principal and the chain
of delegation are carried by OAuth Transaction Tokens {{TXN-TOKENS}},
whose `sub`, `scope`, `tctx`, and `rctx` claims describe who the call
is for, what it may do, and the context of the originating request, or by an access token obtained through Token Exchange
{{RFC8693}} with its `act` chain. Fine-grained permissions can use
Rich Authorization Requests {{RFC9396}}.

The token itself travels in the application message, for example in a
PACE dialog's `credentials` or in MCP `_meta`. The Action Statement
carries its hash in `tth` or `oth`, so the evidence shows which grant
the action relied on without repeating it. The Transaction Token's
`txn` value SHOULD be recorded in the statement, so that evidence and
the operator's own transaction logs can be correlated.

## Capability: SD-Card

An agent's capability card, for example an SD-Card {{SD-CARD}}, says
what the agent claims it can do. TRAIL uses it as follows:

* The card's `cnf` key SHOULD be the agent's WIT key. One key then
  covers discovery, workload authentication, relay authorization, and
  evidence.
* The issuer-signed card, with only its public disclosures and no key
  binding JWT, MAY be included in the cell's `state` snapshot. It is
  the same for every subscriber, so relays can cache it.
* Context-specific disclosures and the key binding JWT are
  per-verifier. They MUST travel inside the dialog (for example in PACE
  OPEN or ACCEPT), not on a cached track. The key binding JWT's
  `interaction_id` SHOULD be the dialog identifier.
* An Action Statement records the `sd_hash` of the presentation the
  counterparty relied on in `sdh`. "B acted because A presented
  capability X, issued by registry R" can then be checked later.

## Attestation: RATS {#attestation}

The attestation level ({{tab-levels}}) of an edge key MUST come from an
Attestation Result {{RFC9334}} produced by a verifier the receiver
trusts, bound to the edge's WIT key, for example as in
{{WIMSE-ATTEST}}. An EAT Attestation Result {{RATS-EAR}} {{RFC9711}}
wrapped in a Conceptual Message Wrapper {{RFC9999}} is RECOMMENDED.
Trustworthiness claims {{RATS-AR4SI}} map to levels by deployment
policy. The Attestation Result is included in the edge's `state`
snapshot and refreshed before it expires. An edge MUST NOT
state its own level.

# Authorization {#authz}

Every party authenticates to relays with a C4M token {{MoQ-C4M}} bound
with DPoP {{RFC9449}}. The token's `cnf.jkt` MUST be the thumbprint of
the party's WIT key, so relays authorize by workload identity and one
key serves every purpose.

One rule covers publishing: a party's token grants PUBLISH_NAMESPACE
and PUBLISH only under its own prefix `("trail", S, P)`. Subscription
rights depend on role:

| Role | SUBSCRIBE / SUBSCRIBE_TRACKS / FETCH under |
|---|---|
| Cell edge | `state` of its authorities, operator, and peers; its own `("trail", S, X, "to") -- P` tracks |
| Authority, witness | `("trail", S)` prefix, or the parties in its grant; its own `to` tracks |
| Archive | `log` and `state` tracks it stores |
{: #tab-roles title="C4M scopes per role"}

In C4M terms, a cell edge with pseudonym P gets a scope with actions
PUBLISH_NAMESPACE and PUBLISH and namespace match `["trail", S, P]`,
and a second scope with action SUBSCRIBE and track match on its own
pseudonym under the `to` namespaces of its authorities.

MOQT does not copy authorization tokens from SUBSCRIBE_TRACKS into the
PUBLISH messages it produces. A relay MUST check each forwarded PUBLISH
against the scope of the token presented with the SUBSCRIBE_TRACKS.

Relays are not trusted for TRAIL's guarantees, but they can help
enforce them:

* C4M tokens for agents SHOULD carry `moqt-reval` and a status claim
  that refers to the issuer's Token Status List. A relay that
  revalidates then stops forwarding for a suspended agent, which adds a
  network-level kill switch to the edge's own.
* An edge relay can check that call Objects carry TRAIL_STATEMENT. It
  cannot verify the statement. A C4M claim that requires a given
  Property on every published Object would let relays enforce presence
  ({{open-issues}}).

C4M and DPoP authenticate sessions and MOQT actions hop by hop. The
Action Statement authenticates the content of each action end to end.
They are complementary and use the same key.

# Security Considerations

This section follows {{RFC3552}}.

## Sources of the Guarantees

| Guarantee | Comes from | Not from |
|---|---|---|
| Origin | Statement signature by an edge key with an Attestation Result | TLS, which covers one hop; Secure Objects, whose keys are shared |
| Integrity and binding | Signature over `log-ref` and `object-ref` | Relay behavior |
| Completeness | Merkle log, dense Object IDs, witness cosignatures | Relay delivery or caching |
| Non-equivocation | Witnesses that refuse to cosign conflicting roots | The operator |
| Liveness of oversight | The two freshness rules | Relay heartbeats |
| Confidentiality of claims | SD-CWT redaction and COSE-HPKE to each recipient | Namespace access control |
| Access at the edge | C4M with DPoP | Hop-by-hop enforcement beyond the first relay |
{: #tab-trust title="Sources of TRAIL guarantees"}

## Parties and Threats

Agent:
: Can ask to act. Cannot sign, reach a relay, or skip policy.

Operator:
: Runs cells and chooses relays. Cannot edit statements after the fact
  or show different histories to different witnesses without producing
  proof of equivocation. Cannot quietly stop publishing, because
  checkpoints stop being cosigned and the freshness rules pause the
  cell. Can choose which authorities to grant; that choice is itself
  in the operator's signed state.

Relay:
: Can route, cache, delay, or drop Objects, and can change mutable
  Properties. Cannot forge statements or move them, and dropped
  statements show up as gaps.

Authority:
: Can read its granted classes, issue directives within its grant, and
  cosign. Cannot read other classes, because disclosures are encrypted
  to each recipient. Cannot act outside its grant.

Hardware root:
: Vouches for the keys of `tee` and `isolated` edges. It is the one
  root that every party must trust. Deployments SHOULD accept roots
  from more than one vendor.

## Shared Keys

Secure Objects keys are shared by a track's publisher and subscribers,
so any key holder can produce Objects that decrypt. TRAIL never relies
on them for origin. Disclosures use public-key encryption to each
recipient, so authorities cannot forge disclosures for each other.

## Replay

A replayed call Object carries a statement identifier the receiver has
already accepted. Receivers MUST keep the identifiers they accept for
at least the maximum statement age they allow, and reject repeats for
actions that are not idempotent. Even when a replay is not caught, the
`object-ref` binding means it can only reproduce the same Object at the
same Location.

## Denial of Service

The freshness rules let an attacker who can block state or
cosignatures pause agents. Deployments SHOULD provide redundant relay
paths for `state` and `to` tracks, SHOULD use more than one witness,
and SHOULD scope pausing by action impact. Prefix subscriptions near
the root of the namespace can cause large numbers of PUBLISH messages,
and relays SHOULD limit them through authorization.

## Limits

A signed statement shows what happened and who sent it, not whether it
was harmful. A signed, approved message can still carry an injection.
Covert channels through timing, patterns in ordinary output, or public
content are not blocked and can only be found by audit. Interfaces
through which a hardware enforcement point exposes its signing key are
not standardized.

# Privacy Considerations {#privacy}

A system that can prove what every agent did can also show what every
person using those agents did. TRAIL follows {{RFC6973}} and makes
broad surveillance technically hard:

* Namespaces contain only rotating pseudonyms. WIMSE identifiers,
  dialog identifiers, and principals appear only in signed payloads,
  and sensitive claims are redacted.
* Authority rules run in the cell as directives. By default only
  verdicts, counts, and Merkle roots leave the cell.
* Disclosures go only to recipients entitled to a class, encrypted to
  each recipient, and every disclosure is itself a logged statement.
* Disclosure keys and archived disclosures SHOULD be tied to a purpose
  and a retention period, and destroyed when it ends.
* `state` snapshots map pseudonyms to identities. Relays MUST restrict
  subscription to them by C4M scope, and deployments SHOULD encrypt
  snapshot claims that are not needed by every subscriber.
* Publishers SHOULD pad statements to limit size analysis.

Timing and who coordinates with whom remain visible to relays.

# IANA Considerations

## MOQ Properties

IANA is requested to register the following entry in the "MOQ
Properties" registry. Until a value is assigned, implementations can
use a value from the application-specific range 0x3800 to 0x3FFF.

| Type | Name | Scope | Reference |
|---|---|---|---|
| TBD1 (odd) | TRAIL_STATEMENT | Object | This document |

## Media Types

IANA is requested to register `application/trail-statement+cbor` and
`application/trail-state+cbor`.

## TRAIL Registries

IANA is requested to create the registries "TRAIL Statement Kinds",
"TRAIL Directive Types", and "TRAIL Cite Types", with the initial
values in {{fig-statement}} and {{fig-state}} and a Specification
Required policy.

--- back

# HTTP Fallback {#http-fallback}

When an agent call travels over HTTP, the edge signs the request with
WIMSE HTTP Message Signatures {{WIMSE-HTTPSIG}} {{RFC9421}} on the
agent's behalf, and adds a `Trail-Statement` header field whose value
is the base64url-encoded Action Statement. The header MUST be a covered
component. `object-ref` is replaced by a hash of the covered
components. The edge MUST request a signed response, and the
responder's edge returns its own statement in a `Trail-Statement`
response header, which serves as the receipt. Both statements are
appended to the respective logs on MOQT as usual.

# Related Work {#related}

Several individual drafts address per-action evidence for agents or
workloads over HTTP:

* {{WIMSE-AUTHZ-EVIDENCE}} defines signed "Permits" committing to a
  request before dispatch, anchored in SCITT. A TRAIL `gate-request`
  plus `approval` pair plays the same role.
* {{WIMSE-AGENT-EVIDENCE}} points from a WIT to an evidence service
  and defines hash-chained execution receipts.
* {{SCITT-LOCAL-LOG}} defines a producer-operated log with signed
  checkpoints inside SCITT statements, which is close to the TRAIL log.
* {{AGENTPROTO-VERIFIER}} discusses verifier evidence in agentproto.
* {{AAUTH}} signs every agent request with HTTP Message Signatures and
  keeps a mission log.

TRAIL differs in that the evidence is produced outside the agent at an
attested edge, is delivered live over MOQT to parties the publisher
does not choose, and carries control (holds, gates, admission) back on
the same tracks. TRAIL aims to converge with these drafts on the
statement envelope, which all of them base on COSE or SCITT.

# Examples

## A Gated Payment Across Operators

1. Operator O grants financial regulator F the directive types
   `require-gate` and `hold` for action class "payment", and data class
   "finance". The grant is in O's state.
2. F's state contains `require-gate` for class "payment" above a
   threshold, and `require-level` of `isolated` for the same class.
3. Agent A in cell C_A, in dialog D, asks to send a payment instruction
   to agent B over A2A on its PACE track.
4. Edge A appends a `gate-request` statement, with the amount redacted,
   and sends it with the "finance" disclosure encrypted to F on
   `("trail", S, C_A, "to") -- F`.
5. F appends an `approval` statement and sends it on
   `("trail", S, F, "to") -- C_A`.
6. Edge A appends a `send` statement that cites the approval and
   publishes the call Object with TRAIL_STATEMENT.
7. Edge B checks the statement, A's level, and A's cosigned checkpoint,
   delivers the call, and appends a `recv` statement citing A's.
8. A health authority watching the same swarm sees that a payment
   happened, but cannot open the amount.
9. At the end of the epoch, F cosigns both cells' checkpoints.

## A Monitored Service

Regulator R publishes `require-witness` with itself as witness, a
maximum age of 60 seconds, and service Z as counterparty. Z's edge
rejects any call whose sender's latest checkpoint lacks an R
cosignature younger than 60 seconds. R cosigns only logs it is
receiving in full. If R stops watching, calls to Z stop within a
minute.

## A Swarm Size Cap

Admission issuer Q never has more than 1,000 valid admission tokens for
swarm S. Q's `require-admission` directive makes peers reject calls
from agents without a valid token hash. An agent spawned without a
token has nobody to talk to. Retiring an agent releases its token, Q
marks it `invalid`, and any witness can count valid tokens in Q's log.

# Open Issues {#open-issues}

* A C4M claim that requires a given Property on every Object published
  under a scope, so that edge relays can enforce the presence of
  TRAIL_STATEMENT.
* A Token Status List delivery profile over MOQT, for use beyond TRAIL.
* Whether TRAIL_STATEMENT should be standardized as a general "per
  Object proof" Property for MOQT, with TRAIL as one user, which would
  amount to the first non-HTTP binding of the WIMSE Workload Proof
  Token.
* Coordination with WIMSE and SCITT individual drafts on per-action
  evidence, to avoid incompatible statement formats.
* Measurements of state propagation latency, and of log throughput at
  1,000 to 10,000 agents per swarm.

# Acknowledgments
{:numbered="false"}

This document draws on discussions of PACE, MCP over MOQT, C4M, and
SD-Card.
