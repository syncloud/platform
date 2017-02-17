#!/usr/bin/env bash

BACKEND_MODE=$1

BACKEND_DIR="www/public/js/backend.${BACKEND_MODE}"
ACTIVE_BACKEND_DIR="www/public/js/backend"

rm -rf $ACTIVE_BACKEND_DIR

cp -a $BACKEND_DIR $ACTIVE_BACKEND_DIR