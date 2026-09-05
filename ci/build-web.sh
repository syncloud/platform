#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )
cd "${DIR}"

mkdir -p build/snap/web

cd "${DIR}/web/platform"
npm config ls -l
npm config set fetch-retry-mintimeout 200000
npm config set fetch-retry-maxtimeout 1200000
npm install
npm update browserslist
npm run test
npm run lint
npm run build
cp -r dist "${DIR}/build/snap/web/platform"

cd "${DIR}/web/login"
npm install
npm run build
cp -r dist "${DIR}/build/snap/web/login"
