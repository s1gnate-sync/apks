package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	path                 = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
	term                 = "xterm"
	lang                 = "C"
	tempDir              = ""
	progName             = "(build)"
	packageDir           = ""
	configPath           = ""
	command    *exec.Cmd = nil
)

func printErr(params ...string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", progName, strings.Join(params, " "))
}

func setupTempDir() {
	tempDir, _ = os.MkdirTemp(packageDir, "temp*")

	if tempDir == "" {
		printErr("can't create temporary dir", "packageDir="+packageDir)
		os.Exit(1)
	}
}

func setupPackage() {
	path := os.Getenv("packagedir")
	if path == "" {
		printErr("env variable packagedir is undefined")
		os.Exit(1)
	}

	if _, err := os.Stat(path); err != nil {
		printErr("invalid value: ", err.Error(), "packagedir="+path)
		os.Exit(1)
	}

	config := filepath.Join(path, "config.yaml")
	if _, err := os.Stat(config); err != nil {
		printErr("no config: ", err.Error(), "config="+config)
		os.Exit(1)
	}

	configPath = config
	packageDir = path
}

func setupCommand() {
	command = exec.Command(
		"bwrap",
		"--clearenv",
		"--setenv", "PATH", path,
		"--setenv", "TERM", term,
		"--setenv", "HOME", tempDir,
		"--setenv", "LC_ALL", lang,
		"--bind", "/", "/",
		"--bind", tempDir, "/tmp",
		"--dev", "/dev",
		"--proc", "/proc",
		"--",
		"melange",
		"build",
		configPath,
		"--rm",
		"--create-build-log",
		"--generate-index=0",
		"--log-level", "info",
		"--runner", "bubblewrap",
		"--out-dir", filepath.Join(packageDir, "packages"),
		"--signing-key", os.Getenv("key"),
	)

	arch := os.Getenv("arch")
	if arch == "" {
		arch = runtime.GOARCH
	}
	command.Args = append(command.Args, "--arch", arch)

	srcDir := filepath.Join(packageDir, "src")
	if _, err := os.Stat(srcDir); err == nil {
		command.Args = append(command.Args, "--source-dir", srcDir)
	} else {
		command.Args = append(command.Args, "--empty-workspace")
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
}

func init() {
	setupPackage()
	setupTempDir()
	setupCommand()
}

func clearTmp() {
	if err := os.RemoveAll(tempDir); err != nil {
		// printErr("temporary dir not removed:", err.Error(), "tempdir="+tempDir)

		filepath.Walk(tempDir, func(path string, info fs.FileInfo, err error) error {
			if err == nil {
				os.Chmod(path, 0o777)
			}
			return err
		})

		if err := os.RemoveAll(tempDir); err != nil {
			printErr("temporary dir not removed:", err.Error(), "tempdir="+tempDir)
		}
	}
}

func main() {
	printErr(command.Args...)

	err := command.Run()
	if err != nil {
		printErr("command failed:", err.Error())
		clearTmp()
		os.Exit(1)
	}

	clearTmp()
}
