#!/bin/bash -ex

if [[ -z "$2" ]]; then
    echo "usage $0 app file"
    exit 1
fi

APP=$1
BACKUP_FILE=$2

EXTRACT_DIR=$(mktemp -d)

BACKUP_SIZE=$(stat --printf="%s" ${BACKUP_FILE})

TEMP_SPACE_LEFT=$(df -B 1 --output=avail ${EXTRACT_DIR} | tail -1)
TEMP_SPACE_NEEDED=$(( ${BACKUP_SIZE} * 10 ))

echo "temp space left: ${TEMP_SPACE_LEFT}"
echo "temp space needed: ${TEMP_SPACE_NEEDED}"

if [[ ${TEMP_SPACE_NEEDED} -gt ${TEMP_SPACE_LEFT} ]]; then
    echo "not enough temp space for the restore"
    exit 1
fi

tar -C ${EXTRACT_DIR} -xf ${BACKUP_FILE}
ls -la ${EXTRACT_DIR}
APP_DIR=/var/snap/$APP

snap stop $APP

rm -rf ${APP_DIR}/current/*
cp -R --preserve ${EXTRACT_DIR}/current/. ${APP_DIR}/current/
rm -rf ${APP_DIR}/common/*
cp -R --preserve ${EXTRACT_DIR}/common/. ${APP_DIR}/common/

snap start $APP

rm -rf ${EXTRACT_DIR}
