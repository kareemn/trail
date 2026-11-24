// Package config reads an edge's JSON configuration.
//
// Until state snapshots are implemented (draft 9.1), peers' keys come from
// this file instead of from their published state. That is a stand-in, not
// part of TRAIL.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kareemn/trail/edge/internal/keys"
)

// Config is one edge.
type Config struct {
	// Identity is the edge's WIMSE identifier, the record's iss.
	Identity string `json:"identity"`
	// Agent is the identifier of the agent behind this edge, the record's sub.
	Agent string `json:"agent"`
	// StateDir holds the edge's key and trail. The sandbox never allows it.
	StateDir string `json:"state_dir"`
	// Listen is where the agent reaches the edge, normally 127.0.0.1:0.
	Listen string `json:"listen"`
	// Inbound, if set, accepts calls for the agent from other edges.
	Inbound *Inbound `json:"inbound,omitempty"`
	// Routes are the only destinations the agent can reach.
	Routes []Route `json:"routes"`
	// Peers are other edges whose records this edge accepts.
	Peers []Peer `json:"peers,omitempty"`
	// Sandbox is what the agent may touch on the host.
	Sandbox Sandbox `json:"sandbox"`
	// EpochRecords is how many records an epoch holds before a checkpoint.
	EpochRecords int `json:"epoch_records,omitempty"`
	// MaxSkewSeconds bounds the created time of inbound calls (14.4).
	MaxSkewSeconds int `json:"max_skew_seconds,omitempty"`
}

// Inbound accepts calls from other edges and delivers them to the agent.
type Inbound struct {
	Listen string `json:"listen"`
	// Agent is where the agent itself listens, inside the sandbox.
	Agent string `json:"agent"`
}

// Route maps a path prefix on the edge to an upstream.
type Route struct {
	// Prefix is the path the agent uses, such as /mcp/github.
	Prefix string `json:"prefix"`
	// Upstream is the real endpoint, such as https://api.githubcopilot.com/mcp.
	Upstream string `json:"upstream"`
	// Kind is send, tool-call, or model-call.
	Kind string `json:"kind"`
	// Class is the deployment-defined action class.
	Class *uint `json:"class,omitempty"`
	// Aud is the counterparty's identifier. Set it when the counterparty
	// runs TRAIL; its edge must then answer with a signed record.
	Aud string `json:"aud,omitempty"`
	// Env is the variable the agent reads to find this route.
	Env string `json:"env,omitempty"`
	// Headers are credentials the edge adds. Values come from the edge's
	// own environment, never the agent's.
	Headers map[string]Secret `json:"headers,omitempty"`
}

// Secret is a header value taken from the edge's environment.
type Secret struct {
	Env    string `json:"env"`
	Format string `json:"format,omitempty"` // e.g. "Bearer %s"
}

// Peer is another edge's identifier and public key.
type Peer struct {
	ID  string   `json:"id"`
	JWK keys.JWK `json:"jwk"`
}

// Sandbox lists what the agent may access.
type Sandbox struct {
	Read     []string `json:"read,omitempty"`
	Write    []string `json:"write,omitempty"`
	ReadFile []string `json:"read_file,omitempty"`
	// Workdir grants read-write access to the directory the agent starts in.
	Workdir bool `json:"workdir,omitempty"`
	// Env lists variables from the edge's environment passed to the agent.
	Env []string `json:"env,omitempty"`
}

// Load reads a configuration file and fills defaults.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if c.Identity == "" || c.Agent == "" || c.StateDir == "" {
		return nil, fmt.Errorf("%s: identity, agent, and state_dir are required", path)
	}
	c.StateDir = expand(c.StateDir, filepath.Dir(path))
	for i := range c.Sandbox.Read {
		c.Sandbox.Read[i] = expand(c.Sandbox.Read[i], filepath.Dir(path))
	}
	for i := range c.Sandbox.Write {
		c.Sandbox.Write[i] = expand(c.Sandbox.Write[i], filepath.Dir(path))
	}
	for i := range c.Sandbox.ReadFile {
		c.Sandbox.ReadFile[i] = expand(c.Sandbox.ReadFile[i], filepath.Dir(path))
	}
	if c.Listen == "" {
		c.Listen = "127.0.0.1:0"
	}
	if c.EpochRecords == 0 {
		c.EpochRecords = 64
	}
	if c.MaxSkewSeconds == 0 {
		c.MaxSkewSeconds = 300
	}
	for _, r := range c.Routes {
		if !strings.HasPrefix(r.Prefix, "/") || r.Upstream == "" {
			return nil, fmt.Errorf("%s: route %q needs a /prefix and an upstream", path, r.Prefix)
		}
		switch r.Kind {
		case "send", "tool-call", "model-call":
		default:
			return nil, fmt.Errorf("%s: route %q: kind must be send, tool-call, or model-call", path, r.Prefix)
		}
	}
	return &c, nil
}

// expand resolves ~ and paths relative to the config file.
func expand(p, base string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			p = filepath.Join(home, p[2:])
		}
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(base, p)
	}
	return filepath.Clean(p)
}
