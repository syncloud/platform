#!/usr/bin/env bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

NAME=platform
USER=platform

if [ ! -f uwsgi/uwsgi.tar.gz ]; then
  ./uwsgi/build.sh
else
  echo "skipping uwsgi build"
fi

if [ ! -f nginx/nginx.tar.gz ]; then
  ./nginx/build.sh
else
  echo "skipping nginx build"
fi

if [ ! -f openldap/openldap.tar.gz ]; then
  echo "no openldap build, get one from 3rdparty"
  exit 1
else
  echo "skipping openldap build"
fi

if ! jekyll -v; then
  echo "installing jekyll"
  apt-get -y install ruby ruby-dev make gcc nodejs
  gem install jekyll --no-rdoc --no-ri
fi

rm -rf build
mkdir -p build/${NAME}

echo "copying files"
cp -r bin build/${NAME}
cp -r config build/${NAME}
cd www
rm -rf _site
jekyll build
cd ..
cp -r www build/${NAME}
cp -r socket build/${NAME}

echo "extracting nginx"
tar xzf nginx/nginx.tar.gz -C build/${NAME}
echo "extracting uwsgi"
tar xzf uwsgi/uwsgi.tar.gz -C build/${NAME}
echo "extracting openldap"
tar xzf openldap/openldap.tar.gz -C build/${NAME}
rm -rf ${NAME}.tar.gz
echo "zipping"
tar cpzf ${NAME}.tar.gz -C build/ ${NAME}

echo "app: ${NAME}.tar.gz"