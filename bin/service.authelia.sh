#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
rm -rf /var/snap/platform/current/authelia.socket
exec ${DIR}/authelia/authelia.sh \
  --config /var/snap/platform/current/config/authelia/config.yml \
  --config.experimental.filters template
