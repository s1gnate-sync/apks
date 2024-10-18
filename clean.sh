#!/bin/sh

for dir in $(find -type d -iname "temp[0-9]*"); do
  rm -fr "$dir"
done
