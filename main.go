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
	ignoreVCS = kingpin.Flag("ignore-vcs",
		"ignore VCS directories like .git and .hg").Bool()

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
		ignoreVCS:   *ignoreVCS,
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

	sort.SliceStable(nodes, func(i, j int) bool {
		return opt.sorter(&nodes[i], &nodes[j])
	})

	for _, node := range nodes {
		fmt.Println(node.format(opt))
	}
}
