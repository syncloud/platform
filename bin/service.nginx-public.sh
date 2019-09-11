#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 [start|stop]"
    exit 1
fi

case $1 in
start)
    ${DIR}/nginx/sbin/nginx -t -c ${SNAP_COMMON}/config.runtime/nginx/nginx.conf -g 'error_log '${SNAP_COMMON}'/log/nginx_public_error.log warn;'
    exec $DIR/nginx/sbin/nginx -c ${SNAP_COMMON}/config.runtime/nginx/nginx.conf -g 'error_log '${SNAP_COMMON}'/log/nginx_public_error.log warn;'
    ;;
post-start)
    timeout 5 /bin/bash -c 'until echo > /dev/tcp/localhost/80; do sleep 1; done'
    ;;
reload)
    $DIR/nginx/sbin/nginx -c ${SNAP_COMMON}/config.runtime/nginx/nginx.conf -s reload -g 'error_log '${SNAP_COMMON}'/log/nginx_public_error.log warn;'
    ;;
stop)
    $DIR/nginx/sbin/nginx -c ${SNAP_COMMON}/config.runtime/nginx/nginx.conf -s stop -g 'error_log '${SNAP_COMMON}'/log/nginx_public_error.log warn;'
    ;;
*)
    echo "not valid command"
    exit 1
    ;;
esac


