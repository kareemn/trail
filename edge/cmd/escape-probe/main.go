// Command escape-probe checks the Phase 0 question: does nono leave the
// agent exactly one way out, the edge's localhost port?
//
// Run it inside the sandbox. It tries every way out it knows and reports
// each one. The edge port must be reachable; everything else must fail.
// It exits 1 if any result differs from what is expected.
//
// With -serve it instead acts as a stand-in edge: it listens on the given
// address outside the sandbox and answers the edge's health path.
//
// Under "trail-edge run", -edge and -edge-pid default to the TRAIL_EDGE_ADDR
// and TRAIL_EDGE_PID variables the edge gives the agent.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// healthPath matches edge.HealthPath.
const healthPath = "/.well-known/trail-edge"

type check struct {
	name   string
	wantOK bool
	run    func() error
}

func main() {
	serve := flag.String("serve", "", "act as the stand-in edge on this address, e.g. 127.0.0.1:7777")
	edge := flag.String("edge", envOr("TRAIL_EDGE_ADDR", "127.0.0.1:7777"), "the edge's address; must be reachable")
	other := flag.String("other-local", "127.0.0.1:7778", "a local port that is not the edge; must not be reachable")
	secret := flag.String("secret", "", "a file the edge owns; must not be readable")
	pidDefault, _ := strconv.Atoi(os.Getenv("TRAIL_EDGE_PID"))
	edgePID := flag.Int("edge-pid", pidDefault, "the edge's process id; signalling it must fail")
	flag.Parse()

	if *serve != "" {
		os.Exit(runServer(*serve))
	}

	timeout := 3 * time.Second
	dial := func(network, addr string) func() error {
		return func() error {
			c, err := net.DialTimeout(network, addr, timeout)
			if err != nil {
				return err
			}
			defer c.Close()
			if network == "udp" {
				// UDP dial never fails; a write and a read show whether
				// packets really leave.
				_ = c.SetDeadline(time.Now().Add(timeout))
				if _, err := c.Write(dnsQuery()); err != nil {
					return err
				}
				buf := make([]byte, 512)
				_, err = c.Read(buf)
				return err
			}
			return nil
		}
	}

	checks := []check{
		{"http to the edge " + *edge, true, func() error {
			c := &http.Client{Timeout: timeout}
			resp, err := c.Get("http://" + *edge + healthPath)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			b, _ := io.ReadAll(resp.Body)
			if strings.TrimSpace(string(b)) != "trail-edge" {
				return fmt.Errorf("unexpected answer %q", b)
			}
			return nil
		}},
		{"tcp to another local port " + *other, false, dial("tcp", *other)},
		{"tcp to 1.1.1.1:443", false, dial("tcp", "1.1.1.1:443")},
		{"tcp to [2606:4700:4700::1111]:443", false, dial("tcp", "[2606:4700:4700::1111]:443")},
		{"udp dns to 8.8.8.8:53", false, dial("udp", "8.8.8.8:53")},
		{"dns lookup of example.com", false, func() error {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			_, err := net.DefaultResolver.LookupHost(ctx, "example.com")
			return err
		}},
		{"tcp to metadata 169.254.169.254:80", false, dial("tcp", "169.254.169.254:80")},
	}
	if *secret != "" {
		checks = append(checks, check{"read edge secret " + *secret, false, func() error {
			_, err := os.ReadFile(*secret)
			return err
		}})
	}
	if *edgePID > 0 {
		checks = append(checks, check{fmt.Sprintf("signal 0 to edge pid %d", *edgePID), false, func() error {
			return syscall.Kill(*edgePID, 0)
		}})
	}

	failed := 0
	for _, c := range checks {
		err := c.run()
		ok := err == nil
		verdict := "PASS"
		if ok != c.wantOK {
			verdict = "FAIL"
			failed++
		}
		want := "blocked"
		if c.wantOK {
			want = "reachable"
		}
		detail := "reached"
		if err != nil {
			detail = err.Error()
		}
		fmt.Printf("%s  want %-9s  %s: %s\n", verdict, want, c.name, detail)
	}
	if failed > 0 {
		fmt.Printf("%d check(s) failed\n", failed)
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}

func runServer(addr string) int {
	fmt.Fprintf(os.Stderr, "stand-in edge listening on %s (pid %d)\n", addr, os.Getpid())
	err := http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "trail-edge")
	}))
	fmt.Fprintln(os.Stderr, err)
	return 1
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// dnsQuery is a minimal DNS query for example.com, type A.
func dnsQuery() []byte {
	q := []byte{0x12, 0x34, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for _, label := range []string{"example", "com"} {
		q = append(q, byte(len(label)))
		q = append(q, label...)
	}
	return append(q, 0x00, 0x00, 0x01, 0x00, 0x01)
}
