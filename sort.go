package main

type sorter func(a, b *Node) bool

func nameSorter(a, b *Node) bool {
	return a.StatPath < b.StatPath
}

func (s sorter) reverse() sorter {
	return func(a, b *Node) bool {
		return !s(a, b)
	}
}

func (s sorter) bySize() sorter {
	return func(a, b *Node) bool {
		if a.Size < b.Size {
			return true
		}

		return s(a, b)
	}
}

func (s sorter) byMtime() sorter {
	return func(a, b *Node) bool {
		if a.ModTime().Before(b.ModTime()) {
			return true
		}

		return s(a, b)
	}
}
