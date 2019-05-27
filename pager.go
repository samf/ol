package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/mattn/go-isatty"
)

const (
	// DefaultPager is used when PAGER is not in the environment
	DefaultPager = "less"
)

func pageOut(opt options, nodes []Node) error {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		pageFinish(opt, os.Stdout, nodes)
		return nil
	}

	pagerEnv, ok := os.LookupEnv("PAGER")
	if !ok {
		pagerEnv = DefaultPager
	}

	pagerCmd := exec.Command(pagerEnv)
	pagerCmd.Stdout = os.Stdout
	pager, err := pagerCmd.StdinPipe()
	if err != nil {
		pageFinish(opt, os.Stdout, nodes)
		return err
	}
	err = pagerCmd.Start()
	if err != nil {
		pageFinish(opt, os.Stdout, nodes)
		return err
	}

	pageFinish(opt, pager, nodes)
	pager.Close()
	return pagerCmd.Wait()
}

func pageFinish(opt options, out io.Writer, nodes []Node) error {
	for _, node := range nodes {
		fmt.Fprintln(out, node.format(opt))
	}

	return nil
}
