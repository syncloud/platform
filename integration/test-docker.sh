#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}


if [[ -z "$1" || -z "$2" || -z "$3" || -z "$4" ]]; then
    echo "usage $0 redirect_user redirect_password redirect_domain release"
    exit 1
fi

./docker.sh

apt-get install sshpass
#ssh-keygen -f "/root/.ssh/known_hosts" -R [localhost]:2222

if [[ -n "$TEAMCITY_VERSION" ]]; then
    TC="export TEAMCITY_VERSION=\"$TEAMCITY_VERSION\" ; "
fi

sshpass -p "syncloud" ssh -o StrictHostKeyChecking=no root@localhost -p 2222 "$TC /test/integration/unit-test.sh"
sshpass -p "syncloud" ssh -o StrictHostKeyChecking=no root@localhost -p 2222 "/test/integration/pip-install.sh"
sshpass -p "syncloud" ssh -o StrictHostKeyChecking=no root@localhost -p 2222 "/test/integration/binary-install.py"
sshpass -p "syncloud" ssh -o StrictHostKeyChecking=no root@localhost -p 2222 "$TC py.test -s /test/integration/verify.py --email=$1 --password=$2 --domain=$3 --release=$4"