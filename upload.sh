#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

app=platform
branch=$1
build_number=$2
installer=$3
bucket=apps.syncloud.org
ARCH=$(uname -m)

mkdir -p /opt/app
SAMCMD=/opt/app/sam/bin/sam

FILE_NAME=${app}-${build_number}-${ARCH}.tar.gz
if [ $installer == "snapd" ]; then
  ARCH=$(dpkg-architecture -q DEB_HOST_ARCH)
  FILE_NAME=${app}_${build_number}_${ARCH}.snap
fi

if [ "${branch}" == "master" ] || [ "${branch}" == "stable" ] ; then

  s3cmd put $FILE_NAME s3://${bucket}/apps/$FILE_NAME
  
  if [ "${branch}" == "stable" ]; then
    branch=rc
  fi

  ${SAMCMD} release $branch $branch --override ${app}=${build_number}

fi

