// Package record builds, signs, and verifies Action Records (draft section 3).
//
// An Action Record is a COSE_Sign1 whose payload is a deterministically
// encoded CBOR map. This package covers the members the edge uses today;
// the MOQT call-ref and SD-CWT redaction are not implemented yet.
package record

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/go-cose"
)

// MediaType is the typ of an Action Record (3.1).
const MediaType = "application/trail-record+cbor"

// Kind is the record kind (3.2).
type Kind uint

const (
	KindSend Kind = iota + 1
	KindRecv
	KindToolCall
	KindModelCall
	KindStateWrite
	KindStateRead
	KindGateRequest
	KindApproval
	KindHoldApplied
	KindDisclosure
	KindCheckpoint
	KindMediaSeal
	KindLifecycle
	KindAdmissionRelease
	KindRestart
	KindCosignature
)

var kindNames = map[Kind]string{
	KindSend: "send", KindRecv: "recv", KindToolCall: "tool-call", KindModelCall: "model-call",
	KindStateWrite: "state-write", KindStateRead: "state-read", KindGateRequest: "gate-request",
	KindApproval: "approval", KindHoldApplied: "hold-applied", KindDisclosure: "disclosure",
	KindCheckpoint: "checkpoint", KindMediaSeal: "media-seal", KindLifecycle: "lifecycle",
	KindAdmissionRelease: "admission-release", KindRestart: "restart", KindCosignature: "cosignature",
}

func (k Kind) String() string {
	if s, ok := kindNames[k]; ok {
		return s
	}
	return fmt.Sprintf("kind(%d)", uint(k))
}

// ParseKind maps a kind name from configuration to a Kind.
func ParseKind(s string) (Kind, error) {
	for k, n := range kindNames {
		if n == s {
			return k, nil
		}
	}
	return 0, fmt.Errorf("unknown record kind %q", s)
}

// Outcome is the outcome member (3.2).
type Outcome uint

const (
	Allowed Outcome = iota
	Denied
	Held
	Gated
	Completed
	Failed
	Degraded
	Unsigned
)

var outcomeNames = []string{"allowed", "denied", "held", "gated", "completed", "failed", "degraded", "unsigned"}

func (o Outcome) String() string {
	if int(o) < len(outcomeNames) {
		return outcomeNames[o]
	}
	return fmt.Sprintf("outcome(%d)", uint(o))
}

// CiteType is the cite-type (3.2).
type CiteType uint

const (
	CiteInput CiteType = iota + 1
	CitePeerRecord
	CiteMediaGroup
	CiteDirective
	CiteApproval
)

// TrailRef is [trail-id, leaf].
type TrailRef struct {
	_       struct{} `cbor:",toarray"`
	TrailID []byte
	Leaf    uint64
}

// HTTPRef is the http-ref form of call-ref: [2, components-hash, created].
type HTTPRef struct {
	_              struct{} `cbor:",toarray"`
	Type           uint
	ComponentsHash []byte
	Created        uint64
}

// Cite is [cite-type, ref, hash]. Only trail-ref refs are used so far.
type Cite struct {
	_    struct{} `cbor:",toarray"`
	Type CiteType
	Ref  TrailRef
	Hash []byte
}

// ContextHashes binds the tokens presented with a call, by hash (3.2).
type ContextHashes struct {
	WTH []byte            `cbor:"1,keyasint,omitempty"`
	TTH []byte            `cbor:"2,keyasint,omitempty"`
	OTH map[string][]byte `cbor:"3,keyasint,omitempty"`
	SDH []byte            `cbor:"4,keyasint,omitempty"`
	ADM []byte            `cbor:"5,keyasint,omitempty"`
}

// Payload is the trail-record map.
type Payload struct {
	Version  uint           `cbor:"1,keyasint"`
	Kind     Kind           `cbor:"2,keyasint"`
	TrailRef TrailRef       `cbor:"3,keyasint"`
	CallRef  *HTTPRef       `cbor:"4,keyasint,omitempty"`
	Policy   []byte         `cbor:"5,keyasint"`
	Outcome  *Outcome       `cbor:"6,keyasint,omitempty"`
	Aud      string         `cbor:"8,keyasint,omitempty"`
	Context  *ContextHashes `cbor:"9,keyasint,omitempty"`
	Cites    []Cite         `cbor:"10,keyasint,omitempty"`
	Body     any            `cbor:"11,keyasint,omitempty"`
	Class    *uint          `cbor:"13,keyasint,omitempty"`
}

// O returns a pointer to an outcome, for building payloads.
func O(o Outcome) *Outcome { return &o }

// U returns a pointer to a uint, for building payloads.
func U(u uint) *uint { return &u }

var (
	encMode cbor.EncMode
	decMode cbor.DecMode
)

func init() {
	var err error
	if encMode, err = cbor.CoreDetEncOptions().EncMode(); err != nil {
		panic(err)
	}
	if decMode, err = (cbor.DecOptions{DefaultMapType: nil}).DecMode(); err != nil {
		panic(err)
	}
}

// Marshal encodes v with deterministic CBOR (RFC 8949 section 4.2.1).
func Marshal(v any) ([]byte, error) { return encMode.Marshal(v) }

// EmptyPolicy is the policy claim when the edge applied no directives: the
// SHA-256 of the deterministic encoding of an empty array (9.1).
func EmptyPolicy() []byte {
	enc, _ := Marshal([]any{})
	sum := sha256.Sum256(enc)
	return sum[:]
}

