#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

./docker.sh

apt-get install sshpass
sshpass -p "syncloud" ssh -o StrictHostKeyChecking=no root@localhost -p 2222 "/test/integration/pip-install.sh"
sshpass -p "syncloud" ssh -o StrictHostKeyChecking=no root@localhost -p 2222 "/test/integration/binary-install.py"