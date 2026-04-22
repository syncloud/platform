#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 [start|stop]"
    exit 1
fi

case $1 in
start)
    SOCKET=/var/snap/platform/current/authelia.socket
    rm -f ${SOCKET}
    ${DIR}/nginx/bin/nginx.sh -t -c /var/snap/platform/current/nginx.conf -e stderr
    ${DIR}/nginx/bin/nginx.sh -c /var/snap/platform/current/nginx.conf -e stderr
    chmod 0777 ${SOCKET}
    ;;
reload)
    $DIR/nginx/bin/nginx.sh -c /var/snap/platform/current/nginx.conf -s reload -e stderr
    ;;
*)
    echo "not valid command"
    exit 1
    ;;
esac

