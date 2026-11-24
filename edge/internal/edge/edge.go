// Package edge is the TRAIL edge: the only way traffic leaves or enters
// the cell (draft 1.6, 4).
//
// Outbound, the agent calls the edge on a localhost port. The edge matches
// the call to a configured route, writes an Action Record to its trail,
// signs the request, adds Trail-Record, and sends it. Anything without a
// route is denied and recorded. Inbound, the edge verifies the caller's
// record and signature, records a recv that cites it, delivers the call to
// the agent, and answers with its own signed record.
//
// Not implemented yet: state snapshots and directives (9), witnesses (8.3),
// gates, admission, and disclosure (10), SD-CWT redaction (3.4), streamed
// responses (6.1), and the MOQT binding (5).
package edge

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/kareemn/trail/edge/internal/config"
	"github.com/kareemn/trail/edge/internal/httpsig"
	"github.com/kareemn/trail/edge/internal/keys"
	"github.com/kareemn/trail/edge/internal/record"
	"github.com/kareemn/trail/edge/internal/trail"
)

const maxBody = 32 << 20

// lookupEnv reads credentials from the edge's own environment.
var lookupEnv = os.Getenv

// HealthPath answers on the agent-facing port so a sandboxed process can
// check that its one way out works.
const HealthPath = "/.well-known/trail-edge"

type peer struct {
	id  string
	pub *ecdsa.PublicKey
}

// Edge is one running edge.
type Edge struct {
	cfg    *config.Config
	key    *ecdsa.PrivateKey
	keyID  string
	signer *record.Signer
	trail  *trail.Trail
	routes []config.Route
	peers  map[string]peer // by base64url kid
	client *http.Client
	logger *log.Logger
	now    func() time.Time

	mu        sync.Mutex
	sinceCkpt int
	seen      map[string]time.Time // record ids accepted inbound
}

// New opens the edge's key and trail.
func New(cfg *config.Config, logger *log.Logger) (*Edge, error) {
	key, err := keys.LoadOrCreate(cfg.StateDir)
	if err != nil {
		return nil, err
	}
	tr, err := trail.Open(cfg.StateDir, record.ID)
	if err != nil {
		return nil, err
	}
	e := &Edge{
		cfg:    cfg,
		key:    key,
		keyID:  keys.KeyID(&key.PublicKey),
		trail:  tr,
		peers:  map[string]peer{},
		client: &http.Client{Timeout: 5 * time.Minute, CheckRedirect: noRedirects},
		logger: logger,
		now:    time.Now,
		seen:   map[string]time.Time{},
	}
	e.signer = &record.Signer{Key: key, Kid: keys.Thumbprint(&key.PublicKey), Iss: cfg.Identity, Now: e.now}
	for _, p := range cfg.Peers {
		pub, err := keys.ParseJWK(p.JWK)
		if err != nil {
			return nil, fmt.Errorf("peer %s: %w", p.ID, err)
		}
		e.peers[keys.KeyID(pub)] = peer{id: p.ID, pub: pub}
	}
	e.routes = append(e.routes, cfg.Routes...)
	sort.Slice(e.routes, func(i, j int) bool { return len(e.routes[i].Prefix) > len(e.routes[j].Prefix) })
	return e, nil
}

func noRedirects(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }

// KeyID is the edge key's thumbprint in base64url.
func (e *Edge) KeyID() string { return e.keyID }

// PublicKey is the edge's public key.
func (e *Edge) PublicKey() *ecdsa.PublicKey { return &e.key.PublicKey }

// Routes returns the configured routes.
func (e *Edge) Routes() []config.Route { return e.routes }

// Close writes a final checkpoint and closes the trail.
func (e *Edge) Close() error {
	if e.trail.Size() > 0 {
		if err := e.Checkpoint(); err != nil {
			e.logger.Printf("final checkpoint: %v", err)
		}
	}
	return e.trail.Close()
}

