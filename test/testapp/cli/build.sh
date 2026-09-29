#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
BUILD=${DIR}/../build

cd "${DIR}"

for i in 1 2 3; do go mod download && break || sleep 5; done
CGO_ENABLED=0 go build -o "${BUILD}/meta/hooks/install" ./cmd/install
CGO_ENABLED=0 go build -o "${BUILD}/meta/hooks/configure" ./cmd/configure
CGO_ENABLED=0 go build -o "${BUILD}/bin/cli" ./cmd/cli
CGO_ENABLED=0 go build -o "${BUILD}/bin/backend" ./cmd/backend
