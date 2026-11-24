// Package httpsig carries Action Records on HTTP messages (draft 6.1).
//
// The edge signs each message with an RFC 9421 HTTP Message Signature
// labeled "trail", covering the components the draft requires and the
// Trail-Record field. The record's http-ref holds a hash of the signature
// base for the covered components other than trail-record. This package
// implements only what TRAIL needs: one signature, ecdsa-p256-sha256,
// with the components listed below.
package httpsig

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Header names.
const (
	TrailRecord    = "Trail-Record"
	ContentDigest  = "Content-Digest"
	SignatureInput = "Signature-Input"
	Signature      = "Signature"
	label          = "trail"
)

var std = base64.StdEncoding

// Digest returns the RFC 9530 Content-Digest value for a body.
func Digest(body []byte) string {
	sum := sha256.Sum256(body)
	return "sha-256=:" + std.EncodeToString(sum[:]) + ":"
}

// Message is the part of a request or response the signature covers.
type Message struct {
	Request   bool
	Method    string // request only
	TargetURI string // request only
	Status    int    // response only
	Header    http.Header
}

// Components returns the components to cover, in order, given whether the
// message has content (6.1). trail-record is always last.
func Components(m Message, hasContent bool) []string {
	var c []string
	if m.Request {
		c = []string{"@method", "@target-uri"}
	} else {
		c = []string{"@status"}
	}
	if hasContent {
		c = append(c, "content-digest")
	}
	return append(c, "trail-record")
}

func componentValue(m Message, name string) (string, error) {
	switch name {
	case "@method":
		return m.Method, nil
	case "@target-uri":
		return m.TargetURI, nil
	case "@status":
		return strconv.Itoa(m.Status), nil
	}
	src := m.Header.Values(name)
	if len(src) == 0 {
		return "", fmt.Errorf("httpsig: missing %s", name)
	}
	vals := make([]string, len(src))
	for i, v := range src {
		vals[i] = strings.TrimSpace(v)
	}
	return strings.Join(vals, ", "), nil
}

// ComponentsHash is the components-hash of an http-ref: the SHA-256 of the
// signature base lines for the covered components other than trail-record,
// without the @signature-params line.
func ComponentsHash(m Message, components []string) ([]byte, error) {
	var b strings.Builder
	for _, c := range components {
		if c == "trail-record" {
			continue
		}
		v, err := componentValue(m, c)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(&b, "%q: %s\n", c, v)
	}
	sum := sha256.Sum256([]byte(b.String()))
	return sum[:], nil
}

func paramsLine(components []string, created int64, keyID string) string {
	q := make([]string, len(components))
	for i, c := range components {
		q[i] = strconv.Quote(c)
	}
	return fmt.Sprintf("(%s);created=%d;tag=%q;keyid=%q", strings.Join(q, " "), created, label, keyID)
}

func base(m Message, components []string, params string) (string, error) {
	var b strings.Builder
	for _, c := range components {
		v, err := componentValue(m, c)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&b, "%q: %s\n", c, v)
	}
	fmt.Fprintf(&b, "\"@signature-params\": %s", params)
	return b.String(), nil
}

// Sign adds Signature-Input and Signature to m.Header. Trail-Record and,
// if there is content, Content-Digest must already be set.
func Sign(m Message, components []string, created int64, key *ecdsa.PrivateKey, keyID string) error {
	params := paramsLine(components, created, keyID)
	b, err := base(m, components, params)
	if err != nil {
		return err
	}
	sum := sha256.Sum256([]byte(b))
	r, s, err := ecdsa.Sign(rand.Reader, key, sum[:])
	if err != nil {
		return err
	}
	sig := make([]byte, 64)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:])
	m.Header.Set(SignatureInput, label+"="+params)
	m.Header.Set(Signature, label+"=:"+std.EncodeToString(sig)+":")
	return nil
}

// Params is what a verifier learns from Signature-Input.
type Params struct {
	Components []string
	Created    int64
	KeyID      string
	raw        string
}

var (
	reInput   = regexp.MustCompile(`(?:^|,\s*)trail=(\([^)]*\)[^,]*)`)
	reSig     = regexp.MustCompile(`(?:^|,\s*)trail=:([A-Za-z0-9+/=]+):`)
	reComp    = regexp.MustCompile(`"([^"]+)"`)
	reCreated = regexp.MustCompile(`;created=(\d+)`)
	reKeyID   = regexp.MustCompile(`;keyid="([^"]*)"`)
	reTag     = regexp.MustCompile(`;tag="([^"]*)"`)
)

// ParseInput reads the "trail" signature's parameters from a header.
func ParseInput(h http.Header) (*Params, error) {
	mm := reInput.FindStringSubmatch(h.Get(SignatureInput))
	if mm == nil {
		return nil, errors.New("httpsig: no trail signature")
	}
	raw := mm[1]
	p := &Params{raw: raw}
	inner := raw[1:strings.Index(raw, ")")]
	for _, c := range reComp.FindAllStringSubmatch(inner, -1) {
		p.Components = append(p.Components, c[1])
	}
	if c := reCreated.FindStringSubmatch(raw); c != nil {
		p.Created, _ = strconv.ParseInt(c[1], 10, 64)
	} else {
		return nil, errors.New("httpsig: no created parameter")
	}
	if k := reKeyID.FindStringSubmatch(raw); k != nil {
		p.KeyID = k[1]
	} else {
		return nil, errors.New("httpsig: no keyid parameter")
	}
	if t := reTag.FindStringSubmatch(raw); t == nil || t[1] != label {
		return nil, errors.New("httpsig: tag is not \"trail\"")
	}
	return p, nil
}

// Verify checks the "trail" signature on m against pub and that it covers
// what the draft requires.
func Verify(m Message, p *Params, pub *ecdsa.PublicKey, hasContent bool) error {
	have := map[string]bool{}
	for _, c := range p.Components {
		have[c] = true
	}
	for _, need := range Components(m, hasContent) {
		if !have[need] {
			return fmt.Errorf("httpsig: signature does not cover %s", need)
		}
	}
	mm := reSig.FindStringSubmatch(m.Header.Get(Signature))
	if mm == nil {
		return errors.New("httpsig: no trail signature value")
	}
	sig, err := std.DecodeString(mm[1])
	if err != nil || len(sig) != 64 {
		return errors.New("httpsig: bad signature encoding")
	}
	b, err := base(m, p.Components, p.raw)
	if err != nil {
		return err
	}
	sum := sha256.Sum256([]byte(b))
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:])
	if !ecdsa.Verify(pub, sum[:], r, s) {
		return errors.New("httpsig: signature does not verify")
	}
	return nil
}

// EncodeRecord is the Trail-Record field value: an sf-binary item.
func EncodeRecord(rec []byte) string { return ":" + std.EncodeToString(rec) + ":" }

// DecodeRecord reads the Trail-Record field value.
func DecodeRecord(v string) ([]byte, error) {
	v = strings.TrimSpace(v)
	if len(v) < 2 || v[0] != ':' || v[len(v)-1] != ':' {
		return nil, errors.New("httpsig: Trail-Record is not a byte sequence")
	}
	return std.DecodeString(v[1 : len(v)-1])
}
