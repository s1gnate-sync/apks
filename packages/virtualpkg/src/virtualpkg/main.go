package main

import (
	"fmt"
	"os"

	"virtualpkg/hook"
	"virtualpkg/list"
	"virtualpkg/pkg"
	"virtualpkg/sync"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s ACTION [ARGS...]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Actions: \n  \"post-commit|pre-commit\" - runs commit hook, \n  \"cache\" - saves all virtual packages to current directory, \n  \"list\" - list installed virtual packages, \n  \"build\" - creates new package and writes it to stdout")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	action := os.Args[1]
	if action == "pre-commit" {
		hook.MainPre(os.Args[1:])
	} else if action == "post-commit" {
		hook.MainPost(os.Args[1:])
	} else if action == "cache" {
		sync.Main(os.Args[1:])
	} else if action == "list" {
		list.Main(os.Args[1:])
	} else if action == "build" {
		pkg.Main(os.Args[1:])
	} else {
		usage()
	}
}
