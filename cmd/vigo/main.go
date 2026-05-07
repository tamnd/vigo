// Command vigo launches the Vigo terminal IDE.
//
// vigo is a 100% faithful, modernized port of Borland's Turbo Vision IDE,
// reimagined as a self-hosted Go IDE written in Go. See README.md and the
// docs/ tree for the project specification.
package main

import (
	"flag"
	"fmt"
	"os"
)

// Version is the build version of the vigo binary. It is overridden at link
// time by the release workflow via -ldflags.
var Version = "0.0.1-dev"

func main() {
	flagSet := flag.NewFlagSet("vigo", flag.ContinueOnError)
	showVersion := flagSet.Bool("version", false, "print version and exit")
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}
	if *showVersion {
		fmt.Printf("vigo %s\n", Version)
		return
	}

	// The real TUI is wired up in the v0.1 foundation PR.
	fmt.Printf("vigo %s — bootstrap binary; v0.1 foundation lands soon.\n", Version)
}
