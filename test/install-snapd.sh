#!/bin/bash -xe

VERSION=$1
if [ -z "${VERSION}" ]; then
  for i in 1 2 3 4 5 6 7 8 9 10; do
    VERSION=$(curl -fsS http://apps.syncloud.org/releases/stable/snapd2.version) && break
    echo "curl failed (attempt $i), retrying in 5s..."
    sleep 5
  done
fi
[ -n "${VERSION}" ] || { echo "failed to resolve snapd version after retries"; exit 1; }
ARCH=$(dpkg --print-architecture)
SNAPD=snapd-${VERSION}-${ARCH}.tar.gz

cd /tmp
rm -rf "${SNAPD}"
rm -rf snapd
wget --tries=10 --waitretry=5 --retry-connrefused http://apps.syncloud.org/apps/"${SNAPD}" --progress=dot:giga
tar xzvf "${SNAPD}"
mkdir -p /var/lib/snapd/snaps
./snapd/install.sh
