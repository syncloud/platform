#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
SOCKET=/var/snap/platform/current/authelia.socket
rm -rf ${SOCKET}
${DIR}/authelia/authelia.sh \
  --config /var/snap/platform/current/config/authelia/config.yml \
  --config.experimental.filters template &
PID=$!
while [ ! -S ${SOCKET} ]; do sleep 0.1; done
chmod 0777 ${SOCKET}
wait $PID
