#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

app=platform
branch=$1
build_number=$2
bucket=apps.syncloud.org
arch=$(uname -m)

mkdir -p /opt/app
SAMCMD=/opt/app/sam/bin/sam

if [ ! -f ${SAMCMD} ]; then
    ${DIR}/install-sam.sh 85 stable
fi

if [ ! -f /usr/bin/s3cmd ]; then
    ${DIR}/install-s3cmd.sh
fi

if [ "${branch}" == "master" ] || [ "${branch}" == "stable" ] ; then
  
  s3cmd put ${app}-${build_number}-${arch}.tar.gz s3://${bucket}/apps/${app}-${build_number}-${arch}.tar.gz
  
  if [ "${branch}" == "stable" ]; then
    branch=rc
  fi

  ${SAMCMD} release $branch $branch --override ${app}=${build_number}

fi

