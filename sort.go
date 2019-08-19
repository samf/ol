package main

import "syscall"

type sorter func(a, b *Node) bool

func nameSorter(a, b *Node) bool {
	return a.StatPath < b.StatPath
}

func (s sorter) reverse() sorter {
	return func(a, b *Node) bool {
		return s(b, a)
	}
}

func (s sorter) bySize() sorter {
	return func(a, b *Node) bool {
		// largest first
		switch {
		case *a.Size > *b.Size:
			return true
		case *a.Size < *b.Size:
			return false
		}

		return s(a, b)
	}
}

func (s sorter) byMtime() sorter {
	return func(a, b *Node) bool {
		switch {
		// newest first
		case a.ModTime().After(b.ModTime()):
			return true
		case a.ModTime().Before(b.ModTime()):
			return false
		}

		return s(a, b)
	}
}

func (s sorter) byAtime() sorter {
	return func(a, b *Node) bool {
		astat := a.Sys().(*syscall.Stat_t)
		bstat := b.Sys().(*syscall.Stat_t)

		switch {
		// newest first
		case astat.Atimespec.Sec < bstat.Atimespec.Sec:
			return true
		case astat.Atimespec.Sec > bstat.Atimespec.Sec:
			return false
		case astat.Atimespec.Nsec < bstat.Atimespec.Nsec:
			return true
		case astat.Atimespec.Nsec > bstat.Atimespec.Nsec:
			return false
		}

		return s(a, b)
	}
}

func (s sorter) byCtime() sorter {
	return func(a, b *Node) bool {
		astat := a.Sys().(*syscall.Stat_t)
		bstat := b.Sys().(*syscall.Stat_t)

		switch {
		// newest first
		case astat.Ctimespec.Sec < bstat.Ctimespec.Sec:
			return true
		case astat.Ctimespec.Sec > bstat.Ctimespec.Sec:
			return false
		case astat.Ctimespec.Nsec < bstat.Ctimespec.Nsec:
			return true
		case astat.Ctimespec.Nsec > bstat.Ctimespec.Nsec:
			return false
		}

		return s(a, b)
	}
}
