// Command trail-edge is a TRAIL edge for agents in nono's OS sandbox.
// See ../../README.md.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/kareemn/trail/edge/internal/config"
	"github.com/kareemn/trail/edge/internal/edge"
	"github.com/kareemn/trail/edge/internal/keys"
	"github.com/kareemn/trail/edge/internal/sandbox"
)

const version = "0.1.0-dev"

func usage() {
	fmt.Fprintf(os.Stderr, `trail-edge %s

Usage:
  trail-edge run    -config edge.json -- <agent command> [args]
      Start the edge and run the agent in nono's OS sandbox. The agent's
      only way out is the edge.
  trail-edge serve  -config edge.json
      Run the edge without starting an agent (for an agent or tool that
      runs elsewhere and listens on inbound.agent).
  trail-edge key    -config edge.json
      Print this edge's peer entry (identifier and public key) for other
      edges' configs.
  trail-edge verify -config edge.json
      Check every record, trail position, and checkpoint in the trail.
  trail-edge version
`, version)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case sandbox.ExecCommand:
		err = sandbox.Exec()
		fmt.Fprintln(os.Stderr, "trail-edge sandbox:", err)
		os.Exit(126)
	case "run":
		var code int
		code, err = run(os.Args[2:])
		if err == nil {
			os.Exit(code)
		}
	case "serve":
		err = serve(os.Args[2:])
	case "key":
		err = key(os.Args[2:])
	case "verify":
		err = verify(os.Args[2:])
	case "version":
		fmt.Printf("trail-edge %s (nono %s)\n", version, sandbox.Version())
	case "-h", "-help", "--help", "help":
		usage()
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "trail-edge:", err)
		os.Exit(1)
	}
}

func loadConfig(name string, args []string) (*config.Config, []string, error) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	path := fs.String("config", "edge.json", "edge configuration file")
	fs.Parse(args)
	cfg, err := config.Load(*path)
	if err != nil {
		return nil, nil, err
	}
	if err := os.MkdirAll(cfg.StateDir, 0o700); err != nil {
		return nil, nil, err
	}
	return cfg, fs.Args(), nil
}

func newLogger(cfg *config.Config) *log.Logger {
	return log.New(os.Stderr, "trail-edge["+cfg.Identity+"] ", log.Ltime)
}

// start opens the edge and its listeners.
func start(ctx context.Context, cfg *config.Config, logger *log.Logger) (*edge.Edge, net.Listener, chan error, error) {
	e, err := edge.New(cfg, logger)
	if err != nil {
		return nil, nil, nil, err
	}
	errc := make(chan error, 2)
	out, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		e.Close()
		return nil, nil, nil, err
	}
	go func() { errc <- edge.Serve(ctx, out, e.OutboundHandler()) }()
	if cfg.Inbound != nil {
		in, err := net.Listen("tcp", cfg.Inbound.Listen)
		if err != nil {
			out.Close()
			e.Close()
			return nil, nil, nil, err
		}
		logger.Printf("inbound on %s, delivering to %s", in.Addr(), cfg.Inbound.Agent)
		go func() { errc <- edge.Serve(ctx, in, e.InboundHandler()) }()
	}
	logger.Printf("agent-facing port %s, key %s", out.Addr(), e.KeyID())
	return e, out, errc, nil
}

