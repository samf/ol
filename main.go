// Package main is a CLI for the ol command
package main

import (
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
	sortReverse = kingpin.Flag("reverse", "reverse the sort order").Short('r').
			Bool()

	dirs = kingpin.Arg("dirs", "dirs to scan").Default(".").Strings()
)

func main() {
	var nodes []Node

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
	opt.valid()

	for _, dir := range opt.dirs {
		dirnodes, err := GetNodes(dir, opt)
		kingpin.FatalIfError(err, "traversing %v", dir)
		nodes = append(nodes, dirnodes...)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return opt.sorter(&nodes[i], &nodes[j])
	})

	err := pageOut(opt, nodes)
	kingpin.FatalIfError(err, "paging output")
}
