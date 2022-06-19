package main

import (
	"fmt"
	"io"
	"samf/ctree"
	"strconv"
	"sync"
	"time"
)

type node struct {
	ctree.Node

	time *time.Time
	size *int64
}

var (
	now  time.Time
	once sync.Once
)

func (n *node) printNode(ctx *Context, w Walk, out io.Writer) {
	width := ctx.width - 1 // one for margin

	dsize, taken := n.truncSize(w)
	width -= taken + 1 // one for margin after this field (same for all fields)

	dtime, taken := n.truncTime(w)
	width -= taken + 1

	dpath := n.truncPath(w, width)
	fmt.Fprintf(out, "%v %v %v\n", dsize, dtime, dpath)
}

func (n *node) truncSize(w Walk) (string, int) {
	size := *n.size
	unit := "B"

	if size > 1024 {
		size /= 1024
		unit = "K"
	}
	if size > 1024 {
		size /= 1024
		unit = "M"
	}
	if size > 1024 {
		size /= 1024
		unit = "G"
	}
	if size > 1024 {
		size /= 1024
		unit = "T"
	}

	width := 4
	return fmt.Sprintf("%*v%v", width, size, unit), width
}

func (n *node) truncTime(w Walk) (string, int) {
	var (
		width int
		value string
	)

	once.Do(func() {
		now = time.Now()
	})

	age := now.Sub(*n.time)
	switch {
	case w.TimeFixed && w.TimeStrict:
	case w.TimeFixed:
	case w.TimeStrict:
		value = fmt.Sprintf("%v", age)
	default:
		value, width = friendlyDuration(age)
	}
	return fmt.Sprintf("%*v", width, value), width
}

func (n *node) truncPath(w Walk, width int) string {
	path := n.Path()
	if chopped := len(path) - width; width > 0 && chopped > 0 {
		path = path[chopped:]
	}
	return path
}

func friendlyDuration(age time.Duration) (string, int) {
	var (
		value string
		unit  string
		width = 4 // allows for 9999 weeks, which is ~192 years
	)
	switch {
	case age > 7*24*time.Hour:
		value = strconv.Itoa(int(age / (7 * 24 * time.Hour)))
		unit = "w"
	case age > 24*time.Hour:
		value = strconv.Itoa(int(age / (24 * time.Hour)))
		unit = "d"
	case age > time.Hour:
		value = strconv.Itoa(int(age / time.Hour))
		unit = "h"
	case age > time.Minute:
		value = strconv.Itoa(int(age / time.Minute))
		unit = "m"
	default:
		value = strconv.Itoa(int(age / time.Second))
		unit = "s"
	}
	return fmt.Sprintf("%*v%v", width, value, unit), width + 1 // +1 for unit
}
