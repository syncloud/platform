#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 app [--include-data]"
    exit 1
fi

APP=$1
INCLUDE_DATA=${1:-no}

STORAGE_DIR=/data
BACKUP_NAME=${APP}_`date +"%Y%m%d_%H%M%S"`
BASE_DIR=${STORAGE_DIR}/platform/backup/${APP}
BACKUP_DIR=${BASE_DIR}/${BACKUP_NAME}
APP_DIR=/var/snap/$APP
APP_CURRENT_DIR=current
APP_COMMON_DIR=/var/snap/$APP/common

APP_DATA_DIR=${STORAGE_DIR}/$APP

APP_DATA_SIZE=0

STORAGE_SPACE_LEFT=$(df --output=avail ${STORAGE_DIR} | tail -1)
STORAGE_SPACE_NEEDED=$(( ${APP_DATA_SIZE} * 10 ))

echo "space left on storage: ${STORAGE_SPACE_LEFT}"
echo "space needed: ${STORAGE_SPACE_NEEDED}"

if [[ ${STORAGE_SPACE_NEEDED} -gt ${STORAGE_SPACE_LEFT} ]]; then
    echo "not enaugh space on storage for the backup"
    exit 1
fi

mkdir -p ${BASE_DIR}
tar xf ${BASE_DIR}/${BACKUP_NAME}.tar.gz -C ${BASE_DIR}
snap stop $APP
rm -rf ${APP_DIR}/current/*
mv ${BACKUP_DIR}/current/* ${APP_DIR}/current/
rm -rf ${APP_DIR}/common/*
mv ${BACKUP_DIR}/common/* ${APP_DIR}/common/
if [[ "${INCLUDE_DATA}" == "--include-data" ]]; then
    mv ${BACKUP_DIR}/data/* ${APP_DATA_DIR}/
fi
snap start $APP

rm -rf ${BACKUP_DIR}
