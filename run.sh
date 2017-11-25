#!/usr/bin/env bash

set -euf -o pipefail

NUM=0

rm -rf data/
mkdir data/

while read url; do
  echo "${url}"
  wget -O data/${NUM}.lz4 "${url}"
  lz4 -d "data/${NUM}.lz4" "data/${NUM}"
  rm "data/${NUM}.lz4" 
  go run parse.go "data/${NUM}" > "data/parsed-${NUM}"
  rm "data/${NUM}" 
  NUM=$((NUM+1))
done <$1
