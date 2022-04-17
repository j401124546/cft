#!/bin/bash
DIR=$(pwd)

if [ -d "${DIR}/output" ]
then
  rm -rf "${DIR}/output"
fi

mkdir "${DIR}/output"
mkdir "${DIR}/output/config"

go build -o "${DIR}/output/app"
cp "${DIR}/config/application.json" "${DIR}/output/config"

git log --author="jiayuke" --pretty=tformat: --numstat | awk '{ add += $1; subs += $2; loc += $1 - $2 } END { printf "added lines: %s, removed lines : %s, total lines: %s\n", add, subs, loc }'