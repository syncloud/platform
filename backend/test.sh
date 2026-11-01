#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )
BIN=${DIR}/build/snap/bin
HOOKS=${DIR}/build/snap/meta/hooks

for cmd in backend api cli; do
  "${BIN}/${cmd}" -h
done

for hook in install post-refresh configure; do
  "${HOOKS}/${hook}" -h
done
