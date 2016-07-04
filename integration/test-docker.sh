#!/bin/bash

APP_ARCHIVE_PATH=$(realpath "$4")
echo ${APP_ARCHIVE_PATH}

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

export TMPDIR=/tmp
export TMP=/tmp
export DEBIAN_FRONTEND=noninteractive

if [ "$#" -eq 7 ]; then
    TEST_SUITE=$7
else
    TEST_SUITE="verify"
fi

if [ "$#" -lt 6 ]; then
    echo "usage $0 redirect_user redirect_password redirect_domain app_archive_path sam_version release"
    exit 1
fi

SAM_VERSION=$5
RELEASE=$6

./docker.sh ${SAM_VERSION} ${RELEASE}

apt-get install -y sshpass
pip2 install -r ${DIR}/../src/dev_requirements.txt
py.test -x -s ${TEST_SUITE}.py --email=$1 --password=$2 --domain=$3 --app-archive-path=${APP_ARCHIVE_PATH}