// Append signs a record at the next leaf of the trail and appends it.
func (e *Edge) Append(p *record.Payload) (raw, id []byte, err error) {
	raw, id, err = e.trail.Append(func(trailID []byte, leaf uint64) ([]byte, error) {
		p.TrailRef = record.TrailRef{TrailID: trailID, Leaf: leaf}
		return e.signer.Sign(e.cfg.Agent, p)
	})
	if err != nil {
		return nil, nil, err
	}
	e.mu.Lock()
	e.sinceCkpt++
	due := e.sinceCkpt >= e.cfg.EpochRecords
	e.mu.Unlock()
	if due {
		if err := e.Checkpoint(); err != nil {
			e.logger.Printf("checkpoint: %v", err)
		}
	}
	return raw, id, nil
}

// TreeHead is the checkpoint body (8.2).
type TreeHead struct {
	_        struct{} `cbor:",toarray"`
	TrailID  []byte
	TreeSize uint64
	RootHash []byte
}

// Checkpoint closes the current epoch.
func (e *Edge) Checkpoint() error {
	_, _, err := e.trail.AppendCheckpoint(func(trailID []byte, size uint64, root []byte) ([]byte, error) {
		p := &record.Payload{
			Kind:     record.KindCheckpoint,
			TrailRef: record.TrailRef{TrailID: trailID, Leaf: size},
			Body:     TreeHead{TrailID: trailID, TreeSize: size, RootHash: root},
		}
		return e.signer.Sign(e.cfg.Agent, p)
	})
	if err == nil {
		e.mu.Lock()
		e.sinceCkpt = 0
		e.mu.Unlock()
	}
	return err
}

// Lifecycle records an event about the agent, such as start or exit.
func (e *Edge) Lifecycle(body map[string]any) {
	if _, _, err := e.Append(&record.Payload{Kind: record.KindLifecycle, Outcome: record.O(record.Completed), Body: body}); err != nil {
		e.logger.Printf("lifecycle record: %v", err)
	}
}

// ---------------------------------------------------------------- outbound

// OutboundHandler serves the agent-facing port.
func (e *Edge) OutboundHandler() http.Handler { return http.HandlerFunc(e.serveOutbound) }

