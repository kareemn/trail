package record

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/go-cose"
)

func TestSignVerifyAndID(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	s := &Signer{Key: key, Kid: []byte("kid"), Iss: "spiffe://t/edge"}
	raw, err := s.Sign("spiffe://t/agent", &Payload{Kind: KindSend, Outcome: O(Allowed), TrailRef: TrailRef{TrailID: make([]byte, 16), Leaf: 3}})
	if err != nil {
		t.Fatal(err)
	}
	d, err := Verify(raw, func([]byte) (*ecdsa.PublicKey, error) { return &key.PublicKey, nil })
	if err != nil {
		t.Fatal(err)
	}
	if d.Iss != "spiffe://t/edge" || d.Sub != "spiffe://t/agent" || d.Payload.TrailRef.Leaf != 3 || *d.Payload.Outcome != Allowed {
		t.Fatalf("decoded %+v", d)
	}

	// Adding an unprotected header (an inclusion proof, say) must not
	// change the identifier (3.2).
	var msg cose.Sign1Message
	if err := msg.UnmarshalCBOR(raw); err != nil {
		t.Fatal(err)
	}
	msg.Headers.RawUnprotected = nil
	msg.Headers.Unprotected = cose.UnprotectedHeader{int64(396): []byte("proof")}
	withProof, err := msg.MarshalCBOR()
	if err != nil {
		t.Fatal(err)
	}
	id2, err := ID(withProof)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(withProof, raw) || !bytes.Equal(id2, d.ID) {
		t.Fatal("identifier changed with the unprotected header")
	}

	// The payload is deterministically encoded.
	var m map[int]cbor.RawMessage
	if err := cbor.Unmarshal(msg.Payload, &m); err != nil {
		t.Fatal(err)
	}
	again, _ := Marshal(m)
	if !bytes.Equal(again, msg.Payload) {
		t.Fatal("payload is not in deterministic encoding")
	}
}

func TestVerifyRejectsWrongKey(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	other, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	s := &Signer{Key: key, Kid: []byte("kid"), Iss: "i"}
	raw, _ := s.Sign("s", &Payload{Kind: KindSend})
	if _, err := Verify(raw, func([]byte) (*ecdsa.PublicKey, error) { return &other.PublicKey, nil }); err == nil {
		t.Fatal("verified with the wrong key")
	}
}
