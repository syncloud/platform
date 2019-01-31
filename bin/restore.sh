#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$2" ]]; then
    echo "usage $0 app file [--include-data]"
    exit 1
fi

APP=$1
BACKUP_FILE=$2
INCLUDE_DATA=${3:-no}

STORAGE_DIR=/data

EXTRACT_DIR=${STORAGE_DIR}/platform/backup/${APP}

APP_DIR=/var/snap/$APP
APP_CURRENT_DIR=current
APP_COMMON_DIR=/var/snap/$APP/common

APP_DATA_DIR=${STORAGE_DIR}/$APP

APP_DATA_SIZE=$(stat --printf="%s" ${BACKUP_FILE})

STORAGE_SPACE_LEFT=$(df --output=avail ${STORAGE_DIR} | tail -1)
STORAGE_SPACE_NEEDED=$(( ${APP_DATA_SIZE} * 10 ))

echo "space left on storage: ${STORAGE_SPACE_LEFT}"
echo "space needed: ${STORAGE_SPACE_NEEDED}"

if [[ ${STORAGE_SPACE_NEEDED} -gt ${STORAGE_SPACE_LEFT} ]]; then
    echo "not enaugh space on storage for the backup"
    exit 1
fi

mkdir -p ${EXTRACT_DIR}
tar -C ${EXTRACT_DIR} -xf ${BACKUP_FILE}
ls -la ${EXTRACT_DIR}
snap stop $APP
rm -rf ${APP_DIR}/current/*
cp -R ${EXTRACT_DIR}/current/. ${APP_DIR}/current/
rm -rf ${APP_DIR}/common/*
cp -R ${EXTRACT_DIR}/common/. ${APP_DIR}/common/
if [[ "${INCLUDE_DATA}" == "--include-data" ]]; then
    cp -R ${EXTRACT_DIR}/data/. ${APP_DATA_DIR}/
fi
snap start $APP

rm -rf ${EXTRACT_DIR}
