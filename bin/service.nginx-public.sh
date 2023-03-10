#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

${DIR}/nginx/sbin/nginx -t -c /var/snap/platform/current/nginx.conf -e stderr
exec $DIR/nginx/sbin/nginx -c /var/snap/platform/current/nginx.conf -e stderr


