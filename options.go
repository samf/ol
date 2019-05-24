package main

import "golang.org/x/sys/unix"

const (
	defaultRows = 24
	defaultCols = 80
)

type options struct {
	debug     bool
	workers   int
	ignoreVCS bool

	sortSize    bool
	sortMtime   bool
	sortReverse bool

	rows int
	cols int

	filter filter
	sorter sorter
}

func (opt *options) valid() error {
	// deal with filters
	opt.filter = noopFilter
	if opt.ignoreVCS {
		opt.filter = opt.filter.noGit().noHG()
	}

	opt.sorter = nameSorter
	if opt.sortSize {
		opt.sorter = opt.sorter.bySize()
	}
	if opt.sortMtime {
		opt.sorter = opt.sorter.byMtime()
	}
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

	return nil
}