func run(args []string) (int, error) {
	cfg, argv, err := loadConfig("run", args)
	if err != nil {
		return 0, err
	}
	if len(argv) == 0 {
		return 0, errors.New("run: no agent command; use: trail-edge run -config edge.json -- <command>")
	}
	logger := newLogger(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	e, out, _, err := start(ctx, cfg, logger)
	if err != nil {
		return 0, err
	}
	defer e.Close()

	addr := out.Addr().(*net.TCPAddr)
	edgeURL := "http://" + addr.String()
	wd, err := os.Getwd()
	if err != nil {
		return 0, err
	}
	spec := &sandbox.Spec{
		EdgePort:  uint16(addr.Port),
		Read:      append(sandbox.SystemRead(), cfg.Sandbox.Read...),
		Write:     cfg.Sandbox.Write,
		ReadFile:  cfg.Sandbox.ReadFile,
		Argv:      argv,
		Dir:       wd,
		MustBlock: []string{cfg.StateDir},
	}
	if cfg.Sandbox.Workdir {
		spec.Write = append(spec.Write, wd)
	}
	summary, err := sandbox.Summary(spec)
	if err != nil {
		return 0, err
	}
	cmd, err := sandbox.Command(spec, agentEnv(cfg, edgeURL, addr.String()))
	if err != nil {
		return 0, err
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	sum := sha256.Sum256([]byte(summary))
	e.Lifecycle(map[string]any{
		"event":   "start",
		"argv":    argv,
		"sandbox": "nono " + sandbox.Version(),
		// The sandbox's rules, by hash; the text is in the edge's log.
		"sandbox-summary-sha256": sum[:],
		"edge-port":              addr.Port,
	})
	logger.Printf("sandbox:\n%s", summary)
	if err := cmd.Start(); err != nil {
		e.Lifecycle(map[string]any{"event": "start-failed", "error": err.Error()})
		return 0, err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for s := range sigs {
			cmd.Process.Signal(s)
		}
	}()
	err = cmd.Wait()
	signal.Stop(sigs)
	code := 0
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code = exitErr.ExitCode()
	} else if err != nil {
		code = 1
	}
	e.Lifecycle(map[string]any{"event": "exit", "exit_code": code})
	return code, nil
}

// agentEnv is the agent's whole environment: a few ordinary variables,
// what the config passes through, and where each route is. Nothing else
// from the edge's environment, so the agent never sees a credential.
func agentEnv(cfg *config.Config, edgeURL, edgeAddr string) []string {
	pass := []string{"PATH", "HOME", "USER", "LOGNAME", "SHELL", "TERM", "LANG", "LC_ALL", "TMPDIR", "TZ"}
	pass = append(pass, cfg.Sandbox.Env...)
	var env []string
	for _, k := range pass {
		if v, ok := os.LookupEnv(k); ok {
			env = append(env, k+"="+v)
		}
	}
	for _, r := range cfg.Routes {
		if r.Env != "" {
			env = append(env, r.Env+"="+edgeURL+r.Prefix)
		}
	}
	// Clients that honor proxy variables reach the edge too. A CONNECT
	// through it is refused and recorded.
	for _, k := range []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy", "ALL_PROXY"} {
		env = append(env, k+"="+edgeURL)
	}
	env = append(env, "NO_PROXY=", "no_proxy=", "TRAIL_EDGE_URL="+edgeURL, "TRAIL_EDGE_ADDR="+edgeAddr,
		"TRAIL_EDGE_PID="+strconv.Itoa(os.Getpid()))
	return env
}

func serve(args []string) error {
	cfg, _, err := loadConfig("serve", args)
	if err != nil {
		return err
	}
	logger := newLogger(cfg)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	e, _, errc, err := start(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer e.Close()
	select {
	case <-ctx.Done():
		return nil
	case err := <-errc:
		return err
	}
}

func key(args []string) error {
	cfg, _, err := loadConfig("key", args)
	if err != nil {
		return err
	}
	k, err := keys.LoadOrCreate(cfg.StateDir)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(config.Peer{ID: cfg.Identity, JWK: keys.PublicJWK(&k.PublicKey)}, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func verify(args []string) error {
	cfg, _, err := loadConfig("verify", args)
	if err != nil {
		return err
	}
	k, err := keys.LoadOrCreate(cfg.StateDir)
	if err != nil {
		return err
	}
	if err := edge.VerifyTrail(cfg.StateDir, &k.PublicKey, os.Stdout); err != nil {
		return err
	}
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("trail verifies")
	return nil
}
