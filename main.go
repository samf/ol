// Package main is a CLI for the ol command
package main

import (
	"fmt"

	"github.com/alecthomas/kingpin"
)

var (
	debug = kingpin.Flag("debug", "enable debug mode").Bool()
)

func main() {
	fmt.Println("hello world")
	fmt.Printf("debug %v", debug)
}
