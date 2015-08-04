#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

if [ "$#" -ne 7 ]; then
    echo "usage $0 redirect_user redirect_password redirect_domain release app_version app_arch sam_version"
    exit 1
fi

./docker.sh $7

apt-get install sshpass
SSH="sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost"
SCP="sshpass -p syncloud scp -o StrictHostKeyChecking=no -P 2222"

${SCP} ${DIR}/../platform-${5}-${6}.tar.gz root@localhost:/

py.test -s verify.py --email=$1 --password=$2 --domain=$3 --release=$4 --app-version=$5 --arch=$6

${SCP} root@localhost:/opt/data/platform/log/\* .