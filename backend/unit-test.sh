#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )

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

cd "${DIR}/visual-diff"
check_format
