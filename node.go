package main

import (
	"fmt"

	"github.com/samf/racewalk/v2"
)

// Node represents a single file, in the UNIX sense where a file can be a
// regular file, directory, symlink, etc.
type Node struct {
	racewalk.FileNode

	Parent *Node
	Size   int64
}

func (node Node) String() string {
	return fmt.Sprintf("%5d %s", node.Size, node.StatPath)
}

func makeNode(fnode racewalk.FileNode) Node {
	return Node{
		FileNode: fnode,
		Size:     fnode.Size(),
	}
}