func (e *Edge) serveOutbound(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == HealthPath && !r.URL.IsAbs() {
		fmt.Fprintln(w, "trail-edge")
		return
	}
	if r.Method == http.MethodConnect {
		// The edge cannot see or sign a tunnel, so it never opens one (4).
		e.deny(w, "CONNECT "+r.Host, "tunnels are not allowed; use a route")
		return
	}
	route, upstream, ok := e.match(r)
	if !ok {
		target := r.URL.Path
		if r.URL.IsAbs() {
			target = r.URL.String()
		}
		e.deny(w, r.Method+" "+target, "no route")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBody+1))
	if err != nil || len(body) > maxBody {
		http.Error(w, "trail-edge: request body too large or unreadable", http.StatusRequestEntityTooLarge)
		return
	}
	out, err := http.NewRequestWithContext(r.Context(), r.Method, upstream.String(), bytes.NewReader(body))
	if err != nil {
		http.Error(w, "trail-edge: "+err.Error(), http.StatusBadGateway)
		return
	}
	copyHeaders(out.Header, r.Header)
	ctxHashes := e.addSecrets(out.Header, route)

	kind, _ := record.ParseKind(route.Kind)
	p := &record.Payload{
		Kind:    kind,
		Outcome: record.O(record.Allowed),
		Aud:     route.Aud,
		Class:   route.Class,
		Context: ctxHashes,
		Body:    callBody(r.Method, upstream, body),
	}

	if kind == record.KindModelCall {
		// Model endpoints don't run TRAIL: record, don't attach (4).
		if _, _, err := e.Append(p); err != nil {
			http.Error(w, "trail-edge: trail append failed", http.StatusInternalServerError)
			return
		}
		e.forward(w, out)
		return
	}

	created := e.now().Unix()
	hasContent := len(body) > 0
	if hasContent {
		out.Header.Set(httpsig.ContentDigest, httpsig.Digest(body))
	}
	msg := httpsig.Message{Request: true, Method: out.Method, TargetURI: upstream.String(), Header: out.Header}
	comps := httpsig.Components(msg, hasContent)
	ch, err := httpsig.ComponentsHash(msg, comps)
	if err != nil {
		http.Error(w, "trail-edge: "+err.Error(), http.StatusInternalServerError)
		return
	}
	p.CallRef = &record.HTTPRef{Type: 2, ComponentsHash: ch, Created: uint64(created)}

	// Commit point: the record is durably in the trail before the call
	// leaves (6.4).
	rec, recID, err := e.Append(p)
	if err != nil {
		http.Error(w, "trail-edge: trail append failed", http.StatusInternalServerError)
		return
	}
	out.Header.Set(httpsig.TrailRecord, httpsig.EncodeRecord(rec))
	if err := httpsig.Sign(msg, comps, created, e.key, e.keyID); err != nil {
		http.Error(w, "trail-edge: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := e.client.Do(out)
	if err != nil {
		e.recordResult(p, recID, nil, record.Failed, map[string]any{"error": err.Error()})
		http.Error(w, "trail-edge: upstream: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil || len(respBody) > maxBody {
		e.recordResult(p, recID, nil, record.Failed, map[string]any{"error": "response too large or unreadable"})
		http.Error(w, "trail-edge: upstream response too large or unreadable", http.StatusBadGateway)
		return
	}

	if route.Aud == "" {
		// The counterparty doesn't run TRAIL: there is no record to check.
		e.recordResult(p, recID, nil, record.Unsigned, map[string]any{"status": resp.StatusCode, "digest": httpsig.Digest(respBody)})
		writeResponse(w, resp, respBody)
		return
	}

	peerRec, err := e.verifyResponse(resp, respBody, route.Aud, recID)
	if err != nil {
		e.recordResult(p, recID, nil, record.Failed, map[string]any{"status": resp.StatusCode, "error": err.Error()})
		http.Error(w, "trail-edge: response refused: "+err.Error(), http.StatusBadGateway)
		return
	}
	e.recordResult(p, recID, peerRec, record.Completed, map[string]any{"status": resp.StatusCode})
	writeResponse(w, resp, respBody)
}

// match finds the route for a request: by path prefix on the edge, or, for
// a forward-proxy request, by upstream prefix.
func (e *Edge) match(r *http.Request) (config.Route, *url.URL, bool) {
	if r.URL.IsAbs() {
		// A client that honors the proxy variables sends even calls for
		// the edge's own routes in proxy form.
		if la, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok && r.URL.Host == la.String() {
			local := *r
			u := *r.URL
			u.Scheme, u.Host = "", ""
			local.URL = &u
			return e.match(&local)
		}
		full := r.URL.String()
		for _, rt := range e.routes {
			if strings.HasPrefix(full, rt.Upstream) {
				u, err := url.Parse(full)
				return rt, u, err == nil
			}
		}
		return config.Route{}, nil, false
	}
	for _, rt := range e.routes {
		if r.URL.Path == rt.Prefix || strings.HasPrefix(r.URL.Path, strings.TrimSuffix(rt.Prefix, "/")+"/") {
			u, err := url.Parse(strings.TrimSuffix(rt.Upstream, "/") + strings.TrimPrefix(r.URL.Path, rt.Prefix))
			if err != nil {
				return rt, nil, false
			}
			u.RawQuery = r.URL.RawQuery
			return rt, u, true
		}
	}
	return config.Route{}, nil, false
}

// deny refuses a call that has no route, and records the refusal.
func (e *Edge) deny(w http.ResponseWriter, target, reason string) {
	_, _, err := e.Append(&record.Payload{
		Kind:    record.KindSend,
		Outcome: record.O(record.Denied),
		Body:    map[string]any{"target": target, "reason": reason},
	})
	if err != nil {
		e.logger.Printf("deny record: %v", err)
	}
	e.logger.Printf("denied %s: %s", target, reason)
	http.Error(w, "trail-edge: denied: "+reason, http.StatusForbidden)
}

// addSecrets adds a route's credentials from the edge's own environment
// and returns their hashes for context-hashes (3.2).
func (e *Edge) addSecrets(h http.Header, rt config.Route) *record.ContextHashes {
	if len(rt.Headers) == 0 {
		return nil
	}
	oth := map[string][]byte{}
	for name, s := range rt.Headers {
		v := lookupEnv(s.Env)
		if v == "" {
			e.logger.Printf("route %s: %s is empty; header %s not added", rt.Prefix, s.Env, name)
			continue
		}
		if s.Format != "" {
			v = fmt.Sprintf(s.Format, v)
		}
		h.Set(name, v)
		sum := sha256.Sum256([]byte(v))
		oth[strings.ToLower(name)] = sum[:]
	}
	if len(oth) == 0 {
		return nil
	}
	return &record.ContextHashes{OTH: oth}
}

// recordResult appends the record that closes an outbound call: a recv
// that cites the responder's record when there is one.
func (e *Edge) recordResult(sent *record.Payload, sentID []byte, peerRec *record.Decoded, o record.Outcome, body map[string]any) {
	body["request"] = map[string]any{"leaf": sent.TrailRef.Leaf, "id": sentID}
	p := &record.Payload{Kind: record.KindRecv, Outcome: record.O(o), Aud: sent.Aud, Class: sent.Class, Body: body}
	if peerRec != nil {
		p.Cites = []record.Cite{{Type: record.CitePeerRecord, Ref: peerRec.Payload.TrailRef, Hash: peerRec.ID}}
	}
	if _, _, err := e.Append(p); err != nil {
		e.logger.Printf("result record: %v", err)
	}
}

// verifyResponse checks the responder's record: signature, binding to the
// response, issuer, and that it cites our request's record (4).
func (e *Edge) verifyResponse(resp *http.Response, body []byte, aud string, reqID []byte) (*record.Decoded, error) {
	d, err := e.verifyMessage(httpsig.Message{Status: resp.StatusCode, Header: resp.Header}, body)
	if err != nil {
		return nil, err
	}
	if d.Iss != aud {
		return nil, fmt.Errorf("record issued by %q, expected %q", d.Iss, aud)
	}
	for _, c := range d.Payload.Cites {
		if c.Type == record.CitePeerRecord && record.Equal(c.Hash, reqID) {
			return d, nil
		}
	}
	return nil, errors.New("response record does not cite the request")
}

// verifyMessage checks a message's Trail-Record and its signature, and
// that the record's call-ref matches the message as received.
func (e *Edge) verifyMessage(m httpsig.Message, body []byte) (*record.Decoded, error) {
	v := m.Header.Get(httpsig.TrailRecord)
	if v == "" {
		return nil, errors.New("no Trail-Record")
	}
	raw, err := httpsig.DecodeRecord(v)
	if err != nil {
		return nil, err
	}
	var pr peer
	d, err := record.Verify(raw, func(kid []byte) (*ecdsa.PublicKey, error) {
		p, ok := e.peers[base64.RawURLEncoding.EncodeToString(kid)]
		if !ok {
			return nil, errors.New("record signed by an unknown key")
		}
		pr = p
		return p.pub, nil
	})
	if err != nil {
		return nil, err
	}
	if d.Iss != pr.id {
		return nil, fmt.Errorf("record iss %q does not match key holder %q", d.Iss, pr.id)
	}
	params, err := httpsig.ParseInput(m.Header)
	if err != nil {
		return nil, err
	}
	if params.KeyID != base64.RawURLEncoding.EncodeToString(d.Kid) {
		return nil, errors.New("signature keyid differs from record kid")
	}
	hasContent := len(body) > 0
	if hasContent && m.Header.Get(httpsig.ContentDigest) != httpsig.Digest(body) {
		return nil, errors.New("Content-Digest does not match the content")
	}
	if err := httpsig.Verify(m, params, pr.pub, hasContent); err != nil {
		return nil, err
	}
	cr := d.Payload.CallRef
	if cr == nil || cr.Type != 2 {
		return nil, errors.New("record has no http call-ref")
	}
	ch, err := httpsig.ComponentsHash(m, params.Components)
	if err != nil {
		return nil, err
	}
	if !record.Equal(ch, cr.ComponentsHash) || int64(cr.Created) != params.Created {
		return nil, errors.New("record call-ref does not match this message")
	}
	return d, nil
}

// ----------------------------------------------------------------- inbound

// InboundHandler accepts calls from other edges for the agent.
func (e *Edge) InboundHandler() http.Handler { return http.HandlerFunc(e.serveInbound) }

func (e *Edge) serveInbound(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBody+1))
	if err != nil || len(body) > maxBody {
		http.Error(w, "trail-edge: request body too large or unreadable", http.StatusRequestEntityTooLarge)
		return
	}
	target := "http://" + r.Host + r.URL.RequestURI()
	msg := httpsig.Message{Request: true, Method: r.Method, TargetURI: target, Header: r.Header}

	refuse := func(reason string) {
		e.Append(&record.Payload{
			Kind:    record.KindRecv,
			Outcome: record.O(record.Denied),
			Body:    map[string]any{"target": r.Method + " " + r.URL.RequestURI(), "reason": reason},
		})
		e.logger.Printf("refused inbound %s %s: %s", r.Method, r.URL.RequestURI(), reason)
		http.Error(w, "trail-edge: refused: "+reason, http.StatusForbidden)
	}

	// Steps 1 to 3: a record must be present, signed by a known peer, and
	// bound to this call.
	d, err := e.verifyMessage(msg, body)
	if err != nil {
		refuse(err.Error())
		return
	}
	if d.Payload.Aud != "" && d.Payload.Aud != e.cfg.Identity {
		refuse(fmt.Sprintf("record is for %q", d.Payload.Aud))
		return
	}
	// Step 4 (level, holds, admission, freshness) needs state: not yet.
	// Step 5: replay.
	created := time.Unix(int64(d.Payload.CallRef.Created), 0)
	skew := time.Duration(e.cfg.MaxSkewSeconds) * time.Second
	if now := e.now(); created.Before(now.Add(-skew)) || created.After(now.Add(skew)) {
		refuse("created is outside the allowed clock skew")
		return
	}
	if !e.markSeen(d.ID, created, skew) {
		refuse("replayed record")
		return
	}

	// Step 6: record the receipt, citing the sender's record, then deliver.
	cite := []record.Cite{{Type: record.CitePeerRecord, Ref: d.Payload.TrailRef, Hash: d.ID}}
	inClass := d.Payload.Class
	if _, _, err := e.Append(&record.Payload{
		Kind: record.KindRecv, Outcome: record.O(record.Allowed), Aud: d.Iss, Class: inClass, Cites: cite,
		Body: callBody(r.Method, r.URL, body),
	}); err != nil {
		http.Error(w, "trail-edge: trail append failed", http.StatusInternalServerError)
		return
	}

	agentURL := "http://" + e.cfg.Inbound.Agent + r.URL.RequestURI()
	in, err := http.NewRequestWithContext(r.Context(), r.Method, agentURL, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "trail-edge: "+err.Error(), http.StatusBadGateway)
		return
	}
	copyHeaders(in.Header, r.Header)
	resp, err := e.client.Do(in)
	status, respBody, respHeader := http.StatusBadGateway, []byte("trail-edge: agent unreachable\n"), http.Header{}
	outcome := record.Failed
	if err == nil {
		defer resp.Body.Close()
		respBody, err = io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
		if err == nil && len(respBody) <= maxBody {
			status, outcome = resp.StatusCode, record.Completed
			copyHeaders(respHeader, resp.Header)
		} else {
			respBody = []byte("trail-edge: agent response too large or unreadable\n")
		}
	}

	// The response carries our record, citing the request's (4).
	created2 := e.now().Unix()
	hasContent := len(respBody) > 0
	if hasContent {
		respHeader.Set(httpsig.ContentDigest, httpsig.Digest(respBody))
	}
	rmsg := httpsig.Message{Status: status, Header: respHeader}
	comps := httpsig.Components(rmsg, hasContent)
	ch, err := httpsig.ComponentsHash(rmsg, comps)
	if err != nil {
		http.Error(w, "trail-edge: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rec, _, err := e.Append(&record.Payload{
		Kind: record.KindSend, Outcome: record.O(outcome), Aud: d.Iss, Class: inClass, Cites: cite,
		CallRef: &record.HTTPRef{Type: 2, ComponentsHash: ch, Created: uint64(created2)},
		Body:    map[string]any{"status": status},
	})
	if err != nil {
		http.Error(w, "trail-edge: trail append failed", http.StatusInternalServerError)
		return
	}
	respHeader.Set(httpsig.TrailRecord, httpsig.EncodeRecord(rec))
	if err := httpsig.Sign(rmsg, comps, created2, e.key, e.keyID); err != nil {
		http.Error(w, "trail-edge: "+err.Error(), http.StatusInternalServerError)
		return
	}
	for k, vs := range respHeader {
		w.Header()[k] = vs
	}
	w.WriteHeader(status)
	w.Write(respBody)
}

