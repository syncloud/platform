#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ "$#" -lt 4 ]; then
    echo "usage $0 domain device_host package suite"
    exit 1
fi

DOMAIN=$1
DEVICE_HOST=$2
PACKAGE=$3
TEST_SUITE=$4

APP=platform

APP_ARCHIVE_PATH=$(realpath "$PACKAGE")

pip2 install -r ${DIR}/dev_requirements.txt

#fix dns
device_ip=$(getent hosts ${DEVICE_HOST} | awk '{ print $1 }')
echo "$device_ip $DOMAIN.syncloud.info" >> /etc/hosts
echo "$device_ip $APP.$DOMAIN.syncloud.info" >> /etc/hosts
echo "$device_ip app.$DOMAIN.syncloud.info" >> /etc/hosts

cd $DIR
xvfb-run -l --server-args="-screen 0, 1024x4096x24" py.test -x -s ${TEST_SUITE} --domain=$DOMAIN --app-archive-path=${APP_ARCHIVE_PATH} --device-host=${DEVICE_HOST}
