#!/bin/sh

set -eu

cd "$(dirname $(readlink -f $0))"

arch=""
case "$(uname -m)" in
  aarch64|arm64)
    arch=arm64
  ;;
  x86_64|amd64)
    arch=amd64
  ;;
  *)
    exit 1
  ;;  
esac


$(CONFIG="$1" $PWD/bin/yaegi.$arch run $PWD/bin/build.go)
