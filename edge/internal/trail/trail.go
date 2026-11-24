// Package trail keeps an edge's trail: its Action Records in order, as an
// RFC 9162 Merkle tree whose leaves are record identifiers (draft 8.1),
// closed by checkpoints (8.2).
//
// Records are stored in an append-only file, each as a 4-byte big-endian
// length and the encoded COSE_Sign1. Append returns only after fsync,
// which is the edge's acknowledged append before a call leaves (6.4).
// Publishing the trail as tlog tiles is not implemented yet.
package trail

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	recordsFile = "records.log"
	trailIDFile = "trail-id"
)

// Trail is one append-only trail.
type Trail struct {
	mu      sync.Mutex
	dir     string
	f       *os.File
	id      []byte
	leaves  [][]byte // record identifiers, in order
	records [][]byte // encoded records, in order
	idFn    func([]byte) ([]byte, error)
}

// Open opens or creates the trail in dir. idFn computes a record's
// identifier from its encoding.
func Open(dir string, idFn func([]byte) ([]byte, error)) (*Trail, error) {
	t := &Trail{dir: dir, idFn: idFn}
	idHex, err := os.ReadFile(filepath.Join(dir, trailIDFile))
	switch {
	case errors.Is(err, os.ErrNotExist):
		if err := t.newID(); err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	default:
		if t.id, err = hex.DecodeString(string(idHex)); err != nil || len(t.id) != 16 {
			return nil, fmt.Errorf("trail: bad trail-id file")
		}
	}
	if err := t.load(); err != nil {
		return nil, err
	}
	t.f, err = os.OpenFile(filepath.Join(dir, recordsFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Trail) newID() error {
	t.id = make([]byte, 16)
	if _, err := rand.Read(t.id); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(t.dir, trailIDFile), []byte(hex.EncodeToString(t.id)), 0o600)
}

func (t *Trail) load() error {
	recs, err := ReadRecords(filepath.Join(t.dir, recordsFile))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, r := range recs {
		id, err := t.idFn(r)
		if err != nil {
			return err
		}
		t.records = append(t.records, r)
		t.leaves = append(t.leaves, id)
	}
	return nil
}

// ReadRecords reads every encoded record from a records file.
func ReadRecords(path string) ([][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var out [][]byte
	for {
		var n uint32
		if err := binary.Read(r, binary.BigEndian, &n); err != nil {
			if errors.Is(err, io.EOF) {
				return out, nil
			}
			return nil, fmt.Errorf("trail: truncated length after %d records", len(out))
		}
		buf := make([]byte, n)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, fmt.Errorf("trail: truncated record %d", len(out))
		}
		out = append(out, buf)
	}
}

// ID returns the trail identifier.
func (t *Trail) ID() []byte { return append([]byte(nil), t.id...) }

// Size returns the number of records in the trail.
func (t *Trail) Size() uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return uint64(len(t.leaves))
}

// Append assigns the next leaf index, asks build for the encoded record
// at that index, and appends it durably. It returns the record and its
// identifier.
func (t *Trail) Append(build func(trailID []byte, leaf uint64) ([]byte, error)) ([]byte, []byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	leaf := uint64(len(t.leaves))
	rec, err := build(t.id, leaf)
	if err != nil {
		return nil, nil, err
	}
	id, err := t.idFn(rec)
	if err != nil {
		return nil, nil, err
	}
	var hdr [4]byte
	binary.BigEndian.PutUint32(hdr[:], uint32(len(rec)))
	if _, err := t.f.Write(append(hdr[:], rec...)); err != nil {
		return nil, nil, err
	}
	if err := t.f.Sync(); err != nil {
		return nil, nil, err
	}
	t.records = append(t.records, rec)
	t.leaves = append(t.leaves, id)
	return rec, id, nil
}

// Head returns the size and RFC 9162 root of the current tree.
func (t *Trail) Head() (uint64, []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return uint64(len(t.leaves)), Root(t.leaves)
}

// Close closes the records file.
func (t *Trail) Close() error { return t.f.Close() }

// Root is the RFC 9162 Merkle Tree Hash of a list of leaves.
func Root(leaves [][]byte) []byte {
	if len(leaves) == 0 {
		sum := sha256.Sum256(nil)
		return sum[:]
	}
	if len(leaves) == 1 {
		h := sha256.New()
		h.Write([]byte{0})
		h.Write(leaves[0])
		return h.Sum(nil)
	}
	k := 1
	for k*2 < len(leaves) {
		k *= 2
	}
	h := sha256.New()
	h.Write([]byte{1})
	h.Write(Root(leaves[:k]))
	h.Write(Root(leaves[k:]))
	return h.Sum(nil)
}

// AppendCheckpoint appends a record built from the current tree head:
// size and root cover every record before it, and its own leaf is size
// (8.2). Holding the lock keeps another record from landing in between.
func (t *Trail) AppendCheckpoint(build func(trailID []byte, size uint64, root []byte) ([]byte, error)) ([]byte, []byte, error) {
	t.mu.Lock()
	size := uint64(len(t.leaves))
	root := Root(t.leaves)
	t.mu.Unlock()
	var stale bool
	rec, id, err := t.Append(func(trailID []byte, leaf uint64) ([]byte, error) {
		if leaf != size {
			stale = true
			return nil, errors.New("trail: tree grew while checkpointing")
		}
		return build(trailID, size, root)
	})
	if stale {
		return t.AppendCheckpoint(build)
	}
	return rec, id, err
}
