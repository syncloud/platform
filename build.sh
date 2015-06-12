#!/usr/bin/env bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

export TMPDIR=/tmp
export TMP=/tmp
NAME=platform
USER=platform

ARCH=x86_64
if [[ -n "$1" ]]; then
    ARCH=$1
fi


function 3rdparty {
  APP=$1
  if [ ! -d 3rdparty ]; then
    mkdir 3rdparty
  fi
  cd 3rdparty
  if [ ! -f ${APP}-${ARCH}.tar.gz ]; then
    wget http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_${APP}_${ARCH}/lastSuccessful/${APP}.tar.gz\
    -O ${APP}-${ARCH}.tar.gz --progress dot:giga
  else
    echo "skipping ${APP}"
  fi
  cd ..
}

3rdparty uwsgi
3rdparty nginx
3rdparty openldap

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
tar xzf 3rdparty/nginx-${ARCH}.tar.gz -C build/${NAME}
echo "extracting uwsgi"
tar xzf 3rdparty/uwsgi-${ARCH}.tar.gz -C build/${NAME}
echo "extracting openldap"
tar xzf 3rdparty/openldap-${ARCH}.tar.gz -C build/${NAME}
rm -rf ${NAME}.tar.gz
echo "zipping"
tar cpzf ${NAME}.tar.gz -C build/ ${NAME}

echo "app: ${NAME}.tar.gz"