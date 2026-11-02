---
title: "TRAIL: Tamper-evident Records for Agent Interactions and Liveness"
abbrev: "TRAIL"
category: std
docname: draft-nassar-agentproto-trail-latest
submissiontype: IETF
number:
date: 2026-10-31
consensus: true
v: 3
workgroup: "Agent Communication Protocols"
keyword:
 - agent communication
 - audit
 - revocation
 - transparency
 - MOQT
 - HTTP
venue:
  group: "Agent Communication Protocols"
  type: "Working Group"
  mail: "agentproto@ietf.org"
  arch: "https://mailarchive.ietf.org/arch/browse/agentproto/"
  subscribe: "https://www.ietf.org/mailman/listinfo/agentproto/"

author:
 -
    fullname: Kareem Nassar
    organization: Cisco
    email: kanassar@cisco.com
 -
    fullname: Suhas Nandakumar
    organization: Cisco
    email: snandaku@cisco.com
 -
    fullname: Cullen Jennings
    organization: Cisco
    email: fluffy@cisco.com

normative:
  RFC2119:
  RFC8174:
  RFC8610:
  RFC8949:
  RFC9052:
  RFC9110:
  RFC9162:
  RFC9338:
  RFC9421:
  RFC9449:
  RFC9530:
  RFC9596:
  RFC9597:
  RFC9651:
  RFC9943:
  MoQ-TRANSPORT: I-D.ietf-moq-transport
  MoQ-C4M: I-D.ietf-moq-c4m
  MoQ-SECOBJ: I-D.ietf-moq-secure-objects
  OAUTH-STATUS-LIST: I-D.ietf-oauth-status-list
  SD-CWT: I-D.ietf-spice-sd-cwt
  COSE-HPKE: I-D.ietf-cose-hpke
  WIMSE-CREDS: I-D.ietf-wimse-workload-creds
  WIMSE-HTTPSIG: I-D.ietf-wimse-http-signature
  TXN-TOKENS: I-D.ietf-oauth-transaction-tokens
  RFC9334:
  RFC8392:
  RATS-EAR: I-D.ietf-rats-ear
  TLOG-TILES:
    title: Static asset-based transparency log
    author:
    - org: C2SP
    date: false
    target: https://c2sp.org/tlog-tiles
  TLOG-CHECKPOINT:
    title: Interoperable transparency log signed tree heads
    author:
    - org: C2SP
    date: false
    target: https://c2sp.org/tlog-checkpoint
informative:
  RFC3552:
  RFC6973:
  RFC8693:
  MCP-MOQT: I-D.jennings-agentproto-mcp-over-moqt
  PACE: I-D.nandakumar-agentproto-pace
  RFC7942:
  RFC9396:
  RFC9711:
  RFC9999:
  RATS-AR4SI: I-D.ietf-rats-ar4si
  WIMSE-ARCH: I-D.ietf-wimse-arch
  WIMSE-WPT: I-D.ietf-wimse-wpt
  WIMSE-AIMS: I-D.ietf-wimse-aims
  WIMSE-ATTEST: I-D.reddy-wimse-workload-attestation
  WIMSE-AUTHZ-EVIDENCE: I-D.munoz-wimse-authorization-evidence
  WIMSE-AGENT-EVIDENCE: I-D.thaha-wimse-agent-evidence-binding
  SCITT-LOCAL-LOG: I-D.mih-scitt-checkpointed-local-log
  AAUTH: I-D.hardt-oauth-aauth-protocol
  AGENTPROTO-VERIFIER: I-D.he-agentproto-verifier-evidence
  TLOG-WITNESS:
    title: HTTP protocol to obtain transparency log witness cosignatures
    author:
    - org: C2SP
    date: false
    target: https://c2sp.org/tlog-witness
  SSE:
    title: 'HTML Living Standard: Server-sent events'
    author:
    - org: WHATWG
    date: false
    target: https://html.spec.whatwg.org/multipage/server-sent-events.html
  A2A:
    title: Agent2Agent Protocol Specification
    author:
    - org: A2A Project
    date: false
    target: https://a2a-protocol.org/latest/specification/
  MCP:
    title: Model Context Protocol Specification
    author:
    - org: Anthropic
    date: false
    target: https://modelcontextprotocol.io/specification/
  RFC9942:
  AGENTPROTO-126:
    title: Minutes of the Agent Communication Protocols (agentproto) BoF, IETF 126
    author:
    - org: IETF
    date: 2026-07
    target: https://datatracker.ietf.org/doc/minutes-126-agentproto-202607230700/
  AUDIT-ARCH: I-D.kuehlewind-audit-architecture
  SD-CARD: I-D.nandakumar-agent-sd-jwt
  RFC9901:

--- abstract

TRAIL makes agent-to-agent and agent-to-tool communication auditable
and revocable. An enforcement point outside the agent, the edge, signs
each call as a COSE Action Record. That one record authenticates the
call, lets the receiver prove what it was sent, and is appended to a
tamper-evident trail that witnesses countersign. The edge lets an
agent's traffic through only while it holds current state from every
authority it answers to and a recent witness countersignature on its
trail. TRAIL defines bindings to Media over QUIC Transport (MOQT) and to
HTTP. Both produce the same evidence. The MOQT binding also binds each
record to exactly one Object, pushes holds to edges as soon as they are
published, and delivers trails to authorized parties the operator did
not enumerate.

--- middle

# Introduction

An agent is a model running in a loop with tools, credentials, and a
goal. It reads untrusted input, acts with real authority, and works
faster than people can watch. Two things are needed outside the agent:

1. An enforcement point that the agent cannot see, reach, or turn off.

2. A record of what the agent did that the operator cannot quietly
   rewrite, that parties outside the operator can check, and that those
   parties can act on while the agent is still running.

Single-node sandboxes and gateways provide the first. The second is
missing as soon as one agent's call reaches an agent on another node.

## Problem Statement {#problem}

Agent calls cross operator boundaries, but the controls and the records
for those calls do not. Today:

* A receiver cannot verify, per call, which workload sent it, on whose
  behalf, under which policy, and whether it was altered in transit by
  an intermediary such as a relay or proxy.

* Records of agent actions are written by the same software and
  operator whose conduct they describe, so they can be dropped,
  reordered, or rewritten without detection.

* Parties with oversight rights cannot observe activity while it
  happens, pause or stop it, or require approval before an action
  takes effect. Revoking an agent is hardest when the agent is no
  longer reachable.

* Limits that apply across many nodes, such as a cap on concurrent
  agents in a swarm, cannot be enforced because no node sees the total.

* Sharing records with one party exposes them to every party, because
  there is no way to disclose part of a record to only those entitled
  to it.

## Goals {#goals}

TRAIL aims to provide the following properties for agent calls that
cross node boundaries:

Integrity:
: Any change to a call, its response, or its record after signing is
  detectable.

Origin:
: Every call and record is attributable to an edge and, through it, to
  a workload identity and a principal.

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
: An edge stops an agent's traffic when it loses contact with the
  authorities it answers to, or when its trail is no longer being
  witnessed.

Live oversight:
: Authorities, including ones the operator did not enumerate, can
  watch, gate, hold, and stop activity while it runs.

Bounded scope:
: Limits that span many nodes can be counted and enforced.

Scoped disclosure:
: Each party sees only the parts of a record it is entitled to.

Media evidence:
: Media produced or consumed by agents can be bound to the same
  records.

## Design Principles {#principles}

One artifact:
: A single signed object, the Action Record, serves as per-message
  proof of possession for the receiver, as evidence for auditors, and,
  when a response cites it, as the counterparty's acknowledgment.

The evidence comes from outside the agent:
: The edge holds the keys and signs. How far a record can be trusted
  follows from where the edge sits, which a verifier attests
  ({{levels}}), not from anything the agent asserts.

Two bindings, one evidence format:
: Agent calls travel over MOQT, through the agentproto MOQT binding
  {{PACE}} and MCP over MOQT {{MCP-MOQT}}, or over HTTP. Records,
  trails, state, and the rules that act on them are the same in both.
  MOQT is preferred for the evidence plane ({{why-moqt}}).

Three track kinds per party:
: Every party publishes records, state, and directed records
  ({{tracks}}).

Two freshness rules:
: Fresh authority state and fresh witness cosignatures keep an edge
  running ({{freshness}}).

## Scope {#scope}

This document specifies the Action Record, the rule that no action
leaves a cell without one, the MOQT and HTTP bindings, how TRAIL
attaches to agent protocols, trails and witnesses, state and
directives, and the profiles of existing identity, delegation,
attestation, and authorization work that TRAIL uses.

This document does not specify:

* Agent application protocols. A2A {{A2A}}, MCP {{MCP}}, and others
  keep their message formats.

* Dialog management. TRAIL consumes the agentproto dialog context; it
  does not define it.

* Agent behavior, decision-making, policy languages, inspector models,
  or the detection of harmful content ({{limits}}).

* New MOQT messages.

## Requirements Language

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and
"OPTIONAL" in this document are to be interpreted as described in
BCP 14 {{RFC2119}} {{RFC8174}} when, and only when, they appear in all
capitals, as shown here.

## Terminology {#terminology}

This document uses the terms Track, Track Namespace, Full Track Name,
Group, Object, Subgroup, Location, Property, Immutable Properties,
Original Publisher, Relay, and Fetch from {{MoQ-TRANSPORT}}; Signed
Statement and Transparency Service from {{RFC9943}}; Workload Identity
Token (WIT) from {{WIMSE-CREDS}}; and Attestation Result and Verifier
from {{RFC9334}}.

Agent:
: Software that uses an AI model in a loop, with tools and credentials,
  to pursue a goal for a principal.

Sandbox:
: The isolated environment an agent runs in, such as a container or a
  microVM. TRAIL does not define sandboxes. It requires that every path
  out of one pass through an edge.

Edge:
: Where traffic leaves a cell. Every call into or out of the cell
  passes through it. It holds the keys, evaluates policy, signs Action
  Records, appends them to its trail, and applies holds. The agent
  cannot reach it. When an edge is split across components, they act
  as one edge.

