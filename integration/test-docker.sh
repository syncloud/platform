#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}


if [[ -z "$1" || -z "$2" || -z "$3" || -z "$4" || -z "$5" || -z "$6" ]]; then
    echo "usage $0 redirect_user redirect_password redirect_domain release platform_version platform_arch"
    exit 1
fi

./docker.sh

apt-get install sshpass
#ssh-keygen -f "/root/.ssh/known_hosts" -R [localhost]:2222

if [[ -n "$TEAMCITY_VERSION" ]]; then
    TC="export TEAMCITY_VERSION=\"$TEAMCITY_VERSION\" ; "
fi

SSH="sshpass -p syncloud ssh -o StrictHostKeyChecking=no root@localhost -p 2222"

#${SSH} "/test/integration/pip-install.sh"
${SSH} "/opt/app/sam/bin/sam --debug update --release $4"
${SSH} "/opt/app/sam/bin/sam --debug install /test/build/platform-${5}-${6}.tar.gz"
${SSH} "$TC /test/integration/unit-test.sh"
${SSH} "$TC /opt/app/platform/python/bin/py.test -s /test/integration/verify.py --email=$1 --password=$2 --domain=$3 --release=$4"