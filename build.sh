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

if [ -z "${1:-}" ]; then
  :
  # rm -fr out
  # for file in $(find .  -maxdepth 1 -iname '*.yaml'); do
  #   $(CONFIG="${file##*/}" $PWD/bin/yaegi.$arch run $PWD/bin/build.go)
  # done
else
  $(CONFIG="$1" $PWD/bin/yaegi.$arch run $PWD/bin/build.go)
fi

for dir in out/apk/*; do
  _dir=$PWD
  cd $dir
  find . -type f -iname '*.apk' | xargs -r $_dir/bin/melange.$arch index 
  cd $_dir
done
