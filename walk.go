package main

import (
	"fmt"
	"os"
	"samf/ctree"
	"syscall"
)

// Walk is the Kong command struct for the 'walk' subcommand
type Walk struct {
	// Paths holds all of the paths that we will walk
	Paths []string `kong:"arg,required,type='existingdir',help='directory(ies) to scan'"`

	// Mtime means sort by mtime
	Mtime bool `kong:"xor='sort',help='recently modified first'"`
	// Size means to sort by size
	Size bool `kong:"short='s',xor='sort',help='largest first'"`
	// Reverse reverses the usual sort order
	Reverse bool `kong:"short='r',help='reverse the sort order'"`

	// Lines gives the maximum number of lines to print
	Lines  int  `kong:"short='l',xor='lines',help='maximum number of lines to show'"`
	Screen bool `kong:"short='S',xor='lines',help='print enough lines to fill the terminal, then stop'"`

	TimeFixed  bool `kong:"help='display fixed time rather than relative to present'"`
	TimeStrict bool `kong:"help='use exact time label rather than human-friendly label'"`
}

// Run runs the walk command
func (w Walk) Run(ctx *Context) error {
	var walks []*ctree.DNode
	for _, dpath := range w.Paths {
		if ctx.verbose {
			fmt.Printf("walking %q\n", dpath)
		}
		root := ctree.NewRoot(dpath)

		walk, err := root.Run()
		if err != nil {
			return fmt.Errorf(
				"failed to walk tree at %q: %w",
				root.Path,
				err,
			)
		}

		if !ctx.quiet {
			for _, err := range walk.Errors() {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}

		walks = append(walks, walk)
	}

	var nodes []*node
	for _, walk := range walks {
		for _, ctNode := range walk.Flatten() {
			fi := ctNode.Info()
			mtime := (*fi).ModTime()
			size := (*fi).Size()
			switch stats := (*fi).Sys().(type) {
			case *syscall.Stat_t:
				size = stats.Blocks * 512
			}
			newNode := &node{
				Node: ctNode,
				time: &mtime,
				size: &size,
			}
			nodes = append(nodes, newNode)
		}
	}

	max := len(nodes)
	if ctx.verbose {
		fmt.Printf("sorting %v objects\n", max)
	}

	switch {
	case w.Size:
		bySize(nodes, w.Reverse)
	case w.Mtime:
		fallthrough
	default:
		byMtime(nodes, w.Reverse)
	}

	max = w.max(ctx, max)
	for i := 0; i < max; i++ {
		nodes[i].printNode(ctx, w, os.Stdout)
	}

	return nil
}

// max computes how many lines to print
func (w Walk) max(ctx *Context, max int) int {
	if w.Screen {
		w.Lines = ctx.height - 2 // 1 for prompt 1 for fuzz
	}

	if w.Lines > max {
		return max
	}
	if w.Lines > 0 {
		return w.Lines
	}

	return max
}
