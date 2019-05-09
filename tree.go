package main

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/samf/racewalk/v2"
)

// GetNodes traverses a tree and returns a slice of Nodes
func GetNodes(root string, opt options) ([]Node, error) {
	var (
		dirParents sync.Map
		wg         sync.WaitGroup
	)
	nodes := []Node{}
	nodeChan := make(chan Node)

	wg.Add(1)
	go func() {
		for node := range nodeChan {
			fmt.Println(node.StatPath)
			nodes = append(nodes, node)
		}
		wg.Done()
	}()

	rwOpt := &racewalk.Options{
		NumWorkers: opt.workers,
		Debug:      opt.debug,
	}

	err := racewalk.Walk(root, rwOpt, func(path string, subdirs,
		entries []racewalk.FileNode) ([]racewalk.FileNode, error) {
		var newsubs []racewalk.FileNode
		parent := makeNode(entries[0])
		dotdot, ok := dirParents.Load(path)
		if ok {
			parent.Parent = dotdot.(*Node)
		}
		nodeChan <- parent

		for _, fnode := range entries[1:] {
			if opt.filter(fnode) {
				continue
			}

			node := makeNode(fnode)
			node.Parent = &parent
			nodeChan <- node
		}

		for _, dir := range subdirs {
			subpath := filepath.Join(path, dir.Name())
			dirParents.LoadOrStore(subpath, &parent)
			if !opt.filter(dir) {
				newsubs = append(newsubs, dir)
			}
		}

		return newsubs, nil
	})
	close(nodeChan)
	if err != nil {
		return nil, err
	}

	wg.Wait()
	return nodes, nil
}
