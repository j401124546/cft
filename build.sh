#!/bin/bash
DIR=$(pwd)
echo $DIR

if [ -d "${DIR}/output" ]
then
  rm -rf "${DIR}/output"
fi

mkdir "${DIR}/output"
mkdir "${DIR}/output/config"

go build -o "${DIR}/output/app"
cp "${DIR}/config/application.json" "${DIR}/output/config"