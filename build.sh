#!/usr/bin/env bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

ROOT=/opt
APP_NAME=syncloud-platform
APP_ROOT=${ROOT}/${APP_NAME}

if [ ! -d nginx/build ]; then
  ./nginx/build.sh
else
  echo "skipping nginx build"
fi

rm -rf ${APP_ROOT}
mkdir ${APP_ROOT}

cp -r config ${APP_ROOT}/
tar xzf nginx/build/nginx.tar.gz -C ${APP_ROOT}/

tar cpzf ${APP_NAME}.tar.gz -C ${ROOT} ${APP_NAME}

