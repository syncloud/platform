#!/bin/bash


if [ "$#" -lt 1 ]; then
    echo "usage $0 release"
    exit 1
fi

ARCH=$(uname -m)
VERSION=89
RELEASE=$1

SAM=sam-${VERSION}-${ARCH}.tar.gz
wget http://apps.syncloud.org/apps/${SAM} --progress=dot:giga
tar xzf $SAM -C /opt/app
/opt/app/sam/bin/sam update --release ${RELEASE}