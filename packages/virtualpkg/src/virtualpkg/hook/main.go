package hook

import (
	"fmt"
	"os"
	"strings"
	"time"
	"virtualpkg/sync"
)

func apkCmdline() []string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", os.Getppid()))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return strings.Split(string(data), "\000")
}

func MainPost(args []string) {
	file, err := os.OpenFile("/var/log/apk-post-commit.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", args[0], err)
		os.Exit(1)
	}

	file.WriteString(fmt.Sprintf("%d\t%s\n", time.Now().Unix(), strings.Join(apkCmdline(), " ")))
	file.Close()

	err = os.Chdir("/etc/apk/virtual")
	if err == nil {
		sync.Main(args)
	}
}

func MainPre(args []string) {
	cmdline := apkCmdline()
	name := cmdline[0]
	cmdline = cmdline[1:]

	if name != "/sbin/apk" {
		fmt.Fprintf(os.Stderr,
			"(%s) %s\n",
			args[0],
			"apk requires execution with an absolute path",
		)
		os.Exit(1)
	}

	add := false
	virtual := ""
	for index, arg := range args {
		if arg == "add" {
			add = true
		}

		if arg == "-t" || arg == "--virtual" {
			if len(args)-index+1 >= 1 {
				virtual = args[index+1]
			}
		}
	}

	if add {
		if virtual == "" {
			fmt.Fprintf(os.Stderr,
				"(%s) apk add: %s\n",
				args[0],
				"missing mandatory flag --virtual",
			)
			os.Exit(1)
		} else if !strings.HasPrefix(virtual, "virtual-") {
			fmt.Fprintf(os.Stderr,
				"(%s) apk add: %s\n",
				args[0],
				"wrong format for --virtual flag, must be 'virtual-...'",
			)
			os.Exit(1)
		}
	}
}
