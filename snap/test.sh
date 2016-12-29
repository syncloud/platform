#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

if [ "$#" -lt 3 ]; then
    echo "usage $0 release snapd_file snap_file"
    exit 1
fi

SAM_VERSION=none
RELEASE=$1
SNAPD_FILE=$2
SNAP_FILE=$3

${DIR}/../integration/docker.sh ${SAM_VERSION} ${RELEASE}

sshpass -p syncloud scp -o StrictHostKeyChecking=no -P 2222 $SNAPD_FILE root@localhost:/snapd.tar.gz

sshpass -p syncloud scp -o StrictHostKeyChecking=no -P 2222 $SNAP_FILE  root@localhost:/platform.snap

sshpass -p syncloud scp -o StrictHostKeyChecking=no -P 2222 install-snapd.sh root@localhost:/

sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost "/install-snapd.sh"

set +e
sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost "snap install /platform.snap --devmode"
exit_code=$?
set -e

sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost journalctl | tail -200

exit $exit_code
