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
		case astat.Atim.Sec < bstat.Atim.Sec:
			return true
		case astat.Atim.Sec > bstat.Atim.Sec:
			return false
		case astat.Atim.Nsec < bstat.Atim.Nsec:
			return true
		case astat.Atim.Nsec > bstat.Atim.Nsec:
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
		case astat.Ctim.Sec < bstat.Ctim.Sec:
			return true
		case astat.Ctim.Sec > bstat.Ctim.Sec:
			return false
		case astat.Ctim.Nsec < bstat.Ctim.Nsec:
			return true
		case astat.Ctim.Nsec > bstat.Ctim.Nsec:
			return false
		}

		return s(a, b)
	}
}
