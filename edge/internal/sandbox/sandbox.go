// Package sandbox starts the agent inside nono's OS sandbox (Landlock on
// Linux, Seatbelt on macOS) with the edge's port as its only way out.
//
// nono's Apply sandboxes the calling process and everything it starts, and
// cannot be undone. The edge must stay outside, so it starts a copy of its
// own binary (the "__sandbox-exec" command), which applies the sandbox to
// itself and then execs the agent. The edge process is never sandboxed and
// the agent never runs unsandboxed.
package sandbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	nono "github.com/always-further/nono-go"
)

// SpecEnv carries the Spec from the edge to the __sandbox-exec child. The
// child removes it before it execs the agent.
const SpecEnv = "TRAIL_SANDBOX_SPEC"

// ExecCommand is the hidden subcommand that applies the sandbox.
const ExecCommand = "__sandbox-exec"

// Spec is what the child needs to sandbox itself.
type Spec struct {
	EdgePort  uint16   `json:"edge_port"`
	Read      []string `json:"read,omitempty"`
	Write     []string `json:"write,omitempty"`
	ReadFile  []string `json:"read_file,omitempty"`
	Argv      []string `json:"argv"`
	Dir       string   `json:"dir"`
	MustBlock []string `json:"must_block,omitempty"` // paths the sandbox must not cover
}

// Supported reports whether nono can sandbox on this platform.
func Supported() bool { return nono.IsSupported() }

// Version is the nono library version.
func Version() string { return nono.Version() }

// SystemRead returns the read-only system paths an ordinary program needs
// to start on this OS. They hold no secrets of the edge.
func SystemRead() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{"/usr", "/bin", "/sbin", "/System", "/Library", "/private/etc", "/private/var/db/timezone", "/dev", "/opt/homebrew"}
	default:
		return []string{"/usr", "/bin", "/sbin", "/lib", "/lib64", "/etc", "/dev", "/proc"}
	}
}

func build(s *Spec) (*nono.CapabilitySet, error) {
	caps := nono.New()
	add := func(paths []string, mode nono.AccessMode, file bool) error {
		for _, p := range paths {
			if _, err := os.Stat(p); err != nil {
				continue // nono refuses missing paths; skip them
			}
			var err error
			if file {
				err = caps.AllowFile(p, mode)
			} else {
				err = caps.AllowPath(p, mode)
			}
			if err != nil {
				return fmt.Errorf("allow %s: %w", p, err)
			}
		}
		return nil
	}
	if err := add(s.Read, nono.AccessRead, false); err != nil {
		return nil, err
	}
	if err := add(s.Write, nono.AccessReadWrite, false); err != nil {
		return nil, err
	}
	if err := add(s.ReadFile, nono.AccessRead, true); err != nil {
		return nil, err
	}
	if err := caps.SetNetworkMode(nono.NetworkProxyOnly); err != nil {
		return nil, err
	}
	if err := caps.SetProxyPort(s.EdgePort); err != nil {
		return nil, err
	}
	for _, p := range s.MustBlock {
		covered, err := caps.PathCovered(p)
		if err != nil {
			return nil, err
		}
		if covered {
			return nil, fmt.Errorf("sandbox would allow %s, which holds the edge's key and trail", p)
		}
	}
	return caps, nil
}

// Summary describes the sandbox a Spec produces, for the lifecycle record.
func Summary(s *Spec) (string, error) {
	caps, err := build(s)
	if err != nil {
		return "", err
	}
	defer caps.Close()
	return caps.Summary(), nil
}

// Command returns the command that runs the agent in the sandbox. The
// caller sets Stdin, Stdout, Stderr and starts it.
func Command(s *Spec, env []string) (*exec.Cmd, error) {
	if !Supported() {
		return nil, errors.New("nono sandboxing is not supported on this platform")
	}
	if _, err := Summary(s); err != nil { // fail early, in the edge
		return nil, err
	}
	self, err := os.Executable()
	if err != nil {
		return nil, err
	}
	spec, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(self, ExecCommand)
	cmd.Dir = s.Dir
	cmd.Env = append(env, SpecEnv+"="+string(spec))
	return cmd, nil
}

// Exec is the __sandbox-exec command: apply the sandbox, then exec the
// agent. It returns only on failure.
func Exec() error {
	// Landlock applies to the calling thread, and execve keeps the calling
	// thread's sandbox, so apply and exec on one locked thread.
	runtime.LockOSThread()
	var s Spec
	if err := json.Unmarshal([]byte(os.Getenv(SpecEnv)), &s); err != nil {
		return fmt.Errorf("%s: %w", SpecEnv, err)
	}
	if len(s.Argv) == 0 {
		return errors.New("no agent command")
	}
	// Resolve the program before the sandbox narrows what can be read.
	path, err := exec.LookPath(s.Argv[0])
	if err != nil {
		return err
	}
	if path, err = filepath.Abs(path); err != nil {
		return err
	}
	caps, err := build(&s)
	if err != nil {
		return err
	}
	if err := nono.Apply(caps); err != nil {
		return fmt.Errorf("apply sandbox: %w", err)
	}
	env := make([]string, 0, len(os.Environ()))
	for _, kv := range os.Environ() {
		if len(kv) > len(SpecEnv) && kv[:len(SpecEnv)+1] == SpecEnv+"=" {
			continue
		}
		env = append(env, kv)
	}
	return syscall.Exec(path, s.Argv, env)
}
