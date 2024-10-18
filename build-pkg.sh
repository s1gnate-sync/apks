#!/bin/sh
set -xeu
dir=$(dirname $(readlink -f $0))
export packagedir="$dir/packages/${1:?specify packagedir}"
eval export $(cat $dir/env)
exec yaegi run -unrestricted -syscall -noautoimport -unsafe  $dir/bin/build.go
