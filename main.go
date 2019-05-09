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
)

func main() {
	kingpin.Parse()
	opt := options{
		debug:   *debug,
		workers: *workers,
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
		fmt.Printf("%q\n", node.Path)
	}
}