Cell:
: Everything behind one edge. Depending on where the edge is
  ({{levels}}), a cell is one sandbox, a confidential VM, or a whole
  machine. The cell is the unit at which TRAIL enforces.

Party:
: Any entity that publishes TRAIL tracks: an edge, an authority, a
  witness, an operator, or an admission issuer. One party can hold
  several roles.

Operator:
: The organization that runs cells and agents, and grants authorities
  their rights.

Principal:
: The person or organization an agent acts for.

Authority:
: A party with rights over some scope that an operator has granted,
  such as a regulator, a counterparty's inspector, or the operator's
  own safety team.

Witness:
: A party that checks that it has received every record in a trail and
  countersigns the trail's checkpoints. It has no rights over actions
  ({{roles}}).

Action Record:
: A signed record by an edge of one action ({{record}}). It is a SCITT
  Signed Statement.

Trail:
: The append-only sequence of Action Records of one party, organized as
  a Merkle tree ({{trails}}).

Checkpoint:
: An Action Record that signs the size and root of its own trail at the
  end of an epoch, covering every record before it.

Epoch:
: The span of a trail that ends with a checkpoint. Over MOQT, an epoch
  is one Group of the records track.

Archive:
: A service that stores trails and serves them to authorized readers
  ({{archives}}).

Evidence plane:
: The tracks, or HTTP resources, that carry trails, state, and directed
  records between parties.

Acknowledgment:
: The responder's Action Record that cites a request's record. It is
  not a COSE Receipt {{RFC9942}}.

Directed record:
: A record published by one party for one other party, which can be
  encrypted to it: gate requests, approvals, admission tokens,
  disclosures, requests for disclosure, and cosignatures ({{tracks}}).

State:
: A party's current signed snapshot, republished on a schedule
  ({{snapshots}}).

Directive:
: A rule in an authority's state snapshot: a policy reference, a hold,
  or a requirement ({{snapshots}}).

Grant:
: A right an operator gives an authority: which directive types it may
  issue, for what scope, and which data classes it may read.

Action class:
: A deployment-defined number that groups actions by impact, such as
  payments. Selectors and freshness budgets select actions by class.

Data class:
: A deployment-defined number that labels a redacted claim, such as an
  amount. Grants list the data classes an authority may read.

Hold:
: A directive that tells edges to stop actions that match a selector.

Gate:
: An action that waits for a signed approval of its request.

Admission issuer, admission token:
: A party that hands out a fixed number of tokens, and the token an
  agent's records cite ({{admission}}).

Attestation level:
: `software`, `tee`, or `isolated`: who can read the edge's key, and
  what can go around the edge ({{levels}}).

Swarm:
: A set of cells and parties that share one TRAIL deployment and one
  top-level namespace.

# Architecture {#architecture}

~~~ ascii-art
 +---------------- Cell A -----------------+
 |  +---------+       +------------------+  |
 |  | Agent   |------>|      Edge A      |--+--> call (MOQT or HTTP)
 |  | (in a   | local |  policy, keys,   |  |    + Action Record
 |  | sandbox)|       |  holds, signing  |  |
 |  +---------+       +------------------+  |
 +-----------------------|---------^--------+
            trail, state |         |  authority state (policy, holds,
            (Action      v         |  status lists), approvals
             Records)              |
                   +-----------------------+
                   |    evidence plane     |
                   | (MOQT relays or HTTP) |
                   +-----------------------+
                     /         |         \
                    v          v          v
               Authority    Witness     Edge B
               (watch,      (verify,    (verify record,
                hold, gate)  cosign)     deliver, cite)
