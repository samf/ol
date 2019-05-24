package main

import (
	"fmt"
	"time"

	"github.com/samf/racewalk/v2"
	"github.com/samf/tier"
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
	avail -= 6 // size
	avail -= 5 // when

	when := node.getTime(opt)

	path := node.StatPath
	if len(path) > avail {
		path = path[:avail]
	}

	return fmt.Sprintf("%5s %4s %s", node.getSize(opt), when, path)
}

func (node Node) getSize(opt options) string {
	size := uint64(node.Size)
	if stat := node.GetStat(); stat != nil {
		realSize := uint64(stat.Blocks) * uint64(stat.Blksize)
		if realSize < size {
			size = realSize
		}
	}

	return tier.Bytes.Make(int64(size)).Short()
}

func (node Node) getTime(opt options) string {
	age := time.Since(node.ModTime())
	return tier.Time.Make(int64(age.Seconds())).Short()
}

func makeNode(fnode racewalk.FileNode) Node {
	return Node{
		FileNode: fnode,
		Size:     fnode.Size(),
	}
}
