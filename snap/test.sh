#!/bin/bash -x

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

if [ "$#" -lt 3 ]; then
    echo "usage $0 snapd_version sam_version release"
    exit 1
fi

SNAPD_VERSION=$1
SAM_VERSION=$2
RELEASE=$3

${DIR}/../integration/docker.sh ${SAM_VERSION} ${RELEASE}

sshpass -p syncloud scp -o StrictHostKeyChecking=no -P 2222 syncloud-platform_16.11_amd64.snap root@localhost:/

sshpass -p syncloud scp -o StrictHostKeyChecking=no -P 2222 snapd root@localhost:/

sshpass -p syncloud scp -o StrictHostKeyChecking=no -P 2222 install-snapd.sh root@localhost:/

sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost "/install-snapd.sh"

sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost "snap install /syncloud-platform_16.11_amd64.snap --devmode"

exit_code=$?

sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost journalctl | tail -200

exit $exit_code