~~~
{: #fig-arch title="TRAIL architecture"}

The agent never holds a signing key. For each consequential action the
edge:

1. checks the action against cached state from its authorities;

2. signs an Action Record bound to the outgoing call;

3. appends the record to its trail; and

4. sends the call with the record attached.

The receiving edge verifies the attached record, delivers the call, and
appends a record that cites it. Once the receiver holds the sender's
state, verification needs no network round trip; on first contact it
needs one subscription or fetch of that state. Authorities and
witnesses read trails. Directives flow back in the authorities' state.

## Attestation Levels {#levels}

The agent always runs in a sandbox. What differs is where the edge sits,
and therefore what the cell contains, who can read the edge's key, and
what can go around the edge. TRAIL names three levels. The level is
never asserted by the edge itself: it is derived from an Attestation
Result {{RFC9334}} about the edge's key ({{attestation}}).

| Level | Where the edge is | What the cell is | Who can read the key | What goes around the edge |
|---|---|---|---|---|
| `software` | A process beside the sandbox on the same host, such as a sidecar proxy or egress gateway | The sandbox and the edge | The host | An escape from the sandbox |
| `tee` | Inside a confidential VM, with its key sealed to that VM | The confidential VM | Not the host | An escape into the VM |
| `isolated` | Separate hardware that all of the machine's traffic crosses, such as a DPU or SmartNIC | The whole machine | Nothing on the host | Only compromise of that hardware |
{: #tab-levels title="Attestation levels"}

A key held in a TPM on an ordinary host cannot be exported, but the
host can still use it to sign. An edge with such a key is at the
`software` level.

The edge is well placed to produce these records. It is on the path, so
the trail is complete; it holds keys outside the agent, so a
compromised agent cannot forge its own records; it appends each record
before the call leaves; and a verifier can attest what holds its key.
At the software level a record still shows integrity and binds the
operator to what it signed, but a compromised host can sign anything.
Authorities can require a higher level for higher-impact actions with
a require-level directive.

## Authorities and Witnesses {#roles}

TRAIL separates two jobs:

* A witness checks the completeness and consistency of trails. It needs
  every record, but only its structure: identifiers, order, and
  redacted claims. It has no rights over actions. Its cosignatures
  affect actions only through an authority's require-witness directive
  ({{freshness}}).

* An authority sets rules, issues holds, approves gated actions, and
  reads the claims its grant covers.

Keeping them apart gives least privilege, lets witnesses check
authorities as well as operators, and lets several witnesses back one
authority. One party MAY hold both roles; a monitor that both checks a
trail and can stop the agent it records is an authority and a witness
at once.

## Dialogs {#dialogs}

Agent interactions are grouped into dialogs: the correlated sequence of
interactions among participants over the course of a task, named by a
dialog identifier and a generation and managed by a controller
{{PACE}}. TRAIL does not manage dialogs. It records the actions taken
inside them:

* Every Action Record for an action inside a dialog carries the dialog
  identifier and the generation the sender last knew ({{payload}}).
  Records of all participants can then be correlated across cells and
  operators.

* Dialog control operations that change the dialog or stop work are
  consequential actions and SHOULD be recorded. Changes to dialog
  membership are then witnessed.

* A hold is applied by the edge regardless of dialog state. It does
  not depend on the dialog's controller.

* Where a dialog requires evidence of approval for a high-risk
  operation, a TRAIL approval record ({{gates}}) is that evidence.

* TRAIL namespaces never contain dialog identifiers. A binding that
  puts them in its own track names, as PACE does, protects them with
  its own relay authorization.

{{pace}} maps these rules onto PACE.

## Why MOQT {#why-moqt}

Every TRAIL property is available over HTTP. MOQT is preferred for the
evidence plane, and for calls where both sides support it. A MOQT
Object has a stable name and cannot be changed once published
{{MoQ-TRANSPORT}}, relays cannot strip Immutable Properties, the trail
and the call share one session, and subscriptions push trails and
state to parties the publisher did not enumerate. {{tab-http-losses}}
compares each of these with HTTP. MOQT also provides:

History and live data under one name:
: The same records track serves live records by SUBSCRIBE and past ones
  by FETCH, from relays or archives.

One authorization system and one key:
: C4M with DPoP {{MoQ-C4M}} authorizes call tracks and TRAIL tracks
  with the same WIT-bound key and the same prefix rule.

Media:
: Agents that see and hear receive audio and video as MOQT Objects in
  the same dialog, so media seals and key release ({{media}}) use the
  same tracks and relays as the rest of the evidence.

# The Action Record {#record}

An Action Record is a COSE_Sign1 {{RFC9052}} structure. It is a SCITT
Signed Statement {{RFC9943}}, so it can also be registered with a
Transparency Service without conversion. It is the same in both
bindings except for the call-ref member.

## Protected Header

The protected header MUST contain:

* alg (1): an asymmetric signature algorithm. ES256 is mandatory to
  implement.

* typ (16) {{RFC9596}}: `application/trail-record+cbor`.

* kid (4): the JWK Thumbprint of the edge key, which is the cnf key of
  the edge's WIT ({{identity}}).

* CWT Claims (15) {{RFC9597}} with at least:

  * iss (1): the edge's WIMSE identifier.

  * sub (2): the identifier of the agent that acted.

  * iat (6): the time of signing.

## Payload {#payload}

The payload is a CBOR {{RFC8949}} map, described in CDDL {{RFC8610}}:

~~~ cddl
trail-record = {
  1 => uint,                  ; version (1)
  2 => kind,
  3 => trail-ref,             ; issuer's trail and leaf index
  ? 4 => call-ref,            ; the call this action produced
  5 => bstr,                  ; policy: hash of directives in force
  ? 6 => outcome,
  ? 7 => dialog-ref,          ; agentproto dialog context
  ? 8 => tstr,                ; aud: counterparty identifier
  ? 9 => context-hashes,      ; credentials relied on, by hash
  ? 10 => [+ cite],           ; what this action depended on
  ? 11 => any,                ; kind-specific body
  ? 12 => tstr,               ; principal (usually redacted)
  ? 13 => uint,               ; action class
  ? 14 => [+ state-ref],      ; authority snapshots the edge held
  * int => any                ; SD-CWT redacted claims and extensions
}

kind = &(
  send: 1, recv: 2, tool-call: 3, model-call: 4,
  state-write: 5, state-read: 6,
  gate-request: 7, approval: 8,
  hold-applied: 9, disclosure: 10,
  checkpoint: 11, media-seal: 12,
  lifecycle: 13, admission-release: 14, restart: 15,
  cosignature: 16
)

trail-ref = [ trail-id: bstr .size 16, leaf: uint ]

call-ref = moqt-ref / http-ref
moqt-ref = [
  1,
  ftn-hash: bstr,             ; SHA-256 of the Full Track Name
  group: uint,
  object: uint,
  payload-hash: bstr          ; SHA-256 of the application payload
]
http-ref = [
  2,
  components-hash: bstr,      ; see the HTTP binding
  created: uint               ; the signature's created parameter
]

outcome = &( allowed: 0, denied: 1, held: 2, gated: 3,
             completed: 4, failed: 5, degraded: 6, unsigned: 7 )

dialog-ref = [ dialog-id: bstr .size 16, generation: uint ]

context-hashes = {
  ? 1 => bstr,                ; wth: hash of the agent's WIT
  ? 2 => bstr,                ; tth: hash of the Transaction Token
  ? 3 => { + tstr => bstr },  ; oth: other context tokens by name
  ? 4 => bstr,                ; sdh: sd_hash of the capability card
  ? 5 => bstr,                ; adm: hash of the admission token
}

state-ref = [ authority: tstr, sequence: uint ]

cite = [ cite-type, ref: (trail-ref / call-ref / media-ref),
         hash: bstr ]
cite-type = &( input: 1, peer-record: 2, media-group: 3,
               directive: 4, approval: 5 )
media-ref = [ ftn-hash: bstr, group: uint, objects: uint ]
~~~
{: #fig-record title="Action Record payload"}

The context-hashes member follows the WIMSE Workload Proof Token
{{WIMSE-WPT}}, which binds a proof to the tokens presented with it by
hash: wth for the WIT, tth for a Transaction Token, and oth for other
context tokens such as an access token. An Action Record is, in effect,
a Workload Proof Token for one call that is also kept as evidence.

A cite's hash is the record identifier of the cited record, except in
a media-group cite, where it is the media root ({{media}}). A directive
cite names the policy bundle record in the authority's trail. The
authority snapshots that the edge held are listed in member 14, not in
cites.

Records and state snapshots MUST use deterministic encoding
({{Section 4.2.1 of RFC8949}}). The record identifier is the SHA-256
hash of the COSE_Sign1 encoded with an empty unprotected header map.
Inclusion proofs and other material MAY be added to the unprotected
header, as in SCITT, without changing the identifier.

## Signature Binding {#binding}

trail-ref binds the record to one position in the issuer's trail, and
call-ref binds it to one call. A moqt-ref names one immutable Object at
one Location; an http-ref covers a hash of the signed components of one
request or response ({{http-calls}}). A relay, proxy, or archive
therefore cannot serve the record at another position, and the record
cannot be attached to another call. Verifiers MUST check both against
where they received the record and the call.

## Selective Disclosure {#sd}

Sensitive claims, such as amounts, principal identifiers, or tool
arguments, SHOULD be redacted with SD-CWT {{SD-CWT}}. The signed record
then holds salted digests, which every verifier can check. Each
disclosure is labeled with a data class. The edge encrypts the
disclosures of a class to each authority whose grant includes that
class with COSE-HPKE {{COSE-HPKE}}, and sends them as directed records
({{tracks}}). Every authority verifies the same signed record, and each
opens only the claims of its classes. Every party that can read the
trail, including witnesses, sees only the redacted record.

This construction is the CBOR equivalent of the SD-JWT construction
{{RFC9901}} used by SD-Card {{SD-CARD}}, so capability cards and
evidence share one disclosure model.

# The Core Rule {#core-rule}

Agent-to-agent and agent-to-tool traffic MUST NOT leave a cell unless
an Action Record for it has been appended to the edge's trail. An edge
MUST drop inbound agent traffic that does not carry a valid Action
Record.

| Destination | Edge behavior |
|---|---|
| Agent or tool reached over MOQT | Check, sign, append to trail, publish the call Object with TRAIL_RECORD |
| Agent or tool reached over HTTP | As above, with the record carried as in {{http-calls}} |
| Model endpoint | Record a model-call record; the payload is not forwarded to the evidence plane |
| Anything else | Deny by default |
{: #tab-core title="Core rule at the edge"}

Agent logic does not change. Routine checks run against state already
cached at the edge ({{state}}) and do not wait on the network. Only
gated actions wait for a round trip.

During incremental deployment, an operator's policy MAY admit unsigned
inbound traffic from named endpoints that do not run TRAIL. The edge
then records a recv record with outcome `unsigned`. A require-record
directive from an authority removes that exception for the senders,
receivers, and action classes it selects, so removing a record in
transit makes the call fail rather than bypass oversight.

The receiving edge MUST:

1. Drop the call if its Action Record is missing.

2. Verify the record signature with the key named by kid, obtained from
   the sender's state ({{state}}).

3. Check that call-ref matches the call as received, and that aud, if
   present, names the receiver.

4. Check the sender against the receiver's requirements: attestation
   level, holds, admission, the status of the tokens named in
   context-hashes, and upstream freshness from the cosigned checkpoint
   in the sender's state ({{freshness}}). The receiver cannot observe
   the sender's downstream freshness directly. The sender declares the
   authority snapshots it held in member 14; a receiver that holds
   those authorities' state MUST check that each named snapshot exists
   and is within the applicable budget.

5. Reject a replayed record ({{replay}}).

6. Deliver the call, and append a recv record that cites the sender's
   record.

A response carries the responder's record, which cites the request's
record. The requester keeps it as the responder's acknowledgment. Each
side can then prove what the other said, which is the basis for
assigning liability between them.

# MOQT Binding {#moqt}

## Requirements on MOQT

Implementations MUST support {{MoQ-TRANSPORT}} draft-22 or later,
including PUBLISH_NAMESPACE, SUBSCRIBE_TRACKS, PUBLISH, SUBSCRIBE with
FILL_PARAMETERS, FETCH, TRACK_STATUS, Immutable Properties, and the
AUTHORIZATION TOKEN parameter with C4M tokens {{MoQ-C4M}}.

Records, state, and directed tracks MUST use Subgroup Delivery Mode,
MUST NOT use Datagram Delivery Mode, and MUST NOT set a delivery
timeout.

An edge MUST NOT carry anything that a TRAIL check depends on in a
mutable Property.

## Namespaces {#namespaces}

Every TRAIL track lives under the publishing party's prefix:

~~~ ascii-art
("trail", S, P)        -- "records" append-only records of party P
("trail", S, P)        -- "state"   current snapshot of party P
("trail", S, P, "to")  -- Q         directed records from P to Q
~~~
{: #fig-ns title="TRAIL namespace layout"}

S is the swarm pseudonym. P and Q are party pseudonyms. All three MUST
be 16-byte values that do not encode identities, locations, or tasks.
The mapping from a pseudonym to a WIMSE identifier is in that party's
state snapshot, which only authorized parties can read
({{authorization}}).

Call Objects stay on their dialog tracks, for example
("pace", "dialogs", {dialog-id}, {agent-id}, "a2a", "task") --
{task-id} {{PACE}}. Media stays on its dialog media tracks.

An authority sees the whole swarm with one SUBSCRIBE_TRACKS on
("trail", S), or one party with SUBSCRIBE_TRACKS on ("trail", S, P).
Because a cell whose records track is not published under ("trail", S)
gets no cosignatures and fails upstream freshness, an operator cannot
start a cell that such a subscription will not see.

## The Three Track Kinds {#tracks}

records:
: The party's trail. Each Object is one Action Record. A Group is one
  epoch and ends with a checkpoint ({{checkpoints}}). Object IDs within
  a Group start at 0 and are consecutive. The publisher MUST NOT use
  the Prior Object ID Gap Property on a records track, so any gap is
  evidence of loss or withholding. A witness MAY ask for an early
  checkpoint with REQUEST_UPDATE and NEW_GROUP_REQUEST, when the track
  sets DYNAMIC_GROUPS.

state:
: The party's current snapshot. Each Group holds exactly one Object, a
  complete signed snapshot ({{snapshots}}). A receiver never merges
  deltas. The track SHOULD use Descending group order. Subscribers join
  with FILL_PARAMETERS selecting the latest Group, so they get current
  state and live updates in one step.

to:
: Directed records addressed to one other party: gate requests and
  approvals, admission tokens, disclosures, and requests for
  disclosure. Payloads that are not public MUST be encrypted to the
  recipient with COSE-HPKE. Every directed record that matters as
  evidence is also appended, in redacted form, to the sender's trail.

## The TRAIL_RECORD Property {#property}

TRAIL defines one Object Property, TRAIL_RECORD (odd type, byte value),
whose value is the Action Record for that Object.

* The edge MUST place it inside Immutable Properties, or in the
  Encrypted Properties List when the call track is protected with
  Secure Objects {{MoQ-SECOBJ}}, so that relays cannot see it.

* Relays forward it unchanged, as they do any unknown Property. A relay
  cannot add it, because MOQT does not let relays add Immutable
  Properties. This is why the edge must be the Original Publisher.

* Because the record is carried with the Object, the receiver needs no
  FETCH and no second subscription to verify a call.

## Commit Point {#moqt-commit}

MOQT does not acknowledge individual Objects. TRAIL therefore defines
the commit point as follows:

* The edge MUST write the record to its records track before it
  releases the call Object. Deployments SHOULD make the edge its own
  first-hop relay, so the record is cached before the call leaves.

* The receiver verifies the record inline (it is in the Object).

* The checkpoint for the epoch MUST include the record. A receiver that
  holds a record whose leaf is covered by a cosigned checkpoint that
  does not include it has proof of equivocation.

* For actions above a policy threshold, the edge MAY wait until a
  witness has cosigned a checkpoint that includes the record before it
  releases the call (synchronous mode).

## Delivery Parameters

| Track | Publisher priority | Group order | Joining |
|---|---|---|---|
| state | 0 | Descending | Latest Group with FILL_PARAMETERS |
| to | 0 | Ascending | Subscribe; FETCH for history |
| records | 128 | Ascending | FETCH from an archive, then subscribe |
| media | 32 to 63 | Descending | Live; FETCH complete Groups for evidence |
{: #tab-delivery title="Recommended delivery parameters"}

MOQT priority applies only within a session and only under congestion.
An edge SHOULD carry its state and to subscriptions on a session
separate from bulk records and media traffic.

# HTTP Binding {#http}

Most agent traffic today is HTTP {{RFC9110}}: MCP {{MCP}} commonly
uses its Streamable HTTP transport, and A2A {{A2A}} uses JSON-RPC over
HTTP. The HTTP binding lets that traffic carry TRAIL without changing
the agents.

## Calls {#http-calls}

The edge signs the request with WIMSE HTTP Message Signatures
{{WIMSE-HTTPSIG}} {{RFC9421}} on the agent's behalf and adds a
Trail-Record header field: an Item Structured Field {{RFC9651}} whose
value is a Byte Sequence holding the Action Record.

~~~ http-message
POST /a2a/tasks HTTP/1.1
Host: agent-b.example
Content-Type: application/json
Content-Digest: sha-256=:X48E9qOokqqrvdts8nOJRJN3OWDUoyWxBf...=:
Trail-Record: :0oRYQ6MBJgN4HmFwcGxpY2F0aW9uL3RyYWlsLXJlY29y...:
Signature-Input: trail=("@method" "@target-uri" "content-digest"
    "trail-record");created=1793318400;tag="trail";
    keyid="NzbLsXh8uDCcd-6MNwXF4W_7noWXFZAfHkxZsRGC9Xs"
Signature: trail=:MEUCIQDk3...:
~~~
{: #fig-http-call title="An HTTP request carrying a record (truncated)"}

The keyid is the JWK Thumbprint of the edge key, the same value as the
record's kid.

* The signature MUST cover `@method`, `@target-uri`, `content-digest`
  {{RFC9530}} when the message has content, and `trail-record`.

* `components-hash` in the record's `http-ref` is the SHA-256 hash of
  the signature base for the covered components other than
  `trail-record`, without the `@signature-params` line. `created`
  equals the signature's `created` parameter.

* The edge MUST request a signed response. The responder's edge returns
  its own record, citing the request's record, in a Trail-Record
  response field covered by its response signature.

* When the response is a stream, such as server-sent events {{SSE}} on
  an MCP Streamable HTTP response, the responder's record is carried in
  the event that completes the exchange and lists the hashes of the
  earlier events.

## Trails as Static Tiles {#http-trails}

A party publishes its trail as a tiled transparency log
{{TLOG-TILES}}. Tiles use the same {{RFC9162}} Merkle tree as TRAIL
trails; entry bundles hold the records in leaf order, and the last
record of each epoch is the checkpoint ({{checkpoints}}). The trail's
checkpoint resource is a signed tree head in the {{TLOG-CHECKPOINT}}
format over the same tree, so existing tlog tooling can read TRAIL
trails directly.

Tiles are immutable once full and can be served from caches. Readers
authenticate as described in {{read-access}} to a server that the
operator whose trail it is does not run, so the operator does not learn
who reads its trail.

A reader finds trails through witnesses. A cell obtains cosignatures by
submitting each checkpoint to its witnesses, for example with
{{TLOG-WITNESS}}, and each witness records the countersignature, with
the cell's trail identifier and base URI, in its own trail. An
authority that reads a witness's trail therefore learns every cell that
witness cosigns, including cells started later.

## State and Directed Records {#http-state}

A party serves its current state snapshot at a `state` resource listed
in its state, with an ETag. Readers poll it with conditional requests
or keep an SSE stream open on which each new snapshot is an event.
Directed records are delivered by POST to the recipient's `inbox`
resource.

The freshness rules are unchanged, because the edge checks the age of
what it holds, not how it arrived. An authority's freshness budget for a
class MUST exceed the edge's poll interval plus the expected delivery
latency, or the edge will pause between polls. A hold therefore takes
effect within one poll interval, instead of on publication.

## Commit Point {#http-commit}

The trail append and the call can be on two systems. The edge MUST
receive a successful acknowledgment of the append before it sends the
call.

## Properties Lost and Their Mitigations {#http-losses}

| Property | MOQT | HTTP | Mitigation |
|---|---|---|---|
| Exact binding | One immutable Object | Hash of the covered components | Sufficient for HTTP calls; WIMSE and AAuth use the same construction |
| Stripping | Relays cannot remove Immutable Properties | An intermediary can drop the header field | require-record ({{core-rule}}): the receiver refuses the call, so stripping fails closed |
| Single commit point | Trail and call on one session | Two systems | Acknowledged append before send |
| Delivery of holds | Pushed on subscription | Poll or SSE | Per-class freshness budgets bound the delay ({{budgets}}) |
| Trails for parties not enumerated | SUBSCRIBE_TRACKS on a prefix | Static tiles, pulled | Discovery through witness trails ({{http-trails}}) |
{: #tab-http-losses title="HTTP binding losses and mitigations"}

## Hybrid Deployment {#hybrid}

Calls and the evidence plane can use different bindings. In a hybrid
deployment, agent calls stay on HTTP and trails, state, and directed
records travel over MOQT, so only edges, authorities, and witnesses
need MOQT. This is the RECOMMENDED first deployment for agents that use
HTTP today. A verifier checks a record the same way whichever binding
carried it.

# Carrying TRAIL in Agent Protocols {#protocols}

TRAIL does not change agent protocols. It attaches to the messages that
their bindings already carry. This section states how for the
agentproto MOQT binding {{PACE}}, for MCP over MOQT {{MCP-MOQT}}, and
for A2A and MCP over HTTP.

## The General Rule

In every mapping, the edge acts for the agent and is the sender of the
agent's messages. For each message that carries a request, a response,
a result, or a control operation, the edge attaches an Action Record:
TRAIL_RECORD over MOQT ({{property}}), or Trail-Record over HTTP
({{http-calls}}). Messages that carry only progress, keep-alive, or
streaming fragments MAY omit it, provided the message that completes
the exchange carries a record whose call-ref covers it and whose body
lists the hashes of the fragments it completes.

The protocol's own identifiers (PACE message_id, A2A task identifiers,
MCP JSON-RPC id, and MCP over MOQT client ID) SHOULD be recorded in
the record body, so that the evidence can be matched to the
protocol's own logs. They are not used for any TRAIL check.

## PACE {#pace}

### Layering {#pace-layers}

TRAIL is not part of the PACE protocol. PACE messages, tracks, and
dialog semantics are unchanged, and PACE works without it. TRAIL runs
at the edge, in the Object path below the PACE MOQT binding: the edge
is the Original Publisher of each Object the binding produces
({{property}}). On the way out, the edge checks the Object against
cached state, signs an Action Record whose call-ref names the Object's
track, Location, and payload, appends the record to its trail, and
publishes the Object with TRAIL_RECORD. On the way in, it verifies the
record before the binding sees the Object. The edge also publishes its
own tracks under ("trail", S, P) on the same relays.

~~~ ascii-art
+------------------------------------+
| Application protocols (A2A, MCP)   |
+------------------------------------+
| PACE dialog management             |
| dialog-id, generation, lifecycle   |
+------------------------------------+
| PACE MOQT binding                  |
| ("pace","dialogs",D,agent) tracks  |
+====================================+
| TRAIL at the edge                  |
| out: check, sign, append, attach   |
|      TRAIL_RECORD                  |--+ "records", "state", "to"
| in:  verify, check, append recv    |  |  under ("trail",S,P)
+====================================+  |
| MOQT: relays, namespaces, C4M      |<-+
+------------------------------------+
| QUIC / WebTransport                |
+------------------------------------+
~~~
{: #fig-pace-layers title="TRAIL in the PACE Object path at the edge"}

PACE correlates a dialog and manages its lifecycle. TRAIL proves what
each participant did in it and enforces rules at the edge. PACE
defines no audit data structures; TRAIL supplies them, keyed by the
PACE dialog identifier.

| Concern | PACE | TRAIL |
|---|---|---|
| Correlation | dialog-id and generation | dialog-ref in each signed record |
| Sender authentication | credentials in control messages; mechanism chosen by the deployment | TRAIL_RECORD signs each request, response, result, and control Object that changes the dialog or stops work, end to end |
| Records | None | Action Records in witnessed trails, on "records" tracks |
| Stopping a participant | REVOKE, sent by the controller | Hold, issued by an authority and applied at the edge |
| Stopping work | CANCEL and CANCEL_ACK | Records both; require-gate stops work before it starts |
| Approval evidence | Required for high-risk operations; format open | Approval record |
| Lost contact | Suspend actions that depend on an unconfirmed context change | Freshness rules pause the edge |
| Relay authorization | C4M scopes per dialog | C4M scopes per party prefix, with the same key |
{: #tab-pace-layers title="Division of work between PACE and TRAIL"}

### Mapping

Dialog tracks:
: Call Objects stay on the participant's dialog tracks, for example
  ("pace", "dialogs", {dialog-id}, {agent-id}, "a2a", "task") --
  {task-id}. TRAIL tracks are separate, under ("trail", S, P). The PACE
  dialog-id and agent-id never appear in TRAIL namespaces.

Dialog context:
: Every record for an Object on a dialog track carries dialog-ref with
  the PACE dialog-id and the generation the sender last knew. These
  are the values of the dialog reference that the application message
  carries under `org.ietf.agentproto/dialog-context`. A receiving edge
  MUST refuse a call whose dialog reference and dialog-ref differ.

Control operations:
: Each control message on a participant's ctrl track that changes the
  dialog or stops work (OPEN, ACCEPT, REJECT, UPDATE, JOIN_REQUEST,
  LEAVE, REVOKE, TERMINATE, CANCEL, and CANCEL_ACK) SHOULD carry a
  lifecycle record. STATE_REQUEST, STATE, and ERROR MAY omit it. PACE
  requires a participant to authenticate the sender of each operation
  before applying it, and relays terminate TLS on each hop. TRAIL_RECORD
  authenticates each control Object end to end, and the trail makes
  changes to dialog membership witnessed.

Payload:
: PACE carries each message as the Object payload, with no PACE
  header. payload-hash covers the whole Object payload.

Credentials:
: Identity assertions, authorization proofs, and SD-Card presentations
  travel in the PACE credentials member, as PACE specifies. The record
  carries their hashes in context-hashes. A key-bound presentation
  travels on ACCEPT, UPDATE, or JOIN_REQUEST, the messages that answer
  a peer's nonce.

Confirmation and evidence:
: Where PACE lets a peer require evidence of approval for a high-risk
  operation, a TRAIL approval record ({{gates}}) is the evidence. The
  receiving edge finds it as a cite of type approval in the sender's
  record.

Holds, REVOKE, and CANCEL:
: When a hold stops a participant, the controller SHOULD follow with
  REVOKE so that the dialog context matches, and the edge MUST keep
  applying the hold until it is lifted. An edge that holds work already
  in progress SHOULD send CANCEL for it, so that peers stop work they
  started for the held agent. A CANCEL_ACK with outcome `irreversible`
  is recorded like any other control message, so the trail shows which
  actions took effect before the stop.

Delegated dialogs:
: A dialog opened to carry out delegated work lists its parents in
  parent_dialog_ids. The lifecycle record of OPEN carries the full
  dialog context, so the trails show which dialogs descend from which.

Discovery:
: The PACE card track and the TRAIL state track are separate. A cell's
  state MAY include the issuer-signed SD-Card that the agent also
  publishes on its card track, so that a receiver can check capability
  claims and key binding from the same subscription it uses for the
  edge's credentials.

Revocation at the relay:
: PACE withdraws a revoked participant's relay authorization, for
  example by refusing C4M revalidation. A hold at the edge stops the
  same participant independently of the relay, so either one is enough
  to stop it.

## MCP over MOQT {#mcp-moqt}

Request and response tracks:
: MCP over MOQT maps each JSON-RPC exchange onto a request track
  published by the client and a response track published by the
  server, in the client's scope (mcp, {server-id}, clients,
  {client-id}, ...). Each request is a new Group, and the response uses
  the same Group ID. The client's edge attaches a record of kind
  tool-call (or send for other requests) to the request Object. The
  server's edge attaches a record that cites the request's record to
  the Object that carries the final result or error.

Track ownership:
: In a client scope the client and the server publish under the same
  namespaces, and only the track name tells a request track from a
  response track. A receiving edge MUST check that the record's iss is
  the edge it expects for that track: the client's edge for a request
  track, the server's edge for a response track. A peer that publishes
  a track belonging to the other is then refused.

Multi-object responses:
: A response Group can hold progress notifications, in one or more
  Subgroups, followed by the final result or error, which has the
  highest Object ID. The Object that carries the final result MUST
  carry a record, and the general rule above applies to the
  notifications, taken in Object ID order. The record travels with its
  Object, so a client that recovers a response by FETCH after a
  reconnection verifies it the same way.

Multi round-trip requests:
: MCP servers do not send requests to clients. A server that needs an
  elicitation, a sampling result, or the client's roots ends the
  response with an input-required result, which carries a record like
  any other result. The client's retry is a new request in a new Group,
  and its record cites the input-required result's record. Sampling
  done by the client to answer it is recorded as a model-call record by
  the client's edge, and the retry cites that record as input. An
  elicitation that a person answers can be gated ({{gates}}).

Cancellation:
: A client cancels a request by publishing notifications/cancelled in
  the request's Group. That Object is a control message and carries a
  record, as does a PACE CANCEL that names the request.

Origin of tool results:
: Secure Objects keys are shared by a track's publisher and
  subscribers, so they do not show which endpoint produced an Object.
  TRAIL_RECORD is a per-publisher signature, and placing it in the
  Encrypted Properties List keeps it hidden from relays.

Standalone and PACE modes:
: In standalone mode there is no dialog. Records omit dialog-ref and
  record the server ID, the client ID, and the JSON-RPC id in the body;
  MCP keeps no session, so there is no session identifier to record.
  In PACE mode the same tracks move under ("pace", "dialogs",
  {dialog-id}, {agent-id}, "mcp", ...), each MCP request carries the
  PACE dialog reference in `_meta`, and records carry dialog-ref as in
  {{pace}}.

Authorization:
: MCP over MOQT presents an OAuth access token, or a C4M token minted
  from the OAuth grant and bound with the DPoP binding of C4M, to the
  relay. Credentials that the MCP server itself verifies travel in each
  request's `_meta`. TRAIL uses the same C4M tokens ({{authorization}}).
  The access token, or a Transaction Token derived from it, is recorded
  by hash in oth or tth.

Catalog, resources, prompts, and skills:
: The catalog, resource, prompt, and skill tracks carry content that
  the server publishes, not actions, and relays serve them from cache.
  Tool definitions and client-executed skill instructions steer the
  model, so an edge SHOULD record a state-read record when an agent
  uses such content, and the later action cites it. Context provenance
  then follows data from a server into a decision. A server-executed
  skill (skills/invoke) and a task are requests like tools/call, and
  are recorded the same way.

## A2A and MCP over HTTP

A2A requests and MCP Streamable HTTP requests carry a Trail-Record
field as in {{http-calls}}. The responder's record is carried in the
response, or, for a streamed response, in the event that completes the
exchange. When the message carries a PACE dialog reference, in MCP
`_meta` or in A2A message metadata, the record carries the same values
in dialog-ref. Records are appended to trails on whichever evidence
plane the deployment uses ({{hybrid}}).

# Trails, Checkpoints, and Witnesses {#trails}

## The Trail

Each party's trail is a Merkle tree as defined in {{RFC9162}}. The
leaves are record identifiers, in trail order. The leaf index is the
leaf value in the record's trail-ref. One trail per cell is the
default. An edge MAY run more than one trail, each with its own
trail-id and records track, for example one per tenant.

An edge that loses its trail state MUST start a new trail with a new
trail-id, MUST record a restart record in the new trail that names the
old trail-id and last known checkpoint, and MUST NOT reuse an old
trail-id. A party changes its pseudonym the same way: it starts a new
trail under the new pseudonym, with a restart record that names the
old trail-id, so that witnesses keep the history connected.

## Checkpoints {#checkpoints}

The last record of every epoch is a checkpoint: an Action Record of
kind checkpoint, signed by the trail's own edge, whose body is a tree
head:

~~~ cddl
tree-head = [
  trail-id: bstr .size 16,
  tree-size: uint,
  root-hash: bstr,            ; RFC 9162 Merkle tree hash
  ? degraded: [+ uint],       ; leaf indexes of degraded actions
  ? media-seals: [+ uint],    ; leaf indexes of media-seal records
]
~~~
{: #fig-tree-head title="Checkpoint body"}

tree-size and root-hash cover every record before the checkpoint. The
checkpoint itself is the leaf at index tree-size, and the next
checkpoint covers it.

A checkpoint is a record made by the trail producer about its own
trail. It is not a COSE Receipt: TRAIL trails are not Transparency
Services, and a party proves inclusion of a record with an {{RFC9162}}
inclusion proof against a cosigned checkpoint. A witness that also
operates a SCITT Transparency Service MAY additionally register records
and issue COSE Receipts {{RFC9942}} for them; TRAIL does not depend on
that.

The epoch length is set by policy.

## Witnesses {#witnesses}

A witness reads a party's trail. For each checkpoint it:

* checks that it holds every record up to the checkpoint's tree size,
  and fetches any it lacks;

* checks that the new tree is consistent with the last tree head it
  cosigned for that trail;

* countersigns the checkpoint with a COSE countersignature {{RFC9338}}
  whose protected header carries CWT Claims with iat, the time of
  cosigning; and

* appends a record of kind cosignature to its own trail, whose body is
  the countersignature with the trail-id and tree size it covers, and
  sends the countersignature to the party as a directed record.

A witness MUST NOT countersign two different roots for the same trail
and tree size. Two such countersignatures are proof that the witness
misbehaved, and two conflicting checkpoints with countersignatures from
different witnesses are proof that the trail's operator equivocated.

A cosignature says more than "the root is consistent". It says that the
witness has received every record in the tree. TRAIL uses this in place
of a separate monitoring lease.

The party copies its latest cosigned checkpoint into its state
snapshot, so that any peer can check it with one subscription or
fetch.

## Archives {#archives}

MOQT relays are not required to cache. Each deployment MUST name at
least one archive for every witness. An archive reads trails and serves
them by FETCH, or as tiles over HTTP. Witnesses and authorities that
miss records fetch them from an archive. A witness MAY be its own
archive.

## What Makes Records Tamper-Evident {#tamper}

Five properties together make changes to records detectable:

1. Each Action Record is signed by the edge key, so changing any byte
   invalidates the signature.

2. Each record is a leaf of its trail's Merkle tree, so editing,
   removing, inserting, or reordering records changes the root.

3. Each checkpoint signs the size and root of the trail, fixing its
   history up to that point.

4. Witnesses cosign checkpoints and never cosign two roots for the same
   size, so a party cannot show different histories to different
   readers without producing proof of equivocation.

5. A counterparty's acknowledgment cites the hash of the record it
   received and is kept in the counterparty's own trail, so a rewritten
   record no longer matches what the counterparty holds.

Tamper evidence is not tamper proofing. An edge that is compromised at
signing time produces validly signed false records. The attestation
level bounds how likely that is, and acknowledgments and witnesses
bound how far a false record can go unnoticed.

# State, Directives, and Freshness {#state}

## State Snapshots {#snapshots}

A party's state carries its current self-description, its directives
(authorities), and its grants (operators). Each snapshot is a
COSE_Sign1 signed by the party's key. Its protected header is the same
as an Action Record's, with typ `application/trail-state+cbor`. Its
payload is:

~~~ cddl
trail-state = {
  1 => uint,                    ; version (1)
  2 => uint,                    ; sequence, always increasing
  3 => uint,                    ; iat
  4 => uint,                    ; exp
  5 => tstr,                    ; WIMSE identifier of the party
  ? 6 => [+ bstr],              ; credentials: WIT, Attestation
                                ;   Result, capability cards
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
  policy: 1,            ; value: hash and location of a policy bundle
  hold: 2,              ; stop matching actions; no value
  require-level: 3,     ; value: level, the minimum
  require-witness: 4,   ; value: require-witness, below
  require-gate: 5,      ; value: [+ tstr], the approvers
  require-admission: 6, ; value: tstr, the admission issuer
  freshness-budget: 7,  ; value: freshness-budget, below
  require-record: 8     ; matching calls carry a record; no value
)

level = &( software: 0, tee: 1, isolated: 2 )

selector = {
  ? 1 => [+ bstr],              ; party pseudonyms
  ? 2 => [+ tstr],              ; agent identifiers or prefixes
  ? 3 => [+ kind],              ; record kinds
  ? 4 => [+ uint],              ; action classes (deployment-defined)
  ? 5 => [+ tstr],              ; counterparty identifiers
}

grant = {
  1 => tstr,                    ; authority's WIMSE identifier
  2 => [+ directive-type],      ; directive types it may issue
  3 => selector,                ; the scope it may issue them for
  4 => [+ uint],                ; data classes it may read
  ? 5 => uint,                  ; longest grace period it may allow
}
~~~
{: #fig-state title="State snapshot"}

A selector matches an action when every member present in it matches.
A member matches when any of its values matches. An absent member
matches everything.

Policy bundles, which can be large, are referenced from directives by
hash. The bundle itself is a record in the authority's trail, and edges
fetch it once. The policy claim of every Action Record is the SHA-256
hash of the deterministically encoded array of
`[authority, directive-id]` pairs for the directives the edge applied,
sorted by their encoding, so "it was allowed" can be checked later.

An authority's directives apply only within the grant its operator
gives it. An edge MUST ignore any directive that falls outside the
issuing authority's grant in the operator's current state snapshot. One
authority therefore cannot stop a whole swarm unless the operator agreed
to that, and the operator's grants are visible to every party that can
read its state.

Holds that apply to a single credential, such as one agent's WIT or one
admission token, SHOULD be expressed as a suspended entry in the
issuer's Token Status List {{OAUTH-STATUS-LIST}}, carried in the
issuer's state. Revocation is an invalid entry. This uses the same
mechanism that OAuth and WIMSE deployments use for credential status.

An edge MUST record a hold-applied record the first time it stops an
action because of a hold.

## How Directives Act {#directive-effects}

Directives differ in when they take effect:

| Directive | Takes effect | Effect |
|---|---|---|
| policy | As soon as the edge holds it | Matching actions are allowed or denied by the policy bundle |
| hold | As soon as the edge holds it | Matching actions stop |
| require-level, require-admission, require-record | As soon as the edge holds it | Calls that do not meet the requirement are refused |
| require-gate | On each matching action | The action waits for an approval ({{gates}}) |
| freshness-budget | When the authority's own state gets too old | Matching action classes stop, or degrade ({{degraded}}) |
| require-witness | When a required cosignature gets too old | The edge pauses covered agents, and receivers refuse their calls |
{: #tab-directives title="How directives act"}

## The Two Freshness Rules {#freshness}

Downstream freshness (authority to cell):
: An edge MUST hold, for every authority that has a grant in its
  operator's state, a snapshot that is within that authority's
  freshness budget for the action class ({{budgets}}). If it does not,
  it MUST NOT perform actions of that class covered by the grant,
  except as allowed by degraded mode ({{degraded}}).

Upstream freshness (cell to witness):
: For every require-witness directive that applies, a cell's latest
  checkpoint MUST carry countersignatures from at least k of the named
  witnesses, each with an iat no older than the directive's max-age.
  If it does
  not, the edge MUST NOT perform covered actions except as allowed by
  degraded mode, and receivers MUST reject calls from that cell that
  the directive covers.

An authority keeps its snapshot fresh by republishing it, even when
nothing changes. This is the heartbeat, and it proves that the
authority can still reach the edge. A witness keeps a cell fresh by
cosigning its checkpoints. This proves that the witness is still
receiving the cell's trail. An authority that cannot deliver a hold can
stop republishing, and the edge stops on its own once the snapshot ages
out; revocation therefore works when the agent is unreachable.

These two rules replace heartbeats, monitoring leases, checkpoint
deadlines, and maximum staleness bounds. For example, the rule "any
agent calling service Z must be monitored by regulator R" is a
require-witness directive from R, with Z as the counterparty in its
selector. If R stops watching, its cosignatures stop, the callers'
checkpoints age out, and Z's edge rejects their calls.

## Freshness Budgets per Action Class {#budgets}

Failing closed on every action turns an authority outage or a relay
partition into an outage of the business. Freshness is therefore set
per action class. A freshness-budget directive sets how old the
authority's own snapshot may get, and what happens when it is too old.
A require-witness directive, which can also select action classes,
sets which witnesses must have cosigned, how many of them, and how
recently. Each limit is set in exactly one place:

~~~ cddl
freshness-budget = {
  1 => uint,                    ; max age of the snapshot
  2 => &( pause: 0, degrade: 1 ),
  ? 3 => uint,                  ; grace period, if degrade
}

require-witness = [             ; require-witness value
  witnesses: [+ tstr],
  k: uint,                      ; how many must have cosigned
  max-age: uint,                ; max age of those cosignatures
]
~~~
{: #fig-budget title="Freshness budget and witness requirement"}

When several budgets or witness requirements apply to one action, the
edge MUST apply the strictest. When no budget applies, the snapshot's
exp is the limit and the behavior is pause. When a require-witness
limit is exceeded, the edge uses the behavior and grace period of the
freshness-budget that applies to the same action class, or pause if
there is none.

## Degraded Mode {#degraded}

An action class whose budget says degrade MAY continue for a bounded
grace period after the budget is exceeded:

* The grace period, and the classes allowed to degrade, MUST come from
  the last snapshot the edge received while it was fresh, so an
  attacker who blocks state cannot widen what continues.

* An operator grant caps the grace period an authority may allow.

* Gated actions and actions covered by a hold never degrade.

* Every degraded action is recorded with outcome degraded, and the next
  checkpoint lists them.

* When the grace period ends, the class pauses.

# Gates, Admission, Disclosure, and Media {#gates-etc}

## Gates {#gates}

For an action matched by a require-gate directive, the edge holds the
action and:

* appends a gate-request record to its trail, whose body is the hash of
  the pending call and any redacted context the approver needs; and

* sends the same record to each approver named in the directive, as a
  directed record, with disclosures for that approver's classes.

The approver publishes an approval record as a directed record for the
cell, and appends it to its own trail. The approval cites the
gate-request record with a cite of type peer-record and states the
decision. The edge continues only after a valid approval from every
named approver, and its send record cites each approval.

A require-gate selector can cover every action an agent takes. Every
such action then waits for a round trip to the approver, and stops if
the approver cannot be reached, because gated actions never degrade.

The approver can be a person. In that case the approval is signed with
a key bound to that person, for example through an OAuth client or a
platform authenticator, and the approval record names that person in a
redacted claim.

## Admission and Counting {#admission}

Limits that span many cells, such as "no more than 1,000 live agents in
swarm S", use admission tokens:

* The admission issuer is a party. It issues each agent an admission
  token: a CWT {{RFC8392}} with cnf set to the agent's WIT key, a unique
  cti, an expiry, and a status claim that points into the issuer's
  Token Status List {{OAUTH-STATUS-LIST}}. The token is sent as a
  directed record for the cell.

* The issuer never has more than the cap of valid tokens at once.

* A require-admission directive tells receivers to accept calls only
  from agents whose records carry the hash of a valid admission token in
  context-hashes.

* When an agent retires, its edge appends an admission-release record.
  The issuer then marks the token invalid and can issue a new one.

* Any witness can count live tokens from the issuer's trail, and the
  same token cited by two agents shows up in the trails.

Runtime limits use the token's expiry, and the C4M moqt-reval claim
({{authorization}}) for relay access. The cap binds only among parties
that check tokens; a server that does not run TRAIL does not refuse
calls without one.

## Disclosure on Request {#disclosure}

An authority that needs a claim it was not sent proactively sends a
disclosure request as a directed record for the cell. If the authority's
grant covers the claim's data class, the edge answers with disclosures
encrypted to the authority, as a directed record, and appends a
disclosure record to its trail. If it does not, the edge MUST refuse
and MUST record the refusal. Every look at protected data, and every
refused attempt, is therefore itself a witnessed record.

## Media Evidence {#media}

Agents that see and hear act on media that is already carried on MOQT
tracks. For each media Group the edge publishes, or receives and
forwards to the agent, it computes a Merkle root over the (Object ID,
payload hash) pairs of the Group, and appends a media-seal record that
names the track, the Group, the Object count, and the root.

An Action Record that depended on media cites it with a media-group
cite that names the track, the Group, and the Object count, with the
root as its hash. Changing one frame breaks the link. Inclusion proofs
let a verifier check single frames.

Live media MAY use Datagram Delivery Mode and delivery timeouts.
Evidence checks MUST use complete Groups fetched from an archive.

For approval before playback, each media Group uses its own Secure
Objects key. The publisher sends the key, encrypted with COSE-HPKE, to
the approver as a directed record. When the approver approves, it sends
the key to the entitled viewers as directed records. Nothing plays
before approval, and each release is a record in the approver's trail.

# Identity, Delegation, and Attestation {#identity}

TRAIL defines no new credential, token, delegation, or attestation
format:

| Need | Existing work | Use in TRAIL |
|---|---|---|
| Workload identity of edges, agents, authorities | WIMSE workload credentials {{WIMSE-CREDS}} | kid, iss, sub |
| Agent keys held outside the agent | WIMSE AIMS {{WIMSE-AIMS}} | The edge holds the agent's WIT key |
| Proof bound to presented tokens | WIMSE Workload Proof Token {{WIMSE-WPT}} | context-hashes |
| Per-request HTTP authentication | WIMSE HTTP Signatures {{WIMSE-HTTPSIG}} | HTTP binding |
| Principal and delegation chain | Transaction Tokens {{TXN-TOKENS}}, Token Exchange {{RFC8693}} | principal, tth, oth |
| Suspension and revocation | Token Status List {{OAUTH-STATUS-LIST}} | Holds on credentials, admission tokens |
| Sender-constrained access | DPoP {{RFC9449}} | C4M tokens, HTTP reads |
| Attestation of the edge key | RATS {{RFC9334}}, EAR {{RATS-EAR}} | Attestation levels |
| Statements and logs | SCITT {{RFC9943}}, {{RFC9162}}, tlog-tiles {{TLOG-TILES}} | Record format, trails, HTTP publication |
{: #tab-reuse title="Existing work used by TRAIL"}

## Identity: WIMSE

WIMSE workload credentials {{WIMSE-CREDS}} are mandatory to implement.

* Each edge, authority, witness, operator, and admission issuer has a
  WIMSE identifier and a WIT. The kid of every TRAIL signature is the
  thumbprint of the signer's WIT cnf key.

* Each agent has its own WIMSE identifier and WIT. The agent's WIT key
  is held by the edge, never by the agent, consistent with
  {{WIMSE-AIMS}}. The record's sub is the agent's identifier and wth is
  the hash of the agent's WIT.

* WITs are published in the party's state snapshot. A receiver reads
  the sender's state once and then verifies every call from that sender
  without further lookups.

The WIMSE architecture {{WIMSE-ARCH}} requires an intermediary that
changes identity context to make the change explicit and auditable. The
edge acts for the agent, and the Action Record makes that change
explicit.

Other identity schemes, such as X.509 certificates or SPIFFE SVIDs, can
be used through a profile that provides the same four things: a key
bound to an operator and an agent, a link to the principal, an
attestation result for the key, and a status mechanism.

## Principal and Delegation: OAuth {#principal}

TRAIL does not define a delegation format. The principal and the chain
of delegation are carried by OAuth Transaction Tokens {{TXN-TOKENS}},
whose sub, scope, tctx, and rctx claims describe who the call is for,
what it may do, and the context of the originating request, or by an
access token obtained through Token Exchange {{RFC8693}} with its act
chain. Fine-grained permissions can use Rich Authorization Requests
{{RFC9396}}.

Every record for an action taken on behalf of a principal MUST bind
both the principal and the delegation chain: member 12 identifies the
principal, usually redacted, and tth or oth carries the hash of the
token that conveys the delegation. The token itself travels in the
application message, for example in a PACE dialog's credentials, in
MCP _meta, or in an HTTP header. The Transaction Token's txn value
SHOULD be recorded, so that evidence and the operator's own transaction
logs can be correlated.

## Capability: SD-Card

An agent's capability card, for example an SD-Card {{SD-CARD}}, says
what the agent claims it can do. TRAIL uses it as follows:

* The card's cnf key SHOULD be the agent's WIT key. One key then covers
  discovery, workload authentication, relay authorization, and
  evidence.

* The issuer-signed card, with only its public disclosures and no key
  binding JWT, MAY be included in the cell's state snapshot.

* Context-specific disclosures and the key binding JWT are
  per-verifier. They MUST travel inside the dialog, bound to the peer's
  nonce ({{pace}}), not on a cached track.

* An Action Record records the sd_hash of the presentation the
  counterparty relied on in sdh.

## Attestation: RATS {#attestation}

The attestation level ({{tab-levels}}) of an edge key MUST come from an
Attestation Result {{RFC9334}} produced by a verifier the receiver
trusts, bound to the edge's WIT key, for example as in
{{WIMSE-ATTEST}}. An EAT Attestation Result {{RATS-EAR}} {{RFC9711}}
wrapped in a Conceptual Message Wrapper {{RFC9999}} is RECOMMENDED.
Trustworthiness claims {{RATS-AR4SI}} map to levels by deployment
policy. The Attestation Result is included in the edge's state snapshot
and refreshed before it expires. An edge MUST NOT state its own level.

# Authorization {#authorization}

## Publishing

Every party authenticates with a token bound with DPoP {{RFC9449}} to
its WIT key: over MOQT, a C4M token {{MoQ-C4M}} whose cnf.jkt is the
thumbprint of the WIT key; over HTTP, an OAuth access token. One rule
covers publishing: a party may publish only its own trail, its own
state, and directed records it sends. Over MOQT, its token grants
PUBLISH_NAMESPACE and PUBLISH only under its own prefix
("trail", S, P).

## Reading Trails and State {#read-access}

Trails, state, and directed records are not public. Only authorized
parties read them, and directed records only by their recipient. Each
swarm has an authorization server that issues the tokens for reading
its trails and state. How a swarm names it is an open issue
({{open-issues}}). That authorization server SHOULD NOT be
controlled by any single operator whose cells it covers, so that the
operator cannot choose who oversees it. Every authority that holds a
grant from an operator in the swarm, and every witness named in a
require-witness directive, MUST be able to obtain read access to the
trails and state it needs.

| Role | Reads |
|---|---|
| Cell edge | State of its authorities, operator, and peers; directed records addressed to it |
| Authority, witness | Trails and state under ("trail", S), or of the parties in its grant; directed records addressed to it |
| Archive | The trails and state it stores |
{: #tab-scopes title="Read access per role"}

In C4M terms, a cell edge with pseudonym P gets a scope with actions
PUBLISH_NAMESPACE and PUBLISH and namespace match `["trail", S, P]`,
and a second scope with action SUBSCRIBE and track match on its own
pseudonym under the to namespaces of its authorities. MOQT does not copy
authorization tokens from SUBSCRIBE_TRACKS into the PUBLISH messages it
produces, so a relay MUST check each forwarded PUBLISH against the
scope of the token presented with the SUBSCRIBE_TRACKS.

## Relays

Relays are not trusted for TRAIL's guarantees, but they can help
enforce them:

* C4M tokens for agents SHOULD carry moqt-reval and a status claim that
  refers to the issuer's Token Status List. A relay that revalidates
  then stops forwarding for a suspended agent, which adds a
  network-level kill switch to the edge's own.

* An edge relay can check that call Objects carry TRAIL_RECORD. It
  cannot verify the record. A C4M claim that requires a given Property
  on every published Object would let relays enforce presence
  ({{open-issues}}).

C4M and DPoP authenticate sessions and actions hop by hop. The Action
Record authenticates the content of each action end to end. They are
complementary and use the same key.

# Implementation Status

RFC Editor: please remove this section before publication.

This section follows {{RFC7942}}. No implementation is reported in this
version. {{open-issues}} lists the measurements the authors intend to
publish.

# Security Considerations

This section follows {{RFC3552}}.

## Sources of the Guarantees

| Guarantee | Comes from | Not from |
|---|---|---|
| Origin | Record signature by an edge key with an Attestation Result | TLS, which covers one hop; Secure Objects, whose keys are shared |
| Integrity and binding | Signature over trail-ref and call-ref | Relay or proxy behavior |
| Completeness | Merkle trail, dense Object IDs, witness cosignatures | Relay delivery or caching |
| Non-equivocation | Witnesses that refuse to cosign conflicting roots | The operator |
| Liveness of oversight | The two freshness rules | Relay heartbeats |
| Confidentiality of claims | SD-CWT redaction and COSE-HPKE to each recipient | Namespace access control |
| Access at the edge | C4M or OAuth with DPoP | Hop-by-hop enforcement beyond the first relay |
{: #tab-guarantees title="Sources of TRAIL guarantees"}

## Parties and Threats

Agent:
: Can ask to act. Cannot sign, reach the network except through the
  edge (subject to {{bypass}}), or skip policy.

Operator:
: Runs cells and chooses relays. Cannot edit records after the fact or
  show different histories to different witnesses without producing
  proof of equivocation. Cannot quietly stop publishing, because
  checkpoints stop being cosigned and the freshness rules pause the
  cell. Can choose which authorities to grant; that choice is itself in
  the operator's signed state.

Relay or proxy:
: Can route, cache, delay, or drop messages, and can change mutable
  Properties. Cannot forge records or move them, and dropped records
  show up as gaps. An HTTP proxy can drop the Trail-Record field; under
  require-record the receiver then refuses the call.

Authority:
: Can read its granted classes and issue directives within its grant.
  Cannot read other classes, because disclosures are encrypted
  to each recipient. Cannot act outside its grant.

Hardware root:
: Vouches for the keys of tee and isolated edges. It is the one root
  that every party must trust. Deployments SHOULD accept roots from
  more than one vendor.

## Shared Keys

Secure Objects keys are shared by a track's publisher and subscribers,
so any key holder can produce Objects that decrypt. TRAIL never relies
on them for origin. Disclosures use public-key encryption to each
recipient, so authorities cannot forge disclosures for each other.

## Replay {#replay}

A replayed call carries a record identifier the receiver has already
accepted. Receivers MUST keep the identifiers they accept for at least
the maximum record age they allow, and reject repeats for actions that
are not idempotent. Over HTTP, receivers MUST also reject a created
value outside their allowed clock skew. Even when a replay is not
caught, the call-ref binding means it can only reproduce the same call.
As PACE requires for its control messages, an edge MUST NOT send a
call that is not idempotent in QUIC 0-RTT data.

## Denial of Service {#dos}

The freshness rules let an attacker who can block state or
cosignatures, or an outage of an authority, a witness, or a relay,
pause agents. The edge cannot tell an outage from an attack. TRAIL
limits the damage instead: freshness budgets are set per action class
({{budgets}}), so an outage first stops high-impact actions; degraded
mode ({{degraded}}) lets classes the last fresh snapshot allowed
continue for a bounded grace period, recorded for later review; and an
operator grant bounds how much any one authority can stop. Deployments
SHOULD provide redundant relay paths for state and directed records,
SHOULD use more than one witness with a k-of-n rule, and SHOULD test
each action class under the loss of each authority and witness. Prefix
subscriptions near the root of the namespace can cause large numbers of
PUBLISH messages, and relays SHOULD limit them through authorization.

## Bypass {#bypass}

Everything in TRAIL depends on the edge seeing every call. At the
software level, an agent that escapes its sandbox, or a compromised
host, can go around the edge. At the tee level, the key is protected,
but an agent that escapes into the confidential VM can still go around
an edge that runs inside it. Only at the isolated level does hardware
in the network path see every packet. Receivers enforce the core rule
inbound, so a call that bypassed the edge is refused by every TRAIL
peer; bypass reaches only endpoints that do not run TRAIL. Authorities
can require the isolated level for high-impact action classes with
require-level.

## Limits {#limits}

TRAIL defines no way to judge whether an action is safe. An authority
judges with its own tools and acts through TRAIL: a gate before the
action, a policy at the edge as it happens, or a hold afterwards. It
sees content only through disclosures its grant covers. A signed,
approved message can still carry an injection. Covert channels through
timing, patterns in ordinary output, or public content are not blocked
and can only be found by audit. Interfaces through which a hardware
enforcement point exposes its signing key are not standardized.

# Privacy Considerations

A system that can prove what every agent did can also show what every
person using those agents did. TRAIL follows {{RFC6973}} and limits
this:

* TRAIL namespaces and resource paths contain only pseudonyms, which
  a party rotates by starting a new trail ({{trails}}). WIMSE
  identifiers, dialog identifiers, and principals appear in TRAIL only
  inside signed payloads, and sensitive claims are redacted. Dialog
  identifiers that a binding such as PACE puts in its own track names
  are covered by that binding's privacy rules.

* Trails and state are readable only by authorized parties
  ({{read-access}}).

* Authority rules run in the cell as directives, rather than on copies
  of the data.

* Disclosures go only to authorities entitled to a class, encrypted to
  each, and every disclosure is itself a record in the trail.

* Disclosure keys and archived disclosures SHOULD be tied to a purpose
  and a retention period, and destroyed when it ends.

* Publishers SHOULD pad records to limit size analysis.

Redacted records are still visible in full to every party that can
read the trail, so their unredacted claims, kinds, timing, and sizes,
and the pattern of who talks to whom, remain visible to those readers
and to relays and proxies. Deployments SHOULD redact every claim that
is not needed by all readers.

# IANA Considerations

## HTTP Field Name

IANA is requested to register Trail-Record in the "Hypertext Transfer
Protocol (HTTP) Field Name Registry" with status permanent, structured
type Item, and this document as reference.

## MOQ Properties

IANA is requested to register the following entry in the "MOQ
Properties" registry. Until a value is assigned, implementations can
use a value from the application-specific range 0x3800 to 0x3FFF.

| Type | Name | Scope | Reference |
|---|---|---|---|
| TBD1 (odd) | TRAIL_RECORD | Object | This document |
{: #tab-iana-moq title="MOQ Property"}

## Media Types

IANA is requested to register `application/trail-record+cbor` and
`application/trail-state+cbor` in the "Media Types" registry, with
the following template for each:

Type name:
: application

Subtype names:
: trail-record+cbor; trail-state+cbor

Required and optional parameters:
: none

Encoding considerations:
: binary (CBOR)

Security considerations:
: see {{security-considerations}}

Interoperability considerations:
: none

Published specification:
: this document

Applications that use this media type:
: TRAIL edges, authorities, witnesses, and verifiers

Fragment identifier considerations:
: as for application/cbor

Change controller:
: IETF

## TRAIL Registries

IANA is requested to create the registries "TRAIL Record Kinds",
"TRAIL Outcomes", "TRAIL Call Reference Types", "TRAIL Directive
Types", "TRAIL Cite Types", "TRAIL Record Members", "TRAIL State
Members", "TRAIL Context Hash Members", and "TRAIL Selector Members",
with the initial values in {{fig-record}} and {{fig-state}} and a
Specification Required policy. The member registries cover the
integer keys of the trail-record, trail-state, context-hashes, and
selector maps, so that extensions do not collide.

--- back

# Related Work {#related-work}

Several individual drafts address per-action evidence for agents or
workloads:

* {{WIMSE-AUTHZ-EVIDENCE}} defines signed "Permits" committing to a
  request before dispatch, anchored in SCITT. A TRAIL gate-request and
  approval pair plays the same role.

* {{WIMSE-AGENT-EVIDENCE}} points from a WIT to an evidence service and
  defines hash-chained execution receipts.

* {{SCITT-LOCAL-LOG}} defines a producer-operated log with signed
  checkpoints inside SCITT statements, which is close to the TRAIL
  trail.

* {{AGENTPROTO-VERIFIER}} discusses verifier evidence in agentproto.

* {{AAUTH}} signs every agent request with HTTP Message Signatures and
  keeps a mission log.

* {{AUDIT-ARCH}} proposes an architecture for auditing agent delegation
  and interaction traceability across domains.

TRAIL differs in four ways: the evidence is produced outside the agent
at an attested edge; the edge acts only while its oversight is fresh;
trails reach authorized parties the operator did not enumerate; and
control (holds, gates, admission) returns on the same evidence plane.
TRAIL aims to converge with these drafts on the record envelope, which
all of them base on COSE or SCITT.

# Examples {#examples}

## A Gated Payment Across Operators

1. Operator O grants financial regulator F the directive types
   require-gate, require-level, and hold for action class "payment",
   and data class "finance". The grant is in O's state.

2. F's state contains require-gate for class "payment" above a
   threshold, naming F as approver, and require-level of `isolated`
   for the same class.

3. Agent A in cell C_A, in dialog D, asks to send a payment instruction
   to agent B.

4. Edge A appends a gate-request record, with the amount redacted, and
   sends it with the "finance" disclosure encrypted to F as a directed
   record.

5. F appends an approval record and sends it to C_A as a directed
   record.

6. Edge A appends a send record that cites the approval and sends the
   call with the record attached.

7. Edge B checks the record, A's level, and A's cosigned checkpoint,
   delivers the call, and appends a recv record citing A's.

8. A health authority that can read the same trails sees that a payment
   happened, but cannot open the amount.

## A Monitored Service

Regulator R publishes require-witness with itself as witness, a maximum
age of 60 seconds, and service Z as counterparty. Z's edge rejects any
call whose sender's latest checkpoint lacks an R cosignature younger
than 60 seconds. R cosigns only trails it is receiving in full. If R
stops watching, calls to Z stop within a minute.

## A Swarm Size Cap

Admission issuer Q never has more than 1,000 valid admission tokens for
swarm S. Q's require-admission directive makes peers reject calls from
agents without a valid token hash. An agent spawned without a token has
nobody to talk to among TRAIL peers. Retiring an agent releases its
token, Q marks it invalid, and any witness can count valid tokens in
Q's trail.

# Requirements Raised at the IETF 126 BoF {#requirements}

The following requirements were raised during the agentproto BoF at
IETF 126, as recorded in the draft minutes {{AGENTPROTO-126}}. The
minutes are cited for the requirements only; their speakers have not
reviewed this document.

| Requirement raised at the BoF | Minutes | TRAIL mechanism |
|---|---|---|
| Revoking a session or agent is hard, especially if the agent is no longer available | Items 7 and 8 | Holds and Token Status List entries at the edge, and freshness rules under which silence also stops the agent ({{freshness}}) |
| Audit requires a durable owner of an agent's actions, not only runtime credentials | Items 7 and 8 | Every record binds the principal and the delegation chain ({{principal}}) |
| The protocol should enable audit signals that applications can use | Items 7 and 8 | Trails delivered live to authorized parties ({{trails}}) |
| Agent ownership mechanisms that can designate legal liability | Item 6 | The responder's acknowledgment cites the request ({{core-rule}}) |
{: #tab-requirements title="BoF requirements and TRAIL mechanisms"}

# Open Issues {#open-issues}

* Which party runs the authorization server for reading a swarm's
  trails and state ({{read-access}}).

* Whether a receiver, or its authority, can require that a sender's
  operator has granted a named authority certain rights before calls
  are accepted. Today only witnesses can be required.

* Whether verifiers publish Attestation Results themselves, in
  addition to edges carrying them in their state.

* Extensions needed by further use cases, tracked in the issue
  tracker: holds that match everything descending from an input or
  from a PACE dialog through its parent_dialog_ids,
  gate responses that substitute an action, approvals signed with a
  passkey, records that cite a model server's attestation, and
  admission tokens that meter spend.

* Default freshness budgets and grace periods per action class.

* Whether one binding should be mandatory to implement for the evidence
  plane.

* Whether the checkpoint and the tlog-checkpoint signed note can be one
  object, and how TRAIL witnesses relate to {{TLOG-WITNESS}}.

* A C4M claim that requires a given Property on every Object published
  under a scope, so that edge relays can enforce the presence of
  TRAIL_RECORD.

* A Token Status List delivery profile over MOQT, for use beyond TRAIL.

* Whether TRAIL_RECORD should be standardized as a general "per Object
  proof" Property for MOQT, with TRAIL as one user, which would amount
  to the first non-HTTP binding of the WIMSE Workload Proof Token.

* Coordination with WIMSE and SCITT individual drafts on per-action
  evidence, to avoid incompatible record formats.

* Measurements, for swarms of 1,000 and 10,000 agents, of the time from
  a hold being published to its enforcement over MOQT and over HTTP,
  the time from an authority going silent to covered agents pausing,
  trail throughput per edge, witness cosignature latency, and relay
  fan-out cost.

# Acknowledgments
{:numbered="false"}

This document draws on discussions of PACE, MCP over MOQT, C4M, and
SD-Card.
