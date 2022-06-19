package main

import "sort"

// sort by mtime
func byMtime(nodes []*node, reverse bool) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].time == nil {
			fi := nodes[i].Info()
			mtime := (*fi).ModTime()
			nodes[i].time = &mtime
		}
		if nodes[j].time == nil {
			fi := nodes[j].Info()
			mtime := (*fi).ModTime()
			nodes[j].time = &mtime
		}

		// order "reversed" on purpose: normal case is
		// newest date first
		res := nodes[j].time.Before(*nodes[i].time)
		if reverse {
			res = !res
		}

		return res
	})
}

// sort by size
func bySize(nodes []*node, reverse bool) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].size == nil {
			fi := nodes[i].Info()
			size := (*fi).Size()
			nodes[i].size = &size
		}
		if nodes[j].size == nil {
			fi := nodes[j].Info()
			size := (*fi).Size()
			nodes[j].size = &size
		}

		res := *nodes[i].size > *nodes[j].size
		if reverse {
			res = !res
		}

		return res
	})
}
