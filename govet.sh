#! /bin/bash
set -e

find . -maxdepth 1 -mindepth 1 -type d \
  | egrep -v '^\.\/vendor|^.\/docs|^\.\/\..*' \
  | xargs -I {} -P 4 go tool vet -all -composites=false -unreachable=false -tests=false -shadow {}
