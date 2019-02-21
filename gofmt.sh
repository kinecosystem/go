#! /bin/bash
set -e

OUTPUT=$(ls -d */ \
  | egrep -v '^vendor|^docs' \
  | xargs -I {} -P 4 gofmt -l {})

if [[ $OUTPUT ]]; then
  printf "gofmt found unformatted files:\n\n"
  echo "$OUTPUT"
  exit 1
fi
