#!/bin/bash

APP_ARCHIVE_PATH=$(realpath "$4")
echo ${APP_ARCHIVE_PATH}

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

export TMPDIR=/tmp
export TMP=/tmp
export DEBIAN_FRONTEND=noninteractive

if [ "$#" -eq 7 ]; then
    TEST_SUITE=$7.py
else
    TEST_SUITE="verify.py test-ui.py"
fi

if [ "$#" -lt 6 ]; then
    echo "usage $0 redirect_user redirect_password redirect_domain app_archive_path sam_version release"
    exit 1
fi

ARCH=$(dpkg-architecture -q DEB_HOST_GNU_CPU)
SAM_VERSION=$5
RELEASE=$6

./docker.sh ${RELEASE}

SAM=sam-${SAM_VERSION}-${ARCH}.tar.gz
if [ ! -f ${SAM} ]; then
  wget http://apps.syncloud.org/apps/${SAM} --progress=dot:giga
else
  echo "skipping sam"
fi
sshpass -p syncloud scp -o StrictHostKeyChecking=no -P 2222 $SAM root@localhost:/sam.tar.gz

sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost "tar xzf /sam.tar.gz -C ${ROOTFS}/opt/app

sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost "/opt/app/sam/bin/sam update --release ${RELEASE}"

apt-get install -y sshpass xvfb firefox

coin --to ${DIR} raw --subfolder geckodriver https://github.com/mozilla/geckodriver/releases/download/v0.9.0/geckodriver-v0.9.0-linux64.tar.gz
mv ${DIR}/geckodriver/geckodriver ${DIR}/geckodriver/wires

pip2 install -r ${DIR}/../src/dev_requirements.txt
xvfb-run --server-args="-screen 0, 1024x4096x24" py.test -x -s ${TEST_SUITE} --email=$1 --password=$2 --domain=$3 --app-archive-path=${APP_ARCHIVE_PATH}