// markSeen remembers an accepted record id for the replay window and
// reports whether it is new.
func (e *Edge) markSeen(id []byte, created time.Time, window time.Duration) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	k := string(id)
	if _, ok := e.seen[k]; ok {
		return false
	}
	cutoff := e.now().Add(-2 * window)
	for s, t := range e.seen {
		if t.Before(cutoff) {
			delete(e.seen, s)
		}
	}
	e.seen[k] = created
	return true
}

// ----------------------------------------------------------------- helpers

var hopHeaders = map[string]bool{
	"Connection": true, "Proxy-Connection": true, "Keep-Alive": true, "Proxy-Authenticate": true,
	"Proxy-Authorization": true, "Te": true, "Trailer": true, "Transfer-Encoding": true, "Upgrade": true,
	"Host": true, "Content-Length": true,
	// The edge sets these itself; whatever the agent or upstream sent is dropped.
	"Trail-Record": true, "Signature": true, "Signature-Input": true, "Content-Digest": true,
}

func copyHeaders(dst, src http.Header) {
	for k, vs := range src {
		if hopHeaders[http.CanonicalHeaderKey(k)] {
			continue
		}
		dst[k] = append([]string(nil), vs...)
	}
}

func (e *Edge) forward(w http.ResponseWriter, out *http.Request) {
	resp, err := e.client.Do(out)
	if err != nil {
		http.Error(w, "trail-edge: upstream: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	h := http.Header{}
	copyHeaders(h, resp.Header)
	for k, vs := range h {
		w.Header()[k] = vs
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func writeResponse(w http.ResponseWriter, resp *http.Response, body []byte) {
	h := http.Header{}
	copyHeaders(h, resp.Header)
	for k, vs := range h {
		w.Header()[k] = vs
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// callBody is the kind-specific body of a call record. It records the
// protocol's own identifiers when the call is JSON-RPC (MCP, A2A), as the
// draft recommends (7.1). Arguments are not recorded until SD-CWT
// redaction exists.
func callBody(method string, u *url.URL, body []byte) map[string]any {
	b := map[string]any{"method": method, "path": u.Path, "content-digest": httpsig.Digest(body)}
	if u.Host != "" {
		b["host"] = u.Host
	}
	var rpc struct {
		JSONRPC string          `json:"jsonrpc"`
		Method  string          `json:"method"`
		ID      json.RawMessage `json:"id"`
		Params  struct {
			Name string `json:"name"`
		} `json:"params"`
	}
	if json.Unmarshal(body, &rpc) == nil && rpc.JSONRPC != "" {
		b["jsonrpc-method"] = rpc.Method
		if len(rpc.ID) > 0 {
			b["jsonrpc-id"] = string(rpc.ID)
		}
		if rpc.Method == "tools/call" && rpc.Params.Name != "" {
			b["tool"] = rpc.Params.Name
		}
	}
	return b
}

// Serve runs a handler on a listener until ctx ends.
func Serve(ctx context.Context, ln net.Listener, h http.Handler) error {
	srv := &http.Server{Handler: h, ReadHeaderTimeout: 30 * time.Second}
	go func() {
		<-ctx.Done()
		sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(sctx)
	}()
	err := srv.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
