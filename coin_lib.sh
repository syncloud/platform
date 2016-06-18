#!/usr/bin/env bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

if [ ! -z "$1" ]; then
  ARCH=$1
else
  ARCH=$(dpkg-architecture -qDEB_HOST_GNU_CPU)
fi

if [ ! -d lib ]; then
  mkdir lib
fi

rm -rf lib/*

coin --to=lib py http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_psutil_${ARCH}/lastSuccessful/psutil-2.1.3-cp27-none-linux_${ARCH}.whl
coin --to=lib py http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_miniupnpc_${ARCH}/lastSuccessful/miniupnpc-1.9-cp27-none-linux_${ARCH}.whl

coin --to=lib py http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_python_ldap_${ARCH}/lastSuccessful/python_ldap-2.4.19-cp27-none-linux_${ARCH}.whl

coin --to=lib raw http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_certbot_${ARCH}/lastSuccessful/certbot-${ARCH}.tar.gz

coin --to=lib raw http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_openssl_${ARCH}/lastSuccessful/openssl-${ARCH}.tar.gz
