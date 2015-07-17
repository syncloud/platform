#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}
PYTHON=${DIR}/../build/platform/python/bin

if [[ -z "$1" || -z "$2" || -z "$3" || -z "$4" || -z "$5" || -z "$6" ]]; then
    echo "usage $0 redirect_user redirect_password redirect_domain release platform_version platform_arch"
    exit 1
fi

./docker.sh

apt-get install sshpass
SSH="sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost"

sshpass -p syncloud scp -o StrictHostKeyChecking=no -P 2222 ${DIR}/../platform-${5}-${6}.tar.gz root@localhost:/

${SSH} "/opt/app/sam/bin/sam --debug update --release $4"
${SSH} "/opt/app/sam/bin/sam --debug install /platform-${5}-${6}.tar.gz"

${PYTHON}/py.test -s verify.py --email=$1 --password=$2 --domain=$3 --release=$4