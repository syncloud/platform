#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

export TMPDIR=/tmp
export TMP=/tmp
export DEBIAN_FRONTEND=noninteractive

if [ "$#" -lt 7 ]; then
    echo "usage $0 redirect_user redirect_password redirect_domain version release [sam|snapd] device_host"
    exit 1
fi

ARCH=$(uname -m)
VERSION=$4
RELEASE=$5
INSTALLER=$6
DOMAIN=$3-${ARCH}-${DRONE_BRANCH}-${INSTALLER}
DEVICE_HOST=$7

APP=platform
GECKODRIVER=0.14.0
FIREFOX=52.0

if [ $ARCH == "x86_64" ]; then
    SNAP_ARCH=amd64
else
    SNAP_ARCH=armhf
fi

if [ $INSTALLER == "snapd" ]; then
    ARCHIVE=${APP}_${VERSION}_${SNAP_ARCH}.snap
else
    ARCHIVE=${APP}-${VERSION}-${ARCH}.tar.gz
fi
APP_ARCHIVE_PATH=$(realpath "$ARCHIVE")

cd ${DIR}

echo ${APP_ARCHIVE_PATH}

if [ "$ARCH" == "x86_64" ]; then
    TEST_SUITE="verify.py test-ui.py"
else
    TEST_SUITE=verify.py
fi

cd ${DIR}

attempts=100
attempt=0

set +e
sshpass -p syncloud ssh -o StrictHostKeyChecking=no root@${DEVICE_HOST} date
while test $? -gt 0
do
  if [ $attempt -gt $attempts ]; then
    exit 1
  fi
  sleep 3
  echo "Waiting for SSH $attempt"
  attempt=$((attempt+1))
  sshpass -p syncloud ssh -o StrictHostKeyChecking=no root@${DEVICE_HOST} date
done
set -e

sshpass -p syncloud scp -o StrictHostKeyChecking=no install-${INSTALLER}.sh root@${DEVICE_HOST}:/installer.sh

sshpass -p syncloud ssh -o StrictHostKeyChecking=no root@${DEVICE_HOST} /installer.sh ${RELEASE}

pip2 install -r ${DIR}/../src/dev_requirements.txt

coin --to ${DIR} raw --subfolder geckodriver https://github.com/mozilla/geckodriver/releases/download/v${GECKODRIVER}/geckodriver-v${GECKODRIVER}-linux64.tar.gz
coin --to ${DIR} raw https://ftp.mozilla.org/pub/firefox/releases/${FIREFOX}/linux-x86_64/en-US/firefox-${FIREFOX}.tar.bz2
curl https://raw.githubusercontent.com/mguillem/JSErrorCollector/master/dist/JSErrorCollector.xpi -o  JSErrorCollector.xpi

#fix dns
device_ip=$(getent hosts ${DEVICE_HOST} | awk '{ print $1 }')
echo "$device_ip $APP.$DOMAIN.syncloud.info" >> /etc/hosts

cat /etc/hosts

xvfb-run -l --server-args="-screen 0, 1024x4096x24" py.test -x -s ${TEST_SUITE} --email=$1 --password=$2 --domain=$DOMAIN --release=$RELEASE --app-archive-path=${APP_ARCHIVE_PATH} --installer=${INSTALLER} --device-host=${DEVICE_HOST}