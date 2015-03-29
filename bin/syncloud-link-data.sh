#!/bin/bash

DATA_DIR=/data

if [ -L "$DATA_DIR" ] && [ "$2" != "force" ] ; then
   echo "$DATA_DIR link exists, use force parameter to override" 1>&2
   exit 1
fi

rm -rf /data
ln -s $1 /data
chown -RL www-data. /data
touch /data/.ocdata
ls -la /data