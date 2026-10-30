#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )
BIN=${DIR}/build/snap/bin
HOOKS=${DIR}/build/snap/meta/hooks

cd "${DIR}/backend"
for i in 1 2 3; do go mod download && break || sleep 5; done

for cmd in backend api cli; do
  CGO_ENABLED=0 go build -o "${BIN}/${cmd}" "./cmd/${cmd}"
done

for hook in install post-refresh configure; do
  CGO_ENABLED=0 go build -o "${HOOKS}/${hook}" "./cmd/${hook}"
done

CGO_ENABLED=0 go build -o "${BIN}/login" ./cmd/login
CGO_ENABLED=0 go build -o "${BIN}/stability" ./cmd/stability

cd "${DIR}/visual-diff"
CGO_ENABLED=0 go build -o visual-diff ./cmd
