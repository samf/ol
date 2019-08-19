package main

import (
	"os"

	"golang.org/x/sys/unix"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	defaultRows = 24
	defaultCols = 80
)

type options struct {
	debug   bool
	workers int

	dirsize  bool
	treesize bool

	vcs     bool
	sameFS  bool
	sameGit bool
	sameHG  bool
	sameVCS bool

	sortSize    bool
	sortMtime   bool
	sortAtime   bool
	sortCtime   bool
	sortReverse bool

	rows int
	cols int

	filter filter
	sorter sorter

	dirs []string
}

func (opt *options) valid() {
	// dirs must be readable dirs
	opt.dirs = *dirs
	for _, dir := range opt.dirs {
		// we use Stat() instead of Lstat(): for CLI args only, we jump can
		// through symlinks
		info, err := os.Stat(dir)
		if err != nil {
			kingpin.FatalIfError(err, "cannot stat directory %v", dir)
		}
		if !info.IsDir() {
			kingpin.Fatalf("%v is not a directory", dir)
		}
	}

	// treesize implies dirsize
	if opt.treesize {
		opt.dirsize = true
	}

	// deal with filters
	// NB: we don't deal with the sameFS() filter here; it has to be done
	// just prior to calling racewalk.Walk()
	opt.filter = noopFilter
	if !opt.vcs {
		opt.filter = opt.filter.noGit().noHG()
	}
	if opt.sameVCS {
		opt.sameGit = true
		opt.sameHG = true
	}
	if opt.sameGit {
		opt.filter = opt.filter.oneGit()
	}
	if opt.sameHG {
		opt.filter = opt.filter.oneHG()
	}

	// deal with sorters, always starting with nameSorter
	opt.sorter = nameSorter
	if opt.sortSize {
		opt.sorter = opt.sorter.bySize()
	}
	if opt.sortMtime {
		opt.sorter = opt.sorter.byMtime()
	}
	// reverse must be applied last!
	if opt.sortReverse {
		opt.sorter = opt.sorter.reverse()
	}

	// try to get terminal size
	ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err == nil {
		opt.rows = int(ws.Row)
		opt.cols = int(ws.Col)
	} else {
		opt.rows = defaultRows
		opt.cols = defaultCols
	}
}
