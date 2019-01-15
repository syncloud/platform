#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 app"
    exit 1
fi

APP=$1

BACKUP_NAME=${APP}_`date +"%Y%m%d_%H%M%S"`
BASE_DIR=/data/platform/backup/${APP}
BACKUP_DIR=${BASE_DIR}/${BACKUP_NAME}

mkdir -p ${BACKUP_DIR}

snap stop $APP
cp -r /var/snap/$APP/current/ ${BACKUP_DIR}/
cp -r /var/snap/$APP/common/ ${BACKUP_DIR}/
snap start $APP
tar czf ${BACKUP_NAME}.tar.gz -C ${BASE_DIR} ${BACKUP_NAME}
rm -rf ${BACKUP_DIR}
