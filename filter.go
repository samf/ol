package main

import (
	"os"
	"path/filepath"

	"github.com/samf/racewalk/v2"
)

type filter func(racewalk.FileNode) bool

func noopFilter(racewalk.FileNode) bool {
	return false
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
