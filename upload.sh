#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

app=platform
branch=$1
build_number=$2
installer=$3
bucket=apps.syncloud.org
ARCH=$(uname -m)

mkdir -p /opt/app

ARCH=$(dpkg-architecture -q DEB_HOST_ARCH)
FILE_NAME=${app}_${build_number}_${ARCH}.snap

if [ "${branch}" == "master" ] || [ "${branch}" == "stable" ] ; then

  s3cmd put $FILE_NAME s3://${bucket}/apps/$FILE_NAME
  
  if [ "${branch}" == "stable" ]; then
    branch=rc
  fi

  printf ${build_number} > ${app}.version
  s3cmd put ${app}.version s3://${bucket}/releases/${branch}/${app}.version

fi