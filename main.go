// Package main is a CLI for the ol command
package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug = kingpin.Flag("debug", "enable debug mode").Bool()
)

func main() {
	fmt.Println("hello world")
	fmt.Printf("debug %v\n", *debug)

	nodes, err := GetNodes(".")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
		return
	}

	for _, node := range nodes {
		fmt.Printf("%q\n", node.Path)
	}
}
