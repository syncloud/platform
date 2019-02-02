#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$2" ]]; then
    echo "usage $0 app file"
    exit 1
fi

APP=$1
BACKUP_FILE=$2

TEMP_DIR=$(mktemp -d)

APP_BASE_DIR=/var/snap/$APP
APP_CURRENT_DIR=${APP_BASE_DIR}/current
APP_COMMON_DIR=${APP_BASE_DIR}/common

APP_CURRENT_SIZE=$(du -s ${APP_CURRENT_DIR} | cut -f1)
APP_COMMON_SIZE=$(du -s ${APP_COMMON_DIR} | cut -f1)

TEMP_SPACE_LEFT=$(df --output=avail ${TEMP_DIR} | tail -1)
TEMP_SPACE_NEEDED=$(( (${APP_CURRENT_SIZE} + ${APP_COMMON_SIZE}) * 2 ))

echo "temp space left: ${TEMP_SPACE_LEFT}"
echo "temp space needed: ${TEMP_SPACE_NEEDED}"

if [[ ${TEMP_SPACE_NEEDED} -gt ${TEMP_SPACE_LEFT} ]]; then
    echo "not enaugh temp space for the backup"
    exit 1
fi

snap stop $APP
mkdir ${TEMP_DIR}/current
cp -R --preserve ${APP_CURRENT_DIR}/. ${BACKUP_DIR}/current

mkdir ${TEMP_DIR}/common
cp -R --preserve ${APP_COMMON_DIR}/. ${BACKUP_DIR}/common

snap start $APP
tar czf ${BACKUP_FILE} -C ${TEMP_DIR} .
rm -rf ${TEMP_DIR}
