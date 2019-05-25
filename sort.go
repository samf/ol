package main

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
		switch {
		case a.Size < b.Size:
			return true
		case a.Size > b.Size:
			return false
		}

		return s(a, b)
	}
}

func (s sorter) byMtime() sorter {
	return func(a, b *Node) bool {
		switch {
		case a.ModTime().Before(b.ModTime()):
			return true
		case a.ModTime().After(b.ModTime()):
			return false
		}

		return s(a, b)
	}
}
