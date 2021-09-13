#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

exec $DIR/backend unix ${SNAP_DATA}/backend.sock
