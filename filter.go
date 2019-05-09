package main

import "github.com/samf/racewalk/v2"

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
