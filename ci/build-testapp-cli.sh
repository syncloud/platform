#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )
BUILD=${DIR}/test/testapp/build

cd "${DIR}/test/testapp/cli"

CGO_ENABLED=0 go build -o "${BUILD}/meta/hooks/install" ./cmd/install
CGO_ENABLED=0 go build -o "${BUILD}/meta/hooks/configure" ./cmd/configure
CGO_ENABLED=0 go build -o "${BUILD}/bin/cli" ./cmd/cli
CGO_ENABLED=0 go build -o "${BUILD}/bin/backend" ./cmd/backend
