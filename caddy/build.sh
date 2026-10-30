#!/bin/sh -ex

DIR=$( cd "$( dirname "$0" )" && pwd )
CADDY_VERSION=$1
L4_VERSION=$2
BUILD_DIR=${DIR}/../build/snap/caddy
mkdir -p ${BUILD_DIR}
export CGO_ENABLED=0

cd ${DIR}/matcher
go test -trimpath ./...

for attempt in 1 2 3 4 5; do
  if xcaddy build v${CADDY_VERSION} --with github.com/mholt/caddy-l4=github.com/cyberb/caddy-l4@v${L4_VERSION} --with github.com/syncloud/platform/caddy/matcher=${DIR}/matcher --output ${BUILD_DIR}/caddy; then
    break
  fi
  if [ $attempt -eq 5 ]; then
    exit 1
  fi
  sleep 10
done
