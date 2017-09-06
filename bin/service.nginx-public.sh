#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 [start|stop]"
    exit 1
fi

case $2 in
start)
    $DIR/nginx/sbin/nginx -c ${SNAP_COMMON}/config.runtime/nginx/nginx/public.conf
    ;;
reload)
    $DIR/nginx/sbin/nginx -c ${SNAP_COMMON}/config.runtime/nginx/nginx/public.conf -s reload
    ;;
stop)
    $DIR/nginx/sbin/nginx -c ${SNAP_COMMON}/config.runtime/nginx/nginx/public.conf -s stop
    ;;
*)
    echo "not valid command"
    exit 1
    ;;
esac


