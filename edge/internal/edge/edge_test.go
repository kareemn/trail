package edge

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kareemn/trail/edge/internal/config"
	"github.com/kareemn/trail/edge/internal/httpsig"
	"github.com/kareemn/trail/edge/internal/keys"
	"github.com/kareemn/trail/edge/internal/record"
	"github.com/kareemn/trail/edge/internal/trail"
)

// pair builds Edge A (whose agent calls out) and Edge B (which fronts an
// echo agent), wired as peers. It returns A's agent-facing server.
type pair struct {
	a, b        *Edge
	aOut, bIn   *httptest.Server
	agentB      *httptest.Server
	dirA, dirB  string
	cfgA, cfgB  *config.Config
	lastAgentRq *http.Request
}

func newPair(t *testing.T) *pair {
	t.Helper()
	p := &pair{dirA: t.TempDir(), dirB: t.TempDir()}
	logger := log.New(io.Discard, "", 0)
	p.agentB = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.lastAgentRq = r
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.Write(append([]byte(`{"echo":`), append(body, '}')...))
	}))
	t.Cleanup(p.agentB.Close)

	keyA, _ := keys.LoadOrCreate(p.dirA)
	keyB, _ := keys.LoadOrCreate(p.dirB)

	p.cfgB = &config.Config{
		Identity: "spiffe://test/edge/b", Agent: "spiffe://test/agent/b", StateDir: p.dirB,
		Inbound:      &config.Inbound{Agent: strings.TrimPrefix(p.agentB.URL, "http://")},
		Peers:        []config.Peer{{ID: "spiffe://test/edge/a", JWK: keys.PublicJWK(&keyA.PublicKey)}},
		EpochRecords: 4, MaxSkewSeconds: 300,
	}
	var err error
	if p.b, err = New(p.cfgB, logger); err != nil {
		t.Fatal(err)
	}
	p.bIn = httptest.NewServer(p.b.InboundHandler())
	t.Cleanup(p.bIn.Close)

	p.cfgA = &config.Config{
		Identity: "spiffe://test/edge/a", Agent: "spiffe://test/agent/a", StateDir: p.dirA,
		Routes: []config.Route{
			{Prefix: "/a2a/b", Upstream: p.bIn.URL, Kind: "send", Aud: "spiffe://test/edge/b"},
			{Prefix: "/tool/plain", Upstream: p.agentB.URL, Kind: "tool-call",
				Headers: map[string]config.Secret{"Authorization": {Env: "TEST_TOKEN", Format: "Bearer %s"}}},
		},
		Peers:        []config.Peer{{ID: "spiffe://test/edge/b", JWK: keys.PublicJWK(&keyB.PublicKey)}},
		EpochRecords: 4, MaxSkewSeconds: 300,
	}
	if p.a, err = New(p.cfgA, logger); err != nil {
		t.Fatal(err)
	}
	p.aOut = httptest.NewServer(p.a.OutboundHandler())
	t.Cleanup(p.aOut.Close)
	return p
}

func (p *pair) setUpstream(prefix, u string) {
	for i := range p.a.routes {
		if p.a.routes[i].Prefix == prefix {
			p.a.routes[i].Upstream = u
			return
		}
	}
	panic("no route " + prefix)
}

func decodeAll(t *testing.T, dir string, pub interface{}) []*record.Decoded {
	t.Helper()
	recs, err := trail.ReadRecords(filepath.Join(dir, "records.log"))
	if err != nil {
		t.Fatal(err)
	}
	var out []*record.Decoded
	for _, raw := range recs {
		d, err := record.Verify(raw, func([]byte) (*ecdsaPub, error) { return pub.(*ecdsaPub), nil })
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, d)
	}
	return out
}

func kinds(ds []*record.Decoded) string {
	var s []string
	for _, d := range ds {
		k := d.Payload.Kind.String()
		if d.Payload.Outcome != nil {
			k += ":" + d.Payload.Outcome.String()
		}
		s = append(s, k)
	}
	return strings.Join(s, " ")
}

