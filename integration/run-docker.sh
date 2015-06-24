#!/bin/bash

if [[ -z "$1" || -z "$2" || -z "$3" || -z "$4" ]]; then
    echo "usage $0 redirect_user redirect_password redirect_domain release"
    exit 1
fi

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

./docker.sh

apt-get install sshpass

SSH="sshpass -p syncloud ssh -o StrictHostKeyChecking=no root@localhost -p 2222"

${SSH} "/opt/app/sam/bin/sam --debug update --release $4"
${SSH} "/opt/app/sam/bin/sam --debug install /test/build/platform-local-x86_64.tar.gz"