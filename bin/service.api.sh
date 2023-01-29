#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

exec $DIR/api unix ${SNAP_DATA}/api.sock
