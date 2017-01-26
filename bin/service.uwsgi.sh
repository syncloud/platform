#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 [internal|public]"
    exit 1
fi

export LD_LIBRARY_PATH=$DIR/python/lib

exec $DIR/uwsgi/bin/uwsgi --ini ${SNAP_COMMON}/config/uwsgi/$1.ini


