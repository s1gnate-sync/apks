package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	Path = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
	Term = "xterm"
	Lang = "C"
	Home = "/home/build"
)

func main() {
	exe, _ := os.Executable()
	if exe != "" {
		os.Chdir(filepath.Join(filepath.Dir(exe), ".."))
	}

	prefix, _ := filepath.Abs(".")
	if prefix == "" {
		prefix = "."
	}

	config := os.Getenv("CONFIG")
	if config == "" {
		os.Exit(1)
	}

	srcArg := "--empty-workspace"
	if _, err := os.Stat(filepath.Join(prefix, config, "melange.yaml")); err == nil {
		config = filepath.Join(prefix, config, "melange.yaml")
		if _, err := os.Stat(filepath.Join(prefix, config, "src")); err == nil {
			srcArg = "--src-dir " + filepath.Join(prefix, config, "src")
		}
	} else {
		config = filepath.Join(prefix, config+".yaml")
	}

	outDir := filepath.Join(prefix, "out")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "(buildroot) %s: %s\n", outDir, err)
	}

	tempRoot := filepath.Join(outDir, "tmp")
	if err := os.MkdirAll(tempRoot, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "(buildroot) %s: %s\n", tempRoot, err)
	}

	tempDir, err := os.MkdirTemp(tempRoot, "*")
	if err != nil {
		tempDir = tempRoot
	}

	if err = os.MkdirAll(tempDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "(buildroot) %s: %s\n", tempDir, err)
	}

	melangeCacheDir := filepath.Join(outDir, "cache", "melange")
	if err = os.MkdirAll(melangeCacheDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "(buildroot) %s: %s\n", melangeCacheDir, err)
	}

	apkoCacheDir := filepath.Join(outDir, "cache", "apko")
	if err = os.MkdirAll(apkoCacheDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "(buildroot) %s: %s\n", apkoCacheDir, err)
	}


	 fmt.Println(strings.Join([]string{
		"bwrap",
		"--clearenv",
		"--setenv", "PATH", Path,
		"--setenv", "TERM", Term,
		"--setenv", "HOME", Home,
		"--setenv", "LC_ALL", Lang,
		"--bind", "/", "/",
		"--bind", tempDir, "/tmp",
		"--dev", "/dev",
		"--proc", "/proc",
		filepath.Join(prefix, "bin", "melange." + runtime.GOARCH),
		"build",
		srcArg,
		"--log-level", "info",
		"--runner", "bubblewrap",
		"--out-dir", outDir,
		"--cache-dir", melangeCacheDir,
		"--arch", runtime.GOARCH,
		"--apk-cache-dir", apkoCacheDir,
		"--workspace-dir", filepath.Join(tempDir, "workspace"),
		"--guest-dir", filepath.Join(tempDir, "guest"),
		config,
	}, " "))
}
