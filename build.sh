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

if [ ! -d 3rdparty ]; then
  mkdir 3rdparty
fi

cd 3rdparty

if [ ! -f uwsgi-${ARCH}.tar.gz ]; then
  wget http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_uwsgi_${ARCH}/lastSuccessful/uwsgi.tar.gz\
  -O uwsgi-${ARCH}.tar.gz --progress dot:giga
else
  echo "skipping uwsgi build"
fi

if [ ! -f nginx-${ARCH}.tar.gz ]; then
  wget http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_nginx_${ARCH}/lastSuccessful/nginx.tar.gz\
  -O nginx-${ARCH}.tar.gz --progress dot:giga
else
  echo "skipping nginx build"
fi

if [ ! -f openldap-${ARCH}.tar.gz ]; then
  wget http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_openldap_${ARCH}/lastSuccessful/openldap.tar.gz\
  -O openldap-${ARCH}.tar.gz --progress dot:giga
else
  echo "skipping openldap build"
fi

cd ..

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