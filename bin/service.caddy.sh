#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
CONFIG=/var/snap/platform/current/Caddyfile
ADMIN=unix//var/snap/platform/current/caddy-admin.sock
export XDG_DATA_HOME=/var/snap/platform/current/caddy

case $1 in
start)
    rm -f /var/snap/platform/current/authelia.socket /var/snap/platform/current/caddy-admin.sock
    touch /var/snap/platform/current/custom-proxy.caddy
    exec ${DIR}/caddy/caddy run --config ${CONFIG} --adapter caddyfile
    ;;
reload)
    exec ${DIR}/caddy/caddy reload --force --config ${CONFIG} --adapter caddyfile --address ${ADMIN}
    ;;
*)
    echo "usage $0 [start|reload]"
    exit 1
    ;;
esac
