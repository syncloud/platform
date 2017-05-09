#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

export TMPDIR=/tmp
export TMP=/tmp
export DEBIAN_FRONTEND=noninteractive

if [ "$#" -lt 8 ]; then
    echo "usage $0 redirect_user redirect_password redirect_domain app_archive_path installer_version release [all|test_suite] [sam|snapd]"
    exit 1
fi

APP_ARCHIVE_PATH=$(realpath "$4")
INSTALLER_VERSION=$5
RELEASE=$6
TEST=$7
INSTALLER=$8

GECKODRIVER=0.14.0
FIREFOX=50.0

echo ${APP_ARCHIVE_PATH}

if [ "$TEST" == "all" ]; then
    TEST_SUITE="verify.py test-ui.py"
else
    TEST_SUITE=${TEST}.py
fi

cd ${DIR}
./docker.sh ${RELEASE}

sshpass -p syncloud scp -o StrictHostKeyChecking=no -P 2222 install-${INSTALLER}.sh root@localhost:/installer.sh

sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost /installer.sh ${INSTALLER_VERSION} ${RELEASE}

apt-get install -y sshpass xvfb

coin --to ${DIR} raw --subfolder geckodriver https://github.com/mozilla/geckodriver/releases/download/v${GECKODRIVER}/geckodriver-v${GECKODRIVER}-linux64.tar.gz
coin --to ${DIR} raw https://ftp.mozilla.org/pub/firefox/releases/${FIREFOX}/linux-x86_64/en-US/firefox-${FIREFOX}.tar.bz2
curl https://raw.githubusercontent.com/mguillem/JSErrorCollector/master/dist/JSErrorCollector.xpi -o  JSErrorCollector.xpi

pip2 install -r ${DIR}/../src/dev_requirements.txt
xvfb-run -l --server-args="-screen 0, 1024x4096x24" py.test -x -s ${TEST_SUITE} --email=$1 --password=$2 --domain=$3 --app-archive-path=${APP_ARCHIVE_PATH} --installer=${INSTALLER}