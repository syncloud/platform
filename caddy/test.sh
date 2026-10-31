#!/bin/sh -ex

DIR=$( cd "$( dirname "$0" )" && pwd )
BUILD_DIR=${DIR}/../build/snap/caddy

${BUILD_DIR}/caddy version
${BUILD_DIR}/caddy list-modules | grep '^layer4.matchers.syncloud_stream$'
${BUILD_DIR}/caddy list-modules | grep '^caddy.listeners.layer4$'
