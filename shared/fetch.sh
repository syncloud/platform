#!/bin/sh -ex

DIR=$( cd "$( dirname "$0" )/.." && pwd )
BUILD=$1
ARCH=$2
URL=http://${ARTIFACT_HOST}:8081/files/platform/${BUILD}-shared/${ARCH}.tar.gz

cd ${DIR}
for attempt in 1 2 3 4 5; do
  if wget -q ${URL} -O shared.tar.gz; then
    break
  fi
  if [ $attempt -eq 5 ]; then
    exit 1
  fi
  sleep 10
done
tar xzf shared.tar.gz
rm shared.tar.gz
