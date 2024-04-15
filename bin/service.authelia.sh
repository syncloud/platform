#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
exec ${DIR}/auth/authelia.sh \
  --config /var/snap/platform/current/config/authelia/config.yml \
  --config.experimental.filters expand-env,template
