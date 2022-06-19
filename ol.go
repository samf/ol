package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"golang.org/x/crypto/ssh/terminal"
)

// Context is for passing other context down to subcommands
type Context struct {
	quiet   bool
	verbose bool
	width   int
	height  int
}

// CLI defines the CLI
var CLI struct {
	// Walk is the walk subcommand
	Walk Walk `kong:"cmd,default='withargs'"`

	// Quiet surpresses warnings and non-fatal errors
	Quiet bool `kong:"short='q',xor='loud',help='surpress warnings and non-fatal errors'"`
	// Verbose gives more status output
	Verbose bool `kong:"short='v',xor='loud',help='give more verbose status output'"`
}

func main() {
	ctx := kong.Parse(&CLI)

	cmdCtx := newContext()

	err := ctx.Run(cmdCtx)
	ctx.FatalIfErrorf(err)
}

func newContext() *Context {
	width, height, err := terminal.GetSize(0)
	if err != nil {
		if CLI.Verbose {
			fmt.Fprintf(
				os.Stderr,
				"cannot determine screen dimensions: %v",
				err.Error(),
			)
		}

		// do nothing; some things will fail
		width = 0
		height = 0
	}

	return &Context{
		quiet:   CLI.Quiet,
		verbose: CLI.Verbose,
		width:   width,
		height:  height,
	}
}
