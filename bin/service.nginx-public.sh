#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$2" ]]; then
    echo "usage $0 [internal|public] [start|stop]"
    exit 1
fi

case $2 in
start)
    $DIR/nginx/sbin/nginx -c ${SNAP_COMMON}/config/nginx/nginx-$1.conf
    ;;
reload)
    $DIR/nginx/sbin/nginx -c ${SNAP_COMMON}/config/nginx/nginx-$1.conf -s reload
    ;;
stop)
    $DIR/nginx/sbin/nginx -c ${SNAP_COMMON}/config/nginx/nginx-$1.conf -s stop
    ;;
*)
    echo "not valid command"
    exit 1
    ;;
esac


