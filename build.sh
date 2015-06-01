#!/usr/bin/env bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

NAME=platform
USER=www-data

if [ ! -d uwsgi/build ]; then
  ./uwsgi/build.sh
else
  echo "skipping uwsgi build"
fi

if [ ! -d nginx/build ]; then
  ./nginx/build.sh
else
  echo "skipping nginx build"
fi

#if [ ! -d avahi/build ]; then
#  ./avahi/build.sh
#else
#  echo "skipping avahi build"
#fi

rm -rf build
mkdir -p build/${NAME}

cp -r bin build/${NAME}
cp -r config build/${NAME}
cp -r www build/${NAME}
cp -r socket build/${NAME}
chown -R ${USER}. build/${NAME}/socket

#tar xzf avahi/avahi.tar.gz -C build/${NAME}
tar xzf nginx/build/nginx.tar.gz -C build/${NAME}
cp -rf uwsgi/install/uwsgi build/${NAME}
rm -rf ${NAME}.tar.gz
tar cpzf ${NAME}.tar.gz -C build/ ${NAME}

echo "app: ${NAME}.tar.gz"