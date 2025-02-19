#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 [start|stop]"
    exit 1
fi

case $1 in
start)
    ${DIR}/nginx/bin/nginx -t -c /var/snap/platform/current/nginx.conf -e stderr
    exec $DIR/nginx/bin/nginx -c /var/snap/platform/current/nginx.conf -e stderr
    ;;
reload)
    $DIR/nginx/bin/nginx -c /var/snap/platform/current/nginx.conf -s reload -e stderr
    ;;
*)
    echo "not valid command"
    exit 1
    ;;
esac

