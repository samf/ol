package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func setup() string {
	dirpath := fmt.Sprintf("%v-%v", time.Now().Unix(), os.Getpid())
	err := os.Mkdir(dirpath, 0777)
	if err != nil {
		panic(err)
	}

	return dirpath
}

func cleanup(dirpath string) {
	err := os.RemoveAll(dirpath)
	if err != nil {
		panic(err)
	}
}

func TestHelloWorld(t *testing.T) {
}
