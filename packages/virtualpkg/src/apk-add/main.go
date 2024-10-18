package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func add() {
}

func main() {
	args := os.Args[1:]
	if len(os.Args) == 0 {
		return
	}

	name := ""
	if args[0] == "-t" || args[0] == "--virtual" {
		if len(args) == 1 {
			os.Exit(1)
		}

		name = args[1]
		args = args[2:]
	}

	if len(args) == 0 {
		return
	}

	if name == "" {
		for _, arg := range args {
			addVirtual(arg, arg)
		}
	} else {
		addVirtual(name, args...)
	}
}

func addVirtual(name string, args ...string) {
	deps := map[string]bool{}

	name = fmt.Sprintf("virtual-%s", strings.TrimPrefix(name, "virtual-"))
	for _, arg := range args {
		deps[arg] = true
	}

	data, _ := exec.Command("/sbin/apk", "info", name, "--depends").Output()
	entries := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(entries) > 0 {
		for _, arg := range entries[1:] {
			arg = strings.TrimSpace(arg)
			if arg != "" {
				deps[arg] = true
			}
		}
	}

	cmd := exec.Command("/sbin/apk", "add", "--virtual", name)
	for dep := range deps {
		cmd.Args = append(cmd.Args, dep)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
	}
}
