#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
SOCKET=/var/snap/platform/current/authelia-internal.socket
rm -f ${SOCKET}
exec ${DIR}/authelia/authelia.sh \
  --config /var/snap/platform/current/config/authelia/config.yml \
  --config.experimental.filters template
