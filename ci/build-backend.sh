#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )
BIN=${DIR}/build/snap/bin
HOOKS=${DIR}/build/snap/meta/hooks

check_format() {
  local unformatted
  unformatted=$(gofmt -l .)
  if [[ -n "${unformatted}" ]]; then
    echo "not gofmt formatted:"
    echo "${unformatted}"
    exit 1
  fi
}

cd "${DIR}/backend"
for i in 1 2 3; do go mod download && break || sleep 5; done
check_format
go test ./... -coverprofile cover.out
go tool cover -func cover.out

for cmd in backend api cli; do
  CGO_ENABLED=0 go build -o "${BIN}/${cmd}" "./cmd/${cmd}"
  "${BIN}/${cmd}" -h
done

for hook in install post-refresh configure; do
  CGO_ENABLED=0 go build -o "${HOOKS}/${hook}" "./cmd/${hook}"
  "${HOOKS}/${hook}" -h
done

CGO_ENABLED=0 go build -o "${BIN}/login" ./cmd/login
CGO_ENABLED=0 go build -o "${BIN}/stability" ./cmd/stability

cd "${DIR}/visual-diff"
check_format
CGO_ENABLED=0 go build -o visual-diff ./cmd
