#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ "$#" -lt 6 ]; then
    echo "usage $0 redirect_user redirect_password redirect_domain release device_host package"
    exit 1
fi

ARCH=$(uname -m)
RELEASE=$4
DOMAIN=$3-${ARCH}-${DRONE_BRANCH}
DEVICE_HOST=$5
ARCHIVE=$6

APP=platform

if [ ${ARCH} == "x86_64" ]; then
    TEST_SUITE="verify.py test-ui.py"
else
    TEST_SUITE=verify.py
fi

APP_ARCHIVE_PATH=$(realpath "$ARCHIVE")

cd ${DIR}

echo ${APP_ARCHIVE_PATH}

cd ${DIR}

attempts=100
attempt=0

set +e
sshpass -p syncloud ssh -o StrictHostKeyChecking=no root@${DEVICE_HOST} date
while test $? -gt 0
do
  if [ ${attempt} -gt ${attempts} ]; then
    exit 1
  fi
  sleep 3
  echo "Waiting for SSH $attempt"
  attempt=$((attempt+1))
  sshpass -p syncloud ssh -o StrictHostKeyChecking=no root@${DEVICE_HOST} date
done
set -e

#sshpass -p syncloud scp -o StrictHostKeyChecking=no install-snapd.sh root@${DEVICE_HOST}:/installer.sh
#sshpass -p syncloud ssh -o StrictHostKeyChecking=no root@${DEVICE_HOST} /installer.sh ${RELEASE}
sshpass -p syncloud ssh -o StrictHostKeyChecking=no root@${DEVICE_HOST} snap remove platform

pip2 install -r ${DIR}/../requirements.txt

pip2 install -r ${DIR}/../src/dev_requirements.txt

#fix dns
device_ip=$(getent hosts ${DEVICE_HOST} | awk '{ print $1 }')
echo "$device_ip $DOMAIN.syncloud.info" >> /etc/hosts
echo "$device_ip $APP.$DOMAIN.syncloud.info" >> /etc/hosts
echo "$device_ip app.$DOMAIN.syncloud.info" >> /etc/hosts

cat /etc/hosts
py.test --fixtures
xvfb-run -l --server-args="-screen 0, 1024x4096x24" py.test -x -s ${TEST_SUITE} \
    --email=$1 --password=$2 --domain=${DOMAIN} --release=${RELEASE} \
    --app-archive-path=${APP_ARCHIVE_PATH} --device-host=${DEVICE_HOST}