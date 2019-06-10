package main

import (
	"os"
	"path/filepath"
	"sync"
	"syscall"

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
	filter := opt.filter

	wg.Add(1)
	go func() {
		for node := range nodeChan {
			nodes = append(nodes, node)
		}
		wg.Done()
	}()

	rwOpt := &racewalk.Options{
		NumWorkers: opt.workers,
		Debug:      opt.debug,
	}

	if opt.sameFS {
		finfo, err := os.Stat(root)
		if err != nil {
			return nil, err
		}

		if stat := finfo.Sys(); stat != nil {
			ustat, ok := stat.(*syscall.Stat_t)
			if ok {
				filter = filter.sameFS(ustat)
			}
		}
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
			if filter(fnode) {
				continue
			}

			node := makeNode(fnode)
			node.Parent = &parent
			nodeChan <- node
		}

		for _, dir := range subdirs {
			subpath := filepath.Join(path, dir.Name())
			dirParents.LoadOrStore(subpath, &parent)
			if !filter(dir) {
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
