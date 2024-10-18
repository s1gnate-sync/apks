#!/bin/sh

set -xeu

cd virtualpkg
CGO_ENABLED=0 go build -o virtualpkg -trimpath
cd ..

cd apk-add
CGO_ENABLED=0 go build -o apk-add -trimpath
cd ..

mkdir -p $DESTDIR/usr/bin/ 
cp virtualpkg/virtualpkg  $DESTDIR/usr/bin/virtualpkg
cp apk-add/apk-add  $DESTDIR/usr/bin/apk-add

mkdir -p $DESTDIR/etc/apk/virtual
mkdir -p $DESTDIR/etc/apk/commit_hooks.d

ln -s /usr/bin/virtualpkg $DESTDIR/etc/apk/commit_hooks.d/virtualpkg-hook


