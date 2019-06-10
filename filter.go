package main

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/samf/racewalk/v2"
)

type filter func(racewalk.FileNode) bool

func noopFilter(racewalk.FileNode) bool {
	return false
}

func (f filter) sameFS(orig *syscall.Stat_t) filter {
	return func(n racewalk.FileNode) bool {
		if stat := n.GetStat(); stat != nil {
			if orig.Dev != stat.Dev {
				return true
			}
		}
		if f == nil {
			return false
		}

		return f(n)
	}
}

func (f filter) noGit() filter {
	return func(n racewalk.FileNode) bool {
		switch {
		case n.Name() == ".git":
			return true
		case f == nil:
			return false
		}

		return f(n)
	}
}

func (f filter) noHG() filter {
	return func(n racewalk.FileNode) bool {
		switch {
		case n.Name() == ".hg":
			return true
		case f == nil:
			return false
		}

		return f(n)
	}
}

func (f filter) oneGit() filter {
	return func(n racewalk.FileNode) bool {
		gitpath := filepath.Join(n.StatPath, ".git")
		switch _, err := os.Stat(gitpath); {
		case err == nil:
			return true
		case f == nil:
			return false
		}

		return f(n)
	}
}

func (f filter) oneHG() filter {
	return func(n racewalk.FileNode) bool {
		hgpath := filepath.Join(n.StatPath, ".hg")
		switch _, err := os.Stat(hgpath); {
		case err == nil:
			return true
		case f == nil:
			return false
		}

		return f(n)
	}
}