// The "records on both sides" block: A1 send, B1 recv cites A1, B2 send
// cites A1, A2 recv cites B2.
func TestRecordsOnBothSides(t *testing.T) {
	p := newPair(t)
	body := `{"jsonrpc":"2.0","id":7,"method":"tasks/send","params":{"q":"hi"}}`
	resp, err := http.Post(p.aOut.URL+"/a2a/b/tasks", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 || !strings.Contains(string(got), `"q":"hi"`) {
		t.Fatalf("status %d body %s", resp.StatusCode, got)
	}
	if resp.Header.Get(httpsig.TrailRecord) != "" {
		t.Error("agent A should not see TRAIL headers")
	}
	if p.lastAgentRq.Header.Get(httpsig.TrailRecord) != "" {
		t.Error("agent B should not see TRAIL headers")
	}

	a := decodeAll(t, p.dirA, &p.a.key.PublicKey)
	b := decodeAll(t, p.dirB, &p.b.key.PublicKey)
	if kinds(a) != "send:allowed recv:completed" {
		t.Fatalf("A trail: %s", kinds(a))
	}
	if kinds(b) != "recv:allowed send:completed" {
		t.Fatalf("B trail: %s", kinds(b))
	}
	a1, a2, b1, b2 := a[0], a[1], b[0], b[1]
	if !bytes.Equal(b1.Payload.Cites[0].Hash, a1.ID) || !bytes.Equal(b2.Payload.Cites[0].Hash, a1.ID) {
		t.Error("B1 and B2 must cite A1")
	}
	if !bytes.Equal(a2.Payload.Cites[0].Hash, b2.ID) {
		t.Error("A2 must cite B2")
	}
	if b1.Payload.Aud != "spiffe://test/edge/a" || a1.Payload.Aud != "spiffe://test/edge/b" {
		t.Error("aud not set")
	}
	if m, _ := a1.Payload.Body.(map[any]any); m["jsonrpc-method"] != "tasks/send" {
		t.Errorf("A1 body %v", a1.Payload.Body)
	}
}

func TestNoRouteIsDeniedAndRecorded(t *testing.T) {
	p := newPair(t)
	resp, err := http.Get(p.aOut.URL + "/elsewhere")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status %d", resp.StatusCode)
	}
	// A CONNECT through the proxy variables is refused too.
	req, _ := http.NewRequest(http.MethodConnect, p.aOut.URL, nil)
	req.Host = "example.com:443"
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("CONNECT status %d", resp.StatusCode)
	}
	if k := kinds(decodeAll(t, p.dirA, &p.a.key.PublicKey)); k != "send:denied send:denied" {
		t.Fatalf("A trail: %s", k)
	}
}

func TestCredentialsAddedAndHashed(t *testing.T) {
	p := newPair(t)
	lookupEnv = func(k string) string {
		if k == "TEST_TOKEN" {
			return "s3cret"
		}
		return ""
	}
	t.Cleanup(func() { lookupEnv = os.Getenv })
	resp, err := http.Get(p.aOut.URL + "/tool/plain/x")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if got := p.lastAgentRq.Header.Get("Authorization"); got != "Bearer s3cret" {
		t.Fatalf("Authorization %q", got)
	}
	a := decodeAll(t, p.dirA, &p.a.key.PublicKey)
	if kinds(a) != "tool-call:allowed recv:unsigned" {
		t.Fatalf("A trail: %s", kinds(a))
	}
	if a[0].Payload.Context == nil || len(a[0].Payload.Context.OTH["authorization"]) != 32 {
		t.Fatal("record must carry the credential's hash in oth")
	}
}

