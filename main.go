// Package main is a CLI for the ol command
package main

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug   = kingpin.Flag("debug", "enable debug mode").Bool()
	workers = kingpin.Flag("workers", "number of workers").Short('w').
		Default("0").Int()
	vcs = kingpin.Flag("vcs",
		"descend into VCS directories like .git and .hg").Bool()
	sameGit = kingpin.Flag("same-git", "don't descend into new git repos").
		Bool()
	sameHG = kingpin.Flag("same-hg", "don't descend into new mercurial repos").
		Bool()
	sameVCS = kingpin.Flag("same-vcs",
		"don't descend into new repos (same as --same-git and --same-hg)").
		Bool()

	sortSize  = kingpin.Flag("size", "sort by size").Short('s').Bool()
	sortMtime = kingpin.Flag("mtime", "sort by modification time").Short('m').
			Bool()
	sortReverse = kingpin.Flag("reverse", "reverse the sort order").Short('r').Bool()
)

func main() {
	kingpin.Parse()
	opt := options{
		debug:       *debug,
		workers:     *workers,
		vcs:         *vcs,
		sameGit:     *sameGit,
		sameHG:      *sameHG,
		sameVCS:     *sameVCS,
		sortSize:    *sortSize,
		sortMtime:   *sortMtime,
		sortReverse: *sortReverse,
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

	sort.Slice(nodes, func(i, j int) bool {
		return opt.sorter(&nodes[i], &nodes[j])
	})

	err = pageOut(opt, nodes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
