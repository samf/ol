// Package main is a CLI for the ol command
package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug   = kingpin.Flag("debug", "enable debug mode").Bool()
	workers = kingpin.Flag("workers", "number of workers").Short('w').
		Default("0").Int()
	ignoreVCS = kingpin.Flag("ignore-vcs",
		"ignore VCS directories like .git and .hg").Bool()
)

func main() {
	kingpin.Parse()
	opt := options{
		debug:     *debug,
		workers:   *workers,
		ignoreVCS: *ignoreVCS,
	}
	err := opt.valid()
	if err != nil {
		return
	}

	nodes, err := GetNodes(".", opt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
		return
	}

	for _, node := range nodes {
		fmt.Println(node.format(opt))
	}
}
