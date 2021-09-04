#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$2" ]]; then
    echo "usage $0 [internal|public] [stop|start]"
    exit 1
fi

case $2 in
start)
    exec $DIR/python/bin/uwsgi --ini /snap/platform/current/config/uwsgi/"$1".ini
    ;;
stop)
    exec $DIR/python/bin/uwsgi --stop /snap/platform/current/uwsgi."$1".pid
    ;;
*)
    echo "not valid command"
    exit 1
    ;;
esac

