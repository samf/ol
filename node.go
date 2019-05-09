package main

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/samf/racewalk/v2"
)

// Node represents a single file, in the UNIX sense where a file can be a
// regular file, directory, symlink, etc.
type Node struct {
	racewalk.FileNode

	Parent *Node
	Size   int64
}

func (node Node) format(opt options) string {
	avail := opt.cols - 1
	avail -= 6  // size
	avail -= 17 // when

	when := node.getTime(opt)

	path := node.StatPath
	if len(path) > avail {
		path = path[:avail]
	}

	return fmt.Sprintf("%6s %15s %s", node.getSize(opt), when, path)
}

func (node Node) getSize(opt options) string {
	size := uint64(node.Size)
	if stat := node.GetStat(); stat != nil {
		realSize := uint64(stat.Blocks) * uint64(stat.Blksize)
		if realSize < size {
			size = realSize
		}
	}

	return humanize.Bytes(size)
}

func (node Node) getTime(opt options) string {
	return humanize.Time(node.ModTime())
}

func makeNode(fnode racewalk.FileNode) Node {
	return Node{
		FileNode: fnode,
		Size:     fnode.Size(),
	}
}
