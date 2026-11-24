package trail

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

func h(parts ...[]byte) []byte {
	d := sha256.New()
	for _, p := range parts {
		d.Write(p)
	}
	return d.Sum(nil)
}

// RFC 9162 section 2.1.1: leaf hash 0x00||d, node hash 0x01||l||r, and
// the split at the largest power of two below n.
func TestRoot(t *testing.T) {
	a, b, c := []byte("a"), []byte("b"), []byte("c")
	la, lb, lc := h([]byte{0}, a), h([]byte{0}, b), h([]byte{0}, c)
	if got := Root([][]byte{a}); !bytes.Equal(got, la) {
		t.Error("one leaf")
	}
	if got, want := Root([][]byte{a, b}), h([]byte{1}, la, lb); !bytes.Equal(got, want) {
		t.Error("two leaves")
	}
	if got, want := Root([][]byte{a, b, c}), h([]byte{1}, h([]byte{1}, la, lb), lc); !bytes.Equal(got, want) {
		t.Error("three leaves")
	}
	if got := Root(nil); !bytes.Equal(got, h()) {
		t.Error("empty tree")
	}
}
