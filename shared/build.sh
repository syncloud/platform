#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )
CADDY_VERSION=$1
L4_VERSION=$2
XCADDY_VERSION=$3
OUT=${DIR}/artifact/shared
export CGO_ENABLED=0
export GOARM=7

cd "${DIR}"
go install github.com/caddyserver/xcaddy/cmd/xcaddy@v${XCADDY_VERSION}
./backend/unit-test.sh
mkdir -p "${OUT}"

for arch in amd64 arm64 arm; do
  export GOARCH=${arch}
  rm -rf build/snap/bin build/snap/meta build/snap/caddy build/snap/frp test/testapp/build test/streamapp/build test/api/api.test test/externalapp/externalapp visual-diff/visual-diff
  ./caddy/build.sh ${CADDY_VERSION} ${L4_VERSION}
  ./frp/build.sh
  ./backend/build.sh
  ./test/testapp/cli/build.sh
  (cd test/api && go test -c -o api.test)
  (cd test/externalapp && go build -o externalapp)
  tar czf "${OUT}/${arch}.tar.gz" build/snap test/testapp/build test/streamapp/build test/api/api.test test/externalapp/externalapp visual-diff/visual-diff
done
ls -la "${OUT}"
