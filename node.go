package main

import (
	"fmt"
	"os"
	"time"

	"github.com/samf/racewalk/v2"
	"github.com/samf/tier"
)

// Node represents a single file, in the UNIX sense where a file can be a
// regular file, directory, symlink, etc.
type Node struct {
	racewalk.FileNode

	Parent *Node
	Size   *int64
}

func (node Node) format(opt options) string {
	avail := opt.cols - 1
	avail -= 6 // size
	avail -= 5 // when

	when := node.getTime(opt)

	path := node.getPath(opt)
	if len(path) > avail {
		path = path[:avail]
	}

	return fmt.Sprintf("%5s %4s %s", node.getSize(opt), when, path)
}

func sizeUp(nodes []Node, opt options) {
	for _, node := range nodes {
		if node.IsDir() {
			continue
		}

		node.fixFileSize()
	}

	if !opt.dirsize {
		return
	}

	for _, node := range nodes {
		if node.IsDir() {
			continue
		}

		if node.Parent != nil {
			*node.Parent.Size += *node.Size
		}
	}

	if !opt.treesize {
		return
	}

	for _, node := range nodes {
		if !node.IsDir() {
			continue
		}

		for n := node.Parent; n != nil; n = n.Parent {
			*n.Size += *node.Size
		}
	}
}

func (node Node) fixFileSize() {
	size := *node.Size
	if stat := node.GetStat(); stat != nil {
		realSize := int64(stat.Blocks) * int64(stat.Blksize)
		if realSize < size {
			*node.Size = realSize
		}
	}
}

func (node Node) getSize(opt options) string {
	return tier.Bytes.Make(*node.Size).Short()
}

func (node Node) getTime(opt options) string {
	age := time.Since(node.ModTime())
	return tier.Time.Make(int64(age.Seconds())).Short()
}

func (node Node) getPath(opt options) string {
	path := node.StatPath

	switch mode := node.Mode(); {
	case mode.IsRegular():
	case mode.IsDir():
		path += "/"
	case mode&os.ModeSymlink != 0:
		path += "@"
	case mode&os.ModeNamedPipe != 0:
		path += "|"
	}

	return path
}

func makeNode(fnode racewalk.FileNode) Node {
	size := new(int64)
	*size = fnode.Size()
	return Node{
		FileNode: fnode,
		Size:     size,
	}
}
