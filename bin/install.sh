#!/bin/sh
set -ex

dir="$(readlink -f $(dirname $0))"
tmp=$(mktemp -d)
cd $tmp
trap 'rm -fr $tmp' EXIT

wget https://github.com/chainguard-dev/melange/releases/download/v0.13.6/melange_0.13.6_linux_arm64.tar.gz

wget https://github.com/chainguard-dev/melange/releases/download/v0.13.6/melange_0.13.6_linux_amd64.tar.gz

wget https://github.com/traefik/yaegi/releases/download/v0.16.1/yaegi_v0.16.1_linux_arm64.tar.gz

wget https://github.com/traefik/yaegi/releases/download/v0.16.1/yaegi_v0.16.1_linux_amd64.tar.gz

tar xvf yaegi_v0.16.1_linux_amd64.tar.gz
mv yaegi $dir/yaegi.amd64

tar xvf yaegi_v0.16.1_linux_arm64.tar.gz
mv yaegi $dir/yaegi.arm64

tar xvf melange_0.13.6_linux_arm64.tar.gz
mv melange_0.13.6_linux_arm64/melange $dir/melange.arm64

tar xvf melange_0.13.6_linux_amd64.tar.gz
mv melange_0.13.6_linux_amd64/melange $dir/melange.amd64

