#!/bin/sh -ex

DIR=$( cd "$( dirname "$0" )" && pwd )
BUILD_DIR=${DIR}/../build/snap/frp
mkdir -p ${BUILD_DIR}

FRP_REPO=https://github.com/cyberb/frp.git
FRP_REF=webserver-unix-socket

SRC=${DIR}/../build/frp-src
rm -rf ${SRC}
git clone --depth 1 --branch ${FRP_REF} ${FRP_REPO} ${SRC}
cd ${SRC}
CGO_ENABLED=0 go build -tags noweb -o ${BUILD_DIR}/frpc ./cmd/frpc
${BUILD_DIR}/frpc --version
