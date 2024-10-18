package pkg

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func CreateVirtualPackage(name, version string, deps ...string) ([]byte, error) {
	if version == "" {
		version = fmt.Sprintf("%d", time.Now().UnixMicro())
	}

	info := []byte(strings.Join([]string{
		fmt.Sprintf("pkgname = virtual-%s", strings.TrimPrefix(name, "virtual-")),
		fmt.Sprintf("pkgver = %s", version),
		fmt.Sprintf("depend = %s", strings.Join(deps, " ")),
		"arch = noarch",
		"url = ",
		"size = 0",
		"pkgdesc = ",
	}, "\n"))

	buf := bytes.NewBuffer(nil)
	gzipWriter := gzip.NewWriter(buf)
	tarWriter := tar.NewWriter(gzipWriter)

	err := tarWriter.WriteHeader(&tar.Header{
		Name: ".PKGINFO",
		Mode: 0o644,
		Size: int64(binary.Size(info)),
		Uid:  0,
		Gid:  0,
	})

	if err == nil {
		_, err = tarWriter.Write(info)
	}

	tarWriter.Close()
	gzipWriter.Close()

	return buf.Bytes(), err
}

func Main(args []string) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s NAME PACKAGE [PACKAGE...]\n", os.Args[0])
		os.Exit(1)
	}

	data, err := CreateVirtualPackage(args[1], "", args[2:]...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", args[0], err)
		os.Exit(1)
	}

	os.Stdout.Write(data)
	os.Stdout.Sync()
}
