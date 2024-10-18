package list

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Row struct {
	Name    string
	Version string
	Deps    []string
}

func Main(args []string) {
	file, err := os.Open("/lib/apk/db/installed")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", args[0], err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	item := Row{}
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			if strings.HasPrefix(item.Name, "virtual-") && len(item.Deps) > 0 {
				println(strings.Join(item.Deps, "\n"))
			}
			item = Row{}
			continue
		}

		key, value, ok := strings.Cut(text, ":")
		if !ok {
			fmt.Fprintf(os.Stderr, "%s: unexpected value '%s'\n", args[0], item)
			os.Exit(1)
		}

		if key == "P" {
			item.Name = value
		} else if key == "D" {
			item.Deps = strings.Split(value, " ")
		} else if key == "V" {
			item.Version = value
		}
	}
}
