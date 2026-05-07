// Command vigo launches the Vigo terminal IDE.
//
// vigo is a faithful Go port of Borland's Turbo Vision IDE, with the goal of
// hosting a self-built Go editor and debugger inside a recognizably classic
// blue desktop. See README.md and the docs/ tree for the specifications.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tamnd/vigo/app"
	"github.com/tamnd/vigo/vio"
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

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "vigo: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	screen := vio.NewTcellScreen()
	if err := screen.Init(); err != nil {
		return fmt.Errorf("init screen: %w", err)
	}
	defer screen.Fini()

	a := app.New(screen)
	return a.Run()
}
