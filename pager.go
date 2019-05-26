package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/mattn/go-isatty"
)

const (
	// DefaultPager is used when PAGER is not in the environment
	DefaultPager = "less"
)

// InstallPager returns an io.WriteCloser and an error. If no pager is
// installed, (nil, nil) will be returned, and os.Stdout should be used.
// The pager should be Closed().
func InstallPager(opt options) (io.WriteCloser, error) {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return nil, nil
	}

	pagerEnv, ok := os.LookupEnv("PAGER")
	if !ok {
		pagerEnv = DefaultPager
	}

	pagerCmd := exec.Command(pagerEnv)
	pagerCmd.Stdout = os.Stdout
	pager, err := pagerCmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	pagerCmd.Start()

	return pager, nil
}
