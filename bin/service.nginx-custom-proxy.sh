#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 [start|stop]"
    exit 1
fi

case $1 in
start)
    rm -f /var/snap/platform/current/custom-proxy.socket
    ${DIR}/nginx/bin/nginx.sh -t -c /var/snap/platform/current/custom-proxy.conf -e stderr
    exec $DIR/nginx/bin/nginx.sh -c /var/snap/platform/current/custom-proxy.conf -e stderr
    ;;
reload)
    $DIR/nginx/bin/nginx.sh -c /var/snap/platform/current/custom-proxy.conf -s reload -e stderr
    ;;
*)
    echo "not valid command"
    exit 1
    ;;
esac
