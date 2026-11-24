// Package keys holds the edge's signing key and computes JWK thumbprints.
//
// The draft names every TRAIL key by the JWK Thumbprint (RFC 7638) of its
// public key: kid in an Action Record and keyid in an HTTP signature (3.1,
// 6.1). The private key never leaves the edge's state directory, which the
// sandbox does not allow.
package keys

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
)

// JWK is the public part of a P-256 key in JWK form.
type JWK struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

const keyFile = "edge.key"

var b64 = base64.RawURLEncoding

// LoadOrCreate returns the edge key in dir, creating it on first use.
func LoadOrCreate(dir string) (*ecdsa.PrivateKey, error) {
	path := filepath.Join(dir, keyFile)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		der, err := x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			return nil, err
		}
		pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
		if err := os.WriteFile(path, pemBytes, 0o600); err != nil {
			return nil, err
		}
		return key, nil
	}
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("%s: no PEM block", path)
	}
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := k.(*ecdsa.PrivateKey)
	if !ok || key.Curve != elliptic.P256() {
		return nil, fmt.Errorf("%s: not a P-256 key", path)
	}
	return key, nil
}

// PublicJWK returns the JWK of a P-256 public key.
func PublicJWK(pub *ecdsa.PublicKey) JWK {
	return JWK{
		Kty: "EC",
		Crv: "P-256",
		X:   b64.EncodeToString(pad32(pub.X.Bytes())),
		Y:   b64.EncodeToString(pad32(pub.Y.Bytes())),
	}
}

// ParseJWK returns the public key a JWK describes.
func ParseJWK(j JWK) (*ecdsa.PublicKey, error) {
	if j.Kty != "EC" || j.Crv != "P-256" {
		return nil, fmt.Errorf("unsupported JWK %s/%s", j.Kty, j.Crv)
	}
	x, err := b64.DecodeString(j.X)
	if err != nil {
		return nil, err
	}
	y, err := b64.DecodeString(j.Y)
	if err != nil {
		return nil, err
	}
	pub := &ecdsa.PublicKey{Curve: elliptic.P256(), X: new(big.Int).SetBytes(x), Y: new(big.Int).SetBytes(y)}
	if !pub.Curve.IsOnCurve(pub.X, pub.Y) {
		return nil, errors.New("JWK point is not on P-256")
	}
	return pub, nil
}

// Thumbprint is the RFC 7638 JWK Thumbprint (SHA-256) of a public key.
func Thumbprint(pub *ecdsa.PublicKey) []byte {
	j := PublicJWK(pub)
	// RFC 7638: required members only, lexicographic order, no spaces.
	canon := fmt.Sprintf(`{"crv":%q,"kty":%q,"x":%q,"y":%q}`, j.Crv, j.Kty, j.X, j.Y)
	sum := sha256.Sum256([]byte(canon))
	return sum[:]
}

// KeyID is the thumbprint as the base64url text used in HTTP keyid.
func KeyID(pub *ecdsa.PublicKey) string { return b64.EncodeToString(Thumbprint(pub)) }

// MarshalJWK is a convenience for printing a public key.
func MarshalJWK(pub *ecdsa.PublicKey) ([]byte, error) {
	return json.MarshalIndent(PublicJWK(pub), "", "  ")
}

func pad32(b []byte) []byte {
	if len(b) >= 32 {
		return b
	}
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}
