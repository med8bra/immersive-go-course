package main

import (
	"flag"
	"fmt"
	"go-ls/cmd"
)

var directory string = *flag.String("directory", ".", "Directory to list")

func main() {
	flag.Parse()

	err := cmd.Ls(directory)
	if err != nil {
		fmt.Printf("go-ls: %s", err)
	}
}
