// Command trail-edge is a TRAIL edge for agents that run in a nono
// sandbox. See ../../README.md for the plan. Nothing is implemented yet.
package main

import (
	"fmt"
	"os"
)

const version = "0.0.0-dev"

func usage() {
	fmt.Fprintf(os.Stderr, `trail-edge %s

Usage:
  trail-edge version
  trail-edge run [flags] -- <agent command>   (not implemented yet)
`, version)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "version":
		fmt.Println(version)
	case "run":
		fmt.Fprintln(os.Stderr, "trail-edge run: not implemented yet (Phase 1)")
		os.Exit(2)
	default:
		usage()
		os.Exit(2)
	}
}