func TestInboundRefusals(t *testing.T) {
	p := newPair(t)
	// No record at all (core rule: drop).
	resp, err := http.Post(p.bIn.URL+"/x", "text/plain", strings.NewReader("hi"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("unsigned call: status %d", resp.StatusCode)
	}

	// Capture a valid signed call from A, then replay it and tamper with it.
	var captured *http.Request
	var capturedBody []byte
	spy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedBody, _ = io.ReadAll(r.Body)
		captured = r.Clone(r.Context())
		r.Body = io.NopCloser(bytes.NewReader(capturedBody))
		r.Host = strings.TrimPrefix(p.bIn.URL, "http://")
		p.b.InboundHandler().ServeHTTP(w, r)
	}))
	defer spy.Close()
	p.setUpstream("/a2a/b", spy.URL)
	// The spy changes the host, so the signed target URI no longer matches
	// what B reconstructs: B must refuse.
	resp, _ = http.Post(p.aOut.URL+"/a2a/b/t", "application/json", strings.NewReader(`{"a":1}`))
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("redirected call: status %d", resp.StatusCode)
	}

	send := func(body []byte, host string) int {
		req, _ := http.NewRequest(http.MethodPost, p.bIn.URL+"/t", bytes.NewReader(body))
		for k, v := range captured.Header {
			req.Header[k] = v
		}
		req.Host = host
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		return resp.StatusCode
	}
	spyHost := strings.TrimPrefix(spy.URL, "http://")
	// Same signed call, body changed: Content-Digest and call-ref break.
	if code := send([]byte(`{"a":2}`), spyHost); code != http.StatusForbidden {
		t.Fatalf("tampered body: status %d", code)
	}

	// Control: the untouched call to the host it was signed for is accepted,
	// so the refusals above were for the reasons tested. Sent again, it is
	// a replay.
	if code := send(capturedBody, spyHost); code != http.StatusOK {
		t.Fatalf("untouched call: status %d", code)
	}
	if code := send(capturedBody, spyHost); code != http.StatusForbidden {
		t.Fatalf("replayed call: status %d", code)
	}

	b := decodeAll(t, p.dirB, &p.b.key.PublicKey)
	want := "recv:denied recv:denied recv:denied recv:allowed checkpoint send:completed recv:denied" // epoch is 4 records
	if k := kinds(b); k != want {
		t.Fatalf("B trail: %s\nwant:     %s", k, want)
	}
	reasons := []string{"no Trail-Record", "signature does not verify", "Content-Digest", "", "", "", "replayed"}
	for i, r := range reasons {
		if r == "" {
			continue
		}
		m, _ := b[i].Payload.Body.(map[any]any)
		if got, _ := m["reason"].(string); !strings.Contains(got, r) {
			t.Errorf("record %d reason %q, want %q", i, got, r)
		}
	}
}

func TestReplayRefused(t *testing.T) {
	p := newPair(t)
	var hdr http.Header
	var body []byte
	tap := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)
		hdr = r.Header.Clone()
		r.Body = io.NopCloser(bytes.NewReader(body))
		p.b.InboundHandler().ServeHTTP(w, r)
	}))
	defer tap.Close()
	p.setUpstream("/a2a/b", tap.URL)
	resp, _ := http.Post(p.aOut.URL+"/a2a/b/t", "application/json", strings.NewReader(`{"a":1}`))
	if resp.StatusCode != 200 {
		t.Fatalf("first call: %d", resp.StatusCode)
	}
	if hdr.Get(httpsig.TrailRecord) == "" {
		t.Fatal("tap saw no record")
	}
	req, _ := http.NewRequest(http.MethodPost, tap.URL+"/t", bytes.NewReader(body))
	req.Header = hdr
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("replay: status %d", resp.StatusCode)
	}
}

func TestCheckpointsAndTamper(t *testing.T) {
	p := newPair(t)
	for i := 0; i < 5; i++ {
		resp, err := http.Post(p.aOut.URL+"/a2a/b/t", "application/json", strings.NewReader(`{"n":1}`))
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
	}
	p.a.Close()
	var out bytes.Buffer
	if err := VerifyTrail(p.dirA, &p.a.key.PublicKey, &out); err != nil {
		t.Fatalf("%v\n%s", err, out.String())
	}
	if c := strings.Count(out.String(), "checkpoint"); c < 3 {
		t.Fatalf("want at least 3 checkpoints, got %d\n%s", c, out.String())
	}

	// Flip one byte inside the third record's payload.
	path := filepath.Join(p.dirA, "records.log")
	data, _ := os.ReadFile(path)
	recs, _ := trail.ReadRecords(path)
	off := 0
	for i := 0; i < 2; i++ {
		off += 4 + len(recs[i])
	}
	data[off+4+len(recs[2])/2] ^= 0xff
	os.WriteFile(path, data, 0o600)
	if err := VerifyTrail(p.dirA, &p.a.key.PublicKey, nil); err == nil {
		t.Fatal("tampered trail verified")
	}
}
