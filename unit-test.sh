#!/usr/bin/env bash

APP_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )"  && pwd )
PYTHON=${APP_DIR}/build/platform/python/bin

cd ${APP_DIR}/src

${PYTHON}/py.test --cov syncloud test