// Signer signs records for one edge.
type Signer struct {
	Key *ecdsa.PrivateKey
	Kid []byte // JWK thumbprint of Key's public key
	Iss string // the edge's identifier
	Now func() time.Time
}

// Sign returns the encoded COSE_Sign1 for a payload, with sub naming the
// agent that acted.
func (s *Signer) Sign(sub string, p *Payload) ([]byte, error) {
	if p.Version == 0 {
		p.Version = 1
	}
	if p.Policy == nil {
		p.Policy = EmptyPolicy()
	}
	payload, err := Marshal(p)
	if err != nil {
		return nil, err
	}
	now := time.Now
	if s.Now != nil {
		now = s.Now
	}
	ph := cose.ProtectedHeader{}
	ph.SetAlgorithm(cose.AlgorithmES256)
	ph[cose.HeaderLabelType] = MediaType
	ph[cose.HeaderLabelKeyID] = s.Kid
	ph[cose.HeaderLabelCWTClaims] = map[any]any{
		int64(1): s.Iss,
		int64(2): sub,
		int64(6): now().Unix(),
	}
	// Encode the protected header deterministically so the record
	// identifier is stable across implementations.
	rawProt, err := Marshal(map[any]any(ph))
	if err != nil {
		return nil, err
	}
	msg := &cose.Sign1Message{
		Headers: cose.Headers{RawProtected: mustBstr(rawProt), Protected: ph},
		Payload: payload,
	}
	signer, err := cose.NewSigner(cose.AlgorithmES256, s.Key)
	if err != nil {
		return nil, err
	}
	if err := msg.Sign(rand.Reader, nil, signer); err != nil {
		return nil, err
	}
	return msg.MarshalCBOR()
}

func mustBstr(b []byte) cbor.RawMessage {
	enc, err := Marshal(b)
	if err != nil {
		panic(err)
	}
	return enc
}

// Decoded is a record after signature verification.
type Decoded struct {
	Raw     []byte
	ID      []byte
	Kid     []byte
	Iss     string
	Sub     string
	Iat     int64
	Payload Payload
}

// KeyLookup returns the public key for a kid, or an error if unknown.
type KeyLookup func(kid []byte) (*ecdsa.PublicKey, error)

// Verify decodes a record, checks its typ and signature, and computes its
// identifier.
func Verify(raw []byte, lookup KeyLookup) (*Decoded, error) {
	var msg cose.Sign1Message
	if err := msg.UnmarshalCBOR(raw); err != nil {
		return nil, fmt.Errorf("record: %w", err)
	}
	ph := msg.Headers.Protected
	if typ, _ := ph[cose.HeaderLabelType].(string); typ != MediaType {
		return nil, fmt.Errorf("record: typ %v is not %s", ph[cose.HeaderLabelType], MediaType)
	}
	alg, err := ph.Algorithm()
	if err != nil || alg != cose.AlgorithmES256 {
		return nil, fmt.Errorf("record: unsupported alg %v", alg)
	}
	kid, _ := ph[cose.HeaderLabelKeyID].([]byte)
	if len(kid) == 0 {
		return nil, errors.New("record: no kid")
	}
	pub, err := lookup(kid)
	if err != nil {
		return nil, err
	}
	verifier, err := cose.NewVerifier(cose.AlgorithmES256, pub)
	if err != nil {
		return nil, err
	}
	if err := msg.Verify(nil, verifier); err != nil {
		return nil, fmt.Errorf("record: signature: %w", err)
	}
	d := &Decoded{Raw: raw, Kid: kid}
	if claims, ok := ph[cose.HeaderLabelCWTClaims].(map[any]any); ok {
		d.Iss, _ = claimString(claims, 1)
		d.Sub, _ = claimString(claims, 2)
		d.Iat, _ = claimInt(claims, 6)
	}
	if d.Iss == "" || d.Sub == "" {
		return nil, errors.New("record: CWT claims need iss and sub")
	}
	if err := decMode.Unmarshal(msg.Payload, &d.Payload); err != nil {
		return nil, fmt.Errorf("record: payload: %w", err)
	}
	if d.ID, err = ID(raw); err != nil {
		return nil, err
	}
	return d, nil
}

// ID is the record identifier: the SHA-256 of the COSE_Sign1 encoded with
// an empty unprotected header map (3.2).
func ID(raw []byte) ([]byte, error) {
	var arr []cbor.RawMessage
	var tag cbor.RawTag
	body := raw
	if err := decMode.Unmarshal(raw, &tag); err == nil && tag.Number == 18 {
		body = tag.Content
	}
	if err := decMode.Unmarshal(body, &arr); err != nil || len(arr) != 4 {
		return nil, errors.New("record: not a COSE_Sign1")
	}
	emptyMap, _ := Marshal(map[any]any{})
	arr[1] = emptyMap
	inner, err := Marshal(arr)
	if err != nil {
		return nil, err
	}
	enc, err := Marshal(cbor.RawTag{Number: 18, Content: inner})
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(enc)
	return sum[:], nil
}

// Equal reports whether two identifiers match.
func Equal(a, b []byte) bool { return bytes.Equal(a, b) }

func claimString(m map[any]any, label int64) (string, bool) {
	for k, v := range m {
		if n, ok := toInt(k); ok && n == label {
			s, ok := v.(string)
			return s, ok
		}
	}
	return "", false
}

func claimInt(m map[any]any, label int64) (int64, bool) {
	for k, v := range m {
		if n, ok := toInt(k); ok && n == label {
			return toInt(v)
		}
	}
	return 0, false
}

func toInt(v any) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case uint64:
		return int64(n), true
	case int:
		return int64(n), true
	}
	return 0, false
}
