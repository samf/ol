package main

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/samf/racewalk/v2"
)

// GetNodes traverses a tree and returns a slice of Nodes
func GetNodes(root string) ([]Node, error) {
	var (
		wg   sync.WaitGroup
		tree sync.Map
	)
	nodes := []Node{}
	nodeChan := make(chan Node)

	wg.Add(1)
	go func() {
		for node := range nodeChan {
			fmt.Println(node.Path)
			nodes = append(nodes, node)
		}
		wg.Done()
	}()

	opt := &racewalk.Options{
		NumWorkers: 0,
		Debug:      true,
	}

	err := racewalk.Walk(root, opt, func(path string, dirs,
		others []racewalk.FileNode) ([]racewalk.FileNode, error) {
		var parent *Node
		node := Node{}
		iparent, ok := tree.Load(path)
		if ok {
			parent = iparent.(*Node)
		}

		for _, fnode := range append(dirs, others...) {
			path := filepath.Join(path, fnode.Name())
			node = Node{
				FileNode: fnode,
				Parent:   parent,
				Path:     path,
				Size:     fnode.Size(),
			}
			tree.LoadOrStore(path, &node)
			nodeChan <- node
		}
		return dirs, nil
	})
	close(nodeChan)
	if err != nil {
		return nil, err
	}

	wg.Wait()
	return nodes, nil
}
