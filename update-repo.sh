#!/bin/sh
set -eu
base=$(dirname $(readlink -f $0))
for packagedir in $(find $(dirname $(readlink -f $0))/packages  -maxdepth 1 -type d | tail -n+2); do

  if [ "${1:-}" = "--force" ]; then 
    rm -fr $packagedir/packages
  fi

  if [ ! -d $packagedir/packages ]; then
    $base/build-pkg.sh $(basename $packagedir)
  fi


  dir=$packagedir/packages

  for pkg in $(find $packagedir/packages  -maxdepth 2 -type f  -iname "*.apk" | cut -c $(( ${#dir} + 2 ))-); do
    mkdir -p $base/repository/$(dirname $pkg)
    cp "$dir/$pkg" "$base/repository/$pkg"
  done
done

. $base/env

for archdir in $(find $base/repository -maxdepth 1 -type d | tail -n+2); do
  melange index --output $archdir/APKINDEX.tar.gz --signing-key "$signingkey" $archdir/*.apk 

  rm -f "$archdir/index.html"
done

rm -f "$base/repository/index.html"

set +xeu

python3 -m http.server -b 127.0.0.1 -d $base &
pid="$$"

sleep 1

curl -L http://127.0.0.1:8000/repository/ -o $base/repository/index.html
curl -L http://127.0.0.1:8000/repository/aarch64/ -o $base/repository/aarch64/index.html

jobs -p | xargs -n 1 kill -9




