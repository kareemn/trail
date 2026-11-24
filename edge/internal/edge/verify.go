package edge

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/fxamacker/cbor/v2"
	"github.com/kareemn/trail/edge/internal/record"
	"github.com/kareemn/trail/edge/internal/trail"
)

// VerifyTrail checks a trail the way a witness would before cosigning
// (8.3), using only the records and the edge's public key: every
// signature, every trail-ref, and every checkpoint's tree head. It writes
// one line per record to out when out is not nil.
func VerifyTrail(dir string, pub *ecdsa.PublicKey, out io.Writer) error {
	recs, err := trail.ReadRecords(filepath.Join(dir, "records.log"))
	if err != nil {
		return err
	}
	lookup := func([]byte) (*ecdsa.PublicKey, error) { return pub, nil }
	var leaves [][]byte
	var trailID []byte
	problems := 0
	for i, raw := range recs {
		d, err := record.Verify(raw, lookup)
		if err != nil {
			fmt.Fprintf(orDiscard(out), "%5d  BAD  %v\n", i, err)
			problems++
			id, _ := record.ID(raw)
			leaves = append(leaves, id)
			continue
		}
		note := ""
		p := d.Payload
		if trailID == nil {
			trailID = p.TrailRef.TrailID
		}
		if !bytes.Equal(p.TrailRef.TrailID, trailID) || p.TrailRef.Leaf != uint64(i) {
			note = fmt.Sprintf("  BAD trail-ref (leaf %d)", p.TrailRef.Leaf)
			problems++
		}
		if p.Kind == record.KindCheckpoint {
			var th TreeHead
			enc, _ := cbor.Marshal(p.Body)
			if cbor.Unmarshal(enc, &th) != nil {
				note += "  BAD checkpoint body"
				problems++
			} else if th.TreeSize != uint64(i) || !bytes.Equal(th.RootHash, trail.Root(leaves)) {
				note += "  BAD checkpoint root"
				problems++
			} else {
				note += fmt.Sprintf("  covers %d records, root %x", th.TreeSize, th.RootHash[:8])
			}
		}
		outcome := ""
		if p.Outcome != nil {
			outcome = p.Outcome.String()
		}
		fmt.Fprintf(orDiscard(out), "%5d  %-11s %-9s %s%s\n", i, p.Kind, outcome, summarize(p), note)
		leaves = append(leaves, d.ID)
	}
	if problems > 0 {
		return fmt.Errorf("%d problem(s) in %d records", problems, len(recs))
	}
	if len(recs) == 0 {
		return errors.New("trail is empty")
	}
	return nil
}

func summarize(p record.Payload) string {
	parts := ""
	if m, ok := p.Body.(map[any]any); ok {
		for _, k := range []string{"event", "method", "host", "path", "tool", "jsonrpc-method", "target", "reason", "status", "exit_code"} {
			if v, ok := m[k]; ok {
				parts += fmt.Sprintf("%s=%v ", k, v)
			}
		}
	}
	if p.Aud != "" {
		parts += "aud=" + p.Aud + " "
	}
	for _, c := range p.Cites {
		parts += fmt.Sprintf("cites=%x ", c.Hash[:6])
	}
	return parts
}

func orDiscard(w io.Writer) io.Writer {
	if w == nil {
		return io.Discard
	}
	return w
}
