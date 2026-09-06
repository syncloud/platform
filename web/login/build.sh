#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
OUT=${DIR}/../../build/snap/web

cd "${DIR}"
npm install
npm run build

mkdir -p "${OUT}"
cp -r dist "${OUT}/login"
